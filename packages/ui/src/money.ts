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
export function formatMoney(
  amount: string | number | null | undefined,
  currency = "BRL",
): string | null {
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
/**
 * formatPercent renders a rate in pt-BR ("12.80" → "12,80%"). Rates arrive as
 * decimal STRINGS whose precision is itself information — a comissão quoted as
 * "13.00" is not the same statement as "13" — so a string keeps exactly the
 * decimal places the API sent and only the separator changes. A numeric rate has
 * no such source precision, so it is rounded to at most `fractionDigits` and
 * trailing zeros are dropped (25 → "25%", -65.49 → "-65,49%"). Absent or
 * unparseable returns null: the caller renders "—", never a fabricated "0%".
 */
export function formatPercent(
  rate: string | number | null | undefined,
  fractionDigits = 2,
): string | null {
  if (rate === null || rate === undefined || rate === "") return null;
  const value = typeof rate === "number" ? rate : Number(rate);
  if (!Number.isFinite(value)) return null;
  if (typeof rate === "string") return `${rate.trim().replace(".", ",")}%`;
  const digits = new Intl.NumberFormat("pt-BR", {
    minimumFractionDigits: 0,
    maximumFractionDigits: fractionDigits,
  }).format(value);
  return `${digits}%`;
}

/**
 * formatSignedPercent renders a delta-vs-market rate in pt-BR with an explicit
 * "+" on non-negative values ("-3.2" → "-3,2%", "3.2" → "+3,2%"). Null for an
 * absent/unparseable rate, same honest-unknown contract as formatPercent.
 */
export function formatSignedPercent(
  rate: string | number | null | undefined,
  fractionDigits = 1,
): string | null {
  const formatted = formatPercent(rate, fractionDigits);
  if (formatted === null) return null;
  return formatted.startsWith("-") ? formatted : `+${formatted}`;
}

export function formatMoneyOr(
  amount: string | number | null | undefined,
  fallback = "—",
  currency = "BRL",
): string {
  return formatMoney(amount, currency) ?? fallback;
}
