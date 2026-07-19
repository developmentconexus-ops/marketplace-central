# HUB HANDOFF — Wave C (MIS-004-mvp-demo)

**Escrito por:** HUB v2 (`local_7c41fdd8-dd6f-487c-a7ea-1b511e214fa9`), 2026-07-19 (T-1, demo 2026-07-20).
**Para:** hub novo que assume Wave C, booted via skill `harness-hub` (harness plugin **0.3.3** — última instalada; cache tem 0.1.0→0.3.3, binding = HARNESS-CORE 0.3.3 + docs/HARNESS-PROFILE.md + mission Parallel Execution Plan).

---

## 1. Estado da missão

- **main tip no momento deste handoff:** `19b5f142` (ledger D-96). **M05-FAIXA JÁ MERGED @a4ddd106** — §3 resolvido, não é mais in-flight. Re-verifica o tip no boot.
- **Waves A+B fechadas:** M-01, M-02, M-03, M-04, M-05, M-07, M-08 = **CLOSED**. Jornada central verde ponta-a-ponta (import → vínculos → sinais mercado → simulador → pedidos).
- **Wave C (teu escopo):** M-06 ∥ M-09.
  - **M-06 produto-detalhe** = `planned`. **FE-only** (`apps/web/src/pages/produto/**` novo + `routes/produto.tsx` pós-seam; SEM backend/OpenAPI/SDK novo — consome clients existentes). Alvo real da demo (fecha narrativa: clicar produto → header + veredito + abas Anúncios vinculados/Estoque). Deps: M-01/M-02/M-04/M-05 todas CLOSED. Edge M-05→M-06 = `listings.ts` campos de sinal (já em main).
  - **M-09 dashboard-demo** = `planned`. **CORTÁVEL por design** (mission §131: "só inicia com jornada central verde"). Jornada já verde sem ele. Recomendação v2: **cortar** salvo folga real T-1. Superfícies: `modules/dashboard/**`, OpenAPI `/dashboard*`, `sdk-runtime/src/dashboard.ts`, `pages/dashboard/**`.

## 2. REGRA CRÍTICA — sole-committer no main

Doutrina: UM só hub escreve em `main` (seam compartilhado). **Não podem dois hubs mergear simultâneo.** Sequência combinada com operator:
- HUB v2 (eu) **retém só o CHIP-M05-FAIXA até o merge** (é meu, já rodando). Ao mergear, **stand down** e passo main pro hub novo.
- Hub novo forka chips de Wave C do main **no boot** (re-verifica tip — pode já ter avançado pós-M05-FAIXA) e dispara. M-06 (`pages/produto`) é disjunto de M05-FAIXA (`ListingDetailPanel` + backend market/listings) → rodam paralelo sem colisão.
- **Merge order:** v2 mergeia M05-FAIXA primeiro; hub novo mergeia M-06 depois. Zero race. Confirmar via evento antes de mergear se houver dúvida de timing.

## 3. ~~Trabalho EM VOO~~ — RESOLVIDO (nada in-flight sobra pra ti)

- **CHIP-M05-FAIXA MERGED @a4ddd106** (D-96, P7 QA PASS). Faixa min/mediana/max + own-seller exclusion dinâmica + migration **0070** (max_valid) + labels honestos — tudo em main. Worktree dazzling-mclaren-0d1dc4 stand-down; hub v2 remove após a sessão do chip encerrar. **Migration 0070 consumida** — hub novo pega **0071+** pra qualquer migration de Wave C (M-06 não tem migration; M-09 teria se não for cortado). Nenhum chip do hub v2 continua vivo: main é todo teu no boot.

## 4. Restrições binding (verbatim-strength)

- **Codex MORTO até 2026-07-25** (pós-demo) — quota estourada. Gates dual = **Claude-only**: cold Opus (contexto fresco) + adversarial sonnet (tentar REFUTAR, incl. revert empírico de guards), **AGREEMENT obrigatório**. Implementação = sonnet fallback. NÃO redispatch em codex.
- **Subagent crews = sonnet** (não fable — limite semanal).
- ZERO Mercado Livre writes. Chips nunca bootam server / bind :8080 / load .env — dev stack é seam do hub (`docker compose`); precisa live drive = REQUEST.
- Never push sem permissão explícita do operator. `git branch -d` nunca `-D`. Sem reset/revert/stash/clean/WSL/cold-clone/cache-purge. GOCACHE=path absoluto; `cd apps/server_core` pra go (gomodcache root pollution trap). Dep change = REQUEST.
- Hub commita só em `main` após branch-guard fail-closed (`b=$(git rev-parse --abbrev-ref HEAD); [ "$b" != "main" ] && exit 1`).
- FE-chip: worktree sem node_modules → junction a main + config vitest throwaway (fs.allow + setupFiles absolutos), DELETE pre-commit. FE lane = vitest+build (raw tsc baseline=2 conhecido).
- Superseded denylist (profile §10): NÃO invocar mnfs feature-execution/planning do cache codex 0.1.0 stale.
- Design: DESIGN-REFERENCE @main 8144238 = binding 1:1 QA visual wave B+ (paper+green, Instrument Sans + IBM Plex Mono).
- Só QA passa milestone (P6 dual gate + P7 browser live-drive). Todo ruling → HUB-LEDGER.md row. Evidência em .mnfs artifacts — não escrito = não aconteceu.

## 5. Infra / fatos operacionais

- Stack: `docker compose` (backend :8080, frontend :5174, postgres :5435). Mount atual = main dir (`.:/workspace`). Pós-merge FE ⇒ SEMPRE `docker compose restart frontend`.
- Installation demo: `inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2`, tenant `tenant_default`. Nosso seller_id ML = 691607102 (resolver dinâmico, nunca hardcode).
- 31 vínculos exact_ean resolvidos (D-95). Coleta full-catalog: 7 codprods com price_evidence OK (Deca/Docol), 17 NO_CANDIDATE honestos (placas Sinalize sem catálogo ML).
- Ledger corrente até **D-95** (@ffd8d2af).

## 6. Endereço do hub v2 (relay)

Se chegar evento de chip estranho, é provável que seja meu (M05-FAIXA). Relay pra `local_7c41fdd8-...` até eu mergear e sair.
