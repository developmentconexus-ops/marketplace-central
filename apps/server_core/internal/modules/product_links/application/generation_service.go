package application

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type ProductMatcher interface {
	FindProductsForLinking(ctx context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error)
}

// LinkAutoApprover applies the one automatic path (D-121-2): a candidate whose
// CODPROD and EAN resolved the same product. It is an interface so generation
// stays testable without a workflow store, and optional so a deployment that
// has not wired it still generates candidates — it just approves none.
type LinkAutoApprover interface {
	AutoApproveCandidate(ctx context.Context, input AutoApproveCandidateInput) (bool, error)
}

type GenerationService struct {
	snapshots       ports.ListingSnapshotReader
	matcher         ProductMatcher
	store           ports.LinkCandidateStore
	identityAnchors ports.ProviderIdentityAnchorReader
	autoApprover    LinkAutoApprover
	now             func() time.Time
}

type GenerationServiceConfig struct {
	Snapshots       ports.ListingSnapshotReader
	Matcher         ProductMatcher
	Store           ports.LinkCandidateStore
	IdentityAnchors ports.ProviderIdentityAnchorReader
	AutoApprover    LinkAutoApprover
	Now             func() time.Time
}

type GenerateLinkCandidatesInput struct {
	InstallationID string
	Limit          int
}

type ListLinkCandidatesInput struct {
	InstallationID string
	Limit          int
}

func NewGenerationService(cfg GenerationServiceConfig) *GenerationService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &GenerationService{
		snapshots:       cfg.Snapshots,
		matcher:         cfg.Matcher,
		store:           cfg.Store,
		identityAnchors: cfg.IdentityAnchors,
		autoApprover:    cfg.AutoApprover,
		now:             now,
	}
}

func (s *GenerationService) GenerateLinkCandidates(ctx context.Context, input GenerateLinkCandidatesInput) (domain.LinkCandidateGenerationResult, error) {
	if s.snapshots == nil || s.matcher == nil || s.store == nil || s.identityAnchors == nil {
		return domain.LinkCandidateGenerationResult{}, errors.New("PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED")
	}

	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return domain.LinkCandidateGenerationResult{}, errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}

	// An absent limit means the WHOLE installation, never a 20-row default page.
	// Generation replaces the candidate set for the identities it processed, so a
	// capped run leaves every uncapped listing without a candidate — the operator
	// then reads thousands of anúncios as "sem vínculo" when the matcher was
	// simply never asked about them. Callers that want a sample pass one.
	snapshots, err := s.snapshots.ListListingSnapshots(ctx, installationID, input.Limit)
	if err != nil {
		return domain.LinkCandidateGenerationResult{}, err
	}
	identityAnchors, err := s.resolveIdentityAnchors(snapshots)
	if err != nil {
		return domain.LinkCandidateGenerationResult{}, err
	}

	identities := make([]domain.ListingIdentity, 0, len(snapshots))
	candidates := make([]domain.LinkCandidate, 0, len(snapshots))
	pendingApprovals := make([]autoApproval, 0, len(snapshots))
	for _, snapshot := range snapshots {
		identities = append(identities, domain.ListingIdentity{
			InstallationID:      installationID,
			ProviderItemID:      snapshot.ProviderItemID,
			ProviderVariationID: snapshot.ProviderVariationID,
		})

		generated, approvals, err := s.generateForSnapshot(ctx, snapshot, identityAnchors[strings.TrimSpace(snapshot.ProviderCode)])
		if err != nil {
			return domain.LinkCandidateGenerationResult{}, classifyMatcherError(err)
		}
		candidates = append(candidates, generated...)
		pendingApprovals = append(pendingApprovals, approvals...)
	}

	if err := s.store.ReplaceLinkCandidates(ctx, installationID, identities, candidates); err != nil {
		return domain.LinkCandidateGenerationResult{}, err
	}

	// Auto-approval runs only after the candidates are persisted: the link it
	// writes points at a candidate row, and approving against a candidate set
	// that failed to store would leave a link whose evidence does not exist.
	if s.autoApprover != nil {
		// Every corroborated candidate is attempted. One listing that fails to
		// approve is not a reason to leave the rest of the batch unlinked, and
		// the failures are still reported rather than swallowed.
		var approvalErr error
		for _, approval := range pendingApprovals {
			if _, err := s.autoApprover.AutoApproveCandidate(ctx, AutoApproveCandidateInput{
				Candidate:            approval.Candidate,
				CollisionsAtDecision: approval.Collisions,
			}); err != nil {
				approvalErr = errors.Join(approvalErr, err)
			}
		}
		if approvalErr != nil {
			return domain.LinkCandidateGenerationResult{}, approvalErr
		}
	}

	return domain.LinkCandidateGenerationResult{
		InstallationID: installationID,
		GeneratedCount: len(candidates),
		Items:          candidates,
	}, nil
}

