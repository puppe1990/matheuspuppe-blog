package ogcard

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func TestRender_SiteCardIsOpenGraphSize(t *testing.T) {
	img, err := Render(SiteCard("sintaxe.apps.gestaobem.com"))
	if err != nil {
		t.Fatal(err)
	}
	b := img.Bounds()
	if b.Dx() != Width || b.Dy() != Height {
		t.Fatalf("site card = %dx%d, want %dx%d", b.Dx(), b.Dy(), Width, Height)
	}
}

func TestRender_ArticleCardDiffersFromSiteByTitle(t *testing.T) {
	site, err := Render(SiteCard("sintaxe.example.com"))
	if err != nil {
		t.Fatal(err)
	}
	article, err := Render(ArticleCard(models.Article{
		Title:           "O Primeiro Ensaio",
		Category:        "IA & Modelos",
		ReadTimeMinutes: 8,
		Author:          models.Author{Name: "Matheus Puppe"},
	}, "sintaxe.example.com"))
	if err != nil {
		t.Fatal(err)
	}
	if pngEqual(t, site, article) {
		t.Fatal("article card must paint the title, not the site masthead")
	}
}

func TestRender_DifferentTitlesPaintDifferentCards(t *testing.T) {
	a, err := Render(Card{Title: "Fila por capítulo", Kicker: "Sistemas", Host: "sintaxe.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Render(Card{Title: "Chips quânticos", Kicker: "Hardware", Host: "sintaxe.example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if pngEqual(t, a, b) {
		t.Fatal("cards with different titles must not be identical")
	}
}

func TestRender_UsesJournalPalette(t *testing.T) {
	img, err := Render(SiteCard("sintaxe.apps.gestaobem.com"))
	if err != nil {
		t.Fatal(err)
	}
	if !hasApprox(img, 22, 21, 19, 18) {
		t.Error("missing ink #161513")
	}
	if !hasApprox(img, 239, 233, 222, 24) {
		t.Error("missing foam #efe9de")
	}
	if !hasApprox(img, 245, 158, 11, 24) {
		t.Error("missing copper #f59e0b")
	}
}

func TestRender_HasNoLimePremultiplyArtifact(t *testing.T) {
	img, err := Render(SiteCard("sintaxe.apps.gestaobem.com"))
	if err != nil {
		t.Fatal(err)
	}
	if limePixels(img) > 0 {
		t.Fatalf("found %d lime pixels — RGBA alpha was premultiplied", limePixels(img))
	}
}

func TestRender_LongPortugueseTitleDoesNotPanic(t *testing.T) {
	img, err := Render(Card{
		Title:  "Tradução de EPUBs com LLMs: fila resiliente, retry por capítulo e modelos baratos",
		Kicker: "Engenharia de Sistemas",
		Byline: "Matheus Puppe  ·  8 min",
		Host:   "sintaxe.apps.gestaobem.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != Width {
		t.Fatalf("width = %d", img.Bounds().Dx())
	}
}

func TestHostOf_UsesHostname(t *testing.T) {
	if got := HostOf("https://sintaxe.apps.gestaobem.com"); got != "sintaxe.apps.gestaobem.com" {
		t.Fatalf("HostOf = %q", got)
	}
	if got := HostOf(""); got != defaultHost {
		t.Fatalf("empty HostOf = %q", got)
	}
}

func TestWritePNG_Encodes1200x630(t *testing.T) {
	var buf bytes.Buffer
	if err := WritePNG(&buf, SiteCard("sintaxe.apps.gestaobem.com")); err != nil {
		t.Fatal(err)
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != Width || cfg.Height != Height {
		t.Fatalf("png = %dx%d, want %dx%d", cfg.Width, cfg.Height, Width, Height)
	}
}

func pngEqual(t *testing.T, a, b image.Image) bool {
	t.Helper()
	ab, bb := a.Bounds(), b.Bounds()
	if ab != bb {
		return false
	}
	for y := ab.Min.Y; y < ab.Max.Y; y += 4 {
		for x := ab.Min.X; x < ab.Max.X; x += 4 {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	return true
}

func hasApprox(img image.Image, r, g, b, slop uint8) bool {
	bounds := img.Bounds()
	want := color.RGBA{R: r, G: g, B: b, A: 255}
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 3 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 3 {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			if pa>>8 < 200 {
				continue
			}
			got := color.RGBA{R: uint8(pr >> 8), G: uint8(pg >> 8), B: uint8(pb >> 8), A: 255}
			if near(got.R, want.R, slop) && near(got.G, want.G, slop) && near(got.B, want.B, slop) {
				return true
			}
		}
	}
	return false
}

func limePixels(img image.Image) int {
	bounds := img.Bounds()
	n := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			if pa>>8 < 200 {
				continue
			}
			r, g, b := uint8(pr>>8), uint8(pg>>8), uint8(pb>>8)
			if g > 140 && r < 80 && b < 80 {
				n++
			}
		}
	}
	return n
}

func near(got, want, slop uint8) bool {
	d := int(got) - int(want)
	if d < 0 {
		d = -d
	}
	return d <= int(slop)
}
