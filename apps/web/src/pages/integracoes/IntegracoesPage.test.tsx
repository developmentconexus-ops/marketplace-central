import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MarketplaceCentralClientError } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { InstallationProvider } from "../../app/InstallationContext";
import { IntegracoesPage } from "./IntegracoesPage";

// jsdom in this workspace does not instantiate Web Storage under the worker's
// Vitest runner. Keep the storage-shaped fixture here so the test can prove
// that the assortment interaction never writes a browser preference.
if (typeof globalThis.localStorage === "undefined") {
  const stores = [new Map<string, string>(), new Map<string, string>()];
  const define = (name: string, value: unknown) =>
    Object.defineProperty(Storage.prototype, name, { configurable: true, writable: true, value });
  define("getItem", function (this: Storage, key: string) {
    const store = (this as Storage & { __store: Map<string, string> }).__store;
    return store.has(key) ? store.get(key)! : null;
  });
  define("setItem", function (this: Storage, key: string, value: string) {
    (this as Storage & { __store: Map<string, string> }).__store.set(key, String(value));
  });
  define("removeItem", function (this: Storage, key: string) {
    (this as Storage & { __store: Map<string, string> }).__store.delete(key);
  });
  define("clear", function (this: Storage) {
    (this as Storage & { __store: Map<string, string> }).__store.clear();
  });
  define("key", function (this: Storage, index: number) {
    return Array.from((this as Storage & { __store: Map<string, string> }).__store.keys())[index] ?? null;
  });
  const storageEntries: Array<["localStorage" | "sessionStorage", Map<string, string>]> = [
    ["localStorage", stores[0]],
    ["sessionStorage", stores[1]],
  ];
  for (const [name, store] of storageEntries) {
    const storage = Object.create(Storage.prototype) as Storage & { __store: Map<string, string> };
    storage.__store = store;
    Object.defineProperty(storage, "length", { configurable: true, get: () => store.size });
    Object.defineProperty(globalThis, name, { configurable: true, value: storage });
    Object.defineProperty(window, name, { configurable: true, value: storage });
  }
}

const createErpImport = vi.fn();
const getErpImport = vi.fn();
const listErpImports = vi.fn();
const getActiveSource = vi.fn();
const setActiveSource = vi.fn();
const getSellableAssortment = vi.fn();
const setSellableAssortment = vi.fn();
const getCatalogAssortmentCounts = vi.fn();
const listIntegrationInstallations = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    createErpImport: (...args: unknown[]) => createErpImport(...args),
    getErpImport: (...args: unknown[]) => getErpImport(...args),
    listErpImports: (...args: unknown[]) => listErpImports(...args),
    getActiveSource: (...args: unknown[]) => getActiveSource(...args),
    setActiveSource: (...args: unknown[]) => setActiveSource(...args),
    getSellableAssortment: (...args: unknown[]) => getSellableAssortment(...args),
    setSellableAssortment: (...args: unknown[]) => setSellableAssortment(...args),
    getCatalogAssortmentCounts: (...args: unknown[]) => getCatalogAssortmentCounts(...args),
    listIntegrationInstallations: (...args: unknown[]) => listIntegrationInstallations(...args),
  }),
}));

function activeSourceConfig(source: string) {
  return { active_source: source, source_kind: source === "sankhya" ? "live_read_through" : "upload_snapshot", set_at: "2026-07-24T10:00:00Z", set_by: null };
}

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <MemoryRouter>
      <QueryClientProvider client={queryClient}>
        {/* IntegracoesPage now mounts ConnectionHealthCard, which reads
            useInstallation() (InstallationContext.tsx). In the real app
            AppRouter.tsx:60 wraps every route in InstallationProvider; this
            test rendered the page standalone before that dependency existed. */}
        <InstallationProvider>
          <IntegracoesPage />
        </InstallationProvider>
      </QueryClientProvider>
    </MemoryRouter>,
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

