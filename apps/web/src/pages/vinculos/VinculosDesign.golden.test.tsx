import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { VinculosPage } from "./VinculosPage";

// GOLDEN (day-1, MIS-004-C10 gap #2). Locks the design contract for /vinculos:
//   1. anúncio-cêntrica 9-col table (ANÚNCIO·CANAL·PRODUTO SUGERIDO·SKU HUB·
//      IDENTIFICADO POR·CONFIANÇA·MOTIVO·AÇÃO) with the EXEMPLO-IO row.
//   2. honest NO_CANDIDATE row (conf 0, "sem candidato", Criar produto / Ignorar) — ADR-17.
//   3. paper+green tokens only — never literal slate/blue/emerald/red/amber-<n> off-theme.
//
// F-05 (D-122) renamed three of those headers. "ANÚNCIO ML" → "ANÚNCIO" is the
// neutral structural label. "SKU ML" → "CANAL" is a correction as much as a
// neutralization: that cell has always rendered `provider_code`, which the wire
// fills with the marketplace SLUG ("mercado_livre") — the candidate contract
// carries no seller SKU at all. "GTIN" → "IDENTIFICADO POR" replaces a
// single-anchor reading with the SET of anchors that decided, which supersedes
// it. The column count and order are unchanged.

const listProductLinkCandidates = vi.fn();
const listProductLinkWorkflows = vi.fn();
const listErpImports = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listProductLinkCandidates: (...a: unknown[]) => listProductLinkCandidates(...a),
    listProductLinkWorkflows: (...a: unknown[]) => listProductLinkWorkflows(...a),
    listErpImports: (...a: unknown[]) => listErpImports(...a),
    approveProductLinkCandidate: vi.fn(),
    rejectProductLinkListing: vi.fn(),
    previewProductLinkBatch: vi.fn(),
    applyProductLinkBatch: vi.fn(),
  }),
}));

vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function base(overrides: Partial<ProductLinkCandidateItem>): ProductLinkCandidateItem {
  return {
    candidate_id: "cand_x",
    installation_id: "inst_1",
    provider_code: "mercado_livre",
    provider_item_id: "MLB123",
    state: "exact_ean",
    match_input: "ean",
    match_value: "7890000000001",
    confidence: 92,
    confidence_band: "ALTA",
    // A single EAN resolving a single product is the CONFIRM queue under D-121
    // — it waits for a human yes. The fixture said REVIEW, which is the
    // anchors-disagree state and would have named no deciding anchor at all.
    match_status: "CONFIRM",
    reasons: [{ anchor: "EAN idêntico", direction: "FOR", detail: "100%" }],
    created_at: "2026-07-19T12:00:00Z",
    updated_at: "2026-07-19T12:00:00Z",
    ...overrides,
  };
}

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

// Off-theme literal Tailwind color utilities carry a numeric shade (slate-200,
// blue-600, emerald-100, red-700, amber-800). The design tokens never do
// (bg-surface, text-ink, bg-accent-soft, text-amber, text-warn). So any
// `-<color>-<digit>` in the markup is an off-theme regression.
const OFF_THEME = /\b(?:bg|text|border|ring|divide|from|to|hover:bg|hover:text|hover:border)-(?:slate|blue|emerald|red|amber|indigo|gray|zinc|green|sky|violet)-\d/;

