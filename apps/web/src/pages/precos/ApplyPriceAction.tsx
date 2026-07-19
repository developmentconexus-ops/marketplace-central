import { useState } from "react";
import type { JSX } from "react";
import { useMutation } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { ErrorState } from "@marketplace-central/ui";
import type { CreateMutationRequest, MutationPreview } from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";

export interface ApplyPriceActionProps {
  /** Installation that owns the listing; empty until an account is resolved. */
  installationId: string;
  /** ML listing to reprice, or null when the product has no linked anúncio. */
  listingId: string | null;
  /** Simulated price to apply (decimal string); empty disables the action. */
  newPrice: string;
}

/**
 * ApplyPriceAction stages a price change as a `price_update` protocol and takes
 * it NO FURTHER than `previewed` — the surface exposes no approve/apply control
 * on purpose. The `price_update` executor writes to Mercado Livre; on the demo
 * that step is forbidden. The point is to show governance: a protocol parked at
 * `previewed`, awaiting human approval, with zero writes leaving the platform.
 */
export function ApplyPriceAction({ installationId, listingId, newPrice }: ApplyPriceActionProps): JSX.Element {
  const client = useClient();
  const [result, setResult] = useState<MutationPreview | null>(null);

  const apply = useMutation({
    mutationFn: async (): Promise<MutationPreview> => {
      const request: CreateMutationRequest = {
        installation_id: installationId,
        type: "price_update",
        actor: "operator_supplied_unverified",
        intent: { new_price: { amount: newPrice, currency: "BRL" } },
        selection: { mode: "explicit", listing_ids: [listingId as string] },
      };
      const created = await client.createMutation(request);
      // Cap at previewed: build the preview, then stop. Never approve/execute.
      return client.previewMutation(created.protocol_id);
    },
    onSuccess: (preview) => setResult(preview),
  });

  const priceMissing = newPrice.trim() === "";
  const noListing = !listingId;
  const disabled = noListing || priceMissing || !installationId || apply.isPending;

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-ink">Aplicar preço</h3>
        <button
          type="button"
          data-testid="apply-trigger"
          disabled={disabled}
          onClick={() => apply.mutate()}
          className="rounded-md bg-accent px-3 py-1.5 text-sm font-medium text-accent-ink hover:bg-accent/90 disabled:opacity-40"
        >
          {apply.isPending ? "Preparando…" : "Aplicar preço"}
        </button>
      </div>

      {noListing ? (
        <p data-testid="apply-hint" className="text-xs text-muted">
          Selecione um produto com anúncio vinculado para preparar a aplicação.
        </p>
      ) : (
        <p className="text-xs text-muted">
          A aplicação para no estado <span className="font-mono text-ink">previewed</span> — aprovação é manual e fora
          desta tela.
        </p>
      )}

      {apply.isError ? (
        <ErrorState onRetry={() => apply.mutate()} detail="Não foi possível preparar a aplicação do preço." />
      ) : null}

      {result ? (
        <div
          data-testid="apply-result"
          role="status"
          className="flex flex-col gap-1 rounded-md border border-border bg-surface-2 p-3 text-sm"
        >
          <div className="flex items-center gap-2">
            <span data-testid="apply-state" className="rounded bg-amber-soft px-1.5 py-0.5 text-xs font-medium text-amber">
              {result.state}
            </span>
            <Link to={`/protocolos/${result.protocol_id}`} className="font-mono text-accent hover:underline">
              {result.protocol_id}
            </Link>
          </div>
          <p className="text-xs text-muted">
            Aguardando aprovação — nada foi enviado ao Mercado Livre.
          </p>
        </div>
      ) : null}
    </div>
  );
}
