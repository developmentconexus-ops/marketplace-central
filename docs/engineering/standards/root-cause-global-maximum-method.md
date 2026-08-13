# Root-Cause / Global-Maximum Engineering Method

> **Status:** canonical engineering standard  
> **Scope:** Marketplace Central engineering work  
> **Last verified:** 2026-08-13

## Binding principle

> **Always simplify the code, never simplify correctness. Find the root cause, test whether the proposed solution is only a local maximum, and prefer the global structure that makes the defect class unrepresentable or mechanically impossible at the strongest reasonable boundary.**

This method governs non-trivial engineering work. It is not a mandate for maximum sophistication. It distinguishes **essential complexity** from **accidental complexity** and removes the latter without weakening invariants.

## What this method prevents

1. **Patch-on-patch:** fixing a visible symptom while preserving the structural fact that made it possible.
2. **Local-maximum optimization:** improving a workaround, legacy seam or faulty foundation instead of replacing it.
3. **False simplification:** deleting enforcement/proof because it is inconvenient while the invalid state remains representable.
4. **Overengineering:** creating a framework, registry, parser, compatibility layer or second source of truth when a stronger existing boundary can enforce the property.
5. **Infinite review:** continuing hypothetical hardening after the material property is solved and proved.

## Definitions

### Symptom

The observable failure or finding: wrong response, race, duplicate implementation, CI escape, unsafe write, broken import boundary, contract drift or similar evidence.

A symptom says where a defect appeared; it does not determine where the fix belongs.

### Root cause

The structural fact that made the defect possible.

Examples:

- not “this handler forgot validation,” but “runtime request shapes and OpenAPI are independent authorities”;
- not “this module wrote the wrong tenant,” but “tenant identity can be invented/defaulted at several boundaries”;
- not “this marketplace mapping leaked,” but “provider DTOs are legally importable outside the provider boundary”;
- not “this table has two writers,” but “data ownership is not structurally assigned to one authority.”

### Target property / invariant

A statement that must remain true for every valid implementation, independent of current code shape.

Examples:

- each business fact has one owning context and one write authority;
- provider wire models cannot become domain contracts;
- an unknown monetary/operational fact cannot be represented as a plausible known zero;
- application code cannot write another context’s state directly;
- one API authority generates or mechanically validates all runtime/client contract shapes.

### Local maximum

A solution that improves the current implementation but preserves the structural limitation that produced the defect class.

A local maximum is legal only when explicitly transitional, with a named successor and deletion condition.

### Global maximum

The best sustainable structure for the system’s actual constraints: it resolves the root cause, converges authorities, preserves required invariants, minimizes accidental complexity and makes invalid states impossible or mechanically detectable at the strongest reasonable boundary.

Global maximum does **not** mean maximum abstraction, maximum code, big-tech infrastructure without need, perfect future-proofing or indefinite redesign.

### Essential vs accidental complexity

**Essential complexity** comes from real domain/system constraints: multi-tenancy, external-provider uncertainty, transactional consistency, authorization, asynchronous delivery, temporal facts, auditability, idempotency and contract evolution.

**Accidental complexity** comes from implementation choices: duplicate policies, parallel abstractions, hand-synced registries, unnecessary compatibility, redundant lifecycle code, speculative providers and repeated enforcement of the same property.

Remove accidental complexity; do not erase essential complexity.

### YAGNI

YAGNI removes speculative capability, not required correctness.

YAGNI may justify deleting unused extensibility, unsupported compatibility paths, future-vendor frameworks, configuration nobody sets, abstractions with no real boundary and duplicate enforcement.

YAGNI does **not** justify deleting an invariant, fail-closed default, proof for a reachable failure mode or a boundary that prevents a known defect class.

## Enforcement hierarchy

Prefer the strongest reasonable mechanism:

1. **Structure / public API makes the invalid state unrepresentable**
2. **Type system makes it invalid**
3. **Database/schema constraint makes it invalid**
4. **Runtime boundary fails closed**
5. **Test proves reachable behavior**
6. **Static guard/lint detects the violation**
7. **Documentation/convention**

This is a preference order, not dogma. A lower layer is valid when a stronger layer cannot express the complete property without disproportionate cost or destroying legitimate flexibility. Record why.

## Required decision flow

### 1. Observe and reproduce

State the actual evidence. Use the correct measurement universe. Do not start from the proposed patch.

For integration/runtime claims, distinguish:

- `contract_validated` — local behavior/types/fakes proved;
- `integration_validated` — real dependency/runtime proved;
- `blocked_for_real_validation` — exact blocker named.

A fake cannot prove credentials, network, provider semantics, database policy or deployment behavior.

### 2. Identify the root cause

Ask:

- What structural fact made this possible?
- Can the same fact produce other symptoms?
- Are several findings evidence of one shared cause?
- Would the finding still matter if the current directory/package layout were the opposite?

Repeated findings around the same construct trigger mandatory local-vs-global review.

### 3. State the target property

Write the invariant independently of current implementation.

Bad: “move `connectors` into `adapters`.”  
Good: “provider-specific protocol knowledge is unreachable from business contexts; consumer-owned ports are the only dependency from business semantics toward an external provider.”

Bad: “replace this SQL.”  
Good: “one context is the only write authority for this business state.”

### 4. Name authority and boundary

Identify:

