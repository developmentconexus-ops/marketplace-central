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

  // A fonte é o que explica por que um campo do espelho está "—": a exportação
  // Sankhya carrega custo e estoque, o catálogo do cliente não.
  it("names the source of each import so two similar protocols are distinguishable", async () => {
    listErpImports.mockResolvedValue({
      items: [
        summary({ import_id: "imp_1", protocol: "#001-E", source: "xlsx" }),
        summary({ import_id: "imp_2", protocol: "#002-E", source: "catalogo_cliente" }),
      ],
    });

    renderSection();

    expect(await screen.findByText("Planilha Sankhya")).toBeInTheDocument();
    expect(screen.getByText("Catálogo do cliente")).toBeInTheDocument();
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
      "EMPTY_CODPROD — CODPROD ausente",
    );
    expect(screen.getByTestId("erp-import-rejected-rows")).toHaveTextContent("1 linha: 7");
  });

  it("summarises a repeated issue by code instead of listing every row", async () => {
    listErpImports.mockResolvedValue({
      items: [summary({ import_id: "imp_3", protocol: "#003-E", warning_count: 12 })],
    });
    getErpImport.mockResolvedValue(
      detail({
        import_id: "imp_3",
        protocol: "#003-E",
        warnings: Array.from({ length: 12 }, (_unused, index) => ({
          row: index + 1,
          code: "MISSING_CUSTO",
          detail: "custo is unknown",
          offending_value: "",
        })),
      }),
    );

    renderSection();

    fireEvent.click(await screen.findByRole("button", { name: "Ver detalhes" }));

    const warnings = await screen.findByTestId("erp-import-warnings");
    // The total is exact; only the row sample is bounded.
    expect(warnings).toHaveTextContent("12 linhas: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 e mais 2");
    expect(warnings.querySelectorAll("li")).toHaveLength(1);
  });

  it("renders an honest empty state when there is no import history", async () => {
    listErpImports.mockResolvedValue({ items: [] });

    renderSection();

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });
});
