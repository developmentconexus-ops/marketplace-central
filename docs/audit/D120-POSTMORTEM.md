# D-120 Postmortem — Auditoria pós-demo (2026-07-20)

Status: **DRAFT para verificação conjunta**. Cada achado tem evidência `file:line` e um passo
"Como verificar". Nada aqui é opinião sem citação; achados inferidos (não observados ao vivo)
estão marcados `[INFERIDO]`.

Origem: 5 auditorias read-only paralelas (anúncios, pedidos, mercado, integração/arquitetura,
uso da API ML) na sessão de replanejamento D-120.

---

## 0. Achado-mestre: árvore/stack desatualizados

| ID | Achado | Evidência | Como verificar |
|----|--------|-----------|----------------|
| A0.1 | Working tree em detached HEAD `74104f79`, **37 commits atrás de `main` (`138aac3d`)**. `/mercado`, faturado (0074), scan-paging e `/integracoes` real só existem em `main`. | `git log --oneline HEAD..main \| wc -l` → 37 | `git status`, `git merge-base HEAD main` |
| A0.2 | Stack `:8080` monta `.:/workspace` (cwd-relative) — se serviu este checkout, a demo rodou código velho. Padrão deploy-gap já provado no CHIP-PED-FIX (comissão "—" no pedido golden). | docker-compose mount; memória `dev-stack-mount-cwd-relative` | Dentro do container: conferir commit servido vs `main`; comparar resposta raw de endpoint com código de `main` |

**Ação imediata:** confirmar commit servido na demo; re-apontar checkout + stack para `main`.

---

## 1. Causas sistêmicas (S1–S5)

### S1 — Sem estratégia de sincronização
Persistimos listings e orders em Postgres, mas **nada mantém o banco atualizado**: sem webhook,
sem scheduler de dados, sem reconciliação. Dado entra por clique manual e apodrece.

| Evidência | Como verificar |
|-----------|----------------|
| Únicos tickers em `composition/root.go:574-576`: OAuth refresh (5min), state cleanup (1h), fee sync (15min — grava só 16%/22% hardcoded, `fee_sync.go:13-55`). Nenhum puxa listings/orders/mercado. | `grep -rn "Ticker\|Scheduler" apps/server_core/internal/composition/` |
| `webhook_receive` declarado como capability (`auth_adapter.go:64`, `runtime_capability.go:21`) mas **nenhuma rota** `/webhooks`/`/notifications` registrada no mux; nenhuma chamada de subscription a tópicos ML. | Enumerar `mux.Handle` em `root.go`; `grep -rn "webhook" apps/server_core/internal --include=*.go` |

### S2 — Tudo um-a-um sequencial, sem resiliência
| Evidência | Como verificar |
|-----------|----------------|
| Zero uso de multiget `GET /items?ids=` no repo inteiro. | `grep -rn "ids=" apps/server_core/internal/modules/connectors` |
| Fan-outs são `for` sequencial sem goroutine/errgroup: listings `capability_adapter.go:244-254`, orders `capability_adapter.go:459-473`, catalog children `catalog_offers_reader.go:56-79`. | `grep -rn "go func\|errgroup\|WaitGroup" internal/modules/connectors internal/modules/market internal/modules/orders` → vazio |
| 429 reconhecido (`capability_adapter.go:607-609`) mas **sem retry, sem backoff, sem Retry-After, sem rate limiter**. `http.Client{Timeout: 15s}` puro (`capability_adapter.go:62`). | `grep -rn "Retry-After\|backoff" apps/server_core/internal` → só política OAuth |
| Refresh de 2006 anúncios = 21×`/items/search` + 2006×`GET /items/{id}` = **2027 round-trips sequenciais** por clique. | Contar chamadas em `ListListings` + `ingestion.go:34-71` (pageSize 100, `root.go:623`) |

### S3 — Enrichment no read, não no ingest
| Evidência | Como verificar |
|-----------|----------------|
| Lista de pedidos: `handleReadList` → `Enrich()` faz até **3 calls ML sequenciais POR pedido** (shipment + costs + carrier: `shipping_reader.go:147-211`) em `enrich_service.go:108-130`. FE agrava puxando até 20 páginas (`PedidosPage.tsx:32,41-56`). ~10s explicados. | Cronometrar `GET /orders` com N pedidos com shipping_id; ler loop |
| O próprio código nomeia a classe do bug ("FINDING-M08-LIST-TIMEOUT", `enrich_service.go:101-107`) — fiscal foi movido pro detalhe, shipment ficou no list. | Ler comentário |

