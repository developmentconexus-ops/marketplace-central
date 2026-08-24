import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifact = resolve(root, 'qualification/d6-r2-wireframes/b12-notification-routing-settings.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifact), 'B12 rendered artifact missing: qualification/d6-r2-wireframes/b12-notification-routing-settings.html');

const html = readFileSync(artifact, 'utf8');

assert(html.includes('data-p8-status="candidate"'), 'B12 must identify itself as P8 CANDIDATE');
assert(!html.includes('data-p8-status="locked"'), 'B12 must not self-claim operator LOCK');
assert(html.includes('data-surface="R129"'), 'B12 must identify the R129 routing Settings surface');
assert(html.includes('id="notificationSettings"'), 'B12 notification Settings host missing');
assert(html.includes('data-required-permission="notifications.manage"'), 'B12 must be conditioned by notifications.manage');
assert(!html.includes('data-required-permission="access.read"'), 'B12 must not require access.read');
assert(!html.includes('role_keys'), 'B12 must not disclose role_keys');

const expectedNavLabels = [
  'Visão geral', 'Preparação', 'Anúncios', 'Preços', 'Disponibilidade',
  'Visão operacional', 'Vendas', 'Expedição', 'Pós-venda', 'Performance',
  'Mercado', 'Economia', 'Trabalho', 'Aprovações', 'Configurações',
];
for (const label of expectedNavLabels) {
  assert(html.includes(`>${label}</button>`) || html.includes(`>${label}</a>`), `locked B00 IA label missing from B12 shell: ${label}`);
}
assert(!/class="[^"]*nav-item[^"]*"[^>]*>\s*Notificações\s*</u.test(html), 'B12 must not turn Notifications into a global sidebar mass');
assert(html.includes('data-settings-destination="notifications"'), 'B12 must remain a local Settings destination');

for (const operationId of ['ListNotificationRoutes', 'ListNotificationRouteRecipientCandidates', 'SetNotificationRoute']) {
  assert(html.includes(`data-operation="${operationId}"`), `B12 operation trace missing: ${operationId}`);
}
assert(html.includes('data-route-owner="PersonalNotifications"'), 'route state must remain owned by PersonalNotifications');
assert(html.includes('data-recipient-owner="IdentityAccess"'), 'recipient discovery must remain owned by IdentityAccess');

const kinds = [
  'MARKETPLACE_INSTALLATION_ATTENTION',
  'AVAILABILITY_ATTENTION',
  'ECONOMIC_RECONCILIATION_ATTENTION',
  'NEW_MARKETPLACE_SALE',
  'SALE_ATTENTION',
  'MATERIALIZATION_ATTENTION',
  'FULFILLMENT_ACTIONABLE',
  'FULFILLMENT_ATTENTION',
  'SHIPMENT_EXCEPTION',
  'POST_SALE_ATTENTION',
];
for (const kind of kinds) {
  assert(html.includes(`data-route-kind="${kind}"`), `fixed ORG_ROUTED row missing: ${kind}`);
}

assert(html.includes('data-route-state="configured"'), 'B12 must demonstrate CONFIGURED state');
assert(html.includes('data-route-state="unconfigured"'), 'B12 must demonstrate UNCONFIGURED state');
assert(html.includes('Sem configuração de destinatários'), 'UNCONFIGURED must use the approved human wording');
assert(!html.includes('data-configured-empty="true"'), 'CONFIGURED([]) is forbidden');
assert(!html.includes('data-route-switch="disabled"'), 'generic disabled-switch semantics are forbidden');

assert(html.includes('id="routeEditor"'), 'inline route editor missing');
assert(html.includes('data-editor-mode="inline-row"'), 'B12 editor must expand inline in the affected row');
assert(!/<dialog\b/i.test(html), 'B12 must not use a mandatory modal baseline');
assert(!html.includes('data-editor-mode="drawer"'), 'B12 must not use a drawer-only editor');
assert(html.includes('data-action="save-route"'), 'save route action missing');
assert(html.includes('data-action="remove-configuration"'), 'remove configuration action missing');
assert(html.includes('data-desired-state="unconfigured"'), 'remove configuration must write UNCONFIGURED desired state');
assert(html.includes('data-precondition="If-Match"'), 'SetNotificationRoute must carry If-Match stale-write protection');

assert(html.includes('id="recipientCandidates"'), 'recipient candidate selector missing');
assert(html.includes('data-candidate-fields="principal_id display_name"'), 'candidate projection must stay principal_id + display_name only');
assert(html.includes('id="loadMoreCandidates"'), 'cursor continuation for recipient candidates missing');
assert(html.includes('data-cursor-source="next_cursor"'), 'candidate continuation must be cursor-based');
assert(!html.includes('data-filter="candidate-search"'), 'candidate server-side search is not admitted in B12 baseline');
assert(!html.includes('data-field="permission"'), 'recipient selector must not disclose Permissions');
assert(!html.includes('data-field="role"'), 'recipient selector must not disclose roles');
assert(html.includes('data-candidate-authority="discovery-not-authorization"'), 'candidate presence must not be treated as route-write authorization');

assert(html.includes('data-recipient-state="historical-ineligible"'), 'historical ineligible recipient treatment missing');
assert(html.includes('Destinatário não elegível · configuração anterior'), 'historical ineligible recipient must not expose an opaque ID');
assert(html.includes('data-failure="validation-rejected"'), 'save-time recipient validation rejection state missing');
assert(html.includes('data-failure="stale-route-write"'), 'stale route-write state missing');
assert(html.includes('data-knowledge-state="unavailable"'), 'route-read unavailable state missing');

for (const scenario of ['configured', 'unconfigured', 'historical-ineligible', 'validation-rejected', 'stale-route-write', 'request-unavailable', 'organization-switch', 'candidate-pagination']) {
  assert(html.includes(`value="${scenario}"`), `deterministic B12 scenario missing: ${scenario}`);
}

for (const fn of ['openRouteEditor', 'saveRoute', 'removeConfiguration', 'invalidateRoutingContext']) {
  assert(html.includes(`function ${fn}`), `B12 interaction function missing: ${fn}`);
}
assert(html.includes('data-responsive-law="mobile-stacked-route-editor"'), 'mobile B12 structural law missing');

console.log('d6_r_b12_status=CANDIDATE');
console.log('d6_r_b12_route_slots=10/10');
console.log('d6_r_b12_editor=INLINE_ROW');
console.log('d6_r_b12_permission=notifications.manage');
console.log('d6_r_b12_access_read=FORBIDDEN');
console.log('d6_r_b12_candidate_disclosure=MINIMAL');
console.log('d6_r_b12_configured_empty=FORBIDDEN');
console.log('d6_r_b12_wireframe=PASS');
