package application

import (
	"context"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubSnapshotReader struct {
	installationID string
	limit          int
	snapshots      []productlinksdomain.ListingSnapshot
}

func (s *stubSnapshotReader) ListListingSnapshots(_ context.Context, installationID string, limit int) ([]productlinksdomain.ListingSnapshot, error) {
	s.installationID = installationID
	s.limit = limit
	return s.snapshots, nil
}

type stubProductMatcher struct {
	results map[string][]internalreaddomain.ProductCandidate
}

func canonicalIDPtr(id int) *internalreaddomain.InternalProductID {
	canonicalID := internalreaddomain.InternalProductID(id)
	return &canonicalID
}

func (s *stubProductMatcher) FindProductsForLinking(_ context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	switch {
	case input.SellerSKU != nil:
		return s.results["sku:"+*input.SellerSKU], nil
	case input.EAN != nil:
		return s.results["ean:"+*input.EAN], nil
	case input.Title != nil:
		return s.results["title:"+*input.Title], nil
	default:
		return nil, nil
	}
}

type stubCandidateStore struct {
	installationID string
	identities     []productlinksdomain.ListingIdentity
	candidates     []productlinksdomain.LinkCandidate
}

func (s *stubCandidateStore) ReplaceLinkCandidates(_ context.Context, installationID string, identities []productlinksdomain.ListingIdentity, candidates []productlinksdomain.LinkCandidate) error {
	s.installationID = installationID
	s.identities = append([]productlinksdomain.ListingIdentity(nil), identities...)
	s.candidates = append([]productlinksdomain.LinkCandidate(nil), candidates...)
	return nil
}

func (s *stubCandidateStore) ListLinkCandidates(_ context.Context, _ string, _ int) ([]productlinksdomain.LinkCandidate, error) {
	return append([]productlinksdomain.LinkCandidate(nil), s.candidates...), nil
}

func (s *stubCandidateStore) GetLinkCandidate(_ context.Context, candidateID string) (productlinksdomain.LinkCandidate, bool, error) {
	for _, candidate := range s.candidates {
		if candidate.CandidateID == candidateID {
			return candidate, true, nil
		}
	}
	return productlinksdomain.LinkCandidate{}, false, nil
}

func TestGenerateLinkCandidatesUsesExactSKUFirst(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 20, 0, 0, 0, time.UTC)
	snapshots := &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB1",
		SellerSKU:      "SKU-1",
		EAN:            "789",
		Title:          "Produto A",
		FetchedAt:      now.Add(-time.Minute),
	}}}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:SKU-1": {{InternalProductID: canonicalIDPtr(101), ProductID: 9001, Name: "Produto Interno", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
		"ean:789":   {{InternalProductID: canonicalIDPtr(101), ProductID: 9002, Name: "Produto Interno", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
	}}
	store := &stubCandidateStore{}
	svc := NewGenerationService(GenerationServiceConfig{Snapshots: snapshots, Matcher: matcher, Store: store, Now: func() time.Time { return now }})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1", Limit: 5})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}

	if snapshots.installationID != "inst-1" || snapshots.limit != 5 {
		t.Fatalf("snapshots called with installation=%q limit=%d", snapshots.installationID, snapshots.limit)
	}
	if result.GeneratedCount != 1 || len(store.candidates) != 1 {
		t.Fatalf("generated=%d stored=%d, want 1", result.GeneratedCount, len(store.candidates))
	}
	if got := store.candidates[0].State; got != productlinksdomain.LinkCandidateStateExactSKU {
		t.Fatalf("state=%s, want exact_sku", got)
	}
	if got := store.candidates[0].InternalProductID; got == nil || *got != 101 {
		t.Fatalf("internal product id=%v, want canonical 101", got)
	}
}

func TestGenerateLinkCandidatesFallsBackToExactEAN(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 20, 5, 0, 0, time.UTC)
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-1",
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB2",
			EAN:            "999",
			Title:          "Produto B",
			FetchedAt:      now.Add(-time.Minute),
		}}},
		Matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
			"ean:999": {{InternalProductID: canonicalIDPtr(202), ProductID: 9202, Name: "Produto EAN", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
		}},
		Store: &stubCandidateStore{},
		Now:   func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if result.Items[0].State != productlinksdomain.LinkCandidateStateExactEAN {
		t.Fatalf("state=%s, want exact_ean", result.Items[0].State)
	}
}

