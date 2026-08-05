import { describe, expect, it, vi } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { createMarketplaceCentralClient } from "./index";
import type {
  CatalogProductFact,
  CatalogProductFactPage,
  CanonicalCatalogProduct,
  IntegrationProviderDefinition,
  ListingGroupPage,
  ListingPage,
  ListingReadModel,
  ListingDetail,
  ListingSummary,
  DashboardSummary,
  SyncRunPage,
  MarketObservation,
  MarketObservationPage,
  MarketReference,
  MarketReferencePage,
  CategoryAttribute,
  CategoryAttributesResponse,
  PricingCalcProfile,
  PricingDifalListResponse,
  PricingDifalOverrideResponse,
  PricingScenario,
  PricingDecomposeResponse,
  PricingSolveResponse,
} from "./index";

describe("sdk runtime", () => {
  it("requests an encoded ERP import chain path", async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        protocol: "#101-E",
        importados: 4,
        vinculados: 2,
        enfileirados: 3,
        queue_read_at: "2026-07-27T15:34:56Z",
      }),
    });
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });

    const chain = await client.getErpImportChain("id/with space");

    expect(fetchImpl).toHaveBeenCalledWith(
      "http://localhost:8080/erp/imports/id%2Fwith%20space/chain",
      expect.any(Object),
    );
    expect(chain.enfileirados).toBe(3);
    expect(chain.queue_read_at).toBe("2026-07-27T15:34:56Z");
  });

  it("models canonical CODPROD products with nullable source facts", () => {
    const unknown: CanonicalCatalogProduct = {
      internal_product_id: 1001,
      name: "Produto teste",
      ean: "7890000000000",
      manufacturer_reference: "REF-1001",
      seller_sku: "SELLER-1001",
      brand_name: null,
      product_group_name: null,
      ncm: null,
      quality_flags: ["invalid_ean"],
      cost_amount: { source: "sankhya", value: null, quality: "unknown", observed_at: null, quality_reason: "missing_cost" },
      price_amount: { source: "sankhya", value: 0, quality: "current", observed_at: "2026-07-13T12:00:00Z", quality_reason: null },
      stock_quantity: { source: "sankhya", value: 0, quality: "current", observed_at: "2026-07-13T12:00:00Z", quality_reason: null },
    };
    expect(unknown.internal_product_id).toBe(1001);
    expect(unknown.cost_amount.value).toBeNull();
    expect(unknown.price_amount.value).toBe(0);
  });

  it("keeps canonical nullable fields required across OpenAPI and SDK", () => {
    const openapi = readFileSync(resolve(process.cwd(), "../../contracts/api/marketplace-central.openapi.yaml"), "utf8");
    const sdk = readFileSync(resolve(process.cwd(), "src/index.ts"), "utf8");
    const factSchema = openapi.slice(openapi.indexOf("    CanonicalNumericSourceFact:"), openapi.indexOf("    CanonicalCatalogProduct:"));
    const productSchema = openapi.slice(openapi.indexOf("    CanonicalCatalogProduct:"), openapi.indexOf("    CatalogProduct:"));
    const catalogFactSchema = openapi.slice(openapi.indexOf("    CatalogProductFact:"), openapi.indexOf("    CatalogQuantityFact:"));
    const canonicalSdk = sdk.slice(sdk.indexOf("export interface CanonicalCatalogProduct"), sdk.indexOf("export interface CatalogProductFact"));
    const catalogFactSdk = sdk.slice(sdk.indexOf("export interface CatalogProductFact"), sdk.indexOf("export interface CatalogProductFactPage"));
    const catalogFact: CatalogProductFact = {
      internal_product_id: 1001,
      reference: null,
      description: null,
      ean: null,
      manufacturer_reference: "REF-1001",
      brand_name: null,
      ncm: null,
      quality_flags: ["invalid_ean", "ean_collision"],
      active: true,
      sellable_stock: { quantity: null, quality: [] },
      current_price: { amount: null, currency: "BRL", quality: [] },
      cost: { amount: null, currency: "BRL", quality: [] },
    };

    expect(factSchema).toContain("required: [source, value, quality, observed_at, quality_reason]");
    expect(factSchema).toMatch(/quality_reason:\s*\n\s*type: string\s*\n\s*nullable: true/);
    expect(productSchema).toContain("required: [internal_product_id, name, ean, manufacturer_reference, seller_sku, brand_name, product_group_name, cost_amount, price_amount, stock_quantity, ncm, quality_flags]");
    for (const field of ["ean", "manufacturer_reference", "seller_sku", "brand_name", "product_group_name"]) {
      expect(productSchema).toMatch(new RegExp(`${field}:\\s*\\n\\s*type: string\\s*\\n\\s*nullable: true`));
      expect(canonicalSdk).toMatch(new RegExp(`\\n  ${field}: string \\| null;`));
    }
    expect(productSchema).toMatch(/ncm:\s*\n\s*type: string\s*\n\s*nullable: true\s*\n\s*pattern: '\^\[0-9\]\{8\}\$'/);
    expect(productSchema).toMatch(/quality_flags:\s*\n\s*type: array\s*\n\s*items:\s*\n\s*\$ref: '#\/components\/schemas\/QualityFlag'/);
    expect(canonicalSdk).toContain("ncm: string | null;");
    expect(canonicalSdk).toContain("quality_flags: QualityFlag[];");
    expect(catalogFactSchema).toContain("required: [internal_product_id, reference, description, ean, active, sellable_stock, current_price, cost, manufacturer_reference, brand_name, ncm, quality_flags]");
    expect(catalogFactSchema).toMatch(/reference:\s*\n\s*type: string\s*\n\s*nullable: true\s*\n\s*deprecated: true\s*\n\s*description:.*manufacturer_reference/);
    for (const field of ["ean", "manufacturer_reference", "brand_name", "ncm"]) {
      expect(catalogFactSchema).toMatch(new RegExp(`${field}:\\s*\\n\\s*type: string\\s*\\n\\s*nullable: true`));
      expect(catalogFactSdk).toMatch(new RegExp(`\\n  ${field}: string \\| null;`));
    }
    expect(catalogFactSchema).toMatch(/quality_flags:\s*\n\s*type: array\s*\n\s*items:\s*\n\s*\$ref: '#\/components\/schemas\/QualityFlag'/);
    expect(catalogFactSdk).toContain("quality_flags: QualityFlag[];");
    expect(catalogFactSdk).toMatch(/@deprecated[\s\S]*\n\s*reference: string \| null;/);
    expect(catalogFact.quality_flags).toEqual(["invalid_ean", "ean_collision"]);
    expect(sdk).toMatch(/\n  quality_reason: string \| null;/);

    // quality_flags must admit every value the wire can carry — the reader passes
    // the closed QualityFlag domain through unfiltered (catalog_page.go emits the
    // `complete` default; the single-product reader propagates missing_product /
    // ambiguous_product). An enum narrowed to [invalid_ean, ean_collision] would
    // make conformant clients reject the majority of real responses.
    const qualityFlagSchema = openapi.slice(openapi.indexOf("    QualityFlag:"), openapi.indexOf("    CanonicalCatalogProduct:"));
    const qualityFlagSdk = sdk.slice(sdk.indexOf("export type QualityFlag"), sdk.indexOf("export interface CanonicalCatalogProduct"));
    for (const flag of ["complete", "missing_product", "ambiguous_product", "invalid_ean", "ean_collision"]) {
      expect(qualityFlagSchema).toContain(`- ${flag}`);
      expect(qualityFlagSdk).toContain(`"${flag}"`);
    }
    const wireFact: CatalogProductFact = { ...catalogFact, quality_flags: ["complete", "ean_collision"] };
    const wireCanonical: CanonicalCatalogProduct = {
      internal_product_id: 42,
      name: "PARAFUSO",
      ean: null,
      manufacturer_reference: null,
      seller_sku: null,
      brand_name: null,
      product_group_name: null,
      ncm: null,
      quality_flags: ["missing_product", "ambiguous_product"],
      cost_amount: { source: "oracle", value: null, quality: "unknown", observed_at: null, quality_reason: "cost unavailable" },
      price_amount: { source: "oracle", value: null, quality: "unknown", observed_at: null, quality_reason: "price unavailable" },
      stock_quantity: { source: "oracle", value: null, quality: "unknown", observed_at: null, quality_reason: "stock unavailable" },
    };
    expect(wireFact.quality_flags).toContain("complete");
    expect(wireCanonical.quality_flags).toContain("ambiguous_product");
  });

  it("lists catalog facts with the IC-01 page envelope", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const response = {
      items: [],
      next_cursor: null,
      page_size: 0,
      as_of: "2026-07-14T12:00:00Z",
    } satisfies CatalogProductFactPage;
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify(response), { status: 200 });
      },
    });

    const page = await client.listCatalogProductFacts({ cursor: "MTIz", limit: 25 });
    const searchPage = await client.searchCatalogProductFacts({ q: "PARAFUSO", limit: 50 });
    // Search pages: the cursor from a previous page of matches must reach the wire.
    await client.searchCatalogProductFacts({ q: "PARAFUSO", cursor: "NDU2", limit: 50 });

    expect(String(requests[0].input)).toBe("http://localhost:8080/catalog/products?cursor=MTIz&limit=25");
    expect(String(requests[1].input)).toBe("http://localhost:8080/catalog/products/search?q=PARAFUSO&limit=50");
    expect(String(requests[2].input)).toBe(
      "http://localhost:8080/catalog/products/search?q=PARAFUSO&cursor=NDU2&limit=50",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(page.next_cursor).toBeNull();
    expect(searchPage.page_size).toBe(0);
  });

  it("reads listing pages, grouped pages, details, and summaries with server filter syntax", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const page = { items: [], next_cursor: null, page_size: 0, as_of: "2026-07-15T12:00:00Z" } satisfies ListingPage;
    const listing = {
      listing_id: "install/1~provider/1~-",
      installation_id: "install/1",
      provider: "mercado_livre",
      provider_listing_id: "MLB1",
      variation_id: null,
      title: "Parafuso sextavado M6",
      listing_type: { code: "gold_special", label: "Clássico" },
      status: "active",
      link: { state: "resolved", product_id: "prod/1", seller_sku: "SKU-1" },
      price: { amount: "100.00", currency: "BRL" },
      published_quantity: 5,
      sync_state: "synced",
      sync_error: null,
      pending_issue: null,
      cost: { amount: "60.00", currency: "BRL" },
      below_margin_worst_case: false,
      icms_worst_case_by_uf: null,
      fetched_at: "2026-07-15T12:00:00Z",
    } satisfies ListingReadModel;
    const groups = { groups: [], next_cursor: null, page_size: 0, as_of: "2026-07-15T12:00:00Z" } satisfies ListingGroupPage;
    const summary = {
      total: 0,
      active: 0,
      paused: 0,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 0,
        below_margin_worst_case: null,
        margin_unknown: null,
      },
      as_of: "2026-07-15T12:00:00Z",
    } satisfies ListingSummary;
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        const path = String(input);
        const detail = { ...listing, timeline: [] } satisfies ListingDetail;
        const body = path.includes("/summary") ? summary : path.includes("/by-product") ? groups : path.includes("~") ? detail : page;
        return new Response(JSON.stringify(body), { status: 200, headers: { "Content-Type": "application/json" } });
      },
    });

    const list = await client.listListings({
      installation_id: "install/1",
      status: "active",
      sync_state: "stale",
      exception: "unlinked",
      has_exception: false,
      limit: 25,
      cursor: "cursor/1",
    });
    const grouped = await client.listListingsByProduct({ installation_id: "install/1", link_state: "resolved" });
    const detail = await client.getListing("install/1~provider/1~-");
    const typedSummary = await client.getListingsSummary("install/1");

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/listings?installation_id=install%2F1&filter.status=active&filter.sync_state=stale&filter.exception=unlinked&filter.has_exception=false&limit=25&cursor=cursor%2F1",
    );
    expect(String(requests[1].input)).toBe(
      "http://localhost:8080/listings/by-product?installation_id=install%2F1&filter.link_state=resolved",
    );
    expect(String(requests[2].input)).toBe("http://localhost:8080/listings/install%2F1~provider%2F1~-");
    expect(String(requests[3].input)).toBe("http://localhost:8080/listings/summary?installation_id=install%2F1");
    expect(requests.every((request) => request.init?.method === "GET")).toBe(true);
    expect(list.page_size).toBe(0);
    expect(grouped.groups).toEqual([]);
    expect(detail.timeline).toEqual([]);
    expect(typedSummary.exceptions.margin_unknown).toBeNull();
  });

  it("rejects non-positive product-link CODPROD before transport", () => {
    const fetchImpl = vi.fn();
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });
    const base = {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_item_id: "MLB1",
      actor: { actor_type: "operator", actor_id: "op-1" },
    };
    for (const internal_product_id of [-1, 0]) {
      expect(() => client.manualResolveProductLink({ ...base, internal_product_id })).toThrow("invalid_identity");
    }
    expect(fetchImpl).not.toHaveBeenCalled();
  });
  it("listIntegrationProviders calls /integrations/providers and parses Amazon LWA and interactive Shopee metadata", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const response = {
      items: [
        {
          provider_code: "amazon",
          tenant_id: "system",
          family: "marketplace",
          display_name: "Amazon",
          auth_strategy: "lwa",
          install_mode: "interactive",
          metadata: {
            country: "BR",
            rollout_stage: "wave_2",
            execution_mode: "available",
            fee_source: "api_sync",
            baseline_commission_percent: 16.0,
            baseline_fixed_fee_amount: 0,
            credential_schema: [
              { key: "lwa_client_id", label: "LWA Client ID", secret: false },
              { key: "lwa_client_secret", label: "LWA Client Secret", secret: true },
            ],
            docs_url: "https://developer-docs.amazon.com/sp-api",
          },
          declared_capabilities: [],
          is_active: true,
          created_at: "2026-04-09T00:00:00Z",
          updated_at: "2026-04-09T00:00:00Z",
        },
        {
          provider_code: "shopee",
          tenant_id: "system",
          family: "marketplace",
          display_name: "Shopee",
          auth_strategy: "shopee_partner",
          install_mode: "interactive",
          metadata: {
            country: "BR",
            rollout_stage: "v1",
            execution_mode: "available",
            fee_source: "seed",
            baseline_commission_percent: 14.5,
            baseline_fixed_fee_amount: 2.5,
            credential_schema: [
              { key: "partner_id", label: "Partner ID", secret: false },
              { key: "partner_key", label: "Partner Key", secret: true },
              { key: "shop_id", label: "Shop ID", secret: false },
            ],
          },
          declared_capabilities: ["pricing_fee_sync"],
          is_active: true,
          created_at: "2026-04-09T00:00:00Z",
          updated_at: "2026-04-09T00:00:00Z",
        },
      ],
    } satisfies { items: IntegrationProviderDefinition[] };
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify(response), { status: 200, headers: { "Content-Type": "application/json" } });
      },
    });

    const result = await client.listIntegrationProviders();

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/providers");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items).toHaveLength(2);

    const amazon = result.items.find((item) => item.provider_code === "amazon");
    const shopee = result.items.find((item) => item.provider_code === "shopee");

    expect(amazon).toBeDefined();
    expect(amazon?.auth_strategy).toBe("lwa");
    expect(amazon?.metadata?.execution_mode).toBe("available");
    expect(amazon?.metadata?.rollout_stage).toBe("wave_2");
    expect(amazon?.metadata?.baseline_commission_percent).toBe(16.0);
    expect(amazon?.metadata?.credential_schema?.[0]).toMatchObject({
      key: "lwa_client_id",
      secret: false,
    });

    expect(shopee).toBeDefined();
    expect(shopee?.metadata?.execution_mode).toBe("available");
    expect(shopee?.metadata?.rollout_stage).toBe("v1");
  });

  it("starts integration authorize flow with auth_url and expires_in", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            provider_code: "mercado_livre",
            state: "opaque-state",
            auth_url: "https://auth.example.com/authorize",
            expires_in: 300,
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.startIntegrationAuthorization("inst-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/auth/authorize");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.auth_url).toBe("https://auth.example.com/authorize");
    expect(result.expires_in).toBe(300);
  });

  it("lists integration installations with connection snapshot", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                installation_id: "inst-1",
                tenant_id: "tenant_default",
                provider_code: "mercado_livre",
                family: "marketplace",
                display_name: "Mercado Livre",
                status: "connected",
                health_status: "healthy",
                external_account_id: "691607102",
                external_account_name: "METALNOBREACABAMENTOS",
                active_credential_id: "cred-1",
                last_verified_at: "2026-07-08T12:00:00Z",
                connection: {
                  state: "connected",
                  health: "healthy",
                  provider_code: "mercado_livre",
                  external_account_id: "691607102",
                  external_account_name: "METALNOBREACABAMENTOS",
                  auth_strategy: "oauth2",
                  last_verified_at: "2026-07-08T12:00:00Z",
                  expires_at: "2026-07-08T18:00:00Z",
                  next_action: "none",
                },
                runtime_capabilities: [
                  {
                    code: "account_probe",
                    state: "available",
                    executable: true,
                    live_validated: true,
                    local_validated: true,
                  },
                ],
                created_at: "2026-07-08T12:00:00Z",
                updated_at: "2026-07-08T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listIntegrationInstallations();

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].connection.external_account_name).toBe("METALNOBREACABAMENTOS");
    expect(result.items[0].runtime_capabilities[0].code).toBe("account_probe");
  });

  it("starts integration reauthorization flow with auth_url and expires_in", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            provider_code: "mercado_livre",
            state: "reauth-state",
            auth_url: "https://auth.example.com/reauthorize",
            expires_in: 600,
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.startIntegrationReauthorization("inst-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/reauth/authorize");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.state).toBe("reauth-state");
    expect(result.auth_url).toBe("https://auth.example.com/reauthorize");
    expect(result.expires_in).toBe(600);
  });

  it("throws structured error for integration reauthorization failures", async () => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(
          JSON.stringify({
            error: { code: "conflict", message: "reauthorization not allowed" },
          }),
          { status: 409, headers: { "Content-Type": "application/json" } },
        ),
    });

    await expect(client.startIntegrationReauthorization("inst-1")).rejects.toMatchObject({
      status: 409,
      error: { code: "conflict", message: "reauthorization not allowed" },
    });
  });

  it("submits integration credentials as json and parses auth status", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            status: "connected",
            health_status: "warning",
            provider_code: "mercado_livre",
            external_account_id: "acct-1",
            external_account_name: "Acct 1",
            connection: {
              state: "connected",
              health: "warning",
              provider_code: "mercado_livre",
              external_account_id: "acct-1",
              external_account_name: "Acct 1",
              auth_strategy: "api_key",
              next_action: "none",
            },
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.submitIntegrationCredentials("inst-1", {
      api_key: "secret-key",
      metadata: { environment: "sandbox" },
      credentials: { username: "user", password: "pass" },
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/auth/credentials");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
    expect(requests[0].init?.body).toBe(
      JSON.stringify({
        api_key: "secret-key",
        metadata: { environment: "sandbox" },
        credentials: { username: "user", password: "pass" },
      }),
    );
    expect(result.status).toBe("connected");
    expect(result.health_status).toBe("warning");
    expect(result.connection.external_account_name).toBe("Acct 1");
  });

  it("throws structured error for integration credential submission failures", async () => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(
          JSON.stringify({
            error: { code: "invalid_request", message: "missing credentials" },
          }),
          { status: 400, headers: { "Content-Type": "application/json" } },
        ),
    });

    await expect(
      client.submitIntegrationCredentials("inst-1", {
        credentials: { username: "user" },
      }),
    ).rejects.toMatchObject({
      status: 400,
      error: { code: "invalid_request", message: "missing credentials" },
    });
  });

  it("disconnects integration installation and parses auth status", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            status: "disconnected",
            health_status: "critical",
            connection: {
              state: "disconnected",
              health: "critical",
              provider_code: "mercado_livre",
              external_account_id: "",
              external_account_name: "",
              auth_strategy: "unknown",
              next_action: "authorize",
            },
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.disconnectIntegrationInstallation("inst-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/disconnect");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.status).toBe("disconnected");
    expect(result.health_status).toBe("critical");
    expect(result.connection.next_action).toBe("authorize");
  });

  it("throws structured error for disconnect failures", async () => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(
          JSON.stringify({
            error: { code: "conflict", message: "installation already disconnected" },
          }),
          { status: 409, headers: { "Content-Type": "application/json" } },
        ),
    });

    await expect(client.disconnectIntegrationInstallation("inst-1")).rejects.toMatchObject({
      status: 409,
      error: { code: "conflict", message: "installation already disconnected" },
    });
  });

  it("gets integration auth status", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            status: "connected",
            health_status: "healthy",
            provider_code: "mercado_livre",
            external_account_id: "acct-1",
            external_account_name: "Acct 1",
            connection: {
              state: "connected",
              health: "healthy",
              provider_code: "mercado_livre",
              external_account_id: "acct-1",
              external_account_name: "Acct 1",
              auth_strategy: "oauth2",
              next_action: "none",
            },
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.getIntegrationAuthStatus("inst-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/auth/status");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.status).toBe("connected");
    expect(result.health_status).toBe("healthy");
    expect(result.connection.auth_strategy).toBe("oauth2");
  });

  it("probes integration account", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            provider_code: "mercado_livre",
            provider_account_id: "691607102",
            provider_account_name: "METALNOBREACABAMENTOS",
            site_id: "MLB",
            status: "active",
            fetched_at: "2026-07-08T12:00:00Z",
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.probeIntegrationAccount("inst-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/probes/account");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.provider_account_name).toBe("METALNOBREACABAMENTOS");
  });

  it("reads integration listings probe", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                provider_code: "mercado_livre",
                provider_item_id: "MLB123",
                title: "Produto teste",
                fetched_at: "2026-07-08T12:00:00Z",
                variations: [
                  {
                    provider_variation_id: "VAR-1",
                    available_quantity: 5,
                  },
                ],
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.probeIntegrationListings("inst-1", 3);

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/probes/listings?limit=3");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].variations?.[0].provider_variation_id).toBe("VAR-1");
  });

  it("reads integration orders probe", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                provider_code: "mercado_livre",
                provider_order_id: "ORDER-1",
                fetched_at: "2026-07-08T12:00:00Z",
                items: [
                  {
                    provider_item_id: "MLB123",
                    quantity: 2,
                  },
                ],
                payments: [
                  {
                    payment_id: "PAY-1",
                    amount: 120.5,
                  },
                ],
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.probeIntegrationOrders("inst-1", 1);

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/probes/orders?limit=1");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].payments[0].payment_id).toBe("PAY-1");
  });

  it("reads integration fee quote probe", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            provider_code: "mercado_livre",
            site_id: "MLB",
            category_id: "",
            listing_type_id: "gold_special",
            price_amount: 100,
            currency_id: "BRL",
            commission_percent: 11,
            fixed_fee_amount: 6,
            fetched_at: "2026-07-08T12:00:00Z",
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.probeIntegrationFeeQuote("inst-1", {
      listing_type_id: "gold_special",
      price: 100,
      currency_id: "BRL",
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/integrations/installations/inst-1/probes/fee-quote?listing_type_id=gold_special&price=100&currency_id=BRL");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.fixed_fee_amount).toBe(6);
  });

  it("imports product link listing snapshots", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            imported_count: 2,
            items: [
              {
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_item_id: "MLB123",
                provider_variation_id: "VAR-1",
                seller_sku: "SKU-1",
                fetched_at: "2026-07-08T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.importProductLinkListingSnapshots({
      installation_id: "inst-1",
      limit: 5,
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/product-links/listing-snapshots/imports");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
    expect(requests[0].init?.body).toBe(JSON.stringify({ installation_id: "inst-1", limit: 5 }));
    expect(result.imported_count).toBe(2);
    expect(result.items[0].provider_variation_id).toBe("VAR-1");
  });

  it("refreshes listings and preserves the active operation on conflict", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    let conflict = false;
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        if (conflict) {
          return new Response(JSON.stringify({ error: { code: "refresh_in_progress", message: "busy", details: { operation_run_id: "op_active" } } }), { status: 409, headers: { "Content-Type": "application/json" } });
        }
        return new Response(JSON.stringify({ operation_run_id: "op_new" }), { status: 202, headers: { "Content-Type": "application/json" } });
      },
    });

    await expect(client.refreshListings({ installation_id: "inst_test" })).resolves.toEqual({ operation_run_id: "op_new" });
    expect(String(requests[0].input)).toBe("http://localhost:8080/listings/refresh");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.body).toBe(JSON.stringify({ installation_id: "inst_test" }));

    conflict = true;
    await expect(client.refreshListings({ installation_id: "inst_test" })).rejects.toMatchObject({ status: 409, error: { code: "refresh_in_progress", details: { operation_run_id: "op_active" } } });
  });

  it("generates product link candidates as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            generated_count: 1,
            items: [
              {
                candidate_id: "cand-1",
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_item_id: "MLB123",
                provider_variation_id: "VAR-1",
                internal_product_id: 101,
                internal_product_name: "Produto teste",
                internal_reference_code: "SKU-1",
                state: "exact_sku",
                match_input: "seller_sku",
                match_value: "SKU-1",
                source_snapshot_fetched_at: "2026-07-08T11:59:00Z",
                created_at: "2026-07-08T12:00:00Z",
                updated_at: "2026-07-08T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.generateProductLinkCandidates({
      installation_id: "inst-1",
      limit: 10,
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/product-links/link-candidates/generations");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
    expect(requests[0].init?.body).toBe(JSON.stringify({ installation_id: "inst-1", limit: 10 }));
    expect(result.generated_count).toBe(1);
    expect(result.items[0].state).toBe("exact_sku");
    expect(result.items[0].match_input).toBe("seller_sku");
  });

  it("lists product link candidates with installation filter", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                candidate_id: "cand-1",
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_item_id: "MLB123",
                state: "unresolved",
                match_input: "none",
                created_at: "2026-07-08T12:00:00Z",
                updated_at: "2026-07-08T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listProductLinkCandidates("inst-1", 25);

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/product-links/link-candidates?installation_id=inst-1&limit=25",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].candidate_id).toBe("cand-1");
    expect(result.items[0].state).toBe("unresolved");
  });

  it("lists product link workflows with installation filter", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                identity: {
                  installation_id: "inst-1",
                  provider_item_id: "MLB123",
                  provider_variation_id: "VAR-1",
                },
                current_link: {
                  installation_id: "inst-1",
                  provider_code: "mercado_livre",
                  provider_item_id: "MLB123",
                  provider_variation_id: "VAR-1",
                  state: "resolved",
                  internal_product_id: 101,
                  internal_product_name: "Produto teste",
                  internal_reference_code: "SKU-1",
                  created_at: "2026-07-08T12:00:00Z",
                  updated_at: "2026-07-08T12:10:00Z",
                },
                candidates: [],
                audit: [],
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listProductLinkWorkflows("inst-1", 10);

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/product-links/link-workflows?installation_id=inst-1&limit=10",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].identity.provider_item_id).toBe("MLB123");
    expect(result.items[0].current_link?.state).toBe("resolved");
  });

  it("approves product link candidates as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            link: {
              installation_id: "inst-1",
              provider_code: "mercado_livre",
              provider_item_id: "MLB123",
              state: "resolved",
              source_candidate_id: "cand-1",
              internal_product_id: 101,
              internal_product_name: "Produto teste",
              internal_reference_code: "SKU-1",
              created_at: "2026-07-08T12:00:00Z",
              updated_at: "2026-07-08T12:10:00Z",
            },
            audit: {
              audit_id: "audit-1",
              installation_id: "inst-1",
              provider_code: "mercado_livre",
              provider_item_id: "MLB123",
              action: "approve_candidate",
              actor: {
                actor_type: "operator",
                actor_id: "leandro",
                actor_name: "Leandro",
              },
              previous_state: "unresolved",
              next_state: "resolved",
              next_internal_product_id: 101,
              created_at: "2026-07-08T12:10:00Z",
            },
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.approveProductLinkCandidate({
      candidate_id: "cand-1",
      reason: "Exact EAN match",
      actor: { actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" },
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/product-links/link-resolutions/approve-candidate");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.body).toBe(
      JSON.stringify({
        candidate_id: "cand-1",
        reason: "Exact EAN match",
        actor: { actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" },
      }),
    );
    expect(result.link.internal_product_id).toBe(101);
    expect(result.audit.action).toBe("approve_candidate");
  });

  it("lists inventory stock risks with filters", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                identity: {
                  installation_id: "inst-1",
                  provider_item_id: "MLB123",
                },
                provider_code: "mercado_livre",
                link_state: "resolved",
                state: "oversell",
                actionability: "actionable",
                actionable: true,
                provider_quantity: 9,
                recommended_quantity: 7,
                policy_id: "stock-seguro-default",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listInventoryStockRisks({
      installation_id: "inst-1",
      state: "oversell",
      actionability: "actionable",
      limit: 20,
    });

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/inventory/stock-risks?installation_id=inst-1&state=oversell&actionability=actionable&limit=20",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].state).toBe("oversell");
    expect(result.items[0].recommended_quantity).toBe(7);
  });

  it("applies inventory manual stock action as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            action: {
              action_id: "act-1",
              state: "applied",
              trigger: "manual",
              provider_ref: {
                tenant_id: "tenant_default",
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_account_id: "acct-1",
                provider_item_id: "MLB123",
              },
              requested_quantity: 7,
              recommended_quantity: 7,
              policy_id: "stock-seguro-default",
              operator: {
                actor_type: "operator",
                actor_id: "leandro",
                actor_name: "Leandro",
              },
              idempotency_key: "act-1",
              audit_events: [],
              created_at: "2026-07-09T12:00:00Z",
              updated_at: "2026-07-09T12:00:00Z",
            },
            risk: {
              identity: {
                installation_id: "inst-1",
                provider_item_id: "MLB123",
              },
              provider_code: "mercado_livre",
              link_state: "resolved",
              state: "healthy",
              actionability: "blocked",
              actionable: false,
              provider_quantity: 7,
              recommended_quantity: 7,
              policy_id: "stock-seguro-default",
            },
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.applyInventoryManualStockAction({
      stock_action_id: "act-1",
      installation_id: "inst-1",
      provider_item_id: "MLB123",
      requested_quantity: 7,
      reason: "Operator apply",
      approval: {
        approved: true,
        approved_at: "2026-07-09T12:00:00Z",
        operator: {
          actor_type: "operator",
          actor_id: "leandro",
          actor_name: "Leandro",
        },
      },
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/inventory/stock-actions/manual-apply");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.body).toBe(
      JSON.stringify({
        stock_action_id: "act-1",
        installation_id: "inst-1",
        provider_item_id: "MLB123",
        requested_quantity: 7,
        reason: "Operator apply",
        approval: {
          approved: true,
          approved_at: "2026-07-09T12:00:00Z",
          operator: {
            actor_type: "operator",
            actor_id: "leandro",
            actor_name: "Leandro",
          },
        },
      }),
    );
    expect(result.action.state).toBe("applied");
    expect(result.risk.state).toBe("healthy");
  });

  it("imports marketplace orders as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            imported_count: 1,
            skipped_count: 0,
            items: [
              {
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_order_id: "2001",
                provider_status: "paid",
                provider_updated_at: "2026-07-09T12:00:00Z",
                fetched_at: "2026-07-09T12:00:00Z",
                items: [
                  {
                    provider_item_id: "MLB123",
                    quantity: 1,
                    link_quality: "resolved",
                  },
                ],
                payments: [],
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.importMarketplaceOrders({ installation_id: "inst-1", limit: 1 });

    expect(String(requests[0].input)).toBe("http://localhost:8080/orders/import");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.body).toBe(JSON.stringify({ installation_id: "inst-1", limit: 1 }));
    expect(result.imported_count).toBe(1);
    expect(result.items[0].provider_order_id).toBe("2001");
  });

  it("lists marketplace orders with normalized snapshots", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                installation_id: "inst-1",
                provider_code: "mercado_livre",
                provider_order_id: "2001",
                provider_status: "paid",
                provider_updated_at: "2026-07-09T12:00:00Z",
                fetched_at: "2026-07-09T12:00:00Z",
                items: [
                  {
                    provider_item_id: "MLB123",
                    quantity: 1,
                    link_quality: "resolved",
                    internal_product_id: 321,
                  },
                ],
                payments: [
                  {
                    provider_payment_id: "pay-1",
                    total_paid_amount: 125.92,
                  },
                ],
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listMarketplaceOrders("inst-1", 5);

    expect(String(requests[0].input)).toBe("http://localhost:8080/orders?installation_id=inst-1&limit=5");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.items[0].items[0].link_quality).toBe("resolved");
    expect(result.items[0].payments[0].total_paid_amount).toBe(125.92);
  });

  it("lists canonical orders with cursor and filter query parameters", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                provider_order_id: "order-1",
                provider_code: "mercado_livre",
                status: "paid",
                provider_status_detail: "",
                buyer_nickname: null,
                total: null,
                currency: null,
                fulfillment: null,
                nf_state: null,
                created_at: null,
                provider_created_at: null,
                provider_closed_at: null,
                provider_updated_at: null,
                items: [],
                payments: [],
              },
            ],
            next_cursor: null,
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listOrders({
      installation_id: "inst/1",
      cursor: "cursor+with/slash=",
      status: "paid & shipped",
      date_from: "2026-07-01",
      date_to: "2026-07-16T12:00:00Z",
      q: "buyer/order?1",
    });

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/orders?installation_id=inst%2F1&cursor=cursor%2Bwith%2Fslash%3D&status=paid+%26+shipped&date_from=2026-07-01&date_to=2026-07-16T12%3A00%3A00Z&q=buyer%2Forder%3F1",
    );
    const read: OrderRead = result.items[0];
    expect(read.provider_order_id).toBe("order-1");
    expect(result.items[0].buyer_nickname).toBeNull();
    expect(result.items[0].total).toBeNull();
    const nullableCurrency: string | null = result.items[0].currency;
    const nullableFulfillment: string | null = result.items[0].fulfillment;
    const nullableNFState: string | null = result.items[0].nf_state;
    expect(nullableCurrency).toBeNull();
    expect(nullableFulfillment).toBeNull();
    expect(nullableNFState).toBeNull();
  });

  it("gets an order with both identifiers URL-encoded", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            provider_order_id: "order/with ?chars",
            provider_code: "mercado_livre",
            status: "paid",
            provider_status_detail: "",
            buyer_nickname: null,
            total: null,
            currency: null,
            fulfillment: null,
            nf_state: null,
            created_at: null,
            provider_created_at: null,
            provider_closed_at: null,
            provider_updated_at: null,
            items: [],
            payments: [],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.getOrder("installation/1", "order/with ?chars");

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/orders/order%2Fwith%20%3Fchars?installation_id=installation%2F1",
    );
    const read: OrderRead = result;
    expect(read.provider_code).toBe("mercado_livre");
    expect(result.provider_order_id).toBe("order/with ?chars");
  });

  it("imports profitability margin inputs as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            imported_count: 2,
            items: [
              {
                installation_id: "inst-1",
                provider_order_id: "2001",
                scope: "item",
                kind: "revenue",
                amount: 19.9,
                currency: "BRL",
                source_system: "orders",
                quality: "complete",
                created_at: "2026-07-09T12:00:00Z",
                updated_at: "2026-07-09T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.importProfitabilityMarginInputs({ installation_id: "inst-1", limit: 1 });

    expect(String(requests[0].input)).toBe("http://localhost:8080/profitability/margin-inputs/import");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.imported_count).toBe(2);
    expect(result.items[0].kind).toBe("revenue");
  });

  it("lists profitability margin inputs", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                installation_id: "inst-1",
                provider_order_id: "2001",
                scope: "item",
                kind: "cost",
                currency: "BRL",
                source_system: "internal_read",
                quality: "missing",
                quality_reason: "missing_cost",
                created_at: "2026-07-09T12:00:00Z",
                updated_at: "2026-07-09T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listProfitabilityMarginInputs("inst-1", 5);

    expect(String(requests[0].input)).toBe("http://localhost:8080/profitability/margin-inputs?installation_id=inst-1&limit=5");
    expect(result.items[0].quality).toBe("missing");
  });

  it("creates profitability manual adjustment as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            adjustment_id: "adj-1",
			idempotency_key: "manual-adjustment-1",
            installation_id: "inst-1",
            provider_order_id: "2001",
            scope: "order",
            category: "freight",
            amount: 12.5,
            currency: "BRL",
            reason: "manual freight",
            actor: { actor_type: "operator", actor_id: "leandro" },
            created_at: "2026-07-09T12:00:00Z",
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.createProfitabilityManualAdjustment({
      installation_id: "inst-1",
		idempotency_key: "manual-adjustment-1",
      provider_order_id: "2001",
      scope: "order",
      category: "freight",
      amount: 12.5,
      reason: "manual freight",
      actor: { actor_type: "operator", actor_id: "leandro" },
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/profitability/manual-adjustments");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual({
      installation_id: "inst-1",
		idempotency_key: "manual-adjustment-1",
      provider_order_id: "2001",
      scope: "order",
      category: "freight",
      amount: 12.5,
      reason: "manual freight",
      actor: { actor_type: "operator", actor_id: "leandro" },
    });
    expect(result.category).toBe("freight");
  });

  it("calculates profitability profit snapshots as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "inst-1",
            calculated_count: 1,
            items: [
              {
                installation_id: "inst-1",
                provider_order_id: "2001",
                scope: "item",
                currency: "BRL",
                realization_state: "realized",
                revenue_amount: 20,
                sale_fee_amount: 2,
                cost_amount: 10,
                tax_amount: 1,
                adjustment_amount: 1,
                contribution_amount: 8,
                margin_percent: 40,
                quality: "complete",
                flags: [],
                created_at: "2026-07-09T12:00:00Z",
                updated_at: "2026-07-09T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.calculateProfitabilityProfitSnapshots({ installation_id: "inst-1", limit: 1 });

    expect(String(requests[0].input)).toBe("http://localhost:8080/profitability/profit-snapshots/calculate");
    expect(requests[0].init?.method).toBe("POST");
    expect(result.items[0].realization_state).toBe("realized");
    expect(result.items[0].contribution_amount).toBe(8);
  });

  it("lists profitability profit snapshots", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            items: [
              {
                installation_id: "inst-1",
                provider_order_id: "2001",
                scope: "item",
                currency: "BRL",
                realization_state: "not_realized",
                contribution_amount: null,
                margin_percent: null,
                quality: "not_realized",
                flags: ["order_cancelled"],
                created_at: "2026-07-09T12:00:00Z",
                updated_at: "2026-07-09T12:00:00Z",
              },
            ],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const result = await client.listProfitabilityProfitSnapshots("inst-1", 5);

    expect(String(requests[0].input)).toBe("http://localhost:8080/profitability/profit-snapshots?installation_id=inst-1&limit=5");
    expect(result.items[0].realization_state).toBe("not_realized");
    expect(result.items[0].quality).toBe("not_realized");
    expect(result.items[0].flags).toContain("order_cancelled");
    expect(result.items[0].contribution_amount).toBeNull();
    expect(result.items[0].margin_percent).toBeNull();
  });

  it("builds canonical pricing simulation requests", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ items: [] }), { status: 200 });
      },
    });

    await client.listPricingSimulations();

    expect(String(requests[0].input)).toBe("http://localhost:8080/pricing/simulations");
    expect(requests[0].init?.method).toBe("GET");
  });

  it("posts marketplace account payload as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            tenant_id: "tenant_default",
            account_id: "acct-1",
            channel_code: "mercado_livre",
            display_name: "Mercado Livre",
            status: "active",
            connection_mode: "api",
          }),
          { status: 201, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    await client.createMarketplaceAccount({
      account_id: "acct-1",
      channel_code: "mercado_livre",
      display_name: "Mercado Livre",
      connection_mode: "api",
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/marketplaces/accounts");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
    expect(requests[0].init?.body).toBe(
      JSON.stringify({
        account_id: "acct-1",
        channel_code: "mercado_livre",
        display_name: "Mercado Livre",
        connection_mode: "api",
      }),
    );
  });

  it("posts run pricing simulation payload as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            simulation_id: "sim-1",
            tenant_id: "tenant_default",
            product_id: "prod-1",
            account_id: "acct-1",
            margin_amount: 10.5,
            margin_percent: 15.0,
            status: "healthy",
          }),
          { status: 201, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    await client.runPricingSimulation({
      simulation_id: "sim-1",
      product_id: "prod-1",
      account_id: "acct-1",
      base_price_amount: 100.0,
      cost_amount: 60.0,
      commission_percent: 0.16,
      fixed_fee_amount: 5.0,
      shipping_amount: 10.0,
      min_margin_percent: 0.10,
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/pricing/simulations");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
  });

  it("posts batch simulation payload as json", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ items: [] }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      },
    });

    await client.runBatchSimulation({
      product_ids: ["prod-1"],
      policy_ids: ["policy-1"],
      origin_cep: "01001000",
      destination_cep: "20040002",
      price_source: "my_price",
      price_overrides: { "prod-1::policy-1": 123.45 },
    });

    expect(String(requests[0].input)).toBe("http://localhost:8080/pricing/simulations/batch");
    expect(requests[0].init?.method).toBe("POST");
    expect(requests[0].init?.headers).toEqual({ "Content-Type": "application/json" });
    expect(requests[0].init?.body).toBe(
      JSON.stringify({
        product_ids: ["prod-1"],
        policy_ids: ["policy-1"],
        origin_cep: "01001000",
        destination_cep: "20040002",
        price_source: "my_price",
        price_overrides: { "prod-1::policy-1": 123.45 },
      }),
    );
  });

  it("gets melhor envio status", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ connected: true }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      },
    });

    const result = await client.getMelhorEnvioStatus();

    expect(String(requests[0].input)).toBe("http://localhost:8080/connectors/melhor-envio/status");
    expect(requests[0].init?.method).toBe("GET");
    expect(result.connected).toBe(true);
  });

  it("createMutation posts the command request", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP-000001" }), { status: 201 });
      },
    });
    const request = {
      installation_id: "install-1",
      type: "price_update" as const,
      actor: "operator-1",
      intent: { price: { amount: "49.90", currency: "BRL" } },
      selection: { listing_ids: ["MLB-1"] },
    };

    await client.createMutation(request);

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual(request);
  });

  it("previewMutation posts to the encoded protocol path", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP/1", items: [] }), { status: 200 });
      },
    });

    await client.previewMutation("MP/1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations/MP%2F1/preview");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual({});
  });

  it("approveMutation posts the execute confirmation", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP-1" }), { status: 200 });
      },
    });

    await client.approveMutation("MP-1", { execute: true });

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations/MP-1/approve");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual({ execute: true });
  });

  it("cancelMutation posts to the protocol command path", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP-1" }), { status: 200 });
      },
    });

    await client.cancelMutation("MP-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations/MP-1/cancel");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual({});
  });

  it("retryMutationFailures posts to the protocol command path", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP-2" }), { status: 200 });
      },
    });

    await client.retryMutationFailures("MP-1");

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations/MP-1/retry");
    expect(requests[0].init?.method).toBe("POST");
    expect(JSON.parse(String(requests[0].init?.body))).toEqual({});
  });

  it("listMutations gets an encoded filtered protocol page", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ items: [{ protocol_id: "MP-1" }], next_cursor: null, page_size: 1 }), { status: 200 });
      },
    });

    const page = await client.listMutations({
      installation_id: "install/1",
      state: "partially_failed",
      type: "price_update",
      limit: 200,
      cursor: "opaque+/=",
    });

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/mutations?installation_id=install%2F1&state=partially_failed&type=price_update&limit=200&cursor=opaque%2B%2F%3D",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(page.items[0].protocol_id).toBe("MP-1");
    expect(page.next_cursor).toBeNull();
  });

  it("getMutation gets an encoded protocol path with typed response", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ protocol_id: "MP/1?" }), { status: 200 });
      },
    });

    const protocol = await client.getMutation("MP/1?");

    expect(String(requests[0].input)).toBe("http://localhost:8080/mutations/MP%2F1%3F");
    expect(requests[0].init?.method).toBe("GET");
    expect(protocol.protocol_id).toBe("MP/1?");
  });

  it("listMutationItems gets an encoded filtered item page", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ items: [{ seq: 1 }], next_cursor: "next", page_size: 1 }), { status: 200 });
      },
    });

    const page = await client.listMutationItems("MP/1", {
      state: "failed",
      limit: 50,
      cursor: "opaque+/=",
    });

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/mutations/MP%2F1/items?state=failed&limit=50&cursor=opaque%2B%2F%3D",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(page.items[0].seq).toBe(1);
    expect(page.next_cursor).toBe("next");
  });

  it("throws parsed error payload on non-ok response", async () => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(
          JSON.stringify({
            error: { code: "invalid_request", message: "invalid account" },
          }),
          { status: 400, headers: { "Content-Type": "application/json" } },
        ),
    });

    await expect(
      client.createMarketplaceAccount({
        account_id: "",
        channel_code: "mercado_livre",
        display_name: "Mercado Livre",
        connection_mode: "api",
      }),
    ).rejects.toMatchObject({
      status: 400,
      error: { code: "invalid_request", message: "invalid account" },
    });
  });

  it("throws parsed error payload on non-ok GET response", async () => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(
          JSON.stringify({
            error: { code: "not_found", message: "no products found" },
          }),
          { status: 404, headers: { "Content-Type": "application/json" } },
        ),
    });

    await expect(client.listCatalogProducts()).rejects.toMatchObject({
      status: 404,
      error: { code: "not_found", message: "no products found" },
    });
  });

  it("uses exact encoded scope for assisted Sankhya current and candidate reads", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        const candidates = String(input).endsWith("/candidates?installation_id=install%2F1");
        return new Response(
          JSON.stringify(
            candidates
              ? { installation_id: "install/1", provider_order_id: "order/1", candidates: [] }
              : {
                  installation_id: "install/1",
                  provider_order_id: "order/1",
                  internal_document_id: 31301,
                  lines: [],
                  audit: {
                    event_id: "event-1",
                    event_type: "confirmed",
                    actor_type: "operator_supplied_unverified",
                    actor_id: "operator-1",
                    reason: "exact",
                    source_at: "2026-07-13T17:00:00Z",
                    recorded_at: "2026-07-13T17:01:00Z",
                    configuration_revision: "cfg-1",
                    evidence_state: "exact",
                    evidence_reference: "sankhya-linkage:v1:safe",
                  },
                },
          ),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });

    const current = await client.getAssistedSankhyaLinkage("install/1", "order/1");
    const candidates = await client.listAssistedSankhyaLinkageCandidates("install/1", "order/1");

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/orders/order%2F1/sankhya-linkage?installation_id=install%2F1",
    );
    expect(String(requests[1].input)).toBe(
      "http://localhost:8080/orders/order%2F1/sankhya-linkage/candidates?installation_id=install%2F1",
    );
    expect(current.audit.evidence_reference).toBe("sankhya-linkage:v1:safe");
    expect(candidates.candidates).toEqual([]);
  });

  it("confirms assisted Sankhya linkage with operator facts only", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(
          JSON.stringify({
            installation_id: "install-1",
            provider_order_id: "order-1",
            internal_document_id: 31301,
            lines: [],
            audit: {
              event_id: "server-event-1",
              event_type: "confirmed",
              actor_type: "operator_supplied_unverified",
              actor_id: "operator-7",
              reason: "exact selection",
              source_at: "2026-07-13T17:00:00Z",
              recorded_at: "2026-07-13T17:01:00Z",
              configuration_revision: "cfg-runtime-1",
              evidence_state: "exact",
              evidence_reference: "sankhya-linkage:v1:safe",
            },
            lineage: [{ mpc_line_id: "mpl_1", internal_document_id: 31301, internal_line_number: 1, state: "none", descendants: [] }],
          }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      },
    });
    const request = {
      selected_document_id: 31301,
      selections: [{ mpc_line_id: "mpl_1", internal_document_id: 31301, internal_line_number: 1 }],
      actor_id: "operator-7",
      reason: "exact selection",
      source_at: "2026-07-13T17:00:00Z",
      idempotency_key: "idem-1",
    };

    const result = await client.confirmAssistedSankhyaLinkage("install-1", "order-1", request);
    const sent = JSON.parse(String(requests[0].init?.body));

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/orders/order-1/sankhya-linkage/confirm?installation_id=install-1",
    );
    expect(Object.keys(sent).sort()).toEqual(
      ["actor_id", "idempotency_key", "reason", "selected_document_id", "selections", "source_at"].sort(),
    );
    expect(sent).not.toHaveProperty("tenant_id");
    expect(sent).not.toHaveProperty("event_id");
    expect(sent).not.toHaveProperty("configuration_revision");
    expect(sent).not.toHaveProperty("evidence_reference");
    expect(sent).not.toHaveProperty("actor_type");
    expect(sent).not.toHaveProperty("external_order_key");
    expect(result.audit.event_id).toBe("server-event-1");
    expect(result.lineage[0].state).toBe("none");
  });

  it("reads dashboard summary and sync runs with encoded query values and nullable fields", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const summary = {
      sync_errors: null,
      pending_links: 2,
      below_margin: null,
      missing_gtin: 1,
      orders_today: null,
      orders_7d: 4,
      last_sync_at: null,
      degraded: ["linkage"],
      as_of: "2026-07-16T12:01:00Z",
    } satisfies DashboardSummary;
    const runs = {
      items: [{
        operation_run_id: "run/1",
        installation_id: "install/1",
        module: "listing/read",
        status: "running",
        result_code: "",
        failure_code: "",
        translated_error_code: "",
        attempt_count: 1,
        duration_ms: 0,
        started_at: "2026-07-16T12:00:00Z",
        finished_at: null,
      }],
      next_cursor: null,
      page_size: 1,
      as_of: "2026-07-16T12:01:00Z",
    } satisfies SyncRunPage;
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify(String(input).includes("/dashboard/") ? summary : runs), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      },
    });

    const dashboard = await client.getDashboardSummary("install/1");
    const syncRuns = await client.listSyncRuns({
      installation_id: "install/1",
      cursor: "cursor/1",
      limit: 25,
      module: "listing/read",
      status: "running",
    });

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/dashboard/summary?installation_id=install%2F1",
    );
    expect(String(requests[1].input)).toBe(
      "http://localhost:8080/sync/runs?installation_id=install%2F1&cursor=cursor%2F1&limit=25&filter.module=listing%2Fread&filter.status=running",
    );
    expect(requests.every((request) => request.init?.method === "GET")).toBe(true);
    expect(dashboard.sync_errors).toBeNull();
    expect(dashboard.last_sync_at).toBeNull();
    expect(syncRuns.items[0].finished_at).toBeNull();
  });

  it("gets category attributes with an encoded path segment and preserves provider order", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const attributes: CategoryAttribute[] = [
      {
        id: "color",
        name: "Color",
        required: false,
        value_type: "string",
        values: null,
        constraints: null,
      },
      {
        id: "brand",
        name: "Brand",
        required: true,
        value_type: "enum",
        values: [{ id: "brand-b", name: "Brand B" }, { id: "brand-a", name: "Brand A" }],
        constraints: { max_length: 80 },
      },
    ];
    const response: CategoryAttributesResponse = {
      category_id: "MLB/123?draft",
      attributes,
    };
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify(response), { status: 200 });
      },
    });

    const result = await client.getCategoryAttributes("MLB/123?draft");

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/listings/categories/MLB%2F123%3Fdraft/attributes",
    );
    expect(requests[0].init?.method).toBe("GET");
    expect(result.attributes).toEqual(attributes);
    expect(result.attributes.map((attribute) => attribute.id)).toEqual(["color", "brand"]);
    expect(result.attributes[0].constraints).toBeNull();
    expect(result.attributes[1].constraints).toEqual({ max_length: 80 });
  });

  it.each([
    [404, "category_not_found"],
    [502, "provider_unavailable"],
  ] as const)("preserves category attribute error code %s", async (status, code) => {
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () =>
        new Response(JSON.stringify({ error: { code, message: "category attributes unavailable" } }), { status }),
    });

    await expect(client.getCategoryAttributes("MLB/category")).rejects.toMatchObject({
      status,
      error: { code },
    });
  });

  it("lists market observations and references with encoded ordered ids", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        return new Response(JSON.stringify({ items: [], as_of: "2026-07-16T12:00:00Z" }), { status: 200 });
      },
    });

    await client.listMarketObservations("install/1", ["listing/2", "listing/1"]);
    await client.listMarketReferences(["product/2", "product/1"]);

    expect(String(requests[0].input)).toBe(
      "http://localhost:8080/market/observations?installation_id=install%2F1&listing_ids=listing%2F2%2Clisting%2F1",
    );
    expect(String(requests[1].input)).toBe(
      "http://localhost:8080/market/references?product_ids=product%2F2%2Cproduct%2F1",
    );
    expect(requests.every((request) => request.init?.method === "GET")).toBe(true);
  });

  it("round-trips honest-empty market evidence with nullable provenance", async () => {
    const observation: MarketObservation = {
      listing_id: "listing-empty",
      our_sale_price: null,
      competitive_status: null,
      winner_price: null,
      competitive_target: null,
      catalog_offer_price: null,
      catalog_stats: null,
      evidence_state: "no_price_evidence",
      captured_at: null,
      source: null,
    };
    const reference: MarketReference = {
      product_id: "product-empty",
      catalog_product_id: null,
      match_state: "no_candidate",
      match_method: null,
      catalog_stats: null,
      evidence_state: "no_price_evidence",
      captured_at: null,
      source: null,
    };
    const responses: [MarketObservationPage, MarketReferencePage] = [
      { items: [observation], as_of: "2026-07-16T12:00:00Z" },
      { items: [reference], as_of: "2026-07-16T12:00:00Z" },
    ];
    let responseIndex = 0;
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () => new Response(JSON.stringify(responses[responseIndex++]), { status: 200 }),
    });

    const observationResult = await client.listMarketObservations("install-1", ["listing-empty"]);
    const referenceResult = await client.listMarketReferences(["product-empty"]);

    expect(observationResult.items[0]).toEqual(observation);
    expect(observationResult.items[0].our_sale_price).toBeNull();
    expect(observationResult.items[0].captured_at).toBeNull();
    expect(observationResult.items[0].source).toBeNull();
    expect(referenceResult.items[0]).toEqual(reference);
    expect(referenceResult.items[0].catalog_product_id).toBeNull();
    expect(referenceResult.items[0].captured_at).toBeNull();
    expect(referenceResult.items[0].source).toBeNull();
  });

  it("keeps the six market observation signals distinct", () => {
    const observation = {
      listing_id: "listing-1",
      our_sale_price: { amount: "100.00", currency: "BRL" },
      competitive_status: "competing",
      winner_price: { amount: "98.00", currency: "BRL" },
      competitive_target: { amount: "99.00", currency: "BRL" },
      catalog_offer_price: { amount: "101.00", currency: "BRL" },
      catalog_stats: {
        min: "97.00",
        p25: "98.00",
        median: "100.00",
        trimmed_mean: "100.25",
        p75: "102.00",
        max: "105.00",
        n_offers: 8,
        n_sellers: 6,
      },
      evidence_state: "observed",
      captured_at: "2026-07-16T12:00:00Z",
      source: "official_api",
    } satisfies MarketObservation;

    expect(observation.our_sale_price.amount).toBe("100.00");
    expect(observation.winner_price.amount).toBe("98.00");
    expect(observation.competitive_target.amount).toBe("99.00");
    expect(observation.catalog_offer_price.amount).toBe("101.00");
    expect(observation.catalog_stats.median).toBe("100.00");
  });

  it("drives the IC-04 pricing calculator endpoints (profile, difal, scenarios, decompose, solve)", async () => {
    const requests: Array<{ input: RequestInfo | URL; init?: RequestInit }> = [];
    const profile: PricingCalcProfile = {
      regime: "SIMPLES",
      aliquota_pct: "4",
      limiar_verde_pct: "18",
      limiar_amarelo_pct: "10",
      tarifa_full: null,
      difal_enabled: false,
      difal_destino_uf: null,
      origem: "operator",
    };
    const difalList: PricingDifalListResponse = {
      items: [
        {
          uf: "SP",
          interna_pct: "18.00",
          interestadual_pct: "12",
          efetivo_pct: "6.00",
          origem_versao: "padrao-2026",
          override: null,
        },
      ],
      disclaimer: "seed padrão 2026 — não é orientação fiscal",
    };
    const overrideResp: PricingDifalOverrideResponse = {
      persisted: true,
      rate: { ...difalList.items[0], efetivo_pct: "7.00", override: { interna_pct: "19.00", updated_at: "2026-07-18T00:00:00Z", actor: "op" } },
    };
    const scenario: PricingScenario = { id: "sc1", name: "Cenário", payload: { preco: "100.00" }, created_at: "2026-07-18T00:00:00Z" };
    const decompose: PricingDecomposeResponse = {
      decomposition: {
        preco: "100.00",
        comissao: "12.00",
        taxa_fixa: "0",
        frete: "15.00",
        imposto: "4.00",
        difal: null,
        tarifa_full: null,
        custo: null,
        margem_valor: null,
        margem_pct: null,
        componentes_desconhecidos: ["custo_erp"],
      },
      blocking_state: "SEM_CUSTO",
    };
    const solve: PricingSolveResponse = {
      reached: false,
      preco: null,
      ceiling_pct: "80.00",
      desconhecidos: [],
      blocking_state: null,
      code: "UNREACHABLE_TARGET",
    };

    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async (input, init) => {
        requests.push({ input, init });
        const path = String(input);
        const method = init?.method ?? "GET";
        let body: unknown = {};
        if (path.endsWith("/pricing/profile")) body = profile;
        else if (path.endsWith("/pricing/difal")) body = difalList;
        else if (path.includes("/pricing/difal/")) body = overrideResp;
        else if (path.endsWith("/pricing/scenarios")) body = method === "POST" ? scenario : { items: [scenario] };
        else if (path.endsWith("/pricing/decompose")) body = decompose;
        else if (path.endsWith("/pricing/solve")) body = solve;
        if (method === "DELETE") return new Response(null, { status: 204 });
        return new Response(JSON.stringify(body), { status: method === "POST" && path.endsWith("/scenarios") ? 201 : 200 });
      },
    });

    const gotProfile = await client.getPricingProfile();
    const putProfile = await client.putPricingProfile(profile);
    const gotDifal = await client.listPricingDifal();
    const putDifal = await client.putPricingDifalOverride("SP", { interna_pct: "19.00", actor: "op" });
    const listScenarios = await client.listPricingScenarios();
    const createScenario = await client.createPricingScenario({ id: "sc1", name: "Cenário", payload: { preco: "100.00" } });
    await client.deletePricingScenario("sc1");
    const gotDecompose = await client.pricingDecompose({ preco: "100.00", comissao_pct: "12", modalidade: "classico", product_id: 90001 });
    const gotSolve = await client.pricingSolveTarget({ margem_alvo_pct: "95", comissao_pct: "16", modalidade: "classico", custo: "10.00", frete_produto: "15.00" });

    const url = (i: number) => String(requests[i].input);
    const verb = (i: number) => requests[i].init?.method ?? "GET";
    expect(url(0)).toBe("http://localhost:8080/pricing/profile");
    expect(verb(0)).toBe("GET");
    expect(verb(1)).toBe("PUT");
    expect(url(2)).toBe("http://localhost:8080/pricing/difal");
    expect(url(3)).toBe("http://localhost:8080/pricing/difal/SP");
    expect(verb(3)).toBe("PUT");
    expect(url(4)).toBe("http://localhost:8080/pricing/scenarios");
    expect(verb(5)).toBe("POST");
    expect(url(6)).toBe("http://localhost:8080/pricing/scenarios/sc1");
    expect(verb(6)).toBe("DELETE");
    expect(url(7)).toBe("http://localhost:8080/pricing/decompose");
    expect(url(8)).toBe("http://localhost:8080/pricing/solve");

    expect(gotProfile.origem).toBe("operator");
    expect(putProfile.aliquota_pct).toBe("4");
    expect(gotDifal.disclaimer).toBe("seed padrão 2026 — não é orientação fiscal");
    expect(gotDifal.items[0].override).toBeNull();
    expect(putDifal.persisted).toBe(true);
    expect(putDifal.rate.efetivo_pct).toBe("7.00");
    expect(listScenarios.items).toHaveLength(1);
    expect(createScenario.id).toBe("sc1");
    // ADR-17: unknown cost stays null (never fabricated 0) and surfaces as SEM_CUSTO.
    expect(gotDecompose.decomposition.custo).toBeNull();
    // Task A7: imposto is nullable now but still round-trips its legacy value
    // unchanged when the wire sends one (negative control for the null case below).
    expect(gotDecompose.decomposition.imposto).toBe("4.00");
    expect(gotDecompose.blocking_state).toBe("SEM_CUSTO");
    // Unreachable target is a 200 with an explicit code, not a thrown error.
    expect(gotSolve.reached).toBe(false);
    expect(gotSolve.code).toBe("UNREACHABLE_TARGET");
  });

  // Task A7: PricingDecomposition.imposto widened from string to string | null
  // (OpenAPI + SDK, same commit as the backend DTO/handler change). null means
  // the D-41 per-cell fiscal path is active for this item — the SDK must pass
  // it through untouched, never coercing to "0.00"/"" (ADR-17).
  it("passes a null pricing decompose imposto through untouched (nullable field parity)", async () => {
    const decompose: PricingDecomposeResponse = {
      decomposition: {
        preco: "100.00",
        comissao: "12.00",
        taxa_fixa: "0",
        frete: "15.00",
        imposto: null,
        difal: "8.00",
        tarifa_full: null,
        custo: "40.00",
        margem_valor: null,
        margem_pct: null,
        componentes_desconhecidos: [],
      },
      blocking_state: null,
    };
    const client = createMarketplaceCentralClient({
      baseUrl: "http://localhost:8080",
      fetchImpl: async () => new Response(JSON.stringify(decompose), { status: 200 }),
    });
    const got = await client.pricingDecompose({ preco: "100.00", comissao_pct: "12", modalidade: "classico" });
    expect(got.decomposition.imposto).toBeNull();
  });
});

