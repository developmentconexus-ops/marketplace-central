import { Button } from "./Button";

interface ErrorStateProps {
  onRetry: () => void;
  detail?: string;
}

export function ErrorState({ onRetry, detail }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-start gap-3 text-sm text-red-700">
      <p role="alert">{detail ? <>Erro ao carregar. {detail}</> : "Erro ao carregar."}</p>
      <Button variant="secondary" onClick={onRetry}>
        Tentar novamente
      </Button>
    </div>
  );
}
