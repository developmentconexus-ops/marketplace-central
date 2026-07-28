import type { ProductLinkCandidateItem, ProductLinkReason } from "@marketplace-central/sdk-runtime";

/**
 * Test-only fixture builder for `/vinculos`.
 *
 * WHY THIS EXISTS. Four reviewer rounds and two author sweeps kept finding the
 * same defect class: a fixture asserting a layout off a candidate the backend
 * cannot emit. Each occurrence was point-fixed, and each fix was followed by
 * another instance — one of them inside the file the previous fix had edited,
 * two of them found by the sweep built to end the class, and two more found by
 * the round that reviewed that sweep. The failure was never carelessness about
 * a particular field; it was the METHOD. Every check was a human reading the Go
 * generator looking for one thing, and each reading narrowed the world to what
 * it was already looking for:
 *
 *   - the first read stopped at the scoring switch and missed the finalizing step;
 *   - the sweep built to catch that checked one anchor (`marca`) and called the
 *     whole array swept, missing a removed anchor and an impossible sentence;
 *   - the enumeration run to prove the sweep exhaustive used `anchor: "[a-z_]*"`,
 *     whose character class silently dropped every capitalized and accented
 *     anchor, and reported one violation where there were five.
 *
 * So this file is not another checklist. An impossible candidate is meant to be
 * UNWRITEABLE rather than detectable: `wireCandidate` throws, in the test that
 * builds it, naming the generator line the fixture contradicts.
 *
 * WHAT IT CHECKS, and each rule's source (paths relative to
 * `apps/server_core/internal/modules/`):
 *
 *  1. Anchor vocabulary. `product_links/.../generation_service.go` only ever
 *     emits anchors from `connectors/ports/marketplace_capability.go`'s
 *     `knownIdentityAnchors`. `wireFixtures.guard.test.ts` reads that Go source
 *     and fails if this list drifts from it.
 *  2. One reason per anchor. `appendProviderDeclaredUnavailableReasons` keys its
 *     absence bookkeeping by anchor and overwrites in place (:686-696), so an
 *     anchor cannot appear twice.
 *  3. `side` only on INCOMPARABLE. `missingMatchedAnchorReason` sets a side only
 *     on the two INCOMPARABLE branches (:631-643); FOR/AGAINST/UNAVAILABLE
 *     leave it empty and `json:"side,omitempty"` drops it from the payload.
 *  4. `marca` always present. It is a declared anchor for every provider
 *     (`identity_anchor_adapter.go:28-35` walks all four), it never carries a
 *     FOR/AGAINST signal, and `classifyProviderIdentityAnchor` emits for it on
 *     both branches — UNAVAILABLE when unsupplied (:700-702), INCOMPARABLE when
 *     supplied (:706-708, since `marca` has no case in `identityAnchorValues`).
 *     `resolveIdentityAnchors` ABORTS generation when the declaration does not
 *     resolve (:149-169), so there is no candidate anywhere without it.
 *  5. For `mercado_livre` — the only capability declaration in this tree,
 *     `{seller_sku, ean, title}` at `mercado_livre/capability_adapter.go:90` —
 *     `marca` is UNAVAILABLE with that exact sentence.
 *  6. The score triple. Every `(confidence, confidence_band, match_status)` the
 *     generator can assign is enumerated below; anything else is unreachable.
 *  7. NO_CANDIDATE shape. It is emitted only by `applyUnresolvedScore`, reached
 *     only through `newCandidate(..., Unresolved, MatchInputNone, ...)` (:215,
 *     :379), so the state and match_input follow from the status.
 *
 * WHAT IT DOES NOT CHECK — stated so the next reader does not mistake a green
 * build for a proof this file cannot give:
 *
 *  - the `detail` wording of non-`marca` reasons (they are built with runtime
 *    values — codprod ids, match counts — so there is no closed set to compare);
 *  - that the seeded FOR/AGAINST set matches `state` (e.g. `exact_sku` implying
 *    a `seller_sku` FOR);
 *  - the `title` suppression condition (:709-711 needs listing and product
 *    values this side cannot see), so a MISSING `title` reason is accepted.
 *
 * Those are gaps, not guarantees. They are the honest remainder after the rules
 * that can be stated exactly.
 */

/** Mirrors `knownIdentityAnchors`; the guard test proves it still does. */
export const KNOWN_IDENTITY_ANCHORS = ["seller_sku", "ean", "title", "marca"] as const;

