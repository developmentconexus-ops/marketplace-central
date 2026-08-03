# F-A1b — Reautorização da conta ML: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Dar ao operador uma ação que funcione quando a tela disser "Precisa reautorizar", e impedir que um callback de OAuth que falha deixe credencial e sessão órfãs no banco.

**Architecture:** O caminho de reautorização já existe inteiro no backend e está contratado até o SDK — falta o botão no web. O segundo defeito é ordem de escrita no `HandleCallback`: credencial e sessão são gravadas antes da atualização da installation que pode violar um índice único. Consertamos com uma checagem nomeada **antes** de qualquer escrita.

**Tech Stack:** Go 1.2x (`apps/server_core`), React + TypeScript + Vitest + Testing Library (`apps/web`), Postgres, harness PowerShell (`scripts/harness.ps1`).

---

## Contexto: como este plano nasceu

Em 2026-08-03, antes de fazer um live drive da fatia F-A1 (falha de token ML visível), o operador ia **revogar a autorização do app na conta ML real** para produzir um `invalid_grant` de verdade. Antes disso testamos o caminho de volta — reautorizar clicando "Conectar" enquanto ainda conectado. O teste reprovou, e é por isso que este plano existe.

Log do backend, verbatim:

```
16:36:47 INFO  integrations.installations action=list result=200 count=1
16:36:47 INFO  integrations.installations action=create_draft result=201 installation_id=inst-mercado_livre-6acd6d66-e3e0-4c6d-b312-ab83286c0e2a
16:36:48 INFO  integrations.auth.start   action=start_authorize result=200
16:36:52 WARN  integrations.auth.callback action=handle_callback result=302
               error="ERROR: duplicate key value violates unique constraint
               \"uq_integration_installations_active_provider_account\" (SQLSTATE 23505)"
16:36:56 INFO  integrations.installations action=list result=200 count=2
```

Sobrou no banco: installation `pending_connection` sem `external_account_id` nem `active_credential_id`, mais uma credencial `is_active=t` e uma `auth_session` `state=valid` apontando para ela. Duas contas ML na tela. O estado órfão foi apagado à mão depois (3 linhas), então **o banco do dev stack já está limpo**: 1 installation `connected`, 1 credencial ativa, 1 sessão.

**Consequência que trava a F-A1:** revogar no ML não apaga linha nenhuma nossa. A installation `connected` continuaria lá e todo clique em "Conectar" bateria no mesmo 23505 — o operador ficaria sem API do ML **sem caminho de volta pelo produto**. A F-A1 entrega um aviso ("Precisa reautorizar") sem saída. Este plano é pré-requisito do live drive da F-A1 (Task 7 de `2026-08-03-fa1-fa2-token-visivel-plan.md`).

## Medição (§12.3) — oito respostas com `file:line`

| # | Pergunta | Resposta |
|---|---|---|
| 1 | Quem já faz isso | `StartReauth` em `application/auth_flow_service.go:618`; guard de conta `hasReauthAccountMismatch` em `application/auth_flow_service.go:355` |
| 2 | Quem consome | **Zero call-sites no web.** Rota `POST /integrations/installations/{id}/reauth/authorize` em `transport/auth_handler.go:180`; OpenAPI em `contracts/api/marketplace-central.openapi.yaml:1283`; SDK `startIntegrationReauthorization` em `packages/sdk-runtime/src/index.ts:2184`. `grep -rn startIntegrationReauthorization apps/web/src` não bate nada |
| 3 | Qual contrato muda | **Nenhum.** A operação já está no OpenAPI e no SDK. Nada a regenerar |
| 4 | Produtor real do valor | `applyAuthResult` em `application/auth_flow_service.go:775` grava `status` + `external_account_id` e é quem viola `uq_integration_installations_active_provider_account` (`migrations/0017_oauth_credential_lifecycle.sql:29-32`, índice parcial sobre `(tenant_id, provider_code, external_account_id)` onde `status NOT IN ('disconnected','failed') AND external_account_id <> ''`) |
| 5 | Teste que pega regressão | **Nenhum.** `application/auth_flow_service_security_test.go` cobre state/nonce/assinatura; nenhum teste reautoriza installation já conectada nem exercita conta já vinculada |
| 6 | Dívida registrada que toca | Nenhuma. D-A/D-B/D-C nascem aqui |
| 7 | O que o provider entrega | Revogação da autorização no ML → `HTTP 400` com `{"error":"invalid_grant"}` (doc oficial do ML; mesma string que `classifyRefreshHTTPError` já casa em `adapters/mercadolivre/auth_adapter.go`) |
| 8 | Local maximum | Reauth existe ponta a ponta menos o botão; `List` já existe em `authFlowInstallationStore` (`application/auth_flow_service.go:130`) e serve à checagem sem SQL novo |

## Anti-redundância (§12.4)

Nenhuma task cria cálculo, helper ou endpoint novo sem esta linha:

| O que a task faz | O que já existe | Por que não serve sozinho |
|---|---|---|
| T1 botão de reautorizar | `startIntegrationReauthorization` (`sdk-runtime/src/index.ts:2184`) | Existe e funciona; nenhum componente do web o chama |
| T1 rótulo "Reautorizar" | `nextActionLabel` (`ConnectionHealthCard.tsx:33`) | É `<span>`, texto morto; vira o rótulo do botão, sem string nova |
| T2 checagem de conta já vinculada | `List` (`auth_flow_service.go:130`), `hasReauthAccountMismatch` (`:355`) | O guard existente protege a installation **atual** contra conta trocada; não protege **outra** installation contra a mesma conta |
| T2 erro nomeado | `domain/errors.go:24` (`ErrReauthAccountMismatch`) | Erro diferente: aqui a conta bate, o que colide é o vínculo em outra installation |
| T3 motivo na tela | `buildIntegrationsRedirectPath` (`transport/auth_handler.go:416`) já monta `?auth=failed` | Nada no web lê esse parâmetro; e `failed` não diz **por quê** |

## Três defeitos

