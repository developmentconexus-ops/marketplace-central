package application

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	productlinksports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
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

type stubIdentityAnchorReader struct {
	declarations map[string][]productlinksports.ProviderIdentityAnchor
	calls        map[string]int
	errors       map[string]error
}

type recordingAutoApprover struct {
	candidates []productlinksdomain.LinkCandidate
}

func (a *recordingAutoApprover) AutoApproveCandidate(_ context.Context, input AutoApproveCandidateInput) (bool, error) {
	a.candidates = append(a.candidates, input.Candidate)
	return true, nil
}

func (s *stubIdentityAnchorReader) ProviderIdentityAnchors(providerCode string) ([]productlinksports.ProviderIdentityAnchor, error) {
	if s.calls == nil {
		s.calls = map[string]int{}
	}
	s.calls[providerCode]++
	if err := s.errors[providerCode]; err != nil {
		return nil, err
	}
	return s.declarations[providerCode], nil
}

func mercadoLivreIdentityAnchorReader() *stubIdentityAnchorReader {
	return &stubIdentityAnchorReader{declarations: map[string][]productlinksports.ProviderIdentityAnchor{
		"mercado_livre": {
			{Anchor: "seller_sku", Supplied: true},
			{Anchor: "ean", Supplied: true},
			{Anchor: "title", Supplied: true},
			{Anchor: "marca", Supplied: false},
		},
	}}
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
		Snapshots:       snapshots,
		Matcher:         &stubProductMatcher{},
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
	})

	if _, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-1"}); err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if snapshots.limit != 0 {
		t.Fatalf("snapshot reader limit = %d, want 0 (unbounded)", snapshots.limit)
	}
}

func TestGenerateLinkCandidatesUsesProviderDeclarationForUnavailableReasons(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 0, 0, 0, time.UTC)
	reader := &stubIdentityAnchorReader{declarations: map[string][]productlinksports.ProviderIdentityAnchor{
		"provider-a": {
			{Anchor: "seller_sku", Supplied: false},
			{Anchor: "ean", Supplied: true},
			{Anchor: "title", Supplied: true},
			{Anchor: "marca", Supplied: false},
		},
	}}
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-declaration", ProviderCode: "provider-a", ProviderItemID: "item-a",
		SellerSKU: "SKU-A", Title: "Produto A", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:SKU-A": {{InternalProductID: canonicalIDPtr(101), Name: "Produto A"}},
	}}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots:       &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
		Matcher:         matcher,
		Store:           &stubCandidateStore{},
		IdentityAnchors: reader,
		Now:             func() time.Time { return now },
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: snapshot.InstallationID})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("result.Items=%#v, want exactly one candidate", result.Items)
	}
	reasons := result.Items[0].Reasons
	assertUniqueReasonAnchors(t, reasons)
	if _, ok := findReason(reasons, "title", productlinksdomain.LinkCandidateReasonDirectionUnavailable); ok {
		t.Fatalf("reasons=%#v, want supplied title without UNAVAILABLE reason", reasons)
	}
	marca, ok := findReason(reasons, "marca", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok || marca.Detail != "provider não fornece a âncora marca" {
		t.Fatalf("marca reason=%#v, want provider declaration detail", marca)
	}
	ean, ok := findReason(reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
	if !ok || ean.Side != productlinksdomain.LinkCandidateReasonSideBoth || ean.Detail != "sem EAN para corroborar o CODPROD" {
		t.Fatalf("ean reason=%#v, want supported-but-empty incomparable detail", ean)
	}
	if marca.Detail == ean.Detail {
		t.Fatalf("marca and ean details are equal: %q", marca.Detail)
	}
}

func TestGenerateLinkCandidatesKeepsEveryDeclaredAnchorVisible(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 27, 9, 0, 0, 0, time.UTC)
	for name, tc := range map[string]struct {
		declared  bool
		title     string
		name      string
		direction productlinksdomain.LinkCandidateReasonDirection
		side      productlinksdomain.LinkCandidateReasonSide
	}{
		"provider does not declare": {
			declared: false, title: "Listing title", name: "Product name",
			direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable,
		},
		"listing value is empty": {
			declared: true, title: "", name: "Product name",
			direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
			side:      productlinksdomain.LinkCandidateReasonSideProvider,
		},
		"ERP value is empty": {
			declared: true, title: "Listing title", name: "",
			direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
			side:      productlinksdomain.LinkCandidateReasonSideERP,
		},
		"both values are empty": {
			declared: true, title: "", name: "",
			direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
			side:      productlinksdomain.LinkCandidateReasonSideBoth,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if hardNegative, detail := detectHardNegative(tc.title, tc.name); hardNegative {
				t.Fatalf("detectHardNegative(%q, %q)=(true, %q), want false for C2 fixture", tc.title, tc.name, detail)
			}
			product := internalreaddomain.ProductCandidate{
				InternalProductID: canonicalIDPtr(801), Name: tc.name,
				ReferenceCode: stringPtr("SKU-ANCHOR"), EAN: stringPtr("EAN-ANCHOR"),
			}
			snapshot := productlinksdomain.ListingSnapshot{
				InstallationID: "inst-f02", ProviderCode: "provider-f02", ProviderItemID: "item-" + name,
				SellerSKU: "SKU-ANCHOR", EAN: "EAN-ANCHOR", Title: tc.title, FetchedAt: now,
			}
			matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
				"sku:SKU-ANCHOR": {product}, "ean:EAN-ANCHOR": {product},
			}}
			reader := &stubIdentityAnchorReader{declarations: map[string][]productlinksports.ProviderIdentityAnchor{
				"provider-f02": {
					{Anchor: "seller_sku", Supplied: true}, {Anchor: "ean", Supplied: true},
					{Anchor: "title", Supplied: tc.declared}, {Anchor: "marca", Supplied: false},
				},
			}}
			svc := NewGenerationService(GenerationServiceConfig{
				Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
				Matcher:   matcher, Store: &stubCandidateStore{}, IdentityAnchors: reader,
				Now: func() time.Time { return now },
			})

			result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: snapshot.InstallationID})
			if err != nil {
				t.Fatalf("GenerateLinkCandidates() error = %v", err)
			}
			if len(result.Items) != 1 {
				t.Fatalf("result.Items=%#v, want one concordant candidate", result.Items)
			}
			reason, ok := findReason(result.Items[0].Reasons, "title", tc.direction)
			if !ok {
				t.Fatalf("reasons=%#v, want title %s present", result.Items[0].Reasons, tc.direction)
			}
			if reason.Side != tc.side {
				t.Fatalf("title reason=%#v, want side %q", reason, tc.side)
			}
		})
	}
}

func TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 27, 10, 0, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-d1", ProviderCode: "mercado_livre", ProviderItemID: "item-d1",
		SellerSKU: "SKU-D1", EAN: "EAN-LISTING", Title: "Produto D1", FetchedAt: now,
	}
	product := internalreaddomain.ProductCandidate{
		InternalProductID: canonicalIDPtr(802), Name: "Produto D1", EAN: stringPtr("EAN-ERP"),
	}
	candidate := generateSingle(t, snapshot, &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:SKU-D1": {product},
	}}, now)

	for _, reason := range candidate.Reasons {
		if reason.Anchor == "ean" {
			return
		}
	}
	t.Fatalf("reasons=%#v, want ean anchor to remain present", candidate.Reasons)
}

func TestUnresolvedDeclaredAnchorReadsProviderValueWithoutERPProduct(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		ean              string
		wantIncomparable bool
		wantSide         productlinksdomain.LinkCandidateReasonSide
	}{
		"listing value empty": {
			wantIncomparable: true,
			wantSide:         productlinksdomain.LinkCandidateReasonSideProvider,
		},
		"listing value present": {
			ean: "EAN-D2", wantIncomparable: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			now := time.Date(2026, 7, 27, 10, 5, 0, 0, time.UTC)
			snapshot := productlinksdomain.ListingSnapshot{
				InstallationID: "inst-d2", ProviderCode: "provider-d2", ProviderItemID: "item-d2",
				EAN: tc.ean, Title: "Produto D2", FetchedAt: now,
			}
			reader := &stubIdentityAnchorReader{declarations: map[string][]productlinksports.ProviderIdentityAnchor{
				"provider-d2": {{Anchor: "ean", Supplied: true}},
			}}
			svc := NewGenerationService(GenerationServiceConfig{
				Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
				Matcher:   &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{}},
				Store:     &stubCandidateStore{}, IdentityAnchors: reader,
				Now: func() time.Time { return now },
			})
			result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: snapshot.InstallationID})
			if err != nil {
				t.Fatalf("GenerateLinkCandidates() error = %v", err)
			}
			if len(result.Items) != 1 {
				t.Fatalf("result.Items=%#v, want one unresolved candidate", result.Items)
			}

			if tc.wantIncomparable {
				reason, ok := findReason(result.Items[0].Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
				if !ok || reason.Side != tc.wantSide {
					t.Fatalf("reasons=%#v, want ean INCOMPARABLE with side %q", result.Items[0].Reasons, tc.wantSide)
				}
				return
			}
			if _, ok := findReason(result.Items[0].Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable); ok {
				t.Fatalf("reasons=%#v, want seeded ean sem correspondência reason to stand", result.Items[0].Reasons)
			}
			seeded, ok := findReason(result.Items[0].Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
			if !ok || seeded.Detail != "ean sem correspondência" {
				t.Fatalf("reasons=%#v, want seeded ean sem correspondência reason", result.Items[0].Reasons)
			}
		})
	}
}

