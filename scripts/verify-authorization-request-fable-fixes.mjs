import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const p9Path = join(root, 'docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md');
const temp = mkdtempSync(join(tmpdir(), 'mpc-authorization-request-fable-proof-'));
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const REVIEW_KINDS = ['listing_intent', 'price_intent', 'business_order_intent', 'invoicing_intent'];
const EXPECTED_P9_OPS = [
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
const VALIDITY_UNAVAILABLE_TYPE = 'https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable';
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
function component(document, name) {
  const raw = document.components?.schemas?.[name];
  assert(raw, `schema not found: ${name}`);
  return resolveRef(document, raw);
}
function operations(document) {
  const result = [];
  for (const rawPathItem of Object.values(document.paths ?? {})) {
    const pathItem = resolveRef(document, rawPathItem);
    for (const [method, operation] of Object.entries(pathItem ?? {})) {
      if (HTTP_METHODS.has(method)) result.push(operation);
    }
  }
  return result;
}
function kindOf(document, schema) {
  return resolveRef(document, resolveRef(document, schema)?.properties?.kind)?.const;
}
function assertDecisionPairing(document) {
  const decision = component(document, 'AuthorizationDecision');
  assert((decision.oneOf ?? []).length === 4, `AuthorizationDecision pairing must have 4 branches, found ${(decision.oneOf ?? []).length}`);
  const kinds = [];
  for (const raw of decision.oneOf ?? []) {
    const branch = resolveRef(document, raw);
    const targetKind = kindOf(document, branch.properties?.target);
    const reviewKind = kindOf(document, branch.properties?.review_basis);
    assert(targetKind && reviewKind && targetKind === reviewKind, `AuthorizationDecision target/review mismatch: ${targetKind}/${reviewKind}`);
    kinds.push(targetKind);
  }
  sameSet(kinds, REVIEW_KINDS, 'AuthorizationDecision pairing kinds');
}
function actionableSchemas(document) {
  const result = [component(document, 'ActionableAuthorizationRequestListItem')];
  const view = component(document, 'ActionableAuthorizationRequestView');
  for (const raw of view.oneOf ?? []) result.push(resolveRef(document, raw));
  return result;
}
function assertNoDeadActionableCorrelationFields(document) {
  for (const schema of actionableSchemas(document)) {
    for (const field of ['requester_or_initiator_principal_id', 'predecessor_authorization_request_id']) {
      assert(!Object.hasOwn(schema.properties ?? {}, field), `actionable AuthorizationRequest wire leaked no-consumer field: ${field}`);
    }
  }
}
function decideOperation(document) {
  return operations(document).find((operation) => operation.operationId === 'CreateAuthorizationDecision');
}
function validityUnavailableSchema(document) {
  const decide = decideOperation(document);
  assert(decide, 'CreateAuthorizationDecision missing');
  const response = resolveRef(document, decide.responses?.['503']);
  const schema = response?.content?.['application/problem+json']?.schema;
  assert(schema, 'CreateAuthorizationDecision 503 Problem schema missing');
  return resolveRef(document, schema);
}
function assertValidityUnavailableDiscriminator(document) {
  const schema = validityUnavailableSchema(document);
  const type = resolveRef(document, schema.properties?.type)?.const;
  const status = resolveRef(document, schema.properties?.status)?.const;
  assert(type === VALIDITY_UNAVAILABLE_TYPE, `authorization validity 503 must use specific Problem type, found ${type}`);
  assert(status === 503, `authorization validity Problem status must be 503, found ${status}`);
}
function p9OperationHomes() {
  const doc = readFileSync(p9Path, 'utf8');
  const homes = [...doc.matchAll(/P9_OP_HOME:([A-Za-z0-9]+)/g)].map((match) => match[1]);
  sameSet(homes, EXPECTED_P9_OPS, 'P9 operation homes');
  return homes;
}
function assertP9HomesExistInOad(document, homes = p9OperationHomes()) {
  const actual = new Set(operations(document).map((operation) => operation.operationId));
  for (const home of homes) assert(actual.has(home), `P9 operation home does not exist in canonical OAD: ${home}`);
}
function expectFailure(name, fn) {
  let failed = false;
  try { fn(); } catch { failed = true; }
  assert(failed, `Fable bounded-fix negative control unexpectedly passed: ${name}`);
  negativeControls++;
}
function negativeProof(document) {
  expectFailure('Decision pairing can be stripped', () => {
    const candidate = clone(document);
    component(candidate, 'AuthorizationDecision').oneOf = [];
    assertDecisionPairing(candidate);
  });
  expectFailure('dead actionable correlation field can return', () => {
    const candidate = clone(document);
    component(candidate, 'ActionableAuthorizationRequestListItem').properties.requester_or_initiator_principal_id = { type: 'string' };
    assertNoDeadActionableCorrelationFields(candidate);
  });
  expectFailure('semantic validity 503 can fall back to about:blank', () => {
    const candidate = clone(document);
    validityUnavailableSchema(candidate).properties.type = { const: 'about:blank' };
    assertValidityUnavailableDiscriminator(candidate);
  });
  expectFailure('P9 operation home can drift from canonical OAD', () => {
    assertP9HomesExistInOad(document, [...p9OperationHomes(), 'DefinitelyNotAProductOperation']);
  });
  assert(negativeControls === 4, `Fable bounded-fix negative-control count must be 4, found ${negativeControls}`);
}

try {
  const bundle = join(temp, 'product.json');
  const spec = npxSpec(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundle]);
  run(spec.command, spec.args);
  const document = JSON.parse(readFileSync(bundle, 'utf8'));
  assertDecisionPairing(document);
  assertNoDeadActionableCorrelationFields(document);
  assertValidityUnavailableDiscriminator(document);
  assertP9HomesExistInOad(document);
  negativeProof(document);
  console.log('authorization_request_fable_decision_pairing=PASS');
  console.log('authorization_request_fable_actionable_projection=MINIMAL');
  console.log('authorization_request_fable_503_discriminator=SPECIFIC_PROBLEM_TYPE');
  console.log('authorization_request_fable_p9_wire_crosscheck=10/10');
  console.log(`authorization_request_fable_negative_controls=${negativeControls}/4`);
  console.log('authorization_request_fable_bounded_fixes=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
