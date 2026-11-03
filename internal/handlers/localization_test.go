package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
)

// localeMiddleware mirrors the app wiring for direct handler tests: /en paths
// force lang=en (pathLocale) and i18n resolves the catalog onto the request.
func localeMiddleware() func(http.Handler) http.Handler {
	return i18n.LocaleMiddleware(map[string]*i18n.Catalog{
		"en": appi18n.NewCatalog("en"),
		"pt": appi18n.NewCatalog("pt"),
	}, "pt")
}

// localizedRequest marks the request as English like the /en routes do.
func localizedRequest(target, locale string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	query := req.URL.Query()
	query.Set("lang", locale)
	req.URL.RawQuery = query.Encode()
	return req
}

func serveLocalized(t *testing.T, h http.HandlerFunc, locale, target string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	localeMiddleware()(h).ServeHTTP(w, localizedRequest(target, locale))
	return w
}

func englishSeedStore(t *testing.T) store.Store {
	t.Helper()
	st := seedTestStore(t)
	seedFixtureCatalog(t, st)
	return st
}

func TestHome_EnglishEdition(t *testing.T) {
	st := englishSeedStore(t)
	h := NewHomeHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	rec := serveLocalized(t, h.ServeHTTP, "en", "/en")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"The First Essay", "Essays", "Essays on systems", "AI &amp; Models"} {
		if !strings.Contains(body, want) {
			t.Errorf("english home missing %q", want)
		}
	}
	// An essay without an English edition falls back to the Portuguese copy.
	if !strings.Contains(body, "Ensaio de Sistemas") {
		t.Errorf("english home lost the fallback essay")
	}
}

func TestArticle_EnglishSlug(t *testing.T) {
	st := englishSeedStore(t)
	h := NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	rec := serveLocalized(t, func(w http.ResponseWriter, r *http.Request) { h.Show(w, r, "first-essay") }, "en", "/en/artigo/first-essay")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"The First Essay", "Key points", "First takeaway.", "Published March 28, 2026"} {
		if !strings.Contains(body, want) {
			t.Errorf("english article missing %q", want)
		}
	}
	// The toggle points back at the Portuguese edition.
	if !strings.Contains(body, `href="/artigo/primeiro-ensaio"`) {
		t.Errorf("english article lost the PT toggle: %s", toggleSnippet(body))
	}
}

func TestArticle_LocaleSlugsAreStrict(t *testing.T) {
	st := englishSeedStore(t)
	h := NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	english := serveLocalized(t, func(w http.ResponseWriter, r *http.Request) { h.Show(w, r, "primeiro-ensaio") }, "en", "/en/artigo/primeiro-ensaio")
	if english.Code != http.StatusNotFound {
		t.Fatalf("portuguese slug under /en: status = %d, want 404", english.Code)
	}

	portuguese := serveLocalized(t, func(w http.ResponseWriter, r *http.Request) { h.Show(w, r, "first-essay") }, "pt", "/artigo/first-essay")
	if portuguese.Code != http.StatusNotFound {
		t.Fatalf("english slug under /: status = %d, want 404", portuguese.Code)
	}
}

func TestArticle_PortugueseTogglePointsToEnglish(t *testing.T) {
	st := englishSeedStore(t)
	h := NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	rec := serveLocalized(t, func(w http.ResponseWriter, r *http.Request) { h.Show(w, r, "primeiro-ensaio") }, "pt", "/artigo/primeiro-ensaio")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `href="/en/artigo/first-essay"`) {
		t.Fatalf("portuguese article lost the EN toggle: %s", toggleSnippet(rec.Body.String()))
	}
	if !strings.Contains(rec.Body.String(), `aria-current="true">PT</span>`) {
		t.Fatal("portuguese edition does not mark PT as the active language")
	}
	// Snippets cite the public repository, file and line included.
	if !strings.Contains(rec.Body.String(), `https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/app/app.go#L149-L163`) {
		t.Fatal("code block lost its GitHub citation")
	}
	if !strings.Contains(rec.Body.String(), `>internal/app/app.go:149</a>`) {
		t.Fatal("code block lost the source label")
	}
	// Code is highlighted on the server: chroma classes travel with the page.
	if !strings.Contains(rec.Body.String(), `class="kd">func`) {
		t.Fatal("code block rendered without syntax highlighting")
	}
}

func TestArticle_WithoutTranslationHasNoToggle(t *testing.T) {
	st := englishSeedStore(t)
	h := NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	rec := serveLocalized(t, func(w http.ResponseWriter, r *http.Request) { h.Show(w, r, "ensaio-de-sistemas") }, "pt", "/artigo/ensaio-de-sistemas")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), `href="/en/artigo/`) {
		t.Fatal("article without an English edition rendered a language toggle")
	}
}

