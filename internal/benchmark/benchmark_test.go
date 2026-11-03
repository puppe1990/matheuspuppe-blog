package benchmark

import "testing"

func TestNormalize_ClampsAndSnaps(t *testing.T) {
	cases := []struct {
		in, wantReqs int
	}{
		{25000, 25000}, // off the slider grid but valid, like the shipped default
		{0, MinReqs},
		{1000, MinReqs},
		{200000, MaxReqs},
		{99999, 99999},
	}
	for _, c := range cases {
		got := Normalize(Params{Reqs: c.in, PayloadKB: 4})
		if got.Reqs != c.wantReqs {
			t.Errorf("Normalize(%d).Reqs = %d, want %d", c.in, got.Reqs, c.wantReqs)
		}
	}
}

func TestNormalize_PayloadKB(t *testing.T) {
	got := Normalize(Params{Reqs: 25000, PayloadKB: 99})
	if got.PayloadKB != MaxPayloadKB {
		t.Errorf("PayloadKB = %d, want %d", got.PayloadKB, MaxPayloadKB)
	}
}

// Expected values computed by hand from the shipped formulas (25k req/s, 4 KB).
func TestRuntimes_DefaultParams(t *testing.T) {
	runtimes := Runtimes(Params{Reqs: 25000, PayloadKB: 4})
	if len(runtimes) != 4 {
		t.Fatalf("got %d runtimes, want 4", len(runtimes))
	}

	cases := []struct {
		index    int
		name     string
		p99      float64
		mem, cpu int
		perWatt  int
	}{
		{0, "Rust (Tokio + Axum)", 0.56, 72, 26, 6127},
		{1, "Go (Gin / net/http)", 3.56, 226, 40, 3788},
		{2, "Node.js (Fastify)", 11.6, 434, 55, 2370},
		{3, "Python (FastAPI + Uvicorn)", 46.94, 689, 66, 1715},
	}
	for _, c := range cases {
		rt := runtimes[c.index]
		if rt.Name != c.name {
			t.Errorf("runtime %d name = %q, want %q", c.index, rt.Name, c.name)
		}
		if rt.P99MS != c.p99 {
			t.Errorf("%s p99 = %v, want %v", c.name, rt.P99MS, c.p99)
		}
		if rt.MemMB != c.mem {
			t.Errorf("%s mem = %d, want %d", c.name, rt.MemMB, c.mem)
		}
		if rt.CPUPercent != c.cpu {
			t.Errorf("%s cpu = %d, want %d", c.name, rt.CPUPercent, c.cpu)
		}
		if rt.PerWatt != c.perWatt {
			t.Errorf("%s perWatt = %d, want %d", c.name, rt.PerWatt, c.perWatt)
		}
	}
}

func TestRuntimes_P99FloorHolds(t *testing.T) {
	runtimes := Runtimes(Params{Reqs: MinReqs, PayloadKB: MinPayloadKB})
	if runtimes[0].P99MS < 0.4 {
		t.Errorf("rust p99 = %v, want >= 0.4", runtimes[0].P99MS)
	}
	if runtimes[1].P99MS < 1.8 {
		t.Errorf("go p99 = %v, want >= 1.8", runtimes[1].P99MS)
	}
	if runtimes[2].P99MS < 5.5 {
		t.Errorf("node p99 = %v, want >= 5.5", runtimes[2].P99MS)
	}
	if runtimes[3].P99MS < 12.0 {
		t.Errorf("python p99 = %v, want >= 12.0", runtimes[3].P99MS)
	}
}

func TestRuntimes_BarPercents(t *testing.T) {
	runtimes := Runtimes(Params{Reqs: 25000, PayloadKB: 4})
	if got := runtimes[0].P99Pct; got != 1.9 {
		t.Errorf("rust p99 bar = %v, want 1.9", got)
	}
	if got := runtimes[0].MemPct; got != 9.0 {
		t.Errorf("rust mem bar = %v, want 9.0", got)
	}
	for _, rt := range runtimes {
		if rt.P99Pct < 0 || rt.P99Pct > 100 || rt.MemPct < 0 || rt.MemPct > 100 {
			t.Errorf("%s bar percent out of range: p99=%v mem=%v", rt.Name, rt.P99Pct, rt.MemPct)
		}
	}
}

func TestRuntimes_OrderIsStable(t *testing.T) {
	runtimes := Runtimes(Params{Reqs: 50000, PayloadKB: 16})
	langs := []string{"Rust", "Go", "JavaScript / Node", "Python"}
	for i, lang := range langs {
		if runtimes[i].Lang != lang {
			t.Errorf("runtime %d lang = %q, want %q", i, runtimes[i].Lang, lang)
		}
	}
}
