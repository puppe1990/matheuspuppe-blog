package handlers

import (
	"errors"
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// AuthorHandler serves one columnist's page with their published essays.
type AuthorHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewAuthorHandler(views *view.Renderer, st store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *AuthorHandler {
	return &AuthorHandler{views: views, store: st, site: site, catalog: catalog, cfg: cfg}
}

func (h *AuthorHandler) Serve(w http.ResponseWriter, r *http.Request, slug string) {
	locale := localeOf(r)
	author, err := h.store.FindAuthorBySlug(r.Context(), slug)
	if errors.Is(err, store.ErrNotFound) {
		renderNotFoundPage(w, r, h.views, h.cfg, h.site)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	articles, err := h.store.ListArticlesByAuthor(r.Context(), author.ID)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	localized := make([]models.Article, len(articles))
	for i, article := range articles {
		localized[i] = editorial.Localize(article, locale)
	}

	writeView(w, r, h.views, h.cfg, "app", "autor", amarraData(r, h.site, map[string]any{
		"Title":     author.Name + " · Sintaxe",
		"ActiveNav": "sobre",
		"Author":    editorial.LocalizeAuthor(author, locale),
		"Articles":  newArticleViews(articles, bookmarksFromRequest(r), locale),
		"OGImage":   author.AvatarURL,
	}), http.StatusOK)
}
