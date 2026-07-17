//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
	marketapp "marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	markettransport "marketplace-central/apps/server_core/internal/modules/market/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

type marketContractHarness struct {
	pool          *pgxpool.Pool
	repoA         *marketpostgres.Repository
	repoB         *marketpostgres.Repository
	handler       http.Handler
	tenantA       string
	tenantB       string
	installationA string
	now           time.Time
}

func TestMarketContractVirginDatabase(t *testing.T) {
	h := newMarketContractHarness(t)
	listingIDs := []string{"listing-virgin-b-" + h.tenantA, "listing-virgin-a-" + h.tenantA}
	status, body, raw := marketGET(t, h, "/market/observations?installation_id="+url.QueryEscape(h.installationA)+"&listing_ids="+url.QueryEscape(strings.Join(listingIDs, ",")))
	if status != http.StatusOK {
		t.Fatalf("observations status=%d body=%s", status, raw)
	}
	observationItems := marketItems(t, body)
	if got := marketItemIDs(observationItems, "listing_id"); fmt.Sprint(got) != fmt.Sprint(listingIDs) {
		t.Fatalf("observation ids=%v, want input order %v", got, listingIDs)
	}
	for _, item := range observationItems {
		assertNoPriceObservation(t, item, raw)
	}
	if body["as_of"] != h.now.Format(time.RFC3339) {
		t.Fatalf("observations as_of=%v, want %s", body["as_of"], h.now.Format(time.RFC3339))
	}

	productIDs := []string{"product-virgin-b-" + h.tenantA, "product-virgin-a-" + h.tenantA}
	status, body, raw = marketGET(t, h, "/market/references?product_ids="+url.QueryEscape(strings.Join(productIDs, ",")))
	if status != http.StatusOK {
		t.Fatalf("references status=%d body=%s", status, raw)
	}
	referenceItems := marketItems(t, body)
	if got := marketItemIDs(referenceItems, "product_id"); fmt.Sprint(got) != fmt.Sprint(productIDs) {
		t.Fatalf("reference ids=%v, want input order %v", got, productIDs)
	}
	for _, item := range referenceItems {
		assertNoPriceReference(t, item, raw)
	}
}

func TestMarketContractRoundTripAndTenantIsolation(t *testing.T) {
	h := newMarketContractHarness(t)
	seedMarketContractRows(t, h)

	listingIDs := []string{"listing-missing-" + h.tenantA, "listing-round-trip-" + h.tenantA, "listing-tenant-b-only-" + h.tenantA}
	status, body, raw := marketGET(t, h, "/market/observations?installation_id="+url.QueryEscape(h.installationA)+"&listing_ids="+url.QueryEscape(strings.Join(listingIDs, ",")))
	if status != http.StatusOK {
		t.Fatalf("observations status=%d body=%s", status, raw)
	}
	items := marketItems(t, body)
	if got := marketItemIDs(items, "listing_id"); fmt.Sprint(got) != fmt.Sprint(listingIDs) {
		t.Fatalf("observation ids=%v, want input order %v", got, listingIDs)
	}
	assertNoPriceObservation(t, items[0], raw)
	assertRoundTripObservation(t, items[1])
	assertNoPriceObservation(t, items[2], raw)
	if strings.Contains(raw, "999.99") || strings.Contains(raw, "TENANT_B_SECRET") {
		t.Fatalf("tenant B observation leaked: %s", raw)
	}

	productIDs := []string{"product-tenant-b-only-" + h.tenantA, "product-round-trip-" + h.tenantA, "product-missing-" + h.tenantA}
	status, body, raw = marketGET(t, h, "/market/references?product_ids="+url.QueryEscape(strings.Join(productIDs, ",")))
	if status != http.StatusOK {
		t.Fatalf("references status=%d body=%s", status, raw)
	}
	items = marketItems(t, body)
	if got := marketItemIDs(items, "product_id"); fmt.Sprint(got) != fmt.Sprint(productIDs) {
		t.Fatalf("reference ids=%v, want input order %v", got, productIDs)
	}
	assertNoPriceReference(t, items[0], raw)
	assertRoundTripReference(t, items[1])
	assertNoPriceReference(t, items[2], raw)
	if strings.Contains(raw, "TENANT_B_SECRET") || strings.Contains(raw, "tenant-b-method") {
		t.Fatalf("tenant B reference leaked: %s", raw)
	}

	_, err := h.pool.Exec(context.Background(), `
		INSERT INTO market_observations (tenant_id, listing_id, evidence_state, captured_at, source)
		VALUES ($1, $2, 'observed', $3, $4)
	`, h.tenantA, "listing-null-source-"+h.tenantA, h.now, nil)
	if err == nil {
		t.Fatal("INSERT with null source succeeded, want database rejection")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23502" {
		t.Fatalf("INSERT error=%v, want SQLSTATE 23502", err)
	}
}

