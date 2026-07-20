package tenant_config

import (
	"testing"

	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

func TestActiveSource_Valid(t *testing.T) {
	cases := []struct {
		name   string
		source ActiveSource
		want   bool
	}{
		{"sankhya", SourceSankhya, true},
		{"xlsx", SourceXLSX, true},
		{"catalogo_cliente", SourceCatalogoCliente, true},
		{"empty", ActiveSource(""), false},
		{"bogus", ActiveSource("bogus"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.source.Valid(); got != tc.want {
				t.Errorf("ActiveSource(%q).Valid() = %v, want %v", string(tc.source), got, tc.want)
			}
		})
	}
}

func TestActiveSource_DefaultKind(t *testing.T) {
	cases := []struct {
		name   string
		source ActiveSource
		want   sourcekind.SourceKind
	}{
		{"sankhya", SourceSankhya, sourcekind.LiveReadThrough},
		{"xlsx", SourceXLSX, sourcekind.UploadSnapshot},
		{"catalogo_cliente", SourceCatalogoCliente, sourcekind.UploadSnapshot},
		{"invalid", ActiveSource("bogus"), sourcekind.SourceKind("")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.source.DefaultKind(); got != tc.want {
				t.Errorf("ActiveSource(%q).DefaultKind() = %q, want %q", string(tc.source), string(got), string(tc.want))
			}
		})
	}
}