func TestProviderDeclaredUnmodelledAnchorIsIncomparableWithoutSide(t *testing.T) {
	t.Parallel()
	brand := "Marca ERP"
	product := internalreaddomain.ProductCandidate{BrandName: &brand}
	reasons := appendProviderDeclaredUnavailableReasons(nil, newProviderIdentityAnchorComparison(
		productlinksdomain.ListingSnapshot{},
		[]productlinksports.ProviderIdentityAnchor{{Anchor: "marca", Supplied: true}},
		&product,
	))
	if len(reasons) != 1 || reasons[0].Direction != productlinksdomain.LinkCandidateReasonDirectionIncomparable || reasons[0].Side != "" {
		t.Fatalf("reasons=%#v, want one side-less INCOMPARABLE marca reason", reasons)
	}
}

func TestProductLinkReasonSideJSONIsOnlyPresentForIncomparable(t *testing.T) {
	t.Parallel()
	for _, direction := range []productlinksdomain.LinkCandidateReasonDirection{
		productlinksdomain.LinkCandidateReasonDirectionFor,
		productlinksdomain.LinkCandidateReasonDirectionAgainst,
		productlinksdomain.LinkCandidateReasonDirectionUnavailable,
	} {
		payload, err := json.Marshal(productlinksdomain.LinkCandidateReason{Anchor: "ean", Direction: direction})
		if err != nil {
			t.Fatalf("json.Marshal(%s) error = %v", direction, err)
		}
		if strings.Contains(string(payload), `"side"`) {
			t.Fatalf("direction %s payload=%s, want side key absent", direction, payload)
		}
	}
	payload, err := json.Marshal(productlinksdomain.LinkCandidateReason{
		Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
		Side: productlinksdomain.LinkCandidateReasonSideERP,
	})
	if err != nil {
		t.Fatalf("json.Marshal(INCOMPARABLE) error = %v", err)
	}
	if !strings.Contains(string(payload), `"side":"erp"`) {
		t.Fatalf("INCOMPARABLE payload=%s, want side key present", payload)
	}
}

func stringPtr(value string) *string { return &value }

func TestProviderUnavailableReasonPrecedenceKeepsObservedEvidence(t *testing.T) {
	t.Parallel()
	reasons := appendProviderDeclaredUnavailableReasons(
		[]productlinksdomain.LinkCandidateReason{{
			Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "ean observado",
		}},
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{{Anchor: "ean", Supplied: false}}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{{
		Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "ean observado",
	}}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want observed evidence %#v", reasons, want)
	}
}

func TestProviderUnavailableReasonPrecedenceKeepsDeclarationOverSeedUnavailable(t *testing.T) {
	t.Parallel()
	reasons := appendProviderDeclaredUnavailableReasons(
		[]productlinksdomain.LinkCandidateReason{{
			Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência",
		}},
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{{Anchor: "ean", Supplied: false}}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{{
		Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "provider não fornece a âncora ean",
	}}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want declaration reason %#v", reasons, want)
	}
}

func TestProviderUnavailableReasonEmitsDeclarationWithoutSeed(t *testing.T) {
	t.Parallel()
	reasons := appendProviderDeclaredUnavailableReasons(
		nil,
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{{Anchor: "marca", Supplied: false}}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{{
		Anchor: "marca", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "provider não fornece a âncora marca",
	}}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want declaration reason %#v", reasons, want)
	}
}

func TestProviderSuppliedReasonPromotesFirstSeedUnavailable(t *testing.T) {
	t.Parallel()
	reasons := appendProviderDeclaredUnavailableReasons(
		[]productlinksdomain.LinkCandidateReason{
			{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "first"},
			{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "second"},
		},
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{{Anchor: "ean", Supplied: true}}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{{
		Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
		Side: productlinksdomain.LinkCandidateReasonSideProvider, Detail: "first",
	}}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want first seed promoted to provider-side INCOMPARABLE %#v", reasons, want)
	}
}

