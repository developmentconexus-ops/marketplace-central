import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { useState } from "react";
import { MemoryRouter, Route, Routes, useLocation, useNavigate } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { InstallationGate, InstallationProvider, useInstallation } from "./InstallationContext";

const { listIntegrationInstallations } = vi.hoisted(() => ({
  listIntegrationInstallations: vi.fn(),
}));

vi.mock("./ClientContext", () => ({
  useClient: () => ({ listIntegrationInstallations }),
}));

const installations = [
  { installation_id: "inst_test", display_name: "Conta teste" },
  { installation_id: "inst_second", display_name: "Conta secundária" },
];

function Harness({ routes = false }: { routes?: boolean }) {
  const context = useInstallation();
  const location = useLocation();
  const navigate = useNavigate();
  return (
    <>
      <output data-testid="installation">{context.installationId}</output>
      <output data-testid="status">{context.status}</output>
      <output data-testid="location">{`${location.pathname}${location.search}`}</output>
      <button onClick={() => context.setInstallationId("inst_second")}>Trocar instalação</button>
      {routes ? (
        <>
          <button onClick={() => navigate("/dois" + location.search)}>Rota dois</button>
          <button onClick={() => navigate("/tres" + location.search)}>Rota três</button>
        </>
      ) : null}
    </>
  );
}

function renderContext(initialEntry: string, routes = false) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <InstallationProvider>
          <div data-testid="shell" />
          <InstallationGate>
            <Routes>
              <Route path="*" element={<Harness routes={routes} />} />
            </Routes>
          </InstallationGate>
        </InstallationProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("InstallationContext", () => {
  beforeEach(() => {
    listIntegrationInstallations.mockReset();
    listIntegrationInstallations.mockResolvedValue({ items: installations });
  });

  it("falls back from an unknown installation and rewrites the URL", async () => {
    renderContext("/anuncios?installation=inst_ghost&tab=pendencia");

    expect(await screen.findByTestId("installation")).toHaveTextContent("inst_test");
    await waitFor(() =>
      expect(screen.getByTestId("location")).toHaveTextContent(
        "/anuncios?installation=inst_test&tab=pendencia",
      ),
    );
  });

  it("keeps a known installation without rewriting the URL", async () => {
    renderContext("/anuncios?installation=inst_second&tab=pendencia");

    expect(await screen.findByTestId("installation")).toHaveTextContent("inst_second");
    expect(screen.getByTestId("location")).toHaveTextContent(
      "/anuncios?installation=inst_second&tab=pendencia",
    );
  });

  it("changes selection without losing the route, query params, or refetching", async () => {
    renderContext("/anuncios?installation=inst_test&tab=pendencia&q=camiseta");
    await screen.findByTestId("installation");

    fireEvent.click(screen.getByRole("button", { name: "Trocar instalação" }));

    expect(screen.getByTestId("installation")).toHaveTextContent("inst_second");
    expect(screen.getByTestId("location")).toHaveTextContent(
      "/anuncios?installation=inst_second&tab=pendencia&q=camiseta",
    );
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("uses one request across three routes and a selection change", async () => {
    renderContext("/um?installation=inst_test", true);
    await screen.findByTestId("installation");

    fireEvent.click(screen.getByRole("button", { name: "Rota dois" }));
    fireEvent.click(screen.getByRole("button", { name: "Trocar instalação" }));
    fireEvent.click(screen.getByRole("button", { name: "Rota três" }));

    expect(screen.getByTestId("location")).toHaveTextContent("/tres?installation=inst_second");
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("reports and renders the empty state", async () => {
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    renderContext("/anuncios");

    expect(await screen.findByText("Conecte uma conta em Integrações")).toBeInTheDocument();
    expect(screen.getByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.getByTestId("shell")).toBeInTheDocument();
  });

  it("reports an error and recovers with exactly one retry request", async () => {
    listIntegrationInstallations
      .mockRejectedValueOnce(new Error("indisponível"))
      .mockResolvedValueOnce({ items: installations });
    renderContext("/anuncios");

    expect(await screen.findByRole("alert")).toHaveTextContent("Erro ao carregar.");
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);

    await act(async () => {
      fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));
    });

    expect(await screen.findByTestId("installation")).toHaveTextContent("inst_test");
    expect(screen.getByTestId("status")).toHaveTextContent("ready");
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(2);
  });
});
