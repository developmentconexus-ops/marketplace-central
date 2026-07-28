Status: implementado e commitado em `4b76a28`.

Alterações apenas no `write_set`:

- Nova página/rota `/importacoes`.
- Histórico removido de `/vinculos` e mantido em `/integracoes`.
- Link “Importações” adicionado ao menu.
- Import corrigido após o `git mv`.
- Teste novo criado; teste existente preservado.

Validação:

- Vitest: `could-not-run (sandbox)` — erro de acesso do esbuild; não reexecutado.
- Typecheck: executado; somente os 15 erros preexistentes listados no cartão.
- `git diff --check`: executado com sucesso.
- Commit: executado com sucesso.

G1: `/importacoes` não depende de conta marketplace; os hooks chamam apenas `listErpImports()` e `getErpImport()`.

G2: escolhi o subtítulo “Histórico de importações do ERP por protocolo. Somente leitura.” e reutilizei `ImportacaoSection`, sem nova abstração.

G3: a futura rota `/importacoes/:importId` permanece livre para a próxima fatia.

Não verificado: suíte Vitest verde fora do sandbox, build e teste específico de wiring do `AppRouter` (fora do `write_set`).