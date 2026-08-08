import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { VinculosPage } from "./VinculosPage";
import { MARCA_UNAVAILABLE_DETAIL, wireCandidate } from "./wireFixtures";

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
  // The `listErpImports` mock that stood here was load-bearing until this branch
  // merged `main`: `VinculosPage.tsx` rendered `<ImportacaoSection />`, and the
  // previous comment said so — and said, in the same breath, that it "becomes
  // wrong the moment this branch merges". That moment is the merge commit
  // `293c1485`, which brought CHIP-IMPORT-CHAIN's move of the component to
  // `pages/importacoes/` (`4b76a287`) and took the two render lines out of
  // `VinculosPage.tsx` with it. Nothing under `pages/vinculos/` reaches
  // `useErpImports` any more; the only importers are `pages/importacoes/` and
  // `pages/integracoes/`.
  //
  // So the mock was deleted, not re-worded, and the deletion is the measurement:
  // a mock that is genuinely dead cannot turn this file red by leaving. It did
  // not.
  //
  // The `beforeEach` went with it rather than being kept and re-purposed: the
  // two surviving mocks were never reset here, and adding resets for them would
  // be a behaviour change smuggled in under a comment fix.

  it("renders tabs and shows the queue KPIs once loaded", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        // The 29th fixture, and the last one outside the mechanism. It carried
        // `95/ALTA` with `match_status: "REVIEW"` and `reasons: []` — three
        // impossibilities at once: ALTA is assigned at exactly one site and only
        // on the ACCEPT path (generation_service.go:503-505), and no candidate
        // is reason-less, because every declared anchor without a FOR/AGAINST
        // signal is emitted by the finalizer.
        //
        // It survived four reviewer rounds and this chip's own hand sweep by
        // living in a file the sweep did not own. That is the argument for a
        // throwing constructor over a checklist: the checklist has a scope, and
        // the fixture was outside it.
        wireCandidate({
          candidate_id: "cand_1",
          provider_item_id: "MLB1",
          state: "exact_sku",
          match_input: "seller_sku",
          match_status: "ACCEPT",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [
            {
              anchor: "seller_sku",
              direction: "FOR",
              detail: "seller_sku resolve exato para codprod",
            },
            { anchor: "ean", direction: "FOR", detail: "ean corrobora o mesmo codprod" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
      ],
    });
    listProductLinkWorkflows.mockResolvedValue({ items: [] });

    renderPage();

    expect(screen.getByRole("tab", { name: "Fila" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Resolvidos" })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "Vincular" })).toBeInTheDocument();
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

  it("renders only resolved links as real rows in the resolvidos tab, ignoring other states", async () => {
    listProductLinkCandidates.mockResolvedValue({ items: [] });
    listProductLinkWorkflows.mockResolvedValue({
      items: [
        {
          identity: { installation_id: "inst_1", provider_item_id: "MLB1" },
          current_link: {
            state: "resolved",
            updated_at: "2020-01-01T00:00:00Z",
            internal_product_id: 111,
            internal_product_name: "Produto Y",
          },
          candidates: [],
          audit: [
            { audit_id: "aud_1", next_state: "resolved", created_at: "2020-01-01T00:00:00Z" },
          ],
        },
        {
          identity: { installation_id: "inst_1", provider_item_id: "MLB2" },
          current_link: { state: "conflict", updated_at: "2020-01-01T00:00:00Z" },
          candidates: [],
          audit: [],
        },
        {
          identity: { installation_id: "inst_1", provider_item_id: "MLB3" },
          candidates: [],
          audit: [],
        },
      ],
    });

    renderPage();

    fireEvent.click(screen.getByRole("tab", { name: "Resolvidos" }));

    // Exactly one real resolved row (conflict + unresolved are excluded).
    const rows = await screen.findAllByTestId("resolvido-row");
    expect(rows).toHaveLength(1);
    expect(screen.getByText("MLB1")).toBeInTheDocument();
    expect(screen.getByText("Produto Y")).toBeInTheDocument();
    expect(screen.getByText("Vinculado ✓")).toBeInTheDocument();
    expect(screen.queryByText("MLB2")).not.toBeInTheDocument();
    // Desfazer is enabled because a resolving audit entry exists.
    expect(screen.getByRole("button", { name: "Desfazer" })).toBeEnabled();
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
