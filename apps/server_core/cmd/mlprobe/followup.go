package main

// Follow-up probes (round 2) — run with: go run ./cmd/mlprobe -round2
// F1: listings count by status (34 vs 2006 discrepancy)
// F2: third-party item access variants (403 root cause — attributes? auth? endpoint?)
// F3: single GET /shipments/{id} with x-format-new (handling deadline location)
// F4: tariff fixed-fee threshold refine (between 10 and 19)

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (p *probe) round2(sellerID string, compItemID string, shipmentID string, categoryID string) {
	// F1: paging.total per status via regular (non-scan) search
	for _, status := range []string{"", "active", "paused", "closed"} {
		path := "/users/" + sellerID + "/items/search?limit=1"
		label := "all"
		if status != "" {
			path += "&status=" + status
			label = status
		}
		st, _, body, dur, _ := p.do(path, nil)
		var res map[string]any
		_ = json.Unmarshal(body, &res)
		p.save("F1-listing-count-"+label, "total de anúncios por status (34 vs 2006?)", "GET", path, st, dur, nil,
			map[string]any{"paging": res["paging"]})
	}

	// F2: third-party item access variants
	if compItemID != "" {
		variants := []struct {
			name string
			path string
			auth bool
		}{
			{"F2a-third-party-multiget-noattrs", "/items?ids=" + compItemID, true},
			{"F2b-third-party-single-auth", "/items/" + compItemID, true},
			{"F2c-third-party-single-noauth", "/items/" + compItemID, false},
			{"F2d-third-party-attrs-safe", "/items?ids=" + compItemID + "&attributes=id,price,available_quantity,listing_type_id,shipping", true},
			{"F2e-third-party-attrs-sold", "/items?ids=" + compItemID + "&attributes=id,sold_quantity", true},
		}
		for _, v := range variants {
			var st int
			var body []byte
			var dur int64
			if v.auth {
				st, _, body, dur, _ = p.do(v.path, nil)
			} else {
				st, body, dur = p.doNoAuth(v.path)
			}
			parsed := parseBody(body)
			// summarize: status per envelope + sold_quantity presence
			notes := []string{}
			if arr, ok := parsed.([]any); ok && len(arr) > 0 {
				if env, ok := arr[0].(map[string]any); ok {
					notes = append(notes, fmt.Sprintf("envelope code=%v", env["code"]))
					if b, ok := env["body"].(map[string]any); ok {
						_, hasSold := b["sold_quantity"]
						notes = append(notes, fmt.Sprintf("sold_quantity present=%v value=%v", hasSold, b["sold_quantity"]))
					}
				}
			}
			if m, ok := parsed.(map[string]any); ok {
				_, hasSold := m["sold_quantity"]
				notes = append(notes, fmt.Sprintf("sold_quantity present=%v value=%v", hasSold, m["sold_quantity"]))
			}
			p.save(v.name, "acesso a item de TERCEIRO: qual variante funciona?", "GET", v.path, st, dur, nil, truncateBody(parsed), notes...)
		}
	}

	// F3: single shipment with x-format-new — where is the handling deadline?
	if shipmentID != "" {
		st, h, body, dur, _ := p.do("/shipments/"+shipmentID, map[string]string{"x-format-new": "true"})
		parsed := sanitize(parseBody(body))
		notes := []string{}
		if m, ok := parsed.(map[string]any); ok {
			var keys []string
			for k := range m {
				keys = append(keys, k)
			}
			notes = append(notes, "top-level keys: "+strings.Join(keys, ","))
		}
		p.save("F3-shipment-single-xformatnew", "GET /shipments/{id} x-format-new: campos de prazo de despacho", "GET", "/shipments/"+shipmentID, st, dur, h, parsed, notes...)
	}

	// F4: refine fixed-fee threshold between 10 and 19
	type point struct {
		Price   float64 `json:"price"`
		Status  int     `json:"status"`
		Details any     `json:"sale_fee_details"`
	}
	var sweep []point
	for _, pr := range []float64{11, 12, 12.5, 13, 14, 15, 16, 17, 18} {
		path := fmt.Sprintf("/sites/MLB/listing_prices?price=%.2f&category_id=%s&listing_type_id=gold_special", pr, categoryID)
		st, _, body, _, _ := p.do(path, nil)
		var res any
		_ = json.Unmarshal(body, &res)
		pt := point{Price: pr, Status: st}
		if m, ok := res.(map[string]any); ok {
			pt.Details = m["sale_fee_details"]
		} else if arr, ok := res.([]any); ok && len(arr) > 0 {
			if m, ok := arr[0].(map[string]any); ok {
				pt.Details = m["sale_fee_details"]
			}
		}
		sweep = append(sweep, pt)
	}
	p.save("F4-tariff-threshold-refine", "limiar exato da tarifa fixa (entre 10 e 19)", "GET", "/sites/MLB/listing_prices?price=11..18", 200, 0, nil, sweep)
}

func (p *probe) doNoAuth(path string) (int, []byte, int64) {
	req, err := newGetNoAuth(mlBase + path)
	if err != nil {
		return 0, nil, 0
	}
	resp, err := p.http.Do(req)
	if err != nil {
		return 0, nil, 0
	}
	defer resp.Body.Close()
	body := readCapped(resp)
	return resp.StatusCode, body, 0
}

func truncateBody(v any) any {
	raw, err := json.Marshal(v)
	if err != nil || len(raw) < 3000 {
		return v
	}
	return string(raw[:3000]) + "...[truncated]"
}
