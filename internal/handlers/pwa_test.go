package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
)

// TestHome_PWAAssetsAreVersioned protege o cache-bust: /static responde com
// max-age de um ano, então uma troca de ícone só chega ao navegador (e ao PWA
// instalado) se a URL mudar.
func TestHome_PWAAssetsAreVersioned(t *testing.T) {
	st := seedTestStore(t)
	h := NewHomeHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	body := rec.Body.String()
	for _, want := range []string{
		`rel="manifest" href="/static/manifest.webmanifest?v=2"`,
		`rel="apple-touch-icon" href="/static/icons/icon.png?v=2"`,
		`href="/static/favicon.svg?v=2"`,
		`href="/static/icons/favicon-32.png?v=2"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("head without the versioned PWA asset %q", want)
		}
	}
}
