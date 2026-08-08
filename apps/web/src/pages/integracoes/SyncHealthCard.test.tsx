import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MarketplaceCentralClientError, type SyncHealth } from "@marketplace-central/sdk-runtime";
import { render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { InstallationProvider } from "../../app/InstallationContext";
import { IntegracoesPage } from "./IntegracoesPage";
import { SyncHealthCard } from "./SyncHealthCard";

// jsdom in this workspace does not instantiate Web Storage under the worker's
// Vitest runner (same gotcha as IntegracoesPage.test.tsx — reused verbatim,
// not forked, so any component under IntegracoesPage that touches storage
// does not blow up here too).
if (typeof globalThis.localStorage === "undefined") {
  const stores = [new Map<string, string>(), new Map<string, string>()];
  const define = (name: string, value: unknown) =>
    Object.defineProperty(Storage.prototype, name, { configurable: true, writable: true, value });
  define("getItem", function (this: Storage, key: string) {
    const store = (this as Storage & { __store: Map<string, string> }).__store;
    return store.has(key) ? store.get(key)! : null;
  });
  define("setItem", function (this: Storage, key: string, value: string) {
    (this as Storage & { __store: Map<string, string> }).__store.set(key, String(value));
  });
  define("removeItem", function (this: Storage, key: string) {
    (this as Storage & { __store: Map<string, string> }).__store.delete(key);
  });
  define("clear", function (this: Storage) {
    (this as Storage & { __store: Map<string, string> }).__store.clear();
  });
  define("key", function (this: Storage, index: number) {
    return (
      Array.from((this as Storage & { __store: Map<string, string> }).__store.keys())[index] ?? null
    );
  });
  const storageEntries: Array<["localStorage" | "sessionStorage", Map<string, string>]> = [
    ["localStorage", stores[0]],
    ["sessionStorage", stores[1]],
  ];
  for (const [name, store] of storageEntries) {
    const storage = Object.create(Storage.prototype) as Storage & { __store: Map<string, string> };
    storage.__store = store;
    Object.defineProperty(storage, "length", { configurable: true, get: () => store.size });
    Object.defineProperty(globalThis, name, { configurable: true, value: storage });
    Object.defineProperty(window, name, { configurable: true, value: storage });
  }
}

// Fixed "now" so relative-time assertions ("há 5 min") are deterministic.
// Only Date.now is stubbed (not fake timers), so react-query's internal
// setTimeout-driven retry/refetch and Testing Library's waitFor keep running
// on real timers.
const NOW = new Date("2026-08-01T12:00:00Z").getTime();

// Full client surface mocked (not just getSyncHealth) so the isolation test
// can prove the OTHER cards render their normal success content while only
// this card's fetch fails — a partial mock would make every other card fail
// too, for an unrelated reason, and the isolation assertion would be vacuous.
const createErpImport = vi.fn();
const getErpImport = vi.fn();
const listErpImports = vi.fn();
const getActiveSource = vi.fn();
const setActiveSource = vi.fn();
const getSellableAssortment = vi.fn();
const setSellableAssortment = vi.fn();
const getCatalogAssortmentCounts = vi.fn();
const listIntegrationInstallations = vi.fn();
const getSyncHealth = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    createErpImport: (...args: unknown[]) => createErpImport(...args),
    getErpImport: (...args: unknown[]) => getErpImport(...args),
    listErpImports: (...args: unknown[]) => listErpImports(...args),
    getActiveSource: (...args: unknown[]) => getActiveSource(...args),
    setActiveSource: (...args: unknown[]) => setActiveSource(...args),
    getSellableAssortment: (...args: unknown[]) => getSellableAssortment(...args),
    setSellableAssortment: (...args: unknown[]) => setSellableAssortment(...args),
    getCatalogAssortmentCounts: (...args: unknown[]) => getCatalogAssortmentCounts(...args),
    listIntegrationInstallations: (...args: unknown[]) => listIntegrationInstallations(...args),
    getSyncHealth: (...args: unknown[]) => getSyncHealth(...args),
  }),
}));

function activeSourceConfig(source: string) {
  return {
    active_source: source,
    source_kind: "live_read_through",
    set_at: "2026-07-24T10:00:00Z",
    set_by: null,
  };
}

