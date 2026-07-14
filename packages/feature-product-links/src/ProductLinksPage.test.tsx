import { MemoryRouter } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { createWebQueryClient, linkageQueryKeys, queryKeyNamespaces } from "@marketplace-central/web-query";
import { ProductLinksPage, type ProductLinksClient } from "./ProductLinksPage";

function makeClient(overrides: Partial<ProductLinksClient> = {}): ProductLinksClient {
  const workflows = [
    {
      identity: {
        installation_id: "inst-1",
        provider_item_id: "MLB-CONFLICT",
      },
      candidates: [
        {
          candidate_id: "cand-conflict",
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_item_id: "MLB-CONFLICT",
          internal_product_id: 11,
          internal_product_name: "Produto Conflito",
          internal_reference_code: "REF-11",
          state: "conflict" as const,
          match_input: "title" as const,
          created_at: "2026-07-08T12:00:00Z",
          updated_at: "2026-07-08T12:00:00Z",
        },
      ],
      audit: [],
    },
    {
      identity: {
        installation_id: "inst-1",
        provider_item_id: "MLB-UNRESOLVED",
      },
      candidates: [
        {
          candidate_id: "cand-unresolved",
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_item_id: "MLB-UNRESOLVED",
          state: "unresolved" as const,
          match_input: "none" as const,
          created_at: "2026-07-08T12:00:00Z",
          updated_at: "2026-07-08T12:00:00Z",
        },
      ],
      audit: [],
    },
    {
      identity: {
        installation_id: "inst-1",
        provider_item_id: "MLB-RESOLVED",
      },
      current_link: {
        installation_id: "inst-1",
        provider_code: "mercado_livre",
        provider_item_id: "MLB-RESOLVED",
        state: "resolved" as const,
        internal_product_id: 22,
        internal_product_name: "Produto Resolvido",
        internal_reference_code: "REF-22",
        created_at: "2026-07-08T12:00:00Z",
        updated_at: "2026-07-08T12:05:00Z",
      },
      candidates: [],
      audit: [
        {
          audit_id: "audit-1",
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_item_id: "MLB-RESOLVED",
          action: "approve_candidate" as const,
          actor: { actor_type: "operator", actor_id: "lea", actor_name: "Lea" },
          previous_state: "unresolved" as const,
          next_state: "resolved" as const,
          created_at: "2026-07-08T12:05:00Z",
        },
      ],
    },
    {
      identity: {
        installation_id: "inst-1",
        provider_item_id: "MLB-REJECTED",
      },
      current_link: {
        installation_id: "inst-1",
        provider_code: "mercado_livre",
        provider_item_id: "MLB-REJECTED",
        state: "rejected" as const,
        created_at: "2026-07-08T12:00:00Z",
        updated_at: "2026-07-08T12:05:00Z",
      },
      candidates: [],
      audit: [],
    },
  ];

  return {
    listIntegrationInstallations: vi.fn().mockResolvedValue({
      items: [
        {
          installation_id: "inst-1",
          tenant_id: "tenant_default",
          provider_code: "mercado_livre",
          family: "marketplace" as const,
          display_name: "Mercado Livre Primary",
          status: "connected" as const,
          health_status: "healthy" as const,
          external_account_id: "691607102",
          external_account_name: "Metal Nobre",
          connection: {
            state: "connected" as const,
            health: "healthy" as const,
            provider_code: "mercado_livre",
            external_account_id: "691607102",
            external_account_name: "Metal Nobre",
            auth_strategy: "oauth2" as const,
            next_action: "none" as const,
          },
          runtime_capabilities: [],
          created_at: "2026-07-08T12:00:00Z",
          updated_at: "2026-07-08T12:00:00Z",
        },
      ],
    }),
    listProductLinkWorkflows: vi.fn().mockResolvedValue({ items: workflows }),
    approveProductLinkCandidate: vi.fn().mockResolvedValue({}),
    rejectProductLinkListing: vi.fn().mockResolvedValue({}),
    manualResolveProductLink: vi.fn().mockResolvedValue({}),
    ...overrides,
  };
}

// Production defaults retry once; an error-state assertion would have to wait out
// the retry delay, so tests keep those defaults but disable the retry.
function createTestQueryClient() {
  const queryClient = createWebQueryClient();
  const defaults = queryClient.getDefaultOptions();
  queryClient.setDefaultOptions({ ...defaults, queries: { ...defaults.queries, retry: false } });
  return queryClient;
}

function renderPage(clientOverrides: Partial<ProductLinksClient> = {}, queryClient = createWebQueryClient()) {
  const client = makeClient(clientOverrides);
  render(
    <MemoryRouter initialEntries={["/product-links?installation=inst-1"]}>
      <QueryClientProvider client={queryClient}>
        <ProductLinksPage client={client} />
      </QueryClientProvider>
    </MemoryRouter>,
  );
  return { client, queryClient };
}

