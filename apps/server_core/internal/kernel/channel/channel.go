// Package channel names a sales channel and an account on it. The code is data,
// never a Go enum: adding a marketplace must not require editing this package.
package channel

import (
	"errors"
	"strings"
)

var (
	// ErrEmptyCode is returned when a channel code carries no content.
	ErrEmptyCode = errors.New("channel: code is empty")
	// ErrEmptyExternal is returned when an account has no external identifier.
	ErrEmptyExternal = errors.New("channel: external account identifier is empty")
)

// Code identifies a sales channel, lower-cased and trimmed.
type Code struct {
	value string
}

// ParseCode builds a Code from free text, rejecting blank input.
func ParseCode(s string) (Code, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))
	if trimmed == "" {
		return Code{}, ErrEmptyCode
	}
	return Code{value: trimmed}, nil
}

// String returns the code, or the empty string for the zero value.
func (c Code) String() string { return c.value }

// IsZero reports whether this is the zero value.
func (c Code) IsZero() bool { return c.value == "" }

// AccountRef points at one seller account on one channel.
type AccountRef struct {
	channel  Code
	external string
}

// NewAccountRef builds an AccountRef, rejecting a zero code or a blank external id.
func NewAccountRef(c Code, external string) (AccountRef, error) {
	if c.IsZero() {
		return AccountRef{}, ErrEmptyCode
	}
	trimmed := strings.TrimSpace(external)
	if trimmed == "" {
		return AccountRef{}, ErrEmptyExternal
	}
	return AccountRef{channel: c, external: trimmed}, nil
}

// Channel returns the channel this account belongs to.
func (a AccountRef) Channel() Code { return a.channel }

// External returns the channel's own identifier for this account.
func (a AccountRef) External() string { return a.external }