- **D-A — não existe ação de reautorização no produto.** `connect()` (`apps/web/src/pages/integracoes/IntegracoesPage.tsx:513`) só reaproveita installation em `pending_connection`/`draft`; para conta já conectada, "Conectar" significa *adicionar outra conta*. O comentário em `:492` já contava com um caminho de reauth que nunca ganhou UI.
- **D-B — `HandleCallback` não é atômico.** Ordem: `saveCredential` (`:359`) → `authSessions.Upsert` (`:364`) → `applyAuthResult` (`:375`). Erro no terceiro deixa os dois primeiros gravados.
- **D-C — o motivo da falha do callback nunca chega à tela.** `transport/auth_handler.go:76` redireciona para `?auth=failed` sem razão, e `grep -rn "searchParams" apps/web/src/pages/integracoes` não bate nada — nem o `failed` é lido. `LegacyRedirect` (`apps/web/src/app/LegacyRedirect.tsx:10`) preserva a query, então o parâmetro chega; falta quem leia.

## Escopo declarado — o que este plano NÃO faz

**D-B é consertado pela causa alcançável, não por transação.** A correção genérica seria uma porta de unidade-de-trabalho atravessando três repositórios (credencial, sessão, installation), cada um com seu pool — mudança arquitetural que não cabe aqui. A T2 fecha a única causa que reproduzimos e que o operador alcança pela tela. **Fica em aberto:** qualquer outra falha depois de `saveCredential` ainda deixa escrita parcial. Registrar como dívida `D-45` no fechamento (ver Task 6).

Abandonar o fluxo no meio (clicar "Reautorizar" e não concluir no ML) deixa a installation em `pending_connection`, porque `StartAuthorize` aplica um snapshot pendente em `auth_flow_service.go:320`. Isso é comportamento existente e **não** é regressão desta fatia; `connect()` reaproveita installations pendentes, então não acumula. Não consertar aqui.

## Fake não é evidência

Ordem do operador, 2026-08-03: *"não quero mock, fallback, validação sempre real para realmente validar, não fazer teste apenas pra passar"*.

- **A fatia não fecha com T1–T3 verdes.** Só fecha com T4 e T5 verdes.
- Dois testes abaixo são **controles negativos**: `reauth of the same account is allowed` (T2) e `nao mostra motivo quando o callback nao falhou` (T3).
- **Toda asserção de não-chamada precisa de must-fail por injeção**, e precisa rodar depois de um flush de microtask: em `IntegracoesPage.test.tsx` uma asserção síncrona sobre `createIntegrationInstallation` passou mesmo com `fireEvent.click` injetado, porque o handler só chega lá depois de um `await`. Ver `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx` (commit `aafd014a`) para o padrão correto.

## File Structure

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `apps/web/src/pages/integracoes/ConnectionHealthCard.tsx` | Modificar | Renderizar botão "Reautorizar" quando `next_action === "reauth"` e redirecionar para o `auth_url` |
| `apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx` | Modificar | Ramo do botão, ramo sem botão, e falha da chamada |
| `apps/server_core/internal/modules/integrations/domain/errors.go` | Modificar | Sentinela `ErrProviderAccountAlreadyLinked` |
| `apps/server_core/internal/modules/integrations/application/auth_flow_service.go` | Modificar | Checagem antes de qualquer escrita no `HandleCallback` |
| `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go` | Modificar | Ramo que recusa, controle negativo que permite, e prova de não-escrita |
| `apps/server_core/internal/modules/integrations/transport/auth_handler.go` | Modificar | Levar o código do erro no redirect do callback |
| `apps/server_core/internal/modules/integrations/transport/auth_handler_test.go` | Modificar | Redirect carrega o código; sucesso não carrega |
| `apps/web/src/pages/integracoes/IntegracoesPage.tsx` | Modificar | Ler `?auth` e `?reason` e mostrar o motivo |
| `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx` | Modificar | Mostra o motivo; controle negativo sem motivo |
| `apps/server_core/tests/integration/integrations_reauth_test.go` | Criar | Prova contra Postgres real: reauth atualiza a mesma linha; conta já vinculada é recusada sem escrever |

## Comandos

Sempre de dentro de `apps/server_core` para os testes Go (o `.gitignore` da raiz engole `.gocache` se rodar de fora):

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/integrations/... -count=1
```

Web:

```bash
cd apps/web && npx vitest run src/pages/integracoes/
```

Lane hermética (auto-provisiona o Postgres e injeta `MPC_TEST_DATABASE_URL`; **não** passe `-DatabaseUrl`, isso lança `HPG_EXTERNAL_TARGET_FORBIDDEN` de propósito):

```bash
pwsh scripts/harness.ps1 -Command integration
```

> **Vermelho pré-existente:** `TestListingsReadContractEndToEnd` (`tests/integration/listings_read_test.go:181`) já falha na main e faz a lane reportar `status=blocked`. Não é seu. Confirme que o `failure_token=test=` do **seu** teste não aparece; é assim que se lê verde nesta lane hoje.

**Ordem das tasks:** 1 botão · 2 checagem antes da escrita · 3 motivo visível · 4 prova real · 5 drive ao vivo · 6 fechamento.

---

### Task 1: Botão de reautorizar em `/integracoes`

O rótulo "Reautorizar" já é renderizado como texto morto em `ConnectionHealthCard.tsx:33`. Vira botão. Zero backend, zero contrato.

**Files:**
- Modify: `apps/web/src/pages/integracoes/ConnectionHealthCard.tsx`
- Test: `apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx`

- [ ] **Step 1: Escreva os testes que falham**

Acrescente ao fim de `ConnectionHealthCard.test.tsx`, dentro do `describe` existente. Leia o topo do arquivo primeiro: ele já monta um `QueryClientProvider` e mocka `useClient`/`useInstallation`; **reuse os helpers que já estiverem lá em vez de criar novos**, e só adicione ao mock de `useClient` o que faltar.

```tsx
  it("oferece reautorizar quando o proximo passo e reauth e manda o browser para o ML", async () => {
    startIntegrationReauthorization.mockResolvedValue({
      auth_url: "https://auth.mercadolivre.com.br/authorization?x=1",
      expires_in: 600,
    });
    renderCard({
      installation_id: "inst-1",
      status: "requires_reauth",
      connection: {
        state: "needs_reauth",
        next_action: "reauth",
        reauth_reason: "INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400",
      },
    });

    const button = screen.getByTestId("connection-reauth-inst-1");
    expect(button).toHaveTextContent("Reautorizar");

    fireEvent.click(button);

    await waitFor(() => expect(startIntegrationReauthorization).toHaveBeenCalledWith("inst-1"));
    await waitFor(() =>
      expect(assignedUrl).toBe("https://auth.mercadolivre.com.br/authorization?x=1"),
    );
  });

  // Controle negativo: uma conta saudavel nao pode ganhar um botao que dispara
  // OAuth. Sem este teste, renderizar o botao sempre passaria no teste acima.
  it("nao oferece reautorizar quando a conexao esta saudavel", async () => {
    renderCard({
      installation_id: "inst-1",
      status: "connected",
      connection: { state: "connected", next_action: "none", reauth_reason: "" },
    });

    expect(await screen.findByTestId("connection-health-inst-1")).toBeInTheDocument();
    expect(screen.queryByTestId("connection-reauth-inst-1")).not.toBeInTheDocument();
    expect(startIntegrationReauthorization).not.toHaveBeenCalled();
  });

  it("mostra erro nomeado e nao navega quando a reautorizacao falha", async () => {
    startIntegrationReauthorization.mockRejectedValue(new Error("boom"));
    renderCard({
      installation_id: "inst-1",
      status: "requires_reauth",
      connection: { state: "needs_reauth", next_action: "reauth", reauth_reason: "" },
    });

    fireEvent.click(screen.getByTestId("connection-reauth-inst-1"));

    expect(await screen.findByTestId("connection-reauth-error-inst-1")).toHaveTextContent(
      "Não foi possível iniciar a reautorização.",
    );
    expect(assignedUrl).toBe("");
  });
