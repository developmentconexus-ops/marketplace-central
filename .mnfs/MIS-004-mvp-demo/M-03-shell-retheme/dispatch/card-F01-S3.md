SLICE CARD — F01-S3 · feature F-01-theme-tokens-fonts · milestone M-03

depends_on: [F01-S2] + a CHIP-OWNED dependency add (already done before you start).

PRECONDITION (chip guarantees before dispatch — verify, do not perform): the two packages
`@fontsource/instrument-sans` and `@fontsource/ibm-plex-mono` are ALREADY in apps/web/package.json and
installed in node_modules (hub grant D-03). Verify with:
`node -e "require.resolve('@fontsource/instrument-sans/400.css'); require.resolve('@fontsource/ibm-plex-mono/400.css'); console.log('ok')"`
If that fails, STOP and report BLOCKED — do NOT run `npm install`, do NOT edit package.json or the lockfile
(dependency management is chip-owned; those files are frozen for you).

goal: Wire the self-hosted @fontsource fonts into the app bundle so the "Instrument Sans" / "IBM Plex Mono"
families (already referenced by the CSS variables from F01-S1) resolve locally, and remove the Google Fonts
CDN so the offline demo makes ZERO external font requests.

complexity: standard.

write_set (edit ONLY these two files):
- apps/web/src/main.tsx     (add @fontsource weight-CSS imports)
- apps/web/index.html       (remove the 3 Google Fonts <link> elements)

HARD CONSTRAINTS:
- Do NOT edit apps/web/package.json or the lockfile (chip owns the dep add; already committed).
- Do NOT run npm install / npm ci.
- Do NOT edit apps/web/src/index.css (F01-S1 owns it; the family vars are already correct there).
- Do NOT add @font-face by hand — @fontsource ships the @font-face rules; you only import its CSS.

REQUIRED CHANGES:

1) apps/web/src/main.tsx — add exactly these six imports (weights the design uses: Instrument Sans
   400/500/600/700, IBM Plex Mono 400/500), grouped immediately ABOVE the existing `import "./index.css";`
   on line 1:
```
import "@fontsource/instrument-sans/400.css";
import "@fontsource/instrument-sans/500.css";
import "@fontsource/instrument-sans/600.css";
import "@fontsource/instrument-sans/700.css";
import "@fontsource/ibm-plex-mono/400.css";
import "@fontsource/ibm-plex-mono/500.css";
```
   (@fontsource static packages register `@font-face` for family names exactly "Instrument Sans" and
   "IBM Plex Mono", matching the `--font-sans`/`--font-mono` vars set in index.css. No other change to main.tsx.)

2) apps/web/index.html — DELETE all three Google Fonts `<link>` elements from `<head>`:
   - `<link rel="preconnect" href="https://fonts.googleapis.com" />`
   - `<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />`
   - `<link href="https://fonts.googleapis.com/css2?family=Inter...&family=JetBrains+Mono...&display=swap" rel="stylesheet" />`
   Leave the pre-paint theme boot script (added by F01-S2) and everything else intact. Add nothing to replace
   the links — the fonts now come from the bundle via step 1.

validation_kind: build / typecheck / grep-assertion (no new unit test — the bundle build + no-CDN grep ARE
the verification; do not fabricate a test).

commands (run from worktree root; capture output):
- `npx tsc --noEmit -p apps/web/tsconfig.json`                         (clean)
- `npm run build -w @marketplace-central/web`                          (must succeed; proves the .woff2 assets
   resolve and bundle — inspect output for hashed font assets / no resolution error)
- `rg -n '@fontsource/(instrument-sans|ibm-plex-mono)/[0-9]{3}\.css' apps/web/src/main.tsx`   (6 matches)
- `rg -n 'fonts\.googleapis|fonts\.gstatic|Inter:wght|JetBrains\+Mono' apps/web/index.html`   → EXPECT NO MATCHES
- `rg -n 'data-theme' apps/web/index.html`   → boot script from F01-S2 STILL present (you did not remove it)

expected_artifacts: six @fontsource weight imports in main.tsx above ./index.css; zero Google Fonts links in
index.html; green tsc; green web build bundling the local font assets with no external font request; the
F01-S2 boot script untouched.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run)
and captured output; confirm you did NOT touch package.json/lockfile/index.css; anything you did NOT verify.
