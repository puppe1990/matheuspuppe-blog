package store

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/testdata"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestStore_Migrations(t *testing.T) {
	_ = newTestStore(t)
}

func seedOneArticle(t *testing.T, s *SQLiteStore) models.Article {
	t.Helper()
	ctx := context.Background()
	f := testdata.New()
	author := testdata.Author(f, 1)
	id, err := s.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	author.ID = id

	article := testdata.Article(f, 1, author, time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC))
	article.Title = "A Nova Fronteira dos Agentes"
	article.Slug = "fronteira-agentes"
	article.Tags = []string{"IA", "Agentes"}
	article.Claps = 342
	article.Views = 4890
	article.Tier = 1
	id, err = s.UpsertArticle(ctx, article)
	if err != nil {
		t.Fatal(err)
	}
	article.ID = id
	return article
}

func TestUpsertAuthor_IdempotentBySlug(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	author := testdata.Author(testdata.New(), 1)

	id, err := s.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	author.Name = "Nome Atualizado"
	idAgain, err := s.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	if id != idAgain {
		t.Fatalf("ids differ: %d then %d", id, idAgain)
	}

	authors, err := s.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 || authors[0].Name != "Nome Atualizado" {
		t.Fatalf("authors = %+v", authors)
	}
}

func TestUpsertArticle_RoundTripsBlocksAndAuthor(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	seeded := seedOneArticle(t, s)

	got, err := s.FindArticleBySlug(ctx, "fronteira-agentes")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID == 0 || got.Title != seeded.Title {
		t.Fatalf("got %+v", got)
	}
	if got.Author.Name != seeded.Author.Name || got.Author.Slug != seeded.Author.Slug {
		t.Fatalf("author = %+v, want %+v", got.Author, seeded.Author)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "IA" {
		t.Fatalf("tags = %v", got.Tags)
	}
	if len(got.Blocks) != len(seeded.Blocks) || got.Blocks[0].Type != "paragraph" {
		t.Fatalf("blocks = %v", got.Blocks)
	}
	if got.Claps != 342 || got.Views != 4890 || got.Tier != 1 {
		t.Fatalf("counters = claps %d views %d tier %d", got.Claps, got.Views, got.Tier)
	}
	if !got.PublishedAt.Equal(seeded.PublishedAt) {
		t.Fatalf("publishedAt = %v, want %v", got.PublishedAt, seeded.PublishedAt)
	}
}

func TestFindArticleBySlug_Missing(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.FindArticleBySlug(context.Background(), "nao-existe"); err == nil {
		t.Fatal("expected error for missing slug")
	}
}

func TestUpsertArticle_IdempotentAndKeepsCounters(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	article := seedOneArticle(t, s)

	if err := s.IncrementArticleViews(ctx, article.ID, 1); err != nil {
		t.Fatal(err)
	}
	article.Title = "Título Revisado"
	article.Claps = 0 // a reseed must not zero the live counter
	article.Views = 0
	if _, err := s.UpsertArticle(ctx, article); err != nil {
		t.Fatal(err)
	}

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 1 {
		t.Fatalf("got %d articles, want 1", len(articles))
	}
	if articles[0].Title != "Título Revisado" {
		t.Fatalf("title = %q", articles[0].Title)
	}
	if articles[0].Claps != 342 || articles[0].Views != 4891 {
		t.Fatalf("counters = claps %d views %d, want 342/4891", articles[0].Claps, articles[0].Views)
	}
}

func TestListArticles_NewestFirst(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	f := testdata.New()
	author := testdata.Author(f, 1)
	id, err := s.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	author.ID = id

	base := time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC)
	for _, article := range testdata.Articles(f, 5, author, base) {
		if _, err := s.UpsertArticle(ctx, article); err != nil {
			t.Fatal(err)
		}
	}

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 5 {
		t.Fatalf("got %d articles, want 5", len(articles))
	}
	newest := base.AddDate(0, 0, -1) // faker article n=1 sits one day before base
	if !articles[0].PublishedAt.Equal(newest) {
		t.Fatalf("first = %v, want newest %v", articles[0].PublishedAt, newest)
	}
	for i := 1; i < len(articles); i++ {
		if articles[i].PublishedAt.After(articles[i-1].PublishedAt) {
			t.Fatalf("order broken at %d: %v before %v", i, articles[i-1].PublishedAt, articles[i].PublishedAt)
		}
	}
}

