# CHIP-ANCHORS — rung LIVE (U1-U3), hub-run por R-8

Rodado **antes do merge**, contra o código do chip, sem copiar credencial nenhuma para a worktree:
o compose sobe do checkout primário (onde `.env` resolve) com um override que troca **apenas** a
origem do bind `/workspace` para `.claude/worktrees/chip-anchors`. Entradas de volume fundem por
caminho de container, então o resto da stack fica igual e a árvore do chip não é tocada.

```
docker compose -f docker-compose.yml -f <scratchpad>/compose.chip-anchors.yml up -d backend frontend
docker exec marketplace-central-backend-1 ls .../product_links/adapters/connectors
  identity_anchor_adapter.go      ← confirma que o container serve o código do chip
```

Conta ML conectada, installation `inst-mercado_livre-7e0d2125-…`, 34 anúncios reais.
Candidatos **regenerados** via `POST /product-links/link-candidates/generations` →
`generated_count: 38` (R3: nada de backfill; motivo novo só nasce na geração nova).

## U1 — motivo `marca`/`refforn` deixa de ser suposição do núcleo — **PASS**

Antes (candidato persistido às 11:25, código pré-chip):

```
{"anchor":"marca","direction":"UNAVAILABLE","detail":"marca inexistente no lado provider"}
```

Depois da regeneração (19:19, código do chip), na tela em `/vinculos`, motivo expandido:

```
marca: provider não fornece a âncora marca
refforn: provider não fornece a âncora refforn
```

A frase agora **deriva da declaração**: `capability_adapter.go:90` declara
`[seller_sku, ean, title]`, e o adapter do `product_links` marca como não-fornecida toda âncora do
vocabulário que a declaração não contém. Provider 2 que declarar `marca` vira `Supplied` e a linha
desaparece sozinha — que era o ponto do chip.

## U2 — toalheiro deixa de ser reprovado por contradição de dimensão — **PASS**

Anúncio `MLB4735326915` "Toalheiro Simples Soul Zen 50cm Cromado" × produto ERP `sku:33698`
"SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO":

| | estado | banda | conf. | motivo dominante |
|---|---|---|---|---|
| antes | rejeitado | BAIXA | 25% | `Título hard-negative: medidas` |
| agora | `exact_sku` CONFIRM | MEDIA | 70% | `✓ SKU — seller_sku resolve exato para codprod` |

`50cm` (anúncio) e `500MM` (ERP) pararam de ler como contradição. O guard continua bloqueante —
as linhas de EAN ambíguo do mesmo print seguem `BAIXA 20%` com `âncora ambígua`.

## U3 — `/vinculos` → Resolvidos mostra os 29, KPI concorda com o banco — **PASS**

```
psql: select state, count(*) from product_links group by state;  →  resolved | 29
GET /product-links/link-workflows?installation_id=…   (SEM limit, default)
  items = 34,  itens com current_link = 29
DOM da aba Resolvidos: document.querySelectorAll('table tbody tr').length  →  29
KPI "Resolvidos hoje" → 29
```

O default compartilhado (20) truncava 9 vínculos. Três fontes independentes concordam em 29:
banco, resposta da API sem `limit`, e linhas renderizadas.

## Veredicto

**U1-U3 PASS.** A rung LIVE está satisfeita; o chip pode escrever a linha `LIVE-VERIFIED:` no pack
citando este arquivo (o marcador tem de estar no pack porque é lá que o merge-gate lê).

## Achados do hub — nenhum bloqueia o merge, nenhum é do chip

1. **`refforn` no vocabulário de âncoras não carrega informação.** `refforn` é a referência do
   fornecedor no NOSSO ERP; provider nenhum vai declarar que fornece isso. O resultado é uma linha
   `UNAVAILABLE` permanente em toda linha de todo provider, para sempre. `marca` é diferente — é
   conceito dos dois lados e o ML tem atributo de marca, ele só não chega no nosso snapshot. Ou
   seja: o vocabulário mistura âncora **comparável entre os lados** com campo **interno**. Candidato
   a corretivo depois do merge, não agora.
2. **Âncora declarada que o núcleo não sabe comparar some em silêncio**
   (`generation_service.go:639`, `if anchor.Supplied { continue }`). Um `mpn`/`gtin14` do
   marketplace 2 não vira motivo nenhum. Não é mentira, é omissão — primo do ADR-17. Já registrado
   no aceite do C2.
3. **FE, dono é M-06**: a coluna "SKU ML" renderiza `provider_code` (`mercado_livre`) em toda linha
   — não é SKU nenhum. O honesto é `match_value` quando `match_input == "seller_sku"`, senão "—".
   Já está na carta do F-05.
