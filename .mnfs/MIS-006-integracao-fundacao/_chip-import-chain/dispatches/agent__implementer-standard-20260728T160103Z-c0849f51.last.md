Status: implementado e commitado em `1bffdcf`.

Alterações: todos os 9 caminhos estão dentro do `write_set`; nenhum caminho não declarado foi alterado. O teste de `ImportacaoSection` foi envolvido em `MemoryRouter` e recebeu uma asserção para `/importacoes/imp_1`.

Validação:

- Typecheck: executado; somente erros pré-existentes, nenhum nos arquivos alterados.
- `git diff --check`: passou.
- Vitest: `could-not-run (sandbox)` — esbuild retornou `Access is denied`.
- Commit: executado com sucesso.

G1: hook page-local é correto; `web-query` mantém utilitários/chaves compartilhados, enquanto hooks ERP existentes já seguem em `pages/vinculos`.

G2: usei chave local para a cadeia e mensagem genérica explícita para falhas não-404; 4xx não fazem retry.

G3: nenhuma seam nomeada foi bloqueada.

Não verificado localmente: execução dos testes Vitest e `ImportacoesPage.test.tsx`; o chip deve executar a validação fora do sandbox.