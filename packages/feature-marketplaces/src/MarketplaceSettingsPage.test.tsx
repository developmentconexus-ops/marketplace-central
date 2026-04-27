import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import type { IntegrationInstallation, IntegrationProviderDefinition } from "@marketplace-central/sdk-runtime";
import { MarketplaceSettingsPage } from "./MarketplaceSettingsPage";

const mockListProviders = vi.fn();
const mockListInstallations = vi.fn();
const mockCreateInstallation = vi.fn();
const mockStartAuthorization = vi.fn();
const mockStartReauthorization = vi.fn();
const mockSubmitCredentials = vi.fn();
const mockDisconnectInstallation = vi.fn();
const mockStartFeeSync = vi.fn();
const mockListAccounts = vi.fn();
const mockListPolicies = vi.fn();
const mockNavigateToAuthUrl = vi.fn();

const mockClient = {
  listIntegrationProviders: mockListProviders,
  listIntegrationInstallations: mockListInstallations,
  createIntegrationInstallation: mockCreateInstallation,
  startIntegrationAuthorization: mockStartAuthorization,
  startIntegrationReauthorization: mockStartReauthorization,
  submitIntegrationCredentials: mockSubmitCredentials,
  disconnectIntegrationInstallation: mockDisconnectInstallation,
  startIntegrationFeeSync: mockStartFeeSync,
  listMarketplaceAccounts: mockListAccounts,
  listMarketplacePolicies: mockListPolicies,
} as const;