func TestMarketContractTwoHundredOneIDs(t *testing.T) {
	h := newMarketContractHarness(t)
	ids := make([]string, 201)
	for i := range ids {
		ids[i] = fmt.Sprintf("listing-%03d-%s", i, h.tenantA)
	}

	status, body, raw := marketGET(t, h, "/market/observations?installation_id="+url.QueryEscape(h.installationA)+"&listing_ids="+url.QueryEscape(strings.Join(ids, ",")))
	if status != http.StatusUnprocessableEntity || marketErrorCode(body) != "too_many_ids" {
		t.Fatalf("observations status=%d code=%q body=%s, want 422 too_many_ids", status, marketErrorCode(body), raw)
	}

	for i := range ids {
		ids[i] = fmt.Sprintf("product-%03d-%s", i, h.tenantA)
	}
	status, body, raw = marketGET(t, h, "/market/references?product_ids="+url.QueryEscape(strings.Join(ids, ",")))
	if status != http.StatusUnprocessableEntity || marketErrorCode(body) != "too_many_ids" {
		t.Fatalf("references status=%d code=%q body=%s, want 422 too_many_ids", status, marketErrorCode(body), raw)
	}
}

func TestMarketContractMalformedListingID(t *testing.T) {
	h := newMarketContractHarness(t)
	listingIDs := []string{"inst_" + h.tenantA + "~MLB123456~-", "garbage listing id"}
	status, body, raw := marketGET(t, h, "/market/observations?installation_id="+url.QueryEscape(h.installationA)+"&listing_ids="+url.QueryEscape(strings.Join(listingIDs, ",")))
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s, want malformed ID to remain item-level", status, raw)
	}
	items := marketItems(t, body)
	if len(items) != 2 || items[1]["listing_id"] != listingIDs[1] || items[1]["evidence_state"] != string(domain.EvidenceStateNoPriceEvidence) {
		t.Fatalf("malformed item=%v, want item-level no_price_evidence", items)
	}
}

func newMarketContractHarness(t *testing.T) *marketContractHarness {
	t.Helper()
	pool, _ := testpostgres.OpenPool(t, "market_contract")
	token := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	h := &marketContractHarness{
		pool:          pool,
		tenantA:       "market-contract-a-" + token,
		tenantB:       "market-contract-b-" + token,
		installationA: "market-installation-a-" + token,
		now:           time.Date(2026, time.July, 16, 18, 0, 0, 0, time.UTC),
	}
	h.repoA = marketpostgres.NewRepository(pool, h.tenantA)
	h.repoB = marketpostgres.NewRepository(pool, h.tenantB)
	mux := httpx.NewRouteClassMux()
	readService := marketapp.NewReadService(h.repoA, h.repoA, func() time.Time { return h.now })
	markettransport.NewHandler(readService).Register(mux)
	h.handler = mux
	return h
}

