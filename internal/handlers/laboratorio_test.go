package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
)

func getLaboratorio(t *testing.T, target string) string {
	t.Helper()
	h := NewLaboratorioHandler(setupTestViews(t), testSite(), appi18n.DefaultCatalog(), cais.Config{})
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET %s = %d, want 200", target, rr.Code)
	}
	return rr.Body.String()
}

func TestLaboratorioHandler_Defaults(t *testing.T) {
	body := getLaboratorio(t, "/laboratorio")
	for _, needle := range []string{
		"Simulador de runtimes",
		"Requisições por segundo",
		"Rust (Tokio",
		"Go (Gin / net/http)",
		"Node.js (Fastify)",
		"Python (FastAPI",
		"25.000",
		"6.127",
		"Sobre o modelo",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("body missing %q", needle)
		}
	}
}

func TestLaboratorioHandler_QueryDrivesTheSimulation(t *testing.T) {
	body := getLaboratorio(t, "/laboratorio?reqs=100000&payload=32")
	if !strings.Contains(body, "100.000") {
		t.Errorf("reqs label not applied")
	}
	if !strings.Contains(body, "32 KB") {
		t.Errorf("payload label not applied")
	}
	// Rust p99 at max load: 0.3 + 0.8 + 0.5 = 1.6 ms.
	if !strings.Contains(body, "1.6 ms") {
		t.Errorf("rust p99 missing at max load")
	}
}

func TestLaboratorioHandler_ClampsOutOfRange(t *testing.T) {
	body := getLaboratorio(t, "/laboratorio?reqs=999999&payload=0")
	if !strings.Contains(body, "100.000") {
		t.Errorf("reqs out of range not clamped to max")
	}
	if !strings.Contains(body, "1 KB") {
		t.Errorf("payload out of range not clamped to min")
	}
}
