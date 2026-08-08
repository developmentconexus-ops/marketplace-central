import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import {
  MarketplaceCentralClientError,
  type MutationPreview,
  type MutationProtocol,
  type MutationType,
} from "@marketplace-central/sdk-runtime";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { MutationPreviewModal } from "./MutationPreviewModal";

const createMutation = vi.fn();
const previewMutation = vi.fn();
const cancelMutation = vi.fn();
const approveMutation = vi.fn();
const getMutation = vi.fn();
const listMutationItems = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    createMutation,
    previewMutation,
    cancelMutation,
    approveMutation,
    getMutation,
    listMutationItems,
  }),
}));

const draft: MutationProtocol = {
  protocol_id: "MP-000042",
  installation_id: "inst_test",
  type: "price_update",
  state: "draft",
  actor: "operator_supplied_unverified",
  intent: {},
  selection: {},
  totals: {},
  source_as_of: null,
  retried_from: null,
  created_at: "2026-07-17T12:00:00Z",
  previewed_at: null,
  approved_at: null,
  finished_at: null,
};

const preview: MutationPreview = {
  ...draft,
  state: "previewed",
  totals: { items: 2, previewed: 2, failed: 0 },
  previewed_at: "2026-07-17T12:00:01Z",
  items: [
    {
      seq: 1,
      item_id: "MP-000042:listing_1",
      listing_id: "listing_1",
      idempotency_key: "MP-000042:listing_1",
      before: { price: { amount: "40.00", currency: "BRL" } },
      after: { price: { amount: "49.90", currency: "BRL" } },
      state: "previewed",
      failure: null,
      applied_at: null,
    },
    {
      seq: 2,
      item_id: "MP-000042:listing_2",
      listing_id: "listing_2",
      idempotency_key: "MP-000042:listing_2",
      before: null,
      after: { status: "paused" },
      state: "previewed",
      failure: null,
      applied_at: null,
    },
  ],
};

function renderModal(type: MutationType, onClose = vi.fn()) {
  const queryClient = new QueryClient({ defaultOptions: { mutations: { retry: false } } });
  render(
    <QueryClientProvider client={queryClient}>
      <MutationPreviewModal
        open
        type={type}
        installationId="inst_test"
        selectedIds={["listing_1", "listing_2"]}
        onClose={onClose}
      />
    </QueryClientProvider>,
  );
  return { onClose };
}

function submit() {
  fireEvent.click(screen.getByRole("button", { name: "Gerar prévia" }));
}