describe("Vínculos design golden", () => {
  beforeEach(() => {
    listErpImports.mockReset();
    listErpImports.mockResolvedValue({ items: [] });
    listProductLinkWorkflows.mockReset();
    listProductLinkWorkflows.mockResolvedValue({ items: [] });
    listProductLinkCandidates.mockReset();
  });

  it("renders the EXEMPLO-IO row in anúncio-cêntrica 9-col layout", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        base({
          candidate_id: "cand_1",
          provider_item_id: "MLB123",
          provider_code: "mercado_livre",
          internal_product_id: 1001,
          internal_product_name: "Produto Y",
        }),
      ],
    });

    renderPage();

    const row = await screen.findByTestId("queue-row");
    const cells = within(row);
    expect(cells.getByText("MLB123")).toBeInTheDocument(); // ANÚNCIO
    // CANAL — the provider, by display name. The raw wire slug never renders.
    expect(cells.getByText("Mercado Livre")).toBeInTheDocument();
    expect(cells.queryByText("mercado_livre")).not.toBeInTheDocument();
    expect(cells.getByText("Produto Y")).toBeInTheDocument(); // PRODUTO SUGERIDO
    expect(cells.getByText("1001")).toBeInTheDocument(); // SKU HUB
    // IDENTIFICADO POR — the anchor that decided. A single EAN resolving a
    // single product is the confirmation queue (D-121), so it names EAN alone;
    // "CODPROD + EAN" is reserved for the corroborated ACCEPT.
    expect(cells.getByTestId("identificado-por")).toHaveTextContent("EAN");
    expect(cells.getByText("92%")).toBeInTheDocument(); // CONFIANÇA
    // MOTIVO — compact chip: motivo (anchor) sempre visível, detail no tooltip
    // (IC-01: % nunca aparece sozinho). Forma completa fica no title e na
    // expansão "+N" / no drawer.
    const motivoChip = cells.getByText("✓ EAN idêntico");
    expect(motivoChip).toBeInTheDocument();
    expect(motivoChip).toHaveAttribute("title", "EAN idêntico: 100%");
    expect(cells.getByRole("button", { name: "Vincular" })).toBeInTheDocument(); // AÇÃO

    // Column headers present (anúncio-cêntrica order).
    for (const header of ["Anúncio", "Canal", "Produto sugerido", "SKU HUB", "Identificado por", "Confiança", "Motivo", "Ação"]) {
      expect(screen.getByRole("columnheader", { name: header })).toBeInTheDocument();
    }
  });

  it("renders an honest NO_CANDIDATE row (ADR-17): conf 0, sem candidato, Criar produto / Ignorar", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        base({
          candidate_id: "cand_nc",
          provider_item_id: "MLB999",
          provider_code: "xyz-9",
          match_status: "NO_CANDIDATE",
          match_input: "none",
          match_value: undefined,
          confidence: 0,
          confidence_band: "BAIXA",
          internal_product_id: undefined,
          internal_product_name: undefined,
          reasons: [],
        }),
      ],
    });

    renderPage();

    const row = await screen.findByTestId("queue-row");
    const cells = within(row);
    expect(cells.getByText("Sem candidato")).toBeInTheDocument();
    // Nothing decided a NO_CANDIDATE row → "—", never a fabricated anchor.
    expect(cells.queryByTestId("identificado-por")).not.toBeInTheDocument();
    // No fabricated confidence: never a "0%" verde chip; band chip absent.
    expect(cells.queryByText("ALTA")).not.toBeInTheDocument();
    // Negative actions only.
    expect(cells.getByRole("button", { name: "Criar produto" })).toBeInTheDocument();
    expect(cells.getByRole("button", { name: "Ignorar" })).toBeInTheDocument();
    expect(cells.queryByRole("button", { name: "Vincular" })).not.toBeInTheDocument();
  });

  it("uses paper+green tokens only — no off-theme literal color classes anywhere on the page", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        base({ candidate_id: "cand_1", internal_product_id: 1001, internal_product_name: "Produto Y" }),
        base({
          candidate_id: "cand_nc",
          provider_item_id: "MLB999",
          match_status: "NO_CANDIDATE",
          confidence: 0,
          confidence_band: "BAIXA",
          reasons: [],
        }),
      ],
    });

    const { container } = renderPage();
    await screen.findAllByTestId("queue-row");

    const offenders = Array.from(container.querySelectorAll<HTMLElement>("[class]"))
      .map((el) => el.getAttribute("class") ?? "")
      .filter((cls) => OFF_THEME.test(cls));
    expect(offenders).toEqual([]);
  });
});
