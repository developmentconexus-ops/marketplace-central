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
