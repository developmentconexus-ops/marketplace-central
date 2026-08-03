import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { ConnectionHealthCard } from "./ConnectionHealthCard";

const mockUseInstallation = vi.fn();
vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => mockUseInstallation(),
}));

// ConnectionHealthCard calls useQueryClient() (needed for ErrorState's
// required onRetry, same idiom as InstallationContext.tsx:107), so it must
// render inside a QueryClientProvider — same pattern as
// SyncHealthCard.test.tsx:89.
function renderCard() {
  const queryClient = new QueryClient();
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
});
