package ports

import (
	"errors"
	"testing"
	"time"
)

func TestMalformedRunCursorRejected(t *testing.T) {
	startedAt := time.Date(2026, 7, 16, 12, 30, 0, 123, time.UTC)
	valid, err := (RunCursor{StartedAt: startedAt, OperationRunID: "run-2"}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeRunCursor(valid)
	if err != nil || !decoded.StartedAt.Equal(startedAt) || decoded.OperationRunID != "run-2" {
		t.Fatalf("decoded=%+v err=%v", decoded, err)
	}

	for _, encoded := range []string{
		"",
		"not-base64",
		"eyJ2IjoyLCJraW5kIjoicnVuIiwic3RhcnRlZF9hdCI6IjIwMjYtMDctMTZUMTI6MzA6MDBaIiwib3BlcmF0aW9uX3J1bl9pZCI6InJ1bi0yIn0=",
		"eyJ2IjoxLCJraW5kIjoibGlzdGluZyIsInN0YXJ0ZWRfYXQiOiIyMDI2LTA3LTE2VDEyOjMwOjAwWiIsIm9wZXJhdGlvbl9ydW5faWQiOiJydW4tMiJ9",
		"eyJ2IjoxLCJraW5kIjoicnVuIiwic3RhcnRlZF9hdCI6ImJhZCIsIm9wZXJhdGlvbl9ydW5faWQiOiJydW4tMiJ9",
		"eyJ2IjoxLCJraW5kIjoicnVuIiwic3RhcnRlZF9hdCI6IjIwMjYtMDctMTZUMTI6MzA6MDBaIiwib3BlcmF0aW9uX3J1bl9pZCI6IiJ9",
	} {
		t.Run(encoded, func(t *testing.T) {
			_, err := DecodeRunCursor(encoded)
			if !errors.Is(err, ErrInvalidRunCursor) {
				t.Fatalf("err=%v", err)
			}
		})
	}
}