const baseProviders: IntegrationProviderDefinition[] = [
  {
    provider_code: "mercado_livre",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Mercado Livre",
    auth_strategy: "oauth2",
    install_mode: "interactive",
    metadata: { execution_mode: "available", rollout_stage: "v1", baseline_commission_percent: 0.16 },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
  {
    provider_code: "magalu",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Magalu",
    auth_strategy: "oauth2",
    install_mode: "interactive",
    metadata: { execution_mode: "available", rollout_stage: "v1", baseline_commission_percent: 0.16 },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
  {
    provider_code: "shopee",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Shopee",
    auth_strategy: "shopee_partner",
    install_mode: "interactive",
    metadata: { execution_mode: "available", rollout_stage: "v1", fee_source: "seed" },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
  {
    provider_code: "amazon",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Amazon Brasil",
    auth_strategy: "lwa",
    install_mode: "interactive",
    metadata: { execution_mode: "available", rollout_stage: "v1", baseline_commission_percent: 0.12 },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
  {
    provider_code: "leroy_merlin",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Leroy Merlin",
    auth_strategy: "api_key",
    install_mode: "manual",
    metadata: { execution_mode: "available", rollout_stage: "wave_2", baseline_commission_percent: 0.18 },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
  {
    provider_code: "madeira_madeira",
    tenant_id: "system",
    family: "marketplace",
    display_name: "Madeira Madeira",
    auth_strategy: "token",
    install_mode: "manual",
    metadata: { execution_mode: "available", rollout_stage: "wave_2", baseline_commission_percent: 0.15 },
    declared_capabilities: ["pricing_fee_sync"],
    is_active: true,
    created_at: "2026-04-25T00:00:00Z",
    updated_at: "2026-04-25T00:00:00Z",
  },
];

const connectedInstallation: IntegrationInstallation = {
  installation_id: "inst-amz",
  tenant_id: "t1",
  provider_code: "amazon",
  family: "marketplace",
  display_name: "Amazon Brasil",
  status: "connected",
  health_status: "healthy",
  external_account_id: "seller-1",
  external_account_name: "Seller 1",
  created_at: "2026-04-25T00:00:00Z",
  updated_at: "2026-04-25T00:00:00Z",
};

beforeEach(() => {
  mockListProviders.mockReset();
  mockListInstallations.mockReset();
  mockCreateInstallation.mockReset();
  mockStartAuthorization.mockReset();
  mockStartReauthorization.mockReset();
  mockSubmitCredentials.mockReset();
  mockDisconnectInstallation.mockReset();
  mockStartFeeSync.mockReset();
  mockListAccounts.mockReset();
  mockListPolicies.mockReset();
  mockNavigateToAuthUrl.mockReset();

  mockListProviders.mockResolvedValue({ items: baseProviders });
  mockListInstallations.mockResolvedValue({ items: [] });
  mockListAccounts.mockResolvedValue({ items: [] });
  mockListPolicies.mockResolvedValue({ items: [] });
});

describe("MarketplaceSettingsPage", () => {
  it("renders all six providers from integration catalog", async () => {
    render(<MarketplaceSettingsPage client={mockClient as never} navigateToAuthUrl={mockNavigateToAuthUrl} />);
    await waitFor(() => expect(screen.getByText("Mercado Livre")).toBeInTheDocument());
    expect(screen.getByText("Magalu")).toBeInTheDocument();
    expect(screen.getByText("Shopee")).toBeInTheDocument();
    expect(screen.getByText("Amazon Brasil")).toBeInTheDocument();
    expect(screen.getByText("Leroy Merlin")).toBeInTheDocument();
    expect(screen.getByText("Madeira Madeira")).toBeInTheDocument();
  });

  it("shows loading state first", () => {
    const { container } = render(<MarketplaceSettingsPage client={mockClient as never} />);
    expect(container.querySelector("[aria-busy='true']")).not.toBeNull();
  });

  it("shows error banner when loading fails", async () => {
    mockListProviders.mockRejectedValue({ error: { message: "catalog unavailable" } });
    render(<MarketplaceSettingsPage client={mockClient as never} navigateToAuthUrl={mockNavigateToAuthUrl} />);
    await waitFor(() => expect(screen.getByText(/catalog unavailable/i)).toBeInTheDocument());
  });

  it("maps connected installation status onto provider card", async () => {
    mockListInstallations.mockResolvedValue({ items: [connectedInstallation] });
    render(<MarketplaceSettingsPage client={mockClient as never} />);
    await waitFor(() => expect(screen.getByText("Amazon Brasil")).toBeInTheDocument());
    expect(screen.getByLabelText("Status: connected")).toBeInTheDocument();
  });

  it("starts interactive authorize flow for Shopee", async () => {
    const shopeeInstallation: IntegrationInstallation = {
      installation_id: "inst-shopee",
      tenant_id: "t1",
      provider_code: "shopee",
      family: "marketplace",
      display_name: "Shopee",
      status: "disconnected",
      health_status: "warning",
      external_account_id: "",
      external_account_name: "",
      created_at: "2026-04-25T00:00:00Z",
      updated_at: "2026-04-25T00:00:00Z",
    };
    mockListInstallations.mockResolvedValue({ items: [shopeeInstallation] });
    mockStartAuthorization.mockResolvedValue({ auth_url: "https://auth.test/shopee" });

    render(<MarketplaceSettingsPage client={mockClient as never} navigateToAuthUrl={mockNavigateToAuthUrl} />);
    await waitFor(() => expect(screen.getByText("Shopee")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Shopee" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "Authorize" })).toBeInTheDocument());
    expect(screen.queryByRole("button", { name: "Submit credentials" })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Authorize" }));

    await waitFor(() => expect(mockStartAuthorization).toHaveBeenCalledWith("inst-shopee"));
    expect(mockNavigateToAuthUrl).toHaveBeenCalledWith("https://auth.test/shopee");
  });

  it("creates installation and starts auth for interactive provider", async () => {
    const created: IntegrationInstallation = {
      installation_id: "inst-ml",
      tenant_id: "t1",
      provider_code: "mercado_livre",
      family: "marketplace",
      display_name: "Mercado Livre",
      status: "pending_connection",
      health_status: "warning",
      external_account_id: "",
      external_account_name: "",
      created_at: "2026-04-25T00:00:00Z",
      updated_at: "2026-04-25T00:00:00Z",
    };
    mockCreateInstallation.mockResolvedValue(created);
    mockStartAuthorization.mockResolvedValue({ auth_url: "https://auth.test" });

    render(<MarketplaceSettingsPage client={mockClient as never} navigateToAuthUrl={mockNavigateToAuthUrl} />);
    await waitFor(() => expect(screen.getByText("Mercado Livre")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Mercado Livre" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "Create installation" })).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Create installation" }));

    await waitFor(() => expect(mockCreateInstallation).toHaveBeenCalled());
    await waitFor(() => expect(mockStartAuthorization).toHaveBeenCalledWith("inst-ml"));
    expect(mockNavigateToAuthUrl).toHaveBeenCalledWith("https://auth.test");
  });

  it("starts auth again for disconnected interactive provider", async () => {
    const disconnectedInstallation: IntegrationInstallation = {
      ...connectedInstallation,
      status: "disconnected",
      health_status: "warning",
      external_account_id: "",
      external_account_name: "",
    };
    mockListInstallations.mockResolvedValue({ items: [disconnectedInstallation] });
    mockStartAuthorization.mockResolvedValue({ auth_url: "https://auth.test/again" });

    render(<MarketplaceSettingsPage client={mockClient as never} navigateToAuthUrl={mockNavigateToAuthUrl} />);
    await waitFor(() => expect(screen.getByText("Amazon Brasil")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Amazon Brasil" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "Authorize" })).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Authorize" }));

    await waitFor(() => expect(mockStartAuthorization).toHaveBeenCalledWith("inst-amz"));
    expect(mockNavigateToAuthUrl).toHaveBeenCalledWith("https://auth.test/again");
  });

  it("disables fee sync when installation is pending connection", async () => {
    const pendingInstallation: IntegrationInstallation = {
      ...connectedInstallation,
      status: "pending_connection",
      health_status: "warning",
      external_account_id: "",
      external_account_name: "",
    };
    mockListInstallations.mockResolvedValue({ items: [pendingInstallation] });

    render(<MarketplaceSettingsPage client={mockClient as never} />);
    await waitFor(() => expect(screen.getByText("Amazon Brasil")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Amazon Brasil" }));

    await waitFor(() => expect(screen.getByRole("button", { name: "Run fee sync" })).toBeDisabled());
  });

  it("enables fee sync when installation is connected", async () => {
    mockListInstallations.mockResolvedValue({ items: [connectedInstallation] });

    render(<MarketplaceSettingsPage client={mockClient as never} />);
    await waitFor(() => expect(screen.getByText("Amazon Brasil")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Amazon Brasil" }));

    await waitFor(() => expect(screen.getByRole("button", { name: "Run fee sync" })).toBeEnabled());
  });

  it("shows action error banner when authorize fails", async () => {
    const pendingInstallation: IntegrationInstallation = {
      ...connectedInstallation,
      status: "pending_connection",
      health_status: "warning",
      external_account_id: "",
      external_account_name: "",
    };
    mockListInstallations.mockResolvedValue({ items: [pendingInstallation] });
    mockStartAuthorization.mockRejectedValue({
      status: 400,
      error: {
        code: "INTEGRATIONS_AUTH_CONFIGURATION_INVALID",
        message: "INTEGRATIONS_AUTH_CONFIGURATION_INVALID",
      },
    });

    render(<MarketplaceSettingsPage client={mockClient as never} />);
    await waitFor(() => expect(screen.getByText("Amazon Brasil")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Open Amazon Brasil" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "Authorize" })).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Authorize" }));

    await waitFor(() =>
      expect(
        screen.getByText(/Amazon authorization config is invalid/i),
      ).toBeInTheDocument(),
    );
  });
});