func (s *GenerationService) resolveIdentityAnchors(snapshots []domain.ListingSnapshot) (map[string][]ports.ProviderIdentityAnchor, error) {
	resolved := make(map[string][]ports.ProviderIdentityAnchor, len(snapshots))
	for _, snapshot := range snapshots {
		providerCode := strings.TrimSpace(snapshot.ProviderCode)
		if _, ok := resolved[providerCode]; ok {
			continue
		}
		if providerCode == "" {
			return nil, unavailableIdentityAnchorsError(providerCode, errors.New("provider code is empty"))
		}
		declaration, err := s.identityAnchors.ProviderIdentityAnchors(providerCode)
		if err != nil {
			return nil, unavailableIdentityAnchorsError(providerCode, err)
		}
		if declaration == nil {
			return nil, unavailableIdentityAnchorsError(providerCode, errors.New("identity anchor declaration is nil"))
		}
		resolved[providerCode] = slices.Clone(declaration)
	}
	return resolved, nil
}

func unavailableIdentityAnchorsError(providerCode string, cause error) error {
	return fmt.Errorf("PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE for provider %q: %w", providerCode, cause)
}

func (s *GenerationService) ListLinkCandidates(ctx context.Context, input ListLinkCandidatesInput) ([]domain.LinkCandidate, error) {
	if s.store == nil {
		return nil, errors.New("PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return nil, errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	return s.store.ListLinkCandidates(ctx, installationID, limit)
}

func (s *GenerationService) generateForSnapshot(ctx context.Context, snapshot domain.ListingSnapshot, identityAnchors []ports.ProviderIdentityAnchor) ([]domain.LinkCandidate, []autoApproval, error) {
	now := s.now().UTC()

	skuMatches, err := s.findProducts(ctx, snapshot.SellerSKU, domain.LinkCandidateMatchInputSellerSKU)
	if err != nil {
		return nil, nil, err
	}
	eanMatches, err := s.findProducts(ctx, snapshot.EAN, domain.LinkCandidateMatchInputEAN)
	if err != nil {
		return nil, nil, err
	}

	exactCandidates := buildExactCandidates(snapshot, skuMatches, eanMatches, identityAnchors, now)
	if len(exactCandidates) > 0 {
		return exactCandidates, autoApprovals(exactCandidates, skuMatches), nil
	}

	titleMatches, err := s.findProducts(ctx, snapshot.Title, domain.LinkCandidateMatchInputTitle)
	if err != nil {
		return nil, nil, err
	}
	if len(titleMatches.Products) > 0 {
		return buildCandidatesFromProducts(snapshot, titleMatches.Products, domain.LinkCandidateStateTitleMatch, domain.LinkCandidateMatchInputTitle, snapshot.Title, identityAnchors, now), nil, nil
	}

	unresolved := newCandidate(snapshot, domain.LinkCandidateStateUnresolved, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)
	applyUnresolvedScore(&unresolved, identityAnchors)
	return []domain.LinkCandidate{unresolved}, nil, nil
}

// autoApproval is a candidate the generator judged corroborated, carrying the
// collision count the generator itself read at that moment.
type autoApproval struct {
	Candidate domain.LinkCandidate
	// Pointer so "did not read a count" stays distinguishable from "counted
	// nothing". Only the generator that actually read the anchor's match set
	// may fill it (ADR-17).
	Collisions *int
}

// autoApprovals selects the candidates eligible for the automatic path. ACCEPT
// is reachable only from buildConcordantCandidate — CODPROD and EAN resolving
// the SAME single product with no hard negative — so this reads the generator's
// own verdict instead of re-deriving it. The collision count is likewise the
// generator's own reading (len of the anchor's match set at candidate time),
// never a second count computed here.
func autoApprovals(candidates []domain.LinkCandidate, skuMatches productMatchResult) []autoApproval {
	approvals := make([]autoApproval, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.MatchStatus != domain.LinkCandidateMatchStatusAccept {
			continue
		}
		collisions := len(skuMatches.Products)
		approvals = append(approvals, autoApproval{Candidate: candidate, Collisions: &collisions})
	}
	return approvals
}

type productMatchResult struct {
	InputValue string
	Products   []internalreaddomain.ProductCandidate
	HasSignal  bool
}

func (s *GenerationService) findProducts(ctx context.Context, value string, matchInput domain.LinkCandidateMatchInput) (productMatchResult, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return productMatchResult{}, nil
	}

	input := internalreadports.FindProductsInput{}
	switch matchInput {
	case domain.LinkCandidateMatchInputSellerSKU:
		input.SellerSKU = &trimmed
	case domain.LinkCandidateMatchInputEAN:
		input.EAN = &trimmed
	case domain.LinkCandidateMatchInputTitle:
		input.Title = &trimmed
	default:
		return productMatchResult{}, nil
	}

	products, err := s.matcher.FindProductsForLinking(ctx, input)
	if err != nil {
		return productMatchResult{}, err
	}

	filtered := make([]internalreaddomain.ProductCandidate, 0, len(products))
	for _, product := range products {
		if _, ok := canonicalProductID(product); !ok || internalreaddomain.HasQualityFlag(product.QualityFlags, internalreaddomain.QualityMissingProduct) {
			continue
		}
		filtered = append(filtered, product)
	}

	return productMatchResult{
		InputValue: trimmed,
		Products:   filtered,
		HasSignal:  len(filtered) > 0,
	}, nil
}

func buildExactCandidates(snapshot domain.ListingSnapshot, skuMatches, eanMatches productMatchResult, identityAnchors []ports.ProviderIdentityAnchor, now time.Time) []domain.LinkCandidate {
	if len(skuMatches.Products) == 1 && len(eanMatches.Products) == 1 {
		skuProductID, _ := canonicalProductID(skuMatches.Products[0])
		eanProductID, _ := canonicalProductID(eanMatches.Products[0])
		if skuProductID != eanProductID {
			conflictProducts := uniqueProducts(append(slices.Clone(skuMatches.Products), eanMatches.Products...))
			return buildConflictCandidates(snapshot, conflictProducts, skuMatches, eanMatches, identityAnchors, now)
		}
		// seller_sku + ean concordant on the same codprod: IC-01 A2's viable
		// auto-ACCEPT pair (proxy for "2 independent anchors").
		return []domain.LinkCandidate{buildConcordantCandidate(snapshot, skuMatches, eanMatches, identityAnchors, now)}
	}

	if len(skuMatches.Products) > 1 || len(eanMatches.Products) > 1 {
		return buildCollisionCandidates(snapshot, skuMatches, eanMatches, identityAnchors, now)
	}

	if len(skuMatches.Products) == 1 {
		return buildCandidatesFromProducts(snapshot, skuMatches.Products, domain.LinkCandidateStateExactSKU, domain.LinkCandidateMatchInputSellerSKU, skuMatches.InputValue, identityAnchors, now)
	}
	if len(eanMatches.Products) == 1 {
		return buildCandidatesFromProducts(snapshot, eanMatches.Products, domain.LinkCandidateStateExactEAN, domain.LinkCandidateMatchInputEAN, eanMatches.InputValue, identityAnchors, now)
	}
	return nil
}

func buildConflictCandidates(snapshot domain.ListingSnapshot, products []internalreaddomain.ProductCandidate, skuMatches, eanMatches productMatchResult, identityAnchors []ports.ProviderIdentityAnchor, now time.Time) []domain.LinkCandidate {
	candidates := make([]domain.LinkCandidate, 0, len(products))
	for _, product := range products {
		productID, _ := canonicalProductID(product)
		matchInput := domain.LinkCandidateMatchInputSellerSKU
		matchValue := skuMatches.InputValue
		ownAnchor := "seller_sku"
		conflictAnchor := "ean"
		conflictingIDs := productIDsExcept(eanMatches.Products, productID)
		if !containsProduct(skuMatches.Products, productID) {
			matchInput = domain.LinkCandidateMatchInputEAN
			matchValue = eanMatches.InputValue
			ownAnchor = "ean"
			conflictAnchor = "seller_sku"
			conflictingIDs = productIDsExcept(skuMatches.Products, productID)
		}
		candidate := newCandidate(snapshot, domain.LinkCandidateStateConflict, matchInput, matchValue, product, now)
		applyConflictScore(&candidate, ownAnchor, conflictAnchor, productID, conflictingIDs, identityAnchors)
		candidates = append(candidates, candidate)
	}
	if len(candidates) == 0 {
		candidate := newCandidate(snapshot, domain.LinkCandidateStateConflict, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)
		applyUnresolvedScore(&candidate, identityAnchors)
		candidates = append(candidates, candidate)
	}
	return candidates
}

// buildCollisionCandidates handles an anchor that resolved MORE than one
// product. The anchor decided nothing — it named a set — so every product in
// that set is offered for review and none of them is corroborated (D-121:
// segurança > cobertura). A second anchor that did resolve a single product is
// still offered, but blocked from confirmation: corroboration against an
// ambiguous anchor is not corroboration.
func buildCollisionCandidates(snapshot domain.ListingSnapshot, skuMatches, eanMatches productMatchResult, identityAnchors []ports.ProviderIdentityAnchor, now time.Time) []domain.LinkCandidate {
	anchors := []struct {
		result     productMatchResult
		matchInput domain.LinkCandidateMatchInput
		state      domain.LinkCandidateState
		name       string
		other      string
	}{
		{skuMatches, domain.LinkCandidateMatchInputSellerSKU, domain.LinkCandidateStateExactSKU, "seller_sku", "ean"},
		{eanMatches, domain.LinkCandidateMatchInputEAN, domain.LinkCandidateStateExactEAN, "ean", "seller_sku"},
	}
	candidates := make([]domain.LinkCandidate, 0, len(skuMatches.Products)+len(eanMatches.Products))
	for index, anchor := range anchors {
		otherCount := len(anchors[1-index].result.Products)
		for _, product := range uniqueProducts(anchor.result.Products) {
			productID, _ := canonicalProductID(product)
			candidate := newCandidate(snapshot, domain.LinkCandidateStateConflict, anchor.matchInput, anchor.result.InputValue, product, now)
			if len(anchor.result.Products) > 1 {
				applyCollisionScore(&candidate, anchor.name, len(anchor.result.Products), productIDsExcept(anchor.result.Products, productID), identityAnchors)
			} else {
				applyAmbiguousCorroborationScore(&candidate, anchor.name, productID, anchor.other, otherCount, identityAnchors)
			}
			candidates = append(candidates, candidate)
		}
	}
	if len(candidates) == 0 {
		candidate := newCandidate(snapshot, domain.LinkCandidateStateUnresolved, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)
		applyUnresolvedScore(&candidate, identityAnchors)
		candidates = append(candidates, candidate)
	}
	return candidates
}

func buildCandidatesFromProducts(snapshot domain.ListingSnapshot, products []internalreaddomain.ProductCandidate, state domain.LinkCandidateState, matchInput domain.LinkCandidateMatchInput, matchValue string, identityAnchors []ports.ProviderIdentityAnchor, now time.Time) []domain.LinkCandidate {
	candidates := make([]domain.LinkCandidate, 0, len(products))
	for _, product := range uniqueProducts(products) {
		candidate := newCandidate(snapshot, state, matchInput, matchValue, product, now)
		applySingleAnchorScore(&candidate, snapshot, state, product, identityAnchors)
		candidates = append(candidates, candidate)
	}
	return candidates
}

func newCandidate(snapshot domain.ListingSnapshot, state domain.LinkCandidateState, matchInput domain.LinkCandidateMatchInput, matchValue string, product internalreaddomain.ProductCandidate, now time.Time) domain.LinkCandidate {
	var productID *int
	canonicalID, hasCanonicalID := canonicalProductID(product)
	if hasCanonicalID {
		productID = &canonicalID
	}
	referenceCode := ""
	if product.ReferenceCode != nil {
		referenceCode = strings.TrimSpace(*product.ReferenceCode)
	}
	fetchedAt := snapshot.FetchedAt.UTC()
	candidateID := buildCandidateID(snapshot, state, matchInput, canonicalID)
	return domain.LinkCandidate{
		CandidateID:             candidateID,
		InstallationID:          snapshot.InstallationID,
		ProviderCode:            snapshot.ProviderCode,
		ProviderItemID:          snapshot.ProviderItemID,
		ProviderVariationID:     snapshot.ProviderVariationID,
		InternalProductID:       productID,
		InternalProductName:     strings.TrimSpace(product.Name),
		InternalReferenceCode:   referenceCode,
		State:                   state,
		MatchInput:              matchInput,
		MatchValue:              strings.TrimSpace(matchValue),
		SourceSnapshotFetchedAt: &fetchedAt,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
}

func buildCandidateID(snapshot domain.ListingSnapshot, state domain.LinkCandidateState, matchInput domain.LinkCandidateMatchInput, productID int) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%d",
		strings.TrimSpace(snapshot.InstallationID),
		strings.TrimSpace(snapshot.ProviderItemID),
		strings.TrimSpace(snapshot.ProviderVariationID),
		state,
		matchInput,
		productID,
	)
}

