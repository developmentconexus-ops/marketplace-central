# HUB gate pack — CHIP-IMPORT-CHAIN — the markers the hub owns

Branch: `chip/import-chain` · merged HEAD `997c87a7` · code `1e4c297e` · gate ran on `c2229f26`

This pack carries ONLY the two markers that are the hub's to write, plus the string-level
re-verification of the conditions. The chip's own pack is the authority on the implementation and
is not recopied here (R-14): `.mnfs/MIS-006-integracao-fundacao/_chip-import-chain/EVIDENCE.md`.
The gate itself is recorded verbatim in `GATE-P6.md`, with both seats' raw output in
`p6-opus-r1.md` and `p6-sol-r1-raw.md`.

P6-DUAL-GATE: AGREEMENT

What the marker means here, stated so nobody has to infer it. Two reading seats ran on the same
frozen prompt, blind to each other — a cold Opus with no Bash and no Write, and GPT-5.6 Sol medium
under `codex exec --sandbox read-only`. Both returned **APROVA**. Neither raised a blocking
finding. They converged independently on the same ressalva (`queue_read_at` unguarded by type) and
both independently contradicted the prior cold review's claim that the field was already guarded.
The hub raised three conditions on that basis; all three are discharged in `1e4c297e` and
re-verified BY STRING at merge, not by line and not by the chip's word:

| # | Condition | String verified in the merged tree |
|---|---|---|
| 1 | `queue_read_at` type-guarded | `ImportChainPanel.tsx:27-33` — `typeof value !== "string" \|\| value.trim().length === 0` BEFORE `formatDateTime` |
| 2 | known zero locked | `ImportChainPanel.test.tsx:142-156` — `enfileirados: 0` → `toHaveTextContent("0")` AND `not.toHaveTextContent("—")` |
| 3 | 404 not attributed by status | `grep "status === 404"` under `pages/importacoes/` = **0 hits**; discriminator is `candidate.error === "import_not_found"` (`ImportChainPanel.tsx:12`), with a second test for a 404 carrying no body |

The chip ran the must-fail itself after its codex worker reported vitest `could-not-run (sandbox)`
— conditions 1 and 3 failed on the discriminating assertion before the fix, condition 2 passes in
both worlds by design because it is a lock, not a repair.

LIVE-VERIFIED: `/importacoes` and `/importacoes/:importId` driven in the browser against the live
dev stack on real ERP import data, hub-owned L2, filed at `7099d9f`. The chain rendered from the
server response; the screen's `vinculados = 0` was confirmed TRUE against direct SQL rather than
assumed, with the leading-zero subcount discarded by two independent measurements.

## What this pack does NOT claim

Three things stayed unproven and are inherited as declared limitations, not laundered:

1. **`vinculados` non-zero was never exercised.** The zero on screen is proven true; the path with
   a real count was not walked.
2. **Mount with literally zero installations is not proven.** The L2 drove a tenant WITH an ML
   installation and used the ungated routes' clean `href` as the observable — that shows the
   routing decision, not the no-marketplace state. Both reading seats say it works by reading
   (`InstallationProvider` always renders children; only `InstallationGate` blocks; `/importacoes`
   is outside the block). The hub concedes the Sol side is right that this is an indirect
   observable, against its own L2 write-up.
3. **Malformed payload against a conforming server is unreachable by construction** — OpenAPI
   declares all five fields `required` and `protocol` carries a `CHECK` constraint `NOT NULL`. The
   unit test is the only possible proof of that half and was accepted as such.

## Reviewer errors recorded, so nobody chases ghosts

Both are written up in `GATE-P6.md`. In short: the Sol side's `I7`/`I8` NOT-EVIDENCED rest on the
hub pointing it at the wrong tree (hub's instrument error, not the chip's), and the Sol side
FABRICATED a `file:line` citation for `apps/web/vitest.chip.config.ts`, a path that exists in no
tree — probable origin is `chip.md:113-114`, which ORDERS the file deleted, read as an
observation. Neither touches the date-coercion finding, which that seat proved by executing.

## I9 — discharged by the hub's executing seat

Neither reading seat can run anything, so both left I9 NOT-PROVEN for the same structural reason.
The hub ran the ladder at the chip's own HEAD and recorded it in `GATE-P6.md` §I9. First run under
the three-seat gate ratified in `docs/HARNESS-PROFILE.md` §11.
