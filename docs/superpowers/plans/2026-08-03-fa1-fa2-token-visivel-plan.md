# F-A1 + F-A2 — Falha de token ML visível e lote que não aborta

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Uma falha de refresh do token do Mercado Livre passa a ser classificada, persistida, logada e mostrada na tela — hoje ela é invisível nas quatro camadas e o app pode estar morto com `/integracoes` verde.

**Architecture:** Nenhuma peça nova. O adapter ML passa a distinguir `invalid_grant` (terminal) de erro transitório em vez de devolver sempre `ErrRefreshProviderError`; `RefreshCredential` ganha um caminho de erro que grava `refresh_failure_code`/`consecutive_failures`/`next_retry_at` usando a `RefreshPolicy` que já existe e hoje tem zero chamadores, e degrada a instalação para `degraded`/`requires_reauth` — estados que `ProjectConnectionSnapshot` já projeta corretamente para `next_action: retry|reauth`. O ticker passa a logar e continuar o lote, copiando o padrão de `mutations/background/poller.go`. O frontend lê `health_status` + `connection` das installations que `useInstallation()` já busca.

**Tech Stack:** Go 1.x (`apps/server_core`), `log/slog`, `internal/platform/apierror`, React + TanStack Query (`apps/web`), `@marketplace-central/sdk-runtime`.

**Sem mudança de contrato.** `IntegrationInstallation.health_status`, `IntegrationInstallation.status: "requires_reauth" | "degraded"`, `IntegrationConnectionSnapshot.next_action: "reauth" | "retry"` e `.reauth_reason` já existem no OpenAPI e no SDK (`packages/sdk-runtime/src/index.ts:157`, `:148-156`, `:477-488`). Nenhuma task toca `contracts/api/marketplace-central.openapi.yaml` nem `packages/sdk-runtime`. Se alguma task fizer você querer mexer lá, pare: o campo já existe.

---

## Fake não é evidência

Ordem do operador, 2026-08-03: *"não quero mock, fallback, validação sempre real para realmente validar, não fazer teste apenas pra passar"*. Doutrina do repo (AGENTS.md): mocks provam comportamento de contrato, nunca integração viva.

Regra desta fatia, então:

- Os testes com `flowAdapter`/`flowInstallationStore`/`flowAuthWriter` das Tasks 1–4 provam **uma coisa só**: qual ramo o código escolhe dado um erro. Isso é legítimo e é tudo que eles provam. Nenhum deles é critério de aceite.
- O que os fakes **não podem** provar, e por isso tem prova própria: (a) que a linha realmente entra em `auth_sessions` e `integration_installations` — o `flowInstallationStore` é código meu escrevendo num `map`, o SQL nunca roda; (b) que `ListExpiringSessions` de fato pula a conta enquanto `next_retry_at` está no futuro — esse é um `WHERE` de verdade num Postgres de verdade; (c) que o operador vê a falha.
- (a) e (b) são a **Task 5**, na lane hermética de integração, contra Postgres real. (c) é a **Task 7**, no browser, contra a conta ML real.
- **Uma fatia não fecha com Task 1–4 verdes.** Só fecha com 5 e 7 verdes.
- Todo teste tem que ter sido visto **vermelho pelo motivo certo** antes de ficar verde. Cada task abaixo tem um step "rode e veja falhar" com a mensagem esperada. Se o vermelho vier com outra mensagem, o teste está medindo outra coisa — pare e conserte o teste, não a produção.
- Três dos testes abaixo são **controles negativos**: `anything else stays provider error` (Task 1), `SingleTransientFailureDoesNotDegrade` (Task 3), `RunOnceReturnsErrorWhenListingFails` (Task 4). Eles têm que passar desde o começo. Se um deles falhar junto com os outros, o fake está quebrado e nenhum verde daquela task vale.
- O único `fallback` do plano é `logger == nil → slog.Default()` na Task 4. Não é fallback de dado: é nil-guard de dependência, cópia literal de `mutations/background/poller.go:26`. Nenhum valor de negócio tem default em lugar nenhum deste plano — desconhecido continua desconhecido (ADR-17).

**Ordem das tasks:** 1 adapter · 2 persiste sessão · 3 degrada instalação · 4 ticker · 5 integração real · 6 card · 7 drive ao vivo.

## Anti-redundância (§12.4 do spec)

Cada peça abaixo já existe. Nenhuma task cria substituto.

| Preciso de | Já existe em | Por que serve |
|---|---|---|
| Classificar erro de refresh | `internal/modules/integrations/domain/refresh_policy.go:52` `ClassifyRefreshError` | Correto; só não tem chamador de produção |
| Backoff exponencial | `domain/refresh_policy.go:32` `BackoffDuration(attempt int)` | Base 30s, teto 15min, clampado |
| Gravar falha na sessão | `application/auth_service.go:12-31` `UpsertAuthSessionInput` | Já aceita `RefreshFailureCode`, `ConsecutiveFailures`, `NextRetryAt` |
| Respeitar backoff na varredura | `adapters/postgres/auth_session_repo.go:87` | `AND (next_retry_at IS NULL OR next_retry_at <= now())` já está no WHERE |
| Projetar estado degradado para a tela | `domain/connection_snapshot.go:102-127` `ProjectConnectionSnapshot` | `RequiresReauth` → `needs_reauth`+`reauth`; `Degraded` → `degraded`+`retry` |
| Escrever o snapshot degradado | `application/auth_flow_service.go:492-514` (idiom do `Disconnect`) | Monta `Installation` com novo status, projeta, chama `ApplyConnectionSnapshot` |
| Log estruturado + lote que continua | `internal/modules/mutations/background/poller.go:26,66,71` | `logger *slog.Logger` injetado com fallback `slog.Default()`; erro por item loga e o loop segue |
| Installations no frontend | `apps/web/src/app/InstallationContext.tsx:30-34` `useInstallation()` | Já busca `listIntegrationInstallations()` e expõe `installations` |
| Card de saúde na `/integracoes` | `apps/web/src/pages/integracoes/SyncHealthCard.tsx` | Padrão de tom/badge/`formatRelative` a espelhar |

---

## File Structure

**Backend**
- `internal/modules/integrations/adapters/mercadolivre/auth_adapter.go` — mapeia a resposta HTTP de erro do ML para a sentinela de domínio certa. Vocabulário de provider (`invalid_grant`) fica aqui, nunca no domínio (ADR-C4).
- `internal/modules/integrations/application/auth_flow_service.go` — caminho de erro de `RefreshCredential`: persiste falha, degrada instalação, devolve o erro original.
- `internal/modules/integrations/background/refresh_ticker.go` — log-and-continue no lote, log do erro do tick.
- `internal/composition/root.go:683` — passa o logger ao ticker.

**Frontend**
- `apps/web/src/pages/integracoes/ConnectionHealthCard.tsx` — **novo**; card que mostra o estado de conexão de cada installation.
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx` — monta o card.

**Testes**
- `internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go` — unidade, fakes
- `internal/modules/integrations/application/auth_flow_service_test.go` — unidade, fakes
- `internal/modules/integrations/background/refresh_ticker_test.go` — unidade, fakes
- `apps/server_core/tests/integration/integrations_refresh_failure_test.go` — **novo**, Postgres real, lane `integration`
- `apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx` — unidade de render

---

## Comandos

Go, **sempre de dentro de `apps/server_core`**, com `GOCACHE` absoluto (doutrina do harness):

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/... -run TestX -v
```

Frontend, da raiz do repo:

```bash
npm run test --workspace apps/web -- ConnectionHealthCard
```

---

### Task 1: Adapter ML distingue token morto de erro transitório

Hoje `RefreshToken` devolve `ErrRefreshProviderError` para **todo** status >= 400 (`auth_adapter.go:243`). `ClassifyRefreshError` classifica isso como **transitório**, então um refresh token revogado seria retentado para sempre e nunca chegaria a `requires_reauth`. Ligar a política sem consertar isso deixaria a política errada, não só desligada.

O ML devolve, no corpo do erro do token endpoint, `{"error":"invalid_grant", ...}` quando o refresh token é inválido/revogado, e HTTP 429 quando está limitando taxa.

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:242-244`
- Test: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go`

