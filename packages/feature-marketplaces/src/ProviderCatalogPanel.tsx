import { useMemo, useState } from "react";
import type { CredentialField, IntegrationInstallation, IntegrationProviderDefinition } from "@marketplace-central/sdk-runtime";
import { X } from "lucide-react";
import { StatusBadge } from "./components/StatusBadge";

interface ProviderCatalogPanelProps {
  provider: IntegrationProviderDefinition | null;
  installation: IntegrationInstallation | null;
  pending: boolean;
  onClose: () => void;
  onConnect: () => Promise<void>;
  onAuthorize: () => Promise<void>;
  onReauthorize: () => Promise<void>;
  onSubmitCredentials: (credentials: Record<string, string>) => Promise<void>;
  onDisconnect: () => Promise<void>;
  onFeeSync: () => Promise<void>;
}

export function ProviderCatalogPanel({
  provider,
  installation,
  pending,
  onClose,
  onConnect,
  onAuthorize,
  onReauthorize,
  onSubmitCredentials,
  onDisconnect,
  onFeeSync,
}: ProviderCatalogPanelProps) {
  const [credentials, setCredentials] = useState<Record<string, string>>({});

  const credentialSchema = useMemo(() => {
    if (!provider) {
      return [];
    }
    const rawSchema: unknown = provider.metadata?.credential_schema;
    if (!Array.isArray(rawSchema)) {
      return [];
    }
    return rawSchema.filter((item): item is CredentialField => {
      if (typeof item !== "object" || item === null) {
        return false;
      }
      const value = item as Record<string, unknown>;
      return typeof value.key === "string" && typeof value.label === "string";
    });
  }, [provider]);

  if (!provider) {
    return null;
  }

  const status = installation?.status ?? String(provider.metadata?.execution_mode ?? "available");
  const blocked = provider.metadata?.execution_mode === "blocked";
  const installMode = provider.install_mode;
  const noInstallation = installation === null;
  const feeSyncEnabled = installation?.status === "connected" || installation?.status === "degraded";

  async function handlePrimaryAction() {
    if (blocked) {
      return;
    }
    if (noInstallation) {
      await onConnect();
      return;
    }
    if (installation.status === "pending_connection" || installation.status === "disconnected") {
      await onAuthorize();
      return;
    }
    if (installation.status === "requires_reauth") {
      await onReauthorize();
      return;
    }
  }

  const primaryLabel = blocked
    ? "Unavailable"
    : noInstallation
      ? "Create installation"
      : installation.status === "pending_connection" || installation.status === "disconnected"
        ? "Authorize"
        : installation.status === "requires_reauth"
          ? "Reauthorize"
          : "Open";

  return (
    <aside
      role="dialog"
      aria-label="Provider details"
      className="fixed right-0 top-0 h-full w-[400px] overflow-y-auto border-l border-slate-200 bg-white p-5 shadow-2xl"
    >
      <div className="mb-4 flex items-start justify-between">
        <div>
          <p className="text-lg font-semibold text-slate-900">{provider.display_name}</p>
          <p className="text-xs text-slate-500">{provider.provider_code}</p>
        </div>
        <button type="button" onClick={onClose} aria-label="Close panel" className="rounded-lg p-2 text-slate-500 hover:bg-slate-100">
          <X className="h-4 w-4" />
        </button>
      </div>

      <div className="space-y-3">
        <StatusBadge status={status} />
        <p className="text-sm text-slate-600">
          Auth strategy: <span className="font-medium text-slate-900">{provider.auth_strategy}</span>
        </p>
        <p className="text-sm text-slate-600">
          Install mode: <span className="font-medium text-slate-900">{installMode}</span>
        </p>
        {blocked && (
          <p className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">
            {String(provider.metadata?.unavailable_reason ?? "Provider currently unavailable.")}
          </p>
        )}
      </div>

      {!blocked && installMode !== "interactive" && credentialSchema.length > 0 && (
        <div className="mt-5 space-y-3">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-500">Credentials</p>
          {credentialSchema.map((field) => (
            <div key={field.key}>
              <label className="mb-1 block text-xs text-slate-600">{field.label}</label>
              <input
                type={field.secret ? "password" : "text"}
                value={credentials[field.key] ?? ""}
                onChange={(event) => setCredentials((current) => ({ ...current, [field.key]: event.target.value }))}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm"
              />
            </div>
          ))}
          <button
            type="button"
            onClick={() => onSubmitCredentials(credentials)}
            disabled={pending || noInstallation}
            className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white disabled:opacity-40"
          >
            Submit credentials
          </button>
        </div>
      )}

      <div className="mt-6 grid grid-cols-2 gap-2">
        <button
          type="button"
          onClick={handlePrimaryAction}
          disabled={pending || blocked || (installation?.status === "connected" || installation?.status === "degraded")}
          className="rounded-lg bg-orange-500 px-3 py-2 text-sm font-semibold text-white disabled:opacity-40"
        >
          {primaryLabel}
        </button>
        <button
          type="button"
          onClick={onFeeSync}
          disabled={pending || !feeSyncEnabled}
          className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-semibold text-slate-700 disabled:opacity-40"
        >
          Run fee sync
        </button>
        <button
          type="button"
          onClick={onDisconnect}
          disabled={pending || noInstallation}
          className="col-span-2 rounded-lg border border-red-300 px-3 py-2 text-sm font-semibold text-red-700 disabled:opacity-40"
        >
          Disconnect
        </button>
      </div>
    </aside>
  );
}
