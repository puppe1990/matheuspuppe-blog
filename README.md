# Sintaxe — Jornal de Tecnologia

O blog editorial [Sintaxe](https://matheuspuppe-blog.apps.gestaobem.com) reconstruído com [Amarra-cais](https://github.com/puppe1990/amarra-cais): Go, HTML server-rendered (Amarra Views + Drive), Tailwind e SQLite. Os 6 ensaios, 6 colunistas e comentários vivem em `internal/db/content*.go` e são publicados no boot por seeds idempotentes.

## O que tem aqui

- **Capa** (`/`) com lead + destaques + grade, busca, cadernos com contagem e filtro por tempo de leitura.
- **Leitor** (`/artigo/{slug}`) com blocos ricos (takeaways numerados, quotes, código copiável, benchmark interativo), progresso de leitura, modo de foco, tamanho de fonte, aplausos, compartilhamento, comentários e contador de views.
- **Biblioteca pessoal** (`/salvos`) guardada no cookie `sintaxe_salvos` (HttpOnly), com toggle por artigo.
- **Playground** (`/laboratorio`): simulador de runtimes (Rust/Go/Node/Python) — a matemática é server-side (`internal/benchmark`), os sliders só atualizam rótulos e disparam um GET.
- **Manifesto** (`/sobre`), **página do colunista** (`/autor/{slug}`), **boletim** por e-mail (SQLite) e `sitemap.xml` / `rss.xml` / `robots.txt` sem cookies.

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
- `internal/benchmark` (TDD): as 4 curvas do simulador com valores conferidos à mão.
- `internal/store` + `internal/db`: migrations, upserts idempotentes, seeds que preservam claps/views.
- `internal/handlers`: capa, leitor (422/303), salvos por cookie, laboratório, sobre/autor, newsletter, sitemap/RSS/robots e 404 desenhado.

## Deploy

`cleat__deploy app=matheuspuppe-blog ref=main` — o Cleat instala Go + Node, roda `npm ci` (build do Tailwind) e `go build ./cmd/server`. O `.cleat_deploy/deploy.json` fixa `runtime: golang`, `binaries: [server]`, `memory_max_mb: 256`.

## Estrutura

| Path                  | Responsabilidade                                              |
| --------------------- | ------------------------------------------------------------- |
| `cmd/server/`         | Entry point (bootstrap + seeds + devreload)                   |
| `internal/app/`       | Bootstrap, `registerRoutes`, 404 desenhado                    |
| `internal/handlers/`  | HTTP handlers (`view.Write`) + view-models                    |
| `internal/editorial/` | Regras editoriais puras (filtros, tiers, formatação)          |
| `internal/benchmark/` | Curvas do simulador de runtimes                               |
| `internal/models/`    | Row structs (SQL) e payloads                                  |
| `internal/store/`     | SQLite + migrations                                           |
| `internal/db/`        | Conteúdo publicado + `RunSeeds` idempotente                   |
| `internal/testdata/`  | Fixtures gofakeit                                             |
| `web/templates/`      | Layout, partials e páginas Amarra                             |
| `web/static/`         | CSS (Tailwind), JS (`amarra.js`, `sintaxe.js`), PWA e imagens |
