import type { ErpImportIssue, ErpImportStatus, ErpImportSummary } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { useState } from "react";
import { useErpImportDetail, useErpImportsList } from "./useErpImports";

const statusLabels: Record<ErpImportStatus, string> = {
  COMPLETED: "Concluída",
  REJECTED: "Rejeitada",
};

const statusClasses: Record<ErpImportStatus, string> = {
  COMPLETED: "bg-accent-soft text-accent-ink",
  REJECTED: "bg-warn-soft text-warn",
};

function StatusBadge({ status }: { status: ErpImportStatus }) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ${statusClasses[status]}`}>
      {statusLabels[status]}
    </span>
  );
}

function IssueList({ title, issues, testId }: { title: string; issues: ErpImportIssue[]; testId: string }) {
  if (issues.length === 0) return null;
  return (
    <div className="mt-2">
      <p className="text-xs font-semibold text-ink">{title}</p>
      <ul className="mt-1 flex flex-col gap-1" data-testid={testId}>
        {issues.map((issue, index) => (
          <li
            key={`${issue.row}-${issue.code}-${index}`}
            className="rounded-control border border-border bg-surface-2 px-2 py-1 text-xs text-muted"
          >
            <span className="font-medium">Linha {issue.row}</span> — {issue.code}: {issue.detail}
            {issue.offending_value ? (
              <span className="text-faint"> (valor: {issue.offending_value})</span>
            ) : null}
          </li>
        ))}
      </ul>
    </div>
  );
}

function ImportRowDetail({ importId }: { importId: string }) {
  const detailQuery = useErpImportDetail(importId);

  if (detailQuery.isPending) return <LoadingState />;
  if (detailQuery.isError) {
    return <ErrorState onRetry={() => void detailQuery.refetch()} />;
  }

  const detail = detailQuery.data;
  if (!detail) return <EmptyState />;

  const hasIssues = detail.rejected_rows.length > 0 || detail.warnings.length > 0;

  return (
    <div className="mt-3 border-t border-border pt-3" data-testid="erp-import-detail">
      {hasIssues ? (
        <>
          <IssueList title="Rejeitados" issues={detail.rejected_rows} testId="erp-import-rejected-rows" />
          <IssueList title="Avisos" issues={detail.warnings} testId="erp-import-warnings" />
        </>
      ) : (
        <p className="text-xs text-faint">Sem linhas rejeitadas ou avisos nesta importação.</p>
      )}
    </div>
  );
}

function ImportRow({ item }: { item: ErpImportSummary }) {
  const [expanded, setExpanded] = useState(false);

  return (
    <li className="rounded-card border border-border bg-surface p-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <span className="font-mono font-medium text-ink">{item.protocol}</span>
          <StatusBadge status={item.status} />
        </div>
        <button
          type="button"
          className="rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
          aria-expanded={expanded}
          onClick={() => setExpanded((value) => !value)}
        >
          {expanded ? "Ocultar detalhes" : "Ver detalhes"}
        </button>
      </div>
      <dl className="mt-2 grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-muted sm:grid-cols-4">
        <div>
          <dt className="text-faint">Aceitos</dt>
          <dd className="font-mono font-medium text-ink">{item.accepted_count}</dd>
        </div>
        <div>
          <dt className="text-faint">Rejeitados</dt>
          <dd className="font-mono font-medium text-ink">{item.rejected_count}</dd>
        </div>
        <div>
          <dt className="text-faint">Avisos</dt>
          <dd className="font-mono font-medium text-ink">{item.warning_count}</dd>
        </div>
        <div>
          <dt className="text-faint">Importado em</dt>
          <dd className="font-medium text-ink">{item.imported_at}</dd>
        </div>
      </dl>
      {expanded ? <ImportRowDetail importId={item.import_id} /> : null}
    </li>
  );
}

export function ImportacaoSection() {
  const importsQuery = useErpImportsList();

  const items = importsQuery.data?.items ?? [];

  return (
    <section aria-labelledby="importacao-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="importacao-title" className="text-sm font-semibold text-ink">
        Importação
      </h2>
      <p className="mt-1 text-xs text-faint">
        Histórico de importações do ERP (protocolo). Somente leitura — sem upload nesta seção.
      </p>

      <div className="mt-3">
        {importsQuery.isPending ? (
          <LoadingState />
        ) : importsQuery.isError ? (
          <ErrorState onRetry={() => void importsQuery.refetch()} />
        ) : items.length === 0 ? (
          <EmptyState />
        ) : (
          <ul className="flex flex-col gap-2">
            {items.map((item) => (
              <ImportRow key={item.import_id} item={item} />
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}
