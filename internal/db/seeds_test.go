package db

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
	"github.com/puppe1990/matheuspuppe-blog/internal/testdata"
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

func TestRunSeeds_PublishesTheFirstEdition(t *testing.T) {
	s := newTestStore(t)
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	authors, err := s.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 {
		t.Fatalf("authors = %d, want only the columnist", len(authors))
	}
	author := authors[0]
	if author.Slug != "matheus-puppe" || author.Name != "Matheus Puppe" {
		t.Fatalf("author = %+v", author)
	}
	if author.AvatarURL != "/static/img/author-matheus.jpg" {
		t.Fatalf("avatar = %q", author.AvatarURL)
	}
	if author.ArticleCount != 1 {
		t.Fatalf("article count = %d, want 1", author.ArticleCount)
	}

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 1 {
		t.Fatalf("articles = %d, want the first edition", len(articles))
	}
	if articles[0].Slug != "traducao-epubs-llms-fila-por-capitulo" {
		t.Fatalf("article = %q, want the translation essay", articles[0].Slug)
	}
	if articles[0].Author.Slug != "matheus-puppe" {
		t.Fatalf("article author = %q", articles[0].Author.Slug)
	}
	if len(articles[0].Blocks) == 0 {
		t.Fatalf("first essay has no content blocks")
	}
	if articles[0].SlugEN != "epub-translation-llms-resilient-queue" {
		t.Fatalf("english slug = %q", articles[0].SlugEN)
	}
	if articles[0].TitleEN == "" || articles[0].LeadParagraphEN == "" || len(articles[0].BlocksEN) == 0 {
		t.Fatalf("english edition incomplete: %+v", articles[0])
	}
	if author.RoleEN == "" || author.BioEN == "" {
		t.Fatalf("author english fields empty: %+v", author)
	}
}

func TestPrimeiroEnsaio_CitesThePublicRepo(t *testing.T) {
	const repo = "https://github.com/puppe1990/epub-translator"

	countCitations := func(blocks []models.ContentBlock) (codeLinks, repoLinks int) {
		for _, block := range blocks {
			if strings.HasPrefix(block.SourceURL, repo+"/blob/") {
				codeLinks++
			}
			if block.SourceURL == repo {
				repoLinks++
			}
		}
		return codeLinks, repoLinks
	}

	codeLinks, repoLinks := countCitations(primeiroEnsaio().Blocks)
	if codeLinks != 4 {
		t.Fatalf("citations to the repo history = %d, want 2 codes + 2 comments", codeLinks)
	}
	if repoLinks != 1 {
		t.Fatalf("repo links = %d, want 1", repoLinks)
	}

	englishCode, englishRepo := countCitations(primeiroEnsaioEn().Blocks)
	if englishCode != codeLinks || englishRepo != repoLinks {
		t.Fatalf("english citations = %d codes + %d repo, want %d + %d", englishCode, englishRepo, codeLinks, repoLinks)
	}
}

func TestRunSeeds_PrunesCatalogRowsThatLeftTheSeeds(t *testing.T) {
	s := newTestStore(t)
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// An orphan catalog row written outside the seeds: author, essay and a
	// reader comment hanging on the essay.
	f := testdata.New()
	orphanAuthor := testdata.Author(f, 7)
	authorID, err := s.UpsertAuthor(ctx, orphanAuthor)
	if err != nil {
		t.Fatal(err)
	}
	orphanAuthor.ID = authorID
	orphanEssay := testdata.Article(f, 7, orphanAuthor, time.Now().UTC().AddDate(0, 0, -1))
	essayID, err := s.UpsertArticle(ctx, orphanEssay)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddComment(ctx, testdata.Comment(f, essayID)); err != nil {
		t.Fatal(err)
	}

	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 1 || articles[0].Slug != "traducao-epubs-llms-fila-por-capitulo" {
		t.Fatalf("catalog after reseed = %+v, want only the seeded essay", articles)
	}
	authors, err := s.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 || authors[0].Slug != "matheus-puppe" {
		t.Fatalf("orphan author survived the reseed: %+v", authors)
	}
	var comments int
	if err := s.DB().QueryRow(`SELECT COUNT(*) FROM comments`).Scan(&comments); err != nil {
		t.Fatal(err)
	}
	if comments != 0 {
		t.Fatalf("comments of the removed essay survived: %d", comments)
	}
	if _, err := s.FindArticleBySlug(ctx, orphanEssay.Slug); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("removed essay still resolves: %v", err)
	}
}

func TestRunSeeds_Idempotent(t *testing.T) {
	s := newTestStore(t)
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	if err := RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	authors, err := s.ListAuthors(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(authors) != 1 {
		t.Fatalf("reseed duplicated the author: %d", len(authors))
	}
}
