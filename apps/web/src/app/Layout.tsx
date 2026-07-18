import { Outlet } from "react-router-dom";
import { InstallationGate } from "./InstallationContext";
import { Header } from "./Header";

export function Layout() {
  return (
    <div className="flex min-h-screen flex-col bg-bg text-ink">
      <Header />
      <main className="flex-1 overflow-auto p-4 lg:p-6">
        <InstallationGate>
          <Outlet />
        </InstallationGate>
      </main>
    </div>
  );
}
