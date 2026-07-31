import type {
  ActiveSourceName,
  ErpImportSourceInput,
  SellableAssortmentConfig,
  SetSellableAssortmentRequest,
} from "@marketplace-central/sdk-runtime";
import { EmptyState, ErrorState, LoadingState } from "@marketplace-central/ui";
import {
  useActiveSourceQuery,
  useCatalogAssortmentCountsQuery,
  useSellableAssortmentQuery,
  useSetActiveSourceMutation,
  useSetSellableAssortmentMutation,
} from "@marketplace-central/web-query";
import { useId, useRef, useState, type ChangeEvent, type DragEvent } from "react";
import { useClient } from "../../app/ClientContext";
import { ImportacaoSection } from "../importacoes/ImportacaoSection";
import { useErpImportDetail } from "../vinculos/useErpImports";
import { useErpImportUpload, type ErpImportUploadError } from "./useErpImportUpload";

type SourceOption = {
  value: ErpImportSourceInput;
  label: string;
  hint: string;
};

// Catálogo do cliente is the lenient path (the customer's raw product export
// carries no CUSTO/ESTOQUE_FISICO — those stay honest-unknown, ADR-17). Sankhya
// is the strict cost+stock path that rejects a workbook missing required
// columns. Default to the lenient path: it is the demo's live import.
const SOURCE_OPTIONS: SourceOption[] = [
  {
    value: "catalogo_cliente",
    label: "Catálogo do cliente",
    hint: "Planilha de produtos do cliente (custo e estoque podem faltar — ficam como desconhecido).",
  },
  {
    value: "xlsx",
    label: "Sankhya (custo + estoque)",
    hint: "Exportação Sankhya com colunas de custo e estoque obrigatórias.",
  },
];

const ACCEPT = ".xlsx,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";

function isXlsx(file: File): boolean {
  return file.name.toLowerCase().endsWith(".xlsx");
}

function uploadErrorMessage(error: ErpImportUploadError): string {
  switch (error.kind) {
    case "invalid_file":
      return "Arquivo inválido. Envie um .xlsx exportado do ERP.";
    case "missing_required_column":
      return error.column
        ? `Coluna obrigatória ausente: “${error.column}”. Selecione a fonte correta ou use “Catálogo do cliente”.`
        : "Coluna obrigatória ausente. Selecione a fonte correta ou use “Catálogo do cliente”.";
    case "duplicate_file":
      return error.protocol
        ? `Este arquivo já foi importado (protocolo ${error.protocol}). Nenhuma ação necessária.`
        : "Este arquivo já foi importado. Nenhuma ação necessária.";
    case "import_in_progress":
      return "Outra importação está em andamento. Aguarde a conclusão e tente novamente.";
    case "deadline_exceeded":
      return "A importação passou do tempo limite do servidor. O arquivo não foi importado — divida a planilha e tente novamente.";
    case "internal_error":
      return "Falha inesperada ao importar. Tente novamente.";
    case "network_error":
    default:
      return "Não foi possível conectar ao servidor. Verifique a conexão e tente novamente.";
  }
}

function IssueList({ title, issues, testId }: { title: string; issues: { row: number; code: string; detail: string; offending_value?: string | null }[]; testId: string }) {
  if (issues.length === 0) return null;
  return (
    <div className="mt-2">
      <p className="text-xs font-semibold text-ink">{title}</p>
      <ul className="mt-1 flex max-h-64 flex-col gap-1 overflow-auto" data-testid={testId}>
        {issues.map((issue, index) => (
          <li
            key={`${issue.row}-${issue.code}-${index}`}
            className="rounded-control border border-border bg-surface-2 px-2 py-1 text-xs text-muted"
          >
            <span className="font-medium">Linha {issue.row}</span> — {issue.code}: {issue.detail}
            {issue.offending_value ? <span className="text-faint"> (valor: {issue.offending_value})</span> : null}
          </li>
        ))}
      </ul>
    </div>
  );
}

