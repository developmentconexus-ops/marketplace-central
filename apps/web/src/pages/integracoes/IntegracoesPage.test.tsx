import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { getActiveErpSource, setActiveErpSource } from "@marketplace-central/web-query";
import { IntegracoesPage } from "./IntegracoesPage";

const createErpImport = vi.fn();
const getErpImport = vi.fn();
const listErpImports = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    createErpImport: (...args: unknown[]) => createErpImport(...args),
    getErpImport: (...args: unknown[]) => getErpImport(...args),
    listErpImports: (...args: unknown[]) => listErpImports(...args),
  }),
}));

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <IntegracoesPage />
    </QueryClientProvider>,
  );
}

function selectFile(name = "catalogo.xlsx") {
  const input = screen.getByTestId("erp-import-file-input") as HTMLInputElement;
  const file = new File(["binary"], name, {
    type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  });
  fireEvent.change(input, { target: { files: [file] } });
  return file;
}

const detailFixture = {
  import_id: "imp_1",
  protocol: "#004-E",
  file_sha256: "abc",
  source: "catalogo_cliente" as const,
  imported_at: "2026-07-20T10:00:00Z",
  status: "COMPLETED" as const,
  accepted_count: 2012,
  rejected_count: 0,
  warning_count: 4024,
  rejected_rows: [],
  warnings: [{ row: 1, code: "MISSING_CUSTO", detail: "custo is unknown", offending_value: null }],
};

describe("IntegracoesPage", () => {
  beforeEach(() => {
    createErpImport.mockReset();
    getErpImport.mockReset();
    listErpImports.mockReset();
    listErpImports.mockResolvedValue({ items: [] });
    setActiveErpSource("xlsx");
  });

  it("renders the plataforma-config heading", () => {
    renderPage();
    expect(screen.getByRole("heading", { name: "Configuração da plataforma" })).toBeInTheDocument();
  });

  it("posts the upload with the default catálogo-do-cliente source and shows the result summary", async () => {
    createErpImport.mockResolvedValue({ import_id: "imp_1", protocol: "#004-E", status: "COMPLETED" });
    getErpImport.mockResolvedValue(detailFixture);
    renderPage();

    const file = selectFile();
    expect(screen.getByTestId("erp-import-file-name")).toHaveTextContent("catalogo.xlsx");
    fireEvent.click(screen.getByTestId("erp-import-submit"));

    await waitFor(() => expect(createErpImport).toHaveBeenCalledTimes(1));
    expect(createErpImport).toHaveBeenCalledWith(file, "catalogo_cliente", "catalogo.xlsx");

    const result = await screen.findByTestId("erp-import-result");
    expect(within(result).getByTestId("result-accepted")).toHaveTextContent("2012");
    expect(within(result).getByTestId("result-warnings")).toHaveTextContent("4024");
    expect(getErpImport).toHaveBeenCalledWith("imp_1");
  });

  it("posts the strict xlsx source when Sankhya is selected", async () => {
    createErpImport.mockResolvedValue({ import_id: "imp_2", protocol: "#005-E", status: "COMPLETED" });
    getErpImport.mockResolvedValue({ ...detailFixture, import_id: "imp_2" });
    renderPage();

    fireEvent.click(screen.getByLabelText(/Sankhya \(custo/));
    const file = selectFile("sankhya.xlsx");
    fireEvent.click(screen.getByTestId("erp-import-submit"));

    await waitFor(() => expect(createErpImport).toHaveBeenCalledTimes(1));
    expect(createErpImport).toHaveBeenCalledWith(file, "xlsx", "sankhya.xlsx");
  });

  it("shows a duplicate message on a 409 duplicate_file and does not surface a result", async () => {
    createErpImport.mockRejectedValue({ status: 409, body: { error: "duplicate_file", protocol: "#003-E" } });
    renderPage();
    selectFile();
    fireEvent.click(screen.getByTestId("erp-import-submit"));

    const error = await screen.findByTestId("erp-import-error");
    expect(error).toHaveTextContent("#003-E");
    expect(screen.queryByTestId("erp-import-result")).not.toBeInTheDocument();
  });

  it("shows the missing-column name on a 422 missing_required_column", async () => {
    createErpImport.mockRejectedValue({ status: 422, body: { error: "missing_required_column", column: "CUSTO" } });
    renderPage();
    selectFile();
    fireEvent.click(screen.getByTestId("erp-import-submit"));

    const error = await screen.findByTestId("erp-import-error");
    expect(error).toHaveTextContent("CUSTO");
  });

  it("rejects a non-xlsx file locally without hitting the network", () => {
    renderPage();
    const input = screen.getByTestId("erp-import-file-input") as HTMLInputElement;
    const file = new File(["x"], "notes.txt", { type: "text/plain" });
    fireEvent.change(input, { target: { files: [file] } });

    expect(screen.getByRole("alert")).toHaveTextContent(/\.xlsx/);
    expect(screen.getByTestId("erp-import-submit")).toBeDisabled();
    expect(createErpImport).not.toHaveBeenCalled();
  });

  it("defaults the Fonte ativa selector to Sankhya and persists a switch to Catálogo do cliente", () => {
    renderPage();
    const xlsx = screen.getByTestId("active-source-xlsx") as HTMLInputElement;
    const catalogo = screen.getByTestId("active-source-catalogo_cliente") as HTMLInputElement;
    expect(xlsx.checked).toBe(true);
    expect(catalogo.checked).toBe(false);

    fireEvent.click(catalogo);

    expect(catalogo.checked).toBe(true);
    expect(xlsx.checked).toBe(false);
    expect(getActiveErpSource()).toBe("catalogo_cliente");
  });

  it("keeps the Mercado Livre connect affordance inert (disabled, no network)", () => {
    renderPage();
    const connect = screen.getByTestId("provider-connect-ml");
    expect(connect).toBeDisabled();
    expect(connect).toHaveTextContent("disponível em breve");
    expect(createErpImport).not.toHaveBeenCalled();
  });
});
