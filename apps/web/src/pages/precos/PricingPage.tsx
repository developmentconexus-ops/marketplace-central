import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import type {
  CatalogProductFact,
  PricingCalcInput,
  PricingCalcProfile,
} from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { DecompositionPanel } from "./DecompositionPanel";
import { SolverPanel } from "./SolverPanel";
import { ParamsDrawer } from "./ParamsDrawer";
import { DifalDrawer } from "./DifalDrawer";
import { MarketComparison } from "./MarketComparison";
import { ApplyPriceAction } from "./ApplyPriceAction";
import { ScenariosPanel } from "./ScenariosPanel";
import { ptBrMoneyToDot } from "./ptbrDecimal";

/** Seed padrão thresholds/regime used until the operator's profile loads (or if it errors). */
const DEFAULT_PROFILE: PricingCalcProfile = {
  regime: "SIMPLES",
  aliquota_pct: "4",
  limiar_verde_pct: "18",
  limiar_amarelo_pct: "10",
  tarifa_full: null,
  difal_enabled: false,
  difal_destino_uf: null,
  origem: "seed",
};

/**
 * Sale modalidades. The ML commission is NOT hardcoded here: an absent
 * `comissao_pct` lets the backend resolver chain run (degrau 3 live COTACAO
 * quote → degrau 4 tenant PADRAO). Sending a fixed pct would force MANUAL
 * override (degrau 0) and hide the real tariff, so we only send `modalidade`.
 */
export const MODALIDADES = [
  { key: "classico", label: "Clássico" },
  { key: "premium", label: "Premium" },
  { key: "full", label: "Full" },
] as const;

export type ModalidadeKey = (typeof MODALIDADES)[number]["key"];

function productLabel(p: CatalogProductFact): string {
  return p.description ?? p.manufacturer_reference ?? `#${p.internal_product_id}`;
}

/**
 * PricingPage is the /precos IC-04 simulator shell. It loads the calc profile
 * and the catalog products, owns the selected product + price/modalidade state,
 * and lays out the regions the later slices flesh out: the per-component
 * decomposition, the parâmetros/DIFAL drawers, the market comparison, and the
 * "aplicar preço" action. Editing the profile (or a DIFAL override) invalidates
 * the decomposition so the breakdown recomputes against the new server-side
 * profile. All money/percentages stay decimal strings end-to-end (never float).
 */
