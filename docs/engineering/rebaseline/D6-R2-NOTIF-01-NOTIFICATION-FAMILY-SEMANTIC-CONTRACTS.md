# NOTIF-01 — Notification Family Semantic Contracts

> **Status:** CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Parent gate:** [D1-R Producer Edges & Notification Routing Boundary Correction](D6-R2-NOTIF-01-D1-R-PRODUCER-ROUTING-BOUNDARY-CORRECTION.md)
> **Accepted inputs:** operator-approved [D0-R Trigger-Scope Correction](D6-R2-NOTIF-01-D0-R-TRIGGER-SCOPE-CORRECTION.md) + operator-approved [Trigger/Audience Census](D6-R2-NOTIF-01-REFERENCE-STUDY.md)
> **Purpose:** define exactly what each admitted Launch-V1 `NotificationKind` means to a human before D1-R is ratified and before D2-R models identity/state
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Contract law

A `NotificationKind` is not an enum label in search of semantics. Each family must have one bounded human meaning and one source-owned birth transition.

Every contract below defines:

```text
human meaning
attention value
source owner
birth transition
not-this boundary
audience strategy + consumer
deep-link home
repeat law
suppression / overlap law
conceptual user message
```

The conceptual message is **not final copy authority**. It exists only to prove that the semantic contract can be explained in user language without leaking provider/runtime mechanics.

A family must split before D2 if materially different subtypes require different source authority, human job, audience strategy, deep-link surface, repeat law or interaction. Do not hide such differences behind a generic `reason:string` or routing DSL.

---

# 2. Family contracts

## F01 — `MARKETPLACE_INSTALLATION_ATTENTION`

**Human meaning**  
A marketplace channel/account that the Organization relies on has crossed into a materially degraded or non-operable condition that a marketplace operator/admin should notice.

**Attention value**  
The condition may block or materially constrain normal marketplace operation across many downstream objects; waiting for the user to discover it inside Settings is insufficient.

**Source owner**  
`Marketplace Portfolio`.

**Birth transition**  
Marketplace Portfolio commits a new installation-level attention occurrence after provider/integration evidence has been interpreted into accepted Portfolio meaning, for example effective authorization/capability/operability becoming materially blocked or an accepted installation-posture signal crossing a material attention boundary.

**Not this**

```text
raw provider webhook/token error
one transient transport failure
routine credential refresh
polling the same unresolved condition
ordinary installation configuration change
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace operator/admin humans for this kind.

**Deep-link home**  
`/configuracoes/canais/:marketplaceInstallationId`.

**Repeat law**  
One Notification per distinct material attention occurrence. Re-reading the same unresolved posture never retriggers. A genuinely resolved condition that later re-enters a new material attention episode may create a new occurrence.

**Suppression / overlap**  
If the exact occurrence immediately produces Work already assigned to recipient `P`, the approved per-recipient Work-replacement suppression may suppress this source-family item for `P`; other routed recipients remain eligible.

**Conceptual user message**  
“Canal Mercado Livre precisa de atenção — a conta não está plenamente operável.”

---

## F02 — `OFFERING_ASYNC_ACTION_RESULT`

**Human meaning**  
A consequential ListingIntent/PriceIntent action started by this human did not finish conclusively in the initiating interaction and later reached a material outcome the initiator needs to know.

**Attention value**  
The user should not have to keep the editor open or repeatedly poll an intent to discover whether a consequential marketplace action eventually converged, failed or became uncertain.

**Source owner**  
`Marketplace Offering Operations`.

**Birth transition**  
A human-initiated action with exact initiator lineage was previously non-terminal for that human, and Offering later commits a materially new accepted outcome such as `converged`, `rejected`, `ambiguous` or `divergent` after its authoritative convergence/reconciliation semantics.

**Not this**

```text
synchronous terminal result already returned to the initiator
background/automation action with no exact human initiator
routine provider reread with no new material outcome
raw provider 2xx/202
```

**Audience / consumer**  
`DIRECT_SOURCE` → exact initiating human Principal owned by Offering lineage.

**Deep-link home**  
The exact ListingIntent, PriceIntent or Listing/Price workspace that owns the action.

**Repeat law**  
A terminal outcome is not repeated. A later **materially different** outcome in the same uncertainty lineage can be new awareness: e.g. `ambiguous` followed later by authoritative `converged` is a legitimate second occurrence because uncertainty has actually been resolved.

**Suppression / overlap**  
If the same result occurrence produces Work already assigned to the initiator, `WORK_ASSIGNMENT` may replace this family for that human under the approved suppression law.

**Conceptual user message**  
“Alteração de preço concluída” / “Publicação não convergiu e precisa de atenção.”

---

## F03 — `AVAILABILITY_ATTENTION`

**Human meaning**  
Availability Control can no longer safely maintain the accepted sellable-availability outcome for a marketplace target without human attention.

**Attention value**  
Availability failures can create oversell/undersell or stale channel exposure; material divergence should not depend on someone continuously watching the Availability page.

**Source owner**  
`Availability Control`.

**Birth transition**  
Availability commits a new owner-qualified attention occurrence: synchronization/convergence becomes materially blocked, ambiguous or divergent, or the accepted availability objective cannot currently be maintained.

**Not this**

```text
ordinary quantity change
successful synchronization
known zero availability
routine periodic refresh
one provider transport retry before owner meaning changes
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace/availability operators.

