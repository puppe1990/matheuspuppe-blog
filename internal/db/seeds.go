package db

import (
	"context"
	"fmt"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// RunSeeds upserts the Sintaxe editorial catalog: authors, essays and starter
// comments. Idempotent by design — content refreshes by slug while live reader
// counters (claps, views) and reader-written comments survive a reseed. The
// seed set is the source of truth: catalog rows that left the code are pruned
// (with their comments) in the same run.
func RunSeeds(s store.Store) error {
	ctx := context.Background()

	// cais:seeds
	authorIDs, err := upsertSeedAuthors(ctx, s)
	if err != nil {
		return err
	}

	articleIDs, err := upsertSeedArticles(ctx, s, authorIDs)
	if err != nil {
		return err
	}

	if err := s.PruneCatalog(ctx, slugsOf(articleIDs), slugsOf(authorIDs)); err != nil {
		return fmt.Errorf("prune catalog: %w", err)
	}

	return seedStarterComments(ctx, s, articleIDs)
}

func slugsOf(ids map[string]int64) []string {
	slugs := make([]string, 0, len(ids))
	for slug := range ids {
		slugs = append(slugs, slug)
	}
	return slugs
}

func upsertSeedAuthors(ctx context.Context, s store.Store) (map[string]int64, error) {
	authorIDs := map[string]int64{}
	for _, author := range seedAuthors() {
		id, err := s.UpsertAuthor(ctx, author)
		if err != nil {
			return nil, fmt.Errorf("seed author %q: %w", author.Slug, err)
		}
		authorIDs[author.Slug] = id
	}
	return authorIDs, nil
}

func upsertSeedArticles(ctx context.Context, s store.Store, authorIDs map[string]int64) (map[string]int64, error) {
	articleIDs := map[string]int64{}
	for _, article := range seedArticles() {
		authorID, ok := authorIDs[article.Author.Slug]
		if !ok {
			return nil, fmt.Errorf("seed article %q references unknown author %q", article.Slug, article.Author.Slug)
		}
		article.Author.ID = authorID
		id, err := s.UpsertArticle(ctx, article)
		if err != nil {
			return nil, fmt.Errorf("seed article %q: %w", article.Slug, err)
		}
		articleIDs[article.Slug] = id
	}
	return articleIDs, nil
}

// seedStarterComments only fills articles that have no comments yet, so a
// reseed never duplicates the conversation. The emptiness check runs once per
// article — checking per comment would skip the second comment of a pair.
func seedStarterComments(ctx context.Context, s store.Store, articleIDs map[string]int64) error {
	byArticle := map[string][]models.Comment{}
	order := []string{}
	for _, seed := range seedComments() {
		if _, ok := articleIDs[seed.ArticleSlug]; !ok {
			return fmt.Errorf("seed comment references unknown article %q", seed.ArticleSlug)
		}
		if _, seen := byArticle[seed.ArticleSlug]; !seen {
			order = append(order, seed.ArticleSlug)
		}
		byArticle[seed.ArticleSlug] = append(byArticle[seed.ArticleSlug], seed.Comment)
	}

	for _, slug := range order {
		articleID := articleIDs[slug]
		comments, err := s.ListComments(ctx, articleID)
		if err != nil {
			return fmt.Errorf("seed comments of %q: %w", slug, err)
		}
		if len(comments) > 0 {
			continue
		}
		for _, comment := range byArticle[slug] {
			comment.ArticleID = articleID
			if _, err := s.AddComment(ctx, comment); err != nil {
				return fmt.Errorf("seed comment of %q: %w", slug, err)
			}
		}
	}
	return nil
}
