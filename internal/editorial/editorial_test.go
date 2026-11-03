package editorial

import (
	"testing"
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func fakeAuthor() models.Author {
	return models.Author{ID: 1, Slug: "helena", Name: "Helena Vance", Role: "Diretora de Engenharia & Ensaísta", Bio: "Escreve sobre sociologia do código."}
}

func fakeArticles() []models.Article {
	return []models.Article{
		{
			ID:              1,
			Slug:            "artigo-ia",
			Title:           "A Nova Fronteira dos Agentes Autônomos",
			Subtitle:        "Por que o raciocínio superou a escala pura",
			LeadParagraph:   "Modelos de linguagem estão migrando de geração estocástica para planejamento multi-passos.",
			Category:        "IA & Modelos",
			ReadTimeMinutes: 7,
			PublishedAt:     time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC),
			Author:          fakeAuthor(),
			Tags:            []string{"IA", "Agentes", "LLM"},
			Blocks: []models.ContentBlock{
				{Type: "paragraph", Text: "Primeiro parágrafo com dez palavras para teste de contagem."},
				{Type: "takeaways", Items: []string{"Ponto um", "Ponto dois"}},
			},
			Claps: 120,
			Views: 1400,
			Tier:  1,
		},
		{
			ID:              2,
			Slug:            "artigo-rust",
			Title:           "Migrando de Go para Rust em 2026",
			Subtitle:        "Zero garbage collection em serviços de baixa latência",
			LeadParagraph:   "Como eliminamos spikes de p99 em infraestrutura de pagamentos.",
			Category:        "Engenharia de Sistemas",
			ReadTimeMinutes: 4,
			PublishedAt:     time.Date(2026, 2, 10, 14, 30, 0, 0, time.UTC),
			Author:          fakeAuthor(),
			Tags:            []string{"Rust", "Go", "Backend"},
			Blocks:          []models.ContentBlock{{Type: "paragraph", Text: "Texto do segundo artigo sobre Rust."}},
			Claps:           85,
			Views:           950,
			Tier:            2,
		},
		{
			ID:              3,
			Slug:            "artigo-chips",
			Title:           "O Fim da Lei de Moore",
			Subtitle:        "Aceleradores de silício especializados",
			LeadParagraph:   "A física dos nós de 2nm impulsionou a era dos chiplets 3D.",
			Category:        "Hardware & Chips",
			ReadTimeMinutes: 6,
			PublishedAt:     time.Date(2026, 1, 5, 9, 0, 0, 0, time.UTC),
			Author:          fakeAuthor(),
			Tags:            []string{"Hardware", "Silício"},
			Claps:           10,
			Views:           100,
			Tier:            2,
		},
		{
			ID:              4,
			Slug:            "artigo-extra",
			Title:           "Caderno Extra",
			Subtitle:        "Despachos diversos",
			LeadParagraph:   "Um ensaio padrão da grade.",
			Category:        "Opinião & Futuro",
			ReadTimeMinutes: 3,
			PublishedAt:     time.Date(2025, 12, 1, 8, 0, 0, 0, time.UTC),
			Author:          fakeAuthor(),
			Tags:            []string{"Carreira"},
			Claps:           1,
			Views:           10,
			Tier:            3,
		},
	}
}

func TestFilter_SearchMatchesTitle(t *testing.T) {
	got := Filter(fakeArticles(), FilterOptions{Search: "Agentes"})
	if len(got) != 1 || got[0].Slug != "artigo-ia" {
		t.Fatalf("got %d articles (%v), want only artigo-ia", len(got), slugs(got))
	}
}

func TestFilter_SearchMatchesAuthorAndTags(t *testing.T) {
	byAuthor := Filter(fakeArticles(), FilterOptions{Search: "helena vance"})
	if len(byAuthor) != 4 {
		t.Fatalf("author search got %d, want 4", len(byAuthor))
	}
	byTag := Filter(fakeArticles(), FilterOptions{Search: "rust"})
	if len(byTag) != 1 || byTag[0].Slug != "artigo-rust" {
		t.Fatalf("tag search got %v, want artigo-rust", slugs(byTag))
	}
}

func TestFilter_Category(t *testing.T) {
	got := Filter(fakeArticles(), FilterOptions{Category: "Engenharia de Sistemas"})
	if len(got) != 1 || got[0].Slug != "artigo-rust" {
		t.Fatalf("got %v, want artigo-rust", slugs(got))
	}
}