func uniqueProducts(products []internalreaddomain.ProductCandidate) []internalreaddomain.ProductCandidate {
	seen := make(map[int]struct{}, len(products))
	unique := make([]internalreaddomain.ProductCandidate, 0, len(products))
	for _, product := range products {
		productID, ok := canonicalProductID(product)
		if !ok {
			continue
		}
		if _, ok := seen[productID]; ok {
			continue
		}
		seen[productID] = struct{}{}
		unique = append(unique, product)
	}
	return unique
}

func containsProduct(products []internalreaddomain.ProductCandidate, productID int) bool {
	for _, product := range products {
		canonicalID, ok := canonicalProductID(product)
		if ok && canonicalID == productID {
			return true
		}
	}
	return false
}

func canonicalProductID(product internalreaddomain.ProductCandidate) (int, bool) {
	if product.InternalProductID == nil || *product.InternalProductID <= 0 {
		return 0, false
	}
	return int(*product.InternalProductID), true
}

// --- IC-01 Amendment A2 scoring (confidence/band/status/reasons) ---
//
// seller_sku and ean are the only cross-side anchors available against
// provider data (A2); title ranks only, never accepts. Hard-negative title
// lexical checks (kit/combo/cor/voltagem) cap
// any match at BAIXA/REJECT even when EAN/SKU agree ("contradição vence
// EAN").