**Deep-link home**  
`/disponibilidade` with the exact affected subject/Installation context.

**Repeat law**  
No repeat while the same attention occurrence remains unresolved. Resolution followed by a later distinct failure episode may generate new awareness. A materially distinct later attention class for the same subject may also qualify if the owner commits a new occurrence.

**Suppression / overlap**  
Per-recipient Work replacement applies when the exact Availability occurrence creates assigned Work.

**Conceptual user message**  
“Disponibilidade não está convergindo para a publicação X.”

---

## F04 — `ECONOMIC_RECONCILIATION_ATTENTION`

**Human meaning**  
A commercial-economic attribution/reconciliation cannot be resolved automatically from accepted evidence and now requires bounded human resolution.

**Attention value**  
A reconciliation item can remain economically misleading or incomplete until someone resolves it; ordinary metric movement, however, is analysis rather than interruption-worthy awareness.

**Source owner**  
`Commercial Economics`.

**Birth transition**  
Economics commits a new attribution/reconciliation state whose accepted semantics require human resolution rather than ordinary automatic interpretation.

**Not this**

```text
margin changed
ROAS changed
normal realized-vs-expected variance
new Sale Economics record
routine recalculation
```

**Audience / consumer**  
`ORG_ROUTED` → configured commercial/marketplace managers responsible for reconciliation.

**Deep-link home**  
`/economia/reconciliacao` with the exact attribution/reconciliation context.

**Repeat law**  
One occurrence per distinct unresolved reconciliation issue. Recalculation/polling does not retrigger the same issue. A new issue or a resolved issue that later re-enters a distinct unresolved episode may notify again.

**Suppression / overlap**  
Per-recipient Work replacement applies if the exact reconciliation occurrence is materialized as already-assigned Work.

**Conceptual user message**  
“Reconciliação econômica precisa de revisão — atribuição não pôde ser resolvida automaticamente.”

---

## F05 — `NEW_MARKETPLACE_SALE`

**Human meaning**  
A new marketplace demand has been authoritatively recognized by MPC and the configured operations humans should know that a new Sale now exists.

**Attention value**  
A new Sale begins a real operating lifecycle and may require timely commercial/operational follow-up even while the user is elsewhere in the product.

**Source owner**  
`Marketplace Sales`.

**Birth transition**  
Marketplace Sales commits the **first authoritative confirmation** of one canonical source-qualified Marketplace Sale in MPC after provider evidence is authoritatively reread/interpreted.

**Not this**

