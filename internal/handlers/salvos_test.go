package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

func newSavedHandler(t *testing.T) (*SavedHandler, store.Store) {
	t.Helper()
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	return NewSavedHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{}), st
}

func getSaved(t *testing.T, target string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	h, _ := newSavedHandler(t)
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestSavedHandler_EmptyLibrary(t *testing.T) {
	rr := getSaved(t, "/salvos")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Artigos salvos") {
		t.Errorf("missing saved counter header")
	}
	if !strings.Contains(body, "Nenhum artigo salvo") {
		t.Errorf("missing empty state")
	}
}

func TestSavedHandler_ShowsOnlySavedArticles(t *testing.T) {
	cookie := &http.Cookie{Name: bookmarkCookie, Value: "ensaio-de-seguranca"}
	rr := getSaved(t, "/salvos", cookie)
	body := rr.Body.String()
	if !strings.Contains(body, "Ensaio de Segurança") {
		t.Errorf("saved essay missing from the library")
	}
	if strings.Contains(body, "O Primeiro Ensaio") {
		t.Errorf("unsaved essay leaked into the library")
	}
	if !strings.Contains(body, `data-saved="true"`) {
		t.Errorf("saved card not marked")
	}
}

func TestSavedHandler_SearchWithinLibrary(t *testing.T) {
	cookie := &http.Cookie{Name: bookmarkCookie, Value: "ensaio-de-seguranca,ensaio-de-opiniao"}
	rr := getSaved(t, "/salvos?q=opini%C3%A3o", cookie)
	body := rr.Body.String()
	if !strings.Contains(body, "Ensaio de Opinião") {
		t.Errorf("search miss inside the library")
	}
	if strings.Contains(body, "Ensaio de Segurança") {
		t.Errorf("search leaked a non-matching saved essay")
	}
}

func TestSavedHandler_ToggleAddsAndRemoves(t *testing.T) {
	h, _ := newSavedHandler(t)
	slug := "ensaio-de-seguranca"

	req, rr := postForm(t, "/salvos/"+slug, url.Values{"next": {"/?q=rust"}})
	h.Toggle(rr, req, slug)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/?q=rust" {
		t.Fatalf("Location = %q, want the sanitized next", got)
	}
	cookies := rr.Result().Cookies()
	if len(cookies) == 0 || !strings.Contains(cookies[0].Value, slug) {
		t.Fatalf("cookie missing the toggled slug: %+v", cookies)
	}
	if !cookies[0].HttpOnly {
		t.Fatalf("bookmark cookie must be HttpOnly")
	}

	// Toggling again with the saved cookie removes it.
	req2, rr2 := postForm(t, "/salvos/"+slug, url.Values{}, cookies[0])
	h.Toggle(rr2, req2, slug)
	after := rr2.Result().Cookies()
	if len(after) == 0 || strings.Contains(after[0].Value, slug) {
		t.Fatalf("second toggle did not remove the slug: %+v", after)
	}
}

func TestSavedHandler_ToggleUnknownSlugIs404(t *testing.T) {
	h, _ := newSavedHandler(t)
	req, rr := postForm(t, "/salvos/nao-existe", url.Values{})
	h.Toggle(rr, req, "nao-existe")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestSavedHandler_ToggleRejectsCrossSiteNext(t *testing.T) {
	h, _ := newSavedHandler(t)
	req, rr := postForm(t, "/salvos/"+leadSlug, url.Values{"next": {"//evil.example/phish"}})
	h.Toggle(rr, req, leadSlug)
	if got := rr.Header().Get("Location"); got != "/" {
		t.Fatalf("Location = %q, want / for a cross-site next", got)
	}
}
