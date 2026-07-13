package application

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	snapshots, err := s.snapshots.ListListingSnapshots(ctx, installationID, limit)
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

	return []domain.LinkCandidate{newCandidate(snapshot, domain.LinkCandidateStateUnresolved, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now)}, nil
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
		if !containsProduct(skuMatches.Products, productID) {
			matchInput = domain.LinkCandidateMatchInputEAN
			matchValue = eanMatches.InputValue
		}
		candidates = append(candidates, newCandidate(snapshot, domain.LinkCandidateStateConflict, matchInput, matchValue, product, now))
	}
	if len(candidates) == 0 {
		candidates = append(candidates, newCandidate(snapshot, domain.LinkCandidateStateConflict, domain.LinkCandidateMatchInputNone, "", internalreaddomain.ProductCandidate{}, now))
	}
	return candidates
}

func buildCandidatesFromProducts(snapshot domain.ListingSnapshot, products []internalreaddomain.ProductCandidate, state domain.LinkCandidateState, matchInput domain.LinkCandidateMatchInput, matchValue string, now time.Time) []domain.LinkCandidate {
	candidates := make([]domain.LinkCandidate, 0, len(products))
	for _, product := range uniqueProducts(products) {
		candidates = append(candidates, newCandidate(snapshot, state, matchInput, matchValue, product, now))
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
