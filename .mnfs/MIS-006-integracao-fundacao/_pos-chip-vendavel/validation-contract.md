# CHIP-VENDAVEL — validation contract

Verdito só por QA live-drive (browser, persona fresca) + P6 dual gate. curl-only = FAIL.
Toda contagem lida na tela é confrontada com SQL no mesmo banco antes de virar PASS.

| ID | critério | evidência mínima |
|---|---|---|
| VC-1 | Card "Sortimento vendável" em `/integracoes` com 3 toggles; estado persiste no BANCO por tenant: alterar toggle, limpar localStorage+sessionStorage, recarregar → estado mantido e `localStorage.length === 0` | screenshot/DOM + SQL da linha em tenant_config antes/depois |
| VC-2 | Linha "Resultado: N de M" do card e chip "Vendáveis N de M" do catálogo batem com SQL rodando a MESMA regra no mesmo banco. O alvo é a CONCORDÂNCIA tela↔SQL na fonte que a tela está lendo, não um número fixo — o número de referência não vive aqui (ver nota abaixo) | contagem DOM + query e resultado colados |
| VC-3 | `/catalogo` abre filtrado por default; "ver todos" mostra a população cheia; produto vendável com estoque 0 aparece com badge quando `only_em_estoque` desligado — nunca some do sortimento por badge | DOM antes/depois do toggle de tela |
| VC-4 | Geração de vínculos ignora produto não-vendável. **Must-fail nomeado**: com a regra ativa, candidato de produto USOPROD≠'R' NÃO nasce; desligando `only_revenda` no banco e regerando, PODE nascer — o mesmo teste falha se o filtro for removido do `MirrorMatcher` | teste de integração + failure_token=test= da lane |
| VC-5 | Q4 do sync pinado em `CODEMP IN (1,2)`: teste que falha se o pin sair da query; `usoprod`/`ad_ecommerce` populados no espelho após sync (contagens > 0 e coerentes com Q1) | teste + SQL do espelho pós-sync (sync live = passo hub pós-merge) |
| VC-6 | Import xlsx SEM as colunas novas continua aceitando (colunas opcionais, ficam NULL e passam o filtro por honest-unknown); xlsx COM as colunas as popula | teste com as duas planilhas |
| VC-7 | Zero erro de console e zero request 4xx/5xx nas telas tocadas (`/integracoes`, `/catalogo`, `/vinculos`); `tsc` e vitest verdes na raiz do worktree; ladder Go verde | logs da varredura + saídas das lanes |

Nota do VC-2 — onde vive o número (F3, ratificado pelo hub em A-26 ponto 2):

Este arquivo carregava `N≈3.822`. Aquele número foi medido SEM o corte de `CODLOCAL` e SEM
reserva, e foi superado antes de qualquer tela existir. `ADDENDUM-01-codlocal.md:53` e A-3
(`BATCH-PLAN.md:1133`) ratificaram **2.923**, e `BATCH-PLAN.md:1188` restringe: 2.923 é o
número do caminho **Oracle live apenas** — o caminho do espelho tem contagem própria, porque
o espelho só carrega o que o último sync trouxe. Um alvo fixo neste contrato reprovaria o
espelho por estar certo. Prosa falsa em contrato se deleta; a emenda já existe e é a fonte.

Regras herdadas (doutrina, não repetir errado):
- Verde de integração só é evidência depois do vermelho NOMEAR o teste (must-fail primeiro).
- Asserção de presença não pega valor errado — assere o VALOR esperado, não `!= ""`.
- Varredura se valida no CHAMADOR.
- String de tela em pt-BR do operador; nenhum marcador de dev tipo `(unproved)` em copy.
