package db

import "github.com/puppe1990/matheuspuppe-blog/internal/models"

// englishEdition carries the translated catalog entry for an essay.
type englishEdition struct {
	Slug     string
	Title    string
	Subtitle string
	Lead     string
	Tags     []string
	Blocks   []models.ContentBlock
}

// primeiroEnsaioEn is the English edition of the first essay, slug included.
func primeiroEnsaioEn() englishEdition {
	return englishEdition{
		Slug:     "epub-translation-llms-resilient-queue",
		Title:    "Translating EPUBs with LLMs: a resilient queue, per-chapter retry and cheap models",
		Subtitle: "A deploy lands mid-translation, the model drops a tag and the bill grows. Running a book translator taught me how to design the queue, the cancellation and the cost per chapter.",
		Lead:     "A book is not a request. It is hundreds of XHTML documents, each taking minutes of inference, and a server that restarts in the middle of the work. The first version of my translator treated the book as a single job: a deploy in the middle of the afternoon was enough to lose hours of translation. The design that survived lives in a small application, with SQLite and one file per book, and solves the three problems every LLM service has: interruption, validation and cost.",
		Tags:     []string{"LLM", "EPUB", "Queues", "Reliability", "Cost"},
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"The chapter is the atomic unit: retry, pause, cost and restart recovery live in it, not in the book.",
				"In production the deploy lands mid-flight: SIGTERM needs a drain that returns the running chapter to the queue and leaves the book marked as translating for the next boot.",
				"The model is the least reliable part of the pipeline: validate the response by contract — tags, ratio, language — before spending the next attempt.",
			}},
			{Type: "heading2", Text: "A book is not a request"},
			{Type: "paragraph", Text: "An EPUB is a zip with an ordered list of XHTML documents. The translator opens the package, reads the spine and turns each document into a chapter: one row in the chapters table, with a unique path inside the book, the source text, a cl100k token count and the execution state. That decision changes everything that follows. Retry stops being \"try the book again\" and becomes \"try chapter 14\"; a failing document does not bring down the other 199; and the cost of each call becomes observable per unit, because tokens are measured before any request."},
			{Type: "paragraph", Text: "The temptation to send the whole book in a single call is strong: less code, one stream, one result. In practice the model's context becomes the product's limit, retry gets expensive, and any failure in the middle destroys the entire translation. Translating per document is more tedious to build and infinitely cheaper to operate."},
			{Type: "heading2", Text: "Three states, one queue and a pause that is not a failure"},
			{Type: "paragraph", Text: "A chapter moves through four states: pending, running, done and failed. The queue is simple on purpose: pending and failed chapters, in document order. The runner keeps one goroutine per book, so two chapters of the same book never compete for the same file or the same database row. The retry button on a row takes another route: the chosen chapter jumps the queue and becomes the next one after the current work, without restarting the whole book."},
			{Type: "paragraph", Text: "Failure is a first-class state. A chapter that blew through the context, came back with the wrong tags or hit a provider outage ends as failed, with the error written on the row. The book is only marked failed if some chapter stayed that way. This keeps the interface honest: the table says exactly what happened, and retry points at the right problem."},
			{Type: "paragraph", Text: "Pausing is a different story, and the distinction matters. Pause cancels the in-flight call, waits for the execution to finish and returns only the running chapter to the queue. The book becomes paused, finished chapters stay intact, and the interrupted chapter is not counted as a failure. When the reader resumes, the queue continues where it stopped."},
			{Type: "heading2", Text: "The deploy in the middle of a chapter"},
			{Type: "paragraph", Text: "The first bug showed up on a routine deploy. systemd sends SIGTERM and waits fifteen seconds; the HTTP server, in turn, allowed five seconds for connections to close. Except the book page keeps an SSE stream open the whole time to show progress: that stream never ends on its own, the shutdown blew through the deadline, the process exited with an error and killed the translation request in flight. The chapter stayed marked as running forever, and the interface lied until someone restarted the server again. In the week of the incident, a pause clicked after deploy 1887 left the book stuck for hours, because the paused state was written over work nobody had paused."},
			{Type: "paragraph", Text: "The fix has two parts. First, the runner needs a drain that recognises shutdown: cancel each book's context, wait up to two seconds for the work to finish and return chapters stuck in running to the queue. Second, shutdown order starts to matter: drain first, then the HTTP server, and SQLite last. Without that order the drain writes to a database that is already closed."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "Shutdown order: drain the queue before closing the database",
				Code: `case <-ctx.Done():
    // SIGTERM: hand the running chapter back to the queue.
    if a.runner != nil {
        _ = a.runner.Drain(2 * time.Second)
    }
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := a.server.Shutdown(shutdownCtx); err != nil {
        _ = a.server.Close()
    }
    _ = a.store.Close()
    return nil`,
			},
			{Type: "paragraph", Text: "The detail that makes the drain work is what it does not do: it never marks the book as paused. Draining is about the process surviving, not about a user decision. The book stays marked as translating, and the next boot understands that it must be resumed."},
			{Type: "quote", Text: "A deploy is not a rare event: it is the weather the service lives in. Code has to be written to be interrupted.", AuthorQuote: "Matheus Puppe"},
			{Type: "heading2", Text: "The boot that reconciles"},
			{Type: "paragraph", Text: "On startup, before serving any request, the runner cleans up whatever the previous process left behind. Chapters stuck in running go back to the queue; paused books keep their pause, and whatever was pending in them parks instead of showing up as open work. Only books marked as translating start running again. The rule is simple and worth repeating: the boot resumes what the user asked for, never what they paused."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "ResumeInterrupted: the boot only resumes what was translating",
				Code: `func (r *Runner) ResumeInterrupted() error {
    if err := r.store.ResetRunningChapters(); err != nil {
        return err
    }
    if err := r.store.ParkPausedBookQueue(); err != nil {
        return err
    }
    books, err := r.store.ListBooks()
    if err != nil {
        return err
    }
    for _, book := range books {
        if book.Status == "translating" {
            r.Start(book.ID)
        }
    }
    return nil
}`,
			},
			{Type: "paragraph", Text: "That pair — drain on SIGTERM and reconciliation on boot — turns a deploy into something invisible for whoever is translating a book. The most the reader sees is the running chapter line blink and progress continue from the next pending one."},
			{Type: "heading2", Text: "Validation before the next attempt"},
			{Type: "paragraph", Text: "Translating XHTML is not swapping text for text: the model returns markup. If it eats a tag, the whole chapter is invalid. The defence is to mask the tags before the call, leaving placeholders in their place, and check the response by contract: the same normalised tag sequence as the source, text that is actually translated, a word ratio between 0.3 and 2.5, no preamble and no mixed chapter. A ratio outside that range almost always means a truncated answer or the model returning the original."},
			{Type: "paragraph", Text: "When validation fails, the service tries again — three attempts, with a reminder appended to the conversation naming what broke. But not every error deserves a retry. A context overflow and a 503 \"no upstreams\" from OpenRouter are not invalid translations: they are the provider out of capacity that minute, and insisting only burns time and money. Those errors surface to the chapter as an immediate failure."},
			{Type: "quote", Text: "The original XHTML is the source of truth; the model only translates the text. If it dropped a closing paragraph tag, the right move is to put it back — not to accept the broken page.", AuthorQuote: "Comment in the keepIfDropped code"},
			{Type: "paragraph", Text: "That is why masking has a second half. Any tag whose placeholder the model omitted is put back in the right position, in the original order. One index lost a closing paragraph tag — 401 tags in the source, 400 in the response — and the chapter was saved with the page intact because the source calls the shots."},
			{Type: "heading2", Text: "Cheap models, cost per chapter"},
			{Type: "paragraph", Text: "The translator speaks the OpenAI chat-completions protocol and nothing else. Each provider is a profile with a fixed base URL — OpenAI, Gemini, Groq, NVIDIA, OpenRouter and a custom one — and only one is active at a time, with its own stored key. Configuration is resolved on every call, so switching providers in the middle of a book takes effect on the next chapter without restarting anything. Each chapter row records provider, model, start, finish and duration of the last execution; configuration is frozen before the request so the table never lies afterwards."},
			{Type: "paragraph", Text: "The budget is per chapter for a practical reason: it is the only unit where comparing models is not self-deception. Defaults are deliberately cheap — a Gemini flash lite, a gpt-oss through Groq — because most of the work is long-form prose translation, not hard reasoning. When a chapter matters more, the model changes and the next row already carries the new one; what was translated is not retranslated. There is no automatic provider fallback when the selected one fails. That is an explicit choice: a silent fallback hides the failure and changes the bill without anyone noticing."},
			{Type: "paragraph", Text: "The size of each request is also a budget. Every provider has a character ceiling per call: Gemini takes 12,000 and returns up to 64k tokens; the NVIDIA DeepSeek stays at 6,000, because the gateway answers 504 close to four minutes of buffered completion; Riva works with 6,000 to fit inside its 8k context. The client has a ten-minute timeout per call, enough for providers that only send headers at the end of generation. And the Groq and OpenRouter model catalogs are cached for fifteen minutes: choosing a model is a query, not an event."},
			{Type: "heading2", Text: "What I would carry into the next service"},
			{Type: "paragraph", Text: "Five rules survived this project. First: write the state before calling the provider — the running row must exist before the first byte, or the restart has nothing to reconcile. Second: make the unit of work the smallest thing worth repeating; expensive retry is retry nobody uses. Third: cancellation is cooperative and needs a budget — a two-second drain inside a fifteen-second SIGTERM solves more than any infinite retry. Fourth: boot is part of the design; every stateful service needs a function that cleans up what the previous process left behind. Fifth: measure cost in the unit you retry; without that, the invoice is the only telemetry."},
			{Type: "paragraph", Text: "The translator is not a distributed system today: no external queue, and it runs on SQLite. Still, every one of these decisions came from a real incident — a deploy, a stuck chapter, an omitted tag, an unexpected bill. It is the kind of engineering that never shows up in the demo and shows up entirely on the day the service has to keep going on its own."},
		},
	}
}