```

Para esses testes existirem você precisa de três coisas no topo do arquivo. Acrescente **só o que ainda não estiver lá**:

```tsx
const startIntegrationReauthorization = vi.fn();

// window.location.href não é atribuível em jsdom sem isso. Capturar em vez de
// navegar é o que torna a asserção de navegação possível.
let assignedUrl = "";
beforeEach(() => {
  assignedUrl = "";
  startIntegrationReauthorization.mockReset();
  Object.defineProperty(window, "location", {
    configurable: true,
    value: {
      get href() {
        return assignedUrl;
      },
      set href(value: string) {
        assignedUrl = value;
      },
    },
  });
});
```

e, no mock de `../../app/ClientContext`, o método:

```tsx
    startIntegrationReauthorization: (...args: unknown[]) => startIntegrationReauthorization(...args),
```

`renderCard(installation)` é o helper que monta uma installation e renderiza `<ConnectionHealthCard />`. Se o arquivo ainda não tiver um, escreva-o com esta forma exata, reusando o `QueryClientProvider` que já existe no arquivo:

```tsx
function renderCard(overrides: Record<string, unknown>) {
  const installation = {
    installation_id: "inst-1",
    provider_code: "mercado_livre",
    display_name: "Mercado Livre",
    family: "marketplace",
    status: "connected",
    health_status: "healthy",
    connection: { state: "connected", next_action: "none", reauth_reason: "" },
    ...overrides,
  };
  useInstallationMock.mockReturnValue({ installations: [installation], status: "ready" });
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ConnectionHealthCard />
    </QueryClientProvider>,
  );
}
```

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/web && npx vitest run src/pages/integracoes/ConnectionHealthCard.test.tsx
```

Esperado: os três novos falham. O primeiro e o terceiro por `Unable to find an element by: [data-testid="connection-reauth-inst-1"]`. O segundo **passa** desde já — é controle negativo, e passar agora é correto: ele só ganha valor depois que o botão existir.

- [ ] **Step 3: Implemente**

Em `ConnectionHealthCard.tsx`, troque a linha do rótulo de ação (`{action ? <span className="text-faint">Ação: {action}</span> : null}`, hoje em `:70`) por um bloco que vira botão quando a ação é reautorizar. Adicione ao topo do arquivo:

```tsx
import { useState } from "react";
import { useClient } from "../../app/ClientContext";
```

e substitua o corpo de `InstallationRow`:

```tsx
function InstallationRow({ installation }: { installation: IntegrationInstallation }) {
  const client = useClient();
  const connection = installation.connection;
  const tone = stateTone[connection.state];
  const action = nextActionLabel[connection.next_action];
  const reason = connection.reauth_reason?.trim() ?? "";
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // O backend inteiro do reauth já existe e está contratado
  // (transport/auth_handler.go:180, OpenAPI:1283, sdk-runtime/src/index.ts:2184).
  // O único elo que faltava era este clique. Reautorizar pina a MESMA
  // installation, então o callback atualiza a linha existente em vez de tentar
  // inserir uma segunda para o mesmo seller — que é o 23505 que este plano
  // nasceu para matar.
  async function reauth() {
    setBusy(true);
    setError(null);
    try {
      const start = await client.startIntegrationReauthorization(installation.installation_id);
      window.location.href = start.auth_url;
    } catch {
      setError("Não foi possível iniciar a reautorização.");
      setBusy(false);
    }
  }

  return (
    <div
      className="flex flex-col gap-1 rounded-control border border-border bg-surface-2 px-3 py-2 text-xs"
      data-testid={`connection-health-${installation.installation_id}`}
    >
      <div className="flex items-center justify-between gap-3">
        <span className="font-medium text-ink">{installation.display_name}</span>
        <span
          className={`inline-flex items-center gap-1 whitespace-nowrap rounded-pill px-2 py-0.5 font-medium ${toneBadgeClassName[tone]}`}
        >
          <span className="h-1.5 w-1.5 rounded-pill bg-current" aria-hidden="true" />
          {stateLabel[connection.state]}
        </span>
      </div>
      {connection.next_action === "reauth" ? (
        <button
          type="button"
          onClick={() => void reauth()}
          disabled={busy}
          data-testid={`connection-reauth-${installation.installation_id}`}
          className="mt-1 self-start rounded-control border border-border px-2 py-1 font-medium text-ink disabled:opacity-60"
        >
          {busy ? "Abrindo…" : nextActionLabel.reauth}
        </button>
      ) : action ? (
        <span className="text-faint">Ação: {action}</span>
      ) : null}
      {error ? (
        <span className="text-warn" data-testid={`connection-reauth-error-${installation.installation_id}`}>
          {error}
        </span>
      ) : null}
      {/* O motivo é o erro cru do provider. Mostrar cru é deliberado: é o único
          diagnóstico que existe e traduzi-lo apagaria o código que o operador
          precisa citar num chamado. */}
      {reason ? (
        <span className="break-words text-faint" data-testid={`connection-health-reason-${installation.installation_id}`}>
          {reason}
        </span>
      ) : null}
    </div>
  );
}
```

- [ ] **Step 4: Rode e veja passar**

```bash
cd apps/web && npx vitest run src/pages/integracoes/
```

