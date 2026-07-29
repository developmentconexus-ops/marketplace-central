import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type {
  ProductLinkAuditEntry,
  ProductLinkWorkflowItem,
} from "@marketplace-central/sdk-runtime";
import { render, screen, within } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ResolvidosTab } from "./ResolvidosTab";

const listProductLinkWorkflows = vi.fn();
const undoProductLinkResolution = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listProductLinkWorkflows: (...args: unknown[]) => listProductLinkWorkflows(...args),
    undoProductLinkResolution: (...args: unknown[]) => undoProductLinkResolution(...args),
  }),
}));

function auditEntry(overrides: Partial<ProductLinkAuditEntry>): ProductLinkAuditEntry {
  return {
    audit_id: "aud_1",
    installation_id: "inst_1",
    provider_code: "mercado_livre",
    provider_item_id: "MLB1",
    action: "approve_candidate",
    actor: { actor_type: "operator", actor_id: "op_1" },
    previous_state: "unresolved",
    next_state: "resolved",
    created_at: "2026-07-20T12:00:00Z",
    ...overrides,
  } as ProductLinkAuditEntry;
}

function workflow(
  providerItemId: string,
  audit: ProductLinkAuditEntry[],
): ProductLinkWorkflowItem {
  return {
    identity: { installation_id: "inst_1", provider_item_id: providerItemId },
    current_link: {
      installation_id: "inst_1",
      provider_code: "mercado_livre",
      provider_item_id: providerItemId,
      state: "resolved",
      internal_product_id: 1001,
      internal_product_name: "Produto Y",
      created_at: "2026-07-20T12:00:00Z",
      updated_at: "2026-07-20T12:00:00Z",
    },
    candidates: [],
    audit,
  };
}

function renderTab() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ResolvidosTab installationId="inst_1" />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  listProductLinkWorkflows.mockReset();
  undoProductLinkResolution.mockReset();
});

/**
 * F-04 — the auto-approved badge.
 *
 * The predicate is `actor.actor_type === "system"` on the audit entry that
 * RESOLVED the link, not the `rule_matched = exact_ean_unique AND actor =
 * system` pair the M-06 brief names: `rule_matched` is absent from the wire, and
 * that exact pair is forbidden by the CHECK at 0082_product_link_decisions.sql:54.
 * Full reasoning lives on `isAutoResolved` in useVinculosResolved.ts.
 */
describe("ResolvidosTab — badge de auto-aprovado", () => {
  it("shows the badge for a link resolved by the system", async () => {
    listProductLinkWorkflows.mockResolvedValue({
      items: [
        workflow("MLB_AUTO", [
          auditEntry({
            audit_id: "aud_auto",
            provider_item_id: "MLB_AUTO",
            actor: { actor_type: "system", actor_id: "auto_linker" },
          }),
        ]),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("resolvido-row");
    expect(within(row).getByTestId("auto-resolved-badge")).toHaveTextContent("Automático");
    // The state chip is not replaced by the badge — both facts stay on screen.
    expect(within(row).getByText("Vinculado ✓")).toBeInTheDocument();
  });

  it("shows NO badge for a link resolved by an operator", async () => {
    listProductLinkWorkflows.mockResolvedValue({
      items: [
        workflow("MLB_MANUAL", [
          auditEntry({
            audit_id: "aud_manual",
            provider_item_id: "MLB_MANUAL",
            actor: { actor_type: "operator", actor_id: "op_1", actor_name: "Operador" },
          }),
        ]),
      ],
    });

    renderTab();

    const row = await screen.findByTestId("resolvido-row");
    expect(within(row).queryByTestId("auto-resolved-badge")).not.toBeInTheDocument();
    expect(within(row).getByText("Vinculado ✓")).toBeInTheDocument();
  });

  it("shows NO badge, and does not crash or print undefined, for a pre-M-05 link with no resolving audit entry", async () => {
    // milestone.md:213's pre-M-05 record: the link is resolved but nothing in the
    // audit trail moved it there. Absence of the badge reads "not automatic",
    // which is the truth available; a fabricated badge would not be (ADR-17).
    listProductLinkWorkflows.mockResolvedValue({
      items: [workflow("MLB_LEGACY", [])],
    });

    renderTab();

    const row = await screen.findByTestId("resolvido-row");
    expect(within(row).queryByTestId("auto-resolved-badge")).not.toBeInTheDocument();
    expect(row.textContent).not.toContain("undefined");
    // Desfazer is honestly disabled — there is no resolution to reverse.
    expect(within(row).getByRole("button", { name: "Desfazer" })).toBeDisabled();
  });

  it("keeps the structural column label neutral of provider (F-05)", async () => {
    listProductLinkWorkflows.mockResolvedValue({ items: [workflow("MLB1", [])] });

    renderTab();

    await screen.findByTestId("resolvido-row");
    expect(screen.getByRole("columnheader", { name: "Anúncio" })).toBeInTheDocument();
    expect(screen.queryByRole("columnheader", { name: "Anúncio ML" })).not.toBeInTheDocument();
  });
});
