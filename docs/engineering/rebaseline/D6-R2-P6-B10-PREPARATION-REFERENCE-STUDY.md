# D6-R2 P6 — B10 Preparation Reference Study

> **Status:** DERIVED / REVALIDATED — B10 structure confirmed; operator-language P8 candidate pending walkthrough
> **Block:** B10 Preparation / R10
> **Method profile:** `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD.md + FRONTEND-METHOD.md`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Human job and accepted boundary

B10 supports:

```text
exact Organization + Marketplace Installation
→ search admitted source products
→ select exact SourceInstance + native product key
→ inspect readiness + publication requirements
→ resolve/clear correspondence when needed and authorized
→ re-read readiness after consequential correspondence effect
→ continue to ListingIntent only when current pre-ListingIntent state permits
```

B10 does not own source master data, a generic PIM, provider mutation, ListingIntent editing, Offering dispatchability, or Authorization.

## 2. Reference evidence

P6 studied mature products only for task structure, never as MPC authority:

- Google Merchant Center — search/triage → affected product → issue detail:  
  <https://support.google.com/merchants/answer/12476548>
- Akeneo — product grid → readiness/completeness → exact product detail:  
  <https://help.akeneo.com/en_US/product-readiness/product-readiness-overview>
- Amazon Seller Central — inventory search/filter → status/reason → exact listing/product correction context:  
  <https://sellercentral.amazon.com/seller-forums/discussions/t/a8234e1f-f065-44d8-a4c4-f23cf92b3f4b>
- Mirakl/Lowe's — validation feedback remains tied to exact product/import identity:  
  <https://seller.lowes.com/mirakl-seller-documentation/>

Shared useful pattern:

```text
search/filter
→ compact issue/readiness signal
→ exact subject detail
→ explicit correction/revalidation boundary
```

MPC-specific rejection: none of those products authorize MPC to edit arbitrary source attributes, invent readiness scores, or expose raw provider DTO/field bags.

## 3. Structure decision

Three structures were considered:

```text
A persistent split-view
B progressive full inline expansion
C search/list → selected exact-subject detail
```

**C remains selected.**

- A has no evidence-backed need for high-frequency rapid cycling, is vertically fragile for deep requirements, and complicates responsive behavior.
- B weakens exact source identity and recovery/revalidation context.
- C preserves scanability in search and gives requirements/correspondence enough durable detail space.

**CURRENT STRUCTURE CONFIRMED.**

**P7 layout hypotheses remain NOT TRIGGERED.** Rendering rejected alternatives would manufacture ambiguity rather than resolve a real Product question.

## 4. P7 feasibility disposition

| Question | Disposition |
| --- | --- |
| Required fields / summaries | PRESENT-IN-AUTHORITY through source search, ProductChannelReadiness and PublicationRequirements |
| Identity sources | PRESENT-IN-AUTHORITY: Organization + exact Installation + SourceInstance + native product key |
| Pagination / scale | PRESENT-IN-AUTHORITY: cursor/limit; no total-count inference |
| Sort / filter | admitted query/source narrowing only; additional generic sort rejected by YAGNI |
| Preview / content truth | owner-local readiness/requirements/correspondence projections only |
| Material writes | Resolve/Clear correspondence only; ListingIntent creation is the downstream boundary |

## 5. P8 structural obligations

The browser-operable candidate must prove:

- locked B00 shell/IA remains intact;
- exact Organization and Marketplace Installation behavior;
- Organization switch invalidates Installation context;
- material search works;
- known-empty != unknown/unavailable;
- exact source identity remains technically preserved;
- complete applicable publication-requirement census is inspectable;
- correspondence effects require an authoritative re-read;
- ListingIntent is an explicit unopened downstream boundary;
- mobile drawer supports Escape and focus return;
- no frontend-created business truth.

## 6. Lock-time assumption A01

A01: real high-frequency rapid cycling across many products is not yet evidenced. It influenced rejection of persistent split-view but does not block the current search→detail candidate by itself.

**A01 disposition: PENDING OPERATOR.**

Before B10 can LOCK, the operator selects exactly one:

```text
ACCEPT_FOR_LOCK_WITH_LATER_PROBE
BLOCK_LOCK
```

## 7. First P8 candidate result

The first operated candidate exposed a material simplification: one generic `Categoria / atributos` bucket could not let the operator inspect the full marketplace-specific requirement basis.

The initial interpretation was that accepted narrative authority already covered the distinction and therefore no upstream repair was necessary. Subsequent wire-level proof falsified that conclusion.

