Status: implementado e commitado.

Commit: `1e4c297 fix(web): harden import chain reads`

Alterações no `write_set`:

- Guard de tipo para `queue_read_at`, reutilizando `formatDateTime`.
- 404 só é “não encontrada” com `error: "import_not_found"`.
- Testes para data inválida, zero conhecido e 404 genérico.

Validação:

- `tsc`: executou; apenas erros preexistentes, nenhum nos arquivos alterados.
- Vitest: `could-not-run (sandbox)` — bloqueio de acesso ao `vitest.config.ts`/esbuild.
- `git diff --check`: passou.
- Commit: executado com sucesso.

G1: o guard fica no consumidor; depois dele, `formatDateTime` continua sendo usado. A função pertence ao pacote compartilhado `packages/web-query` e também é consumida por `ImportacaoSection`.

G2: mantive valores numéricos reais sem impor regra adicional de inteiro/não-negativo.

G3: a condição 3 depende do contrato hook/SDK, mas não altera arquivos nem seams de outro chip.

Os três arquivos `.mnfs/.../dispatches/*` já estavam não rastreados e foram preservados, sem alteração.