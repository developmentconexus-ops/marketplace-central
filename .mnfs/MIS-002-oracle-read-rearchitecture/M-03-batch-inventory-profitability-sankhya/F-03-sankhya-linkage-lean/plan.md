# F-03 Sankhya linkage lean — execution plan

1. Add a local, uncommitted `oraclebatch.Chunks` stub matching F-01's fixed
   signature so this worktree can compile; never include it in the commit.
2. Refactor `FindCandidates` and `ListDescendants` to retain input/config
   validation only and remove request-time `ValidateConfiguration` calls.
3. Replace per-candidate line reads with chunked IN-list reads, merge rows by
   document ID, and retain the existing line-limit and identity checks.
4. Route every raw Oracle cause in the Sankhya adapter through
   `safeOracleCause`; add count, no-validation, 700-ID chunk, startup, and
   redaction regression tests. Assert there is no linkage freshness/cache
   parameter or path.
5. Wire one startup validation call in the composition root, run targeted
   tests, grep evidence, then the full server test suite. Record exact command
   output in `validation.md` and commit only feature files.
