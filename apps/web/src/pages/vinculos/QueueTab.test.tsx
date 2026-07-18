import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
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
    confidence: 0.5,
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
          confidence: 0.95,
          confidence_band: "ALTA",
          reasons: [{ anchor: "SKU idêntico", direction: "FOR", detail: "100%" }],
        }),
        candidate({
          candidate_id: "cand_2",
          provider_item_id: "MLB2",
          internal_product_id: 222,
          confidence: 0.6,
          confidence_band: "MEDIA",
          reasons: [{ anchor: "Título parcial", direction: "AGAINST", detail: "62%" }],
        }),
        candidate({
          candidate_id: "cand_3",
          provider_item_id: "MLB3",
          internal_product_id: 333,
          confidence: 0.3,
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

    // IC-01: the motivo (anchor) and the % render together in the same chip.
    expect(screen.getByText("SKU idêntico: 100%")).toBeInTheDocument();
    expect(screen.getByText("Título parcial: 62%")).toBeInTheDocument();
    // A reason with no % still renders its motivo (never blank / never a bare value).
    expect(screen.getByText("EAN: sem EAN")).toBeInTheDocument();
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
    // Approve is disabled — there is nothing to approve.
    expect(screen.getByRole("button", { name: "Aprovar" })).toBeDisabled();
  });

  it("removes the item from the queue after an approve + page-local invalidation", async () => {
    listProductLinkCandidates
      .mockResolvedValueOnce({
        items: [
          candidate({
            candidate_id: "cand_1",
            provider_item_id: "MLB1",
            internal_product_id: 111,
            confidence: 0.95,
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

    const approveButton = await screen.findByRole("button", { name: "Aprovar" });
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
          confidence: 0.95,
          confidence_band: "ALTA",
          reasons: [{ anchor: "SKU idêntico", direction: "FOR", detail: "100%" }],
        }),
      ],
    });

    renderTab(["/?candidate=cand_1"]);

    // The drawer (DetailPanel complementary region) is open on first render.
    expect(await screen.findByRole("complementary", { name: "MLB1" })).toBeInTheDocument();
    expect(screen.getByTestId("drawer-candidate")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Aprovar este candidato" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Rejeitar anúncio" })).toBeInTheDocument();
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