func TestProviderSuppliedDeclarationPromotesSeededAbsenceWithoutDuplicating(t *testing.T) {
	t.Parallel()
	seed := []productlinksdomain.LinkCandidateReason{
		{Anchor: "seller_sku", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "seller_sku observado"},
		{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
	}
	reasons := appendProviderDeclaredUnavailableReasons(
		seed,
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{
			{Anchor: "seller_sku", Supplied: true},
			{Anchor: "ean", Supplied: true},
		}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{
		{Anchor: "seller_sku", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "seller_sku observado"},
		{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionIncomparable,
			Side: productlinksdomain.LinkCandidateReasonSideProvider, Detail: "ean sem correspondência"},
	}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want promoted seed %#v", reasons, want)
	}
}

func TestProviderUnavailableReasonOrderingIsStable(t *testing.T) {
	t.Parallel()
	reasons := appendProviderDeclaredUnavailableReasons(
		[]productlinksdomain.LinkCandidateReason{
			{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "ean sem correspondência"},
			{Anchor: "title", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "title observado"},
		},
		newProviderIdentityAnchorComparison(productlinksdomain.ListingSnapshot{}, []productlinksports.ProviderIdentityAnchor{
			{Anchor: "marca", Supplied: false},
			{Anchor: "ean", Supplied: false},
			{Anchor: "seller_sku", Supplied: false},
		}, nil),
	)

	want := []productlinksdomain.LinkCandidateReason{
		{Anchor: "ean", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "provider não fornece a âncora ean"},
		{Anchor: "title", Direction: productlinksdomain.LinkCandidateReasonDirectionFor, Detail: "title observado"},
		{Anchor: "marca", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "provider não fornece a âncora marca"},
		{Anchor: "seller_sku", Direction: productlinksdomain.LinkCandidateReasonDirectionUnavailable, Detail: "provider não fornece a âncora seller_sku"},
	}
	if !reflect.DeepEqual(reasons, want) {
		t.Fatalf("reasons=%#v, want stable order %#v", reasons, want)
	}
}

func TestTitleMatchHardNegativeKeepsTitleForAndAgainstInSeedOrder(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 10, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-title-hard-negative", ProviderCode: "mercado_livre", ProviderItemID: "MLB-TITLE-HARD-NEGATIVE",
		Title: "Kit 3 Toalheiro Simples 50cm Cromado",
	}
	products := []internalreaddomain.ProductCandidate{{
		InternalProductID: canonicalIDPtr(90001), ProductID: 90001,
		Name: "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO",
	}}

	candidates := buildCandidatesFromProducts(
		snapshot,
		products,
		productlinksdomain.LinkCandidateStateTitleMatch,
		productlinksdomain.LinkCandidateMatchInputTitle,
		"toalheiro simples",
		mercadoLivreIdentityAnchorReader().declarations["mercado_livre"],
		now,
	)
	if len(candidates) == 0 {
		t.Fatal("buildCandidatesFromProducts() returned zero candidates, want the canonical product candidate")
	}
	candidate := candidates[0]
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject {
		t.Fatalf("match status=%q, want REJECT", candidate.MatchStatus)
	}

	var titleReasons []productlinksdomain.LinkCandidateReason
	for _, reason := range candidate.Reasons {
		if reason.Anchor == "title" {
			titleReasons = append(titleReasons, reason)
		}
	}
	if len(titleReasons) != 2 {
		t.Fatalf("title reasons=%#v, want FOR then AGAINST", titleReasons)
	}
	if titleReasons[0].Direction != productlinksdomain.LinkCandidateReasonDirectionFor {
		t.Fatalf("first title reason=%#v, want FOR", titleReasons[0])
	}
	if titleReasons[1].Direction != productlinksdomain.LinkCandidateReasonDirectionAgainst {
		t.Fatalf("second title reason=%#v, want AGAINST", titleReasons[1])
	}
	if !strings.Contains(titleReasons[1].Detail, "kit") {
		t.Fatalf("title AGAINST detail=%q, want hard negative named", titleReasons[1].Detail)
	}
}

func TestNormalizeDimensionTokenFailsClosed(t *testing.T) {
	t.Parallel()
	key, parsed := normalizeDimensionToken("abcmm")
	if parsed {
		t.Fatal("normalizeDimensionToken(abcmm) parsed=true, want false")
	}
	if key != "abcmm" {
		t.Fatalf("normalizeDimensionToken(abcmm) key=%q, want raw token", key)
	}
}

func TestGenerateLinkCandidatesResolvesIdentityAnchorsOncePerProvider(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 5, 0, 0, time.UTC)
	reader := &stubIdentityAnchorReader{declarations: map[string][]productlinksports.ProviderIdentityAnchor{
		"provider-a": {{Anchor: "seller_sku", Supplied: true}},
		"provider-b": {{Anchor: "seller_sku", Supplied: true}},
	}}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{
			{InstallationID: "inst-memo", ProviderCode: " provider-a ", ProviderItemID: "item-a", FetchedAt: now},
			{InstallationID: "inst-memo", ProviderCode: "provider-a", ProviderItemID: "item-a2", FetchedAt: now},
			{InstallationID: "inst-memo", ProviderCode: "provider-b", ProviderItemID: "item-b", FetchedAt: now},
		}},
		Matcher:         &stubProductMatcher{},
		Store:           &stubCandidateStore{},
		IdentityAnchors: reader,
		Now:             func() time.Time { return now },
	})

	if _, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-memo"}); err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if reader.calls["provider-a"] != 1 || reader.calls["provider-b"] != 1 {
		t.Fatalf("identity anchor calls=%#v, want one lookup for each distinct trimmed provider", reader.calls)
	}
}

