# P5 Reconciliation — r01

```yaml
id: P5-RECON-R01
type: planning-review
status: complete
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
lifecycle_scope: support
```

Reconciliação do audit Sol P5 r01 (`p5-sol-decomposition-audit-r01.md`, verdict **NEEDS-FOLD**,
19 findings BLOCKING). Todos os 19 avaliados **VALID**. Nenhum exigiu autoridade do operador:
P5-F-12 muda MECANISMO preservando o resultado aprovado no P3 (bulk aprovação local, zero ML);
P5-F-18 foi resolvido com evidência nova (R-04). Disposição finding a finding:

| Finding | Disposição | Artefatos alterados |
| --- | --- | --- |
| P5-F-01 (wave B × edge M-07→M-08) | FOLDED — parágrafo de ordem intra-wave B: trecho decomposição de M-08 F-01 só inicia após M-07 F-01 publicar ports (gate do hub); edge `M-07 F-01 → M-08 F-01 (trecho decomposição)` | `mission.md` §Parallel Execution Plan |
| P5-F-02 (edges faltantes) | FOLDED — edges `M-01→M-08` (custo GetCostAsOf) e `M-05→M-06` (listings.ts sinais) adicionados | `mission.md` DAG |
| P5-F-03 (edge agregado falso p/ M-09) | FOLDED — `{M-01..M-08}→M-09` substituído por `{M-01, M-03, M-04, M-05, M-08}→M-09` (produtores reais) | `mission.md` DAG |
| P5-F-04 (path `db/migrations/` inexistente) | FOLDED — todos os owned paths de migration corrigidos p/ `apps/server_core/migrations/`; fixture nomeada `apps/server_core/internal/platform/migrate/runner_test.go` | `mission.md` matriz; briefs M-01 F-02/F-03, M-02 F-02/F-04, M-04 F-01, M-07 F-01, M-08 F-01 |
| P5-F-05 (M-01 F-01 sem superfície OpenAPI/SDK) | FOLDED — F-01 passa a ownar seção OpenAPI do schema de produto catalog + **additive-lock grant** nos tipos catalog de `sdk-runtime/src/index.ts` (ADR-12, mesmo commit) | `mission.md` matriz M-01; `M-01/F-01/feature.md` |
| P5-F-06 (REFERENCIA inválida → refforn) | FOLDED — REFERENCIA inválida ⇒ `ean: null` + warning (tratada AUSENTE, nunca migra de campo); `refforn` SÓ de TGFPRO.REFFORN (IC-01) | `M-01/F-01/feature.md` |
| P5-F-07 (lifecycle import ≠ IC-02) | FOLDED — POST síncrono **201** `{import_id, protocol, status}`, protocolo `#NNN-E`, statuses `COMPLETED\|REJECTED`, `EMPTY_DESCRPROD` adicionado, `INVALID_EAN`/`INVALID_NCM` = warnings campo-ausente (não rejeição), EARS all-rejected⇒REJECTED | `M-01/F-02/feature.md` |
| P5-F-08 (reservado unknown → físico numérico) | FOLDED — ESTOQUE_RESERVADO ausente ⇒ disponível = DESCONHECIDO propagado (nunca físico-como-número); físico segue consultável como físico | `M-01/F-03/feature.md` |
| P5-F-09 (negativa dura rebaixada a REVIEW) | FOLDED — negativa dura (kit/combo, cor, medida, voltagem) ⇒ **REJECT mesmo com EAN igual**; contradição não-dura ⇒ teto REVIEW; fixture explícita | `M-02/F-03/feature.md` |
| P5-F-10 (chaves `product_ids` ≠ IC-03) | FOLDED — chave de produto = `codprod` (`?codprods=`, POST `{codprod}`) | `M-02/F-04/feature.md` |
| P5-F-11 (modelo de execução da coleta adiado) | FOLDED — decisão ratificada NO PLANO: POST /market/collections SÍNCRONO delimitado, **200** com sumário (sem tabela de job, sem polling); M-06 ajustado p/ await+refetch | `M-02/F-04/feature.md`; `M-06/F-01/feature.md` |
| P5-F-12 (`link_apply` insatisfazível p/ batch) | FOLDED (mudança de mecanismo, resultado aprovado preservado) — batch-preview dry-run + batch apply LOCAIS dentro de `product_links` com auditoria própria; envelope mutations fora; racional registrado (selection_resolver sem payload por item; stub `MPC_PROVIDER_WRITES_ENABLED=false` desliga writer local — `composition/root.go:568-588`) | `mission.md` matriz M-04; `M-04/milestone.md`; `M-04/F-01` + `F-02/feature.md` |
| P5-F-13 (transporte M-05→market indeciso) | FOLDED — M-02 F-04 publica port Go `market.EvidenceReader` (batch); M-05 F-01 consome via port, NUNCA HTTP self-call | `mission.md` matriz M-02; `M-02/F-04` + `M-05/F-01/feature.md` |
| P5-F-14 (menu ⚙ sem Configurações→DIFAL) | FOLDED — menu ⚙ = IC-05 exato: Configurações (item DIFAL = deep-link `/precos?params=1`), Integrações, Catálogo, Estoque, + Vínculos (entrada secundária); M-07 F-02 trata `?params=1` abrindo o drawer de parâmetros | `M-03/F-02/feature.md`; `M-07/F-02/feature.md` |
| P5-F-15 (M-03 F-03 sem Interaction Model) | FOLDED — seção Interaction Model adicionada (primitivas stateless; DataTable/DetailDrawer controlados pela página; zero fetch; tema via CSS vars) | `M-03/F-03/feature.md` |
| P5-F-16 (campos de evidência IC-03 não propagados) | FOLDED — `match_status` + `n_offers`/`n_sellers` no `market_signal` (M-05 F-01), no drawer de Anúncios (M-05 F-02) e na comparação de mercado do Simulador (M-07 F-02) | `M-05/F-01` + `F-02`; `M-07/F-02/feature.md` |
| P5-F-17 (drift IC-04 em M-07) | FOLDED — 200 `UNREACHABLE_TARGET` (era NO_SOLUTION); 404 `UF_NOT_FOUND` (era 422 INVALID_UF); 422 `INVALID_RATE` fora 0–35%; regime default `SIMPLES` 4% aplicado com origem `default` (era "nunca assume"); `tarifa_full` nullable no CalcProfile + modalidades `classico\|premium\|full` no Decompose; port `DifalForUF(uf) → {efetivo_pct, versao}`; override persiste só Δ>0,049pp; `origem_versao: "padrao-2026"` | `M-07/F-01` + `F-02/feature.md` |
| P5-F-18 (seed interna_pct sem dataset) | FOLDED com evidência nova — `research/difal-interna-rates-2026.md` (R-04) produzido por external-researcher: 27 UFs, fonte citada por linha, 26 verified + MS verify-at-execution (disputado 17%×19% — override do operador esperado); IC-04 referencia R-04; M-07 F-01 lista R-04 como input do seed | `research/difal-interna-rates-2026.md` (novo); IC-04; `M-07/F-01/feature.md` |
| P5-F-19 (Cancelados funcional reintroduzido) | FOLDED — tab Cancelados = DISABLED "em breve" (sem fetch); projeção/filtro de cancelados NÃO servido no MVP (MIS-005); exclusão dos KPIs mantida | `M-08/F-01` + `F-02/feature.md` |

## Mudanças fora dos findings (mesma rodada)

- IC-03: enum `blocking_state` explicitado (`NO_CANDIDATE|NO_PRICE_EVIDENCE|INSUFFICIENT_MARKET|SEM_CUSTO`) — feito antes do audit, mantido.
- `mission.md` risco R10 (executor `price_update` habilitado; gate = teto `previewed` + runbook) — mantido; Sol não contestou.

## Rerun

Protocolo exige rerun do audit após mudança em DAG/ownership/briefs. Manifest r02
(`p5-input-r02.sha256`) congelado após esta reconciliação; Sol P5 r02 dispatchado com a mesma
cerimônia (gpt-5.6-sol / medium / read-only / OS-process). PASS do r02 fecha o P5.
