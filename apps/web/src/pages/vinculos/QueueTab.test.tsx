import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { QueueTab } from "./QueueTab";

const listProductLinkCandidates = vi.fn();
const approveProductLinkCandidate = vi.fn();
const rejectProductLinkListing = vi.fn();
const previewProductLinkBatch = vi.fn();
const applyProductLinkBatch = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listProductLinkCandidates: (...args: unknown[]) => listProductLinkCandidates(...args),
    approveProductLinkCandidate: (...args: unknown[]) => approveProductLinkCandidate(...args),
    rejectProductLinkListing: (...args: unknown[]) => rejectProductLinkListing(...args),
    previewProductLinkBatch: (...args: unknown[]) => previewProductLinkBatch(...args),
    applyProductLinkBatch: (...args: unknown[]) => applyProductLinkBatch(...args),
  }),
}));

function candidate(overrides: Partial<ProductLinkCandidateItem>): ProductLinkCandidateItem {
  return {
    candidate_id: "cand_x",
    installation_id: "inst_1",
    provider_code: "mercado_livre",
    provider_item_id: "MLB_X",
    state: "title_match",
    match_input: "title",
    confidence: 50,
    confidence_band: "MEDIA",
    match_status: "REVIEW",
    reasons: [],
    created_at: "2026-07-18T12:00:00Z",
    updated_at: "2026-07-18T12:00:00Z",
    ...overrides,
  };
}

function renderTab(initialEntries: string[] = ["/"]) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={initialEntries}>
        <QueueTab installationId="inst_1" />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  listProductLinkCandidates.mockReset();
  approveProductLinkCandidate.mockReset();
  rejectProductLinkListing.mockReset();
  previewProductLinkBatch.mockReset();
  applyProductLinkBatch.mockReset();
});

