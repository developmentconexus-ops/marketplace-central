import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, fireEvent, render, screen } from "@testing-library/react";
import type { MutationItem, MutationProtocol } from "@marketplace-central/sdk-runtime";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ProtocoloPage } from "./ProtocoloPage";

const { getMutation, listMutationItems, invalidateAfterMutation } = vi.hoisted(() => ({
  getMutation: vi.fn(),
  listMutationItems: vi.fn(),
  invalidateAfterMutation: vi.fn(),
}));

vi.mock("../../app/ClientContext", () => ({ useClient: () => ({ getMutation, listMutationItems }) }));
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

function mount(queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })) {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={["/protocolos/MP-000042"]}>
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
    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    expect(screen.getByRole("alert")).toHaveTextContent("Não foi possível carregar o protocolo.");
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();
  });

  it("retoma polling no acesso direto, para no terminal e invalida uma vez", async () => {
    getMutation.mockResolvedValueOnce(applying).mockResolvedValue({ ...applying, state: "applied", finished_at: "2026-07-17T12:00:03Z" });
    mount();

    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    expect(screen.getByText("applying")).toBeInTheDocument();
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
    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });

    expect(screen.getByRole("heading", { name: "Protocolo MP-000042" })).toBeInTheDocument();
    expect(screen.getByText("Pausar anúncios")).toBeInTheDocument();
    expect(screen.getByText("operator_supplied_unverified")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "MP-000001" })).toHaveAttribute("href", "/protocolos/MP-000001");
    expect(screen.getByText("—")).toBeInTheDocument();
    expect(screen.getByText("MLB-1")).toBeInTheDocument();
    expect(screen.getByText(/indisponível/i)).toBeInTheDocument();
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
    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    fireEvent.click(screen.getByRole("button", { name: "Próxima página" }));
    await act(async () => { await vi.runOnlyPendingTimersAsync(); });
    expect(listMutationItems).toHaveBeenLastCalledWith("MP-000042", { limit: 50, cursor: "cursor-2" });
  });
});
