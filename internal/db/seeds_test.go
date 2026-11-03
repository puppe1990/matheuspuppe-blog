package db

import (
	"context"
	"testing"

	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

func newTestStore(t *testing.T) *store.SQLiteStore {
	t.Helper()
	s, err := store.NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestRunSeeds_PublishesTheCatalog(t *testing.T) {
	s := newTestStore(t)
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 6 {
		t.Fatalf("got %d articles, want 6", len(articles))
	}
	// Newest first: the lead essay is the March 28 agents piece.
	if articles[0].Slug != "fronteira-modelos-raciocinio-agentes-autonomos" {
		t.Fatalf("first article = %q", articles[0].Slug)
	}
	if articles[0].Author.Name != "Dra. Camila Siqueira" {
		t.Fatalf("author join = %+v", articles[0].Author)
	}
	if len(articles[0].Blocks) == 0 || len(articles[0].Tags) == 0 {
		t.Fatalf("article content missing: blocks %d tags %d", len(articles[0].Blocks), len(articles[0].Tags))
	}

	authors, err := s.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 6 {
		t.Fatalf("got %d authors, want 6", len(authors))
	}
	for _, author := range authors {
		if author.ArticleCount == 0 {
			t.Fatalf("author %q has no articles", author.Slug)
		}
	}

	lead, err := s.FindArticleBySlug(ctx, articles[0].Slug)
	if err != nil {
		t.Fatal(err)
	}
	comments, err := s.ListComments(ctx, lead.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Fatalf("got %d comments on the lead essay, want 2", len(comments))
	}
	rust, err := s.FindArticleBySlug(ctx, "go-para-rust-sistemas-criticos-2026")
	if err != nil {
		t.Fatal(err)
	}
	comments, err = s.ListComments(ctx, rust.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 {
		t.Fatalf("got %d comments on the rust essay, want 1", len(comments))
	}
}

func TestRunSeeds_IdempotentAndKeepsCounters(t *testing.T) {
	s := newTestStore(t)
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	lead, err := s.FindArticleBySlug(ctx, "fronteira-modelos-raciocinio-agentes-autonomos")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddArticleClaps(ctx, lead.ID, 10); err != nil {
		t.Fatal(err)
	}
	if err := s.IncrementArticleViews(ctx, lead.ID, 99); err != nil {
		t.Fatal(err)
	}

	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	leadAfter, err := s.FindArticleBySlug(ctx, lead.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if leadAfter.Claps != lead.Claps+10 || leadAfter.Views != lead.Views+99 {
		t.Fatalf("counters reset: claps %d views %d", leadAfter.Claps, leadAfter.Views)
	}
	comments, err := s.ListComments(ctx, lead.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Fatalf("reseed duplicated comments: %d", len(comments))
	}
}