// An unconfigured identity-anchor reader must refuse BEFORE any work, including on an
// empty snapshot batch. Empty is the case that hides the defect: the resolve loop never
// runs, so a missing guard returns a clean zero result and the caller reads "nothing to
// link" instead of "this engine was never wired". A non-empty batch would deref nil and
// at least announce itself.
func TestGenerateLinkCandidatesRefusesWithoutIdentityAnchorReader(t *testing.T) {
	t.Parallel()
	store := &stubCandidateStore{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots:       &stubSnapshotReader{snapshots: nil},
		Matcher:         &stubProductMatcher{},
		Store:           store,
		IdentityAnchors: nil,
	})

	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-unwired"})
	if err == nil {
		t.Fatalf("GenerateLinkCandidates() error = nil, want PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED (result=%#v)", result)
	}
	if err.Error() != "PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED" {
		t.Fatalf("GenerateLinkCandidates() error = %q, want PRODUCT_LINKS_CANDIDATE_ENGINE_NOT_CONFIGURED", err.Error())
	}
	if store.installationID != "" {
		t.Fatalf("store written with installationID=%q, want no write at all", store.installationID)
	}
}

func TestGenerateLinkCandidatesFailsWhenIdentityAnchorDeclarationUnavailable(t *testing.T) {
	t.Parallel()
	reader := &stubIdentityAnchorReader{
		declarations: map[string][]productlinksports.ProviderIdentityAnchor{},
		errors:       map[string]error{"provider-failed": fmt.Errorf("declaration backend unavailable")},
	}
	store := &stubCandidateStore{}
	approver := &recordingAutoApprover{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-failure", ProviderCode: "provider-failed", ProviderItemID: "item-failed",
		}}},
		Matcher:         &stubProductMatcher{},
		Store:           store,
		AutoApprover:    approver,
		IdentityAnchors: reader,
	})

	_, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-failure"})
	if err == nil || !strings.Contains(err.Error(), "PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE") || !strings.Contains(err.Error(), "provider-failed") {
		t.Fatalf("GenerateLinkCandidates() error = %v, want named provider declaration failure", err)
	}
	if store.installationID != "" {
		t.Fatalf("store installationID=%q, want no persistence", store.installationID)
	}
	if len(approver.candidates) != 0 {
		t.Fatalf("auto-approved candidates=%#v, want none", approver.candidates)
	}
}

func TestGenerateLinkCandidatesFailsWhenProviderCodeIsEmpty(t *testing.T) {
	t.Parallel()

	reader := &stubIdentityAnchorReader{
		declarations: map[string][]productlinksports.ProviderIdentityAnchor{},
		errors:       map[string]error{},
	}
	store := &stubCandidateStore{}
	approver := &recordingAutoApprover{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-empty-provider", ProviderCode: "", ProviderItemID: "item-empty-provider",
		}}},
		Matcher:         &stubProductMatcher{},
		Store:           store,
		AutoApprover:    approver,
		IdentityAnchors: reader,
	})

	_, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-empty-provider"})
	if err == nil || !strings.Contains(err.Error(), "PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE") || !strings.Contains(err.Error(), "provider code is empty") {
		t.Fatalf("GenerateLinkCandidates() error = %v, want empty provider code failure", err)
	}
	if store.installationID != "" {
		t.Fatalf("store installationID=%q, want no persistence", store.installationID)
	}
	if len(approver.candidates) != 0 {
		t.Fatalf("auto-approved candidates=%#v, want none", approver.candidates)
	}
}

func TestGenerateLinkCandidatesFailsWhenIdentityAnchorDeclarationIsNil(t *testing.T) {
	t.Parallel()

	reader := &stubIdentityAnchorReader{
		declarations: map[string][]productlinksports.ProviderIdentityAnchor{},
		errors:       map[string]error{},
	}
	store := &stubCandidateStore{}
	approver := &recordingAutoApprover{}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{{
			InstallationID: "inst-nil-declaration", ProviderCode: "provider-missing", ProviderItemID: "item-missing-declaration",
		}}},
		Matcher:         &stubProductMatcher{},
		Store:           store,
		AutoApprover:    approver,
		IdentityAnchors: reader,
	})

	_, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-nil-declaration"})
	if err == nil || !strings.Contains(err.Error(), "PRODUCT_LINKS_PROVIDER_IDENTITY_ANCHORS_UNAVAILABLE") || !strings.Contains(err.Error(), "identity anchor declaration is nil") {
		t.Fatalf("GenerateLinkCandidates() error = %v, want nil declaration failure", err)
	}
	if store.installationID != "" {
		t.Fatalf("store installationID=%q, want no persistence", store.installationID)
	}
	if len(approver.candidates) != 0 {
		t.Fatalf("auto-approved candidates=%#v, want none", approver.candidates)
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
	svc := NewGenerationService(GenerationServiceConfig{Snapshots: snapshots, Matcher: matcher, Store: store, IdentityAnchors: mercadoLivreIdentityAnchorReader(), Now: func() time.Time { return now }})

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
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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
		Store:           store,
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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
		Matcher:         &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{}},
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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
				Store:           store,
				IdentityAnchors: mercadoLivreIdentityAnchorReader(),
				Now:             func() time.Time { return now },
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

	candidates := buildCandidatesFromProducts(snapshot, products, productlinksdomain.LinkCandidateStateTitleMatch, productlinksdomain.LinkCandidateMatchInputTitle, "stable", mercadoLivreIdentityAnchorReader().declarations["mercado_livre"], now)
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
// the anchor model. marca must ALWAYS appear as UNAVAILABLE (ADR-17)
// regardless of band/status. `refforn` used to sit beside it and left the
// cross-side vocabulary in D-A (D-122) — no marketplace ever supplies it.

