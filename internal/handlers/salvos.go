package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// SavedHandler serves the personal library kept in the sintaxe_salvos cookie.
type SavedHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewSavedHandler(views *view.Renderer, st store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *SavedHandler {
	return &SavedHandler{views: views, store: st, site: site, catalog: catalog, cfg: cfg}
}

func (h *SavedHandler) t(r *http.Request, key string) string {
	return i18n.CatalogOr(r, h.catalog).T(key)
}

func (h *SavedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	locale := localeOf(r)
	saved := bookmarksFromRequest(r)
	if saved == nil {
		// Distinguish "never saved" from "show everything": an empty (non-nil)
		// set narrows the filter to zero rows.
		saved = []string{}
	}

	articles, err := h.store.ListArticles(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	// Bookmarks are keyed by the canonical (Portuguese) slug; keep them before
	// localizing so the library works across languages.
	bookmarked := make([]models.Article, 0, len(saved))
	for _, article := range articles {
		if containsSlug(saved, article.Slug) {
			bookmarked = append(bookmarked, article)
		}
	}
	localized := make([]models.Article, len(bookmarked))
	for i, article := range bookmarked {
		localized[i] = editorial.Localize(article, locale)
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	category := r.URL.Query().Get("categoria")
	if category == "" {
		category = "Todos"
	}
	maxMinutes := intParam(r, "tempo", 0)

	filtering := editorial.Filtering(search, category, maxMinutes)
	filtered := editorial.Filter(localized, editorial.FilterOptions{
		Search:     search,
		Category:   category,
		MaxMinutes: maxMinutes,
	})
	lead, secondary, grid := editorial.Segment(filtered, filtering)

	writeView(w, r, h.views, h.cfg, "app", "salvos", amarraData(r, h.site, map[string]any{
		"Title":          h.t(r, "salvos.title") + " · Sintaxe",
		"ActiveNav":      "salvos",
		"SavedCount":     len(saved),
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

// Toggle flips one article in the saved cookie and returns the reader to the
// page they came from (sanitized to same-site paths).
func (h *SavedHandler) Toggle(w http.ResponseWriter, r *http.Request, slug string) {
	if _, err := h.store.FindArticleBySlug(r.Context(), slug); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			renderNotFoundPage(w, r, h.views, h.cfg, h.site)
			return
		}
		httpx.ServerError(w, err, h.cfg)
		return
	}

	saved := editorial.ToggleBookmark(bookmarksFromRequest(r), slug)
	writeBookmarkCookie(w, saved, h.cfg.CookieSecure())
	http.Redirect(w, r, sanitizeNext(r.FormValue("next")), http.StatusSeeOther)
}
