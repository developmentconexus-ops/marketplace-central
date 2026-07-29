# CHIP-VENDAVEL — regra de sortimento vendável (context pack)

Base de despacho: `main @ cf741f68`. Chip interino pós-MIS-006, ordenado pelo operador antes do
planning da MIS-007. Track único — nenhum chip concorrente em voo.

## Objetivo

O operador abre `/catalogo` e vê 10.538 produtos, a maioria morta para venda. A regra de
"vendável" vira configuração por tenant com efeito em todas as leituras relevantes:

- Card **"Sortimento vendável"** em `/integracoes` (ao lado do card de fonte ativa): 3 toggles
  persistidos no banco + linha "Resultado: N de M produtos" computada ao vivo.
- `/catalogo` abre filtrado com chip contador "Vendáveis N de M" + botão "ver todos" (escape de
  tela, não muda config). Produto do sortimento com estoque zero = badge, nunca corte extra.
- Geração de vínculos considera só vendáveis.

## Regra ratificada (operador, 2026-07-29 — medida na base viva METALPRD)

| toggle | semântica | default | corte medido |
|---|---|---|---|
| `only_revenda` | `TGFPRO.USOPROD = 'R'` | **true** | 10.538 → 10.007 |
| `only_em_estoque` | estoque disponível > 0 (ver EMENDA A7 abaixo para o recorte) | **true** | ver A7 |
| `only_ecommerce` | `TGFPRO.AD_ECOMMERCE = 'S'` | **false** | flag rala hoje (606); cliente vai passar a manter |

Fatos que NÃO são graus de liberdade do chip:
- "Tipo de estoque" (`TIPCONTEST`) é pista falsa — tipo de CONTROLE, não vendabilidade. Fora.
- `AD_ECOMMERCE` entra como coluna e toggle desligado, não como filtro duro.
- Critério com dado ausente NÃO exclui (honest-unknown): linha xlsx sem `usoprod` passa o
  filtro de revenda. Ausência ≠ reprovação.
- F-LINK-1 refutado: espelho Sankhya casa 31/31 EANs da conta ML (exato, sem normalizar);
  xlsx casa 0. Nada de normalização de EAN neste chip.

## EMENDA A7 — recorte de localização de estoque (ratificada pelo operador 2026-07-29, pós-despacho)

Verbatim do operador, dito na sessão do chip: *"CODLOCAL é legal: 10101 é estoque revenda, 10108 é
show room não é pra contar, e 10102 é outlet — esses que vendem"*. Perguntado se valia só para o
filtro ou também para o número em tela: **os dois**.

**Disponível vendável = `CODEMP IN (1,2) AND CODLOCAL IN (10101, 10102) AND CODPARC = 0`.**
10108 (show room) fora. Consequências fixadas:
- O pin do Q4 do sync é COMPLETO (empresa **e** local); o must-fail falha se qualquer um dos dois sair.
- `catalog_page.go` hoje usa `CODLOCAL = 10101` → passa a `IN (10101, 10102)` (outlet entra; o
  operador quer essa mudança de número).
- Espelho, catálogo vivo e contador N/M usam a MESMA definição de disponível — divergência = VC-2 falha.
- Lista de CODLOCAL é constante comentada no código; nada de UI de escolha de local (YAGNI).

**O número 3.822 do pack original está MORTO** — foi medido sem o recorte de CODLOCAL. Até a medição
viva voltar do especialista Oracle, o VC-2 se descarrega por concordância (tela == SQL da mesma regra
no mesmo banco), nunca por constante.

### A7 fechada — db-consult respondido e medido ao vivo (METALPRD 2026-07-29, READ-ONLY)

Predicado RATIFICADO, autoridade final do chip:

```sql
p.ATIVO='S' AND p.CODPROD>0 AND p.USOPROD='R'
AND EXISTS (SELECT 1 FROM METALPRD.TGFEST e
  WHERE e.CODPROD=p.CODPROD AND e.CODPARC=0
    AND e.CODEMP IN (1,2) AND e.CODLOCAL IN (10101,10102)
  GROUP BY e.CODPROD HAVING SUM(NVL(e.ESTOQUE,0)-NVL(e.RESERVADO,0))>0)
```

- **Subtrair RESERVADO: SIM.** 871 linhas reservadas / 37.428 un no escopo; **354 produtos zeram só por
  reserva**. "Tem estoque pra vender" = disponível.
- **`CODPARC = 0`: correto e defensivo.** No recorte é o único CODPARC com estoque (4.091 linhas,
  109.494 un); em outras empresas/locais existe estoque de terceiro, então o predicado fica.
