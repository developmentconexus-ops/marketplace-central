import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { AnunciosPage } from "./AnunciosPage";

const listListings = vi.fn();
const getListingsSummary = vi.fn(() =>
  Promise.resolve({
    total: 1,
    active: 1,
    paused: 0,
    exceptions: { below_margin_worst_case: null, margin_unknown: null },
  }),
);

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({ listListings: (...args: unknown[]) => listListings(...args), getListingsSummary }),
}));

vi.mock("../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function renderPage(search: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/anuncios${search}`]}>
        <AnunciosPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("AnunciosPage error recovery", () => {
  it("clears filter params on retry only for invalid_filter errors", async () => {
    listListings.mockRejectedValueOnce({
      status: 400,
      error: { code: "invalid_filter", message: "invalid" },
    });
    listListings.mockResolvedValue({ items: [] });
    renderPage("?installation=inst_1&filter.exception=sync_error&q=abc");

    const [retry] = await screen.findAllByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("0 anúncio(s) encontrado(s).")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", q: "abc" });
  });

  it("retries with filters intact when the error is not invalid_filter", async () => {
    listListings.mockRejectedValueOnce({
      status: 500,
      error: { code: "internal", message: "boom" },
    });
    listListings.mockResolvedValue({ items: [] });
    renderPage("?installation=inst_1&filter.exception=sync_error");

    const [retry] = await screen.findAllByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("0 anúncio(s) encontrado(s).")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", exception: "sync_error" });
  });
});
