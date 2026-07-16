import { NavLink, Outlet, useLocation } from "react-router-dom";
import {
  ActivitySquare,
  Calculator,
  LayoutDashboard,
  Link2,
  Package,
  ReceiptText,
  ShieldCheck,
  Tags,
} from "lucide-react";
import { InstallationGate, useInstallation } from "./InstallationContext";

const navItems = [
  { to: "/", label: "Visão geral", icon: LayoutDashboard },
  { to: "/catalogo", label: "Catálogo", icon: Package },
  { to: "/anuncios", label: "Anúncios", icon: Tags },
  { to: "/vinculos", label: "Vínculos & Import.", icon: Link2 },
  { to: "/estoque", label: "Estoque", icon: ShieldCheck },
  { to: "/precos", label: "Preços & Simulador", icon: Calculator },
  { to: "/pedidos", label: "Pedidos", icon: ReceiptText },
  { to: "/integracoes", label: "Integrações & Sync", icon: ActivitySquare },
];

export function Layout() {
  const location = useLocation();
  const { installationId, setInstallationId, installations, status } = useInstallation();
  const selectedInstallation = installations.find(
    (installation) => installation.installation_id === installationId,
  );

  return (
    <div className="flex min-h-screen flex-col overflow-hidden lg:h-screen lg:flex-row">
      {/* Sidebar */}
      <aside className="shrink-0 flex flex-col border-b border-slate-700 lg:w-60 lg:border-b-0 lg:border-r" style={{ backgroundColor: "#0F172A" }}>
        <div className="px-5 py-5 border-b border-slate-700">
          <span className="text-white font-semibold text-sm tracking-wide">Marketplace Central</span>
        </div>

        <nav className="flex-1 overflow-x-auto px-3 py-3 lg:overflow-visible lg:py-4">
          <div className="flex gap-1 lg:block lg:space-y-1">
          {navItems.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={{ pathname: to, search: location.search }}
              end={to === "/"}
              className={({ isActive }) =>
                `flex min-w-fit items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors duration-150 lg:min-w-0 ${
                  isActive
                    ? "bg-blue-600 text-white"
                    : "text-slate-400 hover:bg-slate-800 hover:text-white"
                }`
              }
            >
              <Icon size={16} />
              {label}
            </NavLink>
          ))}
          </div>
        </nav>

        <div className="hidden px-5 py-4 border-t border-slate-700 lg:block">
          <p className="text-xs text-slate-500">v0.1.0</p>
        </div>
      </aside>

      {/* Main area */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <header className="flex h-14 shrink-0 items-center border-b border-slate-200 bg-white px-4 lg:px-6">
          <h1 className="text-sm font-medium text-slate-700">Marketplace Central</h1>
          <div className="ml-auto flex items-center gap-2">
            {status === "ready" && selectedInstallation ? (
              <span className="flex h-9 max-w-64 items-center gap-1 rounded-full border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 shadow-sm transition-colors focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-100 hover:border-slate-300">
                ML:
                <select
                  aria-label="Selecionar instalação"
                  className="appearance-none bg-transparent outline-none"
                  value={installationId}
                  onChange={(event) => setInstallationId(event.target.value)}
                >
                  {installations.map((installation) => (
                    <option key={installation.installation_id} value={installation.installation_id}>
                      {installation.display_name}
                    </option>
                  ))}
                </select>
                <span aria-hidden="true">▾</span>
              </span>
            ) : status === "empty" ? (
              <span className="text-xs text-slate-500">Conecte uma conta em Integrações</span>
            ) : (
              <span aria-hidden="true" className="h-9 w-24" />
            )}
            <span className="rounded-full border border-slate-200 bg-slate-50 px-3 py-2 text-xs font-medium text-slate-600">
              Empresa
            </span>
          </div>
        </header>
        <main className="flex-1 overflow-auto p-4 lg:p-6">
          <InstallationGate>
            <Outlet />
          </InstallationGate>
        </main>
      </div>
    </div>
  );
}
