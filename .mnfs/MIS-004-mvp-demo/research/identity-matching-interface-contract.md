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
