SLICE CARD — F03-S1 · feature F-03-shared-primitives · milestone M-03

depends_on: [F01-S2] (theme layer committed; semantic token utilities available from index.css @theme).

goal: Add the canonical IC-04 `MarginChip` primitive: a stateless margin classifier that renders the margin
percentage in a pill colored by a two-threshold band, with HONEST unknown handling (never renders 0% or a
green/healthy affordance for missing/invalid input). No global state, no data fetching.

complexity: standard.

write_set (create ONLY these two files):
- packages/ui/src/MarginChip.tsx        (new — the component; NOT yet exported from the barrel)
- packages/ui/src/MarginChip.test.tsx   (new — failing-test-first, then green)

NOTE: do NOT touch packages/ui/src/index.ts in this slice — the barrel export is added later in F03-S6 (single
barrel writer). MarginChip is importable by path for its test.

PUBLIC API (fixed — do not add props, do not add a second abstraction):
```ts
interface MarginChipProps {
  marginPct: number | null;
  thresholds?: { healthy: number; tight: number };
}
export function MarginChip(props: MarginChipProps): JSX.Element
```

BEHAVIOR (exact):
- Default thresholds: `healthy = 18`, `tight = 10`. Classification of a FINITE numeric `marginPct`:
  - healthy  ⇔  `marginPct >= healthy`         (inclusive; default ≥18)
  - tight    ⇔  `marginPct >= tight && marginPct < healthy`   (default ≥10 and <18)
  - warn     ⇔  `marginPct < tight`            (default <10)
- UNKNOWN band ⇔ `marginPct` is `null`, `NaN`, `+Infinity`, or `-Infinity`. Unknown renders the neutral
  placeholder `—` (em dash) and MUST NEVER render `0%`, a number, or a healthy/green affordance. This is the
  integrity rule — fail honest.
- Custom valid thresholds override the defaults. A thresholds object is VALID only if both `healthy` and
  `tight` are finite numbers AND `healthy > tight`. If thresholds is provided but INVALID (non-finite value,
  or `healthy <= tight` inverted), FALL BACK to the defaults (18/10) AND emit a development-only warning
  (guard with `import.meta.env.DEV` so it is stripped from production). Do not throw.
- Rendered text: healthy/tight/warn → the margin followed by `%` (e.g. `25%`, `18%`, `12%`, `3%`); unknown → `—`.

STYLING (use the semantic token utilities from index.css @theme — NEVER raw Tailwind palette like emerald/red):
- pill shape: `rounded-pill`, small inline-flex, mono or sans per the other chips.
- healthy → accent tokens (e.g. `bg-accent-soft text-accent-ink`).
- tight   → amber tokens  (e.g. `bg-amber-soft text-amber`).
- warn    → warn tokens   (e.g. `bg-warn-soft text-warn`).
- neutral/unknown → muted/faint tokens (e.g. `bg-surface-2 text-muted` or `text-faint`), visually distinct
  from all three value bands.
Read packages/ui/src/Badge.tsx for the STRUCTURAL idiom of a small pill primitive (className composition,
export style) and match it — but use the SEMANTIC token classes above, not Badge's current raw-palette classes.
Do not duplicate a class-merge helper if one already exists in packages/ui/src (cite path:line and reuse).

REQUIRED TESTS (MarginChip.test.tsx — @testing-library/react render; failing-first then green):
- `marginPct={null}` → renders `—`, does NOT render `0%` or any digit, carries the neutral class (assert a
  distinct neutral marker, NOT a healthy one).
- `marginPct={25}` → `25%`, healthy band (accent).
- `marginPct={18}` → `18%`, healthy (boundary inclusive ≥18).
- `marginPct={12}` → `12%`, tight band (amber; ≥10 <18).
- `marginPct={3}`  → `3%`, warn band (warn; <10).
- `marginPct={NaN}` and `marginPct={Infinity}` → neutral `—`, never `0%`.
- custom valid thresholds (e.g. `{healthy:30, tight:20}`) reclassify (e.g. 25 becomes tight, not healthy).
- inverted/non-finite thresholds (e.g. `{healthy:10, tight:20}`) → fall back to defaults (25 → healthy) and a
  dev warning is issued.
Assert BAND via semantic class selection (e.g. the element's className contains the accent/amber/warn/neutral
token class), not hardcoded hex.

validation_kind: unit test / typecheck.

commands (run from worktree root; capture output — SANDBOX NOTE in bindings applies to vitest):
- `npm run test -w @marketplace-central/web -- ../../packages/ui/src/MarginChip.test.tsx`   (all cases green)
- `npx tsc --noEmit -p apps/web/tsconfig.json`   (expect ONLY the known pre-existing TS2688-'node' env error;
   anything else is yours)

expected_artifacts: MarginChip.tsx with the fixed API + honest-unknown + two-threshold-with-fallback behavior
using semantic token classes; MarginChip.test.tsx covering null/25/18/12/3/NaN/Infinity/custom/inverted; green
vitest for that file; tsc clean except the known env error. Barrel NOT touched.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run)
and captured output; confirm barrel/index.ts NOT touched and no raw-palette classes used; anything you did NOT
verify. Note (G2) any non-trivial decision (e.g. the class-merge approach) with a 1-line alternative.
