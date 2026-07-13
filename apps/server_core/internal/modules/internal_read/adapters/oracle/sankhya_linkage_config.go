package oracle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	sankhyaLinkageEnabledEnv               = "MPC_ASSISTED_SANKHYA_LINKAGE_ENABLED"
	sankhyaLinkageSchemaEnv                = "MPC_ASSISTED_SANKHYA_LINKAGE_SCHEMA"
	sankhyaLinkageHeaderFieldEnv           = "MPC_ASSISTED_SANKHYA_LINKAGE_HEADER_FIELD"
	sankhyaLinkageConfigurationRevisionEnv = "MPC_ASSISTED_SANKHYA_LINKAGE_CONFIGURATION_REVISION"
	sankhyaLinkageUniquenessAttestationEnv = "MPC_ASSISTED_SANKHYA_LINKAGE_UNIQUENESS_ATTESTATION_ID"
	sankhyaLinkageCandidateLimitEnv        = "MPC_ASSISTED_SANKHYA_LINKAGE_CANDIDATE_LIMIT"
	sankhyaLinkageCandidateLineLimitEnv    = "MPC_ASSISTED_SANKHYA_LINKAGE_CANDIDATE_LINE_LIMIT"
	sankhyaLinkageLineageLimitEnv          = "MPC_ASSISTED_SANKHYA_LINKAGE_LINEAGE_LIMIT"
)

var ErrSankhyaLinkageRuntimeUnavailable = errors.New("assisted sankhya linkage runtime unavailable")

// SankhyaLinkageRuntimeConfig is the complete, server-owned activation
// contract. The evidence reference is an opaque digest of the exact validated
// deployment facts; it never exposes the configured Oracle identifier or the
// administrator attestation.
type SankhyaLinkageRuntimeConfig struct {
	ReaderConfig      SankhyaLinkageConfig
	EvidenceReference string
}

func LoadSankhyaLinkageRuntimeConfigFromEnv(getenv func(string) string) (SankhyaLinkageRuntimeConfig, error) {
	if getenv == nil {
		return SankhyaLinkageRuntimeConfig{}, runtimeConfigError(sankhyaLinkageEnabledEnv)
	}
	enabled, err := requiredBool(getenv, sankhyaLinkageEnabledEnv)
	if err != nil || !enabled {
		return SankhyaLinkageRuntimeConfig{}, runtimeConfigError(sankhyaLinkageEnabledEnv)
	}

	config := SankhyaLinkageConfig{
		Schema:                  requiredValue(getenv, sankhyaLinkageSchemaEnv),
		HeaderFieldName:         requiredValue(getenv, sankhyaLinkageHeaderFieldEnv),
		ConfigurationRevision:   requiredValue(getenv, sankhyaLinkageConfigurationRevisionEnv),
		ExpectedOriginTOP:       sankhyaOriginTOP,
		ExpectedDestinationTOP:  sankhyaDestinationTOP,
		UniquenessAttestationID: requiredValue(getenv, sankhyaLinkageUniquenessAttestationEnv),
	}
	config.CandidateLimit, err = requiredPositiveInt(getenv, sankhyaLinkageCandidateLimitEnv)
	if err != nil {
		return SankhyaLinkageRuntimeConfig{}, runtimeConfigError(sankhyaLinkageCandidateLimitEnv)
	}
	config.CandidateLineLimit, err = requiredPositiveInt(getenv, sankhyaLinkageCandidateLineLimitEnv)
	if err != nil {
		return SankhyaLinkageRuntimeConfig{}, runtimeConfigError(sankhyaLinkageCandidateLineLimitEnv)
	}
	config.LineageLimit, err = requiredPositiveInt(getenv, sankhyaLinkageLineageLimitEnv)
	if err != nil {
		return SankhyaLinkageRuntimeConfig{}, runtimeConfigError(sankhyaLinkageLineageLimitEnv)
	}
	if err := validateSankhyaLinkageConfig(config); err != nil {
		return SankhyaLinkageRuntimeConfig{}, fmt.Errorf("%w: invalid deployment contract", ErrSankhyaLinkageRuntimeUnavailable)
	}

	return SankhyaLinkageRuntimeConfig{
		ReaderConfig:      config,
		EvidenceReference: sankhyaLinkageEvidenceReference(config),
	}, nil
}

func requiredValue(getenv func(string) string, key string) string {
	value := getenv(key)
	if value != strings.TrimSpace(value) {
		return ""
	}
	return value
}

func requiredBool(getenv func(string) string, key string) (bool, error) {
	raw := requiredValue(getenv, key)
	if raw == "" {
		return false, errors.New("missing boolean")
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, errors.New("malformed boolean")
	}
	return value, nil
}

func requiredPositiveInt(getenv func(string) string, key string) (int, error) {
	raw := requiredValue(getenv, key)
	if raw == "" || strings.HasPrefix(raw, "+") {
		return 0, errors.New("missing integer")
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("malformed integer")
	}
	return value, nil
}

func sankhyaLinkageEvidenceReference(config SankhyaLinkageConfig) string {
	payload := strings.Join([]string{
		"sankhya-linkage-v1",
		config.Schema,
		config.HeaderFieldName,
		config.ConfigurationRevision,
		config.UniquenessAttestationID,
		strconv.Itoa(config.CandidateLimit),
		strconv.Itoa(config.CandidateLineLimit),
		strconv.Itoa(config.LineageLimit),
	}, "\x00")
	digest := sha256.Sum256([]byte(payload))
	return "sankhya-linkage:v1:" + hex.EncodeToString(digest[:])
}

func runtimeConfigError(key string) error {
	return fmt.Errorf("%w: %s", ErrSankhyaLinkageRuntimeUnavailable, key)
}
