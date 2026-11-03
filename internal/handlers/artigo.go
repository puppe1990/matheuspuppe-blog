package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/flash"
	"github.com/puppe1990/amarra-cais/pkg/cais/httpx"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"
	"github.com/puppe1990/amarra-cais/pkg/cais/validate"

	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
	"github.com/puppe1990/matheuspuppe-blog/internal/models"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

// commentForm keeps what the reader typed so a 422 re-renders the same essay
// with the fields filled.
type commentForm struct {
	Name    string
	Role    string
	Content string
}

type ArticleHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewArticleHandler(views *view.Renderer, st store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *ArticleHandler {
	return &ArticleHandler{views: views, store: st, site: site, catalog: catalog, cfg: cfg}
}

func (h *ArticleHandler) Show(w http.ResponseWriter, r *http.Request, slug string) {
	article, err := h.store.FindArticleBySlug(r.Context(), slug)
	if errors.Is(err, store.ErrNotFound) {
		h.renderNotFound(w, r)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	// Views are analytics: a failed bump must not break the read, but the page
	// still shows the visit it just registered.
	if err := h.store.IncrementArticleViews(r.Context(), article.ID, 1); err == nil {
		article.Views++
	}
	h.render(w, r, article, commentForm{}, nil, http.StatusOK)
}

func (h *ArticleHandler) CommentPost(w http.ResponseWriter, r *http.Request, slug string) {
	article, err := h.store.FindArticleBySlug(r.Context(), slug)
	if errors.Is(err, store.ErrNotFound) {
		h.renderNotFound(w, r)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	form := commentForm{
		Name:    strings.TrimSpace(r.FormValue("nome")),
		Role:    strings.TrimSpace(r.FormValue("cargo")),
		Content: strings.TrimSpace(r.FormValue("conteudo")),
	}
	errs := form.validate()
	if errs.Any() {
		h.render(w, r, article, form, errs, http.StatusUnprocessableEntity)
		return
	}

	role := form.Role
	if role == "" {
		role = "Leitor Sintaxe"
	}
	comment := models.Comment{
		ArticleID:  article.ID,
		AuthorName: form.Name,
		AuthorRole: role,
		Content:    form.Content,
		CreatedAt:  time.Now().UTC(),
	}
	if _, err := h.store.AddComment(r.Context(), comment); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	flash.Set(w, "notice", "Seu comentário foi publicado.", h.cfg.CookieSecure())
	http.Redirect(w, r, "/artigo/"+article.Slug, http.StatusSeeOther)
}

func (h *ArticleHandler) ClapPost(w http.ResponseWriter, r *http.Request, slug string) {
	article, err := h.store.FindArticleBySlug(r.Context(), slug)
	if errors.Is(err, store.ErrNotFound) {
		h.renderNotFound(w, r)
		return
	}
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	// The shipped rule keeps a minimum increment of one clap per request.
	delta := editorial.AddClaps(0, 1)
	if _, err := h.store.AddArticleClaps(r.Context(), article.ID, delta); err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	http.Redirect(w, r, "/artigo/"+article.Slug, http.StatusSeeOther)
}

func (form commentForm) validate() validate.FieldErrors {
	var errs validate.FieldErrors
	switch {
	case form.Name == "":
		errs.Add("nome", "Informe seu nome.")
	case len([]rune(form.Name)) < 2:
		errs.Add("nome", "O nome precisa de ao menos 2 caracteres.")
	}
	switch {
	case form.Content == "":
		errs.Add("conteudo", "Escreva sua reflexão antes de publicar.")
	case len([]rune(form.Content)) < 4:
		errs.Add("conteudo", "O comentário precisa de ao menos 4 caracteres.")
	}
	return errs
}

func (h *ArticleHandler) renderNotFound(w http.ResponseWriter, r *http.Request) {
	renderNotFoundPage(w, r, h.views, h.cfg, h.site)
}

func (h *ArticleHandler) render(w http.ResponseWriter, r *http.Request, article models.Article, form commentForm, errs validate.FieldErrors, status int) {
	comments, err := h.store.ListComments(r.Context(), article.ID)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}

	writeView(w, r, h.views, h.cfg, "app", "artigo", amarraData(r, h.site, map[string]any{
		"Title":         article.Title,
		"OGImage":       article.CoverImage,
		"OGDescription": article.Subtitle,
		"ActiveNav":     "artigos",
		"Article":       newArticleView(article, bookmarksFromRequest(r)),
		"Comments":      newCommentViews(comments),
		"CommentCount":  len(comments),
		"Errors":        errs,
		"Form":          form,
	}), status)
}
