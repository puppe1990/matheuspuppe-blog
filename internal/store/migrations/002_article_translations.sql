-- 002_article_translations.sql — English versions of the editorial catalog.
-- up
ALTER TABLE authors ADD COLUMN role_en TEXT NOT NULL DEFAULT '';
ALTER TABLE authors ADD COLUMN bio_en TEXT NOT NULL DEFAULT '';

ALTER TABLE articles ADD COLUMN slug_en TEXT NOT NULL DEFAULT '';
ALTER TABLE articles ADD COLUMN title_en TEXT NOT NULL DEFAULT '';
ALTER TABLE articles ADD COLUMN subtitle_en TEXT NOT NULL DEFAULT '';
ALTER TABLE articles ADD COLUMN lead_paragraph_en TEXT NOT NULL DEFAULT '';
ALTER TABLE articles ADD COLUMN tags_en TEXT NOT NULL DEFAULT '[]';
ALTER TABLE articles ADD COLUMN blocks_en TEXT NOT NULL DEFAULT '[]';

CREATE UNIQUE INDEX idx_articles_slug_en ON articles(slug_en) WHERE slug_en <> '';

-- down
DROP INDEX IF EXISTS idx_articles_slug_en;
ALTER TABLE articles DROP COLUMN blocks_en;
ALTER TABLE articles DROP COLUMN tags_en;
ALTER TABLE articles DROP COLUMN lead_paragraph_en;
ALTER TABLE articles DROP COLUMN subtitle_en;
ALTER TABLE articles DROP COLUMN title_en;
ALTER TABLE articles DROP COLUMN slug_en;
ALTER TABLE authors DROP COLUMN bio_en;
ALTER TABLE authors DROP COLUMN role_en;
