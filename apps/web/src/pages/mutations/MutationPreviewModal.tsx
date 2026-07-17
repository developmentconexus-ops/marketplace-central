import { useMutation } from "@tanstack/react-query";
import type { CreateMutationRequest, MutationPreview, MutationProtocol, MutationType } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { failureCopy } from "@marketplace-central/web-query";
import { useState } from "react";
import { useClient } from "../../app/ClientContext";
import { MutationIntentForm } from "./MutationIntentForm";
import { MutationItemsTable } from "./MutationItemsTable";
import { mutationError, mutationTypeLabels, presentMutationValue } from "./mutationPresentation";

type InteractionStep = "intent" | "previewing" | "preview-shown" | "approving" | "applying" | "terminal" | "error";

export interface MutationPreviewModalProps {
  open: boolean;
  type: MutationType;
  installationId: string;
  selectedIds: string[];
  onClose: () => void;
}

export function MutationPreviewModal({ open, type, installationId, selectedIds, onClose }: MutationPreviewModalProps) {
  const client = useClient();
  const [step, setStep] = useState<InteractionStep>("intent");
  const [draft, setDraft] = useState<MutationProtocol | null>(null);
  const [preview, setPreview] = useState<MutationPreview | null>(null);
  const [errorCode, setErrorCode] = useState<string | null>(null);

  const createOperation = useMutation({ mutationFn: (request: CreateMutationRequest) => client.createMutation(request) });
  const previewOperation = useMutation({ mutationFn: (protocolId: string) => client.previewMutation(protocolId) });
  const cancelOperation = useMutation({ mutationFn: (protocolId: string) => client.cancelMutation(protocolId) });

  if (!open) return null;

  const createPreview = async (intent: Record<string, unknown>) => {
    setStep("previewing");
    setErrorCode(null);
    let created: MutationProtocol | null = null;
    try {
      created = await createOperation.mutateAsync({
        installation_id: installationId,
        type,
        actor: "operator_supplied_unverified",
        intent,
        selection: { mode: "explicit", listing_ids: selectedIds },
      });
      setDraft(created);
      const response = await previewOperation.mutateAsync(created.protocol_id);
      setPreview(response);
      setStep("preview-shown");
    } catch (error) {
      const { code } = mutationError(error);
      if (created && code === "selection_too_large") {
        await cancelOperation.mutateAsync(created.protocol_id).catch(() => undefined);
        setDraft(null);
        setPreview(null);
      }
      setErrorCode(code);
      setStep("error");
    }
  };

  const cancel = async () => {
    if (!draft) {
      onClose();
      return;
    }
    await cancelOperation.mutateAsync(draft.protocol_id);
    setDraft(null);
    setPreview(null);
    onClose();
  };

  const resetToIntent = () => {
    setErrorCode(null);
    setStep("intent");
  };

  const retryAfterError = async () => {
    if (!draft) {
      resetToIntent();
      return;
    }
    setStep("previewing");
    setErrorCode(null);
    try {
      const response = await previewOperation.mutateAsync(draft.protocol_id);
      setPreview(response);
      setStep("preview-shown");
    } catch (error) {
      setErrorCode(mutationError(error).code);
      setStep("error");
    }
  };

  return (
    <div className="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/50 p-4" role="presentation">
      <section role="dialog" aria-modal="true" aria-labelledby="mutation-modal-title" className="flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-xl bg-white shadow-lg">
        <header className="border-b border-slate-200 px-5 py-4">
          <h2 id="mutation-modal-title" className="text-lg font-semibold text-slate-950">
            {mutationTypeLabels[type] ?? "Alterar anúncios"}
          </h2>
          <p className="mt-1 text-sm text-slate-600">{selectedIds.length} anúncio(s) selecionado(s)</p>
        </header>

        <div className="min-h-0 flex-1 overflow-y-auto p-5">
          {step === "intent" ? <MutationIntentForm type={type} onSubmit={(intent) => void createPreview(intent)} /> : null}
          {step === "previewing" ? <LoadingState /> : null}
          {step === "preview-shown" && preview ? <PreviewContent preview={preview} /> : null}
          {step === "error" && errorCode === "selection_too_large" ? (
            <ErrorState
              detail={`${selectedIds.length} anúncios selecionados; reduza a seleção ou refine o filtro para gerar uma nova prévia.`}
              onRetry={resetToIntent}
            />
          ) : null}
          {step === "error" && errorCode !== "selection_too_large" ? (
            <ErrorState detail={failureCopy(errorCode ?? "internal")} onRetry={() => void retryAfterError()} />
          ) : null}
        </div>

        <footer className="flex flex-wrap justify-end gap-2 border-t border-slate-200 px-5 py-4">
          {(step === "intent" || step === "preview-shown") ? (
            <button type="button" onClick={() => void cancel()} disabled={cancelOperation.isPending} className="rounded-lg border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 disabled:opacity-50">
              Cancelar
            </button>
          ) : null}
          {step === "intent" ? (
            <button type="submit" form="mutation-intent-form" className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700">
              Gerar prévia
            </button>
          ) : null}
          {step === "preview-shown" ? (
            <button type="button" disabled className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-50">
              Confirmar e aplicar
            </button>
          ) : null}
        </footer>
      </section>
    </div>
  );
}

function PreviewContent({ preview }: { preview: MutationPreview }) {
  const totals = preview.totals;
  return (
    <div className="space-y-5">
      <div className="flex flex-wrap gap-2" aria-label="Totais da prévia">
        <Total value={totals.items} label="itens" />
        <Total value={totals.previewed} label="previstos" />
        <Total value={totals.failed} label="falhas" />
      </div>
      <MutationItemsTable items={preview.items} />
    </div>
  );
}

function Total({ value, label }: { value: unknown; label: string }) {
  if (value === undefined || value === null) return null;
  return <span className="rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700">{presentMutationValue(value)} {label}</span>;
}
