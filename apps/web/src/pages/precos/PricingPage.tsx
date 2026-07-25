import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { ErrorState, LoadingState, formatPercent } from "@marketplace-central/ui";
import type {
  CatalogProductFact,
  PricingCalcInput,
  PricingCalcProfile,
} from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { useInstallation } from "../../app/InstallationContext";
import { DecompositionPanel } from "./DecompositionPanel";
import { SolverPanel } from "./SolverPanel";
import { ParamsDrawer } from "./ParamsDrawer";
import { DifalDrawer } from "./DifalDrawer";
import { MarketComparison } from "./MarketComparison";
import { QuickChips } from "./QuickChips";
import { PricingMatrix } from "./PricingMatrix";
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

/** Server-side cap on the explicit-ids catalog ask (OpenAPI `ids`: 1..100). */
const CATALOG_IDS_MAX = 100;

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

  // Deep link: /precos?produto=<codprod> opens the simulator on that product —
  // this is where /mercado's "Simular" lands. The id may sit outside the priced
  // block the list query fetches, so it is resolved on its own below.
  const deepLinkProductId = Number.parseInt(searchParams.get("produto") ?? "", 10);
  const requestedProductId = Number.isFinite(deepLinkProductId) ? deepLinkProductId : null;

  const [selectedId, setSelectedId] = useState<number | null>(requestedProductId);
  const [modalidade, setModalidade] = useState<ModalidadeKey>("premium");
  const [precoInput, setPrecoInput] = useState<string>("");
  const [solverTarget, setSolverTarget] = useState<string>("");
  const [pickerOpen, setPickerOpen] = useState<boolean>(false);
  // Deep link: /precos?params=1 opens the Parâmetros drawer on mount.
  const [paramsOpen, setParamsOpen] = useState<boolean>(searchParams.get("params") === "1");
  const [difalOpen, setDifalOpen] = useState<boolean>(false);

  const profileQuery = useQuery({
    queryKey: ["pricing", "profile"],
    queryFn: () => client.getPricingProfile(),
  });
  const profile = profileQuery.data ?? DEFAULT_PROFILE;

  // The workspace's selected account (never items[0]: an abandoned authorization
  // leaves a pending_connection installation that sorts first and carries no
  // listings, which would empty this whole screen). See InstallationProvider.
  const { installationId, status: installationStatus } = useInstallation();

  // The analysis set is exactly the products that HAVE a market to be priced
  // against: those carrying a resolved link to an ML listing. It is derived from
  // the links (the domain fact), never from the catalog id layout.
  const linkedListingsQuery = useQuery({
    queryKey: ["pricing", "linked-listings", installationId],
    queryFn: () =>
      client.listListings({ installation_id: installationId, link_state: "resolved", limit: 100 }),
    enabled: installationId !== "",
  });

  // Deep-linked ids ride along in the same ask so /mercado's "Simular" resolves
  // even for a product with no listing yet — one request, no id-layout tricks.
  const analysisIds = useMemo<number[]>(() => {
    const ids: number[] = [];
    const seen = new Set<number>();
    const add = (raw: number) => {
      if (!Number.isFinite(raw) || raw <= 0 || seen.has(raw)) return;
      seen.add(raw);
      ids.push(raw);
    };
    if (requestedProductId !== null) add(requestedProductId);
    for (const listing of linkedListingsQuery.data?.items ?? []) {
      add(Number.parseInt(listing.link.product_id ?? "", 10));
    }
    return ids.slice(0, CATALOG_IDS_MAX);
  }, [linkedListingsQuery.data, requestedProductId]);

  const productsQuery = useQuery({
    queryKey: ["pricing", "catalog-facts", analysisIds],
    queryFn: () => client.catalogProductFactsByIds({ ids: analysisIds }),
    enabled: analysisIds.length > 0,
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

  // The ids ask answers with exactly the products the ACTIVE source carries: an
  // id with no row there is simply absent, which is what makes the "produto fora
  // da lista" notice below honest instead of a silent substitution.
  const products = productsQuery.data?.items ?? [];

  const selected = useMemo<CatalogProductFact | null>(() => {
    if (products.length === 0) return null;
    // No explicit pick yet ⇒ default to the first product. An explicit id that
    // isn't in the loaded list ⇒ null (honest empty), NEVER a silent fallback to
    // products[0] — that would show a different product's decomposition.
    if (selectedId === null) return products[0];
    return products.find((p) => p.internal_product_id === selectedId) ?? null;
  }, [products, selectedId]);

  // "Aplicar preço" targets the ML listing of the selected product.
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

  // An explicit product pick the active source can't answer for (a deep link to a
  // product outside it, or a scenario reloaded after a source flip) ⇒ tell the
  // operator instead of silently decomposing a different product at this price.
  // Derived (not stored) so it can't go stale against a still-loading catalog.
  // While the lookup is in flight the product is unknown, not missing — claiming
  // "não encontrado" there would be a lie that resolves itself.
  const analysisPending =
    installationStatus === "loading" ||
    (installationId !== "" && linkedListingsQuery.isPending) ||
    (analysisIds.length > 0 && productsQuery.isPending);
  const productMissing = selectedId !== null && selected === null && !analysisPending;

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

  // Header resumo — a compact echo of the profile that drives the calc surface
  // (design §Simulador). Regime label + alíquota, plus a DIFAL flag when active.
  const regimeLabel = profile.regime === "SIMPLES" ? "Simples" : "L. Presumido";
  const resumoParams =
    `custo ERP · comissões ML · frete por peso · ${regimeLabel} ${profile.aliquota_pct}%` +
    (profile.difal_enabled ? " · DIFAL" : "");

  // ML free-shipping rule (static, not fabricated): the seller pays freight at/above
  // R$ 79; below it the buyer pays but a R$ 6,50 fixed fee applies.
  const precoNum = Number.parseFloat(precoForApi);
  const freteFreeTier = Number.isFinite(precoNum) && precoNum >= 79;
  const freteNota = freteFreeTier
    ? "≥ R$ 79: frete grátis obrigatório — o vendedor paga o envio."
    : "< R$ 79: o comprador paga o frete, mas incide a taxa fixa de R$ 6,50.";

  const simName = selected ? productLabel(selected) : "";
  const simSku = selected?.reference ?? (selected ? `#${selected.internal_product_id}` : "");

  return (
    <div className="flex flex-col gap-3.5 p-6">
      <header className="flex flex-wrap items-center gap-2.5">
        <h1 className="text-[22px] font-bold tracking-tight text-ink">Preços &amp; Simulador</h1>

        {/* Produtos picker pill — the products currently in the analysis; clicking
            one selects it (mirrors the design's "Produtos: N ▾"). */}
        <div className="relative">
          <button
            type="button"
            data-testid="produtos-pill"
            aria-expanded={pickerOpen}
            onClick={() => setPickerOpen((v) => !v)}
            className={`cursor-pointer whitespace-nowrap rounded-pill border px-3 py-1 text-[12.5px] font-semibold ${
              pickerOpen ? "border-accent bg-accent-soft text-accent-ink" : "border-border bg-surface text-muted"
            }`}
          >
            Produtos: {products.length} ▾
          </button>
          {pickerOpen ? (
            <div
              data-testid="produtos-dropdown"
              className="absolute left-0 top-9 z-20 flex w-[340px] flex-col gap-2 rounded-card border border-border bg-surface p-3.5 shadow-lg"
            >
              <div className="text-[10.5px] font-semibold tracking-wider text-faint">NA ANÁLISE</div>
              {products.length === 0 ? (
                <p className="text-xs text-muted">Nenhum produto vinculado a um anúncio no momento.</p>
              ) : (
                products.map((p) => (
                  <button
                    key={p.internal_product_id}
                    type="button"
                    data-testid={`produtos-option-${p.internal_product_id}`}
                    onClick={() => {
                      setSelectedId(p.internal_product_id);
                      setPickerOpen(false);
                    }}
                    className={`flex items-center gap-2 rounded-control px-1 py-0.5 text-left text-[12.5px] hover:bg-surface-2 ${
                      selected?.internal_product_id === p.internal_product_id ? "text-accent-ink" : "text-ink"
                    }`}
                  >
                    <span className="font-mono text-[11px] text-faint">{p.reference ?? `#${p.internal_product_id}`}</span>
                    <span className="flex-1 truncate">{productLabel(p)}</span>
                  </button>
                ))
              )}
            </div>
          ) : null}
        </div>

        {/* Shipping/destino context — honest: the DIFAL destino UF the calc uses,
            or "—" when none is configured (never a fabricated CEP, ADR-17). */}
        <span
          data-testid="cep-chip"
          className="whitespace-nowrap rounded-pill border border-border px-3 py-1 text-[12.5px] text-muted"
        >
          CEP → {profile.difal_destino_uf ?? "—"}
        </span>

        <span className="ml-auto flex items-center gap-2">
          <span className="whitespace-nowrap text-[11.5px] text-faint">{resumoParams}</span>
          <button
            type="button"
            data-testid="params-trigger"
            data-deep-linked={paramsOpen ? "1" : undefined}
            onClick={() => setParamsOpen(true)}
            className={`cursor-pointer whitespace-nowrap rounded-control border px-3 py-1.5 text-[12.5px] font-semibold ${
              paramsOpen
                ? "border-accent bg-accent-soft text-accent-ink"
                : "border-border text-muted hover:border-accent hover:text-accent-ink"
            }`}
          >
            ⚙ Parâmetros de cálculo
          </button>
        </span>
      </header>

      {profileQuery.isError ? (
        <p role="alert" className="rounded-control bg-warn-soft px-3 py-2 text-sm text-warn">
          Não foi possível carregar os parâmetros de cálculo — usando o padrão.
        </p>
      ) : null}

      <div className="flex flex-wrap items-start gap-3.5">
        <div className="min-w-[480px] flex-1">
          {analysisPending ? (
            <LoadingState />
          ) : products.length === 0 ? (
            /* No product to price is a FACT, not a blank table: say which of the
               two causes it is instead of rendering empty headers (ADR-17). */
            <p
              role="status"
              data-testid="pricing-empty"
              className="rounded-card border border-border bg-surface px-4 py-6 text-sm text-muted"
            >
              {analysisIds.length === 0
                ? "Nenhum anúncio vinculado a um produto — vincule em Vínculos para simular preços."
                : "Os produtos vinculados não estão na fonte ativa do catálogo — troque a fonte em Integrações ou reimporte."}
            </p>
          ) : (
            <PricingMatrix
              products={products}
              selectedId={selectedId}
              onSelect={setSelectedId}
              modalidade={modalidade}
              profile={profile}
              installationId={installationId}
            />
          )}
        </div>

        {/* The simular panel is the PRESERVED single-product surface, now folded into
            ONE 380px card (design parity). It renders while a product is selected
            (products[0] on mount) OR when a reloaded scenario's product is missing —
            so the honest "product not in list" notice is never swallowed. */}
        {selected !== null || productMissing ? (
          <aside
            aria-label={`Simular · ${simName}`}
            className="w-[380px] shrink-0 overflow-hidden rounded-card border border-border bg-surface"
          >
            <div className="flex items-center gap-2 border-b border-border bg-surface-2 px-4 py-2.5">
              <b className="text-[13px] text-ink">Simular · {simName}</b>
              <span className="ml-auto font-mono text-[11px] text-faint">{simSku}</span>
            </div>

            <div className="flex flex-col gap-3.5 px-4 py-3.5">
              {productMissing ? (
                <p role="alert" data-testid="scenario-reload-notice" className="rounded-control bg-warn-soft px-3 py-2 text-sm text-warn">
                  Este produto não está na fonte ativa do catálogo — selecione outro produto ou troque a fonte em Integrações.
                </p>
              ) : null}

              {/* Preço + Margem desejada grid */}
              <div className="grid grid-cols-2 gap-2.5">
                <div>
                  <div className="mb-1 text-[11px] text-faint">Preço de venda</div>
                  <div className="flex items-center gap-1 rounded-control border-[1.5px] border-accent px-2.5 py-1.5">
                    <span className="text-xs text-faint">R$</span>
                    <input
                      inputMode="decimal"
                      value={precoInput}
                      onChange={(e) => setPrecoInput(e.target.value)}
                      placeholder={selected?.current_price.amount ?? "0,00"}
                      aria-label="Preço de venda"
                      className="w-full border-none bg-transparent font-mono text-[15px] font-semibold text-ink outline-none"
                    />
                  </div>
                </div>
                <div>
                  <div className="mb-1 text-[11px] text-faint">Margem desejada</div>
                  <div className="flex items-center gap-1 rounded-control border border-border px-2.5 py-1.5">
                    <input
                      inputMode="decimal"
                      value={solverTarget}
                      onChange={(e) => setSolverTarget(e.target.value)}
                      placeholder="0,0"
                      aria-label="Margem alvo"
                      className="w-full border-none bg-transparent text-right font-mono text-[15px] font-semibold text-ink outline-none"
                    />
                    <span className="text-xs text-faint">%</span>
                  </div>
                </div>
              </div>

              {/* Quick-seed chips (mediana / menor conc. / cobrir X%) */}
              <QuickChips
                productId={selected ? String(selected.internal_product_id) : null}
                coverPct={profile.limiar_verde_pct}
                onSeedPreco={(amount) => setPrecoInput(amount)}
                onSeedTarget={(pct) => setSolverTarget(pct)}
              />

              {/* MODALIDADE segmented — active tier shows its resolved margin %,
                  the others stay honest "—" (never a fabricated per-modalidade rate). */}
              <div>
                <div className="mb-1.5 text-[10.5px] font-semibold tracking-wider text-faint">MODALIDADE</div>
                <div className="grid grid-cols-3 gap-1.5">
                  {MODALIDADES.map((m) => {
                    const active = modalidade === m.key;
                    const pct = active ? decomposeQuery.data?.decomposition.margem_pct ?? null : null;
                    return (
                      <button
                        key={m.key}
                        type="button"
                        onClick={() => setModalidade(m.key)}
                        aria-pressed={active}
                        className={`cursor-pointer rounded-control border-[1.5px] px-2 py-1.5 text-center ${
                          active ? "border-accent bg-accent-soft" : "border-border-2 hover:bg-surface-2"
                        }`}
                      >
                        <div className={`text-[11px] font-semibold ${active ? "text-accent-ink" : "text-muted"}`}>{m.label}</div>
                        <div className="mt-0.5 font-mono text-xs font-bold text-ink">{formatPercent(pct) ?? "—"}</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Per-component decomposition (mono breakdown card) */}
              <div data-testid="region-decomposicao">
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
                    tarifa={decomposeQuery.data.tarifa}
                  />
                ) : (
                  <p className="text-sm text-muted">Selecione um produto e um preço para simular.</p>
                )}
              </div>

              {/* Margem-alvo → preço solver (input lifted to the grid above) */}
              <div data-testid="region-solver">
                <SolverPanel
                  productId={selected ? selected.internal_product_id : null}
                  modalidade={modalidade}
                  target={solverTarget}
                  onTargetChange={setSolverTarget}
                  hideInput
                />
              </div>

              {/* ML free-shipping rule note */}
              <div
                data-testid="frete-nota"
                className={`rounded-control px-2.5 py-1.5 text-[11.5px] ${
                  freteFreeTier ? "bg-amber-soft text-amber" : "bg-accent-soft text-accent-ink"
                }`}
              >
                {freteNota}
              </div>

              <div data-testid="region-comparacao">
                <MarketComparison productId={selected ? String(selected.internal_product_id) : null} />
              </div>

              <div data-testid="region-aplicar">
                <ApplyPriceAction installationId={installationId} listingId={listingId} newPrice={precoForApi} />
              </div>

              <div data-testid="region-cenarios">
                <ScenariosPanel payload={scenarioPayload} onReload={applyScenario} />
              </div>
            </div>
          </aside>
        ) : null}
      </div>

      <ParamsDrawer
        open={paramsOpen}
        profile={profile}
        defaults={DEFAULT_PROFILE}
        modalidade={modalidade}
        comissaoTarifa={decomposeQuery.data?.tarifa?.comissao}
        onSave={(next) => saveProfile.mutate(next)}
        onClose={() => setParamsOpen(false)}
        onOpenDifal={() => {
          setParamsOpen(false);
          setDifalOpen(true);
        }}
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