```text
orders_v2/webhook arrived
same Sale was polled again
Sale fields changed later
shipment changed
cancellation/hold/fraud stop
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace-operations humans; the Organization may also choose fulfillment humans if that is genuinely its operating model, but Permission alone never subscribes them.

**Deep-link home**  
`/vendas/:nativeSaleKey` under the exact Marketplace Installation context.

**Repeat law**  
Exactly one `NEW_MARKETPLACE_SALE` awareness occurrence per canonical Sale identity. Later Sale changes belong to other admitted families or remain in the Sale workspace.

**Suppression / overlap**  
No routine Work is fabricated for a normal new Sale. If the exact first-confirmation occurrence exceptionally also creates Work already assigned to one routed human, the generic per-recipient Work-replacement law may apply only with explicit source-occurrence correlation.

**Conceptual user message**  
“Nova venda recebida — Mercado Livre · Pedido #…”

---

## F06 — `SALE_ATTENTION`

**Human meaning**  
A Sale that existed in MPC has entered a Sale-owned condition that materially changes whether/how normal processing may safely continue.

**Attention value**  
Safe-handling stops such as cancellation/hold/fraud-risk-like conditions can make continuing ordinary fulfillment materially wrong; the operator should notice promptly.

**Source owner**  
`Marketplace Sales`.

**Birth transition**  
Marketplace Sales commits a new owner-qualified safe-handling attention occurrence for a canonical Sale, such as a cancellation/hold/fraud-risk-like stop or another accepted Sale-owned condition that materially changes handling.

**Not this**

```text
new Sale creation
ordinary address/amount/status refresh with no handling consequence
Shipment exception
PostSaleResolution-owned action requirement
raw provider status string
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace-operations humans.

**Deep-link home**  
`/vendas/:nativeSaleKey`.

**Repeat law**  
The same unresolved Sale attention occurrence does not repeat. A later distinct safe-handling occurrence may notify again. Different source subtypes stay in this family only while they retain the same owner, audience, human job and Sale-detail continuation.

**Suppression / overlap**  
Per-recipient Work replacement applies. **Global-coherence candidate:** when the same underlying consequence immediately creates a richer `PostSaleResolution` occurrence and the same human is routed to `POST_SALE_ATTENTION`, the two families should not create duplicate awareness for that recipient; see §4.2.

**Conceptual user message**  
“Venda #… precisa de atenção antes de continuar — processamento seguro foi interrompido.”

---

## F07 — `MATERIALIZATION_ATTENTION`

**Human meaning**  
The MPC-owned business-system materialization for a Sale cannot safely converge/progress automatically and the backoffice/marketplace operator responsible for materialization needs to intervene or understand the exception.

**Attention value**  
Ambiguous/rejected/divergent/external-required business-order or invoicing materialization can block downstream operation while remaining invisible to a user focused on Sales or Fulfillment.

**Source owner**  
`Business-System Materialization`.

**Birth transition**  
Materialization commits a new Business Order or Invoicing materialization attention occurrence such as `ambiguous`, materially `rejected`, `divergent` or `external-required` where human awareness is needed. Normal successful convergence is not an alert.

**Not this**

```text
successful TOP/order materialization
successful 313→306 progression
one transport retry before owner interpretation
Fulfillment not yet actionable merely because upstream work is still normal
raw Sankhya response
```

**Audience / consumer**  
`ORG_ROUTED` → one configured **business-system materialization/backoffice operations audience**. Baseline routing does not vary by hidden materialization subreason.

**Deep-link home**  
`/vendas/:nativeSaleKey` → materialization/invoicing region.

**Repeat law**  
Same unresolved materialization occurrence does not repeat. A later distinct materialization failure/attention episode can notify again.

**Suppression / overlap**  
Per-recipient Work replacement applies. Fulfillment humans should receive `FULFILLMENT_ACTIONABLE`/`FULFILLMENT_ATTENTION` when Fulfillment-owned meaning changes rather than requiring reason-dependent routing inside this kind.

**Conceptual user message**  
“Pedido/faturamento da venda #… precisa de atenção — materialização não convergiu.”

---

## F08 — `FULFILLMENT_ACTIONABLE`

**Human meaning**  
A FulfillmentExecution that previously could not legitimately begin physical execution has now crossed the first accepted boundary where the fulfillment team can actually start.