describe("ProductLinksPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders workflow states after loading", async () => {
    renderPage();
    expect(screen.getByText(/loading product link workflows/i)).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText("MLB-CONFLICT")).toBeInTheDocument());
    expect(screen.getAllByText("Conflict").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Unresolved").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Resolved").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Rejected").length).toBeGreaterThan(0);
  });

  it("renders empty state when there are no installations", async () => {
    renderPage({
      listIntegrationInstallations: vi.fn().mockResolvedValue({ items: [] }),
      listProductLinkWorkflows: vi.fn().mockResolvedValue({ items: [] }),
    });
    await waitFor(() => expect(screen.getByText(/no installations available/i)).toBeInTheDocument());
  });

  it("renders error state when workflow loading fails", async () => {
    renderPage(
      {
        listProductLinkWorkflows: vi.fn().mockRejectedValue({ error: { message: "workflow exploded" } }),
      },
      createTestQueryClient(),
    );
    await waitFor(() => expect(screen.getByText("workflow exploded")).toBeInTheDocument());
  });

  it("approves a candidate and refreshes workflows", async () => {
    const { client, queryClient } = renderPage();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    await waitFor(() => expect(screen.getByText("MLB-CONFLICT")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: /approve ref-11/i }));
    await waitFor(() =>
      expect(client.approveProductLinkCandidate).toHaveBeenCalledWith(
        expect.objectContaining({
          candidate_id: "cand-conflict",
          actor: expect.objectContaining({ actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" }),
        }),
      ),
    );
    expect(client.listProductLinkWorkflows).toHaveBeenCalledTimes(2);
    expect(invalidateQueries).toHaveBeenCalledTimes(2);
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.linkage });
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.catalog });
  });

  it("rejects a listing identity", async () => {
    const { client, queryClient } = renderPage();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    await waitFor(() => expect(screen.getByText("MLB-CONFLICT")).toBeInTheDocument());
    fireEvent.click(screen.getAllByRole("button", { name: /reject listing/i })[0]);
    await waitFor(() =>
      expect(client.rejectProductLinkListing).toHaveBeenCalledWith(
        expect.objectContaining({
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_item_id: "MLB-CONFLICT",
        }),
      ),
    );
    expect(invalidateQueries).toHaveBeenCalledTimes(2);
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.linkage });
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.catalog });
  });

  it("manually resolves a workflow", async () => {
    const { client, queryClient } = renderPage();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    await waitFor(() => expect(screen.getByText("MLB-CONFLICT")).toBeInTheDocument());
    fireEvent.click(screen.getAllByRole("button", { name: /manual resolve/i })[0]);
    fireEvent.change(screen.getByLabelText("Internal product ID MLB-CONFLICT"), { target: { value: "77" } });
    fireEvent.change(screen.getByLabelText("Internal product name MLB-CONFLICT"), { target: { value: "Produto Manual" } });
    fireEvent.change(screen.getByLabelText("Reference code MLB-CONFLICT"), { target: { value: "REF-77" } });
    fireEvent.change(screen.getByLabelText("Reason MLB-CONFLICT"), { target: { value: "Validated by operator" } });
    fireEvent.click(screen.getByRole("button", { name: /save manual resolution/i }));
    await waitFor(() =>
      expect(client.manualResolveProductLink).toHaveBeenCalledWith(
        expect.objectContaining({
          installation_id: "inst-1",
          provider_item_id: "MLB-CONFLICT",
          internal_product_id: 77,
          internal_product_name: "Produto Manual",
          internal_reference_code: "REF-77",
          reason: "Validated by operator",
        }),
      ),
    );
    expect(invalidateQueries).toHaveBeenCalledTimes(2);
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.linkage });
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.catalog });
  });

  it("registers linkage as uncached and does not invalidate after a failed write", async () => {
    const { queryClient } = renderPage({
      approveProductLinkCandidate: vi.fn().mockRejectedValue(new Error("approval failed")),
    });
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");

    await waitFor(() => expect(screen.getByText("MLB-CONFLICT")).toBeInTheDocument());
    const cachedQuery = queryClient.getQueryCache().find({ queryKey: linkageQueryKeys.workflows("inst-1") });
    expect(cachedQuery?.options.staleTime).toBe(0);
    expect(cachedQuery?.options.gcTime).toBe(0);
    fireEvent.click(screen.getByRole("button", { name: /approve ref-11/i }));

    await waitFor(() => expect(screen.getByText("approval failed")).toBeInTheDocument());
    expect(invalidateQueries).not.toHaveBeenCalled();
  });
});
