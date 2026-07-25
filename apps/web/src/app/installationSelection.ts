import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";

// An installation that never completed its authorization has no account, no
// listings and no orders — selecting it can only show empty screens. Those rows
// stay on /integracoes (where the operator resumes or gives up on them) but are
// kept out of every account selector. The selected one is always kept so the
// control can never lose its own value.
//
// Shared by the header selector and the per-screen account selectors (Stock
// Seguro), so the operator sees ONE account list across the workspace.
export function selectableInstallations(
  installations: IntegrationInstallation[],
  selectedId: string,
): IntegrationInstallation[] {
  const neverAuthorized = new Set(["draft", "pending_connection"]);
  return installations.filter(
    (installation) =>
      installation.installation_id === selectedId ||
      !installation.status ||
      !neverAuthorized.has(installation.status),
  );
}
