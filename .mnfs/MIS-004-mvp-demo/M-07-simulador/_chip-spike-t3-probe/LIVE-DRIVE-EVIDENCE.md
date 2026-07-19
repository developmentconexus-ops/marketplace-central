# CHIP-SPIKE-T3 — Live drive evidence (hub, 2026-07-19)

Probe: `GET /integrations/installations/{id}/probes/catalog-match` @ main d4b5ef4a, backend container image 9126632, installation `inst-mercado_livre-d373dc64-...` (tenant_default), 8 produtos reais do xlsx ERP (import #003-E) + 1 re-teste com fallback `q`. Raw responses: `live-drive-results.jsonl`.

## Resultados

| Produto (codprod) | Entrada | catalog_hits | domain_discovery | fee_quote | Flags |
|---|---|---|---|---|---|
| 17304 CHAVE ALLEN JG (EAN real 7891504123035) | ean+price 699 | 1 (MLB27460123, MLB-WRENCHES) | 3 | ✅ MLB236055: Clássico 11% / Premium 14% | category_predita, buy_box_null |
| 3733 CHAVE FENDA VDE (EAN real 7891504918075) | ean+price 189 | 1 | 3 | ✅ | category_predita, buy_box_null |
| 21122 CHAVE COMBINADA 24MM (EAN real 4060833651202) | ean+price 49 | 2 | 3 | ✅ | category_predita, buy_box_null |
| 17733 FACA DESOSSA (EAN real 7891060995985) | ean+price 35 (SEM q) | 0 | 0 | ❌ null | no_catalog_hit — honesto |
| 17733 re-teste | ean+q+price | 0 | 1 (MLB-KITCHEN_KNIVES) | ✅ | q-fallback resolve categoria |
| 859 (EAN teste 1234567832999) | ean+price | 0 | 0 | ❌ null | no_catalog_hit — honesto |
| 554 (EAN teste 1234567829517) | ean+price | 0 | 0 | ❌ null | no_catalog_hit — honesto |
| 2002 "Invalid EAN Product" (sem EAN) | q+price | 0 | 0 | ❌ null | honesto (nome sem sentido) |
| 90004 PUXADOR FENG (sem EAN) | q+price 99 | 0 | 1 | ✅ predita | category_predita |

## Leitura do hub

- **Hit-rate catálogo com EAN real: 3/4 (75%).** EANs de teste/inválidos e nomes sem sentido → null honesto em tudo (ADR-17 correto, zero fabricação).
- **Fallback `q` (domain_discovery) fecha o gap:** com ean+q juntos, 4/4 EAN-reais resolvem categoria→fee. **Regra de design pro CHIP-T2: tier-3 SEMPRE envia ean E q(descrprod) na mesma sonda.**
- **Comissão real ≠ default: 11%/14% (Kit de Chaves @R$699) vs 13%/16% do tier-4.** Prova viva de que a escada importa — default superestima 2pp nessa categoria.
- **🔴 F1 CONFIRMADO LIVE: buy_box null em 9/9 chamadas** → categoria vem 100% de domain_discovery (predita). CHIP-T2 obrigado a investigar se buy_box_winner.category_id existe no payload real (path errado? permissão? campo top-level?) — 5ª candidata da classe campo-fabricado.
- **🟡 NOVO: mojibake no nome do catálogo** ("PeÃ§as" = UTF-8 double-decode) em catalog_hits[].name — bug de decode no catalog_match_reader; corrigir no CHIP-T2 (cosmético na sonda, real se nome for exibido/persistido).
- **Frete: fora do escopo da sonda** (sem dims no xlsx — product_enrichments manual). Segue plano design §2.

## Veredito

**GATE CHIP-T2 = GO.** Tier-3 COTAÇÃO-MATCH viável com ean+q combinados; comportamento honesto nos misses; comissão exata por categoria comprovada. Condições T2: ean+q sempre juntos, investigar buy_box category_id live (F1), fix mojibake, currency provider-sourced se sair de MLB (F4).
