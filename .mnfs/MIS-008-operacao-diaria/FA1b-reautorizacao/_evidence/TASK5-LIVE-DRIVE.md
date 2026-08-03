# Task 5 — Live drive contra a conta ML real

Executado ao vivo em 2026-08-03, contra a conta ML real do operador
(`external_account_id=691607102`, installation
`inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0`).

## Pré-condição: stack repontada pra main

O backend em execução estava servindo o worktree da F-00
(`compose.f00.yml`, container criado `2026-08-03T20:53:29Z`) — não o tip da
`main` com a F-A1b mergeada. Reconstruído e recriado antes de qualquer ação
na conta real:

```
docker compose build backend   # imagem nova, confirmado por grep dentro da
                                # imagem: ensureProviderAccountUnlinked e
                                # INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED
                                # presentes em auth_flow_service.go / errors.go
docker compose up -d backend   # container recriado 2026-08-03T21:04:35Z
```

Sem esse passo o live drive teria testado código errado (mesma classe de
falha documentada em memória como "binário velho faz o live drive mentir").

## Step 1 — Operador revoga

Ação do operador, na conta real: Conta ML → Configurações → Segurança →
Opções de segurança e recuperação de senha → Ferramentas de segurança →
aplicações conectadas → remover o app. Confirmado pelo operador ("aplicação
removida").

## Step 2 — Refresh forçado

```sql
UPDATE integration_auth_sessions
   SET access_token_expires_at = now() - interval '1 minute',
       state = 'expiring',
       next_retry_at = NULL
 WHERE installation_id = 'inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0';
```

Log do ticker (21:15:02Z, dentro da janela de 5min/10min esperada):

```
2026/08/03 21:15:02 ERROR integrations refresh ticker item failed installation_id=inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0 err="INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400 body={\"message\":\"Error validating grant. Your authorization code or refresh token may be expired or it was already used\",\"error\":\"invalid_grant\",\"status\":400,\"cause\":[]}"
```

## Step 3 — Persistência e degradação confirmadas

```
 state          | refresh_failure_code                                    | consecutive_failures | next_retry_at
 refresh_failed | INTEGRATIONS_REFRESH_TOKEN_INVALID: ...invalid_grant... | 1                     | 2026-08-03 22:15:01 (+1h, cooldown terminal)

 status           | health_status | next_action | reauth_reason
 requires_reauth  | critical      | reauth      | INTEGRATIONS_REFRESH_TOKEN_INVALID: ...invalid_grant...
```

Bate com o esperado no plano, linha a linha.

## Step 4 — Tela mostra estado e botão

`/integracoes` aberta pelo domínio do túnel
(`https://multiradial-unironically-nieves.ngrok-free.dev/integracoes`).
Dump de acessibilidade confirma:

```
Mercado Livre (cliente)
Precisa reautorizar
[botão] Reautorizar
INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400 body={"message":"Error validating grant...invalid_grant"...}
```

Screenshot não obtido — a ferramenta de browser (headless) recusou
compositar frame ("Browser pane is not displayed"); a árvore de
acessibilidade + dump de texto acima serve como prova equivalente do
conteúdo renderizado.

## Step 5 — Botão devolve a conta

Clique em "Reautorizar" (disparado pela sessão) redirecionou para
`mercadolivre.com` (tela de login real). Operador completou login+consentimento
na própria conta ("Pronto, feito"). Retorno confirmado:

- Tela: `/integracoes` sem `?auth=failed`, card "Mercado Livre (cliente)" →
  "Conectado".
- Banco:

```sql
SELECT count(*) FROM integration_installations;                       -- 1
SELECT count(*) FILTER (WHERE is_active) FROM integration_credentials; -- 1
SELECT state, consecutive_failures, next_retry_at
  FROM integration_auth_sessions WHERE installation_id = '...';        -- valid, 0, NULL
```

Nenhuma segunda installation — o defeito original que esta fatia existe pra
matar não reapareceu.

## Step 6 — Controle positivo do caminho de recusa

Clique em "Conectar" (cria installation nova) autorizando a MESMA conta ML.
A sessão do ML ainda estava autenticada, então o consentimento foi automático
(sem tela de login) — recusa observada direto no callback:

- Tela: "A autorização não foi concluída: `INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED`"
  (aviso da Task 3, código exato).
- Banco, antes da limpeza:

```
installation_id ...b41c0a85...  status=pending_connection  external_account_id=NULL
integration_credentials: ativas=1, total=27   -- nenhuma credencial nova ativa
```

Só a installation rascunho foi criada — comportamento declarado fora de
escopo no plano. Apagada em seguida (autorizado pelo operador):

```sql
DELETE FROM integration_installations WHERE installation_id = 'inst-mercado_livre-b41c0a85-5daa-4451-a084-d85bdfe4ecbf';
```

## Estado final

```
installation_id                                          | status    | health_status | external_account_id
inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0   | connected | healthy       | 691607102
```

1 installation, 1 credencial ativa, saudável. Task 5 provada de ponta a
ponta contra a conta real: revogação → degradação detectada e visível →
reautorização sem duplicar → recusa de vínculo duplicado provada.

## Veredito

Task 5 — **CONCLUÍDA**. Todos os 7 steps do plano executados e conferidos
contra o comportamento esperado. F-A1b fecha sem pendência de validação.
