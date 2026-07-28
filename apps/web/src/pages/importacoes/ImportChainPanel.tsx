import { ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { formatDateTime } from "@marketplace-central/web-query";
import { useErpImportChain } from "./useErpImportChain";

interface ImportChainPanelProps {
  importId: string;
}

function isImportNotFoundError(error: unknown): boolean {
  if (typeof error !== "object" || error === null) return false;
  const candidate = error as { status?: unknown; error?: unknown };
  return candidate.status === 404 || candidate.error === "import_not_found";
}

function renderCounter(value: unknown, hint: string) {
  return typeof value === "number" && Number.isFinite(value) ? value : <UnknownValue hint={hint} />;
}

export function ImportChainPanel({ importId }: ImportChainPanelProps) {
  const chainQuery = useErpImportChain(importId);
  const errorDetail = isImportNotFoundError(chainQuery.error)
    ? "Importação não encontrada."
    : "Não foi possível carregar a cadeia da importação.";

  return (
    <section className="rounded-card border border-border bg-surface p-4">
      <h2 className="text-sm font-semibold text-ink">Cadeia da importação</h2>
      <p className="mt-1 text-xs text-faint">importados → vinculados → enfileirados, lidos do servidor.</p>
      <div className="mt-3">
        {chainQuery.isPending ? (
          <LoadingState />
        ) : chainQuery.isError || !chainQuery.data ? (
          <div data-testid="erp-import-chain-error">
            <ErrorState detail={errorDetail} onRetry={() => void chainQuery.refetch()} />
          </div>
        ) : (
          <div data-testid="erp-import-chain">
            <p className="text-xs text-muted">
              Protocolo <span className="font-mono font-medium text-ink">{chainQuery.data.protocol}</span>
            </p>
            <dl className="mt-3 grid grid-cols-3 gap-3">
              <div>
                <dt className="text-xs text-faint">Produtos do import</dt>
                <dd data-testid="erp-import-chain-importados" className="mt-0.5 font-mono text-sm font-semibold text-ink">
                  {renderCounter(chainQuery.data.importados, "produtos do import desconhecidos")}
                </dd>
              </div>
              <div>
                <dt className="text-xs text-faint">Vinculados</dt>
                <dd data-testid="erp-import-chain-vinculados" className="mt-0.5 font-mono text-sm font-semibold text-ink">
                  {renderCounter(chainQuery.data.vinculados, "produtos vinculados desconhecidos")}
                </dd>
              </div>
              <div>
                <dt className="text-xs text-faint">Enfileirados</dt>
                <dd data-testid="erp-import-chain-enfileirados" className="mt-0.5 font-mono text-sm font-semibold text-ink">
                  {renderCounter(chainQuery.data.enfileirados, "produtos enfileirados desconhecidos")}
                </dd>
              </div>
            </dl>
            <p className="mt-3 text-xs text-faint">
              Fila lida em: {formatDateTime(chainQuery.data.queue_read_at) ?? <UnknownValue hint="momento da leitura da fila desconhecido" />}
            </p>
          </div>
        )}
      </div>
    </section>
  );
}
