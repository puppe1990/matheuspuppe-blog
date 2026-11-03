package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

const leadSlug = "fronteira-modelos-raciocinio-agentes-autonomos"

func newArticleHandler(t *testing.T) (*ArticleHandler, store.Store) {
	t.Helper()
	st := seedTestStore(t)
	return NewArticleHandler(setupTestViews(t), st, testSite(), appi18n.DefaultCatalog(), cais.Config{}), st
}

func showArticle(t *testing.T, slug string) *httptest.ResponseRecorder {
	t.Helper()
	h, _ := newArticleHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/artigo/"+slug, nil)
	rr := httptest.NewRecorder()
	h.Show(rr, req, slug)
	return rr
}

func TestArticleHandler_ShowRendersEssay(t *testing.T) {
	rr := showArticle(t, leadSlug)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	for _, needle := range []string{
		"A Nova Fronteira dos Agentes Autônomos",
		"Uma dissecação arquitetural do test-time compute",
		"Dra. Camila Siqueira",
		"Conclusões Fundamentais &amp; Invariants",
		"O Esgotamento do Pré-Treino Passivo",
		"#Inteligência Artificial",
		"Discussão da Comunidade Técnica",
		"Thiago Nogueira",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("body missing %q", needle)
		}
	}
}

func TestArticleHandler_ShowRendersBenchmarkBlock(t *testing.T) {
	body := showArticle(t, leadSlug).Body.String()
	if !strings.Contains(body, "SWE-Bench Pro") {
		t.Errorf("body missing the interactive benchmark block")
	}
	if !strings.Contains(body, "Precisão Lógica") || !strings.Contains(body, "88%") {
		t.Errorf("body missing benchmark metric values")
	}
}

func TestArticleHandler_ShowIncrementsViews(t *testing.T) {
	h, st := newArticleHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/artigo/"+leadSlug, nil)
	h.Show(httptest.NewRecorder(), req, leadSlug)

	article, err := st.FindArticleBySlug(context.Background(), leadSlug)
	if err != nil {
		t.Fatal(err)
	}
	if article.Views != 4891 {
		t.Errorf("views = %d, want 4891", article.Views)
	}
}

func TestArticleHandler_ShowMissingIs404(t *testing.T) {
	rr := showArticle(t, "nao-existe")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "id=\"amarra-main\"") {
		t.Errorf("404 body must render the journal chrome")
	}
}

func TestArticleHandler_CommentPostPersists(t *testing.T) {
	h, st := newArticleHandler(t)
	form := url.Values{
		"nome":     {"Thiago Nogueira"},
		"cargo":    {"Engenheiro de ML"},
		"conteudo": {"Excelente análise sobre verificação formal de esquemas."},
	}
	req, rr := postForm(t, "/artigo/"+leadSlug+"/comentarios", form)
	h.CommentPost(rr, req, leadSlug)
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "/artigo/"+leadSlug {
		t.Fatalf("Location = %q", got)
	}

	article, err := st.FindArticleBySlug(context.Background(), leadSlug)
	if err != nil {
		t.Fatal(err)
	}
	comments, err := st.ListComments(context.Background(), article.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 3 {
		t.Fatalf("comments = %d, want 3 (2 seeded + 1 new)", len(comments))
	}
	if comments[0].AuthorName != "Thiago Nogueira" {
		t.Fatalf("newest comment = %+v", comments[0])
	}
}

func TestArticleHandler_CommentPostDefaultsReaderRole(t *testing.T) {
	h, st := newArticleHandler(t)
	form := url.Values{
		"nome":     {"Leitora Anônima"},
		"conteudo": {"Comentário sem cargo informado."},
	}
	req, rr := postForm(t, "/artigo/"+leadSlug+"/comentarios", form)
	h.CommentPost(rr, req, leadSlug)

	article, _ := st.FindArticleBySlug(context.Background(), leadSlug)
	comments, _ := st.ListComments(context.Background(), article.ID)
	if comments[0].AuthorRole != "Leitor Sintaxe" {
		t.Fatalf("role = %q, want Leitor Sintaxe", comments[0].AuthorRole)
	}
}

func TestArticleHandler_CommentPostValidation(t *testing.T) {
	h, st := newArticleHandler(t)
	form := url.Values{"nome": {""}, "conteudo": {"curto"}}
	req, rr := postForm(t, "/artigo/"+leadSlug+"/comentarios", form)
	h.CommentPost(rr, req, leadSlug)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "Informe seu nome") {
		t.Errorf("body missing the name error, got: %.300s", body)
	}
	article, _ := st.FindArticleBySlug(context.Background(), leadSlug)
	comments, _ := st.ListComments(context.Background(), article.ID)
	if len(comments) != 2 {
		t.Fatalf("invalid comment was stored: %d comments", len(comments))
	}
}

func TestArticleHandler_ClapPost(t *testing.T) {
	h, st := newArticleHandler(t)
	req, rr := postForm(t, "/artigo/"+leadSlug+"/aplausos", url.Values{})
	h.ClapPost(rr, req, leadSlug)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	article, _ := st.FindArticleBySlug(context.Background(), leadSlug)
	if article.Claps != 343 {
		t.Fatalf("claps = %d, want 343", article.Claps)
	}
}

func TestArticleHandler_ShowScheduledIs404(t *testing.T) {
	h, st := newArticleHandler(t)
	scheduled := seedScheduledEssay(t, st)

	rr := httptest.NewRecorder()
	h.Show(rr, httptest.NewRequest(http.MethodGet, "/artigo/"+scheduled.Slug, nil), scheduled.Slug)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("scheduled essay = %d, want 404 before its day", rr.Code)
	}
}
