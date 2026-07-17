package domain

import (
	"errors"
	"fmt"
)

type ProtocolState string

const (
	ProtocolStateDraft           ProtocolState = "draft"
	ProtocolStatePreviewed       ProtocolState = "previewed"
	ProtocolStateApproved        ProtocolState = "approved"
	ProtocolStateApplying        ProtocolState = "applying"
	ProtocolStateApplied         ProtocolState = "applied"
	ProtocolStatePartiallyFailed ProtocolState = "partially_failed"
	ProtocolStateFailedPreserved ProtocolState = "failed_preserved"
	ProtocolStateCancelled       ProtocolState = "cancelled"
)

type ItemState string

const (
	ItemStatePreviewed ItemState = "previewed"
	ItemStateApproved  ItemState = "approved"
	ItemStateApplying  ItemState = "applying"
	ItemStateApplied   ItemState = "applied"
	ItemStateFailed    ItemState = "failed"
	ItemStateSkipped   ItemState = "skipped"
)

type ProtocolType string

const (
	ProtocolTypePriceUpdate   ProtocolType = "price_update"
	ProtocolTypeStockCorrect  ProtocolType = "stock_correct"
	ProtocolTypeLinkApply     ProtocolType = "link_apply"
	ProtocolTypeListingPause  ProtocolType = "listing_pause"
	ProtocolTypeListingResync ProtocolType = "listing_resync"
	ProtocolTypeListingEdit   ProtocolType = "listing_edit"
	ProtocolTypeListingCreate ProtocolType = "listing_create"
)

var enabledProtocolTypes = map[ProtocolType]struct{}{
	ProtocolTypePriceUpdate:   {},
	ProtocolTypeStockCorrect:  {},
	ProtocolTypeLinkApply:     {},
	ProtocolTypeListingPause:  {},
	ProtocolTypeListingResync: {},
	ProtocolTypeListingEdit:   {},
}

func ProtocolTypeEnabled(t ProtocolType) bool {
	_, ok := enabledProtocolTypes[t]
	return ok
}

const InvalidStateCode = "invalid_state"

var (
	ErrInvalidStateTransition = errors.New(InvalidStateCode)
	ErrItemOutcomesRequired   = errors.New("item outcomes are required")
	ErrItemOutcomeInvalid     = errors.New("item outcome is not terminal")
)

type InvalidStateTransitionError struct {
	Code string
	From ProtocolState
	To   ProtocolState
}

func (e *InvalidStateTransitionError) Error() string {
	return fmt.Sprintf("%s: transition from %q to %q is not allowed", e.Code, e.From, e.To)
}

func (e *InvalidStateTransitionError) Unwrap() error { return ErrInvalidStateTransition }

var legalProtocolTransitions = map[ProtocolState]map[ProtocolState]struct{}{
	ProtocolStateDraft: {
		ProtocolStatePreviewed: {},
		ProtocolStateCancelled: {},
	},
	ProtocolStatePreviewed: {
		ProtocolStatePreviewed: {},
		ProtocolStateApproved:  {},
		ProtocolStateCancelled: {},
	},
	ProtocolStateApproved: {
		ProtocolStateApplying: {},
	},
	ProtocolStateApplying: {
		ProtocolStateApplied:         {},
		ProtocolStatePartiallyFailed: {},
		ProtocolStateFailedPreserved: {},
	},
}

func TransitionProtocolState(from, to ProtocolState) error {
	if targets, ok := legalProtocolTransitions[from]; ok {
		if _, ok := targets[to]; ok {
			return nil
		}
	}
	return &InvalidStateTransitionError{Code: InvalidStateCode, From: from, To: to}
}

func TerminalProtocolState(items []ItemState) (ProtocolState, error) {
	if len(items) == 0 {
		return "", ErrItemOutcomesRequired
	}

	failed := 0
	for _, state := range items {
		switch state {
		case ItemStateApplied, ItemStateSkipped:
		case ItemStateFailed:
			failed++
		default:
			return "", fmt.Errorf("%w: %q", ErrItemOutcomeInvalid, state)
		}
	}

	switch failed {
	case 0:
		return ProtocolStateApplied, nil
	case len(items):
		return ProtocolStateFailedPreserved, nil
	default:
		return ProtocolStatePartiallyFailed, nil
	}
}
