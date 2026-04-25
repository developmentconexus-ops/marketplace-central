import type { IntegrationInstallation, IntegrationProviderDefinition } from "@marketplace-central/sdk-runtime";
import { MarketplaceIcon } from "./components/MarketplaceIcon";
import { StatusBadge } from "./components/StatusBadge";

interface ProviderCatalogCardProps {
  provider: IntegrationProviderDefinition;
  installation: IntegrationInstallation | null;
  onSelect: (provider: IntegrationProviderDefinition) => void;
}

function resolveStatus(provider: IntegrationProviderDefinition, installation: IntegrationInstallation | null): string {
  const executionMode = provider.metadata?.execution_mode ?? "available";
  if (executionMode === "blocked") {
    return "blocked";
  }
  if (installation) {
    return installation.status;
  }
  if (executionMode === "planned") {
    return "planned";
  }
  return "available";
}

function capabilitiesPreview(capabilities: string[]): string {
  if (capabilities.length === 0) {
    return "No declared capabilities";
  }
  return capabilities.slice(0, 3).join(", ");
}

export function ProviderCatalogCard({ provider, installation, onSelect }: ProviderCatalogCardProps) {
  const status = resolveStatus(provider, installation);
  const baseline = provider.metadata?.baseline_commission_percent;
  const baselineText = typeof baseline === "number" ? `${(baseline * 100).toFixed(1)}%` : "n/a";

  return (
    <button
      type="button"
      onClick={() => onSelect(provider)}
      className="w-full rounded-lg border border-slate-200 bg-white p-4 text-left shadow-sm transition-all hover:border-slate-300 hover:shadow-md"
      aria-label={`Open ${provider.display_name}`}
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <MarketplaceIcon code={provider.provider_code} size={30} />
          <div>
            <p className="text-sm font-semibold text-slate-900">{provider.display_name}</p>
            <p className="text-xs text-slate-500">{provider.auth_strategy}</p>
          </div>
        </div>
        <StatusBadge status={status} />
      </div>
      <div className="mt-3 text-xs text-slate-600">
        <p>Capabilities: {capabilitiesPreview(provider.declared_capabilities)}</p>
        <p className="mt-1">Baseline commission: {baselineText}</p>
      </div>
    </button>
  );
}
