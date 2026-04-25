import { useCallback, useEffect, useMemo, useState } from "react";
import type {
  CreateIntegrationInstallationRequest,
  IntegrationInstallation,
  IntegrationProviderDefinition,
  MarketplaceAccount,
  MarketplacePolicy,
  SubmitIntegrationCredentialsRequest,
} from "@marketplace-central/sdk-runtime";
import { ProviderCatalogCard } from "./ProviderCatalogCard";
import { ProviderCatalogPanel } from "./ProviderCatalogPanel";

interface MarketplaceClient {
  listIntegrationProviders: () => Promise<{ items: IntegrationProviderDefinition[] }>;
  listIntegrationInstallations: () => Promise<{ items: IntegrationInstallation[] }>;
  createIntegrationInstallation: (req: CreateIntegrationInstallationRequest) => Promise<IntegrationInstallation>;
  startIntegrationAuthorization: (installationId: string) => Promise<{ auth_url: string }>;
  startIntegrationReauthorization: (installationId: string) => Promise<{ auth_url: string }>;
  submitIntegrationCredentials: (
    installationId: string,
    req: SubmitIntegrationCredentialsRequest,
  ) => Promise<unknown>;
  disconnectIntegrationInstallation: (installationId: string) => Promise<unknown>;
  startIntegrationFeeSync: (installationId: string) => Promise<unknown>;
  listMarketplaceAccounts: () => Promise<{ items: MarketplaceAccount[] }>;
  listMarketplacePolicies: () => Promise<{ items: MarketplacePolicy[] }>;
}

interface MarketplaceSettingsPageProps {
  client: MarketplaceClient;
  navigateToAuthUrl?: (url: string) => void;
}

function createInstallationID(providerCode: string): string {
  return `inst-${providerCode}-${crypto.randomUUID()}`;
}

function defaultNavigateToAuthUrl(url: string) {
  window.location.assign(url);
}