Esperado: todos os arquivos de `src/pages/integracoes/` verdes. Se `IntegracoesPage.test.tsx` quebrar, é porque `ConnectionHealthCard` agora usa `useClient` — acrescente `startIntegrationReauthorization` ao mock de `useClient` daquele arquivo também, **sem apagar asserção nenhuma**. Se alguma asserção alheia ficar impossível, restaure a intenção contra outro observável e prove por injeção (padrão em `IntegracoesPage.test.tsx`, commit `aafd014a`); apagar não é opção.

- [ ] **Step 5: Must-fail do controle negativo**

Verde do controle negativo não vale nada se ele nunca puder reprovar. Troque temporariamente a condição do botão para `true`:

```tsx
      {true ? (
```

Rode de novo. Esperado: `nao oferece reautorizar quando a conexao esta saudavel` FALHA. Restaure `connection.next_action === "reauth"` e confirme verde.

- [ ] **Step 6: Commit**

```bash
git add apps/web/src/pages/integracoes/ConnectionHealthCard.tsx apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx
git commit -m "feat(web): give the operator a working reauthorize action"
```

---

### Task 2: Recusar conta já vinculada ANTES de escrever

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/domain/errors.go`
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go`
- Test: `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`

- [ ] **Step 1: Escreva os testes que falham**

Acrescente ao fim de `auth_flow_service_test.go`. Leia os fakes que já existem no arquivo antes (`flowInstallationStore`, `flowAuthWriter`, e o construtor de serviço usado pelos testes vizinhos) e **monte o fixture com eles**, sem inventar fake novo.

```go
// Duas installations do mesmo provider, uma delas já vinculada ao seller que o
// callback acabou de devolver. Antes desta task o serviço gravava credencial e
// sessão e só então batia no índice único, deixando as duas órfãs.
func TestHandleCallbackRefusesAccountLinkedToAnotherInstallationBeforeWriting(t *testing.T) {
	h := newCallbackHarness(t, "seller-1")
	h.installations.items = append(h.installations.items, domain.Installation{
		InstallationID:      "inst-old",
		ProviderCode:        "mercado_livre",
		ExternalAccountID:   "seller-1",
		Status:              domain.InstallationStatusConnected,
	})

	_, err := h.svc.HandleCallback(context.Background(), h.input)

	if !errors.Is(err, domain.ErrProviderAccountAlreadyLinked) {
		t.Fatalf("err = %v, want ErrProviderAccountAlreadyLinked", err)
	}
	if got := len(h.credentials.rotated); got != 0 {
		t.Fatalf("credenciais gravadas = %d, want 0 (a recusa tem que vir antes da escrita)", got)
	}
	if got := len(h.sessions.sessions); got != 0 {
		t.Fatalf("sessoes gravadas = %d, want 0 (a recusa tem que vir antes da escrita)", got)
	}
	if got := len(h.installations.connectionSnapshots); got != 0 {
		t.Fatalf("snapshots aplicados = %d, want 0", got)
	}
}

// Controle negativo: reautorizar a MESMA installation com a MESMA conta é o
// caminho feliz desta fatia inteira. Se a checagem recusar isso, ela quebrou o
// que veio consertar.
func TestHandleCallbackAllowsReauthOfTheSameInstallation(t *testing.T) {
	h := newCallbackHarness(t, "seller-1")
	h.installations.items[0].ExternalAccountID = "seller-1"
	h.installations.items[0].Status = domain.InstallationStatusRequiresReauth

	if _, err := h.svc.HandleCallback(context.Background(), h.input); err != nil {
		t.Fatalf("HandleCallback err = %v, want nil", err)
	}
	if got := len(h.installations.connectionSnapshots); got != 1 {
		t.Fatalf("snapshots aplicados = %d, want 1", got)
	}
}

// Uma installation desconectada NÃO segura a conta: o índice único de
// 0017_oauth_credential_lifecycle.sql:31 exclui 'disconnected' e 'failed', e a
// checagem tem que usar exatamente o mesmo predicado, senão ela recusa o que o
// banco aceitaria.
func TestHandleCallbackIgnoresDisconnectedInstallationHoldingTheAccount(t *testing.T) {
	h := newCallbackHarness(t, "seller-1")
	h.installations.items = append(h.installations.items, domain.Installation{
		InstallationID:    "inst-old",
		ProviderCode:      "mercado_livre",
		ExternalAccountID: "seller-1",
		Status:            domain.InstallationStatusDisconnected,
	})

	if _, err := h.svc.HandleCallback(context.Background(), h.input); err != nil {
		t.Fatalf("HandleCallback err = %v, want nil", err)
	}
}
```

`newCallbackHarness(t, providerAccountID)` é o fixture. Escreva-o no mesmo arquivo, ligando os fakes existentes. O ponto sensível é o estado do OAuth: `HandleCallback` começa por `verifyAndConsumeCallbackState`, então o harness tem que produzir um `state` válido — a forma mais barata e honesta é **chamar `StartAuthorize` primeiro** e usar o `state` que ele devolve, exatamente como faz `TestStartAuthorizeHandleCallbackAcceptsPersistedSignedState` em `auth_flow_service_security_test.go:163`. Copie a montagem de lá; não reimplemente assinatura de state à mão.

```go
type callbackHarness struct {
	svc           *AuthFlowService
	installations *flowInstallationStore
	credentials   *flowCredentialRotator
	sessions      *flowAuthWriter
	input         HandleCallbackInput
}
```

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/integrations/application/... -count=1 -run TestHandleCallback
```

Esperado: `TestHandleCallbackRefusesAccountLinkedToAnotherInstallationBeforeWriting` falha na compilação com `undefined: domain.ErrProviderAccountAlreadyLinked`. Depois de adicionar a sentinela (Step 3, primeira metade), ele passa a falhar por comportamento: `err = <nil>, want ErrProviderAccountAlreadyLinked`. **Observe as duas falhas**; a segunda é a que importa.

- [ ] **Step 3: Implemente**

Em `domain/errors.go`, ao lado de `ErrReauthAccountMismatch` (`:24`):

```go
	// A conta do provider bate com a que o operador autorizou, mas OUTRA
	// installation ativa já a detém. Distinto de ErrReauthAccountMismatch, onde
	// quem não bate é a conta desta installation.
	ErrProviderAccountAlreadyLinked = errors.New("INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED")
```

Em `application/auth_flow_service.go`, dentro de `HandleCallback`, logo **depois** do guard de mismatch (`:355-357`) e **antes** de `saveCredential` (`:359`):

```go
	if err := s.ensureProviderAccountUnlinked(ctx, inst, payload.ProviderAccountID); err != nil {
		return AuthStatus{}, err
	}
