package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// FilterLink is one pressable filter chip (category or reading time).
type FilterLink struct {
	Label  string
	Href   string
	Count  int
	Active bool
}

type HomeHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewHomeHandler(views *view.Renderer, st store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *HomeHandler {
	return &HomeHandler{views: views, store: st, site: site, catalog: catalog, cfg: cfg}
}

func (h *HomeHandler) t(r *http.Request, key string) string {
	return i18n.CatalogOr(r, h.catalog).T(key)
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	articles, err := h.store.ListArticles(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	category := r.URL.Query().Get("categoria")
	if category == "" {
		category = "Todos"
	}
	maxMinutes := intParam(r, "tempo", 0)

	saved := bookmarksFromRequest(r)
	filtering := editorial.Filtering(search, category, maxMinutes)
	filtered := editorial.Filter(articles, editorial.FilterOptions{
		Search:     search,
		Category:   category,
		MaxMinutes: maxMinutes,
		OnlySlugs:  nil,
	})
	lead, secondary, grid := editorial.Segment(filtered, filtering)

	writeView(w, r, h.views, h.cfg, "app", "home", amarraData(r, h.site, map[string]any{
		"Title":          h.t(r, "home.title"),
		"ActiveNav":      "artigos",
		"EditionDate":    editorial.FormatDate(time.Now()),
		"TotalArticles":  len(articles),
		"TotalResults":   len(filtered),
		"Filtering":      filtering,
		"Search":         search,
		"ActiveCategory": category,
		"MaxMinutes":     maxMinutes,
		"CategoryLinks":  categoryLinks(articles, search, category, maxMinutes),
		"TimeLinks":      timeLinks(search, category, maxMinutes),
		"ClearHref":      editorial.CategoryHref(category, "", maxMinutes),
		"Lead":           leadArticleView(lead, saved),
		"Secondary":      newArticleViews(secondary, saved),
		"Grid":           newArticleViews(grid, saved),
	}), http.StatusOK)
}

// categoryLinks builds the segmented caderno bar with live counts; the active
// search and reading-time filters survive the category switch.
func categoryLinks(articles []models.Article, search, activeCategory string, maxMinutes int) []FilterLink {
	counts := editorial.CategoryCounts(articles)
	links := make([]FilterLink, 0, len(editorial.Categories))
	for _, category := range editorial.Categories {
		links = append(links, FilterLink{
			Label:  category,
			Href:   editorial.CategoryHref(category, search, maxMinutes),
			Count:  counts[category],
			Active: category == activeCategory,
		})
	}
	return links
}

func timeLinks(search, activeCategory string, maxMinutes int) []FilterLink {
	options := []struct {
		label   string
		minutes int
	}{
		{"Todos", 0},
		{"≤ 5 min", 5},
		{"≤ 8 min", 8},
	}
	links := make([]FilterLink, 0, len(options))
	for _, option := range options {
		links = append(links, FilterLink{
			Label:  option.label,
			Href:   editorial.CategoryHref(activeCategory, search, option.minutes),
			Active: maxMinutes == option.minutes,
		})
	}
	return links
}
