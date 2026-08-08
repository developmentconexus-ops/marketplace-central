package main

// Round 4 — MIS-007 re-measurement of M3 and M5 after the operator reactivated
// listings (M1 round was captured while the account was in "férias" mode, so
// all 34 listings read as paused/under_review — a transitory state, not the
// API's real behavior). Re-run with: go run ./cmd/mlprobe -round4
// Also includes the bonus sub_status[] check requested alongside the remeasure.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (p *probe) round4(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("=== M3b/bonus: live status + sub_status check via /items?ids= ===")
	activeIDs, activeCatalog := p.liveStatusCheck(ctx, pool)

	fmt.Println("=== M3b: /item/{id}/performance on live-active non-catalog items ===")
	p.m3b(activeIDs, activeCatalog)

	fmt.Println("=== M5b: price_to_win on the 10 catalog-linked listings, now active ===")
	p.m5b(ctx, pool, activeCatalog)

	fmt.Println("round4 done")
}

// liveStatusCheck hits /items?ids= (live) for all DB-known listing ids, in
// batches of 20 (the confirmed M8 hard limit), and reports which are
// genuinely `active` right now plus their sub_status[]. Returns the live
// active ids split into (non-catalog, catalog-linked).
func (p *probe) liveStatusCheck(ctx context.Context, pool *pgxpool.Pool) (activeNonCatalog []string, activeCatalog map[string]string) {
	all := allListingIDs(ctx, pool)
	catalogMap := catalogListingIDs(ctx, pool)
	activeCatalog = map[string]string{}

	type liveRow struct {
		ID               string `json:"id"`
		Status           string `json:"status"`
		SubStatus        []any  `json:"sub_status"`
		CatalogProductID any    `json:"catalog_product_id"`
	}
	var allRows []liveRow

	for start := 0; start < len(all); start += 20 {
		end := start + 20
		if end > len(all) {
			end = len(all)
		}
		batch := all[start:end]
		path := "/items?ids=" + strings.Join(batch, ",") + "&attributes=id,status,sub_status,catalog_product_id"
		st, h, body, dur, _ := p.do(path, nil)
		var arr []any
		_ = json.Unmarshal(body, &arr)
		for _, e := range arr {
			env, ok := e.(map[string]any)
			if !ok {
				continue
			}
			b, ok := env["body"].(map[string]any)
			if !ok {
				continue
			}
			id, _ := b["id"].(string)
			status, _ := b["status"].(string)
			var sub []any
			if s, ok := b["sub_status"].([]any); ok {
				sub = s
			}
			allRows = append(allRows, liveRow{ID: id, Status: status, SubStatus: sub, CatalogProductID: b["catalog_product_id"]})
			if status == "active" {
				if cp, has := catalogMap[id]; has && cp != "" {
					activeCatalog[id] = cp
				} else {
					activeNonCatalog = append(activeNonCatalog, id)
				}
			}
		}
		p.save(fmt.Sprintf("M3bonus-live-status-batch-%d", start/20+1), "status/sub_status AO VIVO via /items?ids= (Postgres estava velho: conta em férias na rodada anterior)", "GET", path, st, dur, h, allRows[start:min(end, len(allRows))])
	}

	activeCount, nonEmptySubStatus := 0, 0
	var subStatusSample []string
	for _, r := range allRows {
		if r.Status == "active" {
			activeCount++
			if len(r.SubStatus) > 0 {
				nonEmptySubStatus++
				subStatusSample = append(subStatusSample, fmt.Sprintf("%s: %v", r.ID, r.SubStatus))
			}
		}
	}
	notes := []string{
		fmt.Sprintf("%d/%d listings are LIVE active now (was 0 in the previous round — conta estava em férias)", activeCount, len(allRows)),
		fmt.Sprintf("%d/%d active listings have non-empty sub_status[]", nonEmptySubStatus, activeCount),
	}
	if len(subStatusSample) > 0 {
		notes = append(notes, "non-empty sub_status samples: "+strings.Join(subStatusSample, " | "))
	} else {
		notes = append(notes, "sub_status[] is EMPTY on every active listing observed")
	}
	p.save("M3bonus-substatus-summary", "sub_status[] vem vazio quando o anúncio está active?", "GET", "/items?ids=... (agregado)", 200, 0, nil, allRows, notes...)

	return activeNonCatalog, activeCatalog
}