**Attention value**  
This is a human handoff from upstream sales/materialization prerequisites into physical work; the fulfillment operator should not need to poll for new executable work continuously.

**Source owner**  
`Fulfillment Lifecycle`.

**Birth transition**  
Fulfillment commits the **first** owner-qualified transition of a `FulfillmentExecution` into an actionable physical-execution state.

**Not this**

```text
Sale merely exists
Business Order merely converged
separation recorded
conference recorded
packing recorded
dispatch handoff recorded
execution was already actionable and was reread
```

**Audience / consumer**  
`ORG_ROUTED` → configured fulfillment/dispatch operators.

**Deep-link home**  
`/expedicao/execucoes/:fulfillmentExecutionId`.

**Repeat law**  
At most once for the first actionable entry of that FulfillmentExecution. A later blocked→unblocked resume does not recreate `FULFILLMENT_ACTIONABLE` by default; later material problems belong to `FULFILLMENT_ATTENTION` and ordinary resumed progression belongs in the specialized queue unless a future consumer proves otherwise.

**Suppression / overlap**  
Per-recipient Work replacement applies if this exact entry occurrence creates assigned Work, although routine actionable fulfillment should normally remain ordinary operational work rather than fabricated Work.

**Conceptual user message**  
“Venda #… está pronta para iniciar a expedição.”

---

## F09 — `FULFILLMENT_ATTENTION`

**Human meaning**  
An existing physical execution has crossed a material attention boundary that requires the fulfillment/dispatch team to notice a risk, discrepancy or blocker outside routine progression.

**Attention value**  
Dispatch risk, physical discrepancy or provider-readiness blockage can require timely intervention and should not be hidden among normal checkpoint activity.

**Source owner**  
`Fulfillment Lifecycle`.

**Birth transition**  
Fulfillment commits a new owner-qualified material attention occurrence, e.g. dispatch due/at-risk, conference discrepancy, blocked physical readiness, or a provider-requirement closure failure that materially obstructs the accepted fulfillment path.

**Not this**

```text
every separation/conference/packing checkpoint
normal dispatch handoff
normal completion
routine queue aging with no owner-declared attention transition
raw provider requirement payload
```

**Audience / consumer**  
`ORG_ROUTED` → configured fulfillment/dispatch operators.

**Deep-link home**  
`/expedicao/execucoes/:fulfillmentExecutionId`.

**Repeat law**  
No repeat for the same unresolved attention occurrence/subtype. A later materially distinct attention transition can produce another Notification. The family remains grouped because admitted subtypes share the same owner, audience strategy, human job and execution-detail continuation; if future evidence breaks that equivalence, split the kind instead of adding generic routing rules.

**Suppression / overlap**  
Per-recipient Work replacement applies when the exact Fulfillment occurrence creates assigned Work.

**Conceptual user message**  
“Expedição da venda #… precisa de atenção — prazo/execução entrou em risco.”

---

## F10 — `SHIPMENT_EXCEPTION`

**Human meaning**  
A dispatched/in-transit Shipment has entered a material delivery exception that marketplace/fulfillment operations should notice.

**Attention value**  
The operational locus has moved from physical preparation to delivery observation; a real exception may require seller action even when nobody is viewing the Shipment list.

**Source owner**  
`Fulfillment Lifecycle` for accepted source-qualified Shipment observation/exception meaning.

**Birth transition**  
Fulfillment commits a new source-qualified Shipment exception occurrence after provider evidence is interpreted into an accepted material exception/outcome.

**Not this**

