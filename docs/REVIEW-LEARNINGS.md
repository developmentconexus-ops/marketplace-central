# Review Learnings — marketplace-central

Reviewer false-positive corrections (REVIEW-STANDARD.md §10). Loaded into every code-review
dispatch. Team-wide patterns ONLY — never one-off exceptions. One line each:
`YYYY-MM-DD · pattern · why it is fine here`.

- 2026-07-15 · godror "timezone" WARNING at Oracle connect · not a failure — signals a LIVE
  session (config parse + auth succeeded); ping-failed-count=0 is the health signal.
- 2026-07-15 · HPG_MIGRATION_FAILED with migrations_first=-1 in fresh worktree · build died on
  empty .gomodcache under GOPROXY=off, not a SQL/migration defect — warm the cache first.
- 2026-07-15 · 3D000 "database does not exist" on first CREATE DATABASE in integration lane ·
  postgres first-boot init restart race; the lane's retry loop absorbs it.
- 2026-07-16 · pgx `**string` (pointer-to-pointer) scan targets for nullable columns · the
  module's established nullable-scan idiom (ADR-17 honest-null), not an error — M-01 slice-1
  reviewer flagged it as a defect, dismissed with evidence.
- 2026-07-16 · `page_size` reported as `len(page.Items)` (actual returned count, not the
  requested limit) · IC-02-mandated response convention — M-01 slice-4 reviewer FP class.
