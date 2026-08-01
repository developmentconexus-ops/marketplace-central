// Command mlprobe runs the read-only ML API live-test round T1–T12 defined in
// docs/design/ML-API-QUERY-CATALOG.md against the connected Mercado Livre account.
//
// It is a throwaway diagnostic tool: it resolves the active credential the same
// way the server does (Postgres + AES-GCM local key), performs only GET requests,
// and writes one evidence JSON file per test to docs/design/evidence/ml-api/.
// The access token is never written to evidence or stdout; buyer PII is masked.
//
// Run inside the backend container:
//
//	docker exec -w /workspace/apps/server_core marketplace-central-backend-1 go run ./cmd/mlprobe
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto"
)

const (
	mlBase      = "https://api.mercadolibre.com"
	evidenceDir = "/workspace/docs/design/evidence/ml-api"
	callDelay   = 200 * time.Millisecond
)

var piiKeys = map[string]bool{
	"name": true, "first_name": true, "last_name": true, "nickname_buyer": true,
	"email": true, "phone": true, "alternative_phone": true, "number": true,
	"doc_number": true, "receiver_name": true, "receiver_phone": true,
	"street_name": true, "street_number": true, "address_line": true,
	"comment": true, "zip_code": true, "business_name": true,
	"identification": true, "state_tax_id": true,
}

type evidence struct {
	Test       string      `json:"test"`
	Purpose    string      `json:"purpose"`
	Method     string      `json:"method"`
	URL        string      `json:"url"`
	ReqHeaders []string    `json:"request_headers"`
	Status     int         `json:"status"`
	DurationMS int64       `json:"duration_ms"`
	RespHeader http.Header `json:"response_headers_subset,omitempty"`
	Body       any         `json:"body"`
	Notes      []string    `json:"notes,omitempty"`
	FetchedAt  string      `json:"fetched_at"`
}

type probe struct {
	http  *http.Client
	token string
	notes []string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "mlprobe:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	dbURL := os.Getenv("MC_DATABASE_URL")
	key := os.Getenv("MPC_ENCRYPTION_KEY")
	if dbURL == "" || key == "" {
		return fmt.Errorf("MC_DATABASE_URL / MPC_ENCRYPTION_KEY missing")
	}
	if err := os.MkdirAll(evidenceDir, 0o755); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	token, installationID, tenantID, err := resolveToken(ctx, pool, key)
	if err != nil {
		return err
	}
	fmt.Printf("token resolved for installation=%s tenant=%s (token redacted)\n", installationID, tenantID)

	p := &probe{http: &http.Client{Timeout: 30 * time.Second}, token: token}

	if len(os.Args) > 1 && os.Args[1] == "-sla" {
		for _, path := range []string{"/shipments/47564931010/sla", "/shipments/47543615644/sla"} {
			st, h, body, dur, _ := p.do(path, map[string]string{"x-format-new": "true"})
			p.save("F5-shipment-sla-"+path[11:22], "prazo de despacho (handling deadline) via /sla", "GET", path, st, dur, h, sanitize(parseBody(body)))
		}
		return nil
	}

	if len(os.Args) > 1 && os.Args[1] == "-round2" {
		me, st, err := p.getJSON("/users/me", nil)
		if err != nil || st != 200 {
			return fmt.Errorf("token invalid? /users/me status=%d err=%v", st, err)
		}
		sellerID := fmt.Sprintf("%.0f", me["id"].(float64))
		// ids observed in round 1 evidence: competitor item (T8 403), shipment (T3/T6), category (T8a)
		p.round2(sellerID, "MLB4638031549", "47564931010", "MLB420116")
		fmt.Println("round2 done")
		return nil
	}

	if len(os.Args) > 1 && os.Args[1] == "-round3" {
		me, st, err := p.getJSON("/users/me", nil)
		if err != nil || st != 200 {
			return fmt.Errorf("token invalid? /users/me status=%d err=%v", st, err)
		}
		sellerID := fmt.Sprintf("%.0f", me["id"].(float64))
		p.round3(ctx, pool, sellerID)
		return nil
	}

	if len(os.Args) > 1 && os.Args[1] == "-round3-m9" {
		p.m9(ctx, pool)
		fmt.Println("m9 rerun done")
		return nil
	}

	if len(os.Args) > 1 && os.Args[1] == "-round4" {
		p.round4(ctx, pool)
		return nil
	}

	// T0 identity — also validates the token before everything else.
	me, st, err := p.getJSON("/users/me", nil)
	if err != nil || st != 200 {
		return fmt.Errorf("token invalid? /users/me status=%d err=%v — run credential refresh first", st, err)
	}
	sellerID := fmt.Sprintf("%.0f", me["id"].(float64))
	p.save("T0-identity", "identidade da conta conectada (valida token)", "GET", "/users/me", 200, 0, nil, sanitize(me), "seller_id="+sellerID)
	fmt.Println("T0 ok, seller", sellerID)

	ourItemIDs := listingIDs(ctx, pool, tenantID, installationID, 25)

	p.t1(ourItemIDs)
	orders := p.t2t3t4(sellerID)
	p.t5t6(orders)
	p.t7(sellerID)
	catID := p.t8t9(ctx, pool, tenantID, ourItemIDs)
	p.t10(catID)
	p.t11(sellerID)
	p.t12(ourItemIDs)

	fmt.Println("done — evidence in", evidenceDir)
	return nil
}

