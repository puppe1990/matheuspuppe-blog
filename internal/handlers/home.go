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

func (h *HomeHandler) t(r *http.Request, key string, args ...any) string {
	return i18n.CatalogOr(r, h.catalog).T(key, args...)
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	locale := localeOf(r)
	articles, err := h.store.ListArticles(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	localized := make([]models.Article, len(articles))
	for i, article := range articles {
		localized[i] = editorial.Localize(article, locale)
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	category := r.URL.Query().Get("categoria")
	if category == "" {
		category = "Todos"
	}
	maxMinutes := intParam(r, "tempo", 0)

	saved := bookmarksFromRequest(r)
	filtering := editorial.Filtering(search, category, maxMinutes)
	filtered := editorial.Filter(localized, editorial.FilterOptions{
		Search:     search,
		Category:   category,
		MaxMinutes: maxMinutes,
	})
	lead, secondary, grid := editorial.Segment(filtered, filtering)

	writeView(w, r, h.views, h.cfg, "app", "home", amarraData(r, h.site, map[string]any{
		"Title":          h.t(r, "home.title"),
		"ActiveNav":      "artigos",
		"EditionDate":    editorial.FormatDateIn(time.Now(), locale),
		"TotalArticles":  len(localized),
		"TotalResults":   len(filtered),
		"Filtering":      filtering,
		"Search":         search,
		"ActiveCategory": category,
		"MaxMinutes":     maxMinutes,
		"CategoryLinks":  categoryLinksFor(localized, search, category, maxMinutes, locale),
		"TimeLinks":      timeLinksFor(search, category, maxMinutes, locale, func(key string) string { return h.t(r, key) }),
		"ClearHref":      localePath(locale, editorial.CategoryHref(category, "", maxMinutes)),
		"Lead":           leadArticleView(lead, saved, locale),
		"Secondary":      newArticleViews(secondary, saved, locale),
		"Grid":           newArticleViews(grid, saved, locale),
	}), http.StatusOK)
}
