import { useState } from "react";
import type { JSX } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import type { PricingScenario } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";

export interface ScenariosPanelProps {
  /** The current simulation state to persist when a scenario is saved. */
  payload: Record<string, unknown>;
  /** Called with a saved scenario's payload when the operator reloads it. */
  onReload: (payload: Record<string, unknown>) => void;
}

/** Derive a stable id from the scenario name so re-saving a name overwrites it. */
function scenarioId(name: string): string {
  return name
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

/**
 * ScenariosPanel persists, reloads and deletes named pricing scenarios via the
 * IC-04 scenarios CRUD. Saving snapshots the current simulation `payload`;
 * reloading hands that payload back to the page (which re-applies it); deleting
 * removes it. The scenario id is derived from the name so re-saving overwrites
 * rather than duplicating.
 */
export function ScenariosPanel({ payload, onReload }: ScenariosPanelProps): JSX.Element {
  const client = useClient();
  const queryClient = useQueryClient();
  const [name, setName] = useState<string>("");

  const scenariosQuery = useQuery({
    queryKey: ["pricing", "scenarios"],
    queryFn: () => client.listPricingScenarios(),
  });

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: ["pricing", "scenarios"] });
  };

  const save = useMutation({
    mutationFn: () =>
      client.createPricingScenario({ id: scenarioId(name), name: name.trim(), payload }),
    onSuccess: () => {
      setName("");
      invalidate();
    },
  });

  const remove = useMutation({
    mutationFn: (id: string) => client.deletePricingScenario(id),
    onSuccess: () => invalidate(),
  });

  const items = scenariosQuery.data?.items ?? [];
  const saveDisabled = name.trim() === "" || save.isPending;

  return (
    <div className="flex flex-col gap-3">
      <h3 className="text-sm font-semibold text-ink">Cenários</h3>

      <div className="flex flex-wrap items-end gap-2">
        <label className="flex flex-col text-xs text-muted">
          Nome do cenário
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="ex.: Clássico agressivo"
            aria-label="Nome do cenário"
            className="mt-1 w-52 rounded-md border border-border bg-surface px-2 py-1 text-sm text-ink"
          />
        </label>
        <button
          type="button"
          data-testid="scenario-save"
          disabled={saveDisabled}
          onClick={() => save.mutate()}
          className="rounded-md border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2 disabled:opacity-40"
        >
          {save.isPending ? "Salvando…" : "Salvar cenário"}
        </button>
      </div>

      {scenariosQuery.isLoading ? (
        <LoadingState />
      ) : scenariosQuery.isError ? (
        <ErrorState
          onRetry={() => void scenariosQuery.refetch()}
          detail="Falha ao carregar os cenários."
        />
      ) : items.length === 0 ? (
        <p data-testid="scenarios-empty" className="text-xs text-muted">
          Nenhum cenário salvo ainda — simule e salve para comparar depois.
        </p>
      ) : (
        <ul className="flex flex-col gap-1">
          {items.map((s: PricingScenario) => (
            <li
              key={s.id}
              data-testid={`scenario-row-${s.id}`}
              className="flex items-center justify-between rounded-md border border-border bg-surface-2 px-2 py-1.5"
            >
              <span className="truncate text-sm text-ink">{s.name}</span>
              <span className="flex shrink-0 gap-2">
                <button
                  type="button"
                  data-testid={`scenario-reload-${s.id}`}
                  onClick={() => onReload(s.payload)}
                  className="rounded-md px-2 py-1 text-xs text-accent hover:bg-surface"
                >
                  Recarregar
                </button>
                <button
                  type="button"
                  data-testid={`scenario-delete-${s.id}`}
                  disabled={remove.isPending}
                  onClick={() => remove.mutate(s.id)}
                  className="rounded-md px-2 py-1 text-xs text-warn hover:bg-warn-soft disabled:opacity-40"
                >
                  Excluir
                </button>
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
