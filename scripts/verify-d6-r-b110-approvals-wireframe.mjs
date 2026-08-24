import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifact = resolve(root, 'qualification/d6-r2-wireframes/b110-approvals.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifact), 'B110 rendered artifact missing: qualification/d6-r2-wireframes/b110-approvals.html');

const html = readFileSync(artifact, 'utf8');

assert(html.includes('data-p8-status="candidate"'), 'B110 must identify itself as P8 CANDIDATE');
assert(!html.includes('data-p8-status="locked"'), 'B110 must not self-claim operator LOCK');
assert(html.includes('data-block="B110"'), 'B110 block identity missing');
assert(html.includes('data-surface="R110-R111"'), 'B110 must identify the approvals collection/detail surface family');
assert(html.includes('id="approvalsWorkspace"'), 'B110 approvals workspace host missing');

const expectedNavLabels = [
  'Visão geral', 'Preparação', 'Anúncios', 'Preços', 'Disponibilidade',
  'Visão operacional', 'Vendas', 'Expedição', 'Pós-venda', 'Performance',
  'Mercado', 'Economia', 'Trabalho', 'Aprovações', 'Configurações',
];
for (const label of expectedNavLabels) {
  assert(html.includes(`>${label}</button>`) || html.includes(`>${label}</a>`), `locked B00 IA label missing from B110 shell: ${label}`);
}
assert(!/>\s*Histórico de aprovações\s*<\/a>/u.test(html), 'B110 must not create a second global approvals-history destination');

assert(html.includes('data-lens="actionable"'), 'B110 actionable local lens missing');
assert(html.includes('data-lens="history"'), 'B110 history local lens missing');
assert(html.includes('data-actionable-permission="governance.decide"'), 'actionable lens must be conditioned by governance.decide');
assert(html.includes('data-history-permission="governance.read"'), 'history lens must be conditioned by governance.read');
assert(html.includes('data-permission-law="independent"'), 'governance.decide and governance.read must remain independent');

for (const operationId of [
  'ListMyActionableAuthorizationRequests',
  'GetMyActionableAuthorizationRequest',
  'CreateAuthorizationDecision',
  'ListAuthorizationDecisions',
  'GetAuthorizationDecision',
]) {
  assert(html.includes(`data-operation="${operationId}"`), `B110 operation trace missing: ${operationId}`);
}

assert(html.includes('data-representation="structured-list"'), 'actionable approvals must use a structured list baseline');
assert(!html.includes('data-representation="table"'), 'actionable approvals must not default to a comparison table');
assert(!html.includes('data-filter="search"'), 'B110 must not invent approval search');
assert(!html.includes('data-total-count'), 'B110 must not invent approval totals');
assert(!html.includes('data-bulk-action'), 'B110 must not invent bulk approval actions');
assert(html.includes('data-cursor-source="next_cursor"'), 'actionable continuation must remain cursor-based');

assert(html.includes('data-request-route="/org/:organizationId/aprovacoes/solicitacoes/:authorizationRequestId"'), 'actionable request deep-link route missing');
assert(html.includes('data-decision-route="/org/:organizationId/aprovacoes/decisoes/:authorizationDecisionId"'), 'historical decision route missing');
assert(html.includes('data-entry-kind="AUTHORIZATION_ACTION_REQUIRED"'), 'F13 entry scenario missing');
assert(html.includes('data-source-ref="AuthorizationRequestRef"'), 'F13 must deep-link through AuthorizationRequestRef');
assert(html.includes('data-notification-authority="awareness-not-capability"'), 'Notification must not become a capability token');

for (const kind of ['listing_intent', 'price_intent', 'business_order_intent', 'invoicing_intent']) {
  assert(html.includes(`data-review-basis-kind="${kind}"`), `typed authorization review-basis example missing: ${kind}`);
}
assert(!html.includes('data-review-basis-kind="generic"'), 'generic authorization review basis is forbidden');
assert(!html.includes('data-review-payload="generic"'), 'generic authorization payload is forbidden');

assert(html.includes('data-action="approve-request"'), 'approve action missing');
assert(html.includes('data-action="reject-request"'), 'reject action missing');
assert(html.includes('data-confirmation-mode="inline"'), 'decision confirmation must remain inline with decision context');
assert(!/<dialog\b/i.test(html), 'B110 must not require a modal confirmation baseline');
assert(html.includes('data-precondition="If-Match"'), 'decision write must carry AuthorizationRequest If-Match');
assert(html.includes('data-retry-carrier="Idempotency-Key"'), 'decision write must carry Idempotency-Key');
assert(html.includes('data-client-body="outcome-only"'), 'frontend decision body must remain outcome-only');
assert(!html.includes('data-client-target-etag'), 'frontend must not send a governed target ETag');

assert(html.includes('data-failure="stale-decision"'), '412 stale-decision recovery state missing');
assert(html.includes('data-failure="validity-unavailable"'), '503 material-validity unavailable state missing');
assert(html.includes('data-auto-retry-decision="forbidden"'), 'consequential decision must not auto-retry');
assert(html.includes('data-state="decision-recorded"'), 'decision success state missing');
assert(html.includes('data-action="open-origin"'), 'post-decision source continuation missing');
assert(html.includes('data-source-open-law="reauthorize-current-owner"'), 'source continuation must reauthorize current owner truth');

assert(html.includes('data-history-filter="decided_from"'), 'history decided_from filter missing');
assert(html.includes('data-history-filter="decided_before"'), 'history decided_before filter missing');
assert(html.includes('data-history-mode="immutable-decision"'), 'historical detail must be immutable Decision truth');
assert(html.includes('data-history-actions="none"'), 'historical detail must not expose decision controls');

for (const scenario of [
  'actionable-populated',
  'actionable-empty',
  'request-price',
  'request-listing',
  'request-business-order',
  'request-invoicing',
  'inline-confirmation',
  'stale-decision',
  'validity-unavailable',
  'decision-recorded',
  'history-list',
  'history-detail',
  'decide-only',
  'read-only',
  'f13-deep-link',
  'organization-switch',
]) {
  assert(html.includes(`value="${scenario}"`), `deterministic B110 scenario missing: ${scenario}`);
}

for (const fn of [
  'showLens',
  'openActionableRequest',
  'beginDecisionConfirmation',
  'submitDecision',
  'handleStaleDecision',
  'handleValidityUnavailable',
  'showDecisionHistory',
  'invalidateApprovalContext',
]) {
  assert(html.includes(`function ${fn}`), `B110 interaction function missing: ${fn}`);
}

assert(html.includes('data-responsive-law="mobile-stacked-approval-actions"'), 'mobile B110 structural law missing');

console.log('d6_r_b110_status=CANDIDATE');
console.log('d6_r_b110_lenses=ACTIONABLE+HISTORY');
console.log('d6_r_b110_permissions=governance.decide|governance.read_INDEPENDENT');
console.log('d6_r_b110_review_basis=4/4');
console.log('d6_r_b110_confirmation=INLINE');
console.log('d6_r_b110_stale_and_503=EXPLICIT');
console.log('d6_r_b110_notification=AWAWARENESS_NOT_CAPABILITY');
console.log('d6_r_b110_wireframe=PASS');
