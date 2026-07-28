import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ProductLinkCandidateItem } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { QueueTab } from "./QueueTab";
import { MARCA_UNAVAILABLE_DETAIL, driftCandidate, wireCandidate } from "./wireFixtures";

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

// Every fixture below is built through `wireCandidate`, which THROWS on a
// candidate the generator cannot emit (see wireFixtures.ts for why the previous
// arrangement — a hand-written literal per test — failed four reviewer rounds).
// A deliberately impossible shape goes through `driftCandidate`, which requires
// the reason to be written down.
const candidate = wireCandidate;

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
          state: "exact_sku",
          match_status: "ACCEPT",
          match_input: "seller_sku",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "100%" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
        candidate({
          candidate_id: "cand_2",
          provider_item_id: "MLB2",
          internal_product_id: 222,
          state: "exact_ean",
          match_status: "CONFIRM",
          match_input: "ean",
          confidence: 60,
          confidence_band: "MEDIA",
          reasons: [
            { anchor: "title", direction: "AGAINST", detail: "62%" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
        candidate({
          candidate_id: "cand_3",
          provider_item_id: "MLB3",
          internal_product_id: 333,
          state: "conflict",
          confidence: 20,
          confidence_band: "BAIXA",
          reasons: [
            { anchor: "ean", direction: "UNAVAILABLE", detail: "sem EAN" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
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
    expect(screen.getByText("✓ SKU")).toHaveAttribute("title", "seller_sku: 100%");
    expect(screen.getByText("✕ Título: 62%")).toBeInTheDocument();
    // A row whose only signals are UNAVAILABLE still shows them — ranking,
    // never filtering (never blank / never a bare value).
    expect(screen.getByText("– EAN")).toHaveAttribute("title", "ean: sem EAN");

    // Regression guard: confidence is already an integer 0-100 percentage (OpenAPI
    // ProductLinkCandidate.confidence). A candidate with confidence: 95 must render
    // exactly "95%" — never "9500%" (double-scaled) and never "0.95%" (unscaled raw).
    expect(screen.getByText("95%")).toBeInTheDocument();
    expect(screen.queryByText("9500%")).not.toBeInTheDocument();
  });

  it("collapses extra motivos behind +N and drops none of them (expansion shows the full wire form)", async () => {
    // A four-anchor row, which is the full width of the vocabulary: the anchor
    // that resolved, the one that found nothing on the ERP side, the one with no
    // value to compare, and the one `mercado_livre` never supplies. Four full
    // sentences per row is what made the table unreadable; only the two
    // highest-ranked stay on screen.
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({
          candidate_id: "cand_confirm",
          provider_item_id: "MLB4",
          internal_product_id: 10741,
          internal_product_name: "PUXADOR FENG",
          state: "exact_sku",
          match_status: "CONFIRM",
          match_input: "seller_sku",
          confidence: 70,
          confidence_band: "MEDIA",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato para codprod" },
            { anchor: "title", direction: "INCOMPARABLE", side: "erp", detail: "produto ERP sem título comparável" },
            { anchor: "ean", direction: "UNAVAILABLE", detail: "sem EAN para corroborar o CODPROD" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    // Collapsed: the FOR anchor plus the highest-ranked remaining one, then +2.
    expect(screen.getByText("✓ SKU")).toBeInTheDocument();
    expect(screen.queryByText(/sem EAN para corroborar/)).not.toBeInTheDocument();
    const toggle = screen.getByRole("button", { name: "Mostrar todos os 4 motivos" });
    expect(toggle).toHaveTextContent("+2");

    fireEvent.click(toggle);

    // Expanded: every motivo, in its full wire form — nothing was dropped.
    for (const full of [
      "seller_sku: seller_sku resolve exato para codprod",
      "title (falta no ERP): produto ERP sem título comparável",
      "ean: sem EAN para corroborar o CODPROD",
      `marca: ${MARCA_UNAVAILABLE_DETAIL}`,
    ]) {
      expect(screen.getByText(full)).toBeInTheDocument();
    }

    // And the action stays reachable in the same row.
    expect(within(row).getByRole("button", { name: "Vincular" })).toBeInTheDocument();
  });

  // `anchorShortLabels` still carries an entry for `refforn`, an anchor D-A
  // REMOVED from the vocabulary. That entry is not dead code and it is not a
  // leftover: D-A decided that already-persisted reasons are not migrated, so a
  // row decided before the removal still holds a `refforn` motivo, and dropping
  // the label would degrade real audit data to a raw machine name.
  //
  // This is the test that makes that argument checkable. It is a drift fixture
  // rather than a wire one because the GENERATOR can no longer emit this — only
  // the database can — and `wireCandidate` is right to reject it: the previous
  // arrangement had `refforn` inside a fixture captioned "the real CONFIRM
  // shape", which asserted the opposite of the truth.
  it("keeps the label of an anchor D-A removed, for rows decided before the removal", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        driftCandidate(
          "persisted audit data, not generator output: D-A removed `refforn` from knownIdentityAnchors " +
            "without migrating reasons already written, so this row can be read from the DB and never re-emitted.",
          {
            candidate_id: "cand_legacy",
            provider_item_id: "MLB_LEGACY",
            internal_product_id: 10741,
            reasons: [
              { anchor: "refforn", direction: "UNAVAILABLE", detail: "refforn inexistente no lado provider" },
            ],
          },
        ),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("queue-row");
    // The historical anchor is typeset, not printed raw...
    expect(within(row).getByText("– Ref. forn.")).toBeInTheDocument();
    // ...and the machine name never reaches the operator.
    expect(row.textContent).not.toContain("refforn:");
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
  // This one is a DRIFT fixture, and saying so is the correction of an earlier
  // claim. `resolveIdentityAnchors` aborts generation unless the declaration
  // resolves (generation_service.go:149-169), and `identity_anchor_adapter.go`
  // walks all four anchors, so EVERY candidate of EVERY provider carries a
  // `marca` reason — UNAVAILABLE when unsupplied, INCOMPARABLE when supplied. A
  // row with only two reasons, both INCOMPARABLE, is therefore not emitted by
  // any declaration: the earlier note calling it "producible under a capability
  // declaration" was false about THIS reason array, not merely incomplete.
  //
  // It is kept, and kept unproducible, because it is the boundary of the V2
  // invariant — the extreme where a filter wearing a ranking's comment leaves
  // the cell empty. The test immediately after it drives the SAME defect from a
  // reason set today's backend really emits, which is what actually discharges
  // the criterion.
  it("keeps a motivo on screen for a row whose reasons are ALL INCOMPARABLE (ADR-17)", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        driftCandidate(
          "the all-INCOMPARABLE boundary of V2. No declaration emits it: every candidate carries a `marca` " +
            "reason (resolveIdentityAnchors :149-169 + identity_anchor_adapter.go:28-35), so a two-reason " +
            "all-INCOMPARABLE array is unreachable. Held as the extreme case of the ranking invariant.",
          {
            candidate_id: "cand_inc",
            provider_item_id: "MLB_INC",
            internal_product_id: 444,
            internal_product_name: "PUXADOR FENG",
            reasons: [
              { anchor: "seller_sku", direction: "INCOMPARABLE", side: "provider", detail: "anúncio sem seller_sku" },
              { anchor: "ean", direction: "INCOMPARABLE", side: "erp", detail: "produto ERP sem ean cadastrado" },
            ],
          },
        ),
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
          // `marca` is the UNAVAILABLE one, and it has to be: mercado_livre
          // DECLARES title supplied (capability_adapter.go:90), so a `title`
          // UNAVAILABLE carrying "provider não fornece a âncora title" — which
          // is what this fixture said before — contradicts the one declaration
          // in the tree. `marca` is the anchor that sentence is true of.
          reasons: [
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
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
    expect(screen.getByText(`marca: ${MARCA_UNAVAILABLE_DETAIL}`)).toBeInTheDocument();
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
        driftCandidate(
          "the side-less INCOMPARABLE branch (generation_service.go:706-708) fires only for an anchor a " +
            "provider DECLARES supplied but `identityAnchorValues` cannot read — `marca` is the live case, " +
            "and mercado_livre, the only adapter here, declares it unsupplied. No declaration in this tree " +
            "emits it; the SDK marks `side` optional, and the screen must not fill the gap in.",
          {
            candidate_id: "cand_noside",
            provider_item_id: "MLB_NOSIDE",
            internal_product_id: 666,
            reasons: [
              { anchor: "marca", direction: "INCOMPARABLE", detail: "não foi possível comparar a âncora marca" },
            ],
          },
        ),
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
          // NO_CANDIDATE comes only from applyUnresolvedScore, reached only
          // through newCandidate(..., Unresolved, MatchInputNone, ...), so the
          // state and the match_input follow from the status rather than being
          // free fields — and the row still carries reasons: `reasons: []` was
          // the shape of a candidate no path builds.
          state: "unresolved",
          match_status: "NO_CANDIDATE",
          match_input: "none",
          confidence: 0,
          confidence_band: "BAIXA",
          reasons: [
            { anchor: "seller_sku", direction: "INCOMPARABLE", side: "erp", detail: "seller_sku sem correspondência" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
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
            // 95/ALTA is assigned at exactly one site — buildConcordantCandidate
            // (:503-505) — and only on the ACCEPT path, so the band cannot be
            // raised without the status that earns it.
            state: "exact_ean",
            match_status: "ACCEPT",
            match_input: "ean",
            confidence: 95,
            confidence_band: "ALTA",
            reasons: [
              { anchor: "ean", direction: "FOR", detail: "EAN idêntico" },
              { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato para codprod" },
              { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
            ],
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
          state: "exact_sku",
          match_status: "ACCEPT",
          match_input: "seller_sku",
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "100%" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
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
          confidence: 95,
          confidence_band: "ALTA",
          reasons: [
            { anchor: "ean", direction: "FOR", detail: "EAN idêntico" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
        // Single anchor resolved a single product → confirmation queue, one anchor.
        candidate({
          candidate_id: "cand_confirm_sku",
          provider_item_id: "MLB_CONFIRM",
          state: "exact_sku",
          match_status: "CONFIRM",
          match_input: "seller_sku",
          confidence: 70,
          confidence_band: "MEDIA",
          reasons: [
            { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
        }),
        // Title match: a FOR reason on screen, but nothing decided.
        candidate({
          candidate_id: "cand_title",
          provider_item_id: "MLB_TITLE",
          state: "title_match",
          match_status: "REVIEW",
          match_input: "title",
          reasons: [
            { anchor: "title", direction: "FOR", detail: "match por título" },
            { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
          ],
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
        driftCandidate(
          "the API ships a sixth match_status before the SDK is regenerated. The cast is the POINT: " +
            "compile-time exhaustiveness cannot see a value the wire invents at runtime.",
          {
            candidate_id: "cand_drift",
            provider_item_id: "MLB_DRIFT",
            internal_product_id: 111,
            internal_product_name: "Parafuso A",
            match_status: "PENDING_REVIEW" as ProductLinkCandidateItem["match_status"],
            match_input: "ean",
            reasons: [{ anchor: "ean", direction: "FOR", detail: "EAN idêntico" }],
          },
        ),
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
        driftCandidate(
          "the API ships a fifth reason direction before the SDK is regenerated — the exact failure a " +
            "Record<Union, …> is blind to, one layer down from the one this chip fixes.",
          {
            candidate_id: "cand_dir",
            provider_item_id: "MLB_DIR",
            internal_product_id: 111,
            internal_product_name: "Parafuso A",
            reasons: [
              { anchor: "ean", direction: "PARTIAL" as ProductLinkCandidateItem["reasons"][number]["direction"], detail: "ean parcial" },
              { anchor: "seller_sku", direction: "FOR", detail: "seller_sku resolve exato" },
            ],
          },
        ),
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
        driftCandidate(
          "the API ships a fourth confidence_band before the SDK is regenerated. 72/CRITICA is unreachable " +
            "by construction — ALTA is assigned at one site and only with ACCEPT — which is the point: the " +
            "band the screen must survive is one no scoring path here can produce.",
          {
            candidate_id: "cand_band",
            provider_item_id: "MLB_BAND",
            internal_product_id: 111,
            internal_product_name: "Parafuso A",
            confidence: 72,
            confidence_band: "CRITICA" as ProductLinkCandidateItem["confidence_band"],
          },
        ),
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

  // The row degrading honestly is HALF the surface. The drawer renders the same
  // band from its own component, and a second copy of a lookup table is exactly
  // how the first `direction` drift survived: hardening one copy and declaring
  // the class closed. The operator reaches this by clicking the row they just
  // read, so the two surfaces disagreeing is the normal path, not a corner.
  it("survives the same unknown confidence_band in the DRAWER, not only in the row", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        driftCandidate(
          "the API ships a fourth confidence_band before the SDK is regenerated. 72/CRITICA is unreachable " +
            "by construction — ALTA is assigned at one site and only with ACCEPT — which is the point: the " +
            "band the screen must survive is one no scoring path here can produce.",
          {
            candidate_id: "cand_band",
            provider_item_id: "MLB_BAND",
            internal_product_id: 111,
            internal_product_name: "Parafuso A",
            confidence: 72,
            confidence_band: "CRITICA" as ProductLinkCandidateItem["confidence_band"],
          },
        ),
      ],
    });

    renderTab(["/?candidate=cand_band"]);

    const drawer = await screen.findByTestId("drawer-candidate");
    expect(within(drawer).getByText("CRITICA")).toBeInTheDocument();
    expect(within(drawer).getByText("72%")).toBeInTheDocument();
    expect(drawer.textContent).not.toContain("undefined");
    for (const el of Array.from(drawer.querySelectorAll<HTMLElement>("[class]"))) {
      expect(el.getAttribute("class") ?? "").not.toContain("undefined");
    }
  });

  // Injectivity has to hold on what the OPERATOR sees, not on the intermediate
  // string. HTML collapses runs of whitespace and trims the edges, so a typeset
  // form that differs only by spacing is one name on screen — and `registry.go`
  // `buildDefinitions` dedupes provider codes by exact string equality, so
  // `amazon_marketplace` and `amazon__marketplace` are both registrable at the
  // same time. Two marketplaces wearing one name is wrong information.
  it("does not let two provider codes collapse onto one name through whitespace", async () => {
    listProductLinkCandidates.mockResolvedValue({
      items: [
        candidate({ candidate_id: "cand_a", provider_item_id: "MLB_A", provider_code: "amazon_marketplace" }),
        candidate({ candidate_id: "cand_b", provider_item_id: "MLB_B", provider_code: "amazon__marketplace" }),
        candidate({ candidate_id: "cand_c", provider_item_id: "MLB_C", provider_code: "_amazon" }),
        candidate({ candidate_id: "cand_d", provider_item_id: "MLB_D", provider_code: "amazon" }),
      ],
    });

    renderTab();
    await screen.findAllByTestId("queue-row");

    // `getAllByText` normalizes whitespace exactly the way the browser paints
    // it, which is the whole point: if two codes reach the same painted string,
    // this returns two nodes for a name that should identify one provider.
    expect(screen.getAllByText("Amazon Marketplace")).toHaveLength(1);
    expect(screen.getAllByText("Amazon")).toHaveLength(1);
    // The forms that cannot be typeset injectively stay verbatim — an ugly slug
    // the operator can act on beats a pretty name that may be another provider.
    expect(screen.getByText("amazon__marketplace")).toBeInTheDocument();
    expect(screen.getByText("_amazon")).toBeInTheDocument();
  });
});
