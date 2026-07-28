import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ErpImportChain } from "@marketplace-central/sdk-runtime";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ImportChainPanel } from "./ImportChainPanel";

const getErpImportChain = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getErpImportChain: (...args: unknown[]) => getErpImportChain(...args),
  }),
}));

function renderPanel(importId = "imp_1") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ImportChainPanel importId={importId} />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  getErpImportChain.mockReset();
});

describe("ImportChainPanel", () => {
  it("renders a known protocol verbatim", async () => {
    getErpImportChain.mockResolvedValue({
      protocol: "#001-E",
      importados: 1,
      vinculados: 1,
      enfileirados: 1,
      queue_read_at: "2026-07-18T12:00:00Z",
    });

    renderPanel();

    expect(await screen.findByText("#001-E", { exact: true })).toBeInTheDocument();
  });

  it("consumes the chain returned by the server", async () => {
    getErpImportChain.mockResolvedValue({
      protocol: "#137-E",
      importados: 137,
      vinculados: 42,
      enfileirados: 9,
      queue_read_at: "2026-07-18T12:00:00Z",
    });

    renderPanel("imp_real");

    expect(await screen.findByTestId("erp-import-chain-importados")).toHaveTextContent("137");
    expect(screen.getByTestId("erp-import-chain-vinculados")).toHaveTextContent("42");
    expect(screen.getByTestId("erp-import-chain-enfileirados")).toHaveTextContent("9");
    expect(getErpImportChain).toHaveBeenCalledWith("imp_real");
  });

  it("renders a missing counter as unknown without blanking the card", async () => {
    const chain = {
      protocol: "#138-E",
      importados: 138,
      enfileirados: 8,
      queue_read_at: "2026-07-18T12:00:00Z",
    } as Partial<ErpImportChain> as ErpImportChain;
    getErpImportChain.mockResolvedValue(chain);

    renderPanel();

    expect(await screen.findByTestId("erp-import-chain-vinculados")).toHaveTextContent("—");
    expect(screen.getByTestId("erp-import-chain-vinculados")).not.toHaveTextContent("0");
    expect(screen.getByTestId("erp-import-chain-importados")).toHaveTextContent("138");
    expect(screen.getByTestId("erp-import-chain-enfileirados")).toHaveTextContent("8");
  });

  it("renders a null counter as unknown without blanking the card", async () => {
    const chain = {
      protocol: "#139-E",
      importados: 139,
      vinculados: 41,
      enfileirados: null,
      queue_read_at: "2026-07-18T12:00:00Z",
    } as unknown as ErpImportChain;
    getErpImportChain.mockResolvedValue(chain);

    renderPanel();

    expect(await screen.findByTestId("erp-import-chain-enfileirados")).toHaveTextContent("—");
    expect(screen.getByTestId("erp-import-chain-enfileirados")).not.toHaveTextContent("0");
    expect(screen.getByTestId("erp-import-chain-importados")).toHaveTextContent("139");
    expect(screen.getByTestId("erp-import-chain-vinculados")).toHaveTextContent("41");
  });

  it("renders an absent protocol as unknown without blanking the counters", async () => {
    const chain = {
      importados: 140,
      vinculados: 40,
      enfileirados: 7,
      queue_read_at: "2026-07-18T12:00:00Z",
    } as Partial<ErpImportChain> as ErpImportChain;
    getErpImportChain.mockResolvedValue(chain);

    renderPanel();

    expect(await screen.findByText("—", { exact: true })).toBeInTheDocument();
    expect(screen.getByTestId("erp-import-chain-importados")).toHaveTextContent("140");
    expect(screen.getByTestId("erp-import-chain-vinculados")).toHaveTextContent("40");
    expect(screen.getByTestId("erp-import-chain-enfileirados")).toHaveTextContent("7");
  });

  it("renders an empty or whitespace protocol as unknown", async () => {
    getErpImportChain.mockResolvedValue({
      protocol: "   ",
      importados: 1,
      vinculados: 1,
      enfileirados: 1,
      queue_read_at: "2026-07-18T12:00:00Z",
    });

    renderPanel();

    expect(await screen.findByText("—", { exact: true })).toBeInTheDocument();
  });

  it("renders a not-found error without a chain", async () => {
    getErpImportChain.mockRejectedValue({ status: 404, error: "import_not_found" });

    renderPanel();

    expect(await screen.findByTestId("erp-import-chain-error")).toHaveTextContent("não encontrada");
    expect(screen.queryByTestId("erp-import-chain")).toBeNull();
  });

  it("renders a generic error for server failures without a chain", async () => {
    getErpImportChain.mockRejectedValue({ status: 500, error: "internal_error" });

    renderPanel();

    // Um 5xx é transitório, então o hook faz UMA nova tentativa (o `retry: false`
    // do QueryClient do teste não vale: a opção do próprio hook prevalece). O
    // estado de erro só assenta depois do retryDelay de ~1s.
    const error = await screen.findByTestId("erp-import-chain-error", {}, { timeout: 5000 });
    expect(error).toHaveTextContent("Não foi possível carregar a cadeia da importação.");
    expect(error).not.toHaveTextContent("não encontrada");
    expect(screen.queryByTestId("erp-import-chain")).toBeNull();
  });
});
