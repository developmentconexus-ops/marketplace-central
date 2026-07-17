# F-03 selection-preview-api — validation record

Chip: CHIP-M03 · branch `chip/m-03-mutation-envelope-writes` · closed 2026-07-16.

## Slices → commits → reviews

| Slice | Commits | Review (ledger) | Verdict |
| --- | --- | --- | --- |
| S1 shared IC-02 filter grammar + selection resolver | 184d11b0 + 2fc49a41 (G2 note) | D-42/D-42b | ACCEPT-W-C, closed (G2 dependency-direction note) |
| S2 create + preview services (ReplacePreview single-tx, source time = max FetchedAt) | f348cc07 + 21cc5346 (actor gate wired into Create) + dfc21f2f (G2 DRY note) | D-44/D-44b | ACCEPT-W-C, closed |
| S3 approve + cancel (execute strict bool, 15:00 boundary, SQL-predicate cancel) | c24e036a (scope amendment after honest BLOCKED D-45) + 06dc01e4 | D-45→D-46/D-47 | ACCEPT (redispatch after honest block) |
| S4 retry clone + read pagination (allocateProtocolID shared helper) | 1d6433d5 → ACCEPT-W-C → 7e2670d1 (direct typed error; allocation dedup) | D-48/D-49/D-49b | ACCEPT after fix |
| S5 command HTTP endpoints (5 POSTs, 14-code error table) | a900cca6 → REJECT → 66b0944a (domain.ProtocolTypeEnabled single authority + Create gate) + 44e3a1f (G2 note) | D-50/D-51/D-51b | ACCEPT after fix |
| S6 query HTTP endpoints (keyset cursors mirroring listings) | 5b26a2e1 → ACCEPT-W-C → ad512b10 (invalid_filter/invalid_cursor split + honest test name) | D-52/D-53/D-53b | ACCEPT after fix |
| S7 command contract (OpenAPI + SDK atomic, 5 POSTs) | faa377cf | D-54/D-55 | ACCEPT (no conditions) |
| S8 query contract (OpenAPI + SDK atomic, 3 GETs; exactly 8 IC-03 ops) | dc79483b + 3599e864 (clarification line) | D-56/D-57 | ACCEPT (no conditions) |
| S9 route registration + lifecycle/error-matrix integration | d2fff81a + 8045b92b (assertion key) + 7cbdee15 (honest lane fakes) → ACCEPT-W-C → ed8afcbd (listing_create filter e2e) | D-58/D-59/D-59b | ACCEPT after fix |

## Contract evidence (M03-C01–C09, C12 served)

- Eight IC-03 HTTP operations mounted and integration-proven: POST /mutations,
  /{id}/preview, /{id}/approve, /{id}/cancel, /{id}/retry-failures; GET /mutations,
  /mutations/{id}, /mutations/{id}/items. Route smoke on both lanes (root_test.go);
  full lifecycle create→preview→approve→poll→terminal over ephemeral PostgreSQL
  (lifecycle_test.go, stub writer lane).
- Selection: single IC-02 grammar authority `domain.SetFilterValue` shared by listings URL
  and mutation JSON entry points; resolver stops at 2001 (selection_too_large cap 2000);
  explicit + filter modes; strict JSON decode with trailing-content check.
- Preview: ReplacePreview single tx (row lock + state check + delete + insert + update);
  ADR-17 source time = max non-null FetchedAt, any missing → typed
  `source_time_unavailable`, never now(); caps enforced post-resolver pre-reads.
- Approve/cancel: only literal `execute:true` (strict *bool decode, unknown-field +
  trailing-EOF rejected); staleness strictly >15min (exactly 15:00 approves); cancel legal
  from draft|previewed only, enforced by tenant-scoped SQL predicate + RowsAffected==1
  (concurrency arbiter).
- Retry: eligibility partially_failed|failed_preserved with ≥1 item failure in
  {provider_rate_limited, provider_unavailable, stale_source}; clone = new immutable
  protocol (retried_from set, state=draft, totals 0, timestamps+source_as_of NULL,
  ZERO item/idempotency rows) — integration-proven; nothing_to_retry 409.
- Enabled-type single authority `domain.ProtocolTypeEnabled` (6 types; listing_create
  excluded): gates transport decode, Service.Create insert, query type filter
  (→ 400 invalid_filter, integration-asserted e2e), and contract enum documentation.