## 8. Downstream falsification result

B10 demonstrated that the accepted publication meaning was richer than the machine-readable Product OAD realization. The affected downstream scope was paused and the smallest owner — Product OAD wire realization — was reopened.

This became the PR #68 prerequisite rather than a frontend workaround.

## 9. Source truth vs ListingIntent truth

The Product distinction remains:

```text
source evidence
!=
ListingIntent desired publication value
```

In particular:

```text
missing source != publication impossible
```

A missing source fact stays missing. Later, ListingIntent/Offering may author a permitted resolution such as `FOLLOW_SOURCE(candidate_key)` or `EXPLICIT_OVERRIDE(PublicationValue)` without rewriting source truth. B10 does not decide final draft dispatchability and does not edit the override.

## 10. Superseded pre-PR68 disposition

The earlier B10 candidate used shorthand states such as `met` and treated `required/recommended/conditional` as an applicability dimension. That representation is superseded by the accepted PR #68 wire contract.

The prior statement `UPSTREAM FINDING: NONE` is superseded for this exact wire issue.

## 11. PR #68 integrated — bounded B10 rebaseline

PR #68 is integrated into `main` at:

`ed3d164b0574b7950c2c7467d150c89576bba1ec`

**UPSTREAM FINDING: RESOLVED.**

The bounded rebaseline preserves the existing B10 search→detail structure and changes only how accepted publication truth is represented.

Canonical distinctions now carried by P8:

```text
requirement_class != applicability

requirement_class
= required / recommended / optional / conditional

applicability
= current / draft_dependent

source_evidence.state
= known / missing / conflicting / unknown / unavailable / unsupported
```

Each requirement also carries a bounded `value_spec`, `not_applicable_allowed`, and source evidence compatible with the requirement value family. The seven value-spec families are:

```text
text
exact_decimal
boolean
option
text_list
option_list
number_unit
```

Known/conflicting source evidence preserves opaque candidate identity. Conflicting evidence preserves multiple candidates rather than choosing one. `not_applicable` remains an explicit override value only; source evidence does not synthesize it.

The P8 candidate must preserve this machine truth without turning the wire vocabulary into the primary user language. The technical contract remains encoded and testable even when the operator surface translates it.

### 11.1 LOCK impact sweep

The PR #68 change affects B10's PublicationRequirements projection only. Existing LOCKED blocks do not consume that wire meaning:

```text
B00 | UNAFFECTED
B01 | UNAFFECTED
B00-R2 | UNAFFECTED
B11 | UNAFFECTED
B12 | UNAFFECTED
B110 | UNAFFECTED
```

No existing LOCK is reopened.

## 12. Operator REVISE — human language projection

The operator rejected the first post-PR68 P8 candidate as operationally unreadable even though its contract representation was technically accurate.

Root cause:

```text
wire-contract vocabulary
→ copied into primary screen language
→ operator must understand backend/API semantics to perform a normal preparation task
```

This conflicts with the accepted frontend-method rule **human needs before screens**. The frontend must project accepted truth into the user's job language; backend coherence does not authorize backend-shaped UX.

The Product/OAD authority remains sufficient. This is a bounded frontend P8 correction, not a new upstream finding.

### 12.1 Human questions the primary surface must answer

```text
O que o marketplace está pedindo?
O que já está atendido?
O que está faltando ou não pôde ser verificado?
Qual informação existe hoje?
O que o operador precisa fazer agora?
```

Primary table language therefore becomes:

```text
Requisito do marketplace
Exigência
Situação
Informação atual
O que fazer
```

Wire states are translated into operator states such as:

```text
Atendido
Falta informação
Há informações diferentes
Não foi possível verificar
Não disponível na fonte
```

Technical identities, revision keys, source-candidate keys and downstream resolution vocabulary move behind explicit secondary disclosure for support/audit. They remain encoded in the artifact and verifier; they are not deleted or reinterpreted.

### 12.2 P8 disposition after operator REVISE

The search→detail structure is unchanged. PR #68 semantics are unchanged. Existing LOCK impact remains unchanged.

The revised candidate must remain browser-operable and **CANDIDATE / NOT LOCKED** until the operator operates this human-language version.

`Operator walkthrough: PENDING`.

`P8 status: CANDIDATE / NOT LOCKED`.

A01 remains **PENDING OPERATOR** with the same two lawful dispositions:

```text
ACCEPT_FOR_LOCK_WITH_LATER_PROBE
BLOCK_LOCK
```

Only after explicit operator LOCK may B10 proceed to P9.
