package application

import (
	"context"
	"strings"
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

func TestGenerateLinkCandidatesWithoutLimitReadsEverySnapshot(t *testing.T) {
	t.Parallel()
	// A capped generation run leaves every listing beyond the cap without a
	// candidate, which the operator reads as "sem vínculo" — so an absent limit
	// asks the reader for the whole installation (0), not a 20-row page.
	snapshots := &stubSnapshotReader{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: snapshots,
		Matcher:   &stubProductMatcher{},
		Store:     &stubCandidateStore{},
	})

	if _, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"}); err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if snapshots.limit != 0 {
		t.Fatalf("snapshot reader limit = %d, want 0 (unbounded)", snapshots.limit)
	}
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

// --- IC-01 Amendment A2 fixtures (PLAN-M04 §4) ---
//
// Each fixture drives GenerateLinkCandidates end-to-end and asserts the
// deterministic confidence_band + match_status + key reasons[] entries from
// the anchor model. marca/refforn must ALWAYS appear as UNAVAILABLE
// (ADR-17) regardless of band/status.

func findReason(reasons []productlinksdomain.LinkCandidateReason, anchor string, direction productlinksdomain.LinkCandidateReasonDirection) (productlinksdomain.LinkCandidateReason, bool) {
	for _, reason := range reasons {
		if reason.Anchor == anchor && reason.Direction == direction {
			return reason, true
		}
	}
	return productlinksdomain.LinkCandidateReason{}, false
}

func assertMandatoryUnavailableReasons(t *testing.T, reasons []productlinksdomain.LinkCandidateReason) {
	t.Helper()
	if _, ok := findReason(reasons, "marca", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want marca UNAVAILABLE (ADR-17, always emitted)", reasons)
	}
	if _, ok := findReason(reasons, "refforn", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want refforn UNAVAILABLE (ADR-17, always emitted)", reasons)
	}
}

func generateSingle(t *testing.T, snapshot productlinksdomain.ListingSnapshot, matcher *stubProductMatcher, now time.Time) productlinksdomain.LinkCandidate {
	t.Helper()
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
		Matcher:   matcher,
		Store:     &stubCandidateStore{},
		Now:       func() time.Time { return now },
	})
	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: snapshot.InstallationID})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("result.Items=%#v, want exactly one candidate", result.Items)
	}
	return result.Items[0]
}