func seedMarketContractRows(t *testing.T, h *marketContractHarness) {
	t.Helper()
	older := time.Date(2026, time.July, 16, 15, 0, 0, 0, time.UTC)
	newer := time.Date(2026, time.July, 16, 17, 0, 0, 0, time.UTC)
	listingID := "listing-round-trip-" + h.tenantA
	if err := h.repoA.AppendObservations(context.Background(), []domain.MarketObservation{
		marketObservation(listingID, older, "1.01", domain.CompetitiveStatusWinning, domain.SourceOfficialAPI),
		marketObservationWithSignals(listingID, newer, domain.SourceManual),
	}); err != nil {
		t.Fatalf("append tenant A observations: %v", err)
	}
	if err := h.repoB.AppendObservations(context.Background(), []domain.MarketObservation{
		marketObservation("listing-round-trip-"+h.tenantA, h.now, "999.99", domain.CompetitiveStatusNotListed, domain.SourceVendor),
		marketObservation("listing-tenant-b-only-"+h.tenantA, h.now, "888.88", domain.CompetitiveStatusNotListed, domain.SourceVendor),
	}); err != nil {
		t.Fatalf("append tenant B observations: %v", err)
	}

	if err := h.repoA.AppendReferences(context.Background(), []domain.MarketReference{
		marketReference("product-round-trip-"+h.tenantA, older, domain.MatchStateReview, domain.SourceOfficialAPI),
		marketReferenceWithSignals("product-round-trip-"+h.tenantA, newer, domain.SourceManual),
	}); err != nil {
		t.Fatalf("append tenant A references: %v", err)
	}
	if err := h.repoB.AppendReferences(context.Background(), []domain.MarketReference{
		marketReferenceWithSecret("product-round-trip-"+h.tenantA, h.now, domain.SourceVendor),
		marketReferenceWithSecret("product-tenant-b-only-"+h.tenantA, h.now, domain.SourceVendor),
	}); err != nil {
		t.Fatalf("append tenant B references: %v", err)
	}
}

func marketObservation(listingID string, capturedAt time.Time, amount string, status domain.CompetitiveStatus, source domain.Source) domain.MarketObservation {
	return domain.MarketObservation{
		ListingID:         listingID,
		OurSalePrice:      &domain.Money{Amount: amount, Currency: "BRL"},
		CompetitiveStatus: &status,
		EvidenceState:     domain.EvidenceStateObserved,
		CapturedAt:        &capturedAt,
		Source:            &source,
	}
}

func marketObservationWithSignals(listingID string, capturedAt time.Time, source domain.Source) domain.MarketObservation {
	status := domain.CompetitiveStatusCompeting
	return domain.MarketObservation{
		ListingID:         listingID,
		OurSalePrice:      &domain.Money{Amount: "101.01", Currency: "BRL"},
		CompetitiveStatus: &status,
		WinnerPrice:       &domain.Money{Amount: "202.02", Currency: "BRL"},
		CompetitiveTarget: &domain.Money{Amount: "303.03", Currency: "BRL"},
		CatalogOfferPrice: &domain.Money{Amount: "404.04", Currency: "BRL"},
		CatalogStats:      &domain.CatalogStats{Min: "1.01", P25: "2.02", Median: "3.03", TrimmedMean: "4.04", P75: "5.05", Max: "6.06", NOffers: 7, NSellers: 6},
		EvidenceState:     domain.EvidenceStateObserved,
		CapturedAt:        &capturedAt,
		Source:            &source,
	}
}

func marketReference(productID string, capturedAt time.Time, matchState domain.MatchState, source domain.Source) domain.MarketReference {
	catalogProductID := "catalog-" + productID
	matchMethod := "semantic-gate"
	stats := marketStats()
	return domain.MarketReference{
		ProductID: productID, CatalogProductID: &catalogProductID, MatchState: matchState, MatchMethod: &matchMethod,
		CatalogStats: &stats, EvidenceState: domain.EvidenceStateObserved, CapturedAt: &capturedAt, Source: &source,
	}
}

