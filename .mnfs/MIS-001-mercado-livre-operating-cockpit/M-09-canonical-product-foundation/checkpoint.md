# M-09 Terminal Checkpoint

- status: `externally_blocked`
- milestone task: `019f5d00-0b82-7b61-9920-32c7bd490333`
- hub task: `019f5cf6-8c9f-7321-ba07-f5b5b5e6bc77`
- frozen SHA: `a1d4aedd9cd4ccd41550d9233a35d9795073e7b7`
- normalization commit: PASS; only the three authorized wording files and
  directly required evidence/checkpoint changed.
- fixed-SHA review: PASS with no findings.
- M-09-C03: PASS; the exact active-path residue scan has zero matches.
- reused evidence: targeted/full Go, SDK/OpenAPI 40/40, governed runner contract
  14/14, and M-09-C01/C02/C04 remain accepted because the normalization diff is
  non-executable.
- M-09-C05: BLOCKED before Oracle connection. The governed runner remained in
  `docker build`; no live Oracle test container or authentication attempt was
  reached.
- runtime diagnosis: `Invoke-LiveOracleDockerProcess` redirects stdout and
  stderr, drains stdout to EOF before reading stderr, and has no timeout. Docker
  build progress can fill stderr while PowerShell waits for stdout, producing
  the observed build-phase deadlock.
- safety: the stuck verifier and its exact Docker client/wrapper chain were
  stopped; no Oracle/provider/database write, credential rejection/exposure,
  source edit, SHA drift, or destructive Git action occurred.
- evidence:
  - `normalization/validation.md`
  - `_fixed-sha-qa/deterministic-qa.md`
  - `qa-inventory-clock/validation.md`
- authorization boundary: the Portfolio packet forbade Oracle-runner changes;
  no runner correction was attempted.
- next owner: Portfolio Hub to authorize a bounded runner I/O/timeout correction
  and contract proof, followed by fixed-SHA review and the governed read-only
  Oracle C05 lane at a new frozen SHA.
