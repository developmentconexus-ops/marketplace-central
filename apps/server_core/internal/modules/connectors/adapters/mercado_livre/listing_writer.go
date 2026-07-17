package mercadolivre

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type mlListingWriteRequest struct {
	Status     string                    `json:"status,omitempty"`
	Attributes []domain.ListingAttribute `json:"attributes,omitempty"`
}

type mlListingWriteResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (a *CapabilityAdapter) UpdateListing(ctx context.Context, request domain.ListingWriteRequest) (domain.ListingWriteResult, error) {
	request, err := validateListingWriteRequest(request)
	if err != nil {
		return domain.ListingWriteResult{}, err
	}
	account := domain.ProviderAccountRef{TenantID: request.TenantID, InstallationID: request.InstallationID, ProviderCode: "mercado_livre"}
	token, err := a.accessToken(ctx, account)
	if err != nil {
		return domain.ListingWriteResult{}, err
	}
	payload := mlListingWriteRequest{}
	if request.Action == domain.ListingWritePause {
		payload.Status = "paused"
	} else {
		payload.Attributes = request.Attributes
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return domain.ListingWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, "failed to serialize listing write payload")
	}
	resp, raw, err := a.doRawWithIdempotency(ctx, account, token, http.MethodPut, "/items/"+url.PathEscape(request.ListingID), bytes.NewReader(body), request.IdempotencyKey)
	if err != nil {
		return domain.ListingWriteResult{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return domain.ListingWriteResult{}, listingProviderError(resp.StatusCode, raw)
	}

	var provider mlListingWriteResponse
	if err := json.Unmarshal(raw, &provider); err != nil {
		return domain.ListingWriteResult{}, domain.NewCapabilityError(domain.ErrorCode("CONNECTORS_INTERNAL"), "provider listing update returned an invalid response")
	}
	status := strings.ToLower(strings.TrimSpace(provider.Status))
	result := domain.ListingWriteResult{ListingID: request.ListingID, IdempotencyKey: request.IdempotencyKey, Action: request.Action, Result: domain.WriteResultApplied, Message: "provider listing update applied"}
	// Applied only on the expected post-write state; any other or unknown
	// provider status is an honest rejection, never defaulted to success.
	expected := "paused"
	if request.Action == domain.ListingWriteEdit {
		expected = "active"
	}
	if status != expected {
		result.Result, result.Message = domain.WriteResultRejected, "provider listing state is "+firstNonEmpty(status, "unknown")
	}
	return result, nil
}

func validateListingWriteRequest(request domain.ListingWriteRequest) (domain.ListingWriteRequest, error) {
	request.TenantID, request.InstallationID = strings.TrimSpace(request.TenantID), strings.TrimSpace(request.InstallationID)
	request.ListingID, request.IdempotencyKey = strings.TrimSpace(request.ListingID), strings.TrimSpace(request.IdempotencyKey)
	if request.TenantID == "" || request.InstallationID == "" || request.ListingID == "" || request.IdempotencyKey == "" {
		return domain.ListingWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "listing write reference is incomplete")
	}
	if request.Action == domain.ListingWritePause {
		if len(request.Attributes) != 0 {
			return domain.ListingWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "pause does not accept attributes")
		}
		return request, nil
	}
	if request.Action != domain.ListingWriteEdit || len(request.Attributes) == 0 {
		return domain.ListingWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "listing write action or attributes are invalid")
	}
	seen := map[string]struct{}{}
	for i := range request.Attributes {
		request.Attributes[i].ID, request.Attributes[i].ValueName = strings.TrimSpace(request.Attributes[i].ID), strings.TrimSpace(request.Attributes[i].ValueName)
		if request.Attributes[i].ID == "" || request.Attributes[i].ValueName == "" {
			return domain.ListingWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "listing attribute id and value are required")
		}
		if _, duplicate := seen[request.Attributes[i].ID]; duplicate {
			return domain.ListingWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "listing attribute ids must be unique")
		}
		seen[request.Attributes[i].ID] = struct{}{}
	}
	return request, nil
}

func listingProviderError(status int, raw []byte) error {
	message := sanitizedProviderMessage("listing update", raw, status)
	code := domain.ErrorCode("CONNECTORS_INTERNAL")
	switch {
	case status == 401 || status == 403:
		code = domain.ErrCodeProviderAuth
	case status == 429:
		code = domain.ErrCodeProviderRateLimited
	case status >= 400 && status < 500:
		code = domain.ErrCodeProviderValidation
	case status >= 500:
		code = domain.ErrCodeProviderTransient
	}
	return domain.NewCapabilityError(code, message)
}
