import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { MutationPreview, MutationProtocol } from "@marketplace-central/sdk-runtime";
import { ApplyPriceAction } from "./ApplyPriceAction";

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
  created_at: "2026-07-18T12:00:00Z",
  previewed_at: null,
  approved_at: null,
  finished_at: null,
};

const preview: MutationPreview = {
  ...draft,
  state: "previewed",
  previewed_at: "2026-07-18T12:00:01Z",
  items: [
    {
      seq: 1,
      item_id: "it_1",
      listing_id: "MLB3758134295",
      idempotency_key: "k1",
      before: { price: "89.00" },
      after: { price: "95.00" },
      state: "previewed",
      failure: null,
      applied_at: null,
    },
  ],
};

const createMutation = vi.fn((_input: unknown) => Promise.resolve(draft));
const previewMutation = vi.fn(() => Promise.resolve(preview));
// The write-path executor must never be reachable from this surface.
const approveMutation = vi.fn(() => Promise.reject(new Error("approve must not be called")));

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({ createMutation, previewMutation, approveMutation }),
}));

function renderAction(props: Partial<Parameters<typeof ApplyPriceAction>[0]> = {}) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <ApplyPriceAction
          installationId="inst_test"
          listingId="MLB3758134295"
          newPrice="95.00"
          {...props}
        />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ApplyPriceAction", () => {
  beforeEach(() => {
    createMutation.mockClear();
    previewMutation.mockClear();
    approveMutation.mockClear();
  });

  it("disables apply and explains when no linked listing is available", () => {
    renderAction({ listingId: null });
    const button = screen.getByRole("button", { name: /aplicar preço/i });
    expect(button).toBeDisabled();
    expect(screen.getByTestId("apply-hint")).toBeInTheDocument();
  });

  it("disables apply when there is no price to apply", () => {
    renderAction({ newPrice: "" });
    expect(screen.getByRole("button", { name: /aplicar preço/i })).toBeDisabled();
  });

  it("creates a price_update protocol and caps it at previewed — never approves", async () => {
    renderAction();
    fireEvent.click(screen.getByRole("button", { name: /aplicar preço/i }));

    await waitFor(() => expect(createMutation).toHaveBeenCalledTimes(1));
    expect(createMutation.mock.calls[0][0]).toMatchObject({
      installation_id: "inst_test",
      type: "price_update",
      intent: { new_price: { amount: "95.00", currency: "BRL" } },
      selection: { mode: "explicit", listing_ids: ["MLB3758134295"] },
    });
    // Cap at previewed: preview is reached, approve is never invoked.
    await waitFor(() => expect(previewMutation).toHaveBeenCalledWith("MP-000042"));
    expect(approveMutation).not.toHaveBeenCalled();
  });

  it("shows the previewed result with a protocol link and the no-write reassurance", async () => {
    renderAction();
    fireEvent.click(screen.getByRole("button", { name: /aplicar preço/i }));

    const result = await screen.findByTestId("apply-result");
    expect(within(result).getByTestId("apply-state")).toHaveTextContent("previewed");
    const link = within(result).getByRole("link", { name: /MP-000042/ });
    expect(link).toHaveAttribute("href", "/protocolos/MP-000042");
    expect(result).toHaveTextContent(/nada foi enviado ao Mercado Livre/i);
  });

  it("exposes no approve/apply-to-ML control anywhere on the surface", async () => {
    renderAction();
    fireEvent.click(screen.getByRole("button", { name: /aplicar preço/i }));
    await screen.findByTestId("apply-result");
    expect(screen.queryByRole("button", { name: /aprovar|approve|enviar ao mercado/i })).toBeNull();
  });

  it("surfaces an error when the protocol cannot be created", async () => {
    createMutation.mockImplementationOnce(() => Promise.reject(new Error("boom")));
    renderAction();
    fireEvent.click(screen.getByRole("button", { name: /aplicar preço/i }));
    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(approveMutation).not.toHaveBeenCalled();
  });
});