describe("QueueTab", () => {
  it("renders rows with distinct bands and anchor chips showing motivo together with the %", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_1",
          provider_item_id: "MLB1",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [{ anchor: "SKU idêntico", direction: "FOR", detail: "100%" }],
        }),
        candidate({
          candidate_id: "cand_2",
          provider_item_id: "MLB2",
          internal_product_id: 222,
          confidence: 60,
          confidence_band: "MEDIA",
          reasons: [{ anchor: "Título parcial", direction: "AGAINST", detail: "62%" }],
        }),
        candidate({
          candidate_id: "cand_3",
          provider_item_id: "MLB3",
          internal_product_id: 333,
          confidence: 30,
          confidence_band: "BAIXA",
          reasons: [{ anchor: "EAN", direction: "UNAVAILABLE", detail: "sem EAN" }],
        }),
      ],
    });

    renderTab();

    await waitFor(() => {
      expect(screen.getAllByTestId("queue-row")).toHaveLength(3);
    });

    // distinct bands
    expect(screen.getByText("ALTA")).toBeInTheDocument();
    expect(screen.getByText("MEDIA")).toBeInTheDocument();
    expect(screen.getByText("BAIXA")).toBeInTheDocument();

    // IC-01: the motivo (anchor) is always on screen. In the compact table cell a
    // FOR/UNAVAILABLE detail lives in the tooltip (the % is already its own
    // column), while an AGAINST detail — the reason the vínculo is blocked —
    // stays inline. A bare % never renders on its own either way.
    expect(screen.getByText("✓ SKU idêntico")).toHaveAttribute("title", "SKU idêntico: 100%");
    expect(screen.getByText("✕ Título parcial: 62%")).toBeInTheDocument();
    // A row whose only signal is UNAVAILABLE still shows that motivo — ranking,
    // never filtering (never blank / never a bare value).
    expect(screen.getByText("– EAN")).toHaveAttribute("title", "EAN: sem EAN");

    // Regression guard: confidence is already an integer 0-100 percentage (OpenAPI
    // ProductLinkCandidate.confidence). A candidate with confidence: 95 must render
    // exactly "95%" — never "9500%" (double-scaled) and never "0.95%" (unscaled raw).
    expect(screen.getByText("95%")).toBeInTheDocument();
    expect(screen.queryByText("9500%")).not.toBeInTheDocument();
  });

  it("collapses extra motivos behind +N and drops none of them (expansion shows the full wire form)", async () => {
    // The real CONFIRM shape: 2 corroborating anchors + the 2 internal-only
    // anchors that are always UNAVAILABLE. Four full sentences per row is what
    // made the table unreadable; only the two decisive ones stay on screen.
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_confirm",
          provider_item_id: "MLB4",
          internal_product_id: 10741,
          internal_product_name: "PUXADOR FENG",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato para codprod" },
            { anchor: "ean", direction: "UNAVAILABLE", detail: "sem EAN para corroborar o CODPROD" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: "marca inexistente no lado provider" },
            { anchor: "refforn", direction: "UNAVAILABLE", detail: "refforn inexistente no lado provider" },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    // Collapsed: the FOR anchor plus the highest-ranked remaining one, then +2.
    expect(screen.getByText("✓ SKU")).toBeInTheDocument();
    expect(screen.queryByText(/marca inexistente/)).not.toBeInTheDocument();
    const toggle = screen.getByRole("button", { name: "Mostrar todos os 4 motivos" });
    expect(toggle).toHaveTextContent("+2");

    fireEvent.click(toggle);

    // Expanded: every motivo, in its full wire form — nothing was dropped.
    for (const full of [
      "seller_sku: seller_sku resolve exato para codprod",
      "ean: sem EAN para corroborar o CODPROD",
      "marca: marca inexistente no lado provider",
      "refforn: refforn inexistente no lado provider",
    ]) {
      expect(screen.getByText(full)).toBeInTheDocument();
    }

    // And the action stays reachable in the same row.
    expect(within(row).getByRole("button", { name: "Vincular" })).toBeInTheDocument();
  });

  it("renders an honest NO_CANDIDATE row instead of a blank row", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_nc",
          provider_item_id: "MLB_NC",
          match_status: "NO_CANDIDATE",
          confidence: 0,
          confidence_band: "BAIXA",
          reasons: [],
        }),
      ],
    });

    renderTab();

    expect(await screen.findByText("Sem candidato")).toBeInTheDocument();
    const row = screen.getByTestId("queue-row");
    expect(row.getAttribute("data-match-status")).toBe("NO_CANDIDATE");
    // ADR-17 honest negative: no "Vincular" (nothing to link); Criar produto is a
    // disabled honest affordance (no create seam); only Ignorar is actionable.
    expect(screen.queryByRole("button", { name: "Vincular" })).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Criar produto" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Ignorar" })).toBeEnabled();
  });

  it("removes the item from the queue after an approve + page-local invalidation", async () => {
    listProductLinkCandidates
      .mockResolvedValueOnce({
        items: [
          candidate({
            candidate_id: "cand_1",
            provider_item_id: "MLB1",
            internal_product_id: 111,
            confidence: 95,
            confidence_band: "ALTA",
          }),
        ],
      })
      .mockResolvedValue({ items: [] });
    approveProductLinkCandidate.mockResolvedValue({
      link: {},
      audit: {},
    });

    renderTab();

    const approveButton = await screen.findByRole("button", { name: "Vincular" });
    fireEvent.click(approveButton);

    await waitFor(() => {
      expect(approveProductLinkCandidate).toHaveBeenCalledWith(
        expect.objectContaining({ candidate_id: "cand_1" }),
      );
    });

    // After the mutation settles, the page-local invalidate refetches the queue,
    // which now returns an empty list → the item leaves the queue.
    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.queryByText("MLB1")).not.toBeInTheDocument();
  });

  it("opens the drawer from a ?candidate= deep link (survives F5)", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_1",
          provider_item_id: "MLB1",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [{ anchor: "SKU idêntico", direction: "FOR", detail: "100%" }],
        }),
      ],
    });

    renderTab(["/?candidate=cand_1"]);

    // The drawer (DetailPanel complementary region) is open on first render.
    expect(await screen.findByRole("complementary", { name: "MLB1" })).toBeInTheDocument();
    expect(screen.getByTestId("drawer-candidate")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Vincular este candidato" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Ignorar anúncio" })).toBeInTheDocument();
  });

  it("bulk selects rows, previews as a pure dry-run, applies the valid subset, and clears selection + refetches the queue", async () => {
    listProductLinkCandidates
      .mockResolvedValueOnce({
        items: [
          candidate({ candidate_id: "cand_1", provider_item_id: "MLB1", internal_product_id: 111 }),
          candidate({ candidate_id: "cand_2", provider_item_id: "MLB2", internal_product_id: 222 }),
        ],
      })
      // After apply, the page-local queue refetches — the applied item leaves the fila.
      .mockResolvedValue({
        items: [candidate({ candidate_id: "cand_2", provider_item_id: "MLB2", internal_product_id: 222 })],
      });
    previewProductLinkBatch.mockResolvedValue({
      items: [
        { candidate_id: "cand_1", status: "OK" },
        { candidate_id: "cand_2", status: "FAILED", cause: "ALREADY_RESOLVED" },
      ],
    });
    applyProductLinkBatch.mockResolvedValue({
      batch_id: "batch_1",
      applied: [{ candidate_id: "cand_1" }],
      failed: [{ candidate_id: "cand_2", cause: "ALREADY_RESOLVED" }],
    });

    renderTab();

    const checkboxes = await screen.findAllByRole("checkbox");
    fireEvent.click(checkboxes[0]);
    fireEvent.click(checkboxes[1]);

    expect(await screen.findByText("2 selecionado(s)")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Pré-visualizar em lote" }));

    // Opening the preview is a pure dry-run — apply must not fire yet.
    await waitFor(() => {
      expect(previewProductLinkBatch).toHaveBeenCalledWith({
        approvals: [{ candidate_id: "cand_1" }, { candidate_id: "cand_2" }],
      });
    });
    expect(applyProductLinkBatch).not.toHaveBeenCalled();

    fireEvent.click(await screen.findByRole("button", { name: "Prosseguir só com válidos" }));

    await waitFor(() => {
      expect(applyProductLinkBatch).toHaveBeenCalledWith(
        expect.objectContaining({ approvals: [{ candidate_id: "cand_1" }] }),
      );
    });

    // Itemized applied + failed feedback renders (partial failure is normal).
    expect(await screen.findByText("1 aplicado(s), 1 falha(s)")).toBeInTheDocument();
    expect(screen.getByTestId("batch-result-failed")).toHaveTextContent("cand_2: ALREADY_RESOLVED");
    expect(screen.getByTestId("batch-result-applied")).toHaveTextContent("cand_1: aplicado");

    // Selection cleared post-apply.
    expect(screen.queryByText(/selecionado\(s\)/)).not.toBeInTheDocument();

    // Page-local queue was invalidated + refetched (2xx applied item left the fila).
    await waitFor(() => {
      expect(listProductLinkCandidates).toHaveBeenCalledTimes(2);
    });
    expect(screen.queryByText("MLB1")).not.toBeInTheDocument();
    expect(screen.getByText("MLB2")).toBeInTheDocument();
  });
});
