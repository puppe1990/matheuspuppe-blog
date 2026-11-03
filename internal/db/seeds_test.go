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
	if author.ArticleCount != 2 {
		t.Fatalf("article count = %d, want 2", author.ArticleCount)
	}

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 2 {
		t.Fatalf("articles = %d, want the published edition", len(articles))
	}
	if articles[0].Slug != "captura-de-trafego-do-browser-para-agentes" {
		t.Fatalf("newest article = %q, want the network listener essay", articles[0].Slug)
	}

	translation := articles[1]
	if translation.Slug != "traducao-epubs-llms-fila-por-capitulo" {
		t.Fatalf("article = %q, want the translation essay", translation.Slug)
	}
	if translation.Author.Slug != "matheus-puppe" {
		t.Fatalf("article author = %q", translation.Author.Slug)
	}
	if len(translation.Blocks) == 0 {
		t.Fatalf("first essay has no content blocks")
	}
	if translation.SlugEN != "epub-translation-llms-resilient-queue" {
		t.Fatalf("english slug = %q", translation.SlugEN)
	}
	if translation.TitleEN == "" || translation.LeadParagraphEN == "" || len(translation.BlocksEN) == 0 {
		t.Fatalf("english edition incomplete: %+v", translation)
	}

	listener := articles[0]
	if listener.SlugEN != "browser-traffic-capture-for-agents" {
		t.Fatalf("network listener english slug = %q", listener.SlugEN)
	}
	if len(listener.Blocks) == 0 || len(listener.BlocksEN) == 0 {
		t.Fatalf("network listener missing blocks: pt=%d en=%d", len(listener.Blocks), len(listener.BlocksEN))
	}
	if listener.Category != "Segurança & Criptografia" {
		t.Fatalf("network listener category = %q", listener.Category)
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

func TestNetworkListener_CitesThePublicRepo(t *testing.T) {
	const repo = "https://github.com/puppe1990/network-listener"

	citations, repoLinks := 0, 0
	for _, block := range networkListener().Blocks {
		if strings.HasPrefix(block.SourceURL, repo) {
			citations++
			if block.SourceURL == repo {
				repoLinks++
			}
		}
	}
	if citations != 6 {
		t.Fatalf("citations = %d, want 3 codes + 2 quotes + 1 repo", citations)
	}
	if repoLinks != 1 {
		t.Fatalf("repo links = %d, want 1", repoLinks)
	}

	english := networkListenerEn().Blocks
	englishCitations := 0
	for _, block := range english {
		if strings.HasPrefix(block.SourceURL, repo) {
			englishCitations++
		}
	}
	if englishCitations != citations {
		t.Fatalf("english citations = %d, want %d", englishCitations, citations)
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
	if len(articles) != 2 {
		t.Fatalf("catalog after reseed = %d articles, want the two seeded essays", len(articles))
	}
	slugs := map[string]bool{}
	for _, article := range articles {
		slugs[article.Slug] = true
	}
	if !slugs["traducao-epubs-llms-fila-por-capitulo"] || !slugs["captura-de-trafego-do-browser-para-agentes"] {
		t.Fatalf("catalog after reseed = %+v, want both seeded essays", slugs)
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
