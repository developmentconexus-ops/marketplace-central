// Package tenant carries the tenant identifier. It is a kernel member because
// every context scopes every row and every query by it, with the same meaning
// and the same invariant everywhere: a tenant is never empty and never guessed.
package tenant

import (
	"errors"
	"strings"
)

// ErrEmpty is returned when a tenant identifier carries no content.
var ErrEmpty = errors.New("tenant: identifier is empty")

// ID is a tenant identifier. The field is unexported so no caller outside this
// package can build one by struct literal and skip validation.
type ID struct {
	value string
}

// Parse builds an ID, rejecting empty and blank input.
func Parse(s string) (ID, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return ID{}, ErrEmpty
	}
	return ID{value: trimmed}, nil
}

// String returns the identifier, or the empty string for the zero value.
func (i ID) String() string { return i.value }

// IsZero reports whether this is the zero value rather than a parsed tenant.
func (i ID) IsZero() bool { return i.value == "" }