func buildConcordantCandidate(snapshot domain.ListingSnapshot, skuMatches, eanMatches productMatchResult, identityAnchors []ports.ProviderIdentityAnchor, now time.Time) domain.LinkCandidate {
	product := skuMatches.Products[0]
	candidate := newCandidate(snapshot, domain.LinkCandidateStateExactSKU, domain.LinkCandidateMatchInputSellerSKU, skuMatches.InputValue, product, now)

	reasons := []domain.LinkCandidateReason{
		{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "seller_sku resolve exato para codprod"},
		{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "ean corrobora o mesmo codprod (unproved)"},
	}
	if hardNeg, detail := detectHardNegative(snapshot.Title, product.Name); hardNeg {
		candidate.Confidence = 25
		candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
		candidate.MatchStatus = domain.LinkCandidateMatchStatusReject
		reasons = append(reasons, domain.LinkCandidateReason{Anchor: "title", Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: detail})
	} else {
		candidate.Confidence = 95
		candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandAlta
		candidate.MatchStatus = domain.LinkCandidateMatchStatusAccept
	}
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
	return candidate
}

func applySingleAnchorScore(candidate *domain.LinkCandidate, snapshot domain.ListingSnapshot, state domain.LinkCandidateState, product internalreaddomain.ProductCandidate, identityAnchors []ports.ProviderIdentityAnchor) {
	var reasons []domain.LinkCandidateReason
	var confidence int
	var band domain.LinkCandidateConfidenceBand
	var status domain.LinkCandidateMatchStatus

	switch state {
	case domain.LinkCandidateStateExactSKU:
		// D-121-2: one anchor resolved one product and nothing contradicts it,
		// but nothing corroborates it either — that is the confirmation queue,
		// not an automatic link and not a review. The warning names the anchor
		// that is missing, so the operator knows what they are supplying by
		// confirming.
		confidence, band, status = 70, domain.LinkCandidateConfidenceBandMedia, domain.LinkCandidateMatchStatusConfirm
		eanDetail := "sem EAN para corroborar o CODPROD"
		if strings.TrimSpace(snapshot.EAN) != "" {
			eanDetail = "sem EAN para corroborar o CODPROD: o EAN do anúncio não casa nenhum produto"
		}
		reasons = []domain.LinkCandidateReason{
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "seller_sku resolve exato para codprod"},
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: eanDetail},
		}
	case domain.LinkCandidateStateExactEAN:
		confidence, band, status = 60, domain.LinkCandidateConfidenceBandMedia, domain.LinkCandidateMatchStatusConfirm
		skuDetail := "sem CODPROD para corroborar o EAN"
		if strings.TrimSpace(snapshot.SellerSKU) != "" {
			skuDetail = "sem CODPROD para corroborar o EAN: o seller_sku do anúncio não casa nenhum produto"
		}
		reasons = []domain.LinkCandidateReason{
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "ean corrobora codprod (unproved)"},
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: skuDetail},
		}
	case domain.LinkCandidateStateTitleMatch:
		confidence, band, status = 35, domain.LinkCandidateConfidenceBandBaixa, domain.LinkCandidateMatchStatusReview
		reasons = []domain.LinkCandidateReason{
			{Anchor: "title", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "match por título (ranking-only, nunca ACCEPT)"},
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "seller_sku sem correspondência"},
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
		}
	default:
		applyUnresolvedScore(candidate, identityAnchors)
		return
	}

	if hardNeg, detail := detectHardNegative(snapshot.Title, product.Name); hardNeg {
		confidence, band, status = 25, domain.LinkCandidateConfidenceBandBaixa, domain.LinkCandidateMatchStatusReject
		reasons = append(reasons, domain.LinkCandidateReason{Anchor: "title", Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: detail})
	}

	candidate.Confidence = confidence
	candidate.ConfidenceBand = band
	candidate.MatchStatus = status
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
}

