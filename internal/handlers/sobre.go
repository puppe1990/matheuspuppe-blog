package handlers

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// SobreHandler serves the editorial manifesto and the author board.
type SobreHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewSobreHandler(views *view.Renderer, st store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *SobreHandler {
	return &SobreHandler{views: views, store: st, site: site, catalog: catalog, cfg: cfg}
}

func (h *SobreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authors, err := h.store.ListAuthors(r.Context())
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	writeView(w, r, h.views, h.cfg, "app", "sobre", amarraData(r, h.site, map[string]any{
		"Title":     "Manifesto Editorial · Sintaxe",
		"ActiveNav": "sobre",
		"Authors":   authors,
	}), http.StatusOK)
}
