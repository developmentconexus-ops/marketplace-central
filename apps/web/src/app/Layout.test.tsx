import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter, Route, Routes, useLocation } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { InstallationProvider } from "./InstallationContext";
import { Layout } from "./Layout";

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

function LocationOutput() {
  const location = useLocation();
  return <output data-testid="location">{`${location.pathname}${location.search}`}</output>;
}

function renderLayout(initialEntry = "/") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <InstallationProvider>
          <Routes>
            <Route element={<Layout />}>
              <Route path="*" element={<><LocationOutput /><div data-testid="page-content">Página protegida</div></>} />
            </Route>
          </Routes>
        </InstallationProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("Layout", () => {
  beforeEach(() => {
    listIntegrationInstallations.mockReset();
    listIntegrationInstallations.mockResolvedValue({ items: installations });
  });

  it("renders the eight sidebar labels in the contracted order", async () => {
    renderLayout();
    await screen.findByRole("combobox", { name: "Selecionar instalação" });

    const labels = within(screen.getByRole("navigation"))
      .getAllByRole("link")
      .map((link) => link.textContent);

    expect(labels).toEqual([
      "Visão geral",
      "Catálogo",
      "Anúncios",
      "Vínculos & Import.",
      "Estoque",
      "Preços & Simulador",
      "Pedidos",
      "Integrações & Sync",
    ]);
    expect(labels).not.toContain("Mercado");
    expect(labels).not.toContain("Classifications");
    expect(labels).not.toContain("Marketplaces");
  });

  it("shows the selected installation in the ML pill", async () => {
    renderLayout();

    const selector = await screen.findByRole("combobox", { name: "Selecionar instalação" });
    expect(selector).toHaveValue("inst_test");
    expect(screen.getByRole("option", { name: "Conta teste", selected: true })).toBeInTheDocument();
    expect(selector.closest("span")).toHaveTextContent(/^ML:.*▾$/);
  });

  it("changes installation while preserving the pathname and unrelated query params", async () => {
    renderLayout("/anuncios?installation=inst_test&tab=pendencia&q=camiseta");

    const selector = await screen.findByRole("combobox", { name: "Selecionar instalação" });
    fireEvent.change(selector, { target: { value: "inst_second" } });

    await waitFor(() =>
      expect(screen.getByTestId("location")).toHaveTextContent(
        "/anuncios?installation=inst_second&tab=pendencia&q=camiseta",
      ),
    );
  });

  it("preserves the installation query param when following a sidebar link", async () => {
    renderLayout("/anuncios?installation=inst_test&tab=pendencia");
    await screen.findByRole("combobox", { name: "Selecionar instalação" });

    fireEvent.click(screen.getByRole("link", { name: "Catálogo" }));

    expect(screen.getByTestId("location")).toHaveTextContent(
      "/catalogo?installation=inst_test&tab=pendencia",
    );
  });

  it("keeps the shell visible and gates page content when there are no installations", async () => {
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    renderLayout("/anuncios");

    expect(await screen.findAllByText("Conecte uma conta em Integrações")).not.toHaveLength(0);
    expect(screen.getByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.getByRole("navigation")).toBeInTheDocument();
    expect(screen.queryByRole("combobox", { name: "Selecionar instalação" })).not.toBeInTheDocument();
    expect(screen.queryByTestId("page-content")).not.toBeInTheDocument();
  });
});