describe("MutationPreviewModal", () => {
  beforeEach(() => {
    createMutation.mockReset();
    previewMutation.mockReset();
    cancelMutation.mockReset();
    approveMutation.mockReset();
    getMutation.mockReset();
    listMutationItems.mockReset();
    createMutation.mockResolvedValue(draft);
    previewMutation.mockResolvedValue(preview);
    cancelMutation.mockResolvedValue({ ...draft, state: "cancelled" });
    approveMutation.mockResolvedValue({
      ...draft,
      state: "applying",
      approved_at: "2026-07-17T12:00:02Z",
    });
    getMutation.mockResolvedValue({
      ...draft,
      state: "applying",
      approved_at: "2026-07-17T12:00:02Z",
    });
    listMutationItems.mockResolvedValue({ items: [], next_cursor: null, page_size: 50 });
  });

  it("mantém a confirmação indisponível até previewMutation resolver", async () => {
    let resolvePreview!: (value: MutationPreview) => void;
    previewMutation.mockReturnValue(
      new Promise((resolve) => {
        resolvePreview = resolve;
      }),
    );
    renderModal("listing_pause");

    submit();

    expect(await screen.findByText("Carregando…")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Confirmar e aplicar" })).not.toBeInTheDocument();
    resolvePreview(preview);
    expect(await screen.findByRole("button", { name: "Confirmar e aplicar" })).toBeDisabled();
  });

  const intents: Array<{
    type: MutationType;
    fill: () => void;
    intent: Record<string, unknown>;
  }> = [
    {
      type: "price_update",
      fill: () =>
        fireEvent.change(screen.getByLabelText("Novo preço"), { target: { value: "49.90" } }),
      intent: { new_price: { amount: "49.90", currency: "BRL" } },
    },
    {
      type: "stock_correct",
      fill: () =>
        fireEvent.change(screen.getByLabelText("Quantidade a publicar"), {
          target: { value: "7" },
        }),
      intent: { publish_quantity: 7 },
    },
    { type: "listing_pause", fill: () => undefined, intent: {} },
    { type: "listing_resync", fill: () => undefined, intent: {} },
    {
      type: "link_apply",
      fill: () => {
        fireEvent.change(screen.getByLabelText("Ação de vínculo"), {
          target: { value: "manual_resolve" },
        });
        fireEvent.change(screen.getByLabelText("ID do produto"), { target: { value: "PROD-10" } });
      },
      intent: { action: "manual_resolve", product_id: "PROD-10" },
    },
    {
      type: "listing_edit",
      fill: () => {
        fireEvent.change(screen.getByLabelText("ID do atributo"), { target: { value: "COLOR" } });
        fireEvent.change(screen.getByLabelText("Valor do atributo"), { target: { value: "Azul" } });
      },
      intent: { attributes: [{ id: "COLOR", value_name: "Azul" }] },
    },
  ];

  it.each(intents)("envia o intent IC-03 exato para $type", async ({ type, fill, intent }) => {
    renderModal(type);
    fill();
    submit();

    await waitFor(() => expect(previewMutation).toHaveBeenCalledWith("MP-000042"));
    expect(createMutation).toHaveBeenCalledWith({
      installation_id: "inst_test",
      type,
      actor: "operator_supplied_unverified",
      intent,
      selection: { mode: "explicit", listing_ids: ["listing_1", "listing_2"] },
    });
  });

  it("renderiza totais e linhas antes→depois retornados pelo servidor", async () => {
    renderModal("listing_pause");
    submit();

    expect(await screen.findByText("2 itens")).toBeInTheDocument();
    expect(screen.getByText("2 previstos")).toBeInTheDocument();
    expect(screen.getByText("0 falhas")).toBeInTheDocument();
    expect(screen.getByText(/40\.00/)).toBeInTheDocument();
    expect(screen.getByText(/49\.90/)).toBeInTheDocument();
    expect(screen.getByTitle("valor anterior não informado")).toBeInTheDocument();
    expect(screen.getByText(/paused/)).toBeInTheDocument();
  });

  it("fecha localmente ao cancelar antes da criação do rascunho", () => {
    const { onClose } = renderModal("listing_pause");
    fireEvent.click(screen.getByRole("button", { name: "Cancelar" }));

    expect(onClose).toHaveBeenCalledOnce();
    expect(cancelMutation).not.toHaveBeenCalled();
  });

  it("cancela uma vez o rascunho existente antes de fechar", async () => {
    const { onClose } = renderModal("listing_pause");
    submit();
    await screen.findByRole("button", { name: "Confirmar e aplicar" });
    fireEvent.click(screen.getByRole("button", { name: "Cancelar" }));

    await waitFor(() => expect(cancelMutation).toHaveBeenCalledTimes(1));
    expect(cancelMutation).toHaveBeenCalledWith("MP-000042");
    await waitFor(() => expect(onClose).toHaveBeenCalledOnce());
  });

  it("cancela e descarta o rascunho quando a seleção excede o limite", async () => {
    previewMutation.mockRejectedValue(
      new MarketplaceCentralClientError(422, "selection_too_large", "too many", {}),
    );
    renderModal("listing_pause");
    submit();

    expect(await screen.findByRole("alert")).toHaveTextContent("2 anúncios selecionados");
    expect(screen.getByRole("alert")).toHaveTextContent("reduza a seleção ou refine o filtro");
    await waitFor(() => expect(cancelMutation).toHaveBeenCalledTimes(1));
    expect(screen.queryByRole("button", { name: "Confirmar e aplicar" })).not.toBeInTheDocument();
  });

  it("exige confirmação explícita antes de aprovar", async () => {
    renderModal("listing_pause");
    submit();

    const approve = await screen.findByRole("button", { name: "Confirmar e aplicar" });
    expect(approve).toBeDisabled();
    fireEvent.click(screen.getByRole("checkbox", { name: "Confirmo que revisei a prévia" }));
    expect(approve).toBeEnabled();
  });

  it("aprova uma vez diante de dois cliques imediatos", async () => {
    let resolveApprove!: (value: MutationProtocol) => void;
    approveMutation.mockReturnValue(
      new Promise((resolve) => {
        resolveApprove = resolve;
      }),
    );
    renderModal("listing_pause");
    submit();
    await screen.findByRole("button", { name: "Confirmar e aplicar" });
    fireEvent.click(screen.getByRole("checkbox", { name: "Confirmo que revisei a prévia" }));
    const approve = screen.getByRole("button", { name: "Confirmar e aplicar" });

    fireEvent.click(approve);
    fireEvent.click(approve);

    await waitFor(() => expect(approveMutation).toHaveBeenCalledTimes(1));
    expect(approveMutation).toHaveBeenCalledWith("MP-000042", { execute: true });
    resolveApprove({ ...draft, state: "applying" });
  });

  it("volta à prévia sem reaprová-la quando ela expira", async () => {
    approveMutation.mockRejectedValue(
      new MarketplaceCentralClientError(409, "preview_stale", "stale", {}),
    );
    renderModal("listing_pause");
    submit();
    await screen.findByRole("button", { name: "Confirmar e aplicar" });
    fireEvent.click(screen.getByRole("checkbox", { name: "Confirmo que revisei a prévia" }));
    fireEvent.click(screen.getByRole("button", { name: "Confirmar e aplicar" }));

    expect(await screen.findByText("Prévia expirada. Gere novamente.")).toBeInTheDocument();
    expect(
      screen.getByRole("checkbox", { name: "Confirmo que revisei a prévia" }),
    ).not.toBeChecked();
    expect(screen.queryByRole("button", { name: "Confirmar e aplicar" })).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Gerar prévia novamente" })).toBeEnabled();
    expect(approveMutation).toHaveBeenCalledTimes(1);
  });

  it("mostra o resultado terminal com cópia segura e link do protocolo", async () => {
    const terminal = {
      ...draft,
      state: "partially_failed" as const,
      finished_at: "2026-07-17T12:00:03Z",
    };
    approveMutation.mockResolvedValue(terminal);
    getMutation.mockResolvedValue(terminal);
    listMutationItems.mockResolvedValue({
      items: [
        {
          ...preview.items[0],
          state: "failed",
          failure: { code: "provider_validation", message_pt: "texto cru proibido" },
        },
      ],
      next_cursor: null,
      page_size: 50,
    });
    renderModal("listing_pause");
    submit();
    await screen.findByRole("button", { name: "Confirmar e aplicar" });
    fireEvent.click(screen.getByRole("checkbox", { name: "Confirmo que revisei a prévia" }));
    fireEvent.click(screen.getByRole("button", { name: "Confirmar e aplicar" }));

    expect(await screen.findByText("Rejeitado pela validação do marketplace.")).toBeInTheDocument();
    expect(screen.queryByText("texto cru proibido")).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Ver protocolo" })).toHaveAttribute(
      "href",
      "/protocolos/MP-000042",
    );
  });
});
