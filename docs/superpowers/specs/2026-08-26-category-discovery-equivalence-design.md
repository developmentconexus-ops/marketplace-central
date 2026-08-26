# Category Discovery, Equivalence Memory & Organization Product View — Evaluation / Global Maximum Design

> **Status:** DESIGN FOR OPERATOR ADJUDICATION — no authority changed yet
> **Trigger:** operator-raised category/taxonomy question during B23; verification confirmed the product-linking model is decided (D2 §12, D4, D4-R1) but category discovery/equivalence is absent everywhere
> **Method:** engineering-method Global Maximum + frontend method §3.10A
> **Affected if ratified:** D4/D4-R1 (provider category evidence port) + D5 W1/OAD (new Q reads) + B23 publication-context region
> **Confirmed NOT reopened:** the no-PIM / evidence-as-stored-truth model (D2/D3/D4-R1/D7: source evidence is durably stored in MPC PostgreSQL with provenance; rejected only the editable mirror-master and the universal product table)

## 1. Evidence

- **Mercado Livre** (developers portal): hierarchical category tree per site; `domain_discovery`/**category predictor** returns the most probable category (+ its attributes) from a product title; leaf-category correctness drives which attributes exist.
- **Hubs (AnyMarket)**: category prediction plus per-marketplace category linking is a headline capability (claimed ~50% linking-effort reduction); catalog mastering is delegated to the ERP/backoffice when one exists.
- **Repository:** correspondence (product↔listing) is decided with corroboration safety (D2 §12); `publication_context` is chosen per intent; no category browse/search/suggestion operation exists in the OAD; no org-level category equivalence memory exists.

## 2. Capability dispositions proposed for adjudication

### C1 — Marketplace publication-context discovery (browse/search + source-based suggestion)

**Proposed: UPSTREAM FINDING — RATIFY.** The B23 context region currently hides a material interaction behind a fixture select; a real operator cannot pick from thousands of leaf categories unaided.

Shape (smallest sufficient):

- one new Readiness-owned Q read, `SearchPublicationContexts` (`readiness.read`, H/A/S): exact Organization + Installation (+ optional `q` text or the source-qualified product as suggestion basis) → typed candidates `{category_key, product_type_key?, display_name, path_presentation[], suggestion_basis: provider_prediction | text_search | organization_history, confidence?}` with honest known/unknown/unavailable population states, cursor-bounded;
- provider prediction (ML category predictor / domain_discovery) enters as D4 evidence through the existing Readiness publication-requirements port family; suggestion is **presentation/read only** — the chosen context still lands in the intent as canonical keys via the existing write;
- **surface impact: +1 Product operation** (106 → 107; Permissions/Principal kinds unchanged) — explicitly adjudicated, not smuggled.

### C2 — Organization category equivalence memory ("tabela de categorias")

**Proposed: UPSTREAM FINDING — RATIFY, as derived suggestion only.** When a second product of the same source category is prepared, the operator should see "sua organização publicou produtos desta categoria de origem em: Torneiras (ML) — usado em 12 anúncios".

Shape (deliberately write-free):

- the memory is **derived from decisions that already exist** (resolved intents/listings of the organization: source category evidence ↔ chosen publication context); it surfaces inside C1's candidates as `suggestion_basis: organization_history`;
- **no manual de-para editor, no equivalence write API, no automatic application** in Product 1.0 — suggestion never becomes authority (D2 §12 corroboration law extends: a history match alone never auto-selects);
- zero new operations beyond C1; zero new Permissions.

### C3 — Organization product view ("Meus produtos")

**Proposed: DEFERRED (Product reason), operator may overrule.** The entry job (find the exact product to prepare/publish) is covered by the locked B10 search; a browsable all-products projection is a convenience whose scale/pagination/refresh semantics deserve their own bounded increment after B23 closes, over the same stored evidence (D3 projection class). Deferral is recorded with a reopen trigger: if operating B10/B20 proves discovery-by-search insufficient for a real catalog size, this becomes the next ratified finding.

## 3. Rejected in all cases

MPC-owned category taxonomy (any level count) as authority; editable product mirror-master; automatic category application without human confirmation; universal cross-marketplace category ontology; provider DTO leakage past the D4 port.

## 4. Consequence for B23

C1/C2 are material to the B23 publication-context region: per the frontend-method blocking law, the B23 walkthrough stays on hold until this adjudication lands. After ratification + repair, the context region renders discovery/suggestion honestly (search, provider suggestion with confidence, organization-history suggestion with usage provenance, explicit human choice) and B23 P8 resumes.

## 5. Proof obligations (upon ratification)

OAD lint + census moves to 107/31 with the one adjudicated Q read; projection verifier gains the typed candidate/suggestion shapes and negative controls (suggestion never a write carrier, history never auto-selects, unknown/unavailable never collapse); B23 wireframe verifier gains the context-discovery states.