- who owns the business meaning;
- which artifact is the source of truth;
- which component owns lifecycle/writes;
- which public contract is legal;
- the strongest reasonable enforcement boundary.

A solution that creates a second authority is presumed wrong until justified.

### 5. Evaluate credible candidates

For each candidate answer:

- Does it remove the root cause or only this symptom?
- What defect class remains representable afterward?
- What complexity does it add now?
- What complexity does it avoid later?
- Is the added complexity essential or accidental?
- Does it preserve a legacy seam only because the seam already exists?

### 6. Choose one legal outcome

Exactly one:

1. **Restructure now** — implement the global-maximum structure in the current work.
2. **Transitional solution** — bounded local maximum with named successor/deletion gate.
3. **Stop and split prerequisite** — the correct fix crosses the current boundary; do not patch around it.
4. **Current structure confirmed** — architecture is sound and a local correction is appropriate.

### 7. Define proof before implementation

Specify how the target property will be demonstrated:

- RED/counterfactual proof when useful;
- GREEN positive behavior;
- type/compile failure for illegal dependency/state;
- schema/constraint proof;
- generation/diff check for derived artifacts;
- integration/live proof when an external behavior is claimed;
- broader regression only when the changed boundary warrants it.

Proof validates the property, not a spelling of the implementation.

### 8. Implement and simplify

After correctness is established, subtract:

- duplicate paths;
- obsolete compatibility;
- transitional mechanisms whose successor landed;
- dead abstractions;
- redundant authorities;
- guards that no longer protect a distinct reachable property.

### 9. Close with evidence

Complete means:

- root-cause disposition explicit;
- target property holds;
- authority is unambiguous;
- relevant proof is green/non-vacuous;
- no known material contradiction remains;
- transitional debt has a named exit;
- findings are resolved, disproved or explicitly deferred with ownership.

## Structural inversion test

For every structural conclusion, state:

> **What part of this conclusion would still be true if the current implementation were the opposite in every respect?**

Current directory shape, migration cost, old ADR wording, import graph or schema are valid **current-state evidence** and sequencing inputs. They are not sufficient arguments that a target boundary is correct.

Target structure must be justified by domain semantics, user-observable behavior, named failure modes, system constraints and enforceability.

## Control versus artifact

A control that exists but does not fire on the path a change actually takes is absent.

Examples:

- a script invoked by no required path;
- a database policy bypassed by the runtime role;
- a generated file with no regeneration/diff gate;
- a test file outside test discovery;
- a guard with no negative fixture proving it can fail.

Do not count artifacts; prove effects.

## Guard / lint / verifier policy

A custom guard is justified only when all are true:

1. it protects a material property;
2. the failure is reachable or has occurred;
3. a stronger structure/type/schema/runtime mechanism cannot reasonably express the complete property yet;
4. a standard existing tool does not already enforce it;
5. maintenance cost is lower than recurring defect risk;
6. a negative fixture proves the guard can fire.

Repeated syntax-specific patches to a guard are a signal to re-evaluate the enforcement boundary rather than indefinitely hardening spelling recognition.

## Transitional enforcement

Every transitional mechanism records:

- property protected now;
- why the global maximum cannot land in the current slice;
- named successor;
- milestone/stage that removes it;
- deletion as part of that successor’s definition of done.

Without an exit condition, treat it as permanent architecture and hold it to the permanent bar.

## Data/contract/integration rules for this repository

During the architecture rebaseline, these are cross-cutting reasoning constraints:

- unknown is not zero/default;
- provider payload is not domain data merely because it contains useful fields;
- an external identifier is not automatically canonical identity;
- a read model/projection is not an authority;
- a successful provider HTTP response is not necessarily convergence;
- timeout after a potentially accepted write is an ambiguous outcome, not automatic failure/retry;
- polling, webhook, cursor and freshness semantics belong to the business capability that understands what “complete/current” means; generic platform code owns only mechanics;
- frontend convenience does not justify a second business policy/contract authority;
- current migrations/tables/routes are evidence to inventory, not obligations to preserve when there are no production compatibility constraints.

## Engineering Decision Record

Use a proportional written record for non-trivial decisions:

```markdown
## Engineering Decision Record

### Symptom / evidence

### Root cause

### Target property

### Authority and boundary

### Local-maximum candidate

### Global-maximum candidate

### Decision
Restructure now | Transitional solution | Stop and split prerequisite | Current structure confirmed

### Enforcement
Strongest reasonable layer and why.

### Proof
RED/GREEN/type/schema/static/integration evidence.

### Transitional exit
Successor + deletion condition, or N/A.
```

Formal ADR/spec is needed when the decision changes durable architecture, contracts, ownership, security, persistence or cross-context semantics. Tiny factual/mechanical work can use a short record.

## Review and convergence

Review is adversarial but bounded.

- A finding is evidence to verify, not an instruction to obey automatically.
- Repeated findings at the same architectural altitude trigger root-cause analysis, not another patch round.
- Optional hardening does not become blocking without a material property at risk.
- Convergence between reviewers is weak evidence; unexplained divergence is evidence that a question remains unsettled.

Stop reviewing when root cause, target property, authority, boundary, trade-offs and proof are settled; remaining findings are mechanical/non-material; and no material contradiction remains.

**Global maximum is not permission for endless perfection search. Reopen a settled decision only on a new material finding or changed constraint.**