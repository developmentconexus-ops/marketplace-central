import { Link, useParams } from "react-router-dom";
import { ImportChainPanel } from "./ImportChainPanel";

export function ImportacaoDetailPage() {
  const { importId } = useParams();

  return (
    <section aria-labelledby="importacao-detail-title" className="mx-auto flex max-w-5xl flex-col gap-[14px]">
      <header>
        <h1 id="importacao-detail-title" className="text-[22px] font-bold tracking-tight text-ink">
          Detalhe da importação
        </h1>
        <Link
          to="/importacoes"
          className="text-sm font-medium text-accent-ink underline decoration-accent-soft underline-offset-2"
        >
          Voltar para importações
        </Link>
      </header>
      {importId ? (
        <ImportChainPanel importId={importId} />
      ) : (
        <p role="alert" className="text-sm text-warn">
          Identificador da importação não informado.
        </p>
      )}
    </section>
  );
}
