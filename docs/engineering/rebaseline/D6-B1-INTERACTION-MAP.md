# D6-B1 — Mapa de Interações do Frontend

> **Status:** CANDIDATO DERIVADO — App Shell aprovado; mapa rederivado após D6-R1 / OAD 99 operações · 30 Permissions; ainda não é ratificação de D6-B1
> **Pai:** `D6-FRONTEND.md`
> **Autoridade de wire:** `contracts/api/product/openapi.yaml`
> **Escopo:** frontend Product 1.0; sem mecânicas D7 e sem implementação Product
> **Revisado:** 2026-08-21

## 1. Objetivo

Este artefato prova a cadeia:

```text
necessidade do usuário
  → navegação / estado de tela
  → semantic owner
  → operação/capability Product canônica
  → Permission ordinária
  → tratamento de identidade, conhecimento e segurança
```

As 99 operações Product **não** significam 99 telas. Uma tela pode compor leituras de vários owners, desde que essa composição nunca ganhe write authority.

## 2. Arquitetura de informação aprovada e refinada

```text
VISÃO GERAL

OPERAÇÕES
  Preparação
  Publicações
  Disponibilidade
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
```

### 2.1 Contexto global

`Organization` é o único workspace global. `organization_id` é autoridade de escopo; `display_name` é apresentação.

Trocar Organization:

- re-resolve as queries server-owned;
- limpa Marketplace Installation incompatível;
- nunca reaproveita Permission ou seleção de conta como autoridade implícita.

### 2.2 Contexto de Marketplace Installation

| Modo | Uso |
| --- | --- |
| `organization-wide` | a operação Product não tem dimensão Installation |
| `all-or-exact` | a operação admite filtro Installation opcional |
| `exact-required` | o contrato exige Installation exata; não existe merge sintético |

Listings, Sales, Shipments e as quatro leituras de Performance usam Installation exata quando seus paths assim exigem.

### 2.3 Performance não é Mercado nem Economia

```text
Performance → Como nossa participação está performando?
Mercado     → O que está acontecendo competitivamente fora de nós?
Economia    → O que esses fatos significam economicamente?
```

O workspace estratégico pode compor os três, além de Offering/Sales/Availability, mas continua somente leitura.

### 2.4 Todos os canais

Não há KPI global baseline de “Todos os marketplaces”. Métricas de providers diferentes não são somadas por nome. Uma futura visão multi-marketplace apresenta resultados por Installation até que equivalência de measurement basis seja realmente provada.

## 3. Gramática de navegação frontend

Isto é identidade de navegação, não escolha de router nem Product path:

```text
/org/:organizationId/visao-geral
/org/:organizationId/preparacao
/org/:organizationId/publicacoes/*
/org/:organizationId/disponibilidade
/org/:organizationId/performance/*
/org/:organizationId/mercado
/org/:organizationId/economia/*
/org/:organizationId/vendas/*
/org/:organizationId/expedicao/*
/org/:organizationId/pos-venda/*
/org/:organizationId/trabalho/*
/org/:organizationId/aprovacoes/*
/org/:organizationId/configuracoes/*
```

Recursos externos mantêm a qualificação de origem em navegação e requests: Installation + native key para Listing/Sale/Shipment. O URL nunca é business authority por si só.

## 4. Fluxos de usuário

### F1 — Adicionar e autorizar uma conta de marketplace

```text
Configurações / Canais
  → ListMarketplaceInstallations
  → Adicionar canal (Mercado Livre é o único atual)
  → CreateMarketplaceInstallation [Idempotency-Key]
  → se necessário: cerimônia OAuth D5 Technical Ingress para Installation exata
  → GetMarketplaceInstallation
```

Quando Performance de mídia for habilitada e o provider expuser mais de um advertiser elegível, a seleção/binding é cerimônia **Technical Non-Product** sob humano atual + `portfolio.manage`; nunca Product `ConnectAds` nem escolha automática por nome/primeiro resultado.

### F2 — Preparar um produto e submeter uma intenção de publicação

