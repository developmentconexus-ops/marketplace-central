import { useQuery } from "@tanstack/react-query";
import type { MutationProtocol } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import { failureCopy, mutationsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";

export function MutationResultSummary({ protocol }: { protocol: MutationProtocol }) {
  const client = useClient();
  const items = useQuery({
    queryKey: mutationsQueryKeys.items(protocol.protocol_id),
    queryFn: () => client.listMutationItems(protocol.protocol_id),
    staleTime: QUERY_STALE_TIME.mutations,
  });

  return (
    <div className="space-y-4">
      <div>
        <h3 className="text-base font-semibold text-slate-950">Aplicação concluída</h3>
        <p className="mt-1 text-sm text-slate-600">Estado do protocolo: {protocol.state}</p>
      </div>
      {items.isLoading ? <LoadingState /> : null}
      {items.isError ? (
        <ErrorState onRetry={() => void items.refetch()} detail="Não foi possível carregar os resultados." />
      ) : null}
      {items.data?.items.length === 0 ? <EmptyState /> : null}
      {items.data?.items.length ? (
        <ul className="divide-y divide-slate-100 rounded-lg border border-slate-200">
          {items.data.items.map((item) => (
            <li key={item.item_id} className="px-4 py-3 text-sm">
              <span className="font-medium text-slate-900">{item.listing_id}</span>
              {item.failure ? <p className="mt-1 text-red-700">{failureCopy(failureCode(item.failure))}</p> : null}
            </li>
          ))}
        </ul>
      ) : null}
      <a className="inline-flex text-sm font-medium text-blue-700 hover:underline" href={`/protocolos/${protocol.protocol_id}`}>
        Ver protocolo
      </a>
    </div>
  );
}

function failureCode(failure: Record<string, unknown>): string {
  return typeof failure.code === "string" ? failure.code : "internal";
}
