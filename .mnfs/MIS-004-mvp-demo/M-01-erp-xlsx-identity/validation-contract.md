# Milestone Validation Contract

```yaml
id: M-01
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-01-erp-xlsx-identity.

## QA Level

Dual gate frio (Opus + Sol medium, concordância) + QA live-drive fresh (P1b). Sem UI própria neste milestone — live-drive = exercício HTTP/API + inspeção DB no stack local.

## Required Outcome

Planilha .xlsx importada com protocolo e relatório de rejeição; identidade CODPROD/EAN/REFFORN correta per IC-01; Reader port servindo custo/estoque/identidade do snapshot; Oracle intacto.

## Criteria

## Criterion: Import xlsx completo com protocolo
ID: M-01-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST /erp/imports com a planilha exemplo real (multipart) contra o stack local
- Expected: **201 síncrono** com body `{import_id, protocol, status}`; `protocol` no formato `#NNN-E`; `status` ∈ {COMPLETED, REJECTED}; re-POST do MESMO arquivo (hash igual) não duplica snapshot (409 `duplicate_file`, IC-02; hub ruling D-13 2026-07-18: código de erro wire = flat-lowercase conforme família ratificada em OpenAPI+SDK — truth order OpenAPI>prose; prose original dizia `DUPLICATE_FILE`); logs do server na janela do import SEM payload cru da planilha (grep por valores de custo/descrição do fixture ⇒ zero hits em log)
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §import (transcript request/response + SELECT import_protocols)
Blocking failure: 202/job assíncrono, protocolo fora do formato, ou import duplicado criando segundo snapshot
Blocking failure observed: No
Owner: QA Validator

## Criterion: Relatório de rejeição honesto por linha
ID: M-01-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST /erp/imports com fixture contendo: linha DESCRPROD vazio, linha REFERENCIA checksum inválido, linha NCM malformado, linha válida
- Expected: linha DESCRPROD vazio ⇒ REJEITADA com motivo `EMPTY_DESCRPROD` inspecionável; REFERENCIA inválida ⇒ linha IMPORTADA com warning `INVALID_EAN` e `ean: null` (valor tratado como ausente); NCM malformado ⇒ warning `INVALID_NCM` + `ncm: null`; arquivo 100% rejeitado ⇒ `status: REJECTED` com relatório completo
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §rejeicao (response + relatório)
Blocking failure: warning virando rejeição (ou vice-versa), motivo genérico sem código, ou valor inválido persistido como se válido
Blocking failure observed: No
Owner: QA Validator

## Criterion: Identidade IC-01 no reader e na API
ID: M-01-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: rodar testes de contrato do reader (fixtures IC-01: REFERENCIA GTIN válida, inválida, vazia; REFFORN presente/ausente; colisão EAN entre CODPRODs) + GET /catalog de produto importado
- Expected: REFERENCIA válida-GTIN → `ean` populado; inválida/vazia → `ean: null` + warning, `refforn` INALTERADO (exclusivamente TGFPRO.REFFORN); GET responde `ean|refforn|marca|ncm: string|null` — null quando desconhecido, nunca string vazia; colisão EAN ⇒ ambos legíveis + colisão sinalizada
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §identidade (transcript testes + response GET)
Blocking failure: REFERENCIA inválida migrando p/ refforn, ou string vazia fabricada no lugar de null
Blocking failure observed: No
Owner: QA Validator

## Criterion: Reader port servindo do snapshot
ID: M-01-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: após import COMPLETED, exercitar `GetCostAsOf` e `GetSellableStock` (teste de integração hermético) p/ item com ESTOQUE_RESERVADO presente e outro com ausente
- Expected: custo = valor da planilha com source time do import; reservado presente ⇒ disponível = físico − reservado; reservado AUSENTE ⇒ disponível = DESCONHECIDO propagado (físico segue consultável como físico) — nunca físico-como-disponível
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §reader (transcript)
Blocking failure: disponível fabricado quando reservado ausente, ou custo sem source time
Blocking failure observed: No
Owner: QA Validator

## Criterion: Oracle path intacto
ID: M-01-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: lane live-oracle (testes de contrato cobrindo ambas as fontes) + build completo, `GOCACHE=.gocache`
- Expected: exit 0; testes do reader Oracle passam com a NOVA semântica REFERENCIA (ean null p/ inválida) — regressão coberta
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §oracle (transcript lane)
Blocking failure: lane vermelha ou semântica antiga (REFERENCIA→refforn) sobrevivendo no path Oracle
Blocking failure observed: No
Owner: QA Validator

## Criterion: Migrações e seams
ID: M-01-C06
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `ls apps/server_core/migrations/ | grep -E '^004[5-9]'` + rodar `runner_test.go`; diff do chip vs matriz de ownership
- Expected: migrações novas SOMENTE no bloco 0045–0049; fixture count = contagem real; diff toca só superfícies exclusivas M-01 + grant aditivo (`sdk-runtime/src/index.ts` tipos catalog, mesmo commit do OpenAPI, ADR-12); OpenAPI seções `/erp/imports*` + schema catalog atualizadas no mesmo commit do SDK
- Actual:
- Artifact: `M-01-erp-xlsx-identity/validation-result.md` §seams (diff summary)
Blocking failure: migração fora do bloco, fixture divergente, ou write fora do ownership/grant
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-01-erp-xlsx-identity/validation-result.md` com seções import, rejeicao, identidade, reader, oracle, seams.
- Fixture xlsx exemplo + fixture de linhas inválidas versionadas junto à evidência.
- Dual gate: verdicts Opus + Sol registrados antes do QA live-drive.

## Blocking Failures

Qualquer blocking acima ⇒ FAIL; correção pelo chip (máx. 2 tentativas), depois escalação hub.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close das features do chip M-01).
- Next action: aguardar F-01/F-02/F-03 aceitas; rodar critérios.
- Required files/evidence: `M-01-erp-xlsx-identity/validation-result.md`.
- Blockers or open decisions: none.