```text
Preparação / Installation exata
  → SearchSourceProductsForMarketplace
  → selecionar SourceInstance + native_product_key explícitos
  → GetProductChannelReadiness
  → GetPublicationRequirements
  → Resolve/ClearProductChannelCorrespondence quando necessário
  → CreateListingIntentDraft [Idempotency-Key]
  → Get/UpdateListingIntent [ETag / If-Match]
  → CreateListingIntentMedia quando necessário
  → CreatePriceIntent separadamente
  → GetSellableAvailability para target pre-creation
  → opcional EvaluatePriceScenario
  → SubmitListingIntent com revisão atual
  → observar accepted/pending/rejected/ambiguous/convergence sem inferir sucesso pelo transporte
```

A busca sem `source_instance_id` pesquisa fontes Readiness-admitidas; nunca escolhe uma fonte default escondida.

### F3 — Investigar Performance da conta

```text
Estratégia e Inteligência / Performance / Installation exata
  → selecionar período atual
  → opcional selecionar período comparativo equivalente
  → GetMarketplacePerformanceSummary
  → exibir tráfego, atividade de vendas e resumo de mídia
  → preservar complete/partial/unknown/unavailable/unsupported
  → exibir delta somente quando o Product declara comparable
```

O frontend não recalcula provider CVR/ROAS, não transforma evidência parcial em agregado completo e não confunde evidência histórica preservada com fato authored pelo MPC.

### F4 — Investigar Performance de uma publicação

```text
Performance / Publicações
  → ListMarketplaceListingPerformance
  → manter também Listings com Performance unknown/unavailable
  → selecionar Listing source-qualified
  → GetMarketplaceListingPerformance
  → cruzar visualmente, quando permitido, com Offering/Availability/Market/Economics
```

`FAMILY`, `CATALOG GROUP` e `CAMPAIGN` de Retail Media não são atribuídos a um Listing sem evidência suficiente.

### F5 — Analisar Mídia sem virar Ads Manager

```text
Performance / Mídia / Installation exata
  → ListRetailMediaPerformance
  → preservar scope:
       campaign | listing | marketplace_catalog_group | marketplace_family_group
  → exibir provider measurement basis / attribution / coverage
```

Não existem ações de criar campanha, alterar budget/bid/targeting/creative, pausar campanha ou “otimizar automaticamente”.

### F6 — Analisar mercado/economia e criar PriceIntent

```text
subject explícito
  → Get/ListCompetitivePosition
  → ListComparableOffers
  → Get/ListExpectedEconomics
  → EvaluatePriceScenario
  → CreatePriceIntent [Idempotency-Key]
  → GetPriceIntent / convergence
```

Market/Economics/Performance podem informar a decisão; Offering continua dono do PriceIntent.

### F7 — Observar e configurar Disponibilidade

Operação:

```text
Disponibilidade
  → ListSellableAvailability
  → GetSellableAvailability
```

Configuração:

```text
Configurações / Disponibilidade
  → List/Get/Create/Update/DeactivateInventorySource
  → GetEffectiveAvailabilityAllocationScopePolicy
  → UpdateAvailabilityAllocationScopePolicy
```

Sem `SetAvailableQuantity`, `SyncAvailability` ou AvailabilityIntent manual.

### F8 — Operar uma venda pelo ciclo cross-owner

```text
Vendas / Installation exata
  → ListMarketplaceSales
  → GetMarketplaceSale
  → ResolveSaleSellingEntityAttribution quando necessário

Detalhe da venda — composição read-only
  → GetSaleEconomics
  → List/GetBusinessOrderIntent
  → Get/ResolveBusinessSystemPartyResolution
  → GetDestinationRealization
  → List/GetInvoicingIntent
  → List/GetFulfillmentExecution
  → List/GetPostSaleResolution
```

Nenhum comando Sankhya, `CreateBusinessOrderIntent`, `CreateInvoicingIntent` ou retry genérico nasce da tela.

### F9 — Executar expedição física e observar Shipment

```text
Expedição / Execuções
  → ListFulfillmentExecutions
  → GetFulfillmentExecution
  → RecordSeparation
  → RecordPhysicalConference
  → RecordPacking
  → RecordDispatchHandoff
  → List/GetFulfillmentArtifacts quando permitido

Expedição / Envios / Installation exata
  → ListShipments
  → GetShipment
```

