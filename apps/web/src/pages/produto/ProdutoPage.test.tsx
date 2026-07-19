import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it } from "vitest";
import { ProdutoPage } from "./ProdutoPage";

function renderAt(path: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/catalogo/produtos/:productId" element={<ProdutoPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProdutoPage", () => {
  it("renders the three tabs and defaults to Veredicto", async () => {
    renderAt("/catalogo/produtos/90008");

    expect(await screen.findByRole("tab", { name: "Veredicto" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Anúncios vinculados" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Estoque" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "true");
  });

  it("restores the active tab from the ?tab= deep link", async () => {
    renderAt("/catalogo/produtos/90008?tab=estoque");

    expect(await screen.findByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "false");
  });

  it("updates the URL tab param when a tab is clicked", async () => {
    renderAt("/catalogo/produtos/90008");

    fireEvent.click(await screen.findByRole("tab", { name: "Estoque" }));

    expect(screen.getByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "false");
  });

  it("renders a not-found ErrorState with a link to /catalogo for an invalid productId", async () => {
    renderAt("/catalogo/produtos/abc");

    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Voltar para o catálogo" })).toHaveAttribute("href", "/catalogo");
    expect(screen.queryByRole("tab")).not.toBeInTheDocument();
  });
});
