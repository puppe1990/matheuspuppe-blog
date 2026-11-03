package ogcard

import (
	_ "embed"
	"fmt"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/Fraunces144pt-SemiBold.ttf
var frauncesTTF []byte

//go:embed fonts/PlusJakartaSans-Medium.ttf
var jakartaTTF []byte

var (
	fontsOnce    sync.Once
	frauncesFont *opentype.Font
	jakartaFont  *opentype.Font
	fontsErr     error
)

func loadFonts() error {
	fontsOnce.Do(func() {
		frauncesFont, fontsErr = opentype.Parse(frauncesTTF)
		if fontsErr != nil {
			fontsErr = fmt.Errorf("parse Fraunces: %w", fontsErr)
			return
		}
		jakartaFont, fontsErr = opentype.Parse(jakartaTTF)
		if fontsErr != nil {
			fontsErr = fmt.Errorf("parse Plus Jakarta Sans: %w", fontsErr)
		}
	})
	return fontsErr
}

func serifFace(size float64) (font.Face, error) {
	return opentype.NewFace(frauncesFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
}

func sansFace(size float64) (font.Face, error) {
	return opentype.NewFace(jakartaFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
}