// ResultSummary loads the freshly-created import's detail so the operator sees
// the real aceitos/rejeitados/avisos counts and the issue rows — the POST
// response only carries the protocol/status, not the tallies.
function ResultSummary({ importId, protocol }: { importId: string; protocol: string }) {
  const detailQuery = useErpImportDetail(importId);

  if (detailQuery.isPending) return <LoadingState />;
  if (detailQuery.isError) return <ErrorState onRetry={() => void detailQuery.refetch()} />;
  const detail = detailQuery.data;
  if (!detail) return <EmptyState />;

  return (
    <div
      className="mt-3 rounded-card border border-accent/40 bg-accent-soft p-3"
      role="status"
      data-testid="erp-import-result"
    >
      <p className="text-sm font-semibold text-accent-ink">
        Importação concluída — protocolo <span className="font-mono">{protocol}</span>
      </p>
      <dl className="mt-2 grid grid-cols-3 gap-x-4 gap-y-1 text-xs text-muted">
        <div>
          <dt className="text-faint">Aceitos</dt>
          <dd className="font-mono font-medium text-ink" data-testid="result-accepted">
            {detail.accepted_count}
          </dd>
        </div>
        <div>
          <dt className="text-faint">Rejeitados</dt>
          <dd className="font-mono font-medium text-ink" data-testid="result-rejected">
            {detail.rejected_count}
          </dd>
        </div>
        <div>
          <dt className="text-faint">Avisos</dt>
          <dd className="font-mono font-medium text-ink" data-testid="result-warnings">
            {detail.warning_count}
          </dd>
        </div>
      </dl>
      <IssueList title="Rejeitados" issues={detail.rejected_rows} testId="result-rejected-rows" />
      <IssueList title="Avisos" issues={detail.warnings} testId="result-warning-rows" />
    </div>
  );
}

