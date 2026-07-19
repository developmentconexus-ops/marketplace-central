package domain

// difalSeedRow holds only the R-04 interna_pct per UF (M07-plan.md P2-planner
// :92). InterestadualPct/EfetivoPct are derived by buildDifalSeed, never
// hand-duplicated, so the seed can't drift from the interestadual rule.
type difalSeedRow struct {
	UF         string
	InternaPct string
}

// rawDifalSeed is the R-04 27-UF interna_pct table, ordered by UF.
var rawDifalSeed = []difalSeedRow{
	{"AC", "19.0"},
	{"AL", "20.5"}, // post-01/04/2026 Lei 9.776/2025 rate; seed uses the new value
	{"AM", "20.0"},
	{"AP", "18.0"},
	{"BA", "20.5"},
	{"CE", "20.0"},
	{"DF", "20.0"},
	{"ES", "17.0"},
	{"GO", "19.0"},
	{"MA", "23.0"},
	{"MG", "18.0"},
	// MS: R-04 seeds 17.0; a separate source cites 19% (DISPUTED). Seed
	// stays 17.0 pending operator verification — this is a
	// verify-at-execution note (M07-plan.md SECTION 0), not a planning
	// unknown. An operator override, once persisted (S5/S6), takes
	// precedence via DifalRate.Override.
	{"MS", "17.0"},
	{"MT", "17.0"},
	{"PA", "19.0"},
	{"PB", "20.0"},
	{"PE", "20.5"},
	{"PI", "22.5"},
	{"PR", "19.5"},
	{"RJ", "20.0"},
	{"RN", "20.0"},
	{"RO", "19.5"},
	{"RR", "20.0"},
	{"RS", "17.0"},
	{"SC", "17.0"}, // origem of the interestadual 12/7 rule (IC-04)
	{"SE", "19.0"},
	{"SP", "18.0"},
	{"TO", "20.0"},
}

// DifalSeed is the R-04 27-UF DIFAL seed table (origem_versao
// "padrao-2026"), ordered by UF, with no active override.
var DifalSeed = buildDifalSeed()

func buildDifalSeed() []DifalRate {
	seed := make([]DifalRate, 0, len(rawDifalSeed))
	for _, row := range rawDifalSeed {
		interestadual := InterestadualPct(row.UF)
		seed = append(seed, DifalRate{
			UF:               row.UF,
			InternaPct:       row.InternaPct,
			InterestadualPct: interestadual,
			EfetivoPct:       computeEfetivoPct(row.InternaPct, interestadual),
			OrigemVersao:     OrigemVersaoPadrao2026,
		})
	}
	return seed
}
