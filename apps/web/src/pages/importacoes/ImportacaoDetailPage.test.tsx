import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ImportacaoDetailPage } from "./ImportacaoDetailPage";

const getErpImportChain = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getErpImportChain: (...args: unknown[]) => getErpImportChain(...args),
  }),
}));

beforeEach(() => {
  getErpImportChain.mockReset();
});

describe("ImportacaoDetailPage", () => {
  it("passes the route import id to the chain panel", async () => {
    getErpImportChain.mockResolvedValue({
      protocol: "#001-E",
      importados: 1,
      vinculados: 1,
      enfileirados: 0,
      queue_read_at: "2026-07-18T12:00:00Z",
    });
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <MemoryRouter initialEntries={["/importacoes/imp_1"]}>
        <QueryClientProvider client={queryClient}>
          <Routes>
            <Route path="/importacoes/:importId" element={<ImportacaoDetailPage />} />
          </Routes>
        </QueryClientProvider>
      </MemoryRouter>,
    );

    expect(await screen.findByTestId("erp-import-chain")).toBeInTheDocument();
    expect(getErpImportChain).toHaveBeenCalledWith("imp_1");
  });
});