export function PricingPage() {
  const client = useClient();
  const queryClient = useQueryClient();
  const [searchParams] = useSearchParams();

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [modalidade, setModalidade] = useState<ModalidadeKey>("premium");
  const [precoInput, setPrecoInput] = useState<string>("");
  // Deep link: /precos?params=1 opens the Parâmetros drawer on mount.
  const [paramsOpen, setParamsOpen] = useState<boolean>(searchParams.get("params") === "1");
  const [difalOpen, setDifalOpen] = useState<boolean>(false);

  const profileQuery = useQuery({
    queryKey: ["pricing", "profile"],
    queryFn: () => client.getPricingProfile(),
  });
  const profile = profileQuery.data ?? DEFAULT_PROFILE;

  const productsQuery = useQuery({
    queryKey: ["pricing", "catalog-facts"],
    queryFn: () => client.listCatalogProductFacts({ limit: 50 }),
  });

  const difalQuery = useQuery({
    queryKey: ["pricing", "difal"],
    queryFn: () => client.listPricingDifal(),
    enabled: difalOpen,
  });

  // Profile edits + DIFAL overrides both re-derive the server-side decomposition.
  const invalidateDecompose = () => {
    void queryClient.invalidateQueries({ queryKey: ["pricing", "decompose"] });
  };

  const saveProfile = useMutation({
    mutationFn: (next: PricingCalcProfile) => client.putPricingProfile(next),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["pricing", "profile"] });
      invalidateDecompose();
      setParamsOpen(false);
    },
  });

  const overrideDifal = useMutation({
    mutationFn: (vars: { uf: string; interna: string }) =>
      client.putPricingDifalOverride(vars.uf, { interna_pct: vars.interna, actor: "operator" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["pricing", "difal"] });
      invalidateDecompose();
    },
  });

  const products = productsQuery.data?.items ?? [];
  const selected = useMemo<CatalogProductFact | null>(() => {
    if (products.length === 0) return null;
    // No explicit pick yet ⇒ default to the first product. An explicit id that
    // isn't in the loaded list ⇒ null (honest empty), NEVER a silent fallback to
    // products[0] — that would show a different product's decomposition.
    if (selectedId === null) return products[0];
    return products.find((p) => p.internal_product_id === selectedId) ?? null;
  }, [products, selectedId]);

  // "Aplicar preço" targets an ML listing under an installation — resolve both
  // (installation app-wide, listing by the selected product) for the action.
  const installationsQuery = useQuery({
    queryKey: ["pricing", "installations"],
    queryFn: () => client.listIntegrationInstallations(),
  });
  const installationId = installationsQuery.data?.items[0]?.installation_id ?? "";

  const listingQuery = useQuery({
    queryKey: ["pricing", "listing", installationId, selected?.internal_product_id ?? null],
    queryFn: () =>
      client.listListingsByProduct({
        installation_id: installationId,
        product_id: String(selected!.internal_product_id),
        limit: 1,
      }),
    enabled: installationId !== "" && selected !== null,
  });
  const listingId = listingQuery.data?.groups[0]?.listings[0]?.listing_id ?? null;

  // Working price: explicit input wins, else the product's current price. `preco`
  // stays RAW (pt-BR, what the operator sees / scenarios round-trip); `precoForApi`
  // is the dot-decimal the Go parser expects, normalized only at the SDK send point.
  const preco = precoInput.trim() !== "" ? precoInput.trim() : selected?.current_price.amount ?? "";
  const precoForApi = ptBrMoneyToDot(preco);

  // A scenario snapshots the working simulation; reloading re-applies it to the
  // page state (product, modalidade, price) so the operator can compare setups.
  const scenarioPayload: Record<string, unknown> = {
    preco,
    modalidade,
    product_id: selected?.internal_product_id ?? null,
  };
  const applyScenario = (payload: Record<string, unknown>) => {
    if (typeof payload.product_id === "number") setSelectedId(payload.product_id);
    if (typeof payload.modalidade === "string" && MODALIDADES.some((m) => m.key === payload.modalidade)) {
      setModalidade(payload.modalidade as ModalidadeKey);
    }
    if (typeof payload.preco === "string") setPrecoInput(payload.preco);
  };

  // An explicit product pick that the loaded catalog can't resolve (e.g. a
  // scenario reloaded for a product outside the fetched page) ⇒ tell the
  // operator instead of silently decomposing a different product at this price.
  // Derived (not stored) so it can't go stale against a still-loading catalog.
  const productMissing = selectedId !== null && products.length > 0 && selected === null;

  // comissao_pct is deliberately OMITTED so the backend resolver chain runs
  // (live COTACAO degrau 3 → PADRAO degrau 4). Sending it = MANUAL override.
  const decomposeInput: PricingCalcInput | null = selected && preco !== ""
    ? {
        preco: precoForApi,
        modalidade,
        product_id: selected.internal_product_id,
      }
    : null;

  const decomposeQuery = useQuery({
    queryKey: ["pricing", "decompose", decomposeInput],
    queryFn: () => client.pricingDecompose(decomposeInput as PricingCalcInput),
    enabled: decomposeInput !== null,
  });

  return (
    <div className="flex gap-4 p-4">
      <div className="flex min-w-0 flex-1 flex-col gap-4">
        <header className="flex items-center justify-between">
          <h1 className="text-lg font-semibold text-ink">Preços &amp; Simulador</h1>
          <div className="flex gap-2">
            <button
              type="button"
              data-testid="params-trigger"
              data-deep-linked={paramsOpen ? "1" : undefined}
              onClick={() => setParamsOpen(true)}
              className="rounded-md border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2"
            >
              ⚙ Parâmetros de cálculo
            </button>
            <button
              type="button"
              data-testid="difal-trigger"
              onClick={() => setDifalOpen(true)}
              className="rounded-md border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2"
            >
              DIFAL por UF
            </button>
          </div>
        </header>

        {profileQuery.isError ? (
          <p role="alert" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
            Não foi possível carregar os parâmetros de cálculo — usando o padrão.
          </p>
        ) : null}

        <div className="grid grid-cols-1 gap-4 lg:grid-cols-[280px_1fr]">
          <aside aria-label="Produtos em análise" className="rounded-lg border border-border bg-surface p-3">
            {productsQuery.isLoading ? (
              <LoadingState />
            ) : (
              <ul className="flex flex-col gap-1">
                {products.map((p) => (
                  <li key={p.internal_product_id}>
                    <button
                      type="button"
                      onClick={() => setSelectedId(p.internal_product_id)}
                      aria-pressed={selected?.internal_product_id === p.internal_product_id}
                      className={`w-full rounded-md px-2 py-1.5 text-left text-sm ${
                        selected?.internal_product_id === p.internal_product_id
                          ? "bg-accent-soft text-accent-ink"
                          : "text-ink hover:bg-surface-2"
                      }`}
                    >
                      {productLabel(p)}
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </aside>

          <section aria-label={`Simular · ${selected ? productLabel(selected) : ""}`} className="flex flex-col gap-4">
            {productMissing ? (
              <p role="alert" data-testid="scenario-reload-notice" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
                O produto deste cenário não está na lista atual — selecione outro produto ou recarregue o catálogo.
              </p>
            ) : null}

            <div className="flex flex-wrap items-end gap-4 rounded-lg border border-border bg-surface p-3">
              <label className="flex flex-col text-xs text-muted">
                Preço de venda
                <span className="mt-1 flex items-center gap-1 text-ink">
                  R$
                  <input
                    inputMode="decimal"
                    value={precoInput}
                    onChange={(e) => setPrecoInput(e.target.value)}
                    placeholder={selected?.current_price.amount ?? "0,00"}
                    className="w-28 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink"
                    aria-label="Preço de venda"
                  />
                </span>
              </label>

              <fieldset className="flex flex-col text-xs text-muted">
                <legend className="mb-1">Modalidade</legend>
                <div className="flex gap-1">
                  {MODALIDADES.map((m) => (
                    <button
                      key={m.key}
                      type="button"
                      onClick={() => setModalidade(m.key)}
                      aria-pressed={modalidade === m.key}
                      className={`rounded-md px-2 py-1 text-sm ${
                        modalidade === m.key ? "bg-accent-soft text-accent-ink" : "text-ink hover:bg-surface-2"
                      }`}
                    >
                      {m.label}
                    </button>
                  ))}
                </div>
              </fieldset>
            </div>

            <div data-testid="region-decomposicao" className="rounded-lg border border-border bg-surface p-3">
              {decomposeQuery.isLoading ? (
                <LoadingState />
              ) : decomposeQuery.isError ? (
                <ErrorState onRetry={() => void decomposeQuery.refetch()} detail="Falha ao decompor o preço." />
              ) : decomposeQuery.data ? (
                <DecompositionPanel
                  decomposition={decomposeQuery.data.decomposition}
                  profile={profile}
                  blockingState={decomposeQuery.data.blocking_state}
                  difalUf={profile.difal_destino_uf}
                />
              ) : (
                <p className="text-sm text-muted">Selecione um produto e um preço para simular.</p>
              )}
            </div>

            <div data-testid="region-solver" className="rounded-lg border border-border bg-surface p-3">
              <h3 className="mb-2 text-sm font-semibold text-ink">Margem alvo → preço</h3>
              <SolverPanel
                productId={selected ? selected.internal_product_id : null}
                modalidade={modalidade}
              />
            </div>

            <div data-testid="region-comparacao" className="rounded-lg border border-border bg-surface p-3">
              <MarketComparison productId={selected ? String(selected.internal_product_id) : null} />
            </div>

            <div data-testid="region-aplicar" className="rounded-lg border border-border bg-surface p-3">
              <ApplyPriceAction installationId={installationId} listingId={listingId} newPrice={precoForApi} />
            </div>

            <div data-testid="region-cenarios" className="rounded-lg border border-border bg-surface p-3">
              <ScenariosPanel payload={scenarioPayload} onReload={applyScenario} />
            </div>
          </section>
        </div>
      </div>

      <ParamsDrawer
        open={paramsOpen}
        profile={profile}
        onSave={(next) => saveProfile.mutate(next)}
        onClose={() => setParamsOpen(false)}
        saving={saveProfile.isPending}
      />

      <DifalDrawer
        open={difalOpen}
        data={difalQuery.data}
        isLoading={difalQuery.isLoading}
        isError={difalQuery.isError}
        onOverride={(uf, interna) => overrideDifal.mutate({ uf, interna })}
        onRetry={() => void difalQuery.refetch()}
        onClose={() => setDifalOpen(false)}
        savingUf={overrideDifal.isPending ? overrideDifal.variables?.uf ?? null : null}
      />
    </div>
  );
}
