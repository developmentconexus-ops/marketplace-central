package domain

import "time"

// ProductLinkDecisionRule names the evidence a link was approved on (E10).
type ProductLinkDecisionRule string

const (
	// DecisionRuleConcordantCodprodEAN — the listing's CODPROD and its EAN both
	// named the same product. Corroboration by two independent anchors is the
	// ONLY rule the system may apply on its own (D-121-2, operator-ratified).
	DecisionRuleConcordantCodprodEAN ProductLinkDecisionRule = "concordant_codprod_ean"
	// DecisionRuleExactCodprodUnique / DecisionRuleExactEANUnique — a single
	// anchor matched exactly one product. That is enough to PROPOSE the product
	// and ask the operator to confirm, never enough to decide alone, so these
	// rules only ever appear with actor='operator'.
	DecisionRuleExactCodprodUnique ProductLinkDecisionRule = "exact_codprod_unique"
	DecisionRuleExactEANUnique     ProductLinkDecisionRule = "exact_ean_unique"
	// DecisionRuleManual — the operator named the product themselves, with no
	// anchor carrying the decision.
	DecisionRuleManual ProductLinkDecisionRule = "manual"
)

const (
	DecisionActorSystem   = "system"
	DecisionActorOperator = "operator"
)

// ProductLinkDecision is one E10 audit row: why a link was approved, by whom,
// and how many products the deciding anchor matched at that moment. It is
// written in the same transaction as the link's state transition — an approved
// link whose decision row is missing would leave the trail a guess.
type ProductLinkDecision struct {
	DecisionID          string                  `json:"decision_id"`
	InstallationID      string                  `json:"installation_id"`
	ProviderItemID      string                  `json:"provider_item_id"`
	ProviderVariationID string                  `json:"provider_variation_id,omitempty"`
	LinkID              string                  `json:"link_id"`
	RuleMatched         ProductLinkDecisionRule `json:"rule_matched"`
	Actor               string                  `json:"actor"`
	// CollisionsAtDecision is how many products the deciding anchor matched.
	// nil when no collision count was read for this decision (a manual
	// resolution decides no anchor at all) — ADR-17: absent, never 0, which
	// would claim the anchor matched nothing.
	CollisionsAtDecision *int      `json:"collisions_at_decision,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	// SupersededBy is the decision that replaced this one. Empty while this is
	// the decision in force.
	SupersededBy string `json:"superseded_by,omitempty"`
}

// LinkID renders the listing identity as the link identifier E10 records.
// product_links has no surrogate id: the link IS its natural key.
func LinkID(installationID, providerItemID, providerVariationID string) string {
	return installationID + "|" + providerItemID + "|" + providerVariationID
}

// IsAutomatic reports whether the system took this decision unaided.
func (d ProductLinkDecision) IsAutomatic() bool { return d.Actor == DecisionActorSystem }
