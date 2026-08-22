# D8 — Live Probe Evidence Record

> **Status:** EXECUTED / RECORDED — 2026-08-22
> **Protocol:** [D8 Controlled Live Probe Protocol](D8-LIVE-PROBE-PROTOCOL.md)
> **Parent ledger:** [D8-R1 §3.1](D8-R1-PROOF-CLOSURE-COHERENCE.md)
> **Candidate SHA at execution:** `6ed73bc007eed12bb04e8d8d8e5490e4123632d6`
> **Executed from:** credentialed operator environment (local repository `.env` + operator Sankhya Gateway credential file + encrypted Installation credential store). No secret material appears in this record.

## Common preflight

```text
executed_at (UTC)        2026-08-22 13:45–14:15
organization/tenant      opaque single-tenant mirror store (docker Postgres)
marketplace installation inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0
credential binding       GET /users/me id = 691607102 == installation external_account_id (match)
credential rotation      expired oauth2 credential v50 refreshed once; rotated token persisted as v51
                         BEFORE first authenticated use (single-use refresh preserved)
sankhya transport        sanctioned API Gateway only (client_credentials + X-Token, Bearer TTL 300s);
                         Direct Oracle not used
docs recheck             current ML sync/modify doc (es_AR mirror) reverified immediately before P1:
                         available_quantity writable via PUT /items (0 ⇒ paused/out_of_stock);
                         price-only PUT rejection / price-ignored-with-warning rule tied to active
                         pricing automation; SelecaoDocumentoSP.faturar (MGECOM) request shape
                         reverified against current developer.sankhya.com.br reference
operator context         TOP 313 currently used as TEST lane for in-progress own-e-commerce
                         implementation (operator statement 2026-08-22); real reference flows are the
                         early-August 313→306 progressions
```

---

## P1 — Mercado Livre Price/Availability effect — `EXECUTED_AND_RECORDED`

**Verdict: `PASS_CONVERGED` (both meanings, separately).**

Fixture: existing active seller-owned Item `MLB4834219830` ↔ User Product `MLBU783470824` (1:1, no variations, no `inventory_id`, non-catalog, `me2 xd_drop_off`, `gold_pro`).

Preflight (read-only):

- item: `active`, price 29.90 BRL, `available_quantity` 3, `sold_quantity` 1;
- UP stock: single location `selling_address` qty 3, `stock_mode countable`, `x-version: 1`;
- pricing automation: `404 automation_not_found` (write not automation-controlled);
- promotions: all `candidate` status only — no active promotion constraint.

Availability (one write): `PUT /items/MLB4834219830 {"available_quantity": 2}` → 200. Authoritative reread: item qty 2 **and** UP stock surface qty 2 with `x-version: 1→2`. Blast radius: single Item↔UP, single location.

Price (one write, deliberately price-only per protocol §5.3): `PUT /items/MLB4834219830 {"price": 30.4}` → 200, no price warning (only unrelated `shipping.lost_me1_by_user` notices). Authoritative reread: item price/base_price 30.4 **and** buyer-facing `sale_price` surface amount 30.4.

Restoration (fresh effects): price→29.90 then qty→3; rereads converged (UP stock 3; sale_price 29.90). `sold_quantity` stayed 1 throughout — no intervening sale; no pause/out-of-stock transition occurred.

Architecture implication: the accepted D4/D4-R1 Item-path write lane for the selected first lane (UP model, single `selling_address`, no multi-origin) controls **both** Offering Price meaning and Availability meaning with authoritative-reread convergence. No reopen.

---

## P2 — selected-lane fiscal / invoice / label progression (Mercado Livre) — `OPERATOR_RATIFIED_REDEFER`

