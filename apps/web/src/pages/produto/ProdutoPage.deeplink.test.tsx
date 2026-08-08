import type { JSX } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes, useLocation } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { InstallationProvider } from "../../app/InstallationContext";
import { ProdutoPage } from "./ProdutoPage";

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getCatalogProduct: () => new Promise(() => undefined),
    listIntegrationInstallations: () => new Promise(() => undefined),
  }),
}));

function LocationProbe(): JSX.Element {
  const location = useLocation();
  return <div data-testid="location-probe" data-search={location.search} />;
}

function renderAt(path: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <InstallationProvider>
          <LocationProbe />
          <Routes>
            <Route path="/catalogo/produtos/:productId" element={<ProdutoPage />} />
          </Routes>
        </InstallationProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProdutoPage deep-link + F5 restore", () => {
  it("activates the Anúncios panel from ?tab=anuncios", async () => {
    renderAt("/catalogo/produtos/90008?tab=anuncios");

    expect(await screen.findByRole("tab", { name: "Anúncios vinculados" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute(
      "aria-selected",
      "false",
    );
    expect(screen.getByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "false");
  });

  it("keeps Anúncios active after a remount at the same URL (F5-restore analogue)", async () => {
    const { unmount } = renderAt("/catalogo/produtos/90008?tab=anuncios");
    expect(await screen.findByRole("tab", { name: "Anúncios vinculados" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    unmount();

    renderAt("/catalogo/produtos/90008?tab=anuncios");
    expect(await screen.findByRole("tab", { name: "Anúncios vinculados" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
  });

  it("falls back to Veredicto for a bogus ?tab= value without crashing", async () => {
    renderAt("/catalogo/produtos/90008?tab=bogus");

    expect(await screen.findByRole("tab", { name: "Veredicto" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    expect(screen.getByRole("tab", { name: "Anúncios vinculados" })).toHaveAttribute(
      "aria-selected",
      "false",
    );
    expect(screen.getByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "false");
  });

  it("updates the location search to tab=estoque when the Estoque tab is clicked", async () => {
    renderAt("/catalogo/produtos/90008");

    fireEvent.click(await screen.findByRole("tab", { name: "Estoque" }));

    expect(screen.getByTestId("location-probe")).toHaveAttribute("data-search", "?tab=estoque");
  });
});
