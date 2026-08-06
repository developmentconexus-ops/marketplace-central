// Package sourcekind holds the SourceKind classification shared by the product
// source port (internal_read/ports) and the per-tenant active-source config
// (tenant_config). It is a dedicated, dependency-free package so both sides can
// reference the type without an import cycle (ADR-023 §5).
package sourcekind

// SourceKind distinguishes the two shapes a product data source can take.
//
//   - UploadSnapshot: the source materializes a full snapshot per protocol
//     (xlsx, catalogo_cliente). It has an import protocol and a completed-snapshot
//     lifecycle.
//   - LiveReadThrough: the source is read through live with no protocol/snapshot
//     (sankhya/oracle).
//
// It is intentionally NOT an extension of erp_import's ImportSource enum, which
// only models upload-shaped sources and cannot represent live-read-through
// without misfitting the domain (refactor-inventory-backend.md §9).
type SourceKind string

const (
	// UploadSnapshot is the kind for protocol/snapshot-shaped sources (xlsx,
	// catalogo_cliente).
	UploadSnapshot SourceKind = "upload_snapshot"
	// LiveReadThrough is the kind for protocol-less live sources (sankhya).
	LiveReadThrough SourceKind = "live_read_through"
)

// Valid reports whether k is one of the two ratified kinds. Anything else
// (including the zero value) is invalid — callers never silently coerce an
// unknown kind to a default.
func (k SourceKind) Valid() bool {
	switch k {
	case UploadSnapshot, LiveReadThrough:
		return true
	default:
		return false
	}
}

// String returns the wire form of the kind.
func (k SourceKind) String() string { return string(k) }

// Of classifies a source by its wire name ("xlsx", "catalogo_cliente",
// "sankhya"). The two upload sources are the only snapshot-shaped ones; every
// other source is read through live. Consumers that only carry the source name
// (e.g. a read fact's Source.System) use this instead of restating the list.
func Of(system string) SourceKind {
	switch system {
	case "xlsx", "catalogo_cliente":
		return UploadSnapshot
	default:
		return LiveReadThrough
	}
}