func TestCommentPost_EnglishRedirectsToEnglishSlug(t *testing.T) {
	st := englishSeedStore(t)
	h := NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.NewCatalog("pt"), cais.Config{Env: "test"})

	form := url.Values{"nome": {"Reader"}, "cargo": {"Engineer"}, "conteudo": {"Great essay, thanks."}}
	req, rec := postForm(t, "/en/artigo/first-essay/comentarios", form)
	req.URL.RawQuery = "lang=en"

	localeMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { h.CommentPost(w, r, "first-essay") })).ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if location := rec.Header().Get("Location"); location != "/en/artigo/first-essay" {
		t.Fatalf("location = %q, want the English slug", location)
	}
}

func TestSitemap_ListsTheEnglishEdition(t *testing.T) {
	st := englishSeedStore(t)
	h := NewSitemapHandler(st, testSite(), cais.Config{Env: "test"})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))
	body := rec.Body.String()
	for _, want := range []string{"https://sintaxe.example.com/en", "/en/artigo/first-essay", "/en/sobre"} {
		if !strings.Contains(body, want) {
			t.Errorf("sitemap missing %q", want)
		}
	}
}

func TestSeeds_EnglishEditionRoundTrip(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	author, err := st.UpsertAuthor(ctx, models.Author{
		Slug:   "matheuspuppe",
		Name:   "Matheus Puppe",
		Role:   "Engenheiro",
		Bio:    "Bio.",
		RoleEN: "Engineer",
		BioEN:  "Bio.",
	})
	if err != nil {
		t.Fatal(err)
	}
	article := models.Article{
		Slug:            "ensaio-pt",
		Title:           "Ensaio",
		Subtitle:        "Sub.",
		LeadParagraph:   "Lead.",
		Category:        "IA & Modelos",
		ReadTimeMinutes: 5,
		PublishedAt:     day(2026, 3, 1, 9),
		Author:          models.Author{ID: author},
		Tags:            []string{"IA"},
		Blocks:          []models.ContentBlock{{Type: "paragraph", Text: "Corpo."}},
		SlugEN:          "essay-en",
		TitleEN:         "Essay",
		SubtitleEN:      "Sub.",
		LeadParagraphEN: "Lead.",
		TagsEN:          []string{"AI"},
		BlocksEN:        []models.ContentBlock{{Type: "paragraph", Text: "Body."}},
	}
	if _, err := st.UpsertArticle(ctx, article); err != nil {
		t.Fatal(err)
	}

	byEN, err := st.FindArticleBySlugEN(ctx, "essay-en")
	if err != nil {
		t.Fatalf("FindArticleBySlugEN: %v", err)
	}
	if byEN.Slug != "ensaio-pt" || byEN.TitleEN != "Essay" || len(byEN.BlocksEN) != 1 {
		t.Fatalf("english roundtrip = %+v", byEN)
	}
	if byEN.Author.RoleEN != "Engineer" || byEN.Author.BioEN != "Bio." {
		t.Fatalf("author english roundtrip = %+v", byEN.Author)
	}
	if _, err := st.FindArticleBySlug(ctx, "essay-en"); err == nil {
		t.Fatal("english slug resolved through the portuguese lookup")
	}
}

func TestLocalePathHelpers(t *testing.T) {
	if got := localePath("en", "/"); got != "/en" {
		t.Errorf("localePath root = %q", got)
	}
	if got := localePath("en", "/?q=fila"); got != "/en?q=fila" {
		t.Errorf("localePath with query = %q", got)
	}
	if got := localePath("pt", "/artigo/x"); got != "/artigo/x" {
		t.Errorf("localePath pt = %q", got)
	}
	if got := alternatePath("/en/artigo/epub-translation-llms-resilient-queue", "en"); got != "/artigo/epub-translation-llms-resilient-queue" {
		t.Errorf("alternatePath en→pt = %q", got)
	}
	if got := alternatePath("/sobre", "pt"); got != "/en/sobre" {
		t.Errorf("alternatePath pt→en = %q", got)
	}
	if got := alternatePath("/", "pt"); got != "/en" {
		t.Errorf("alternatePath root = %q", got)
	}
	if got := currentPath(httptest.NewRequest(http.MethodGet, "/en/artigo/x?lang=en", nil)); got != "/en/artigo/x" {
		t.Errorf("currentPath strips the injected lang: %q", got)
	}
	if got := currentPath(httptest.NewRequest(http.MethodGet, "/en?categoria=IA&lang=en", nil)); got != "/en?categoria=IA" {
		t.Errorf("currentPath keeps reader filters: %q", got)
	}
	if got := currentPath(httptest.NewRequest(http.MethodGet, "/salvos?q=fila", nil)); got != "/salvos?q=fila" {
		t.Errorf("currentPath pt = %q", got)
	}
}

func toggleSnippet(body string) string {
	index := strings.Index(body, `hreflang="`)
	if index < 0 {
		return "(no alternate link)"
	}
	end := index + 200
	if end > len(body) {
		end = len(body)
	}
	return body[index:end]
}
