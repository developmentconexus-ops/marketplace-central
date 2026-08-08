import { Link } from "react-router-dom";

interface Atalho {
  key: string;
  label: string;
  href: string;
}

// Static shortcuts per the design screen's rail. Targets verified against the routes actually
// registered in AppRouter.tsx (no /simulador — the pricing route is /precos). The spreadsheet
// upload lives on /integracoes; /vinculos only shows what a completed import produced, so
// pointing "Importar planilha" there left the operator on a screen with no way to import.
const atalhos: Atalho[] = [
  { key: "importar", label: "Importar planilha", href: "/integracoes" },
  { key: "simular", label: "Simular", href: "/precos" },
  { key: "vinculos", label: "Vínculos", href: "/vinculos" },
];

export function Atalhos() {
  return (
    <section
      aria-labelledby="atalhos-title"
      className="rounded-card border border-border bg-surface p-4"
    >
      <h2 id="atalhos-title" className="text-sm font-semibold text-ink">
        Atalhos
      </h2>
      <ul className="mt-2 flex flex-col gap-2">
        {atalhos.map((atalho) => (
          <li key={atalho.key}>
            <Link
              to={atalho.href}
              className="block rounded-lg px-2 py-1.5 text-sm font-medium text-accent-ink transition-colors hover:bg-accent-soft focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            >
              {atalho.label}
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
