package domain

import "testing"

// r04Seed is the golden R-04 27-UF interna_pct table (P2-planner.md :92).
// The map key order is irrelevant here; DifalSeed's own order is asserted
// separately.
var r04Seed = map[string]string{
	"AC": "19.0",
	"AL": "20.5",
	"AM": "20.0",
	"AP": "18.0",
	"BA": "20.5",
	"CE": "20.0",
	"DF": "20.0",
	"ES": "17.0",
	"GO": "19.0",
	"MA": "23.0",
	"MG": "18.0",
	"MS": "17.0", // DISPUTED — see difal_seed.go comment; seed stays 17.0
	"MT": "17.0",
	"PA": "19.0",
	"PB": "20.0",
	"PE": "20.5",
	"PI": "22.5",
	"PR": "19.5",
	"RJ": "20.0",
	"RN": "20.0",
	"RO": "19.5",
	"RR": "20.0",
	"RS": "17.0",
	"SC": "17.0",
	"SE": "19.0",
	"SP": "18.0",
	"TO": "20.0",
}

var interestadual12Set = map[string]bool{
	"MG": true, "PR": true, "RJ": true, "RS": true, "SC": true, "SP": true,
}

func TestDifalSeedHasExactly27EntriesOrderedByUF(t *testing.T) {
	if len(DifalSeed) != 27 {
		t.Fatalf("len(DifalSeed) = %d, want 27", len(DifalSeed))
	}
	for i := 1; i < len(DifalSeed); i++ {
		if DifalSeed[i-1].UF >= DifalSeed[i].UF {
			t.Fatalf("DifalSeed not ordered by UF: %q before %q", DifalSeed[i-1].UF, DifalSeed[i].UF)
		}
	}
}

func TestDifalSeedMatchesR04Exactly(t *testing.T) {
	seen := map[string]bool{}
	for _, rate := range DifalSeed {
		want, ok := r04Seed[rate.UF]
		if !ok {
			t.Fatalf("DifalSeed contains unexpected UF %q", rate.UF)
		}
		if rate.InternaPct != want {
			t.Fatalf("UF %s: InternaPct = %q, want %q (R-04)", rate.UF, rate.InternaPct, want)
		}
		seen[rate.UF] = true
	}
	for uf := range r04Seed {
		if !seen[uf] {
			t.Fatalf("R-04 UF %q missing from DifalSeed", uf)
		}
	}
}

func TestDifalSeedInterestadualPct(t *testing.T) {
	for _, rate := range DifalSeed {
		want := "7"
		if interestadual12Set[rate.UF] {
			want = "12"
		}
		if rate.InterestadualPct != want {
			t.Fatalf("UF %s: InterestadualPct = %q, want %q", rate.UF, rate.InterestadualPct, want)
		}
	}
}

func TestDifalSeedEfetivoPctIsMaxInternaMinusInterestadualZeroFloor(t *testing.T) {
	cases := map[string]string{
		"MG": "6.00",
		"PR": "7.50",
		"RJ": "8.00",
		"RS": "5.00",
		"SC": "5.00",
		"SP": "6.00",
	}
	byUF := map[string]DifalRate{}
	for _, rate := range DifalSeed {
		byUF[rate.UF] = rate
	}
	for uf, want := range cases {
		got := byUF[uf].EfetivoPct
		if got != want {
			t.Fatalf("UF %s: EfetivoPct = %q, want %q", uf, got, want)
		}
	}
}

func TestDifalSeedOrigemVersaoAndDisclaimer(t *testing.T) {
	for _, rate := range DifalSeed {
		if rate.OrigemVersao != "padrao-2026" {
			t.Fatalf("UF %s: OrigemVersao = %q, want padrao-2026", rate.UF, rate.OrigemVersao)
		}
	}
	if DisclaimerSeedPadrao2026 != "seed padrão 2026 — não é orientação fiscal" {
		t.Fatalf("DisclaimerSeedPadrao2026 = %q, want the mandated disclaimer text", DisclaimerSeedPadrao2026)
	}
}