```text
shipment provider webhook arrived
normal tracking update
normal dispatch
normal delivered outcome
re-poll of the same provider exception
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace/fulfillment operations humans.

**Deep-link home**  
`/expedicao/envios/:nativeShipmentKey` under exact source context.

**Repeat law**  
Same exception episode does not repeat. A resolved Shipment that later enters a new distinct exception episode, or a materially distinct later exception occurrence, may notify again.

**Suppression / overlap**  
Per-recipient Work replacement applies when the exact Shipment exception creates already-assigned Work.

**Conceptual user message**  
“Problema no envio da venda #… — entrega entrou em exceção.”

---

## F11 — `POST_SALE_ATTENTION`

**Human meaning**  
A seller-relevant cancellation/return/refund/claim consequence is now being coordinated by Post-Sale Resolution and requires awareness or action from post-sale operations.

**Attention value**  
Post-sale consequences can create deadlines, reverse physical/economic work and customer-facing obligations that deserve attention even when the user is outside the Post-Sale workspace.

**Source owner**  
`Post-Sale Resolution`.

**Birth transition**  
Post-Sale commits a new resolution-owned attention occurrence: a `PostSaleResolution` opens for a material consequence, or accepted provider/source evidence creates a new action requirement, consequence or due-date attention transition inside an existing resolution.

**Not this**

```text
Sale-owned safe-handling stop before Post-Sale owns consequence coordination
ordinary refund evidence with no new Post-Sale attention meaning
every status update inside an existing resolution
raw claim/return webhook
```

**Audience / consumer**  
`ORG_ROUTED` → configured marketplace/post-sale operators.

**Deep-link home**  
`/pos-venda/:postSaleResolutionId`.

**Repeat law**  
Opening one resolution is one occurrence. Later materially distinct action/due-date/consequence transitions in the same resolution may create new occurrences; rereading the same pending action does not.

**Suppression / overlap**  
Per-recipient Work replacement applies. **Global-coherence candidate:** if the exact consequence simultaneously qualifies `SALE_ATTENTION` only as a precursor and `POST_SALE_ATTENTION` as the richer owned continuation, Personal Notifications should prefer the Post-Sale awareness for recipients routed to both; see §4.2.

**Conceptual user message**  
“Pós-venda da venda #… precisa de atenção — há uma devolução/reclamação a tratar.”

---

## F12 — `WORK_ASSIGNMENT`

**Human meaning**  
A material actionable Work item has been assigned or reassigned specifically to this human, creating direct personal responsibility.

**Attention value**  
Assignment is the strongest personal-awareness case: the user should know that a responsibility is now theirs without needing to watch the Work queue continuously.

**Source owner**  
`Operational Work`.

**Birth transition**  
Work commits assignment/reassignment to exact human Principal `P`.

**Not this**

```text
Work created but not assigned to P
hold/resume
escalation with no new assignment
Work updated while still assigned to P
Work closed
```

**Audience / consumer**  
`DIRECT_SOURCE` → exact assigned/reassigned human Principal.

**Deep-link home**  
`/trabalho/:workId`.

**Repeat law**  
Every distinct assignment occurrence is distinct awareness. `A → B → A` creates two separate historical assignment Notifications for A; replay of the same assignment occurrence creates only one.

**Suppression / overlap**  
This is the **preferred personal item** when the exact underlying source occurrence also would route a generic source-family alert to the same newly assigned human. It never suppresses source alerts for other recipients.

**Conceptual user message**  
“Trabalho atribuído a você — revisar divergência da venda #…”

---

## F13 — `AUTHORIZATION_ACTION_REQUIRED`

**Human meaning**  
A governed action is currently pending and this human is, according to Governance’s current authority semantics, one of the exact Principals who can legitimately decide it.

**Attention value**  
A decision request can block consequential product work while the approver is elsewhere; the approver should not have to poll the Approvals queue to learn that their authority is needed.

**Source owner**  
`Controlled Action Governance`.

**Birth transition**  
Governance commits an authorization state in which one governed target/action is pending and Principal `P` becomes currently valid to decide it according to accepted delegation/grant/authority semantics.

**Not this**

```text
P merely has approvals.read
P could theoretically approve some action class but no decision is pending
same pending decision is reread
expired/revoked authority with no current decision eligibility
```

**Audience / consumer**  
`OWNER_DERIVED` → exact currently valid decision Principal set resolved by Governance.

**Deep-link home**  
`/aprovacoes/:authorizationDecisionId` and/or the exact governed target context admitted by Governance.

**Repeat law**  
One occurrence per decision–Principal actionable-eligibility episode. Unchanged pending eligibility does not repeat. If Governance later removes eligibility and subsequently establishes a genuinely new eligibility episode for the same still-pending/new decision, that is a distinct source occurrence to be modeled downstream.

**Suppression / overlap**  
Per-recipient Work replacement may apply only if the exact Governance action-required occurrence is explicitly correlated to Work already assigned to that same Principal; Notifications must never infer equivalence from titles or target IDs alone.

**Conceptual user message**  
“Aprovação precisa da sua decisão — alteração de preço da publicação X.”

---

## F14 — `AUTHORIZATION_DECISION_RESULT`

**Human meaning**  
A human who requested/initiated a governed action now needs to know that the authorization decision has been committed and what its disposition is.

**Attention value**  
The decision commonly completes outside the initiating interaction; the requester should not need to keep checking Approvals to discover whether the blocked action can proceed or was denied.

**Source owner**  
`Controlled Action Governance`.

**Birth transition**  
Governance commits the authorization decision and owns exact requester/initiator human lineage for a recipient who did not already receive the same terminal result as the direct actor in the same decision interaction.

**Not this**

```text
every audit/history update
recipient has approvals.read but did not request/initiate the action
same committed decision reread
requester is the deciding actor and already received the terminal decision synchronously
```

**Audience / consumer**  
`DIRECT_SOURCE` → exact requester/initiating human Principal owned by Governance lineage.

**Deep-link home**  
`/aprovacoes/:authorizationDecisionId` plus the governed target continuation.

**Repeat law**  
One result awareness per committed authorization decision for the exact recipient. A later separate authorization decision/reopen must have distinct source identity; the same committed result never retriggers from polling.

**Suppression / overlap**  
Per-recipient Work replacement may apply if the exact decision-result occurrence creates already-assigned Work. This family does not replace `OFFERING_ASYNC_ACTION_RESULT`: Governance reports authorization disposition; the action-owning domain later reports consequential external-effect outcome if that outcome was asynchronous and materially separate.

**Conceptual user message**  
“Sua solicitação foi aprovada” / “Sua solicitação foi rejeitada.”

---

# 3. Cross-family lifecycle example

One Sale can legitimately create several Notifications because they represent **different human handoffs**, not a synthetic global Sale workflow:

```text
Marketplace Sales
  NEW_MARKETPLACE_SALE
        ↓