```

E o helper, depois de `HandleCallback`:

```go
// ensureProviderAccountUnlinked recusa o callback quando outra installation
// ativa já detém a conta do provider, ANTES de qualquer escrita.
//
// Sem isso a violação só aparecia em applyAuthResult (:775), depois de
// saveCredential e do upsert da sessão — e como as três escritas não
// compartilham transação, a falha deixava credencial ativa e sessão válida
// penduradas numa installation que nunca ficou conectada.
//
// O predicado é o do índice parcial uq_integration_installations_active_provider_account
// (migrations/0017_oauth_credential_lifecycle.sql:29-32) e tem que continuar
// igual a ele: se divergir, recusamos o que o banco aceitaria, ou pior,
// deixamos passar o que ele recusa.
func (s *AuthFlowService) ensureProviderAccountUnlinked(
	ctx context.Context,
	inst domain.Installation,
	providerAccountID string,
) error {
	accountID := strings.TrimSpace(providerAccountID)
	if accountID == "" {
		// Índice parcial não cobre external_account_id vazio; não inventamos
		// regra que o banco não tem.
		return nil
	}

	installations, err := s.installations.List(ctx)
	if err != nil {
		return err
	}

	for _, other := range installations {
		if other.InstallationID == inst.InstallationID {
			continue
		}
		if other.ProviderCode != inst.ProviderCode {
			continue
		}
		if strings.TrimSpace(other.ExternalAccountID) != accountID {
			continue
		}
		if other.Status == domain.InstallationStatusDisconnected ||
			other.Status == domain.InstallationStatusFailed {
			continue
		}
		return fmt.Errorf("%w: installation_id=%s", domain.ErrProviderAccountAlreadyLinked, other.InstallationID)
	}

	return nil
}
```

Confirme que `strings` e `fmt` já estão importados no arquivo (estão, `fmt` é usado em `:365`).

- [ ] **Step 4: Rode e veja passar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/integrations/... -count=1
```

Esperado: `ok` em todos os pacotes de `integrations`.

- [ ] **Step 5: Must-fail por injeção**

Mova a chamada de `ensureProviderAccountUnlinked` para **depois** de `saveCredential`. Rode de novo. Esperado: `TestHandleCallbackRefusesAccountLinkedToAnotherInstallationBeforeWriting` FALHA com `credenciais gravadas = 1, want 0 (a recusa tem que vir antes da escrita)` — o erro certo sendo devolvido pelo motivo errado. Restaure a posição original e confirme verde. **Sem este passo o teste não distingue "recusou" de "recusou sem sujar".**

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/integrations/domain/errors.go apps/server_core/internal/modules/integrations/application/auth_flow_service.go apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go
git commit -m "fix(integrations): refuse an already-linked provider account before writing"
```

---

### Task 3: O motivo da falha do callback chega à tela

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/transport/auth_handler.go`
- Test: `apps/server_core/internal/modules/integrations/transport/auth_handler_test.go`
- Modify: `apps/web/src/pages/integracoes/IntegracoesPage.tsx`
- Test: `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx`

- [ ] **Step 1: Teste do backend que falha**

Acrescente a `auth_handler_test.go`, reusando o `stubAuthFlow` que já existe em `:44`:

```go
func TestAuthHandlerCallbackCarriesTheFailureCodeInTheRedirect(t *testing.T) {
	flow := &stubAuthFlow{callbackErr: fmt.Errorf("%w: installation_id=inst-old", domain.ErrProviderAccountAlreadyLinked)}
	rec := httptest.NewRecorder()
	newTestAuthHandlerMux(flow).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=c&state=s", nil))

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	location, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatalf("Location invalido: %v", err)
	}
	if got := location.Query().Get("auth"); got != "failed" {
		t.Fatalf("auth = %q, want failed", got)
	}
	if got := location.Query().Get("reason"); got != "INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED" {
		t.Fatalf("reason = %q, want INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED", got)
	}
}

// Controle negativo: sucesso nao pode carregar reason. Sem ele, gravar reason
// sempre passaria no teste acima.
func TestAuthHandlerCallbackSuccessCarriesNoReason(t *testing.T) {
	flow := &stubAuthFlow{}
	rec := httptest.NewRecorder()
	newTestAuthHandlerMux(flow).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=c&state=s", nil))

	location, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatalf("Location invalido: %v", err)
	}
	if got := location.Query().Get("reason"); got != "" {
		t.Fatalf("reason = %q, want vazio", got)
	}
}
```