func TestDifalSeedNoUFHasNegativeEfetivo(t *testing.T) {
	// interna is always >= 0 and the clamp formula is max(interna-inter,0);
	// this guards against a future seed edit introducing a negative value.
	for _, rate := range DifalSeed {
		v, err := ParseRat(rate.EfetivoPct)
		if err != nil {
			t.Fatalf("UF %s: EfetivoPct %q is not a decimal string: %v", rate.UF, rate.EfetivoPct, err)
		}
		if v.Sign() < 0 {
			t.Fatalf("UF %s: EfetivoPct %q is negative", rate.UF, rate.EfetivoPct)
		}
	}
}

func TestDifalForUFReflectsSeedWhenNoOverride(t *testing.T) {
	var sp DifalRate
	for _, rate := range DifalSeed {
		if rate.UF == "SP" {
			sp = rate
		}
	}
	result := DifalForUF(sp)
	if result.EfetivoPct != "6.00" {
		t.Fatalf("DifalForUF(SP).EfetivoPct = %q, want 6.00", result.EfetivoPct)
	}
	if result.Versao != "padrao-2026" {
		t.Fatalf("DifalForUF(SP).Versao = %q, want padrao-2026", result.Versao)
	}
}

func TestDifalForUFReflectsActiveOverride(t *testing.T) {
	var ba DifalRate
	for _, rate := range DifalSeed {
		if rate.UF == "BA" {
			ba = rate
		}
	}
	ba.Override = &DifalOverride{
		InternaPct: "21.0",
		UpdatedAt:  "2026-07-18T00:00:00Z",
		Actor:      "operator-1",
	}
	// BA interestadual is 7 (not in the 12% set): 21.0 - 7 = 14.00 (2dp wire).
	result := DifalForUF(ba)
	if result.EfetivoPct != "14.00" {
		t.Fatalf("DifalForUF(BA override) EfetivoPct = %q, want 14.00", result.EfetivoPct)
	}
}

// TestDifalForUFOverridePreserves2dpPrecision proves the 2dp wire precision is
// load-bearing: an operator override at 2dp (19.05) must survive into the
// frozen DifalForUF port output (consumed by M-08), not round to 1dp.
func TestDifalForUFOverridePreserves2dpPrecision(t *testing.T) {
	var ce DifalRate // CE is 7% group (interna 20.0)
	for _, rate := range DifalSeed {
		if rate.UF == "CE" {
			ce = rate
		}
	}
	ce.Override = &DifalOverride{
		InternaPct: "19.05",
		UpdatedAt:  "2026-07-18T00:00:00Z",
		Actor:      "operator-1",
	}
	// 19.05 - 7 = 12.05; at 1dp this would wrongly collapse to "12.1".
	result := DifalForUF(ce)
	if result.EfetivoPct != "12.05" {
		t.Fatalf("DifalForUF(CE override 19.05) EfetivoPct = %q, want 12.05 (2dp must survive)", result.EfetivoPct)
	}
}

// TestDifalSeedSevenPercentGroupEfetivo asserts every non-12% UF's derived
// efetivo equals interna-7 at 2dp — guards the derivation across all 21 rows,
// not just the six 12%-group goldens above. Uses the production
// computeEfetivoPct (same package) as the oracle.
func TestDifalSeedSevenPercentGroupEfetivo(t *testing.T) {
	for _, rate := range DifalSeed {
		if interestadual12Set[rate.UF] {
			continue
		}
		want := computeEfetivoPct(rate.InternaPct, "7")
		if rate.EfetivoPct != want {
			t.Fatalf("UF %s (7%% group): EfetivoPct = %q, want %q (interna-7 @2dp)", rate.UF, rate.EfetivoPct, want)
		}
	}
}