func resolveToken(ctx context.Context, pool *pgxpool.Pool, key string) (token, installationID, tenantID string, err error) {
	row := pool.QueryRow(ctx, `
		SELECT i.installation_id, i.tenant_id, c.encrypted_payload
		FROM integration_installations i
		JOIN integration_credentials c
		  ON c.tenant_id = i.tenant_id AND c.installation_id = i.installation_id
		 AND c.is_active = true AND c.revoked_at IS NULL
		WHERE i.provider_code = 'mercado_livre' AND i.status = 'connected'
		ORDER BY c.version DESC LIMIT 1`)
	var payload []byte
	if err = row.Scan(&installationID, &tenantID, &payload); err != nil {
		return "", "", "", fmt.Errorf("no connected mercado_livre credential: %w", err)
	}
	svc, err := crypto.NewLocalKeyService(key, "local-key-v1")
	if err != nil {
		return "", "", "", err
	}
	decoded, _, err := svc.DecryptJSON(payload)
	if err != nil {
		return "", "", "", fmt.Errorf("decrypt: %w", err)
	}
	tok, _ := decoded["access_token"].(string)
	if strings.TrimSpace(tok) == "" {
		return "", "", "", fmt.Errorf("access_token missing in payload")
	}
	return tok, installationID, tenantID, nil
}

func listingIDs(ctx context.Context, pool *pgxpool.Pool, tenantID, installationID string, n int) []string {
	rows, err := pool.Query(ctx, `
		SELECT DISTINCT provider_listing_id FROM listings
		WHERE tenant_id = $1 AND installation_id = $2 AND status = 'active'
		LIMIT $3`, tenantID, installationID, n)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// ---- HTTP helpers ----

func (p *probe) do(path string, headers map[string]string) (int, http.Header, []byte, int64, error) {
	time.Sleep(callDelay)
	req, err := http.NewRequest("GET", mlBase+path, nil)
	if err != nil {
		return 0, nil, nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	start := time.Now()
	resp, err := p.http.Do(req)
	dur := time.Since(start).Milliseconds()
	if err != nil {
		return 0, nil, nil, dur, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	return resp.StatusCode, resp.Header, body, dur, nil
}

func (p *probe) getJSON(path string, headers map[string]string) (map[string]any, int, error) {
	st, _, body, _, err := p.do(path, headers)
	if err != nil {
		return nil, st, err
	}
	var m map[string]any
	_ = json.Unmarshal(body, &m)
	return m, st, nil
}

func parseBody(body []byte) any {
	var v any
	if json.Unmarshal(body, &v) == nil {
		return v
	}
	s := string(body)
	if len(s) > 2000 {
		s = s[:2000] + "...[truncated]"
	}
	return s
}

func (p *probe) save(test, purpose, method, url string, status int, durMS int64, respHeaders http.Header, body any, notes ...string) {
	reqHeaders := []string{"Authorization: Bearer [REDACTED]"}
	ev := evidence{
		Test: test, Purpose: purpose, Method: method, URL: url,
		ReqHeaders: reqHeaders, Status: status, DurationMS: durMS,
		RespHeader: pickHeaders(respHeaders), Body: body, Notes: notes,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}
	raw, _ := json.MarshalIndent(ev, "", "  ")
	fn := filepath.Join(evidenceDir, test+".json")
	_ = os.WriteFile(fn, raw, 0o644)
	fmt.Printf("%s status=%d -> %s\n", test, status, fn)
}

func pickHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	out := http.Header{}
	for _, k := range []string{"X-Ratelimit-Remaining", "X-Ratelimit-Limit", "Retry-After", "Content-Type", "X-Request-Id"} {
		if v := h.Get(k); v != "" {
			out.Set(k, v)
		}
	}
	return out
}

// sanitize masks PII values recursively, preserving structure.
func sanitize(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			if piiKeys[strings.ToLower(k)] {
				out[k] = maskValue(val)
			} else {
				out[k] = sanitize(val)
			}
		}
		return out
	case []any:
		for i := range t {
			t[i] = sanitize(t[i])
		}
		return t
	default:
		return v
	}
}

