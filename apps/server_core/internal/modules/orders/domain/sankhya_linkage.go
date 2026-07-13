package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	ErrInvalidSankhyaLinkage  = errors.New("invalid sankhya linkage")
	ErrSankhyaLinkageConflict = errors.New("sankhya linkage conflict")
	ErrSankhyaLinkageNotFound = errors.New("sankhya linkage not found")
)

type LinkageEventType string

const LinkageEventConfirmed LinkageEventType = "confirmed"

type LinkageEvidenceState string

const (
	LinkageEvidenceUnknown LinkageEvidenceState = "unknown"
	LinkageEvidenceExact   LinkageEvidenceState = "exact"
)

type LinkageScope struct {
	TenantID        string
	InstallationID  string
	ProviderOrderID string
}

type InternalDocumentIdentity struct {
	DocumentID int64
}

type InternalDocumentLineIdentity struct {
	DocumentID int64
	LineNumber int
}

type SankhyaLineMapping struct {
	MPCLineID MPCLineID
	Origin    InternalDocumentLineIdentity
}

type LinkageAudit struct {
	EventID               string
	EventType             LinkageEventType
	ActorType             string
	ActorID               string
	Reason                string
	IdempotencyKey        string
	SourceAt              time.Time
	RecordedAt            time.Time
	PreviousEventID       string
	ConfigurationRevision string
	EvidenceState         LinkageEvidenceState
	EvidenceReference     string
}

type SankhyaLinkage struct {
	Scope            LinkageScope
	ExternalOrderKey string
	Header           InternalDocumentIdentity
	Lines            []SankhyaLineMapping
	Audit            LinkageAudit
}

func ExternalOrderKeyFor(scope LinkageScope) string {
	return "ml:v1:" + scope.InstallationID + ":" + scope.ProviderOrderID
}

func (linkage SankhyaLinkage) Validate() error {
	invalid := func(field string) error {
		return fmt.Errorf("%w: %s", ErrInvalidSankhyaLinkage, field)
	}
	if strings.TrimSpace(linkage.Scope.TenantID) == "" {
		return invalid("tenant_id")
	}
	if strings.TrimSpace(linkage.Scope.InstallationID) == "" {
		return invalid("installation_id")
	}
	if strings.TrimSpace(linkage.Scope.ProviderOrderID) == "" {
		return invalid("provider_order_id")
	}
	if linkage.ExternalOrderKey != ExternalOrderKeyFor(linkage.Scope) {
		return invalid("external_order_key")
	}
	if linkage.Header.DocumentID <= 0 {
		return invalid("internal_document_id")
	}
	if len(linkage.Lines) == 0 {
		return invalid("lines")
	}
	lineIDs := make(map[MPCLineID]struct{}, len(linkage.Lines))
	origins := make(map[InternalDocumentLineIdentity]struct{}, len(linkage.Lines))
	for _, line := range linkage.Lines {
		if !line.MPCLineID.Valid() {
			return invalid("mpc_line_id")
		}
		if line.Origin.DocumentID != linkage.Header.DocumentID || line.Origin.LineNumber <= 0 {
			return invalid("internal_document_line")
		}
		if _, duplicate := lineIDs[line.MPCLineID]; duplicate {
			return invalid("duplicate_mpc_line_id")
		}
		if _, duplicate := origins[line.Origin]; duplicate {
			return invalid("duplicate_internal_document_line")
		}
		lineIDs[line.MPCLineID] = struct{}{}
		origins[line.Origin] = struct{}{}
	}
	if strings.TrimSpace(linkage.Audit.EventID) == "" {
		return invalid("event_id")
	}
	if linkage.Audit.EventType != LinkageEventConfirmed {
		return invalid("event_type")
	}
	if strings.TrimSpace(linkage.Audit.ActorType) == "" || strings.TrimSpace(linkage.Audit.ActorID) == "" {
		return invalid("actor")
	}
	if strings.TrimSpace(linkage.Audit.Reason) == "" {
		return invalid("reason")
	}
	if strings.TrimSpace(linkage.Audit.IdempotencyKey) == "" {
		return invalid("idempotency_key")
	}
	if linkage.Audit.SourceAt.IsZero() || linkage.Audit.RecordedAt.IsZero() {
		return invalid("audit_time")
	}
	if strings.TrimSpace(linkage.Audit.ConfigurationRevision) == "" {
		return invalid("configuration_revision")
	}
	if linkage.Audit.EvidenceState != LinkageEvidenceUnknown && linkage.Audit.EvidenceState != LinkageEvidenceExact {
		return invalid("evidence_state")
	}
	return nil
}

// SemanticallyEqual compares the operator intent guarded by an idempotency
// key. Event identity and repository recording time are generated audit facts,
// so a retry may supply different values for those two fields and still return
// the already persisted result.
func (linkage SankhyaLinkage) SemanticallyEqual(other SankhyaLinkage) bool {
	if linkage.Scope != other.Scope || linkage.ExternalOrderKey != other.ExternalOrderKey || linkage.Header != other.Header {
		return false
	}
	leftAudit := linkage.Audit
	rightAudit := other.Audit
	if leftAudit.EventType != rightAudit.EventType ||
		leftAudit.ActorType != rightAudit.ActorType ||
		leftAudit.ActorID != rightAudit.ActorID ||
		leftAudit.Reason != rightAudit.Reason ||
		leftAudit.IdempotencyKey != rightAudit.IdempotencyKey ||
		!samePersistedTime(leftAudit.SourceAt, rightAudit.SourceAt) ||
		leftAudit.PreviousEventID != rightAudit.PreviousEventID ||
		leftAudit.ConfigurationRevision != rightAudit.ConfigurationRevision ||
		leftAudit.EvidenceState != rightAudit.EvidenceState ||
		leftAudit.EvidenceReference != rightAudit.EvidenceReference {
		return false
	}
	left := canonicalLineMappings(linkage.Lines)
	right := canonicalLineMappings(other.Lines)
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func samePersistedTime(left, right time.Time) bool {
	return left.UTC().Truncate(time.Microsecond).Equal(right.UTC().Truncate(time.Microsecond))
}

func canonicalLineMappings(lines []SankhyaLineMapping) []SankhyaLineMapping {
	result := append([]SankhyaLineMapping(nil), lines...)
	sort.Slice(result, func(left, right int) bool {
		if result[left].MPCLineID != result[right].MPCLineID {
			return result[left].MPCLineID < result[right].MPCLineID
		}
		if result[left].Origin.DocumentID != result[right].Origin.DocumentID {
			return result[left].Origin.DocumentID < result[right].Origin.DocumentID
		}
		return result[left].Origin.LineNumber < result[right].Origin.LineNumber
	})
	return result
}
