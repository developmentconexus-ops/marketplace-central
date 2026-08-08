import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, renderHook } from "@testing-library/react";
import type { MutationProtocol } from "@marketplace-central/sdk-runtime";
import type { ReactNode } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useMutationProtocol } from "./useMutationProtocol";

const { getMutation, invalidateAfterMutation } = vi.hoisted(() => ({
  getMutation: vi.fn(),
  invalidateAfterMutation: vi.fn(),
}));

vi.mock("../../app/ClientContext", () => ({ useClient: () => ({ getMutation }) }));
vi.mock("@marketplace-central/web-query", async (importOriginal) => ({
  ...(await importOriginal<typeof import("@marketplace-central/web-query")>()),
  invalidateAfterMutation,
}));

const applying: MutationProtocol = {
  protocol_id: "MP-000042",
  installation_id: "inst_test",
  type: "listing_pause",
  state: "applying",
  actor: "operator_supplied_unverified",
  intent: {},
  selection: {},
  totals: {},
  source_as_of: null,
  retried_from: null,
  created_at: "2026-07-17T12:00:00Z",
  previewed_at: "2026-07-17T12:00:01Z",
  approved_at: "2026-07-17T12:00:02Z",
  finished_at: null,
};

describe("useMutationProtocol", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    getMutation.mockReset();
    invalidateAfterMutation.mockReset();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("consulta a cada 2000 ms enquanto não terminal e para no terminal", async () => {
    const terminal = {
      ...applying,
      state: "applied" as const,
      finished_at: "2026-07-17T12:00:03Z",
    };
    getMutation.mockResolvedValueOnce(applying).mockResolvedValue(terminal);
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const wrapper = ({ children }: { children: ReactNode }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
    const { result } = renderHook(() => useMutationProtocol("MP-000042"), { wrapper });

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(result.current.data?.state).toBe("applying");
    await act(async () => {
      await vi.advanceTimersByTimeAsync(1999);
    });
    expect(getMutation).toHaveBeenCalledTimes(1);
    await act(async () => {
      await vi.advanceTimersByTimeAsync(1);
    });
    expect(getMutation).toHaveBeenCalledTimes(2);
    await act(async () => {
      await vi.runAllTimersAsync();
    });
    expect(result.current.data?.state).toBe("applied");
    expect(invalidateAfterMutation).toHaveBeenCalledWith(queryClient, "listing_pause");
    await act(async () => {
      await vi.advanceTimersByTimeAsync(4000);
    });
    expect(getMutation).toHaveBeenCalledTimes(2);
  });

  it("invalida uma vez por protocolo entre consumidores do mesmo QueryClient", async () => {
    const terminal = {
      ...applying,
      state: "failed_preserved" as const,
      finished_at: "2026-07-17T12:00:03Z",
    };
    getMutation.mockResolvedValue(terminal);
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const wrapper = ({ children }: { children: ReactNode }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
    renderHook(() => useMutationProtocol("MP-000042"), { wrapper });
    renderHook(() => useMutationProtocol("MP-000042"), { wrapper });

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    expect(invalidateAfterMutation).toHaveBeenCalledTimes(1);
  });

  it("retoma a consulta pelo retry da Query após uma falha transitória", async () => {
    getMutation.mockRejectedValueOnce(new Error("rede indisponível")).mockResolvedValue(applying);
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: 1, retryDelay: 10 } },
    });
    const wrapper = ({ children }: { children: ReactNode }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
    const { result } = renderHook(() => useMutationProtocol("MP-000042"), { wrapper });

    await act(async () => {
      await vi.advanceTimersByTimeAsync(0);
    });
    await act(async () => {
      await vi.runOnlyPendingTimersAsync();
    });
    expect(result.current.data?.state).toBe("applying");
    expect(getMutation).toHaveBeenCalledTimes(2);
  });
});
