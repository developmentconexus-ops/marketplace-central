# F-02-integracoes-health-section

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-09
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-09 sync-observability.

## Brief

Seção "Saúde do sync" na IntegracoesPage: card novo `SyncHealthCard` (arquivo próprio em
`apps/web/src/pages/integracoes/`), montado em `IntegracoesPage.tsx` (ordem atual de seções
`:558-574`: ActiveSourceCard → SellableAssortmentCard → UploadCard → ProviderConnectCard →
ImportacaoSection; SyncHealthCard entra após ProviderConnectCard — junto do contexto ML).
Consome `getSyncHealth` (F-01) via padrão da página (hooks web-query ou `useClient()`
direto como `listIntegrationInstallations` `:508` — spec segue o idioma da página). Por
entidade: nome, última sincronização (relativa: "há 5min"/"nunca"), badge estado
(verde=sucesso recente, vermelho=falhas consecutivas com last_error em tooltip,
cinza=nunca), phase quando presente. Bloco webhook: payload SEMPRE traz o shape IC-05;
discriminador de estado inicial = `last_notification_at === null` → subseção mostra o
FATO: "nenhuma notificação recebida" (NUNCA veredito de configuração — payload `{null,0,0}`
é idêntico p/ não-configurado e configurado-silencioso; IC-05 pina, F-r07-3); com
timestamp presente → "última notificação / pendentes / descartadas 24h". Tokens DESIGN-REFERENCE @8144238.

EARS:
- While health responde entities, when card renderiza, cada entity shall mostrar badge +
  timestamp relativo coerente com o payload.
- While entity nunca rodou (nulls), when card renderiza, the linha shall mostrar `nunca` em
  cinza — nunca "0 min atrás", nunca verde.
- While webhook block está no estado canônico inicial (last_notification_at null), when
  card renderiza, the subseção shall dizer "nenhuma notificação recebida" — o fato, nunca
  o veredito "não configurado" (F-r07-3); não esconder, não inventar — pending/dropped 0
  do estado inicial NÃO renderizam como atividade.

## Inputs

F-01 (SDK method + payload IC-05); fato #4 de `research/p5-prerequisites.md`
(IntegracoesPage estrutura/idioma de data
fetching verbatim); DESIGN-REFERENCE @8144238.

## Expected Output

SyncHealthCard + mount de 1 linha + (se idioma pedir) hook novo em web-query — par com o
SDK do F-01 já publicado.

## Constraints

- IntegracoesPage.tsx: SÓ a linha de mount (card é arquivo próprio — colisão zero com
  seções existentes).
- Refetch: polling leve (30s) OU manual — spec decide contra idioma web-query da página;
  sem websocket.
- Timestamps relativos com título absoluto (hover) — auditabilidade.
- tsc verde; teste de página existente (`IntegracoesPage.test.tsx`) NÃO quebra.

## Inputs/Outputs

- Input: payload `getSyncHealth` — shape BINDING em IC-05 §Required Outputs (`entities[]`:
  entity, last_success_at GREATEST, last_incremental_at, consecutive_failures, phase,
  last_error; bloco `webhook`: last_notification_at, pending, dropped_24h — F-r08-4) — a
  spec não re-decide shape.
- Output (render states): por entidade — verde (sucesso recente) / vermelho (falhas
  consecutivas + last_error em tooltip) / cinza "nunca" (nulls — nunca "0 min atrás");
  bloco webhook — estado inicial `{null,0,0}` → FATO "nenhuma notificação recebida"
  (nunca veredito de configuração — F-r07-3); com timestamp → última notificação /
  pendentes / descartadas 24h. (Seção adicionada — auditoria P5 r07 F-r07-6.)

## Interaction Model

Card read-only (sem ações nesta missão). Estado de erro do fetch: card mostra erro nomeado
(envelope apierror + hasCode — idioma pós-CHIP-ERROR-UNIFY), nunca some silencioso.

## Negative Scenarios

- Health 500 → card com estado de erro visível; resto da página INTACTO (card isola o
  fetch).
- Entity desconhecida no payload (futura) → renderiza genérica (não filtra por lista
  hardcoded).

## Ownership

- Owned paths: `apps/web/src/pages/integracoes/SyncHealthCard.tsx` (novo), 1 linha em
  `IntegracoesPage.tsx`, hook web-query se necessário.
- Forbidden paths: outras seções da página; SDK (F-01 publicou); BE.
- Parallel-safe with: none — depends on F-01.

## Validation Expectations

- Fixture por estado (verde/vermelho/nunca/webhook-inicial) — 4 renders provados.
- Fixture negativa F-r04-1: entidade incremental-only (full velho, incremental recente) →
  badge VERDE + timestamp relativo recente — NUNCA "há N dias"/cinza (o payload já traz
  `last_success_at = GREATEST(full, incremental)` — IC-05).
- tsc + teste de página verdes.
- Live-drive hub: products real na tela com timestamp verdadeiro (confere com SELECT de
  sync_state).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md`; screenshot-métrica.
- Blockers or open decisions: none.
