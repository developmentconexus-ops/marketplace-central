import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, fireEvent, render, screen } from "@testing-library/react";
import type { MutationItem, MutationProtocol } from "@marketplace-central/sdk-runtime";
import { MemoryRouter, Route, Routes, useLocation } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ProtocoloPage } from "./ProtocoloPage";

const { getMutation, listMutationItems, retryMutationFailures, invalidateAfterMutation } = vi.hoisted(() => ({
  getMutation: vi.fn(),
  listMutationItems: vi.fn(),
  retryMutationFailures: vi.fn(),
  invalidateAfterMutation: vi.fn(),
}));

vi.mock("../../app/ClientContext", () => ({ useClient: () => ({ getMutation, listMutationItems, retryMutationFailures }) }));
vi.mock("@marketplace-central/web-query", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@marketplace-central/web-query")>()),
  invalidateAfterMutation,
}));

const applying: MutationProtocol = {
  protocol_id: "MP-000042", installation_id: "inst_test", type: "listing_pause", state: "applying",
  actor: "operator_supplied_unverified", intent: {}, selection: {}, totals: {}, source_as_of: null,
  retried_from: "MP-000001", created_at: "2026-07-17T12:00:00Z", previewed_at: "2026-07-17T12:00:01Z",
  approved_at: "2026-07-17T12:00:02Z", finished_at: null,
};

const failedItem: MutationItem = {
  seq: 1, item_id: "item-1", listing_id: "MLB-1", idempotency_key: "idem-1",
  before: null, after: { status: "paused" }, state: "failed",
  failure: { code: "provider_unavailable", message_pt: "não usar", message_provider: "raw secreto" }, applied_at: null,
};

function LocationProbe() {
  const location = useLocation();
  return <div data-testid="location">{location.pathname}{location.search}</div>;
}

