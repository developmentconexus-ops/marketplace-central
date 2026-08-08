import { ImportacaoSection } from "./ImportacaoSection";

export function ImportacoesPage() {
  return (
    <section
      aria-labelledby="importacoes-title"
      className="mx-auto flex max-w-5xl flex-col gap-[14px]"
    >
      <header>
        <h1 id="importacoes-title" className="text-[22px] font-bold tracking-tight text-ink">
          Importações
        </h1>
        <p className="mt-1 text-sm text-muted">
          Histórico de importações do ERP por protocolo. Somente leitura.
        </p>
      </header>
      <ImportacaoSection />
    </section>
  );
}
