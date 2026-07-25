import type { JSX } from "react";
import { Link } from "react-router-dom";
import type { ListingReadModel } from "@marketplace-central/sdk-runtime";
import { UnknownValue } from "@marketplace-central/ui";
import { DASH, formatMoney, formatPosition } from "./mercadoFormatters";

// Reprecificação grid — exact column track + min-width from Mercado.dc.html (min-w 1020):
// ANÚNCIO | MEU PREÇO | MENOR CONC. | MEDIANA | POSIÇÃO | VENDAS 30D | MARGEM ATUAL |
// SE IGUALAR MENOR | SUGESTÃO | (action)
const GRID_COLS =
  "minmax(170px,1.3fr) 84px 84px 84px 66px 76px 90px 110px minmax(150px,1fr) 120px";

const HEAD = [
  "ANÚNCIO",
  "MEU PREÇO",
  "MENOR CONC.",
  "MEDIANA",
  "POSIÇÃO",
  "VENDAS 30D",
  "MARGEM ATUAL",
  "SE IGUALAR MENOR",
  "SUGESTÃO",
  "",
];

export interface RepricingTableProps {
  rows: ListingReadModel[];
}

/**
 * Money cell — honest dash when the amount is missing (ADR-17), never "R$ —".
 */
function MoneyCell({ value, muted = false }: { value: string | null; muted?: boolean }): JSX.Element {
  if (value === null) return <UnknownValue />;
  return <span className={`font-mono ${muted ? "text-muted" : ""}`}>{value}</span>;
}

/**
 * The radar's "Reprecificação" tab: one row per OUR listing. Every value is a
 * real ListingReadModel field or its market_signal — MEU PREÇO (price), MENOR
 * CONC. (signal.min_valid), MEDIANA (signal.median), POSIÇÃO (signal.position),
 * VENDAS 30D (sales_30d). MARGEM ATUAL / SE IGUALAR MENOR / SUGESTÃO have no
 * backing field (the margin + repricing engine is M-07-owned), so they render
 * the honest UnknownValue, never a fabricated margin or suggestion (ADR-17).
 * "Aplicar" / "Simular" are inert demo affordances — no live ML write (D-57).
 */
export function RepricingTable({ rows }: RepricingTableProps): JSX.Element {
  return (
    <>
      <div className="overflow-x-auto rounded-card border border-border bg-surface">
        <div style={{ minWidth: 1020 }}>
          <div
            className="grid bg-surface-2 px-[14px] py-[9px] text-[11px] tracking-[0.04em] text-faint"
            style={{ gridTemplateColumns: GRID_COLS }}
          >
            {HEAD.map((h, i) => (
              <span key={h || `blank-${i}`}>{h}</span>
            ))}
          </div>
          {rows.map((r) => {
            const signal = r.market_signal ?? null;
            const meu = formatMoney(r.price);
            const menor = formatMoney(signal?.min_valid ?? null);
            const mediana = formatMoney(signal?.median ?? null);
            const pos = formatPosition(signal?.position ?? null);
            const leader = signal?.position?.rank === 1;
            return (
              <div
                key={r.listing_id}
                className="grid items-center border-t border-border-2 px-[14px] py-[10px] text-[12.5px]"
                style={{ gridTemplateColumns: GRID_COLS }}
              >
                {/* A variação é o que distingue linhas que compartilham MLB e título: sem ela um
                    anúncio com 4 variações lê como 4 linhas idênticas e o operador não sabe qual
                    reprecificar. Mesmo rótulo de AnunciosTable para os dois lerem igual. */}
                <span className="overflow-hidden text-ellipsis whitespace-nowrap pr-2">
                  <span className="font-mono text-[11px] text-faint">{r.provider_listing_id}</span>{" "}
                  {r.title}
                  {r.variation_id ? (
                    <span className="ml-1 font-mono text-[11px] text-faint">var. {r.variation_id}</span>
                  ) : null}
                </span>
                <MoneyCell value={meu} />
                <MoneyCell value={menor} muted />
                <MoneyCell value={mediana} muted />
                <span className={leader ? "font-bold" : ""}>{pos ?? <UnknownValue />}</span>
                <span className="font-mono text-muted">
                  {r.sales_30d == null ? <UnknownValue /> : r.sales_30d}
                </span>
                {/* MARGEM ATUAL — no field; margin engine is M-07-owned → honest dash (ADR-17). */}
                <span className="justify-self-start rounded-pill bg-surface-2 px-[9px] py-px text-[11px] font-semibold text-faint">
                  <UnknownValue hint="margem — M-07" />
                </span>
                {/* SE IGUALAR MENOR — no field → honest dash. */}
                <span className="justify-self-start rounded-pill bg-surface-2 px-[9px] py-px text-[11px] font-semibold text-faint">
                  <UnknownValue hint="simulação de margem — M-07" />
                </span>
                {/* SUGESTÃO — repricing recommendation is M-07-owned → honest dash. */}
                <span className="pr-2 text-[12px]">
                  <UnknownValue hint="sugestão de preço — M-07" />
                </span>
                <span className="flex gap-[5px]">
                  <button
                    type="button"
                    disabled
                    title="ajuste de preço via fila de sync — em breve"
                    className="cursor-not-allowed whitespace-nowrap rounded-[7px] bg-accent px-[10px] py-1 text-[11.5px] font-semibold text-white opacity-60"
                  >
                    Aplicar
                  </button>
                  {/* "Simular" hands the linked ERP product to the /precos
                      simulator. A listing with no resolved product link has
                      nothing to simulate, so the control stays inert there
                      rather than opening the simulator on the wrong product. */}
                  {r.link.product_id ? (
                    <Link
                      to={`/precos?produto=${encodeURIComponent(r.link.product_id)}`}
                      title="abrir no simulador de preços"
                      className="whitespace-nowrap rounded-[7px] border border-border px-[10px] py-1 text-[11.5px] text-ink hover:bg-surface-2"
                    >
                      Simular
                    </Link>
                  ) : (
                    <button
                      type="button"
                      disabled
                      title="anúncio sem produto vinculado — nada a simular"
                      className="cursor-not-allowed whitespace-nowrap rounded-[7px] border border-border px-[10px] py-1 text-[11.5px] text-muted"
                    >
                      Simular
                    </button>
                  )}
                </span>
              </div>
            );
          })}
        </div>
      </div>
      <p className="text-[11.5px] text-faint">
        margem, "se igualar menor" e sugestão dependem do motor de precificação (M-07){DASH}exibidos
        como "{DASH}" até então · "Aplicar" mudará o preço via fila de sync com preview e protocolo ·
        posição vem da busca pública do ML
      </p>
    </>
  );
}