func marketReferenceWithSignals(productID string, capturedAt time.Time, source domain.Source) domain.MarketReference {
	return marketReference(productID, capturedAt, domain.MatchStateAccept, source)
}

func marketReferenceWithSecret(productID string, capturedAt time.Time, source domain.Source) domain.MarketReference {
	catalogProductID := "TENANT_B_SECRET"
	matchMethod := "tenant-b-method"
	return domain.MarketReference{
		ProductID: productID, CatalogProductID: &catalogProductID, MatchState: domain.MatchStateReject, MatchMethod: &matchMethod,
		EvidenceState: domain.EvidenceStateObserved, CapturedAt: &capturedAt, Source: &source,
	}
}

func marketStats() domain.CatalogStats {
	return domain.CatalogStats{Min: "1.01", P25: "2.02", Median: "3.03", TrimmedMean: "4.04", P75: "5.05", Max: "6.06", NOffers: 7, NSellers: 6}
}

func assertRoundTripObservation(t *testing.T, item map[string]any) {
	t.Helper()
	if item["evidence_state"] != string(domain.EvidenceStateObserved) || item["competitive_status"] != string(domain.CompetitiveStatusCompeting) || item["source"] != string(domain.SourceManual) {
		t.Fatalf("observation state/status/source=%v/%v/%v", item["evidence_state"], item["competitive_status"], item["source"])
	}
	for field, amount := range map[string]string{"our_sale_price": "101.01", "winner_price": "202.02", "competitive_target": "303.03", "catalog_offer_price": "404.04"} {
		value, ok := item[field].(map[string]any)
		if !ok || value["amount"] != amount || value["currency"] != "BRL" {
			t.Fatalf("observation %s=%v, want %s BRL", field, item[field], amount)
		}
	}
	assertStats(t, item["catalog_stats"])
	if item["captured_at"] != "2026-07-16T17:00:00Z" {
		t.Fatalf("captured_at=%v", item["captured_at"])
	}
}

func assertRoundTripReference(t *testing.T, item map[string]any) {
	t.Helper()
	if item["evidence_state"] != string(domain.EvidenceStateObserved) || item["match_state"] != string(domain.MatchStateAccept) || item["match_method"] != "semantic-gate" || item["source"] != string(domain.SourceManual) {
		t.Fatalf("reference state/match/method/source=%v/%v/%v/%v", item["evidence_state"], item["match_state"], item["match_method"], item["source"])
	}
	if item["catalog_product_id"] != "catalog-"+item["product_id"].(string) {
		t.Fatalf("catalog_product_id=%v", item["catalog_product_id"])
	}
	assertStats(t, item["catalog_stats"])
	if item["captured_at"] != "2026-07-16T17:00:00Z" {
		t.Fatalf("captured_at=%v", item["captured_at"])
	}
}

func assertStats(t *testing.T, value any) {
	t.Helper()
	stats, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("catalog_stats=%v, want object", value)
	}
	for field, want := range map[string]string{"min": "1.01", "p25": "2.02", "median": "3.03", "trimmed_mean": "4.04", "p75": "5.05", "max": "6.06"} {
		if stats[field] != want {
			t.Fatalf("catalog_stats[%s]=%v, want %s", field, stats[field], want)
		}
	}
	if stats["n_offers"] != float64(7) || stats["n_sellers"] != float64(6) {
		t.Fatalf("catalog_stats counts=%v/%v", stats["n_offers"], stats["n_sellers"])
	}
}

