# Interface Contract

```yaml
id: IC-02
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

Arquivo .xlsx do cliente ↔ módulo `erp_import` ↔ port `internal_read/ports.Reader` (consumidores: product_links, pricing, catalog, FE).

## Why This Contract Exists

Empresa da demo não estará conectada ao Sankhya; planilha é a fonte ERP. Colunas, validação e degradação honesta precisam ser idênticas entre parser, snapshot e Reader adapter.

## File Format (contrato de colunas)

Primeira aba; linha 1 = cabeçalho; matching de coluna por NOME (case-insensitive, trim, acentos ignorados).

| Coluna | Obrigatória | Tipo/validação | Ausente/inválida ⇒ |
| --- | --- | --- | --- |
| CODPROD | sim | string não vazia, única no arquivo | linha rejeitada (`EMPTY_CODPROD`/`DUPLICATE_CODPROD`) |
| DESCRPROD | sim | string não vazia | linha rejeitada (`EMPTY_DESCRPROD`) |
| CUSTO | sim | decimal > 0 (vírgula OU ponto) | linha rejeitada (`INVALID_CUSTO`) — custo NUNCA vira 0 |
| ESTOQUE_FISICO | sim | inteiro ≥ 0 | linha rejeitada (`INVALID_ESTOQUE`) |
| ESTOQUE_RESERVADO | não | inteiro ≥ 0 | unknown (disponível fica unknown — propaga, não zera) |
| EAN | não | GTIN checksum válido | warning `INVALID_EAN`, tratado ausente ⇒ matching máx. REVIEW (IC-01) |
| REFFORN | não | string | unknown |
| MARCA | não | string | unknown |
| NCM | não | 8 dígitos | warning `INVALID_NCM`, tratado ausente ⇒ GetTaxInputs unknown |

Coluna obrigatória AUSENTE no cabeçalho ⇒ import inteiro rejeitado (`MISSING_REQUIRED_COLUMN`, 422).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `POST /erp/imports` | operador (runbook/API) | multipart .xlsx | 201 `{import_id, protocol, status}` | processa síncrono (volume demo); protocolo formato M-03-like `#NNN-E` |
| `GET /erp/imports/{id}` | UI/QA | id | relatório: status, counts, linhas rejeitadas [{row, code, detail}], warnings | — |
| `GET /erp/imports` | UI/QA | — | lista newest-first | — |

## Persistence Expectations

- `erp_import_protocols`: id, tenant_id, file_hash (sha256), source `xlsx`, imported_at (UTC), status `COMPLETED`|`REJECTED`, accepted_count, rejected_count.
- `erp_import_products`: protocol_id FK, campos IC-01 + custo + estoques. Import = snapshot COMPLETO (não incremental); Reader serve o último protocolo `COMPLETED` do tenant. Runtime independe do arquivo original.

## Reader Adapter Mapping (subset consumido no MVP)

- `FindProductsForLinking`/catálogo ← identity fields; `GetSellableStock` ← estoque_fisico/reservado (reservado unknown propaga); `GetCostAsOf` ← custo, source time = imported_at (custo da planilha = custo de reposição p/ margem — assunção registrada, reversível); `GetTaxInputs` ← ncm quando presente, senão unknown/erro honesto — NUNCA zero-default (paridade com `Unavailable*` readers).
- Seleção de fonte (oracle|xlsx) por configuração; Oracle permanece caminho alternativo intacto.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| não é .xlsx / aba vazia / ilegível | 400 | `INVALID_FILE` | |
| cabeçalho sem coluna obrigatória | 422 | `MISSING_REQUIRED_COLUMN` | body nomeia coluna |
| todas as linhas rejeitadas | 201 + status `REJECTED` | — | relatório completo mesmo assim |
| import concorrente do mesmo tenant | 409 | `IMPORT_IN_PROGRESS` | |
| re-upload de arquivo idêntico (mesmo `file_hash` sha256, mesmo tenant) | 409 | `DUPLICATE_FILE` | duplicata ≠ concorrência; body aponta o import original: `{import_id, protocol}` do protocolo existente |

## Must Not Decide In Feature Execution

Nomes/obrigatoriedade de colunas, códigos de rejeição, semântica full-snapshot, mapeamento Reader.

## Validation Impact

Import de planilha exemplo real (fixture não sensível aprovada pelo operador); caso negativo: arquivo sem CUSTO ⇒ 422; linha com CUSTO=0 ⇒ rejeitada com motivo; produto sem EAN ⇒ REVIEW-only no matching.