function UploadCard() {
  const [file, setFile] = useState<File | null>(null);
  const [source, setSource] = useState<ErpImportSourceInput>("catalogo_cliente");
  const [dragging, setDragging] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const fieldId = useId();

  const upload = useErpImportUpload();
  const created = upload.data;

  function pickFile(next: File | null) {
    setLocalError(null);
    upload.reset();
    if (next && !isXlsx(next)) {
      setFile(null);
      setLocalError("Formato não suportado. Envie um arquivo .xlsx.");
      return;
    }
    setFile(next);
  }

  function onInputChange(event: ChangeEvent<HTMLInputElement>) {
    pickFile(event.target.files?.[0] ?? null);
  }

  function onDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    setDragging(false);
    pickFile(event.dataTransfer.files?.[0] ?? null);
  }

  function onSubmit() {
    if (!file) return;
    upload.mutate({ file, source });
  }

  return (
    <section aria-labelledby={`${fieldId}-title`} className="rounded-card border border-border bg-surface p-4">
      <h2 id={`${fieldId}-title`} className="text-sm font-semibold text-ink">
        Importar catálogo
      </h2>
      <p className="mt-1 text-xs text-faint">
        Envie a planilha de produtos (.xlsx). Selecione a fonte para escolher entre a importação do catálogo
        do cliente ou a exportação Sankhya com custo e estoque.
      </p>

      <fieldset className="mt-3">
        <legend className="text-xs font-semibold text-ink">Fonte</legend>
        <div className="mt-1 flex flex-col gap-2 sm:flex-row">
          {SOURCE_OPTIONS.map((option) => (
            <label
              key={option.value}
              className={`flex flex-1 cursor-pointer flex-col gap-0.5 rounded-control border p-2.5 text-xs ${
                source === option.value ? "border-accent bg-accent-soft" : "border-border bg-surface-2"
              }`}
            >
              <span className="flex items-center gap-2">
                <input
                  type="radio"
                  name={`${fieldId}-source`}
                  value={option.value}
                  checked={source === option.value}
                  onChange={() => {
                    setSource(option.value);
                    upload.reset();
                  }}
                />
                <span className="font-medium text-ink">{option.label}</span>
              </span>
              <span className="pl-6 text-faint">{option.hint}</span>
            </label>
          ))}
        </div>
      </fieldset>

      <div
        className={`mt-3 flex flex-col items-center justify-center gap-2 rounded-card border border-dashed p-6 text-center ${
          dragging ? "border-accent bg-accent-soft" : "border-border bg-surface-2"
        }`}
        onDragOver={(event) => {
          event.preventDefault();
          setDragging(true);
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={onDrop}
        data-testid="erp-import-dropzone"
      >
        <input
          ref={inputRef}
          type="file"
          accept={ACCEPT}
          className="sr-only"
          onChange={onInputChange}
          aria-label="Selecionar arquivo .xlsx"
          data-testid="erp-import-file-input"
        />
        <p className="text-xs text-muted">
          {file ? (
            <span className="font-medium text-ink" data-testid="erp-import-file-name">
              {file.name}
            </span>
          ) : (
            "Arraste o arquivo aqui ou selecione"
          )}
        </p>
        <button
          type="button"
          className="rounded-control border border-border bg-surface px-3 py-1.5 text-xs font-medium text-ink hover:bg-surface-2"
          onClick={() => inputRef.current?.click()}
        >
          Selecionar arquivo
        </button>
      </div>

      <div className="mt-3 flex items-center gap-3">
        <button
          type="button"
          className="rounded-control bg-accent px-4 py-2 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-50"
          disabled={!file || upload.isPending}
          onClick={onSubmit}
          data-testid="erp-import-submit"
        >
          {upload.isPending ? "Importando…" : "Importar"}
        </button>
        {upload.isPending ? <span className="text-xs text-muted">Processando planilha…</span> : null}
      </div>

      {localError ? (
        <p className="mt-2 text-xs text-warn" role="alert">
          {localError}
        </p>
      ) : null}
      {upload.isError && upload.error ? (
        <p className="mt-2 rounded-control border border-warn/40 bg-warn-soft px-2 py-1.5 text-xs text-warn" role="alert" data-testid="erp-import-error">
          {uploadErrorMessage(upload.error)}
        </p>
      ) : null}
      {created ? <ResultSummary importId={created.import_id} protocol={created.protocol} /> : null}
    </section>
  );
}

// ActiveSourceCard picks which source the whole platform reads: every
// mirror-backed read resolves it per tenant from the database, and the product
// sync scheduler runs the adapter it names. The selection is therefore server
// state (PUT /config/active-source), not a browser preference — while it lived
// in localStorage the labels changed here and the API kept serving the other
// source.
const ACTIVE_SOURCE_OPTIONS: { value: ActiveSourceName; label: string; hint: string }[] = [
  { value: "sankhya", label: "Sankhya (ao vivo)", hint: "Lê o ERP Sankhya direto, sem planilha." },
  { value: "xlsx", label: "Planilha Sankhya", hint: "Última planilha Sankhya importada — custo e estoque." },
  { value: "catalogo_cliente", label: "Catálogo do cliente", hint: "Catálogo importado do cliente — custo e estoque aparecem como “—”." },
];

const SOURCE_UNSET_COPY = "Nenhuma fonte definida ainda — escolha a fonte que o app vai ler.";

// A workspace that never chose a source has no row at all, and the server fails
// closed on the read (400 unknown_erp_source). That is not a broken platform,
// it is the first state of every install — and this card is the only place the
// source gets chosen, so it must stay usable and say what is missing instead of
// reporting a read failure.
function isSourceUnsetError(error: unknown): boolean {
  if (!error || typeof error !== "object") return false;
  const { status, error: code } = error as { status?: number; error?: unknown };
  return status === 400 && code === "unknown_erp_source";
}

function ActiveSourceCard() {
  const client = useClient();
  const activeSourceQuery = useActiveSourceQuery(client);
  const setActiveSource = useSetActiveSourceMutation(client);
  const activeSource = activeSourceQuery.data?.active_source;
  const readFailed = activeSourceQuery.isError && !isSourceUnsetError(activeSourceQuery.error);
  // "No source" is any settled state that produced no source: the fail-closed
  // 400, and also a read that ended without data or a reportable error. Both
  // mean the same thing to the operator — nothing is configured, choose one.
  const sourceUnset = !activeSource && !activeSourceQuery.isFetching && !readFailed;
  const fieldId = useId();
  return (
    <section aria-labelledby={`${fieldId}-active-title`} className="rounded-card border border-border bg-surface p-4">
      <h2 id={`${fieldId}-active-title`} className="text-sm font-semibold text-ink">
        Fonte ativa
      </h2>
      <p className="mt-1 text-xs text-faint">
        Selecione a fonte que o app inteiro lê: catálogo, custo, estoque, preço e a sincronização automática.
      </p>
      {/* Disabled only while a request is actually in flight: keying this off
          `isPending` kept the selector dead forever once the read failed, which
          is exactly the state a fresh workspace starts in. */}
      <fieldset className="mt-3" disabled={activeSourceQuery.isFetching || setActiveSource.isPending}>
        <legend className="sr-only">Fonte ativa de dados</legend>
        <div className="flex flex-col gap-2 sm:flex-row" data-testid="active-source-selector">
          {ACTIVE_SOURCE_OPTIONS.map((option) => (
            <label
              key={option.value}
              className={`flex flex-1 cursor-pointer flex-col gap-0.5 rounded-control border p-2.5 text-xs ${
                activeSource === option.value ? "border-accent bg-accent-soft" : "border-border bg-surface-2"
              }`}
            >
              <span className="flex items-center gap-2">
                <input
                  type="radio"
                  name={`${fieldId}-active-source`}
                  value={option.value}
                  checked={activeSource === option.value}
                  onChange={() => setActiveSource.mutate(option.value)}
                  data-testid={`active-source-${option.value}`}
                />
                <span className="font-medium text-ink">{option.label}</span>
              </span>
              <span className="pl-6 text-faint">{option.hint}</span>
            </label>
          ))}
        </div>
      </fieldset>
      {/* No selection is shown until the server answers, and a failed read or
          write says so: guessing a source here would tell the operator the app
          is reading data it is not. */}
      {activeSourceQuery.isFetching && !activeSourceQuery.data ? (
        <p className="mt-2 text-xs text-faint">Carregando fonte ativa…</p>
      ) : null}
      {sourceUnset ? (
        <p className="mt-2 text-xs text-faint" data-testid="active-source-unset">
          {SOURCE_UNSET_COPY}
        </p>
      ) : null}
      {readFailed ? (
        <p className="mt-2 text-xs text-warn" role="alert">
          Não foi possível ler a fonte ativa configurada.
        </p>
      ) : null}
      {setActiveSource.isError ? (
        <p className="mt-2 text-xs text-warn" role="alert" data-testid="active-source-error">
          Não foi possível trocar a fonte ativa. A seleção continua a anterior.
        </p>
      ) : null}
    </section>
  );
}

const SELLABLE_ASSORTMENT_OPTIONS: {
  key: keyof SellableAssortmentConfig;
  label: string;
  hint: string;
}[] = [
  {
    key: "only_revenda",
    label: "Somente produtos de revenda",
    hint: "Mantém no sortimento apenas o que o ERP classifica como revenda.",
  },
  {
    key: "only_em_estoque",
    label: "Somente com estoque disponível",
    hint: "Considera o saldo disponível no CD e na loja física. Com a regra desligada, produtos sem saldo continuam no sortimento com o aviso Sem estoque.",
  },
  {
    key: "only_ecommerce_eligible",
    label: "Somente elegíveis ao e-commerce",
    hint: "Remove os produtos que o ERP marcou como fora do e-commerce. Produtos ainda sem definição continuam no sortimento.",
  },
];

function SellableAssortmentCard() {
  const client = useClient();
  const activeSourceQuery = useActiveSourceQuery(client);
  const assortmentQuery = useSellableAssortmentQuery(client);
  const countsQuery = useCatalogAssortmentCountsQuery(client, activeSourceQuery.data?.active_source);
  const setAssortment = useSetSellableAssortmentMutation(client);
  const assortment = assortmentQuery.data;
  const sourceUnset =
    isSourceUnsetError(activeSourceQuery.error) ||
    isSourceUnsetError(assortmentQuery.error) ||
    isSourceUnsetError(countsQuery.error);
  const fieldId = useId();

  function setOption(key: keyof SetSellableAssortmentRequest) {
    if (!assortment) return;
    const next = {
      ...assortment,
      [key]: !assortment[key],
    } as SetSellableAssortmentRequest;
    setAssortment.mutate(next);
  }

  return (
    <section aria-labelledby={`${fieldId}-assortment-title`} className="rounded-card border border-border bg-surface p-4">
      <h2 id={`${fieldId}-assortment-title`} className="text-sm font-semibold text-ink">
        Sortimento vendável
      </h2>
      <p className="mt-1 text-xs text-faint">
        Define quais produtos entram no sortimento que o app considera vendável.
      </p>
      <fieldset className="mt-3" disabled={assortmentQuery.isFetching || setAssortment.isPending || !assortment}>
        <legend className="sr-only">Regras do sortimento vendável</legend>
        <div className="flex flex-col gap-2" data-testid="sellable-assortment-selector">
          {SELLABLE_ASSORTMENT_OPTIONS.map((option) => (
            <label
              key={option.key}
              className={`flex cursor-pointer flex-col gap-0.5 rounded-control border p-2.5 text-xs ${
                assortment?.[option.key] ? "border-accent bg-accent-soft" : "border-border bg-surface-2"
              }`}
            >
              <span className="flex items-center gap-2">
                {assortment ? (
                  <input
                    type="checkbox"
                    name={`${fieldId}-${option.key}`}
                    checked={assortment[option.key]}
                    onChange={() => setOption(option.key)}
                    data-testid={`sellable-assortment-${option.key.replaceAll("_", "-")}`}
                  />
                ) : null}
                <span className="font-medium text-ink">{option.label}</span>
              </span>
              <span className="pl-6 text-faint">{option.hint}</span>
            </label>
          ))}
        </div>
      </fieldset>
      {sourceUnset ? (
        <p className="mt-2 text-xs text-faint" data-testid="sellable-assortment-source-unset">
          {SOURCE_UNSET_COPY}
        </p>
      ) : null}
      {assortmentQuery.isError && !sourceUnset ? (
        <p className="mt-2 text-xs text-warn" role="alert">
          Não foi possível ler a regra de sortimento configurada.
        </p>
      ) : null}
      {setAssortment.isError ? (
        <p className="mt-2 text-xs text-warn" role="alert">
          Não foi possível salvar a regra de sortimento. A configuração continua a anterior.
        </p>
      ) : null}
      {countsQuery.data && !sourceUnset ? (
        <p className="mt-2 text-xs text-faint">
          Resultado: {countsQuery.data.sellable_count} de {countsQuery.data.total_count} produtos
        </p>
      ) : null}
      {countsQuery.isError && !sourceUnset ? (
        <p className="mt-2 text-xs text-warn" role="alert">
          Não foi possível calcular quantos produtos entram no sortimento.
        </p>
      ) : null}
    </section>
  );
}

// ProviderConnectCard starts the Mercado Livre OAuth connect flow. A NEW
// installation is created per connect (the reauth path pins the existing
// installation to its original ML account via the account-mismatch guard, so a
// different account MUST get its own installation); the browser is then sent to
// the ML authorization URL and returns via the ngrok callback. Existing
// installations and their data are untouched — the new account appears as an
// extra entry in the installation selector.
function ProviderConnectCard() {
  const client = useClient();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function connect() {
    setBusy(true);
    setError(null);
    try {
      // Reuse an ML installation still waiting for its first authorization;
      // otherwise create a fresh one so the new account never collides with an
      // already-connected installation's account guard.
      const existing = await client.listIntegrationInstallations();
      // The never-authorized states are carried by `status` — there is no
      // top-level `state` field, so the old lookup never matched and every
      // click minted another installation the operator could not get rid of.
      const pending = existing.items.find(
        (i) =>
          i.provider_code === "mercado_livre" &&
          (i.status === "pending_connection" || i.status === "draft"),
      );
      let installationId = pending?.installation_id;
      if (!installationId) {
        installationId = `inst-mercado_livre-${crypto.randomUUID()}`;
        await client.createIntegrationInstallation({
          installation_id: installationId,
          provider_code: "mercado_livre",
          family: "marketplace",
          display_name: "Mercado Livre (cliente)",
        });
      }
      const start = await client.startIntegrationAuthorization(installationId);
      window.location.href = start.auth_url;
    } catch {
      setError("Não foi possível iniciar a conexão. Tente novamente.");
      setBusy(false);
    }
  }

  return (
    <section aria-labelledby="provider-connect-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="provider-connect-title" className="text-sm font-semibold text-ink">
        Conectar marketplace
      </h2>
      <p className="mt-1 text-xs text-faint">Conecte uma conta de marketplace para sincronizar anúncios e pedidos.</p>
      <div className="mt-3 flex items-center justify-between rounded-control border border-border bg-surface-2 p-3">
        <span className="text-sm font-medium text-ink">Mercado Livre</span>
        <button
          type="button"
          disabled={busy}
          onClick={connect}
          className="rounded-control border border-border bg-surface px-3 py-1.5 text-xs font-semibold text-ink hover:bg-surface-2 disabled:cursor-not-allowed disabled:text-muted"
          data-testid="provider-connect-ml"
        >
          {busy ? "Redirecionando…" : "Conectar"}
        </button>
      </div>
      {error ? <p className="mt-2 text-xs text-danger">{error}</p> : null}
    </section>
  );
}

export function IntegracoesPage() {
  return (
    <section aria-labelledby="integracoes-title" className="mx-auto flex max-w-5xl flex-col gap-[14px]">
      <header>
        <h1 id="integracoes-title" className="text-[22px] font-bold tracking-tight text-ink">
          Configuração da plataforma
        </h1>
        <p className="mt-1 text-sm text-muted">Importe o catálogo de produtos e conecte marketplaces.</p>
      </header>
      <ActiveSourceCard />
      <SellableAssortmentCard />
      <UploadCard />
      <ProviderConnectCard />
      <ImportacaoSection />
    </section>
  );
}