func assertNoPriceObservation(t *testing.T, item map[string]any, raw string) {
	t.Helper()
	if item["evidence_state"] != string(domain.EvidenceStateNoPriceEvidence) {
		t.Fatalf("observation evidence_state=%v", item["evidence_state"])
	}
	for _, field := range []string{"our_sale_price", "competitive_status", "winner_price", "competitive_target", "catalog_offer_price", "catalog_stats", "captured_at", "source"} {
		value, present := item[field]
		if !present || value != nil {
			t.Fatalf("observation field %s=%v present=%v, want explicit null", field, value, present)
		}
		if !strings.Contains(raw, `"`+field+`":null`) {
			t.Fatalf("raw body lacks explicit null %s: %s", field, raw)
		}
	}
}

func assertNoPriceReference(t *testing.T, item map[string]any, raw string) {
	t.Helper()
	if item["evidence_state"] != string(domain.EvidenceStateNoPriceEvidence) || item["match_state"] != string(domain.MatchStateNoCandidate) {
		t.Fatalf("reference state=%v match_state=%v", item["evidence_state"], item["match_state"])
	}
	for _, field := range []string{"catalog_product_id", "match_method", "catalog_stats", "captured_at", "source"} {
		value, present := item[field]
		if !present || value != nil {
			t.Fatalf("reference field %s=%v present=%v, want explicit null", field, value, present)
		}
		if !strings.Contains(raw, `"`+field+`":null`) {
			t.Fatalf("raw body lacks explicit null %s: %s", field, raw)
		}
	}
}

func marketGET(t *testing.T, h *marketContractHarness, path string) (int, map[string]any, string) {
	t.Helper()
	recorder := httptest.NewRecorder()
	h.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	raw := recorder.Body.String()
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s: status=%d body=%q: %v", path, recorder.Code, raw, err)
	}
	return recorder.Code, body, raw
}

func marketItems(t *testing.T, body map[string]any) []map[string]any {
	t.Helper()
	items, ok := body["items"].([]any)
	if !ok {
		t.Fatalf("body items=%v", body["items"])
	}
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		decoded, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("item=%v", item)
		}
		result = append(result, decoded)
	}
	return result
}

func marketItemIDs(items []map[string]any, key string) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item[key].(string))
	}
	return ids
}

func marketErrorCode(body map[string]any) string {
	errorBody, ok := body["error"].(map[string]any)
	if !ok {
		return ""
	}
	code, _ := errorBody["code"].(string)
	return code
}

type marketSourceFile struct {
	path    string
	content string
}

func TestMarketContractAbsenceProof(t *testing.T) {
	repoRoot := marketRepoRoot(t)
	files := marketProductionGoFiles(t, repoRoot)

	collectorViolations := marketCollectorViolations(files)
	if len(collectorViolations) != 0 {
		t.Fatalf("production CollectorPort implementation found: %v", collectorViolations)
	}
	if len(marketCollectorViolations([]marketSourceFile{{path: "internal/modules/market/fake.go", content: "func (c fake) CollectObservations(ctx context.Context, installationID string, listingIDs []string) ([]domain.MarketObservation, error) { return nil, nil }"}})) == 0 {
		t.Fatal("collector absence scan did not reject a temporary production method violation")
	}

	insertHits := marketInsertHits(files)
	wantInsertFiles := []string{
		"/internal/modules/market/adapters/postgres/observation_repository.go",
		"/internal/modules/market/adapters/postgres/reference_repository.go",
	}
	gotInsertFiles := make([]string, 0, len(insertHits))
	for _, hit := range insertHits {
		gotInsertFiles = append(gotInsertFiles, hit.path)
	}
	sort.Strings(gotInsertFiles)
	sort.Strings(wantInsertFiles)
	if fmt.Sprint(gotInsertFiles) != fmt.Sprint(wantInsertFiles) {
		t.Fatalf("market INSERT hits=%v, want only adapter files %v", gotInsertFiles, wantInsertFiles)
	}
	if len(marketInsertHits([]marketSourceFile{{path: "internal/modules/market/seed.go", content: `INSERT INTO market_observations`}})) == 0 {
		t.Fatal("market INSERT scan did not reject a temporary seed violation")
	}

	forbiddenHits := marketForbiddenHits(files)
	if len(forbiddenHits) != 0 {
		t.Fatalf("forbidden scraping/ML-market dependency found: %v", forbiddenHits)
	}
	if len(marketForbiddenHits([]marketSourceFile{{path: "internal/modules/market/fake.go", content: `import "github.com/PuerkitoBio/goquery"`}})) == 0 {
		t.Fatal("forbidden dependency scan did not reject a temporary goquery violation")
	}

	rootPath := filepath.Join(repoRoot, "apps", "server_core", "internal", "composition", "root.go")
	rootBytes, err := os.ReadFile(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	if violations := marketRootWiringViolations(string(rootBytes)); len(violations) != 0 {
		t.Fatalf("market collector/scheduler wiring found in root.go: %v", violations)
	}
	if len(marketRootWiringViolations("marketapp.NewCollectionService(collector)\nmarketScheduler.Start()")) == 0 {
		t.Fatal("composition absence scan did not reject temporary collector/scheduler wiring")
	}
}

func marketRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		candidate := filepath.Join(dir, "apps", "server_core", "internal", "modules", "market")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("repository root not found from %s", dir)
		}
		dir = parent
	}
}

