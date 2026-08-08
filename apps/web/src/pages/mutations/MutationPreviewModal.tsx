import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  hasCode,
  type CreateMutationRequest,
  type MutationPreview,
  type MutationProtocol,
  type MutationType,
} from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { failureCopy, mutationsQueryKeys } from "@marketplace-central/web-query";
import { useEffect, useRef, useState } from "react";
import { useClient } from "../../app/ClientContext";
import { MutationIntentForm } from "./MutationIntentForm";
import { MutationItemsTable } from "./MutationItemsTable";
import { MutationResultSummary } from "./MutationResultSummary";
import { mutationError, mutationTypeLabels, presentMutationValue } from "./mutationPresentation";
import { isMutationTerminal, useMutationProtocol } from "./useMutationProtocol";

type InteractionStep =
  "intent" | "previewing" | "preview-shown" | "approving" | "applying" | "terminal" | "error";

type PreviewError =
  { kind: "preview_stale" } | { kind: "selection_too_large" } | { kind: "other"; code: string };

const classifyPreviewError = (error: unknown): PreviewError =>
  hasCode(error, "selection_too_large")
    ? { kind: "selection_too_large" }
    : hasCode(error, "preview_stale")
      ? { kind: "preview_stale" }
      : { kind: "other", code: mutationError(error).code };

const displayCode = (e: PreviewError | null): string =>
  e === null ? "internal" : e.kind === "other" ? e.code : e.kind;

export interface MutationPreviewModalProps {
  open: boolean;
  type: MutationType;
  installationId: string;
  selectedIds: string[];
  onClose: () => void;
}

