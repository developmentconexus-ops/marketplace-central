import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MarketplaceCentralClientError } from "@marketplace-central/sdk-runtime";
import { queryKeyNamespaces } from "@marketplace-central/web-query";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ListingsRefreshControl } from "./ListingsRefreshControl";

const refreshListings = vi.fn();
const listIntegrationOperationRuns = vi.fn();

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({ refreshListings, listIntegrationOperationRuns }),
}));

function run(operation_run_id: string, status: "queued" | "running" | "succeeded" | "failed" | "cancelled") {
  return { operation_run_id, status, operation_type: "listings_refresh" };
}

/** The control polls the account's runs on mount; an account with no run answers empty. */
const noRuns = { items: [] };

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
      .mockResolvedValueOnce(noRuns)
      .mockResolvedValueOnce({ items: [run("run_1", "queued")] })
      .mockResolvedValueOnce({ items: [run("run_1", "running")] });
    renderControl();

    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText("na fila")).toBeInTheDocument();
    expect(refreshListings).toHaveBeenCalledWith({ installation_id: "inst_1" });
    expect(listIntegrationOperationRuns).toHaveBeenCalledWith("inst_1");
    await advancePolling();
    expect(await screen.findByText(/em andamento/)).toBeInTheDocument();
  });

  it("attaches to a refresh_in_progress run without rendering an error", async () => {
    refreshListings.mockRejectedValue(
      new MarketplaceCentralClientError(409, "refresh_in_progress", "refresh já em andamento", {
        operation_run_id: "run_active",
      }),
    );
    listIntegrationOperationRuns
      .mockResolvedValueOnce(noRuns)
      .mockResolvedValue({ items: [run("run_active", "running")] });
    renderControl();

    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));

    expect(await screen.findByText(/em andamento/)).toBeInTheDocument();
    expect(screen.queryByText("Falha ao iniciar atualização.")).not.toBeInTheDocument();
    expect(listIntegrationOperationRuns).toHaveBeenCalledWith("inst_1");
  });

  it("stops polling after the observed run succeeds", async () => {
    refreshListings.mockResolvedValue({ operation_run_id: "run_1" });
    listIntegrationOperationRuns
      .mockResolvedValueOnce(noRuns)
      .mockResolvedValueOnce({ items: [run("run_1", "running")] })
      .mockResolvedValueOnce({ items: [run("run_1", "succeeded")] });
    renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    expect(await screen.findByText(/em andamento/)).toBeInTheDocument();

    await advancePolling();
    expect(await screen.findByText("concluído")).toBeInTheDocument();
    const callsAtSuccess = listIntegrationOperationRuns.mock.calls.length;
    await advancePolling();
    await advancePolling();
    expect(listIntegrationOperationRuns).toHaveBeenCalledTimes(callsAtSuccess);
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
    // Runs belong to an account, so the fake answers per installation: inst_2 has none.
    listIntegrationOperationRuns.mockImplementation((id: string) =>
      Promise.resolve(id === "inst_1" ? { items: [run("run_1", "running")] } : noRuns),
    );
    const { rerender, queryClient } = renderControl();
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    expect(await screen.findByText(/em andamento/)).toBeInTheDocument();

    rerender(
      <QueryClientProvider client={queryClient}>
        <ListingsRefreshControl installationId="inst_2" />
      </QueryClientProvider>,
    );
    await waitFor(() => expect(screen.queryByText(/em andamento/)).not.toBeInTheDocument());
    const callsBefore = listIntegrationOperationRuns.mock.calls.length;
    await advancePolling();
    expect(listIntegrationOperationRuns).toHaveBeenCalledTimes(callsBefore);
  });

  it("shows a start error and re-enables retry for non-409 failures", async () => {
    refreshListings.mockRejectedValue({ status: 500, error: { code: "internal" } });
    listIntegrationOperationRuns.mockResolvedValue(noRuns);
    renderControl();
    const button = screen.getByRole("button", { name: "Atualizar" });
    fireEvent.click(button);

    expect(await screen.findByText("Falha ao iniciar atualização.")).toBeInTheDocument();
    expect(button).toBeEnabled();
    // No run was created, so nothing is tracked and no status pill is claimed.
    expect(screen.queryByText(/em andamento|na fila/)).not.toBeInTheDocument();
  });

  it("adopts a refresh already running on the account when the screen mounts", async () => {
    listIntegrationOperationRuns.mockResolvedValue({ items: [run("run_elsewhere", "running")] });
    renderControl();

    // The operator reloaded mid-refresh: the pill must show the live run and the
    // button must stay disabled instead of inviting a click that 409s.
    expect(await screen.findByText(/em andamento/)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Atualizar" })).toBeDisabled();
    expect(refreshListings).not.toHaveBeenCalled();
  });

  it("reports how long the running refresh has been going", async () => {
    listIntegrationOperationRuns.mockResolvedValue({
      items: [
        {
          ...run("run_elsewhere", "running"),
          started_at: new Date(Date.now() - 185_000).toISOString(),
        },
      ],
    });
    renderControl();

    expect(await screen.findByText("em andamento há 3 min")).toBeInTheDocument();
  });
});
