import { formatAsOf } from "@marketplace-central/web-query";

export function FreshnessIndicator({ asOf }: { asOf: string | null | undefined }) {
  return (
    <span className="text-muted text-xs font-mono" aria-label="Atualização dos dados">
      {formatAsOf(asOf)}
    </span>
  );
}