export function MutationPreviewModal({
  open,
  type,
  installationId,
  selectedIds,
  onClose,
}: MutationPreviewModalProps) {
  const client = useClient();
  const queryClient = useQueryClient();
  const [step, setStep] = useState<InteractionStep>("intent");
  const [draft, setDraft] = useState<MutationProtocol | null>(null);
  const [preview, setPreview] = useState<MutationPreview | null>(null);
  const [previewError, setPreviewError] = useState<PreviewError | null>(null);
  const [confirmed, setConfirmed] = useState(false);
  const approving = useRef(false);

  const createOperation = useMutation({
    mutationFn: (request: CreateMutationRequest) => client.createMutation(request),
  });
  const previewOperation = useMutation({
    mutationFn: (protocolId: string) => client.previewMutation(protocolId),
  });
  const cancelOperation = useMutation({
    mutationFn: (protocolId: string) => client.cancelMutation(protocolId),
  });
  const approveOperation = useMutation({
    mutationFn: (protocolId: string) => client.approveMutation(protocolId, { execute: true }),
  });

  if (!open) return null;

  const createPreview = async (intent: Record<string, unknown>) => {
    setStep("previewing");
    setPreviewError(null);
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
      const classified = classifyPreviewError(error);
      if (created && classified.kind === "selection_too_large") {
        await cancelOperation.mutateAsync(created.protocol_id).catch(() => undefined);
        setDraft(null);
        setPreview(null);
      }
      setPreviewError(classified);
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
    setPreviewError(null);
    setStep("intent");
  };

  const retryAfterError = async () => {
    if (!draft) {
      resetToIntent();
      return;
    }
    setStep("previewing");
    setPreviewError(null);
    try {
      const response = await previewOperation.mutateAsync(draft.protocol_id);
      setPreview(response);
      setStep("preview-shown");
    } catch (error) {
      setPreviewError(classifyPreviewError(error));
      setStep("error");
    }
  };

  const approve = async () => {
    if (!draft || !confirmed || approving.current) return;
    approving.current = true;
    setStep("approving");
    setPreviewError(null);
    try {
      const protocol = await approveOperation.mutateAsync(draft.protocol_id);
      queryClient.setQueryData(mutationsQueryKeys.detail(draft.protocol_id), protocol);
      setStep(isMutationTerminal(protocol.state) ? "terminal" : "applying");
    } catch (error) {
      if (hasCode(error, "preview_stale") && error.status === 409) {
        setConfirmed(false);
        setPreviewError({ kind: "preview_stale" });
        setStep("preview-shown");
      } else {
        setPreviewError(classifyPreviewError(error));
        setStep("error");
      }
    } finally {
      approving.current = false;
    }
  };

  return (
    <div
      className="fixed inset-0 z-40 flex items-center justify-center bg-slate-950/50 p-4"
      role="presentation"
    >
      <section
        role="dialog"
        aria-modal="true"
        aria-labelledby="mutation-modal-title"
        className="flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-xl bg-white shadow-lg"
      >
        <header className="border-b border-slate-200 px-5 py-4">
          <h2 id="mutation-modal-title" className="text-lg font-semibold text-slate-950">
            {mutationTypeLabels[type] ?? "Alterar anúncios"}
          </h2>
          <p className="mt-1 text-sm text-slate-600">
            {selectedIds.length} anúncio(s) selecionado(s)
          </p>
        </header>

        <div className="min-h-0 flex-1 overflow-y-auto p-5">
          {step === "intent" ? (
            <MutationIntentForm type={type} onSubmit={(intent) => void createPreview(intent)} />
          ) : null}
          {step === "previewing" ? <LoadingState /> : null}
          {step === "preview-shown" && preview ? (
            <div className="space-y-5">
              <PreviewContent preview={preview} />
              {previewError?.kind === "preview_stale" ? (
                <p className="text-sm font-medium text-red-700">Prévia expirada. Gere novamente.</p>
              ) : null}
              <label className="flex items-center gap-2 text-sm font-medium text-slate-700">
                <input
                  type="checkbox"
                  checked={confirmed}
                  disabled={previewError?.kind === "preview_stale"}
                  onChange={(event) => setConfirmed(event.target.checked)}
                />
                Confirmo que revisei a prévia
              </label>
            </div>
          ) : null}
          {step === "approving" ? <LoadingState /> : null}
          {(step === "applying" || step === "terminal") && draft ? (
            <ApplicationResult
              protocolId={draft.protocol_id}
              onTerminal={() => setStep("terminal")}
            />
          ) : null}
          {step === "error" && previewError?.kind === "selection_too_large" ? (
            <ErrorState
              detail={`${selectedIds.length} anúncios selecionados; reduza a seleção ou refine o filtro para gerar uma nova prévia.`}
              onRetry={resetToIntent}
            />
          ) : null}
          {step === "error" && previewError?.kind !== "selection_too_large" ? (
            <ErrorState
              detail={failureCopy(displayCode(previewError))}
              onRetry={() => void retryAfterError()}
            />
          ) : null}
        </div>

        <footer className="flex flex-wrap justify-end gap-2 border-t border-slate-200 px-5 py-4">
          {step === "intent" || step === "preview-shown" ? (
            <button
              type="button"
              onClick={() => void cancel()}
              disabled={cancelOperation.isPending}
              className="rounded-lg border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 disabled:opacity-50"
            >
              Cancelar
            </button>
          ) : null}
          {step === "intent" ? (
            <button
              type="submit"
              form="mutation-intent-form"
              className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
            >
              Gerar prévia
            </button>
          ) : null}
          {step === "preview-shown" ? (
            previewError?.kind === "preview_stale" ? (
              <button
                type="button"
                onClick={() => void retryAfterError()}
                disabled={previewOperation.isPending}
                className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-50"
              >
                Gerar prévia novamente
              </button>
            ) : (
              <button
                type="button"
                onClick={() => void approve()}
                disabled={!confirmed || approveOperation.isPending}
                className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white disabled:cursor-not-allowed disabled:opacity-50"
              >
                Confirmar e aplicar
              </button>
            )
          ) : null}
        </footer>
      </section>
    </div>
  );
}

function ApplicationResult({
  protocolId,
  onTerminal,
}: {
  protocolId: string;
  onTerminal: () => void;
}) {
  const protocol = useMutationProtocol(protocolId);
  const terminal = protocol.data ? isMutationTerminal(protocol.data.state) : false;

  useEffect(() => {
    if (terminal) onTerminal();
  }, [onTerminal, terminal]);

  if (protocol.isError)
    return (
      <ErrorState
        onRetry={() => void protocol.refetch()}
        detail="Não foi possível acompanhar a aplicação."
      />
    );
  if (!protocol.data) return <LoadingState />;
  if (terminal) return <MutationResultSummary protocol={protocol.data} />;
  return <LoadingState />;
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
  return (
    <span className="rounded-full bg-slate-100 px-3 py-1 text-sm font-medium text-slate-700">
      {presentMutationValue(value)} {label}
    </span>
  );
}