Visibility de botão não substitui Principal-kind nem qualificação física server-owned.

### F10 — Coordenar pós-venda

```text
Pós-venda
  → ListPostSaleResolutions
  → GetPostSaleResolution
  → CreatePostSaleResolution [Idempotency-Key]
```

Sem close/refund/cancel genérico; closure depende de evidência dos owners responsáveis.

### F11 — Tratar Trabalho operacional

```text
Trabalho
  → ListWork
  → GetWork
  → AssignWork / ClearWorkAssignment
  → HoldWork / ResumeWork / EscalateWork
```

Work não cria nem fecha source truth e não vira command bus.

### F12 — Aprovar e administrar acesso/configuração

Governance:

```text
Aprovações
  → List/Get/CreateAuthorizationDecision

Configurações / Delegações
  → List/Establish/Update/RevokeAuthorizationDelegation
```

Acesso:

```text
Shell → GetCurrentAccessContext
Configurações / Acesso
  → ListOrganizationMembers
  → ListAccessRoles
  → AssignAccessRole
  → RevokeAccessRole
```

Governance não concede Permission do target e AccessRole não vira IdP/provider role.

## 5. Inventário de estados de tela/rota

| ID | Tela / estado | Contexto de canal | Owners | Home de operações Product | Permission(s) |
| --- | --- | --- | --- | --- | --- |
| S00 | Shell / acesso atual | organization-wide | IdentityAccess | `GetCurrentAccessContext` | authenticated |
| S01 | Visão geral | organization-wide | composição Portfolio/Economics/Work | reads independentes conforme Permission | independentes |
| S10 | Preparação | exact-required | ProductChannelReadiness | 5 operações Readiness | `readiness.read/manage` |
| S20 | Publicações | exact-required | Offering | `ListMarketplaceListings` | `offering.read` |
| S21 | Publicação / Operação | exact-required | Offering | `GetMarketplaceListing` + painéis opcionais owner-local | `offering.read` + painéis |
| S22 | Publicação / Performance | exact-required | MarketplacePerformanceIntelligence | `GetMarketplaceListingPerformance` | `performance.read` |
| S23 | Intenções de publicação | all-or-exact | Offering | `ListListingIntents` | `offering.read` |
| S24 | Editor / detalhe de ListingIntent | target-explicit | Offering + reads contextuais | Get/Create/Update/Discard/Submit/Media | `offering.read`, `listing.manage` + reads |
| S25 | Intenções de preço | all-or-exact | Offering | List/Get/CreatePriceIntent | `offering.read`, `price.manage` |
| S30 | Disponibilidade | all-or-exact | Availability | List/GetSellableAvailability | `availability.read` |
| S40 | Performance / Resumo | exact-required | MarketplacePerformanceIntelligence | `GetMarketplacePerformanceSummary` | `performance.read` |
| S41 | Performance / Publicações | exact-required | MarketplacePerformanceIntelligence | `ListMarketplaceListingPerformance` | `performance.read` |
| S42 | Performance / Mídia | exact-required | MarketplacePerformanceIntelligence | `ListRetailMediaPerformance` | `performance.read` |
| S50 | Mercado | all-or-exact / subject-explicit | MarketIntelligence | List/GetCompetitivePosition, ListComparableOffers | `market.read` |
| S60 | Economia / Prevista | all-or-exact | CommercialEconomics | List/GetExpectedEconomics, EvaluatePriceScenario | `economics.read` |
| S61 | Economia / Realizada | all-or-exact | CommercialEconomics | List/GetSaleEconomics, GetEconomicPerformanceSummary | `economics.read` |
| S62 | Economia / Reconciliação | organization-wide | CommercialEconomics | List/Get/ResolveEconomicAttribution | `economics.read/reconcile` |
| S70 | Vendas | exact-required | MarketplaceSales | `ListMarketplaceSales` | `sales.read` |
| S71 | Venda / detalhe composto | exact subject | Sales + owners relacionados | `GetMarketplaceSale`, resolve attribution + reads owner-local | component Permissions |
| S72 | Pedidos ERP / materialização | organization-wide/contextual | BusinessSystemMaterialization | 5 operações BusinessOrder/Party/Destination | `materialization.read/resolve` |
| S73 | Faturamento | organization-wide/contextual | BusinessSystemMaterialization | List/GetInvoicingIntent | `materialization.read` |
| S80 | Expedição / Execuções | all-or-exact contextual | Fulfillment | `ListFulfillmentExecutions` | `fulfillment.read` |
| S81 | Execução de expedição | execution-explicit | Fulfillment | Get + 4 checkpoints + 2 artifact reads | `fulfillment.read/execute` |
| S82 | Envios | exact-required | Fulfillment | `ListShipments` | `fulfillment.read` |
| S83 | Envio / detalhe | exact-required | Fulfillment | `GetShipment` | `fulfillment.read` |
| S90 | Pós-venda | all-or-exact contextual | PostSaleResolution | List/CreatePostSaleResolution | `post_sale.read/manage` |
| S91 | Resolução pós-venda | ID + Sale qualificada | PostSaleResolution | `GetPostSaleResolution` | `post_sale.read` |
| S100 | Trabalho | organization-wide | OperationalWork | `ListWork` | `work.read` |
| S101 | Trabalho / detalhe | Work ID | OperationalWork | Get + 5 capabilities de coordenação | `work.read/manage` |
| S110 | Aprovações | organization-wide | ControlledActionGovernance | `ListAuthorizationDecisions` | `governance.read` |
| S111 | Decisão / contexto de aprovação | target/revision-explicit | ControlledActionGovernance | Get/CreateAuthorizationDecision | `governance.read/decide` |
| S120 | Configurações / Canais | organization-wide | MarketplacePortfolio | List/CreateMarketplaceInstallation + Technical Ingress | `portfolio.read/manage` |
| S121 | Configurações / Conta | exact Installation | MarketplacePortfolio | Get/Update/DeactivateMarketplaceInstallation | `portfolio.read/manage` |
| S122 | Configurações / Entidades vendedoras | organization-wide | MarketplacePortfolio | `ListSellingEntities` | `portfolio.read` |
| S123 | Configurações / Acesso | organization-wide | IdentityAccess | member/role lists + assign/revoke | `access.read/manage` |
| S124 | Configurações / Disponibilidade | organization-wide | Availability | 5 InventorySource + 2 policy ops | `availability.read/manage` |
| S125 | Configurações / Expedição | organization-wide | Fulfillment | 5 FulfillmentNode + 2 target ops | `fulfillment.read/manage` |
| S126 | Configurações / Política comercial | organization-wide | CommercialEconomics | Get/UpdateCommercialPolicy | `economics.read`, `economics.policy.manage` |
| S127 | Configurações / Delegações | organization-wide | ControlledActionGovernance | 4 AuthorizationDelegation ops | `governance.manage` |

