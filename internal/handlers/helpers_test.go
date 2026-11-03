package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
	"github.com/puppe1990/matheuspuppe-blog/internal/testdata"
)

func testSite() meta.Site {
	return meta.Site{AppName: "Sintaxe", AppURL: "https://sintaxe.example.com"}
}

func projectRoot(t *testing.T) string {
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

func setupTestViews(t *testing.T) *view.Renderer {
	t.Helper()
	root := projectRoot(t)
	fsys := os.DirFS(filepath.Join(root, "web", "templates"))
	r, err := view.Load(fsys, appi18n.NewCatalog("pt"), appi18n.NewCatalog("en"), appi18n.NewCatalog("pt"))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func setupTestStore(t *testing.T) store.Store {
	t.Helper()
	s, err := store.NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// seedTestStore plants only the columnist; handler tests render their own
// fixtures, so assertions never depend on production editorial copy.
func seedTestStore(t *testing.T) store.Store {
	t.Helper()
	s := setupTestStore(t)
	if _, err := s.UpsertAuthor(context.Background(), models.Author{
		Slug:      "matheus-puppe",
		Name:      "Matheus Puppe",
		Role:      "Engenheiro de software & fundador da Gestão Bem",
		AvatarURL: "/static/img/author-matheus.jpg",
	}); err != nil {
		t.Fatal(err)
	}
	return s
}

// seedScheduledEssay stores one future-dated essay under a distinctive title,
// the fixture for the scheduling tests.
func seedScheduledEssay(t *testing.T, st store.Store) models.Article {
	t.Helper()
	ctx := context.Background()
	f := testdata.New()
	author := testdata.Author(f, 9)
	authorID, err := st.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	author.ID = authorID

	article := testdata.Article(f, 9, author, time.Now().UTC().AddDate(0, 0, 3))
	article.PublishedAt = time.Now().UTC().AddDate(0, 0, 3)
	article.Title = "Ensaio Agendado Confidencial"
	if _, err := st.UpsertArticle(ctx, article); err != nil {
		t.Fatal(err)
	}
	return article
}

func day(year int, month time.Month, dayOfMonth, hour int) time.Time {
	return time.Date(year, month, dayOfMonth, hour, 0, 0, 0, time.UTC)
}

// seedFixtureCatalog gives handler tests six deterministic essays (one per
// caderno) signed by the seeded columnist, plus the two starter comments on the
// lead. The real catalog stays empty until the first edition is written.
func seedFixtureCatalog(t *testing.T, st store.Store) models.Article {
	t.Helper()
	ctx := context.Background()
	author, err := st.FindAuthorBySlug(ctx, "matheus-puppe")
	if err != nil {
		t.Fatalf("fixture catalog needs the seeded author: %v", err)
	}

	articles := []models.Article{
		{
			Slug:            "primeiro-ensaio",
			Title:           "O Primeiro Ensaio",
			Subtitle:        "Uma dissecação do sistema, do cais ao mar aberto.",
			LeadParagraph:   "O primeiro parágrafo do ensaio principal.",
			Category:        "IA & Modelos",
			ReadTimeMinutes: 8,
			PublishedAt:     day(2026, time.March, 28, 9),
			Author:          author,
			Tags:            []string{"Inteligência Artificial", "Arquitetura"},
			Tier:            1,
			Claps:           342,
			Views:           4890,
			SlugEN:          "first-essay",
			TitleEN:         "The First Essay",
			SubtitleEN:      "A dissection of the system, from the pier to the open sea.",
			LeadParagraphEN: "The first paragraph of the main essay.",
			TagsEN:          []string{"Artificial Intelligence", "Architecture"},
			BlocksEN: []models.ContentBlock{
				{Type: "takeaways", Items: []string{
					"First takeaway.",
					"Second takeaway.",
				}},
				{Type: "heading2", Text: "The Central Problem"},
				{Type: "paragraph", Text: "A body paragraph with a technical argument."},
				{Type: "quote", Text: "The right question is worth more than the quick answer.", AuthorQuote: "Matheus Puppe"},
			},
			Blocks: []models.ContentBlock{
				{Type: "takeaways", Items: []string{
					"Primeira conclusão do ensaio.",
					"Segunda conclusão do ensaio.",
					"Terceira conclusão do ensaio.",
				}},
				{Type: "heading2", Text: "O Problema Central"},
				{Type: "paragraph", Text: "Um parágrafo de corpo com argumento técnico."},
				{Type: "quote", Text: "A pergunta certa vale mais que a resposta rápida.", AuthorQuote: "Matheus Puppe"},
				{Type: "code", Language: "go", Code: "func main() {}", Caption: "Um exemplo em Go",
					SourceURL:   "https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/app/app.go#L149-L163",
					SourceLabel: "internal/app/app.go:149"},
				{Type: "interactive-benchmark", Benchmark: &models.BenchmarkData{
					Title:       "Benchmark do Ensaio",
					Description: "Comparativo entre abordagens.",
					Metrics: []models.BenchmarkMetric{
						{Label: "Precisão Lógica", ValueA: 42, ValueB: 88, Unit: "%"},
					},
				}},
			},
		},
		{Slug: "ensaio-de-sistemas", Title: "Ensaio de Sistemas", Subtitle: "Fila, retry e previsibilidade.", LeadParagraph: "Um sistema previsível começa por métricas honestas.", Category: "Engenharia de Sistemas", ReadTimeMinutes: 6, PublishedAt: day(2026, time.March, 24, 14), Author: author, Tags: []string{"Sistemas"}, Tier: 2},
		{Slug: "ensaio-de-hardware", Title: "Ensaio de Hardware", Subtitle: "Telemetria de um adaptador flaky.", LeadParagraph: "Chips e monitores têm histórias.", Category: "Hardware & Chips", ReadTimeMinutes: 7, PublishedAt: day(2026, time.March, 20, 11), Author: author, Tags: []string{"Hardware"}, Tier: 2},
		{Slug: "ensaio-de-design", Title: "Ensaio de Design", Subtitle: "Slop visual e tipografia.", LeadParagraph: "Interfaces com intenção.", Category: "Design & UX", ReadTimeMinutes: 5, PublishedAt: day(2026, time.March, 16, 16), Author: author, Tags: []string{"Design"}, Tier: 3},
		{Slug: "ensaio-de-seguranca", Title: "Ensaio de Segurança", Subtitle: "Criptografia pós-quântica na prática.", LeadParagraph: "Chaves maiores mudam o handshake.", Category: "Segurança & Criptografia", ReadTimeMinutes: 6, PublishedAt: day(2026, time.March, 11, 13), Author: author, Tags: []string{"Segurança", "Criptografia"}, Tier: 3},
		{Slug: "ensaio-de-opiniao", Title: "Ensaio de Opinião", Subtitle: "O ofício de escrever ensaios.", LeadParagraph: "Opinião com profundidade.", Category: "Opinião & Futuro", ReadTimeMinutes: 5, PublishedAt: day(2026, time.March, 5, 8), Author: author, Tags: []string{"Carreira"}, Tier: 3},
	}
	for i := range articles {
		articles[i].CoverImage = "/static/img/tech_ai_agents_lead.jpg"
		if _, err := st.UpsertArticle(ctx, articles[i]); err != nil {
			t.Fatalf("fixture article %q: %v", articles[i].Slug, err)
		}
	}

	lead, err := st.FindArticleBySlug(ctx, "primeiro-ensaio")
	if err != nil {
		t.Fatalf("fixture lead: %v", err)
	}
	starter := []models.Comment{
		{ArticleID: lead.ID, AuthorName: "Leitor Um", AuthorRole: "Engenheiro de ML", Content: "Primeiro comentário da comunidade.", CreatedAt: day(2026, time.March, 28, 14), Claps: 18},
		{ArticleID: lead.ID, AuthorName: "Leitor Dois", AuthorRole: "Tech Lead", Content: "Segundo comentário da comunidade.", CreatedAt: day(2026, time.March, 29, 10), Claps: 9},
	}
	for _, comment := range starter {
		if _, err := st.AddComment(ctx, comment); err != nil {
			t.Fatalf("fixture comment: %v", err)
		}
	}
	return lead
}

// postForm builds the request and the recorder, so handler tests can call
// handlers directly without the router wiring.
func postForm(t *testing.T, target string, form url.Values, cookies ...*http.Cookie) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return req, httptest.NewRecorder()
}
