import type { ProductLinkWorkflowItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { formatDateTime } from "../pedidos/pedidosFormatters";
import { isAutoResolved, resolutionAuditId, useVinculosResolved } from "./useVinculosResolved";

export interface ResolvidosTabProps {
  installationId: string;
}

function ResolvidoRow({
  item,
  onUndo,
  pending,
}: {
  item: ProductLinkWorkflowItem;
  onUndo: (auditId: string) => void;
  pending: boolean;
}) {
  const link = item.current_link;
  const auditId = resolutionAuditId(item);
  const autoResolved = isAutoResolved(item);

  return (
    <tr className="align-top text-ink" data-testid="resolvido-row">
      {/* ANÚNCIO ML */}
      <td className="px-3 py-3">
        <div className="font-mono text-sm font-medium text-ink">{item.identity.provider_item_id}</div>
      </td>

      {/* PRODUTO vinculado */}
      <td className="px-3 py-3">
        <div className="font-medium text-ink">
          {link?.internal_product_name ? link.internal_product_name : <UnknownValue hint="sem descrição no ERP" />}
        </div>
      </td>

      {/* SKU HUB */}
      <td className="px-3 py-3">
        <span className="font-mono text-sm text-ink">
          {link?.internal_product_id === undefined ? <UnknownValue hint="sem CODPROD" /> : link.internal_product_id}
        </span>
      </td>

      {/* Estado (+ badge de auto-aprovado, F-04) */}
      <td className="px-3 py-3">
        <div className="flex flex-wrap items-center gap-1">
          <span className="inline-flex whitespace-nowrap rounded-full bg-accent-soft px-2 py-0.5 text-xs font-medium text-accent-ink">
            Vinculado ✓
          </span>
          {/* Ausência do badge significa "não foi automático" — verdade para
              todo vínculo manual e para todo registro anterior ao M-05, que não
              tem entrada de auditoria resolvedora. Nunca um badge fabricado. */}
          {autoResolved ? (
            <span
              data-testid="auto-resolved-badge"
              title="Vínculo resolvido automaticamente pelo sistema"
              className="inline-flex whitespace-nowrap rounded-full border border-border bg-surface-2 px-2 py-0.5 text-xs font-medium text-muted"
            >
              Automático
            </span>
          ) : null}
        </div>
      </td>

      {/* Resolvido em */}
      <td className="px-3 py-3">
        <span className="text-xs text-muted">
          {/* The wire carries RFC-3339 UTC; the operator reads pt-BR local time.
              An unparseable/absent timestamp stays honest (—), never fabricated. */}
          {formatDateTime(link?.updated_at) ?? <UnknownValue />}
        </span>
      </td>

      {/* Desfazer */}
      <td className="px-3 py-3 text-right">
        <button
          type="button"
          disabled={pending || auditId === undefined}
          title={auditId === undefined ? "Sem registro de auditoria para desfazer" : undefined}
          className="rounded-control border border-border px-2.5 py-1.5 text-xs font-medium text-muted hover:bg-surface-2 hover:text-ink disabled:cursor-not-allowed disabled:opacity-50"
          onClick={() => auditId && onUndo(auditId)}
        >
          Desfazer
        </button>
      </td>
    </tr>
  );
}

export function ResolvidosTab({ installationId }: ResolvidosTabProps) {
  const { resolvedQuery, items, undo } = useVinculosResolved(installationId);

  if (resolvedQuery.isPending) return <LoadingState />;
  if (resolvedQuery.isError) return <ErrorState onRetry={() => void resolvedQuery.refetch()} />;
  if (items.length === 0) return <EmptyState />;

  return (
    <div className="mt-3 overflow-x-auto">
      <table className="w-full min-w-[720px] border-collapse text-left text-sm">
        <caption className="sr-only">Vínculos resolvidos</caption>
        <thead className="border-b border-border bg-surface-2 text-xs font-medium tracking-[0.04em] text-faint">
          <tr>
            {/* Rótulo estrutural neutro de provider (F-05). */}
            <th className="px-3 py-3" scope="col">Anúncio</th>
            <th className="px-3 py-3" scope="col">Produto</th>
            <th className="px-3 py-3" scope="col">SKU HUB</th>
            <th className="px-3 py-3" scope="col">Estado</th>
            <th className="px-3 py-3" scope="col">Resolvido em</th>
            <th className="px-3 py-3 text-right" scope="col">Ação</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border-2">
          {items.map((item) => (
            <ResolvidoRow
              key={`${item.identity.provider_item_id}-${item.identity.provider_variation_id ?? ""}`}
              item={item}
              onUndo={(auditId) => undo.mutate(auditId)}
              pending={undo.isPending}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
}
