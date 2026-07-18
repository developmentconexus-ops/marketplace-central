import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { ErpImportDetail, ErpImportSummary } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ImportacaoSection } from "./ImportacaoSection";

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

function detail(overrides: Partial<ErpImportDetail>): ErpImportDetail {
  return {
    ...summary({}),
    rejected_rows: [],
    warnings: [],
    ...overrides,
  };
}

function renderSection() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ImportacaoSection />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  listErpImports.mockReset();
  getErpImport.mockReset();
});

describe("ImportacaoSection", () => {
  it("renders import rows with protocol, status, and counts", async () => {
    listErpImports.mockResolvedValue({
      items: [
        summary({
          import_id: "imp_1",
          protocol: "#001-E",
          status: "COMPLETED",
          accepted_count: 42,
          rejected_count: 0,
          warning_count: 2,
        }),
        summary({
          import_id: "imp_2",
          protocol: "#002-E",
          status: "REJECTED",
          accepted_count: 5,
          rejected_count: 3,
          warning_count: 0,
        }),
      ],
    });

    renderSection();

    expect(await screen.findByText("#001-E")).toBeInTheDocument();
    expect(screen.getByText("Concluída")).toBeInTheDocument();
    expect(screen.getByText("#002-E")).toBeInTheDocument();
    expect(screen.getByText("Rejeitada")).toBeInTheDocument();
    expect(screen.getByText("42")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
  });

  it("expands a row to fetch and list rejected rows", async () => {
    listErpImports.mockResolvedValue({
      items: [
        summary({
          import_id: "imp_2",
          protocol: "#002-E",
          status: "REJECTED",
          accepted_count: 5,
          rejected_count: 1,
          warning_count: 0,
        }),
      ],
    });
    getErpImport.mockResolvedValue(
      detail({
        import_id: "imp_2",
        protocol: "#002-E",
        status: "REJECTED",
        rejected_rows: [
          { row: 7, code: "EMPTY_CODPROD", detail: "CODPROD ausente", offending_value: "" },
        ],
        warnings: [],
      }),
    );

    renderSection();

    const expandButton = await screen.findByRole("button", { name: "Ver detalhes" });
    fireEvent.click(expandButton);

    await waitFor(() => {
      expect(getErpImport).toHaveBeenCalledWith("imp_2");
    });

    expect(await screen.findByTestId("erp-import-rejected-rows")).toHaveTextContent(
      "Linha 7 — EMPTY_CODPROD: CODPROD ausente",
    );
  });

  it("renders an honest empty state when there is no import history", async () => {
    listErpImports.mockResolvedValue({ items: [] });

    renderSection();

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });
});
