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

// primeiroEnsaioEn é a edição inglesa do primeiro ensaio, com slug próprio. Os
// trechos de código são os do repositório público, citados no commit.
func primeiroEnsaioEn() englishEdition {
	return englishEdition{
		Slug:     "epub-translation-llms-resilient-queue",
		Title:    "Translating EPUBs with LLMs: a resilient queue, per-chapter retry and cheap models",
		Subtitle: "A deploy lands mid-translation, the model drops a tag and the bill grows. Running a book translator taught me how to design the queue, the cancellation and the cost per chapter.",
		Lead:     "A book is not a request. It is hundreds of XHTML documents, each taking minutes of inference, and a server that restarts in the middle of the work. The first version of my translator treated the book as a single job: a deploy in the middle of the afternoon was enough to lose hours of translation. The design that survived lives in a small application, with SQLite and one file per book, and solves the three problems that an LLM service has: interruption, validation and cost.",
		Tags:     []string{"LLM", "EPUB", "Queues", "Reliability", "Cost"},
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"The chapter is the unit of work: retry, pause, cost and restart recovery live in it, not in the book.",
				"The deploy lands mid-translation: SIGTERM needs a drain that returns the running chapter to the queue and leaves the book marked as translating for the next boot.",
				"The model is the least reliable part of the pipeline. Validate the response by contract (tags, ratio, language) before spending the next attempt.",
			}},
			{Type: "heading2", Text: "A book is not a request"},
			{Type: "paragraph", Text: "An EPUB is a zip with an ordered list of XHTML documents. The translator opens the package, reads the spine and turns each document into a chapter: one row in the chapters table, with a unique path inside the book, the source text, a cl100k token count and the execution state. Retry stops being \"try the book again\" and becomes \"try chapter 14\"; a failing document does not bring down the other 199; and the cost of each call becomes observable per unit, because tokens are measured before any request."},
			{Type: "paragraph", Text: "Sending the whole book in a single call is tempting: less code, one stream, one result. In practice the model's context becomes the product's limit, retry gets expensive, and any failure in the middle destroys the entire translation. Translating per document is more work to build and much less work to run."},
			{Type: "heading2", Text: "The four states of a chapter, the queue and the pause"},
			{Type: "paragraph", Text: "A chapter moves through four states: pending, running, done and failed. The queue is simple on purpose: pending and failed chapters, in document order. The runner keeps one goroutine per book, so two chapters of the same book never compete for the same file or the same database row. The retry button on a row takes another route: the chosen chapter jumps the queue and becomes the next one after the current work, without restarting the whole book."},
			{Type: "paragraph", Text: "Failure is a normal state too. A chapter that blew through the context, came back with the wrong tags or hit a provider outage ends as failed, with the error written on the row. The book is only marked failed if some chapter stayed that way. The interface stays honest: the table says what happened and retry points at the right problem."},
			{Type: "paragraph", Text: "Pausing works differently: it cancels the in-flight call, waits for the execution to finish and returns only the running chapter to the queue. The book becomes paused, finished chapters stay intact, and the interrupted chapter is not counted as a failure. When the reader resumes, the queue continues where it stopped."},
			{Type: "heading2", Text: "The deploy in the middle of a chapter"},
			{Type: "paragraph", Text: "The first bug showed up on a routine deploy. systemd sends SIGTERM and waits fifteen seconds; the HTTP server, in turn, allowed five seconds for connections to close. Except the book page keeps an SSE stream open the whole time to show progress: that stream never ends on its own, the shutdown blew through the deadline, the process exited with an error and killed the translation request in flight. The chapter stayed marked as running forever, and the interface lied until someone restarted the server again. In the week of the incident, a pause clicked after deploy 1887 left the book stuck for hours, because the paused state was written over work nobody had paused."},
			{Type: "paragraph", Text: "The fix joined two things. The runner gained a drain that recognises shutdown: cancel each book's context, wait up to two seconds for the work to finish and return chapters stuck in running to the queue. And shutdown order started to matter: drain first, then the HTTP server, and SQLite last. Without that order the drain writes to a database that is already closed."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "Shutdown order: drain the queue before closing the database",
				Code: `select {
case <-ctx.Done():
	if a.runner != nil {
		_ = a.runner.Drain(2 * time.Second)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		_ = a.server.Close()
	}
	select {
	case <-errCh:
	case <-time.After(time.Second):
	}
	_ = a.store.Close()
	return nil
case err := <-errCh:
	_ = a.store.Close()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}`,
				SourceURL:   "https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/app/app.go#L148-L170",
				SourceLabel: "internal/app/app.go:148",
			},
			{Type: "paragraph", Text: "The drain does not mark the book as paused. Draining is about the process surviving, not about obeying a user: the book stays marked as translating, and the next boot understands that it must be resumed."},
			{Type: "heading2", Text: "The boot that reconciles"},
			{Type: "paragraph", Text: "On startup, before serving any request, the runner cleans up whatever the previous process left behind. Chapters stuck in running go back to the queue; paused books keep their pause, and whatever was pending in them parks instead of showing up as open work. Only books marked as translating start running again."},
			{
				Type:        "quote",
				Text:        "O boot de um deploy não arranca o que o usuário não pediu.",
				AuthorQuote: "Comment in the repository code",
				SourceURL:   "https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/pipeline/runner.go#L44-L47",
				SourceLabel: "internal/pipeline/runner.go:44",
			},
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
				SourceURL:   "https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/pipeline/runner.go#L48-L65",
				SourceLabel: "internal/pipeline/runner.go:48",
			},
			{Type: "paragraph", Text: "That pair, drain on SIGTERM and reconciliation on boot, is what makes a deploy go almost unnoticed by whoever is translating a book. The most the reader sees is the running chapter line blink and progress continue from the next pending one."},
			{Type: "heading2", Text: "Validation before the next attempt"},
			{Type: "paragraph", Text: "Translating XHTML is not swapping text for text: the model returns markup. If it eats a tag, the whole chapter is invalid. The defence is to mask the tags before the call, leaving placeholders in their place, and check the response by contract: the same normalised tag sequence as the source, text that is actually translated, a word ratio between 0.3 and 2.5, no preamble and no mixed chapter. A ratio outside that range almost always means a truncated answer or the model returning the original."},
			{Type: "paragraph", Text: "When validation fails, the service tries again, three attempts, with a reminder appended to the conversation naming what broke. Not every error deserves a retry: a context overflow and a 503 \"no upstreams\" from OpenRouter are not invalid translations, they are the provider out of capacity that minute, and insisting only burns time and money. Those errors surface to the chapter as an immediate failure."},
			{
				Type:        "quote",
				Text:        "O XHTML original é a fonte da verdade; o modelo só traduz o texto. Sem isso o INDEX perde um </p> (cap. 222: fonte 401, resposta 400).",
				AuthorQuote: "Comment in the repository code",
				SourceURL:   "https://github.com/puppe1990/epub-translator/blob/a4fa7ce/internal/translate/mask.go#L31-L36",
				SourceLabel: "internal/translate/mask.go:31",
			},
			{Type: "paragraph", Text: "That is why masking has a second half: any tag whose placeholder the model omitted goes back to the right position, in the original order. The INDEX chapter that arrived with one closing tag missing was saved with the page intact, because the source calls the shots."},
			{Type: "heading2", Text: "Cheap models, cost per chapter"},
			{Type: "paragraph", Text: "The translator speaks the OpenAI chat-completions protocol and nothing else. Each provider is a profile with a fixed base URL (OpenAI, Gemini, Groq, NVIDIA, OpenRouter and a custom one); only one is active at a time, with its own stored key. Configuration is resolved on every call, so switching providers in the middle of a book takes effect on the next chapter without restarting anything. Each chapter row records provider, model, start, finish and duration of the last execution; configuration is frozen before the request so the table never lies afterwards."},
			{Type: "paragraph", Text: "The budget is per chapter for a practical reason: it is the only unit where comparing models is not self-deception. Defaults are deliberately cheap (a Gemini flash lite, a gpt-oss through Groq), because most of the work is long-form prose translation, not hard reasoning. When a chapter matters more, the model changes and the next row already carries the new one; what was translated is not retranslated. There is no automatic provider fallback when the selected one fails. That is an explicit choice: a silent fallback hides the failure and changes the bill without anyone noticing."},
			{Type: "paragraph", Text: "The size of each piece also has a ceiling. Every provider has a character limit per call: Gemini takes 12,000 and returns up to 65,536 tokens; NVIDIA's DeepSeek stays at 6,000, because the gateway answers 504 close to four minutes of buffered completion; Riva works with 6,000 to fit inside its 8k-token window, half of it reserved for the output. The client has a ten-minute timeout per call, enough for providers that only send headers at the end of generation. And the Groq and OpenRouter model catalogs are cached for fifteen minutes: choosing a model is just a query."},
			{Type: "heading2", Text: "What I would carry into the next service"},
			{Type: "paragraph", Text: "Writing the state before calling the provider looks like an implementation detail, but it is what gives the restart something to reconcile: the running row must exist before the first byte. And boot is part of the design, because every stateful service needs a function that cleans up what the previous process left behind."},
			{Type: "paragraph", Text: "Choosing the unit of work is also choosing what you will repeat; expensive retry is retry nobody uses. Cancellation is cooperative and needs a budget: two seconds of drain inside a fifteen-second SIGTERM solves more than any infinite retry. And measure cost in the same unit you retry, because without that the invoice is the only telemetry."},
			{
				Type:        "paragraph",
				Text:        "The translator is not a distributed system today, has no external queue and runs on SQLite. Every one of these decisions came from a real incident: a deploy, a stuck chapter, an omitted tag, an unexpected bill. The code is open at",
				SourceURL:   "https://github.com/puppe1990/epub-translator",
				SourceLabel: "github.com/puppe1990/epub-translator",
			},
		},
	}
}