`newTestAuthHandlerMux(flow)` e `stubAuthFlow.callbackErr` podem não existir com esses nomes. **Leia `auth_handler_test.go:44` e os testes de callback vizinhos (`:306-430`) e use a montagem que já está lá**, ajustando os nomes; não crie um segundo jeito de montar o handler no mesmo arquivo.

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/integrations/transport/... -count=1 -run TestAuthHandlerCallback
```

Esperado: o primeiro falha com `reason = "", want INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED`. O segundo passa (controle negativo).

- [ ] **Step 3: Implemente o backend**

Em `auth_handler.go`, troque o corpo do `if err != nil` do callback (hoje `:74-78`):

```go
	if err != nil {
		// O código do erro vai na URL; o texto cru NÃO. Um erro embrulhado com
		// %w carrega installation_id e mensagem de driver, e URL vira log,
		// histórico e print de tela. O operador precisa do código para citar
		// num chamado — o resto está no log do servidor.
		slog.Warn("integrations.auth.callback", "action", "handle_callback", "result", "302", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		http.Redirect(w, r, buildOAuthCallbackRedirectURL(buildIntegrationsFailedRedirectPath(err)), http.StatusFound)
		return
	}
```

E, ao lado de `buildIntegrationsRedirectPath` (`:416`):

```go
// callbackFailureCodes lista as sentinelas cujo código pode ir para a URL.
// Lista fechada de propósito: erro não previsto vira "unknown" em vez de
// vazar texto de driver na barra de endereço.
var callbackFailureCodes = []error{
	domain.ErrProviderAccountAlreadyLinked,
	domain.ErrReauthAccountMismatch,
}

func buildIntegrationsFailedRedirectPath(err error) string {
	reason := "unknown"
	for _, candidate := range callbackFailureCodes {
		if errors.Is(err, candidate) {
			reason = candidate.Error()
			break
		}
	}
	query := url.Values{}
	query.Set("auth", "failed")
	query.Set("reason", reason)
	return "/integrations?" + query.Encode()
}
```

Confirme que `errors` e o pacote `domain` estão importados no arquivo; adicione se faltar.

- [ ] **Step 4: Backend verde**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/integrations/... -count=1
```

- [ ] **Step 5: Teste do frontend que falha**

`LegacyRedirect` (`apps/web/src/app/LegacyRedirect.tsx:10`) preserva a query, então `/integrations?auth=failed&reason=...` chega em `/integracoes` com os parâmetros. Falta quem leia. Acrescente a `IntegracoesPage.test.tsx` — o arquivo já renderiza dentro de `MemoryRouter`, então dá para injetar a URL por `initialEntries`:

```tsx
  it("mostra o motivo quando o callback do OAuth falhou", async () => {
    renderPageAt("/integracoes?auth=failed&reason=INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED");

    expect(await screen.findByTestId("oauth-callback-error")).toHaveTextContent(
      "INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED",
    );
  });

  // Controle negativo: sem falha na URL nao pode aparecer aviso nenhum.
  it("nao mostra aviso de callback quando a URL nao traz falha", async () => {
    renderPageAt("/integracoes");

    expect(await screen.findByTestId("provider-connect-ml")).toBeInTheDocument();
    expect(screen.queryByTestId("oauth-callback-error")).not.toBeInTheDocument();
  });
```

`renderPageAt(url)` é `renderPage()` com `initialEntries`. Refatore o `renderPage()` existente para receber a URL com default `"/integracoes"` em vez de duplicá-lo:

```tsx
function renderPage(initialEntry = "/integracoes") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <MemoryRouter initialEntries={[initialEntry]}>
      <QueryClientProvider client={queryClient}>
        <InstallationProvider>
          <IntegracoesPage />
        </InstallationProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}
const renderPageAt = (url: string) => renderPage(url);
```

- [ ] **Step 6: Implemente o frontend**

Em `IntegracoesPage.tsx`, no componente de página (não no `ProviderConnectCard`), leia a query e renderize o aviso acima dos cards:

```tsx
import { useSearchParams } from "react-router-dom";
```

```tsx
  const [searchParams] = useSearchParams();
  const callbackFailureReason =
    searchParams.get("auth") === "failed" ? (searchParams.get("reason") ?? "unknown") : null;
```

```tsx
      {/* O código cru é deliberado, mesma decisão do reauth_reason no
          ConnectionHealthCard: é o único diagnóstico que existe e traduzi-lo
          apagaria o código que o operador precisa citar num chamado. */}
      {callbackFailureReason ? (
        <div
          data-testid="oauth-callback-error"
          className="rounded-control border border-border bg-warn-soft px-3 py-2 text-xs text-warn"
        >
          A autorização não foi concluída: {callbackFailureReason}
        </div>
      ) : null}
```

- [ ] **Step 7: Frontend verde e must-fail**

```bash
cd apps/web && npx vitest run src/pages/integracoes/
```

Depois, injete o defeito: troque `searchParams.get("auth") === "failed"` por `true`. Esperado: `nao mostra aviso de callback quando a URL nao traz falha` FALHA. Restaure e confirme verde.

- [ ] **Step 8: Commit**

```bash
git add apps/server_core/internal/modules/integrations/transport/auth_handler.go apps/server_core/internal/modules/integrations/transport/auth_handler_test.go apps/web/src/pages/integracoes/IntegracoesPage.tsx apps/web/src/pages/integracoes/IntegracoesPage.test.tsx
git commit -m "fix(integrations): surface why the OAuth callback failed"
```

---

### Task 4: Prova real contra Postgres

As Tasks 1–3 provaram ramificação contra fakes. Duas coisas continuam não provadas e são as que sustentam a fatia:

1. Que reautorizar a **mesma** installation atualiza a linha existente sem violar `uq_integration_installations_active_provider_account`. Nenhum fake conhece índice parcial.
2. Que o predicado da checagem da Task 2 é **o mesmo** do índice. Se divergirem, o fake nunca conta.

**Files:**
- Create: `apps/server_core/tests/integration/integrations_reauth_test.go`

- [ ] **Step 1: Escreva o teste**

Copie a montagem de fixture do vizinho `tests/integration/integrations_refresh_failure_test.go`, que já compila e usa os mesmos repositórios. Atenção aos nomes reais das tabelas: `integration_installations` e `integration_auth_sessions` (com prefixo).

```go
//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Prova contra Postgres real as duas coisas que os fakes das Tasks 1-3 não
// alcançam: que reautorizar a MESMA installation não viola o índice parcial, e
// que uma installation 'disconnected' de fato libera a conta — que é o
// predicado exato que ensureProviderAccountUnlinked replica.
func TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_reauth")
	testpostgres.SeedProvider(t, pool, testpostgres.ProviderFixture{Code: "mercado_livre", DisplayName: "Mercado Livre"})

	repo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	ctx := context.Background()
	now := time.Now().UTC()
	account := fmt.Sprintf("seller-%d", now.UnixNano())
	first := fmt.Sprintf("inst-reauth-a-%d", now.UnixNano())
	second := fmt.Sprintf("inst-reauth-b-%d", now.UnixNano())

	t.Cleanup(func() {
		for _, id := range []string{first, second} {
			if _, err := pool.Exec(ctx, `DELETE FROM integration_auth_sessions WHERE installation_id = $1`, id); err != nil {
				t.Errorf("cleanup integration_auth_sessions %s: %v", id, err)
			}
			if _, err := pool.Exec(ctx, `DELETE FROM integration_installations WHERE installation_id = $1`, id); err != nil {
				t.Errorf("cleanup integration_installations %s: %v", id, err)
			}
		}
	})

	create := func(id string, status integrationsdomain.InstallationStatus) error {
		return repo.CreateInstallation(ctx, integrationsdomain.Installation{
			InstallationID: id,
			TenantID:       cfg.DefaultTenantID,
			ProviderCode:   "mercado_livre",
			Family:         integrationsdomain.IntegrationFamilyMarketplace,
			DisplayName:    "Mercado Livre (teste)",
			Status:         status,
			HealthStatus:   integrationsdomain.HealthStatusHealthy,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	connect := func(id string, status integrationsdomain.InstallationStatus) error {
		inst := integrationsdomain.Installation{
			InstallationID:    id,
			TenantID:          cfg.DefaultTenantID,
			ProviderCode:      "mercado_livre",
			ExternalAccountID: account,
			Status:            status,
			HealthStatus:      integrationsdomain.HealthStatusHealthy,
		}
		snapshot := integrationsdomain.ProjectConnectionSnapshot(
			inst, integrationsdomain.AuthStrategyOAuth2, nil, "",
		)
		return repo.ApplyConnectionSnapshot(ctx, id, snapshot, "")
	}

	if err := create(first, integrationsdomain.InstallationStatusPendingConnection); err != nil {
		t.Fatalf("CreateInstallation first: %v", err)
	}
	if err := connect(first, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("primeira conexao: %v", err)
	}

	// --- 1. Reautorizar a MESMA installation continua sendo uma linha só -----
	if err := connect(first, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("reautorizacao da mesma installation: %v (o indice parcial nao deveria disparar)", err)
	}
	if got := countActiveForAccount(t, pool, account); got != 1 {
		t.Fatalf("linhas ativas com a conta = %d, want 1", got)
	}

	// --- 2. Controle positivo: uma SEGUNDA installation ativa é recusada -----
	// Sem este passo, o passo 3 seria vacuoso: se o índice não existisse, tudo
	// passaria e pareceria "liberado".
	if err := create(second, integrationsdomain.InstallationStatusPendingConnection); err != nil {
		t.Fatalf("CreateInstallation second: %v", err)
	}
	if err := connect(second, integrationsdomain.InstallationStatusConnected); err == nil {
		t.Fatal("uma segunda installation ativa com a mesma conta foi aceita: o indice unico nao esta valendo")
	}

	// --- 3. Desconectar a primeira LIBERA a conta ---------------------------
	// É exatamente o predicado que ensureProviderAccountUnlinked replica; se o
	// banco e o código discordarem aqui, um dos dois está errado.
	if err := connect(first, integrationsdomain.InstallationStatusDisconnected); err != nil {
		t.Fatalf("desconectar a primeira: %v", err)
	}
	if err := connect(second, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("segunda installation apos desconectar a primeira: %v (o indice exclui 'disconnected')", err)
	}
}

func countActiveForAccount(t *testing.T, pool interface {
	QueryRow(context.Context, string, ...any) interface{ Scan(...any) error }
}, account string) int {
	t.Helper()
	var count int
	err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM integration_installations
		 WHERE external_account_id = $1
		   AND status NOT IN ('disconnected', 'failed')
	`, account).Scan(&count)
	if err != nil {
		t.Fatalf("contagem: %v", err)
	}
	return count
}
```

> A assinatura de `countActiveForAccount` acima usa uma interface estrutural para não depender do tipo concreto do pool. Se o `pgxpool.Pool` não satisfizer, troque o parâmetro por `*pgxpool.Pool` e importe `github.com/jackc/pgx/v5/pgxpool` — é o que o vizinho `integrations_refresh_failure_test.go` faz ao chamar `pool.Exec` direto.

- [ ] **Step 2: Rode a lane**

```bash
pwsh scripts/harness.ps1 -Command integration
```

Esperado no cabeçalho: `target=ephemeral-postgres`, `key=MPC_TEST_DATABASE_URL`, `migrations=embedded`, `migrations_second=0`.

Como ler o resultado: `TestListingsReadContractEndToEnd` já falha na main e faz a lane reportar `status=blocked`. **Seu teste está verde se `failure_token=test=TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount` NÃO aparecer na saída.**

- [ ] **Step 3: Injete o defeito e veja o vermelho**

Verde sozinho não prova que o teste mede alguma coisa. Comente o índice em `migrations/0017_oauth_credential_lifecycle.sql:29-32` e rode a lane de novo — a lane aplica as migrations num banco novo, então o índice some de verdade.

Esperado: `failure_token=test=TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount` aparece, com a mensagem `uma segunda installation ativa com a mesma conta foi aceita: o indice unico nao esta valendo`. Restaure a migration, rode outra vez e confirme que o token do seu teste sumiu. **Sem este passo o teste não conta.**

- [ ] **Step 4: Commit**

```bash
git add apps/server_core/tests/integration/integrations_reauth_test.go
git commit -m "test(integration): prove reauth keeps one row and disconnect releases the account"
```

---

### Task 5: Live drive — o gate da fatia

Nada acima prova que o operador consegue voltar. Unidade prova ramo, integração prova SQL; só esta task prova o caminho inteiro contra a conta ML real. **A fatia não fecha sem ela, e a F-A1 não faz o live drive dela antes desta.**

**Pré-condições, todas medidas antes de tocar em qualquer coisa:**

1. **O binário rodando tem que ser posterior aos commits desta fatia.** Compare:

```bash
docker ps --format '{{.Names}}\t{{.CreatedAt}}' | grep backend
git log -1 --format='%H %ci'
```

Se o container for mais velho que o último commit, **pare** — o drive rodaria contra código antigo e "não mostrou nada" leria como ausência de defeito. O entrypoint faz `go run` do fonte bind-montado (`docker/dev/backend-entrypoint.sh`), então basta recriar:

```bash
docker compose up -d --force-recreate backend
```

Confirme o código dentro do container antes de seguir:

```bash
MSYS_NO_PATHCONV=1 docker exec marketplace-central-backend-1 sh -c 'grep -c ensureProviderAccountUnlinked /workspace/apps/server_core/internal/modules/integrations/application/auth_flow_service.go'
```

Esperado: `1` ou mais.

2. **O túnel do OAuth tem que estar de pé.** O callback do ML volta pelo domínio reservado derivado de `MPC_OAUTH_REDIRECT_URI`:

```bash
docker compose --profile oauth up -d ngrok
```

O ngrok **não** sobe junto com o stack (`profiles: [oauth]` em `docker-compose.yml:73`). Confirme o processo e que o domínio responde:

```bash
MSYS_NO_PATHCONV=1 docker exec marketplace-central-ngrok-1 sh -c 'cat /proc/1/cmdline | tr "\0" " "'
```

3. **Estado do banco antes**, para poder comparar depois:

```bash
MSYS_NO_PATHCONV=1 docker exec marketplace-central-postgres-1 sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT installation_id, status, health_status, external_account_id FROM integration_installations;" -c "SELECT count(*) FILTER (WHERE is_active) AS ativas, count(*) AS total FROM integration_credentials;"'
```

- [ ] **Step 1: O operador revoga a autorização no ML**

**Esta ação é do operador, na conta real dele. Não a execute e não a peça sem confirmar que ele já leu o custo.** Caminho:

> Conta → Configurações → Segurança → *Opções de segurança e recuperação de senha* → Ferramentas de segurança → aplicações conectadas → remover o app.

Custo: enquanto revogado, nenhuma chamada ao ML funciona. O refresh token do ML é de **uso único e rotaciona a cada refresh**, então não há volta sem reautorizar. É exatamente por isso que as Tasks 1–4 vêm antes.

- [ ] **Step 2: Force o refresh e espere um tick**

O ticker roda a cada 5 minutos com janela de 10 (`internal/composition/root.go:683`, `background/refresh_ticker.go`). Traga a sessão para dentro da janela:

```sql
UPDATE integration_auth_sessions
   SET access_token_expires_at = now() - interval '1 minute',
       state = 'expiring',
       next_retry_at = NULL
 WHERE installation_id = '<inst-id-ml>';
