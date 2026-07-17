# F-02-erp-import-module

```yaml
id: F-02
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

Módulo novo `erp_import`: upload de planilha .xlsx do Sankhya, validação por linha contra o contrato de colunas IC-02, persistência full-snapshot com protocolo e relatório de rejeição inspecionável. API `/erp/imports` + client SDK. O import NÃO decide identidade além do que IC-02/IC-01 fixam (checksum via validador do catalog, F-01).

## Inputs

- IC-02 (`research/erp-xlsx-import-interface-contract.md`) — contrato de colunas, códigos de rejeição, endpoints, tabelas, semântica full-snapshot.
- IC-01 — validador GTIN (assinatura publicada; import consome).
- Migrations: bloco 0045–0049 em `apps/server_core/migrations/`; fixture de contagem `apps/server_core/internal/platform/migrate/runner_test.go` (bump).
- Padrão de módulo existente (qualquer módulo em `apps/server_core/internal/modules/` como referência estrutural, ex.: `orders` p/ import de arquivo).

## Expected Output

- Tabelas `erp_import_protocols` (id, tenant_id, file_hash sha256, filename, protocolo `#NNN-E`, status `COMPLETED|REJECTED`, imported_at, counts) e `erp_import_products` (protocol_id, codprod, descrprod, custo, estoque_fisico, estoque_reservado, ean, refforn, marca, ncm) — migrations no bloco 0045+ (`apps/server_core/migrations/`).
- `POST /erp/imports` (multipart xlsx) ⇒ processamento SÍNCRONO, **201** com `{import_id, protocol, status}` (protocolo formato `#NNN-E`); `GET /erp/imports` lista; `GET /erp/imports/{id}` detalhe com relatório de rejeições (linha, coluna, código, valor ofensivo) e warnings.
- Códigos EXATOS do IC-02: rejeição de linha — `EMPTY_CODPROD`, `DUPLICATE_CODPROD`, `EMPTY_DESCRPROD`, `INVALID_CUSTO`, `INVALID_ESTOQUE`; **warnings (linha importa, campo tratado como ausente)** — `INVALID_EAN`, `INVALID_NCM`; arquivo — 422 `MISSING_REQUIRED_COLUMN`, 400 `INVALID_FILE`, 409 `IMPORT_IN_PROGRESS`, duplicata exata por `file_hash` ⇒ 409 `DUPLICATE_FILE` com body `{import_id, protocol}` do protocolo existente.
- Seção OpenAPI `/erp/imports*` + `packages/sdk-runtime/src/erpImport.ts` (client tipado novo; export no barrel = seam do hub).
- EARS: While nenhum import roda para o tenant, when POST com xlsx válido chega, the sistema shall processar sincronamente, criar protocolo, validar linhas, persistir snapshot completo e responder 201 status COMPLETED com contagens. While arquivo tem linha inválida, when validação roda, the sistema shall rejeitar SÓ a linha com código+motivo e importar as demais (COMPLETED com rejected_count>0). While TODAS as linhas são rejeitadas, when validação termina, the sistema shall responder 201 com status REJECTED (protocolo persiste com relatório; nenhum produto vira snapshot). While coluna obrigatória ausente, when validação inicia, the sistema shall falhar o arquivo inteiro com 422 MISSING_REQUIRED_COLUMN sem persistir produto algum.

## Inputs/Outputs

Shapes, campos e códigos: IC-02 §Column Contract + §Endpoints (não restatar aqui — IC é a fonte). Coluna EAN passa pelo validador GTIN de F-01: inválido ⇒ **warning `INVALID_EAN` + campo tratado como ausente (`ean: null`)** — NUNCA rejeição de linha; REFERENCIA-derivação é papel do reader (F-01), não do import.

## Negative Scenarios

- .xlsx corrompido/planilha vazia/extensão errada ⇒ 400 `INVALID_FILE`, nenhum protocolo COMPLETED.
- Import concorrente mesmo tenant ⇒ 409 `IMPORT_IN_PROGRESS`.
- CUSTO ≤ 0 ou não-numérico ⇒ linha rejeitada `INVALID_CUSTO` (custo desconhecido NÃO vira 0 — ADR-17).
- ESTOQUE_FISICO negativo ⇒ `INVALID_ESTOQUE`.
- Colunas opcionais ausentes ⇒ campos null (unknown-propagation), sem rejeição.
- Re-upload de arquivo idêntico (mesmo sha256) ⇒ 409 `DUPLICATE_FILE` com `{import_id, protocol}` do protocolo original (IC-02).

## Constraints

- Full-snapshot: cada COMPLETED substitui o anterior como "latest"; sem merge incremental no MVP.
- `tenant_id` em toda query. Payload de planilha nunca logado cru (PII/custo). `POST /erp/imports` reusa o guard HTTP existente do server_core — NENHUM middleware/superfície de auth novo (baseline da missão).
- Sem parser xlsx novo por capricho: dependência nova = REQUEST ao hub (doutrina). Verificar lib existente no repo primeiro; se nenhuma, REQUEST.
- Não tocar `internal_read`/`catalog` (F-01) nem adapter/seleção (F-03).

## Ownership

- Owned paths: `apps/server_core/internal/modules/erp_import/**` (novo, exceto `adapter/` de F-03), `apps/server_core/migrations/0045*–0047*`, fixture `apps/server_core/internal/platform/migrate/runner_test.go` (bump), seção `/erp/imports*` do OpenAPI, `packages/sdk-runtime/src/erpImport.ts`.
- Forbidden paths: `modules/internal_read/**`, `modules/catalog/**` (F-01), composition root (F-03/hub), `sdk-runtime/src/index.ts` (barrel = hub), demais módulos/migrations.
- Parallel-safe with: F-01 (disjoint: módulos, OpenAPI section, migrations exclusivos deste F).

## Validation Expectations

- Transcript: POST planilha exemplo (≥50 linhas, ≥3 inválidas, ≥1 EAN checksum-inválido) ⇒ 201 síncrono com `{import_id, protocol, status}`; GET detalhe ⇒ status COMPLETED, `imported_count`/`rejected_count` exatos, cada rejeição com {linha, código, valor}, warning INVALID_EAN presente com a linha IMPORTADA e `ean: null`.
- Transcript 422 MISSING_REQUIRED_COLUMN (planilha sem CODPROD) e 400 INVALID_FILE (arquivo .txt renomeado).
- Query nas tabelas mostrando snapshot persistido com `custo` decimal exato de uma linha conhecida e `ean: null` para linha sem EAN.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-01).
- Next action: criar `spec.md`; verificar lib xlsx existente ANTES de qualquer REQUEST de dependência.
- Required files/evidence: `validation.md` com transcripts acima.
- Blockers or open decisions: none.
