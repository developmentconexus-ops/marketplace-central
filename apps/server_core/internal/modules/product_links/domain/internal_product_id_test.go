package domain

import "testing"

func TestInternalProductIDMustBePositive(t *testing.T) {
	for _, id := range []int{-1, 0} {
		if err := ValidateInternalProductID(id); err == nil {
			t.Fatalf("ValidateInternalProductID(%d) unexpectedly succeeded", id)
		}
	}
	if err := ValidateInternalProductID(1); err != nil {
		t.Fatalf("ValidateInternalProductID(1) = %v", err)
	}
}
