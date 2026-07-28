import type { ProductLinkCandidateItem, ProductLinkReason } from "@marketplace-central/sdk-runtime";

/**
 * Test-only fixture builder for `/vinculos`.
 *
 * WHY THIS EXISTS. Five reviewer rounds and two author sweeps kept finding the
 * same defect class: a fixture asserting a layout off a candidate the backend
 * cannot emit. Each occurrence was point-fixed, and each fix was followed by
 * another instance — one of them inside the file the previous fix had edited,
 * two of them found by the sweep built to end the class, two more found by the
 * round that reviewed that sweep, and nine more found by the round that reviewed
 * THIS FILE. The failure was never carelessness about a particular field; it was
 * the METHOD. Every check was a human reading the Go generator looking for one
 * thing, and each reading narrowed the world to what it was already looking for:
 *
 *   - the first read stopped at the scoring switch and missed the finalizing step;
 *   - the sweep built to catch that checked one anchor (`marca`) and called the
 *     whole array swept, missing a removed anchor and an impossible sentence;
 *   - the enumeration run to prove the sweep exhaustive used `anchor: "[a-z_]*"`,
 *     whose character class silently dropped every capitalized and accented
 *     anchor, and reported one violation where there were five;
 *   - and the first version of THIS FILE validated the score triple only, so a
 *     `(exact_ean, ACCEPT, ean)` fixture — a combination no site can produce —
 *     passed the constructor that documented itself as making it unwriteable.
 *
 * That last one is the reason for the shape below. A PARTIAL tuple is a partial
 * guard, and a partial guard under a total sentence is worse than no guard: it
 * spends the reader's trust on a claim it does not discharge. So the table now
 * carries the WHOLE producible tuple per producing site, read one site at a time
 * out of `generation_service.go` — `(confidence, band, status, state,
 * match_input)` — and the constructor matches all five together. Every rule that
 * used to be a hand-written special case is a row: the `NO_CANDIDATE` shape
 * check and the `mercado_livre` provider check are both gone, dissolved into
 * data. A fixture fixed by hand comes back next round; a table does not.
 *
 * WHAT IT CHECKS, and each rule's source (paths relative to
 * `apps/server_core/internal/modules/`):
 *
 *  1. Anchor vocabulary. `product_links/.../generation_service.go` only ever
 *     emits anchors from `connectors/ports/marketplace_capability.go`'s
 *     `knownIdentityAnchors`. `wireFixtures.guard.test.ts` reads that Go source
 *     and fails if this list drifts from it.
 *  2. One ABSENCE per anchor, and never an absence beside a signal on the same
 *     anchor. `appendProviderDeclaredUnavailableReasons` keys its absence
 *     bookkeeping by anchor and overwrites in place (:666-669, :687-697), and
 *     drops any absence whose anchor already carries a FOR/AGAINST (:652-655,
 *     :663-664). Signals themselves are NOT deduped — `title` FOR and `title`
 *     AGAINST ride the same candidate out of the hard-negative path.
 *  3. `side` only on INCOMPARABLE. `classifyProviderIdentityAnchor` sets a side
 *     only on INCOMPARABLE branches (:715, :723, :726, :728); UNAVAILABLE (:704)
 *     and the `!readable` INCOMPARABLE (:711) leave it empty, and
 *     `json:"side,omitempty"` drops it from the payload.
 *  4. Provider capability. `resolveIdentityAnchors` ABORTS generation when a
 *     provider's declaration does not resolve (:149-169), so a provider with no
 *     declaration in this tree produces NO candidate at all — not a candidate
 *     with defaults. `DECLARED_PROVIDER_CAPABILITIES` below is the whole set of
 *     declarations that exist, and a `provider_code` absent from it is rejected.
 *  5. Every anchor the declaring provider does NOT supply is UNAVAILABLE with
 *     exactly `provider não fornece a âncora <anchor>` — `classifyProviderIdentityAnchor`
 *     returns that sentence by `fmt.Sprintf` over the anchor name (:704), so the
 *     rule is per-anchor and not a `marca` special case.
 *  6. The producible tuple AND its signal set. Every
 *     `(confidence, band, status, state, match_input)` the generator can assign
 *     is enumerated below, one row per producing site, each row carrying the
 *     FOR/AGAINST reasons that site writes; anything else is unreachable. The
 *     signal set is part of the row because the same block that assigns the
 *     score writes the signals — a tuple alone does not identify the site when
 *     two sites share it. This subsumes the old NO_CANDIDATE shape
 *     rule, which was also WRONG: it demanded `state === "unresolved"`, and
 *     `buildConflictCandidates` emits NO_CANDIDATE with `state === "conflict"`
 *     when the conflict set comes back empty (:340-341). The old check would
 *     have rejected a candidate the backend really does emit.
 *
 * WHAT IT DOES NOT CHECK — stated so the next reader does not mistake a green
 * build for a proof this file cannot give:
 *
 *  - the `detail` wording of non-declared-absence reasons (they are built with
 *    runtime values — codprod ids, match counts — so there is no closed set);
 *  - the DIRECTION and `side` of an absence on a SUPPLIED anchor. Rule 6 pins
 *    which anchors carry signals; which of INCOMPARABLE/UNAVAILABLE the
 *    finalizer then picks for the rest, and with which side, depends on listing
 *    and product values no fixture here carries (:631-645, :702-729). Note in
 *    particular that UNAVAILABLE is NOT exclusive to the unsupplied-anchor
 *    branch: `missingMatchedAnchorReason:642` emits it for a supplied anchor on
 *    its `default:` arm;
 *  - the `title` suppression condition (:719-721 needs listing and product values
 *    this side cannot see), so a MISSING `title` reason is accepted.
 *
 * Those are gaps, not guarantees. They are the honest remainder after the rules
 * that can be stated exactly.
 */