// applyConflictScore scores the CODPROD≠EAN conflict. Neither anchor wins
// (AC-08, ratified by the operator): both products are offered for review with
// the disagreement spelled out, and the decision is the operator's.
func applyConflictScore(candidate *domain.LinkCandidate, ownAnchor, conflictAnchor string, ownProductID int, conflictingIDs []int, identityAnchors []ports.ProviderIdentityAnchor) {
	candidate.Confidence = 20
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusReview
	reasons := []domain.LinkCandidateReason{
		{Anchor: ownAnchor, Direction: domain.LinkCandidateReasonDirectionFor, Detail: fmt.Sprintf("%s aponta codprod %d", ownAnchor, ownProductID)},
		{Anchor: conflictAnchor, Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: fmt.Sprintf("%s aponta codprod %s (conflito, nenhuma âncora vence)", conflictAnchor, formatProductIDs(conflictingIDs))},
	}
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
}

// applyCollisionScore scores one product out of a set an ambiguous anchor
// matched. The anchor named N products, so it identified none of them.
func applyCollisionScore(candidate *domain.LinkCandidate, anchor string, matchCount int, otherIDs []int, identityAnchors []ports.ProviderIdentityAnchor) {
	candidate.Confidence = 20
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusReview
	reasons := []domain.LinkCandidateReason{
		{Anchor: anchor, Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: fmt.Sprintf("%s casa %d produtos no ERP (colisão: também codprod %s) ⇒ âncora ambígua", anchor, matchCount, formatProductIDs(otherIDs))},
	}
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
}

