import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import type {
  CatalogProductFact,
  PricingCalcInput,
} from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";

/** ML commission by modalidade (comes from ML — not operator-editable). */
export const MODALIDADES = [
  { key: "classico", label: "Clássico", comissaoPct: "12" },
  { key: "premium", label: "Premium", comissaoPct: "17" },
  { key: "full", label: "Full", comissaoPct: "17" },
] as const;

export type ModalidadeKey = (typeof MODALIDADES)[number]["key"];

function comissaoFor(modalidade: ModalidadeKey): string {
  return MODALIDADES.find((m) => m.key === modalidade)?.comissaoPct ?? "17";
}

function productLabel(p: CatalogProductFact): string {
  return p.description ?? p.manufacturer_reference ?? `#${p.internal_product_id}`;
}

/**
 * PricingPage is the /precos IC-04 simulator shell. It loads the calc profile
 * and the catalog products, owns the selected product + price/modalidade state,
 * and lays out the regions the later slices flesh out: the per-component
 * decomposition, the parâmetros/DIFAL drawers (trigger here), the market
 * comparison, and the "aplicar preço" action. All money/percentages stay as
 * decimal strings end-to-end (never float).
 */
export function PricingPage() {
  const client = useClient();
  const [searchParams] = useSearchParams();

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [modalidade, setModalidade] = useState<ModalidadeKey>("premium");
  const [precoInput, setPrecoInput] = useState<string>("");

  const profileQuery = useQuery({
    queryKey: ["pricing", "profile"],
    queryFn: () => client.getPricingProfile(),
  });

  const productsQuery = useQuery({
    queryKey: ["pricing", "catalog-facts"],
    queryFn: () => client.listCatalogProductFacts({ limit: 50 }),
  });

  const products = productsQuery.data?.items ?? [];
  const selected = useMemo<CatalogProductFact | null>(() => {
    if (products.length === 0) return null;
    const byId = products.find((p) => p.internal_product_id === selectedId);
    return byId ?? products[0];
  }, [products, selectedId]);

  // Working price: explicit input wins, else the product's current price.
  const preco = precoInput.trim() !== "" ? precoInput.trim() : selected?.current_price.amount ?? "";

  const decomposeInput: PricingCalcInput | null = selected && preco !== ""
    ? {
        preco,
        comissao_pct: comissaoFor(modalidade),
        modalidade,
        product_id: selected.internal_product_id,
      }
    : null;

  const decomposeQuery = useQuery({
    queryKey: ["pricing", "decompose", decomposeInput],
    queryFn: () => client.pricingDecompose(decomposeInput as PricingCalcInput),
    enabled: decomposeInput !== null,
  });

  const paramsDeepLinked = searchParams.get("params") === "1";

  return (
    <div className="flex flex-col gap-4 p-4">
      <header className="flex items-center justify-between">
        <h1 className="text-lg font-semibold text-ink">Preços &amp; Simulador</h1>
        <button
          type="button"
          data-testid="params-trigger"
          data-deep-linked={paramsDeepLinked ? "1" : undefined}
          className="rounded-md border border-border px-3 py-1.5 text-sm text-ink hover:bg-surface-2"
        >
          ⚙ Parâmetros de cálculo
        </button>
      </header>

      {profileQuery.isError ? (
        <p role="alert" className="rounded-md bg-warn-soft px-3 py-2 text-sm text-warn">
          Não foi possível carregar os parâmetros de cálculo — usando o padrão.
        </p>
      ) : null}

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-[280px_1fr]">
        <aside aria-label="Produtos em análise" className="rounded-lg border border-border bg-surface p-3">
          {productsQuery.isLoading ? (
            <LoadingState label="Carregando produtos" />
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
              <LoadingState label="Calculando decomposição" />
            ) : decomposeQuery.isError ? (
              <ErrorState title="Falha ao decompor o preço" />
            ) : decomposeQuery.data ? (
              <p className="font-mono text-sm text-ink">
                Retorno /un{" "}
                <span data-testid="decomp-margem">
                  {decomposeQuery.data.decomposition.margem_valor ?? "—"}
                </span>
              </p>
            ) : (
              <p className="text-sm text-muted">Selecione um produto e um preço para simular.</p>
            )}
          </div>

          <div data-testid="region-comparacao" className="rounded-lg border border-border bg-surface p-3">
            <p className="text-sm text-muted">Comparação de mercado</p>
          </div>

          <div data-testid="region-aplicar" className="rounded-lg border border-border bg-surface p-3">
            <p className="text-sm text-muted">Aplicar preço</p>
          </div>
        </section>
      </div>
    </div>
  );
}