/** Mirrors `knownIdentityAnchors`; the guard test proves it still does. */
export const KNOWN_IDENTITY_ANCHORS = ["seller_sku", "ean", "title", "marca"] as const;

/** `mercado_livre` declares these three supplied; `marca` is the one it does not. */
export const MERCADO_LIVRE_SUPPLIED_ANCHORS = ["seller_sku", "ean", "title"] as const;

/**
 * The sentence `classifyProviderIdentityAnchor` emits for a declared-but-unsupplied
 * anchor, `fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor)` (:704).
 * Per-anchor by construction — the previous constant was `marca`-only, which is
 * why the rule that used it could only ever check one anchor.
 */
export function unavailableDetail(anchor: string): string {
  return `provider não fornece a âncora ${anchor}`;
}

/** Kept for the fixtures that name it; it is `unavailableDetail("marca")`. */
export const MARCA_UNAVAILABLE_DETAIL = unavailableDetail("marca");

/**
 * Every capability declaration that EXISTS in this tree. Not a list of providers
 * the screen supports — a list of adapters that declare `IdentityAnchors`, which
 * is the precondition `resolveIdentityAnchors` (:149-169) enforces before any
 * candidate is generated at all. A provider missing here is not "unchecked"; it
 * is a provider whose generation aborts, so no wire row of it can exist.
 */
export const DECLARED_PROVIDER_CAPABILITIES: Record<string, { supplied: readonly string[]; from: string }> = {
  mercado_livre: {
    supplied: MERCADO_LIVRE_SUPPLIED_ANCHORS,
    from: "connectors/adapters/mercado_livre/capability_adapter.go:90",
  },
};

/**
 * One FOR/AGAINST reason a producing site emits. These are the SIGNALS: the
 * reasons the site writes itself, as opposed to the absences the finalizer adds
 * for anchors that carry no signal.
 *
 * The split is the finalizer's, not a convenience: `:652-655` builds `hasSignal`
 * from FOR/AGAINST only, and `:663-664` then DROPS every absence reason whose
 * anchor has a signal. So the signal set is fixed by the scoring site and the
 * absence set is whatever is left of the declared anchors.
 */
type ReasonSignal = {
  anchor: string;
  direction: "FOR" | "AGAINST";
};

