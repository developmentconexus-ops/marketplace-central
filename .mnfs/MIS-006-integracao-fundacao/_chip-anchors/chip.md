# CHIP-ANCHORS — âncoras de vínculo declaradas pelo provider (+ D-A, D-B)

```yaml
id: CHIP-ANCHORS
type: chip
mission: MIS-006 (hardening — NÃO é milestone da DAG)
status: dispatched
base_sha: 917f7bb58e385847fba5612201823f9db48791c6
branch: (chip cria a sua, prefixo `chip/`)
validation_level: QA-0
```

## Por que este chip existe

O gerador de candidatos de vínculo assume a **forma do Mercado Livre dentro do núcleo**.
`generation_service.go` chama `mandatoryUnavailableReasons()` e emite, em TODA linha de TODO
provider, dois motivos:

```go
{Anchor: "marca",   Direction: UNAVAILABLE, Detail: "marca inexistente no lado provider"},
{Anchor: "refforn", Direction: UNAVAILABLE, Detail: "refforn inexistente no lado provider"},
```

Hoje isso é verdade por acidente (só existe um provider e ele realmente não fornece esses
campos). No marketplace 2 vira **mentira em produção**: o núcleo vai afirmar que o provider não
fornece marca para um provider que fornece. É o buraco que estamos cavando, e o custo de sair
dele cresce com cada consumidor novo desses motivos.

A missão MIS-006 é fundação de produto/ERP/xlsx — capability multi-marketplace está **fora do
escopo dela**. Por isso este trabalho é chip de hardening sobre `main`, não milestone nova: emendar
a DAG ratificada no P7 por trabalho que não é da missão seria pior.

Dois defeitos PRE-EXISTENTES do mesmo módulo entram no mesmo chip porque vivem nos mesmos dois
arquivos (verificados presentes no BASE-SHA `e3c081ae`, ou seja **não** são regressão do M-05 —
o M-05 só os tornou visíveis ao criar 29 vínculos onde havia 1):

- **D-A** — `hardNegativeDimension()` (`generation_service.go:680-696`) normaliza vírgula e
  espaço mas **não normaliza unidade**: `50cm` e `500MM` são lidos como dimensões diferentes, o
  guard de hard-negative declara contradição e reprova um par correto. Caso real na conta
  conectada: `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` (BAIXA 25%, motivo
  `Título hard-negative: medidas`).
- **D-B** — `ListLinkWorkflows` (`resolution_service.go:368-380`) aplica **um limite
  compartilhado** a duas listas independentes (candidatos e vínculos) e `limit*5` aos audits. O
  default é 20 (`transport/http_handler.go:251`). Com 29 vínculos resolvidos, 9 nunca chegam na
  resposta e a KPI de `/vinculos` apresenta 20 como se fosse o total.

## Escopo

### F-01 — Âncoras como capability declarada pelo provider

- O lado provider **declara** quais âncoras de identidade consegue fornecer. O registro de
  capability já existe: `connectors/ports/marketplace_capability.go:9-20`
  (`CapabilityListingRead`, `CapabilityStockRead`, …) — a declaração de âncoras se pendura nessa
  mesma forma, não numa invenção paralela.
- `product_links` **consome via port próprio** (hexagonal, mesma forma de
  `ProductSourceAdapter.Kind()` do M-02): port em `product_links/ports/`, adapter que lê a
  declaração do provider, wiring na composition. O gerador itera a declaração.
- `mandatoryUnavailableReasons()` **morre**. Nenhum nome de âncora fica hardcoded no gerador.
- Depois do fix, `UNAVAILABLE` precisa distinguir dois estados que hoje são um só:
  (a) o provider **não fornece** essa âncora; (b) o provider fornece e **este anúncio não tem
  valor**. A distinção vai no `detail` (texto), **nunca** num quarto valor de `direction` —
  `direction` é enum de wire (ver Rulings).

### F-02 — D-A: unidade canônica no hard-negative de dimensão

- Canonicalizar para uma unidade única (mm) antes de comparar: `cm`×10, `m`×1000, `pol`×25,4.
- O guard **continua bloqueante**. Ele existe para "contradição vence EAN"
  (`generation_service.go:436-444`) e é a única coisa que impede um par CODPROD+EAN concordante
  errado de auto-aprovar. Afrouxá-lo é falha desta feature, não conserto.

### F-03 — D-B: limites independentes em `ListLinkWorkflows`

- Candidatos, vínculos e audits passam a ter limites independentes; o limite de uma lista não
  trunca a outra.
- O **guard de truncamento no FE é do M-06**, não deste chip (`VinculosPage.tsx` é write-set
  exclusivo do M-06). Este chip faz a API devolver o conjunto certo; a tela honesta é do M-06.

## Non-Scope (aparecer no diff = escopo vazando)

- **Qualquer arquivo em `apps/web/`** — o FE inteiro é write-set exclusivo do M-06 (colunas
  provider-neutras, "Identificado por", badge de provider, guard de truncamento). Zero exceção.