**Total derivado:** 40 estados de tela/rota (`S00`–`S127` não contíguos) para 99 operações; número de estados não é meta de implementação.

## 6. Cobertura exata das 99 operações

| Authority / grupo | Quantidade | Home principal |
| --- | ---: | --- |
| Identity / access | 5 | S00, S123 |
| Marketplace Portfolio | 6 | S120–S122 |
| Product & Channel Readiness | 5 | S10 |
| Offering (Listing + ListingIntent + PriceIntent) | 12 | S20–S25 |
| Availability | 9 | S30, S124 |
| Market Intelligence | 3 | S50 |
| **Marketplace Performance Intelligence** | **4** | **S22, S40–S42** |
| Commercial Economics | 11 | S60–S62, S126 |
| Controlled Action Governance | 7 | S110–S111, S127 |
| Marketplace Sales | 3 | S70–S71 |
| Business-System Materialization | 7 | S72–S73 |
| Fulfillment + artifacts + Shipment | 17 | S80–S83, S125 |
| Post-Sale | 3 | S90–S91 |
| Operational Work | 7 | S100–S101 |
| **Total** | **99** | **100% mapeado** |

Nenhuma nova operação existe apenas para uma tela. `Strategy Workspace`, Overview e Sale detail permanecem composições client-side de Qs owner-native.

## 7. Leis de estado/segurança por interação

### 7.1 Conhecimento

