import type { JSX } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import { ErrorState } from "@marketplace-central/ui";
import { applyProdutoTab, parseProdutoTab, type ProdutoTab } from "./productQueryState";

const tabs: Array<{ value: ProdutoTab; label: string }> = [
  { value: "veredicto", label: "Veredicto" },
  { value: "anuncios", label: "Anúncios vinculados" },
  { value: "estoque", label: "Estoque" },
];

function TabPanel({ tab, active, children }: { tab: ProdutoTab; active: boolean; children: JSX.Element }) {
  if (!active) return null;
  return (
    <div role="tabpanel" id={`produto-panel-${tab}`} aria-labelledby={`produto-tab-${tab}`}>
      {children}
    </div>
  );
}

/**
 * ProdutoPage mounts at /catalogo/produtos/:productId. This slice (S1) only
 * scaffolds the route: param parse + not-found guard, the hand-rolled
 * tablist with `?tab=` deep-link state, and placeholder panel bodies.
 * S2-S5 replace the placeholders with real widgets.
 */
export function ProdutoPage(): JSX.Element {
  const { productId = "" } = useParams();
  const numericId = Number(productId);
  const [searchParams, setSearchParams] = useSearchParams();

  if (productId === "" || !Number.isFinite(numericId)) {
    return (
      <section aria-labelledby="produto-not-found-title" className="mx-auto flex max-w-3xl flex-col gap-3">
        <h1 id="produto-not-found-title" className="text-2xl font-semibold tracking-tight text-ink">
          Produto não encontrado
        </h1>
        <ErrorState detail="Identificador de produto inválido." onRetry={() => void 0} />
        <Link to="/catalogo" className="font-medium text-accent underline decoration-accent-soft underline-offset-2 hover:text-accent-ink">
          Voltar para o catálogo
        </Link>
      </section>
    );
  }

  const activeTab = parseProdutoTab(searchParams);

  const selectTab = (tab: ProdutoTab) => {
    setSearchParams(applyProdutoTab(searchParams, tab));
  };

  return (
    <section aria-labelledby="produto-title" className="mx-auto flex max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <p className="text-sm font-medium text-accent">Workspace de produto</p>
        <h1 id="produto-title" className="text-2xl font-semibold tracking-tight text-ink">
          Produto {productId}
        </h1>
      </header>

      <div role="tablist" aria-label="Seções do produto" className="flex flex-wrap gap-2 border-b border-border pb-4">
        {tabs.map((tab) => (
          <button
            key={tab.value}
            type="button"
            role="tab"
            id={`produto-tab-${tab.value}`}
            aria-selected={activeTab === tab.value}
            aria-controls={`produto-panel-${tab.value}`}
            className={`rounded-lg border px-3 py-2 text-sm font-medium transition-colors ${
              activeTab === tab.value
                ? "border-accent bg-accent text-white"
                : "border-border bg-surface text-muted hover:border-border-2 hover:bg-surface-2"
            }`}
            onClick={() => selectTab(tab.value)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <TabPanel tab="veredicto" active={activeTab === "veredicto"}>
        <p className="text-sm text-muted">Veredicto — em breve.</p>
      </TabPanel>
      <TabPanel tab="anuncios" active={activeTab === "anuncios"}>
        <p className="text-sm text-muted">Anúncios vinculados — em breve.</p>
      </TabPanel>
      <TabPanel tab="estoque" active={activeTab === "estoque"}>
        <p className="text-sm text-muted">Estoque — em breve.</p>
      </TabPanel>
    </section>
  );
}