type ProducibleSite = {
  confidence: number;
  band: ProductLinkCandidateItem["confidence_band"];
  status: ProductLinkCandidateItem["match_status"];
  state: ProductLinkCandidateItem["state"];
  matchInput: ProductLinkCandidateItem["match_input"];
  signals: readonly ReasonSignal[];
  from: string;
};

const FOR = "FOR" as const;
const AGAINST = "AGAINST" as const;

/**
 * Every `(confidence, band, status, state, match_input)` `generation_service.go`
 * can assign, one row per PRODUCING SITE — the `newCandidate(...)` call that
 * fixes `state`/`match_input` paired with the scoring call that fixes the triple.
 * Lifted from the assignment sites, never inferred from the band thresholds in
 * the design doc: the code is what ships.
 *
 * Why site-paired rather than two independent lists: `state` and `match_input`
 * are decided at `newCandidate` and the triple is decided afterwards by the
 * scorer, so a fixture can be built from a real state and a real triple and
 * still be impossible — which is exactly the ACCEPT fixture that survived five
 * rounds. `ACCEPT` is assigned at ONE site (`buildConcordantCandidate` :505) and
 * that site's `newCandidate` (:491) hardcodes `ExactSKU`/`SellerSKU`, so the
 * pairing is the fact and the two lists were the fiction.
 *
 * `signals` extends the same argument one step further, and for the same reason.
 * Validating the tuple while leaving `reasons[]` free let two impossible
 * fixtures through the round-6 guard LEGALLY — a `title/AGAINST` under
 * `60/MEDIA/CONFIRM`, when the only site that appends `title/AGAINST` (:562)
 * rewrites the triple to `25/BAIXA/REJECT` in the same statement, and a
 * `20/BAIXA/REVIEW` with no signal at all, when both sites that assign that
 * triple emit one or two. Score and reasons are not independent: at every
 * scoring site they are written by the same block, so a table that constrains
 * one and not the other is a partial guard again, just one level down.
 *
 * NOTE the same tuple can appear on more than one row. `25/BAIXA/REJECT` with
 * `exact_sku`/`seller_sku` is reachable both from the concordant hard-negative
 * (:497-501, which already carried an `ean` FOR) and from the single-anchor
 * hard-negative (:560-562, which never did). The tuple does not distinguish
 * them; the signals do, so producibility is "some site matches BOTH", never
 * "some site matches the tuple".
 */