Business-System Materialization
  [normal success → no Notification]
        ↓
Fulfillment Lifecycle
  FULFILLMENT_ACTIONABLE
        ↓
  routine checkpoints → no Notifications
        ↓
  FULFILLMENT_ATTENTION   only if material attention transition occurs
        ↓
Shipment observation
  SHIPMENT_EXCEPTION      only if delivery materially diverges
        ↓
Post-Sale Resolution
  POST_SALE_ATTENTION     only when post-sale consequence/action becomes material
```

A different human audience can be configured at each ORG_ROUTED handoff without creating one cross-owner workflow authority.

---

# 4. Global coherence review

## 4.1 Result — 14 families remain sufficient

The semantic pass finds no need to add a fifteenth Launch-V1 family. The accepted 14 remain sufficient **provided two boundary clarifications are ratified** below.

The families group multiple source subtypes only where they share:

```text
same semantic owner
same audience strategy
same primary human consumer
same deep-link/workspace continuation
same repeat/suppression posture
```

This keeps `SALE_ATTENTION` and `FULFILLMENT_ATTENTION` bounded without turning either into a generic reason/routing DSL. If later a subtype violates one of those equivalences, split the family through the smallest D0/D1 reopen.

## 4.2 Coherence finding A — Sale vs Post-Sale duplicate awareness

One external business development can legitimately produce two source-owned facts:

```text
Marketplace Sales
→ Sale safe handling must stop

