package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
)

func getPage(t *testing.T, handler http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func TestSobreHandler_RendersManifestoAndAuthors(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	h := NewSobreHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})
	rr := getPage(t, h, "/sobre")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	for _, needle := range []string{
		"Engenharia com profundidade",
		"Quem escreve",
		"Matheus Puppe",
		"ensaios",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("body missing %q", needle)
		}
	}
}

func TestAuthorHandler_RendersColumnist(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	h := NewAuthorHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})
	rr := httptest.NewRecorder()
	h.Serve(rr, httptest.NewRequest(http.MethodGet, "/autor/matheus-puppe", nil), "matheus-puppe")

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Matheus Puppe") {
		t.Errorf("body missing the columnist name")
	}
	if !strings.Contains(body, "O Primeiro Ensaio") {
		t.Errorf("body missing the columnist essay")
	}
}

func TestAuthorHandler_MissingIs404(t *testing.T) {
	st := seedTestStore(t)
	h := NewAuthorHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})
	rr := httptest.NewRecorder()
	h.Serve(rr, httptest.NewRequest(http.MethodGet, "/autor/nao-existe", nil), "nao-existe")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestNotFoundHandler_RendersDesigned404(t *testing.T) {
	h := NewNotFoundHandler(setupTestViews(t), testSite(), cais.Config{})
	rr := getPage(t, h, "/qualquer-coisa")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Esta página não consta na edição") {
		t.Errorf("body missing the designed 404 copy")
	}
}

func TestNewsletterHandler_Subscribes(t *testing.T) {
	st := setupTestStore(t)
	h := NewNewsletterHandler(st, cais.Config{})
	req, rr := postForm(t, "/newsletter", url.Values{"email": {"Leitor@Sintaxe.tech"}, "next": {"/sobre"}})
	h.Post(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/sobre" {
		t.Fatalf("Location = %q", got)
	}
	count, err := st.CountSubscribers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("subscribers = %d, want 1", count)
	}
}

func TestNewsletterHandler_RejectsInvalidEmail(t *testing.T) {
	st := setupTestStore(t)
	h := NewNewsletterHandler(st, cais.Config{})
	req, rr := postForm(t, "/newsletter", url.Values{"email": {"nao-e-email"}})
	h.Post(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/" {
		t.Fatalf("Location = %q, want sanitized fallback", got)
	}
	count, _ := st.CountSubscribers(context.Background())
	if count != 0 {
		t.Fatalf("invalid email was stored")
	}
}

func TestSitemapHandler_ServesCrawlersWithoutCookies(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	h := NewSitemapHandler(st, testSite(), cais.Config{})
	rr := getPage(t, h, "/sitemap.xml")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "xml") {
		t.Fatalf("Content-Type = %q", ct)
	}
	if len(rr.Result().Cookies()) != 0 {
		t.Fatalf("sitemap must not set cookies (#342)")
	}
	body := rr.Body.String()
	if !strings.Contains(body, "https://sintaxe.example.com/artigo/primeiro-ensaio") {
		t.Errorf("sitemap missing an article URL")
	}
	if !strings.Contains(body, "<lastmod>2026-03-28</lastmod>") {
		t.Errorf("sitemap missing lastmod")
	}
}

func TestRobotsHandler_PointsToSitemap(t *testing.T) {
	rr := getPage(t, RobotsHandler(testSite()), "/robots.txt")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Sitemap: https://sintaxe.example.com/sitemap.xml") {
		t.Errorf("robots.txt missing the sitemap pointer: %q", body)
	}
	if len(rr.Result().Cookies()) != 0 {
		t.Fatalf("robots.txt must not set cookies (#342)")
	}
}

func TestFeedHandler_ServesRSS(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	h := NewFeedHandler(st, testSite(), cais.Config{})
	rr := getPage(t, h, "/rss.xml")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "rss") {
		t.Fatalf("Content-Type = %q", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<title>O Primeiro Ensaio</title>") {
		t.Errorf("feed missing the lead essay item")
	}
	if strings.Count(body, "<item>") != 6 {
		t.Errorf("feed items = %d, want 6", strings.Count(body, "<item>"))
	}
}

func TestSitemapAndFeed_HideScheduledEssays(t *testing.T) {
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	scheduled := seedScheduledEssay(t, st)

	sitemap := NewSitemapHandler(st, testSite(), cais.Config{})
	if body := getPage(t, sitemap, "/sitemap.xml").Body.String(); strings.Contains(body, scheduled.Slug) {
		t.Errorf("sitemap leaked the scheduled essay")
	}

	feed := NewFeedHandler(st, testSite(), cais.Config{})
	if body := getPage(t, feed, "/rss.xml").Body.String(); strings.Contains(body, "Ensaio Agendado Confidencial") {
		t.Errorf("feed leaked the scheduled essay")
	}

	// And the librarian's page stays empty until the day arrives.
	author := NewAuthorHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{})
	rr := httptest.NewRecorder()
	author.Serve(rr, httptest.NewRequest(http.MethodGet, "/autor/"+scheduled.Author.Slug, nil), scheduled.Author.Slug)
	if rr.Code != http.StatusOK {
		t.Fatalf("author page = %d, want 200", rr.Code)
	}
	if body := rr.Body.String(); strings.Contains(body, "Ensaio Agendado Confidencial") {
		t.Errorf("author page leaked the scheduled essay")
	}
}
