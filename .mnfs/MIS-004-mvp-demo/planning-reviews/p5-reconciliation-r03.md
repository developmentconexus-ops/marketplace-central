# P5 Reconciliation — r03

```yaml
id: P5-RECON-R03
type: planning-review
status: complete
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
lifecycle_scope: support
```

Reconciliação do audit Sol P5 r03 (`p5-sol-decomposition-audit-r03.md`, manifest r03
`8448f39e…`, verdict **NEEDS-FOLD**: 7/10 folds r02 CLOSED; 3 residuais NOT-CLOSED + 1 finding
novo BLOCKING). Todos avaliados **VALID**; nenhum exigiu autoridade do operador. Causa-raiz dos
residuais: fold r02 corrigiu o locus apontado mas não varreu ocorrências irmãs do mesmo valor
em artefatos vizinhos (milestone vs feature; brief consumidor vs produtor; ADR vs ADR).
Correção de método aplicada: cada fold agora termina com grep do token rejeitado sobre
`mission.md` + `M-*/**` + `research/*` (não só `M-*/**`).

| Finding | Reabre | Disposição | Artefatos alterados |
| --- | --- | --- | --- |
| P5-R2-01 residual | P5-F-05 | FOLDED — Ownership do MILESTONE M-01 agora nomeia o grant: exclusive surfaces incluem seção OpenAPI do schema catalog (F-01, aditivo); bullet dedicado "Additive-lock grant (F-01): tipos catalog em `sdk-runtime/src/index.ts`, aditivo-only, mesmo commit, liberado no close"; seam lock do barrel qualificado "fora do grant acima" | `M-01/milestone.md` §Ownership & Concurrency |
| P5-R2-02 residual | P5-F-06 | FOLDED — Negative Scenario reescrito: REFERENCIA vazia/whitespace ⇒ `ean: null`; `refforn` NÃO afetado (exclusivamente TGFPRO.REFFORN, null só quando REFFORN ausente) — REFERENCIA em branco não anula REFFORN válido | `M-01/F-01/feature.md` §Negative Scenarios |
| P5-R2-03 residual | P5-F-10 | FOLDED — brief consumidor M-05 F-01 §Inputs: "agregados/veredictos por `codprod`" (token `codprods` eliminado); sweep confirmou zero ocorrências fora de artefatos de review históricos | `M-05/F-01/feature.md` §Inputs |
| P5-R2-05 residual | P5-F-12 | FOLDED — ADR-01 coluna Must-preserve qualificada: "único write path com alvo PROVIDER (estado LOCAL = módulo dono, ADR-08)"; sweep: nenhuma ocorrência não-qualificada restante | `mission.md` ADR-01 |
| P5-R3-01 | novo | FOLDED — ⚙ menu = SOMENTE os 4 itens IC-05 (Configurações→DIFAL, Integrações, Catálogo, Estoque); `Vínculos` FORA da navegação global (nem pill nem ⚙); entrada primária = tela Anúncios; rota /vinculos registrada e alcançável por deep-link (satisfaz o yes-if de Sol integralmente) | `M-03/F-02/feature.md` §Expected Output |

Nota: ocorrências de `codprods`/`único write path` não-qualificado dentro de
`planning-reviews/*` r01–r03 são citações históricas em artefatos de review imutáveis — não são
write-set executável e não se dobram.

## Rerun

Manifest r04 (`p5-input-r04.sha256`) congelado após esta reconciliação; Sol P5 r04 dispatchado
(mesma cerimônia). PASS fecha o P5.
