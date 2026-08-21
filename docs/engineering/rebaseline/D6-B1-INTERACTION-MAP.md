# D6-B1 — Mapa de Interações do Frontend

> **Status:** OPERATOR-RATIFIED — mapa corrigido após D6-R1 / OAD 99 operações · 30 Permissions
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
  → selecionar Sale source-qualified
  → GetMarketplaceSale
  → GetSaleEconomics quando permitido
  → ListBusinessOrderIntents / GetBusinessOrderIntent
  → CreateBusinessOrderIntent quando requerido pelo owner
  → observar PartyResolution e DestinationRealization
  → resolver somente pelas operações Product admitidas
  → observar InvoicingIntent / FulfillmentExecution sem inventar sync/manual close
```

A tela pode compor Sales, Economics, Materialization e Fulfillment; nenhum deles transfere sua write authority para a tela ou para outro owner.

### F9 — Executar expedição preservando autoridade física

```text
Expedição
  → List/GetFulfillmentExecution
  → List/GetFulfillmentNode
  → RecordSeparation
  → RecordPhysicalConference
  → RecordPacking
  → RecordDispatchHandoff
  → List/Get artifacts
  → List/Get Shipment quando houver Installation
```

Botão visível nunca prova capacidade física. O servidor revalida Principal kind, Permission, qualificação física quando requerida e revisão atual.

### F10 — Operar Pós-venda

```text
Pós-venda
  → List/GetPostSaleResolution
  → CreatePostSaleResolution
  → UpdatePostSaleResolution
```

Não existe ação genérica de fechar resolução; lifecycle permanece conforme Product authority.

### F11 — Operar Trabalho

```text
Trabalho
  → List/GetWork
  → AssignWork
  → ClearWorkAssignment
  → HoldWork
  → ResumeWork
  → EscalateWork
```

Não existe `CloseWork`. O frontend não cria um workflow engine paralelo.

### F12 — Governança / Aprovações

```text
Aprovações
  → List/GetAuthorizationDecision
  → CreateAuthorizationDecision
  → DecideAuthorizationDecision
  → List/GetAuthorizationDelegation
  → Establish/Update/RevokeAuthorizationDelegation
```

Uma decisão autoriza segundo Governance; ela não executa automaticamente a operação alvo nem concede a Permission ausente.

## 5. Estados de tela essenciais

A implementação visual pode variar; estes estados semânticos não.

| Tela/área | Estado essencial |
| --- | --- |
| Visão Geral | loading / empty / partial / ready / access-limited |
| Preparação | source selection / readiness incomplete / ready / provider unavailable |
| Publicações | draft / submitted / accepted / pending / rejected / ambiguous |
| Disponibilidade | known / unknown / unavailable / policy/configuration state |
| Performance Resumo | complete / partial / unknown / unavailable / unsupported; comparable / insufficient / not-comparable |
| Performance Publicações | known Listing population even when performance evidence unavailable |
| Performance Mídia | available population distinct from unavailable/unsupported access/evidence |
| Mercado | comparable position present / unavailable / insufficient comparability |
| Economia | expected / actual / attribution unresolved / reconciled where Product admits |
| Vendas | sale observed / attribution unresolved / materialization pending / fulfillment state |
| Expedição | physical checkpoint progression + qualification/precondition failures |
| Pós-venda | active Product lifecycle only; no client-created close state |
| Trabalho | assignment/hold/escalation states only as Product exposes |
| Aprovações | pending / decided / expired/inapplicable only where Product exposes |
| Configurações | permission-limited / provider capability unavailable / technical ingress required |

## 6. Cobertura por grupo de operações

| Grupo | Operações cobertas |
| --- | ---: |
| Identity / access | 5 |
| Portfolio | 6 |
| Readiness | 5 |
| Offering | 12 |
| Availability | 9 |
| Market | 3 |
| Performance | 4 |
| Economics | 11 |
| Governance | 7 |
| Sales | 3 |
| Materialization | 7 |
| Fulfillment + artifacts + Shipment | 17 |
| Post-Sale | 3 |
| Work | 7 |
| **Total** | **99** |

## 7. Falsificadores frontend

O desenho para e volta ao owner apropriado se qualquer implementação exigir um destes atalhos:

1. inferir Organization atual sem `GetCurrentAccessContext`/contexto explícito;
2. somar métricas entre Installations/providers como equivalentes sem prova;
3. mostrar Amazon/Shopee como conectáveis sem D4 evidence atual;
4. criar Dashboard/Strategy/Analytics como Product authority;
5. tratar menu oculto como autorização;
6. reconstruir CVR/ROAS/provider metrics no browser;
7. colapsar campaign/family/catalog/listing scope em uma identidade genérica;
8. transformar `partial` em completo ou unknown/unavailable em zero;
9. chamar evidência provider-preserved de MPC-authored fact;
10. permitir Ads writes/optimization sob Performance;
11. criar time-series/granularity/Metric DSL não admitido;
12. criar `signals[]`, recommendations, opportunity score ou AI explanation como Product truth;
13. mover Price write para Market/Economics/Performance;
14. criar direct Sankhya/Oracle access no frontend;
15. usar Work como generic workflow engine;
16. fazer Governance Decision executar a operação alvo;
17. tratar visibilidade de botão como Permission/physical qualification;
18. inferir total de uma coleção paginada sem total contractual;
19. tratar `412`/ETag stale como erro genérico e sobrescrever silenciosamente;
20. retry automático de efeito potencialmente aceito/ambíguo;
21. transformar Settings em owner;
22. apagar estado semântico no mobile/responsive;
23. vazar DTO/provider endpoint/advertiser id como Product identity;
24. assumir source instance default em Readiness;
25. inferir Marketplace Installation pelo primeiro disponível;
26. gerar publicação/listing a partir de dados incompletos sem owner semantics;
27. transformar OAuth callback/provider connect em Product operation;
28. permitir A/S usar ação H-only apenas porque possui Permission;
29. atribuir physical fact a S sem qualificação current server-owned;
30. atribuir Retail Media family/catalog/campaign a Listing sem evidence exata;
31. exibir delta quando comparação é `insufficient_evidence`/`not_comparable`;
32. esconder `unavailable/unsupported` como coleção vazia normal.

## 8. Resultado

D6-B1 prova que as 99 operações atuais possuem homes de interação coerentes sem criar telas/owners artificiais. O App Shell, o mapa de interação corrigido e os low-fi wireframes foram operator-ratified após revisão independente e bounded fixes.

Qualquer implementação deve continuar preservando as leis acima; esta ratificação não seleciona runtime, router de servidor, banco, worker, deployment ou implementação Product.
