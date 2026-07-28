import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ErpImportSummary } from "@marketplace-central/sdk-runtime";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ImportacoesPage } from "./ImportacoesPage";

const listErpImports = vi.fn();
const getErpImport = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listErpImports: (...args: unknown[]) => listErpImports(...args),
    getErpImport: (...args: unknown[]) => getErpImport(...args),
  }),
}));

function summary(overrides: Partial<ErpImportSummary>): ErpImportSummary {
  return {
    import_id: "imp_1",
    protocol: "#001-E",
    file_sha256: "abc",
    source: "xlsx",
    imported_at: "2026-07-18T12:00:00Z",
    status: "COMPLETED",
    accepted_count: 10,
    rejected_count: 0,
    warning_count: 0,
    ...overrides,
  };
}

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ImportacoesPage />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  listErpImports.mockReset();
  getErpImport.mockReset();
});

describe("ImportacoesPage", () => {
  it("renders the import history below the page heading", async () => {
    listErpImports.mockResolvedValue({
      items: [summary({ protocol: "#009-E" })],
    });

    renderPage();

    expect(screen.getByRole("heading", { name: "Importações" })).toBeInTheDocument();
    expect(await screen.findByText("#009-E")).toBeInTheDocument();
  });
});
