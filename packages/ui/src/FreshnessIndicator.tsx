import { formatAsOf, formatDateTime } from "@marketplace-central/web-query";

/**
 * Rótulo de frescor de um fato. Único componente de idade do produto — a cópia
 * que vivia em @marketplace-central/web-query foi removida (D-49), porque duas
 * cópias com aria-label diferentes prendiam cada teste a qual delas o arquivo
 * tinha importado.
 *
 * Mostra idade relativa ("há 3 h"), que é o que se lê de relance, e guarda o
 * instante absoluto no title para o operador cruzar com um log.
 */
export function FreshnessIndicator({ asOf }: { asOf: string | null | undefined }) {
  const absolute = formatDateTime(asOf);
  return (
    <span
      className="text-muted text-xs font-mono"
      aria-label="Atualização dos dados"
      {...(absolute === null ? {} : { title: absolute })}
    >
      {formatAsOf(asOf)}
    </span>
  );
}
