# F-03-reader-adapter-selection

```yaml
id: F-03
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

Adapter que implementa os Reader ports internos (custo/estoque/identidade — mesmos ports que o reader Oracle serve hoje via `internal_read`) sobre o último snapshot COMPLETED do `erp_import`, mais seleção de fonte por configuração. Demo roda com fonte xlsx; Oracle permanece selecionável e intacto.

## Inputs

- IC-02 §Reader mapping — `GetCostAsOf` ← `custo` com source time = `imported_at`; `GetSellableStock` = `estoque_fisico − estoque_reservado` quando reservado presente (incl. 0 explícito); coluna ESTOQUE_RESERVADO ausente ⇒ **disponível = DESCONHECIDO** (propaga unknown; NUNCA servir físico como número de disponível — ADR-17/IC-02). Físico permanece consultável como físico.
- F-01 (semântica identidade), F-02 (tabelas/snapshot).
- Interfaces Reader existentes em `internal_read` (assinaturas atuais — adapter conforma, não altera).
- Composition root do server (`root.go`/main wiring) — registro aditivo.

## Expected Output

- `modules/erp_import/adapter/**`: implementação dos Reader ports lendo o último protocolo COMPLETED do tenant.
- Config `MC_ERP_SOURCE` = `xlsx` | `oracle` (env, default `oracle` — comportamento atual preservado por default; demo seta `xlsx`). Seleção no composition root, aditiva.
- Entrada em `contracts/governance/modules.json` para `erp_import` (validada em worktree limpo; entra via merge do chip — nunca pré-merge na main).
- EARS: While existe protocolo COMPLETED e fonte=xlsx, when `GetCostAsOf(codprod, t)` é chamado, the sistema shall responder custo do snapshot com source time `imported_at`. While NENHUM protocolo COMPLETED existe e fonte=xlsx, when qualquer Reader port é chamado, the sistema shall retornar erro tipado `ErrNoErpSnapshot` (consumidor mostra desconhecido — nunca zero). While fonte=oracle, when Reader ports são chamados, the sistema shall usar o caminho Oracle atual sem mudança de comportamento.

## Negative Scenarios

- CODPROD inexistente no snapshot ⇒ not-found tipado (não zero).
- `t` anterior ao `imported_at` do único snapshot ⇒ `ErrNoErpSnapshot` (não há evidência para aquele momento — as-of honesto).
- Snapshot FAILED mais recente que o COMPLETED ⇒ serve o COMPLETED (FAILED nunca vira fonte).

## Constraints

- ADR-17 unknown-propagation em tudo.
- Zero mudança de assinatura nos Reader ports (senão M-07/M-08 quebram) — conformar, não redesenhar.
- Governance: entry do modules.json só via merge do chip (memória: pré-merge quebra main).

## Ownership

- Owned paths: `apps/server_core/internal/modules/erp_import/adapter/**`, composition root (linhas aditivas de registro/seleção), `contracts/governance/modules.json` (entry única), `apps/server_core/migrations/0048*–0049*` se índice/ajuste necessário.
- Forbidden paths: `modules/internal_read/**` além de consumir interfaces (mudança de assinatura = ESCALATION), Oracle adapter, demais módulos.
- Parallel-safe with: none — depends on F-01 (semântica) + F-02 (tabelas) artifacts.

## Validation Expectations

- Teste de contrato: mesmo cenário rodado contra fake Oracle e contra snapshot xlsx retorna shapes idênticos (prova de substituibilidade).
- Transcript: fonte=xlsx sem import ⇒ `ErrNoErpSnapshot` propagado até API como estado desconhecido (payload mostrando null/status, não 0).
- Transcript: pós-import, `GetCostAsOf` retorna custo exato de linha conhecida + `imported_at` como source time.
- Lane live-oracle verde no branch (Oracle path intacto).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-01, após F-01+F-02).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` com transcripts acima.
- Blockers or open decisions: none.