- [ ] **Step 1: Escreva o teste que falha**

`auth_adapter_test.go` já existe e já é `package mercadolivre`. Acrescente o código abaixo ao fim dele e garanta que `context`, `errors`, `io`, `net/http`, `strings`, `testing` e o pacote `domain` estejam no bloco de imports — sem duplicar os que já estiverem lá.

```go
type refreshRoundTripper struct {
	status int
	body   string
}

func (r refreshRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: r.status,
		Body:       io.NopCloser(strings.NewReader(r.body)),
		Header:     make(http.Header),
	}, nil
}

func TestRefreshTokenClassifiesProviderErrors(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{
			name:   "invalid_grant is terminal",
			status: http.StatusBadRequest,
			body:   `{"error":"invalid_grant","message":"invalid refresh token"}`,
			want:   domain.ErrRefreshTokenInvalid,
		},
		{
			name:   "rate limit is its own sentinel",
			status: http.StatusTooManyRequests,
			body:   `{"error":"too_many_requests"}`,
			want:   domain.ErrRefreshRateLimited,
		},
		{
			name:   "anything else stays provider error",
			status: http.StatusInternalServerError,
			body:   `{"error":"internal_error"}`,
			want:   domain.ErrRefreshProviderError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			adapter := &Adapter{cfg: Config{
				ClientID:     "cid",
				ClientSecret: "secret",
				TokenURL:     "https://example.invalid/oauth/token",
				HTTPClient:   &http.Client{Transport: refreshRoundTripper{status: tc.status, body: tc.body}},
			}}

			_, err := adapter.RefreshToken(context.Background(), "rt-1")
			if !errors.Is(err, tc.want) {
				t.Fatalf("RefreshToken err = %v, want errors.Is(_, %v)", err, tc.want)
			}
		})
	}
}
```

`Config` (`auth_adapter.go:21-28`) e `Adapter{cfg Config}` (`:30-32`) têm exatamente essa forma — verificado. `TokenURL` aponta para um host inválido de propósito: o `RoundTripper` intercepta antes de qualquer rede sair. Nenhum teste desta fatia toca a internet.

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/adapters/mercadolivre/ -run TestRefreshTokenClassifiesProviderErrors -v
```

Esperado: FAIL nos dois primeiros subtestes com `RefreshToken err = INTEGRATIONS_REFRESH_PROVIDER_ERROR: status=400 body=... , want errors.Is(_, INTEGRATIONS_REFRESH_TOKEN_INVALID)`. O terceiro subteste passa — ele é o **controle negativo**: prova que o teste não está aprovando tudo.

- [ ] **Step 3: Implemente**

Em `auth_adapter.go`, substitua o bloco de erro de `RefreshToken` (linhas 242-244):

```go
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: status=%d body=%s", domain.ErrRefreshProviderError, resp.StatusCode, readProviderErrorBody(resp))
	}
```

por:

```go
	if resp.StatusCode >= 400 {
		body := readProviderErrorBody(resp)
		return nil, fmt.Errorf("%w: status=%d body=%s", classifyRefreshHTTPError(resp.StatusCode, body), resp.StatusCode, body)
	}
```

e acrescente, logo abaixo de `readProviderErrorBody` (que termina na linha 307):

```go
// classifyRefreshHTTPError traduz a resposta de erro do endpoint de token do ML
// para a sentinela de domínio correspondente. O vocabulário do provider
// ("invalid_grant") vive aqui e só aqui: o domínio não conhece nome de erro de
// marketplace nenhum (ADR-C4).
//
// Sem essa tradução todo erro >= 400 virava ErrRefreshProviderError, que
// ClassifyRefreshError (domain/refresh_policy.go:52) considera TRANSITÓRIO —
// um refresh token revogado seria retentado para sempre e a conta nunca
// chegaria a requires_reauth.
func classifyRefreshHTTPError(status int, body string) error {
	// O ML responde 400 com {"error":"invalid_grant"} para refresh token
	// inválido, revogado ou já usado. Nenhum retry conserta: só reautorização.
	if strings.Contains(body, "invalid_grant") {
		return domain.ErrRefreshTokenInvalid
	}
	if status == http.StatusTooManyRequests {
		return domain.ErrRefreshRateLimited
	}
	return domain.ErrRefreshProviderError
}
```

`strings`, `net/http` e `fmt` já estão importados no arquivo.

- [ ] **Step 4: Rode e veja passar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/adapters/mercadolivre/ -run TestRefreshTokenClassifiesProviderErrors -v
```

Esperado: PASS nos três subtestes.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go
git commit -m "fix(integrations): classify ML refresh errors instead of collapsing to provider error

Every HTTP >= 400 from the ML token endpoint returned ErrRefreshProviderError,
which ClassifyRefreshError treats as transient. A revoked refresh token would
have been retried forever and never reached requires_reauth. invalid_grant now
maps to ErrRefreshTokenInvalid (terminal) and 429 to ErrRefreshRateLimited."
```

---

### Task 2: `RefreshCredential` persiste a falha

Hoje todo caminho de erro é `return AuthStatus{}, err` — nada é escrito. `auth_flow_service.go:440-441` é o caso que importa: o adapter falhou e ninguém registra.

Esta task grava a falha na `auth_sessions`. A degradação da instalação é a Task 3 — separada de propósito, para que cada teste isole uma escrita.

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go:440-442`
- Test: `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`

- [ ] **Step 1: Prepare o fake do adapter para falhar**

`flowAdapter` (`auth_flow_service_test.go:139-148`) hoje só sabe ter sucesso. Adicione um campo de erro ao struct:

```go
type flowAdapter struct {
	providerCode  string
	startInput    StartAuthorizeAdapterInput
	callbackInput HandleCallbackAdapterInput
	refreshInput  RefreshCredentialAdapterInput
	callback      CredentialPayload
	apiKey        CredentialPayload
	refresh       CredentialPayload
	refreshErr    error
	refreshCalls  int
}
```

e, no método `Refresh` do fake, devolva-o antes do payload de sucesso:

```go
	if s.refreshErr != nil {
		return CredentialPayload{}, s.refreshErr
	}
```

Nenhum teste existente seta `refreshErr`, então o zero-value `nil` preserva o comportamento atual de todos eles.

- [ ] **Step 2: Escreva o helper de setup**

O setup de `TestRefreshCredentialRotatesAndResetsFailures` (`:638-707`) é inline e vai ser reusado por cinco testes novos. Extraia-o. Acrescente ao arquivo de teste:

