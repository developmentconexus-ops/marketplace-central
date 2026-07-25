package application

import (
	"context"
	"errors"
	"fmt"
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

type GenerationService struct {
	snapshots ports.ListingSnapshotReader
	matcher   ProductMatcher
	store     ports.LinkCandidateStore
	now       func() time.Time
}

type GenerationServiceConfig struct {
	Snapshots ports.ListingSnapshotReader
	Matcher   ProductMatcher
	Store     ports.LinkCandidateStore
	Now       func() time.Time
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
		snapshots: cfg.Snapshots,
		matcher:   cfg.Matcher,
		store:     cfg.Store,
		now:       now,
	}
}

func (s *GenerationService) GenerateLinkCandidates(ctx context.Context, input GenerateLinkCandidatesInput) (domain.LinkCandidateGenerationResult, error) {
	if s.snapshots == nil || s.matcher == nil || s.store == nil {
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

	identities := make([]domain.ListingIdentity, 0, len(snapshots))
	candidates := make([]domain.LinkCandidate, 0, len(snapshots))
	for _, snapshot := range snapshots {
		identities = append(identities, domain.ListingIdentity{
			InstallationID:      installationID,
			ProviderItemID:      snapshot.ProviderItemID,
			ProviderVariationID: snapshot.ProviderVariationID,
		})

		generated, err := s.generateForSnapshot(ctx, snapshot)
		if err != nil {
			return domain.LinkCandidateGenerationResult{}, classifyMatcherError(err)
		}
		candidates = append(candidates, generated...)
	}

	if err := s.store.ReplaceLinkCandidates(ctx, installationID, identities, candidates); err != nil {
		return domain.LinkCandidateGenerationResult{}, err
	}

	return domain.LinkCandidateGenerationResult{
		InstallationID: installationID,
		GeneratedCount: len(candidates),
		Items:          candidates,
	}, nil
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

func (s *GenerationService) generateForSnapshot(ctx context.Context, snapshot domain.ListingSnapshot) ([]domain.LinkCandidate, error) {
	now := s.now().UTC()

	skuMatches, err := s.findProducts(ctx, snapshot.SellerSKU, domain.LinkCandidateMatchInputSellerSKU)
	if err != nil {
		return nil, err
	}
	eanMatches, err := s.findProducts(ctx, snapshot.EAN, domain.LinkCandidateMatchInputEAN)
	if err != nil {
		return nil, err
	}

	exactCandidates := buildExactCandidates(snapshot, skuMatches, eanMatches, now)
	if len(exactCandidates) > 0 {
		return exactCandidates, nil
	}

	titleMatches, err := s.findProducts(ctx, snapshot.Title, domain.LinkCandidateMatchInputTitle)
	if err != nil {
		return nil, err
	}
	if len(titleMatches.Products) > 0 {
		return buildCandidatesFromProducts(snapshot, titleMatches.Products, domain.LinkCandidateStateTitleMatch, domain.LinkCandidateMatchInputTitle, snapshot.Title, now), nil
	}

	unresolved := newCandidate(snapshot, domain.LinkCandidateStateUnresolved, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)
	applyNoCandidateScore(&unresolved)
	return []domain.LinkCandidate{unresolved}, nil
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

func buildExactCandidates(snapshot domain.ListingSnapshot, skuMatches, eanMatches productMatchResult, now time.Time) []domain.LinkCandidate {
	if len(skuMatches.Products) == 1 && len(eanMatches.Products) == 1 {
		skuProductID, _ := canonicalProductID(skuMatches.Products[0])
		eanProductID, _ := canonicalProductID(eanMatches.Products[0])
		if skuProductID != eanProductID {
			conflictProducts := uniqueProducts(append(slices.Clone(skuMatches.Products), eanMatches.Products...))
			return buildConflictCandidates(snapshot, conflictProducts, skuMatches, eanMatches, now)
		}
		// seller_sku + ean concordant on the same codprod: IC-01 A2's viable
		// auto-ACCEPT pair (proxy for "2 independent anchors").
		return []domain.LinkCandidate{buildConcordantCandidate(snapshot, skuMatches, eanMatches, now)}
	}

	if len(skuMatches.Products) > 1 || len(eanMatches.Products) > 1 {
		conflictProducts := uniqueProducts(append(slices.Clone(skuMatches.Products), eanMatches.Products...))
		return buildConflictCandidates(snapshot, conflictProducts, skuMatches, eanMatches, now)
	}

	if len(skuMatches.Products) == 1 {
		return buildCandidatesFromProducts(snapshot, skuMatches.Products, domain.LinkCandidateStateExactSKU, domain.LinkCandidateMatchInputSellerSKU, skuMatches.InputValue, now)
	}
	if len(eanMatches.Products) == 1 {
		return buildCandidatesFromProducts(snapshot, eanMatches.Products, domain.LinkCandidateStateExactEAN, domain.LinkCandidateMatchInputEAN, eanMatches.InputValue, now)
	}
	return nil
}

func buildConflictCandidates(snapshot domain.ListingSnapshot, products []internalreaddomain.ProductCandidate, skuMatches, eanMatches productMatchResult, now time.Time) []domain.LinkCandidate {
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
		applyConflictScore(&candidate, ownAnchor, conflictAnchor, productID, conflictingIDs)
		candidates = append(candidates, candidate)
	}
	if len(candidates) == 0 {
		candidate := newCandidate(snapshot, domain.LinkCandidateStateConflict, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)
		applyNoCandidateScore(&candidate)
		candidates = append(candidates, candidate)
	}
	return candidates
}

func buildCandidatesFromProducts(snapshot domain.ListingSnapshot, products []internalreaddomain.ProductCandidate, state domain.LinkCandidateState, matchInput domain.LinkCandidateMatchInput, matchValue string, now time.Time) []domain.LinkCandidate {
	candidates := make([]domain.LinkCandidate, 0, len(products))
	for _, product := range uniqueProducts(products) {
		candidate := newCandidate(snapshot, state, matchInput, matchValue, product, now)
		applySingleAnchorScore(&candidate, snapshot, state, product)
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
// provider data (A2); title ranks only, never accepts. marca/refforn exist
// solely on the internal side and are therefore ALWAYS surfaced as
// UNAVAILABLE reasons (ADR-17 — motivo sempre visível, never silently
// dropped). Hard-negative title lexical checks (kit/combo/cor/voltagem) cap
// any match at BAIXA/REJECT even when EAN/SKU agree ("contradição vence
// EAN").

func buildConcordantCandidate(snapshot domain.ListingSnapshot, skuMatches, eanMatches productMatchResult, now time.Time) domain.LinkCandidate {
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
	candidate.Reasons = append(reasons, mandatoryUnavailableReasons()...)
	return candidate
}

func applySingleAnchorScore(candidate *domain.LinkCandidate, snapshot domain.ListingSnapshot, state domain.LinkCandidateState, product internalreaddomain.ProductCandidate) {
	var reasons []domain.LinkCandidateReason
	var confidence int
	var band domain.LinkCandidateConfidenceBand
	var status domain.LinkCandidateMatchStatus

	switch state {
	case domain.LinkCandidateStateExactSKU:
		confidence, band, status = 70, domain.LinkCandidateConfidenceBandMedia, domain.LinkCandidateMatchStatusReview
		eanDetail := "EAN ausente ⇒ máximo REVIEW"
		if strings.TrimSpace(snapshot.EAN) != "" {
			eanDetail = "EAN presente mas sem correspondência ⇒ máximo REVIEW"
		}
		reasons = []domain.LinkCandidateReason{
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "seller_sku resolve exato para codprod"},
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: eanDetail},
		}
	case domain.LinkCandidateStateExactEAN:
		confidence, band, status = 60, domain.LinkCandidateConfidenceBandMedia, domain.LinkCandidateMatchStatusReview
		reasons = []domain.LinkCandidateReason{
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "ean corrobora codprod (unproved)"},
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "seller_sku sem correspondência"},
		}
	case domain.LinkCandidateStateTitleMatch:
		confidence, band, status = 35, domain.LinkCandidateConfidenceBandBaixa, domain.LinkCandidateMatchStatusReview
		reasons = []domain.LinkCandidateReason{
			{Anchor: "title", Direction: domain.LinkCandidateReasonDirectionFor, Detail: "match por título (ranking-only, nunca ACCEPT)"},
			{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "seller_sku sem correspondência"},
			{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
		}
	default:
		applyNoCandidateScore(candidate)
		return
	}

	if hardNeg, detail := detectHardNegative(snapshot.Title, product.Name); hardNeg {
		confidence, band, status = 25, domain.LinkCandidateConfidenceBandBaixa, domain.LinkCandidateMatchStatusReject
		reasons = append(reasons, domain.LinkCandidateReason{Anchor: "title", Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: detail})
	}

	candidate.Confidence = confidence
	candidate.ConfidenceBand = band
	candidate.MatchStatus = status
	candidate.Reasons = append(reasons, mandatoryUnavailableReasons()...)
}

func applyConflictScore(candidate *domain.LinkCandidate, ownAnchor, conflictAnchor string, ownProductID int, conflictingIDs []int) {
	candidate.Confidence = 20
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusReject
	reasons := []domain.LinkCandidateReason{
		{Anchor: ownAnchor, Direction: domain.LinkCandidateReasonDirectionFor, Detail: fmt.Sprintf("%s aponta codprod %d", ownAnchor, ownProductID)},
		{Anchor: conflictAnchor, Direction: domain.LinkCandidateReasonDirectionAgainst, Detail: fmt.Sprintf("%s aponta codprod %s (conflito)", conflictAnchor, formatProductIDs(conflictingIDs))},
	}
	candidate.Reasons = append(reasons, mandatoryUnavailableReasons()...)
}

func applyNoCandidateScore(candidate *domain.LinkCandidate) {
	candidate.Confidence = 0
	candidate.ConfidenceBand = domain.LinkCandidateConfidenceBandBaixa
	candidate.MatchStatus = domain.LinkCandidateMatchStatusNoCandidate
	reasons := []domain.LinkCandidateReason{
		{Anchor: "seller_sku", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "seller_sku sem correspondência"},
		{Anchor: "ean", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
	}
	candidate.Reasons = append(reasons, mandatoryUnavailableReasons()...)
}

// mandatoryUnavailableReasons always surfaces marca/refforn as UNAVAILABLE:
// they exist only on the internal side and can never be computed against
// provider data, but ADR-17 forbids silently omitting the fact (never
// silent zero/default).
func mandatoryUnavailableReasons() []domain.LinkCandidateReason {
	return []domain.LinkCandidateReason{
		{Anchor: "marca", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "marca inexistente no lado provider"},
		{Anchor: "refforn", Direction: domain.LinkCandidateReasonDirectionUnavailable, Detail: "refforn inexistente no lado provider"},
	}
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
	// unit (mm/cm/pol/" and bare metre "m"). Matched on the lowercased
	// title.
	hardNegativeDimensionPattern = regexp.MustCompile(`\d+\s*x\s*\d+(?:\s*x\s*\d+)?|\d+(?:[.,]\d+)?\s*(?:mm|cm|pol|")|\d+(?:[.,]\d+)?\s*m\b`)
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

	if stDim, ok := hardNegativeDimension(st, stOriginal); ok {
		if inDim, ok2 := hardNegativeDimension(in, inOriginal); ok2 && stDim != inDim {
			return true, fmt.Sprintf("hard-negative: medida/dimensão divergente %s≠%s", stDim, inDim)
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
func hardNegativeDimension(lowered, original string) (string, bool) {
	tokens := hardNegativeDimensionPattern.FindAllString(lowered, -1)
	normalized := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.ReplaceAll(token, " ", "")
		token = strings.ReplaceAll(token, ",", ".")
		normalized = append(normalized, token)
	}
	for _, size := range hardNegativeSizePattern.FindAllString(original, -1) {
		normalized = append(normalized, strings.ToLower(size))
	}
	if len(normalized) == 0 {
		return "", false
	}
	slices.Sort(normalized)
	return strings.Join(slices.Compact(normalized), "|"), true
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
