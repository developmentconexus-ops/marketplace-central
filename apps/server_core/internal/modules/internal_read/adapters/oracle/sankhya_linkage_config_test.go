package oracle

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadSankhyaLinkageRuntimeConfigRequiresEveryExplicitSetting(t *testing.T) {
	base := validSankhyaLinkageEnvironment()
	keys := []string{
		sankhyaLinkageEnabledEnv,
		sankhyaLinkageSchemaEnv,
		sankhyaLinkageHeaderFieldEnv,
		sankhyaLinkageConfigurationRevisionEnv,
		sankhyaLinkageUniquenessAttestationEnv,
		sankhyaLinkageCandidateLimitEnv,
		sankhyaLinkageCandidateLineLimitEnv,
		sankhyaLinkageLineageLimitEnv,
	}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			values := cloneEnvironment(base)
			delete(values, key)
			if _, err := LoadSankhyaLinkageRuntimeConfigFromEnv(values.get); !errors.Is(err, ErrSankhyaLinkageRuntimeUnavailable) {
				t.Fatalf("error = %v, want unavailable", err)
			}
		})
	}
}

func TestLoadSankhyaLinkageRuntimeConfigRejectsDisabledMalformedAndOutOfRangeValues(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{sankhyaLinkageEnabledEnv, "false"},
		{sankhyaLinkageEnabledEnv, "yes"},
		{sankhyaLinkageSchemaEnv, "metalprd"},
		{sankhyaLinkageHeaderFieldEnv, " AD_MPC_ORDER_KEY"},
		{sankhyaLinkageConfigurationRevisionEnv, " "},
		{sankhyaLinkageCandidateLimitEnv, "0"},
		{sankhyaLinkageCandidateLimitEnv, "1001"},
		{sankhyaLinkageCandidateLineLimitEnv, "501"},
		{sankhyaLinkageLineageLimitEnv, "501"},
		{sankhyaLinkageLineageLimitEnv, "1.5"},
	}
	for _, test := range tests {
		t.Run(test.key+"="+test.value, func(t *testing.T) {
			values := cloneEnvironment(validSankhyaLinkageEnvironment())
			values[test.key] = test.value
			if _, err := LoadSankhyaLinkageRuntimeConfigFromEnv(values.get); !errors.Is(err, ErrSankhyaLinkageRuntimeUnavailable) {
				t.Fatalf("error = %v, want unavailable", err)
			}
		})
	}
}

func TestLoadSankhyaLinkageRuntimeConfigUsesFixedTOPsAndOpaqueStableEvidence(t *testing.T) {
	values := validSankhyaLinkageEnvironment()
	got, err := LoadSankhyaLinkageRuntimeConfigFromEnv(values.get)
	if err != nil {
		t.Fatalf("LoadSankhyaLinkageRuntimeConfigFromEnv() error = %v", err)
	}
	if got.ReaderConfig.ExpectedOriginTOP != 313 || got.ReaderConfig.ExpectedDestinationTOP != 306 {
		t.Fatalf("TOPs = %d/%d", got.ReaderConfig.ExpectedOriginTOP, got.ReaderConfig.ExpectedDestinationTOP)
	}
	if !strings.HasPrefix(got.EvidenceReference, "sankhya-linkage:v1:") ||
		strings.Contains(got.EvidenceReference, values[sankhyaLinkageHeaderFieldEnv]) ||
		strings.Contains(got.EvidenceReference, values[sankhyaLinkageUniquenessAttestationEnv]) {
		t.Fatalf("unsafe evidence reference = %q", got.EvidenceReference)
	}
	again, err := LoadSankhyaLinkageRuntimeConfigFromEnv(values.get)
	if err != nil || again.EvidenceReference != got.EvidenceReference {
		t.Fatalf("evidence references differ: %q/%q err=%v", got.EvidenceReference, again.EvidenceReference, err)
	}
	values[sankhyaLinkageConfigurationRevisionEnv] = "rev-2"
	changed, err := LoadSankhyaLinkageRuntimeConfigFromEnv(values.get)
	if err != nil || changed.EvidenceReference == got.EvidenceReference {
		t.Fatalf("changed evidence reference = %q err=%v", changed.EvidenceReference, err)
	}
}

type linkageEnvironment map[string]string

func (e linkageEnvironment) get(key string) string { return e[key] }

func validSankhyaLinkageEnvironment() linkageEnvironment {
	return linkageEnvironment{
		sankhyaLinkageEnabledEnv:               "true",
		sankhyaLinkageSchemaEnv:                "METALPRD",
		sankhyaLinkageHeaderFieldEnv:           "AD_MPC_ORDER_KEY",
		sankhyaLinkageConfigurationRevisionEnv: "rev-1",
		sankhyaLinkageUniquenessAttestationEnv: "attestation-admin-1",
		sankhyaLinkageCandidateLimitEnv:        "20",
		sankhyaLinkageCandidateLineLimitEnv:    "100",
		sankhyaLinkageLineageLimitEnv:          "100",
	}
}

func cloneEnvironment(source linkageEnvironment) linkageEnvironment {
	result := linkageEnvironment{}
	for key, value := range source {
		result[key] = value
	}
	return result
}
