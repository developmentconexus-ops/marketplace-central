import { Link, NavLink, useLocation } from "react-router-dom";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { useInstallation } from "./InstallationContext";
import { useTheme } from "./theme/useTheme";

// Every Mercado Livre installation carries the same display_name, so the selector
// showed a list of identical entries and the operator could not tell which
// account was selected — or that some of them were abandoned authorizations with
// no data at all. Name the connected account, and say so when it is not connected.
function installationLabel(installation: IntegrationInstallation): string {
  const account = installation.external_account_name?.trim();
  if (account) return account;
  // Only claim "não conectada" when the API actually reported a non-connected
  // status; an absent status is unknown, and stamping it would be a fabrication.
  if (installation.status && installation.status !== "connected") {
    return `${installation.display_name} (não conectada)`;
  }
  return installation.display_name;
}

// An installation that never completed its authorization has no account, no
// listings and no orders — selecting it can only show empty screens. Those rows
// stay on /integracoes (where the operator resumes or gives up on them) but are
// kept out of the account selector. The selected one is always kept so the
// control can never lose its own value.
function selectableInstallations(
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

const enabledPillClass =
  "inline-flex items-center rounded-pill px-3 py-1.5 text-sm font-medium transition-colors";

export function Header() {
  const location = useLocation();
  const { toggleTheme } = useTheme();
  const { installationId, setInstallationId, installations, status } = useInstallation();
  const selectedInstallation = installations.find(
    (installation) => installation.installation_id === installationId,
  );

  return (
    <header className="flex h-14 shrink-0 items-center gap-4 overflow-x-auto border-b border-border bg-surface px-4 lg:px-6">
      <span className="shrink-0 text-sm font-semibold tracking-wide text-ink">
        Marketplace Central
      </span>

      <nav aria-label="Navegação principal" className="flex shrink-0 items-center gap-2 overflow-x-auto">
        <NavLink
          to={{ pathname: "/", search: location.search }}
          end
          className={({ isActive }) =>
            `${enabledPillClass} ${
              isActive
                ? "bg-accent-soft text-accent-ink"
                : "text-muted hover:bg-surface-2 hover:text-ink"
            }`
          }
        >
          Visão geral
        </NavLink>
        <NavLink
          to={{ pathname: "/anuncios", search: location.search }}
          className={({ isActive }) =>
            `${enabledPillClass} ${
              isActive
                ? "bg-accent-soft text-accent-ink"
                : "text-muted hover:bg-surface-2 hover:text-ink"
            }`
          }
        >
          Anúncios
        </NavLink>
        <NavLink
          to={{ pathname: "/mercado", search: location.search }}
          className={({ isActive }) =>
            `${enabledPillClass} ${
              isActive
                ? "bg-accent-soft text-accent-ink"
                : "text-muted hover:bg-surface-2 hover:text-ink"
            }`
          }
        >
          Mercado
        </NavLink>
        <NavLink
          to={{ pathname: "/precos", search: location.search }}
          className={({ isActive }) =>
            `${enabledPillClass} ${
              isActive
                ? "bg-accent-soft text-accent-ink"
                : "text-muted hover:bg-surface-2 hover:text-ink"
            }`
          }
        >
          Simulador
        </NavLink>
        <NavLink
          to={{ pathname: "/pedidos", search: location.search }}
          className={({ isActive }) =>
            `${enabledPillClass} ${
              isActive
                ? "bg-accent-soft text-accent-ink"
                : "text-muted hover:bg-surface-2 hover:text-ink"
            }`
          }
        >
          Pedidos
        </NavLink>
        <span className="inline-flex shrink-0 cursor-not-allowed items-center gap-1.5 rounded-pill px-3 py-1.5 text-sm font-medium text-faint">
          Repasses
          <span className="rounded-pill bg-surface-2 px-1.5 py-0.5 text-[10px] uppercase tracking-wide text-faint">
            em breve
          </span>
        </span>
      </nav>

      <div className="ml-auto flex shrink-0 items-center gap-2">
        <button
          type="button"
          aria-label="Alternar tema"
          className="rounded-pill border border-border bg-surface px-3 py-2 text-sm text-ink transition-colors hover:bg-surface-2"
          onClick={toggleTheme}
        >
          ◐
        </button>

        {status === "ready" && selectedInstallation ? (
          <span className="flex h-9 max-w-64 items-center gap-1 rounded-pill border border-border bg-surface px-3 text-sm font-medium text-ink focus-within:border-accent">
            ML:
            <select
              aria-label="Selecionar instalação"
              className="appearance-none bg-transparent text-ink outline-none"
              value={installationId}
              onChange={(event) => setInstallationId(event.target.value)}
            >
              {selectableInstallations(installations, installationId).map((installation) => (
                <option key={installation.installation_id} value={installation.installation_id}>
                  {installationLabel(installation)}
                </option>
              ))}
            </select>
            <span aria-hidden="true">▾</span>
          </span>
        ) : status === "empty" ? (
          <span className="text-xs text-muted">Conecte uma conta em Integrações</span>
        ) : (
          <span aria-hidden="true" className="h-9 w-24" />
        )}

        <details className="relative">
          <summary
            aria-label="Abrir menu"
            className="flex cursor-pointer items-center rounded-pill border border-border bg-surface px-3 py-2 text-sm text-ink"
          >
            ⚙
          </summary>
          <div
            role="group"
            aria-label="Menu de configurações"
            className="absolute right-0 top-full z-10 mt-2 w-48 rounded-card border border-border bg-surface p-2 shadow"
          >
            <div className="text-sm text-ink">
              <div className="px-3 py-2 font-medium">Configurações</div>
              <Link
                to="/precos?params=1"
                className="block rounded-control px-3 py-2 text-ink hover:bg-surface-2"
              >
                DIFAL
              </Link>
            </div>
            <Link
              to="/integracoes"
              className="block rounded-control px-3 py-2 text-ink hover:bg-surface-2"
            >
              Integrações
            </Link>
            <Link
              to="/catalogo"
              className="block rounded-control px-3 py-2 text-ink hover:bg-surface-2"
            >
              Catálogo
            </Link>
            <Link
              to="/estoque"
              className="block rounded-control px-3 py-2 text-ink hover:bg-surface-2"
            >
              Estoque
            </Link>
          </div>
        </details>
      </div>
    </header>
  );
}
