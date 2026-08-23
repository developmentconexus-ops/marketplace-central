import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifact = resolve(root, 'qualification/d6-r2-wireframes/b11-notifications-inbox.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifact), 'B11 rendered artifact missing: qualification/d6-r2-wireframes/b11-notifications-inbox.html');

const html = readFileSync(artifact, 'utf8');

assert(html.includes('data-p8-status="candidate"'), 'B11 must identify itself as P8 CANDIDATE');
assert(!html.includes('data-p8-status="locked"'), 'B11 must not self-claim operator LOCK');
assert(html.includes('data-surface="R128"'), 'B11 must identify the full personal Inbox R128 surface');
assert(html.includes('id="fullInbox"'), 'B11 full Inbox host missing');
assert(html.includes('data-representation="structured-list"'), 'B11 must use the approved structured-list representation');
assert(!/<table\b/i.test(html), 'B11 must not introduce a table baseline');

const expectedNavLabels = [
  'Visão geral', 'Preparação', 'Anúncios', 'Preços', 'Disponibilidade',
  'Visão operacional', 'Vendas', 'Expedição', 'Pós-venda', 'Performance',
  'Mercado', 'Economia', 'Trabalho', 'Aprovações', 'Configurações',
];
for (const label of expectedNavLabels) {
  assert(html.includes(`>${label}</button>`) || html.includes(`>${label}</a>`), `locked B00 IA label missing from B11 shell: ${label}`);
}
assert(!/class="[^"]*nav-item[^"]*"[^>]*>\s*Notificações\s*</u.test(html), 'B11 must remain outside the global sidebar');

assert(html.includes('id="notificationBell"'), 'locked B00-R2 bell missing from B11 shell');
assert(html.includes('data-operation="ListMyNotifications"'), 'B11 reads must trace ListMyNotifications');
assert(!html.includes('data-required-permission="notifications.manage"'), 'self Inbox must not require notifications.manage');

for (const lens of ['active', 'archived']) {
  assert(html.includes(`data-archive-lens="${lens}"`), `archive lens missing: ${lens}`);
}
for (const lens of ['all', 'unread', 'read']) {
  assert(html.includes(`data-read-lens="${lens}"`), `read lens missing: ${lens}`);
}
assert(html.includes('id="kindFilter"'), 'NotificationKind filter missing');
assert(html.includes('data-filter="notification_kind"'), 'kind filter must trace notification_kind');

const kinds = [
  'MARKETPLACE_INSTALLATION_ATTENTION',
  'OFFERING_ASYNC_ACTION_RESULT',
  'AVAILABILITY_ATTENTION',
  'ECONOMIC_RECONCILIATION_ATTENTION',
  'NEW_MARKETPLACE_SALE',
  'SALE_ATTENTION',
  'MATERIALIZATION_ATTENTION',
  'FULFILLMENT_ACTIONABLE',
  'FULFILLMENT_ATTENTION',
  'SHIPMENT_EXCEPTION',
  'POST_SALE_ATTENTION',
  'WORK_ASSIGNMENT',
  'AUTHORIZATION_ACTION_REQUIRED',
  'AUTHORIZATION_DECISION_RESULT',
];
for (const kind of kinds) {
  assert(html.includes(`value="${kind}"`), `accepted NotificationKind option missing: ${kind}`);
}

assert(html.includes('id="loadMore"'), 'cursor continuation control missing');
assert(html.includes('data-cursor-source="next_cursor"'), 'continuation must be cursor-based');
assert(!html.includes('data-total-count='), 'B11 must not invent total_count');
assert(!html.includes('data-unread-count='), 'B11 must not invent unread count');
assert(!html.includes('data-action="mark-all-read"'), 'mark-all-read is forbidden');
assert(!html.includes('data-action="bulk-archive"'), 'bulk archive is forbidden');
assert(!html.includes('data-filter="text-search"'), 'text search is not admitted for B11');

assert(html.includes('data-kind="OFFERING_ASYNC_ACTION_RESULT"'), 'B11 must demonstrate F02 typed-result rendering');
assert(html.includes('data-offering-outcome="converged"'), 'F02 typed outcome example missing');
assert(html.includes('data-kind="AUTHORIZATION_DECISION_RESULT"'), 'B11 must demonstrate F14 typed-result rendering');
assert(html.includes('data-authorization-outcome="authorize"'), 'F14 typed outcome example missing');

assert(html.includes('data-action="open-source"'), 'source continuation control missing');
for (const action of ['mark-read', 'mark-unread', 'archive', 'restore']) {
  assert(html.includes(`data-action="${action}"`), `awareness action missing: ${action}`);
}
assert(html.includes('data-operation="UpdateMyNotificationAwarenessState"'), 'awareness mutations must trace UpdateMyNotificationAwarenessState');
assert(html.includes('function handleSourceOpen'), 'source re-authorization simulation missing');
assert(html.includes('function applyAwarenessMutation'), 'awareness mutation simulation missing');

for (const scenario of ['populated', 'known-empty', 'request-unavailable', 'stale-awareness-write', 'source-access-denied', 'organization-switch']) {
  assert(html.includes(`value="${scenario}"`), `deterministic B11 scenario missing: ${scenario}`);
}
assert(html.includes('data-knowledge-state="known-empty"'), 'known-empty state marker missing');
assert(html.includes('data-knowledge-state="unavailable"'), 'unavailable state marker missing');
assert(html.includes('data-failure="stale-awareness-write"'), 'stale awareness-write state missing');
assert(html.includes('data-failure="source-access-denied"'), 'source-denied state missing');

assert(html.includes('function invalidateInboxContext'), 'Organization Inbox invalidation behavior missing');
assert(html.includes('data-responsive-law="mobile-stacked-inbox-actions"'), 'mobile B11 structural law missing');
assert(html.includes('data-not-rendered="B12"'), 'B11 must explicitly keep B12 out of this baseline');
assert(!html.includes('id="notificationSettings"'), 'B12 routing Settings must not be rendered in B11');

console.log('d6_r_b11_status=CANDIDATE');
console.log('d6_r_b11_representation=STRUCTURED_LIST');
console.log('d6_r_b11_notification_kinds=14/14');
console.log('d6_r_b11_totals_bulk_search=FORBIDDEN');
console.log('d6_r_b11_b12=NOT_RENDERED');
console.log('d6_r_b11_wireframe=PASS');
