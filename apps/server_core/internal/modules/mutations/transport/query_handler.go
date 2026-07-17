package transport

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

const (
	mutationQueryDefaultLimit = 50
	mutationQueryMaxLimit     = 200
)

type mutationReadService interface {
	ListProtocols(context.Context, application.ProtocolQuery) (application.ProtocolPage, error)
	GetProtocol(context.Context, string) (ports.Protocol, error)
	ListItems(context.Context, application.ItemQuery) (application.ItemPage, error)
}

type QueryHandler struct{ read mutationReadService }

func NewQueryHandler(read mutationReadService) QueryHandler { return QueryHandler{read: read} }

func (h QueryHandler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("GET /mutations", h.HandleList)
	mux.HandleFunc("GET /mutations/{id}", h.HandleGet)
	mux.HandleFunc("GET /mutations/{id}/items", h.HandleItems)
}

type mutationPageEnvelope struct {
	Items      []map[string]any `json:"items"`
	NextCursor *string          `json:"next_cursor"`
	PageSize   int              `json:"page_size"`
}

func (h QueryHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMutationQueryMethodError(w)
		return
	}
	query, err := parseMutationProtocolQuery(r.URL.Query())
	if err != nil {
		writeMutationError(w, err)
		return
	}
	page, err := h.read.ListProtocols(r.Context(), query)
	if err != nil {
		writeMutationError(w, err)
		return
	}
	items := make([]map[string]any, 0, len(page.Items))
	for _, protocol := range page.Items {
		items = append(items, protocolJSONFrom(protocol))
	}
	next, err := encodeMutationProtocolNextCursor(page.NextCursor)
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, mutationPageEnvelope{Items: items, NextCursor: next, PageSize: len(items)})
}

func (h QueryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMutationQueryMethodError(w)
		return
	}
	protocol, err := h.read.GetProtocol(r.Context(), r.PathValue("id"))
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, protocolJSONFrom(protocol))
}

func (h QueryHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMutationQueryMethodError(w)
		return
	}
	query, err := parseMutationItemQuery(r.URL.Query(), r.PathValue("id"))
	if err != nil {
		writeMutationError(w, err)
		return
	}
	page, err := h.read.ListItems(r.Context(), query)
	if err != nil {
		writeMutationError(w, err)
		return
	}
	items := make([]map[string]any, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mutationItemJSONFrom(item))
	}
	next, err := encodeMutationItemNextCursor(page.NextCursor)
	if err != nil {
		writeMutationError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, mutationPageEnvelope{Items: items, NextCursor: next, PageSize: len(items)})
}

func parseMutationProtocolQuery(values url.Values) (application.ProtocolQuery, error) {
	if err := rejectUnknownMutationQueryKeys(values, map[string]struct{}{"installation_id": {}, "state": {}, "type": {}, "cursor": {}, "limit": {}}); err != nil {
		return application.ProtocolQuery{}, err
	}
	query := application.ProtocolQuery{Limit: mutationQueryDefaultLimit}
	installation, present, err := mutationQueryScalar(values, "installation_id")
	if err != nil {
		return application.ProtocolQuery{}, err
	}
	if !present || strings.TrimSpace(installation) == "" {
		return application.ProtocolQuery{}, errInstallationRequired
	}
	query.InstallationID = installation
	state, present, err := mutationQueryScalar(values, "state")
	if err != nil {
		return application.ProtocolQuery{}, err
	}
	if present {
		query.State, err = parseProtocolState(state)
		if err != nil {
			return application.ProtocolQuery{}, err
		}
	}
	protocolType, present, err := mutationQueryScalar(values, "type")
	if err != nil {
		return application.ProtocolQuery{}, err
	}
	if present {
		query.Type = domain.ProtocolType(protocolType)
		if !isEnabledType(query.Type) {
			return application.ProtocolQuery{}, invalidMutationQuery("type")
		}
	}
	query.Cursor, err = parseMutationProtocolCursor(values)
	if err != nil {
		return application.ProtocolQuery{}, err
	}
	query.Limit, err = parseMutationLimit(values)
	return query, err
}

func parseMutationItemQuery(values url.Values, protocolID string) (application.ItemQuery, error) {
	if err := rejectUnknownMutationQueryKeys(values, map[string]struct{}{"state": {}, "cursor": {}, "limit": {}}); err != nil {
		return application.ItemQuery{}, err
	}
	query := application.ItemQuery{ProtocolID: protocolID, Limit: mutationQueryDefaultLimit}
	state, present, err := mutationQueryScalar(values, "state")
	if err != nil {
		return application.ItemQuery{}, err
	}
	if present {
		query.State, err = parseItemState(state)
		if err != nil {
			return application.ItemQuery{}, err
		}
	}
	query.Cursor, err = parseMutationItemCursor(values)
	if err != nil {
		return application.ItemQuery{}, err
	}
	query.Limit, err = parseMutationLimit(values)
	return query, err
}

