package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// networkListener é o segundo ensaio: como a extensão que grava o tráfego do
// navegador alimenta um agente, com escopo, evidência e replay. Os trechos são
// do repositório público, citados no commit.
func networkListener() models.Article {
	en := networkListenerEn()
	return models.Article{
		Slug:            "captura-de-trafego-do-browser-para-agentes",
		Title:           "Capturando tráfego real do browser para alimentar um agente: a extensão Network Listener",
		Subtitle:        "O agente adivinha menos quando lê o que a página já falou com a API. Como a captura, o escopo, o achado com evidência e o replay se encaixam numa extensão MV3 com SQLite.",
		LeadParagraph:   "Antes de escrever o primeiro cliente de uma API, tem uma fonte de verdade que quase ninguém consulta: a própria aba Network. É lá que estão os headers que a documentação esqueceu, o corpo exato do POST, o erro que o front recebeu. O Network Listener nasceu disso: uma extensão do Chrome que grava esse tráfego num SQLite local, com filtros, marcas de sessão, achados e replay, e um servidor MCP que deixa um agente consultar tudo por id de captura.",
		Category:        "Segurança & Criptografia",
		ReadTimeMinutes: 7,
		PublishedAt:     time.Date(2026, 11, 3, 18, 0, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_ai_agents_lead.jpg",
		Author:          authorRef("matheus-puppe"),
		Tags:            []string{"HTTP", "Extensões", "Segurança", "SQLite", "Agentes"},
		Tier:            2,
		SlugEN:          en.Slug,
		TitleEN:         en.Title,
		SubtitleEN:      en.Subtitle,
		LeadParagraphEN: en.Lead,
		TagsEN:          en.Tags,
		BlocksEN:        en.Blocks,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Metadados saem do webRequest; os corpos vêm de um hook que remenda fetch e XMLHttpRequest na própria página, porque uma extensão MV3 não alcança corpo de requisição por outro caminho.",
				"O escopo é conferido no servidor a cada replay e a cada salto de redirect: sem workspace ativo a resposta é 409, fora dos hosts do workspace é 403.",
				"O que alimenta o agente é o inventário de endpoints (/users/{int}), a marca de sessão e o achado com evidência apontando para capturas; o dump cru fica para o export.",
			}},
			{Type: "heading2", Text: "Por que capturar no browser"},
			{Type: "paragraph", Text: "Documentação de API envelhece no dia do deploy, e o DevTools perde tudo quando a aba fecha. O que a extensão grava é o que a página realmente pediu: url, método, tipo de recurso, status, duração, cache, headers dos dois lados, corpo enviado e corpo recebido, além do tamanho do que foi gravado. Tráfego que não passa pelo browser (curl, outro navegador, app nativo) entra pelo import de HAR no dashboard."},
			{
				Type:        "quote",
				Text:        "Um proxy MITM ampliaria a captura para curl, outros browsers e apps nativos. O preço é certificado local, confiança do sistema, HTTP/2, WebSocket, QUIC e uma superfície de segurança que o MVP não pede.",
				AuthorQuote: "docs/local-proxy-decision.md",
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/docs/local-proxy-decision.md",
				SourceLabel: "docs/local-proxy-decision.md",
			},
			{Type: "paragraph", Text: "O proxy MITM existe só no app desktop, em loopback, restrito aos hosts de um programa ativo e ligado por um diálogo explícito no dashboard; a produção nunca sobe esse listener. A extensão segue como caminho principal, e a decisão de manter o proxy fora do servidor está escrita no próprio repositório."},
			{Type: "heading2", Text: "Metadados pelo webRequest, corpos pelo remendo"},
			{Type: "paragraph", Text: "A extensão escuta webRequest para o que o navegador já sabe: url, método, tipo de recurso, status, headers, timing. Corpo é outra história: a API de webRequest não entrega o que foi enviado nem o que voltou, então o jeito é remendar fetch e XMLHttpRequest dentro da página, num hook que roda no MAIN world e manda relatórios por window.postMessage para a ponte isolada. Todo patch vive em try/catch, porque uma página quebrada por telemetria é pior do que uma captura a menos."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "O hook da página: host fora do escopo nem lê corpo, e o postMessage nunca quebra a página",
				Code: `/** Hosts the user left out never pay the cost of reading a body. */
function shouldReadBody(url: string): boolean {
  return bodyCaptureEnabled() && bodyAllowedFor(getHost(url));
}

function postReport(draft: ReportDraft & Partial<PendingRequest>): void {
  if (!bodyCaptureEnabled()) return;

  try {
    window.postMessage({ channel: BODY_CHANNEL, report: draft }, '*');
  } catch {
    // the page is not ours to break
  }
}`,
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/apps/extension/src/content/body-hook.ts#L33-L51",
				SourceLabel: "apps/extension/src/content/body-hook.ts:33",
			},
			{Type: "paragraph", Text: "Três limites aparecem no caminho. Só conteúdo textual é gravado; imagem, vídeo e blob ficam de fora. Cada corpo é cortado num teto configurável, 128 kB por padrão e 150 kB no máximo, porque o JSON de um reel do InnerTube tem uns 45 kB comprimido e os 8 kB iniciais guardavam só um prefixo dele. E há uma armadilha do Chrome no meio: o InnerTube monta um Request antes de chamar fetch, então init.body chega vazio, e o hook também remenda o construtor de Request para guardar o corpo ali, antes que algum wrapper consuma o stream."},
			{Type: "heading2", Text: "Lote, fila e o servidor que cai"},
			{Type: "paragraph", Text: "A extensão envia em lote a cada dois segundos, com dois tetos por envio: quantidade de registros e 2 MB, porque um corpo grande demais estoura o limite de requisição do servidor. Quando o servidor não responde, a fila fica na máquina (espelhada em chrome.storage.session) e o intervalo cresce em backoff exponencial até 60 segundos: min(60 s, intervalo × 2 ^ falhas). O projeto decidiu não ter circuit breaker; o dashboard mostra o erro e a captura continua gravando local até o envio dar certo."},
			{Type: "paragraph", Text: "Do outro lado, Fastify com better-sqlite3, escutando em 127.0.0.1 por padrão. A retenção apaga o que passou de três dias (configurável de 1 a 90) numa varredura de hora em hora, em lotes de 5.000 registros com yield entre eles. O better-sqlite3 é síncrono, e um DELETE único de um backlog grande travaria o event loop do servidor."},
			{Type: "heading2", Text: "Escopo, recusa e replay"},
			{Type: "paragraph", Text: "Cada captura pertence a um workspace, com lista de hosts permitidos e lista de excluídos. O replay (um POST por captura, com edições de método, url, headers e corpo) recusa antes de sair: sem workspace ativo responde 409, url fora dos hosts do workspace responde 403, url inválida responde 400. O redirect é seguido com redirect manual, no máximo cinco saltos, e cada salto passa de novo pelo teste de escopo; se o Location aponta para fora, a cadeia para e a resposta volta marcada como bloqueada."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "Replay com escopo por salto: o Location é reavaliado antes de seguir",
				Code: `for (let hop = 0; hop < REPLAY_MAX_HOPS; hop += 1) {
  const init = hopInit(state.method, headers, hop === 0 ? state.body : undefined);
  const response = await state.fetchImpl(url, init);
  last = { url, status: response.status, body: await readBody(response), responseHeaders: pickedHeaders(response) };

  const location = response.header('location');
  if (!isRedirect(response.status) || !location) return { ...last, blocked };

  const next = planReplay(new URL(location, url).href, state.workspace);
  if (!next.ok) return { ...last, blocked: next.reason };

  headers = headersForHop(headers, host, next.host);
  url = next.url.href;
  host = next.host;
  blocked = undefined;
}`,
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/apps/server/src/lib/replay.ts#L119-L134",
				SourceLabel: "apps/server/src/lib/replay.ts:119",
			},
			{Type: "paragraph", Text: "Duas regras completam o replay. Só um roda por vez; o segundo pedido leva 429 com a mensagem \"espere o replay em andamento terminar\". E GET e HEAD nunca levam corpo, mesmo que a captura original tenha um. O replay também não é um evento à parte no histórico: ele entra como uma captura nova, marcada com source replay e com o id da original em replayOfId, então a resposta que o servidor devolveu fica ao lado do que o browser tinha feito."},
			{Type: "heading2", Text: "Do dump para o achado"},
			{Type: "paragraph", Text: "Um dump de trinta minutos tem milhares de linhas de videoplayback, e ninguém acha nada ali. O recorte começa pela marca de sessão: o botão de nova sessão grava o instante em que você começa a olhar (os Shorts, o login, o checkout) e o export oferece a sessão inteira como filtro. Depois vem o inventário de endpoints, que agrupa UUID, inteiro e id opaco num padrão: /api/user/1 e /api/user/2 viram /api/user/{int}, com contagem, status vistos, content type e o nome da operação quando o tráfego é GraphQL."},
			{Type: "paragraph", Text: "O achado tem severidade de info a critical e uma lista de evidências. Evidência é referência: o achado aponta para o id da captura e nunca copia o corpo dela, então o relatório lido daqui a um mês ainda bate com o banco. O relatório sai em Markdown, pronto para colar no ticket, e o agente usa as mesmas rotas por MCP: listar requisições, inventário de endpoints, buscar captura, criar achado, anexar evidência, gerar relatório e disparar replay."},
			{Type: "paragraph", Text: "Sobre segredo, a decisão é explícita e está comentada no código: o corpo é gravado como foi enviado, porque o token no payload costuma ser exatamente o que você está depurando. O mascaramento é uma opção e roda na captura, por chave: password, token, authorization, api_key, cartão e cvv viram *** até quatro níveis de profundidade."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "O mascaramento opcional: chaves de credencial viram *** no JSON",
				Code: `/** Deep-ish redaction so credentials in JSON bodies never reach the database. */
export function redactSecrets(value: unknown, depth = 0): unknown {
  if (depth > MAX_REDACT_DEPTH || value === null || typeof value !== 'object') return value;
  if (Array.isArray(value)) return value.map((entry) => redactSecrets(entry, depth + 1));

  const result: Record<string, unknown> = {};
  for (const [key, entry] of Object.entries(value as Record<string, unknown>)) {
    result[key] = REDACTED_KEYS.test(key) ? '***' : redactSecrets(entry, depth + 1);
  }

  return result;
}`,
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/apps/extension/src/lib/body-text.ts#L30-L45",
				SourceLabel: "apps/extension/src/lib/body-text.ts:30",
			},
			{
				Type:        "quote",
				Text:        "Off by default: tokens in payloads are usually exactly what you are debugging.",
				AuthorQuote: "Comentário em apps/extension/src/lib/body-text.ts",
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/apps/extension/src/lib/body-text.ts#L61",
				SourceLabel: "apps/extension/src/lib/body-text.ts:61",
			},
			{Type: "paragraph", Text: "O que protege o banco depois de gravado é cifra em repouso, com chave vinda do ambiente ou do cofre do macOS; sem chave, a cifra fica desligada e o dashboard mostra o estado. Ligar a cifra passa por um fluxo de confirmação (pendente até confirmar ou abortar), e a busca dentro de corpos deixa de olhar o JSON no disco quando ela está ativa. No deploy, a extensão usa token de máquina para gravar e o dashboard usa login com e-mail e senha, com limite de oito falhas por IP e e-mail a cada quinze minutos respondendo 429 com Retry-After."},
			{Type: "heading2", Text: "O que eu levaria para o próximo serviço"},
			{Type: "paragraph", Text: "O escopo é conferido no servidor, inclusive a cada salto de redirect; a extensão é conveniência, não garantia. Toda recusa tem nome e status (409 sem workspace, 403 fora do escopo, 400 url inválida, 429 replay ocupado), e é isso que permite ao agente explicar ao humano o que faltou em vez de tentar de novo às cegas."},
			{Type: "paragraph", Text: "Mutações ficam serializadas: um replay por vez evita duas escritas na mesma sessão e é fácil de auditar depois. A retenção em lotes de 5.000 com yield entre eles é requisito, não otimização, quando o driver do banco é síncrono. E guardar o payload como ele foi enviado, com mascaramento opcional e revelação explícita na leitura, resolveu mais casos do que a política contrária teria resolvido."},
			{
				Type:        "paragraph",
				Text:        "O Network Listener roda em produção atrás de token e login, e o código está aberto em",
				SourceURL:   "https://github.com/puppe1990/network-listener",
				SourceLabel: "github.com/puppe1990/network-listener",
			},
		},
	}
}
