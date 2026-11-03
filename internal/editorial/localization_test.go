package editorial

import (
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func englishArticle() models.Article {
	article := fakeArticles()[0]
	article.SlugEN = "the-first-essay"
	article.TitleEN = "The First Essay"
	article.SubtitleEN = "A dissection of the system."
	article.LeadParagraphEN = "The first paragraph of the main essay."
	article.TagsEN = []string{"Artificial Intelligence"}
	article.BlocksEN = []models.ContentBlock{{Type: "paragraph", Text: "Body in English."}}
	article.Author.RoleEN = "Software engineer"
	article.Author.BioEN = "Writes about software."
	return article
}

func TestLocalize_English(t *testing.T) {
	article := englishArticle()
	got := Localize(article, LocaleEN)
	if got.Slug != "the-first-essay" || got.Title != "The First Essay" {
		t.Fatalf("slug/title not localized: %+v", got)
	}
	if got.Subtitle != article.SubtitleEN || got.LeadParagraph != article.LeadParagraphEN {
		t.Fatalf("prose not localized: %+v", got)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "Artificial Intelligence" {
		t.Fatalf("tags = %v", got.Tags)
	}
	if len(got.Blocks) != 1 || got.Blocks[0].Text != "Body in English." {
		t.Fatalf("blocks = %v", got.Blocks)
	}
	if got.Author.Role != "Software engineer" || got.Author.Bio != "Writes about software." {
		t.Fatalf("author not localized: %+v", got.Author)
	}
}

func TestLocalize_FallsBackToPortuguese(t *testing.T) {
	article := fakeArticles()[0] // no English edition
	for _, locale := range []string{LocalePT, LocaleEN} {
		got := Localize(article, locale)
		if got.Slug != article.Slug || got.Title != article.Title {
			t.Fatalf("%s: untranslated article changed: %+v", locale, got)
		}
	}
}

func TestLocalize_KeepsPortugueseUntouched(t *testing.T) {
	article := englishArticle()
	got := Localize(article, LocalePT)
	if got.Slug != article.Slug || got.Title != article.Title || got.Subtitle != article.Subtitle {
		t.Fatalf("pt localize changed the article: %+v", got)
	}
	if got.Author.Role != article.Author.Role {
		t.Fatalf("pt localize changed the author: %+v", got.Author)
	}
}

func TestHasEnglish(t *testing.T) {
	if HasEnglish(fakeArticles()[0]) {
		t.Fatal("article without EN fields reported as translated")
	}
	if !HasEnglish(englishArticle()) {
		t.Fatal("translated article not detected")
	}
	partial := englishArticle()
	partial.BlocksEN = nil
	if HasEnglish(partial) {
		t.Fatal("translation without body reported as complete")
	}
}

func TestCategoryLabel(t *testing.T) {
	if got := CategoryLabel("IA & Modelos", LocaleEN); got != "AI & Models" {
		t.Fatalf("en label = %q", got)
	}
	if got := CategoryLabel("Engenharia de Sistemas", LocaleEN); got != "Systems Engineering" {
		t.Fatalf("en label = %q", got)
	}
	if got := CategoryLabel("IA & Modelos", LocalePT); got != "IA & Modelos" {
		t.Fatalf("pt label = %q", got)
	}
	if got := CategoryLabel("Desconhecida", LocaleEN); got != "Desconhecida" {
		t.Fatalf("fallback = %q", got)
	}
}

func TestFormatDateIn(t *testing.T) {
	date := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	if got := FormatDateIn(date, LocaleEN); got != "November 3, 2026" {
		t.Fatalf("en date = %q", got)
	}
	if got := FormatDateIn(date, LocalePT); got != "3 de novembro de 2026" {
		t.Fatalf("pt date = %q", got)
	}
}
