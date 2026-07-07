import { describe, expect, it, vi } from "vitest";
import { createMarketplaceCentralClient } from "./index";
import type { IntegrationProviderDefinition } from "./index";

describe("sdk runtime", () => {
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
});

