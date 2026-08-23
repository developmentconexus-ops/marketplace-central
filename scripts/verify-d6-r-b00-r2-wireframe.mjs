import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifact = resolve(root, 'qualification/d6-r2-wireframes/b00-r2-notifications.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifact), 'B00-R2 rendered artifact missing: qualification/d6-r2-wireframes/b00-r2-notifications.html');

const html = readFileSync(artifact, 'utf8');

assert(html.includes('data-p8-status="candidate"'), 'B00-R2 must identify itself as P8 CANDIDATE');
assert(!html.includes('data-p8-status="locked"'), 'B00-R2 must not claim operator LOCK before visual adjudication');

const expectedNavLabels = [
  'Visão geral',
  'Preparação',
  'Anúncios',
  'Preços',
  'Disponibilidade',
  'Visão operacional',
  'Vendas',
  'Expedição',
  'Pós-venda',
  'Performance',
  'Mercado',
  'Economia',
  'Trabalho',
  'Aprovações',
  'Configurações',
];
for (const label of expectedNavLabels) {
  assert(html.includes(`>${label}</button>`) || html.includes(`>${label}</a>`), `locked B00 IA label missing: ${label}`);
}
assert(!/class="[^"]*nav-item[^"]*"[^>]*>\s*Notificações\s*</u.test(html), 'Notifications must not become a global sidebar nav item');

assert(html.includes('id="notificationBell"'), 'notification bell missing');
assert(html.includes('data-operation="ListMyNotifications"'), 'bell/preview must trace ListMyNotifications');
assert(html.includes('data-unread-probe="archive_state=active;read_state=unread;limit=1"'), 'bell unread-presence probe is not the accepted bounded query');
assert(!html.includes('data-unread-count='), 'numeric unread-count contract is forbidden in B00-R2');
assert(!html.includes('data-required-permission="notifications.manage"'), 'self Inbox bell must not require notifications.manage');

for (const scenario of ['unread-present', 'unread-empty', 'knowledge-unavailable', 'organization-switch']) {
  assert(html.includes(`value="${scenario}"`), `deterministic B00-R2 scenario missing: ${scenario}`);
}

assert(html.includes('id="notificationPreview"'), 'bounded recent-Inbox preview missing');
assert(html.includes('data-surface="U01"'), 'preview must identify the U01 surface');
assert(html.includes('data-bounded="true"'), 'preview must remain bounded');
assert(html.includes('id="viewAllNotifications"'), 'Ver todas continuation missing');
assert(html.includes('data-target-block="B11"'), 'Ver todas must point only to later B11, not render it');

assert(html.includes('data-action="open-source"'), 'explicit source-open control missing');
assert(html.includes('data-action="mark-read"'), 'explicit awareness-mutation control missing');
assert(html.includes('data-operation="UpdateMyNotificationAwarenessState"'), 'mark-read must trace UpdateMyNotificationAwarenessState');

assert(html.includes('function closePreview'), 'preview close behavior missing');
assert(html.includes('function invalidateNotificationContext'), 'Organization notification-context invalidation behavior missing');
assert(html.includes('closePreview({ restoreFocus: false });'), 'Organization context invalidation must close the preview before new context is presented');

assert(html.includes('data-responsive-law="mobile-menu-title-bell-context-below"'), 'mobile menu/title/bell + stacked local-context law missing');
assert(html.includes('data-not-rendered="B11 B12"'), 'artifact must explicitly keep B11/B12 out of this P8 baseline');
assert(!html.includes('id="fullInbox"'), 'B11 full Inbox must not be rendered in B00-R2');
assert(!html.includes('id="notificationSettings"'), 'B12 routing Settings must not be rendered in B00-R2');

console.log('d6_r_b00_r2_status=CANDIDATE');
console.log('d6_r_b00_r2_global_sidebar_delta=0');
console.log('d6_r_b00_r2_numeric_unread_count=FORBIDDEN');
console.log('d6_r_b00_r2_b11_b12=NOT_RENDERED');
console.log('d6_r_b00_r2_wireframe=PASS');
