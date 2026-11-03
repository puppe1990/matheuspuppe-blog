package handlers

import (
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/validate"

	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// NewsletterHandler registers weekly-boletim subscribers. It is reachable from
// the dialog on any page, so errors and success both flash and redirect back
// instead of rendering a page of their own.
type NewsletterHandler struct {
	store store.Store
	cfg   cais.Config
}

func NewNewsletterHandler(st store.Store, cfg cais.Config) *NewsletterHandler {
	return &NewsletterHandler{store: st, cfg: cfg}
}

func (h *NewsletterHandler) Post(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	next := sanitizeNext(r.FormValue("next"))

	if err := validate.Email(email); err != nil {
		flash.Set(w, "error", "Informe um email válido para assinar o boletim.", h.cfg.CookieSecure())
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	if err := h.store.UpsertSubscriber(r.Context(), email); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	flash.Set(w, "notice", "Inscrição confirmada! Enviamos uma mensagem de boas-vindas para "+email+".", h.cfg.CookieSecure())
	http.Redirect(w, r, next, http.StatusSeeOther)
}