- **Qualquer mudança de shape em OpenAPI / `packages/sdk-runtime`** — ver Rulings R1.
- **Migração** — nenhuma. Zero arquivo em `migrations/`.
- **Backfill / rewrite dos candidatos já persistidos** — ver Rulings R3.
- Segundo provider de verdade (adapter Shopee/Amazon). Este chip entrega a **capability**; o
  segundo adapter é missão futura. Não escreva adapter novo "de exemplo" no código de produção.
- Renomear `provider_code`/`provider_item_id`/`match_input`/`match_value` no wire — já são
  agnósticos, estão certos.

## Rulings do hub (vinculantes — divergência é `BLOCKED`, nunca decisão unilateral)

- **R1 — o wire não muda de forma.** `reasons[]` continua `{anchor, direction, detail}` e
  `direction` continua o enum `FOR|AGAINST|UNAVAILABLE`. Motivo: o M-06 carrega um
  **contract-lock** na seção chain-read do OpenAPI e vai rodar em seguida; mudança de shape aqui
  colide com ele. Se você concluir que a feature exige mudança de wire, isso é `BLOCKED` com
  evidência — o hub decide, não o chip.
- **R2 — capability é dado, não `switch`.** Zero `if provider_code == "mercado_livre"` (ou
  equivalente) dentro de `product_links`. O gerador não pode saber o nome de nenhum provider.
- **R3 — histórico não se reescreve.** Candidatos/motivos já persistidos ficam como estão:
  nenhum `UPDATE` de motivo, nenhuma migração de backfill. A geração nova produz motivos novos
  no próximo refresh. Mesma doutrina do ADR-01 (protocolo de import é histórico imutável).
- **R4 — regra de terceira rodada** (`docs/HARNESS-PROFILE.md` §11, ratificada): terceiro defeito
  da mesma forma (ou terceira rodada de correção no mesmo critério) para de remendar — nomeia a
  mecânica em uma frase, varre a classe inteira por ferramenta com `file:line` **incluindo os
  sites limpos** (a enumeração é o artefato), prova a CLASSE, e pede opinião adversarial
  independente briefada CONTRA qualquer abstração nova. Point-fix de 3ª rodada sem varredura o
  hub não mergeia.
- **R5 — must-fail obrigatório nos guards.** F-02 e F-01 mexem em código de bloqueio. Cada guard
  precisa de prova de que é load-bearing: reverter o fix faz o teste ficar vermelho. Teste que
  passa com e sem o fix não é evidência.

## Ownership & concurrency (seis eixos)

| Eixo | CHIP-ANCHORS |
|------|--------------|
| Migração | **nenhuma** |
| DB shape | **nenhuma** (lê `link_candidates`/`product_link_*` existentes) |
| Módulo Go | **dono**: `product_links` (domain/application/ports/adapters) + a declaração de âncoras no lado `connectors` |
| `root.go` | **grant additive pré-autorizado**: só a sua região de composition-wiring do port novo. Additive-only, liberado no CLOSED, diff citado no payload. Nenhum outro track escreve `root.go` agora (M-02 fechou; M-06 declara "nenhum `root.go`") |
| Contrato/SDK | **nenhum** — R1. OpenAPI e `packages/sdk-runtime` byte-idênticos |
| FE surface | **nenhum** — exclusivo do M-06 |

- Roda em paralelo com: nada. Serial antes do M-06, **de propósito**: o QA de tela do M-06 tem
  que ver o conjunto final de âncoras. Se rodarem juntos, o operador valida colunas alimentadas
  por motivos que mudam embaixo dele.
- DAG interna: F-01 ∥ F-02 (arquivos vizinhos, mesma unidade — se colidirem textualmente,
  serialize F-02 depois de F-01). F-03 é independente das duas (`resolution_service.go`).

## Ladder + gates

- **L0/L1/L2** por `docs/HARNESS-PROFILE.md` §2 — inclusive `GOCACHE` em caminho **absoluto**
  (D-14) e governance de **worktree limpo detached** com BaseSha de **40 hex**
  (`917f7bb58e385847fba5612201823f9db48791c6`). Allowlist de falhas pré-existentes: **cite**, não
  re-prove.
- **L2 / dev stack é seam do hub.** Nunca suba servidor, nunca ocupe `:8080`/`:5174`, nunca
  carregue `.env*` no ambiente da sessão. Precisa da stack ⇒ `REQUEST`.
- **P6 dual-gate** (Opus cold + Sol medium no diff de SHA fixo) + **P7**: este chip é backend, então
  o live-drive de tela é **rodado pelo hub** com a conta ML conectada (critérios `U1-U3` no
  `validation-contract.md`). O chip entrega o pack; não tenta dirigir browser.
- Codex: a parede de quota resetou em 2026-07-25, então P2/P3 seguem a matriz normal do core §1.

## Evidência exigida

`_chip-anchors/EVIDENCE.md` com: tabela `file:line` por critério, resultados do ladder (cite a
allowlist), veredictos do dual-gate + reconciliação, provas de must-fail (R5), **dispatch ledger**
(fechamento sem dispatch de planner/implementer/reviewer o hub reprova), e o diff do grant de
`root.go` citado explicitamente.