/** `mercado_livre` declares these three supplied; `marca` is the one it does not. */
export const MERCADO_LIVRE_SUPPLIED_ANCHORS = ["seller_sku", "ean", "title"] as const;

export const MARCA_UNAVAILABLE_DETAIL = "provider não fornece a âncora marca";

type Score = {
  confidence: number;
  band: ProductLinkCandidateItem["confidence_band"];
  status: ProductLinkCandidateItem["match_status"];
  from: string;
};

/**
 * Every score triple `generation_service.go` can assign. Lifted from the
 * assignment sites, not inferred from the band thresholds in the design doc —
 * the code is what ships.
 */
const PRODUCIBLE_SCORES: Score[] = [
  { confidence: 95, band: "ALTA", status: "ACCEPT", from: "buildConcordantCandidate :503-505" },
  { confidence: 70, band: "MEDIA", status: "CONFIRM", from: "applySingleAnchorScore exact_sku :528" },
  { confidence: 60, band: "MEDIA", status: "CONFIRM", from: "applySingleAnchorScore exact_ean :539" },
  { confidence: 40, band: "BAIXA", status: "REVIEW", from: "applyAmbiguousCorroborationScore :601-603" },
  { confidence: 35, band: "BAIXA", status: "REVIEW", from: "applySingleAnchorScore title_match :551" },
  { confidence: 25, band: "BAIXA", status: "REJECT", from: "hard-negative branch :498-500 / :566" },
  { confidence: 20, band: "BAIXA", status: "REVIEW", from: "applyConflictScore :575 / applyCollisionScore :588" },
  { confidence: 0, band: "BAIXA", status: "NO_CANDIDATE", from: "applyUnresolvedScore :620-622" },
];

function fail(rule: string, detail: string): never {
  throw new Error(
    `wireCandidate: ${rule}\n  ${detail}\n` +
      "  This fixture cannot be produced by the backend. Either correct it, or — if the\n" +
      "  unproducible shape IS the assertion (a wire-drift probe) — build it with\n" +
      "  driftCandidate(why, overrides), which records why in the test itself.",
  );
}

function assertProducibleReasons(reasons: ProductLinkReason[], providerCode: string): void {
  const seen = new Set<string>();
  for (const reason of reasons) {
    if (!(KNOWN_IDENTITY_ANCHORS as readonly string[]).includes(reason.anchor)) {
      fail(
        `anchor ${JSON.stringify(reason.anchor)} is not in the identity vocabulary.`,
        `Emitted anchors are ${KNOWN_IDENTITY_ANCHORS.join(", ")} (marketplace_capability.go knownIdentityAnchors). ` +
          "`refforn` in particular was REMOVED by D-A because the question answers `no` for every provider present and future.",
      );
    }
    if (seen.has(reason.anchor)) {
      fail(`anchor ${JSON.stringify(reason.anchor)} appears twice.`, "The finalizer keys absences by anchor and overwrites in place (:686-696).");
    }
    seen.add(reason.anchor);

    if (reason.side !== undefined && reason.direction !== "INCOMPARABLE") {
      fail(
        `${reason.anchor} carries a side on a ${reason.direction} reason.`,
        "Only the two INCOMPARABLE branches of missingMatchedAnchorReason set a side (:631-643).",
      );
    }
  }

  const marca = reasons.find((reason) => reason.anchor === "marca");
  if (!marca) {
    fail(
      "no `marca` reason.",
      "Every declared anchor without a FOR/AGAINST signal is emitted by the finalizer, `marca` never carries one, " +
        "and resolveIdentityAnchors aborts generation when the declaration does not resolve (:149-169). " +
        "Every candidate of every provider carries exactly one.",
    );
  }
  if (providerCode === "mercado_livre") {
    if (marca.direction !== "UNAVAILABLE" || marca.detail !== MARCA_UNAVAILABLE_DETAIL) {
      fail(
        `mercado_livre's \`marca\` is ${marca.direction} ${JSON.stringify(marca.detail)}.`,
        `It declares ${MERCADO_LIVRE_SUPPLIED_ANCHORS.join("/")} (capability_adapter.go:90), so marca is unsupplied and ` +
          `classifyProviderIdentityAnchor returns UNAVAILABLE with ${JSON.stringify(MARCA_UNAVAILABLE_DETAIL)} (:700-702).`,
      );
    }
  }
}

