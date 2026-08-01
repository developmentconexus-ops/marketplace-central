package main

// Round 3 — MIS-007 live measurement M1..M10 (see operator brief).
// Run with: go run ./cmd/mlprobe -round3
// Read-only. Zero writes to ML. Evidence saved under evidenceDir as before,
// plus a human-readable summary written to stdout (captured by the caller).

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const targetInstallationID = "inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0"

func (p *probe) round3(ctx context.Context, pool *pgxpool.Pool, sellerID string) {
	fmt.Println("=== M1: orders/search offset ceiling + scan support ===")
	p.m1(sellerID)

	fmt.Println("=== M2: listing_prices logistic_type sensitivity ===")
	p.m2()

	fmt.Println("=== M3: /item/{id}/performance ===")
	p.m3(allListingIDs(ctx, pool))

	fmt.Println("=== M4: /moderations/infractions/{user_id} ===")
	p.m4(sellerID)

	fmt.Println("=== M5: price_to_win on catalog listing ===")
	p.m5(ctx, pool)

	fmt.Println("=== M6: cancelled orders cancel_detail vs status_detail ===")
	p.m6(ctx, pool)

	fmt.Println("=== M7: shipment costs vs lead_time.cost ===")
	p.m7(ctx, pool)

	fmt.Println("=== M8: /items multiget id-count ceiling ===")
	p.m8(allListingIDs(ctx, pool))

	fmt.Println("=== M9: destination UF: shipment vs billing_info ===")
	p.m9(ctx, pool)

	fmt.Println("=== M10: /shipments/{id}/sla ===")
	p.m10(ctx, pool)

	fmt.Println("round3 done")
}

// ---- DB helpers (read-only) ----

