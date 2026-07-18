# Interface Contract

```yaml
id: IC-01
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

catalog ↔ erp_import ↔ product_links ↔ market ↔ FE (todas as telas que exibem identidade/veredicto).

## Why This Contract Exists

Dois workers independentes não podem divergir sobre o que é o SKU canônico, o que é EAN, como um match é aceito e quais estados existem. Defeito estrutural corrente: reader Oracle trata REFERENCIA como referência de fabricante — este contrato fixa a semântica correta (research §4).

## Resources Or Entities

**ProductIdentity** (dono: catalog):
- `codprod` (string, PK canônico, SKU interno — TGFPRO.CODPROD)
- `descrprod` (string)
- `ean` (string nullable — TGFPRO.REFERENCIA; GTIN-8/12/13/14 com checksum válido; inválido ⇒ tratado como ausente + warning)
- `refforn` (string nullable — referência do fabricante, TGFPRO.REFFORN)
- `marca` (string nullable), `ncm` (string nullable, 8 dígitos)

Regra dura: `seller_sku` (ML) resolve SOMENTE para `codprod`. Nunca para `ean` ou `refforn`.

## Matching Gate (determinístico — research §5)

- Auto-ACCEPT exige **2 âncoras independentes** concordando (ex.: EAN+marca, EAN+refforn, refforn+marca). Título fuzzy RANQUEIA candidatos, NUNCA é âncora de auto-ACCEPT.
- **Contradição vence EAN**: qualquer hard negative (kit/combo, cor, medida/dimensão, voltagem divergente) ⇒ REJECT do candidato mesmo com EAN igual (colisões provadas: Doka/Menegotti, Doka/VW).
- EAN ausente ⇒ resultado máximo = REVIEW (nunca auto-ACCEPT).
- Confiança (chip UI): ≥85 verde · 50–84 âmbar · <50 vermelho. Enum de banda (API): `confidence_band: ALTA|MEDIA|BAIXA` (mesmos limiares; consumido por M-04 F-01).

## Enums And Statuses

- `match_status`: `ACCEPT` | `REVIEW` | `REJECT` | `NO_CANDIDATE`
- `price_evidence_status`: `OK` | `NO_PRICE_EVIDENCE` | `INSUFFICIENT_MARKET` (n_sellers válidos < 5)
- `verdict`: emitido SÓ com match ACCEPT + price_evidence OK + custo/fees/frete/imposto conhecidos; senão exibe o estado bloqueante. Nunca R$0, nunca margem zero por default (ADR-17).

## Error Cases

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| codprod inexistente em leitura | 404 | `PRODUCT_NOT_FOUND` | nunca inventar produto |
| match forçado sem âncoras (API de resolução manual é exceção humana) | 422 | `INSUFFICIENT_ANCHORS` | resolução manual grava `resolved_by: operator` |

## Must Preserve

Vocabulário de estados EXATO acima em API, SDK e UI (copy PT-BR pode traduzir rótulo, valor de enum não muda). Evidência do match (âncoras usadas, contradições) persistida com a decisão.

## Must Not Decide In Feature Execution

Lista de âncoras, thresholds de confiança, novos estados, semântica de colunas Sankhya.

## Validation Impact

Fixtures de colisão (Doka/Menegotti, Doka/VW) obrigatórias no resolver (M-02); QA verifica REVIEW-only quando sem EAN (M-01/M-04).

## Amendment A2 — feasible anchor model for provider-side data (2026-07-18, hub ruling D-25)

Field finding (CHIP-M04 D01/D02): listing snapshots carry ONLY `seller_sku + ean + title`
(migração 0022); `marca`/`refforn` existem apenas no lado interno — pares de âncoras
marca/refforn são estruturalmente incomputáveis contra dados de provider hoje. Modelo viável
ratificado (owner ruling; o texto original permanece válido para quando houver enriquecimento
de provider):

- **Âncoras cross-side disponíveis**: `seller_sku` (exact → codprod; regra dura inalterada) e
  `ean` (corroboração; derivado de REFERENCIA, flag `unproved` mantida). `title` segue
  ranking-only. `marca`/`refforn` DEVEM aparecer em `reasons[]` como `UNAVAILABLE` (ADR-17 —
  motivo sempre visível), nunca silenciosamente omitidas.
- **`seller_sku` é âncora ACCEPT-grade** no escopo M-04 (vínculo listing-própria ↔ ERP): é
  mapeamento autorado pelo próprio tenant no ML e resolve somente para codprod. NÃO se aplica
  ao matching de ofertas de mercado (M-02) — ofertas de concorrentes não têm seller_sku
  significativo; para M-02 nada muda neste amendment.
- **Par de auto-ACCEPT viável** = `seller_sku` + `ean` concordando no MESMO codprod, sem hard
  negative ⇒ banda ALTA, `match_status: ACCEPT` (proxy das "2 âncoras independentes").
- **Bandas**: âncora única (inclusive seller_sku sozinho) ⇒ MEDIA / `REVIEW` — a regra "EAN
  ausente ⇒ máximo REVIEW" permanece binding. Title-only ⇒ BAIXA. Conflito SKU/EAN ou hard
  negative de título (kit/combo, cor, medida, voltagem) ⇒ cap BAIXA + reason AGAINST
  (contradição vence EAN, inalterado). Thresholds numéricos ≥85/50–84/<50 inalterados.
- **Fixtures**: M-04 F-01 autora fixtures próprias derivadas deste modelo (≥8 casos, banda +
  reasons exatos, incluindo hard-negative estilo Doka); fixtures de colisão do M-02 seguem
  obrigatórias e inalteradas. Verdade única: qualquer futuro consumo de âncoras no M-02+ lê
  ESTE amendment, não decide localmente.
