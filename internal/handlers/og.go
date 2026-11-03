package handlers

import (
	"errors"
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/ogcard"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// OGHandler paints a 1200×630 PNG with the essay title. Crawlers fetch this
// URL from og:image; it must not set cookies.
type OGHandler struct {
	store store.Store
	site  meta.Site
	cfg   cais.Config
}

func NewOGHandler(st store.Store, site meta.Site, cfg cais.Config) *OGHandler {
	return &OGHandler{store: st, site: site, cfg: cfg}
}

func (h *OGHandler) Serve(w http.ResponseWriter, r *http.Request, slug string) {
	article, locale, err := h.lookup(r, slug)
	if errors.Is(err, store.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	localized := editorial.Localize(article, locale)
	card := ogcard.ArticleCard(localized, ogcard.HostOf(h.site.AppURL))
	card.Kicker = editorial.CategoryLabel(article.Category, locale)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if err := ogcard.WritePNG(w, card); err != nil {
		httpx.ServerError(w, err, h.cfg)
	}
}

func (h *OGHandler) lookup(r *http.Request, slug string) (models.Article, string, error) {
	article, err := h.store.FindArticleBySlug(r.Context(), slug)
	if err == nil {
		return article, editorial.LocalePT, nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return models.Article{}, "", err
	}
	article, err = h.store.FindArticleBySlugEN(r.Context(), slug)
	if err != nil {
		return models.Article{}, "", err
	}
	return article, editorial.LocaleEN, nil
}
