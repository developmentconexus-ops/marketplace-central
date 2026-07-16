package productlinks

import (
	"context"
	"errors"
	"testing"

	mutationsapp "marketplace-central/apps/server_core/internal/modules/mutations/application"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type fakeResolver struct {
	approveCalls int
	manualCalls  int
	rejectCalls  int
	result       productlinksdomain.ProductLinkResolutionResult
	err          error
	workflows    []productlinksdomain.ProductLinkWorkflowItem
}

func (f *fakeResolver) ApproveCandidate(context.Context, productlinksapp.ApproveCandidateInput) (productlinksdomain.ProductLinkResolutionResult, error) {
	f.approveCalls++
	return f.result, f.err
}

func (f *fakeResolver) ManualResolve(context.Context, productlinksapp.ManualResolveInput) (productlinksdomain.ProductLinkResolutionResult, error) {
	f.manualCalls++
	return f.result, f.err
}

func (f *fakeResolver) RejectListing(context.Context, productlinksapp.RejectListingInput) (productlinksdomain.ProductLinkResolutionResult, error) {
	f.rejectCalls++
	return f.result, f.err
}

func (f *fakeResolver) ListLinkWorkflows(context.Context, productlinksapp.ListLinkWorkflowsInput) ([]productlinksdomain.ProductLinkWorkflowItem, error) {
	return f.workflows, nil
}

func TestWriterDelegatesEveryLinkActionToPublishedApplicationAPI(t *testing.T) {
	resolver := &fakeResolver{result: productlinksdomain.ProductLinkResolutionResult{Link: productlinksdomain.ProductLink{State: productlinksdomain.ProductLinkStateResolved}}}
	writer := NewWriter(resolver)

	cases := []struct {
		name   string
		input  mutationsports.LinkageInput
		called *int
	}{
		{name: "approve candidate", input: mutationsports.LinkageInput{InstallationID: "inst-1", Action: mutationsports.LinkageActionApproveCandidate, CandidateID: "cand-1"}, called: &resolver.approveCalls},
		{name: "manual resolve", input: mutationsports.LinkageInput{InstallationID: "inst-1", Action: mutationsports.LinkageActionManualResolve, ProviderCode: "mercado_livre", ProviderItemID: "MLB-1", InternalProductID: 10}, called: &resolver.manualCalls},
		{name: "reject listing", input: mutationsports.LinkageInput{InstallationID: "inst-1", Action: mutationsports.LinkageActionRejectListing, ProviderCode: "mercado_livre", ProviderItemID: "MLB-1"}, called: &resolver.rejectCalls},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outcome, err := writer.ApplyLink(context.Background(), tc.input)
			if err != nil {
				t.Fatalf("ApplyLink() error = %v", err)
			}
			if outcome.State != mutationsdomain.ItemStateApplied {
				t.Fatalf("state = %q, want applied", outcome.State)
			}
			if *tc.called != 1 {
				t.Fatalf("published API calls = %d, want 1", *tc.called)
			}
		})
	}
}

func TestWriterMapsIdenticalAndConflictingResolutionOutcomes(t *testing.T) {
	productID := 10
	resolver := &fakeResolver{workflows: []productlinksdomain.ProductLinkWorkflowItem{{
		Identity:    productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB-1"},
		CurrentLink: &productlinksdomain.ProductLink{State: productlinksdomain.ProductLinkStateResolved, InternalProductID: &productID},
	}}}
	writer := NewWriter(resolver)
	outcome, err := writer.ApplyLink(context.Background(), mutationsports.LinkageInput{Action: mutationsports.LinkageActionManualResolve, InstallationID: "inst-1", ProviderCode: "mercado_livre", ProviderItemID: "MLB-1", InternalProductID: productID})
	if err != nil || outcome.State != mutationsdomain.ItemStateSkipped || resolver.manualCalls != 0 {
		t.Fatalf("identical outcome = %#v, calls=%d, err=%v", outcome, resolver.manualCalls, err)
	}

	otherProductID := 11
	outcome, err = writer.ApplyLink(context.Background(), mutationsports.LinkageInput{Action: mutationsports.LinkageActionManualResolve, InstallationID: "inst-1", ProviderCode: "mercado_livre", ProviderItemID: "MLB-1", InternalProductID: otherProductID})
	if err != nil || outcome.Failure == nil || outcome.Failure.Code != mutationsdomain.FailureCodeConflictRemoteChanged || resolver.manualCalls != 0 {
		t.Fatalf("conflicting outcome = %#v, err=%v", outcome, err)
	}
}

func TestWriterPropagatesResolverError(t *testing.T) {
	want := errors.New("resolver unavailable")
	writer := NewWriter(&fakeResolver{err: want})
	_, err := writer.ApplyLink(context.Background(), mutationsports.LinkageInput{Action: mutationsports.LinkageActionRejectListing, InstallationID: "inst-1", ProviderCode: "mercado_livre", ProviderItemID: "MLB-1"})
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func TestUnresolvedLinkGateReturnsIC03CodeFromAdapter(t *testing.T) {
	writer := NewWriter(&fakeResolver{})
	err := mutationsapp.RequireResolvedLink(context.Background(), mutationsdomain.ProtocolTypePriceUpdate, "inst-1~MLB-1~-", writer.HasResolvedLink)
	var gateErr *mutationsapp.GateError
	if !errors.As(err, &gateErr) || gateErr.Code != mutationsdomain.FailureCodeLinkUnresolved {
		t.Fatalf("error = %v, want link_unresolved", err)
	}
}

func TestHasResolvedLinkIsScopedToTheGatedListing(t *testing.T) {
	productID := 10
	resolver := &fakeResolver{workflows: []productlinksdomain.ProductLinkWorkflowItem{
		{
			Identity: productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB-unlinked"},
		},
		{
			Identity:    productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB-linked"},
			CurrentLink: &productlinksdomain.ProductLink{State: productlinksdomain.ProductLinkStateResolved, InternalProductID: &productID},
		},
	}}
	writer := NewWriter(resolver)

	resolved, err := writer.HasResolvedLink(context.Background(), "inst-1~MLB-unlinked~-")
	if err != nil || resolved {
		t.Fatalf("unlinked listing resolved = %v, err = %v; a sibling's resolved link must not satisfy the gate", resolved, err)
	}

	resolved, err = writer.HasResolvedLink(context.Background(), "inst-1~MLB-linked~-")
	if err != nil || !resolved {
		t.Fatalf("linked listing resolved = %v, err = %v", resolved, err)
	}

	resolved, err = writer.HasResolvedLink(context.Background(), "inst-1~MLB-missing~-")
	if err != nil || resolved {
		t.Fatalf("missing listing resolved = %v, err = %v", resolved, err)
	}

	if _, err = writer.HasResolvedLink(context.Background(), "not-a-listing-id"); err == nil {
		t.Fatal("malformed listing id must error, not default to unresolved silently")
	}
}