### S4 — Features terminam em "em breve" / pontas soltas
| Evidência | Como verificar |
|-----------|----------------|
| Decomposer wired **`nil`** em `composition/root.go:513` → Decomposição/DIFAL/Margem/Retorno = "—" para TODO pedido, sempre (idêntico em HEAD e main). | Ler linha; abrir qualquer pedido |
| Filtro de data de pedidos pronto em backend (`orders/transport/query.go:99-107`), SQL (`order_repo.go:64-65`), SDK (`index.ts:750-757`), OpenAPI — botão FE `disabled` "filtro de período em breve" (`PedidosPage.tsx:178-185`). | Ler as 4 camadas |
| `/mercado` (main): MARGEM ATUAL / SE IGUALAR MENOR / SUGESTÃO / MARGEM EST. / VEREDICTO / VENDAS LÍDER = dash hardcoded (`RepricingTable.tsx:83-93`, `OportunidadesTable.tsx:78-86`; `verdict_label` "Always null from M-02", `market.ts:126-127`). Aba Monitorados sem backend algum. "Atualizar agora" `disabled`. | Ler componentes em `main` |
| Endpoints do pipeline de vínculos (`POST /product-links/listing-snapshots/imports` e `.../link-candidates/generations`, `http_handler.go:89-90`) — **zero callers** em FE e composition; só curl. | `grep -rn` pelas operações em `apps/web/src` |

### S5 — Sem onboarding: conta nova = sistema vazio sem caminho
| Evidência | Como verificar |
|-----------|----------------|
| Callback OAuth (`auth_flow_service.go` ~673-679) só seta `Connected` + salva credencial. Não dispara import de listings, geração de vínculos, nem coleta de mercado. | Ler `applyAuthResult` |
| `InstallationGate` (`InstallationContext.tsx:78-95`) bloqueia TODAS as rotas de tenant sem installation, apontando para `/integracoes` — que em HEAD é placeholder sem botão de conectar. | Subir app com 0 installations |
| Coleta de concorrência: `POST /market/collections` = 1 codprod, síncrono, sem batch ("D-F4-p: no batch POST", `collection_handler.go:13-14`), sem scheduler, sem botão. | Ler handler; grep Ticker no módulo market |

---

## 2. Achados por módulo

### Pedidos
| ID | Achado | Evidência |
|----|--------|-----------|
| P1 | FE **nunca** chama `/orders/import`; sem cron. DB só tem o que foi importado out-of-band. | grep `importMarketplaceOrders` em `apps/web` → 0 matches |
| P2 | Import ML = sempre `offset=0, limit=20 (default), sem sort, sem filtro de data` → pedidos de hoje podem nunca entrar. | `provider_operation_service.go:126-139` nunca seta Cursor/Status/UpdatedAfter; `capability_adapter.go:439-455` |
| P3 | Fila mostra TODOS os buckets: `FilaView.tsx:60-61` só ordena (`sortByUrgency`), sem `.filter()` — Lista filtra certo (`pedidosTabs.ts:30-33`). | Ler componente |
| P4 | Fila (view default) não tem coluna de data; Lista tem (`PedidosTable.tsx:104-108`). | Abrir /pedidos |
| P5 | `tags` do ML persistidas mas nunca expostas no DTO; `pack_id` e `context.channel` **dropados** no DTO do pedido. | `capability_adapter.go:976-1034` |
| P6 | Money fields sem currency (`UnitPrice`, `SaleFeeAmount`, payments — `capability.go:261-278`). | Ler struct |

### Anúncios
| ID | Achado | Evidência |
|----|--------|-----------|
| L1 | `date_created` do ML nunca capturado — `CreatedAt` recebe `fetchedAt` (`mapper.go:82-100`). Filtro por data é impossível: o fato não existe. | Ler mapper |
| L2 | `[INFERIDO]` Truncation do search ML (offset cap) faria `ApplyCompletedPull` **fechar em massa** anúncios ausentes do pull (`repository.go:386-393`), sem erro. Scan-paging (`scroll_id`) inexistente neste checkout (fix em main — confirmar). | grep `scroll_id` na branch servida |
| L3 | Campos dropados do `/items/{id}`: category_id, sold_quantity, sub_status, condition, permalink, thumbnail, tags, catalog_product_id, shipping block, date_created. | `capability_adapter.go:925-938` |
| L4 | Refresh roda em goroutine com `context.Background()` — sem deadline nenhum (`refresh_service.go:64`, `root.go:626`). | Ler wiring |
| L5 | Polling do FE (2s) lista `integration_operation_runs` **sem LIMIT** — cresce para sempre (`operation_run_repo.go:112-138`). | Ler SQL |
| L6 | Status filter do ML (`ListListingsInput.Status`) existe no domain e nunca é enviado (`capability_adapter.go:232-235`). | Ler adapter |

