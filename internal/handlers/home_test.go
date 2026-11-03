package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
)

func newHomeHandler(t *testing.T) *HomeHandler {
	t.Helper()
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	return NewHomeHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})
}

func getHome(t *testing.T, target string, cookies ...*http.Cookie) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	rr := httptest.NewRecorder()
	newHomeHandler(t).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET %s = %d, want 200", target, rr.Code)
	}
	return rr.Body.String()
}

func TestHomeHandler_Returns200AndShell(t *testing.T) {
	body := getHome(t, "/")
	if !strings.Contains(body, `id="amarra-main"`) {
		t.Errorf("body missing #amarra-main")
	}
	if !strings.Contains(body, "SINTAXE") {
		t.Errorf("body missing the brand wordmark, got: %.200s", body)
	}
}

func TestHomeHandler_RendersEditorialTiers(t *testing.T) {
	body := getHome(t, "/")
	// Lead (tier 1), secondary (tier 2) and grid (tier 3) essays all render.
	for _, title := range []string{
		"O Primeiro Ensaio",
		"Ensaio de Sistemas",
		"Ensaio de Hardware",
		"Ensaio de Design",
		"Ensaio de Segurança",
		"Ensaio de Opinião",
	} {
		if !strings.Contains(body, title) {
			t.Errorf("body missing %q", title)
		}
	}
	if !strings.Contains(body, "6 Ensaios Ativos") {
		t.Errorf("body missing the active-essay count")
	}
}

func TestHomeHandler_CategoryFilter(t *testing.T) {
	body := getHome(t, "/?categoria=Hardware+%26+Chips")
	if !strings.Contains(body, "Ensaio de Hardware") {
		t.Errorf("filtered body missing the hardware essay")
	}
	if strings.Contains(body, "Ensaio de Sistemas") {
		t.Errorf("filtered body must not carry other cadernos")
	}
}

func TestHomeHandler_SearchFilter(t *testing.T) {
	body := getHome(t, "/?q=criptografia")
	if !strings.Contains(body, "Ensaio de Segurança") {
		t.Errorf("search miss for the post-quantum essay")
	}
	if strings.Contains(body, "O Primeiro Ensaio") {
		t.Errorf("search leaked an unrelated essay")
	}
}

func TestHomeHandler_TimeFilter(t *testing.T) {
	body := getHome(t, "/?tempo=5")
	if !strings.Contains(body, "Ensaio de Design") {
		t.Errorf("≤5min must include the 5-minute essay")
	}
	if strings.Contains(body, "Ensaio de Sistemas") {
		t.Errorf("≤5min must exclude the 6-minute essay")
	}
}

func TestHomeHandler_CategoryCountsShown(t *testing.T) {
	body := getHome(t, "/")
	if !strings.Contains(body, ">(6)<") {
		t.Errorf("body missing the Todos count (6)")
	}
}

func TestHomeHandler_EmptyState(t *testing.T) {
	body := getHome(t, "/?q=nao-existe-nada")
	if !strings.Contains(body, "Nenhuma publicação encontrada") {
		t.Errorf("body missing the empty state")
	}
}

func TestHomeHandler_FirstEditionEmptyState(t *testing.T) {
	st := seedTestStore(t) // only the columnist, no essays yet
	h := NewHomeHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "A primeira edição está em preparação") {
		t.Errorf("first-edition empty state missing")
	}
	if !strings.Contains(body, "0 Ensaios Ativos") {
		t.Errorf("body missing the zero essay counter")
	}
}

func TestHomeHandler_SavedStateMarksCards(t *testing.T) {
	cookie := &http.Cookie{Name: bookmarkCookie, Value: "primeiro-ensaio"}
	body := getHome(t, "/", cookie)
	if !strings.Contains(body, `data-saved="true"`) {
		t.Errorf("saved cookie did not mark any card")
	}
}

func TestHomeHandler_HidesScheduledEssays(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	seedScheduledEssay(t, st)
	h := NewHomeHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if strings.Contains(body, "Ensaio Agendado Confidencial") {
		t.Errorf("scheduled essay visible in the feed before its day")
	}
	if !strings.Contains(body, "6 Ensaios Ativos") {
		t.Errorf("essay counter must ignore scheduled essays (want 6)")
	}
}