func marketProductionGoFiles(t *testing.T, repoRoot string) []marketSourceFile {
	t.Helper()
	var files []marketSourceFile
	for _, root := range []string{
		filepath.Join(repoRoot, "apps", "server_core", "internal", "modules", "market"),
		filepath.Join(repoRoot, "apps", "server_core", "internal", "composition"),
	} {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(info.Name(), "_test.go") {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, marketSourceFile{path: filepath.ToSlash(path), content: string(content)})
			return nil
		}); err != nil {
			t.Fatalf("scan %s: %v", root, err)
		}
	}
	return files
}

func marketCollectorViolations(files []marketSourceFile) []string {
	method := regexp.MustCompile(`(?m)^\s*func\s*\([^)]*\)\s+Collect(?:Observations|References)\s*\(`)
	var violations []string
	for _, file := range files {
		if strings.HasSuffix(file.path, "/internal/modules/market/ports/collector.go") || strings.HasSuffix(file.path, "/internal/modules/market/application/collection_service.go") {
			continue
		}
		if method.MatchString(file.content) {
			violations = append(violations, file.path)
		}
	}
	return violations
}

func marketInsertHits(files []marketSourceFile) []marketSourceFile {
	insert := regexp.MustCompile(`(?i)INSERT\s+INTO\s+market_`)
	var hits []marketSourceFile
	for _, file := range files {
		if insert.MatchString(file.content) {
			hits = append(hits, file)
		}
	}
	result := make([]marketSourceFile, 0, len(hits))
	for _, hit := range hits {
		result = append(result, marketSourceFile{path: marketAdapterSuffix(hit.path), content: hit.content})
	}
	return result
}

func marketAdapterSuffix(path string) string {
	marker := "/internal/modules/market/"
	if index := strings.Index(path, marker); index >= 0 {
		return path[index:]
	}
	return path
}

func marketForbiddenHits(files []marketSourceFile) []string {
	forbidden := []string{"api.mercadolibre.com/sites", "search?", "colly", "goquery", "chromedp"}
	var hits []string
	for _, file := range files {
		lower := strings.ToLower(file.content)
		for _, value := range forbidden {
			if strings.Contains(lower, value) {
				hits = append(hits, file.path+":"+value)
			}
		}
	}
	return hits
}

func marketRootWiringViolations(source string) []string {
	var violations []string
	if strings.Contains(source, "NewCollectionService") {
		violations = append(violations, "NewCollectionService")
	}
	for _, line := range strings.Split(source, "\n") {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "marketscheduler") || strings.Contains(lower, "marketcollector") || strings.Contains(lower, "marketbackground") {
			violations = append(violations, strings.TrimSpace(line))
		}
	}
	return violations
}
