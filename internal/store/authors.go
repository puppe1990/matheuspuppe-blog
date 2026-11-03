package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func (s *SQLiteStore) ListAuthors(ctx context.Context) ([]models.Author, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT a.id, a.slug, a.name, a.role, a.bio, a.avatar_url, a.twitter, a.github, COUNT(ar.id)
FROM authors a
LEFT JOIN articles ar ON ar.author_id = a.id
GROUP BY a.id
ORDER BY a.id`)
	if err != nil {
		return nil, fmt.Errorf("list authors: %w", err)
	}
	defer func() { _ = rows.Close() }()

	authors := []models.Author{}
	for rows.Next() {
		var author models.Author
		if err := rows.Scan(&author.ID, &author.Slug, &author.Name, &author.Role, &author.Bio,
			&author.AvatarURL, &author.Twitter, &author.GitHub, &author.ArticleCount); err != nil {
			return nil, fmt.Errorf("scan author: %w", err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate authors: %w", err)
	}
	return authors, nil
}

func (s *SQLiteStore) FindAuthorBySlug(ctx context.Context, slug string) (models.Author, error) {
	var author models.Author
	err := s.db.QueryRowContext(ctx, `
SELECT a.id, a.slug, a.name, a.role, a.bio, a.avatar_url, a.twitter, a.github,
       (SELECT COUNT(*) FROM articles ar WHERE ar.author_id = a.id)
FROM authors a
WHERE a.slug = ?`, slug).Scan(&author.ID, &author.Slug, &author.Name, &author.Role,
		&author.Bio, &author.AvatarURL, &author.Twitter, &author.GitHub, &author.ArticleCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Author{}, fmt.Errorf("find author %q: %w", slug, ErrNotFound)
		}
		return models.Author{}, fmt.Errorf("find author %q: %w", slug, err)
	}
	return author, nil
}

func (s *SQLiteStore) UpsertAuthor(ctx context.Context, author models.Author) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
INSERT INTO authors (slug, name, role, bio, avatar_url, twitter, github)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(slug) DO UPDATE SET
  name = excluded.name,
  role = excluded.role,
  bio = excluded.bio,
  avatar_url = excluded.avatar_url,
  twitter = excluded.twitter,
  github = excluded.github
RETURNING id`,
		author.Slug, author.Name, author.Role, author.Bio, author.AvatarURL, author.Twitter, author.GitHub,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert author %q: %w", author.Slug, err)
	}
	return id, nil
}
