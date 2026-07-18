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

  it("mounts the header navigation with the enabled pills as links", async () => {
    renderLayout();
    await screen.findByRole("combobox", { name: "Selecionar instalação" });

    // The retheme moved the primary nav out of Layout into <Header/>. Layout's job is to mount
    // that header; the canonical six-pill contract (order + the disabled "em breve" stubs) is
    // owned and asserted by Header.test.tsx. Here we only assert the integration wiring: the
    // primary <nav> renders the four enabled pills as links. (Native `toContain` — avoids adding
    // a jest-dom matcher, which would trip the known F-ENV-4 TS2339 gap and inflate the count.)
    const navLinkNames = within(screen.getByRole("navigation"))
      .getAllByRole("link")
      .map((link) => link.textContent);

    for (const label of ["Visão geral", "Anúncios", "Simulador", "Pedidos"]) {
      expect(navLinkNames).toContain(label);
    }
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

  it("preserves the query string when following an enabled navigation pill", async () => {
    renderLayout("/anuncios?installation=inst_test&tab=pendencia");
    await screen.findByRole("combobox", { name: "Selecionar instalação" });

    // Enabled pills carry the current search (installation + unrelated params) via
    // `to={{ pathname, search: location.search }}` in Header, so no query is dropped on nav.
    // (Catálogo is intentionally a gear-menu link now, not a pill, per the M-03 nav contract.)
    fireEvent.click(screen.getByRole("link", { name: "Simulador" }));

    expect(screen.getByTestId("location")).toHaveTextContent(
      "/precos?installation=inst_test&tab=pendencia",
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
