import { formatAsOf } from "@marketplace-central/web-query";

export function FreshnessIndicator({ asOf }: { asOf: string | null | undefined }) {
  return <span aria-label="Atualização dos dados">{formatAsOf(asOf)}</span>;
}
