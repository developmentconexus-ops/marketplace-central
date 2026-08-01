import { isApiError, type SyncHealth, type SyncHealthEntity, type SyncHealthWebhook } from "@marketplace-central/sdk-runtime";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import { useSyncHealthQuery } from "@marketplace-central/web-query";
import type { UseQueryResult } from "@tanstack/react-query";
import { useClient } from "../../app/ClientContext";

// Three-state badge discriminator (validation contract M09-C2/M09-C4, audited
// P7 r02 A-8): a STATE read off the payload, never a time cutoff. Red beats
// green even for a recent last_success_at — a run that is failing right now
// is failing, no matter how recently it last worked.
type EntityTone = "green" | "red" | "gray";

function entityTone(entity: SyncHealthEntity): EntityTone {
  if (entity.consecutive_failures > 0) return "red";
  if (entity.last_success_at !== null) return "green";
  return "gray";
}

const toneBadgeClassName: Record<EntityTone, string> = {
  green: "bg-accent-soft text-accent-ink",
  red: "bg-warn-soft text-warn",
  gray: "bg-surface-2 text-faint",
};

// Generic title-case label — no hardcoded entity allowlist (feature brief
// negative scenario: an unknown/future entity name must still render).
function entityLabel(entity: string): string {
  return entity
    .split("_")
    .filter((part) => part.length > 0)
    .map((part) => part[0].toUpperCase() + part.slice(1))
    .join(" ");
}

// Relative pt-BR timestamp, same register as DashboardPage.tsx's
// formatLastImportAge / ListingsRefreshControl.tsx's age label (that helper
// is local to its own module, so this mirrors the idiom rather than importing
// a private function). Callers only invoke this for a NON-NULL timestamp — a
// null last_success_at renders "nunca" instead of ever reaching here, so a
// null timestamp can never be humanized into a fake recency.
function formatRelative(iso: string, now: number = Date.now()): string {
  const at = new Date(iso).getTime();
  if (Number.isNaN(at)) return iso;
  const seconds = Math.max(0, Math.floor((now - at) / 1000));
  if (seconds < 60) return "há menos de 1 min";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `há ${minutes} min`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `há ${hours} h`;
  const days = Math.floor(hours / 24);
  return `há ${days} d`;
}

function entityBadgeLabel(entity: SyncHealthEntity, tone: EntityTone): string {
  if (tone === "red") {
    return entity.consecutive_failures === 1 ? "1 falha" : `${entity.consecutive_failures} falhas`;
  }
  if (tone === "green") return "ok";
  return "nunca";
}

function EntityRow({ entity }: { entity: SyncHealthEntity }) {
  const tone = entityTone(entity);
  return (
    <div
      className="flex items-center justify-between gap-3 rounded-control border border-border bg-surface-2 px-3 py-2 text-xs"
      data-testid={`sync-health-entity-${entity.entity}`}
    >
      <div className="flex flex-col gap-0.5">
        <span className="font-medium text-ink">
          {entityLabel(entity.entity)}
          {entity.phase ? <span className="ml-1.5 font-normal text-faint">({entity.phase})</span> : null}
        </span>
        {entity.last_success_at !== null ? (
          <span className="text-faint" title={entity.last_success_at}>
            {formatRelative(entity.last_success_at)}
          </span>
        ) : (
          <span className="text-faint" data-testid={`sync-health-never-${entity.entity}`}>
            nunca
          </span>
        )}
      </div>
      <span
        className={`inline-flex items-center gap-1 whitespace-nowrap rounded-pill px-2 py-0.5 font-medium ${toneBadgeClassName[tone]}`}
        title={tone === "red" && entity.last_error ? entity.last_error.message : undefined}
        data-testid={`sync-health-badge-${entity.entity}`}
      >
        <span className="h-1.5 w-1.5 rounded-pill bg-current" aria-hidden="true" />
        {entityBadgeLabel(entity, tone)}
      </span>
    </div>
  );
}

function WebhookSection({ webhook }: { webhook: SyncHealthWebhook }) {
  // The payload's initial state {null, 0, 0} is identical for "webhook never
  // configured" and "configured, inbox just quiet" (IC-05, F-r07-3) — the
  // card states the observed FACT only and never claims a configuration
  // verdict it cannot know.
  if (webhook.last_notification_at === null) {
    return (
      <div className="mt-3 border-t border-border pt-3">
        <h3 className="text-xs font-semibold text-ink">Notificações (webhook)</h3>
        <p className="mt-1 text-xs text-faint" data-testid="sync-health-webhook-initial">
          Nenhuma notificação recebida.
        </p>
      </div>
    );
  }
  return (
    <div className="mt-3 border-t border-border pt-3">
      <h3 className="text-xs font-semibold text-ink">Notificações (webhook)</h3>
      <dl className="mt-1 flex flex-col gap-0.5 text-xs" data-testid="sync-health-webhook-active">
        <div>
          <dt className="inline text-faint">última notificação: </dt>
          <dd className="inline text-ink" title={webhook.last_notification_at}>
            {formatRelative(webhook.last_notification_at)}
          </dd>
        </div>
        <div>
          <dt className="inline text-faint">pendentes: </dt>
          <dd className="inline text-ink">{webhook.pending}</dd>
        </div>
        <div>
          <dt className="inline text-faint">descartadas (24h): </dt>
          <dd className="inline text-ink">{webhook.dropped_24h}</dd>
        </div>
      </dl>
    </div>
  );
}

