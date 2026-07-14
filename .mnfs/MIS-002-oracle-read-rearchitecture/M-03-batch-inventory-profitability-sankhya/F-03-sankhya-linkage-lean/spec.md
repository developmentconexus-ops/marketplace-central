# F-03 Sankhya linkage lean — execution specification

## Contract

`SankhyaLinkageReader` keeps the existing `ports.SankhyaLinkageReader`
interface and typed error codes. Configuration validation (ping, metadata, and
uniqueness probe) is performed once by the composition root during startup;
`FindCandidates` and `ListDescendants` perform no validation probes. An
enabled, connected Sankhya runtime that fails validation makes startup fail
with the existing read-error shape and a safe Oracle cause.

Candidate lookup remains ordered by its existing query semantics. Candidate
line loading collects document IDs in candidate order, splits them with the
contract-fixed `oraclebatch.Chunks(ids, 500)`, and issues one IN-list query per
chunk. The query preserves per-document line ordering and the configured
per-candidate limit/overflow sentinel. Chunk failure fails the whole read.

All Oracle driver causes emitted by this adapter use `safeOracleCause`,
preserving cancellation/deadline causes and reducing driver errors to an
Oracle numeric code when available. Linkage reads have no freshness or cache
path; the reader interface does not accept `FreshnessPolicy`.

Canonical identity behavior from `97fd4b58` is outside this change and must
remain unchanged.

## Acceptance criteria

- Request-path candidate and descendant calls issue no Ping or metadata
  queries.
- Candidate line query count is `ceil(candidate_count / 500)` after the one
  candidate query; a 700-candidate proof issues one candidate query plus two
  line queries.
- Per-candidate line identity, ordering, nullable values, and overflow errors
  remain explicit and typed.
- DSN-shaped driver text is absent from returned errors and logs; only the
  numeric Oracle code may be exposed for coded driver errors.
- Linkage remains uncached and fresh on every read.
