package handlers

import (
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/validate"

	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// NewsletterHandler registers weekly-boletim subscribers. It is reachable from
// the dialog on any page, so errors and success both flash and redirect back
// instead of rendering a page of their own.
type NewsletterHandler struct {
	store   store.Store
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewNewsletterHandler(st store.Store, catalog *i18n.Catalog, cfg cais.Config) *NewsletterHandler {
	return &NewsletterHandler{store: st, catalog: catalog, cfg: cfg}
}

func (h *NewsletterHandler) t(r *http.Request, key string, args ...any) string {
	return i18n.CatalogOr(r, h.catalog).T(key, args...)
}

func (h *NewsletterHandler) Post(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	next := sanitizeNext(r.FormValue("next"))

	if err := validate.Email(email); err != nil {
		flash.Set(w, "error", h.t(r, "newsletter.invalid"), h.cfg.CookieSecure())
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	if err := h.store.UpsertSubscriber(r.Context(), email); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	flash.Set(w, "notice", h.t(r, "newsletter.confirmed", email), h.cfg.CookieSecure())
	http.Redirect(w, r, next, http.StatusSeeOther)
}
