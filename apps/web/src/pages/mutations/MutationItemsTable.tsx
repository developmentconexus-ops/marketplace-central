import type { MutationItem } from "@marketplace-central/sdk-runtime";
import { EmptyState, UnknownValue } from "@marketplace-central/ui";
import { presentMutationValue } from "./mutationPresentation";

export function MutationItemsTable({ items }: { items: MutationItem[] }) {
  if (items.length === 0) return <EmptyState />;

  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[620px] border-collapse text-left text-sm">
        <caption className="sr-only">Alterações previstas por anúncio</caption>
        <thead className="border-b border-slate-200 text-xs font-medium text-slate-500">
          <tr>
            <th className="px-3 py-3" scope="col">Anúncio</th>
            <th className="px-3 py-3" scope="col">Antes</th>
            <th className="px-3 py-3" scope="col">Depois</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {items.map((item) => (
            <tr key={item.item_id} className="align-top text-slate-700">
              <td className="px-3 py-3 font-medium text-slate-900">{item.listing_id}</td>
              <td className="px-3 py-3 font-mono text-xs">
                {item.before === null ? <UnknownValue hint="valor anterior não informado" /> : presentMutationValue(item.before)}
              </td>
              <td className="px-3 py-3 font-mono text-xs">{presentMutationValue(item.after)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
