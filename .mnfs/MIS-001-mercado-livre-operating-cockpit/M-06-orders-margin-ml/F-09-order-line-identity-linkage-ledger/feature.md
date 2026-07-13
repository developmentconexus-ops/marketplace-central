# F-09 Stable Order-Line Identity and Linkage Ledger

```yaml
id: F-09
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Outcome

Create the MPC-owned persistence substrate for assisted Sankhya linkage:
immutable Mercado Livre order-line identities and one tenant/installation-
scoped, idempotent, append-only external-order/internal-document mapping.

## Brief

Create the orders-owned durable line-identity and append-only linkage-ledger
substrate described in the Outcome, without Oracle/provider access or a public
transport surface.

## Scope

- Add migration `0033_orders_sankhya_linkage.sql` and update the canonical
  embedded-migration count test.
- Add immutable `mpc_line_id` plus an explicit identity/reconciliation state to
  persisted order items. Refreshes must not delete/regenerate a stable identity.
- Reconcile refreshes using persisted identity. Duplicate/indistinguishable or
  otherwise ambiguous rows remain explicitly non-linkable; attribute tuples and
  mutable `line_no` never become proof.
- Add generic orders-domain linkage identities and an orders-owned repository
  port/Postgres adapter that atomically inserts one order header mapping and all
  explicit line mappings, with unique idempotency and source/audit metadata.
- Keep Oracle names/table DTOs out of domain and public contracts. Postgres may
  store internal document/line numbers, TOP/reference evidence, and safe source
  refs, but this Feature performs no Oracle read and exposes no HTTP route.

## Acceptance criteria

1. Fresh import assigns a stable opaque MPC line ID; equal/newer refresh retains
   it, and stale refresh behavior remains unchanged.
2. Reordered unambiguous lines preserve identity by stable provider attributes;
   indistinguishable duplicates are marked ambiguous/non-linkable rather than
   silently remapped.
3. Existing positional rows receive an explicit legacy reconciliation state;
   migration never claims that `line_no` is immutable proof.
4. One transactional repository call inserts the account-scoped external order
   key, exact internal header identity, and every explicit line mapping. A
   semantically identical idempotent retry returns the existing mapping;
   conflicting reuse fails closed without overwrite.
5. Database constraints prove tenant/installation scoping, one mapping per MPC
   line, no active internal origin-line reuse, and append-only audit metadata
   including actor, reason, source time, idempotency key, config revision, and
   safe evidence state/reference.
6. Focused unit/migration tests pass; Postgres integration proof runs when the
   approved ephemeral target is available. No Oracle/provider write/read occurs.

## Expected Output

- Persisted order items have immutable opaque MPC line IDs and explicit stable,
  legacy-unresolved, or ambiguous reconciliation state.
- Refreshes retain identity only through provider item, variation, and seller
  SKU identity groups; mutable line number, quantity, price, and ambiguous
  tuples are never proof.
- Duplicate identity groups preserve the existing opaque ID set across
  refresh, create IDs only for genuinely new excess rows, and remain explicitly
  ambiguous/non-linkable.
- One tenant/installation-scoped transaction appends a generic internal header
  identity, every exact line identity, idempotency data, and complete audit
  metadata without overwriting history.
- Focused domain, repository, and migration evidence is recorded in this
  feature directory; ephemeral Postgres evidence is run only when approved and
  available.

## Non-goals

- Oracle candidate queries, custom-field metadata validation, `TGFVAR` reads,
  profitability tax consumption, HTTP/OpenAPI/SDK, UI, automatic matching, or
  production enablement.
- Any claim that caller-supplied audit actor is authenticated. Production auth
  and manual-adjustment hardening remain explicitly deferred/failing.
