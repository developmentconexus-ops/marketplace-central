package domain

// Regime is a tenant's tax regime, which determines the default aliquota
// used by a CalcProfile until the tenant overrides it.
type Regime string

const (
	RegimeSimples   Regime = "SIMPLES"
	RegimePresumido Regime = "PRESUMIDO"
)

// OrigemDefault marks a CalcProfile field as the built-in default rather
// than a tenant-set value; IC-04 requires this be explicit in the response,
// never implied by an absent/zero value.
const OrigemDefault = "default"

// Default decimal values (S1 decimal-string contract — never float64).
const (
	defaultSimplesAliquotaPct   = "4"
	defaultPresumidoAliquotaPct = "9.25"
	defaultLimiarVerdePct       = "18"
	defaultLimiarAmareloPct     = "10"
)

// DefaultAliquotaPct returns the regime's default aliquota as a decimal
// string. SIMPLES and any unrecognized regime fall back to the SIMPLES
// default, matching NewDefaultCalcProfile's initial state.
func DefaultAliquotaPct(regime Regime) string {
	if regime == RegimePresumido {
		return defaultPresumidoAliquotaPct
	}
	return defaultSimplesAliquotaPct
}

// CalcProfile is a tenant's pricing calculation configuration (IC-04).
// TarifaFull is nullable: nil means unknown/unset, never a zero Money
// (ADR-17) — a `full`-modalidade decomposition with a nil TarifaFull must
// propagate as an unknown component, not compute against 0.
type CalcProfile struct {
	Regime           Regime
	AliquotaPct      string
	LimiarVerdePct   string
	LimiarAmareloPct string
	TarifaFull       *Money
	DifalEnabled     bool
	DifalDestinoUF   *string
	Origem           string
}

// NewDefaultCalcProfile returns the initial state of a CalcProfile before
// any tenant edit: regime SIMPLES at its default aliquota, default
// thresholds, no tarifa_full, DIFAL off, and Origem explicitly "default"
// (IC-04 — never left to be inferred by the caller).
func NewDefaultCalcProfile() CalcProfile {
	return CalcProfile{
		Regime:           RegimeSimples,
		AliquotaPct:      DefaultAliquotaPct(RegimeSimples),
		LimiarVerdePct:   defaultLimiarVerdePct,
		LimiarAmareloPct: defaultLimiarAmareloPct,
		TarifaFull:       nil,
		DifalEnabled:     false,
		DifalDestinoUF:   nil,
		Origem:           OrigemDefault,
	}
}
