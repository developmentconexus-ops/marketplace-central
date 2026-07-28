Status: Implemented and committed as `aceff011`.

Changed paths exactly match write set:

- [ImportChainPanel.tsx](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/sharp-pike-3387c1/apps/web/src/pages/importacoes/ImportChainPanel.tsx)
- [ImportChainPanel.test.tsx](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/sharp-pike-3387c1/apps/web/src/pages/importacoes/ImportChainPanel.test.tsx)

G1: Kept a second local helper. `renderCounter` has three numeric consumers; protocol has one distinct string-integrity consumer. Generalizing would add an abstraction without a second consumer.

G2: Considered inline guarding and generalizing `renderCounter`; rejected both per the card’s semantic-helper requirement.

G3: No named upcoming seam is blocked; this is isolated to the import-chain panel.

Evidence:

- Focused Vitest: `could-not-run (sandbox)` — esbuild reported `Access is denied` reading `vitest.config.ts`.
- Typecheck: ran; only the 15 documented pre-existing errors appeared, with none in touched files.
- `git diff --check`: ran successfully.
- Commit: ran successfully.

Not verified: focused unit tests and full Vitest suite, due to the sandbox tooling failure. Pre-existing untracked `.mnfs` evidence files were untouched.