// Case 1 — CONCORDANT-ALTA: seller_sku + ean agree on the same codprod, no
// hard-negative. Also the primary auto-ACCEPT proxy.
func TestCase1ConcordantSKUAndEANYieldsAltaAccept(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 0, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx1", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX1",
		SellerSKU: "MLB-SKU-1", EAN: "7891234567895", Title: "Furadeira Bosch 550W", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:MLB-SKU-1":       {{InternalProductID: canonicalIDPtr(100), Name: "Furadeira Bosch 550W"}},
		"ean:7891234567895":   {{InternalProductID: canonicalIDPtr(100), Name: "Furadeira Bosch 550W"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandAlta {
		t.Fatalf("confidence_band=%s, want ALTA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusAccept {
		t.Fatalf("match_status=%s, want ACCEPT", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want seller_sku FOR", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 2 — SKU-ALONE-MEDIA: seller_sku matches, EAN absent from the
// snapshot ⇒ máximo REVIEW (EAN-absent binding, never auto-ACCEPT).
func TestCase2SellerSKUAloneWithoutEANYieldsMediaReview(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 5, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx2", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX2",
		SellerSKU: "MLB-SKU-2", EAN: "", Title: "Parafusadeira 12V", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:MLB-SKU-2": {{InternalProductID: canonicalIDPtr(200), Name: "Parafusadeira 12V"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandMedia {
		t.Fatalf("confidence_band=%s, want MEDIA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
		t.Fatalf("match_status=%s, want REVIEW", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want seller_sku FOR", candidate.Reasons)
	}
	eanReason, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok {
		t.Fatalf("reasons=%#v, want ean UNAVAILABLE", candidate.Reasons)
	}
	if !strings.Contains(eanReason.Detail, "ausente") {
		t.Fatalf("ean reason detail=%q, want it to explain EAN ausente ⇒ máximo REVIEW", eanReason.Detail)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 3 — EAN-ALONE-MEDIA: seller_sku has no match, EAN corroborates
// (unproved) a single codprod.
func TestCase3EANAloneYieldsMediaReview(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 10, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx3", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX3",
		SellerSKU: "SKU-NO-MATCH", EAN: "7890000000017", Title: "Chave de Fenda", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"ean:7890000000017": {{InternalProductID: canonicalIDPtr(300), Name: "Chave de Fenda"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandMedia {
		t.Fatalf("confidence_band=%s, want MEDIA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
		t.Fatalf("match_status=%s, want REVIEW", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR (unproved)", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want seller_sku UNAVAILABLE", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 4 — TITLE-ONLY-BAIXA: seller_sku/ean have no match; title is
// ranking-only and can never grant ACCEPT/REVIEW-grade confidence.
func TestCase4TitleOnlyYieldsBaixaReview(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 15, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx4", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX4",
		Title: "Martelo de Borracha", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"title:Martelo de Borracha": {{InternalProductID: canonicalIDPtr(400), Name: "Martelo de Borracha"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
		t.Fatalf("confidence_band=%s, want BAIXA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
		t.Fatalf("match_status=%s, want REVIEW", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want title FOR (ranking-only)", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want seller_sku UNAVAILABLE", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want ean UNAVAILABLE", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 5 — SKU-EAN-CONFLICT-REJECT: seller_sku and ean point at different
// codprod. Both candidates are BAIXA/REJECT; each cites its own anchor FOR
// and the other anchor AGAINST.
func TestCase5SKUEANConflictYieldsBaixaRejectBothSides(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 20, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx5", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX5",
		SellerSKU: "MLB-SKU-5", EAN: "7899999999994", Title: "Serra Circular", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:MLB-SKU-5":     {{InternalProductID: canonicalIDPtr(500), Name: "Serra Circular 500"}},
		"ean:7899999999994": {{InternalProductID: canonicalIDPtr(501), Name: "Serra Circular 501"}},
	}}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
		Matcher:   matcher,
		Store:     &stubCandidateStore{},
		Now:       func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: snapshot.InstallationID})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("result.Items=%#v, want two conflict candidates", result.Items)
	}
	for _, candidate := range result.Items {
		if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
			t.Fatalf("candidate=%#v, want BAIXA", candidate)
		}
		if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject {
			t.Fatalf("candidate=%#v, want REJECT", candidate)
		}
		assertMandatoryUnavailableReasons(t, candidate.Reasons)
		switch *candidate.InternalProductID {
		case 500:
			if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
				t.Fatalf("candidate 500 reasons=%#v, want seller_sku FOR", candidate.Reasons)
			}
			against, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionAgainst)
			if !ok || !strings.Contains(against.Detail, "501") {
				t.Fatalf("candidate 500 reasons=%#v, want ean AGAINST citing codprod 501", candidate.Reasons)
			}
		case 501:
			if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
				t.Fatalf("candidate 501 reasons=%#v, want ean FOR", candidate.Reasons)
			}
			against, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionAgainst)
			if !ok || !strings.Contains(against.Detail, "500") {
				t.Fatalf("candidate 501 reasons=%#v, want seller_sku AGAINST citing codprod 500", candidate.Reasons)
			}
		default:
			t.Fatalf("unexpected candidate product id: %#v", candidate)
		}
	}
}

// Case 6 — DOKA-HARD-NEGATIVE-REJECT: EAN matches, but the listing title
// signals a kit/combo the internal product is not ⇒ contradição vence EAN.
func TestCase6DokaKitHardNegativeCapsBaixaReject(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 25, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx6", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX6",
		EAN: "7896541230004", Title: "Andaime Doka Kit 12 peças", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"ean:7896541230004": {{InternalProductID: canonicalIDPtr(600), Name: "Escora Menegotti unidade"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
		t.Fatalf("confidence_band=%s, want BAIXA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject {
		t.Fatalf("match_status=%s, want REJECT", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR", candidate.Reasons)
	}
	titleAgainst, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst)
	if !ok || !strings.Contains(titleAgainst.Detail, "kit/combo") {
		t.Fatalf("reasons=%#v, want title AGAINST citing kit/combo divergence", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 7 — VOLTAGE-HARD-NEGATIVE-REJECT: seller_sku + ean would be ALTA,
// but title voltage contradicts the internal product ⇒ capped BAIXA/REJECT.
func TestCase7VoltageHardNegativeCapsBaixaReject(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 30, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx7", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX7",
		SellerSKU: "MLB-SKU-7", EAN: "7897777777770", Title: "Furadeira 220V", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:MLB-SKU-7":     {{InternalProductID: canonicalIDPtr(700), Name: "Furadeira 110V"}},
		"ean:7897777777770": {{InternalProductID: canonicalIDPtr(700), Name: "Furadeira 110V"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
		t.Fatalf("confidence_band=%s, want BAIXA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject {
		t.Fatalf("match_status=%s, want REJECT", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want seller_sku FOR", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR", candidate.Reasons)
	}
	titleAgainst, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst)
	if !ok || !strings.Contains(titleAgainst.Detail, "voltagem") {
		t.Fatalf("reasons=%#v, want title AGAINST citing voltagem divergence", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 8 — NO-CANDIDATE: seller_sku/ean/title all fail to resolve. Honest,
// non-empty reasons — never an empty row.
func TestCase8NoCandidateYieldsZeroConfidenceNoCandidate(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 35, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx8", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX8",
		SellerSKU: "SKU-UNKNOWN", EAN: "0000000000000", Title: "Item Desconhecido", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.Confidence != 0 {
		t.Fatalf("confidence=%d, want 0", candidate.Confidence)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusNoCandidate {
		t.Fatalf("match_status=%s, want NO_CANDIDATE", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want seller_sku UNAVAILABLE", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want ean UNAVAILABLE", candidate.Reasons)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}

// Case 9 — MARCA-REFFORN-UNAVAILABLE (explicit): marca/refforn must appear
// as UNAVAILABLE regardless of band/status (ADR-17, motivo sempre
// visível). Bound to case 1's ALTA/ACCEPT payload with a dedicated
// assertion, per PLAN-M04 §4.
func TestCase9MarcaAndReffornAlwaysUnavailableOnConcordantPayload(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 40, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx9", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX9",
		SellerSKU: "MLB-SKU-9", EAN: "7891234567896", Title: "Furadeira Bosch 550W", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:MLB-SKU-9":     {{InternalProductID: canonicalIDPtr(900), Name: "Furadeira Bosch 550W"}},
		"ean:7891234567896": {{InternalProductID: canonicalIDPtr(900), Name: "Furadeira Bosch 550W"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	marca, ok := findReason(candidate.Reasons, "marca", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok {
		t.Fatalf("reasons=%#v, want {anchor:marca,direction:UNAVAILABLE}", candidate.Reasons)
	}
	refforn, ok := findReason(candidate.Reasons, "refforn", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok {
		t.Fatalf("reasons=%#v, want {anchor:refforn,direction:UNAVAILABLE}", candidate.Reasons)
	}
	_ = marca
	_ = refforn
}

// Case 10 — MEDIDA-HARD-NEGATIVE-REJECT: EAN matches the same codprod, but
// the listing title's dimension contradicts the internal product's
// dimension ⇒ contradição vence EAN (kit/combo, cor, medida, voltagem are
// all binding). Caps BAIXA/REJECT with a title AGAINST reason.
func TestCase10DimensionHardNegativeCapsBaixaReject(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 18, 9, 45, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-fx10", ProviderCode: "mercado_livre", ProviderItemID: "MLB-FX10",
		EAN: "7896541230011", Title: "Lona Impermeável 300x150", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"ean:7896541230011": {{InternalProductID: canonicalIDPtr(1000), Name: "Lona Impermeável 220x100"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
		t.Fatalf("confidence_band=%s, want BAIXA", candidate.ConfidenceBand)
	}
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject {
		t.Fatalf("match_status=%s, want REJECT", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR", candidate.Reasons)
	}
	titleAgainst, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst)
	if !ok || !strings.Contains(titleAgainst.Detail, "medida/dimensão") {
		t.Fatalf("reasons=%#v, want title AGAINST citing medida/dimensão divergence", candidate.Reasons)
	}
	if !strings.Contains(titleAgainst.Detail, "300x150") || !strings.Contains(titleAgainst.Detail, "220x100") {
		t.Fatalf("title AGAINST detail=%q, want it to cite both divergent dimensions", titleAgainst.Detail)
	}
	assertMandatoryUnavailableReasons(t, candidate.Reasons)
}
