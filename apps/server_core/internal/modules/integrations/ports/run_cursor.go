package ports

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"
)

var ErrInvalidRunCursor = errors.New("invalid_run_cursor")

type InvalidRunCursorError struct{}

func (*InvalidRunCursorError) Error() string { return ErrInvalidRunCursor.Error() }
func (*InvalidRunCursorError) Unwrap() error { return ErrInvalidRunCursor }

type RunCursor struct {
	StartedAt      time.Time
	OperationRunID string
}

type runCursorEnvelope struct {
	Version        int    `json:"v"`
	Kind           string `json:"kind"`
	StartedAt      string `json:"started_at"`
	OperationRunID string `json:"operation_run_id"`
}

func (c RunCursor) IsFirstPage() bool { return c.StartedAt.IsZero() && c.OperationRunID == "" }

func (c RunCursor) Encode() (string, error) {
	if c.StartedAt.IsZero() || c.OperationRunID == "" {
		return "", &InvalidRunCursorError{}
	}
	payload, err := json.Marshal(runCursorEnvelope{
		Version:        1,
		Kind:           "run",
		StartedAt:      c.StartedAt.UTC().Format(time.RFC3339Nano),
		OperationRunID: c.OperationRunID,
	})
	if err != nil {
		return "", &InvalidRunCursorError{}
	}
	return base64.StdEncoding.EncodeToString(payload), nil
}

func DecodeRunCursor(encoded string) (RunCursor, error) {
	if encoded == "" {
		return RunCursor{}, &InvalidRunCursorError{}
	}
	payload, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil || base64.StdEncoding.EncodeToString(payload) != encoded {
		return RunCursor{}, &InvalidRunCursorError{}
	}
	var envelope runCursorEnvelope
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return RunCursor{}, &InvalidRunCursorError{}
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return RunCursor{}, &InvalidRunCursorError{}
	}
	startedAt, err := time.Parse(time.RFC3339Nano, envelope.StartedAt)
	if err != nil || envelope.Version != 1 || envelope.Kind != "run" || envelope.OperationRunID == "" || startedAt.IsZero() {
		return RunCursor{}, &InvalidRunCursorError{}
	}
	return RunCursor{StartedAt: startedAt, OperationRunID: envelope.OperationRunID}, nil
}
