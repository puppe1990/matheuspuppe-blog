package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppe1990/amarra-cais/pkg/cais/devlog"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"
	caissqlite "github.com/puppe1990/amarra-cais/pkg/cais/sqlite"
	"github.com/puppe1990/amarra-cais/pkg/cais/sqllog"
	_ "modernc.org/sqlite"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// ErrNotFound marks a lookup miss so handlers can map it to a 404 without
// string matching driver errors.
var ErrNotFound = errors.New("store: not found")

type Store interface {
	Sessions() session.Store
	Ping() error
	DB() *sql.DB
	Close() error

	ListArticles(ctx context.Context) ([]models.Article, error)
	ListArticlesByAuthor(ctx context.Context, authorID int64) ([]models.Article, error)
	FindArticleBySlug(ctx context.Context, slug string) (models.Article, error)
	UpsertArticle(ctx context.Context, article models.Article) (int64, error)
	IncrementArticleViews(ctx context.Context, id int64, delta int) error
	AddArticleClaps(ctx context.Context, id int64, delta int) (int, error)
	PruneCatalog(ctx context.Context, articleSlugs, authorSlugs []string) error

	ListAuthors(ctx context.Context) ([]models.Author, error)
	FindAuthorBySlug(ctx context.Context, slug string) (models.Author, error)
	UpsertAuthor(ctx context.Context, author models.Author) (int64, error)

	ListComments(ctx context.Context, articleID int64) ([]models.Comment, error)
	AddComment(ctx context.Context, comment models.Comment) (int64, error)

	UpsertSubscriber(ctx context.Context, email string) error
	CountSubscribers(ctx context.Context) (int, error)
}

type SQLiteStore struct {
	db *sqllog.DB
}

func NewSQLiteStore(dsn string, env string) (*SQLiteStore, error) {
	if dsn != ":memory:" {
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}

	db, err := sql.Open("sqlite", caissqlite.DSN(dsn))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := applyMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := session.EnsureSQLiteSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("sessions schema: %w", err)
	}

	cfg := sqllog.ConfigForEnv(env)
	if cfg.Enabled {
		cfg.Writer = devlog.MirrorDefault(os.Stdout)
	}
	return &SQLiteStore{db: sqllog.Wrap(db, cfg)}, nil
}

func (s *SQLiteStore) Sessions() session.Store {
	return session.NewSQLiteStore(s.db.Raw())
}

func (s *SQLiteStore) Ping() error {
	return s.db.Raw().Ping()
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db.Raw()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
