import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const temp = mkdtempSync(join(tmpdir(), 'mpc-authorization-request-oad-proof-'));
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const ACTIONABLE_IDS = ['ListMyActionableAuthorizationRequests', 'GetMyActionableAuthorizationRequest'];
const REVIEW_KINDS = ['listing_intent', 'price_intent', 'business_order_intent', 'invoicing_intent'];
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
function requestSchema(document, operation) {
  const schema = resolveRef(document, operation.requestBody)?.content?.['application/json']?.schema;
  return schema ? resolveRef(document, schema) : undefined;
}
function response(document, operation, status = '200') { return resolveRef(document, operation.responses?.[status]); }
function responseSchema(document, operation, status = '200') {
  const schema = response(document, operation, status)?.content?.['application/json']?.schema;
  return schema ? resolveRef(document, schema) : undefined;
}
function component(document, name) {
  const raw = document.components?.schemas?.[name];
  assert(raw, `schema not found: ${name}`);
  return resolveRef(document, raw);
}
function hasHeaderParameter(document, entry, name) {
  return parameters(document, entry).some((p) => p.in === 'header' && p.name === name && p.required === true);
}
function queryNames(document, entry) {
  return parameters(document, entry).filter((p) => p.in === 'query').map((p) => p.name);
}
function assertOperation(entry, { method, path, permission, principals = ['H'] }) {
  assert(entry.method === method, `${entry.operation.operationId} method must be ${method}`);
  assert(entry.path === path, `${entry.operation.operationId} path mismatch: ${entry.path}`);
  assert(entry.operation['x-mpc-semantic-owner'] === 'ControlledActionGovernance', `${entry.operation.operationId} owner must be ControlledActionGovernance`);
  assert(entry.operation['x-mpc-required-permission'] === permission, `${entry.operation.operationId} Permission mismatch`);
  sameSet(entry.operation['x-mpc-principal-kinds'] ?? [], principals, `${entry.operation.operationId} principals`);
}
function kindOf(document, schema) { return resolveRef(document, schema?.properties?.kind)?.const; }
function forbiddenGenericFields(schema, label) {
  for (const name of ['payload', 'metadata', 'data', 'attributes', 'entity_type', 'entity_id', 'fields', 'raw']) {
    assert(!Object.hasOwn(schema?.properties ?? {}, name), `${label} leaked generic field ${name}`);
  }
}
function reviewBranches(document) {
  const review = component(document, 'AuthorizationReviewBasis');
  assert((review.oneOf ?? []).length === 4, `AuthorizationReviewBasis must have 4 branches, found ${(review.oneOf ?? []).length}`);
  const branches = (review.oneOf ?? []).map((branch) => resolveRef(document, branch));
  sameSet(branches.map((branch) => kindOf(document, branch)), REVIEW_KINDS, 'AuthorizationReviewBasis kinds');
  for (const branch of branches) {
    assert(branch.additionalProperties === false, `review basis ${kindOf(document, branch)} must be closed`);
    forbiddenGenericFields(branch, `review basis ${kindOf(document, branch)}`);
  }
  const byKind = new Map(branches.map((branch) => [kindOf(document, branch), branch]));
  sameSet(byKind.get('listing_intent')?.required ?? [], ['kind', 'listing_intent_id', 'source_product', 'target', 'desired'], 'Listing review required');
  sameSet(byKind.get('price_intent')?.required ?? [], ['kind', 'price_intent_id', 'target', 'desired_price'], 'Price review required');
  sameSet(byKind.get('business_order_intent')?.required ?? [], ['kind', 'business_order_intent_id', 'sale_snapshot', 'target_source_instance_id', 'party_resolution', 'destination_realization'], 'BusinessOrder review required');
  sameSet(byKind.get('invoicing_intent')?.required ?? [], ['kind', 'invoicing_intent_id', 'sale_snapshot', 'business_order_snapshot'], 'Invoicing review required');
  return byKind;
}
function targetKind(document, target) { return kindOf(document, resolveRef(document, target)); }
function validateActionableView(document) {
  const view = component(document, 'ActionableAuthorizationRequestView');
  assert((view.oneOf ?? []).length === 4, `ActionableAuthorizationRequestView must have 4 variants, found ${(view.oneOf ?? []).length}`);
  const seen = [];
  for (const raw of view.oneOf ?? []) {
    const branch = resolveRef(document, raw);
    assert(branch.additionalProperties === false, 'actionable request variant must be closed');
    sameSet(branch.required ?? [], ['authorization_request_id', 'target', 'subject_display_label', 'review_basis', 'created_at'], 'actionable request required fields');
    const tKind = targetKind(document, branch.properties?.target);
    const rKind = kindOf(document, resolveRef(document, branch.properties?.review_basis));
    assert(tKind && rKind && tKind === rKind, `actionable target/review mismatch: ${tKind}/${rKind}`);
    seen.push(tKind);
  }
  sameSet(seen, REVIEW_KINDS, 'actionable request variants');
}
function validate(document) {
  const all = operations(document);
  assert(all.length === 106, `AuthorizationRequest D5-R6 Product operation count must be 106, found ${all.length}`);
  assert(new Set(all.map((entry) => entry.operation.operationId)).size === 106, 'operationId values must remain unique');

  const permissions = normalize(all.map((entry) => entry.operation['x-mpc-required-permission']).filter((value) => value && value !== 'authenticated'));
  assert(permissions.length === 31, `ordinary Permission vocabulary must remain 31, found ${permissions.length}`);
  assert(permissions.includes('governance.decide') && permissions.includes('notifications.manage'), 'existing Governance/Notification Permissions missing');
  assert(!permissions.includes('authorization_requests.read'), 'new AuthorizationRequest read Permission is forbidden');

  const list = byId(all, 'ListMyActionableAuthorizationRequests');
  const get = byId(all, 'GetMyActionableAuthorizationRequest');
  const decide = byId(all, 'CreateAuthorizationDecision');

  assertOperation(list, { method: 'get', path: '/organizations/{organization_id}/authorization-requests', permission: 'governance.decide' });
  assertOperation(get, { method: 'get', path: '/organizations/{organization_id}/authorization-requests/{authorization_request_id}', permission: 'governance.decide' });
  assertOperation(decide, { method: 'post', path: '/organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide', permission: 'governance.decide' });
  sameSet(queryNames(document, list), ['limit', 'cursor'], 'actionable AuthorizationRequest list query');
  assert(queryNames(document, get).length === 0, 'actionable AuthorizationRequest detail must not add query controls');

  const getResponse = response(document, get.operation, '200');
  assert(getResponse?.headers?.ETag, 'actionable AuthorizationRequest detail must return ETag');
  const detail = responseSchema(document, get.operation);
  const canonicalView = component(document, 'ActionableAuthorizationRequestView');
  assert(JSON.stringify(detail) === JSON.stringify(canonicalView), 'detail response must be ActionableAuthorizationRequestView');
  validateActionableView(document);

  const listSchema = responseSchema(document, list.operation);
  sameSet(Object.keys(listSchema?.properties ?? {}), ['authorization_requests', 'next_cursor'], 'actionable collection properties');
  const item = resolveRef(document, listSchema?.properties?.authorization_requests?.items);
  sameSet(item?.required ?? [], ['authorization_request_id', 'target', 'subject_display_label', 'created_at'], 'actionable list item required');
  assert(item?.additionalProperties === false, 'actionable list item must be closed');

  assert(hasHeaderParameter(document, decide, 'If-Match'), 'CreateAuthorizationDecision must require If-Match');
  assert(hasHeaderParameter(document, decide, 'Idempotency-Key'), 'CreateAuthorizationDecision must require Idempotency-Key');
  for (const status of ['201', '409', '412', '422', '428']) assert(decide.operation.responses?.[status], `CreateAuthorizationDecision must expose ${status}`);
  const decideBody = requestSchema(document, decide.operation);
  sameSet(decideBody?.required ?? [], ['outcome'], 'CreateAuthorizationDecision body required');
  sameSet(Object.keys(decideBody?.properties ?? {}), ['outcome'], 'CreateAuthorizationDecision body properties');
  assert(decideBody?.additionalProperties === false, 'CreateAuthorizationDecision body must be closed');
  assert(!JSON.stringify(decideBody).includes('etag') && !JSON.stringify(decideBody).includes('target'), 'decision body must not carry target/target ETag');

  const target = component(document, 'AuthorizationTargetRef');
  assert((target.oneOf ?? []).length === 4, 'AuthorizationTargetRef must remain closed 4-way union');
  sameSet((target.oneOf ?? []).map((branch) => targetKind(document, branch)), REVIEW_KINDS, 'AuthorizationTargetRef kinds');
  for (const branch of target.oneOf ?? []) {
    const resolved = resolveRef(document, branch);
    assert(!Object.hasOwn(resolved.properties ?? {}, 'etag'), `AuthorizationTargetRef ${kindOf(document, resolved)} must not contain target ETag`);
  }

  reviewBranches(document);

  const decision = component(document, 'AuthorizationDecision');
  for (const field of ['authorization_decision_id', 'authorization_request_id', 'target', 'review_basis', 'outcome', 'decided_by_principal_id', 'decided_at']) {
    assert((decision.required ?? []).includes(field), `AuthorizationDecision must require ${field}`);
  }
  assert(JSON.stringify(resolveRef(document, decision.properties?.target)) === JSON.stringify(target), 'AuthorizationDecision target must use no-ETag AuthorizationTargetRef');
  assert(JSON.stringify(resolveRef(document, decision.properties?.review_basis)) === JSON.stringify(component(document, 'AuthorizationReviewBasis')), 'AuthorizationDecision must preserve immutable review basis');

  const notification = component(document, 'Notification');
  const notificationBranches = (notification.oneOf ?? []).map((branch) => resolveRef(document, branch));
  const f13 = notificationBranches.find((branch) => kindOf(document, branch) === 'AUTHORIZATION_ACTION_REQUIRED');
  const f14 = notificationBranches.find((branch) => kindOf(document, branch) === 'AUTHORIZATION_DECISION_RESULT');
  assert(f13 && f14, 'F13/F14 Notification branches missing');
  const f13Source = resolveRef(document, f13.properties?.source_ref);
  assert(kindOf(document, f13Source) === 'authorization_request', 'F13 source must be AuthorizationRequestRef');
  assert(Object.hasOwn(f13Source.properties ?? {}, 'authorization_request_id'), 'F13 AuthorizationRequestRef ID missing');
  assert(!JSON.stringify(f13Source).includes('price_intent_id') && !JSON.stringify(f13Source).includes('listing_intent_id'), 'F13 must not fall back to target ref');
  const f14Source = resolveRef(document, f14.properties?.source_ref);
  sameSet((f14Source.oneOf ?? []).map((branch) => targetKind(document, branch)), REVIEW_KINDS, 'F14 target-oriented continuation');

  const workOrigin = component(document, 'WorkOrigin');
  const workKinds = (workOrigin.oneOf ?? []).map((branch) => kindOf(document, resolveRef(document, branch))).filter(Boolean);
  assert(workKinds.includes('authorization_request'), 'WorkOrigin must admit authorization_request');

  for (const forbiddenId of ['CreateAuthorizationRequest', 'InvalidateAuthorizationRequest', 'ReauthorizeAuthorizationRequest', 'ListAllAuthorizationRequests', 'ResolveAuthorizationRequestRecipients', 'CreateNoApproverWork']) {
    assert(!all.some((entry) => entry.operation.operationId === forbiddenId), `forbidden public workflow operation leaked: ${forbiddenId}`);
  }

  return { all, permissions };
}
function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `AuthorizationRequest negative control unexpectedly passed: ${name}`);
  negativeControls++;
}
function negativeProof(document) {
  expectFailure('actionable read widens to governance.read', document, (candidate) => {
    byId(operations(candidate), 'ListMyActionableAuthorizationRequests').operation['x-mpc-required-permission'] = 'governance.read';
  });
  expectFailure('actionable read admits automation', document, (candidate) => {
    byId(operations(candidate), 'GetMyActionableAuthorizationRequest').operation['x-mpc-principal-kinds'] = ['H', 'A'];
  });
  expectFailure('actionable list gains search', document, (candidate) => {
    byId(operations(candidate), 'ListMyActionableAuthorizationRequests').operation.parameters.push({ name: 'query', in: 'query', required: false, schema: { type: 'string' } });
  });
  expectFailure('decision loses request precondition', document, (candidate) => {
    const entry = byId(operations(candidate), 'CreateAuthorizationDecision');
    entry.operation.parameters = entry.operation.parameters.filter((p) => resolveRef(candidate, p)?.name !== 'If-Match');
  });
  expectFailure('decision body regains target', document, (candidate) => {
    const body = requestSchema(candidate, byId(operations(candidate), 'CreateAuthorizationDecision').operation);
    body.properties.target = { $ref: '#/components/schemas/AuthorizationTargetRef' };
  });
  expectFailure('target ETag returns', document, (candidate) => {
    const target = component(candidate, 'AuthorizationTargetRef');
    const branch = resolveRef(candidate, target.oneOf[0]);
    branch.properties.etag = { type: 'string' };
  });
  expectFailure('review basis becomes generic payload', document, (candidate) => {
    const review = component(candidate, 'AuthorizationReviewBasis');
    resolveRef(candidate, review.oneOf[0]).properties.payload = { type: 'object' };
  });
  expectFailure('F13 falls back to target', document, (candidate) => {
    const notification = component(candidate, 'Notification');
    const f13 = notification.oneOf.map((b) => resolveRef(candidate, b)).find((b) => kindOf(candidate, b) === 'AUTHORIZATION_ACTION_REQUIRED');
    f13.properties.source_ref = { $ref: '#/components/schemas/AuthorizationTargetRef' };
  });
  expectFailure('Work loses zero-decider request origin', document, (candidate) => {
    const work = component(candidate, 'WorkOrigin');
    work.oneOf = work.oneOf.filter((b) => kindOf(candidate, resolveRef(candidate, b)) !== 'authorization_request');
  });
  expectFailure('public AuthorizationRequest create leaks', document, (candidate) => {
    candidate.paths['/organizations/{organization_id}/authorization-requests'].post = {
      operationId: 'CreateAuthorizationRequest',
      'x-mpc-operation-class': 'C',
      'x-mpc-required-permission': 'governance.decide',
      'x-mpc-principal-kinds': ['H'],
      'x-mpc-semantic-owner': 'ControlledActionGovernance',
      responses: { '201': { description: 'created' } },
    };
  });
  assert(negativeControls === 10, `AuthorizationRequest negative-control count must be 10, found ${negativeControls}`);
}

try {
  const bundle = join(temp, 'product.json');
  const spec = npxSpec(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundle]);
  run(spec.command, spec.args);
  const document = JSON.parse(readFileSync(bundle, 'utf8'));
  const result = validate(document);
  negativeProof(document);
  console.log(`authorization_request_oad_operations=${result.all.length}/106`);
  console.log(`authorization_request_oad_permissions=${result.permissions.length}/31`);
  console.log(`authorization_request_oad_review_basis=4/4`);
  console.log(`authorization_request_oad_negative_controls=${negativeControls}/10`);
  console.log('authorization_request_oad=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
