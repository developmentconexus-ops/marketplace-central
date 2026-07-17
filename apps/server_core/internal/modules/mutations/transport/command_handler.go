package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type commandService interface {
	Create(context.Context, application.CreateInput) (ports.Protocol, error)
	Preview(context.Context, string) (application.PreviewResult, error)
}

type approvalService interface {
	Approve(context.Context, string, json.RawMessage) (ports.Protocol, error)
	Cancel(context.Context, string) (ports.Protocol, error)
}

type retryService interface {
	Retry(context.Context, string) (ports.Protocol, error)
}

type Handler struct {
	commands commandService
	approval approvalService
	retry    retryService
}

func NewHandler(commands commandService, approval approvalService, retry retryService) Handler {
	return Handler{commands: commands, approval: approval, retry: retry}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("POST /mutations", h.HandleCreate)
	mux.HandleFunc("POST /mutations/{id}/preview", h.HandlePreview)
	mux.HandleFunc("POST /mutations/{id}/approve", h.HandleApprove)
	mux.HandleFunc("POST /mutations/{id}/cancel", h.HandleCancel)
	mux.HandleFunc("POST /mutations/{id}/retry", h.HandleRetry)
}

type createRequest struct {
	InstallationID string              `json:"installation_id"`
	Type           domain.ProtocolType `json:"type"`
	Selection      json.RawMessage     `json:"selection"`
	Intent         json.RawMessage     `json:"intent"`
	Actor          string              `json:"actor"`
}

func (h Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMutationMethodError(w)
		return
	}
	var req createRequest
	raw, err := readBody(r)
	if err != nil || decodeStrict(raw, &req) != nil {
		writeMutationError(w, errInvalidBody)
		return
	}
	if !isEnabledType(req.Type) {
		writeMutationError(w, errTypeNotEnabled)
		return
	}
	if strings.TrimSpace(req.InstallationID) == "" {
		writeMutationError(w, errInstallationRequired)
		return
	}
	if !isJSONObject(req.Intent) {
		writeMutationError(w, errInvalidIntent)
		return
	}
	if !isJSONObject(req.Selection) {
		writeMutationError(w, errInvalidBody)
		return
	}
	protocol, err := h.commands.Create(r.Context(), application.CreateInput{
		InstallationID: req.InstallationID,
		Type:           req.Type,
		Actor:          req.Actor,
		Intent:         req.Intent,
		Selection:      req.Selection,
	})
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, protocolJSONFrom(protocol))
}

func (h Handler) HandlePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMutationMethodError(w)
		return
	}
	if err := validateOptionalBody(r); err != nil {
		writeMutationError(w, err)
		return
	}
	result, err := h.commands.Preview(r.Context(), r.PathValue("id"))
	if err != nil {
		writeMutationError(w, err)
		return
	}
	response := protocolJSONFrom(result.Protocol)
	items := make([]map[string]any, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, mutationItemJSONFrom(item))
	}
	response["items"] = items
	httpx.WriteJSON(w, http.StatusOK, response)
}

func (h Handler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMutationMethodError(w)
		return
	}
	raw, err := readBody(r)
	if err != nil {
		writeMutationError(w, errInvalidBody)
		return
	}
	var req struct {
		Execute *bool `json:"execute"`
	}
	if decodeStrict(raw, &req) != nil {
		writeMutationError(w, errInvalidBody)
		return
	}
	protocol, err := h.approval.Approve(r.Context(), r.PathValue("id"), raw)
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, protocolJSONFrom(protocol))
}

func (h Handler) HandleCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMutationMethodError(w)
		return
	}
	if err := validateOptionalBody(r); err != nil {
		writeMutationError(w, err)
		return
	}
	protocol, err := h.approval.Cancel(r.Context(), r.PathValue("id"))
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, protocolJSONFrom(protocol))
}

func (h Handler) HandleRetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMutationMethodError(w)
		return
	}
	if err := validateOptionalBody(r); err != nil {
		writeMutationError(w, err)
		return
	}
	protocol, err := h.retry.Retry(r.Context(), r.PathValue("id"))
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, protocolJSONFrom(protocol))
}

func readBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func decodeStrict(raw []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("trailing JSON content")
		}
		return err
	}
	return nil
}

func validateOptionalBody(r *http.Request) error {
	raw, err := readBody(r)
	if err != nil {
		return errInvalidBody
	}
	if strings.TrimSpace(string(raw)) == "" {
		return nil
	}
	var empty struct{}
	if err := decodeStrict(raw, &empty); err != nil {
		return errInvalidBody
	}
	return nil
}

func isJSONObject(raw json.RawMessage) bool {
	if len(raw) == 0 || !json.Valid(raw) || string(bytes.TrimSpace(raw)) == "null" {
		return false
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return false
	}
	return object != nil
}

func isEnabledType(protocolType domain.ProtocolType) bool {
	return domain.ProtocolTypeEnabled(protocolType)
}

func protocolJSONFrom(protocol ports.Protocol) map[string]any {
	return map[string]any{
		"protocol_id": protocol.ProtocolID, "installation_id": protocol.InstallationID,
		"type": protocol.Type, "state": protocol.State, "actor": protocol.Actor,
		"intent": protocol.Intent, "selection": protocol.Selection, "totals": protocol.Totals,
		"source_as_of": protocol.SourceAsOf, "retried_from": protocol.RetriedFrom,
		"created_at": protocol.CreatedAt, "previewed_at": protocol.PreviewedAt,
		"approved_at": protocol.ApprovedAt, "finished_at": protocol.FinishedAt,
	}
}

func mutationItemJSONFrom(item ports.MutationItem) map[string]any {
	return map[string]any{
		"seq": item.Seq, "item_id": item.ItemID, "listing_id": item.ListingID,
		"idempotency_key": item.IdempotencyKey, "before": item.Before, "after": item.After,
		"state": item.State, "failure": item.Failure, "applied_at": item.AppliedAt,
	}
}