func (p *probe) m3b(activeNonCatalog []string, activeCatalog map[string]string) {
	if len(activeNonCatalog) == 0 {
		p.save("M3b-item-performance", "/item/{id}/performance em ativo não-catálogo", "GET", "", 0, 0, nil, nil,
			fmt.Sprintf("SKIP: 0 listings live-active AND non-catalog (active-catalog count=%d) — cannot prove success shape without at least one active/non-catalog item", len(activeCatalog)))
		return
	}
	n := min(len(activeNonCatalog), 5)
	for i := 0; i < n; i++ {
		id := activeNonCatalog[i]
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
			for _, want := range []string{"score", "level", "level_wording", "calculated_at", "buckets"} {
				v, has := m[want]
				notes = append(notes, fmt.Sprintf("%s present=%v value_type=%T", want, has, v))
			}
			if buckets, ok := m["buckets"].([]any); ok && len(buckets) > 0 {
				if b0, ok := buckets[0].(map[string]any); ok {
					var bkeys []string
					for k := range b0 {
						bkeys = append(bkeys, k)
					}
					notes = append(notes, "buckets[0] keys: "+strings.Join(bkeys, ","))
					if vars, ok := b0["variables"].([]any); ok && len(vars) > 0 {
						if v0, ok := vars[0].(map[string]any); ok {
							var vkeys []string
							for k := range v0 {
								vkeys = append(vkeys, k)
							}
							notes = append(notes, "buckets[0].variables[0] keys: "+strings.Join(vkeys, ","))
							if rules, ok := v0["rules"].([]any); ok && len(rules) > 0 {
								if r0, ok := rules[0].(map[string]any); ok {
									var rkeys []string
									for k := range r0 {
										rkeys = append(rkeys, k)
									}
									notes = append(notes, "buckets[0].variables[0].rules[0] keys: "+strings.Join(rkeys, ","))
									notes = append(notes, fmt.Sprintf("rules[0]: status=%v progress=%v mode=%v wordings=%v", r0["status"], r0["progress"], r0["mode"], r0["wordings"]))
								}
							}
						}
					}
				}
			}
		}
		p.save(fmt.Sprintf("M3b-item-performance-%d", i+1), "/item/{id}/performance em item AO VIVO active + não-catálogo — score/level/buckets reais", "GET", path, st, dur, h, res, notes...)
	}
}

func (p *probe) m5b(ctx context.Context, pool *pgxpool.Pool, activeCatalog map[string]string) {
	all := catalogListingIDs(ctx, pool) // full set of 10 (for reference/count in notes)
	if len(activeCatalog) == 0 {
		p.save("M5b-price-to-win", "price_to_win?version=v2 nos 10 catalog-linked, agora ativos?", "GET", "", 0, 0, nil, nil,
			fmt.Sprintf("SKIP: 0/%d catalog-linked listings are live-active right now — reactivation did not include the catalog-linked ones (or DB link is stale); result is itself informative", len(all)))
		return
	}
	i := 0
	for id, cp := range activeCatalog {
		i++
		path := "/items/" + id + "/price_to_win?version=v2"
		st, h, body, dur, _ := p.do(path, nil)
		var res map[string]any
		_ = json.Unmarshal(body, &res)
		notes := []string{fmt.Sprintf("catalog_product_id(DB)=%s, live active=true", cp)}
		if res != nil {
			notes = append(notes, fmt.Sprintf("status=%v reason=%v", res["status"], res["reason"]))
			if boosts, ok := res["boosts"].([]any); ok {
				notes = append(notes, fmt.Sprintf("boosts[] len=%d", len(boosts)))
				for bi, b := range boosts {
					if bm, ok := b.(map[string]any); ok {
						notes = append(notes, fmt.Sprintf("boosts[%d]: %v", bi, bm))
					}
				}
			} else {
				notes = append(notes, "boosts is null/absent")
			}
		}
		p.save(fmt.Sprintf("M5b-price-to-win-%d", i), "price_to_win?version=v2 em catalog-linked AO VIVO active (catalog_product_id="+cp+")", "GET", path, st, dur, h, res, notes...)
	}
}
