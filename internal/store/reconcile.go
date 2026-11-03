package store

import (
	"context"
	"fmt"
	"strings"
)

// PruneCatalog reconciles the database with the code-owned seeds: comments
// hanging on removed essays, the removed essays and any author outside the
// seed set are deleted in one transaction, in FK order. Comments on surviving
// essays (reader-written) are untouched. Keeping the seed set as the single
// source of truth means removing an essay is a code edit, not a DB ritual.
func (s *SQLiteStore) PruneCatalog(ctx context.Context, articleSlugs, authorSlugs []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("prune catalog: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	articleCond, articleArgs := notInClause("slug", articleSlugs)
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM comments WHERE article_id IN (SELECT id FROM articles WHERE "+articleCond+")",
		articleArgs...); err != nil {
		return fmt.Errorf("prune comments: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM articles WHERE "+articleCond, articleArgs...); err != nil {
		return fmt.Errorf("prune articles: %w", err)
	}

	authorCond, authorArgs := notInClause("slug", authorSlugs)
	if _, err := tx.ExecContext(ctx, "DELETE FROM authors WHERE "+authorCond, authorArgs...); err != nil {
		return fmt.Errorf("prune authors: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("prune catalog: commit: %w", err)
	}
	return nil
}

// notInClause builds "column NOT IN (?, …)" or an always-true condition when
// the keep set is empty (SQLite rejects an empty IN list).
func notInClause(column string, values []string) (string, []any) {
	if len(values) == 0 {
		return "1 = 1", nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(values)), ",")
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return column + " NOT IN (" + placeholders + ")", args
}
