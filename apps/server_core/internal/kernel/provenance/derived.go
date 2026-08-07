package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Derived builds evidence for a value nobody observed.
//
// NewEvidence cannot serve: it demands one system, one external key and one
// observation time, and a derived value has as many of each as it has inputs.
// Forcing it through NewEvidence would mean picking one input and silently
// discarding the rest — the provenance equivalent of an unknown becoming zero.
//
// ObservedAt is the OLDEST input's time, not the newest and not now: a derived
// number is exactly as fresh as its stalest ingredient.
func Derived(method string, from ...Evidence) (Evidence, error) {
	method = strings.TrimSpace(method)
	if method == "" {
		return Evidence{}, fmt.Errorf("%w: derivation method", ErrIncomplete)
	}
	if len(from) == 0 {
		return Evidence{}, fmt.Errorf("%w: derivation has no inputs", ErrIncomplete)
	}
	refs := make([]string, 0, len(from))
	oldest := time.Time{}
	for _, e := range from {
		if e.IsZero() {
			return Evidence{}, fmt.Errorf("%w: derivation input", ErrIncomplete)
		}
		refs = append(refs, e.Ref())
		if oldest.IsZero() || e.ObservedAt().Before(oldest) {
			oldest = e.ObservedAt()
		}
	}
	sort.Strings(refs)
	h := sha256.New()
	fmt.Fprintf(h, "method=%s\n", method)
	for _, r := range refs {
		fmt.Fprintf(h, "from=%s\n", r)
	}
	return NewEvidence("derived", method, strings.Join(refs, "+"), oldest,
		"sha256:"+hex.EncodeToString(h.Sum(nil)))
}
