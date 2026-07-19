package domain

// Fonte marks where a resolved tariff component came from.
type Fonte string

const (
	// FontePadrao is degrau 4: the tenant/installation config default
	// (TariffDefaults) resolved with no live or per-listing data.
	FontePadrao Fonte = "PADRAO"
	// FonteManual marks a component the caller supplied by hand (a request
	// override), overriding whatever the resolver would have returned.
	// Degrau 0 (not resolved at all), never an estimate.
	FonteManual Fonte = "MANUAL"
	// (future: FonteListing, FonteConta, FonteML for degraus 1-3)
)

// ComponentResolution is one resolved tariff component (commission OR
// shipping). Valor is a decimal string; nil means NO DATA is available for
// this component (ADR-17: unknown operational facts never become zero).
type ComponentResolution struct {
	Valor      *string
	Fonte      Fonte
	Degrau     int
	Estimativa bool
	// Data (source timestamp) intentionally omitted at degrau 4 (config has
	// no per-resolve timestamp); add when degraus 1-3 land. Do NOT stub a
	// fake timestamp.
}

// TariffResolution bundles the resolved components for one solve/decompose.
type TariffResolution struct {
	// Comissao.Valor is a percent string (e.g. "13.00").
	Comissao ComponentResolution
	// Frete.Valor is a money string, or nil when no shipping estimate is
	// known (sem_dados policy, or estimativa policy with no amount set).
	Frete ComponentResolution
}
