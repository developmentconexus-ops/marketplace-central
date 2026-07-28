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

  // V2 — the criterion this chip is judged on. `QueueRow` used to build the
  // collapsed cell by enumerating the directions as string literals
  // (`[...byDirection("AGAINST"), ...byDirection("FOR"),
  // ...byDirection("UNAVAILABLE")]`), which is type-correct and therefore silent
  // once D-B added a fourth direction: a row whose motivos are ALL INCOMPARABLE
  // fell through every branch, `shown` came out empty, and the cell rendered a
  // lone "+2" with zero chips — a filter wearing a ranking's comment (ADR-17).
  //
  // The assertions below are on the RENDERED DOM, never on the internal `shown`
  // array: an assertion about a variable does not prove a screen.
  // PRODUCIBILITY NOTE for the test below, found by the exhaustive fixture sweep
  // (evidence/V-fixture-producibility-sweep.md) and stated rather than left
  // implied: an all-INCOMPARABLE row requires a provider that DECLARES `marca`
  // as supplied, because `marca` has no case in `identityAnchorValues` and then
  // classifies INCOMPARABLE (generation_service.go:711). `mercado_livre`
  // declares only seller_sku/ean/title (capability_adapter.go:90), so for the
  // one provider with an adapter in this tree `marca` is UNAVAILABLE and a
  // 100%-INCOMPARABLE row cannot occur. The fixture is producible under a
  // capability declaration, not under today's only declaration.
  //
  // That is why the test immediately after it exists: it drives the SAME defect
  // from a reason set today's backend really emits.
  it("keeps a motivo on screen for a row whose reasons are ALL INCOMPARABLE (ADR-17)", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_inc",
          provider_item_id: "MLB_INC",
          internal_product_id: 444,
          internal_product_name: "PUXADOR FENG",
          reasons: [
            { anchor: "seller_sku", direction: "INCOMPARABLE", side: "provider", detail: "anúncio sem seller_sku" },
            { anchor: "ean", direction: "INCOMPARABLE", side: "erp", detail: "produto ERP sem ean cadastrado" },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    const chips = within(row).getAllByTestId("motivo-chip");

    // The invariant itself: at least one motivo visible. Never a bare "+N".
    expect(chips.length).toBeGreaterThan(0);
    expect(chips).toHaveLength(2);
    expect(within(row).queryByRole("button", { name: /Mostrar todos os/ })).not.toBeInTheDocument();

    // V4 — `side` is what D-B added on top of the direction: it says WHERE the
    // operator goes to fix the missing value. It reaches the screen from the
    // FIELD, not parsed out of the Portuguese `detail`.
    expect(chips[0]).toHaveTextContent("? SKU (falta no anúncio)");
    expect(chips[1]).toHaveTextContent("? EAN (falta no ERP)");

    // V1 — INCOMPARABLE carries its own glyph and its own token pair. Reusing
    // AGAINST's (warn) would say "blocks"; reusing UNAVAILABLE's (surface-2/
    // faint) would say "the provider never supplies this" — both false here.
    for (const chip of chips) {
      expect(chip).toHaveAttribute("data-direction", "INCOMPARABLE");
      expect(chip.className).toContain("bg-info-soft");
      expect(chip.className).toContain("text-info");
      expect(chip.className).not.toContain("bg-warn-soft");
      expect(chip.className).not.toContain("bg-accent-soft");
      expect(chip.className).not.toContain("bg-surface-2");
    }
  });

  // The V2 defect, driven from a reason set `mercado_livre` REALLY EMITS today —
  // no hypothetical capability declaration anywhere in it.
  //
  // This is verbatim what `applyUnresolvedScore` produces for a listing no
  // anchor could resolve (generation_service.go:620-628): the two anchors that
  // found nothing, INCOMPARABLE on the ERP side because no ERP product matched
  // at all, then the always-declared `marca` arriving UNAVAILABLE because ML
  // does not supply it. It is the single most common row on this screen.
  //
  // Against the old string-literal enumeration the cell showed exactly ONE chip
  // — "– Marca", the one thing nothing can ever be done about — and buried BOTH
  // actionable signals behind "+2". That is the same defect as the empty-cell
  // case, in the form the operator actually meets it: not "no motivo", but "the
  // only useless motivo, promoted".
  it("promotes the actionable absences over the permanent one, on a reason set the backend really emits", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_unresolved",
          provider_item_id: "MLB_UNRES",
          state: "unresolved",
          match_status: "NO_CANDIDATE",
          match_input: "none",
          confidence: 0,
          confidence_band: "BAIXA",
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

    renderTab();

    const row = await screen.findByTestId("queue-row");
    const chips = within(row).getAllByTestId("motivo-chip");

    // Two chips, and they are the two the operator can act on.
    expect(chips).toHaveLength(2);
    expect(chips[0]).toHaveTextContent("? SKU (falta no ERP)");
    expect(chips[1]).toHaveTextContent("? EAN (falta no ERP)");
    for (const chip of chips) {
      expect(chip).toHaveAttribute("data-direction", "INCOMPARABLE");
    }

    // The permanent absence is not dropped — it is RANKED, and lives behind the
    // toggle. Ranking, never filtering.
    expect(within(row).getByRole("button", { name: "Mostrar todos os 3 motivos" })).toBeInTheDocument();
    expect(within(row).queryByText("– Marca")).not.toBeInTheDocument();
  });

  it("ranks INCOMPARABLE above UNAVAILABLE without dropping either (ranking, never filtering)", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_mixed",
          provider_item_id: "MLB_MIX",
          internal_product_id: 555,
          reasons: [
            { anchor: "title", direction: "UNAVAILABLE", detail: "provider não fornece a âncora title" },
            { anchor: "ean", direction: "INCOMPARABLE", side: "both", detail: "anúncio e produto ERP sem ean" },
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato para codprod" },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    const chips = within(row).getAllByTestId("motivo-chip");
    // FOR outranks the two absence states; INCOMPARABLE (actionable) outranks
    // UNAVAILABLE (permanent), so the UNAVAILABLE one is the one deferred to +N.
    expect(chips.map((chip) => chip.getAttribute("data-direction"))).toEqual(["FOR", "INCOMPARABLE"]);
    expect(chips[1]).toHaveTextContent("? EAN (falta nos dois lados)");
    // Nothing was dropped — the remainder is behind the toggle, in full form.
    fireEvent.click(within(row).getByRole("button", { name: "Mostrar todos os 3 motivos" }));
    expect(screen.getByText("title: provider não fornece a âncora title")).toBeInTheDocument();
    expect(screen.getByText("ean (falta nos dois lados): anúncio e produto ERP sem ean")).toBeInTheDocument();
  });

  it("renders an INCOMPARABLE with no `side` on the wire without inventing one", async () => {
    // Real path: classifyProviderIdentityAnchor's `!readable` branch
    // (generation_service.go:711) emits INCOMPARABLE with an EMPTY side, which
    // `json:"side,omitempty"` drops from the payload — `marca` is the live case,
    // since ListingSnapshot has no provider-side brand field. Fabricating a side
    // there would forge the exact datum D-B was created to carry (ADR-17).
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_noside",
          provider_item_id: "MLB_NOSIDE",
          internal_product_id: 666,
          // PRODUCIBILITY: this requires a provider that DECLARES `marca` as
          // supplied — then `marca` reaches `classifyProviderIdentityAnchor`
          // with `Supplied: true`, has no case in `identityAnchorValues`, and
          // classifies INCOMPARABLE with no side (generation_service.go:706-708).
          // `mercado_livre` declares only seller_sku/ean/title
          // (capability_adapter.go:90), so for it `marca` arrives UNAVAILABLE
          // instead. Producible under a capability declaration, not under
          // today's only one — stated because the fixture otherwise reads as
          // "this is what the wire sends".
          reasons: [
            { anchor: "marca", direction: "INCOMPARABLE", detail: "não foi possível comparar a âncora marca" },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    const [chip] = within(row).getAllByTestId("motivo-chip");
    expect(chip).toHaveTextContent("? Marca");
    expect(chip.textContent).not.toMatch(/falta n/);
    expect(chip).toHaveAttribute("title", "marca: não foi possível comparar a âncora marca");
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

  // V10 — neutral vocabulary must not cost the operator the provider datum.
  // Both halves are asserted here because neutralizing the VALUE would be a lie,
  // not neutrality: the anúncio really is on one specific marketplace.
  it("neutralizes the structural labels while keeping WHICH provider the anúncio is on", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_1",
          provider_item_id: "MLB1",
          provider_code: "mercado_livre",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
        }),
      ],
    });

    const { container } = renderTab();
    const row = await screen.findByTestId("queue-row");

    // Structural labels: no provider name in a column header.
    expect(screen.getByRole("columnheader", { name: "Anúncio" })).toBeInTheDocument();
    expect(screen.getByRole("columnheader", { name: "Canal" })).toBeInTheDocument();
    expect(screen.queryByRole("columnheader", { name: "Anúncio ML" })).not.toBeInTheDocument();
    expect(screen.queryByRole("columnheader", { name: "SKU ML" })).not.toBeInTheDocument();

    // Data value: the provider is still on screen — by display name.
    expect(within(row).getByText("Mercado Livre")).toBeInTheDocument();

    // And never the raw wire slug. This is the trap that got CHIP-PED-FILA on
    // four surfaces; the "SKU ML" column was rendering `provider_code`, which
    // the wire fills with the marketplace slug, never a SKU.
    expect(container.textContent).not.toContain("mercado_livre");
  });

  it("typesets an unmapped provider without ever collapsing two distinct codes onto one name", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_1",
          provider_item_id: "SHP1",
          // A provider with no entry in the display-name map. Today only
          // `mercado_livre` is mapped, so the very next marketplace lands here —
          // and V10's rule is about the SCREEN, not about which slug it is.
          provider_code: "shopee",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
        }),
        // Two codes that a naive "split on any separator" would render
        // identically. `registry.go:100-114` dedupes provider codes by exact
        // string equality, so both can be registered at once — and then the
        // operator cannot tell whose listing they are looking at.
        candidate({
          candidate_id: "cand_2",
          provider_item_id: "AMZ1",
          provider_code: "amazon_marketplace",
          internal_product_id: 222,
          internal_product_name: "Parafuso B",
        }),
        candidate({
          candidate_id: "cand_3",
          provider_item_id: "AMZ2",
          provider_code: "amazon-marketplace",
          internal_product_id: 333,
          internal_product_name: "Parafuso C",
        }),
      ],
    });

    const { container } = renderTab();
    const rows = await screen.findAllByTestId("queue-row");
    const [shopee, underscored, hyphenated] = rows;

    // The clean underscore slug is typeset...
    expect(within(shopee).getByText("Shopee")).toBeInTheDocument();
    expect(within(underscored).getByText("Amazon Marketplace")).toBeInTheDocument();

    // ...and the hyphenated code keeps its hyphen, so the two providers never
    // wear the same name. This is the assertion neither a constant fallback
    // ("always Shopee") nor a separator-collapsing typeset can survive.
    expect(within(hyphenated).getByText("Amazon-marketplace")).toBeInTheDocument();
    expect(within(hyphenated).queryByText("Amazon Marketplace")).not.toBeInTheDocument();

    // No lower-case slug of a provider we CAN name reaches the screen.
    expect(container.textContent).not.toContain("shopee");
    expect(container.textContent).not.toContain("amazon_marketplace");
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

  // V6 / D-122 — "Identificado por" shows the anchors that DECIDED, ` + `-joined.
  //
  // The negative case is the one that matters and the one A2-R2 named: `title`
  // produces a FOR reason in the base (generation_service.go:551, "match por
  // título (ranking-only, nunca ACCEPT)") and D-121 routes title-only to REVIEW.
  // Deriving this column from `direction === "FOR"` would print "Título" as the
  // deciding anchor for a candidate title explicitly did NOT decide. The source
  // is the same pair the backend files the decision from — match_status +
  // match_input (resolution_service.go:812-835).
  it("names only the anchors that decided, and names none for a title-only REVIEW row", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        // Corroborated → both anchors, joined.
        candidate({
          candidate_id: "cand_accept",
          provider_item_id: "MLB_ACCEPT",
          state: "exact_ean",
          match_status: "ACCEPT",
          match_input: "ean",
          match_value: "7890000000001",
          reasons: [{ anchor: "ean", direction: "FOR", detail: "EAN idêntico" }],
        }),
        // Single anchor resolved a single product → confirmation queue, one anchor.
        candidate({
          candidate_id: "cand_confirm_sku",
          provider_item_id: "MLB_CONFIRM",
          state: "exact_sku",
          match_status: "CONFIRM",
          match_input: "seller_sku",
          reasons: [{ anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato" }],
        }),
        // Title match: a FOR reason on screen, but nothing decided.
        candidate({
          candidate_id: "cand_title",
          provider_item_id: "MLB_TITLE",
          state: "title_match",
          match_status: "REVIEW",
          match_input: "title",
          reasons: [{ anchor: "title", direction: "FOR", detail: "match por título" }],
        }),
      ],
    });

    renderTab();

    const rows = await screen.findAllByTestId("queue-row");
    const [accept, confirm, title] = rows;

    expect(within(accept).getByTestId("identificado-por")).toHaveTextContent("CODPROD + EAN");
    expect(within(confirm).getByTestId("identificado-por")).toHaveTextContent("CODPROD");
    expect(within(confirm).getByTestId("identificado-por")).not.toHaveTextContent("EAN");

    // The negative case: the title row shows its motivo but names NO anchor.
    expect(within(title).queryByTestId("identificado-por")).not.toBeInTheDocument();
    expect(within(title).getByText("✓ Título")).toBeInTheDocument();
    // ...and the column slot itself is the honest unknown, not a blank cell:
    // index 5 is IDENTIFICADO POR (0 seleção, 1 anúncio, 2 canal, 3 produto,
    // 4 SKU HUB, 5 identificado por). Read positionally because "—" also
    // renders in the SKU HUB cell of this same row.
    expect(title.querySelectorAll("td")[5].textContent).toBe("—");
    expect(title.textContent).not.toContain("Título +");

    // Header renamed per D-122: the single-anchor GTIN reading is superseded.
    expect(screen.getByRole("columnheader", { name: "Identificado por" })).toBeInTheDocument();
    expect(screen.queryByRole("columnheader", { name: "GTIN" })).not.toBeInTheDocument();
  });

  it("survives a match_status the SDK union does not know, degrading one cell instead of the table", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_drift",
          provider_item_id: "MLB_DRIFT",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
          // Wire drift: the API ships a sixth status before the SDK is
          // regenerated. The cast is the POINT of the test — the compile-time
          // exhaustiveness of the maps cannot see a value the wire invents at
          // runtime, so the runtime has to stay honest on its own.
          match_status: "PENDING_REVIEW" as ProductLinkCandidateItem["match_status"],
          match_input: "ean",
          reasons: [{ anchor: "ean", direction: "FOR", detail: "EAN idêntico" }],
        }),
      ],
    });

    renderTab();

    // The table still renders: an unknown status must not take the queue down.
    const row = await screen.findByTestId("queue-row");
    expect(within(row).getByText("Parafuso A")).toBeInTheDocument();

    // And the cell says "unknown" rather than naming an anchor it cannot know.
    expect(within(row).queryByTestId("identificado-por")).not.toBeInTheDocument();
    expect(row.querySelectorAll("td")[5].textContent).toBe("—");
  });

  // Same wire-drift class as the test above, on the OTHER two unions the row
  // indexes by Record. `direction` matters most: the compact cell used to
  // enumerate three literals, which silently DROPPED an unknown direction; the
  // ranking sort that replaced it (the V2 fix) keeps every reason, so the fix
  // made this path reachable. Undoing a silent filter must not install a
  // literal "undefined" in its place.
  it("survives a reason direction the SDK union does not know, without printing undefined", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_dir",
          provider_item_id: "MLB_DIR",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
          reasons: [
            { anchor: "ean", direction: "PARTIAL" as ProductLinkCandidateItem["reasons"][number]["direction"], detail: "ean parcial" },
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato" },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    const chips = within(row).getAllByTestId("motivo-chip");

    // Both reasons keep their place — an unknown direction is not filtered out.
    expect(chips).toHaveLength(2);

    // The known one is unaffected...
    expect(within(row).getByText("✓ SKU")).toBeInTheDocument();

    // ...and the unknown one falls through VERBATIM, exactly as an unknown
    // anchor does (V9): the wire word instead of a glyph we do not have.
    // Never the string "undefined", in the text or in the class attribute.
    expect(within(row).getByText("PARTIAL EAN")).toBeInTheDocument();
    expect(row.textContent).not.toContain("undefined");
    for (const chip of chips) {
      expect(chip.getAttribute("class") ?? "").not.toContain("undefined");
    }

    // The unknown direction ranks LAST, so it never displaces a signal the
    // operator can act on out of the two compact slots.
    expect(chips[0].textContent).toBe("✓ SKU");
    expect(chips[1].textContent).toBe("PARTIAL EAN");
  });

  it("survives a confidence_band the SDK union does not know, without printing undefined", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_band",
          provider_item_id: "MLB_BAND",
          internal_product_id: 111,
          internal_product_name: "Parafuso A",
          confidence: 72,
          confidence_band: "CRITICA" as ProductLinkCandidateItem["confidence_band"],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");

    // The band the wire actually sent, verbatim — not a blank pill, and not a
    // band we would be inventing by falling back to one of the three we know.
    expect(within(row).getByText("CRITICA")).toBeInTheDocument();
    // The confidence itself is a number and is unaffected by the band drift.
    expect(within(row).getByText("72%")).toBeInTheDocument();
    expect(row.textContent).not.toContain("undefined");
    for (const el of Array.from(row.querySelectorAll<HTMLElement>("[class]"))) {
      expect(el.getAttribute("class") ?? "").not.toContain("undefined");
    }
  });
});
