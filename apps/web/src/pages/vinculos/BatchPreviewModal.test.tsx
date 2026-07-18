import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { BatchPreviewModal } from "./BatchPreviewModal";

const previewProductLinkBatch = vi.fn();
const applyProductLinkBatch = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    previewProductLinkBatch: (...args: unknown[]) => previewProductLinkBatch(...args),
    applyProductLinkBatch: (...args: unknown[]) => applyProductLinkBatch(...args),
  }),
}));

function renderModal() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const onClose = vi.fn();
  const onApplied = vi.fn();
  const utils = render(
    <QueryClientProvider client={queryClient}>
      <BatchPreviewModal open candidateIds={["cand_1", "cand_2"]} onClose={onClose} onApplied={onApplied} />
    </QueryClientProvider>,
  );
  return { ...utils, onClose, onApplied };
}

beforeEach(() => {
  previewProductLinkBatch.mockReset();
  applyProductLinkBatch.mockReset();
});

describe("BatchPreviewModal", () => {
  it("fires only the dry-run preview on open — apply is NOT called", async () => {
    previewProductLinkBatch.mockResolvedValue({
      items: [
        { candidate_id: "cand_1", status: "OK" },
        { candidate_id: "cand_2", status: "FAILED", cause: "ALREADY_RESOLVED" },
      ],
    });

    renderModal();

    await waitFor(() => {
      expect(previewProductLinkBatch).toHaveBeenCalledWith({
        approvals: [{ candidate_id: "cand_1" }, { candidate_id: "cand_2" }],
      });
    });

    // The dry-run items render...
    expect(await screen.findByTestId("batch-preview-items")).toBeInTheDocument();
    expect(screen.getByText("cand_1")).toBeInTheDocument();
    expect(screen.getByText("ALREADY_RESOLVED")).toBeInTheDocument();

    // ...but nothing was applied — this is a pure dry-run.
    expect(applyProductLinkBatch).not.toHaveBeenCalled();
  });

  it('applies only the valid ("OK") subset when "prosseguir só com válidos" is clicked', async () => {
    previewProductLinkBatch.mockResolvedValue({
      items: [
        { candidate_id: "cand_1", status: "OK" },
        { candidate_id: "cand_2", status: "FAILED", cause: "ALREADY_RESOLVED" },
      ],
    });
    applyProductLinkBatch.mockResolvedValue({
      batch_id: "batch_1",
      applied: [{ candidate_id: "cand_1" }],
      failed: [],
    });

    const { onApplied } = renderModal();

    const applyButton = await screen.findByRole("button", { name: "Prosseguir só com válidos" });
    fireEvent.click(applyButton);

    await waitFor(() => {
      expect(applyProductLinkBatch).toHaveBeenCalledWith(
        expect.objectContaining({ approvals: [{ candidate_id: "cand_1" }] }),
      );
    });

    await waitFor(() => {
      expect(onApplied).toHaveBeenCalledWith(
        expect.objectContaining({ applied: [{ candidate_id: "cand_1" }] }),
      );
    });
  });
});
