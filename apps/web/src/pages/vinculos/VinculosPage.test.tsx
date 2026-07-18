import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { VinculosPage } from "./VinculosPage";

const listProductLinkCandidates = vi.fn();
const listProductLinkWorkflows = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listProductLinkCandidates: (...args: unknown[]) => listProductLinkCandidates(...args),
    listProductLinkWorkflows: (...args: unknown[]) => listProductLinkWorkflows(...args),
  }),
}));

vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <VinculosPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("VinculosPage", () => {
  it("renders tabs and shows the queue KPIs once loaded", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        {
          candidate_id: "cand_1",
          installation_id: "inst_1",
          provider_code: "mercado_livre",
          provider_item_id: "MLB1",
          state: "exact_sku",
          match_input: "seller_sku",
          confidence: 0.95,
          confidence_band: "ALTA",
          match_status: "REVIEW",
          reasons: [],
          created_at: "2026-07-18T12:00:00Z",
          updated_at: "2026-07-18T12:00:00Z",
        },
      ],
    });
    listProductLinkWorkflows.mockResolvedValue({ items: [] });

    renderPage();

    expect(screen.getByRole("tab", { name: "Fila" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Resolvidos" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "Aprovar" })).toBeInTheDocument();
    });
    expect(screen.getByText("MLB1")).toBeInTheDocument();
  });

  it("shows loading then empty state for the resolvidos tab", async () => {
    listProductLinkCandidates.mockResolvedValue({ items: [] });
    listProductLinkWorkflows.mockResolvedValue({ items: [] });

    renderPage();

    fireEvent.click(screen.getByRole("tab", { name: "Resolvidos" }));

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });

  it("counts only resolved links in the resolvidos tab body, ignoring other states", async () => {
    listProductLinkCandidates.mockResolvedValue({ items: [] });
    listProductLinkWorkflows.mockResolvedValue({
      items: [
        { identity: { installation_id: "inst_1", provider_item_id: "MLB1" }, current_link: { state: "resolved", updated_at: "2020-01-01T00:00:00Z" }, candidates: [], audit: [] },
        { identity: { installation_id: "inst_1", provider_item_id: "MLB2" }, current_link: { state: "conflict", updated_at: "2020-01-01T00:00:00Z" }, candidates: [], audit: [] },
        { identity: { installation_id: "inst_1", provider_item_id: "MLB3" }, candidates: [], audit: [] },
      ],
    });

    renderPage();

    fireEvent.click(screen.getByRole("tab", { name: "Resolvidos" }));

    expect(
      await screen.findByText("1 vínculo(s) — tabela detalhada chega na próxima etapa."),
    ).toBeInTheDocument();
  });

  it("renders an error state with retry for the fila tab", async () => {
    listProductLinkCandidates.mockRejectedValueOnce({ status: 500, error: { code: "internal" } });
    listProductLinkCandidates.mockResolvedValue({ items: [] });
    listProductLinkWorkflows.mockResolvedValue({ items: [] });

    renderPage();

    const retry = await screen.findByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });
});