### Mercado (auditado em `main`)
| ID | Achado | Evidência |
|----|--------|-----------|
| M1 | Oportunidades: unit JÁ é produto ERP (codprod) — `oportunidades.ts:24-26`. Mas label "não vendemos" **não é aplicada**: sem exclusão de codprods vinculados/anunciados. | `buildOppRows` não consulta product_links |
| M2 | "NOSSO PREÇO" na aba Oportunidades = **custo ERP**, rotulado como preço (`OportunidadesTable.tsx:7-8` admite). "SKU" = codprod interno, não referência. | Ler componente |
| M3 | Critério de corte: aggregate `status=OK` exige ≥5 sellers distintos (`aggregation.go:44-58`) — produtos com 1-4 concorrentes somem sem explicação na UI. | Ler agregação |
| M4 | `signal_status` (OK/STALE/SEM_VINCULO/NO_PRICE_EVIDENCE) computado no backend (`signal.go:65-76`) e **não renderizado** em `/mercado` — stale = fresco na tela. Timestamp "coletado" mostrado vem do read de listings, não da coleta de mercado. | Ler `RepricingTable.tsx` / `MercadoPage.tsx` |
| M5 | Exclusão own-seller vira no-op silencioso se `ProviderAccountID` da installation estiver vazio (`market_adapters.go:307-318`, `collection_pipeline_service.go:248`). | Ler resolver |
| M6 | Offers de catálogo dropam `item_id` do concorrente (`catalog_offers_reader.go:29-39`) — agregado nunca rastreável ao anúncio. Buy-box winner decodificado por **2 structs distintos e divergentes** (`catalog_identity_reader.go:31-35` vs `catalog_match_reader.go:160-164`). | Ler structs |

### Integração / arquitetura
| ID | Achado | Evidência |
|----|--------|-----------|
| I1 | Dois eixos chamados "source": `MC_ERP_SOURCE` (boot-time, oracle\|xlsx, backend inteiro) vs `erp_import.source` (por upload, xlsx\|catalogo_cliente). Confusão conceitual raiz do "desconexo". | `root.go` `erpSource()`; `erp_import/domain/import.go` |
| I2 | Em HEAD: `WithActiveSource` nunca chamado → `catalogo_cliente` nunca vira dataset ativo (fix existe em main via IMPORT-FIX — confirmar wiring). | grep `WithActiveSource` |
| I3 | Cache key do internal_read sem tenant nem source (`cache.go` `canonicalKey()`) — mascarado hoje por serem constantes. | Ler função |
| I4 | Single-tenant na prática: `cfg.DefaultTenantID` em ~15 construtores; sem resolução por request. | grep `DefaultTenantID` |
| I5 | Simulador de preços **cego a concorrência**: `pricing/adapters/*` nunca lê market. | Ler readers |
| I6 | Matcher de vínculos: sem EAN nunca passa de REVIEW (SKU→70, título→35); geração 100% desacoplada de aprovação (correto por policy, mas sem trigger automático de geração). | `generation_service.go` |
| I7 | Engrenagem: `<details>` dropdown no `Header.tsx` com links para rotas normais no mesmo Layout — não é shell aninhado; é UX (configurações sem contexto próprio). | Ler Header |
| I8 | Token: `ResolveAccessToken` não checa expiry nem faz refresh-retry; expiry mid-request = erro seco até próximo tick de 5min (`credential_resolver.go:45-79`; ticker engole erros, `refresh_ticker.go:63-64`). | Ler resolver |

### Contra-exemplos (o que ESTÁ certo — preservar)
- Orders → product_links: acoplamento real e funcional (`orders/adapters/productlinks/link_reader.go`).
- Persistência local como fonte da UI (listings/orders/market em Postgres) — arquitetura certa, faltou o motor de sync.
- Anti-corruption entre listings e market via bridge de composition (`market_adapters.go:599-679`).
- Policy de vínculo conservadora (contradição lexical vence EAN; título nunca auto-aceita).
- 4 armadilhas de API ML já corrigidas e documentadas no código (`?context=channel_marketplace`, `?site_id=`, `x-format-new: true`, `?version=v2`).

---

## 3. Mandato de replanejamento (ratificado pelo operador D-120)

Por módulo, nível micro, nesta ordem:
1. Desenhar fluxograma do processo completo.
2. Ler docs ML a fundo (endpoints, campos, limites, webhooks).
3. Testar endpoints live contra conta real — provar que o campo vem antes de modelar.
4. Ratificar contrato DTO + modelo de persistência + mecanismo de alimentação.
5. Implementar.
6. Testar com fixtures multi-página (nunca só live-drive).

Ordem de módulos: **Base/Integração → Pedidos → Anúncios → Mercado.**
Documento irmão: `docs/design/INTEGRATION-DATA-CONTRACT.md` (contrato de dados da base).
