import type { ErpImportIssue, ErpImportStatus, ErpImportSummary } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { formatDateTime } from "@marketplace-central/web-query";
import { useState } from "react";
import { Link } from "react-router-dom";
import { useErpImportDetail, useErpImportsList } from "../vinculos/useErpImports";

const statusLabels: Record<ErpImportStatus, string> = {
  COMPLETED: "Concluída",
  REJECTED: "Rejeitada",
};

const statusClasses: Record<ErpImportStatus, string> = {
  COMPLETED: "bg-accent-soft text-accent-ink",
  REJECTED: "bg-warn-soft text-warn",
};

// Duas fontes de arquivo alimentam o espelho e carregam campos diferentes
// (a exportação Sankhya traz custo e estoque; o catálogo do cliente traz marca e
// NCM). Sem a fonte no histórico, dois protocolos de tamanho parecido são
// indistinguíveis — e é justamente a fonte que explica por que um campo está "—".
const sourceLabels: Record<string, string> = {
  xlsx: "Planilha Sankhya",
  catalogo_cliente: "Catálogo do cliente",
};

function SourceBadge({ source }: { source: string }) {
  return (
    <span className="inline-flex whitespace-nowrap rounded-full bg-surface-2 px-2 py-0.5 text-xs font-medium text-muted">
      {sourceLabels[source] ?? source}
    </span>
  );
}

function StatusBadge({ status }: { status: ErpImportStatus }) {
  return (
    <span className={`inline-flex whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ${statusClasses[status]}`}>
      {statusLabels[status]}
    </span>
  );
}

// A real catalog import produces one issue per missing column per row, so a
// single protocol can carry thousands of them — and rendering that list raw
// buries the three distinct problems it actually describes under 4000 identical
// lines. Group by code so the operator reads "what went wrong, on how many
// rows" first, and show a bounded sample of the affected rows underneath. The
// counts are the exact totals, never a truncated one.
const ISSUE_SAMPLE_SIZE = 10;

function groupIssuesByCode(issues: ErpImportIssue[]): { code: string; issues: ErpImportIssue[] }[] {
  const groups = new Map<string, ErpImportIssue[]>();
  for (const issue of issues) {
    const bucket = groups.get(issue.code);
    if (bucket) bucket.push(issue);
    else groups.set(issue.code, [issue]);
  }
  return [...groups.entries()]
    .map(([code, codeIssues]) => ({ code, issues: codeIssues }))
    .sort((a, b) => b.issues.length - a.issues.length);
}

function IssueList({ title, issues, testId }: { title: string; issues: ErpImportIssue[]; testId: string }) {
  if (issues.length === 0) return null;
  const groups = groupIssuesByCode(issues);
  return (
    <div className="mt-2">
      {/* A contagem aqui é de PROBLEMAS, não de linhas: uma linha sem codprod e sem
          descrição rende dois. Dizer "(6)" ao lado de um KPI "Rejeitados 3" lê como
          contradição, então o número vem rotulado. */}
      <p className="text-xs font-semibold text-ink">
        {title}{" "}
        <span className="font-mono font-normal text-faint">
          ({issues.length} {issues.length === 1 ? "problema" : "problemas"})
        </span>
      </p>
      <ul className="mt-1 flex flex-col gap-2" data-testid={testId}>
        {groups.map((group) => {
          const sample = group.issues.slice(0, ISSUE_SAMPLE_SIZE);
          const remaining = group.issues.length - sample.length;
          return (
            <li key={group.code} className="rounded-control border border-border bg-surface-2 px-2 py-1.5 text-xs">
              <p className="text-ink">
                <span className="font-mono font-medium">{group.code}</span>
                <span className="text-muted"> — {group.issues[0].detail}</span>
              </p>
              <p className="mt-0.5 text-faint">
                {group.issues.length === 1 ? "1 linha" : `${group.issues.length} linhas`}:{" "}
                <span className="font-mono">{sample.map((issue) => issue.row).join(", ")}</span>
                {remaining > 0 ? ` e mais ${remaining}` : null}
              </p>
            </li>
          );
        })}
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
          <SourceBadge source={item.source} />
        </div>
        <div className="flex items-center gap-2">
          <button
            type="button"
            className="rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
            aria-expanded={expanded}
            onClick={() => setExpanded((value) => !value)}
          >
            {expanded ? "Ocultar detalhes" : "Ver detalhes"}
          </button>
          <Link
            to={`/importacoes/${item.import_id}`}
            className="rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink"
          >
            Ver cadeia
          </Link>
        </div>
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
          <dd className="font-medium text-ink">{formatDateTime(item.imported_at) ?? "—"}</dd>
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
