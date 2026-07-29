# CHIP-VENDAVEL — validation contract

Verdito só por QA live-drive (browser, persona fresca) + P6 dual gate. curl-only = FAIL.
Toda contagem lida na tela é confrontada com SQL no mesmo banco antes de virar PASS.

| ID | critério | evidência mínima |
|---|---|---|
| VC-1 | Card "Sortimento vendável" em `/integracoes` com 3 toggles; estado persiste no BANCO por tenant: alterar toggle, limpar localStorage+sessionStorage, recarregar → estado mantido e `localStorage.length === 0` | screenshot/DOM + SQL da linha em tenant_config antes/depois |
| VC-2 | Linha "Resultado: N de M" do card e chip "Vendáveis N de M" do catálogo batem com SQL rodando a MESMA regra no mesmo banco. Descarrega por **concordância** (tela == SQL da mesma regra no mesmo banco) E bate com o número vivo ratificado: **N ≈ 2.923** (medido em METALPRD 2026-07-29 com o predicado A7 fechado; tolerância = drift diário). O 3.822 do pack original está morto — media ESTOQUE bruto, sem reserva e sem recorte de local | contagem DOM + query e resultado colados |
| VC-3 | `/catalogo` abre filtrado por default; "ver todos" mostra a população cheia; produto vendável com estoque 0 aparece com badge quando `only_em_estoque` desligado — nunca some do sortimento por badge | DOM antes/depois do toggle de tela |
| VC-4 | Geração de vínculos ignora produto não-vendável. **Must-fail nomeado**: com a regra ativa, candidato de produto USOPROD≠'R' NÃO nasce; desligando `only_revenda` no banco e regerando, PODE nascer — o mesmo teste falha se o filtro for removido do `MirrorMatcher`. **F-1 (2026-07-29, critério do S7)**: colisão de EAN conta PÓS-corte — gêmeo cortado pela regra NÃO marca o sobrevivente como ambíguo (hoje `MirrorEANCollisionCounts` agrega o espelho inteiro, pré-política, enquanto o mecanismo `len(results)>1` já é pós-corte: os dois têm que concordar). Must-fail: com regra ativa, sobrevivente de par cortado auto-aprova por D-121; sem regra, volta a ambíguo | teste de integração + failure_token=test= da lane |
| VC-5 | Q4 do sync pinado em `CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)` (EMENDA A7): must-fail que falha se **qualquer um dos dois** predicados sair da query — pin parcial não conta; `catalog_page.go` migra de `CODLOCAL = 10101` para o `IN`; `usoprod`/`ad_ecommerce` populados no espelho após sync (contagens > 0 e coerentes com Q1) | teste + SQL do espelho pós-sync (sync live = passo hub pós-merge) |
| VC-6 | Import xlsx SEM as colunas novas continua aceitando (colunas opcionais, ficam NULL e passam o filtro por honest-unknown); xlsx COM as colunas as popula | teste com as duas planilhas |
| VC-7 | Zero erro de console e zero request 4xx/5xx nas telas tocadas (`/integracoes`, `/catalogo`, `/vinculos`); `tsc` verde na raiz; vitest = DUAS suítes EXPLÍCITAS, contadas por linha: workspace `@marketplace-central/web` E `packages/sdk-runtime` (rodada do próprio diretório do pacote — o script `test` da raiz só roda web, a suíte do sdk-runtime NÃO entra nele; `No test files found`/exit 1 = filtro que não casou, nunca verde). Conserto do script da raiz = passo hub PÓS-merge (seam compartilhada, nenhuma fatia é dona); ladder Go verde | logs da varredura + saídas das duas lanes vitest |

Regras herdadas (doutrina, não repetir errado):
- Verde de integração só é evidência depois do vermelho NOMEAR o teste (must-fail primeiro).
- Asserção de presença não pega valor errado — assere o VALOR esperado, não `!= ""`.
- Varredura se valida no CHAMADOR.
- String de tela em pt-BR do operador; nenhum marcador de dev tipo `(unproved)` em copy.
