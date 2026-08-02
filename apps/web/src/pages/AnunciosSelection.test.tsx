import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { AnunciosPage } from "./AnunciosPage";

const listListings = vi.fn();
const getListing = vi.fn();
const createMutation = vi.fn();
const previewMutation = vi.fn();
const getListingsSummary = vi.fn(() =>
  Promise.resolve({
    total: 3,
    active: 3,
    paused: 0,
    exceptions: { below_margin_worst_case: null, margin_unknown: null },
  }),
);
let mockInstallationId = "inst_1";

const baseListing = {
  installation_id: "inst_1",
  provider: "mercado_livre",
  listing_type: { code: "gold_special", label: "Clássico" },
  status: "active" as const,
  link: { state: "resolved" as const, product_id: "product_1", seller_sku: "SKU-1" },
  price: { amount: "129.90", currency: "BRL" },
  published_quantity: 7,
  sync_state: "synced" as const,
  sync_error: null,
  pending_issue: null,
  cost: { amount: "70.00", currency: "BRL" },
  below_margin_worst_case: false,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-16T12:00:00Z",
};

const pageOne = {
  items: [
    { ...baseListing, listing_id: "inst_1~MLB1~-", provider_listing_id: "MLB1", title: "Página 1 A" },
    { ...baseListing, listing_id: "inst_1~MLB2~-", provider_listing_id: "MLB2", title: "Página 1 B" },
  ],
  next_cursor: "c2",
  page_size: 2,
  as_of: "2026-07-16T12:00:00Z",
};

const pageTwo = {
  items: [
    { ...baseListing, listing_id: "inst_1~MLB3~-", provider_listing_id: "MLB3", title: "Página 2 A" },
  ],
  next_cursor: null,
  page_size: 1,
  as_of: "2026-07-16T12:00:00Z",
};

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({
    listListings: (...args: unknown[]) => listListings(...args),
    getListing: (...args: unknown[]) => getListing(...args),
    createMutation: (...args: unknown[]) => createMutation(...args),
    previewMutation: (...args: unknown[]) => previewMutation(...args),
    getListingsSummary,
  }),
}));

vi.mock("../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: mockInstallationId }),
}));

