package ogcard

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestCommittedMastheadPNG(t *testing.T) {
	path := filepath.Join(findModuleRoot(t), "web", "static", "og.png")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != Width || cfg.Height != Height {
		t.Fatalf("og.png = %dx%d, want %dx%d", cfg.Width, cfg.Height, Width, Height)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if limePixels(img) > 0 {
		t.Fatalf("og.png has lime pixels; regenerate with make og")
	}
	if !hasApprox(img, 22, 21, 19, 18) {
		t.Error("og.png missing ink")
	}
	if !hasApprox(img, 245, 158, 11, 24) {
		t.Error("og.png missing copper")
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("go.mod not found")
		}
		wd = parent
	}
}