function assertProducibleScore(item: ProductLinkCandidateItem): void {
  const match = PRODUCIBLE_SCORES.find(
    (score) =>
      score.confidence === item.confidence && score.band === item.confidence_band && score.status === item.match_status,
  );
  if (!match) {
    fail(
      `no scoring path assigns (${item.confidence}, ${item.confidence_band}, ${item.match_status}).`,
      "Producible triples: " + PRODUCIBLE_SCORES.map((s) => `(${s.confidence}, ${s.band}, ${s.status}) ${s.from}`).join("; "),
    );
  }
  if (item.match_status === "NO_CANDIDATE" && (item.state !== "unresolved" || item.match_input !== "none")) {
    fail(
      `NO_CANDIDATE with state=${item.state} match_input=${item.match_input}.`,
      "applyUnresolvedScore is the only path emitting it, reached through newCandidate(..., Unresolved, MatchInputNone, ...) (:215, :379).",
    );
  }
}

/**
 * The default is the `title_match` row, verbatim from `applySingleAnchorScore`
 * (:550-557): the title FOR that ranks but never accepts, the two anchors that
 * found nothing, and `marca`. A default has to be producible too — the previous
 * shared default was `(50, MEDIA, REVIEW)` with `reasons: []`, and MEDIA is only
 * ever assigned with CONFIRM, so every fixture that inherited it was impossible
 * before its own overrides were even applied.
 */
function defaultCandidate(): ProductLinkCandidateItem {
  return {
    candidate_id: "cand_x",
    installation_id: "inst_1",
    provider_code: "mercado_livre",
    provider_item_id: "MLB_X",
    state: "title_match",
    match_input: "title",
    confidence: 35,
    confidence_band: "BAIXA",
    match_status: "REVIEW",
    reasons: [
      { anchor: "title", direction: "FOR", detail: "match por título (ranking-only, nunca ACCEPT)" },
      { anchor: "seller_sku", direction: "INCOMPARABLE", side: "erp", detail: "seller_sku sem correspondência" },
      { anchor: "ean", direction: "INCOMPARABLE", side: "erp", detail: "ean sem correspondência" },
      { anchor: "marca", direction: "UNAVAILABLE", detail: MARCA_UNAVAILABLE_DETAIL },
    ],
    created_at: "2026-07-18T12:00:00Z",
    updated_at: "2026-07-18T12:00:00Z",
  };
}

/** A candidate the backend can actually emit. Throws if it is not one. */
export function wireCandidate(overrides: Partial<ProductLinkCandidateItem> = {}): ProductLinkCandidateItem {
  const item = { ...defaultCandidate(), ...overrides };
  assertProducibleReasons(item.reasons, item.provider_code);
  assertProducibleScore(item);
  return item;
}

/**
 * A candidate THIS CHECKOUT's backend cannot emit, on purpose.
 *
 * Two distinct legitimate uses, and the distinction is worth keeping because
 * they expire differently:
 *
 *  1. WIRE DRIFT — the API ships a union member before the SDK is regenerated.
 *     This is precisely the failure a `Record<Union, …>` cannot see, so the
 *     unproducible shape IS the assertion; validating it away deletes the test.
 *     It never expires: the SDK is always regenerated after the API moves.
 *  2. NO DECLARATION HERE — the shape follows from the generator, but only under
 *     a capability declaration no adapter in this tree makes. `mercado_livre` is
 *     the only one, so every `marca`-supplied branch of the generator is in this
 *     class. It expires the day a second adapter lands, and the fixture becomes
 *     a `wireCandidate`.
 *
 * `why` is required and unused at runtime by design: it forces the reason to be
 * written next to the fixture instead of assumed. The rule it enforces is that
 * "this cannot occur" must be a deliberate, stated choice — the fixtures that
 * started this were unproducible SILENTLY, and read as wire rows. A previous
 * round tried to discharge case 2 with a comment claiming the row was
 * "producible under a capability declaration"; the reviewer's answer was that
 * this is a dodge rather than an honest degrade, because it asserts producible
 * while the finalizer says otherwise. Saying "not producible here, and here is
 * what would have to change" is the honest form of the same fact.
 */
export function driftCandidate(why: string, overrides: Partial<ProductLinkCandidateItem> = {}): ProductLinkCandidateItem {
  if (why.trim().length === 0) {
    throw new Error("driftCandidate: `why` is required — say what drift this fixture stands for.");
  }
  return { ...defaultCandidate(), ...overrides };
}
