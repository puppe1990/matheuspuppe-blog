package handlers

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/matheuspuppe-blog/internal/benchmark"
	"github.com/puppe1990/matheuspuppe-blog/internal/editorial"
)

// runtimeView adds pt-BR formatted numbers and marks the leading stack (lowest
// p99) — the only row that carries the accent.
type runtimeView struct {
	benchmark.Runtime
	PerWattLabel string
	Leader       bool
}

func runtimeViews(runtimes []benchmark.Runtime) []runtimeView {
	leader := 0
	for i, runtime := range runtimes {
		if runtime.P99MS < runtimes[leader].P99MS {
			leader = i
		}
	}
	views := make([]runtimeView, len(runtimes))
	for i, runtime := range runtimes {
		views[i] = runtimeView{
			Runtime:      runtime,
			PerWattLabel: editorial.FormatThousands(runtime.PerWatt),
			Leader:       i == leader,
		}
	}
	return views
}

// LaboratorioHandler serves the runtime simulator; the math is server-side and
// the sliders only refresh labels until a GET form re-renders with new values.
type LaboratorioHandler struct {
	views   *view.Renderer
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewLaboratorioHandler(views *view.Renderer, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *LaboratorioHandler {
	return &LaboratorioHandler{views: views, site: site, catalog: catalog, cfg: cfg}
}

func (h *LaboratorioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := benchmark.Normalize(benchmark.Params{
		Reqs:      intParam(r, "reqs", benchmark.DefaultReqs),
		PayloadKB: intParam(r, "payload", benchmark.DefaultPayloadKB),
	})

	runtimes := benchmark.Runtimes(params)
	views := runtimeViews(runtimes)

	writeView(w, r, h.views, h.cfg, "app", "laboratorio", amarraData(r, h.site, map[string]any{
		"Title":            "Playground & Benchmarks · Sintaxe",
		"ActiveNav":        "playground",
		"Params":           params,
		"ReqsLabel":        editorial.FormatThousands(params.Reqs),
		"Runtimes":         views,
		"MinReqs":          benchmark.MinReqs,
		"MaxReqs":          benchmark.MaxReqs,
		"StepReqs":         benchmark.StepReqs,
		"MinPayloadKB":     benchmark.MinPayloadKB,
		"MaxPayloadKB":     benchmark.MaxPayloadKB,
		"DefaultReqs":      benchmark.DefaultReqs,
		"DefaultPayloadKB": benchmark.DefaultPayloadKB,
	}), http.StatusOK)
}