```go
type refreshFixture struct {
	svc           *AuthFlowService
	installations *flowInstallationStore
	authSessions  *flowAuthWriter
	adapter       *flowAdapter
	now           time.Time
}

// newRefreshFixture monta o mesmo cenário de
// TestRefreshCredentialRotatesAndResetsFailures (auth_flow_service_test.go:638):
// uma installation ML conectada com credencial ativa, sessão de auth encontrada
// e encryptor devolvendo um refresh token. priorFailures é o número de falhas
// consecutivas JÁ registradas na sessão; refreshErr, quando não-nil, faz o
// adapter falhar em vez de rotacionar.
func newRefreshFixture(t *testing.T, priorFailures int, refreshErr error) refreshFixture {
	t.Helper()

	now := time.Unix(1700, 0).UTC()
	expiresAt := time.Unix(2600, 0).UTC()
	installations := &flowInstallationStore{installations: map[string]domain.Installation{
		"inst-ml": {
			InstallationID:     "inst-ml",
			ProviderCode:       "mercado_livre",
			Status:             domain.InstallationStatusConnected,
			HealthStatus:       domain.HealthStatusHealthy,
			ExternalAccountID:  "seller-1",
			ActiveCredentialID: "cred-active",
		},
	}}
	credentials := &flowCredentialRotator{
		activeFound: true,
		activeCredential: domain.Credential{
			CredentialID:     "cred-active",
			InstallationID:   "inst-ml",
			SecretType:       "oauth2",
			EncryptedPayload: []byte("active-ciphertext"),
			EncryptionKeyID:  "key-1",
			IsActive:         true,
		},
	}
	authSessions := &flowAuthWriter{
		sessionFound: true,
		session: domain.AuthSession{
			AuthSessionID:        "auth_inst-ml",
			InstallationID:       "inst-ml",
			ProviderAccountID:    "seller-1",
			State:                domain.AuthStateValid,
			ConsecutiveFailures:  priorFailures,
			AccessTokenExpiresAt: ptrTime(now.Add(time.Minute)),
		},
	}
	encryptor := &flowEncryptor{
		decryptedPayload: map[string]any{
			"type":                "oauth2",
			"access_token":        "old-access",
			"refresh_token":       "old-refresh",
			"provider_account_id": "seller-1",
		},
		decryptKeyID: "key-1",
	}
	adapter := &flowAdapter{
		providerCode: "mercado_livre",
		refreshErr:   refreshErr,
		refresh: CredentialPayload{
			SecretType:        "oauth2",
			AccessToken:       "new-access",
			RefreshToken:      "new-refresh",
			ProviderAccountID: "seller-1",
			ExpiresAt:         &expiresAt,
		},
	}
	svc := mustNewAuthFlowService(t, AuthFlowConfig{
		TenantID:        "tenant_default",
		Installations:   installations,
		Credentials:     credentials,
		AuthSessions:    authSessions,
		OAuthStates:     &securityOAuthStateStore{},
		OAuthStateCodec: roundTripSecurityStateCodec{payloadsByState: map[string]OAuthStatePayload{}},
		Encryptor:       encryptor,
		Clock:           fixedAuthFlowClock{now: now},
		Adapters:        []MarketplaceAuthAdapter{adapter},
	})

	return refreshFixture{
		svc:           svc,
		installations: installations,
		authSessions:  authSessions,
		adapter:       adapter,
		now:           now,
	}
}
```

- [ ] **Step 3: Escreva os testes que falham**

```go
func TestRefreshCredentialPersistsTerminalFailure(t *testing.T) {
	t.Parallel()

	fx := newRefreshFixture(t, 0, domain.ErrRefreshTokenInvalid)

	_, err := fx.svc.RefreshCredential(context.Background(), RefreshCredentialInput{InstallationID: "inst-ml"})
	if !errors.Is(err, domain.ErrRefreshTokenInvalid) {
		t.Fatalf("RefreshCredential err = %v, want ErrRefreshTokenInvalid", err)
	}

	if len(fx.authSessions.sessions) != 1 {
		t.Fatalf("upserts = %d, want 1: a falha de refresh continua invisível", len(fx.authSessions.sessions))
	}
	got := fx.authSessions.sessions[0]
	if got.State != domain.AuthStateRefreshFailed {
		t.Fatalf("State = %q, want %q", got.State, domain.AuthStateRefreshFailed)
	}
	if got.RefreshFailureCode != domain.ErrRefreshTokenInvalid.Error() {
		t.Fatalf("RefreshFailureCode = %q, want %q", got.RefreshFailureCode, domain.ErrRefreshTokenInvalid.Error())
	}
	if got.ConsecutiveFailures != 1 {
		t.Fatalf("ConsecutiveFailures = %d, want 1", got.ConsecutiveFailures)
	}
	if got.NextRetryAt == nil {
		t.Fatal("NextRetryAt = nil: sem next_retry_at o ticker retenta em todo tick")
	}
	// Terminal usa o cooldown inteiro da política, não o backoff curto.
	wantAt := fx.now.Add(domain.DefaultRefreshPolicy().CooldownAfterTerminal)
	if !got.NextRetryAt.Equal(wantAt) {
		t.Fatalf("NextRetryAt = %v, want %v", got.NextRetryAt, wantAt)
	}
}

func TestRefreshCredentialBacksOffTransientFailure(t *testing.T) {
	t.Parallel()

	// Controle: erro transitório NÃO usa o cooldown terminal; usa o backoff
	// exponencial a partir do número de falhas já registrado na sessão.
	fx := newRefreshFixture(t, 2, domain.ErrRefreshProviderError)

	_, err := fx.svc.RefreshCredential(context.Background(), RefreshCredentialInput{InstallationID: "inst-ml"})
	if !errors.Is(err, domain.ErrRefreshProviderError) {
		t.Fatalf("RefreshCredential err = %v, want ErrRefreshProviderError", err)
	}

	if len(fx.authSessions.sessions) != 1 {
		t.Fatalf("upserts = %d, want 1", len(fx.authSessions.sessions))
	}
	got := fx.authSessions.sessions[0]
	if got.ConsecutiveFailures != 3 {
		t.Fatalf("ConsecutiveFailures = %d, want 3 (2 anteriores + 1)", got.ConsecutiveFailures)
	}
	wantAt := fx.now.Add(domain.DefaultRefreshPolicy().BackoffDuration(2))
	if got.NextRetryAt == nil || !got.NextRetryAt.Equal(wantAt) {
		t.Fatalf("NextRetryAt = %v, want %v", got.NextRetryAt, wantAt)
	}
}
```

`errors` já está importado no arquivo de teste; se não estiver, adicione.

- [ ] **Step 4: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/application/ -run 'TestRefreshCredential(PersistsTerminalFailure|BacksOffTransientFailure)' -v
```

Esperado: FAIL com `upserts = 0, want 1: a falha de refresh continua invisível`. Se o vermelho vier de compilação ou de `flowAdapter`, o fixture está errado — conserte o fixture antes de tocar em produção.

- [ ] **Step 5: Implemente**

Em `auth_flow_service.go`, substitua o bloco de erro do adapter (linhas 440-442):

```go
	if err != nil {
		return AuthStatus{}, err
	}
```

por:

```go
	if err != nil {
		// Sem esta escrita a falha some: o ticker descarta o erro, a sessão
		// continua marcada como válida e /integracoes segue verde com o token
		// morto. ListExpiringSessions (adapters/postgres/auth_session_repo.go:87)
		// já filtra por next_retry_at, então gravar aqui é o que ativa o backoff
		// que a política sempre soube calcular e nunca teve quem chamasse.
		if recordErr := s.recordRefreshFailure(ctx, inst, session, err); recordErr != nil {
			return AuthStatus{}, recordErr
		}
		return AuthStatus{}, err
	}
```

E acrescente, logo abaixo de `RefreshCredential` (depois da linha 471):

```go
// recordRefreshFailure persiste a falha de refresh na sessão de auth. O erro
// original é devolvido ao chamador intacto — esta função não o substitui, só o
// torna observável.
func (s *AuthFlowService) recordRefreshFailure(
	ctx context.Context,
	inst domain.Installation,
	session domain.AuthSession,
	cause error,
) error {
	policy := domain.DefaultRefreshPolicy()
	failures := session.ConsecutiveFailures + 1

	// A classe decide o intervalo: erro terminal (refresh token revogado) não
	// tem retry útil, então espera o cooldown inteiro; transitório sobe pelo
	// backoff exponencial a partir das falhas JÁ registradas — passar `failures`
	// aqui puliria o primeiro degrau da escada.
	var retryIn time.Duration
	if domain.ClassifyRefreshError(cause) == domain.ErrorClassTerminal {
		retryIn = policy.CooldownAfterTerminal
	} else {
		retryIn = policy.BackoffDuration(session.ConsecutiveFailures)
	}

	nextRetryAt := s.clock.Now().UTC().Add(retryIn)
	_, err := s.authSessions.Upsert(ctx, UpsertAuthSessionInput{
		AuthSessionID:        firstNonEmpty(session.AuthSessionID, fmt.Sprintf("auth_%s", inst.InstallationID)),
		InstallationID:       inst.InstallationID,
		ProviderAccountID:    firstNonEmpty(session.ProviderAccountID, inst.ExternalAccountID),
		State:                domain.AuthStateRefreshFailed,
		AccessTokenExpiresAt: session.AccessTokenExpiresAt,
		LastVerifiedAt:       session.LastVerifiedAt,
		RefreshFailureCode:   cause.Error(),
		ConsecutiveFailures:  failures,
		NextRetryAt:          &nextRetryAt,
	})
	return err
}
```

> `RefreshFailureCode: cause.Error()` grava a string completa do erro embrulhado (`INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400 body=...`). Isso é intencional: o corpo do ML é o único diagnóstico que temos e a coluna é `text`. Se o teste da Task 2 exigir só a sentinela, ele está errado — ajuste o teste para `strings.HasPrefix`.

No fixture da Task 2, `adapter.refreshErr` recebe a sentinela crua, sem embrulho, então `== domain.ErrRefreshTokenInvalid.Error()` vale nos testes de unidade. Ao vivo (Task 7) o valor gravado virá com `: status=400 body=...` junto — isso é esperado e desejado.

- [ ] **Step 6: Rode e veja passar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/application/ -v
```

