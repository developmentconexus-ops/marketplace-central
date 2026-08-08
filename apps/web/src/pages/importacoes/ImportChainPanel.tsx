import { hasCode, isApiError } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { formatDateTime } from "@marketplace-central/web-query";
import { useErpImportChain } from "./useErpImportChain";

interface ImportChainPanelProps {
  importId: string;
}

function isImportNotFoundError(error: unknown): boolean {
  return isApiError(error) && hasCode(error, "import_not_found");
}

function renderCounter(value: unknown, hint: string) {
  return typeof value === "number" && Number.isFinite(value) ? value : <UnknownValue hint={hint} />;
}

function renderProtocol(value: unknown) {
  return typeof value === "string" && value.trim().length > 0 ? (
    <span className="font-mono font-medium text-ink">{value}</span>
  ) : (
    <UnknownValue hint="protocolo desconhecido" />
  );
}

function renderQueueReadAt(value: unknown) {
  if (typeof value !== "string" || value.trim().length === 0) {
    return <UnknownValue hint="momento da leitura da fila desconhecido" />;
  }

  return formatDateTime(value) ?? <UnknownValue hint="momento da leitura da fila desconhecido" />;
}

export function ImportChainPanel({ importId }: ImportChainPanelProps) {
  const chainQuery = useErpImportChain(importId);
  const errorDetail = isImportNotFoundError(chainQuery.error)
    ? "Importação não encontrada."
    : "Não foi possível carregar o estado da importação.";

  return (
    <section className="rounded-card border border-border bg-surface p-4">
      <h2 className="text-sm font-semibold text-ink">Estado da importação</h2>
      {/* Sem seta entre os três: são duas unidades (linhas do arquivo / produtos internos) e
          nenhum é etapa do outro, então 55 · 0 · 55 é estado normal. */}
      <p className="mt-1 text-xs text-faint">
        Três medidas da mesma importação, em duas unidades — nenhuma é etapa da outra.
      </p>
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
              Protocolo {renderProtocol(chainQuery.data.protocol)}
            </p>
            <dl className="mt-3 grid grid-cols-3 gap-3">
              <div>
                <dt className="text-xs text-faint">Linhas importadas</dt>
                <dd
                  data-testid="erp-import-chain-importados"
                  className="mt-0.5 font-mono text-sm font-semibold text-ink"
                >
                  {renderCounter(chainQuery.data.importados, "linhas importadas desconhecidas")}
                </dd>
              </div>
              <div>
                <dt className="text-xs text-faint">Produtos vinculados</dt>
                <dd
                  data-testid="erp-import-chain-vinculados"
                  className="mt-0.5 font-mono text-sm font-semibold text-ink"
                >
                  {renderCounter(chainQuery.data.vinculados, "produtos vinculados desconhecidos")}
                </dd>
              </div>
              <div>
                <dt className="text-xs text-faint">Linhas na fila de sync</dt>
                <dd
                  data-testid="erp-import-chain-enfileirados"
                  className="mt-0.5 font-mono text-sm font-semibold text-ink"
                >
                  {renderCounter(chainQuery.data.enfileirados, "linhas na fila desconhecidas")}
                </dd>
              </div>
            </dl>
            <p className="mt-3 text-xs text-faint">
              Fila lida em: {renderQueueReadAt(chainQuery.data.queue_read_at)}
            </p>
          </div>
        )}
      </div>
    </section>
  );
}
