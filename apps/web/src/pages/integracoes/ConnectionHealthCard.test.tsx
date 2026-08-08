import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach } from "vitest";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { ConnectionHealthCard } from "./ConnectionHealthCard";

const mockUseInstallation = vi.fn();
const startIntegrationReauthorization = vi.fn();

vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => mockUseInstallation(),
}));

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    startIntegrationReauthorization: (...args: unknown[]) =>
      startIntegrationReauthorization(...args),
  }),
}));

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

// ConnectionHealthCard calls useQueryClient() (needed for ErrorState's
// required onRetry, same idiom as InstallationContext.tsx:107), so it must
// render inside a QueryClientProvider — same pattern as
// SyncHealthCard.test.tsx:89.
function renderCard(overrides?: Record<string, unknown>) {
  if (overrides) {
    const inst = {
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
        reauth_reason: "",
      },
      runtime_capabilities: [],
      created_at: "2026-08-01T00:00:00Z",
      updated_at: "2026-08-01T00:00:00Z",
      ...overrides,
    };
    mockUseInstallation.mockReturnValue({ installations: [inst], status: "ready" });
  }
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ConnectionHealthCard />
    </QueryClientProvider>,
  );
}

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

    renderCard();

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

    renderCard();

    expect(screen.getByTestId("connection-health-inst-1")).toHaveTextContent("Conectado");
    expect(screen.queryByTestId("connection-health-reason-inst-1")).toBeNull();
  });

  it("diz que o estado é desconhecido quando a leitura falha", () => {
    // ADR-17: leitura falhada nunca vira "tudo ok" nem card em branco.
    mockUseInstallation.mockReturnValue({ status: "error", installations: [] });

    renderCard();

    expect(screen.getByTestId("connection-health-unknown")).toBeInTheDocument();
  });

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
        health: "critical",
        provider_code: "mercado_livre",
        external_account_id: "123",
        external_account_name: "loja",
        auth_strategy: "oauth2",
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
      connection: {
        state: "connected",
        health: "healthy",
        provider_code: "mercado_livre",
        external_account_id: "123",
        external_account_name: "loja",
        auth_strategy: "oauth2",
        next_action: "none",
        reauth_reason: "",
      },
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
      connection: {
        state: "needs_reauth",
        health: "critical",
        provider_code: "mercado_livre",
        external_account_id: "123",
        external_account_name: "loja",
        auth_strategy: "oauth2",
        next_action: "reauth",
        reauth_reason: "",
      },
    });

    fireEvent.click(screen.getByTestId("connection-reauth-inst-1"));

    expect(await screen.findByTestId("connection-reauth-error-inst-1")).toHaveTextContent(
      "Não foi possível iniciar a reautorização.",
    );
    expect(assignedUrl).toBe("");
  });
});