func maskValue(v any) any {
	s, ok := v.(string)
	if !ok {
		if v == nil {
			return nil
		}
		return "[MASKED]"
	}
	if len(s) <= 2 {
		return "**"
	}
	return s[:2] + strings.Repeat("*", min(len(s)-2, 12))
}

// ---- tests ----

func (p *probe) t1(ids []string) {
	if len(ids) == 0 {
		p.save("T1-items-multiget", "multiget completo", "GET", "/items?ids=", 0, 0, nil, nil, "SKIP: no listing ids in DB")
		return
	}
	n := min(len(ids), 20)
	path := "/items?ids=" + strings.Join(ids[:n], ",")
	st, h, body, dur, err := p.do(path, nil)
	notes := []string{fmt.Sprintf("%d ids requested (max 20)", n)}
	if err != nil {
		notes = append(notes, "ERR: "+err.Error())
	}
	parsed := parseBody(body)
	// Field-presence summary for first successful item.
	if arr, ok := parsed.([]any); ok && len(arr) > 0 {
		if env, ok := arr[0].(map[string]any); ok {
			if b, ok := env["body"].(map[string]any); ok {
				var keys []string
				for k := range b {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				notes = append(notes, "first item body keys: "+strings.Join(keys, ","))
				for _, want := range []string{"sold_quantity", "date_created", "catalog_product_id", "catalog_listing", "category_id", "tags", "variations", "shipping", "health", "sub_status", "official_store_id", "permalink"} {
					if _, has := b[want]; !has {
						notes = append(notes, "MISSING: "+want)
					}
				}
			}
		}
	}
	p.save("T1-items-multiget", "multiget /items devolve variations/tags/catalog_*/sold_quantity?", "GET", path, st, dur, h, parsed, notes...)
}

// t2t3t4 covers order detail (sale_fee semantics), pack shape, billing info.
// Returns the raw orders list results for reuse by t5/t6.
func (p *probe) t2t3t4(sellerID string) []map[string]any {
	path := "/orders/search?seller=" + sellerID + "&sort=date_desc&limit=50"
	st, h, body, dur, err := p.do(path, nil)
	if err != nil || st != 200 {
		p.save("T2-order-sale-fee", "sale_fee unidade vs linha", "GET", path, st, dur, h, parseBody(body), "orders search failed")
		return nil
	}
	var res map[string]any
	_ = json.Unmarshal(body, &res)
	results, _ := res["results"].([]any)

	var orders []map[string]any
	var multiQty, anyOrder, packOrder map[string]any
	for _, r := range results {
		o, ok := r.(map[string]any)
		if !ok {
			continue
		}
		orders = append(orders, o)
		if anyOrder == nil {
			anyOrder = o
		}
		if packOrder == nil && o["pack_id"] != nil {
			packOrder = o
		}
		if multiQty == nil {
			if items, ok := o["order_items"].([]any); ok {
				for _, it := range items {
					if im, ok := it.(map[string]any); ok {
						if q, ok := im["quantity"].(float64); ok && q > 1 {
							multiQty = o
						}
					}
				}
			}
		}
	}

	// T2: full order detail for a qty>1 order (fallback: any order).
	target := multiQty
	notes := []string{}
	if target == nil {
		target = anyOrder
		notes = append(notes, "NO qty>1 order in last 50 — semantics still unproven for multi-qty; evidence from qty=1 order")
	}
	if target != nil {
		oid := fmt.Sprintf("%.0f", target["id"].(float64))
		opath := "/orders/" + oid
		ost, oh, obody, odur, _ := p.do(opath, nil)
		parsed := sanitize(parseBody(obody))
		if items, ok := target["order_items"].([]any); ok && len(items) > 0 {
			if im, ok := items[0].(map[string]any); ok {
				notes = append(notes, fmt.Sprintf("search-result item[0]: quantity=%v unit_price=%v full_unit_price=%v sale_fee=%v", im["quantity"], im["unit_price"], im["full_unit_price"], im["sale_fee"]))
			}
		}
		p.save("T2-order-sale-fee", "sale_fee unidade vs linha (comparar sale_fee x unit_price x quantity)", "GET", opath, ost, odur, oh, parsed, notes...)

		// discounts endpoint while here
		dst, dh, dbody, ddur, _ := p.do("/orders/"+oid+"/discounts", nil)
		p.save("T2b-order-discounts", "cupom/desconto endpoint separado", "GET", "/orders/"+oid+"/discounts", dst, ddur, dh, sanitize(parseBody(dbody)))
	}

	// T3: pack shape.
	if packOrder != nil {
		pid := fmt.Sprintf("%.0f", packOrder["pack_id"].(float64))
		ppath := "/packs/" + pid
		pst, ph, pbody, pdur, _ := p.do(ppath, nil)
		p.save("T3-pack", "shape de /packs/{id}", "GET", ppath, pst, pdur, ph, sanitize(parseBody(pbody)))
	} else {
		p.save("T3-pack", "shape de /packs/{id}", "GET", "/packs/{id}", 0, 0, nil, nil, "SKIP: no order with pack_id in last 50")
	}

	// T4: billing info (x-version: 2) for up to 2 orders.
	for i, o := range orders {
		if i >= 2 {
			break
		}
		oid := fmt.Sprintf("%.0f", o["id"].(float64))
		bpath := "/orders/" + oid + "/billing_info"
		bst, bh, bbody, bdur, _ := p.do(bpath, map[string]string{"x-version": "2"})
		p.save(fmt.Sprintf("T4-billing-info-%d", i+1), "billing_info x-version:2 shape (PF vs PJ, IE?)", "GET", bpath, bst, bdur, bh, sanitize(parseBody(bbody)), "request header x-version: 2")
	}
	return orders
}

func (p *probe) t5t6(orders []map[string]any) {
	var shipIDs []string
	for _, o := range orders {
		if sh, ok := o["shipping"].(map[string]any); ok {
			if id, ok := sh["id"].(float64); ok {
				shipIDs = append(shipIDs, fmt.Sprintf("%.0f", id))
			}
		}
	}
	if len(shipIDs) == 0 {
		p.save("T5-shipments-multiget", "multiget shipments existe? limite?", "GET", "/shipments?ids=", 0, 0, nil, nil, "SKIP: no shipment ids on orders")
		return
	}
	hdr := map[string]string{"x-format-new": "true"}

	// small multiget
	n := min(len(shipIDs), 2)
	path := "/shipments?ids=" + strings.Join(shipIDs[:n], ",")
	st, h, body, dur, _ := p.do(path, hdr)
	p.save("T5-shipments-multiget", "multiget shipments existe? shape com x-format-new", "GET", path, st, dur, h, sanitize(parseBody(body)), fmt.Sprintf("%d ids", n))

	// limit probe: 21 ids (docs conflict 20 vs 50)
	if len(shipIDs) >= 21 {
		path21 := "/shipments?ids=" + strings.Join(shipIDs[:21], ",")
		st21, h21, body21, dur21, _ := p.do(path21, hdr)
		p.save("T5b-shipments-multiget-21", "limite multiget: 21 ids passa?", "GET", path21, st21, dur21, h21, sanitize(parseBody(body21)))
	} else {
		p.save("T5b-shipments-multiget-21", "limite multiget: 21 ids passa?", "GET", "", 0, 0, nil, nil, fmt.Sprintf("SKIP: only %d shipment ids available", len(shipIDs)))
	}

	// T6: lead_time
	lpath := "/shipments/" + shipIDs[0] + "/lead_time"
	lst, lh, lbody, ldur, _ := p.do(lpath, hdr)
	p.save("T6-lead-time", "estimated_handling_limit presente (fonte SLA da Fila)", "GET", lpath, lst, ldur, lh, sanitize(parseBody(lbody)))

	// bonus: costs with X-Costs-New
	cpath := "/shipments/" + shipIDs[0] + "/costs"
	cst, ch, cbody, cdur, _ := p.do(cpath, map[string]string{"X-Costs-New": "true"})
	p.save("T6b-shipment-costs", "costs com X-Costs-New: quem paga frete", "GET", cpath, cst, cdur, ch, sanitize(parseBody(cbody)))
}

func (p *probe) t7(sellerID string) {
	from := time.Now().AddDate(0, 0, -30).UTC().Format("2006-01-02T15:04:05.000-00:00")
	path := "/orders/search?seller=" + sellerID + "&order.date_last_updated.from=" + from + "&sort=date_desc&limit=5"
	st, h, body, dur, _ := p.do(path, nil)
	var res map[string]any
	_ = json.Unmarshal(body, &res)
	notes := []string{"filter: date_last_updated.from last 30d + sort=date_desc"}
	if pg, ok := res["paging"].(map[string]any); ok {
		notes = append(notes, fmt.Sprintf("paging.total=%v", pg["total"]))
	}
	// keep evidence small: paging + first result dates only
	small := map[string]any{"paging": res["paging"]}
	if rs, ok := res["results"].([]any); ok && len(rs) > 0 {
		if o, ok := rs[0].(map[string]any); ok {
			small["first_result"] = map[string]any{"id": o["id"], "date_created": o["date_created"], "last_updated": o["last_updated"], "status": o["status"]}
		}
	}
	p.save("T7-orders-incremental-filter", "cursor incremental date_last_updated.from funciona?", "GET", path, st, dur, h, small, notes...)
}

// t8t9: competitor sold_quantity via catalog offers + public seller profile.
func (p *probe) t8t9(ctx context.Context, pool *pgxpool.Pool, tenantID string, ourIDs []string) (categoryID string) {
	// find a catalog_product_id from our items (hydrate a few)
	var catalogProductID string
	for _, id := range ourIDs {
		if catalogProductID != "" {
			break
		}
		m, st, _ := p.getJSON("/items/"+id+"?attributes=id,category_id,catalog_product_id,catalog_listing", nil)
		if st != 200 {
			continue
		}
		if categoryID == "" {
			categoryID, _ = m["category_id"].(string)
		}
		if cp, ok := m["catalog_product_id"].(string); ok && cp != "" {
			catalogProductID = cp
		}
	}
	if catalogProductID == "" {
		p.save("T8-competitor-sold-quantity", "sold_quantity de item de terceiro", "GET", "", 0, 0, nil, nil, "SKIP: none of first sampled items has catalog_product_id")
		return categoryID
	}

	opath := "/products/" + catalogProductID + "/items?limit=10"
	ost, oh, obody, odur, _ := p.do(opath, nil)
	p.save("T8a-catalog-offers", "ofertas do catálogo (concorrentes) — item_id/seller_id presentes?", "GET", opath, ost, odur, oh, parseBody(obody))

	var offers map[string]any
	_ = json.Unmarshal(obody, &offers)
	var compItemIDs []string
	var compSellerID string
	if rs, ok := offers["results"].([]any); ok {
		for _, r := range rs {
			if o, ok := r.(map[string]any); ok {
				if iid, ok := o["item_id"].(string); ok {
					compItemIDs = append(compItemIDs, iid)
				}
				if compSellerID == "" {
					if sid, ok := o["seller_id"].(float64); ok {
						compSellerID = fmt.Sprintf("%.0f", sid)
					}
				}
			}
		}
	}
	if len(compItemIDs) > 0 {
		n := min(len(compItemIDs), 10)
		mpath := "/items?ids=" + strings.Join(compItemIDs[:n], ",") + "&attributes=id,seller_id,price,sold_quantity,listing_type_id,shipping,official_store_id,available_quantity"
		mst, mh, mbody, mdur, _ := p.do(mpath, nil)
		p.save("T8-competitor-sold-quantity", "sold_quantity de itens de TERCEIROS via multiget", "GET", mpath, mst, mdur, mh, parseBody(mbody))
	} else {
		p.save("T8-competitor-sold-quantity", "sold_quantity de terceiros", "GET", "", 0, 0, nil, nil, "SKIP: catalog offers returned no item_id")
	}

	if compSellerID != "" {
		upath := "/users/" + compSellerID
		ust, uh, ubody, udur, _ := p.do(upath, nil)
		p.save("T9-public-seller-profile", "reputação/transactions de seller concorrente (público)", "GET", upath, ust, udur, uh, parseBody(ubody))
	} else {
		p.save("T9-public-seller-profile", "perfil público de concorrente", "GET", "", 0, 0, nil, nil, "SKIP: no competitor seller_id")
	}
	return categoryID
}

func (p *probe) t10(categoryID string) {
	if categoryID == "" {
		categoryID = "MLB271599" // fallback: generic category; note in evidence
	}
	type point struct {
		Price      float64 `json:"price"`
		Status     int     `json:"status"`
		SaleFee    any     `json:"sale_fee_amount"`
		FeeDetails any     `json:"sale_fee_details"`
		ListingFee any     `json:"listing_fee_amount"`
	}
	var sweep []point
	prices := []float64{10, 19, 29, 39, 49, 59, 69, 78, 79, 80, 99, 120, 150, 200}
	for _, pr := range prices {
		path := fmt.Sprintf("/sites/MLB/listing_prices?price=%.0f&category_id=%s&listing_type_id=gold_special", pr, categoryID)
		st, _, body, _, _ := p.do(path, nil)
		var res any
		_ = json.Unmarshal(body, &res)
		pt := point{Price: pr, Status: st}
		extract := func(m map[string]any) {
			pt.SaleFee = m["sale_fee_amount"]
			pt.FeeDetails = m["sale_fee_details"]
			pt.ListingFee = m["listing_fee_amount"]
		}
		switch r := res.(type) {
		case map[string]any:
			extract(r)
		case []any:
			if len(r) > 0 {
				if m, ok := r[0].(map[string]any); ok {
					extract(m)
				}
			}
		}
		sweep = append(sweep, pt)
	}
	p.save("T10-tariff-threshold-sweep", "achar limiar tarifa fixa 2026 por faixa de preço (categoria "+categoryID+", gold_special)", "GET", "/sites/MLB/listing_prices?price=...sweep", 200, 0, nil, sweep, fmt.Sprintf("%d price points", len(prices)))
}

func (p *probe) t11(sellerID string) {
	start := time.Now()
	var scrollID string
	total, batches := 0, 0
	var batchNotes []string
	for {
		path := "/users/" + sellerID + "/items/search?search_type=scan&limit=100"
		if scrollID != "" {
			path += "&scroll_id=" + scrollID
		}
		st, _, body, dur, err := p.do(path, nil)
		if err != nil || st != 200 {
			batchNotes = append(batchNotes, fmt.Sprintf("batch %d: status=%d err=%v", batches+1, st, err))
			break
		}
		var res map[string]any
		_ = json.Unmarshal(body, &res)
		results, _ := res["results"].([]any)
		if sid, ok := res["scroll_id"].(string); ok {
			scrollID = sid
		}
		batches++
		total += len(results)
		batchNotes = append(batchNotes, fmt.Sprintf("batch %d: %d ids in %dms", batches, len(results), dur))
		if len(results) == 0 || batches > 40 {
			break
		}
	}
	elapsed := time.Since(start)
	p.save("T11-scan-timing", "scan completo da conta: quantos ids, quanto tempo, scroll_id estável?", "GET", "/users/{me}/items/search?search_type=scan&limit=100", 200, elapsed.Milliseconds(), nil,
		map[string]any{"total_ids": total, "batches": batches, "elapsed_ms": elapsed.Milliseconds(), "batch_log": batchNotes})
}

func (p *probe) t12(ids []string) {
	if len(ids) == 0 {
		p.save("T12-429-behavior", "comportamento 429", "GET", "", 0, 0, nil, nil, "SKIP: no ids")
		return
	}
	// Controlled burst: 100 lightweight GETs, 10 workers, no artificial delay.
	// Not a stress test — stays far under documented per-minute limits.
	var wg sync.WaitGroup
	var mu sync.Mutex
	counts := map[int]int{}
	var h429 http.Header
	var body429 string
	client := &http.Client{Timeout: 15 * time.Second}
	sem := make(chan struct{}, 10)
	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			req, _ := http.NewRequest("GET", mlBase+"/items/"+ids[i%len(ids)]+"?attributes=id,price", nil)
			req.Header.Set("Authorization", "Bearer "+p.token)
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			resp.Body.Close()
			mu.Lock()
			counts[resp.StatusCode]++
			if resp.StatusCode == 429 && h429 == nil {
				h429 = resp.Header
				body429 = string(b)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	notes := []string{fmt.Sprintf("100 requests, 10 workers, %dms total", elapsed.Milliseconds())}
	if h429 == nil {
		notes = append(notes, "429 NOT triggered at this volume — limit is higher; Retry-After presence unproven")
	}
	p.save("T12-429-behavior", "429 real: header Retry-After? corpo?", "GET", "/items/{id}?attributes=id,price ×100", 200, elapsed.Milliseconds(), h429,
		map[string]any{"status_counts": counts, "first_429_body": body429})
}

func newGetNoAuth(url string) (*http.Request, error) {
	return http.NewRequest("GET", url, nil)
}

func readCapped(resp *http.Response) []byte {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	return body
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
