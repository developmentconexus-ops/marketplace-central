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
    // A single EAN resolving a single product is the CONFIRM queue under D-121
    // — it waits for a human yes. The fixture said REVIEW, which is the
    // anchors-disagree state and would have named no deciding anchor at all.
    //
    // The rest of the row is the SAME producible case, taken from the generator
    // that emits it (`applySingleAnchorScore`, state `exact_ean`,
    // generation_service.go:538-547): 60/MEDIA, an `ean` FOR, and the missing
    // corroborating anchor as INCOMPARABLE on the provider side. The old
    // 92/ALTA was not reachable for any status — ALTA is emitted at exactly one
    // place (95, `buildConcordantCandidate`, generation_service.go:503-505) and
    // only for ACCEPT — so the golden was locking a row the backend could never
    // produce. A design gate whose example cannot occur gates a fiction.
    confidence: 60,
    confidence_band: "MEDIA",
    match_status: "CONFIRM",
    reasons: [
      { anchor: "ean", direction: "FOR", detail: "ean corrobora codprod (unproved)" },
      {
        anchor: "seller_sku",
        direction: "INCOMPARABLE",
        side: "provider",
        detail: "sem CODPROD para corroborar o EAN",
      },
      // The THIRD reason every candidate carries, and the one that makes this
      // fixture producible instead of merely plausible.
      // `appendProviderDeclaredUnavailableReasons` walks the provider's DECLARED
      // anchors (generation_service.go:671-698) and emits one for every anchor
      // with no FOR/AGAINST signal. `marca` is always in `KnownIdentityAnchors`
      // (connectors/ports/marketplace_capability.go:40-45), so the adapter always
      // declares it (identity_anchor_adapter.go:28-35) — and `mercado_livre`
      // declares only seller_sku/ean/title as supplied
      // (mercado_livre/capability_adapter.go:90), so `marca` arrives
      // `Supplied: false` and classifies UNAVAILABLE at generation_service.go:704.
      // Every mercado_livre candidate has it. A fixture without it locks a row
      // the backend never emits.
      { anchor: "marca", direction: "UNAVAILABLE", detail: "provider não fornece a âncora marca" },
    ],
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
    // CONFIANÇA — banda E %. The band pill is asserted POSITIVELY, with its
    // token pair: the off-theme sweep below is a negative test and would still
    // pass if the pill lost its tokens entirely.
    expect(cells.getByText("60%")).toBeInTheDocument();
    const bandPill = cells.getByText("MEDIA");
    expect(bandPill).toHaveClass("bg-amber-soft", "text-amber");
    // MOTIVO — compact chip: motivo (anchor) sempre visível, detail no tooltip
    // (IC-01: % nunca aparece sozinho). Forma completa fica no title e na
    // expansão "+N" / no drawer.
    const motivoChip = cells.getByText("✓ EAN");
    expect(motivoChip).toBeInTheDocument();
    expect(motivoChip).toHaveAttribute("title", "ean: ean corrobora codprod (unproved)");
    // ...and the anchor that is MISSING rides beside it with the side that says
    // where to go fix it (D-122/D-B), in its own tokens.
    const incomparableChip = cells.getByText("? SKU (falta no anúncio)");
    expect(incomparableChip).toBeInTheDocument();
    expect(incomparableChip).toHaveAttribute("data-direction", "INCOMPARABLE");
    // ...and the third reason (the always-present `marca`) does NOT get a third
    // chip: the cell caps at two and offers the rest behind "+1". This is the
    // layout the backend actually produces, which is the point of a golden.
    expect(cells.getAllByTestId("motivo-chip")).toHaveLength(2);
    expect(cells.getByRole("button", { name: "Mostrar todos os 3 motivos" })).toBeInTheDocument();
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
          match_status: "NO_CANDIDATE",
          match_input: "none",
          match_value: undefined,
          confidence: 0,
          confidence_band: "BAIXA",
          internal_product_id: undefined,
          internal_product_name: undefined,
          // `applyUnresolvedScore` (generation_service.go:620-628) is the ONLY
          // path that emits NO_CANDIDATE, and it is reached through
          // `newCandidate(..., LinkCandidateStateUnresolved, ...MatchInputNone, ...)`
          // at :215/:379. So `state` is `unresolved` — the fixture inherited
          // `exact_ean` from base(), a pair the generator cannot produce.
          state: "unresolved",
          // And it is never reason-less: the two anchors that found nothing are
          // named (INCOMPARABLE on the ERP side, because no ERP product matched
          // at all — missingMatchedAnchorReason with a nil product,
          // generation_service.go:631-645), plus the always-declared `marca`.
          // "sem candidato" with an empty Motivo cell is exactly the silent row
          // M05-C5 forbids.
          reasons: [
            {
              anchor: "seller_sku",
              direction: "INCOMPARABLE",
              side: "erp",
              detail: "seller_sku sem correspondência",
            },
            { anchor: "ean", direction: "INCOMPARABLE", side: "erp", detail: "ean sem correspondência" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: "provider não fornece a âncora marca" },
          ],
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
        // The corroborated ACCEPT, verbatim from `buildConcordantCandidate`
        // (generation_service.go:495-505). It is here so the sweep still covers
        // the ALTA/accent token path: the EXEMPLO-IO row above is a MEDIA row
        // now, because 92/ALTA was never a producible CONFIRM.
        base({
          candidate_id: "cand_1",
          internal_product_id: 1001,
          internal_product_name: "Produto Y",
          state: "exact_sku",
          match_input: "seller_sku",
          match_status: "ACCEPT",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato para codprod" },
            { anchor: "ean", direction: "FOR", detail: "ean corrobora o mesmo codprod (unproved)" },
            // Same always-declared `marca` as in base(): `buildConcordantCandidate`
            // (generation_service.go:493-507) also finalizes through
            // `appendProviderDeclaredUnavailableReasons`.
            { anchor: "marca", direction: "UNAVAILABLE", detail: "provider não fornece a âncora marca" },
          ],
        }),
        base({
          candidate_id: "cand_nc",
          provider_item_id: "MLB999",
          match_status: "NO_CANDIDATE",
          match_input: "none",
          match_value: undefined,
          confidence: 0,
          confidence_band: "BAIXA",
          internal_product_id: undefined,
          internal_product_name: undefined,
          state: "unresolved",
          reasons: [
            {
              anchor: "seller_sku",
              direction: "INCOMPARABLE",
              side: "erp",
              detail: "seller_sku sem correspondência",
            },
            { anchor: "ean", direction: "INCOMPARABLE", side: "erp", detail: "ean sem correspondência" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: "provider não fornece a âncora marca" },
          ],
        }),
      ],
    });

    const { container } = renderPage();
    const [accept] = await screen.findAllByTestId("queue-row");

    // Positive ALTA coverage. The sweep below only proves no OFF-theme class is
    // present; it would pass just as well if the pill rendered with no tokens
    // at all, so the accent pair is asserted here, on the row that carries it.
    expect(within(accept).getByText("ALTA")).toHaveClass("bg-accent-soft", "text-accent-ink");

    const offenders = Array.from(container.querySelectorAll<HTMLElement>("[class]"))
      .map((el) => el.getAttribute("class") ?? "")
      .filter((cls) => OFF_THEME.test(cls));
    expect(offenders).toEqual([]);
  });
});