const PRODUCIBLE_SITES: ProducibleSite[] = [
  // buildConcordantCandidate — the only ACCEPT in the tree. Both anchors resolve
  // the same product, so both are FOR (:493-496).
  { confidence: 95, band: "ALTA", status: "ACCEPT", state: "exact_sku", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }, { anchor: "ean", direction: FOR }],
    from: "buildConcordantCandidate :491 + :493-496 + :503-505" },
  { confidence: 25, band: "BAIXA", status: "REJECT", state: "exact_sku", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }, { anchor: "ean", direction: FOR }, { anchor: "title", direction: AGAINST }],
    from: "buildConcordantCandidate hard-negative :491 + :497-501" },

  // buildCandidatesFromProducts (:389) → applySingleAnchorScore (:390). The
  // switch keys off the same `state` the candidate was built with, so state and
  // match_input travel together from :212 / :311 / :314. One anchor resolved, so
  // one FOR; the anchors that did not are ABSENCES, not signals, and the
  // finalizer decides their direction.
  { confidence: 70, band: "MEDIA", status: "CONFIRM", state: "exact_sku", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }], from: "applySingleAnchorScore ExactSKU :529 + :534-537" },
  { confidence: 60, band: "MEDIA", status: "CONFIRM", state: "exact_ean", matchInput: "ean",
    signals: [{ anchor: "ean", direction: FOR }], from: "applySingleAnchorScore ExactEAN :539 + :544-547" },
  { confidence: 35, band: "BAIXA", status: "REVIEW", state: "title_match", matchInput: "title",
    signals: [{ anchor: "title", direction: FOR }], from: "applySingleAnchorScore TitleMatch :549 + :550-554" },
  // The hard-negative override at :560-562 replaces the triple of any of the
  // three above, leaves state/match_input untouched, and APPENDS `title AGAINST`
  // to whatever signals that case had already written. On the TitleMatch case
  // that means `title` carries a FOR and an AGAINST at once — the finalizer
  // appends both (:657-661) and dedups only absences.
  { confidence: 25, band: "BAIXA", status: "REJECT", state: "exact_sku", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }, { anchor: "title", direction: AGAINST }],
    from: "applySingleAnchorScore ExactSKU + hard-negative :534-537 + :560-562" },
  { confidence: 25, band: "BAIXA", status: "REJECT", state: "exact_ean", matchInput: "ean",
    signals: [{ anchor: "ean", direction: FOR }, { anchor: "title", direction: AGAINST }],
    from: "applySingleAnchorScore ExactEAN + hard-negative :544-547 + :560-562" },
  { confidence: 25, band: "BAIXA", status: "REJECT", state: "title_match", matchInput: "title",
    signals: [{ anchor: "title", direction: FOR }, { anchor: "title", direction: AGAINST }],
    from: "applySingleAnchorScore TitleMatch + hard-negative :550-554 + :560-562" },

  // buildConflictCandidates (:335) → applyConflictScore (:575-577). `ownAnchor`
  // and `matchInput` are assigned together in the same branch (:322-333), so the
  // FOR is always the match_input anchor and the AGAINST is always the other of
  // the two — the conflict IS between seller_sku and ean.
  { confidence: 20, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }, { anchor: "ean", direction: AGAINST }],
    from: "applyConflictScore :322-336 + :575-581" },
  { confidence: 20, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "ean",
    signals: [{ anchor: "ean", direction: FOR }, { anchor: "seller_sku", direction: AGAINST }],
    from: "applyConflictScore :322-336 + :575-581" },

  // buildCollisionCandidates (:369) → applyCollisionScore (:588-590) when the
  // anchor matched >1 product, applyAmbiguousCorroborationScore (:602-604) when
  // it matched exactly one. Both are built with StateConflict, and `anchor.name`
  // travels with `anchor.matchInput` out of the same struct literal (:360-362).
  // The collision site emits a LONE AGAINST: the anchor named a set, so it
  // corroborates nothing, and there is no FOR to pair it with.
  { confidence: 20, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: AGAINST }], from: "applyCollisionScore :369-371 + :588-593" },
  { confidence: 20, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "ean",
    signals: [{ anchor: "ean", direction: AGAINST }], from: "applyCollisionScore :369-371 + :588-593" },
  { confidence: 40, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "seller_sku",
    signals: [{ anchor: "seller_sku", direction: FOR }, { anchor: "ean", direction: AGAINST }],
    from: "applyAmbiguousCorroborationScore :369-373 + :602-608" },
  { confidence: 40, band: "BAIXA", status: "REVIEW", state: "conflict", matchInput: "ean",
    signals: [{ anchor: "ean", direction: FOR }, { anchor: "seller_sku", direction: AGAINST }],
    from: "applyAmbiguousCorroborationScore :369-373 + :602-608" },

  // applyUnresolvedScore (:621-623) is reached from THREE sites, and they do not
  // agree on `state` — which is the defect the old NO_CANDIDATE special case had.
  // It writes NO signal: both seeded reasons are absences (:624-627), which is
  // the whole point — nothing resolved, so nothing argues either way.
  { confidence: 0, band: "BAIXA", status: "NO_CANDIDATE", state: "unresolved", matchInput: "none",
    signals: [], from: "applyUnresolvedScore :215-216 / :379-380 + :620-627" },
  { confidence: 0, band: "BAIXA", status: "NO_CANDIDATE", state: "conflict", matchInput: "none",
    signals: [], from: "applyUnresolvedScore via conflict set empty :340-341 — state stays conflict" },
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
  const capability = DECLARED_PROVIDER_CAPABILITIES[providerCode];
  if (!capability) {
    fail(
      `provider ${JSON.stringify(providerCode)} declares no marketplace capability in this tree.`,
      `resolveIdentityAnchors aborts generation when a provider's declaration does not resolve (:149-169), so this ` +
        `provider emits no candidate at all. Declared here: ${Object.keys(DECLARED_PROVIDER_CAPABILITIES).join(", ")}. ` +
        "A fixture about how the SCREEN renders an unknown provider is legitimate — it is driftCandidate's case 2.",
    );
  }

  const signalAnchors = new Set(
    reasons.filter((reason) => reason.direction === "FOR" || reason.direction === "AGAINST").map((reason) => reason.anchor),
  );
  const seenAbsence = new Set<string>();
  for (const reason of reasons) {
    if (!(KNOWN_IDENTITY_ANCHORS as readonly string[]).includes(reason.anchor)) {
      fail(
        `anchor ${JSON.stringify(reason.anchor)} is not in the identity vocabulary.`,
        `Emitted anchors are ${KNOWN_IDENTITY_ANCHORS.join(", ")} (marketplace_capability.go knownIdentityAnchors). ` +
          "`refforn` in particular was REMOVED by D-A because the question answers `no` for every provider present and future.",
      );
    }
    // `side` binds on EVERY direction, so it is checked before the split below.
    // Putting the split first would have silently exempted FOR/AGAINST from it.
    if (reason.side !== undefined && reason.direction !== "INCOMPARABLE") {
      fail(
        `${reason.anchor} carries a side on a ${reason.direction} reason.`,
        "Only INCOMPARABLE branches set a side (:637, :640, :715, :723, :726, :728); UNAVAILABLE leaves it empty at both of its " +
          "sites (:642 and :704), and no FOR/AGAINST site writes the field at all.",
      );
    }

    // Absences dedup by anchor; SIGNALS do not. The previous form of this rule
    // rejected any repeated anchor, which is wider than the fact: the
    // TitleMatch hard-negative emits `title` FOR and `title` AGAINST on the same
    // candidate, because :560-562 appends to the signals :550-554 already wrote
    // and the finalizer passes both through untouched (:657-661). A guard that
    // cannot express that site cannot check it either.
    if (reason.direction === "FOR" || reason.direction === "AGAINST") continue;

    if (seenAbsence.has(reason.anchor)) {
      fail(
        `anchor ${JSON.stringify(reason.anchor)} carries two absence reasons.`,
        "The finalizer keys absences by anchor in `absenceIndexes` and overwrites in place (:666-669, :687-697), so at most one survives per anchor.",
      );
    }
    seenAbsence.add(reason.anchor);

    if (signalAnchors.has(reason.anchor)) {
      fail(
        `${reason.anchor} carries both a signal and a ${reason.direction} absence.`,
        "`hasSignal` is built from the FOR/AGAINST reasons (:652-655) and every absence for an anchor in it is DROPPED (:663-664), " +
          "so no emitted candidate pairs the two on one anchor.",
      );
    }

  }

  // Every anchor the provider DECLARES but does not SUPPLY is emitted, with a
  // sentence fixed by the anchor's name. Per-anchor, so a second adapter with a
  // different supplied set is checked by the same rule instead of falling
  // through it — the previous form asked only about `mercado_livre`, which
  // exempted exactly the providers whose producibility is unknown.
  for (const anchor of KNOWN_IDENTITY_ANCHORS) {
    if (capability.supplied.includes(anchor)) continue;
    const reason = reasons.find((candidateReason) => candidateReason.anchor === anchor);
    if (!reason) {
      fail(
        `no \`${anchor}\` reason.`,
        `${providerCode} declares ${capability.supplied.join("/")} supplied (${capability.from}), so ${anchor} is ` +
          "declared-and-unsupplied and the finalizer emits it on every candidate. resolveIdentityAnchors aborts " +
          "generation when the declaration does not resolve (:149-169), so there is no candidate anywhere without it.",
      );
    }
    if (reason.direction !== "UNAVAILABLE" || reason.detail !== unavailableDetail(anchor)) {
      fail(
        `${providerCode}'s \`${anchor}\` is ${reason.direction} ${JSON.stringify(reason.detail)}.`,
        `${anchor} is unsupplied per ${capability.from}, so classifyProviderIdentityAnchor returns UNAVAILABLE with ` +
          `${JSON.stringify(unavailableDetail(anchor))} (:704).`,
      );
    }
  }
}

