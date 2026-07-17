# Milestone Validation Contract

```yaml
id: M-02
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

M-02-price-intel-core.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive. Live-drive = API/DB + **seam de integração REAL**: leitura live contra a installation ML conectada (ver §Blocking Failures — binding declarado).

## Required Outcome

Evidência de mercado persistida no `market`; extensões read-only do adapter ML atrás de ports IC-06; resolver determinístico IC-01; API de coleta+veredicto IC-03 publicada (HTTP + port Go `EvidenceReader`).

## Criteria

## Criterion: Coleta síncrona com sumário
ID: M-02-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST /market/collections p/ conjunto de codprods vinculados (seam live ML — installation real)
- Expected: **200 síncrono** com sumário `{status: COMPLETED|PARTIAL, decisões, contagens, causas}`; NENHUMA tabela de job, nenhum endpoint de polling; falha parcial de rota ⇒ `PARTIAL` com causas nomeadas no body
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §coleta (transcript)
Blocking failure: 202/job/polling, ou falha parcial silenciosa (200 COMPLETED com coleta faltando)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Leituras de evidência IC-03
ID: M-02-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: GET /market/aggregates?codprod= e GET /market/verdicts?codprod= pós-coleta
- Expected: parâmetro exato `codprod`; todo item carrega source, fetched_at, n_offers, n_sellers, match_status; veredicto composto por três enums distintos: `match_status` ∈ {ACCEPT, REVIEW, REJECT, NO_CANDIDATE} (IC-01); `price_evidence_status` ∈ {OK, NO_PRICE_EVIDENCE, INSUFFICIENT_MARKET} (IC-01); `blocking_state` ∈ {NO_CANDIDATE, NO_PRICE_EVIDENCE, INSUFFICIENT_MARKET, SEM_CUSTO} (IC-03); INSUFFICIENT_MARKET disparado com <5 sellers válidos
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §leituras (response bodies)
Blocking failure: campo de evidência ausente, ou estado negativo colapsado em zero/ausência
Blocking failure observed: No
Owner: QA Validator

## Criterion: Port Go EvidenceReader
ID: M-02-C03
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: teste de contrato do port `market.EvidenceReader` (batch: sinais por listing_ids; agregados/veredictos por codprod)
- Expected: shapes idênticos ao contrato HTTP IC-03 (mesmos campos de evidência); consumo interno SEM HTTP self-call (transporte fixado); assinatura publicada e estável p/ M-05
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §port (transcript teste)
Blocking failure: shape divergente do HTTP, ou consumidor forçado a self-call HTTP
Blocking failure observed: No
Owner: QA Validator

## Criterion: Resolver determinístico com hard negatives
ID: M-02-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: rodar fixtures de matching IC-01: gate 2 âncoras, colisão EAN (Doka/Menegotti, Doka/VW), hard negatives (kit/combo, cor, medida, voltagem) com EAN igual, contradição não-dura
- Expected: match exige ≥2 âncoras concordantes; hard negative ⇒ **REJECT mesmo com EAN igual**; contradição não-dura ⇒ teto REVIEW (nunca ACCEPT automático); colisão EAN ⇒ contradição vence EAN; título sozinho NUNCA produz ACCEPT
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §resolver (transcript fixtures, decisão por caso)
Blocking failure: hard negative aceito, ou ACCEPT por âncora única/título
Blocking failure observed: No
Owner: QA Validator

## Criterion: ADR-17 — falha de coleta não zera snapshot
ID: M-02-C05
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: teste negativo: coleta com provider indisponível (adapter forçado a erro) sobre codprod que JÁ tem snapshot válido
- Expected: snapshot anterior INTACTO (SELECT antes/depois idêntico); leitura reporta staleness via fetched_at antigo; response da coleta = PARTIAL/causa, nunca snapshot zerado
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §adr17 (SELECTs + response)
Blocking failure: snapshot válido sobrescrito por zero/vazio após falha
Blocking failure observed: No
Owner: QA Validator

## Criterion: Adapter ML — flag, paginação, telemetria, fallback
ID: M-02-C06
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: lane live-provider-read contra installation real: rotas sale_price, price_to_win, products search/detail, shipments, shipping_options; rota products/{id}/items com flag default OFF e depois ON
- Expected: flag OFF ⇒ rota não chamada + fallback explícito registrado; flag ON ⇒ paginação completa (>1 página quando aplicável) + telemetria da rota visível no log (rota, status, latência); DTO ML NÃO vaza dos adapters (consumidores recebem shapes IC-06); somente métodos GET no log
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §adapter (log excerpt + transcript)
Blocking failure: flag default ON, DTO provider vazando de connectors, ou método de escrita no log
Blocking failure observed: No
Owner: QA Validator

## Criterion: Migrações e seams
ID: M-02-C07
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `ls apps/server_core/migrations/ | grep -E '^005[0-4]'` + `runner_test.go`; diff vs ownership
- Expected: tabelas `market_price_snapshots`, `market_validated_offers`, `market_aggregates`, `market_competitive_signals`, `market_match_decisions` no bloco 0050–0054; fixture = contagem real; diff só em superfícies M-02 (`modules/market/**`, `modules/connectors/**`, `sdk-runtime/src/market.ts`, OpenAPI `/market/*`)
- Actual:
- Artifact: `M-02-price-intel-core/validation-result.md` §seams
Blocking failure: migração fora do bloco ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-02-price-intel-core/validation-result.md` com seções coleta, leituras, port, resolver, adr17, adapter, seams.
- Log do adapter cobrindo C01+C06 (prova read-only).
- Dual gate: verdicts Opus + Sol antes do live-drive.

## Blocking Failures

Seam de integração DECLARADO: C01/C02/C06 são **live-driven** contra a installation ML real do dev stack (credenciais AES via `MPC_ENCRYPTION_KEY`; sem VPN). Mock/stub NÃO satisfaz estes critérios sem autorização explícita do operador registrada aqui. Fixtures do resolver (C04) e teste negativo (C05) são herméticos por design (mocks provam contrato, não integração). Installation ausente ⇒ trigger R6 (runbook), criterion Blocked — não Pass.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01…F-04 do chip M-02).
- Next action: preflight R4 (leitura live cedo na wave A) antes dos critérios completos.
- Required files/evidence: `M-02-price-intel-core/validation-result.md`.
- Blockers or open decisions: none.
