import { NavLink, Outlet } from "react-router-dom";
import { LayoutDashboard, Package, Tags, Store, Calculator, ActivitySquare, Link2, ShieldCheck, ReceiptText } from "lucide-react";

const navItems = [
  { to: "/",                 label: "Dashboard",         icon: LayoutDashboard },
  { to: "/products",         label: "Catalog",           icon: Package },
  { to: "/classifications",  label: "Classifications",   icon: Tags },
  { to: "/marketplaces",     label: "Marketplaces",      icon: Store },
  { to: "/integrations",     label: "Integrations",      icon: ActivitySquare },
  { to: "/product-links",    label: "Product Links",     icon: Link2 },
  { to: "/inventory/stock-seguro", label: "Stock Seguro", icon: ShieldCheck },
  { to: "/orders",           label: "Orders",            icon: ReceiptText },
  { to: "/simulator",        label: "Pricing Simulator", icon: Calculator },
];

export function Layout() {
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
              to={to}
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
        </header>
        <main className="flex-1 overflow-auto p-4 lg:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
