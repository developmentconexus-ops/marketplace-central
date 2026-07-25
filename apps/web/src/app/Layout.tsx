import { Outlet } from "react-router-dom";
import { Header } from "./Header";

// The installation gate lives on the routes that need a connected marketplace
// account, not on the shell: gating the shell also gated /integracoes, the one
// screen where an account gets connected, so a workspace with no account had no
// way to ever get one. See AppRouter.
export function Layout() {
  return (
    <div className="flex min-h-screen flex-col bg-bg text-ink">
      <Header />
      <main className="flex-1 overflow-auto p-4 lg:p-6">
        <Outlet />
      </main>
    </div>
  );
}