func findReason(reasons []productlinksdomain.LinkCandidateReason, anchor string, direction productlinksdomain.LinkCandidateReasonDirection) (productlinksdomain.LinkCandidateReason, bool) {
	for _, reason := range reasons {
		if reason.Anchor == anchor && reason.Direction == direction {
			return reason, true
		}
	}
	return productlinksdomain.LinkCandidateReason{}, false
}

func assertUniqueReasonAnchors(t *testing.T, reasons []productlinksdomain.LinkCandidateReason) {
	t.Helper()
	seen := make(map[string]struct{}, len(reasons))
	for _, reason := range reasons {
		if _, exists := seen[reason.Anchor]; exists {
			t.Fatalf("reasons=%#v, anchor %q appears more than once", reasons, reason.Anchor)
		}
		seen[reason.Anchor] = struct{}{}
	}
}

func assertProviderDeclaredUnavailableReasons(t *testing.T, reasons []productlinksdomain.LinkCandidateReason) {
	t.Helper()
	marca, ok := findReason(reasons, "marca", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok || marca.Detail != "provider não fornece a âncora marca" {
		t.Fatalf("reasons=%#v, want provider-declared marca UNAVAILABLE", reasons)
	}
}

func generateSingle(t *testing.T, snapshot productlinksdomain.ListingSnapshot, matcher *stubProductMatcher, now time.Time) productlinksdomain.LinkCandidate {
	t.Helper()
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots:       &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
		Matcher:         matcher,
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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

func TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 0, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB4735326915",
		SellerSKU: "33698", EAN: "", Title: "Toalheiro Simples Soul Zen 50cm Cromado", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:33698": {{InternalProductID: canonicalIDPtr(33698), Name: "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO"}},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusConfirm {
		t.Fatalf("match_status=%s, want CONFIRM", candidate.MatchStatus)
	}
	if candidate.Confidence != 70 {
		t.Fatalf("confidence=%d, want 70", candidate.Confidence)
	}
	if candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandMedia {
		t.Fatalf("confidence_band=%s, want MEDIA", candidate.ConfidenceBand)
	}
	if reason, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst); ok {
		t.Fatalf("title AGAINST reason=%#v, want none", reason)
	}
	eanReason, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
	if !ok || eanReason.Side != productlinksdomain.LinkCandidateReasonSideBoth || eanReason.Detail != "sem EAN para corroborar o CODPROD" {
		t.Fatalf("ean reason=%#v, want exact missing-EAN incomparable detail", eanReason)
	}
}

func TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 5, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-dimension", ProviderCode: "mercado_livre", ProviderItemID: "MLB-DIM-EQUIV",
		SellerSKU: "DIM-500", EAN: "7890000000508", Title: "Toalheiro Soul 50cm", FetchedAt: now,
	}
	product := internalreaddomain.ProductCandidate{InternalProductID: canonicalIDPtr(500), Name: "TOALHEIRO SOUL 500MM"}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:DIM-500":       {product},
		"ean:7890000000508": {product},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusAccept ||
		candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandAlta {
		t.Fatalf("candidate=%#v, want ALTA/ACCEPT", candidate)
	}
	if reason, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst); ok {
		t.Fatalf("title AGAINST reason=%#v, want none", reason)
	}
}

func TestDimensionCanonicalizationUsesExactMillimetres(t *testing.T) {
	t.Parallel()
	for name, titles := range map[string][2]string{
		"inches":        {"Produto 2pol", "Produto 50.8mm"},
		"metres":        {"Produto 1m", "Produto 100cm"},
		"decimal comma": {"Produto 1,5m", "Produto 150cm"},
		"decimal point": {"Produto 1.5m", "Produto 150cm"},
	} {
		t.Run(name, func(t *testing.T) {
			if hardNegative, detail := detectHardNegative(titles[0], titles[1]); hardNegative {
				t.Fatalf("detectHardNegative(%q, %q)=(true, %q), want equivalent dimensions", titles[0], titles[1], detail)
			}
		})
	}
}