// applyAmbiguousCorroborationScore scores the anchor that DID resolve a single
// product on a listing whose other anchor collided. The single anchor is real,
// but there is nothing to corroborate it against — an ambiguous anchor
// corroborates nothing — so it is review, never confirmation.
func applyAmbiguousCorroborationScore(candidate *domain.LinkCandidate, ownAnchor string, ownProductID int, otherAnchor string, otherCount int, identityAnchors []ports.ProviderIdentityAnchor) {
	candidate.Confidence = 40
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusReview
	reasons := []domain.LinkCandidateReason{
		{Anchor: ownAnchor, Direction: domain.LinkCandidateReasonDirectionFor, Detail: fmt.Sprintf("%s aponta codprod %d", ownAnchor, ownProductID)},
		{Anchor: otherAnchor, Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: fmt.Sprintf("%s casa %d produtos no ERP ⇒ não corrobora", otherAnchor, otherCount)},
	}
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
}

// applyUnresolvedScore scores a listing no anchor could resolve. What M05-C5
// asks for is the REASON: the operator has to see WHY the matcher had nothing
// to offer, never a silent empty row. The status stays NO_CANDIDATE because
// /vinculos keys the ADR-17 affordance off it — the "sem candidato / Criar
// produto / Ignorar" row, and the batch-select guard that stops an operator
// from bulk-approving a listing with no product to approve. Reclassifying it
// as REVIEW would have made those branches dead code on a screen this
// milestone may not edit (apps/web is M-06's seam). Flagged to the hub.
func applyUnresolvedScore(candidate *domain.LinkCandidate, identityAnchors []ports.ProviderIdentityAnchor) {
	candidate.Confidence = 0
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusNoCandidate
	reasons := []domain.LinkCandidateReason{
		{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "seller_sku sem correspondência"},
		{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
	}
	candidate.Reasons = appendProviderDeclaredUnavailableReasons(reasons, identityAnchors)
}

func appendProviderDeclaredUnavailableReasons(reasons []domain.LinkCandidateReason, identityAnchors []ports.ProviderIdentityAnchor) []domain.LinkCandidateReason {
	finalized := make([]domain.LinkCandidateReason, 0, len(reasons)+len(identityAnchors))
	reasonIndexes := make(map[string]int, len(reasons)+len(identityAnchors))
	for _, reason := range reasons {
		index, exists := reasonIndexes[reason.Anchor]
		if !exists {
			reasonIndexes[reason.Anchor] = len(finalized)
			finalized = append(finalized, reason)
			continue
		}
		if finalized[index].Direction == domain.LinkCandidateReasonDirectionUnavailable &&
			reason.Direction != domain.LinkCandidateReasonDirectionUnavailable {
			finalized[index] = reason
		}
	}
	for _, anchor := range identityAnchors {
		if anchor.Supplied {
			continue
		}
		declarationReason := domain.LinkCandidateReason{
			Anchor:    anchor.Anchor,
			Direction: domain.LinkCandidateReasonDirectionUnavailable,
			Detail:    fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor),
		}
		index, exists := reasonIndexes[anchor.Anchor]
		if !exists {
			reasonIndexes[anchor.Anchor] = len(finalized)
			finalized = append(finalized, declarationReason)
			continue
		}
		if finalized[index].Direction == domain.LinkCandidateReasonDirectionUnavailable {
			finalized[index] = declarationReason
		}
	}
	return finalized
}