func TestListArticlesByAuthor(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	f := testdata.New()

	camila := testdata.Author(f, 1)
	other := testdata.Author(f, 2)
	camilaID, err := s.UpsertAuthor(ctx, camila)
	if err != nil {
		t.Fatal(err)
	}
	camila.ID = camilaID
	otherID, err := s.UpsertAuthor(ctx, other)
	if err != nil {
		t.Fatal(err)
	}
	other.ID = otherID

	base := time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC)
	mine := testdata.Articles(f, 2, camila, base)
	for i := range mine {
		if _, err := s.UpsertArticle(ctx, mine[i]); err != nil {
			t.Fatal(err)
		}
	}
	theirs := testdata.Articles(f, 1, other, base)
	if _, err := s.UpsertArticle(ctx, theirs[0]); err != nil {
		t.Fatal(err)
	}

	got, err := s.ListArticlesByAuthor(ctx, camilaID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d articles for the author, want 2", len(got))
	}
	for _, article := range got {
		if article.Author.Slug != camila.Slug {
			t.Fatalf("article %s belongs to %s", article.Slug, article.Author.Slug)
		}
	}
}

func TestListAuthors_CountsArticles(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	f := testdata.New()
	author := testdata.Author(f, 1)
	id, err := s.UpsertAuthor(ctx, author)
	if err != nil {
		t.Fatal(err)
	}
	author.ID = id
	base := time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC)
	for _, article := range testdata.Articles(f, 3, author, base) {
		if _, err := s.UpsertArticle(ctx, article); err != nil {
			t.Fatal(err)
		}
	}

	authors, err := s.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 || authors[0].ArticleCount != 3 {
		t.Fatalf("authors = %+v", authors)
	}
}

func TestFindAuthorBySlug(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	author := testdata.Author(testdata.New(), 1)
	if _, err := s.UpsertAuthor(ctx, author); err != nil {
		t.Fatal(err)
	}

	got, err := s.FindAuthorBySlug(ctx, author.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != author.Name {
		t.Fatalf("got %+v", got)
	}
	if _, err := s.FindAuthorBySlug(ctx, "nao-existe"); err == nil {
		t.Fatal("expected error for missing author")
	}
}

func TestIncrementArticleViews(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	article := seedOneArticle(t, s)

	for i := 0; i < 3; i++ {
		if err := s.IncrementArticleViews(ctx, article.ID, 1); err != nil {
			t.Fatal(err)
		}
	}
	got, err := s.FindArticleBySlug(ctx, article.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if got.Views != 4893 {
		t.Fatalf("views = %d, want 4893", got.Views)
	}
}

func TestAddArticleClaps_ReturnsNewTotal(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	article := seedOneArticle(t, s)

	total, err := s.AddArticleClaps(ctx, article.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if total != 343 {
		t.Fatalf("total = %d, want 343", total)
	}
	total, err = s.AddArticleClaps(ctx, article.ID, 5)
	if err != nil {
		t.Fatal(err)
	}
	if total != 348 {
		t.Fatalf("total = %d, want 348", total)
	}
}

func TestStore_NotFoundErrorsAreDiscoverable(t *testing.T) {
	s := newTestStore(t)
	_, err := s.FindArticleBySlug(context.Background(), "nao-existe")
	if err == nil || !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestMigrations_ContainEditorialTables(t *testing.T) {
	s := newTestStore(t)
	var name string
	err := s.DB().QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='articles'`).Scan(&name)
	if err != nil {
		t.Fatalf("articles table missing: %v", err)
	}
	if !strings.Contains(name, "articles") {
		t.Fatalf("name = %q", name)
	}
}
