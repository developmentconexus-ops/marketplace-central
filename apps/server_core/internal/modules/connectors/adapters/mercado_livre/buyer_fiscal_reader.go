package mercadolivre

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// mlBillingInfoResponse mirrors the documented GET /orders/billing-info/{SITE_ID}/{id}
// response (ML docs: faturamento). Only the ERP-registration fields are decoded; the
// MLA-only blocks (taxes.iibb, attributes, vat) are intentionally ignored (YAGNI — dispatch
// scope). The `invoice_type` field is NOT decoded: ML officially removed it (the integrator
// derives invoice handling from the identification document, not a provider field).
type mlBillingInfoResponse struct {
	Buyer *mlBillingInfoBuyer `json:"buyer"`
}

type mlBillingInfoBuyer struct {
	BillingInfo *mlBillingInfoData `json:"billing_info"`
}

type mlBillingInfoData struct {
	Name           string                   `json:"name"`
	LastName       string                   `json:"last_name"`
	Identification *mlBillingIdentification `json:"identification"`
	Address        *mlBillingAddress        `json:"address"`
}

type mlBillingIdentification struct {
	// Type is OPAQUE (ADR-17): the MLB literal is doc-unverified (documented examples are
	// MLA DNI/CUIT). The decoder renders it as-is and never maps it to a CPF/CNPJ enum.
	Type   string `json:"type"`
	Number string `json:"number"`
}

type mlBillingAddress struct {
	StreetName   string          `json:"street_name"`
	StreetNumber string          `json:"street_number"`
	CityName     string          `json:"city_name"`
	State        *mlBillingState `json:"state"`
	ZipCode      string          `json:"zip_code"`
	CountryID    string          `json:"country_id"`
}

type mlBillingState struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetBuyerFiscalInfo resolves a buyer's fiscal identity via the documented two-step ML flow
// (docs: faturamento): GET /orders/{id} -> buyer.billing_info.id -> GET
// /orders/billing-info/{SITE_ID}/{billing_info_id}. A buyer without billing linkage (no
// billing_info.id) or a 404 on the billing-info call degrades to an empty, honest-absence
// BuyerFiscalInfo (HasData() == false), never an error — the order still renders.
func (a *CapabilityAdapter) GetBuyerFiscalInfo(ctx context.Context, accountRef domain.ProviderAccountRef, providerOrderID string) (domain.BuyerFiscalInfo, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.BuyerFiscalInfo{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.BuyerFiscalInfo{}, mapPricingReaderError(err)
	}
	return a.getBuyerFiscalInfo(ctx, accountRef, token, strings.TrimSpace(providerOrderID))
}

func (a *CapabilityAdapter) getBuyerFiscalInfo(ctx context.Context, accountRef domain.ProviderAccountRef, token, providerOrderID string) (domain.BuyerFiscalInfo, error) {
	// Step 1: GET /orders/{id} -> buyer.billing_info.id (Bearer only). A failure here is a
	// real read error (the order itself could not be read), so it propagates.
	var order mlOrderResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/orders/"+url.PathEscape(providerOrderID), nil, &order); err != nil {
		return domain.BuyerFiscalInfo{}, mapPricingReaderError(err)
	}

	billingInfoID := ""
	if order.Buyer != nil && order.Buyer.BillingInfo != nil {
		billingInfoID = strings.TrimSpace(order.Buyer.BillingInfo.ID)
	}
	if billingInfoID == "" {
		// Buyer carries no billing linkage — honest absence, skip the step-2 call entirely.
		return domain.BuyerFiscalInfo{FetchedAt: a.now()}, nil
	}

	// Step 2: GET /orders/billing-info/{SITE_ID}/{billing_info_id}. Bearer only — NO
	// x-format-new header (that is a shipments-only requirement). Site is the order's own
	// site_id, falling back to the account-configured site before any default so a non-BR
	// tenant is not silently relabelled (unknown != default).
	site := firstNonEmpty(strings.TrimSpace(order.SiteID), a.siteID)
	var billing mlBillingInfoResponse
	path := "/orders/billing-info/" + url.PathEscape(site) + "/" + url.PathEscape(billingInfoID)
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &billing); err != nil {
		// A 404 (a buyer without billing data — ML returns 404 as a normal "none"
		// condition) or a payload we cannot decode is honest absence: degrade to an empty
		// fiscal info so the order still renders, no WARN. Only these two non-fatal codes
		// degrade; auth/transient/rate-limit/validation propagate so the caller degrades
		// AND warns ONCE (never silently swallow a real provider failure).
		if code := domain.ErrorCodeOf(err); code == domain.ErrCodeProviderInvalidReference || code == domain.ErrCodeProviderPayloadInvalid {
			return domain.BuyerFiscalInfo{FetchedAt: a.now()}, nil
		}
		return domain.BuyerFiscalInfo{}, mapPricingReaderError(err)
	}
	return mapBuyerFiscalInfo(billing, a.now()), nil
}

func mapBuyerFiscalInfo(resp mlBillingInfoResponse, fetchedAt time.Time) domain.BuyerFiscalInfo {
	info := domain.BuyerFiscalInfo{FetchedAt: fetchedAt}
	if resp.Buyer == nil || resp.Buyer.BillingInfo == nil {
		return info
	}
	bi := resp.Buyer.BillingInfo
	info.Name = trimmedStringPtr(joinName(bi.Name, bi.LastName))
	if bi.Identification != nil {
		info.DocType = trimmedStringPtr(bi.Identification.Type)
		info.DocNumber = trimmedStringPtr(bi.Identification.Number)
	}
	info.Address = mapBuyerFiscalAddress(bi.Address)
	return info
}

// mapBuyerFiscalAddress maps the ML billing address onto the neutral fiscal address. Every
// field degrades to nil when blank; when ALL fields are blank (e.g. an obfuscated/masked
// address) the whole address is nil — honest absence, never a fabricated-empty object (ADR-17).
func mapBuyerFiscalAddress(addr *mlBillingAddress) *domain.BuyerFiscalAddress {
	if addr == nil {
		return nil
	}
	out := domain.BuyerFiscalAddress{
		StreetName:   trimmedStringPtr(addr.StreetName),
		StreetNumber: trimmedStringPtr(addr.StreetNumber),
		City:         trimmedStringPtr(addr.CityName),
		ZipCode:      trimmedStringPtr(addr.ZipCode),
		CountryID:    trimmedStringPtr(addr.CountryID),
	}
	if addr.State != nil {
		out.StateCode = trimmedStringPtr(addr.State.Code)
		out.StateName = trimmedStringPtr(addr.State.Name)
	}
	if out.StreetName == nil && out.StreetNumber == nil && out.City == nil &&
		out.StateCode == nil && out.StateName == nil && out.ZipCode == nil && out.CountryID == nil {
		return nil
	}
	return &out
}

// joinName joins a first + last name with a single separating space, trimming each part and
// the whole; an all-blank pair yields "" (which trimmedStringPtr turns into a nil Name).
func joinName(first, last string) string {
	return strings.TrimSpace(strings.TrimSpace(first) + " " + strings.TrimSpace(last))
}

// trimmedStringPtr returns a pointer to the trimmed string, or nil when it trims to empty —
// how a masked/blank provider value degrades to honest absence (ADR-17), never a blank string.
func trimmedStringPtr(s string) *string {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil
	}
	return &t
}
