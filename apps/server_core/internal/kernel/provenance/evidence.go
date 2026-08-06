// Package provenance answers "how do we know this?" for every fact the platform
// holds. A fact without evidence is an assertion, and the platform does not
// carry assertions: kernel/fact refuses to build one.
package provenance

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrIncomplete is returned when any part of the evidence is missing.
var ErrIncomplete = errors.New("provenance: incomplete evidence")

// Evidence records where an observation came from and when we saw it.
//
// ObservedAt is when WE collected it. It is deliberately not the same thing as
// when the source says the fact started to hold, nor the legal period the fact
// is valid for — those are separate fields on the facts themselves (§6).
type Evidence struct {
	system      string
	objectKind  string
	externalKey string
	observedAt  time.Time
	payloadHash string
}

// NewEvidence builds Evidence, refusing anything partial.
func NewEvidence(system, objectKind, externalKey string, observedAt time.Time, payloadHash string) (Evidence, error) {
	system = strings.TrimSpace(system)
	objectKind = strings.TrimSpace(objectKind)
	externalKey = strings.TrimSpace(externalKey)
	payloadHash = strings.TrimSpace(payloadHash)

	switch {
	case system == "":
		return Evidence{}, fmt.Errorf("%w: system is empty", ErrIncomplete)
	case objectKind == "":
		return Evidence{}, fmt.Errorf("%w: object kind is empty", ErrIncomplete)
	case externalKey == "":
		return Evidence{}, fmt.Errorf("%w: external key is empty", ErrIncomplete)
	case observedAt.IsZero():
		return Evidence{}, fmt.Errorf("%w: observed_at is the zero time", ErrIncomplete)
	case payloadHash == "":
		return Evidence{}, fmt.Errorf("%w: payload hash is empty", ErrIncomplete)
	}

	return Evidence{
		system:      system,
		objectKind:  objectKind,
		externalKey: externalKey,
		observedAt:  observedAt.UTC(),
		payloadHash: payloadHash,
	}, nil
}

// System returns the source system, e.g. "sankhya" or "mercadolivre".
func (e Evidence) System() string { return e.system }

// ObjectKind returns what kind of thing was observed, e.g. "product".
func (e Evidence) ObjectKind() string { return e.objectKind }

// ExternalKey returns the source system's own key for the object.
func (e Evidence) ExternalKey() string { return e.externalKey }

// ObservedAt returns, in UTC, when we collected the observation.
func (e Evidence) ObservedAt() time.Time { return e.observedAt }

// PayloadHash returns the hash of the raw payload this was read from.
func (e Evidence) PayloadHash() string { return e.payloadHash }

// Ref is the short human-readable pointer, "system/kind:key".
func (e Evidence) Ref() string {
	if e.IsZero() {
		return ""
	}
	return e.system + "/" + e.objectKind + ":" + e.externalKey
}

// IsZero reports whether this is the zero value rather than built evidence.
func (e Evidence) IsZero() bool { return e.system == "" }