// IC-05's error matrix only names the generic envelope (no per-code family
// for /sync/health), so there is no known code to branch on with hasCode —
// isApiError alone is enough to name the real code/message instead of a
// generic "algo deu errado".
function syncHealthErrorDetail(error: unknown): string {
  if (isApiError(error)) {
    return `${error.code}: ${error.message}`;
  }
  return "Falha inesperada ao carregar a saúde do sync.";
}

// Regression guard for the blank-card defect (live-drive, hub-confirmed
// twice): `isLoading ? … : isError ? … : data ? … : null` reads three
// INDEPENDENTLY-readable fields off UseQueryResult and treats them as if
// they always partition into exactly those three cases. They don't —
// react-query v5's own QueryObserverPendingResult (status:"pending",
// data:undefined, isError:false) covers both the ordinary "still fetching"
// shape (fetchStatus:"fetching", isLoading:true) AND the paused/indeterminate
// shape the hub reproduced live (fetchStatus:"paused" or "idle",
// isLoading:false) — a real, typed, reachable combination that isn't any of
// the three named branches, so it silently fell through to `null`.
//
// The fix below switches the discriminant to `query.status`, which
// @tanstack/query-core types as `QueryStatus = 'pending' | 'error' |
// 'success'` — a closed 3-member string-literal union, not three separate
// booleans. An if/else-if/else chain over a closed union is exhaustive BY
// CONSTRUCTION: once "error" and "success" are ruled out, TypeScript's
// control-flow narrowing leaves only "pending" as the remaining possibility,
// so the final `else` is not an untested guess about what's left over — it
// is the only other value the type can hold. Within "pending", fetchStatus
// (`'fetching' | 'paused' | 'idle'`) is used to tell an in-flight first
// fetch from the paused/idle-with-no-data shape, and the latter renders a
// named, honest "state unknown" ErrorState rather than nothing.
//
// ADR-17: an unreadable/indeterminate state must never render as blank,
// zero, or a fabricated "ok" — it must say plainly that it is unknown. This
// wording ("não foi possível carregar o status") is deliberately distinct
// from WebhookSection's "Nenhuma notificação recebida." (M09-C4's honest
// never-notified FACT) — this is a failed/indeterminate READ, not an
// observed fact about the webhook inbox, and must not borrow that copy.
function SyncHealthCardBody({ query }: { query: UseQueryResult<SyncHealth> }) {
  if (query.status === "error") {
    return (
      <div className="mt-2">
        <ErrorState onRetry={() => void query.refetch()} detail={syncHealthErrorDetail(query.error)} />
      </div>
    );
  }
  if (query.status === "success") {
    return (
      <>
        <div className="mt-3 flex flex-col gap-2" data-testid="sync-health-entities">
          {query.data.entities.map((entity) => (
            <EntityRow key={entity.entity} entity={entity} />
          ))}
        </div>
        <WebhookSection webhook={query.data.webhook} />
      </>
    );
  }
  // "error" and "success" have both been ruled out above, so TypeScript
  // narrows query.status to its one remaining member: "pending".
  // A first fetch actually in flight renders the ordinary loading state.
  if (query.fetchStatus === "fetching") {
    return <LoadingState />;
  }
  // fetchStatus is "paused" or "idle" with status still "pending": no data,
  // no error, nothing fetching — this is the exact shape the hub captured
  // live (isLoading:false, isError:false, data:undefined). Name it honestly
  // instead of rendering nothing; retry lets the user re-attempt the read.
  return (
    <div className="mt-2" data-testid="sync-health-unknown">
      <ErrorState
        onRetry={() => void query.refetch()}
        detail="Não foi possível carregar o status da sincronização. Estado desconhecido."
      />
    </div>
  );
}

// SyncHealthCard is read-only (no mutating actions this milestone) and
// isolates its own fetch: the query's loading/error/data states are all
// handled locally, so a failure here renders a card-scoped ErrorState and
// never throws past this component into the rest of IntegracoesPage.
export function SyncHealthCard() {
  const client = useClient();
  const query = useSyncHealthQuery(client);

  return (
    <section aria-labelledby="sync-health-title" className="rounded-card border border-border bg-surface p-4">
      <h2 id="sync-health-title" className="text-sm font-semibold text-ink">
        Saúde do sync
      </h2>
      <SyncHealthCardBody query={query} />
    </section>
  );
}
