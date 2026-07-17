package mercadolivre

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

var (
	priceAmountPattern    = regexp.MustCompile(`^(0|[1-9][0-9]*)(\.[0-9]+)?$`)
	currencyPattern       = regexp.MustCompile(`^[A-Z]{3}$`)
	providerSecretPattern = regexp.MustCompile(`(?i)(token|secret|authorization|password|client_secret)\s*[:=]\s*("[^"]*"|'[^']*'|[^\s,;]+)`)
	bearerPattern         = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)
)

type mlPriceWriteRequest struct {
	Price      json.RawMessage `json:"price"`
	CurrencyID string          `json:"currency_id"`
}

type mlProviderErrorResponse struct {
	Message          string `json:"message"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// UpdatePrice sends an absolute price update to the Mercado Livre item API.
// The canonical request remains in the connectors domain; ML-specific JSON is
// deliberately kept in this adapter.
func (a *CapabilityAdapter) UpdatePrice(ctx context.Context, request domain.PriceWriteRequest) (domain.PriceWriteResult, error) {
	request, err := validatePriceWriteRequest(request)
	if err != nil {
		return domain.PriceWriteResult{}, err
	}

	accountRef := domain.ProviderAccountRef{
		TenantID:       strings.TrimSpace(request.TenantID),
		InstallationID: strings.TrimSpace(request.InstallationID),
		ProviderCode:   "mercado_livre",
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.PriceWriteResult{}, err
	}

	payload, err := json.Marshal(mlPriceWriteRequest{
		Price:      json.RawMessage(request.Price.Amount),
		CurrencyID: request.Price.Currency,
	})
	if err != nil {
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, "failed to serialize price write payload")
	}

	resp, rawBody, err := a.doRawWithIdempotency(ctx, accountRef, token, http.MethodPut, "/items/"+url.PathEscape(request.ListingID), bytes.NewReader(payload), request.IdempotencyKey)
	if err != nil {
		return domain.PriceWriteResult{}, err
	}
	if resp.StatusCode < http.StatusMultipleChoices && resp.StatusCode >= http.StatusOK {
		return domain.PriceWriteResult{
			ListingID:      request.ListingID,
			IdempotencyKey: request.IdempotencyKey,
			Price:          request.Price,
			Result:         domain.WriteResultApplied,
			Message:        "provider write applied",
		}, nil
	}

	message := sanitizedProviderMessage("price update", rawBody, resp.StatusCode)
	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderAuth, message)
	case resp.StatusCode == http.StatusTooManyRequests:
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderRateLimited, message)
	case resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError:
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, message)
	case resp.StatusCode >= http.StatusInternalServerError:
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrCodeProviderTransient, message)
	default:
		return domain.PriceWriteResult{}, domain.NewCapabilityError(domain.ErrorCode("CONNECTORS_INTERNAL"), message)
	}
}

func validatePriceWriteRequest(request domain.PriceWriteRequest) (domain.PriceWriteRequest, error) {
	request.TenantID = strings.TrimSpace(request.TenantID)
	request.InstallationID = strings.TrimSpace(request.InstallationID)
	request.ListingID = strings.TrimSpace(request.ListingID)
	request.IdempotencyKey = strings.TrimSpace(request.IdempotencyKey)
	request.Price.Amount = strings.TrimSpace(request.Price.Amount)
	request.Price.Currency = strings.TrimSpace(request.Price.Currency)

	if request.TenantID == "" || request.InstallationID == "" || request.ListingID == "" || request.IdempotencyKey == "" {
		return domain.PriceWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price write reference is incomplete")
	}
	if !priceAmountPattern.MatchString(request.Price.Amount) {
		return domain.PriceWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price amount must be a positive decimal")
	}
	amount, ok := new(big.Rat).SetString(request.Price.Amount)
	if !ok || amount.Sign() <= 0 {
		return domain.PriceWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price amount must be a positive decimal")
	}
	if !currencyPattern.MatchString(request.Price.Currency) {
		return domain.PriceWriteRequest{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price currency must be an uppercase ISO code")
	}
	return request, nil
}

func sanitizedProviderMessage(operation string, rawBody []byte, statusCode int) string {
	message := ""
	var response mlProviderErrorResponse
	if err := json.Unmarshal(rawBody, &response); err == nil {
		message = firstNonEmpty(response.Message, response.ErrorDescription, response.Error)
	}
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(statusCode)
	}
	if message == "" {
		message = "provider rejected " + operation
	}
	message = strings.Join(strings.Fields(message), " ")
	message = providerSecretPattern.ReplaceAllString(message, "$1=[REDACTED]")
	message = bearerPattern.ReplaceAllString(message, "Bearer [REDACTED]")
	if len(message) > 300 {
		message = message[:300]
	}
	return "provider " + operation + ": " + message
}
