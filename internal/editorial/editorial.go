// Package editorial holds the Sintaxe editorial rules: filtering, tier
// segmentation and category counts. Pure functions over models rows so the
// home page and the saved library share exactly one behavior.
package editorial

import (
	"strings"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// Categories lists the front-page cadernos, "Todos" first. Keep in sync with
// the article seeds; CategoryCounts only reports what the articles use.
var Categories = []string{
	"Todos",
	"IA & Modelos",
	"Engenharia de Sistemas",
	"Hardware & Chips",
	"Design & UX",
	"Segurança & Criptografia",
	"Opinião & Futuro",
}

// FilterOptions mirrors the front-page controls. Zero values mean "off".
type FilterOptions struct {
	Search     string
	Category   string
	MaxMinutes int
	// OnlySlugs, when non-nil, restricts the result to the saved library.
	OnlySlugs []string
}

// Filter applies search text, category, reading-time cap and the saved set.
// Search matches title, subtitle, lead, author name and tags (case-folded).
func Filter(articles []models.Article, opts FilterOptions) []models.Article {
	query := strings.ToLower(strings.TrimSpace(opts.Search))
	only := slugSet(opts.OnlySlugs)

	out := make([]models.Article, 0, len(articles))
	for _, article := range articles {
		if only != nil && !only[article.Slug] {
			continue
		}
		if opts.Category != "" && opts.Category != "Todos" && article.Category != opts.Category {
			continue
		}
		if opts.MaxMinutes > 0 && article.ReadTimeMinutes > opts.MaxMinutes {
			continue
		}
		if query != "" && !matchesSearch(article, query) {
			continue
		}
		out = append(out, article)
	}
	return out
}

func slugSet(slugs []string) map[string]bool {
	if slugs == nil {
		return nil
	}
	set := make(map[string]bool, len(slugs))
	for _, slug := range slugs {
		set[slug] = true
	}
	return set
}

func matchesSearch(article models.Article, query string) bool {
	fields := []string{article.Title, article.Subtitle, article.LeadParagraph, article.Author.Name}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	for _, tag := range article.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

// Segment splits the feed into the 3-tier front page: one lead story, up to two
// secondary features and the remaining grid. While a filter is active the
// editorial tiers are ignored and the feed is read top-down.
func Segment(articles []models.Article, filtering bool) (*models.Article, []models.Article, []models.Article) {
	if len(articles) == 0 {
		return nil, nil, nil
	}
	if filtering {
		lead := articles[0]
		secondary := take(articles[1:], 2)
		grid := articles[1+len(secondary):]
		return &lead, secondary, grid
	}

	leadIdx := -1
	for i, article := range articles {
		if article.Tier == 1 {
			leadIdx = i
			break
		}
	}
	if leadIdx == -1 {
		leadIdx = 0
	}
	lead := articles[leadIdx]

	secondary := make([]models.Article, 0, 2)
	grid := make([]models.Article, 0, len(articles))
	for i, article := range articles {
		if i == leadIdx {
			continue
		}
		if article.Tier == 2 && len(secondary) < 2 {
			secondary = append(secondary, article)
			continue
		}
		grid = append(grid, article)
	}
	return &lead, secondary, grid
}

func take(articles []models.Article, n int) []models.Article {
	if len(articles) < n {
		n = len(articles)
	}
	out := make([]models.Article, n)
	copy(out, articles[:n])
	return out
}

// CategoryCounts totals articles per category; "Todos" carries the grand total.
func CategoryCounts(articles []models.Article) map[string]int {
	counts := map[string]int{"Todos": len(articles)}
	for _, article := range articles {
		counts[article.Category]++
	}
	return counts
}

// Filtering reports whether any front-page control is narrowing the feed.
func Filtering(search, category string, maxMinutes int) bool {
	return strings.TrimSpace(search) != "" || (category != "" && category != "Todos") || maxMinutes > 0
}