func mutationQueryScalar(values url.Values, key string) (string, bool, error) {
	items, present := values[key]
	if !present {
		return "", false, nil
	}
	if len(items) != 1 || strings.TrimSpace(items[0]) != items[0] || items[0] == "" {
		return "", false, invalidMutationQuery(key)
	}
	return items[0], true, nil
}

func rejectUnknownMutationQueryKeys(values url.Values, allowed map[string]struct{}) error {
	for key := range values {
		if _, ok := allowed[key]; !ok {
			return invalidMutationQuery(key)
		}
	}
	return nil
}

func parseMutationLimit(values url.Values) (int, error) {
	raw, present, err := mutationQueryScalar(values, "limit")
	if err != nil || !present {
		return mutationQueryDefaultLimit, err
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > mutationQueryMaxLimit {
		return 0, invalidMutationQuery("limit")
	}
	return limit, nil
}

func parseMutationProtocolCursor(values url.Values) (application.ProtocolCursor, error) {
	raw, present, err := mutationQueryScalar(values, "cursor")
	if err != nil || !present {
		return application.ProtocolCursor{}, err
	}
	var envelope mutationProtocolCursorEnvelope
	if err := decodeMutationCursor(raw, &envelope); err != nil || envelope.Version != 1 || envelope.Kind != "mutation_protocol" || envelope.CreatedAt.IsZero() || envelope.ProtocolID == "" {
		return application.ProtocolCursor{}, invalidMutationQuery("cursor")
	}
	return application.ProtocolCursor{CreatedAt: envelope.CreatedAt, ProtocolID: envelope.ProtocolID}, nil
}

func parseMutationItemCursor(values url.Values) (application.ItemCursor, error) {
	raw, present, err := mutationQueryScalar(values, "cursor")
	if err != nil || !present {
		return application.ItemCursor{}, err
	}
	var envelope mutationItemCursorEnvelope
	if err := decodeMutationCursor(raw, &envelope); err != nil || envelope.Version != 1 || envelope.Kind != "mutation_item" || envelope.Seq < 1 {
		return application.ItemCursor{}, invalidMutationQuery("cursor")
	}
	return application.ItemCursor{Seq: envelope.Seq}, nil
}

func parseProtocolState(value string) (domain.ProtocolState, error) {
	state := domain.ProtocolState(value)
	switch state {
	case domain.ProtocolStateDraft, domain.ProtocolStatePreviewed, domain.ProtocolStateApproved,
		domain.ProtocolStateApplying, domain.ProtocolStateApplied, domain.ProtocolStatePartiallyFailed,
		domain.ProtocolStateFailedPreserved, domain.ProtocolStateCancelled:
		return state, nil
	default:
		return "", invalidMutationQuery("state")
	}
}

func parseItemState(value string) (domain.ItemState, error) {
	state := domain.ItemState(value)
	switch state {
	case domain.ItemStatePreviewed, domain.ItemStateApproved, domain.ItemStateApplying,
		domain.ItemStateApplied, domain.ItemStateFailed, domain.ItemStateSkipped:
		return state, nil
	default:
		return "", invalidMutationQuery("state")
	}
}

type mutationProtocolCursorEnvelope struct {
	Version    int       `json:"v"`
	Kind       string    `json:"kind"`
	CreatedAt  time.Time `json:"created_at"`
	ProtocolID string    `json:"protocol_id"`
}

type mutationItemCursorEnvelope struct {
	Version int    `json:"v"`
	Kind    string `json:"kind"`
	Seq     int    `json:"seq"`
}

func encodeMutationProtocolNextCursor(cursor *application.ProtocolCursor) (*string, error) {
	if cursor == nil {
		return nil, nil
	}
	raw, err := encodeMutationCursor(mutationProtocolCursorEnvelope{Version: 1, Kind: "mutation_protocol", CreatedAt: cursor.CreatedAt, ProtocolID: cursor.ProtocolID})
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func encodeMutationItemNextCursor(cursor *application.ItemCursor) (*string, error) {
	if cursor == nil {
		return nil, nil
	}
	raw, err := encodeMutationCursor(mutationItemCursorEnvelope{Version: 1, Kind: "mutation_item", Seq: cursor.Seq})
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

func encodeMutationCursor(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(payload), nil
}

func decodeMutationCursor(encoded string, target any) error {
	payload, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil || base64.StdEncoding.EncodeToString(payload) != encoded {
		return errors.New("invalid mutation cursor")
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("invalid mutation cursor")
	}
	return nil
}

func invalidMutationQuery(key string) error {
	if key == "cursor" {
		return &requestError{code: "invalid_cursor", message: "cursor inválido"}
	}
	return &requestError{code: "invalid_filter", message: "consulta inválida: parâmetro " + key}
}

func writeMutationQueryMethodError(w http.ResponseWriter) {
	w.Header().Set("Allow", http.MethodGet)
	writeMutationError(w, errMethodNotAllowed)
}
