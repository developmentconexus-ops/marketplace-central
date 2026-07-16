# Marketplace Central — Development Harness (pointer)

**Binding swapped 2026-07-16 at the M-01 milestone boundary** (per core §0 / harness-hub boot
rule). The combined doctrine that previously lived in this file is superseded. Binding is now
THREE layers, read together:

1. **Core (method):** [`docs/HARNESS-CORE.md`](HARNESS-CORE.md) — vendored copy of the
   `harness@mnfs-harness` plugin's `HARNESS-CORE.md`; canonical source
   `Documents/mnfs-harness/harness/HARNESS-CORE.md` (mnfs-harness commit `cd114e6`). The
   vendored copy is updated ONLY at mission/milestone boundaries, by the hub.
2. **Profile (this repo):** [`docs/HARNESS-PROFILE.md`](HARNESS-PROFILE.md) — commands, seams,
   collision-axis instantiation, non-negotiables, human gates, superseded-protocol denylist,
   amendment log.
3. **Mission (current queue):** `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md`
   — `## Parallel Execution Plan` + per-milestone `## Ownership & Concurrency` blocks.

Where the old sections went: §1/§2/§3-method/§4/§7/§8 → core · §5 commands, §6 repo
invariants, seams, denylist → profile · §3 MIS-003 DAG → mission.md Parallel Execution Plan.

Dispatch prompts pin the core + profile paths explicitly (core §2 item h). Skills:
`harness-hub` boots the hub, `harness-worker` binds dispatched sessions, `codex-dispatch`
resolves model/effort/path — all from plugin `harness@mnfs-harness`.
