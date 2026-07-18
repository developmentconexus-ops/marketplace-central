package domain

// DisclaimerSeedPadrao2026 is the mandatory disclaimer every DIFAL surface
// (C02/C05) must render alongside the seed rates — the table is a pricing
// aid, not tax guidance.
const DisclaimerSeedPadrao2026 = "seed padrão 2026 — não é orientação fiscal"

// OrigemVersaoPadrao2026 identifies the R-04 seed table version carried on
// every DifalRate.
const OrigemVersaoPadrao2026 = "padrao-2026"

// interestadual12 is the set of destino UFs whose interestadual rate is
// 12%; every other UF uses 7%. Origem of this rule is fixed SC per IC-04.
var interestadual12 = map[string]bool{
	"MG": true, "PR": true, "RJ": true, "RS": true, "SC": true, "SP": true,
}

// InterestadualPct returns the interestadual DIFAL rate for destino uf as a
// decimal string: "12" for {MG,PR,RJ,RS,SC,SP}, "7" for every other UF.
func InterestadualPct(uf string) string {
	if interestadual12[uf] {
		return "12"
	}
	return "7"
}

// DifalOverride is a tenant operator's override of a UF's interna_pct.
// Modeled here as an optional field on DifalRate so the S3 port only wires
// the interface; persistence (0056 override_* columns) and the >0,049pp
// audit rule land in S5/S6.
type DifalOverride struct {
	InternaPct string
	UpdatedAt  string
	Actor      string
}

// DifalRate is one UF's DIFAL configuration: the R-04 seed interna_pct, the
// derived interestadual_pct, the resulting efetivo_pct, and an optional
// active override.
type DifalRate struct {
	UF               string
	InternaPct       string
	InterestadualPct string
	EfetivoPct       string
	OrigemVersao     string
	Override         *DifalOverride
}

// DifalForUFResult is the frozen IC-04 port output shape:
// DifalForUF(uf) → {efetivo_pct, versao}.
type DifalForUFResult struct {
	EfetivoPct string
	Versao     string
}

// DifalForUF is the pure computation behind the frozen IC-04 port (wired at
// S3): efetivo = max(interna - interestadual, 0). When rate has an active
// Override, the override's interna_pct is used in place of the seed value.
func DifalForUF(rate DifalRate) DifalForUFResult {
	interna := rate.InternaPct
	if rate.Override != nil {
		interna = rate.Override.InternaPct
	}
	return DifalForUFResult{
		EfetivoPct: computeEfetivoPct(interna, rate.InterestadualPct),
		Versao:     rate.OrigemVersao,
	}
}

// computeEfetivoPct applies max(interna-interestadual, 0) and formats the
// result at 2 decimal places. 2dp (not 1) because an operator override may
// set interna_pct at 2dp (e.g. "19.05"); rounding to 1dp would silently
// drop that precision, and this value is FROZEN into the DifalForUF port
// (S3) consumed by M-08. Seed rows are exact at 1dp so 2dp is invisible for
// them ("6.0" → "6.00"); the override path is what forces 2dp here.
func computeEfetivoPct(internaPct, interestadualPct string) string {
	interna, err := ParseRat(internaPct)
	if err != nil {
		panic(err)
	}
	inter, err := ParseRat(interestadualPct)
	if err != nil {
		panic(err)
	}
	diff := interna.Sub(interna, inter)
	if diff.Sign() < 0 {
		diff.SetInt64(0)
	}
	return FormatRatHalfUp(diff, 2)
}
