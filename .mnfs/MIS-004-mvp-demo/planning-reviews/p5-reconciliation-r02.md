# P5 Reconciliation — r02

```yaml
id: P5-RECON-R02
type: planning-review
status: complete
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
lifecycle_scope: support
```

Reconciliação do audit Sol P5 r02 (`p5-sol-decomposition-audit-r02.md`, verdict **NEEDS-FOLD**:
12/19 findings r01 CLOSED; 7 NOT-CLOSED + 2 regressões novas + 1 advisory). Todos os findings
avaliados **VALID** — causa-raiz comum: folds r01 editaram o Brief mas não todas as seções
operativas (Expected Output/Ownership/EARS), ou editaram briefs sem dobrar o IC vinculante.
Nenhum exigiu autoridade do operador. Disposição:

| Finding | Reabre | Disposição | Artefatos alterados |
| --- | --- | --- | --- |
| P5-R2-01 | P5-F-05 | FOLDED — Ownership de M-01 F-01 agora POSSUI a seção OpenAPI do schema de produto catalog + tipos catalog em `sdk-runtime/src/index.ts` (additive-lock grant, ADR-12 mesmo commit); forbidden reescrito (só `erpImport.ts`/exports novos no barrel/migrations); bullet de Expected Output nomeando OpenAPI+SDK | `M-01/F-01/feature.md` §Ownership + §Expected Output |
| P5-R2-02 | P5-F-06 | FOLDED — as 3 ocorrências operativas do mapeamento invertido corrigidas: Reader (inválida ⇒ `ean: null`+warning, `refforn` SÓ de REFFORN), EARS, Negative Scenario ("nunca vai a refforn") | `M-01/F-01/feature.md` §Expected Output, EARS, §Negative Scenarios |
| P5-R2-03 | P5-F-10 | FOLDED — param exato IC-03 `?codprod=` nos dois GETs + wording do port | `M-02/F-04/feature.md` §Expected Output |
| P5-R2-04 | P5-F-11 | FOLDED — IC-03 §Operations dobrado: POST /market/collections = **200 síncrono** com sumário `{status, decisões, contagens, causas}`, sem job/polling (ruling registrado); alinhado com M-02 F-04 e M-06 | `research/market-evidence-read-interface-contract.md` §Operations |
| P5-R2-05 | P5-F-12 | FOLDED — arquitetura de mutação da missão agora carrega a exceção: envelope M-03 = único write path com alvo PROVIDER; estado LOCAL (batch de vínculos M-04) = módulo dono com preview+auditoria; Outcome, Domain Scope §Mutações, Runtime Contract e ADR-08 reescritos consistentes com a matriz | `mission.md` §Outcome, §Domain Scope, §Runtime Contract, ADR-08 |
| P5-R2-06 | P5-F-16 | FOLDED — comparação de mercado do Simulador exibe evidência IC-03 completa: `source`, `fetched_at`/freshness, `n_offers`/`n_sellers`, `match_status` | `M-07/F-02/feature.md` §Expected Output |
| P5-R2-07 | P5-F-17 | FOLDED — `tarifa_full` agora é componente explícito nullable do output de `Decompose` E da fórmula única do IC-04 (debitado só em `full`; 0 explícito nas demais; null em `full` ⇒ UNKNOWN propaga) | `M-07/F-01/feature.md` §Expected Output; IC-04 §Decomposition |
| P5-R2-08 | novo | FOLDED — rota exata IC-04: `GET /pricing/difal` + `PUT /pricing/difal/{uf}` | `M-07/F-01/feature.md` §Expected Output |
| P5-R2-09 | novo | FOLDED — Dependencies do M-09 lista os produtores reais `{M-01, M-03, M-04, M-05, M-08}` e declara M-07 não-produtor | `M-09/milestone.md` §Dependencies |
| P5-R2-10 | ADVISORY | FOLDED mesmo assim (barato) — caveat do R-04 reescrito: tributação por produto FORA do comportamento contratado; IC-04 só tem CalcProfile + override por UF; sem superfície per-SKU no MIS-004 | `research/difal-interna-rates-2026.md` §Caveats item 4 |

Varredura pós-fold: zero ocorrências de `codprods`, `product_ids`, `GET/PUT /pricing/difal`
ou `202` em `M-*/**` (grep desta sessão).

## Rerun

Manifest r03 (`p5-input-r03.sha256`) congelado após esta reconciliação; Sol P5 r03 dispatchado
(mesma cerimônia). PASS fecha o P5.
