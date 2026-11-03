package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/db"
	appi18n "github.com/puppe1990/matheuspuppe-blog/internal/i18n"
	"github.com/puppe1990/matheuspuppe-blog/internal/store"
)

func testSite() meta.Site {
	return meta.Site{AppName: "Sintaxe", AppURL: "https://sintaxe.example.com"}
}

func projectRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("go.mod not found")
		}
		wd = parent
	}
}

func setupTestViews(t *testing.T) *view.Renderer {
	t.Helper()
	root := projectRoot(t)
	fsys := os.DirFS(filepath.Join(root, "web", "templates"))
	r, err := view.Load(fsys, appi18n.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func setupTestStore(t *testing.T) store.Store {
	t.Helper()
	s, err := store.NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// seedTestStore gives handler tests the real published catalog, so assertions
// can name actual essays.
func seedTestStore(t *testing.T) store.Store {
	t.Helper()
	s := setupTestStore(t)
	if err := db.RunSeeds(s); err != nil {
		t.Fatal(err)
	}
	return s
}

// postForm builds the request and the recorder, so handler tests can call
// handlers directly without the router wiring.
func postForm(t *testing.T, target string, form url.Values, cookies ...*http.Cookie) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return req, httptest.NewRecorder()
}