func productIDsExcept(products []internalreaddomain.ProductCandidate, exclude int) []int {
	ids := make([]int, 0, len(products))
	for _, product := range products {
		id, ok := canonicalProductID(product)
		if ok && id != exclude {
			ids = append(ids, id)
		}
	}
	return ids
}

func formatProductIDs(ids []int) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.Itoa(id))
	}
	return strings.Join(parts, ",")
}

var (
	hardNegativeVoltagePattern = regexp.MustCompile(`(\d+)\s*v\b`)
	hardNegativeColorTokens    = []string{"azul", "verde", "vermelho", "preto", "branco", "amarelo", "cinza", "prata", "dourado", "rosa", "laranja", "roxo"}
	// Dimension tokens: NxM(xK) patterns and numeric measurements with a
	// unit (mm/cm/pol/" and bare metre "m"). Matched case-insensitively
	// against the original title so match offsets remain valid UTF-8 offsets.
	hardNegativeDimensionPattern = regexp.MustCompile(`(?i)\d+\s*x\s*\d+(?:\s*x\s*\d+)?|\d+(?:[.,]\d+)?\s*(?:mm|cm|pol|")|\d+(?:[.,]\d+)?\s*m\b`)
	// Clothing/grade sizes are matched CASE-SENSITIVELY on the original
	// title so the metre abbreviation "m" (lowercase, caught above as a
	// measurement) never collides with the "M" grade size.
	hardNegativeSizePattern = regexp.MustCompile(`\b(?:PP|GG|EG|XG|P|M|G)\b`)
)

