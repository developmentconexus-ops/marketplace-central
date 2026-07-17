# F-01-identity-semantics-fix

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-01-erp-xlsx-identity.

## Brief

Corrigir a semântica de identidade em toda a leitura de produto: `REFERENCIA` do Sankhya NÃO é EAN por definição (defeito estrutural, research R-01 §4). Regra canônica IC-01: `ean` = TGFPRO.REFERENCIA **somente quando** passa checksum GTIN-8/12/13/14; REFERENCIA inválida ⇒ `ean: null` + warning registrado (tratada como AUSENTE — nunca migra p/ outro campo). `refforn` vem EXCLUSIVAMENTE de TGFPRO.REFFORN. Implementar validador GTIN único (função pura em `modules/catalog`, dono da identidade), aplicar no reader Oracle (`internal_read`) e expor os campos canônicos (`codprod`, `ean`, `refforn`, `marca`, `ncm`) na API de produto do catalog — incluindo a seção OpenAPI do schema de produto e os tipos catalog do SDK (ADR-12: OpenAPI+SDK no mesmo commit). Nenhum consumidor fora do adapter vê REFERENCIA cru.

## Inputs

- IC-01 (`research/identity-matching-interface-contract.md`) — regra GTIN, campos canônicos, colisões.
- Código atual: `modules/internal_read/**` (reader Oracle), `modules/catalog/**` (API de produto).
- Fixtures de colisão do IC-01 §Error matrix (EAN duplicado entre CODPRODs, REFERENCIA vazia, checksum inválido).

## Expected Output

- Validador GTIN exportado por `catalog` (aceita 8/12/13/14 dígitos, checksum mod-10; input não-numérico ou tamanho errado ⇒ false).
- Reader Oracle mapeia: REFERENCIA válida-GTIN → `ean`; inválida → `ean: null` + warning registrado (valor tratado como AUSENTE — nunca migra p/ `refforn`); `refforn` populado EXCLUSIVAMENTE de TGFPRO.REFFORN.
- Seção OpenAPI do schema de produto do catalog atualizada + tipos catalog do SDK (`packages/sdk-runtime/src/index.ts`, aditivo) no MESMO commit (ADR-12).
- GET de produto (catalog) responde `ean: string|null`, `refforn: string|null`, `marca: string|null`, `ncm: string|null` — null quando desconhecido, nunca string vazia fabricada.
- EARS: While produto tem REFERENCIA que passa checksum GTIN, when qualquer leitura de identidade ocorre, the sistema shall expor o valor em `ean` e manter `refforn` com o valor de REFFORN se existir. While REFERENCIA falha checksum, when leitura ocorre, the sistema shall expor `ean: null` + warning e manter `refforn` apenas com o valor de REFFORN (REFERENCIA inválida descartada da identidade). While dois CODPRODs distintos resolvem o mesmo `ean`, when leitura ocorre, the sistema shall manter ambos e marcar colisão (flag por produto), nunca fundir registros.

## Negative Scenarios

- REFERENCIA vazia/whitespace ⇒ `ean: null` (não string vazia); `refforn` NÃO é afetado — vem exclusivamente de TGFPRO.REFFORN (null só quando REFFORN ausente).
- Checksum inválido com 13 dígitos ⇒ `ean: null` + warning (nunca "quase EAN", nunca vai a `refforn`).
- Colisão de EAN entre CODPRODs ⇒ ambos legíveis + colisão sinalizada; matching downstream trata como contradição (IC-01).

## Constraints

- ADR-17: desconhecido nunca vira zero/vazio.
- Oracle path continua funcional (lane live-oracle passa) — mudança é de MAPEAMENTO, não de conexão.
- Não tocar matching/resolver (M-02 F-03) nem import xlsx (F-02).

## Ownership

- Owned paths: `apps/server_core/internal/modules/internal_read/**`, `apps/server_core/internal/modules/catalog/**` (campos identity + validador), testes desses módulos, seção do schema de produto do catalog no OpenAPI (aditiva), tipos catalog em `packages/sdk-runtime/src/index.ts` (**additive-lock grant** registrado na matriz da missão — aditivo apenas, ADR-12 mesmo commit).
- Forbidden paths: `modules/erp_import/**` (F-02), OpenAPI `/erp/imports*` (F-02), `sdk-runtime/src/erpImport.ts` (F-02), exports novos no barrel além dos tipos catalog do grant (barrel = hub), migrations (F-02), demais módulos.
- Parallel-safe with: F-02 (disjoint: módulos/arquivos distintos; F-02 consome o validador via import — assinatura publicada no IC-01, não negociada em execução).

## Validation Expectations

- Teste unitário do validador: tabela com ≥8 casos (GTIN-8/12/13/14 válidos, checksum errado, 12 dígitos não-UPC, vazio, alfanumérico) — todos com resultado exato esperado.
- Teste de contrato do reader: fixture Oracle com REFERENCIA válida e inválida ⇒ JSON de produto com `ean`/`refforn` exatos por caso, incluindo `ean: null` (não ausente, não "").
- Caso de colisão: dois produtos, mesmo EAN ⇒ ambos retornados com flag de colisão visível no payload.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-01).
- Next action: criar `spec.md` refinando este brief contra IC-01.
- Required files/evidence: `validation.md` da execução com transcripts dos testes acima.
- Blockers or open decisions: none.