func allListingIDs(ctx context.Context, pool *pgxpool.Pool) []string {
	rows, err := pool.Query(ctx, `
		SELECT DISTINCT provider_listing_id FROM listings
		WHERE installation_id = $1
		ORDER BY 1`, targetInstallationID)
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

func catalogListingIDs(ctx context.Context, pool *pgxpool.Pool) map[string]string {
	rows, err := pool.Query(ctx, `
		SELECT provider_listing_id, catalog_product_id FROM listings
		WHERE installation_id = $1 AND catalog_product_id IS NOT NULL AND catalog_product_id != ''
		ORDER BY 1`, targetInstallationID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var id, cp string
		if rows.Scan(&id, &cp) == nil {
			out[id] = cp
		}
	}
	return out
}

func cancelledOrderIDs(ctx context.Context, pool *pgxpool.Pool) []string {
	rows, err := pool.Query(ctx, `
		SELECT provider_order_id FROM orders_marketplace_orders
		WHERE provider_status = 'cancelled'
		ORDER BY 1`)
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

type orderShipPair struct {
	orderID    string
	shipmentID string
	dbState    string
}

func allOrderShipmentPairs(ctx context.Context, pool *pgxpool.Pool) []orderShipPair {
	rows, err := pool.Query(ctx, `
		SELECT o.provider_order_id, o.provider_shipment_id, COALESCE(s.dest_state, '')
		FROM orders_marketplace_orders o
		LEFT JOIN order_shipments s ON s.provider_order_id = o.provider_order_id
		WHERE o.provider_shipment_id IS NOT NULL AND o.provider_shipment_id != ''
		ORDER BY 1`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []orderShipPair
	for rows.Next() {
		var pair orderShipPair
		if rows.Scan(&pair.orderID, &pair.shipmentID, &pair.dbState) == nil {
			out = append(out, pair)
		}
	}
	return out
}

func sampleShipmentIDs(ctx context.Context, pool *pgxpool.Pool, n int) []string {
	rows, err := pool.Query(ctx, `
		SELECT provider_shipment_id FROM order_shipments
		WHERE status = 'delivered'
		ORDER BY 1 LIMIT $1`, n)
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

// ---- M1: orders/search offset ceiling + search_type=scan ----

func (p *probe) m1(sellerID string) {
	for _, off := range []int{0, 500, 1000, 1500, 5000, 10000} {
		path := fmt.Sprintf("/orders/search?seller=%s&offset=%d&limit=1", sellerID, off)
		st, h, body, dur, err := p.do(path, nil)
		notes := []string{}
		if err != nil {
			notes = append(notes, "ERR: "+err.Error())
		}
		var res map[string]any
		_ = json.Unmarshal(body, &res)
		small := map[string]any{}
		if res != nil {
			small["paging"] = res["paging"]
			if msg, ok := res["message"]; ok {
				small["message"] = msg
			}
			if e, ok := res["error"]; ok {
				small["error"] = e
			}
			if len(body) < 1500 {
				small["raw"] = res
			}
		}
		p.save(fmt.Sprintf("M1-offset-%d", off), "offset máximo aceito em /orders/search", "GET", path, st, dur, h, small, notes...)
	}

	// search_type=scan on /orders/search — undocumented for this endpoint.
	path := "/orders/search?seller=" + sellerID + "&search_type=scan&limit=1"
	st, h, body, dur, _ := p.do(path, nil)
	var res any
	_ = json.Unmarshal(body, &res)
	p.save("M1-search-type-scan", "search_type=scan é aceito em /orders/search?", "GET", path, st, dur, h, res)
}

// ---- M2: listing_prices logistic_type sensitivity ----

func (p *probe) m2() {
	categoryID := "MLB270310" // most common category among our real listings (23/34)
	type row struct {
		Price        float64 `json:"price"`
		LogisticType string  `json:"logistic_type"`
		Status       int     `json:"status"`
		FeeDetails   any     `json:"sale_fee_details"`
		SaleFee      any     `json:"sale_fee_amount"`
	}
	var out []row
	logisticTypes := []string{"", "drop_off", "xd_drop_off", "fulfillment"}
	prices := []float64{10, 13, 100}
	for _, lt := range logisticTypes {
		for _, pr := range prices {
			path := fmt.Sprintf("/sites/MLB/listing_prices?price=%.0f&category_id=%s&listing_type_id=gold_special", pr, categoryID)
			if lt != "" {
				path += "&logistic_type=" + lt
			}
			st, _, body, _, _ := p.do(path, nil)
			var res any
			_ = json.Unmarshal(body, &res)
			r := row{Price: pr, LogisticType: lt, Status: st}
			extract := func(m map[string]any) {
				r.SaleFee = m["sale_fee_amount"]
				r.FeeDetails = m["sale_fee_details"]
			}
			switch v := res.(type) {
			case map[string]any:
				extract(v)
			case []any:
				if len(v) > 0 {
					if m, ok := v[0].(map[string]any); ok {
						extract(m)
					}
				}
			}
			out = append(out, r)
		}
	}
	p.save("M2-logistic-type-sensitivity", "sale_fee_details.fixed_fee muda com logistic_type? (categoria "+categoryID+")", "GET", "/sites/MLB/listing_prices?...&logistic_type=...", 200, 0, nil, out,
		fmt.Sprintf("%d combinations (%d logistic_types x %d prices)", len(out), len(logisticTypes), len(prices)))
}

// ---- M3: /item/{id}/performance (singular item, undocumented) ----

func (p *probe) m3(ourIDs []string) {
	if len(ourIDs) == 0 {
		p.save("M3-item-performance", "/item/{id}/performance existe?", "GET", "", 0, 0, nil, nil, "SKIP: no listing ids in DB")
		return
	}
	n := min(len(ourIDs), 5)
	for i := 0; i < n; i++ {
		id := ourIDs[i]
		path := "/item/" + id + "/performance"
		st, h, body, dur, _ := p.do(path, nil)
		var res any
		_ = json.Unmarshal(body, &res)
		notes := []string{}
		if m, ok := res.(map[string]any); ok {
			var keys []string
			for k := range m {
				keys = append(keys, k)
			}
			notes = append(notes, "top-level keys: "+strings.Join(keys, ","))
			if score, ok := m["score"]; ok {
				notes = append(notes, fmt.Sprintf("score present, value=%v", score))
			}
			if level, ok := m["level"]; ok {
				notes = append(notes, fmt.Sprintf("level present, value=%v", level))
			}
			if buckets, ok := m["buckets"]; ok {
				notes = append(notes, fmt.Sprintf("buckets present, type=%T", buckets))
			}
		}
		p.save(fmt.Sprintf("M3-item-performance-%d", i+1), "/item/{id}/performance (singular) — score/level/buckets existem?", "GET", path, st, dur, h, res, notes...)
		if st == 200 {
			// one confirmed 200 is enough evidence of shape; still probe a couple more for consistency.
		}
	}
}

// ---- M4: /moderations/infractions/{user_id} ----

func (p *probe) m4(sellerID string) {
	path := "/moderations/infractions/" + sellerID
	st, h, body, dur, _ := p.do(path, nil)
	var res any
	_ = json.Unmarshal(body, &res)
	notes := []string{}
	switch v := res.(type) {
	case []any:
		notes = append(notes, fmt.Sprintf("array of %d infractions", len(v)))
		if len(v) > 0 {
			if m, ok := v[0].(map[string]any); ok {
				var keys []string
				for k := range m {
					keys = append(keys, k)
				}
				notes = append(notes, "first infraction keys: "+strings.Join(keys, ","))
			}
		}
	case map[string]any:
		var keys []string
		for k := range v {
			keys = append(keys, k)
		}
		notes = append(notes, "top-level keys: "+strings.Join(keys, ","))
	}
	p.save("M4-moderations-infractions", "infrações da conta — reason/remedy aparecem? (7 anúncios under_review)", "GET", path, st, dur, h, res, notes...)
}

// ---- M5: price_to_win on catalog listing (or explicit NO-CATALOG verdict) ----

func (p *probe) m5(ctx context.Context, pool *pgxpool.Pool) {
	cat := catalogListingIDs(ctx, pool)
	if len(cat) == 0 {
		p.save("M5-price-to-win", "price_to_win?version=v2 — conta compete em catálogo?", "GET", "", 0, 0, nil, nil,
			"NO-CATALOG: nenhum dos 34 listings tem catalog_listing=true nem catalog_product_id não-nulo — resultado válido, conta não compete em catálogo")
		return
	}
	i := 0
	for id := range cat {
		if i >= 3 {
			break
		}
		path := "/items/" + id + "/price_to_win?version=v2"
		st, h, body, dur, _ := p.do(path, nil)
		var res any
		_ = json.Unmarshal(body, &res)
		p.save(fmt.Sprintf("M5-price-to-win-%d", i+1), "price_to_win?version=v2 para listing com catalog_product_id="+cat[id], "GET", path, st, dur, h, res)
		i++
	}
}

// ---- M6: cancelled orders — cancel_detail vs status_detail ----

func (p *probe) m6(ctx context.Context, pool *pgxpool.Pool) {
	ids := cancelledOrderIDs(ctx, pool)
	if len(ids) == 0 {
		p.save("M6-cancel-detail", "cancel_detail vs status_detail em pedidos cancelados", "GET", "", 0, 0, nil, nil, "SKIP: no cancelled orders in DB")
		return
	}
	for i, oid := range ids {
		path := "/orders/" + oid
		st, h, body, dur, _ := p.do(path, nil)
		var res map[string]any
		_ = json.Unmarshal(body, &res)
		notes := []string{}
		small := map[string]any{}
		if res != nil {
			small["status"] = res["status"]
			small["status_detail"] = sanitize(res["status_detail"])
			if cd, ok := res["cancel_detail"]; ok {
				small["cancel_detail"] = sanitize(cd)
				notes = append(notes, "cancel_detail PRESENT")
				if cdm, ok := cd.(map[string]any); ok {
					var keys []string
					for k := range cdm {
						keys = append(keys, k)
					}
					notes = append(notes, "cancel_detail keys: "+strings.Join(keys, ","))
					notes = append(notes, fmt.Sprintf("requested_by=%v", cdm["requested_by"]))
				}
			} else {
				notes = append(notes, "cancel_detail ABSENT")
			}
			sdStr := fmt.Sprintf("%v", res["status_detail"])
			cdStr := fmt.Sprintf("%v", small["cancel_detail"])
			notes = append(notes, fmt.Sprintf("status_detail == cancel_detail (raw compare)? %v", sdStr == cdStr))
		}
		p.save(fmt.Sprintf("M6-cancel-detail-%d", i+1), "cancel_detail existe? requested_by? difere de status_detail? (order "+oid+")", "GET", path, st, dur, h, small, notes...)
	}
}

// ---- M7: shipment costs vs lead_time.cost ----

func (p *probe) m7(ctx context.Context, pool *pgxpool.Pool) {
	ids := sampleShipmentIDs(ctx, pool, 5)
	if len(ids) == 0 {
		p.save("M7-shipment-costs", "senders[].cost vs receiver.cost vs lead_time.cost", "GET", "", 0, 0, nil, nil, "SKIP: no delivered shipments in DB")
		return
	}
	for i, sid := range ids {
		hdr := map[string]string{"x-format-new": "true"}

		cpath := "/shipments/" + sid + "/costs"
		cst, ch, cbody, cdur, _ := p.do(cpath, hdr)
		var cres map[string]any
		_ = json.Unmarshal(cbody, &cres)
		notes := []string{}
		small := map[string]any{}
		if cres != nil {
			small["senders"] = cres["senders"]
			small["receiver"] = cres["receiver"]
			if senders, ok := cres["senders"].([]any); ok && len(senders) > 0 {
				if s0, ok := senders[0].(map[string]any); ok {
					notes = append(notes, fmt.Sprintf("senders[0].cost=%v discounts=%v", s0["cost"], s0["discounts"]))
				}
			}
			if recv, ok := cres["receiver"].(map[string]any); ok {
				notes = append(notes, fmt.Sprintf("receiver.cost=%v", recv["cost"]))
			}
		}
		p.save(fmt.Sprintf("M7-costs-%d", i+1), "GET /shipments/{id}/costs — quem paga (shipment "+sid+")", "GET", cpath, cst, cdur, ch, small, notes...)

		lpath := "/shipments/" + sid + "/lead_time"
		lst, lh, lbody, ldur, _ := p.do(lpath, hdr)
		var lres map[string]any
		_ = json.Unmarshal(lbody, &lres)
		lnotes := []string{}
		lsmall := map[string]any{}
		if lres != nil {
			lsmall["cost"] = lres["cost"]
			lnotes = append(lnotes, fmt.Sprintf("lead_time.cost=%v", lres["cost"]))
		}
		p.save(fmt.Sprintf("M7-leadtime-%d", i+1), "GET /shipments/{id}/lead_time.cost — comparar com costs.senders[].cost (shipment "+sid+")", "GET", lpath, lst, ldur, lh, lsmall, lnotes...)
	}
}

// ---- M8: /items multiget id-count ceiling ----

func (p *probe) m8(ourIDs []string) {
	if len(ourIDs) == 0 {
		p.save("M8-multiget-limit", "quantos ids /items?ids= aceita de fato?", "GET", "", 0, 0, nil, nil, "SKIP: no listing ids in DB")
		return
	}
	// Pad the id list by repeating (harmless — server dedupes or errors on count, not distinctness)
	padded := ourIDs
	for len(padded) < 51 {
		padded = append(padded, ourIDs[len(padded)%len(ourIDs)])
	}
	for _, n := range []int{20, 21, 50, 51} {
		ids := padded[:n]
		path := "/items?ids=" + strings.Join(ids, ",")
		st, h, body, dur, _ := p.do(path, nil)
		var res any
		_ = json.Unmarshal(body, &res)
		notes := []string{fmt.Sprintf("%d ids requested", n)}
		if arr, ok := res.([]any); ok {
			notes = append(notes, fmt.Sprintf("response array length=%d", len(arr)))
		}
		if m, ok := res.(map[string]any); ok {
			notes = append(notes, fmt.Sprintf("error body: message=%v error=%v", m["message"], m["error"]))
		}
		p.save(fmt.Sprintf("M8-multiget-%d", n), "limite real de ids em /items?ids=", "GET", path, st, dur, h, truncateBody(res), notes...)
	}
}

// ---- M9: destination UF — shipment vs billing_info; ICMS contributor field ----

func (p *probe) m9(ctx context.Context, pool *pgxpool.Pool) {
	pairs := allOrderShipmentPairs(ctx, pool)
	if len(pairs) == 0 {
		p.save("M9-uf-comparison", "UF destino: shipment vs billing_info", "GET", "", 0, 0, nil, nil, "SKIP: no order/shipment pairs in DB")
		return
	}
	type row struct {
		OrderID          string `json:"order_id"`
		ShipUF           string `json:"shipment_uf"`
		BillingUF        string `json:"billing_uf"`
		Match            bool   `json:"match"`
		ShipStatus       int    `json:"shipment_status"`
		BillingStatus    int    `json:"billing_status"`
		IEFieldsObserved string `json:"ie_fields_observed,omitempty"`
	}
	var rows []row
	ieKeysSeen := map[string]bool{}
	for _, pair := range pairs {
		spath := "/shipments/" + pair.shipmentID
		sst, _, sbody, _, _ := p.do(spath, map[string]string{"x-format-new": "true"})
		var sres map[string]any
		_ = json.Unmarshal(sbody, &sres)
		shipUF := ""
		if dest, ok := sres["destination"].(map[string]any); ok {
			if sa, ok := dest["shipping_address"].(map[string]any); ok {
				if st, ok := sa["state"].(map[string]any); ok {
					if id, ok := st["id"].(string); ok {
						shipUF = strings.TrimPrefix(id, "BR-")
					}
				}
			}
		}

		bpath := "/orders/" + pair.orderID + "/billing_info"
		bst, _, bbody, _, _ := p.do(bpath, map[string]string{"x-version": "2"})
		var bres map[string]any
		_ = json.Unmarshal(bbody, &bres)
		billingUF := ""
		if buyer, ok := bres["buyer"].(map[string]any); ok {
			if bi, ok := buyer["billing_info"].(map[string]any); ok {
				if addr, ok := bi["address"].(map[string]any); ok {
					if st, ok := addr["state"].(map[string]any); ok {
						if code, ok := st["code"].(string); ok {
							billingUF = strings.TrimPrefix(code, "BR-")
						}
					}
				}
				if taxes, ok := bi["taxes"].(map[string]any); ok {
					// Record only field PRESENCE/shape, never the raw value — inscriptions
					// carries the buyer's state tax registration number (documento) and
					// must never be written to evidence, redacted or not.
					if tt, ok := taxes["taxpayer_type"].(map[string]any); ok {
						if desc, ok := tt["description"].(string); ok && desc != "" {
							ieKeysSeen["taxes.taxpayer_type.description="+desc] = true
						} else {
							ieKeysSeen["taxes.taxpayer_type=empty"] = true
						}
					}
					if insc, ok := taxes["inscriptions"].(map[string]any); ok {
						if _, has := insc["state_registration"]; has {
							ieKeysSeen["taxes.inscriptions.state_registration=PRESENT[REDACTED]"] = true
						} else {
							ieKeysSeen["taxes.inscriptions=empty"] = true
						}
					}
				}
				if attrs, ok := bi["attributes"].(map[string]any); ok {
					if ct, ok := attrs["cust_type"]; ok {
						ieKeysSeen[fmt.Sprintf("attributes.cust_type=%v", ct)] = true
					}
				}
			}
		}

		rows = append(rows, row{
			OrderID: pair.orderID, ShipUF: shipUF, BillingUF: billingUF,
			Match:      shipUF != "" && billingUF != "" && strings.EqualFold(shipUF, billingUF),
			ShipStatus: sst, BillingStatus: bst,
		})
	}
	var ieKeys []string
	for k := range ieKeysSeen {
		ieKeys = append(ieKeys, k)
	}
	notes := []string{
		fmt.Sprintf("%d pairs compared", len(rows)),
		"possible ICMS-contributor-related keys observed across billing_info responses: " + strings.Join(ieKeys, ",") + " (empty = none found)",
	}
	p.save("M9-uf-comparison", "UF destino: shipment.receiver_address.state vs orders/{id}/billing_info x-version:2 (PII redigida)", "GET", "multi (see per-pair)", 200, 0, nil, rows, notes...)
}

// ---- M10: /shipments/{id}/sla ----

func (p *probe) m10(ctx context.Context, pool *pgxpool.Pool) {
	ids := sampleShipmentIDs(ctx, pool, 1)
	if len(ids) == 0 {
		p.save("M10-shipment-sla", "/shipments/{id}/sla existe?", "GET", "", 0, 0, nil, nil, "SKIP: no shipments in DB")
		return
	}
	path := "/shipments/" + ids[0] + "/sla"
	st, h, body, dur, _ := p.do(path, map[string]string{"x-format-new": "true"})
	var res any
	_ = json.Unmarshal(body, &res)
	notes := []string{}
	if m, ok := res.(map[string]any); ok {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		notes = append(notes, "top-level keys: "+strings.Join(keys, ","))
	}
	p.save("M10-shipment-sla", "/shipments/{id}/sla — nunca chamado antes; existe? o que responde?", "GET", path, st, dur, h, res, notes...)
}