func TestGenerateLinkCandidatesProducesConflictWhenExactSignalsDisagree(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 20, 10, 0, 0, time.UTC)
	store := &stubCandidateStore{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-1",
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB3",
			SellerSKU:      "SKU-X",
			EAN:            "EAN-X",
			Title:          "Produto C",
			FetchedAt:      now.Add(-time.Minute),
		}}},
		Matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
			"sku:SKU-X": {{InternalProductID: canonicalIDPtr(301), ProductID: 999, Name: "Produto SKU", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
			"ean:EAN-X": {{InternalProductID: canonicalIDPtr(302), ProductID: 999, Name: "Produto EAN", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
		}},
		Store: store,
		Now:   func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if result.GeneratedCount != 2 {
		t.Fatalf("generated=%d, want 2 conflict candidates", result.GeneratedCount)
	}
	for _, candidate := range store.candidates {
		if candidate.State != productlinksdomain.LinkCandidateStateConflict {
			t.Fatalf("candidate=%#v, want conflict", candidate)
		}
		if candidate.InternalProductID == nil || (*candidate.InternalProductID != 301 && *candidate.InternalProductID != 302) {
			t.Fatalf("candidate=%#v, want canonical conflict identity", candidate)
		}
	}
}

func TestGenerateLinkCandidatesFallsBackToTitleMatches(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 20, 15, 0, 0, time.UTC)
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-1",
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB4",
			Title:          "Produto D",
			FetchedAt:      now.Add(-time.Minute),
		}}},
		Matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
			"title:Produto D": {
				{InternalProductID: canonicalIDPtr(401), ProductID: 9401, Name: "Produto D 1", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
				{InternalProductID: canonicalIDPtr(402), ProductID: 9402, Name: "Produto D 2", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
			},
		}},
		Store: &stubCandidateStore{},
		Now:   func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if result.GeneratedCount != 2 {
		t.Fatalf("generated=%d, want 2 title matches", result.GeneratedCount)
	}
	for _, candidate := range result.Items {
		if candidate.State != productlinksdomain.LinkCandidateStateTitleMatch {
			t.Fatalf("candidate=%#v, want title_match", candidate)
		}
	}
}

func TestGenerateLinkCandidatesProducesUnresolvedWhenNothingMatches(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 20, 20, 0, 0, time.UTC)
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-1",
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB5",
			Title:          "Sem Match",
			FetchedAt:      now.Add(-time.Minute),
		}}},
		Matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{}},
		Store:   &stubCandidateStore{},
		Now:     func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if result.GeneratedCount != 1 || result.Items[0].State != productlinksdomain.LinkCandidateStateUnresolved {
		t.Fatalf("result=%#v, want one unresolved candidate", result)
	}
}

func TestGenerateLinkCandidatesRejectsLegacyOnlyAndInvalidCanonicalIDs(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		canonicalID *internalreaddomain.InternalProductID
	}{
		{name: "nil canonical ID"},
		{name: "zero canonical ID", canonicalID: canonicalIDPtr(0)},
		{name: "negative canonical ID", canonicalID: canonicalIDPtr(-1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
			store := &stubCandidateStore{}
			svc := NewGenerationService(GenerationServiceConfig{
				Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
					InstallationID: "inst-legacy", ProviderCode: "mercado_livre", ProviderItemID: "MLB-LEGACY",
					SellerSKU: "LEGACY-ONLY", Title: "Legacy only", FetchedAt: now,
				}}},
				Matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
					"sku:LEGACY-ONLY": {{InternalProductID: tc.canonicalID, ProductID: 777, Name: "Legacy product", QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}}},
				}},
				Store: store,
				Now:   func() time.Time { return now },
			})

			result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-legacy"})
			if err != nil {
				t.Fatalf("GenerateLinkCandidates() error = %v", err)
			}
			if result.GeneratedCount != 1 || len(store.candidates) != 1 {
				t.Fatalf("generated=%d stored=%d, want one unresolved candidate", result.GeneratedCount, len(store.candidates))
			}
			candidate := store.candidates[0]
			if candidate.State != productlinksdomain.LinkCandidateStateUnresolved || candidate.InternalProductID != nil {
				t.Fatalf("candidate=%#v, want unresolved without canonical identity", candidate)
			}
		})
	}
}

func TestCanonicalIdentityControlsDeduplicationAndCandidateID(t *testing.T) {
	t.Parallel()

	snapshot := productlinksdomain.ListingSnapshot{InstallationID: "inst-1", ProviderItemID: "MLB-STABLE"}
	now := time.Date(2026, 7, 13, 12, 30, 0, 0, time.UTC)
	products := []internalreaddomain.ProductCandidate{
		{InternalProductID: canonicalIDPtr(501), ProductID: 1, Name: "First"},
		{InternalProductID: canonicalIDPtr(501), ProductID: 2, Name: "Duplicate canonical identity"},
	}

	candidates := buildCandidatesFromProducts(snapshot, products, productlinksdomain.LinkCandidateStateTitleMatch, productlinksdomain.LinkCandidateMatchInputTitle, "stable", now)
	if len(candidates) != 1 || candidates[0].InternalProductID == nil || *candidates[0].InternalProductID != 501 {
		t.Fatalf("candidates=%#v, want one canonical candidate 501", candidates)
	}
	firstID := newCandidate(snapshot, productlinksdomain.LinkCandidateStateExactSKU, productlinksdomain.LinkCandidateMatchInputSellerSKU, "SKU", products[0], now).CandidateID
	secondID := newCandidate(snapshot, productlinksdomain.LinkCandidateStateExactSKU, productlinksdomain.LinkCandidateMatchInputSellerSKU, "SKU", products[1], now).CandidateID
	if firstID != secondID {
		t.Fatalf("candidate IDs differ by legacy metadata: %q != %q", firstID, secondID)
	}
}
