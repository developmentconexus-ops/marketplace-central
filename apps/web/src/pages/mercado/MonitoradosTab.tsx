import type { JSX } from "react";

// Monitorados grid — exact column track + min-width from Mercado.dc.html (min-w 920):
// TIPO | MONITORADO | LEITURA ATUAL | VARIAÇÃO 7D | ALERTA | CHECADO | (action)
const GRID_COLS =
  "96px minmax(180px,1.4fr) minmax(150px,1.1fr) 120px 130px 100px 96px";

const HEAD = ["TIPO", "MONITORADO", "LEITURA ATUAL", "VARIAÇÃO 7D", "ALERTA", "CHECADO", ""];

// Filter chips per design. Inert: monitoring persists in a daily-snapshot pipeline that is
// not wired yet (API-MAP: Monitorados/variação-7d = snapshots próprios, sem retroativo), so
// there is nothing to filter and no watchlist to mutate — never fabricate monitored rows (ADR-17).
const FILTERS = [
  { key: "todos", label: "Todos" },
  { key: "vendedor", label: "Vendedores" },
  { key: "termo", label: "Termos" },
  { key: "anuncio", label: "Anúncios" },
];

/**
 * The radar's "Monitorados" tab. Watchlists + "variação 7d" require a daily
 * snapshot store with no retroactive history — that pipeline is not built, and
 * there is no endpoint to read monitored entities. Rather than fabricate seller
 * / term / listing rows and 7-day deltas (ADR-17), the tab renders its design
 * chrome (filter chips, "+ Monitorar…", table header, legend) with an honest
 * empty state. All controls are inert demo affordances — no live ML write (D-57).
 */
export function MonitoradosTab(): JSX.Element {
  return (
    <>
      <div className="flex flex-wrap items-center gap-2">
        {FILTERS.map((f) => {
          const on = f.key === "todos";
          return (
            <button
              key={f.key}
              type="button"
              disabled
              title="monitoramento — em breve"
              className={`cursor-not-allowed whitespace-nowrap rounded-pill border-[1.5px] px-[13px] py-1 text-[12px] font-semibold ${
                on
                  ? "border-accent bg-accent-soft text-accent-ink"
                  : "border-border bg-transparent text-muted"
              }`}
            >
              {f.label}
            </button>
          );
        })}
        <button
          type="button"
          disabled
          title="adicionar monitoramento — em breve"
          className="ml-auto cursor-not-allowed whitespace-nowrap rounded-control bg-accent px-[14px] py-[7px] text-[12.5px] font-semibold text-white opacity-60"
        >
          + Monitorar…
        </button>
      </div>
      <div className="overflow-x-auto rounded-card border border-border bg-surface">
        <div style={{ minWidth: 920 }}>
          <div
            className="grid bg-surface-2 px-[14px] py-[9px] text-[11px] tracking-[0.04em] text-faint"
            style={{ gridTemplateColumns: GRID_COLS }}
          >
            {HEAD.map((h, i) => (
              <span key={h || `blank-${i}`}>{h}</span>
            ))}
          </div>
          <div className="flex flex-col gap-1 border-t border-border-2 px-[14px] py-[18px] text-sm text-muted">
            <p>Nenhum monitoramento ativo.</p>
            <p className="text-xs text-faint">
              vendedores, termos e anúncios monitorados dependem do snapshot diário do ML (sem
              histórico retroativo) — disponível em breve
            </p>
          </div>
        </div>
      </div>
      <p className="text-[11.5px] text-faint">
        vendedores, termos de busca e anúncios específicos · leitura pública do ML 1×/dia · alertas
        viram itens do radar de Reprecificação quando cruzam o limite
      </p>
    </>
  );
}
