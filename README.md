# Sintaxe — Jornal de Tecnologia

O blog editorial [Sintaxe](https://sintaxe.apps.gestaobem.com) reconstruído com [Amarra-cais](https://github.com/puppe1990/amarra-cais): Go, HTML server-rendered (Amarra Views + Drive), Tailwind e SQLite. O catálogo editorial vive em `internal/db/content.go` — com o primeiro ensaio publicado e os próximos entrando conforme são escritos — e é publicado no boot por seeds idempotentes e reconciliadores.

## O que tem aqui

- **Capa** (`/`) com lead + destaques + grade, busca, cadernos com contagem e filtro por tempo de leitura.
- **Leitor** (`/artigo/{slug}`) com blocos ricos (takeaways numerados, quotes, código copiável, benchmark interativo), progresso de leitura, modo de foco, tamanho de fonte, aplausos, compartilhamento, comentários e contador de views.
- **Biblioteca pessoal** (`/salvos`) guardada no cookie `sintaxe_salvos` (HttpOnly), com toggle por artigo.
- **Manifesto** (`/sobre`), **página do colunista** (`/autor/{slug}`), **boletim** por e-mail (SQLite) e `sitemap.xml` / `rss.xml` / `robots.txt` sem cookies.
- **Open Graph**: capa em `web/static/og.png` (1200×630). Artigo compartilhado pinta o título em `GET /og/{slug}`. Spec visual em `design/og.html`. Regenerar a capa: `make og`.

## Edição em inglês (`/en`)

O Sintaxe é bilíngue por caminho: a raiz é pt-BR e `/en` é a edição inglesa. O toggle `PT / EN` no header aponta para o caminho equivalente **com o slug traduzido** — `/artigo/traducao-epubs-llms-fila-por-capitulo` ↔ `/en/artigo/epub-translation-llms-resilient-queue`.

- **Conteúdo**: o ensaio carrega a tradução nas colunas `slug_en`, `title_en`, `subtitle_en`, `lead_paragraph_en`, `tags_en`, `blocks_en` (migration `002_article_translations.sql`); role/bio do colunista em `role_en`/`bio_en`. Sem tradução, o `/en` cai na edição portuguesa (menos o slug, que é estrito: slug PT em `/en` responde 404).
- **Slugs são estritos por idioma**: `/artigo/{slug_pt}` não resolve em inglês e `/en/artigo/{slug_en}` não resolve na raiz.
- **Chrome**: todas as strings passam por `t` nos catálogos `internal/i18n/{en,pt}.go`; `/en/*` força o catálogo inglês em `pathLocale` (internal/app) antes do `i18n.LocaleMiddleware`. Datas e rótulos de caderno saem de `internal/editorial/localization.go`.
- **Identidade do leitor**: o cookie `sintaxe_salvos` guarda sempre o slug canônico (pt-BR), então a biblioteca funciona ao trocar de idioma.
- **SEO**: `hreflang` recíproco quando existe tradução, `<html lang>` por catálogo e `sitemap.xml` com as duas URLs.

Para publicar a tradução de um ensaio novo, preencha os campos `*EN` (inclusive `SlugEN`) no arquivo `internal/db/content_articles_*.go` — a edição inglesa entra no mesmo deploy.

## Highlight do código

Os blocos de código são coloridos no servidor com [chroma](https://github.com/alecthomas/chroma) (`internal/highlight`), emitindo classes (`.kd`, `.s`, `.c`…) e nenhum estilo inline; a paleta vive no `input.css` (keywords em cobre, strings em sálvia, números em âmbar, comentários apagados) e acompanha o tema claro/escuro. Sem JS de terceiros, sem flash de código sem cor e sem mexer na CSP.

## Citar código do repositório

Blocos de código e trechos citados apontam para o repo público: `SourceURL` (link `blob/<commit>/arquivo#Lx-Ly`) e `SourceLabel` (`internal/app/app.go:149`) no `models.ContentBlock`. Parágrafos com `SourceURL` ganham o link no fim (é assim que o ensaio mostra `github.com/puppe1990/epub-translator`), e a atribuição de uma citação leva o mesmo link depois do travessão.

## Agendar posts

Escreva o ensaio em `internal/db/content_articles_*.go` com `PublishedAt` no futuro — a visibilidade é por relógio: capa, página do autor, sitemap, RSS e a rota direta escondem o texto até o instante chegar. Como os seeds rodam no boot, o post **aparece sozinho no dia, sem novo deploy**.

```go
PublishedAt: time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC), // 06:00 BRT do dia 10
```

Para meia-noite de Brasília use hora `3` (03:00Z). Um ensaio agendado nunca conta nos "Ensaios Ativos" nem aceita comentário/aplauso antes da hora.

## Quick start

```bash
amarra-cais install    # npm install + go mod tidy (+ build do CSS)
amarra-cais dev        # http://localhost:8791 (LOCALE=pt no .env)
amarra-cais test       # go test ./...
make ci                # test + lint + format-check
```

## Testes

Um comando, headless (SQLite `:memory:`, fixtures com gofakeit em `internal/testdata`):

```bash
go test ./... -race -count=1
```

- `internal/editorial` (TDD): filtros de busca/categoria/tempo, segmentação em tiers, tempo de leitura, bookmarks, datas pt-BR.
- `internal/store` + `internal/db`: migrations, upserts idempotentes, seeds que preservam claps/views.
- `internal/handlers`: capa, leitor (422/303), salvos por cookie, sobre/autor, newsletter, sitemap/RSS/robots e 404 desenhado.

## Deploy

`cleat__deploy app=matheuspuppe-blog ref=main` — o Cleat instala Go + Node, roda `npm ci` (build do Tailwind) e `go build ./cmd/server`. O `.cleat_deploy/deploy.json` fixa `runtime: golang`, `binaries: [server]`, `memory_max_mb: 256`.

## Estrutura

| Path                  | Responsabilidade                                       |
| --------------------- | ------------------------------------------------------ |
| `cmd/server/`         | Entry point (bootstrap + seeds + devreload)            |
| `internal/app/`       | Bootstrap, `registerRoutes`, 404 desenhado             |
| `internal/handlers/`  | HTTP handlers (`view.Write`) + view-models             |
| `internal/editorial/` | Regras editoriais puras (filtros, tiers, formatação)   |
| `internal/ogcard/`    | PNG 1200×630 de compartilhamento (capa e `/og/{slug}`) |

| `internal/models/` | Row structs (SQL) e payloads |
| `internal/store/` | SQLite + migrations |
| `internal/db/` | Conteúdo publicado + `RunSeeds` idempotente |
| `internal/testdata/` | Fixtures gofakeit |
| `web/templates/` | Layout, partials e páginas Amarra |
| `web/static/` | CSS (Tailwind), JS (`amarra.js`, `sintaxe.js`), PWA e imagens |
