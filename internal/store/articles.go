package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// articleSelect joins the author so listing pages render bylines without a
// second query per card. The *_en columns carry the English edition.
const articleSelect = `
SELECT ar.id, ar.slug, ar.title, ar.subtitle, ar.lead_paragraph, ar.category,
       ar.read_time_minutes, ar.published_at, ar.cover_image, ar.tags, ar.blocks,
       ar.claps, ar.views, ar.tier,
       ar.slug_en, ar.title_en, ar.subtitle_en, ar.lead_paragraph_en, ar.tags_en, ar.blocks_en,
       a.id, a.slug, a.name, a.role, a.bio, a.role_en, a.bio_en, a.avatar_url, a.twitter, a.github
FROM articles ar
JOIN authors a ON a.id = ar.author_id`

// publishedNow is the scheduling cutoff. Seeds ship scheduled essays ahead of
// time (published_at in the future) and every read hides them until the
// instant arrives, so a post appears by clock on its day — no redeploy. Stored
// timestamps are RFC3339 UTC, so lexicographic <= is a correct instant compare.
func publishedNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (s *SQLiteStore) ListArticles(ctx context.Context) ([]models.Article, error) {
	rows, err := s.db.QueryContext(ctx, articleSelect+`
WHERE ar.published_at <= ?
ORDER BY ar.published_at DESC, ar.id DESC`, publishedNow())
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanArticles(rows)
}

func (s *SQLiteStore) ListArticlesByAuthor(ctx context.Context, authorID int64) ([]models.Article, error) {
	rows, err := s.db.QueryContext(ctx, articleSelect+`
WHERE ar.author_id = ? AND ar.published_at <= ?
ORDER BY ar.published_at DESC, ar.id DESC`, authorID, publishedNow())
	if err != nil {
		return nil, fmt.Errorf("list articles by author %d: %w", authorID, err)
	}
	defer func() { _ = rows.Close() }()
	return scanArticles(rows)
}

func (s *SQLiteStore) FindArticleBySlug(ctx context.Context, slug string) (models.Article, error) {
	return s.findArticle(ctx, "ar.slug", slug)
}

// FindArticleBySlugEN resolves the English edition by its own slug.
func (s *SQLiteStore) FindArticleBySlugEN(ctx context.Context, slug string) (models.Article, error) {
	return s.findArticle(ctx, "ar.slug_en", slug)
}

func (s *SQLiteStore) findArticle(ctx context.Context, column, slug string) (models.Article, error) {
	rows, err := s.db.QueryContext(ctx, articleSelect+`
WHERE `+column+` = ? AND ar.published_at <= ?`, slug, publishedNow())
	if err != nil {
		return models.Article{}, fmt.Errorf("find article %q: %w", slug, err)
	}
	defer func() { _ = rows.Close() }()

	articles, err := scanArticles(rows)
	if err != nil {
		return models.Article{}, fmt.Errorf("find article %q: %w", slug, err)
	}
	if len(articles) == 0 {
		return models.Article{}, fmt.Errorf("find article %q: %w", slug, ErrNotFound)
	}
	return articles[0], nil
}

