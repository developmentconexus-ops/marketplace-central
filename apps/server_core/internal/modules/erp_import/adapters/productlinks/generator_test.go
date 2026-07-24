package erpproductlinks

import (
	"context"
	"errors"
	"testing"

	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type fakeCandidateEngine struct {
	input  productlinksapp.GenerateLinkCandidatesInput
	called bool
	result productlinksdomain.LinkCandidateGenerationResult
	err    error
}

func (f *fakeCandidateEngine) GenerateLinkCandidates(_ context.Context, input productlinksapp.GenerateLinkCandidatesInput) (productlinksdomain.LinkCandidateGenerationResult, error) {
	f.called = true
	f.input = input
	return f.result, f.err
}

func TestLinkCandidateGeneratorForwardsInstallationAndDefaultLimit(t *testing.T) {
	engine := &fakeCandidateEngine{}

	err := NewLinkCandidateGenerator(engine).GenerateLinkCandidates(context.Background(), "inst-a")
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	if !engine.called {
		t.Fatal("GenerateLinkCandidates() did not call the candidate engine")
	}
	want := productlinksapp.GenerateLinkCandidatesInput{InstallationID: "inst-a", Limit: 0}
	if engine.input != want {
		t.Fatalf("engine input = %#v, want %#v", engine.input, want)
	}
}

func TestLinkCandidateGeneratorPropagatesEngineError(t *testing.T) {
	wantErr := errors.New("candidate generation failed")
	engine := &fakeCandidateEngine{err: wantErr}

	err := NewLinkCandidateGenerator(engine).GenerateLinkCandidates(context.Background(), "inst-a")
	if !errors.Is(err, wantErr) {
		t.Fatalf("GenerateLinkCandidates() error = %v, want wrapped %v", err, wantErr)
	}
}
