// Package models holds the editorial row structs (SQL columns) and the page
// payloads handlers render.
package models

import "time"

// Author is one Sintaxe columnist. ArticleCount is computed on listing pages.
type Author struct {
	ID           int64
	Slug         string
	Name         string
	Role         string
	Bio          string
	AvatarURL    string
	Twitter      string
	GitHub       string
	ArticleCount int
}

// BenchmarkMetric is one row of an interactive-benchmark block.
type BenchmarkMetric struct {
	Label  string  `json:"label"`
	ValueA float64 `json:"value_a"`
	ValueB float64 `json:"value_b"`
	Unit   string  `json:"unit"`
}

// BenchmarkData backs the interactive-benchmark content block.
type BenchmarkData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Metrics     []BenchmarkMetric `json:"metrics"`
}

// ContentBlock is one typed paragraph of an article body. Blocks are stored as
// JSON in articles.blocks; the renderer switches on Type.
type ContentBlock struct {
	Type        string         `json:"type"`
	Text        string         `json:"text,omitempty"`
	Code        string         `json:"code,omitempty"`
	Language    string         `json:"language,omitempty"`
	Caption     string         `json:"caption,omitempty"`
	AuthorQuote string         `json:"author_quote,omitempty"`
	Items       []string       `json:"items,omitempty"`
	Benchmark   *BenchmarkData `json:"benchmark,omitempty"`
}

// Article is one published essay. Tags and Blocks are JSON columns.
type Article struct {
	ID              int64
	Slug            string
	Title           string
	Subtitle        string
	LeadParagraph   string
	Category        string
	ReadTimeMinutes int
	PublishedAt     time.Time
	CoverImage      string
	Author          Author
	Tags            []string
	Blocks          []ContentBlock
	Claps           int
	Views           int
	Tier            int
}

// Comment is one reader reflection under an article.
type Comment struct {
	ID         int64
	ArticleID  int64
	AuthorName string
	AuthorRole string
	Content    string
	CreatedAt  time.Time
	Claps      int
}

// Subscriber is one newsletter subscription.
type Subscriber struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}
