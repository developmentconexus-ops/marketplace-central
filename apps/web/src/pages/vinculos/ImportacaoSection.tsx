import type { ErpImportIssue, ErpImportStatus, ErpImportSummary } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { useState } from "react";
import { useErpImportDetail, useErpImportsList } from "./useErpImports";

const statusLabels: Record<ErpImportStatus, string> = {
  COMPLETED: "Concluída",
  REJECTED: "Rejeitada",
};

const statusClasses: Record<ErpImportStatus, string> = {
  COMPLETED: "bg-emerald-100 text-emerald-800",
  REJECTED: "bg-red-100 text-red-800",
};

function StatusBadge({ status }: { status: ErpImportStatus }) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded px-2 py-0.5 text-xs font-medium ${statusClasses[status]}`}>
      {statusLabels[status]}
    </span>
  );
}

function IssueList({ title, issues, testId }: { title: string; issues: ErpImportIssue[]; testId: string }) {
  if (issues.length === 0) return null;
  return (
    <div className="mt-2">
      <p className="text-xs font-semibold text-slate-700">{title}</p>
      <ul className="mt-1 flex flex-col gap-1" data-testid={testId}>
        {issues.map((issue, index) => (
          <li
            key={`${issue.row}-${issue.code}-${index}`}
            className="rounded border border-slate-200 bg-slate-50 px-2 py-1 text-xs text-slate-700"
          >
            <span className="font-medium">Linha {issue.row}</span> — {issue.code}: {issue.detail}
            {issue.offending_value ? (
              <span className="text-slate-500"> (valor: {issue.offending_value})</span>
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
    <div className="mt-3 border-t border-slate-200 pt-3" data-testid="erp-import-detail">
      {hasIssues ? (
        <>
          <IssueList title="Rejeitados" issues={detail.rejected_rows} testId="erp-import-rejected-rows" />
          <IssueList title="Avisos" issues={detail.warnings} testId="erp-import-warnings" />
        </>
      ) : (
        <p className="text-xs text-slate-500">Sem linhas rejeitadas ou avisos nesta importação.</p>
      )}
    </div>
  );
}

function ImportRow({ item }: { item: ErpImportSummary }) {
  const [expanded, setExpanded] = useState(false);

  return (
    <li className="rounded-lg border border-slate-200 bg-white p-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <span className="font-medium text-slate-900">{item.protocol}</span>
          <StatusBadge status={item.status} />
        </div>
        <button
          type="button"
          className="rounded-lg border border-slate-200 px-2.5 py-1.5 text-xs font-medium text-slate-700 hover:border-slate-300 hover:bg-slate-50"
          aria-expanded={expanded}
          onClick={() => setExpanded((value) => !value)}
        >
          {expanded ? "Ocultar detalhes" : "Ver detalhes"}
        </button>
      </div>
      <dl className="mt-2 grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-slate-600 sm:grid-cols-4">
        <div>
          <dt className="text-slate-500">Aceitos</dt>
          <dd className="font-medium text-slate-900">{item.accepted_count}</dd>
        </div>
        <div>
          <dt className="text-slate-500">Rejeitados</dt>
          <dd className="font-medium text-slate-900">{item.rejected_count}</dd>
        </div>
        <div>
          <dt className="text-slate-500">Avisos</dt>
          <dd className="font-medium text-slate-900">{item.warning_count}</dd>
        </div>
        <div>
          <dt className="text-slate-500">Importado em</dt>
          <dd className="font-medium text-slate-900">{item.imported_at}</dd>
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
    <section aria-labelledby="importacao-title" className="rounded-xl border border-slate-200 bg-white p-4">
      <h2 id="importacao-title" className="text-sm font-semibold text-slate-900">
        Importação
      </h2>
      <p className="mt-1 text-xs text-slate-500">
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
