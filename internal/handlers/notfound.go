package handlers

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"
)

// renderNotFoundPage serves the designed 404 for unmatched routes and for
// articles that do not exist; the router also funnels path-param parse
// failures here once r.NotFound is registered.
func renderNotFoundPage(w http.ResponseWriter, r *http.Request, views *view.Renderer, cfg cais.Config, site meta.Site) {
	writeView(w, r, views, cfg, "app", "notfound", amarraData(r, site, map[string]any{
		"Title":     "Página não encontrada · Sintaxe",
		"ActiveNav": "",
	}), http.StatusNotFound)
}

type NotFoundHandler struct {
	views *view.Renderer
	site  meta.Site
	cfg   cais.Config
}

func NewNotFoundHandler(views *view.Renderer, site meta.Site, cfg cais.Config) *NotFoundHandler {
	return &NotFoundHandler{views: views, site: site, cfg: cfg}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	renderNotFoundPage(w, r, h.views, h.cfg, h.site)
}
