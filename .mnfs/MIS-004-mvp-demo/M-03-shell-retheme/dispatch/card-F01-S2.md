SLICE CARD — F01-S2 · feature F-01-theme-tokens-fonts · milestone M-03

depends_on: [F01-S1] (tokens/@theme already in apps/web/src/index.css on this branch — verify with
`rg -n '@theme' apps/web/src/index.css` before starting; if absent, STOP and report BLOCKED).

goal: Apply the saved theme to `document.documentElement` BEFORE first paint via a synchronous inline
`<head>` script, then expose ONE guarded `useTheme` controller for the header to read/toggle. Default to
`light` under missing, corrupt, or inaccessible storage. Do NOT touch fonts at all in this slice — do NOT
remove the Google Fonts CDN links, do NOT add @font-face or font imports. (The CDN removal + self-hosted
@fontsource wiring is the SEPARATE next slice F01-S3, which is the single writer for those index.html font
links. Leaving the CDN links in place here is intentional and correct.)

complexity: complex (boot path runs before React and must survive storage security/quota failures).

write_set (create/edit ONLY these four files):
- apps/web/index.html                         (edit: add pre-paint boot script ONLY — do NOT touch font links)
- apps/web/src/app/theme/theme.ts             (new: pure theme helpers — single source of truth)
- apps/web/src/app/theme/useTheme.ts          (new: guarded React hook, no provider)
- apps/web/src/app/theme/theme.test.ts        (new: unit + hook tests)

CURRENT STATE (verified):
- apps/web/index.html: `<html lang="pt-BR">`; `<head>` has 3 Google Fonts links (preconnect googleapis,
  preconnect gstatic crossorigin, stylesheet css2?family=Inter...&family=JetBrains+Mono...). Body loads
  the app via `<script type="module" src="/src/main.tsx"></script>`. There is NO existing theme mechanism
  and NO `data-theme` attribute today.
- apps/web/src/app/theme/ does NOT exist — you create it.
- Deps present (do NOT add any): react ^19.2.0, @testing-library/react ^16.3.0 (has `renderHook`),
  jsdom, vitest ^3.2.4.

REQUIRED BEHAVIOR:

1) Storage contract (fixed): key = `marketplace-central-theme`. Valid stored values: exactly `"light"` or
   `"dark"`. Type `Theme = 'light' | 'dark'`. Default = `'light'`.

2) theme.ts — the single source of truth. Export at minimum:
   - `STORAGE_KEY = 'marketplace-central-theme'` and `type Theme = 'light' | 'dark'`.
   - `readStoredTheme(): Theme` — reads localStorage; returns the stored value ONLY if it is exactly
     `'light'`/`'dark'`; on absent, malformed, unsupported, OR `getItem` throwing → returns `'light'`.
   - `applyTheme(theme: Theme): void` — sets `document.documentElement.setAttribute('data-theme', theme)`.
   - `persistTheme(theme: Theme): void` — writes localStorage; if `setItem` throws (quota/security), it
     MUST swallow and return WITHOUT throwing (a failed persist may not break the live theme change).
   - `getInitialTheme(): Theme` — derive the hook's initial state from the ALREADY-BOOTSTRAPPED root
     attribute (`document.documentElement.getAttribute('data-theme')`), falling back to `readStoredTheme()`
     then `'light'`. The hook must NOT be a second independent source of truth.

   NOTE ON try/catch (read before flagging as slop): the storage getItem/setItem guards here are REQUIRED
   by spec — theme is a non-integrity-critical preference where "fail to light / keep live DOM" IS the
   correct behavior. This is NOT the forbidden "blanket recover on an integrity read"; unknown-money rules
   do not apply to a UI theme preference. Keep the catch scoped to the single storage call.

3) Inline pre-paint boot script in index.html `<head>`, placed so it runs synchronously BEFORE
   `/src/main.tsx`. It reads `localStorage.getItem('marketplace-central-theme')`, and sets
   `document.documentElement.setAttribute('data-theme', v === 'dark' ? 'dark' : 'light')` (any non-`'dark'`
   value, or a thrown getItem, yields `'light'`). It must be a plain classic inline `<script>` (NOT
   `type="module"` — module scripts defer and would paint first). It may NOT import theme.ts (a pre-paint
   inline script cannot use ES imports without deferring). The resulting tiny duplication of the key string
   + light/dark check is deliberate and justified — call it out in your report; do not try to "dedupe" it by
   converting the boot to a module.

4) useTheme.ts — a single guarded hook: `useTheme(): { theme: Theme; setTheme(t: Theme): void; toggleTheme(): void }`
   (or equivalent). It initializes state from `getInitialTheme()`, and on change calls `applyTheme` (live DOM,
   always) and `persistTheme` (best-effort). No React context/provider, no second theme store.

5) Do NOT modify the Google Fonts `<link>` elements — they stay in place this slice (F01-S3 owns their removal).
   Your ONLY index.html change is inserting the pre-paint boot script into `<head>`.

validation_kind: unit test / grep-assertion / typecheck.

REQUIRED TEST COVERAGE (theme.test.ts) — use vitest + jsdom; mock localStorage where needed:
- saved `'dark'` → readStoredTheme returns `'dark'`.
- saved `'light'` → returns `'light'`.
- absent (null) → returns `'light'`.
- malformed/unsupported value (e.g. `'blue'`, `''`) → returns `'light'`.
- `getItem` throws → returns `'light'` (no throw escapes).
- `persistTheme` when `setItem` throws → does NOT throw.
- hook (`renderHook`): a theme change updates `document.documentElement`'s `data-theme` IMMEDIATELY, and the
  value survives re-initialization (reconstruct the hook / re-read getInitialTheme → same theme).

commands (run from worktree root; capture output):
- `npm run test -w @marketplace-central/web -- src/app/theme/theme.test.ts`   (all cases green)
- `npx tsc --noEmit -p apps/web/tsconfig.json`                                (clean)
- `rg -n 'data-theme|localStorage|marketplace-central-theme' apps/web/index.html apps/web/src/app/theme`   (present)
- `rg -n 'fonts\.googleapis' apps/web/index.html`   → EXPECT 1 MATCH still present (CDN removal is F01-S3, not this slice)

expected_artifacts: pre-paint inline boot script in index.html head; theme.ts pure helpers with the exact
storage key and light-default guards; useTheme hook initializing from the bootstrapped root attribute with
no provider; theme.test.ts covering all negative-storage cases + immediate-DOM-update + survives-reload;
clean tsc; Google Fonts links UNTOUCHED (removed later by F01-S3).

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run)
and captured output; note the deliberate boot-script duplication and the spec-required storage try/catch;
anything you did NOT verify.
