// Package contracts is Catalog's published surface: everything another context
// is allowed to know about a product. The rest lives under internal/ and the
// compiler is what keeps it there.
package contracts

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// IdentifierKind is the sort of identifier, not its authority. An EAN can be
// missing, duplicated, recycled, or simply wrong at the source; it is evidence
// for a linking decision, never the identity of a product.
type IdentifierKind string

const (
	// IdentifierEAN is a barcode.
	IdentifierEAN IdentifierKind = "ean"
	// IdentifierSKU is the seller's own stock code.
	IdentifierSKU IdentifierKind = "sku"
	// IdentifierSourceCode is the source system's primary code.
	IdentifierSourceCode IdentifierKind = "source_code"
)

var (
	// ErrUnknownIdentifierKind is returned for a kind outside the closed set.
	ErrUnknownIdentifierKind = errors.New("catalog: unknown identifier kind")
	// ErrBlank is returned when a required part carries no content.
	ErrBlank = errors.New("catalog: value is empty")
)

// Identifier is one identifying string of a known kind.
type Identifier struct {
	kind  IdentifierKind
	value string
}

// NewIdentifier validates both the kind and the value.
func NewIdentifier(kind IdentifierKind, value string) (Identifier, error) {
	switch kind {
	case IdentifierEAN, IdentifierSKU, IdentifierSourceCode:
	default:
		return Identifier{}, fmt.Errorf("%w: %q", ErrUnknownIdentifierKind, kind)
	}
	v := strings.TrimSpace(value)
	if v == "" {
		return Identifier{}, fmt.Errorf("%w: identifier of kind %q", ErrBlank, kind)
	}
	return Identifier{kind: kind, value: v}, nil
}

// Kind returns the identifier kind.
func (i Identifier) Kind() IdentifierKind { return i.kind }

// Value returns the identifier string.
func (i Identifier) Value() string { return i.value }

// String renders "kind:value".
func (i Identifier) String() string { return string(i.kind) + ":" + i.value }

// SourceProductKey is the address of an object inside a source system, and the
// only place a source system's own code is allowed to live.
//
// The module tree this replaces declared `type InternalProductID int`, and that
// int WAS Sankhya's CODPROD, described as canonical. While that holds, adding a
// second ERP is not additive — it is an identity collision running through
// every contract. Instance is part of the key for the same reason: two Sankhya
// installations are two sources even when the product code matches.
type SourceProductKey struct {
	tenant      tenant.ID
	system      string
	instance    string
	objectKind  string
	externalKey string
}

// NewSourceProductKey validates every part. There is no partial key.
func NewSourceProductKey(t tenant.ID, system, instance, objectKind, externalKey string) (SourceProductKey, error) {
	if t.IsZero() {
		return SourceProductKey{}, fmt.Errorf("%w: tenant", ErrBlank)
	}
	system = strings.TrimSpace(system)
	instance = strings.TrimSpace(instance)
	objectKind = strings.TrimSpace(objectKind)
	externalKey = strings.TrimSpace(externalKey)
	switch {
	case system == "":
		return SourceProductKey{}, fmt.Errorf("%w: source system", ErrBlank)
	case instance == "":
		return SourceProductKey{}, fmt.Errorf("%w: source instance", ErrBlank)
	case objectKind == "":
		return SourceProductKey{}, fmt.Errorf("%w: object kind", ErrBlank)
	case externalKey == "":
		return SourceProductKey{}, fmt.Errorf("%w: external key", ErrBlank)
	}
	return SourceProductKey{
		tenant:      t,
		system:      system,
		instance:    instance,
		objectKind:  objectKind,
		externalKey: externalKey,
	}, nil
}

// Tenant returns the owning tenant.
func (k SourceProductKey) Tenant() tenant.ID { return k.tenant }

// System returns the source system name.
func (k SourceProductKey) System() string { return k.system }

// Instance returns which installation of that system.
func (k SourceProductKey) Instance() string { return k.instance }

// ObjectKind returns what kind of object the key addresses.
func (k SourceProductKey) ObjectKind() string { return k.objectKind }

// ExternalKey returns the source system's own code.
func (k SourceProductKey) ExternalKey() string { return k.externalKey }

// IsZero reports whether this is the zero value rather than a built key.
func (k SourceProductKey) IsZero() bool { return k.system == "" }

// String is the stable storage and log form.
func (k SourceProductKey) String() string {
	if k.IsZero() {
		return ""
	}
	return strings.Join([]string{
		k.tenant.String(), k.system, k.instance, k.objectKind, k.externalKey,
	}, "/")
}
