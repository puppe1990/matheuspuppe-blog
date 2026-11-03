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

func (h *ArticleHandler) t(r *http.Request, key string, args ...any) string {
	return i18n.CatalogOr(r, h.catalog).T(key, args...)
}

// findArticle resolves the request locale's slug: /en pages look up slug_en
// only, so a Portuguese slug never leaks into the English edition (and the
// other way around).
func (h *ArticleHandler) findArticle(r *http.Request, slug string) (models.Article, error) {
	if localeOf(r) == editorial.LocaleEN {
		return h.store.FindArticleBySlugEN(r.Context(), slug)
	}
	return h.store.FindArticleBySlug(r.Context(), slug)
}

func (h *ArticleHandler) Show(w http.ResponseWriter, r *http.Request, slug string) {
	article, err := h.findArticle(r, slug)
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
	article, err := h.findArticle(r, slug)
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
	errs := form.validate(func(key string) string { return h.t(r, key) })
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

	flash.Set(w, "notice", h.t(r, "artigo.comment_posted"), h.cfg.CookieSecure())
	http.Redirect(w, r, h.articleURL(r, article), http.StatusSeeOther)
}

func (h *ArticleHandler) ClapPost(w http.ResponseWriter, r *http.Request, slug string) {
	article, err := h.findArticle(r, slug)
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
	http.Redirect(w, r, h.articleURL(r, article), http.StatusSeeOther)
}

// articleURL is the reader path in the request locale.
func (h *ArticleHandler) articleURL(r *http.Request, article models.Article) string {
	locale := localeOf(r)
	slug := article.Slug
	if locale == editorial.LocaleEN && article.SlugEN != "" {
		slug = article.SlugEN
	}
	return articlePath(slug, locale)
}

func (form commentForm) validate(tr func(key string) string) validate.FieldErrors {
	var errs validate.FieldErrors
	if form.Name == "" || len([]rune(form.Name)) < 2 {
		errs.Add("nome", tr("artigo.invalid_name"))
	}
	if form.Content == "" || len([]rune(form.Content)) < 4 {
		errs.Add("conteudo", tr("artigo.invalid_content"))
	}
	return errs
}

func (h *ArticleHandler) renderNotFound(w http.ResponseWriter, r *http.Request) {
	renderNotFoundPage(w, r, h.views, h.cfg, h.site)
}

func (h *ArticleHandler) render(w http.ResponseWriter, r *http.Request, article models.Article, form commentForm, errs validate.FieldErrors, status int) {
	locale := localeOf(r)
	comments, err := h.store.ListComments(r.Context(), article.ID)
	if err != nil {
		httpx.ServerError(w, err, h.cfg)
		return
	}
	localized := editorial.Localize(article, locale)
	altURL := ""
	if editorial.HasEnglish(article) {
		if locale == editorial.LocaleEN {
			altURL = articlePath(article.Slug, editorial.LocalePT)
		} else {
			altURL = articlePath(article.SlugEN, editorial.LocaleEN)
		}
	}

	writeView(w, r, h.views, h.cfg, "app", "artigo", amarraData(r, h.site, map[string]any{
		"Title":         localized.Title,
		"OGImage":       "/og/" + localized.Slug,
		"OGArticle":     true,
		"OGDescription": localized.Subtitle,
		"ActiveNav":     "artigos",
		"AltURL":        altURL,
		"Article":       newArticleView(article, bookmarksFromRequest(r), locale),
		"Comments":      newCommentViews(comments, locale),
		"CommentCount":  len(comments),
		"Errors":        errs,
		"Form":          form,
	}), status)
}