function pageElement() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return (
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={["/anuncios?installation=inst_1"]}>
        <AnunciosPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

function renderPage() {
  return render(pageElement());
}

function fillIntent(type: string) {
  if (type === "price_update") {
    fireEvent.change(screen.getByLabelText("Novo preço"), { target: { value: "49.90" } });
  } else if (type === "stock_correct") {
    fireEvent.change(screen.getByLabelText("Quantidade a publicar"), { target: { value: "7" } });
  } else if (type === "listing_edit") {
    fireEvent.change(screen.getByLabelText("ID do atributo"), { target: { value: "COLOR" } });
    fireEvent.change(screen.getByLabelText("Valor do atributo"), { target: { value: "Azul" } });
  }
}

beforeEach(() => {
  mockInstallationId = "inst_1";
  vi.clearAllMocks();
  listListings.mockImplementation((options: { cursor?: string }) =>
    Promise.resolve(options.cursor === "c2" ? pageTwo : pageOne),
  );
  getListing.mockImplementation((listingId: string) =>
    Promise.resolve({
      ...baseListing,
      listing_id: listingId,
      provider_listing_id: "MLB1",
      title: "Página 1 A",
      timeline: [],
    }),
  );
  createMutation.mockResolvedValue({
    protocol_id: "MP-000042",
    installation_id: "inst_1",
    type: "listing_pause",
    state: "draft",
    actor: "operator_supplied_unverified",
    intent: {},
    selection: {},
    totals: {},
    source_as_of: null,
    retried_from: null,
    created_at: "2026-07-17T12:00:00Z",
    previewed_at: null,
    approved_at: null,
    finished_at: null,
  });
  previewMutation.mockResolvedValue({
    protocol_id: "MP-000042",
    installation_id: "inst_1",
    type: "listing_pause",
    state: "previewed",
    actor: "operator_supplied_unverified",
    intent: {},
    selection: {},
    totals: { items: 1, previewed: 1, failed: 0 },
    source_as_of: null,
    retried_from: null,
    created_at: "2026-07-17T12:00:00Z",
    previewed_at: "2026-07-17T12:00:01Z",
    approved_at: null,
    finished_at: null,
    items: [],
  });
});

describe("AnunciosPage selection and pagination", () => {
  it("accumulates opaque listing ids across pages and preserves them when going back", async () => {
    renderPage();

    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));
    fireEvent.click(screen.getByLabelText("Selecionar anúncio Página 1 B"));
    fireEvent.click(screen.getByRole("button", { name: "Próxima" }));
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 2 A"));

    expect(screen.getByText("3 selecionados")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Anterior" }));

    expect(await screen.findByText("Página 1 A")).toBeInTheDocument();
    expect(screen.getByLabelText("Selecionar anúncio Página 1 A")).toBeChecked();
    expect(screen.getByLabelText("Selecionar anúncio Página 1 B")).toBeChecked();
    expect(screen.getByText("3 selecionados")).toBeInTheDocument();
  });

  it("selects and deselects only the current page from the header checkbox", async () => {
    renderPage();

    const selectPage = await screen.findByLabelText("Selecionar todos os anúncios desta página");
    fireEvent.click(selectPage);
    fireEvent.click(screen.getByRole("button", { name: "Próxima" }));
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 2 A"));
    fireEvent.click(screen.getByRole("button", { name: "Anterior" }));

    fireEvent.click(await screen.findByLabelText("Selecionar todos os anúncios desta página"));

    expect(screen.getByLabelText("Selecionar anúncio Página 1 A")).not.toBeChecked();
    expect(screen.getByLabelText("Selecionar anúncio Página 1 B")).not.toBeChecked();
    expect(screen.getByText("1 selecionados")).toBeInTheDocument();
  });

  it("clears selection when the installation changes", async () => {
    const view = renderPage();
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));
    expect(screen.getByText("1 selecionados")).toBeInTheDocument();

    mockInstallationId = "inst_2";
    view.rerender(pageElement());

    await waitFor(() => expect(screen.queryByText("1 selecionados")).not.toBeInTheDocument());
  });

  it("marks the header checkbox indeterminate when only part of the page is selected", async () => {
    renderPage();

    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));

    const header = screen.getByLabelText<HTMLInputElement>("Selecionar todos os anúncios desta página");
    expect(header.indeterminate).toBe(true);
    expect(header).not.toBeChecked();

    fireEvent.click(screen.getByLabelText("Selecionar anúncio Página 1 B"));
    expect(header.indeterminate).toBe(false);
    expect(header).toBeChecked();
  });

  it("disables Próxima on the last page and Anterior on the first", async () => {
    renderPage();
    await screen.findByText("Página 1 A");

    expect(screen.getByRole("button", { name: "Anterior" })).toBeDisabled();
    const next = screen.getByRole("button", { name: "Próxima" });
    expect(next).toBeEnabled();

    fireEvent.click(next);
    await screen.findByText("Página 2 A");

    expect(screen.getByRole("button", { name: "Próxima" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Anterior" })).toBeEnabled();
  });

  it("resets the cursor to the first page when a tab changes", async () => {
    renderPage();
    fireEvent.click(await screen.findByRole("button", { name: "Próxima" }));
    await screen.findByText("Página 2 A");

    fireEvent.click(screen.getByRole("tab", { name: "Ativos" }));

    expect(await screen.findByText("Página 1 A")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", status: "active" });
    expect(screen.getByRole("button", { name: "Anterior" })).toBeDisabled();
  });

  it("opens the detail drawer from a row title and closes it when the installation changes", async () => {
    const view = renderPage();

    fireEvent.click(await screen.findByRole("button", { name: "Página 1 A" }));

    expect(await screen.findByRole("button", { name: "Fechar painel" })).toBeInTheDocument();
    expect(getListing).toHaveBeenCalledWith("inst_1~MLB1~-");

    mockInstallationId = "inst_2";
    view.rerender(pageElement());

    await waitFor(() =>
      expect(screen.queryByRole("button", { name: "Fechar painel" })).not.toBeInTheDocument(),
    );
  });

  it("hides the bulk actions bar at zero selection and reveals it once a row is selected", async () => {
    renderPage();
    await screen.findByText("Página 1 A");

    const bulkActions = ["Atualizar preço", "Corrigir estoque", "Pausar", "Ressincronizar", "Vincular", "Editar"];

    // Design: the bulk bar is gated on selection — no actions render at zero.
    for (const name of bulkActions) {
      expect(screen.queryByRole("button", { name })).not.toBeInTheDocument();
    }

    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));

    expect(screen.getByText("1 selecionados")).toBeInTheDocument();
    for (const name of bulkActions) {
      expect(screen.getByRole("button", { name })).toBeEnabled();
    }
  });

  it.each([
    ["Atualizar preço", "price_update", "Atualizar preço"],
    ["Corrigir estoque", "stock_correct", "Corrigir estoque"],
    ["Pausar", "listing_pause", "Pausar anúncios"],
    ["Ressincronizar", "listing_resync", "Ressincronizar anúncios"],
    ["Vincular", "link_apply", "Aplicar vínculo"],
    ["Editar", "listing_edit", "Editar anúncios"],
  ] as const)("opens %s with mutation type %s", async (action, type, heading) => {
    renderPage();
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));
    fireEvent.click(screen.getByRole("button", { name: action }));

    expect(await screen.findByRole("heading", { name: heading })).toBeInTheDocument();
    fillIntent(type);
    fireEvent.click(screen.getByRole("button", { name: "Gerar prévia" }));

    await waitFor(() => expect(createMutation).toHaveBeenCalledOnce());
    expect(createMutation).toHaveBeenCalledWith({
      installation_id: "inst_1",
      type,
      actor: "operator_supplied_unverified",
      intent: expect.any(Object),
      selection: { mode: "explicit", listing_ids: ["inst_1~MLB1~-"] },
    });
  });

  it("passes the accumulated listing ids as a launch snapshot", async () => {
    renderPage();
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));
    fireEvent.click(screen.getByLabelText("Selecionar anúncio Página 1 B"));
    fireEvent.click(screen.getByRole("button", { name: "Próxima" }));
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 2 A"));
    fireEvent.click(screen.getByRole("button", { name: "Ressincronizar" }));
    fireEvent.click(screen.getByRole("button", { name: "Gerar prévia" }));

    await waitFor(() => expect(createMutation).toHaveBeenCalledOnce());
    expect(createMutation).toHaveBeenCalledWith(expect.objectContaining({
      installation_id: "inst_1",
      type: "listing_resync",
      selection: {
        mode: "explicit",
        listing_ids: ["inst_1~MLB1~-", "inst_1~MLB2~-", "inst_1~MLB3~-"],
      },
    }));
  });

  it("clears selection and closes an unsubmitted modal when installation changes", async () => {
    const view = renderPage();
    fireEvent.click(await screen.findByLabelText("Selecionar anúncio Página 1 A"));
    fireEvent.click(screen.getByRole("button", { name: "Editar" }));
    expect(await screen.findByRole("dialog")).toBeInTheDocument();

    mockInstallationId = "inst_2";
    view.rerender(pageElement());

    await waitFor(() => {
      expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
      expect(screen.queryByText("1 selecionados")).not.toBeInTheDocument();
    });
  });
});