A UI deve distinguir quando alcançável:

```text
complete / partial / unknown / unavailable / unsupported
known-zero / known-empty
stale / current quando a autoridade expõe isso
```

`partial` mostra a limitação do período/coverage; não recebe KPI estilizado como período completo.

### 7.2 Comparação de Performance

- presets de período são estado de navegação;
- requests enviam datas explícitas;
- frontend só mostra delta numérico quando o Product retorna `comparable`;
- `insufficient_evidence` e `not_comparable` são estados visíveis, não zero/delta inventado;
- provider measurement basis e custody histórica permanecem explicáveis.

### 7.3 Mutações consequenciais

- Idempotency-Key permanece estável somente para retry seguro da mesma intake;
- stale precondition não é business rejection;
- ambiguous external effect não ganha botão genérico de “tentar novamente”;
- hidden button não é autorização.

### 7.4 Retail Media

- `advertiser_id` nunca é identidade Product/Installation;
- campaign/catalog/family não vira Listing;
- technical binding não requer `performance.manage` inexistente;
- ausência de binding/provider access/contract admissibility aparece como conhecimento indisponível, não 403 se `performance.read` existe;
- Ads analysis não ganha write controls.

## 8. Falsificadores de frontend

O proof revisado deve tornar visivelmente inválido, entre outros:

1. Organization inferida de marketplace/account/browser;
2. seleção de canal usada como autoridade sem request explícito;
3. “Todos os canais” fundindo collections source-qualified independentes;
4. Amazon/Shopee exibidos como conectáveis hoje;
5. `ConnectMarketplace` Product inventado;
6. SourceInstance default/hardcoded;
7. Product master MPC inventado;
8. screen-shaped `/dashboard`, `/strategy`, `/analytics` ou `/metrics`;
9. `performance.read` implicando Market/Economics/Offering/Sales/Availability;
10. frontend recalculando provider CVR/ROAS como verdade canônica;
11. métrica de FAMILY/CATALOG/CAMPAIGN atribuída a Listing sem prova;
12. evidência parcial exibida como período completo;
13. known-zero confundido com unknown;
14. historical preserved evidence apresentada como MPC-authored source fact;
15. comparação incompatível produzindo delta;
16. Ads management/optimization controls aparecendo;
17. time-series/granularity ou opportunity score inventado sem contrato;
18. IA/MCP aparecendo como autoridade ou operação atual;
19. Market/Economics ganhando Price write;
20. `SetAvailableQuantity` ou Sync/Refresh Product inventado;
21. Sale detail ganhando workflow/write authority cross-owner;
22. comando direto Sankhya/Product materialization inventado;
23. Fulfillment físico confiando em qualificação declarada no cliente;
24. Work close resolvendo source truth;
25. Governance approval executando ação target;
26. route/button visibility virando autorização;
27. uma página paginada virando `total_count` global;
28. stale ETag tratado como erro genérico;
29. ambiguous external effect recebendo blind retry;
30. Settings tornando-se owner;
31. responsive mobile removendo estado/qualificação material;
32. provider DTO/status/AdGroup vazando como ontology Product.

## 9. Conjunto de wireframes representativos

O HTML low-fi revisado deve testar pelo menos:

1. Shell + Visão geral;
2. Preparação;
3. Editor de ListingIntent;
4. **Performance / Resumo**;
5. **Performance / Publicações + detalhe Performance do Listing**;
6. **Performance / Mídia**;
7. Disponibilidade;
8. Venda / composição cross-owner;
9. Expedição física;
10. Trabalho;
11. Configurações / Canais;
12. Economia / Reconciliação;
13. Configurações / Acesso e Aprovações.

Os valores são ilustrativos e o HTML é prova de hierarquia/estado, não runtime/browser/provider.

## 10. Próximo gate

Após substituir o HTML em português:

1. atacar os 32 falsificadores acima;
2. confirmar que nenhuma tela exige nova operação/Permission;
3. confirmar cobertura 99/99 e 30-Permission semantics;
4. submeter interaction map + wireframe revisados ao operador;
5. somente após aprovação, abrir adjudicação de frontend topology/dependencies;
6. não iniciar D7–D9 nem Product implementation.
