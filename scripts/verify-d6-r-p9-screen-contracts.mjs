import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifact = resolve(root, 'docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md');
const repairArtifact = resolve(root, 'docs/engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifact), 'P9 Screen Contracts artifact missing: docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md');
assert(existsSync(repairArtifact), 'D5-R7 W1 carrier repair missing');
const doc = readFileSync(artifact, 'utf8');
const repair = readFileSync(repairArtifact, 'utf8');

assert(doc.includes('**Status:** PROVED / CANONICAL — FINAL P9 BIDIRECTIONAL TRACE PASS'), 'P9 must declare the proved final trace status');
assert(doc.includes('106 Product operations'), 'P9 must bind to canonical 106-operation Product authority');
assert(doc.includes('31 ordinary Permissions'), 'P9 must bind to canonical 31-Permission authority');
assert(doc.includes('query-only remedy remains **SUPERSEDED**'), 'P9 must preserve the old query-only remedy supersession');
assert(doc.includes('P9-F1 factual falsifier = RESOLVED'), 'P9 must explicitly close the original actionable-Governance falsifier');

const expectedSurfaces = ['G00-E', 'U01', 'R128', 'R129', 'R110', 'R111'];
const foundSurfaces = [...doc.matchAll(/P9_SURFACE:([A-Z0-9-]+)/g)].map(match => match[1]);
assert(new Set(foundSurfaces).size === expectedSurfaces.length, `P9 surface marker count mismatch: ${new Set(foundSurfaces).size}/${expectedSurfaces.length}`);
for (const surface of expectedSurfaces) assert(foundSurfaces.includes(surface), `P9 surface trace missing: ${surface}`);

const expectedOps = [
  'ListMyNotifications',
  'UpdateMyNotificationAwarenessState',
  'ListNotificationRoutes',
  'ListNotificationRouteRecipientCandidates',
  'SetNotificationRoute',
  'ListMyActionableAuthorizationRequests',
  'GetMyActionableAuthorizationRequest',
  'CreateAuthorizationDecision',
  'ListAuthorizationDecisions',
  'GetAuthorizationDecision',
];
const foundOps = [...doc.matchAll(/P9_OP_HOME:([A-Za-z0-9]+)/g)].map(match => match[1]);
assert(new Set(foundOps).size === expectedOps.length, `P9 human-operation home count mismatch: ${new Set(foundOps).size}/${expectedOps.length}`);
for (const op of expectedOps) assert(foundOps.includes(op), `P9 human-operation home missing: ${op}`);
assert(!foundOps.includes('ListMyActionableAuthorizations'), 'superseded query-only operation must not have a frontend home');

for (const basis of ['listing_intent', 'price_intent', 'business_order_intent', 'invoicing_intent']) {
  assert(doc.includes(`P9_REVIEW_BASIS:${basis}`), `P9 typed review basis missing: ${basis}`);
}
for (const stateClass of ['SERVER_STATE', 'URL_NAVIGATION_STATE', 'LOCAL_EPHEMERAL', 'GLOBAL_WORKSPACE_CONTEXT']) {
  assert(doc.includes(`P9_STATE_CLASS:${stateClass}`), `P9 client-state class missing: ${stateClass}`);
}

assert(doc.includes('P9_F13:AUTHORIZATION_REQUEST_REF'), 'F13 must trace through AuthorizationRequestRef');
assert(doc.includes('P9_F13_NOTIFICATION_AUTHORITY:AWARENESS_NOT_CAPABILITY'), 'F13 Notification must remain awareness, not capability');
assert(doc.includes('P9_F14:TARGET_ORIENTED'), 'F14 must remain target-oriented');
assert(doc.includes('P9_ZERO_DECIDER:WORK_NOT_AUTHORITY'), 'zero-decider condition must dispose to Work without authority transfer');
assert(doc.includes('P9_WORK_HOME:ListWork+GetWork'), 'zero-decider Work disposition must use existing Work reads');
assert(doc.includes('P9_DECISION_CONCURRENCY:IF_MATCH_AUTHORIZATION_REQUEST'), 'historical P9 snapshot must retain the carrier that D5-R7 explicitly supersedes');
assert(repair.includes('D5R7_W1_CARRIER:TYPED_REQUEST_ETAG'), 'current P9 composition must consume D5-R7 typed request etag');
assert(repair.includes('D5R7_SUPERSEDES:D5R6_P9_D7R_D8R_CARRIER_ONLY'), 'D5-R7 must explicitly supersede only P9 carrier-specific meaning');
assert(repair.includes('D5R7_B110:REVALIDATED_STRUCTURE_UNAFFECTED'), 'B110 LOCK must be revalidated after carrier correction');
assert(doc.includes('P9_DECISION_RETRY:IDEMPOTENCY_KEY_SAME_SEMANTIC_ATTEMPT'), 'decision ambiguous retry must preserve Idempotency-Key semantics');
assert(doc.includes('P9_DECISION_AUTO_RETRY:FORBIDDEN'), 'consequential Decision auto-retry must be forbidden');
assert(doc.includes('P9_SOURCE_OPEN:REAUTHORIZE_CURRENT_OWNER'), 'source continuation must reauthorize current owner');
assert(doc.includes('P9_PERMISSION_LAW:GOVERNANCE_DECIDE_READ_INDEPENDENT'), 'governance.decide and governance.read must remain independent');
assert(doc.includes('P9_ORG_SWITCH:INVALIDATE_SCOPED_STATE'), 'Organization switch invalidation law missing');

const forbiddenControls = [
  'TARGET_ETAG_AS_DECISION_AUTHORITY',
  'GENERIC_REVIEW_PAYLOAD',
  'NOTIFICATION_CAPABILITY_TOKEN',
  'WORK_FALLBACK_APPROVER',
  'F14_AUTHORIZATION_DECISION_ID',
  'SOURCE_OPEN_MARKS_NOTIFICATION_READ',
  'INBOX_TOTAL_OR_UNREAD_COUNT',
  'INBOX_BULK_OR_SEARCH',
  'CANDIDATE_PRESENCE_IS_AUTHORIZATION',
  'CONFIGURED_EMPTY_NOTIFICATION_ROUTE',
  'SECOND_GLOBAL_APPROVAL_HISTORY_DESTINATION',
  'SCREEN_SHAPED_PARALLEL_OWNER_AUTHORITY',
];
for (const control of forbiddenControls) {
  assert(doc.includes(`P9_FORBID:${control}`), `P9 negative control missing: ${control}`);
}

assert(doc.includes('P9_FABLE_GATE:REQUIRED_AFTER_P9'), 'Fable review/adjudication gate must remain explicit');
assert(doc.includes('D7-R remains BLOCKED'), 'historical P9 snapshot must preserve its original pre-Fable gate evidence');

console.log('d6_r_p9_status=PROVED');
console.log('d6_r_p9_surfaces=6/6');
console.log('d6_r_p9_human_ops=10/10');
console.log('d6_r_p9_review_basis=4/4');
console.log('d6_r_p9_decision_concurrency=TYPED_REQUEST_ETAG_VIA_D5R7');
console.log('d6_r_p9_f13=AUTHORIZATION_REQUEST_REF');
console.log('d6_r_p9_f14=TARGET_ORIENTED');
console.log('d6_r_p9_zero_decider=WORK_NOT_AUTHORITY');
console.log('d6_r_p9_negative_controls=12/12');
console.log('d6_r_p9_bidirectional_trace=PASS');