function assortmentToggle(label: string) {
  const option = screen.getByText(label).closest("label");
  if (!option) throw new Error(`Assortment option not found: ${label}`);
  return within(option).getByRole("checkbox") as HTMLInputElement;
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
    getActiveSource.mockReset();
    setActiveSource.mockReset();
    getSellableAssortment.mockReset();
    setSellableAssortment.mockReset();
    getCatalogAssortmentCounts.mockReset();
    listIntegrationInstallations.mockReset();
    // Stateful like the server: a successful PUT changes what the next GET
    // returns, so the blanket post-write invalidation re-reads the new source.
    let stored = "xlsx";
    getActiveSource.mockImplementation(() => Promise.resolve(activeSourceConfig(stored)));
    setActiveSource.mockImplementation((req: { active_source: string }) => {
      stored = req.active_source;
      return Promise.resolve(activeSourceConfig(stored));
    });
    let storedAssortment = {
      only_revenda: true,
      only_em_estoque: false,
      only_ecommerce_eligible: false,
    };
    getSellableAssortment.mockImplementation(() => Promise.resolve(storedAssortment));
    setSellableAssortment.mockImplementation((req: typeof storedAssortment) => {
      storedAssortment = { ...req };
      return Promise.resolve(storedAssortment);
    });
    // Stateful counts catch a missing catalog invalidation: without the refetch, the UI keeps showing the pre-toggle count.
    getCatalogAssortmentCounts.mockImplementation(() =>
      Promise.resolve({ sellable_count: storedAssortment.only_revenda ? 2 : 3, total_count: 4 }),
    );
    listIntegrationInstallations.mockResolvedValue({ items: [] });
  });

  it("renders Sortimento vendável, persists all toggles, and shows the exact live count", async () => {
    const setItem = vi.spyOn(Storage.prototype, "setItem");
    window.localStorage.clear();
    window.sessionStorage.clear();
    renderPage();

    expect(await screen.findByRole("heading", { name: "Sortimento vendável" })).toBeInTheDocument();
    expect(await screen.findByText("Resultado: 2 de 4 produtos")).toBeInTheDocument();

    const revenda = assortmentToggle("Somente produtos de revenda");
    const estoque = assortmentToggle("Somente com estoque disponível");
    const ecommerce = assortmentToggle("Somente elegíveis ao e-commerce");
    expect(revenda.checked, "Somente produtos de revenda must bind to only_revenda").toBe(true);
    expect(estoque.checked, "Somente com estoque disponível must bind to only_em_estoque").toBe(false);
    expect(ecommerce.checked, "Somente elegíveis ao e-commerce must bind to only_ecommerce_eligible").toBe(false);

    fireEvent.click(revenda);
    await waitFor(() =>
      expect(setSellableAssortment).toHaveBeenLastCalledWith({
        only_revenda: false,
        only_em_estoque: false,
        only_ecommerce_eligible: false,
      }),
    );
    await waitFor(() =>
      expect(revenda.checked, "Somente produtos de revenda must bind to only_revenda").toBe(false),
    );

    fireEvent.click(estoque);
    await waitFor(() =>
      expect(setSellableAssortment).toHaveBeenLastCalledWith({
        only_revenda: false,
        only_em_estoque: true,
        only_ecommerce_eligible: false,
      }),
    );
    await waitFor(() =>
      expect(estoque.checked, "Somente com estoque disponível must bind to only_em_estoque").toBe(true),
    );

    fireEvent.click(estoque);
    await waitFor(() =>
      expect(setSellableAssortment).toHaveBeenLastCalledWith({
        only_revenda: false,
        only_em_estoque: false,
        only_ecommerce_eligible: false,
      }),
    );
    await waitFor(() =>
      expect(estoque.checked, "Somente com estoque disponível must bind to only_em_estoque").toBe(false),
    );

    fireEvent.click(ecommerce);
    await waitFor(() =>
      expect(setSellableAssortment).toHaveBeenLastCalledWith({
        only_revenda: false,
        only_em_estoque: false,
        only_ecommerce_eligible: true,
      }),
    );
    await waitFor(() =>
      expect(ecommerce.checked, "Somente elegíveis ao e-commerce must bind to only_ecommerce_eligible").toBe(true),
    );

    expect(window.localStorage.length).toBe(0);
    expect(window.sessionStorage.length).toBe(0);
    expect(setItem).not.toHaveBeenCalled();
    setItem.mockRestore();
  });

  it("pins each assortment toggle to its own server field", async () => {
    renderPage();
    await waitFor(() =>
      expect(
        screen.getAllByRole("checkbox"),
        "Sortimento vendável must render toggles for only_revenda, only_em_estoque, and only_ecommerce_eligible",
      ).toHaveLength(3),
    );
    const revenda = assortmentToggle("Somente produtos de revenda");
    const estoque = assortmentToggle("Somente com estoque disponível");
    const ecommerce = assortmentToggle("Somente elegíveis ao e-commerce");

    expect(revenda.checked, "Somente produtos de revenda must bind to only_revenda").toBe(true);
    fireEvent.click(revenda);
    await waitFor(() =>
      expect(revenda.checked, "Somente produtos de revenda must bind to only_revenda").toBe(false),
    );

    fireEvent.click(estoque);
    await waitFor(() =>
      expect(estoque.checked, "Somente com estoque disponível must bind to only_em_estoque").toBe(true),
    );
    expect(setSellableAssortment).toHaveBeenLastCalledWith({
      only_revenda: false,
      only_em_estoque: true,
      only_ecommerce_eligible: false,
    });
    fireEvent.click(estoque);
    await waitFor(() =>
      expect(estoque.checked, "Somente com estoque disponível must bind to only_em_estoque").toBe(false),
    );

    fireEvent.click(ecommerce);
    await waitFor(() =>
      expect(ecommerce.checked, "Somente elegíveis ao e-commerce must bind to only_ecommerce_eligible").toBe(true),
    );
    expect(setSellableAssortment).toHaveBeenLastCalledWith({
      only_revenda: false,
      only_em_estoque: false,
      only_ecommerce_eligible: true,
    });
  });

  it("reports when assortment counts cannot be calculated instead of showing zero", async () => {
    getCatalogAssortmentCounts.mockRejectedValueOnce(new Error("counts unavailable"));
    renderPage();

    expect(
      await screen.findByText("Não foi possível calcular quantos produtos entram no sortimento."),
    ).toBeInTheDocument();
    expect(screen.queryByText("Resultado: 0 de 0 produtos")).not.toBeInTheDocument();
  });

  it("recalculates the result line after a toggle instead of showing the pre-change count", async () => {
    renderPage();

    expect(await screen.findByText("Resultado: 2 de 4 produtos")).toBeInTheDocument();
    fireEvent.click(assortmentToggle("Somente produtos de revenda"));

    expect(await screen.findByText("Resultado: 3 de 4 produtos")).toBeInTheDocument();
    expect(screen.queryByText("Resultado: 2 de 4 produtos")).not.toBeInTheDocument();
  });

  it("reports a failed assortment read without guessing toggle positions", async () => {
    getSellableAssortment.mockRejectedValueOnce(new Error("assortment unavailable"));
    renderPage();

    expect(await screen.findByText("Não foi possível ler a regra de sortimento configurada.")).toBeInTheDocument();
    expect(
      screen.queryAllByRole("checkbox"),
      "Failed assortment read must render no toggle for only_revenda, only_em_estoque, or only_ecommerce_eligible",
    ).toHaveLength(0);
  });

  it("renders the configure-source state for an unknown assortment source without guessing toggles", async () => {
    getSellableAssortment.mockRejectedValueOnce(
      new MarketplaceCentralClientError(400, "unknown_erp_source", "fonte não configurada", {}),
    );
    renderPage();

    expect(await screen.findByTestId("sellable-assortment-source-unset")).toHaveTextContent(
      "Nenhuma fonte definida ainda — escolha a fonte que o app vai ler.",
    );
    expect(
      screen.queryAllByRole("checkbox"),
      "Unknown assortment source must render no toggle for only_revenda, only_em_estoque, or only_ecommerce_eligible",
    ).toHaveLength(0);
    expect(screen.queryByText("Não foi possível ler a regra de sortimento configurada.")).not.toBeInTheDocument();
  });

  it("renders the configure-source state when counts have no active source", async () => {
    getActiveSource.mockRejectedValueOnce(
      new MarketplaceCentralClientError(400, "unknown_erp_source", "fonte não configurada", {}),
    );
    renderPage();

    expect(await screen.findByTestId("sellable-assortment-source-unset")).toHaveTextContent(
      "Nenhuma fonte definida ainda — escolha a fonte que o app vai ler.",
    );
    expect(screen.queryByText(/Resultado:/)).not.toBeInTheDocument();
    expect(getCatalogAssortmentCounts).not.toHaveBeenCalled();
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
    createErpImport.mockRejectedValue(
      new MarketplaceCentralClientError(409, "duplicate_file", "arquivo já importado", {
        import_id: "imp_003",
        protocol: "#003-E",
      }),
    );
    renderPage();
    selectFile();
    fireEvent.click(screen.getByTestId("erp-import-submit"));

    const error = await screen.findByTestId("erp-import-error");
    expect(error).toHaveTextContent("#003-E");
    expect(screen.queryByTestId("erp-import-result")).not.toBeInTheDocument();
  });

  it("shows the missing-column name on a 422 missing_required_column", async () => {
    createErpImport.mockRejectedValue(
      new MarketplaceCentralClientError(422, "missing_required_column", "coluna obrigatória ausente", {
        column: "CUSTO",
      }),
    );
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

  it("shows the source the server reports and writes a switch back to the server", async () => {
    renderPage();
    const xlsx = (await screen.findByTestId("active-source-xlsx")) as HTMLInputElement;
    const catalogo = screen.getByTestId("active-source-catalogo_cliente") as HTMLInputElement;
    await waitFor(() => expect(xlsx.checked).toBe(true));
    expect(catalogo.checked).toBe(false);

    fireEvent.click(catalogo);

    // The write goes to the server; the selection only moves once it succeeds,
    // because the source decides what the whole platform reads.
    await waitFor(() => expect(setActiveSource).toHaveBeenCalledWith({ active_source: "catalogo_cliente" }));
    await waitFor(() => expect(catalogo.checked).toBe(true));
    expect(xlsx.checked).toBe(false);
  });

  it("offers the live Sankhya source, not only the two upload snapshots", async () => {
    renderPage();
    expect(await screen.findByTestId("active-source-sankhya")).toBeInTheDocument();
  });

  it("keeps the previous source selected and says so when the write fails", async () => {
    setActiveSource.mockRejectedValue({ status: 400, error: { code: "unknown_erp_source" } });
    renderPage();
    const xlsx = (await screen.findByTestId("active-source-xlsx")) as HTMLInputElement;
    await waitFor(() => expect(xlsx.checked).toBe(true));

    fireEvent.click(screen.getByTestId("active-source-catalogo_cliente"));

    expect(await screen.findByTestId("active-source-error")).toBeInTheDocument();
    expect(xlsx.checked).toBe(true);
  });

  // A workspace that never chose a source has no row, and the server fails
  // closed with 400 unknown_erp_source. That is the FIRST state of every
  // install, so this card has to stay usable: it is the only place the source
  // gets chosen, and a disabled selector leaves the whole platform unreadable.
  it("lets the operator choose the source when the tenant has none configured yet", async () => {
    getActiveSource.mockRejectedValue(
      new MarketplaceCentralClientError(400, "unknown_erp_source", "fonte não configurada", {}),
    );
    setActiveSource.mockResolvedValue(activeSourceConfig("sankhya"));
    renderPage();

    const sankhya = (await screen.findByTestId("active-source-sankhya")) as HTMLInputElement;
    await waitFor(() => expect(screen.getByTestId("active-source-unset")).toBeInTheDocument());
    expect(sankhya.closest("fieldset")).not.toBeDisabled();
    expect(sankhya.checked).toBe(false);

    fireEvent.click(sankhya);

    await waitFor(() => expect(setActiveSource).toHaveBeenCalledWith({ active_source: "sankhya" }));
  });

  it("offers the Mercado Livre connect button and starts no flow until it is clicked", () => {
    renderPage();
    const connect = screen.getByTestId("provider-connect-ml");
    expect(connect).toBeEnabled();
    expect(connect).toHaveTextContent("Conectar");
    // listIntegrationInstallations is no longer a signal exclusive to the
    // connect click: ConnectionHealthCard (Task 6, F-A1/F-A2) reads
    // useInstallation(), whose InstallationProvider fetches the list eagerly
    // on mount to render connection health for every installation, connect
    // flow or not. The invariant this test actually cares about — clicking
    // is required before anything starts — is covered by "busy" staying
    // false and the label staying "Conectar" (ProviderConnectCard.connect()
    // flips both synchronously on click, before any await).
  });
});
