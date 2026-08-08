import type { JSX } from "react";
import { FreshnessIndicator, UnknownValue } from "@marketplace-central/ui";
import type { OppRow } from "./oportunidades";
import { formatMoney } from "./mercadoFormatters";

// Oportunidades grid — column track from Mercado.dc.html (min-w 960) + MENOR CONC.
// added per operator D-120: SKU | PRODUTO (ERP) | NOSSO PREÇO (= custo ERP; produto não
// vendido, sem preço de anúncio — label per operator) | MENOR CONC. | MEDIANA ML |
// CONCORRENTES | VENDAS LÍDER 30D | MARGEM EST. | VEREDICTO | ATUALIZADO (F-A3/D-50) | (action)
const GRID_COLS =
  "84px minmax(170px,1.3fr) 76px 90px 90px 110px 120px 110px minmax(150px,1fr) 110px 120px";

const HEAD = [
  "SKU",
  "PRODUTO (ERP)",
  "NOSSO PREÇO",
  "MENOR CONC.",
  "MEDIANA ML",
  "CONCORRENTES",
  "VENDAS LÍDER 30D",
  "MARGEM EST.",
  "VEREDICTO",
  "ATUALIZADO",
  "",
];

const EVIDENCE_NOTE: Record<string, string> = {
  OK: "mercado observado",
  INSUFFICIENT_MARKET: "mercado insuficiente",
  NO_PRICE_EVIDENCE: "sem evidência de preço",
};

export interface OportunidadesTableProps {
  rows: OppRow[];
}

/**
 * The radar's "Oportunidades" tab: ERP catalog products with observed ML demand.
 * CUSTO (ERP fact), MEDIANA ML + CONCORRENTES (market aggregate) are real. MARGEM EST.,
 * VENDAS LÍDER 30D and the recommendation label are NOT backed by any endpoint — the
 * operational margin is M-07-owned (a naive cost-vs-median % omits ML commission, freight
 * and DIFAL, so it is never shown as a decision) — all three render honest "—" (ADR-17).
 * "Criar anúncio" is an inert demo affordance — no live ML write (D-57).
 */
export function OportunidadesTable({ rows }: OportunidadesTableProps): JSX.Element {
  return (
    <>
      <div className="overflow-x-auto rounded-card border border-border bg-surface">
        <div style={{ minWidth: 1050 }}>
          <div
            className="grid bg-surface-2 px-[14px] py-[9px] text-[11px] tracking-[0.04em] text-faint"
            style={{ gridTemplateColumns: GRID_COLS }}
          >
            {HEAD.map((h, i) => (
              <span key={h || `blank-${i}`}>{h}</span>
            ))}
          </div>
          {rows.map((o) => {
            const custo = formatMoney(
              o.costAmount ? { amount: o.costAmount, currency: "BRL" } : null,
            );
            const mediana = formatMoney(o.median);
            const menor = formatMoney(o.minValid);
            return (
              <div
                key={o.sku}
                className="grid items-center border-t border-border-2 px-[14px] py-[10px] text-[12.5px]"
                style={{ gridTemplateColumns: GRID_COLS }}
              >
                <span className="font-mono text-[11.5px] text-faint">{o.sku}</span>
                <span className="overflow-hidden text-ellipsis whitespace-nowrap pr-2 font-medium">
                  {o.name ?? <UnknownValue />}
                </span>
                <span className="font-mono text-muted">
                  {custo === null ? <UnknownValue /> : custo}
                </span>
                <span className="font-mono">{menor === null ? <UnknownValue /> : menor}</span>
                <span className="font-mono">{mediana === null ? <UnknownValue /> : mediana}</span>
                {/* CONCORRENTES = distinct competing sellers (n_sellers, deduped by seller per
                    the IC-03 aggregate contract), NOT the raw offer count — never overstate
                    competition (ADR-17). */}
                <span className="text-muted">
                  {o.nSellers == null ? <UnknownValue /> : o.nSellers}
                </span>
                {/* VENDAS LÍDER 30D — no backing snapshot → honest dash (ADR-17). */}
                <span className="font-mono text-muted">
                  <UnknownValue hint="vendas do líder — sem snapshot" />
                </span>
                {/* MARGEM EST. — operational margin is M-07-owned (commission/freight/DIFAL
                    not in any endpoint) → honest dash, never a decision-colored estimate (ADR-17). */}
                <span className="font-mono text-muted">
                  <UnknownValue hint="margem — M-07" />
                </span>
                {/* VEREDICTO — recommendation label is M-07-owned → honest dash + evidence note. */}
                <span className="pr-2 text-[12px]">
                  <UnknownValue hint="recomendação — M-07" />{" "}
                  <span className="text-[11px] text-faint">
                    {o.evidenceState ? (EVIDENCE_NOTE[o.evidenceState] ?? "") : ""}
                  </span>
                </span>
                <FreshnessIndicator asOf={o.fetchedAt} />
                <button
                  type="button"
                  disabled
                  title="criar anúncio — em breve"
                  className="cursor-not-allowed justify-self-start whitespace-nowrap rounded-[7px] border border-border bg-surface px-[10px] py-1 text-[11.5px] font-semibold text-muted"
                >
                  Criar anúncio
                </button>
              </div>
            );
          })}
        </div>
      </div>
      <p className="text-[11.5px] text-faint">
        cruza catálogo ERP × demanda pública do ML · ordenado por diferença mediana − custo ·
        margem, vendas do líder e recomendação chegam com M-07/snapshots
      </p>
    </>
  );
}
