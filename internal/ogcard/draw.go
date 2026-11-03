package ogcard

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	ink    = color.NRGBA{0x16, 0x15, 0x13, 0xff}
	foam   = color.NRGBA{0xef, 0xe9, 0xde, 0xff}
	copper = color.NRGBA{0xf5, 0x9e, 0x0b, 0xff}
	muted  = color.NRGBA{0xb8, 0xb0, 0xa4, 0xff}
	// copperGhost is the watermark S: NRGBA so the 16% alpha is not treated as
	// premultiplied (color.RGBA would wrap and paint lime).
	copperGhost = color.NRGBA{0xf5, 0x9e, 0x0b, 0x2e}
	copperRule  = color.NRGBA{0xf5, 0x9e, 0x0b, 0x88}
)

func paint(card Card) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), image.NewUniform(ink), image.Point{}, draw.Src)
	paintGlow(img)
	paintGrain(img)
	draw.Draw(img, image.Rect(0, 0, 8, Height), image.NewUniform(copper), image.Point{}, draw.Src)

	if card.Site {
		return paintSite(img, card)
	}
	return paintArticle(img, card)
}

func paintSite(img *image.RGBA, card Card) (image.Image, error) {
	kicker, err := sansFace(15)
	if err != nil {
		return nil, err
	}
	title, err := serifFace(108)
	if err != nil {
		return nil, err
	}
	byline, err := serifFace(28)
	if err != nil {
		return nil, err
	}
	meta, err := sansFace(16)
	if err != nil {
		return nil, err
	}
	mark, err := serifFace(280)
	if err != nil {
		return nil, err
	}

	drawString(img, mark, copperGhost, 820, 500, "S")
	drawTracked(img, kicker, copper, 80, 96, strings.ToUpper(card.Kicker), 4)
	drawString(img, title, foam, 72, 250, card.Title)
	drawString(img, byline, muted, 80, 340, card.Byline)
	drawHairline(img, 80, 400, 720)
	drawFooter(img, meta, card)
	return img, nil
}

func paintArticle(img *image.RGBA, card Card) (image.Image, error) {
	masthead, err := serifFace(28)
	if err != nil {
		return nil, err
	}
	kicker, err := sansFace(14)
	if err != nil {
		return nil, err
	}
	titleSize := articleTitleSize(card.Title)
	title, err := serifFace(titleSize)
	if err != nil {
		return nil, err
	}
	meta, err := sansFace(16)
	if err != nil {
		return nil, err
	}
	badge, err := serifFace(34)
	if err != nil {
		return nil, err
	}

	drawTracked(img, masthead, foam, 80, 86, "SINTAXE", 3)
	paintBadge(img, badge)
	drawHairline(img, 80, 122, Width-80)
	if card.Kicker != "" {
		drawTracked(img, kicker, copper, 80, 168, strings.ToUpper(card.Kicker), 3)
	}

	maxWidth := Width - 160
	lines := wrapLines(title, card.Title, maxWidth, 3)
	y := 230
	lineH := int(titleSize * 1.12)
	for _, line := range lines {
		drawString(img, title, foam, 80, y, line)
		y += lineH
	}
	drawFooter(img, meta, card)
	return img, nil
}

func articleTitleSize(title string) float64 {
	n := len([]rune(title))
	switch {
	case n < 22:
		return 72
	case n < 48:
		return 56
	default:
		return 48
	}
}

func paintBadge(img *image.RGBA, face font.Face) {
	const side = 56
	x, y := Width-80-side, 48
	draw.Draw(img, image.Rect(x, y, x+side, y+side), image.NewUniform(copper), image.Point{}, draw.Src)
	sW := measure(face, "S")
	m := face.Metrics()
	glyphH := (m.Ascent + m.Descent).Round()
	drawString(img, face, ink, x+(side-sW)/2, y+(side-glyphH)/2+m.Ascent.Round()-2, "S")
}

func drawFooter(img *image.RGBA, face font.Face, card Card) {
	const y = 562
	if card.Byline != "" && !card.Site {
		drawString(img, face, muted, 80, y, card.Byline)
	}
	host := card.Host
	if host == "" {
		host = defaultHost
	}
	w := measure(face, host)
	drawString(img, face, muted, Width-80-w, y, host)
}

func drawHairline(img *image.RGBA, x0, y, x1 int) {
	draw.Draw(img, image.Rect(x0, y, x1, y+1), image.NewUniform(copperRule), image.Point{}, draw.Over)
}

func drawString(img *image.RGBA, face font.Face, col color.Color, x, y int, s string) {
	if s == "" {
		return
	}
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: face, Dot: fixed.P(x, y)}
	d.DrawString(s)
}

func drawTracked(img *image.RGBA, face font.Face, col color.Color, x, y int, s string, extra int) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: face, Dot: fixed.P(x, y)}
	for _, r := range s {
		d.DrawString(string(r))
		d.Dot.X += fixed.I(extra)
	}
}

func paintGlow(img *image.RGBA) {
	const cx, cy, maxR = 1040.0, 90.0, 520.0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			d := math.Hypot(dx, dy)
			if d >= maxR {
				continue
			}
			t := (1 - d/maxR)
			blend(img, x, y, copper, t*t*0.18)
		}
	}
}

func paintGrain(img *image.RGBA) {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			n := (x*1973 + y*9277 + x*y*13) & 15
			if n > 1 {
				continue
			}
			delta := 8
			if n == 0 {
				delta = -10
			}
			shift(img, x, y, delta)
		}
	}
}

func blend(img *image.RGBA, x, y int, src color.NRGBA, t float64) {
	if t <= 0 {
		return
	}
	if t > 1 {
		t = 1
	}
	i := img.PixOffset(x, y)
	pix := img.Pix
	pix[i+0] = uint8(float64(pix[i+0])*(1-t) + float64(src.R)*t)
	pix[i+1] = uint8(float64(pix[i+1])*(1-t) + float64(src.G)*t)
	pix[i+2] = uint8(float64(pix[i+2])*(1-t) + float64(src.B)*t)
}

func shift(img *image.RGBA, x, y, delta int) {
	i := img.PixOffset(x, y)
	pix := img.Pix
	for c := 0; c < 3; c++ {
		v := int(pix[i+c]) + delta
		if v < 0 {
			v = 0
		}
		if v > 255 {
			v = 255
		}
		pix[i+c] = uint8(v)
	}
}