/** Order-insensitive, multiplicity-sensitive — `title` FOR + `title` AGAINST is a set of two. */
function signalKeys(signals: readonly ReasonSignal[]): string[] {
  return signals.map((signal) => `${signal.anchor}/${signal.direction}`).sort();
}

function sameSignals(a: readonly ReasonSignal[], b: readonly ReasonSignal[]): boolean {
  const left = signalKeys(a);
  const right = signalKeys(b);
  return left.length === right.length && left.every((key, index) => key === right[index]);
}

function assertProducibleCandidate(item: ProductLinkCandidateItem): void {
  const signals = item.reasons
    .filter((reason) => reason.direction === "FOR" || reason.direction === "AGAINST")
    .map((reason) => ({ anchor: reason.anchor, direction: reason.direction as "FOR" | "AGAINST" }));

  const tupleMatches = PRODUCIBLE_SITES.filter(
    (site) =>
      site.confidence === item.confidence &&
      site.band === item.confidence_band &&
      site.status === item.match_status &&
      site.state === item.state &&
      site.matchInput === item.match_input,
  );

  if (tupleMatches.length === 0) {
    fail(
      `no producing site emits (${item.confidence}, ${item.confidence_band}, ${item.match_status}, ` +
        `state=${item.state}, match_input=${item.match_input}).`,
      "Producible tuples: " +
        PRODUCIBLE_SITES.map(
          (site) => `(${site.confidence}, ${site.band}, ${site.status}, ${site.state}, ${site.matchInput}) ${site.from}`,
        ).join("; "),
    );
  }

  if (!tupleMatches.some((site) => sameSignals(site.signals, signals))) {
    fail(
      `no producing site emits (${item.confidence}, ${item.confidence_band}, ${item.match_status}, ` +
        `state=${item.state}, match_input=${item.match_input}) WITH signals [${signalKeys(signals).join(", ") || "none"}].`,
      `That tuple is producible, but only from ${tupleMatches.length === 1 ? "one site" : `${tupleMatches.length} sites`}, and the ` +
        "FOR/AGAINST reasons are written by the same block that assigns the score, so they are not free to vary: " +
        tupleMatches.map((site) => `[${signalKeys(site.signals).join(", ") || "none"}] ${site.from}`).join("; "),
    );
  }
}

