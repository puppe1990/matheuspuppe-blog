package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/testdata"
)

// seedArticleAt upserts a fake essay with one exact publish instant, which is
// what the scheduling tests vary.
func seedArticleAt(t *testing.T, s *SQLiteStore, author models.Author, n int, at time.Time) models.Article {
	t.Helper()
	article := testdata.Article(testdata.New(), n, author, at)
	article.PublishedAt = at
	if _, err := s.UpsertArticle(context.Background(), article); err != nil {
		t.Fatal(err)
	}
	return article
}

func seedAuthor(t *testing.T, s *SQLiteStore) models.Author {
	t.Helper()
	author := testdata.Author(testdata.New(), 1)
	id, err := s.UpsertAuthor(context.Background(), author)
	if err != nil {
		t.Fatal(err)
	}
	author.ID = id
	return author
}

func TestListArticles_HidesScheduledEssaysUntilTheirDay(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	author := seedAuthor(t, s)
	now := time.Now().UTC()

	past := seedArticleAt(t, s, author, 1, now.AddDate(0, 0, -1))
	today := seedArticleAt(t, s, author, 2, now) // equal instant counts as published
	future := seedArticleAt(t, s, author, 3, now.AddDate(0, 0, 1))

	articles, err := s.ListArticles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, article := range articles {
		got[article.Slug] = true
	}
	if !got[past.Slug] || !got[today.Slug] {
		t.Fatalf("published essays hidden: %v", articles)
	}
	if got[future.Slug] {
		t.Fatalf("scheduled essay leaked into the feed before its day")
	}
	if len(articles) != 2 {
		t.Fatalf("got %d articles, want 2", len(articles))
	}
}

func TestFindArticleBySlug_ScheduledIsNotFoundUntilItsDay(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	author := seedAuthor(t, s)
	future := seedArticleAt(t, s, author, 1, time.Now().UTC().AddDate(0, 0, 2))

	if _, err := s.FindArticleBySlug(ctx, future.Slug); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound for a scheduled essay", err)
	}

	// The row is already stored; only visibility is deferred, so the day
	// arrives without another deploy.
	var count int
	if err := s.DB().QueryRow(`SELECT COUNT(*) FROM articles WHERE slug = ?`, future.Slug).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("scheduled essay not persisted: count = %d", count)
	}
}

func TestListArticlesByAuthor_HidesScheduledEssays(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	author := seedAuthor(t, s)
	now := time.Now().UTC()

	published := seedArticleAt(t, s, author, 1, now.AddDate(0, 0, -2))
	scheduled := seedArticleAt(t, s, author, 2, now.AddDate(0, 0, 2))

	articles, err := s.ListArticlesByAuthor(ctx, author.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(articles) != 1 || articles[0].Slug != published.Slug {
		t.Fatalf("got %d articles (want only %s) — scheduled: %s", len(articles), published.Slug, scheduled.Slug)
	}
}
