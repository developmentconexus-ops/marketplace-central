# F-A1 / F-A2 — falha de token ML visível: pack de evidência

Fatia da Onda 0 da MIS-008. Plano: `docs/superpowers/plans/2026-08-03-fa1-fa2-token-visivel-plan.md`.
Data: 2026-08-03. Branch: `main`.

## Estado: 6 de 7 tasks fechadas. Task 7 (live drive) BLOQUEADA — decisão do operador.

| Task | Commit | O que prova |
|---|---|---|
| 1 — classificação do erro do ML | `8864a37a` | `invalid_grant` vira `ErrRefreshTokenInvalid`, que `ClassifyRefreshError` trata como terminal |
| 2 — persiste a falha | `c82f2e45` | `refresh_failure_code`, `consecutive_failures`, `next_retry_at` são gravados |
| 3 — degrada a installation | `b745d1ba` | terminal → `requires_reauth`/`critical`; acima do limite → `degraded`/`warning` |
| 4 — ticker loga e continua | `29aaa3ed` | falha de item não aborta o lote; erro deixa de ser descartado |
| 5 — prova contra Postgres real | `f9f8f001` | o `WHERE next_retry_at` existe de fato; primeira escrita de `critical` no schema |
| 6 — card em `/integracoes` | `d3a2845c` | estado da conexão chega à tela, com ramo de desconhecido (ADR-17) |
| 6b — asserção restaurada | `aafd014a` | invariante pré-clique reescrita contra observável exclusivo |
| 7 — live drive | — | **bloqueada, ver abaixo** |

## Zero mudança de contrato

`IntegrationConnectionSnapshot` e `IntegrationInstallation.health_status` já existiam no
OpenAPI e no `sdk-runtime`. Nenhum arquivo de contrato foi tocado pela fatia.

## Controles negativos e must-fail observados

Verde sozinho não conta. Cada prova abaixo foi vista **vermelha primeiro**, nomeada.

1. **Task 1** — `anything else stays provider error`: controle negativo, garante que a
   classificação nova não engole todo erro >= 400.
2. **Task 3** — `SingleTransientFailureDoesNotDegrade`: controle negativo, uma falha
   transitória isolada não pode degradar a instalação.
3. **Task 4** — `RunOnceReturnsErrorWhenListingFails`: controle negativo, falha do
   **lote** continua sendo erro; só falha de **item** é engolida.
4. **Task 4 — must-fail por injeção.** O agente que implementou a task não observou o
   vermelho comportamental (implementou e testou junto). Reinjetei o abort antigo em
   `RunOnce` e observei:

   ```
   --- FAIL: TestRunOnceContinuesAfterItemFailure (0.00s)
       refresh_ticker_test.go:72: RunOnce err = boom, want nil (falha de item não é falha do lote)
   FAIL	marketplace-central/apps/server_core/internal/modules/integrations/background	1.275s
   ```

   Restaurado; `git status --porcelain` limpo e pacote verde.

5. **Task 5 — must-fail por injeção.** Cláusula de backoff em `auth_session_repo.go:87`
   substituída por comentário. A lane passou a reportar:

   ```
   failure_token=package=marketplace-central/apps/server_core/tests/integration
   failure_token=test=TestListingsReadContractEndToEnd
   failure_token=test=TestRefreshFailurePersistsAndSuppressesTheSweep
   ```

   Com a cláusula restaurada, o token do meu teste **desaparece** e resta só o pré-existente.

6. **Task 5 — controle positivo interno.** Antes de afirmar que a varredura pula a sessão
   com `next_retry_at` no futuro, o teste afirma que ela **pega** a mesma sessão sem
   backoff. Sem esse passo, lista vazia por tenant errado ou expiry errado passaria por
   backoff funcionando.

7. **Task 6b — must-fail por injeção.** A asserção restaurada foi validada com
   `fireEvent.click(connect)` injetado antes dela:

   ```
   AssertionError: expected "spy" to not be called at all, but actually been called 1 times
   ```

   A versão **síncrona** da mesma asserção passou com o clique injetado — foi por isso que
   ela virou `async` com flush antes do `expect`. Asserção síncrona sobre chamada que só
   ocorre depois de um `await` não falsifica nada.

## Asserção apagada de teste alheio — adjudicada

A Task 6 apagou `expect(listIntegrationInstallations).not.toHaveBeenCalled()` de
`IntegracoesPage.test.tsx` porque o `InstallationProvider` novo busca a lista no mount.
Classe registrada: restaure OU prove por observável.

A justificativa dada ("coberto por botão habilitado e label 'Conectar'") **não fecha** —
esses dois são o estado de mount e continuam verdadeiros mesmo com um `connect()` já em
voo. Restaurada contra os dois observáveis que `connect()` (`IntegracoesPage.tsx:503`)
produz e nada mais na árvore produz: `createIntegrationInstallation` (`:522`) e
`startIntegrationAuthorization` (`:529`). Ambos declarados no mock de `useClient` — antes
não existiam lá, então uma asserção sobre eles teria sido vácuo. Commit `aafd014a`.

## Task 7 — por que está bloqueada

Duas razões, ambas medidas.

**(1) O binário em execução é anterior à fatia.**

```
marketplace-central-backend-1	Up 27 hours	2026-08-02 10:23:30 -0300
8864a37a1ff05eff12fb066feed0c3b817a33202 2026-08-03 12:44:25 -0300
```

O container não contém nenhuma das quatro mudanças de servidor. Contra ele todo observável
dos Steps 2–5 é impossível **por construção**, e um live drive que "não mostrou nada" seria
lido como ausência de defeito. Rebuild do backend é ação de stack: `REQUEST` ao hub.

**(2) Como induzir o `invalid_grant` real.**

Estado medido do dev stack:

```
 installation_id                                         | provider_code | status    | health_status
 inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0 | mercado_livre | connected | healthy

 state | access_token_expires_at        | next_retry_at | refresh_failure_code | consecutive_failures
 valid | 2026-08-03 21:48:54.055485+00  |               |                      | 0
```

Forçar `access_token_expires_at` para o passado exercita só o caminho de **sucesso** — o
refresh token continua bom. Para o vermelho há duas rotas, as duas do operador: revogar a
autorização do app na conta ML real (prova mais forte, custa reautorizar), ou apontar
`TokenURL` para um endpoint que devolve `invalid_grant` no dev stack (mais barato e
reversível, mas é mudança de env + restart).

Nenhuma ação foi tomada na conta real do operador.

## Dívida encontrada e NÃO consertada de passagem

- **Lane de integração vermelha antes da fatia.** `TestListingsReadContractEndToEnd`
  (`tests/integration/listings_read_test.go:181`, último toque `8f458e39`) faz a lane
  reportar `status=blocked`. Reproduz com o meu arquivo removido, e a fatia não tocou
  nenhum arquivo de listings ou catalog. Não investigada — fora de escopo.
- `mapIntegrationError` (`transport/http_handler.go:53`) casa por
  `strings.HasPrefix(err.Error(), "INTEGRATIONS_")`, então erros embrulhados com `%w`
  escapam e viram 400 genérico. Mesma classe do defeito desta fatia.
- `fee_sync_scheduler.go:74` descarta erro com `_, _ =` — mesmo padrão que a Task 4
  corrigiu no ticker de refresh.