function renderCard() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <SyncHealthCard />
    </QueryClientProvider>,
  );
}

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <MemoryRouter>
      <QueryClientProvider client={queryClient}>
        {/* IntegracoesPage now mounts ConnectionHealthCard, which reads
            useInstallation() (InstallationContext.tsx). In the real app
            AppRouter.tsx:60 wraps every route in InstallationProvider; this
            test rendered the page standalone before that dependency existed. */}
        <InstallationProvider>
          <IntegracoesPage />
        </InstallationProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

describe("SyncHealthCard", () => {
  beforeEach(() => {
    createErpImport.mockReset();
    getErpImport.mockReset();
    listErpImports.mockReset();
    listErpImports.mockResolvedValue({ items: [] });
    getActiveSource.mockReset();
    getActiveSource.mockResolvedValue(activeSourceConfig("xlsx"));
    setActiveSource.mockReset();
    getSellableAssortment.mockReset();
    getSellableAssortment.mockResolvedValue({
      only_revenda: true,
      only_em_estoque: false,
      only_ecommerce_eligible: false,
    });
    setSellableAssortment.mockReset();
    getCatalogAssortmentCounts.mockReset();
    getCatalogAssortmentCounts.mockResolvedValue({ sellable_count: 2, total_count: 4 });
    listIntegrationInstallations.mockReset();
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    getSyncHealth.mockReset();
    vi.spyOn(Date, "now").mockReturnValue(NOW);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders a green badge with relative time and the absolute ISO in the title for a healthy entity", async () => {
    const health: SyncHealth = {
      entities: [
        {
          entity: "products",
          last_success_at: "2026-08-01T11:55:00Z",
          last_incremental_at: "2026-08-01T11:55:00Z",
          consecutive_failures: 0,
          phase: "incremental",
          last_error: null,
        },
      ],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const row = await screen.findByTestId("sync-health-entity-products");
    expect(within(row).getByTestId("sync-health-badge-products")).toHaveTextContent("ok");
    expect(within(row).getByTestId("sync-health-badge-products").className).toMatch(
      /text-accent-ink/,
    );
    const timeEl = within(row).getByText("há 5 min");
    expect(timeEl).toHaveAttribute("title", "2026-08-01T11:55:00Z");
    expect(within(row).getByText("(incremental)")).toBeInTheDocument();
  });

  it("renders a red badge with the failure count and the last_error message in a tooltip", async () => {
    const health: SyncHealth = {
      entities: [
        {
          entity: "orders",
          last_success_at: "2026-08-01T09:00:00Z",
          last_incremental_at: null,
          consecutive_failures: 3,
          phase: null,
          last_error: { message: "conexão recusada pelo provedor", at: "2026-08-01T11:59:00Z" },
        },
      ],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const row = await screen.findByTestId("sync-health-entity-orders");
    const badge = within(row).getByTestId("sync-health-badge-orders");
    expect(badge).toHaveTextContent("3 falhas");
    expect(badge.className).toMatch(/text-warn/);
    expect(badge).toHaveAttribute("title", "conexão recusada pelo provedor");
    // phase: null on this fixture — the parenthetical phase suffix must not render at all.
    expect(within(row).queryByText(/\(.*\)/)).not.toBeInTheDocument();
  });

  it("shows the literal 'nunca' in faint styling for an entity that never completed a sync, with no relative-time text", async () => {
    const health: SyncHealth = {
      entities: [
        {
          entity: "listings",
          last_success_at: null,
          last_incremental_at: null,
          consecutive_failures: 0,
          phase: null,
          last_error: null,
        },
      ],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const row = await screen.findByTestId("sync-health-entity-listings");
    expect(within(row).getByTestId("sync-health-never-listings")).toHaveTextContent("nunca");
    expect(within(row).getByTestId("sync-health-badge-listings")).toHaveTextContent("nunca");
    expect(within(row).getByTestId("sync-health-badge-listings").className).toMatch(/text-faint/);
    expect(row.textContent).not.toMatch(/\d+\s*min/);
    expect(row.textContent).not.toMatch(/0 min/);
  });

  // F-r04-1 regression at the UI layer: last_success_at is the backend's
  // already-computed GREATEST(full, incremental). The card must trust it
  // as-is and never re-derive freshness itself — proven here by an entity
  // whose last_success_at is recent while last_incremental_at is much older,
  // the inverse of what a naive "just use the incremental column" bug needs.
  it("trusts a recent last_success_at as fresh/green without recomputing staleness itself", async () => {
    const health: SyncHealth = {
      entities: [
        {
          entity: "products",
          last_success_at: "2026-08-01T11:58:00Z",
          last_incremental_at: "2026-07-20T00:00:00Z",
          consecutive_failures: 0,
          phase: "incremental",
          last_error: null,
        },
      ],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const row = await screen.findByTestId("sync-health-entity-products");
    expect(within(row).getByTestId("sync-health-badge-products")).toHaveTextContent("ok");
    expect(within(row).getByText("há 2 min")).toBeInTheDocument();
    expect(row.textContent).not.toMatch(/há \d+ d/);
    expect(within(row).queryByTestId("sync-health-never-products")).not.toBeInTheDocument();
    expect(within(row).getByText("(incremental)")).toBeInTheDocument();
  });

  it("states the fact 'nenhuma notificação recebida' for the webhook's initial state, never a configuration verdict", async () => {
    const health: SyncHealth = {
      entities: [],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const section = await screen.findByTestId("sync-health-webhook-initial");
    expect(section).toHaveTextContent("Nenhuma notificação recebida.");
    expect(screen.queryByText(/não configurado/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/not configured/i)).not.toBeInTheDocument();
  });

  it("shows notification/pending/dropped counts for the webhook's active state", async () => {
    const health: SyncHealth = {
      entities: [],
      webhook: { last_notification_at: "2026-08-01T11:50:00Z", pending: 2, dropped_24h: 1 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const section = await screen.findByTestId("sync-health-webhook-active");
    expect(section).toHaveTextContent("há 10 min");
    expect(section).toHaveTextContent("2");
    expect(section).toHaveTextContent("1");
    expect(within(section).getByText("há 10 min")).toHaveAttribute("title", "2026-08-01T11:50:00Z");
  });

  it("renders a named error and a retry affordance when the fetch fails, without a blank card", async () => {
    getSyncHealth.mockRejectedValue(
      new MarketplaceCentralClientError(500, "internal_error", "falha ao ler sync_state", {}),
    );
    renderCard();

    const alert = await screen.findByRole("alert");
    expect(alert).toHaveTextContent("internal_error");
    expect(alert).toHaveTextContent("falha ao ler sync_state");
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();
  });

  it("isolates the fetch failure to this card: the rest of IntegracoesPage stays intact with its normal content", async () => {
    getSyncHealth.mockRejectedValue(
      new MarketplaceCentralClientError(500, "internal_error", "falha ao ler sync_state", {}),
    );
    renderPage();

    // The other cards' own successful reads still render their real data —
    // proving this is isolation, not every card failing for unrelated reasons.
    expect(await screen.findByTestId("active-source-xlsx")).toBeInTheDocument();
    await waitFor(() => expect(screen.getByText("Resultado: 2 de 4 produtos")).toBeInTheDocument());
    expect(screen.getByRole("heading", { name: "Conectar marketplace" })).toBeInTheDocument();
    expect(screen.getByTestId("provider-connect-ml")).toBeInTheDocument();

    expect(screen.getByRole("heading", { name: "Saúde do sync" })).toBeInTheDocument();
    expect(await screen.findByText(/internal_error/)).toBeInTheDocument();
  });

  it("renders a generic row for an unknown/future entity name not in any hardcoded list", async () => {
    const health: SyncHealth = {
      entities: [
        {
          entity: "warehouse_sync",
          last_success_at: "2026-08-01T11:59:00Z",
          last_incremental_at: null,
          consecutive_failures: 0,
          phase: null,
          last_error: null,
        },
      ],
      webhook: { last_notification_at: null, pending: 0, dropped_24h: 0 },
    };
    getSyncHealth.mockResolvedValue(health);
    renderCard();

    const row = await screen.findByTestId("sync-health-entity-warehouse_sync");
    expect(row).toHaveTextContent("Warehouse Sync");
    expect(within(row).getByTestId("sync-health-badge-warehouse_sync")).toHaveTextContent("ok");
  });
});