```

Log esperado (não existia antes da fatia F-A1):

```
level=ERROR msg="integrations refresh ticker item failed" installation_id=<inst-id-ml> err="INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400 body=..."
```

- [ ] **Step 3: Confirme a persistência e a degradação**

```sql
SELECT state, refresh_failure_code, consecutive_failures, next_retry_at
  FROM integration_auth_sessions WHERE installation_id = '<inst-id-ml>';

SELECT status, health_status,
       connection_snapshot_json->>'next_action'  AS next_action,
       connection_snapshot_json->>'reauth_reason' AS reauth_reason
  FROM integration_installations WHERE installation_id = '<inst-id-ml>';
```

Esperado: `refresh_failed`, código preenchido, `consecutive_failures = 1`, `next_retry_at` uma hora à frente (cooldown terminal); e `requires_reauth`, `critical`, `reauth`, motivo preenchido.

- [ ] **Step 4: A tela mostra o estado e o botão**

Abra `/integracoes` **pelo domínio do túnel**, não por `localhost:5174` — o cookie de estado do OAuth precisa casar com o host do callback. Esperado: card "Contas conectadas" com badge "Precisa reautorizar", o motivo cru, e o **botão "Reautorizar"** (Task 1). Screenshot.

- [ ] **Step 5: O botão devolve a conta**

Clique em "Reautorizar", complete o consentimento no ML. Esperado: volta para `/integracoes` **sem** `?auth=failed`, o card volta a "Conectado", e:

```sql
SELECT count(*) FROM integration_installations;                       -- 1, não 2
SELECT count(*) FILTER (WHERE is_active) FROM integration_credentials; -- 1, não 2
SELECT state, consecutive_failures, next_retry_at
  FROM integration_auth_sessions WHERE installation_id = '<inst-id-ml>'; -- valid, 0, NULL
