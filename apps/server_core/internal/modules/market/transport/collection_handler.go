package transport

import (
	"errors"
	"net/http"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// handleCollect runs an on-demand, single-product, synchronous collection
// (D-F4-p: no batch POST) and serializes the B.1 envelope.
func (h Handler) handleCollect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	req, err := ParseCollectionRequest(r)
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	summary, err := h.collections.Collect(r.Context(), req.Codprod)
	if err != nil {
		writeCollectionError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newCollectionResponse(summary))
}

// writeCollectionError maps the two named collection sentinels (D-F4-q) to
// their wire status codes; any other error is an honest internal failure.
func writeCollectionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ports.ErrProductNotFound):
		writeMarketError(w, http.StatusNotFound, ports.ErrProductNotFound.Error(), "produto não encontrado", map[string]any{})
	case errors.Is(err, application.ErrCollectionInProgress):
		writeMarketError(w, http.StatusConflict, application.ErrCollectionInProgress.Error(), "coleta em andamento para este produto", map[string]any{})
	default:
		writeMarketError(w, http.StatusInternalServerError, "internal_error", "internal error", map[string]any{})
	}
}

// collectionResponse is the B.1 wire envelope. Field names decisões/
// contagens/causas are verbatim per the ratified amendment; JSON tags carry
// no omitempty so every key always serializes.
type collectionResponse struct {
	Status    string                   `json:"status"`
	Decisoes  []collectionDecisionWire `json:"decisões"`
	Contagens map[string]int           `json:"contagens"`
	Causas    []collectionCauseWire    `json:"causas"`
}

type collectionDecisionWire struct {
	Codprod             string                     `json:"codprod"`
	MatchStatus         domain.MatchStatus         `json:"match_status"`
	PriceEvidenceStatus domain.PriceEvidenceStatus `json:"price_evidence_status"`
	BlockingState       *domain.BlockingState      `json:"blocking_state"`
}

// collectionCauseWire adds codprod (same product for the whole single-product
// run) around the application cause. Detail carries the application's honest,
// MPC-native observability note (operation + provider status) when present, or
// null when the cause has no extra provider context — never a fabricated one.
type collectionCauseWire struct {
	Codprod string  `json:"codprod"`
	Reason  string  `json:"reason"`
	Detail  *string `json:"detail"`
}

func newCollectionResponse(summary application.CollectionSummary) collectionResponse {
	decisoes := make([]collectionDecisionWire, 0, len(summary.Decisoes))
	for _, d := range summary.Decisoes {
		decisoes = append(decisoes, collectionDecisionWire{
			Codprod:             d.Codprod,
			MatchStatus:         d.MatchStatus,
			PriceEvidenceStatus: d.PriceEvidence,
			BlockingState:       d.Blocking,
		})
	}
	causas := make([]collectionCauseWire, 0, len(summary.Causas))
	for _, c := range summary.Causas {
		causas = append(causas, collectionCauseWire{Codprod: summary.Codprod, Reason: string(c.Cause), Detail: c.Detail})
	}
	return collectionResponse{
		Status:    string(summary.Status),
		Decisoes:  decisoes,
		Contagens: summary.Contagens,
		Causas:    causas,
	}
}