Post-Sale Resolution
→ a material consequence/resolution must now be coordinated
```

Both source facts remain valid. But delivering both to the **same human for the same consequence** can create duplicate awareness with no added job value.

Candidate Notification-level precedence:

```text
same bounded underlying consequence correlation
+ recipient P would receive SALE_ATTENTION
+ recipient P would receive POST_SALE_ATTENTION
+ PostSaleResolution already owns the richer continuation
→ suppress SALE_ATTENTION for P
→ deliver POST_SALE_ATTENTION for P
```

This is per-recipient awareness suppression only. It never suppresses the Sale source fact, never closes/changes PostSale, and never applies when the two occurrences represent genuinely different human needs.

**D1-R implication:** if the operator ratifies these contracts, extend the existing bounded suppression ownership from Work replacement to this one proved Sale→Post-Sale precedence case. Do **not** create a generic cross-kind deduplication engine.

## 4.3 Coherence finding B — Materialization audience must not vary by hidden subreason

The earlier census wording allowed `MATERIALIZATION_ATTENTION` consumer language such as “marketplace operations / fulfillment depending on source condition.” That would conflict with the accepted A3 model:

```text
(Organization, NotificationKind)
→ exact configured recipients
```

If one `NotificationKind` silently chooses different routing based on materialization subreason, D1 has accidentally introduced the routing DSL it intended to avoid.

Corrected semantic contract:

> `MATERIALIZATION_ATTENTION` routes to one explicitly configured business-system-materialization/backoffice operations audience. Fulfillment receives its own owner-qualified `FULFILLMENT_ACTIONABLE` / `FULFILLMENT_ATTENTION` occurrences when Fulfillment meaning changes.

**D1-R implication:** ratification should narrow the D1-R wording to this single routing posture.

## 4.4 Offering result vs Governance result are not duplicates

These remain separate:

```text
AUTHORIZATION_DECISION_RESULT
= Governance decided whether the action is authorized

OFFERING_ASYNC_ACTION_RESULT
= Offering later learned the actual consequential marketplace outcome
```

Approval is not execution success. A user may legitimately receive both at different times.

## 4.5 Fulfillment actionable vs attention are not duplicates

These represent different jobs:

```text
FULFILLMENT_ACTIONABLE
= physical work can begin for the first time

FULFILLMENT_ATTENTION
= existing execution later crossed a material risk/blocker boundary
```

Normal checkpoints remain silent.

## 4.6 Infrastructure/provider signals remain outside Product families

No family is born from raw D4/D7 signals. Provider events, retries, queues, worker failures and observability incidents must first become accepted owner meaning or remain technical operations concerns.

---

# 5. Global negative controls

The semantic contract set fails review if it permits any of the following:

1. a kind with no exact human meaning beyond “status changed”;
2. raw webhook/topic/runtime events becoming Product Notifications directly;
3. one kind requiring different audiences through hidden subtype/routing logic;
4. Permission alone becoming recipient responsibility;
5. same unresolved occurrence notifying repeatedly because of polling/replay;
6. routine checkpoint/activity log behavior in the personal Inbox;
7. source-domain truth being copied/reinterpreted by Personal Notifications;
8. source access being granted merely because a Notification exists;
9. suppression changing/deleting source truth or Work/Post-Sale state;
10. a generic cross-kind dedup/rule engine replacing the two bounded proved suppression cases;
11. final UI copy, API payload, event name, DB schema or River mechanics being selected inside this semantic gate;
12. preserving 14 families by force when a future subtype proves materially different audience/job/deep-link semantics.

---

# 6. Candidate D1-R closure effect

If the operator approves this document:

```text
D0-R Product scope                    ACCEPTED
D1-R producer edges                   candidate semantics retained
14 family semantic contracts          ACCEPTED
Materialization routing clarification ACCEPTED
Sale→Post-Sale suppression precedence ACCEPTED
        ↓
D1-R can be amended/ratified as one coherent boundary authority
        ↓
D2-R may finally rederive identity/source/routing state
```

D2-R, D3, OAD, B00 bell/Inbox and Product implementation remain blocked until their own gates.

**Exact next action:** operator adjudicates these fourteen semantic contracts and the two global coherence findings. Do not ratify D1-R or begin D2-R before explicit approval.