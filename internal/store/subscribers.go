package store

import (
	"context"
	"fmt"
	"time"
)

// UpsertSubscriber keeps newsletter signups idempotent: the same address never
// double-registers and keeps its original created_at.
func (s *SQLiteStore) UpsertSubscriber(ctx context.Context, email string) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO subscribers (email, created_at)
VALUES (?, ?)
ON CONFLICT(email) DO UPDATE SET email = excluded.email`,
		email, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("upsert subscriber %q: %w", email, err)
	}
	return nil
}

func (s *SQLiteStore) CountSubscribers(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscribers`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count subscribers: %w", err)
	}
	return count, nil
}
