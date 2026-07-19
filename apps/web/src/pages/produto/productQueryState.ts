export type ProdutoTab = "veredicto" | "anuncios" | "estoque";

const DEFAULT_TAB: ProdutoTab = "veredicto";

export function parseProdutoTab(searchParams: URLSearchParams): ProdutoTab {
  const value = searchParams.get("tab");
  return value === "anuncios" || value === "estoque" ? value : DEFAULT_TAB;
}

export function applyProdutoTab(searchParams: URLSearchParams, tab: ProdutoTab): URLSearchParams {
  const next = new URLSearchParams(searchParams);
  if (tab === DEFAULT_TAB) next.delete("tab");
  else next.set("tab", tab);
  return next;
}
