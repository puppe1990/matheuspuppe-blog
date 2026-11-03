package store

import (
	"context"
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/testdata"
)

func TestComments_AddAndListNewestFirst(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	article := seedOneArticle(t, s)
	f := testdata.New()

	older := testdata.Comment(f, article.ID)
	older.CreatedAt = time.Date(2026, 3, 28, 14, 20, 0, 0, time.UTC)
	newer := testdata.Comment(f, article.ID)
	newer.CreatedAt = time.Date(2026, 3, 29, 10, 5, 0, 0, time.UTC)

	if _, err := s.AddComment(ctx, older); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddComment(ctx, newer); err != nil {
		t.Fatal(err)
	}

	comments, err := s.ListComments(ctx, article.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Fatalf("got %d comments, want 2", len(comments))
	}
	if !comments[0].CreatedAt.Equal(newer.CreatedAt) {
		t.Fatalf("first comment = %v, want newest %v", comments[0].CreatedAt, newer.CreatedAt)
	}
	if comments[0].AuthorName == "" || comments[0].Content == "" {
		t.Fatalf("comment fields empty: %+v", comments[0])
	}
}

func TestComments_ScopedToArticle(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	article := seedOneArticle(t, s)
	f := testdata.New()

	comment := testdata.Comment(f, article.ID)
	if _, err := s.AddComment(ctx, comment); err != nil {
		t.Fatal(err)
	}

	other, err := s.ListComments(ctx, article.ID+99)
	if err != nil {
		t.Fatal(err)
	}
	if len(other) != 0 {
		t.Fatalf("got %d comments for another article, want 0", len(other))
	}
}

func TestSubscribers_UpsertIsIdempotent(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if err := s.UpsertSubscriber(ctx, "leitor@sintaxe.tech"); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertSubscriber(ctx, "leitor@sintaxe.tech"); err != nil {
		t.Fatal(err)
	}

	count, err := s.CountSubscribers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}