// UpsertArticle is the seed write: content fields refresh on slug conflict,
// while claps and views survive because they are live reader counters.
func (s *SQLiteStore) UpsertArticle(ctx context.Context, article models.Article) (int64, error) {
	tags, err := json.Marshal(article.Tags)
	if err != nil {
		return 0, fmt.Errorf("marshal tags for %q: %w", article.Slug, err)
	}
	blocks, err := json.Marshal(article.Blocks)
	if err != nil {
		return 0, fmt.Errorf("marshal blocks for %q: %w", article.Slug, err)
	}
	tagsEN, err := json.Marshal(article.TagsEN)
	if err != nil {
		return 0, fmt.Errorf("marshal tags_en for %q: %w", article.Slug, err)
	}
	blocksEN, err := json.Marshal(article.BlocksEN)
	if err != nil {
		return 0, fmt.Errorf("marshal blocks_en for %q: %w", article.Slug, err)
	}

	var id int64
	err = s.db.QueryRowContext(ctx, `
INSERT INTO articles (slug, title, subtitle, lead_paragraph, category, read_time_minutes,
                      published_at, cover_image, author_id, tags, blocks, claps, views, tier,
                      slug_en, title_en, subtitle_en, lead_paragraph_en, tags_en, blocks_en)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(slug) DO UPDATE SET
  title = excluded.title,
  subtitle = excluded.subtitle,
  lead_paragraph = excluded.lead_paragraph,
  category = excluded.category,
  read_time_minutes = excluded.read_time_minutes,
  published_at = excluded.published_at,
  cover_image = excluded.cover_image,
  author_id = excluded.author_id,
  tags = excluded.tags,
  blocks = excluded.blocks,
  tier = excluded.tier,
  slug_en = excluded.slug_en,
  title_en = excluded.title_en,
  subtitle_en = excluded.subtitle_en,
  lead_paragraph_en = excluded.lead_paragraph_en,
  tags_en = excluded.tags_en,
  blocks_en = excluded.blocks_en
RETURNING id`,
		article.Slug, article.Title, article.Subtitle, article.LeadParagraph, article.Category,
		article.ReadTimeMinutes, article.PublishedAt.UTC().Format(time.RFC3339), article.CoverImage,
		article.Author.ID, string(tags), string(blocks), article.Claps, article.Views, article.Tier,
		article.SlugEN, article.TitleEN, article.SubtitleEN, article.LeadParagraphEN,
		string(tagsEN), string(blocksEN),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert article %q: %w", article.Slug, err)
	}
	return id, nil
}

func (s *SQLiteStore) IncrementArticleViews(ctx context.Context, id int64, delta int) error {
	if delta < 1 {
		delta = 1
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE articles SET views = views + ? WHERE id = ?`, delta, id); err != nil {
		return fmt.Errorf("increment views of article %d: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) AddArticleClaps(ctx context.Context, id int64, delta int) (int, error) {
	if delta < 1 {
		delta = 1
	}
	var total int
	err := s.db.QueryRowContext(ctx, `UPDATE articles SET claps = claps + ? WHERE id = ? RETURNING claps`, delta, id).Scan(&total)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("clap article %d: %w", id, ErrNotFound)
		}
		return 0, fmt.Errorf("clap article %d: %w", id, err)
	}
	return total, nil
}

func scanArticles(rows *sql.Rows) ([]models.Article, error) {
	articles := []models.Article{}
	for rows.Next() {
		var article models.Article
		var publishedAt, tags, blocks, tagsEN, blocksEN string
		err := rows.Scan(
			&article.ID, &article.Slug, &article.Title, &article.Subtitle, &article.LeadParagraph,
			&article.Category, &article.ReadTimeMinutes, &publishedAt, &article.CoverImage,
			&tags, &blocks, &article.Claps, &article.Views, &article.Tier,
			&article.SlugEN, &article.TitleEN, &article.SubtitleEN, &article.LeadParagraphEN,
			&tagsEN, &blocksEN,
			&article.Author.ID, &article.Author.Slug, &article.Author.Name, &article.Author.Role,
			&article.Author.Bio, &article.Author.RoleEN, &article.Author.BioEN,
			&article.Author.AvatarURL, &article.Author.Twitter, &article.Author.GitHub,
		)
		if err != nil {
			return nil, fmt.Errorf("scan article: %w", err)
		}
		if err := hydrateArticle(&article, publishedAt, tags, blocks, tagsEN, blocksEN); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate articles: %w", err)
	}
	return articles, nil
}

func hydrateArticle(article *models.Article, publishedAt, tags, blocks, tagsEN, blocksEN string) error {
	parsed, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		return fmt.Errorf("parse published_at %q of %q: %w", publishedAt, article.Slug, err)
	}
	article.PublishedAt = parsed
	if err := json.Unmarshal([]byte(tags), &article.Tags); err != nil {
		return fmt.Errorf("unmarshal tags of %q: %w", article.Slug, err)
	}
	if err := json.Unmarshal([]byte(blocks), &article.Blocks); err != nil {
		return fmt.Errorf("unmarshal blocks of %q: %w", article.Slug, err)
	}
	if tagsEN != "" {
		if err := json.Unmarshal([]byte(tagsEN), &article.TagsEN); err != nil {
			return fmt.Errorf("unmarshal tags_en of %q: %w", article.Slug, err)
		}
	}
	if blocksEN != "" {
		if err := json.Unmarshal([]byte(blocksEN), &article.BlocksEN); err != nil {
			return fmt.Errorf("unmarshal blocks_en of %q: %w", article.Slug, err)
		}
	}
	return nil
}
