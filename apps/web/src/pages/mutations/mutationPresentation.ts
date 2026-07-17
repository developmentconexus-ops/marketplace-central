import type { MutationType } from "@marketplace-central/sdk-runtime";

export const mutationTypeLabels: Partial<Record<MutationType, string>> = {
  price_update: "Atualizar preço",
  stock_correct: "Corrigir estoque",
  listing_pause: "Pausar anúncios",
  listing_resync: "Ressincronizar anúncios",
  link_apply: "Aplicar vínculo",
  listing_edit: "Editar anúncios",
};

export function presentMutationValue(value: unknown): string {
  if (typeof value === "string") return value;
  if (typeof value === "number" || typeof value === "boolean") return String(value);
  return JSON.stringify(value);
}

export function mutationError(error: unknown): { code: string; message?: string } {
  if (typeof error !== "object" || error === null) return { code: "internal" };
  const body = (error as { error?: { code?: unknown; message?: unknown } }).error;
  return {
    code: typeof body?.code === "string" ? body.code : "internal",
    message: typeof body?.message === "string" ? body.message : undefined,
  };
}
