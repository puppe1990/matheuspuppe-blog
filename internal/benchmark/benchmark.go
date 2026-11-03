// Package benchmark simulates the Sintaxe runtime laboratory. The curves are
// the shipped editorial model (one fixed profile per runtime family) and are
// deterministic, so the page renders the same numbers server-side; the JS only
// refreshes the slider labels.
package benchmark

import "math"

// Slider bounds and defaults, mirrored by the laboratorio page controls.
const (
	MinReqs          = 2000
	MaxReqs          = 100000
	StepReqs         = 2000
	MinPayloadKB     = 1
	MaxPayloadKB     = 32
	DefaultReqs      = 25000
	DefaultPayloadKB = 4
)

// Params is the simulated traffic: requests per second and JSON payload size.
type Params struct {
	Reqs      int
	PayloadKB int
}

// Normalize clamps arbitrary query values into the supported ranges.
func Normalize(p Params) Params {
	p.Reqs = clamp(p.Reqs, MinReqs, MaxReqs)
	p.PayloadKB = clamp(p.PayloadKB, MinPayloadKB, MaxPayloadKB)
	return p
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Runtime is one simulated backend stack with its mechanical response.
type Runtime struct {
	Name       string
	Lang       string
	Color      string // tailwind palette key: amber, cyan, emerald, indigo
	P99MS      float64
	MemMB      int
	CPUPercent int
	PerWatt    int
	P99Pct     float64 // bar width against a 30 ms ceiling
	MemPct     float64 // bar width against an 800 MB ceiling
}

// Runtimes computes the four profiles for the given traffic shape.
func Runtimes(p Params) []Runtime {
	reqs := float64(p.Reqs)
	payload := float64(p.PayloadKB)

	rustCPU := int(math.Min(95, math.Round(15+reqs/100000*45)))
	goCPU := int(math.Min(98, math.Round(25+reqs/100000*60)))
	nodeCPU := int(math.Min(99, math.Round(40+reqs/100000*58)))
	pyCPU := int(math.Min(100, math.Round(55+reqs/100000*45)))

	runtimes := []Runtime{
		{
			Name:       "Rust (Tokio + Axum)",
			Lang:       "Rust",
			Color:      "amber",
			P99MS:      maxF(0.4, round2(0.3+reqs/100000*0.8+payload/32*0.5)),
			MemMB:      int(math.Round(45 + reqs/1000*0.6 + payload*3)),
			CPUPercent: rustCPU,
		},
		{
			Name:       "Go (Gin / net/http)",
			Lang:       "Go",
			Color:      "cyan",
			P99MS:      maxF(1.8, round2(1.5+math.Pow(reqs/40000, 1.8)*4.2+payload/32*2.1)),
			MemMB:      int(math.Round(120 + reqs/1000*2.8 + payload*9)),
			CPUPercent: goCPU,
		},
		{
			Name:       "Node.js (Fastify)",
			Lang:       "JavaScript / Node",
			Color:      "emerald",
			P99MS:      maxF(5.5, round2(4.5+math.Pow(reqs/30000, 2.1)*9.5+payload/32*5.0)),
			MemMB:      int(math.Round(240 + reqs/1000*5.2 + payload*16)),
			CPUPercent: nodeCPU,
		},
		{
			Name:       "Python (FastAPI + Uvicorn)",
			Lang:       "Python",
			Color:      "indigo",
			P99MS:      maxF(12.0, round2(9.0+math.Pow(reqs/20000, 2.3)*22.0+payload/32*9.5)),
			MemMB:      int(math.Round(380 + reqs/1000*8.5 + payload*24)),
			CPUPercent: pyCPU,
		},
	}

	runtimes[0].PerWatt = int(math.Round(reqs / (float64(rustCPU)*0.8 + 20) * 10))
	runtimes[1].PerWatt = int(math.Round(reqs / (float64(goCPU)*0.9 + 30) * 10))
	runtimes[2].PerWatt = int(math.Round(reqs / (float64(nodeCPU)*1.1 + 45) * 10))
	runtimes[3].PerWatt = int(math.Round(reqs / (float64(pyCPU)*1.3 + 60) * 10))

	for i := range runtimes {
		runtimes[i].P99Pct = round1(math.Min(100, runtimes[i].P99MS/30*100))
		runtimes[i].MemPct = round1(math.Min(100, float64(runtimes[i].MemMB)/800*100))
	}
	return runtimes
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func maxF(a, b float64) float64 {
	return math.Max(a, b)
}
