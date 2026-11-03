package i18n

import "testing"

func TestDefaultCatalog_editorialCopy(t *testing.T) {
	c := DefaultCatalog()
	if got := c.T("nav.artigos"); got != "Artigos" {
		t.Errorf("T(nav.artigos) = %q", got)
	}
	if got := c.T("layout.theme_light"); got != "Modo claro" {
		t.Errorf("T(layout.theme_light) = %q", got)
	}
}

func TestNewCatalog_portuguese(t *testing.T) {
	c := NewCatalog("pt-BR")
	if got := c.T("nav.sobre"); got != "Sobre a Redação" {
		t.Errorf("T(nav.sobre) = %q", got)
	}
	if c.HTMLLang() != "pt-BR" {
		t.Errorf("HTMLLang() = %q, want pt-BR", c.HTMLLang())
	}
}
