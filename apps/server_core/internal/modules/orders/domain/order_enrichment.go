package domain

import "strings"

// VinculoStatus reports whether an order has at least one item with a
// resolved product link. It is derived purely from item LinkQuality already
// present on the read model — no external resolver call.
type VinculoStatus string

const (
	VinculoStatusOK         VinculoStatus = "OK"
	VinculoStatusSemVinculo VinculoStatus = "SEM_VINCULO"
)

// MaskedBuyer is the only buyer identification allowed to leave this layer.
// City and UF stay nil until a later slice (e.g. shipment) can honestly
// populate them; Display never carries more than a first name plus the
// initial of a second token.
type MaskedBuyer struct {
	Display string
	City    *string
	UF      *string
}

// MaskBuyer derives a masked buyer identity from a nickname. It never
// fabricates data: an empty or blank nickname yields an empty MaskedBuyer.
func MaskBuyer(nickname string) MaskedBuyer {
	tokens := strings.Fields(nickname)
	if len(tokens) == 0 {
		return MaskedBuyer{}
	}

	display := tokens[0]
	if len(tokens) > 1 {
		if initial, ok := firstRuneUpper(tokens[1]); ok {
			display = tokens[0] + " " + initial + "."
		}
	}
	return MaskedBuyer{Display: display}
}

// DeriveVinculoStatus reports VinculoStatusOK when at least one item carries
// a resolved link quality, otherwise VinculoStatusSemVinculo.
func DeriveVinculoStatus(items []MarketplaceOrderItem) VinculoStatus {
	for _, item := range items {
		if item.LinkQuality == LinkQualityResolved {
			return VinculoStatusOK
		}
	}
	return VinculoStatusSemVinculo
}

func firstRuneUpper(token string) (string, bool) {
	runes := []rune(token)
	if len(runes) == 0 {
		return "", false
	}
	return strings.ToUpper(string(runes[0])), true
}
