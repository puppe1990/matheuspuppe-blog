package handlers

import (
	"bytes"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	"github.com/puppe1990/matheuspuppe-blog/internal/ogcard"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

func newOGHandler(t *testing.T) (*OGHandler, store.Store) {
	t.Helper()
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	return NewOGHandler(st, testSite(), cais.Config{}), st
}

func TestOGHandler_ServesArticlePNG(t *testing.T) {
	h, _ := newOGHandler(t)
	rr := httptest.NewRecorder()
	h.Serve(rr, httptest.NewRequest(http.MethodGet, "/og/"+leadSlug, nil), leadSlug)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if len(rr.Result().Cookies()) != 0 {
		t.Fatal("og image must not set cookies")
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(rr.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Width != ogcard.Width || cfg.Height != ogcard.Height {
		t.Fatalf("png = %dx%d, want %dx%d", cfg.Width, cfg.Height, ogcard.Width, ogcard.Height)
	}
}

func TestOGHandler_MissingIs404(t *testing.T) {
	h, _ := newOGHandler(t)
	rr := httptest.NewRecorder()
	h.Serve(rr, httptest.NewRequest(http.MethodGet, "/og/nao-existe", nil), "nao-existe")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestOGHandler_ScheduledIs404(t *testing.T) {
	h, st := newOGHandler(t)
	scheduled := seedScheduledEssay(t, st)
	rr := httptest.NewRecorder()
	h.Serve(rr, httptest.NewRequest(http.MethodGet, "/og/"+scheduled.Slug, nil), scheduled.Slug)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("scheduled og = %d, want 404", rr.Code)
	}
}

func TestArticleHandler_OpenGraphPointsAtTitleCard(t *testing.T) {
	body := showArticle(t, leadSlug).Body.String()
	wantImage := `content="https://sintaxe.example.com/og/` + leadSlug + `"`
	if !strings.Contains(body, wantImage) {
		t.Errorf("article og:image missing %s", wantImage)
	}
	if !strings.Contains(body, `property="og:type" content="article"`) {
		t.Error("article must set og:type=article")
	}
	if !strings.Contains(body, `property="og:image:width" content="1200"`) {
		t.Error("missing og:image:width")
	}
	if !strings.Contains(body, `name="twitter:image" content="https://sintaxe.example.com/og/`+leadSlug+`"`) {
		t.Error("missing twitter:image pointing at the title card")
	}
	if strings.Contains(body, "/static/img/tech_ai_agents_lead.jpg") && strings.Contains(body, `property="og:image"`) {
		if strings.Contains(body, `property="og:image" content="https://sintaxe.example.com/static/img/tech_ai_agents_lead.jpg"`) {
			t.Error("og:image still uses the cover photo instead of the title card")
		}
	}
}

func TestHomeHandler_OpenGraphUsesMastheadPNG(t *testing.T) {
	body := getHome(t, "/")
	if !strings.Contains(body, `property="og:image" content="https://sintaxe.example.com/static/og.png"`) {
		t.Error("home og:image must be the masthead PNG")
	}
	if !strings.Contains(body, `property="og:image:width" content="1200"`) {
		t.Error("home missing og:image:width")
	}
	if !strings.Contains(body, `name="twitter:card" content="summary_large_image"`) {
		t.Error("home missing twitter:card")
	}
	if !strings.Contains(body, `name="twitter:image" content="https://sintaxe.example.com/static/og.png"`) {
		t.Error("home missing twitter:image")
	}
}
