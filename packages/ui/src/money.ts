// One money formatter for every screen. Before this existed each screen chose:
// /pedidos and /mercado formatted through Intl in pt-BR, while /anuncios, the
// listing drawer, the pricing chips and the classifications table interpolated
// the raw API string, so the same value rendered as "R$ 53,90" on one screen and
// "R$ 53.9" on the next. The API sends decimal STRINGS (never floats) to keep
// cents exact, so the input type is string | number.

const brl = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL",
});

/**
 * formatMoney renders a monetary amount in pt-BR. An amount that is absent or
 * not a number returns null: the caller decides how to render unknown (usually
 * <UnknownValue />), because a fabricated "R$ 0,00" is a lie about the data
 * (ADR-17 honest-unknown).
 */
export function formatMoney(amount: string | number | null | undefined, currency = "BRL"): string | null {
  if (amount === null || amount === undefined || amount === "") return null;
  const value = typeof amount === "number" ? amount : Number(amount);
  if (!Number.isFinite(value)) return null;
  if (currency === "BRL") return brl.format(value);
  return new Intl.NumberFormat("pt-BR", { style: "currency", currency }).format(value);
}

/**
 * formatMoneyOr renders formatMoney's result, falling back to the given
 * placeholder (default the em dash the design uses for unknown values).
 */
export function formatMoneyOr(
  amount: string | number | null | undefined,
  fallback = "—",
  currency = "BRL",
): string {
  return formatMoney(amount, currency) ?? fallback;
}
