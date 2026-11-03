package db

import "github.com/puppe1990/matheuspuppe-blog/internal/models"

// networkListenerEn é a edição inglesa do ensaio sobre o Network Listener. Os
// trechos de código são os do repositório público, citados no commit.
func networkListenerEn() englishEdition {
	return englishEdition{
		Slug:     "browser-traffic-capture-for-agents",
		Title:    "Capturing real browser traffic to feed an agent: the Network Listener extension",
		Subtitle: "An agent guesses less when it reads what the page already told the API. How capture, scope, evidence-backed findings and replay fit together in an MV3 extension with SQLite.",
		Lead:     "Before writing the first client for an API, there is a source of truth almost nobody consults: the Network tab itself. That is where the headers documentation forgot live, along with the exact POST body and the error the frontend received. Network Listener came from that: a Chrome extension that records this traffic into a local SQLite file, with filters, session marks, findings and replay, plus an MCP server that lets an agent query everything by capture id.",
		Tags:     []string{"HTTP", "Browser Extensions", "Security", "SQLite", "Agents"},
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Metadata comes from webRequest; bodies come from a hook that patches fetch and XMLHttpRequest inside the page, because an MV3 extension has no other way to reach a request body.",
				"Scope is enforced on the server on every replay and on every redirect hop: no active workspace answers 409, hosts outside the workspace answer 403.",
				"What feeds the agent is the endpoint inventory (/users/{int}), the session mark and the finding whose evidence points at captures; the raw dump stays for the export.",
			}},
			{Type: "heading2", Text: "Why capture in the browser"},
			{Type: "paragraph", Text: "API documentation ages on deploy day, and DevTools forgets everything when the tab closes. What the extension records is what the page actually asked for: url, method, resource type, status, duration, cache, headers on both sides, request body, response body, and the size of what was stored. Traffic that never touches the browser (curl, another browser, a native app) comes in through HAR import in the dashboard."},
			{
				Type:        "quote",
				Text:        "Um proxy MITM ampliaria a captura para curl, outros browsers e apps nativos. O preço é certificado local, confiança do sistema, HTTP/2, WebSocket, QUIC e uma superfície de segurança que o MVP não pede.",
				AuthorQuote: "docs/local-proxy-decision.md",
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/docs/local-proxy-decision.md",
				SourceLabel: "docs/local-proxy-decision.md",
			},
			{Type: "paragraph", Text: "The MITM proxy exists only in the desktop app, on loopback, restricted to the hosts of one active program and enabled through an explicit dialog in the dashboard; production never starts that listener. The extension remains the main path, and the decision to keep the proxy out of the server is written down in the repository itself."},
			{Type: "heading2", Text: "Metadata from webRequest, bodies from a patch"},
			{Type: "paragraph", Text: "The extension listens on webRequest for what the browser already knows: url, method, resource type, status, headers, timing. Bodies are another story: the webRequest API delivers neither what was sent nor what came back, so the way in is patching fetch and XMLHttpRequest inside the page, in a hook that runs in the MAIN world and posts reports through window.postMessage to the isolated bridge. Every patch sits inside try/catch, because a page broken by telemetry is worse than one capture fewer."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "The page hook: hosts outside scope never read a body, and postMessage never breaks the page",
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
			{Type: "paragraph", Text: "Three limits show up along the way. Only textual content is stored; images, video and blobs stay out. Each body is cut at a configurable ceiling, 128 kB by default and 150 kB at most, because an InnerTube reel JSON is around 45 kB compressed and the first 8 kB only kept a prefix of it. And there is a Chrome trap in the middle: InnerTube builds a Request before calling fetch, so init.body arrives empty, and the hook also patches the Request constructor to keep the body there before some wrapper consumes the stream."},
			{Type: "heading2", Text: "Batch, queue and the server that goes down"},
			{Type: "paragraph", Text: "The extension sends batches every two seconds, with two ceilings per send: record count and 2 MB, because an oversized body would blow the server request limit. When the server does not answer, the queue stays on the machine (mirrored in chrome.storage.session) and the interval grows with exponential backoff up to 60 seconds: min(60 s, interval × 2 ^ failures). The project decided against a circuit breaker; the dashboard shows the error and capture keeps writing locally until a send succeeds."},
			{Type: "paragraph", Text: "On the other side, Fastify with better-sqlite3, listening on 127.0.0.1 by default. Retention deletes anything older than three days (configurable from 1 to 90) in an hourly sweep, in batches of 5.000 rows with a yield between them. better-sqlite3 is synchronous, and a single DELETE over a large backlog would block the server event loop."},
			{Type: "heading2", Text: "Scope, refusals and replay"},
			{Type: "paragraph", Text: "Every capture belongs to a workspace, with an allowed-host list and an excluded-host list. Replay (one POST per capture, with edits to method, url, headers and body) refuses before leaving: no active workspace answers 409, a url outside the workspace hosts answers 403, an invalid url answers 400. Redirects are followed with redirect manual, at most five hops, and each hop goes through the same scope check again; if the Location points outside, the chain stops and the response comes back flagged as blocked."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "Replay with scope per hop: the Location is re-evaluated before following",
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
			{Type: "paragraph", Text: "Two rules complete the replay. Only one runs at a time; a second request gets 429 with the message \"espere o replay em andamento terminar\". And GET and HEAD never carry a body, even when the original capture has one. A replay is also not a separate event in the history: it lands as a new capture, flagged with source replay and the original id in replayOfId, so the response the server gave sits next to what the browser had done."},
			{Type: "heading2", Text: "From dump to finding"},
			{Type: "paragraph", Text: "A thirty-minute dump holds thousands of videoplayback lines, and nobody finds anything in there. Narrowing starts with the session mark: the new-session button records the instant you start looking (the Shorts, the login, the checkout) and the export offers that session as a filter. Then comes the endpoint inventory, which groups UUIDs, integers and opaque ids into a pattern: /api/user/1 and /api/user/2 become /api/user/{int}, with a count, seen statuses, content type and the operation name when the traffic is GraphQL."},
			{Type: "paragraph", Text: "A finding carries a severity from info to critical and a list of evidence. Evidence is a reference: the finding points at the capture id and never copies its body, so the report read a month later still matches the database. The report comes out as Markdown, ready to paste into a ticket, and the agent uses the same routes over MCP: list requests, endpoint inventory, fetch a capture, create a finding, attach evidence, render the report and fire a replay."},
			{Type: "paragraph", Text: "On secrets the decision is explicit and commented in the code: the body is stored as it was sent, because the token in the payload is usually exactly what you are debugging. Masking is an option and runs at capture time, by key: password, token, authorization, api_key, card and cvv become *** down to four levels of depth."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "Optional masking: credential keys become *** in the JSON",
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
				AuthorQuote: "Comment in apps/extension/src/lib/body-text.ts",
				SourceURL:   "https://github.com/puppe1990/network-listener/blob/7c0e4ef/apps/extension/src/lib/body-text.ts#L61",
				SourceLabel: "apps/extension/src/lib/body-text.ts:61",
			},
			{Type: "paragraph", Text: "What protects the database once the row is written is encryption at rest, with a key coming from the environment or the macOS keychain; without a key, encryption stays off and the dashboard shows that state. Turning it on goes through a confirmation flow (pending until confirmed or aborted), and in-body search stops looking at the JSON on disk while it is active. In the deployment, the extension writes with a machine token and the dashboard uses email-and-password login, throttled at eight failures per IP and email every fifteen minutes, answering 429 with Retry-After."},
			{Type: "heading2", Text: "What I would carry into the next service"},
			{Type: "paragraph", Text: "Scope is enforced on the server, including on every redirect hop; the extension is convenience, not a guarantee. Every refusal has a name and a status (409 without a workspace, 403 out of scope, 400 invalid url, 429 replay busy), and that is what lets an agent explain to a human what was missing instead of retrying blind."},
			{Type: "paragraph", Text: "Mutations stay serialized: one replay at a time avoids two writes landing on the same session and is easy to audit later. Retention in batches of 5.000 with a yield between them is a requirement, not an optimization, when the database driver is synchronous. And storing the payload as it was sent, with optional masking and explicit reveal on read, solved more cases than the opposite policy would have."},
			{
				Type:        "paragraph",
				Text:        "Network Listener runs in production behind a token and a login, and the code is open at",
				SourceURL:   "https://github.com/puppe1990/network-listener",
				SourceLabel: "github.com/puppe1990/network-listener",
			},
		},
	}
}
