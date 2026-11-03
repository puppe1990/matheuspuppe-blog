package i18n

import "testing"

func TestDefaultCatalog_englishEdition(t *testing.T) {
	c := DefaultCatalog()
	if got := c.T("nav.artigos"); got != "Essays" {
		t.Errorf("T(nav.artigos) = %q", got)
	}
	if got := c.T("layout.theme_light"); got != "Light mode" {
		t.Errorf("T(layout.theme_light) = %q", got)
	}
	if got := c.T("artigo.published", "November 3, 2026"); got != "Published November 3, 2026" {
		t.Errorf("T(artigo.published) = %q", got)
	}
}

func TestNewCatalog_portuguese(t *testing.T) {
	c := NewCatalog("pt-BR")
	if got := c.T("nav.sobre"); got != "Sobre a redação" {
		t.Errorf("T(nav.sobre) = %q", got)
	}
	if c.HTMLLang() != "pt-BR" {
		t.Errorf("HTMLLang() = %q, want pt-BR", c.HTMLLang())
	}
}