func TestNInOneTitleIsNotReadAsADimension(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 7, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-n-in-one", ProviderCode: "mercado_livre", ProviderItemID: "MLB-N-IN-ONE",
		SellerSKU: "ESCOVA-5-EM-1", EAN: "7890000000126",
		Title: "Escova Secadora Modeladora 5 in 1 30mm Rosa", FetchedAt: now,
	}
	product := internalreaddomain.ProductCandidate{
		InternalProductID: canonicalIDPtr(126),
		Name:              "Escova Secadora 5 em 1 30mm Rosa",
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:ESCOVA-5-EM-1": {product},
		"ean:7890000000126": {product},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.MatchStatus == productlinksdomain.LinkCandidateMatchStatusReject ||
		candidate.ConfidenceBand == productlinksdomain.LinkCandidateConfidenceBandBaixa {
		reason, _ := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst)
		t.Fatalf("candidate=%#v, fabricated dimension contradiction=%q", candidate, reason.Detail)
	}
}

func TestDimensionDisplaySurvivesLengthChangingCaseFold(t *testing.T) {
	t.Parallel()

	hardNegative, detail := detectHardNegative("İnox 50cm", "İnox 40cm")

	if !hardNegative {
		t.Fatal("detectHardNegative()=false, want dimension contradiction")
	}
	if !strings.Contains(detail, "50cm") || !strings.Contains(detail, "40cm") {
		t.Fatalf("detail=%q, want original 50cm and 40cm tokens", detail)
	}
	if strings.Contains(detail, "mm:") || strings.Contains(detail, "ox 5") || strings.Contains(detail, "ox 4") {
		t.Fatalf("detail=%q, want no canonical or truncated fragment", detail)
	}
}

func TestRepeatedMeasurementIsNotRepeatedInTheDetail(t *testing.T) {
	t.Parallel()

	hardNegative, detail := detectHardNegative("Produto 50cm por 50cm", "Produto 40cm")

	if !hardNegative {
		t.Fatal("detectHardNegative()=false, want dimension contradiction")
	}
	if strings.Count(detail, "50cm") != 1 {
		t.Fatalf("detail=%q, want 50cm exactly once", detail)
	}
}

func TestDifferentCanonicalDimensionsStillRejectConcordantCandidate(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 26, 9, 10, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-dimension", ProviderCode: "mercado_livre", ProviderItemID: "MLB-DIM-DIFF",
		SellerSKU: "DIM-DIFF", EAN: "7890000000409", Title: "Toalheiro Soul 50cm", FetchedAt: now,
	}
	product := internalreaddomain.ProductCandidate{InternalProductID: canonicalIDPtr(400), Name: "TOALHEIRO SOUL 40cm"}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:DIM-DIFF":      {product},
		"ean:7890000000409": {product},
	}}

	candidate := generateSingle(t, snapshot, matcher, now)

	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReject ||
		candidate.ConfidenceBand != productlinksdomain.LinkCandidateConfidenceBandBaixa {
		t.Fatalf("candidate=%#v, want BAIXA/REJECT", candidate)
	}
	reason, ok := findReason(candidate.Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst)
	if !ok || !strings.Contains(reason.Detail, "medida/dimensão") ||
		!strings.Contains(reason.Detail, "50cm") || !strings.Contains(reason.Detail, "40cm") {
		t.Fatalf("title AGAINST reason=%#v, want readable original dimensions", reason)
	}
	if strings.Contains(reason.Detail, "mm:") || strings.Contains(reason.Detail, "/1") {
		t.Fatalf("title AGAINST detail=%q, want no canonical representation", reason.Detail)
	}
}

func TestDimensionPresenceAndGradeRulesRemainNonBlocking(t *testing.T) {
	t.Parallel()
	for name, titles := range map[string][2]string{
		"neither side":       {"Camisa lisa", "Camisa básica"},
		"listing side only":  {"Mesa 50cm", "Mesa"},
		"internal side only": {"Mesa", "Mesa 50cm"},
		"lowercase metre":    {"Tecido 1m", "Tecido 100cm"},
	} {
		t.Run(name, func(t *testing.T) {
			if hardNegative, detail := detectHardNegative(titles[0], titles[1]); hardNegative {
				t.Fatalf("detectHardNegative(%q, %q)=(true, %q), want non-blocking", titles[0], titles[1], detail)
			}
		})
	}
	hardNegative, detail := detectHardNegative("Camisa M", "Camisa G")
	if !hardNegative || !strings.Contains(detail, "medida/dimensão") {
		t.Fatalf("detectHardNegative grade result=(%t, %q), want uppercase grades compared", hardNegative, detail)
	}
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
		"sku:MLB-SKU-1":     {{InternalProductID: canonicalIDPtr(100), Name: "Furadeira Bosch 550W"}},
		"ean:7891234567895": {{InternalProductID: canonicalIDPtr(100), Name: "Furadeira Bosch 550W"}},
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
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}

