package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// amarraEssay é o terceiro ensaio: por que o framework existe, o que ele troca
// e o que a aposta custa. Os números vêm do repositório público.
func amarraEssay() models.Article {
	en := amarraEssayEn()
	return models.Article{
		Slug:            "por-que-eu-construi-o-amarra-cais",
		Title:           "Por que eu construí o Amarra-cais",
		Subtitle:        "Um fork do Cais em 5 de setembro, 687 commits e 1.148 testes depois: a aposta de que mini app não precisa de SPA, e o que ela custa.",
		LeadParagraph:   "O Amarra-cais nasceu de um incômodo específico: cada mini app meu carregava uma SPA inteira para mostrar tabela e formulário. Em 5 de setembro de 2026 eu importei o Cais (Inertia + Svelte) como ponto de partida e comecei um fork HTML-first. Router, sessão, CSRF, fila de jobs e SQLite ficaram. O cliente que monta tela saiu.",
		Category:        "Engenharia de Sistemas",
		ReadTimeMinutes: 7,
		PublishedAt:     time.Date(2026, 11, 3, 19, 55, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_interface_design.jpg",
		Author:          authorRef("matheus-puppe"),
		Tags:            []string{"Amarra-cais", "Go", "Arquitetura", "Frameworks", "HTML"},
		Tier:            2,
		SlugEN:          en.Slug,
		TitleEN:         en.Title,
		SubtitleEN:      en.Subtitle,
		LeadParagraphEN: en.Lead,
		TagsEN:          en.Tags,
		BlocksEN:        en.Blocks,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"O começo foi um commit de chore: importar o Cais v0.11 e renomear o módulo. O trabalho que veio depois foi corte por camada: Vite, Svelte, Inertia e a validação duplicada.",
				"A aposta: o handler chama view.Write, o Drive morfa #amarra-main, e sort e filtro são GET. O app não monta estado no cliente; ele pede a página de novo.",
				"A conta fecha com o chato gerado: 23 componentes no kit, g resource e g auth na CLI e um doctor que transforma bug em check (1.148 funções de teste em Go, 203 em JS, 206 issues).",
			}},
			{Type: "heading2", Text: "O que ficou e o que saiu"},
			{Type: "paragraph", Text: "O Cais já tinha o que eu não queria reescrever: router, sessão, CSRF, fila de jobs e SQLite num arquivo só, sem CGO. O Amarra manteve essa base e mudou a camada de cima. O README do projeto resume em três frases: o browser não monta uma SPA, handlers chamam view.Write, o Drive morfa #amarra-main. Não existe Vite, Svelte ou Inertia num app gerado."},
			{
				Type:        "quote",
				Text:        "The browser does not mount a SPA. Handlers call `view.Write`. Drive morphs `#amarra-main`. There is no Vite, Svelte, or Inertia in generated apps.",
				AuthorQuote: "README do amarra-cais",
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/README.md",
				SourceLabel: "README.md",
			},
			{Type: "paragraph", Text: "O preço estava claro desde o começo: migrar um app Cais é reescrita de UI, e isso está escrito na doc (docs/migrate-inertia.md). Não teve conversor automático. A escolha era entre manter dois runtimes para sempre e pagar a reescrita uma vez."},
			{Type: "heading2", Text: "A troca: um morph em vez de um runtime"},
			{Type: "paragraph", Text: "O Drive intercepta clique e submit para o mesmo host por padrão e substitui só o #amarra-main. O resto do HTML vem do servidor: o layout entrega o shell, a página preenche o bloco content. A doc é explícita sobre o que isso compra: primeira carga, curl e crawler recebem a página inteira; navegação vira morph. Não existe atributo para ligar o Drive, ele já vem ligado; quem precisa de visita completa marca data-amarra-skip."},
			{Type: "paragraph", Text: "No lugar de Alpine, Stimulus ou HTMX, o framework traz o kit (23 componentes: form, input, table, modal, tabs, filters, stat, tooltip) e hooks num único script de 117 KB. O app gerado não tem bundler. Sort e filtro são GET (?q=&sort=), o servidor re-renderiza, e o componente com estado (dropdown, dialog, tabs) é hook nativo, não plugin de terceiro. Streaming também cabe: amarra/stream escreve operações nomeadas (append, replace, morph, remove, toast) com alvo por id, e o hub live usa WebSocket para o mesmo contrato."},
			{Type: "heading2", Text: "O chato precisa ser gerado"},
			{Type: "paragraph", Text: "A CLI é Rails-style de propósito: new, install, dev, test, routes, db, doctor, pwa, link/unlink e os geradores (resource, auth, job, stream, component). O g resource emite o kit inteiro, não um CRUD de exemplo com classes soltas. O layout é previsível: os geradores marcam pontos no código (// cais:seeds, // cais:live-views) e o upgrade sabe onde mexer."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "Um handler inteiro: o layout cuida do shell, a página só preenche o content",
				Code: `view.Write(w, r, h.views, view.Page{
  Layout: "app",
  Name:   "dashboard",
  Data: map[string]any{
    "Title":     h.catalog.T("dashboard.title"),
    "Site":      meta.ForRequest(h.site, r),
    "CSRFToken": csrf.TokenFromRequest(r),
    "Flash":     flashMsg,
  },
}, h.cfg)`,
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/website/src/content/docs/docs/how-to/pages-and-views.md#L30-L39",
				SourceLabel: "website/src/content/docs/docs/how-to/pages-and-views.md:29",
			},
			{Type: "heading2", Text: "Cada bug virou um check"},
			{Type: "paragraph", Text: "O doctor existe porque os bugs aconteceram. Cada checagem nasceu de um incidente com número: o install sair 0 com CSS quebrado (#54), styles.css não existir depois do clone (#141), o bundle vendored ficar defasado do go.mod (#320). O comentário fica no código com o número da issue, para a próxima pessoa saber por que aquele erro é fatal."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "O install falha alto: sair 0 aqui deixou gente com app sem estilo (#54)",
				Code: `	// Build Tailwind so styles.css exists after clone/install (#141).
	if _, err := os.Stat(filepath.Join(dir, cssInput)); err == nil {
		_, _ = fmt.Fprintln(c.Out, "→ tailwind build (styles.css)")
		if err := runTailwindBuild(dir, false); err != nil {
			// Fail loudly: exiting 0 here left users on an unstyled app (#54).
			return fmt.Errorf("css build failed: %w — retry: amarra-cais css (NODE_ENV=production skips npm devDependencies)", err)
		}
	}`,
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/internal/cli/commands.go#L48-L55",
				SourceLabel: "internal/cli/commands.go:48",
			},
			{Type: "paragraph", Text: "A cultura de teste acompanha: 1.148 funções de teste em Go e 203 em JS para cerca de 27 mil linhas de código, com o código de teste passando de 29 mil. O AGENTS.md do repo abre dizendo para quem ele é escrito, e não é só para gente: o leitor principal costuma ser um agente, com rg barato, leitura completa caro e arquivo truncado em duas mil linhas. Nome grepável, módulo pequeno e teste headless viram requisito de manutenção quando metade das mudanças passa por agente."},
			{
				Type:        "quote",
				Text:        "Primary reader of this repo is often an LLM agent: cheap `rg`, expensive full reads, file truncation ~2k lines. Optimize for grep-unique names, small modules, and headless tests.",
				AuthorQuote: "AGENTS.md do amarra-cais",
				SourceURL:   "https://github.com/puppe1990/amarra-cais/blob/93c44c3/AGENTS.md#L3",
				SourceLabel: "AGENTS.md:3",
			},
			{Type: "heading2", Text: "O preço, com nome e número"},
			{Type: "paragraph", Text: "Nem tudo é ganho. O shell fora do #amarra-main não é atualizado pelo morph, e isso morde: neste site, trocar o idioma deixava a nav e o toggle no idioma anterior até eu abrir uma issue no próprio framework (#397). O JS do app é cache-first no service worker, então correção em web/static/js só chega com pwa --bump (#399). O data-amarra-skip é uma decisão por link: esquecer dele vira um bug silencioso."},
			{Type: "paragraph", Text: "Também não existe estado rico no cliente, animação coreografada entre telas nem mutação offline com fila própria. Quando o app precisa disso, o Amarra não é a ferramenta. Ele foi desenhado para o caso em que a tela é uma projeção do banco e o servidor é quem decide."},
			{Type: "heading2", Text: "O que os números dizem"},
			{Type: "paragraph", Text: "Dois meses depois do fork: 687 commits, 206 issues, v0.27.0, 20 projetos neste workspace dependendo do módulo (de um tradutor de EPUB a um ERP de locação), 23 componentes no kit e um script de 117 KB no browser. A escala é a minha: cada app com um SQLite, um servidor e um writer. O Amarra não foi desenhado para outra."},
			{Type: "paragraph", Text: "Este site é o exemplo mais direto. O Sintaxe nasceu SPA (React + Vite) e hoje é Amarra: o catálogo editorial é código Go, o agendamento é um published_at no futuro, a edição em inglês é uma coluna a mais. A migração está no histórico do repositório público, com o SPA antigo no commit anterior."},
			{
				Type:        "paragraph",
				Text:        "O framework está aberto em",
				SourceURL:   "https://github.com/puppe1990/amarra-cais",
				SourceLabel: "github.com/puppe1990/amarra-cais",
			},
		},
	}
}
