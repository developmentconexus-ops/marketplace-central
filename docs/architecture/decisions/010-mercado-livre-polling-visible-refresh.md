# ADR-010: Mercado Livre access is polling/GET only; "live" data requires a visible refresh

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-004's MVP demo scope but no document was
ever written under a stable number. It is reconstructed here from Assertion A2 of the
citations of `ADR-09` gathered in
`docs/architecture/decisions/_citations/adr-009-citations.md` — a distinct rule from the
MIS-007 fee-provenance decision (`ADR-009`, this repo's number) that happened to share the
same `ADR-09` token before the renumbering registry existed.

## Context

MIS-004 needed to demo Mercado Livre data (prices, listings, market signals) inside a
three-day build window, without standing up infrastructure the mission had no time to
secure. A webhook receiver needs a public endpoint or tunnel and an async delivery path;
building and hardening that in three days was judged out of reach. The alternative —
fetching data with the same GET requests any authenticated client can make, on demand or
on a local schedule — needed no new exposed surface, but it means the data on screen is
only as fresh as the last fetch, which the operator or the UI has to make visible rather
than implicitly assumed.

## Decision

**MIS-004 accesses Mercado Livre exclusively through polling/GET requests; webhooks are
out of scope for this mission. Because there is no push channel, any screen showing
"live" provider data must show how old that data is and must never present a failed
refresh as if it had succeeded.**

**§1 — No webhooks in MIS-004; refresh is on-demand or locally scheduled GET.** Webhook
delivery is deferred to a later mission.
> `.mnfs/MIS-004-mvp-demo/mission.md:93` — "ADR-09 Polling/GET only | ratificada |
> endpoint público/tunnel em 3 dias | idade da coleta visível; falha de refresh = estado
> stale honesto | refresh visível na UI".
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:68-69` — "### ADR-09:
> Polling/GET only no MVP — Decisão: sem webhooks no MIS-004; refresh por GET/polling
> on-demand. Webhooks completos = MIS-005 M-06."

**§2 — Age of the fetched data is a required, visible field, not an internal detail.**
The UI must show how old the collected data is; a refresh that fails must render as an
honest stale state, never as data silently treated as current.
> `.mnfs/MIS-004-mvp-demo/mission.md:93` — "idade da coleta visível; falha de refresh =
> estado stale honesto | refresh visível na UI."

**§3 — "Live" is a claim that costs a manual or scheduled refresh, never assumed.** The
mission's own accepted-trade-offs list names this as a cost the operator explicitly took
on, not a gap discovered later.
> `.mnfs/MIS-004-mvp-demo/mission.md:98` — "ADR-09 dado \"ao vivo\" exige refresh."
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:71-72` — "Trade-off:
> dados \"ao vivo\" exigem refresh manual/agendado local. Validation impact: refresh
> visível na UI (idade da coleta)."

## Contradictions

**Number collision with an unrelated decision.** The token `ADR-09` also names the
MIS-007 fee-provenance rule reconstructed as this repo's `ADR-009`
(`009-fee-value-carries-provenance.md`) and, more thinly, a MIS-001 proportional-security
rule (one citation, out of scope for this reconstruction). Nothing in the citing text
disambiguates the number; only the mission the citing file belongs to does.

**The candidate number itself churned during MIS-004 planning.** Before ratification,
this same rule was proposed as `ADR-09` by one planning candidate and as `ADR-10` by the
other, reconciled as semantically identical and folded into the mission table under
`ADR-09` — i.e. even inside a single mission's planning process the number was not fixed
until the reconciliation step.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:17` — "Polling/GET
> only; sem webhooks/scheduler no MIS-004 | ADR-09 | ADR-10 | idêntico."

## Rationale

Once a system stops receiving pushes, the freshness of what it shows becomes a fact the
system must track and disclose, not an assumption it can silently make. Treating polled
data as equivalent to pushed data — rendering it without an age indicator, or worse,
leaving a stale value on screen after a refresh silently failed — reintroduces the exact
"looks answered but wasn't measured" defect `ADR-017` names for zeros, applied here to
timestamps: a screen that does not distinguish "just fetched" from "fetch failed an hour
ago" is lying about currency the same way a fabricated zero lies about a magnitude.

## Consequences

- Every MIS-004 screen surfacing Mercado Livre data carries a visible collection-age
  indicator; a refresh action is user-triggered or locally scheduled, never implied.
- A failed refresh must degrade to a visibly stale state — the previous good value stays
  on screen with its true age, it is never replaced by a default or hidden behind a
  generic loading spinner that implies a retry is already in flight.
- Webhook-based, always-fresh delivery is explicitly deferred; MIS-005 owns closing that
  gap (`p3-claude-candidate-r01.md:69`: "Webhooks completos = MIS-005 M-06").

## Alternatives Considered

**Stand up a webhook receiver inside the three-day window.** Rejected: requires a public
endpoint or tunnel and an async delivery/retry path, both judged unbuildable and
unhardenable in the time available (`mission.md:93`, "Prevents: endpoint público/tunnel
em 3 dias").

**Poll silently on a fixed interval and always show the latest cached value as current.**
Rejected: this is the exact failure the visible-refresh half of the rule exists to
prevent — it would present provider data as live when it might be stale by a full polling
interval or by a failed fetch, with no way for the operator to tell.

## Unverified claims

None. All anchors cited for Assertion A2 were read and confirmed to match verbatim at the
cited line.
