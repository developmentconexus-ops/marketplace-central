package domain

import "testing"

func TestRuntimeCapabilityAvailableOnlyWhenExecutable(t *testing.T) {
	capability := RuntimeCapability{
		Code:       RuntimeCapabilityListingRead,
		State:      RuntimeCapabilityStateAvailable,
		Executable: true,
	}

	if !capability.Available() {
		t.Fatalf("expected capability to be available")
	}
}

func TestRuntimeCapabilityUnavailableWhenNotExecutable(t *testing.T) {
	capability := RuntimeCapability{
		Code:       RuntimeCapabilityStockWrite,
		State:      RuntimeCapabilityStateAvailable,
		Executable: false,
	}

	if capability.Available() {
		t.Fatalf("stock write must not be available without executable runtime path")
	}
}

func TestRuntimeCapabilityRejectsEmptyCode(t *testing.T) {
	err := ValidateRuntimeCapability(RuntimeCapability{State: RuntimeCapabilityStateAvailable})
	if err == nil {
		t.Fatalf("expected empty capability code to be invalid")
	}
}

func TestRuntimeCapabilityRunnableWhenDegradedButExecutable(t *testing.T) {
	capability := RuntimeCapability{
		Code:       RuntimeCapabilityFeeQuoteRead,
		State:      RuntimeCapabilityStateDegraded,
		Executable: true,
	}

	if !capability.Runnable() {
		t.Fatalf("expected degraded executable capability to remain runnable")
	}
}
