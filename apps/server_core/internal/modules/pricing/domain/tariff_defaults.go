package domain

import "errors"

// Frete policy values for TariffDefaults.FretePolicy. "estimativa" means
// FreteEstimativaAmount carries a known estimate; "sem_dados" means no
// shipping estimate is available (ADR-17 — never defaulted to 0).
const (
	FretePolicyEstimativa = "estimativa"
	FretePolicySemDados   = "sem_dados"
)

// ErrInvalidFretePolicy is returned when a caller supplies a FretePolicy
// outside {estimativa, sem_dados} (transport maps it to 422).
var ErrInvalidFretePolicy = errors.New("PRICING_INVALID_FRETE_POLICY")

// TariffDefaults is a tenant/installation's pricing tariff configuration
// (CHIP-T1 Slice A): the classic/premium commission percentages (DB-default
// seeded 13.00/16.00 — never Go constants) and the shipping estimate policy.
// All money/percent fields are decimal strings, never float64.
// FreteEstimativaAmount is nullable: nil means no shipping data is known,
// never a zero amount (ADR-17).
type TariffDefaults struct {
	ComissaoClassicoPct   string
	ComissaoPremiumPct    string
	FreteEstimativaAmount *string
	FretePolicy           string
}

// ValidFretePolicy reports whether policy is one of the two allowed values.
func ValidFretePolicy(policy string) bool {
	return policy == FretePolicyEstimativa || policy == FretePolicySemDados
}