Esperado: PASS, incluindo `TestRefreshCredentialRotatesAndResetsFailures` (`:638`) — o caminho de sucesso continua zerando `ConsecutiveFailures` e `NextRetryAt`. Se ele quebrar, você mexeu no caminho de sucesso; desfaça.

- [ ] **Step 7: Commit**

```bash
git add apps/server_core/internal/modules/integrations/application/auth_flow_service.go apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go
git commit -m "feat(integrations): persist refresh failure with policy-driven backoff

RefreshCredential returned bare on adapter error, so no failure was ever
recorded and the session stayed 'valid' with a dead token. It now writes
refresh_failure_code, consecutive_failures and next_retry_at using the
RefreshPolicy that had zero production callers. ListExpiringSessions already
filters on next_retry_at, so the backoff takes effect with no query change."
```

---

### Task 3: Instalação degrada e a razão chega ao snapshot

Gravar a sessão não muda nada na tela: a tela lê `Installation.HealthStatus` e `Installation.ConnectionSnapshot`, e `HealthStatusCritical` nunca é setado em código de produção. Esta task liga as duas pontas.

Regra: erro **terminal** (token revogado) → `InstallationStatusRequiresReauth` + `HealthStatusCritical`; erro **transitório** que passou de `MaxConsecutiveFailures` → `InstallationStatusDegraded` + `HealthStatusWarning`; transitório abaixo do limite → nenhuma mudança de estado (uma falha isolada de rede não deve pintar a conta de vermelho).

`ProjectConnectionSnapshot` (`domain/connection_snapshot.go:102-127`) já mapeia esses dois status para `next_action: reauth` e `next_action: retry` com `ReauthReason` — não escreva projeção nova.

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go` (dentro de `recordRefreshFailure`, criada na Task 2)
- Test: `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`

- [ ] **Step 1: Escreva o teste que falha**

O fake `flowInstallationStore` já acumula todo snapshot aplicado em `connectionSnapshots` (`auth_flow_service_test.go:20,51`) — não adicione campo nem helper novo, leia esse slice.

```go
func TestRefreshCredentialTerminalFailureRequiresReauth(t *testing.T) {
	t.Parallel()

	fx := newRefreshFixture(t, 0, domain.ErrRefreshTokenInvalid)

	_, _ = fx.svc.RefreshCredential(context.Background(), RefreshCredentialInput{InstallationID: "inst-ml"})

	if len(fx.installations.connectionSnapshots) != 1 {
		t.Fatalf("snapshots aplicados = %d, want 1: a tela continua verde com token morto", len(fx.installations.connectionSnapshots))
	}
	snap := fx.installations.connectionSnapshots[0]
	if snap.State != domain.ConnectionStateNeedsReauth {
		t.Fatalf("State = %q, want %q", snap.State, domain.ConnectionStateNeedsReauth)
	}
	if snap.Health != domain.HealthStatusCritical {
		t.Fatalf("Health = %q, want %q", snap.Health, domain.HealthStatusCritical)
	}
	if snap.NextAction != domain.ConnectionNextActionReauth {
		t.Fatalf("NextAction = %q, want %q", snap.NextAction, domain.ConnectionNextActionReauth)
	}
	if snap.ReauthReason == "" {
		t.Fatal("ReauthReason vazio: a tela não teria o que dizer ao operador")
	}
}

func TestRefreshCredentialSingleTransientFailureDoesNotDegrade(t *testing.T) {
	t.Parallel()

	// Controle negativo: uma falha transitória isolada não pode pintar a conta.
	fx := newRefreshFixture(t, 0, domain.ErrRefreshProviderError)

	_, _ = fx.svc.RefreshCredential(context.Background(), RefreshCredentialInput{InstallationID: "inst-ml"})

	if len(fx.installations.connectionSnapshots) != 0 {
		t.Fatalf("snapshots aplicados = %d, want 0 numa única falha transitória", len(fx.installations.connectionSnapshots))
	}
}

func TestRefreshCredentialTransientFailuresOverThresholdDegrade(t *testing.T) {
	t.Parallel()

	fx := newRefreshFixture(t, domain.DefaultRefreshPolicy().MaxConsecutiveFailures, domain.ErrRefreshProviderError)

	_, _ = fx.svc.RefreshCredential(context.Background(), RefreshCredentialInput{InstallationID: "inst-ml"})

	if len(fx.installations.connectionSnapshots) != 1 {
		t.Fatalf("snapshots aplicados = %d, want 1 acima do limite de falhas", len(fx.installations.connectionSnapshots))
	}
	snap := fx.installations.connectionSnapshots[0]
	if snap.State != domain.ConnectionStateDegraded {
		t.Fatalf("State = %q, want %q", snap.State, domain.ConnectionStateDegraded)
	}
	if snap.NextAction != domain.ConnectionNextActionRetry {
		t.Fatalf("NextAction = %q, want %q", snap.NextAction, domain.ConnectionNextActionRetry)
	}
}
```

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/application/ -run 'TestRefreshCredential(TerminalFailureRequiresReauth|SingleTransientFailureDoesNotDegrade|TransientFailuresOverThresholdDegrade)' -v
```

Esperado: FAIL em `TerminalFailureRequiresReauth` e `TransientFailuresOverThresholdDegrade` com `ApplyConnectionSnapshot nunca chamado`. `SingleTransientFailureDoesNotDegrade` **passa** desde já — é o controle negativo que prova que os outros dois falham por causa da escrita ausente, não por causa do fake.

- [ ] **Step 3: Implemente**

Em `recordRefreshFailure`, logo antes do `return err` final, substitua

```go
	return err
```

por

```go
	if err != nil {
		return err
	}

	return s.degradeAfterRefreshFailure(ctx, inst, cause, failures, policy)
}

// degradeAfterRefreshFailure move a instalação para o estado que a tela sabe
// mostrar. Nada aqui projeta snapshot na mão: ProjectConnectionSnapshot
// (domain/connection_snapshot.go:102-127) já mapeia RequiresReauth -> needs_reauth
// + next_action reauth e Degraded -> degraded + next_action retry. O idiom de
// montar uma Installation com o novo status e aplicá-la é o mesmo de Disconnect
// (auth_flow_service.go:492-514).
func (s *AuthFlowService) degradeAfterRefreshFailure(
	ctx context.Context,
	inst domain.Installation,
	cause error,
	failures int,
	policy domain.RefreshPolicy,
) error {
	var status domain.InstallationStatus
	var health domain.HealthStatus

	switch {
	case domain.ClassifyRefreshError(cause) == domain.ErrorClassTerminal:
		// Refresh token revogado: nenhum retry conserta, só reautorização.
		status = domain.InstallationStatusRequiresReauth
		health = domain.HealthStatusCritical
	case failures > policy.MaxConsecutiveFailures:
		// Transitório persistente. Ainda pode voltar sozinho, então degraded
		// (next_action retry) e não requires_reauth.
		status = domain.InstallationStatusDegraded
		health = domain.HealthStatusWarning
	default:
		// Falha transitória isolada não muda o estado da conta. Ela já está
		// registrada na sessão (consecutive_failures) e já tem next_retry_at;
		// pintar a tela de amarelo a cada soluço de rede treinaria o operador
		// a ignorar o aviso justamente quando ele for verdadeiro.
		return nil
	}

	degraded := inst
	degraded.Status = status
	degraded.HealthStatus = health
	degraded.UpdatedAt = s.clock.Now().UTC()

	snapshot := domain.ProjectConnectionSnapshot(
		degraded,
		inferConnectionAuthStrategy(inst.ConnectionSnapshot),
		inst.ConnectionSnapshot.ExpiresAt,
		cause.Error(),
	)
	return s.installations.ApplyConnectionSnapshot(ctx, inst.InstallationID, snapshot, "")
```