**Disposition: `OPERATOR_RATIFIED_REDEFER(first real open Mercado Livre Sale — execute the invoice_pending→label sequence on the next real sale, or during Product implementation behind an explicit beta flag).**

Facts: at execution time the seller had **no open Sale/Shipment** — newest real Sale 2026-08-11 (order `2000017883544496`), and the three newest shipments (`47745061741`, `47706254876`, `47668318460`) all reread `status=delivered`. Protocol law 6 forbids fabricating a buyer transaction and forbids re-invoicing an already-materialized sale; therefore no qualifying fixture existed.

Operator ratification (2026-08-22, verbatim intent): test when a real sale happens or after implementation with the feature flagged beta.

Real-flow correlation evidence retained: early-August real 313→306 progressions exist with série 1 whose values match the real ML sales 1:1 (1719.80 / 209.70 / 169.99), so the production lane itself is live and observed — only the MPC-driven first controlled execution is redeferred.

---

## P3 — first irreversible Sankhya `313 → 306` progression — `EXECUTED_AND_RECORDED`

**Verdict: `PASS_CONVERGED` (bounded: generated nota deliberately left unconfirmed).**

Fixture: TEST e-commerce order `NUNOTA 898263` (TOP 313, `STATUSNOTA=L` confirmed, `PENDENTE=S`, R$122.80, 1 item — CODPROD 12475 × 1 @ 90.40, negotiation type 27, CODEMP 1, partner CODPARC 142907 = PF consumidor final, interstate MG→SP). Operator explicitly authorized execution on this test order **with the constraint that the generated 306 must NOT be confirmed** (confirmation would emit a real fiscal document).

Execution:

1. authoritative reread immediately before dispatch (order L/pending, contact bound — see P5);
2. first dispatch attempt without `serie` → definitive validation rejection `CORE_E02938` "Série '' não pode ser usada com a TOP 306" (no state change; rejection-before-effect does not burn the one-consequential-attempt law);
3. série established from real early-August 306s (`SERIENOTA=1`), one corrected dispatch:
   `SelecaoDocumentoSP.faturar` (MGECOM Gateway) — TOP 306, série 1, `FaturamentoNormal`, dtFaturamento 22/08/2026, single nota 898263 → **accepted**, generated `NUNOTA 898877`, `vlrNotaFat 122.8`, `numNotaFat 0`;
4. **no confirmation service was ever called.**

Authoritative rereads:

- source 898263: `PENDENTE S→N`, `STATUSNOTA=L` unchanged, TOP 313 intact; source item `PENDENTE S→N` (atendido);
- result 898877: distinct identity, TOP 306, `TIPMOV=V`, **`STATUSNOTA=A` (open / NOT confirmed)**, `NUMNOTA=0` (no fiscal number), same partner, item correlation exact (1 × 12475 @ 90.40), `ATUALESTOQUE=-1`, `RESERVA=N` (reservation released into the nota);
- blast radius: the other 10 pending TOP-313 orders reread byte-identical (including the accepted Expected-Tax model order 898307, untouched).

Architecture implication: the accepted D4 B3 lane (MGECOM `SelecaoDocumentoSP.faturar`, 313→306, série/TOP-bound, native NUNOTA as source-qualified external result, no blind retry) is **realizable against the production Gateway**. Bounded residual, explicitly not hidden: nota confirmation / fiscal-document emission (NF-e) remains unexercised by operator constraint and travels with the P2 redefer gate. No reopen.

---

## P4 — native Party create/update — `NOT_TRIGGERED`

Evidence: the selected Sale fixture was already bound to exactly one existing compatible native Party (`CODPARC 142907`, active, cliente) before any MPC action; no zero-match condition existed and no identity-bearing creation evidence was required. Per D8-R1 §3.2 the conditional row records `NOT_TRIGGERED` with these facts. No Party was created or mutated (P5 used the contact sub-entity, not Partner master fields).

---

## P5 — alternate destination/contact realization — `EXECUTED_AND_RECORDED`

**Verdict: `CAPABILITY_NOT_PROVEN` for full alternate-destination realization; contact-reference realization itself converged.**

Preflight: partner 142907 had **zero** contacts; partner master CEP differs from the probe CEP used; no pending order carried `CODPARCDEST`/`CODCONTATO`.

Executed (least-destructive sanctioned representation, per operator choice):

1. one `CRUDServiceProvider.saveRecord` created bounded contact `CODCONTATO=1` ("MPC D8 PROBE DEST ALT") on the test partner with a deliberately distinct CEP — reread converged;
2. one `saveRecord` bound `CODCONTATO=1` to test order 898263 — reread converged;
3. after P3, generated nota 898877 rereads with `CODCONTATO=1` — **contact reference survives fiscal progression**;
4. Partner master fields untouched; unrelated orders untouched.

Material findings:

- the sanctioned `Contato` entity carries **no street-address fields** (fieldset probe: `ENDERECO` invalid descriptor; only CEP/phone/email exist) — a full alternate street destination is **not representable** through the contact surface;
- `CODCONTATO` propagation is a contact reference, **not** a fiscal destination override — exactly the distinction protocol §6.3.4 warns about;
- XML/fiscal-document destination survival is untestable on this fixture because the nota was deliberately left unconfirmed (operator constraint).

Architecture implication: D4 Destination Realization must keep full alternate-destination delivery as **explicit external-required/unsupported** on this SourceInstance unless a future probe proves a sanctioned destination-override surface (e.g. through the fiscal document itself). The safety rails (never overwrite Partner master, never create a second Partner for an address) held. No reopen — this narrows a capability claim rather than contradicting an accepted one.

---

## P6 — unexercised fiscal branch/component — `NOT_TRIGGERED`

D4 §5.14 names the bounded fiscal Unknowns: out-of-state **PJ-contributor** branch; IPI/retained ICMS-ST component visibility; accessory-expense request shape. The selected fixture is a PF **consumidor final** interstate sale (TIPPESSOA=F, CLASSIFICMS=C), single ordinary item, no accessory expense, no IPI/ST component claim — none of the named Unknowns is material to this flow. Per D8-R1 §3.2, `NOT_TRIGGERED` with these fixture facts; the Unknowns remain fail-honest Unknowns.

---

## Residue / cleanup state

- ML fixture fully restored (price 29.90, qty 3, active).
- Sankhya test artifacts remain, owned by the operator's test lane: contact `142907/1`, order 898263 now atendido, unconfirmed nota 898877 (TOP 306, no fiscal number). Their disposal follows ordinary lawful business handling; D8 performed no destructive cleanup.
- Rotated ML credential v51 persisted in the Installation store; local Sankhya bearer token cache deleted after execution.
- No secret, buyer PII or raw provider payload is retained in this record.
