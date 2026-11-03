# Sintaxe (matheuspuppe-blog) — AI Conventions

Jornal de tecnologia da Sintaxe em Amarra-cais. Host `sintaxe.apps.gestaobem.com` (Cleat id 239, slug `matheuspuppe-blog`, porta 4229). Repo `puppe1990/matheuspuppe-blog` (público). O catálogo editorial vive em `internal/db/content.go` (colunista único, Matheus Puppe; o primeiro ensaio já está publicado, os próximos entram conforme são escritos); fixtures de teste vêm de `internal/testdata` (gofakeit, seed fixo) — esse pacote não importa `internal/db`.

Primary reader is often an LLM agent. Prefer small greps, small modules, and headless tests.

## Rule #1: TDD is mandatory

Before writing production code:

1. Write the test (`*_test.go`)
2. Run focused: `go test ./... -v -run TestName`
3. Confirm it **fails** for the right reason
4. Write the **minimal** code to pass
5. Run `amarra-cais test` (ou `go test ./... -race -count=1`)

## Clean code for agents

| Priority | Rule                                                                                                 |
| -------- | ---------------------------------------------------------------------------------------------------- |
| 1        | **Small units** — funções ~4–20 linhas; arquivos alvo 200–300, teto ~500                             |
| 2        | **SRP** — uma razão de mudança por pacote                                                            |
| 3        | **Greppable names** — substantivos de domínio (`editorial`, `salvos`, `autor`); evite `utils`        |
| 4        | **Comments = WHY** — Drive vs full HTML, cookie, Tailwind purge, seeds idempotentes; sem narrar WHAT |
| 5        | **Inject deps** — handlers recebem `Store`, `*view.Renderer`, `cais.Config` via construtor           |
| 6        | **Early returns** — máximo ~2 níveis de indentação                                                   |
| 7        | **Errors with values** — `fmt.Errorf("...: %w", err)`; `store.ErrNotFound` para 404                  |
| 8        | **Headless tests** — SQLite `:memory:`; sem seed manual nos testes (use `db.RunSeeds`)               |

## Layout

| Path                  | Responsabilidade                                                        |
| --------------------- | ----------------------------------------------------------------------- |
| `internal/editorial/` | Filtros, segmentação em tiers, tempo de leitura, bookmarks, datas pt-BR |
| `internal/ogcard/`    | PNG 1200×630 de compartilhamento (capa e `/og/{slug}`)                  |

| `internal/db/` | Conteúdo real + `RunSeeds` idempotente e reconciliador (roda no boot; não zera claps/views; remove linhas fora dos seeds) |
| `internal/store/` | SQLite + migrations (`authors`, `articles`, `comments`, `subscribers`) |
| `internal/handlers/` | Handlers HTML; `viewdata.go` tem `ArticleView`, cookie de salvos e `next` |
| `web/templates/` | `layouts/app.html`, `partials/*` (flat only), `pages/*` |

Markers de geradores (não remover): `registerRoutes`, `Close() error`, `// cais:seeds`, `// cais:live-views`.

## Convenções deste app

- **Sem auth**: tudo público; o único estado pessoal é o cookie `sintaxe_salvos` (HttpOnly, cap de 100 slugs em `editorial`). `next` de todo form passa por `sanitizeNext` (mesmo-site absoluto apenas).
- **Agendamento**: `published_at` no futuro fica invisível (capa, página do autor, sitemap, RSS e rota direta → 404) até o instante — o corte é `published_at <= now` no store (`publishedNow`). Como os seeds rodam no boot, o post **entra sozinho no dia, sem novo deploy**. Hora em UTC: meia-noite de Brasília = `03:00Z`; manhã editorial = `09:00Z` (06:00 BRT).
- **Bilíngue por caminho**: raiz é pt-BR, `/en` é a edição inglesa. `pathLocale` (internal/app) injeta `lang=en` antes do `i18n.LocaleMiddleware`; `localeOf` (handlers) devolve `en`/`pt`. Conteúdo traduzido mora nas colunas `*_en` do artigo (`SlugEN` incluso) e em `role_en`/`bio_en` do autor (`editorial.Localize` faz o swap, com fallback para o português). Slugs são **estritos por idioma** — `/en/artigo/{slug_pt}` e `/artigo/{slug_en}` respondem 404. Toda string de chrome usa `t` com os catálogos `internal/i18n/{en,pt}.go`; datas e cadernos saem de `editorial/localization.go` (`FormatDateIn`, `CategoryLabel`). O cookie `sintaxe_salvos` guarda o slug canônico (pt-BR) para a biblioteca sobreviver à troca de idioma; o toggle `PT/EN` do layout aponta para o caminho equivalente com o slug traduzido.
- **Citações de código**: bloco de código (ou comentário) carrega `SourceURL`/`SourceLabel` em `models.ContentBlock`; o artigo renderiza `internal/…:linha` como link para o repo público, fixado no commit citado (`blob/<sha>/…`). Mudou o trecho no repo? Atualize o sha e o label do bloco.
- **Catálogo = código**: `RunSeeds` reconcilia — artigo/autor fora de `internal/db/content.go` é removido do banco (comentários do artigo removido vão junto, em FK order). Assinantes do boletim são preservados.
- **Templates**: Amarra Views (`view.Write`), kit `<.form>`/`<.input>`/`<.button>`/`<.flash />`. Só **partials flat** carregam; `pages/*/*.html` teria um nível de nesting (não usado aqui). Sem funções custom de template — datas/percentuais são pré-computados nos view-models.
- **Drive**: links/forms são interceptados; POSTs redirecionam 303 ("flash + redirect"). GET de form serializa os campos na query (`?q=&categoria=&tempo=`). Não checar `HX-Request`.
- **CSS/Tailwind**: classes JS-only (`amarra-focus`, `amarra-toast`) e as cores dos runtimes (safelist) não aparecem no HTML escaneado — mantenha no `safelist` do `tailwind.config.js` ou o Tailwind purga a regra. A paleta vem das vars `--amarra-*` em `input.css` (`:root` escuro, `html.light` claro).
- **JS**: `web/static/js/sintaxe.js` é o único script do app; listeners delegados no `document`, re-init de estado no evento `amarra:morphed`.
- **SEO**: `sitemap.xml`, `rss.xml` e `robots.txt` respondem **sem** `Set-Cookie` (#342). A capa usa `web/static/og.png` (1200×630, `make og`). Artigo aponta `og:image` para `GET /og/{slug}` — PNG pintado com o título, sem cookie. Regenerar a capa: `go run ./cmd/genog`.

## Comandos

```bash
amarra-cais dev            # air + tailwind watch (templates recarregam in-process)
go test ./... -race -count=1
amarra-cais routes
amarra-cais db seed        # idempotente; também roda no boot do servidor
make og                    # regenera web/static/og.png
make ci                    # test + lint + format-check
```

## Do not

- Editar templates por request (view.Load uma vez no boot)
- CSS inline (use Tailwind; exceção documentada: larguras JS via `data-bar-width`)
- Mockar o banco (SQLite `:memory:`)
- Passar de ~500 linhas por arquivo sem dividir
- Publicar feature sem teste headless