func TestFilter_TodosMeansAll(t *testing.T) {
	got := Filter(fakeArticles(), FilterOptions{Category: "Todos"})
	if len(got) != 4 {
		t.Fatalf("got %d, want 4", len(got))
	}
}

func TestFilter_MaxMinutes(t *testing.T) {
	got := Filter(fakeArticles(), FilterOptions{MaxMinutes: 5})
	if len(got) != 2 {
		t.Fatalf("got %v, want the 4- and 3-minute pieces", slugs(got))
	}
}

func TestFilter_OnlySlugsRestrictsToSaved(t *testing.T) {
	got := Filter(fakeArticles(), FilterOptions{OnlySlugs: []string{"artigo-ia"}})
	if len(got) != 1 || got[0].Slug != "artigo-ia" {
		t.Fatalf("got %v, want artigo-ia", slugs(got))
	}
	// A saved view still obeys the search box.
	got = Filter(fakeArticles(), FilterOptions{OnlySlugs: []string{"artigo-ia", "artigo-rust"}, Search: "rust"})
	if len(got) != 1 || got[0].Slug != "artigo-rust" {
		t.Fatalf("got %v, want artigo-rust", slugs(got))
	}
}

func TestSegment_EditorialTiers(t *testing.T) {
	lead, secondary, grid := Segment(fakeArticles(), false)
	if lead == nil || lead.Slug != "artigo-ia" {
		t.Fatalf("lead = %v, want artigo-ia (tier 1)", lead)
	}
	if len(secondary) != 2 || secondary[0].Slug != "artigo-rust" || secondary[1].Slug != "artigo-chips" {
		t.Fatalf("secondary = %v, want rust+chips (tier 2)", slugs(secondary))
	}
	if len(grid) != 1 || grid[0].Slug != "artigo-extra" {
		t.Fatalf("grid = %v, want artigo-extra", slugs(grid))
	}
}

func TestSegment_FilteringTakesFirstAsLead(t *testing.T) {
	articles := fakeArticles()
	lead, secondary, grid := Segment(articles, true)
	if lead == nil || lead.Slug != "artigo-ia" {
		t.Fatalf("lead = %v, want first article while filtering", lead)
	}
	if len(secondary) != 2 || len(grid) != 1 {
		t.Fatalf("secondary=%d grid=%d, want 2/1", len(secondary), len(grid))
	}
}

func TestSegment_MissingTierOneFallsBackToFirst(t *testing.T) {
	articles := fakeArticles()[1:]
	lead, _, _ := Segment(articles, false)
	if lead == nil || lead.Slug != "artigo-rust" {
		t.Fatalf("lead = %v, want first article as fallback", lead)
	}
}

func TestSegment_Empty(t *testing.T) {
	lead, secondary, grid := Segment(nil, false)
	if lead != nil || len(secondary) != 0 || len(grid) != 0 {
		t.Fatalf("expected empty segments, got %v %v %v", lead, secondary, grid)
	}
}

func TestReadingTime(t *testing.T) {
	if got := ReadingTime(200); got != 1 {
		t.Errorf("ReadingTime(200) = %d, want 1", got)
	}
	if got := ReadingTime(201); got != 2 {
		t.Errorf("ReadingTime(201) = %d, want 2", got)
	}
	if got := ReadingTime(0); got != 1 {
		t.Errorf("ReadingTime(0) = %d, want 1", got)
	}
}

func TestWordCount(t *testing.T) {
	blocks := fakeArticles()[0].Blocks
	// 9 words in the paragraph + 2 + 2 in the takeaways = 13.
	if got := WordCount(blocks); got != 13 {
		t.Errorf("WordCount = %d, want 13", got)
	}
}

func TestToggleBookmark(t *testing.T) {
	added := ToggleBookmark([]string{"artigo-ia"}, "artigo-rust")
	if len(added) != 2 || added[1] != "artigo-rust" {
		t.Fatalf("added = %v", added)
	}
	removed := ToggleBookmark(added, "artigo-ia")
	if len(removed) != 1 || removed[0] != "artigo-rust" {
		t.Fatalf("removed = %v", removed)
	}
}

