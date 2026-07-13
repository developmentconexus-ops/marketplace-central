# M-09 C05 Runner Correction Checkpoint

- status: `runner_correction_complete_post_commit_oracle_proof_pending`
- milestone task: `019f5d00-0b82-7b61-9920-32c7bd490333`
- runner correction base: `5da9e774bd3acb0ebc6a72b6741a52c5570c7847`
- correction scope: narrow the governed Oracle runner from the full live smoke
  test to the exact product lookup subtest required by M-09-C05.
- selector: `^TestOracleLiveSmoke$/^product_lookup$`.
- proof: governed runner Pester contract 14/14 PASS.
- post-commit governed Oracle evidence target:
  `F-02-oracle-catalog-cutover/_fixed-sha-oracle-evidence.md`; the Milestone
  refreshes this sanitized external/uncommitted file after freezing the runner
  correction SHA.
- limitations: no guessed mapping or provider/Oracle write occurred.
- next: commit this runner correction, freeze that exact SHA, then let the
  Milestone run the governed read-only Oracle lane once and refresh the named
  evidence before fixed-SHA review and proportional QA.