function mount(initialEntry = "/protocolos/MP-000042", queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })) {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <LocationProbe />
        <Routes><Route path="/protocolos/:protocolId" element={<ProtocoloPage />} /></Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProtocoloPage", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    getMutation.mockReset();
    listMutationItems.mockReset();
    retryMutationFailures.mockReset();
    invalidateAfterMutation.mockReset();
    listMutationItems.mockResolvedValue({ items: [], next_cursor: null, page_size: 50 });
  });

  afterEach(() => vi.useRealTimers());

  it("usa os estados compartilhados ao carregar e ao falhar", async () => {
    getMutation.mockReturnValue(new Promise(() => undefined));
    const loading = mount();
    expect(screen.getByRole("status")).toHaveTextContent("Carregando");
    loading.unmount();

    getMutation.mockRejectedValue(new Error("indisponível"));
    mount();
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    expect(screen.getByRole("alert")).toHaveTextContent("Não foi possível carregar o protocolo.");
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();
  });

  it("retoma polling no acesso direto, para no terminal e invalida uma vez", async () => {
    getMutation.mockResolvedValueOnce(applying).mockResolvedValue({ ...applying, state: "applied", finished_at: "2026-07-17T12:00:03Z" });
    mount();

    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    expect(screen.getByText("applying")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Repetir itens com falha" })).not.toBeInTheDocument();
    await act(async () => { await vi.advanceTimersByTimeAsync(2000); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    expect(screen.getByText("applied")).toBeInTheDocument();
    expect(invalidateAfterMutation).toHaveBeenCalledTimes(1);
    await act(async () => { await vi.advanceTimersByTimeAsync(4000); });
    expect(getMutation).toHaveBeenCalledTimes(2);
  });

  it("mostra cabeçalho, valores desconhecidos, falha e detalhe técnico sob demanda", async () => {
    getMutation.mockResolvedValue({ ...applying, state: "failed_preserved", finished_at: "2026-07-17T12:00:03Z" });
    listMutationItems.mockResolvedValue({ items: [failedItem], next_cursor: null, page_size: 50 });
    mount();
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });

    expect(screen.getByRole("heading", { name: "Protocolo MP-000042" })).toBeInTheDocument();
    expect(screen.getByText("Pausar anúncios")).toBeInTheDocument();
    expect(screen.getByText("operator_supplied_unverified")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "MP-000001" })).toHaveAttribute("href", "/protocolos/MP-000001");
    expect(screen.getByText("—")).toBeInTheDocument();
    expect(screen.getByText("MLB-1")).toBeInTheDocument();
    expect(screen.getByText("provider_unavailable")).toBeInTheDocument();
    expect(screen.getByText(/indisponível/i)).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Repetir itens com falha" })).not.toBeInTheDocument();
    expect(screen.queryByText("raw secreto")).not.toBeInTheDocument();
    fireEvent.click(screen.getByText("▸ técnico"));
    expect(screen.getByText("raw secreto")).toBeInTheDocument();
  });

  it("pagina os itens pelo cursor retornado", async () => {
    getMutation.mockResolvedValue({ ...applying, state: "applied" });
    listMutationItems
      .mockResolvedValueOnce({ items: [{ ...failedItem, failure: null, state: "applied" }], next_cursor: "cursor-2", page_size: 1 })
      .mockResolvedValueOnce({ items: [{ ...failedItem, item_id: "item-2", listing_id: "MLB-2" }], next_cursor: null, page_size: 1 });
    mount();
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    fireEvent.click(screen.getByRole("button", { name: "Próxima página" }));
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    expect(listMutationItems).toHaveBeenLastCalledWith("MP-000042", { limit: 50, cursor: "cursor-2" });
  });

  it("repete falhas uma vez, preserva a instalação e abre o novo protocolo", async () => {
    const failed = {
      ...applying,
      state: "failed_preserved" as const,
      finished_at: "2026-07-17T12:00:03Z",
      totals: { items: 1, previewed: 1, applied: 0, failed: 1, skipped: 0 },
    };
    const retried = { ...failed, protocol_id: "MP-000043", retried_from: "MP-000042", state: "applied" as const };
    let resolveRetry!: (protocol: MutationProtocol) => void;
    retryMutationFailures.mockReturnValue(new Promise<MutationProtocol>((resolve) => { resolveRetry = resolve; }));
    getMutation.mockResolvedValueOnce(failed).mockResolvedValue(retried);
    listMutationItems.mockResolvedValue({ items: [failedItem], next_cursor: null, page_size: 50 });
    mount("/protocolos/MP-000042?installation=inst_test&tab=failed");

    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    const retryButton = screen.getByRole("button", { name: "Repetir itens com falha" });
    fireEvent.click(retryButton);
    fireEvent.click(retryButton);
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    expect(retryButton).toBeDisabled();
    expect(retryMutationFailures).toHaveBeenCalledTimes(1);
    expect(retryMutationFailures).toHaveBeenCalledWith("MP-000042");

    await act(async () => {
      resolveRetry(retried);
      await vi.runOnlyPendingTimersAsync();
    });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });

    expect(screen.getByRole("heading", { name: "Protocolo MP-000043" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "MP-000042" })).toHaveAttribute("href", "/protocolos/MP-000042");
    expect(screen.getByTestId("location")).toHaveTextContent("/protocolos/MP-000043?installation=inst_test&tab=failed");
  });

  it("mostra retry quando a falha está fora da página atual", async () => {
    const failed = {
      ...applying,
      state: "failed_preserved" as const,
      finished_at: "2026-07-17T12:00:03Z",
      totals: { items: 2, previewed: 2, applied: 1, failed: 1, skipped: 0 },
    };
    getMutation.mockResolvedValue(failed);
    listMutationItems.mockResolvedValue({ items: [], next_cursor: null, page_size: 50 });
    mount();

    await act(async () => { await vi.runOnlyPendingTimersAsync(); });

    expect(screen.getByRole("button", { name: "Repetir itens com falha" })).toBeInTheDocument();
  });
});