/**
 * The default is the `title_match` row, verbatim from `applySingleAnchorScore`
 * (:548-554): the title FOR that ranks but never accepts, the two anchors that
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

/**
 * A candidate the backend can actually emit.
 *
 * Scope, stated rather than implied — the guarantee is TOTAL over what it names
 * and silent about the rest, and the naming is the point (§11, `f19793d7`: a
 * total guarantee in a docstring is a claim about every input). It throws unless
 * ALL of:
 *
 *  - the provider has a capability declaration in this tree, and every anchor it
 *    declares-but-does-not-supply carries the UNAVAILABLE sentence for that
 *    anchor;
 *  - every reason's anchor is in the identity vocabulary, no anchor carries two
 *    absences, no anchor carries a signal and an absence at once, and no
 *    non-INCOMPARABLE reason carries a `side`;
 *  - `(confidence, band, status, state, match_input)` is a tuple some producing
 *    site assigns, AND that site's FOR/AGAINST signal set is exactly this
 *    fixture's.
 *
 * What it does NOT check, and what would therefore pass: the `detail` sentences
 * of signals, the direction/side of an absence on a SUPPLIED anchor (the
 * finalizer picks those from listing/product values this fixture does not carry),
 * and anything about fields outside the candidate.
 */
export function wireCandidate(overrides: Partial<ProductLinkCandidateItem> = {}): ProductLinkCandidateItem {
  const item = { ...defaultCandidate(), ...overrides };
  assertProducibleReasons(item.reasons, item.provider_code);
  assertProducibleCandidate(item);
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
 *     class, and so is every candidate of any other provider_code: a provider
 *     with no declaration aborts generation (:149-169) rather than emitting
 *     defaults. It expires the day a second adapter lands, the adapter gets a
 *     row in `DECLARED_PROVIDER_CAPABILITIES`, and the fixture becomes a
 *     `wireCandidate`.
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
