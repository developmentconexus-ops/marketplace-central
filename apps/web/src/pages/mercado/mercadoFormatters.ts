import type { ListingMoney, MarketPriceIntelMoney } from "@marketplace-central/sdk-runtime";

/** The em dash is UnknownValue's glyph — an honest unknown, never a fabricated 0/placeholder (ADR-17). */
export const DASH = "—";

/** The operator's minimum-margin policy — a static label in the radar header ("Margem mín: 18%").
 *  It is NOT used to color or rank anything here; operational margin is M-07-owned (ADR-17). */
export const MARGIN_MIN_PCT = 18;

type Money = ListingMoney | MarketPriceIntelMoney | null | undefined;

/**
 * "47.50" → "R$ 47,50" (pt-BR). Returns null for a missing/unparseable amount so the
 * caller renders an honest dash instead of "R$ NaN" (ADR-17).
 */
export function formatMoney(value: Money): string | null {
  if (!value || value.amount == null || value.amount === "") return null;
  const n = Number(value.amount);
  if (!Number.isFinite(n)) return null;
  return `R$ ${n.toLocaleString("pt-BR", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

/** { rank: 9, total: 14 } → "9º/14". null → dash (ADR-17: no fabricated rank). */
export function formatPosition(
  pos: { rank: number; total: number } | null | undefined,
): string | null {
  if (!pos || pos.rank == null || pos.total == null) return null;
  return `${pos.rank}º/${pos.total}`;
}

/** Format the freshest collection timestamp honestly; zero-date / missing → dash (ADR-17). */
export function formatCollectedAt(ts: string | null | undefined): string | null {
  if (!ts || ts.trim() === "" || ts.startsWith("0001-01-01T00:00:00")) return null;
  const d = new Date(ts);
  if (Number.isNaN(d.getTime())) return null;
  return d.toLocaleString("pt-BR", {
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}
