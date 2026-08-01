# P3 reconciliation — r01 (2026-07-31)

Insumos: `p3-claude-candidate-r01.md` (10 ADRs, M-01→M-05) × `p3-opus-counterproposal-r01.md`
(waiver Claude-crew ratificado, seat Opus frio cego, 14 ADRs A-01..A-14, M-01→M-09).
Regra aplicada: `dual-model agreement` só quando a semântica material bate (rationale,
trade-offs, dependências, must-preserve); diferenças respondíveis por evidência resolvidas
aqui com citação; ao STOP só o que altera escopo/decisão irreversível/dependência de
milestone.

## 1. Agreement map (semântica material)

| Tema | Claude | Opus | Veredito |
| --- | --- | --- | --- |
| Núcleo nativo × adapter, teste da fronteira | ADR-01 | (herdado, transversal) | **agreement** — design §3 é autoridade p/ ambos |
| Read = Postgres, zero ML no read como fim | ADR-02 | A-05 (fim idêntico, timing difere — §2.1) | agreement no FIM; diferença no QUANDO (resolvida §2.1) |
| Ingest único idempotente, webhook+scheduler = 2 portas | ADR-03 | A-04 (mais afiado: resource-addressed, enumeração separada) | **agreement**; adota refinamento A-04 (assinatura `Ingest{Order,Listing}(ctx, tenant, installation, providerResourceID)`, enumeradores só produzem IDs) |
| 3 camadas fee + proveniência; degrau-4 vira fallback `config`; seed já morto não é trabalho | ADR-04 | A-09 (ledger enumerado: já-morto / morre / sobrevive / não-adotado) | **agreement**; adota ledger A-09 integral — inclui 2 itens que meu candidato não nomeou: `auth_adapter.go:47-48` baseline 0.16 morre-ou-vira-row-`config`; `FeeSyncScheduler`/`RegisterFeeSyncerFactory` = NÃO adotado (channel_fees é alimentada pelo ingest, nunca pelo fee-syncer morto) |
| Divergences tabela dedicada, ingest-time, auto-resolve, timestamps dos 2 lados | ADR-05 | A-10 (mais afiado: one-open-row-per-(entity,kind), upsert, NOT NULL nos 2 timestamps, CHECK em kind) | **agreement**; adota shape A-10 — resolve o conflito latente append-events × one-open-row ANTES de existirem 2 produtores (cadências 5min × diário) |
| Resiliência centralizada antes de backfill | ADR-06 | A-01 (mais afiado: decorator em `doRawWithHeaders`, token-bucket POR INSTALLATION compartilhado entre goroutines, limite configurável pois fato #11 é `assumed`, opt-out no-retry p/ writes futuros) | **agreement**; adota refinamentos A-01 + must-fail que NOMEIA o tempo esperado (asserção "eventually succeeds" = passe vácuo) |
| Raw DTO persistido | ADR-07 (persistir em listings/orders/order_shipments) | A-03 (idem MENOS billing_info: raw de PII NUNCA persiste; marcador de truncamento) | agreement parcial; **Opus vence por evidência** — §2.2 |
| MASS-CLOSURE morre, absent ≠ closed, marca só após run completo | ADR-08 | A-06 (idem + status = verdade do provider SOMENTE; `ApplyCompletedPull` continua writer ÚNICO, muda semântica não adiciona writer) | **agreement**; adota precisões A-06 |
| Webhook fino + inbox + hint não-confiável + `orders_v2` + IP log-only | ADR-09 | A-11 (idem + always-200 incl. malformado/topic desconhecido — evita tempestade 8-retries fato #7; user_id→installation via credencial própria; classe interativa EXPLÍCITA fora de `registerBatchRoutes`; 1ª superfície pública não-autenticada de escrita) | **agreement**; adota precisões A-11 |
| Migrações 0086+ aditivas, ranges pré-alocados pelo hub | ADR-10 | A-12 (idem + range no brief; regra: migração aplicada nunca renomeia — `schema_migrations` PK = filename; 0021 não se "limpa") | **agreement**; adota A-12 |

Itens só do Opus, sem contraparte no meu candidato — todos respondíveis por evidência P2,
todos ADOTADOS:

- **A-02** endpoints ML novos em arquivos novos, `capability_adapter.go` congelado pós-M-01
  (precedente já existe: `shipping_reader.go`, `buyer_fiscal_reader.go`). Gate mecânico:
  diff-stat ~zero fora do M-01. Mata a colisão mais quente da missão.
- **A-07** cursor terminal com phase explícita, `nil` proibido — `scheduler.go:42-45` APAGA
  cursor nil ⇒ re-backfill infinito silencioso. Meu candidato não viu. Bônus A-07: enum 0075
  já contém `listings`/`orders` — milestone que propuser migração p/ registrar job está errado.
- **A-08** dois Schedulers instanciados (5min/diário), `schedule jsonb` continua não-lido;
  fix do `incremental=false` (`scheduler.go:160`) na fundação — pré-condição do Q4
  (`last_incremental_at` vazio = QA reprova /integracoes).
- **A-13** `listing_variations` ADITIVA, PK de `listings` NÃO muda nesta missão (mudança
  radiaria destrutivamente p/ `product_link_listing_snapshots` 0022 e `product_links` 0025);
  risco nomeado: reestruturar hidratação pode ESFOMEAR `AbsorbProviderSnapshots`
  (`connectors/source.go:54-89`) e degradar re-vínculo — must-fail "snapshot observer starved".
- **A-14** `root.go` + par OpenAPI/SDK = seams serializados pelo hub; SDK é ESCRITO À MÃO
  (client = um objeto literal `index.ts:2113-2330`) ⇒ intuição "SDK gerado paraleliza" é
  INVERTIDA — no máximo 1 milestone com contrato FE em voo.

## 2. Diferenças resolvidas por evidência (sem STOP)

### 2.1 Timing da morte dos 4 sítios — Opus vence, com decisão ratificada PRESERVADA

Meu ADR-02 já mitigava (M-02 persiste shipment na mesma milestone); A-05 generaliza em regra:
**nenhuma fonte viva morre antes do substituto persistido no MESMO merge** + guard allowlist
encolhente (teste arquitetural, 4 entradas, cada milestone deleta a sua no mesmo commit,
zero no fechamento). Evidência da assimetria: sítios C/D (pricing) têm fallback honesto já
existente (`pricing_tariff_defaults` + proveniência `config` — mission.md:36-37) ⇒ podem
morrer cedo; sítios A/B (orders) são produtores ÚNICOS de SLA/bucket/frete/rastreio/
comprador-fiscal (`http_handler.go:549-567`, `:462-467`) ⇒ só morrem com `order_shipments`
populada no mesmo merge. O FIM ratificado (zero ML no read, <2s) não muda; muda o quando —
autoridade de implementação, não de produto. "Onda 0 como chip primeiro" (design §6) vira
"decoupling embutido no milestone que persiste o substituto" — semanticamente o que meu
candidato aprovado já fazia; A-05 só o torna regra verificável.

### 2.2 Raw de billing_info — Opus vence

`billing_info` carrega documento + endereço fiscal do comprador (`buyer_fiscal_reader.go:59-94`,
fato #16). Persistir raw dela viola a barra de PII do AGENTS.md e compõe com a pendência de
scrub já aberta em `docs/design/evidence/ml-api/`. Regra final: raw persiste em
`listings`/`orders`/`order_shipments`; **raw de billing_info NUNCA**; campos tipados sim;
truncamento ganha marcador explícito (prefixo que ainda parseia = mentira). Teste: fixture COM
documento → nenhuma coluna persistida o contém. Vira risco de missão (R-6 abaixo).

### 2.3 Verificação de manifest

Opus reproduziu o mismatch de mission.md (drift disclosed pós-freeze: churn de status +
registro do waiver) e MATCH nos outros 4. Divergência isolada e explicada; sem re-freeze —
o conteúdo que o Opus leu é o autoritativo (gate P1 + correção P2 presentes).

## 3. Diferença MATERIAL — corte de milestones (5 → 9)

Única diferença que altera estrutura/dependência de milestone ⇒ volta ao operador
(condição registrada na própria aprovação do P3: "mudança material volta para você").

Corte Opus re-corta as MESMAS ondas por posse de seam (1 dono por superfície compartilhada):

- **M-02 `sync-core-seam` extraído** — DDL 0086-0088 (`channel_fees`, `divergences`,
  `order_shipments`) + ports + fix `incremental` + guard allowlist, UM dono, antes de
  qualquer produtor. No meu corte, `channel_fees` nascia no M-03(listings) e era consumida
  pelo M-04(orders) — acoplamento cross-milestone num shape compartilhado, exatamente a
  classe de defeito que A-10 nomeia.
- **Onda 0 invertida** — persistir shipment (M-03) É o decoupling de /pedidos (§2.1);
  pricing decouple vira M-07 separado (dono do seam `pricingtarifflive`/`root.go:845-851`,
  distinto do dono de orders); pode rodar cedo via fallback `config` (edge de qualidade de
  dado, não de compilação).
- **Listings dividido** — M-04 ingest (writer, cursor, MASS-CLOSURE, variations) ∥ depois
  M-05 fees camada-2 + divergência + FE /anuncios (contrato browser próprio).
- **M-09 observabilidade solto** — `sync_state` 0075 + `GET /sync/runs` já existem, zero
  consumidor FE; construível desde o dia 1 contra entity `products`.
- **Lanes**: A = M-01∥M-02∥M-09 → B = M-03∥M-04 → C = M-05∥M-06∥M-07 → D = M-08.
  3 workers sustentados vs cadeia serial do meu corte (M-02∥M-03 era o único paralelismo).

Avaliação do reconciliador: corte Opus é superior por evidência — cada lane tem dono
disjunto nos 6 eixos; edges nomeados por artefato; regression window da Onda 0 eliminada
por construção; overhead = 4 fechamentos de milestone a mais (dual-gate + QA cada), pago
pelo paralelismo 3× e pela ausência de re-trabalho de seam. Conteúdo de ESCOPO idêntico ao
aprovado — nada entra, nada sai; muda o fatiamento e a ordem interna.

**Recomendação: adotar corte M-01→M-09 do Opus** com escopo aprovado inalterado.

## 4. Riscos consolidados (fold p/ mission.md no P4)

R-1..R-5 existentes permanecem. Entram do Opus:

- **R-6 (novo, de R-D)** PII escapa via raw-persistência + endpoint público novo.
  Mitigação: regra §2.2 + corpo de webhook nunca alcança tabela de domínio (A-11).
  Trigger: qualquer migração adicionando `raw jsonb` a tabela com dado fiscal/comprador.
- **R-7 (novo, de R-C)** colisão de merge em seam sem dono: numeração de migração, `root.go`,
  SDK escrito à mão + regra same-commit OpenAPI. Mitigação: A-12 ranges no brief, A-14
  constructor local + lock de contrato do hub, ≤1 milestone FE-contract em voo.
  Trigger: 2 briefs despachados sem ranges disjuntos, ou 2 milestones em voo listando OpenAPI.
- R-A/R-B do Opus = já cobertos por A-05/A-06 + R-1/R-3 existentes (não duplicar);
  R-E = coberto por A-10 + R-5.
- Afiados de 2ª linha (entram nos briefs, não na tabela de missão): nil-cursor (A-07),
  starvation `AbsorbProviderSnapshots` (A-13), brief herdando claim obsoleto `fee_sync.go:29`
  (A-09 — critério infalsificável).

## 5. Disposição

- Spine reconciliado = 14 entradas (meus 10 ⊕ refinamentos + A-02/A-07/A-08/A-13/A-14),
  a folhar em mission.md no P4.
- STOP: 1 decisão ao operador — corte de milestones (manter M-01→M-05 aprovado × adotar
  M-01→M-09 Opus). Recomendação: adotar.
- **DECISÃO DO OPERADOR (2026-07-31, STOP P3 r01): ADOTAR M-01→M-09.** Escopo aprovado
  inalterado; corte por seam ratificado com lanes A(M-01∥M-02∥M-09) → B(M-03∥M-04) →
  C(M-05∥M-06∥M-07) → D(M-08).
- Pós-decisão: `planning_phase: architecture`; Sol P3 retroativo (≥2026-08-05) audita ESTE
  artefato + os dois insumos, antes de `planned`.