// Case 2 — SKU-ALONE: seller_sku matches, EAN absent from the snapshot. D-121-2
// moved this to the confirmation queue: one anchor resolved one product and
// nothing contradicts it, but nothing corroborates it either, so a human says
// yes — never the machine.
func TestCase2SellerSKUAloneWithoutEANYieldsMediaConfirm(t *testing.T) {
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
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusConfirm {
		t.Fatalf("match_status=%s, want CONFIRM", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want seller_sku FOR", candidate.Reasons)
	}
	eanReason, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
	if !ok {
		t.Fatalf("reasons=%#v, want ean UNAVAILABLE", candidate.Reasons)
	}
	// M05-C14/AC-11: the warning names the anchor that is missing. A generic
	// "sem corroboração" would not tell the operator what confirming supplies.
	if eanReason.Detail != "sem EAN para corroborar o CODPROD" {
		t.Fatalf("ean reason detail=%q, want the missing anchor named", eanReason.Detail)
	}
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}

// Case 3 — EAN-ALONE-MEDIA: seller_sku has no match, EAN corroborates
// (unproved) a single codprod.
func TestCase3EANAloneYieldsMediaConfirm(t *testing.T) {
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
	if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusConfirm {
		t.Fatalf("match_status=%s, want CONFIRM", candidate.MatchStatus)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionFor); !ok {
		t.Fatalf("reasons=%#v, want ean FOR (unproved)", candidate.Reasons)
	}
	skuReason, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
	if !ok {
		t.Fatalf("reasons=%#v, want seller_sku INCOMPARABLE", candidate.Reasons)
	}
	if !strings.HasPrefix(skuReason.Detail, "sem CODPROD para corroborar o EAN") {
		t.Fatalf("seller_sku reason detail=%q, want the missing anchor named", skuReason.Detail)
	}
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
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
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionIncomparable); !ok {
		t.Fatalf("reasons=%#v, want seller_sku INCOMPARABLE", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionIncomparable); !ok {
		t.Fatalf("reasons=%#v, want ean INCOMPARABLE", candidate.Reasons)
	}
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}

// Case 5 — SKU-EAN-CONFLICT-REJECT: seller_sku and ean point at different
// codprod. Both candidates are BAIXA/REJECT; each cites its own anchor FOR
// and the other anchor AGAINST.
func TestCase5SKUEANConflictYieldsBaixaReviewBothSides(t *testing.T) {
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
		Snapshots:       &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{snapshot}},
		Matcher:         matcher,
		Store:           &stubCandidateStore{},
		IdentityAnchors: mercadoLivreIdentityAnchorReader(),
		Now:             func() time.Time { return now },
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
		// AC-08 (D-121): in a conflict neither anchor wins — both products go to
		// the operator, and neither side is quietly discarded as REJECT.
		if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
			t.Fatalf("candidate=%#v, want REVIEW", candidate)
		}
		assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
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
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
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
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}

// Case 8 — NO ANCHOR RESOLVED: seller_sku/ean/title all fail to resolve.
// M05-C5: the listing carries the reasons spelled out — the operator has to
// see WHY the matcher had nothing, never an empty row. The status stays
// NO_CANDIDATE because /vinculos keys its "sem candidato / Criar produto /
// Ignorar" affordance and its batch-select guard off that value, and apps/web
// is M-06's seam (reconciliation flagged to the hub).
func TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons(t *testing.T) {
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
	if len(candidate.Reasons) == 0 {
		t.Fatal("reasons are empty: M05-C5 requires the operator to see WHY nothing matched")
	}
	if _, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want seller_sku UNAVAILABLE sem correspondência", candidate.Reasons)
	}
	if _, ok := findReason(candidate.Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable); !ok {
		t.Fatalf("reasons=%#v, want ean UNAVAILABLE sem correspondência", candidate.Reasons)
	}
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}

// Case 9 — MARCA-UNAVAILABLE (explicit): marca must appear
// as UNAVAILABLE regardless of band/status (ADR-17, motivo sempre
// visível). Bound to case 1's ALTA/ACCEPT payload with a dedicated
// assertion, per PLAN-M04 §4.
func TestCase9ProviderDeclaredUnavailableAnchorsOnConcordantPayload(t *testing.T) {
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

	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
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
	assertProviderDeclaredUnavailableReasons(t, candidate.Reasons)
}
