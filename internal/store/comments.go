package store

import (
	"context"
	"fmt"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func (s *SQLiteStore) ListComments(ctx context.Context, articleID int64) ([]models.Comment, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, article_id, author_name, author_role, content, created_at, claps
FROM comments
WHERE article_id = ?
ORDER BY created_at DESC, id DESC`, articleID)
	if err != nil {
		return nil, fmt.Errorf("list comments of article %d: %w", articleID, err)
	}
	defer func() { _ = rows.Close() }()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		var createdAt string
		if err := rows.Scan(&comment.ID, &comment.ArticleID, &comment.AuthorName, &comment.AuthorRole,
			&comment.Content, &createdAt, &comment.Claps); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		parsed, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse comment created_at %q: %w", createdAt, err)
		}
		comment.CreatedAt = parsed
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate comments: %w", err)
	}
	return comments, nil
}

func (s *SQLiteStore) AddComment(ctx context.Context, comment models.Comment) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `
INSERT INTO comments (article_id, author_name, author_role, content, created_at, claps)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id`,
		comment.ArticleID, comment.AuthorName, comment.AuthorRole, comment.Content,
		comment.CreatedAt.UTC().Format(time.RFC3339), comment.Claps,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("add comment to article %d: %w", comment.ArticleID, err)
	}
	return id, nil
}