// detectHardNegative compares the provider listing title against the
// internal product name for lexical contradictions (kit/combo, cor,
// medida/dimensão, voltagem). Any hit is a hard negative: contradição
// vence EAN, so a candidate is capped at BAIXA/REJECT even when its
// anchors otherwise agree.
func detectHardNegative(snapshotTitle, internalName string) (bool, string) {
	stOriginal := strings.TrimSpace(snapshotTitle)
	inOriginal := strings.TrimSpace(internalName)
	st := strings.ToLower(stOriginal)
	in := strings.ToLower(inOriginal)
	if st == "" || in == "" {
		return false, ""
	}

	stKit := strings.Contains(st, "kit") || strings.Contains(st, "combo")
	inKit := strings.Contains(in, "kit") || strings.Contains(in, "combo")
	if stKit != inKit {
		return true, "hard-negative: kit/combo divergente entre título do anúncio e produto interno"
	}

	if stVolt, ok := hardNegativeVoltage(st); ok {
		if inVolt, ok2 := hardNegativeVoltage(in); ok2 && stVolt != inVolt {
			return true, fmt.Sprintf("hard-negative: voltagem divergente %sV≠%sV", stVolt, inVolt)
		}
	}

	if stDim, stDisplay, ok := hardNegativeDimension(stOriginal); ok {
		if inDim, inDisplay, ok2 := hardNegativeDimension(inOriginal); ok2 && stDim != inDim {
			return true, fmt.Sprintf("hard-negative: medida/dimensão divergente %s≠%s", stDisplay, inDisplay)
		}
	}

	if stColor, ok := hardNegativeColor(st); ok {
		if inColor, ok2 := hardNegativeColor(in); ok2 && stColor != inColor {
			return true, fmt.Sprintf("hard-negative: cor divergente %s≠%s", stColor, inColor)
		}
	}

	return false, ""
}

func hardNegativeVoltage(text string) (string, bool) {
	match := hardNegativeVoltagePattern.FindStringSubmatch(text)
	if match == nil {
		return "", false
	}
	return match[1], true
}

// hardNegativeDimension builds a normalized, order-independent signature of
// every dimension/measurement/size token in a title. Two titles are a
// dimension contradiction only when BOTH carry dimension tokens and the
// signatures differ — normal titles without measurements are never flagged
// (symmetrical to the voltage/color checks, no false positives).
func hardNegativeDimension(original string) (string, string, bool) {
	type dimensionPair struct {
		key     string
		display string
	}

	tokenIndexes := hardNegativeDimensionPattern.FindAllStringIndex(original, -1)
	pairs := make([]dimensionPair, 0, len(tokenIndexes))
	for _, index := range tokenIndexes {
		originalToken := strings.TrimSpace(original[index[0]:index[1]])
		key, _ := normalizeDimensionToken(originalToken)
		pairs = append(pairs, dimensionPair{key: key, display: originalToken})
	}
	for _, size := range hardNegativeSizePattern.FindAllString(original, -1) {
		pairs = append(pairs, dimensionPair{key: strings.ToLower(size), display: size})
	}
	if len(pairs) == 0 {
		return "", "", false
	}
	slices.SortFunc(pairs, func(a, b dimensionPair) int {
		return strings.Compare(a.key, b.key)
	})
	pairs = slices.CompactFunc(pairs, func(a, b dimensionPair) bool {
		return a.key == b.key
	})
	keys := make([]string, 0, len(pairs))
	display := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		keys = append(keys, pair.key)
		display = append(display, pair.display)
	}
	return strings.Join(keys, "|"), strings.Join(display, "|"), true
}

func normalizeDimensionToken(originalToken string) (string, bool) {
	token := strings.ToLower(strings.TrimSpace(originalToken))
	token = strings.ReplaceAll(token, " ", "")
	token = strings.ReplaceAll(token, ",", ".")
	if strings.Contains(token, "x") {
		return token, true
	}

	unit := ""
	for _, candidate := range []string{"pol", "mm", "cm", `"`, "m"} {
		if strings.HasSuffix(token, candidate) {
			unit = candidate
			break
		}
	}
	value, ok := new(big.Rat).SetString(strings.TrimSuffix(token, unit))
	if !ok {
		return token, false
	}
	switch unit {
	case "cm":
		value.Mul(value, big.NewRat(10, 1))
	case "m":
		value.Mul(value, big.NewRat(1000, 1))
	case "pol", `"`:
		value.Mul(value, big.NewRat(127, 5))
	}
	return "mm:" + value.RatString(), true
}

func hardNegativeColor(text string) (string, bool) {
	for _, token := range hardNegativeColorTokens {
		if strings.Contains(text, token) {
			return token, true
		}
	}
	return "", false
}

func classifyMatcherError(err error) error {
	switch {
	case internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable):
		return errors.New("PRODUCT_LINKS_INTERNAL_READ_UNAVAILABLE")
	case internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorUnsupportedQuery):
		return errors.New("PRODUCT_LINKS_MATCH_QUERY_UNSUPPORTED")
	default:
		return err
	}
}