export function MarketplaceSettingsPage({
  client,
  navigateToAuthUrl = defaultNavigateToAuthUrl,
}: MarketplaceSettingsPageProps) {
  const [providers, setProviders] = useState<IntegrationProviderDefinition[]>([]);
  const [installations, setInstallations] = useState<IntegrationInstallation[]>([]);
  const [accounts, setAccounts] = useState<MarketplaceAccount[]>([]);
  const [policies, setPolicies] = useState<MarketplacePolicy[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [selectedProviderCode, setSelectedProviderCode] = useState<string | null>(null);
  const [actionPending, setActionPending] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setLoadError(null);
    try {
      const [providersRes, installationsRes, accountsRes, policiesRes] = await Promise.all([
        client.listIntegrationProviders(),
        client.listIntegrationInstallations(),
        client.listMarketplaceAccounts(),
        client.listMarketplacePolicies(),
      ]);
      setProviders(providersRes.items);
      setInstallations(installationsRes.items);
      setAccounts(accountsRes.items);
      setPolicies(policiesRes.items);
    } catch (error) {
      const detail = error as { error?: { message?: string } };
      setLoadError(detail?.error?.message ?? "Failed to load provider catalog.");
    } finally {
      setLoading(false);
    }
  }, [client]);

  useEffect(() => {
    void load();
  }, [load]);

  const selectedProvider = useMemo(
    () => providers.find((provider) => provider.provider_code === selectedProviderCode) ?? null,
    [providers, selectedProviderCode],
  );

  const selectedInstallation = useMemo(() => {
    if (!selectedProvider) {
      return null;
    }
    return installations.find((item) => item.provider_code === selectedProvider.provider_code) ?? null;
  }, [selectedProvider, installations]);

  const policyCount = policies.length;
  const accountCount = accounts.length;

  async function ensureInstallation(provider: IntegrationProviderDefinition): Promise<IntegrationInstallation> {
    const existing = installations.find((item) => item.provider_code === provider.provider_code);
    if (existing) {
      return existing;
    }
    const created = await client.createIntegrationInstallation({
      installation_id: createInstallationID(provider.provider_code),
      provider_code: provider.provider_code,
      family: "marketplace",
      display_name: provider.display_name,
    });
    setInstallations((current) => [...current, created]);
    return created;
  }

  async function executeWithPending(action: () => Promise<void>) {
    setActionPending(true);
    try {
      await action();
      await load();
    } finally {
      setActionPending(false);
    }
  }

  async function handleConnect() {
    if (!selectedProvider) {
      return;
    }
    await executeWithPending(async () => {
      const created = await ensureInstallation(selectedProvider);
      const interactive = selectedProvider.install_mode === "interactive";
      if (interactive) {
        const result = await client.startIntegrationAuthorization(created.installation_id);
        if (result.auth_url) {
          navigateToAuthUrl(result.auth_url);
        }
      }
    });
  }

  async function handleAuthorize() {
    if (!selectedInstallation) {
      return;
    }
    await executeWithPending(async () => {
      const result = await client.startIntegrationAuthorization(selectedInstallation.installation_id);
      if (result.auth_url) {
        navigateToAuthUrl(result.auth_url);
      }
    });
  }

  async function handleReauthorize() {
    if (!selectedInstallation) {
      return;
    }
    await executeWithPending(async () => {
      const result = await client.startIntegrationReauthorization(selectedInstallation.installation_id);
      if (result.auth_url) {
        navigateToAuthUrl(result.auth_url);
      }
    });
  }

  async function handleSubmitCredentials(credentials: Record<string, string>) {
    if (!selectedProvider) {
      return;
    }
    await executeWithPending(async () => {
      const installation = await ensureInstallation(selectedProvider);
      const request: SubmitIntegrationCredentialsRequest = {
        credentials,
      };
      if (credentials.api_key) {
        request.api_key = credentials.api_key;
      }
      await client.submitIntegrationCredentials(installation.installation_id, request);
    });
  }

  async function handleDisconnect() {
    if (!selectedInstallation) {
      return;
    }
    await executeWithPending(async () => {
      await client.disconnectIntegrationInstallation(selectedInstallation.installation_id);
    });
  }

  async function handleFeeSync() {
    if (!selectedInstallation) {
      return;
    }
    await executeWithPending(async () => {
      await client.startIntegrationFeeSync(selectedInstallation.installation_id);
    });
  }

  return (
    <div className="min-h-full bg-slate-50">
      <div className="flex items-center justify-between px-6 pb-4 pt-6">
        <div>
          <h1 className="text-xl font-bold text-slate-900">Marketplaces</h1>
          <p className="mt-0.5 text-sm text-slate-500">
            Provider catalog with integration lifecycle and setup status
          </p>
        </div>
        <div className="rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm text-slate-600">
          Accounts: <span className="font-semibold text-slate-900">{accountCount}</span> · Policies:{" "}
          <span className="font-semibold text-slate-900">{policyCount}</span>
        </div>
      </div>

      {loadError && (
        <div className="mx-6 mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
          {loadError}{" "}
          <button type="button" className="ml-1 underline" onClick={() => void load()}>
            Retry
          </button>
        </div>
      )}

      <div
        className="px-6 pb-6 transition-all duration-200"
        style={{ paddingRight: selectedProvider ? "432px" : "24px" }}
      >
        {loading ? (
          <div className="grid gap-4" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))" }} aria-busy="true">
            <div className="h-36 animate-pulse rounded-lg border border-slate-200 bg-white" />
            <div className="h-36 animate-pulse rounded-lg border border-slate-200 bg-white" />
            <div className="h-36 animate-pulse rounded-lg border border-slate-200 bg-white" />
          </div>
        ) : providers.length === 0 ? (
          <div className="rounded-lg border border-dashed border-slate-300 bg-white px-6 py-8 text-sm text-slate-500">
            No providers available yet.
          </div>
        ) : (
          <div className="grid gap-4" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))" }}>
            {providers.map((provider) => {
              const installation = installations.find((item) => item.provider_code === provider.provider_code) ?? null;
              return (
                <ProviderCatalogCard
                  key={provider.provider_code}
                  provider={provider}
                  installation={installation}
                  onSelect={(selected) => setSelectedProviderCode(selected.provider_code)}
                />
              );
            })}
          </div>
        )}
      </div>

      {selectedProvider && (
        <ProviderCatalogPanel
          provider={selectedProvider}
          installation={selectedInstallation}
          pending={actionPending}
          onClose={() => setSelectedProviderCode(null)}
          onConnect={handleConnect}
          onAuthorize={handleAuthorize}
          onReauthorize={handleReauthorize}
          onSubmitCredentials={handleSubmitCredentials}
          onDisconnect={handleDisconnect}
          onFeeSync={handleFeeSync}
        />
      )}
    </div>
  );
}
