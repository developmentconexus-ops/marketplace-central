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
| `only_em_estoque` | estoque disponível > 0 (CODEMP IN (1,2): 1=CD, 2=loja física) | **true** | → **3.822** |
| `only_ecommerce_eligible` | `NVL(AD_ECOMMERCE,'X') <> 'N'` (live) / `IS NULL OR <> 'N'` (espelho) | **false** | corta só o que o ERP marcou `'N'`; sortimento cai 2.923 → ~1.329 |

Fatos que NÃO são graus de liberdade do chip:
- "Tipo de estoque" (`TIPCONTEST`) é pista falsa — tipo de CONTROLE, não vendabilidade. Fora.
- `AD_ECOMMERCE` entra como coluna e toggle desligado, não como filtro duro.
- Critério com dado ausente NÃO exclui (honest-unknown): linha xlsx sem `usoprod` passa o
  filtro de revenda. Ausência ≠ reprovação.
- F-LINK-1 refutado: espelho Sankhya casa 31/31 EANs da conta ML (exato, sem normalizar);
  xlsx casa 0. Nada de normalização de EAN neste chip.

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
4. Filtro aplica em DOIS caminhos de leitura, cada um no seu lugar:
   - catálogo live Sankhya: condições na query de página conforme a regra do tenant;
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

- ATIVO=S: 10.538 · ∧R: 10.007 · ∧estoque>0 (1,2): **3.822** · com EAN: 3.154 (82,5%).
