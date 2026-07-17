import { useQuery } from "@tanstack/react-query";
import type { MutationItem, MutationProtocol } from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState, UnknownValue } from "@marketplace-central/ui";
import { failureCopy, mutationsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useClient } from "../../app/ClientContext";
import { mutationTypeLabels, presentMutationValue } from "./mutationPresentation";
import { useMutationProtocol } from "./useMutationProtocol";

function formatTime(value: string | null) {
  if (value === null) return <UnknownValue />;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("pt-BR", { dateStyle: "short", timeStyle: "short" });
}

function Fact({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-slate-100 pb-3 last:border-0 last:pb-0">
      <dt className="text-sm text-slate-500">{label}</dt>
      <dd className="text-right text-sm font-medium text-slate-900">{children}</dd>
    </div>
  );
}

function ProtocolHeader({ protocol }: { protocol: MutationProtocol }) {
  return (
    <header className="rounded-lg border border-slate-200 bg-white p-5">
      <h1 className="text-xl font-semibold text-slate-950">Protocolo {protocol.protocol_id}</h1>
      <dl className="mt-5 grid gap-3 md:grid-cols-2">
        <Fact label="Tipo">{mutationTypeLabels[protocol.type] ?? protocol.type}</Fact>
        <Fact label="Ator">{protocol.actor}</Fact>
        <Fact label="Criado">{formatTime(protocol.created_at)}</Fact>
        <Fact label="Prévia gerada">{formatTime(protocol.previewed_at)}</Fact>
        <Fact label="Aprovado">{formatTime(protocol.approved_at)}</Fact>
        <Fact label="Finalizado">{formatTime(protocol.finished_at)}</Fact>
        <Fact label="Estado do servidor">{protocol.state}</Fact>
        {protocol.retried_from ? <Fact label="Repetido a partir de"><Link className="text-blue-700 hover:underline" to={`/protocolos/${protocol.retried_from}`}>{protocol.retried_from}</Link></Fact> : null}
      </dl>
    </header>
  );
}

function failureCode(item: MutationItem) {
  return typeof item.failure?.code === "string" ? item.failure.code : "internal";
}

function providerMessage(item: MutationItem) {
  return typeof item.failure?.message_provider === "string" ? item.failure.message_provider : null;
}

function ItemsTable({ items }: { items: MutationItem[] }) {
  const [openItems, setOpenItems] = useState<Set<string>>(() => new Set());
  if (items.length === 0) return <EmptyState />;
  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[760px] border-collapse text-left text-sm">
        <caption className="sr-only">Itens do protocolo</caption>
        <thead className="border-b border-slate-200 text-xs font-medium text-slate-500"><tr>
          <th className="px-3 py-3">Anúncio</th><th className="px-3 py-3">Antes</th><th className="px-3 py-3">Depois</th><th className="px-3 py-3">Resultado</th>
        </tr></thead>
        <tbody className="divide-y divide-slate-100">
          {items.map((item) => {
            const message = providerMessage(item);
            const open = openItems.has(item.item_id);
            return <tr key={item.item_id} className="align-top">
              <td className="px-3 py-3 font-medium text-slate-900">{item.listing_id}</td>
              <td className="px-3 py-3 font-mono text-xs">{item.before === null ? <UnknownValue /> : presentMutationValue(item.before)}</td>
              <td className="px-3 py-3 font-mono text-xs">{presentMutationValue(item.after)}</td>
              <td className="px-3 py-3">
                {item.failure ? <p className="text-red-700">{failureCopy(failureCode(item))}</p> : item.state}
                {message ? <div className="mt-2"><button className="text-sm text-slate-700" type="button" aria-expanded={open} onClick={() => setOpenItems((current) => {
                  const next = new Set(current); open ? next.delete(item.item_id) : next.add(item.item_id); return next;
                })}>▸ técnico</button>{open ? <p className="mt-2 break-words font-mono text-xs">{message}</p> : null}</div> : null}
              </td>
            </tr>;
          })}
        </tbody>
      </table>
    </div>
  );
}

function ProtocolItems({ protocolId }: { protocolId: string }) {
  const client = useClient();
  const [cursor, setCursor] = useState<string | undefined>();
  const [history, setHistory] = useState<Array<string | undefined>>([]);
  const query = useQuery({
    queryKey: [...mutationsQueryKeys.items(protocolId), { cursor }],
    queryFn: () => client.listMutationItems(protocolId, { limit: 50, cursor }),
    staleTime: QUERY_STALE_TIME.mutations,
  });
  if (query.isPending) return <LoadingState />;
  if (query.isError) return <ErrorState detail="Não foi possível carregar os itens." onRetry={() => void query.refetch()} />;
  return <section className="rounded-lg border border-slate-200 bg-white p-5" aria-labelledby="protocol-items-title">
    <h2 id="protocol-items-title" className="mb-4 text-base font-semibold text-slate-950">Itens</h2>
    <ItemsTable items={query.data.items} />
    <nav className="mt-4 flex justify-end gap-2" aria-label="Paginação dos itens">
      <button type="button" disabled={history.length === 0} onClick={() => { const next = [...history]; setCursor(next.pop()); setHistory(next); }}>Página anterior</button>
      <button type="button" disabled={!query.data.next_cursor} onClick={() => { setHistory((value) => [...value, cursor]); setCursor(query.data.next_cursor ?? undefined); }}>Próxima página</button>
    </nav>
  </section>;
}

export function ProtocoloPage() {
  const { protocolId = "" } = useParams();
  const protocol = useMutationProtocol(protocolId);
  if (protocol.isPending) return <LoadingState />;
  if (protocol.isError) return <ErrorState detail="Não foi possível carregar o protocolo." onRetry={() => void protocol.refetch()} />;
  if (!protocol.data) return null;
  return <main className="space-y-5 p-6"><ProtocolHeader protocol={protocol.data} /><ProtocolItems protocolId={protocolId} /></main>;
}
