-- 001_editorial.sql — Sintaxe editorial schema.
-- up
CREATE TABLE authors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  slug TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT '',
  bio TEXT NOT NULL DEFAULT '',
  avatar_url TEXT NOT NULL DEFAULT '',
  twitter TEXT NOT NULL DEFAULT '',
  github TEXT NOT NULL DEFAULT ''
);

CREATE TABLE articles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  slug TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  subtitle TEXT NOT NULL DEFAULT '',
  lead_paragraph TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL,
  read_time_minutes INTEGER NOT NULL DEFAULT 1,
  published_at TEXT NOT NULL,
  cover_image TEXT NOT NULL DEFAULT '',
  author_id INTEGER NOT NULL REFERENCES authors(id),
  tags TEXT NOT NULL DEFAULT '[]',
  blocks TEXT NOT NULL DEFAULT '[]',
  claps INTEGER NOT NULL DEFAULT 0,
  views INTEGER NOT NULL DEFAULT 0,
  tier INTEGER NOT NULL DEFAULT 3
);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_published ON articles(published_at DESC);

CREATE TABLE comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  article_id INTEGER NOT NULL REFERENCES articles(id),
  author_name TEXT NOT NULL,
  author_role TEXT NOT NULL DEFAULT '',
  content TEXT NOT NULL,
  created_at TEXT NOT NULL,
  claps INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_comments_article ON comments(article_id, created_at DESC);

CREATE TABLE subscribers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL
);

-- down
DROP TABLE IF EXISTS subscribers;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS authors;
