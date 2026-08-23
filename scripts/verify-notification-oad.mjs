import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const temp = mkdtempSync(join(tmpdir(), 'mpc-notification-oad-proof-'));
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const NOTIFICATION_KINDS = [
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
const ORG_ROUTED_KINDS = [
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
const NOTIFICATION_IDS = [
  'ListMyNotifications',
  'UpdateMyNotificationAwarenessState',
  'ListNotificationRoutes',
  'ListNotificationRouteRecipientCandidates',
  'SetNotificationRoute',
];
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function clone(value) { return structuredClone(value); }
function normalize(values) { return [...new Set(values)].sort(); }
function sameSet(actual, expected, label) {
  const a = normalize(actual);
  const e = normalize(expected);
  assert(JSON.stringify(a) === JSON.stringify(e), `${label} mismatch\nactual=${JSON.stringify(a)}\nexpected=${JSON.stringify(e)}`);
}
function run(command, args) {
  const result = spawnSync(command, args, { cwd: root, encoding: 'utf8', shell: false });
  if (result.error) fail(`${command} failed to start: ${result.error.message}`);
  if (result.status !== 0) fail([`${command} ${args.join(' ')} failed with exit ${result.status}`, result.stdout?.trim(), result.stderr?.trim()].filter(Boolean).join('\n'));
  return result;
}
function npxSpec(args) {
  const npxCli = process.env.npm_execpath
    ? resolve(dirname(process.env.npm_execpath), 'npx-cli.js')
    : resolve(dirname(process.execPath), 'node_modules/npm/bin/npx-cli.js');
  return process.platform === 'win32'
    ? { command: process.execPath, args: [npxCli, ...args] }
    : { command: 'npx', args };
}
function resolveRef(document, value, seen = new Set()) {
  if (!value || typeof value !== 'object' || typeof value.$ref !== 'string') return value;
  const ref = value.$ref;
  assert(ref.startsWith('#/'), `bundle contains non-local ref: ${ref}`);
  assert(!seen.has(ref), `cyclic ref reached proof resolver: ${ref}`);
  let current = document;
  for (const raw of ref.slice(2).split('/')) {
    const key = raw.replaceAll('~1', '/').replaceAll('~0', '~');
    assert(current && Object.hasOwn(current, key), `unresolved bundled ref: ${ref}`);
    current = current[key];
  }
  return resolveRef(document, current, new Set(seen).add(ref));
}
function operations(document) {
  const result = [];
  for (const [path, rawPathItem] of Object.entries(document.paths ?? {})) {
    const pathItem = resolveRef(document, rawPathItem);
    for (const [method, operation] of Object.entries(pathItem ?? {})) {
      if (HTTP_METHODS.has(method)) result.push({ path, method, operation, pathItem });
    }
  }
  return result;
}
function byId(all, id) {
  const found = all.find((entry) => entry.operation.operationId === id);
  assert(found, `operation not found: ${id}`);
  return found;
}
function parameters(document, entry) {
  return [...(entry.pathItem.parameters ?? []), ...(entry.operation.parameters ?? [])].map((value) => resolveRef(document, value));
}
function responseSchema(document, operation, status = '200') {
  const response = resolveRef(document, operation.responses?.[status]);
  const schema = response?.content?.['application/json']?.schema;
  return schema ? resolveRef(document, schema) : undefined;
}
function schemaBySuffix(document, suffix) {
  const matches = Object.entries(document.components?.schemas ?? {}).filter(([name]) => name === suffix || name.endsWith(`_${suffix}`) || name.endsWith(`-${suffix}`));
  assert(matches.length === 1, `expected exactly one bundled schema ending ${suffix}, found ${matches.map(([name]) => name).join(', ')}`);
  return resolveRef(document, matches[0][1]);
}
function effectiveAccessPermissions(document) {
  const access = responseSchema(document, byId(operations(document), 'GetCurrentAccessContext').operation);
  const org = resolveRef(document, access?.properties?.organizations?.items);
  const permission = resolveRef(document, org?.properties?.permissions?.items);
  assert(Array.isArray(permission?.enum), 'effective access Permission enum missing');
  return permission.enum;
}
function queryMap(document, entry) {
  return new Map(parameters(document, entry).filter((p) => p.in === 'query').map((p) => [p.name, p]));
}
function assertMethodAndPath(entry, method, path) {
  assert(entry.method === method && entry.path === path, `${entry.operation.operationId} wire mismatch: ${entry.method.toUpperCase()} ${entry.path}`);
}
function assertHOnly(entry) { sameSet(entry.operation['x-mpc-principal-kinds'] ?? [], ['H'], `${entry.operation.operationId} principal kinds`); }
function hasHeaderParameter(document, entry, name) { return parameters(document, entry).some((p) => p.in === 'header' && p.name === name && p.required === true); }
function branchByKind(document, schema, kind) {
  const found = (schema.oneOf ?? []).map((branch) => resolveRef(document, branch)).find((branch) => resolveRef(document, branch.properties?.kind)?.const === kind);
  assert(found, `Notification branch missing kind ${kind}`);
  return found;
}

function validate(document) {
  const all = operations(document);
  for (const id of NOTIFICATION_IDS) byId(all, id);
  assert(all.length === 104, `Product operation count must be 104 after NOTIF-01, found ${all.length}`);
  assert(new Set(all.map((entry) => entry.operation.operationId)).size === 104, 'operationId values are not unique');

  const listMine = byId(all, 'ListMyNotifications');
  const updateMine = byId(all, 'UpdateMyNotificationAwarenessState');
  const listRoutes = byId(all, 'ListNotificationRoutes');
  const candidates = byId(all, 'ListNotificationRouteRecipientCandidates');
  const setRoute = byId(all, 'SetNotificationRoute');

  assertMethodAndPath(listMine, 'get', '/organizations/{organization_id}/notifications');
  assertMethodAndPath(updateMine, 'patch', '/organizations/{organization_id}/notifications/{notification_id}');
  assertMethodAndPath(listRoutes, 'get', '/organizations/{organization_id}/notification-routes');
  assertMethodAndPath(candidates, 'get', '/organizations/{organization_id}/notification-route-recipient-candidates');
  assertMethodAndPath(setRoute, 'patch', '/organizations/{organization_id}/notification-routes/{notification_kind}');

  for (const entry of [listMine, updateMine, listRoutes, candidates, setRoute]) assertHOnly(entry);
  assert(listMine.operation['x-mpc-semantic-owner'] === 'PersonalNotifications', 'ListMyNotifications owner must be PersonalNotifications');
  assert(updateMine.operation['x-mpc-semantic-owner'] === 'PersonalNotifications', 'UpdateMyNotificationAwarenessState owner must be PersonalNotifications');
  assert(listRoutes.operation['x-mpc-semantic-owner'] === 'PersonalNotifications', 'ListNotificationRoutes owner must be PersonalNotifications');
  assert(setRoute.operation['x-mpc-semantic-owner'] === 'PersonalNotifications', 'SetNotificationRoute owner must be PersonalNotifications');
  assert(candidates.operation['x-mpc-semantic-owner'] === 'IdentityAccess', 'recipient candidate owner must be IdentityAccess');

  for (const entry of [listMine, updateMine]) assert(entry.operation['x-mpc-required-permission'] === 'authenticated', `${entry.operation.operationId} must have no ordinary Permission`);
  for (const entry of [listRoutes, candidates, setRoute]) assert(entry.operation['x-mpc-required-permission'] === 'notifications.manage', `${entry.operation.operationId} must require notifications.manage`);

  const ordinaryPermissions = normalize(all.map((entry) => entry.operation['x-mpc-required-permission']).filter((value) => value !== 'authenticated'));
  assert(ordinaryPermissions.length === 31, `ordinary Permission vocabulary must be 31, found ${ordinaryPermissions.length}`);
  assert(ordinaryPermissions.includes('notifications.manage'), 'notifications.manage missing from operation Permission vocabulary');
  assert(!ordinaryPermissions.includes('notifications.read'), 'notifications.read must not exist');
  const effectivePermissions = effectiveAccessPermissions(document);
  assert(effectivePermissions.length === 31 && effectivePermissions.includes('notifications.manage') && !effectivePermissions.includes('notifications.read'), 'effective access Permission vocabulary must be exact 31 with notifications.manage only');

  const inboxQuery = queryMap(document, listMine);
  sameSet([...inboxQuery.keys()], ['archive_state', 'read_state', 'notification_kind', 'limit', 'cursor'], 'ListMyNotifications query controls');
  const kindFilter = resolveRef(document, inboxQuery.get('notification_kind')?.schema);
  assert(kindFilter?.type === 'array' && kindFilter.uniqueItems === true && kindFilter.minItems === 1, 'notification_kind filter must be a non-empty unique array');
  sameSet(resolveRef(document, kindFilter.items)?.enum ?? [], NOTIFICATION_KINDS, 'notification_kind filter values');

  const candidateQuery = queryMap(document, candidates);
  sameSet([...candidateQuery.keys()], ['limit', 'cursor'], 'recipient candidate query controls');
  const candidateCollection = responseSchema(document, candidates.operation);
  sameSet(Object.keys(candidateCollection?.properties ?? {}), ['recipient_candidates', 'next_cursor'], 'recipient candidate collection properties');
  const candidate = resolveRef(document, candidateCollection?.properties?.recipient_candidates?.items);
  sameSet(candidate?.required ?? [], ['principal_id', 'display_name'], 'recipient candidate required fields');
  sameSet(Object.keys(candidate?.properties ?? {}), ['principal_id', 'display_name'], 'recipient candidate projection');
  assert(candidate?.additionalProperties === false, 'recipient candidate projection must be closed');

  assert(hasHeaderParameter(document, updateMine, 'If-Match'), 'UpdateMyNotificationAwarenessState must require If-Match');
  assert(!hasHeaderParameter(document, updateMine, 'Idempotency-Key'), 'UpdateMyNotificationAwarenessState must not use Idempotency-Key');
  assert(updateMine.operation.responses?.['412'] && updateMine.operation.responses?.['428'], 'UpdateMyNotificationAwarenessState must expose 412/428 stale-write semantics');
  assert(hasHeaderParameter(document, setRoute, 'If-Match'), 'SetNotificationRoute must require If-Match');
  assert(!hasHeaderParameter(document, setRoute, 'Idempotency-Key'), 'SetNotificationRoute must not use Idempotency-Key');
  assert(setRoute.operation.responses?.['412'] && setRoute.operation.responses?.['428'], 'SetNotificationRoute must expose 412/428 stale-write semantics');

  const orgRouted = schemaBySuffix(document, 'OrgRoutedNotificationKind');
  sameSet(orgRouted.enum ?? [], ORG_ROUTED_KINDS, 'ORG_ROUTED NotificationKind values');
  const allKinds = schemaBySuffix(document, 'NotificationKind');
  sameSet(allKinds.enum ?? [], NOTIFICATION_KINDS, 'NotificationKind values');

  const notification = schemaBySuffix(document, 'Notification');
  assert((notification.oneOf ?? []).length === 14, `Notification must have 14 kind-constrained branches, found ${(notification.oneOf ?? []).length}`);
  const branchKinds = (notification.oneOf ?? []).map((branch) => resolveRef(document, resolveRef(document, branch).properties?.kind)?.const);
  sameSet(branchKinds, NOTIFICATION_KINDS, 'Notification branch kinds');
  const f02 = branchByKind(document, notification, 'OFFERING_ASYNC_ACTION_RESULT');
  const f14 = branchByKind(document, notification, 'AUTHORIZATION_DECISION_RESULT');
  assert((f02.required ?? []).includes('offering_async_result_outcome'), 'F02 must require offering_async_result_outcome');
  sameSet(resolveRef(document, f02.properties?.offering_async_result_outcome)?.enum ?? [], ['converged', 'rejected', 'ambiguous', 'divergent'], 'F02 result values');
  assert((f14.required ?? []).includes('authorization_decision_outcome'), 'F14 must require authorization_decision_outcome');
  sameSet(resolveRef(document, f14.properties?.authorization_decision_outcome)?.enum ?? [], ['authorize', 'reject'], 'F14 result values');
  for (const forbidden of ['result', 'status', 'reason', 'summary', 'payload', 'metadata', 'template_variables']) {
    for (const branch of notification.oneOf ?? []) assert(!Object.hasOwn(resolveRef(document, branch).properties ?? {}, forbidden), `Notification leaked generic field ${forbidden}`);
  }
  const f14Source = resolveRef(document, f14.properties?.source_ref);
  const f14SourceKinds = (f14Source.oneOf ?? []).map((branch) => resolveRef(document, resolveRef(document, branch).properties?.kind)?.const).filter(Boolean);
  sameSet(f14SourceKinds, ['listing_intent', 'price_intent', 'business_order_intent', 'invoicing_intent'], 'F14 AuthorizationTargetRef variants');
  assert(!JSON.stringify(f14Source).includes('authorization_decision_id'), 'F14 source_ref must not retain AuthorizationDecisionRef');

  const routeKindParam = parameters(document, setRoute).find((p) => p.in === 'path' && p.name === 'notification_kind');
  sameSet(resolveRef(document, routeKindParam?.schema)?.enum ?? [], ORG_ROUTED_KINDS, 'SetNotificationRoute path kind values');
  return { all, ordinaryPermissions };
}

function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `NOTIF-01 negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

function negativeProof(document) {
  expectFailure('candidate over-discloses role keys', document, (candidate) => {
    const collection = responseSchema(candidate, byId(operations(candidate), 'ListNotificationRouteRecipientCandidates').operation);
    resolveRef(candidate, collection.properties.recipient_candidates.items).properties.role_keys = { type: 'array', items: { type: 'string' } };
  });
  expectFailure('candidate requires access.read', document, (candidate) => {
    byId(operations(candidate), 'ListNotificationRouteRecipientCandidates').operation['x-mpc-required-permission'] = 'access.read';
  });
  expectFailure('routing write loses stale-write precondition', document, (candidate) => {
    const entry = byId(operations(candidate), 'SetNotificationRoute');
    const index = entry.operation.parameters.findIndex((p) => resolveRef(candidate, p)?.name === 'If-Match');
    entry.operation.parameters.splice(index, 1);
  });
  expectFailure('self mutation gains idempotency key', document, (candidate) => {
    byId(operations(candidate), 'UpdateMyNotificationAwarenessState').operation.parameters.push({ name: 'Idempotency-Key', in: 'header', required: true, schema: { type: 'string' } });
  });
  expectFailure('notifications.read appears', document, (candidate) => {
    const all = operations(candidate);
    byId(all, 'ListMyNotifications').operation['x-mpc-required-permission'] = 'notifications.read';
  });
  expectFailure('F14 points to decision identity', document, (candidate) => {
    const notification = schemaBySuffix(candidate, 'Notification');
    const f14 = branchByKind(candidate, notification, 'AUTHORIZATION_DECISION_RESULT');
    f14.properties.source_ref = { type: 'object', properties: { authorization_decision_id: { type: 'string' } } };
  });
  assert(negativeControls === 6, `NOTIF-01 negative-control count must be 6, found ${negativeControls}`);
}

try {
  const bundle = join(temp, 'product.json');
  const spec = npxSpec(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundle]);
  run(spec.command, spec.args);
  const document = JSON.parse(readFileSync(bundle, 'utf8'));
  const result = validate(document);
  negativeProof(document);
  console.log(`notification_oad_operations=${result.all.length}/104`);
  console.log(`notification_oad_permissions=${result.ordinaryPermissions.length}/31`);
  console.log(`notification_oad_negative_controls=${negativeControls}/6`);
  console.log('notification_oad=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