func TestAddClaps(t *testing.T) {
	if got := AddClaps(10, 1); got != 11 {
		t.Errorf("AddClaps(10,1) = %d", got)
	}
	if got := AddClaps(10, 5); got != 15 {
		t.Errorf("AddClaps(10,5) = %d", got)
	}
	if got := AddClaps(10, 0); got != 11 {
		t.Errorf("AddClaps(10,0) = %d, want minimum increment of 1", got)
	}
}

func TestCalculateRemainingReadingTime(t *testing.T) {
	cases := []struct {
		minutes, percent, want int
	}{
		{10, 0, 10},
		{10, 50, 5},
		{10, 90, 1},
		{10, 100, 0},
		{10, 120, 0},
	}
	for _, c := range cases {
		if got := RemainingReadingTime(c.minutes, c.percent); got != c.want {
			t.Errorf("RemainingReadingTime(%d, %d) = %d, want %d", c.minutes, c.percent, got, c.want)
		}
	}
}

func TestCategoryCounts(t *testing.T) {
	counts := CategoryCounts(fakeArticles())
	if counts["Todos"] != 4 {
		t.Errorf("Todos = %d, want 4", counts["Todos"])
	}
	if counts["IA & Modelos"] != 1 || counts["Engenharia de Sistemas"] != 1 {
		t.Errorf("counts = %v", counts)
	}
}

func TestShareText(t *testing.T) {
	got := ShareText("Artigo Tech", "Dissecação de código")
	if got != "Artigo Tech — Dissecação de código" {
		t.Errorf("ShareText = %q", got)
	}
	if got := ShareText("Só título", ""); got != "Só título" {
		t.Errorf("ShareText without subtitle = %q", got)
	}
}

func TestFormatDate(t *testing.T) {
	got := FormatDate(time.Date(2026, 3, 28, 9, 0, 0, 0, time.UTC))
	if got != "28 de março de 2026" {
		t.Errorf("FormatDate = %q, want 28 de março de 2026", got)
	}
	got = FormatDate(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC))
	if got != "2 de janeiro de 2026" {
		t.Errorf("FormatDate = %q, want 2 de janeiro de 2026", got)
	}
}

func TestFormatThousands(t *testing.T) {
	if got := FormatThousands(4890); got != "4.890" {
		t.Errorf("FormatThousands(4890) = %q", got)
	}
	if got := FormatThousands(999); got != "999" {
		t.Errorf("FormatThousands(999) = %q", got)
	}
	if got := FormatThousands(1234567); got != "1.234.567" {
		t.Errorf("FormatThousands(1234567) = %q", got)
	}
}

func TestParseBookmarkCookie(t *testing.T) {
	got := ParseBookmarkCookie(" artigo-ia , artigo-rust,artigo-ia ,, ")
	want := []string{"artigo-ia", "artigo-rust"}
	if len(got) != len(want) {
		t.Fatalf("ParseBookmarkCookie = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ParseBookmarkCookie = %v, want %v", got, want)
		}
	}
	if got := ParseBookmarkCookie(""); len(got) != 0 {
		t.Errorf("empty cookie = %v, want none", got)
	}
}

func TestEncodeBookmarksRoundTrip(t *testing.T) {
	slugs := []string{"artigo-ia", "artigo-rust"}
	if got := EncodeBookmarks(slugs); got != "artigo-ia,artigo-rust" {
		t.Errorf("EncodeBookmarks = %q", got)
	}
	round := ParseBookmarkCookie(EncodeBookmarks(slugs))
	if len(round) != 2 || round[0] != slugs[0] || round[1] != slugs[1] {
		t.Errorf("round trip = %v", round)
	}
}

func TestClassifyCategoryURL(t *testing.T) {
	href := CategoryHref("IA & Modelos", "rust", 5)
	if href != "/?categoria=IA+%26+Modelos&q=rust&tempo=5" {
		t.Errorf("CategoryHref = %q", href)
	}
	if href := CategoryHref("Todos", "", 0); href != "/" {
		t.Errorf("CategoryHref clean = %q, want /", href)
	}
	if href := CategoryHref("", "", 5); href != "/?tempo=5" {
		t.Errorf("CategoryHref tempo = %q", href)
	}
}

func slugs(articles []models.Article) []string {
	out := make([]string, len(articles))
	for i, a := range articles {
		out[i] = a.Slug
	}
	return out
}