- **Número de aceitação: 2.923 produtos vendáveis.** Sem o recorte de CODLOCAL seriam 3.508 — os
  **585 de diferença** entram indevidos hoje (show room, almoxarifado, quebras…), que é exatamente o
  defeito do sync. Decomposição: 10101 = 2.795 · 10102 = 187 (soma > total por sobreposição).
- **Conciliação com o número morto**: 3.822 usava ESTOQUE bruto, sem reserva e sem local. Mesmo escopo
  bruto com local = 3.277; a reserva tira 354 → 2.923.

**CODLOCAL com disponível > 0 nas empresas 1|2 (para o comentário-constante não mentir):**

| CODLOCAL | TGFLOC.DESCRLOCAL | produtos | vende |
|---|---|---|---|
| 10101 | 1_REVENDA | 8.310 | SIM |
| 10102 | 2_OUTLET | 595 | SIM |
| 10108 | 8_SHOW ROMM [sic no ERP] | 1.574 | não — mostruário |
| 10106 | 6_PENDENCIA PRODUTO FORNECEDOR | 673 | não |
| 10107 | 7_QUEBRAS EM PALETES | 553 | não |
| 10105 | 5_ALMOXARIFADO | 273 | não |
| 10103 | 3_AGUARDANDO CONSERTO | 113 | não |
| 10109 | 9_USO E CONSUMO | 70 | não |

**Whitelist, nunca blacklist** — local interno novo entraria vendendo sozinho se a regra fosse por
exclusão. O código carrega a whitelist como constante comentada.

## Defeito incluído no escopo (achado na investigação)

O Q4 de estoque do sync soma TODAS as empresas — `WHERE CODPARC = 0` sem CODEMP em
[sync.go:252-256](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go).
`estoque_total` do espelho está inflado com 901/501/904. **Pin `AND CODEMP IN (1, 2)`.**
(O catálogo live já usa `est.CODEMP IN (1, 2)` em
[catalog_page.go:135](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go) — o pin torna espelho e catálogo consistentes.)

## Decisões de desenho FIXADAS (não rediscutir)

1. **Config = 3 colunas boolean na infraestrutura `tenant_config` existente** (mesmo módulo e
   handler do `active_source` — [active_source.go](../../../apps/server_core/internal/modules/tenant_config/active_source.go)). Sem tabela nova, sem JSONB, sem módulo novo.
2. **Espelho ganha 2 colunas aditivas**: `usoprod`, `ad_ecommerce` (`products_mirror`,
   nullable). EAN já existe. Migração no bloco **0083–0084** (pré-alocado; precisar de mais =
   `REQUEST migration-number`).
3. Sync Sankhya Q1 ([sync.go:101-110](../../../apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go)) passa a ler `p.USOPROD, p.AD_ECOMMERCE`. Adapter xlsx: colunas
   OPCIONAIS no parser (planilha sem elas continua importando; ficam NULL).
4. Filtro aplica em TRÊS caminhos de leitura (**DR-1**, emendado 2026-07-29 — o pack original
   listava dois e esquecia o catálogo servido pelo espelho):
   - catálogo live Sankhya: condições na query de página conforme a regra do tenant;
   - catálogo servido pelo espelho (`routing.Reader` → leitor de upload → `MirrorCatalogPage`):
     MESMO predicado, MESMA contagem viva. A regra não pode depender da fonte ativa;
   - geração de vínculos: `MirrorMatcher` ([root.go:490](../../../apps/server_core/internal/composition/root.go)) filtra `products_mirror` pela regra.
5. Contador N/M: endpoint leve computado da fonte que a tela usa (não constante, não cache
   manual). M = população sem regra; N = com regra.
6. FE: card em [IntegracoesPage.tsx](../../../apps/web/src/pages/integracoes/IntegracoesPage.tsx); catálogo em `apps/web/src/pages/` (rota `/catalogo`); SDK client no bloco do domínio.

## Fora de escopo (YAGNI — corte do operador)

- Corte por `AD_DTULTVEND` ("sem venda há N meses") — futuro, colunas não entram agora.
- UI de escolha de empresa/armazém — `CODEMP IN (1,2)` é constante comentada no código.
- Qualquer coisa em análise de mercado / MIS-008.
- `TIPCONTEST` no espelho.

## Números de aceitação (base METALPRD 2026-07-29; tolerância = drift diário ~5 produtos)

- ATIVO=S: 10.538 · ∧R: 10.007.
- ∧estoque>0 em CODEMP (1,2) **sem recorte de CODLOCAL**: 3.822 (com EAN 3.154 = 82,5%) —
  número HISTÓRICO, superado pela EMENDA A7. Não usar como constante de aceitação.
- ∧disponível>0 em CODEMP (1,2) **∧ CODLOCAL IN (10101,10102)**, disponível = ESTOQUE − RESERVADO:
  **2.923** ← número de aceitação vivo do VC-2 (tolerância = drift diário).
