import { useSyncExternalStore } from "react";

// The active ERP source is the demo's two-source toggle: xlsx = Sankhya ERP
// (custo + estoque), catalogo_cliente = the prospect catalog imported through
// /integracoes (honest-unknown custo/estoque, ADR-17). It is persisted in
// localStorage so a reload keeps the operator on the dataset they selected, and
// broadcast on a window event so every mounted view (the /integracoes selector
// and the /catalogo hooks) stays in lockstep without prop-drilling.

export type ActiveErpSource = "xlsx" | "catalogo_cliente";

const STORAGE_KEY = "mc.active-erp-source";
const CHANGE_EVENT = "mc:active-erp-source";
export const DEFAULT_ERP_SOURCE: ActiveErpSource = "xlsx";

function isActiveErpSource(value: unknown): value is ActiveErpSource {
  return value === "xlsx" || value === "catalogo_cliente";
}

export function getActiveErpSource(): ActiveErpSource {
  try {
    const stored = window.localStorage.getItem(STORAGE_KEY);
    return isActiveErpSource(stored) ? stored : DEFAULT_ERP_SOURCE;
  } catch {
    return DEFAULT_ERP_SOURCE;
  }
}

export function setActiveErpSource(value: ActiveErpSource): void {
  try {
    window.localStorage.setItem(STORAGE_KEY, value);
  } catch {
    // A blocked localStorage (private mode) still gets the in-session broadcast
    // below; the selection just does not survive a reload. Never throw here.
  }
  if (typeof window !== "undefined") {
    window.dispatchEvent(new CustomEvent(CHANGE_EVENT));
  }
}

function subscribe(callback: () => void): () => void {
  window.addEventListener(CHANGE_EVENT, callback);
  window.addEventListener("storage", callback);
  return () => {
    window.removeEventListener(CHANGE_EVENT, callback);
    window.removeEventListener("storage", callback);
  };
}

export function useActiveErpSource(): [ActiveErpSource, (value: ActiveErpSource) => void] {
  const value = useSyncExternalStore(subscribe, getActiveErpSource, () => DEFAULT_ERP_SOURCE);
  return [value, setActiveErpSource];
}
