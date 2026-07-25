import type { DashboardOverview } from "@marketplace-central/sdk-runtime";
import { Link } from "react-router-dom";

interface AttentionItem {
  key: string;
  label: string;
  count: number;
  href: string;
}

// Only items with a real deep-link target belong in the fila (ruling #2). Pedidos-atrasados is
// deliberately excluded: /pedidos has no URL filter contract to land on (verified in
// anunciosQueryState.ts / AppRouter.tsx — no equivalent parser exists for /pedidos). missing_gtin
// has no deep-link target either, so it is also excluded here (still shown as a KPI card).
function buildItems(summary: DashboardOverview): AttentionItem[] {
  const items: AttentionItem[] = [];
  if (summary.sync_errors !== null && summary.sync_errors > 0) {
    items.push({
      key: "sync_errors",
      label: "Anúncios com erro de sincronização",
      count: summary.sync_errors,
      href: "/anuncios?filter.exception=sync_error",
    });
  }
  if (summary.below_margin !== null && summary.below_margin > 0) {
    items.push({
      key: "below_margin",
      label: "Anúncios abaixo da margem",
      count: summary.below_margin,
      href: "/anuncios?filter.exception=below_margin",
    });
  }
  if (summary.pending_links !== null && summary.pending_links > 0) {
    items.push({
      key: "pending_links",
      label: "Anúncios sem vínculo",
      count: summary.pending_links,
      // Bare /vinculos defaults to the Pendentes/fila tab (verified in vinculos route).
      href: "/vinculos",
    });
  }
  return items;
}

export function FilaDeAtencao({ summary }: { summary: DashboardOverview }) {
  const items = buildItems(summary);

  return (
    <section aria-labelledby="fila-atencao-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="fila-atencao-title" className="text-sm font-semibold text-ink">
        Fila de atenção
      </h2>
      {items.length === 0 ? (
        <p className="mt-2 text-sm text-muted">Nenhum item na fila — tudo em dia.</p>
      ) : (
        <ul className="mt-2 flex flex-col gap-2">
          {items.map((item) => (
            <li key={item.key}>
              <Link
                to={item.href}
                className="flex items-center justify-between rounded-lg px-2 py-1.5 text-sm text-ink transition-colors hover:bg-accent-soft focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              >
                <span>{item.label}</span>
                <span className="font-semibold" style={{ fontFamily: "var(--font-mono)" }}>
                  {item.count}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
