import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { queryKeyNamespaces } from "@marketplace-central/web-query";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ListingsRefreshControl } from "./ListingsRefreshControl";

const refreshListings = vi.fn();
const listIntegrationOperationRuns = vi.fn();

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({ refreshListings, listIntegrationOperationRuns }),
}));

function run(operation_run_id: string, status: "queued" | "running" | "succeeded" | "failed" | "cancelled") {
  return { operation_run_id, status };
}

function renderControl(installationId = "inst_1") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  const invalidateSpy = vi.spyOn(queryClient, "invalidateQueries");
  const view = render(
    <QueryClientProvider client={queryClient}>
      <ListingsRefreshControl installationId={installationId} />
    </QueryClientProvider>,
  );
  return { ...view, queryClient, invalidateSpy };
}

async function advancePolling() {
  await act(async () => {
    await vi.advanceTimersByTimeAsync(2_000);
  });
}

describe("ListingsRefreshControl", () => {
  beforeEach(() => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    refreshListings.mockReset();
    listIntegrationOperationRuns.mockReset();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("observes a 202 run and polls while it is queued or running", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns
      .mockResolvedValueOnce({ items: [run("run_1", "queued")] })
      .mockResolvedValueOnce({ items: [run("run_1", "running")] });
    renderControl();

    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText("na fila")).toBeInTheDocument();
    expect(refreshListings).toHaveBeenCalledWith({ installation_id: "inst_1" });
    expect(listIntegrationOperationRuns).toHaveBeenCalledWith("inst_1");
    await advancePolling();
    expect(await screen.findByText("em andamento")).toBeInTheDocument();
    expect(listIntegrationOperationRuns).toHaveBeenCalledTimes(2);
  });

  it("attaches to a refresh_in_progress run without rendering an error", async () => {
    refreshListings.mockRejectedValue({
      status: 409,
      error: { code: "refresh_in_progress", details: { operation_run_id: "run_active" } },
    });
    listIntegrationOperationRuns.mockResolvedValue({ items: [run("run_active", "running")] });
    renderControl();

    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText("em andamento")).toBeInTheDocument();
    expect(screen.queryByText("Falha ao iniciar atualização.")).not.toBeInTheDocument();
    expect(listIntegrationOperationRuns).toHaveBeenCalledWith("inst_1");
  });

  it("stops polling after the observed run succeeds", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns
      .mockResolvedValueOnce({ items: [run("run_1", "running")] })
      .mockResolvedValueOnce({ items: [run("run_1", "succeeded")] });
    renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    expect(await screen.findByText("em andamento")).toBeInTheDocument();

    await advancePolling();
    expect(await screen.findByText("concluído")).toBeInTheDocument();
    await advancePolling();
    await advancePolling();
    expect(listIntegrationOperationRuns).toHaveBeenCalledTimes(2);
  });

  it("invalidates the listings namespace exactly once for a succeeded run", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns.mockResolvedValue({ items: [run("run_1", "succeeded")] });
    const { invalidateSpy, rerender } = renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText("concluído")).toBeInTheDocument();
    await waitFor(() => expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.listings }));
    rerender(
      <QueryClientProvider client={invalidateSpy.mock.instances[0] as QueryClient}>
        <ListingsRefreshControl installationId="inst_1" />
      </QueryClientProvider>,
    );
    expect(invalidateSpy).toHaveBeenCalledTimes(1);
  });

  it("renders an honest failed state without invalidating listings", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns.mockResolvedValue({ items: [run("run_1", "failed")] });
    const { invalidateSpy } = renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText("Atualização falhou.")).toBeInTheDocument();
    expect(invalidateSpy).not.toHaveBeenCalled();
  });

  it("abandons the old run when the installation changes", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns.mockResolvedValue({ items: [run("run_1", "running")] });
    const { rerender, queryClient } = renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    expect(await screen.findByText("em andamento")).toBeInTheDocument();

    rerender(
      <QueryClientProvider client={queryClient}>
        <ListingsRefreshControl installationId="inst_2" />
      </QueryClientProvider>,
    );
    await waitFor(() => expect(screen.queryByText("em andamento")).not.toBeInTheDocument());
    const callsBefore = listIntegrationOperationRuns.mock.calls.length;
    await advancePolling();
    expect(listIntegrationOperationRuns).toHaveBeenCalledTimes(callsBefore);
  });

  it("shows a start error and re-enables retry for non-409 failures", async () => {
    refreshListings.mockRejectedValue({ status: 500, error: { code: "internal" } });
    renderControl();
    const button = screen.getByRole("button", { name: "Atualizar" });
    fireEvent.click(button);

    expect(await screen.findByText("Falha ao iniciar atualização.")).toBeInTheDocument();
    expect(button).toBeEnabled();
    expect(listIntegrationOperationRuns).not.toHaveBeenCalled();
  });
});
