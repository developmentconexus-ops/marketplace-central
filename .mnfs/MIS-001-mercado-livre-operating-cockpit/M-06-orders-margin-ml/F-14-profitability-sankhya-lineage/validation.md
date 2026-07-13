# F-14 Profitability Sankhya Lineage Validation

## Verdict

Feature validation passed for the bounded fake/unit target. Profitability now
uses one canonical stable MPC line to resolve the current confirmed TOP 313
origin mapping, re-read its exact existing TOP 306 descendants, and invoke the
exact-line tax port once per valid unique descendant. TOP 313 is never passed
to tax. One-to-many components aggregate only when known for every exact
descendant; partial or otherwise incomplete facts retain known amounts where
possible but remain incomplete, and unknown amounts remain nil.

## Context evidence

- Accepted base SHA: `01cabcb0683ed6ace3d56a169c76a64c90b6b2d0`.
- Canonical L2 context compiled with the dispatched Feature contracts,
  validation contract, F-07/F-09/F-11/F-12 specifications, `portfolio-core`,
  `orders-margin`, exact allowed paths, composition shared seam, registered
  command IDs, side effects, and stop conditions.
- Estimated input is 5,586 tokens; every required source carries the compiler's
  explicit L2 overflow reason.
- `context-validate -RequireCurrentBase` passed after implementation because
  every compiled contract source hash and the accepted base remained current.
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-14-profitability-sankhya-lineage/context.json`.

## Implementation evidence

- Orders application adds a read-only current-lineage operation scoped by
  exact tenant, installation, provider order, and stable MPC line. It loads the
  current exact-evidence confirmation, selects exactly one persisted origin,
  validates F-11 configuration, re-reads the uniquely selected persisted TOP
  313 candidate and exact origin line, and passes that line's nullable quantity
  to the descendant reader. A nil source quantity remains nil rather than being
  invented.
- Missing confirmation/mapping or a non-stable line stays `none`; malformed or
  contradictory persistence/lineage stays `conflict`; configuration/source
  failures stay `unavailable`; reader `partial` and `complete` remain intact.
  Focused tests assert no `AppendConfirmation` call.
- Profitability order facts now carry the opaque MPC line ID and reconciliation
  state through a profitability-owned adapter. Blank, malformed, duplicate,
  legacy, or ambiguous line identity cannot resolve tax.
- A profitability-owned lineage port and orders adapter isolate orders DTOs.
  Composition supplies it only when the assisted-linkage service was
  successfully constructed; otherwise existing missing-tax behavior remains.
- Profitability validates the full descendant set before any tax call, rejects
  empty/invalid/duplicate identities, sorts by document/line, and calls the
  existing exact tax port once for each descendant. Returned source identity
  must equal the requested identity.
- Tax provenance records every sorted exact descendant. Each component is
  summed only if non-nil for every descendant. A known partial-lineage sum is
  retained with `partial` quality; any non-universal component is nil/missing,
  so item and order margin remain incomplete.
- No migrations, Oracle adapter, public HTTP/OpenAPI/SDK, UI, Docker, runtime
  configuration, live database/provider, dependency, secret, or PII operation
  changed or ran.

## go-orders-profitability-lineage

Command from `apps/server_core` with repository-local `GOCACHE=.gocache`:

```text
go test -count=1 \
  ./internal/modules/orders/application \
  ./internal/modules/profitability/... \
  ./internal/composition
```

Result: exit 0.

Focused proof covers exact scoped current lineage, no-write resolution, all
five states, stable line projection, exact sorted TOP 306 calls, deterministic
one-to-many provenance, all-descendant component aggregation, partial-known
incompleteness, invalid/duplicate/conflict/unavailable fail-closed behavior,
adapter conversion, existing profitability transport behavior, and composition
build.

## go-repository

Command from `apps/server_core` with repository-local `GOCACHE=.gocache`:

```text
go test -count=1 ./...
```

Result: exit 0. All server commands and Go packages compiled; all repository Go
tests passed without a live integration target.

## git-diff-check

- `git diff --check`: exit 0; Git emitted only existing LF/CRLF conversion
  warnings.
- Changed Feature-owned paths are limited to the dispatched F-14 directory,
  the two exact orders application files, profitability module, and
  composition root.
- Pre-existing untracked `docs/research/**` and `output/**` remain unmodified,
  unstaged, and excluded from this Feature.
- The final index contains only the allowed F-14 paths and passes
  `git diff --cached --check`.

## Acceptance mapping

- F14-AC01: orders focused tests prove exact read-only current lineage, all
  states, and zero ledger append calls.
- F14-AC02: order adapter/application tests prove stable MPC line propagation
  and no tax read for non-stable or unavailable identities.
- F14-AC03: profitability tests prove sorted one-call-per-exact-descendant
  behavior and deterministic multi-line TOP 306 provenance.
- F14-AC04: profitability tests prove component-wise all-descendant sums, nil
  unknowns, and known partial amounts with incomplete quality.
- F14-AC05: focused composition build plus repository-wide Go proof establish
  conditional wiring without public contract or runtime configuration change.

## Correction retry 1

- Original blocking failure: `ResolveCurrentLineage` passed nil expected
  quantity to every F-11 descendant read, so the real reader could only return
  `partial` and production could never complete tax/margin.
- Smallest fix: after configuration validation, re-read candidates by the
  persisted external key, require exactly one persisted TOP 313 header and one
  exact persisted origin line, then pass its nullable quantity to
  `ListDescendants`.
- Fail-closed proof: missing/duplicate selected candidates, absent origin, and
  mismatched candidate-line document return `conflict` before any descendant
  or tax read. A valid exact quantity is asserted at the fake reader boundary
  and permits `complete` lineage.
- Assigned blocking failure resolved: Yes.
- Retry use: 1 of 1 Feature correction attempts.
- Remaining assigned failures: none.

## Residual risk and next

- Fake/unit evidence proves deterministic contract behavior, not live Oracle
  data. F-11 remains the authority that each returned descendant was selected
  under exact destination TOP 306 predicates.
- Milestone review should inspect the fixed Feature commit, then proportional
  QA must decide milestone status. This Feature does not pass the milestone.
