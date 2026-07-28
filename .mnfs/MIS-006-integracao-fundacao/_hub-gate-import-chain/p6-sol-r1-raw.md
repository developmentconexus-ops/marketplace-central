screen coupled to the directory it was promoted out of — defensible (that file is not this chip's to move) but it is 
an unfinished promotion.
  
> ## Criteria I1..I9
  
  | # | Verdict | Proof |
  |---|---|---|
> | I1 route exists + gate decision justified | **PASS (code) / NOT-VERIFIABLE-BY-READING (behavior)** | 
`AppRouter.tsx:68-69` registers both paths outside the gate at `:78-86`; justification comment `:63-66`, verified true 
in attack 5. The contract also demands "mostre a tela montando sem instalação" — no test asserts router placement 
(tests render pages directly, `ImportacoesPage.test.tsx:33-42`), so the behavioral half needs the live drive. |
> | I2 section moved without losing function | **PASS (code+tests) / declaration NOT-EVIDENCED** | Component+test 
present at `pages/importacoes/ImportacaoSection.tsx` / `.test.tsx`; absent from `pages/vinculos/` (Glob). All 5 tests 
and all 14 base assertions preserved, +1 at `.test.tsx:89`. What `/vinculos` shows now is *behaviorally* answered 
(section retained at `IntegracoesPage.tsx:449`, new home `/importacoes`), but the contract requires that be 
**declared** — no EVIDENCE file exists to declare it. |
> | I3 VinculosPage two lines only | **PASS** | Line-by-line comparison vs the base copy in the main checkout: only 
the `:8` import and the `:159` render removed; all remaining lines identical at −1 offset. No CHIP-VINC-NEUTRO 
collision. |
> | I4 consumes `getErpImportChain`, not the old numbers | **PASS** | `useErpImportChain.ts:20` → `index.ts:1901` → 
`/erp/imports/{id}/chain`; counters bound only at `ImportChainPanel.tsx:45,51,57`; no derivation from `listErpImports` 
anywhere (attack 2). |
> | I5 ADR-17: absent → `—`, never zero | **PASS** | Guard `ImportChainPanel.tsx:15-17`; tests `.test.tsx:46-61` 
(absent field) and `:63-79` (null), both asserting `—` **and** `not.toHaveTextContent("0")` at `:58,76`. REPORT-level 
note, out of this chip's scope per validation-contract.md:119: the server always emits an integer 
(`erp_import/transport/http_handler.go:133-138`, OpenAPI `required: [protocol, importados, vinculados, enfileirados, 
queue_read_at]` at `:8080`), so a real-world `0` reaching this screen means the *backend* decided zero — the FE guard 
cannot distinguish that from "unknown". |
> | I6 chain-read error is honest | **PASS** | `ImportChainPanel.tsx:32-35` renders `erp-import-chain-error` and 
nothing else; 404 vs generic split at `:9-13,21-23` keyed on the shape the SDK actually throws (`index.ts:1715` throws 
`{status, error}`), so the 404 branch is real in production, not test-only. Tests `:81-102` assert the chain testid is 
null in both cases. |
> | I7 F-02 verified, not rebuilt | **FAIL (NOT-EVIDENCED)** | The code side is right — `IntegracoesPage.tsx:295-346` 
consumes `useActiveSourceQuery`/`useSetActiveSourceMutation` against `/config/active-source` (`index.ts:1896-1898`) 
and nothing was rebuilt. But the criterion is a *declaration with a verification result and a commit id*, and there is 
**no `EVIDENCE.md` in `_chip-import-chain/`** (Glob returns chip.md, validation-contract.md, DISPATCH-LEDGER.md, 
dispatches/ only). The contract's own words: "FAIL se você declarou 'já existe' sem ter verificado." Nothing is 
declared and no toggle-persists verification is recorded anywhere I can read. |
> | I8 F-03 resolved either way | **NOT-EVIDENCED** | By reading, "satisfied" is demonstrable: the screen needs 
`import_id` for the row link (`ImportacaoSection.tsx:148`) and `ErpImportSummary.import_id` already carries it 
(`erpImport.ts:32`); no contract/SDK change was made or needed. But the contract requires the chip to *say which field 
and where* — no artifact does. |
> | I9 no collateral damage | **NOT-VERIFIABLE-BY-READING / NOT-EVIDENCED** | `tsc` and `vitest` cannot be run by me, 
and no cited before/after counts exist: the two worker `.last.md` files say "Vitest: could-not-run (sandbox)" 
(`…410e32e1.last.md:13`, `…c0849f51.last.md:9`) and the 15 tsc errors are asserted but never listed. 
`DISPATCH-LEDGER.md:36-41` says the chip re-ran vitest and repaired 3 files, but gives no pass/fail counts. 
Separately, a live hazard: `apps/web/vitest.chip.config.ts` **exists on disk** and rewrites every 
`@marketplace-central/*` import (`:16`) despite its own line 4 and chip.md:113-114 requiring deletion before commit; 
it is not in the reconciliation list, so it is presumably untracked, but I cannot confirm that. |
  
> ## What I could NOT verify
  - Any `git` fact: the three commits `4b76a287` / `1bffdcfd` / `67e4a3d`, the true `git diff --name-only` against 
5441fe18, and therefore the completeness of the supplied reconciliation (absence of Go/migration/contract changes 
rests on that list, not on my own diff). I compared against the **main checkout working tree** 
(`C:/Users/leandro.theodoro/Documents/marketplace-central/…`), which the session snapshot puts at `bcab8269` — two hub 
`.mnfs`-only commits above the 5441fe18 floor, so a code-level ancestor for the files I compared, but not literally 
the base sha.
> - Whether `apps/web/vitest.chip.config.ts` is tracked at HEAD. The session-start snapshot showed it as `??` 
(untracked); I have no way to re-check. If it is tracked, I9 is a hard FAIL, not a hazard note.
  - `tsc -p apps/web --noEmit` result and the 15-error baseline; `vitest` counts before/after. Execution required.
> - The whole L2 (validation-contract.md:94-102): `/importacoes` on real data, chain numbers cross-checked against a 
direct DB query, non-existent id → honest error, active-source toggle persisting and invalidating globally, 
light/dark. This is the cash-in of the operator's chain-endpoint waiver (`:86-92`) and the exact place a 
composition-root/decorator wiring defect of the catalog-503 class would surface — no test in this diff can catch it, 
because every test mocks `useClient` directly (`ImportChainPanel.test.tsx:9-13`, `ImportacaoDetailPage.test.tsx:9-13`).
> - Whether the `—` path can ever fire against the real server: the backend and OpenAPI both declare the counters 
required and always-present, so I5's guard is untested against production reality by construction.


codex
The reading side is clean on the core wiring: real app composition is `ClientProvider → QueryClientProvider → AppRouter → InstallationProvider`; the new route reaches the same SDK client used at runtime, and `InstallationProvider` always renders children even when its marketplace-installation query is empty or errors. No lazy-import or context seam diverges from production.

I found no blocker. Two narrow honesty gaps remain only for contract-violating `200` payloads: finite-but-impossible counters are accepted, and `queue_read_at` can be coerced from a non-string. I’m reconciling those as ressalvas, plus separating genuinely verified criteria from evidence I could not independently execute.
Completed fixed-SHA code, contract, runtime-composition, SDK/server, move-integrity, ownership, and prior-evidence review. Preparing the required independent verdict.
  ✓ Read review skill and binding authority documents
  ✓ Inspect exact diff, move integrity, and ownership boundaries
  ✓ Trace runtime composition, routes, API hook, SDK/server shapes, and error states
  ✓ Evaluate prior-review claims and criteria I1–I9
  → Deliver independent gate verdict with evidence and unverifiable items
exec
"C:\\Program Files\\PowerShell\\7\\pwsh.exe" -Command "node -e \"for (const v of [['2026-07-18T12:00:00Z'], true, 1]) { const d = new Date(v); console.log(JSON.stringify(v), Number.isNaN(d.getTime()) ? 'invalid' : d.toISOString()) }\"; git diff --numstat 5441fe18 42d7c2d1 -- apps/web/src/pages/vinculos/ImportacaoSection.test.tsx apps/web/src/pages/importacoes/ImportacaoSection.test.tsx; "'$base=(git show 5441fe18:apps/web/src/pages/vinculos/ImportacaoSection.test.tsx | Select-String '"'expect\\(').Count; "'$head=(git show 42d7c2d1:apps/web/src/pages/importacoes/ImportacaoSection.test.tsx | Select-String '"'expect\\(').Count; \"expect counts: base="'$base head=$head"; git diff --name-only 5441fe18 42d7c2d1 | Where-Object { $_ -match '"'("'^apps/server_core/|migration|''^contracts/|''^packages/sdk-runtime/)'"' }" in C:\Users\leandro.theodoro\Documents\marketplace-central
 succeeded in 8460ms:
["2026-07-18T12:00:00Z"] 2026-07-18T12:00:00.000Z
true 1970-01-01T00:00:00.001Z
1 1970-01-01T00:00:00.001Z
7	3	apps/web/src/pages/{vinculos => importacoes}/ImportacaoSection.test.tsx
expect counts: base=14 head=15

exec
"C:\\Program Files\\PowerShell\\7\\pwsh.exe" -Command '$p='"'.mnfs\\MIS-006-integracao-fundacao\\_chip-import-chain\\L2-hub-live-drive.md'; Select-String -LiteralPath "'$p -Pattern '"'links de navegação|href limpos|sem marketplace' | ForEach-Object {\""'$($_.LineNumber): $($_.Line)"}' in C:\Users\leandro.theodoro\Documents\marketplace-central
 succeeded in 21819ms:
123: código: os links de navegação das rotas **gated** carregam

exec
"C:\\Program Files\\PowerShell\\7\\pwsh.exe" -Command '$l=Get-Content -LiteralPath '"'.mnfs\\MIS-006-integracao-fundacao\\_chip-import-chain\\L2-hub-live-drive.md'; 120..130 | ForEach-Object {'{0}: {1}' -f "'$_,$l[$_-1]}' in C:\Users\leandro.theodoro\Documents\marketplace-central
 succeeded in 15481ms:
120: ## I1 — a rota e o gate de instalação
121: 
122: `/importacoes` monta e renderiza. Evidência observável da decisão de gate, além da leitura de
123: código: os links de navegação das rotas **gated** carregam
124: `?installation=inst-mercado_livre-…`, enquanto `/importacoes`, `/integracoes`, `/catalogo` e
125: `/estoque` são `href` limpos. Importação de ERP não depende de marketplace conectado, e a
126: navegação reflete isso.
127: 
128: ## I7 / F-02 — o toggle de fonte ativa persiste E a app inteira muda de fonte
129: 
130: Este é o teste que o F-02 sempre prometeu e que nenhuma suíte conseguia provar.

  ✓ Read review skill and binding authority documents
  ✓ Inspect exact diff, move integrity, and ownership boundaries
  ✓ Trace runtime composition, routes, API hook, SDK/server shapes, and error states
  ✓ Evaluate prior-review claims and criteria I1–I9
  ✓ Deliver independent gate verdict with evidence and unverifiable items
codex
VERDICT: APROVA

Nenhum defeito bloqueante em `5441fe18..42d7c2d1`. A composição real e o L2 sustentam a fiação do endpoint. Há duas ressalvas restritas a respostas `200` que violem o contrato.

| Severidade | file:line | Token/string | Cenário concreto | É defeito se |
|---|---|---|---|---|
| RESSALVA | `apps/web/src/pages/importacoes/ImportChainPanel.tsx:16` | `Number.isFinite(value)` | Payload inesperado `{ vinculados: -1 }` ou `{ vinculados: 1.5 }` → tela mostra contador impossível. O servidor atual não produz isso: os campos vêm de `count(*)`. | A proteção contra drift deve rejeitar qualquer valor fora do domínio real de contagem; nesse caso use inteiro não negativo e seguro, não apenas finito. Não afeta absent→`—` nem zero conhecido→`0`. |
| RESSALVA | `apps/web/src/pages/importacoes/ImportChainPanel.tsx:70` | `formatDateTime(chainQuery.data.queue_read_at)` | Como o SDK apenas faz cast do JSON, `queue_read_at: ["2026-07-18T12:00:00Z"]`, `true` ou `1` é coercido por `new Date(...)` e vira uma data aparentemente válida, em vez de `—`. | ADR-17/honest-unknown também governa `queue_read_at`; validar `typeof value === "string"` antes do formatter fecha a assimetria. A revisão anterior afirmou que esse campo estava “guarded correctly”, mas não testou coerção. |
| NOTA | `.mnfs/MIS-006-integracao-fundacao/_chip-import-chain/L2-hub-live-drive.md:122` | `href limpos` | O drive montou `/importacoes`, mas links sem `?installation=` não demonstram montagem para um tenant com zero instalações; o ambiente dirigido possuía instalação ML. | I1 exigir prova executada especificamente no estado sem marketplace, como diz o contrato. A leitura estática indica que funcionará: `InstallationProvider` sempre renderiza os filhos e só `InstallationGate` bloqueia. |

Confirmações principais:

- Runtime: `ClientProvider → QueryClientProvider → AppRouter → InstallationProvider`; nenhuma divergência de contexto, lazy import ou cliente paralelo.
- TanStack Query v5 permite `isError` e dados anteriores coexistirem após refetch falho, mas a ternária testa `isError` antes dos dados; erro e cadeia nunca aparecem juntos.
- 404 real é distinguido de falha genérica; todo 4xx encerra sem retry, enquanto 5xx/rede recebem uma tentativa adicional.
- Um `200` com `{}`, objeto parcial ou shape aninhado não fabrica zero: produz `—`; `null` produz o estado de erro.
- O teste movido preserva as 14 asserções originais e adiciona uma; apenas `MemoryRouter` e o link da cadeia mudaram.
- `VinculosPage.tsx` contém exatamente as duas exclusões autorizadas.
- Zero Go, migrations, `contracts/**` ou `packages/sdk-runtime/**` no diff.
- Todas as leituras próprias de `/importacoes` são ERP-side e não recebem `installationId`. A leitura global de instalações no shell não bloqueia a rota.
- O L2 observou a chamada real `/chain`, conferiu `55 / 0 / 55` contra SQL e mostrou 404 sem contadores.

O que não pude verificar:

- Não reexecutei `tsc` nem os 521 testes; o ambiente é read-only. Portanto, as contagens de L0/L1 vêm do `EVIDENCE.md`, sem reprodução independente nesta revisão.
- O L2 e suas consultas SQL foram lidos no artefato do hub, não repetidos.
- Montagem ao vivo com literalmente zero instalações marketplace; o drive usou outro observável.
- `vinculados > 0`; o drive exerceu apenas zero conhecido.
- Payload malformado em produção: o servidor e OpenAPI atuais sempre emitem os cinco campos conformes.
- Ausência real de `protocol`/`queue_read_at`; banco e servidor tornam esses casos inalcançáveis sem drift.

Critérios:

- I1 — NOT-PROVEN: código e justificativa corretos; estado vivo com zero instalações não foi efetivamente dirigido.
- I2 — PASS: move íntegro; `/vinculos` perde, `/integracoes` mantém, `/importacoes` assume.
- I3 — PASS: exatamente duas exclusões em `VinculosPage.tsx`.
- I4 — PASS: consumo real de `getErpImportChain`, confirmado por código e rede no L2.
- I5 — PASS: ausente/nulo→`—`; zero conhecido→`0`; nenhuma fabricação de zero.
- I6 — PASS: 404 específico, demais falhas genéricas, sem cadeia coexistente.
- I7 — PASS: F-02 não foi reconstruído; persistência e invalidação global dirigidas pelo hub.
- I8 — PASS: `import_id` e `getErpImportChain` já existiam; nenhuma extensão contratual necessária.
- I9 — NOT-PROVEN: escopo/ownership PASS; resultados de `tsc`/vitest não foram reproduzidos independentemente.
tokens used
183.331
VERDICT: APROVA

Nenhum defeito bloqueante em `5441fe18..42d7c2d1`. A composição real e o L2 sustentam a fiação do endpoint. Há duas ressalvas restritas a respostas `200` que violem o contrato.

| Severidade | file:line | Token/string | Cenário concreto | É defeito se |
|---|---|---|---|---|
| RESSALVA | `apps/web/src/pages/importacoes/ImportChainPanel.tsx:16` | `Number.isFinite(value)` | Payload inesperado `{ vinculados: -1 }` ou `{ vinculados: 1.5 }` → tela mostra contador impossível. O servidor atual não produz isso: os campos vêm de `count(*)`. | A proteção contra drift deve rejeitar qualquer valor fora do domínio real de contagem; nesse caso use inteiro não negativo e seguro, não apenas finito. Não afeta absent→`—` nem zero conhecido→`0`. |
| RESSALVA | `apps/web/src/pages/importacoes/ImportChainPanel.tsx:70` | `formatDateTime(chainQuery.data.queue_read_at)` | Como o SDK apenas faz cast do JSON, `queue_read_at: ["2026-07-18T12:00:00Z"]`, `true` ou `1` é coercido por `new Date(...)` e vira uma data aparentemente válida, em vez de `—`. | ADR-17/honest-unknown também governa `queue_read_at`; validar `typeof value === "string"` antes do formatter fecha a assimetria. A revisão anterior afirmou que esse campo estava “guarded correctly”, mas não testou coerção. |
| NOTA | `.mnfs/MIS-006-integracao-fundacao/_chip-import-chain/L2-hub-live-drive.md:122` | `href limpos` | O drive montou `/importacoes`, mas links sem `?installation=` não demonstram montagem para um tenant com zero instalações; o ambiente dirigido possuía instalação ML. | I1 exigir prova executada especificamente no estado sem marketplace, como diz o contrato. A leitura estática indica que funcionará: `InstallationProvider` sempre renderiza os filhos e só `InstallationGate` bloqueia. |

Confirmações principais:

- Runtime: `ClientProvider → QueryClientProvider → AppRouter → InstallationProvider`; nenhuma divergência de contexto, lazy import ou cliente paralelo.
- TanStack Query v5 permite `isError` e dados anteriores coexistirem após refetch falho, mas a ternária testa `isError` antes dos dados; erro e cadeia nunca aparecem juntos.
- 404 real é distinguido de falha genérica; todo 4xx encerra sem retry, enquanto 5xx/rede recebem uma tentativa adicional.
- Um `200` com `{}`, objeto parcial ou shape aninhado não fabrica zero: produz `—`; `null` produz o estado de erro.
- O teste movido preserva as 14 asserções originais e adiciona uma; apenas `MemoryRouter` e o link da cadeia mudaram.
- `VinculosPage.tsx` contém exatamente as duas exclusões autorizadas.
- Zero Go, migrations, `contracts/**` ou `packages/sdk-runtime/**` no diff.
- Todas as leituras próprias de `/importacoes` são ERP-side e não recebem `installationId`. A leitura global de instalações no shell não bloqueia a rota.
- O L2 observou a chamada real `/chain`, conferiu `55 / 0 / 55` contra SQL e mostrou 404 sem contadores.

O que não pude verificar:

- Não reexecutei `tsc` nem os 521 testes; o ambiente é read-only. Portanto, as contagens de L0/L1 vêm do `EVIDENCE.md`, sem reprodução independente nesta revisão.
- O L2 e suas consultas SQL foram lidos no artefato do hub, não repetidos.
- Montagem ao vivo com literalmente zero instalações marketplace; o drive usou outro observável.
- `vinculados > 0`; o drive exerceu apenas zero conhecido.
- Payload malformado em produção: o servidor e OpenAPI atuais sempre emitem os cinco campos conformes.
- Ausência real de `protocol`/`queue_read_at`; banco e servidor tornam esses casos inalcançáveis sem drift.

Critérios:

- I1 — NOT-PROVEN: código e justificativa corretos; estado vivo com zero instalações não foi efetivamente dirigido.
- I2 — PASS: move íntegro; `/vinculos` perde, `/integracoes` mantém, `/importacoes` assume.
- I3 — PASS: exatamente duas exclusões em `VinculosPage.tsx`.
- I4 — PASS: consumo real de `getErpImportChain`, confirmado por código e rede no L2.
- I5 — PASS: ausente/nulo→`—`; zero conhecido→`0`; nenhuma fabricação de zero.
- I6 — PASS: 404 específico, demais falhas genéricas, sem cadeia coexistente.
- I7 — PASS: F-02 não foi reconstruído; persistência e invalidação global dirigidas pelo hub.
- I8 — PASS: `import_id` e `getErpImportChain` já existiam; nenhuma extensão contratual necessária.
- I9 — NOT-PROVEN: escopo/ownership PASS; resultados de `tsc`/vitest não foram reproduzidos independentemente.