```

**Uma segunda installation aqui reprova a fatia inteira** — é o defeito original voltando.

- [ ] **Step 6: Controle positivo do caminho de recusa**

Sem isso, o Step 5 passar não prova que a checagem da Task 2 está ligada em produção — só que ninguém a acionou. Com a conta conectada, clique em **"Conectar"** (que cria installation nova) e autorize a **mesma** conta ML. Esperado:

- redirect para `/integracoes?auth=failed&reason=INTEGRATIONS_PROVIDER_ACCOUNT_ALREADY_LINKED`;
- o aviso da Task 3 visível na tela com esse código;
- e, no banco, **nenhuma credencial nova ativa e nenhuma sessão nova** — só a installation rascunho, que é comportamento declarado fora de escopo.

Screenshot. Depois apague o rascunho:

```sql
DELETE FROM integration_installations WHERE installation_id = '<inst-id-rascunho>';
```

- [ ] **Step 7: Grave a evidência**

Screenshots, trechos de log e saídas de SQL vão para `.mnfs/MIS-008-operacao-diaria/FA1b-reautorizacao/_evidence/`. Não escrito = não aconteceu.

---

### Task 6: Fechamento

- [ ] **Step 1: Registre a dívida residual**

Em `.mnfs/HARNESS-DEBTS.md` **não** — esta é dívida de produto. Registre no pack de evidência da fatia e no ledger da missão:

> **D-45 — `HandleCallback` continua sem transação.** `saveCredential` (`auth_flow_service.go:359`) e o upsert da sessão (`:364`) acontecem antes de `applyAuthResult` (`:375`) sem transação compartilhada. A Task 2 fechou a única causa alcançável pela tela (conta já vinculada), mas qualquer outra falha em `applyAuthResult` ainda deixa credencial ativa e sessão válida órfãs. Conserto genérico = porta de unidade-de-trabalho sobre os três repositórios.

E as duas irmãs, encontradas ao medir e deliberadamente não consertadas de passagem:

> **D-46 — `mapIntegrationError` casa por prefixo de string.** `transport/http_handler.go:53` usa `strings.HasPrefix(err.Error(), "INTEGRATIONS_")`, então erro embrulhado com `%w` escapa e vira 400 genérico. Mesma classe do defeito que esta fatia conserta.
>
> **D-47 — `fee_sync_scheduler.go:74` descarta erro com `_, _ =`.** Mesmo padrão que a Task 4 da F-A1 corrigiu no ticker de refresh.

- [ ] **Step 2: Desbloqueie a F-A1**

Edite `docs/superpowers/plans/2026-08-03-fa1-fa2-token-visivel-plan.md`, seção "Task 7", e remova o bloqueio (2), citando o SHA desta fatia. O bloqueio (1) — binário anterior à fatia — continua valendo e tem que ser re-medido no dia.

- [ ] **Step 3: Commit**

```bash
git add .mnfs/MIS-008-operacao-diaria docs/superpowers/plans/2026-08-03-fa1-fa2-token-visivel-plan.md
git commit -m "docs(F-A1b): evidence pack, residual debts D-45..D-47, unblock F-A1 Task 7"
```

---

## Fora de escopo desta fatia

- Transação no `HandleCallback` (D-45 acima).
- `mapIntegrationError` por prefixo de string (D-46).
- `fee_sync_scheduler.go:74` (D-47).
- Suporte a múltiplas contas ML. O produto declara suportar (o comentário em `IntegracoesPage.tsx:492` conta com isso), mas nada nesta fatia testa duas contas **distintas**. Se isso virar requisito, é fatia própria.
- Limpar installations rascunho abandonadas. `connect()` já reaproveita `pending_connection`/`draft`, então não acumulam; some-las é UX, não corretude.
