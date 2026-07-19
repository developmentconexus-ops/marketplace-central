/**
 * pt-BR decimal → dot-decimal normalization for the /precos simulator.
 *
 * The Go backend's decimal parser (`domain.ParseRat`) expects dot-decimal
 * strings (`"150.00"`). The operator, typing on a pt-BR keyboard/locale, enters
 * comma-decimal (`"150,00"`, and for money possibly `"1.500,90"` with '.'
 * thousands separators). Sending the raw operator string straight to the SDK
 * yields a 422 at the parser. These helpers bridge that gap AT THE SDK SEND
 * POINT — inputs keep the raw pt-BR string the operator sees; only the value
 * bound to a request is normalized.
 *
 * The design prototype (docs/design/handoff-2026-07/Simulador.dc.html) strips
 * dots unconditionally because its inputs are always seeded comma-form. Our
 * fields are not: a `fetched`-from-backend `current_price` and a reloaded
 * scenario carry dot-decimal (`"89.00"`, `"77.00"`). So money normalization
 * strips thousands dots ONLY when a comma is present (comma ⇒ dots are
 * thousands); a no-comma string is already dot-decimal and passes through
 * untouched. Rate/percent fields never carry thousands separators, so they only
 * promote a decimal comma.
 */

/**
 * Normalize a pt-BR money string to dot-decimal for the SDK.
 *
 * `"150,00"` → `"150.00"`, `"1.500,90"` → `"1500.90"`, `"77.00"` → `"77.00"`
 * (already dot-decimal), `"1500"` → `"1500"`, `""` → `""`.
 */
export function ptBrMoneyToDot(raw: string): string {
  const s = raw.trim();
  if (s === "") return "";
  if (s.includes(",")) {
    // Canonical pt-BR money: dots are thousands, a single comma is the decimal —
    // so the comma is the LAST separator and appears once. If a dot follows the
    // (last) comma, or there is more than one comma, the string is not pt-BR
    // (e.g. US "1,234.56", or malformed). ADR-17: never fabricate a plausible
    // wrong number — refuse to guess and pass the raw string through so the Go
    // parser rejects it honestly (422) instead of silently mis-scaling it.
    const lastComma = s.lastIndexOf(",");
    const dotAfterComma = s.indexOf(".", lastComma) !== -1;
    const multipleCommas = s.indexOf(",") !== lastComma;
    if (dotAfterComma || multipleCommas) return s;
    // Comma present, unambiguously pt-BR ⇒ strip thousands dots, promote comma.
    return s.replace(/\./g, "").replace(",", ".");
  }
  // No comma ⇒ already dot-decimal (or integer); leave dots alone.
  return s;
}

/**
 * Normalize a pt-BR rate/percent string to dot-decimal for the SDK.
 *
 * `"9,25"` → `"9.25"`, `"18"` → `"18"`, `"20.5"` → `"20.5"` (already
 * dot-decimal), `""` → `""`. Percent fields carry no thousands separators, so a
 * lone comma is the decimal point.
 */
export function ptBrRateToDot(raw: string): string {
  const s = raw.trim();
  if (s === "") return "";
  return s.replace(",", ".");
}