> `ApplyConnectionSnapshot(..., "")` com credential ID vazio é o mesmo argumento que `Disconnect` passa em `:512` — não estamos trocando a credencial, só o estado da conexão.

- [ ] **Step 4: Rode e veja passar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/... -v
```

Esperado: PASS no pacote inteiro.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/application/auth_flow_service.go apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go
git commit -m "feat(integrations): degrade installation state on refresh failure

HealthStatusCritical had no production call-site, so a dead ML token left the
installation reading 'connected/healthy'. A terminal refresh error now moves it
to requires_reauth/critical and a run of transient errors past the policy
threshold to degraded/warning. ProjectConnectionSnapshot already turns both into
next_action reauth/retry with a reason, so no new projection was written."
```

---

### Task 4: O ticker deixa de abortar o lote (F-A2)

`RunOnce` faz `return err` no primeiro item que falha (`refresh_ticker.go:46-47`). Com N contas, a primeira falha impede as N-1 seguintes de sequer tentarem. Hoje o repo tem uma conta só, então o defeito não aparece — e por isso ele precisa ser fechado antes de a segunda existir.

O padrão certo já está no repo: `mutations/background/poller.go:70-72` loga o erro por item e o loop segue. `fee_sync_scheduler.go:74` também continua, mas descartando o erro com `_, _ =` — esse **não** é o padrão a copiar.

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/background/refresh_ticker.go`
- Modify: `apps/server_core/internal/composition/root.go:683`
- Test: `apps/server_core/internal/modules/integrations/background/refresh_ticker_test.go`

- [ ] **Step 1: Escreva o teste que falha**

Os fakes já existem em `refresh_ticker_test.go:12-27`. Estenda-os com dois campos — não crie fakes novos. `refreshSessionStore` ganha `listErr`, `refreshFlow` ganha `failFor` (`inputs` já registra as chamadas):

```go
type refreshSessionStore struct {
	items   []domain.AuthSession
	listErr error
}

func (s refreshSessionStore) ListExpiringSessions(context.Context, time.Duration) ([]domain.AuthSession, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return append([]domain.AuthSession(nil), s.items...), nil
}

type refreshFlow struct {
	inputs  []application.RefreshCredentialInput
	failFor map[string]error
}

func (s *refreshFlow) RefreshCredential(ctx context.Context, input application.RefreshCredentialInput) (application.AuthStatus, error) {
	s.inputs = append(s.inputs, input)
	if err, ok := s.failFor[input.InstallationID]; ok {
		return application.AuthStatus{}, err
	}
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusConnected}, nil
}
```

E os testes:

```go
func TestRunOnceContinuesAfterItemFailure(t *testing.T) {
	t.Parallel()

	sessions := refreshSessionStore{
		items: []domain.AuthSession{
			{InstallationID: "inst-a"},
			{InstallationID: "inst-b"},
			{InstallationID: "inst-c"},
		},
	}
	flow := &refreshFlow{failFor: map[string]error{"inst-a": errors.New("boom")}}

	ticker := NewRefreshTicker(sessions, flow, time.Minute, nil)
	if err := ticker.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce err = %v, want nil (falha de item não é falha do lote)", err)
	}

	if len(flow.inputs) != 3 {
		t.Fatalf("chamadas = %#v, want 3 (a falha de inst-a matou o resto do lote)", flow.inputs)
	}
}

func TestRunOnceReturnsErrorWhenListingFails(t *testing.T) {
	t.Parallel()

	// Controle negativo: falhar em LISTAR é falha do lote inteiro e continua
	// subindo. Sem este teste, "RunOnce sempre retorna nil" passaria vacuoso.
	sessions := refreshSessionStore{listErr: errors.New("db down")}
	ticker := NewRefreshTicker(sessions, &refreshFlow{}, time.Minute, nil)

	if err := ticker.RunOnce(context.Background()); err == nil {
		t.Fatal("RunOnce err = nil, want erro de listagem")
	}
}
```

Adicione `"errors"` aos imports do arquivo de teste.

- [ ] **Step 2: Rode e veja falhar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/background/ -run TestRunOnce -v
```

Esperado: FAIL de compilação primeiro (`NewRefreshTicker` ainda tem 3 parâmetros), e depois `chamadas = [inst-a], want 3`. Corrija a assinatura no Step 3 — os dois testes têm que rodar juntos.

- [ ] **Step 3: Implemente**

Substitua o corpo inteiro de `refresh_ticker.go` a partir da declaração do struct:

```go
type RefreshTicker struct {
	sessions      expiringSessionLister
	flow          credentialRefresher
	logger        *slog.Logger
	interval      time.Duration
	expiresWithin time.Duration
	stop          chan struct{}
}

// NewRefreshTicker aceita logger nil e cai em slog.Default(), mesmo contrato de
// mutations/background/poller.go:26 — chamador de teste não precisa montar um
// logger só para exercitar o loop.
func NewRefreshTicker(sessions expiringSessionLister, flow credentialRefresher, interval time.Duration, logger *slog.Logger) *RefreshTicker {
	if logger == nil {
		logger = slog.Default()
	}
	return &RefreshTicker{
		sessions:      sessions,
		flow:          flow,
		logger:        logger,
		interval:      interval,
		expiresWithin: 10 * time.Minute,
		stop:          make(chan struct{}),
	}
}

func (t *RefreshTicker) RunOnce(ctx context.Context) error {
	sessions, err := t.sessions.ListExpiringSessions(ctx, t.expiresWithin)
	if err != nil {
		// Falhar em LISTAR é falha do lote: não há o que iterar. Sobe.
		return err
	}

	for _, session := range sessions {
		// Falha de UM item não é falha do lote. Abortar aqui fazia a primeira
		// conta com token quebrado impedir todas as seguintes de tentarem —
		// invisível com uma conta só, fatal com duas. Mesmo padrão de
		// mutations/background/poller.go:70-72.
		if _, err := t.flow.RefreshCredential(ctx, application.RefreshCredentialInput{
			InstallationID: session.InstallationID,
		}); err != nil {
			t.logger.Error("integrations refresh ticker item failed",
				"installation_id", session.InstallationID,
				"err", err,
			)
		}
	}
	return nil
}

func (t *RefreshTicker) Start(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stop:
			return
		case <-ticker.C:
			// O erro era descartado com `_ =`: uma listagem falhando em todo
			// tick não deixava rastro nenhum (F-A1).
			if err := t.RunOnce(ctx); err != nil {
				t.logger.Error("integrations refresh ticker pass failed", "err", err)
			}
		}
	}
}
```

Acrescente `"log/slog"` ao bloco de imports do arquivo.

- [ ] **Step 4: Atualize a fiação**

Em `internal/composition/root.go:683`, troque

```go
	go integrationsbg.NewRefreshTicker(authSessionRepo, authFlowSvc, 5*time.Minute).Start(context.Background())
```

por

```go
	go integrationsbg.NewRefreshTicker(authSessionRepo, authFlowSvc, 5*time.Minute, slog.Default()).Start(context.Background())
```

`log/slog` já está importado em `root.go`. Confirme com `grep -n '"log/slog"' internal/composition/root.go`; se não estiver, adicione.

- [ ] **Step 5: Rode e veja passar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./... && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/integrations/... -v
```

Esperado: build limpo e PASS, incluindo `TestRefreshTickerUsesListExpiringSessions` (`refresh_ticker_test.go:29`) — você vai precisar adicionar o quarto argumento `nil` na chamada dele.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/integrations/background/refresh_ticker.go apps/server_core/internal/modules/integrations/background/refresh_ticker_test.go apps/server_core/internal/composition/root.go
git commit -m "fix(integrations): refresh ticker logs and continues instead of aborting the batch

RunOnce returned on the first item error, so one broken account blocked every
other account's refresh, and Start discarded the pass error entirely with \`_ =\`.
Both now log through an injected *slog.Logger, matching
mutations/background/poller.go. Listing failure still aborts the pass: there is
nothing to iterate."
```