- Error model: full transport table integration-asserted (18 matrix subtests) —
  invalid_body/invalid_filter/invalid_cursor/invalid_intent/installation_required 400;
  actor_required/type_not_enabled/empty_selection/selection_too_large/execute_required/
  source_time_unavailable 422; preview_stale/invalid_state/nothing_to_retry 409;
  protocol_not_found 404; method_not_allowed 405; unknown → 500 internal_error, no leak.
  Envelope `{"error":{code,message,details}}` mirrors listings. 12 item-level failure
  codes proven through real HTTP item transcript.
- Snapshot immutability: after preview, underlying listing price change (→999) does not
  alter persisted item `before` (stays 89.00) — integration-proven.
- Pagination: keyset (created_at, protocol_id) DESC / items seq ASC; versioned JSON +
  strict base64 opaque cursors (round-trip + malformed/dup-filter rejection tested);
  limit default 50 / max 200 — all mirroring listings conventions.
- Tenancy: every repository query tenant-scoped; cross-tenant ReplacePreview rejected
  (integration case); detail/items path-id-only semantics honest (no fake tenant param).
- OpenAPI + SDK: S7/S8 atomic commits; contract mirrors implemented transport (code wins);
  exactly 8 operations; ordering documented matching SQL verbatim; SDK
  create/preview/approve/cancel/retryFailures/list/get/listItems with URL-encoding tests;
  vitest 51/51. Preamble + CHIP-SAT sections untouched (diff -U0 proven per slice).

## Lanes

- Unit: green every slice (build/vet + mutations+listings+composition sweeps `-count=1`),
  re-run cold by reviewers and orchestrator pre-pass.
- Integration: GREEN on session pg (fresh `mpc_test_*` dbs, 39 migrations) — S2 preview
  atomicity, S3 cancel round-trip, S4 retry/read tenant-scoped pagination, S9 full
  lifecycle + 18-subtest error matrix; D-59 reviewer re-ran independently on own fresh db.
- `-race`: unavailable on this machine (accepted env limitation, hub decides cross-machine
  run pre-mission-close).

## Live-write posture (security-critical)

Unchanged from F-02: stub lane only throughout F-03; `MPC_PROVIDER_WRITES_ENABLED` gate
untouched (D-59 verified gate tests intact); zero live Mercado Livre writes. Live ML
execution remains gated on explicit operator authorization via hub ESCALATION.

## Carries / hub-queue items out of F-03

- Envelope path leaves `source_as_of` NULL at previewed (IC-03/PD-01 gap) — D-44 finding.
- IC-02 live URL grammar accepts status superset (under_review/inactive/payment_required/
  not_yet_active) vs contract — contract-vs-code conflict, preserved live behavior.
- IC-03 doc-drift: "items cloned" wording vs implemented draft/no-items retry;
  status discrepancies (actor_required 422 vs 400; execute_required 422 vs 400;
  nothing_to_retry 409 vs 422; installation_not_found unimplemented).
- Legacy `POST /mutations/{id}/retry` still mounted by transport Register alongside
  canonical root-mounted `/retry-failures` (pre-existing production wiring) — cleanup.
- apply_repository.go:28-30 hand-rolled state guard: previewed→previewed legal in domain
  table, `domain.TransitionProtocolState` identical today — mechanical consolidation at
  next touch (D-59 corrected the earlier false rationale).
- preview/apply repository split: keep (different ownership — preview owns snapshot/source
  time; apply owns claims/outcomes/terminal) — D-58 recommendation, D-59 concurred.
- ValidateProtocol zero production callers (F-01 carry, still true).
- MutationErrorResponse allOf composition (style, D-55).
- Field findings: F-A index.lock recurrences #7–#13; F-C sandbox write roots; vitest
  sandbox `--configLoader runner` artifact; symbol-level collision lesson (file-disjoint ≠
  symbol-disjoint in same package); harness pg target loader requires literal loopback IP
  (`localhost` → HPG_TARGET_INVALID); PowerShell `$db?` variable-name parsing gotcha.
- F-04 (preview-confirm-ui) remains GATED on hub rebase trigger naming M-02 F-03.
