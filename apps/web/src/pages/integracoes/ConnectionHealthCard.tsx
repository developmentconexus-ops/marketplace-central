import type { IntegrationConnectionSnapshot, IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { installationsQueryKeys } from "@marketplace-central/web-query";
import { useQueryClient } from "@tanstack/react-query";
import { useInstallation } from "../../app/InstallationContext";

// Tom lido do ESTADO do payload, nunca de um corte de tempo — mesmo critério do
// SyncHealthCard (entityTone, SyncHealthCard.tsx:13). Uma conta que precisa de
// reautorização precisa disso agora, por mais recente que seja o last_verified_at.
type ConnectionTone = "green" | "amber" | "red" | "gray";

const stateTone: Record<IntegrationConnectionSnapshot["state"], ConnectionTone> = {
  connected: "green",
  degraded: "amber",
  needs_reauth: "red",
  disconnected: "red",
  pending_connection: "gray",
  draft: "gray",
};

const stateLabel: Record<IntegrationConnectionSnapshot["state"], string> = {
  connected: "Conectado",
  degraded: "Instável",
  needs_reauth: "Precisa reautorizar",
  disconnected: "Desconectado",
  pending_connection: "Aguardando autorização",
  draft: "Rascunho",
};

const nextActionLabel: Record<IntegrationConnectionSnapshot["next_action"], string> = {
  none: "",
  authorize: "Autorizar",
  reauth: "Reautorizar",
  configure: "Configurar",
  retry: "Repetindo automaticamente",
};

// packages/ui não tem classe "danger" de tema (verificado: grep -rn
// "danger-soft" apps/web/src packages/ui/src não bate nada) — só accent/warn/
// surface existem (SyncHealthCard.tsx:19-23). "red" usa a mesma classe de
// "amber" (bg-warn-soft) e se distingue só pelo rótulo, exatamente como o
// plano manda quando a classe de tema não existir.
const toneBadgeClassName: Record<ConnectionTone, string> = {
  green: "bg-accent-soft text-accent-ink",
  amber: "bg-warn-soft text-warn",
  red: "bg-warn-soft text-warn",
  gray: "bg-surface-2 text-faint",
};

function InstallationRow({ installation }: { installation: IntegrationInstallation }) {
  const connection = installation.connection;
  const tone = stateTone[connection.state];
  const action = nextActionLabel[connection.next_action];
  const reason = connection.reauth_reason?.trim() ?? "";

  return (
    <div
      className="flex flex-col gap-1 rounded-control border border-border bg-surface-2 px-3 py-2 text-xs"
      data-testid={`connection-health-${installation.installation_id}`}
    >
      <div className="flex items-center justify-between gap-3">
        <span className="font-medium text-ink">{installation.display_name}</span>
        <span
          className={`inline-flex items-center gap-1 whitespace-nowrap rounded-pill px-2 py-0.5 font-medium ${toneBadgeClassName[tone]}`}
        >
          <span className="h-1.5 w-1.5 rounded-pill bg-current" aria-hidden="true" />
          {stateLabel[connection.state]}
        </span>
      </div>
      {action ? <span className="text-faint">Ação: {action}</span> : null}
      {/* O motivo é o erro cru do provider. Mostrar cru é deliberado: é o único
          diagnóstico que existe e traduzi-lo apagaria o código que o operador
          precisa citar num chamado. */}
      {reason ? (
        <span className="break-words text-faint" data-testid={`connection-health-reason-${installation.installation_id}`}>
          {reason}
        </span>
      ) : null}
    </div>
  );
}

// Consome useInstallation() em vez de abrir um useQuery próprio: a lista já é
// buscada uma vez por installationsQueryKeys.list() em InstallationContext.tsx:31.
export function ConnectionHealthCard() {
  const { installations, status } = useInstallation();
  const queryClient = useQueryClient();

  return (
    <section aria-labelledby="connection-health-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="connection-health-title" className="text-sm font-semibold text-ink">
        Contas conectadas
      </h2>
      {status === "loading" ? <LoadingState /> : null}
      {/* ADR-17: leitura que falhou não pode renderizar verde nem branco. */}
      {status === "error" ? (
        <div className="mt-2" data-testid="connection-health-unknown">
          <ErrorState
            onRetry={() => void queryClient.refetchQueries({ queryKey: installationsQueryKeys.list() })}
            detail="Não foi possível carregar o estado das contas. Estado desconhecido."
          />
        </div>
      ) : null}
      {status !== "loading" && status !== "error" ? (
        <div className="mt-3 flex flex-col gap-2">
          {installations.map((installation) => (
            <InstallationRow key={installation.installation_id} installation={installation} />
          ))}
        </div>
      ) : null}
    </section>
  );
}