---

### Task 5: Prova real contra Postgres — o backoff existe de verdade

As Tasks 2–4 provaram ramificação contra fakes. Duas coisas continuam **não provadas** e são justamente as que sustentam a fatia inteira:

1. Que `refresh_failure_code`/`consecutive_failures`/`next_retry_at` sobrevivem à ida ao banco. O `flowAuthWriter` guarda num slice; o SQL nunca rodou.
2. Que `ListExpiringSessions` de fato **pula** a conta enquanto `next_retry_at` está no futuro. Isso é um `WHERE` real (`auth_session_repo.go:87`) que nenhum fake exercita — e sem ele todo o backoff da Task 2 é decorativo: o ticker retentaria a conta morta a cada 5 minutos.
3. Que `ApplyConnectionSnapshot` com `needs_reauth` grava `health_status = 'critical'` e devolve o `reauth_reason` na releitura. `HealthStatusCritical` nunca foi escrito neste banco; `0016_integrations_foundation.sql` tem CHECKs, e um valor recusado pelo CHECK só aparece aqui.

**Files:**
- Create: `apps/server_core/tests/integration/integrations_refresh_failure_test.go`

- [ ] **Step 1: Escreva o teste**

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

// Prova, contra Postgres real, as duas coisas que os fakes das Tasks 2-4 não
// conseguem provar: que a falha persiste com a forma certa, e que a varredura
// respeita next_retry_at. O segundo é o que faz o backoff existir de fato — se
// o WHERE não filtrasse, o ticker retentaria a conta morta a cada tick e o
// consecutive_failures subiria para sempre.
func TestRefreshFailurePersistsAndSuppressesTheSweep(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_refresh_failure")
	testpostgres.SeedProvider(t, pool, testpostgres.ProviderFixture{Code: "mercado_livre", DisplayName: "Mercado Livre"})

	installationRepo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	sessionRepo := integrationspostgres.NewAuthSessionRepository(pool, cfg.DefaultTenantID)

	ctx := context.Background()
	installationID := fmt.Sprintf("inst-refresh-fail-%d", time.Now().UTC().UnixNano())
	now := time.Now().UTC()

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM auth_sessions WHERE installation_id = $1`, installationID)
		_, _ = pool.Exec(ctx, `DELETE FROM integration_installations WHERE installation_id = $1`, installationID)
	})

	if err := installationRepo.CreateInstallation(ctx, integrationsdomain.Installation{
		InstallationID: installationID,
		TenantID:       cfg.DefaultTenantID,
		ProviderCode:   "mercado_livre",
		Family:         integrationsdomain.IntegrationFamilyMarketplace,
		DisplayName:    "Mercado Livre (teste)",
		Status:         integrationsdomain.InstallationStatusConnected,
		HealthStatus:   integrationsdomain.HealthStatusHealthy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}); err != nil {
		t.Fatalf("CreateInstallation: %v", err)
	}

	// --- 1. A falha persiste com a forma que o domínio escreveu ---------------
	// access_token_expires_at no passado: sem next_retry_at, esta sessão CAI na
	// varredura. É essa qualificação que torna o passo 3 uma prova.
	expiredAt := now.Add(-time.Minute)
	nextRetryAt := now.Add(time.Hour)
	failureCode := integrationsdomain.ErrRefreshTokenInvalid.Error()

	if err := sessionRepo.UpsertAuthSession(ctx, integrationsdomain.AuthSession{
		AuthSessionID:        "auth_" + installationID,
		TenantID:             cfg.DefaultTenantID,
		InstallationID:       installationID,
		State:                integrationsdomain.AuthStateRefreshFailed,
		ProviderAccountID:    "seller-1",
		AccessTokenExpiresAt: &expiredAt,
		NextRetryAt:          &nextRetryAt,
		RefreshFailureCode:   failureCode,
		ConsecutiveFailures:  1,
		CreatedAt:            now,
		UpdatedAt:            now,
	}); err != nil {
		t.Fatalf("UpsertAuthSession: %v", err)
	}

	stored, found, err := sessionRepo.GetAuthSession(ctx, installationID)
	if err != nil || !found {
		t.Fatalf("GetAuthSession found=%v err=%v", found, err)
	}
	if stored.State != integrationsdomain.AuthStateRefreshFailed {
		t.Fatalf("State = %q, want refresh_failed (o CHECK de 0016 aceita esse valor?)", stored.State)
	}
	if stored.RefreshFailureCode != failureCode {
		t.Fatalf("RefreshFailureCode = %q, want %q", stored.RefreshFailureCode, failureCode)
	}
	if stored.ConsecutiveFailures != 1 {
		t.Fatalf("ConsecutiveFailures = %d, want 1", stored.ConsecutiveFailures)
	}
	if stored.NextRetryAt == nil {
		t.Fatal("NextRetryAt = nil depois do round-trip")
	}

	// --- 2. Controle positivo: sem backoff, a varredura PEGA a sessão --------
	// Sem este passo, o passo 3 seria vacuoso: uma lista vazia por qualquer
	// outro motivo (tenant errado, expiry errado) pareceria backoff funcionando.
	clearRetry := integrationsdomain.AuthSession{
		AuthSessionID:        "auth_" + installationID,
		TenantID:             cfg.DefaultTenantID,
		InstallationID:       installationID,
		State:                integrationsdomain.AuthStateExpiring,
		ProviderAccountID:    "seller-1",
		AccessTokenExpiresAt: &expiredAt,
		NextRetryAt:          nil,
		RefreshFailureCode:   failureCode,
		ConsecutiveFailures:  1,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := sessionRepo.UpsertAuthSession(ctx, clearRetry); err != nil {
		t.Fatalf("UpsertAuthSession (sem backoff): %v", err)
	}
	if !sweepContains(t, sessionRepo, installationID) {
		t.Fatal("a varredura NÃO pegou a sessão vencida sem next_retry_at: o teste do backoff seria vacuoso")
	}

	// --- 3. Com next_retry_at no futuro, a varredura pula ---------------------
	withBackoff := clearRetry
	withBackoff.NextRetryAt = &nextRetryAt
	if err := sessionRepo.UpsertAuthSession(ctx, withBackoff); err != nil {
		t.Fatalf("UpsertAuthSession (com backoff): %v", err)
	}
	if sweepContains(t, sessionRepo, installationID) {
		t.Fatal("a varredura pegou a sessão com next_retry_at no futuro: o backoff não existe")
	}

	// --- 4. O estado crítico chega ao banco e volta -------------------------
	snapshot := integrationsdomain.ProjectConnectionSnapshot(
		integrationsdomain.Installation{
			InstallationID: installationID,
			TenantID:       cfg.DefaultTenantID,
			ProviderCode:   "mercado_livre",
			Status:         integrationsdomain.InstallationStatusRequiresReauth,
			HealthStatus:   integrationsdomain.HealthStatusCritical,
		},
		integrationsdomain.AuthStrategyOAuth2,
		&expiredAt,
		failureCode,
	)
	if err := installationRepo.ApplyConnectionSnapshot(ctx, installationID, snapshot, ""); err != nil {
		t.Fatalf("ApplyConnectionSnapshot: %v", err)
	}

	reread, found, err := installationRepo.GetInstallation(ctx, installationID)
	if err != nil || !found {
		t.Fatalf("GetInstallation found=%v err=%v", found, err)
	}
	if reread.HealthStatus != integrationsdomain.HealthStatusCritical {
		t.Fatalf("HealthStatus = %q, want critical (primeira escrita de 'critical' neste banco)", reread.HealthStatus)
	}
	if reread.ConnectionSnapshot.NextAction != integrationsdomain.ConnectionNextActionReauth {
		t.Fatalf("NextAction = %q, want reauth", reread.ConnectionSnapshot.NextAction)
	}
	if reread.ConnectionSnapshot.ReauthReason != failureCode {
		t.Fatalf("ReauthReason = %q, want %q", reread.ConnectionSnapshot.ReauthReason, failureCode)
	}
}

func sweepContains(t *testing.T, repo *integrationspostgres.AuthSessionRepository, installationID string) bool {
	t.Helper()
	sessions, err := repo.ListExpiringSessions(context.Background(), 10*time.Minute)
	if err != nil {
		t.Fatalf("ListExpiringSessions: %v", err)
	}
	for _, s := range sessions {
		if s.InstallationID == installationID {
			return true
		}
	}
	return false
}
```

> `AuthStrategyOAuth2` é o nome esperado da constante em `domain/lifecycle.go`. Confirme com `grep -n "AuthStrategy[A-Z]" internal/modules/integrations/domain/lifecycle.go` e ajuste se o nome for outro. Se `Installation` não tiver campo `TenantID` ou `IntegrationFamilyMarketplace` tiver outro nome, copie do teste vizinho `integrations_credential_repo_test.go:17-30`, que já compila.

- [ ] **Step 2: Rode contra o Postgres real**

A lane se auto-provisiona: `Invoke-Integration` (`scripts/harness.ps1:68-88`) cria ou reusa o container Postgres, aplica as migrations embutidas e injeta `MPC_TEST_DATABASE_URL` no processo filho. **Não** exporte DSN à mão — passar `-DatabaseUrl` lança `HPG_EXTERNAL_TARGET_FORBIDDEN` de propósito. Nenhum `REQUEST` ao hub é necessário.

```bash
pwsh scripts/harness.ps1 -Command integration
```

Saída esperada no cabeçalho: `target=ephemeral-postgres`, `key=MPC_TEST_DATABASE_URL`, `migrations=embedded`, depois `migrations_second=0` (idempotência) e `status=passed`.

Esperado, **antes** das Tasks 2–4 estarem aplicadas: este teste não depende delas — ele exercita repositório, não serviço. Ele tem que passar assim que for escrito. Se falhar no passo 1 com erro de CHECK constraint, o valor `refresh_failed` ou `critical` não é aceito pelo schema e a Task 3 estava errada; pare e reporte.

- [ ] **Step 3: Injete o defeito e veja o vermelho (must-fail)**

Verde sozinho não prova que o teste mede alguma coisa. Comente temporariamente a cláusula de backoff em `auth_session_repo.go:87`:

```sql
    AND (next_retry_at IS NULL OR next_retry_at <= now())
```

Rode a lane de novo. Esperado: FAIL com `a varredura pegou a sessão com next_retry_at no futuro: o backoff não existe`. Restaure a linha e confirme o verde. **Sem este passo o teste não conta.**

- [ ] **Step 4: Commit**

```bash
git add apps/server_core/tests/integration/integrations_refresh_failure_test.go
git commit -m "test(integration): prove refresh backoff against real Postgres

The unit tests use fakes, so neither the persisted failure shape nor the
next_retry_at suppression in ListExpiringSessions was ever exercised against
real SQL — and the whole backoff rests on that WHERE clause. Also the first
write of health_status='critical' in this schema. Includes a positive control
(sweep DOES pick the session with no backoff) so an empty list can't pass for
a working backoff."
```

---

### Task 6: `/integracoes` mostra o estado da conexão

Backend agora produz o sinal. A tela ainda não lê nada: `IntegracoesPage.tsx` só referencia `status` em `:516`, para decidir se reusa uma installation pendente. `health_status` e `connection` nunca são renderizados.

**Files:**
- Create: `apps/web/src/pages/integracoes/ConnectionHealthCard.tsx`
- Modify: `apps/web/src/pages/integracoes/IntegracoesPage.tsx:559-576`
- Test: `apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx`

Sem hook novo: `useInstallation()` (`apps/web/src/app/InstallationContext.tsx:117`) já expõe `installations` e `status`, buscados por `installationsQueryKeys.list()`. Um segundo `useQuery` para a mesma lista seria fetch duplicado.

O teste abaixo faz `vi.mock` do `InstallationContext`. Isso prova **uma coisa só**: que dado um payload, o componente escolhe o rótulo e o tom certos. Não prova que o backend produz esse payload, nem que a página monta o card, nem que o operador enxerga. Isso é a Task 7 e é o gate.

- [ ] **Step 1: Escreva o teste que falha**

Crie `apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx`:

```tsx
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { ConnectionHealthCard } from "./ConnectionHealthCard";

const mockUseInstallation = vi.fn();
vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => mockUseInstallation(),
}));

function installation(overrides: Partial<IntegrationInstallation>): IntegrationInstallation {
  return {
    installation_id: "inst-1",
    tenant_id: "tenant-1",
    provider_code: "mercado_livre",
    family: "marketplace",
    display_name: "Mercado Livre (cliente)",
    status: "connected",
    health_status: "healthy",
    external_account_id: "123",
    external_account_name: "loja",
    connection: {
      state: "connected",
      health: "healthy",
      provider_code: "mercado_livre",
      external_account_id: "123",
      external_account_name: "loja",
      auth_strategy: "oauth2",
      next_action: "none",
    },
    runtime_capabilities: [],
    created_at: "2026-08-01T00:00:00Z",
    updated_at: "2026-08-01T00:00:00Z",
    ...overrides,
  };
}

describe("ConnectionHealthCard", () => {
  it("mostra a razão e a ação quando o token exige reautorização", () => {
    mockUseInstallation.mockReturnValue({
      status: "ready",
      installations: [
        installation({
          status: "requires_reauth",
          health_status: "critical",
          connection: {
            ...installation({}).connection,
            state: "needs_reauth",
            health: "critical",
            next_action: "reauth",
            reauth_reason: "INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400",
          },
        }),
      ],
    });

    render(<ConnectionHealthCard />);

    expect(screen.getByTestId("connection-health-inst-1")).toHaveTextContent("Reautorizar");
    expect(screen.getByTestId("connection-health-reason-inst-1")).toHaveTextContent(
      "INTEGRATIONS_REFRESH_TOKEN_INVALID",
    );
  });

  it("não inventa alarme numa conta saudável", () => {
    mockUseInstallation.mockReturnValue({
      status: "ready",
      installations: [installation({})],
    });

    render(<ConnectionHealthCard />);

    expect(screen.getByTestId("connection-health-inst-1")).toHaveTextContent("Conectado");
    expect(screen.queryByTestId("connection-health-reason-inst-1")).toBeNull();
  });

  it("diz que o estado é desconhecido quando a leitura falha", () => {
    // ADR-17: leitura falhada nunca vira "tudo ok" nem card em branco.
    mockUseInstallation.mockReturnValue({ status: "error", installations: [] });

    render(<ConnectionHealthCard />);

    expect(screen.getByTestId("connection-health-unknown")).toBeInTheDocument();
  });
});
```

- [ ] **Step 2: Rode e veja falhar**

```bash
npm run test --workspace apps/web -- ConnectionHealthCard
```

Esperado: FAIL com `Failed to resolve import "./ConnectionHealthCard"`.

- [ ] **Step 3: Implemente**

Crie `apps/web/src/pages/integracoes/ConnectionHealthCard.tsx`:

```tsx
import type { IntegrationConnectionSnapshot, IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { useInstallation } from "../../app/InstallationContext";

// Tom lido do ESTADO do payload, nunca de um corte de tempo — mesmo critério do
// SyncHealthCard (entityTone, SyncHealthCard.tsx:13). Uma conta que precisa de
// reautorização precisa disso agora, por mais recente que seja o last_verified_at.
type ConnectionTone = "green" | "amber" | "red" | "gray";

const stateTone: Record<IntegrationConnectionSnapshot["state"], ConnectionTone> = {
  connected: "green",
  degraded: "amber",
  needs_reauth: "red",
  disconnected: "red",
  pending_connection: "gray",
  draft: "gray",
};

const stateLabel: Record<IntegrationConnectionSnapshot["state"], string> = {
  connected: "Conectado",
  degraded: "Instável",
  needs_reauth: "Precisa reautorizar",
  disconnected: "Desconectado",
  pending_connection: "Aguardando autorização",
  draft: "Rascunho",
};

const nextActionLabel: Record<IntegrationConnectionSnapshot["next_action"], string> = {
  none: "",
  authorize: "Autorizar",
  reauth: "Reautorizar",
  configure: "Configurar",
  retry: "Repetindo automaticamente",
};

const toneBadgeClassName: Record<ConnectionTone, string> = {
  green: "bg-accent-soft text-accent-ink",
  amber: "bg-warn-soft text-warn",
  red: "bg-danger-soft text-danger",
  gray: "bg-surface-2 text-faint",
};

function InstallationRow({ installation }: { installation: IntegrationInstallation }) {
  const connection = installation.connection;
  const tone = stateTone[connection.state];
  const action = nextActionLabel[connection.next_action];
  const reason = connection.reauth_reason?.trim() ?? "";

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
      {action ? <span className="text-faint">Ação: {action}</span> : null}
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

// Consome useInstallation() em vez de abrir um useQuery próprio: a lista já é
// buscada uma vez por installationsQueryKeys.list() em InstallationContext.tsx:31.
export function ConnectionHealthCard() {
  const { installations, status } = useInstallation();

  return (
    <section aria-labelledby="connection-health-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="connection-health-title" className="text-sm font-semibold text-ink">
        Contas conectadas
      </h2>
      {status === "loading" ? <LoadingState /> : null}
      {/* ADR-17: leitura que falhou não pode renderizar verde nem branco. */}
      {status === "error" ? (
        <div className="mt-2" data-testid="connection-health-unknown">
          <ErrorState detail="Não foi possível carregar o estado das contas. Estado desconhecido." />
        </div>
      ) : null}
      {status !== "loading" && status !== "error" ? (
        <div className="mt-3 flex flex-col gap-2">
          {installations.map((installation) => (
            <InstallationRow key={installation.installation_id} installation={installation} />
          ))}
        </div>
      ) : null}
    </section>
  );
}
```

> Se `ErrorState` exigir `onRetry`, passe `onRetry={undefined}` só se a prop for opcional; caso contrário use o mesmo `queryClient.refetchQueries({ queryKey: installationsQueryKeys.list() })` de `InstallationContext.tsx:107`. Não invente um retry novo.
>
> Se as classes `bg-danger-soft`/`text-danger` não existirem no tema, use `bg-warn-soft text-warn` (que `SyncHealthCard.tsx:21` já usa) e mantenha `red` distinto de `amber` só pelo rótulo. Confirme com `grep -rn "danger-soft" apps/web/src packages/ui/src` antes de escolher.

- [ ] **Step 4: Monte na página**

Em `IntegracoesPage.tsx`, adicione o import no topo e o card entre `ProviderConnectCard` e `SyncHealthCard`:

```tsx
import { ConnectionHealthCard } from "./ConnectionHealthCard";
```

```tsx
      <ProviderConnectCard />
      <ConnectionHealthCard />
      <SyncHealthCard />
```

- [ ] **Step 5: Rode e veja passar**

```bash
npm run test --workspace apps/web -- ConnectionHealthCard
```

Esperado: PASS nos três casos.

```bash
npm run typecheck --workspace apps/web
```

Esperado: sem erro novo. O repo tem `tsc` vermelho herdado em outros módulos — compare com `git stash`/`git stash pop` se aparecer ruído, e só conserte o que estes arquivos introduziram.

- [ ] **Step 6: Commit**

```bash
git add apps/web/src/pages/integracoes/ConnectionHealthCard.tsx apps/web/src/pages/integracoes/ConnectionHealthCard.test.tsx apps/web/src/pages/integracoes/IntegracoesPage.tsx
git commit -m "feat(web): surface connection health and reauth reason on /integracoes

The page read installation.status only to decide whether to reuse a pending row;
health_status and connection were never rendered, so a dead ML token looked
identical to a healthy account. The new card reads the snapshot the backend now
degrades, reuses useInstallation() instead of refetching the list, and names an
unreadable state instead of rendering blank (ADR-17)."
```

---

### Task 7: Verificação ao vivo no browser — o gate da fatia

Nada acima prova que o operador vê a falha. Unidade prova ramo, integração prova SQL; só esta task prova o caminho inteiro contra a conta ML real. **A fatia não fecha sem ela.**

Stack já de pé (verificado 2026-08-03): `marketplace-central-backend-1` healthy em `:8080`, `marketplace-central-frontend-1` em `:5174`, `marketplace-central-postgres-1` healthy em `:5435`. Chip não sobe nem derruba serviço — se o stack cair, é `REQUEST` ao hub, nunca `docker compose up` daqui.

- [ ] **Step 1: Controle positivo — corrompa o refresh token da conta ML**

Numa sessão psql do dev stack, invalide o refresh token guardado da installation ML. O payload é criptografado, então o caminho barato é forçar a expiração do access token para que o ticker tente refrescar já:

```sql
UPDATE auth_sessions
   SET access_token_expires_at = now() - interval '1 minute',
       state = 'expiring',
       next_retry_at = NULL
 WHERE installation_id = '<inst-id-ml>';
```

E aponte o adapter para um token endpoint que devolve `invalid_grant` — ou, se for mais simples, revogue a autorização do app na conta ML e deixe o ML devolver `invalid_grant` de verdade. **A segunda opção exige o operador**: revogação é ação na conta real dele. Peça antes.

- [ ] **Step 2: Espere um tick e leia o log**

O intervalo é 5 min (`root.go:683`). Nos logs do servidor, espere:

```
level=ERROR msg="integrations refresh ticker item failed" installation_id=<inst-id-ml> err="INTEGRATIONS_REFRESH_TOKEN_INVALID: status=400 body=..."
```

Antes desta Onda esse log **não existia** — a ausência dele é o before, a presença é o after.

- [ ] **Step 3: Confirme a persistência**

```sql
SELECT state, refresh_failure_code, consecutive_failures, next_retry_at
  FROM auth_sessions WHERE installation_id = '<inst-id-ml>';
```

Esperado: `state = 'refresh_failed'`, código preenchido, `consecutive_failures = 1`, `next_retry_at` uma hora à frente (cooldown terminal).

```sql
SELECT status, health_status, connection_snapshot->>'next_action', connection_snapshot->>'reauth_reason'
  FROM integration_installations WHERE installation_id = '<inst-id-ml>';
```

Esperado: `requires_reauth`, `critical`, `reauth`, motivo preenchido.

- [ ] **Step 4: Abra `/integracoes` no browser**

Esperado: o card "Contas conectadas" mostra a conta ML com badge vermelho "Precisa reautorizar", linha "Ação: Reautorizar" e o motivo cru. Tire screenshot.

- [ ] **Step 5: Prove que o lote não aborta (F-A2)**

Com a conta ML em falha, o log do tick seguinte tem que mostrar o erro da conta quebrada **e** nenhum `integrations refresh ticker pass failed`. Com uma conta só, isso é o máximo que o live drive prova; a cobertura real do F-A2 é `TestRunOnceContinuesAfterItemFailure`.

- [ ] **Step 6: Restaure**

Reautorize a conta pelo botão "Conectar" da própria página e confirme que o card volta a "Conectado", `consecutive_failures` volta a 0 e `next_retry_at` a NULL — o caminho de sucesso em `auth_flow_service.go:456-468` já faz isso e `TestRefreshCredentialRotatesAndResetsFailures` já o cobre; aqui só confirmamos ao vivo.

- [ ] **Step 7: Grave a evidência**

Screenshots, trechos de log e saídas de SQL vão para o pack em `.mnfs/` da missão. Não escrito = não aconteceu.

---

## Fora de escopo desta fatia

- `F-A3` (coleta de mercado periódica + idade visível) — vive em `market/`, plano próprio.
- `F-00` (scheduler de pedidos) — bloqueado por `D-16`, plano próprio.
- Unificar `mapIntegrationError` (`transport/http_handler.go:47-58`) com o resto do repo: ele casa por `strings.HasPrefix(err.Error(), "INTEGRATIONS_")`, o que **não** casa erro embrulhado com `%w` e devolve 400 para tudo. Real, mas não está no caminho de nenhuma task aqui — nenhum erro de refresh sai por HTTP. Registre como dívida em `.mnfs/HARNESS-DEBTS.md` ou no ledger da missão, não conserte de carona.
- `fee_sync_scheduler.go:74` descarta erro com `_, _ =`. Mesma classe do F-A1, escopo diferente. Registre.
