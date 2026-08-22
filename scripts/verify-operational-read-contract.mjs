import { mkdirSync, mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const temp = mkdtempSync(join(tmpdir(), 'mpc-operational-read-'));
const httpMethods = new Set(['get', 'put', 'post', 'delete', 'patch', 'head', 'options', 'trace']);
const effectStates = ['not_attempted', 'pending', 'accepted', 'rejected', 'ambiguous'];
const convergenceStates = ['pending', 'converged', 'divergent', 'unknown', 'not_applicable'];
const checkpointStates = ['pending', 'recorded'];
const shipmentStates = ['pending', 'ready', 'dispatched', 'delivered', 'exception', 'unknown'];
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function normalize(values) { return [...new Set(values)].sort(); }
function sameSet(actual, expected, label) {
  const a = normalize(actual);
  const e = normalize(expected);
  assert(JSON.stringify(a) === JSON.stringify(e), `${label} mismatch\nactual=${JSON.stringify(a)}\nexpected=${JSON.stringify(e)}`);
}
function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: options.cwd ?? root,
    env: { ...process.env, ...(options.env ?? {}) },
    encoding: 'utf8',
    shell: false,
  });
  if (result.error) fail(`${command} failed to start: ${result.error.message}`);
  if (result.status !== 0) {
    fail([`${command} ${args.join(' ')} failed with exit ${result.status}`, result.stdout?.trim(), result.stderr?.trim()].filter(Boolean).join('\n'));
  }
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
function runNpx(args, options = {}) {
  const spec = npxSpec(args);
  return run(spec.command, spec.args, options);
}
function clone(value) { return structuredClone(value); }
function resolveRef(document, value, seen = new Set()) {
  if (!value || typeof value !== 'object' || typeof value.$ref !== 'string') return value;
  const ref = value.$ref;
  assert(ref.startsWith('#/'), `resolved bundle contains non-local ref: ${ref}`);
  assert(!seen.has(ref), `cyclic ref reached proof resolver: ${ref}`);
  let current = document;
  for (const raw of ref.slice(2).split('/')) {
    const key = raw.replaceAll('~1', '/').replaceAll('~0', '~');
    assert(current && Object.hasOwn(current, key), `unresolved bundled ref: ${ref}`);
    current = current[key];
  }
  return resolveRef(document, current, new Set(seen).add(ref));
}
function schemaBySuffix(document, suffix) {
  const matches = Object.entries(document.components?.schemas ?? {}).filter(([name]) => name === suffix || name.endsWith(`_${suffix}`) || name.endsWith(`-${suffix}`));
  assert(matches.length === 1, `expected exactly one bundled schema ending ${suffix}, found ${matches.map(([name]) => name).join(', ')}`);
  return resolveRef(document, matches[0][1]);
}
function operations(document) {
  const result = [];
  for (const [path, pathItem] of Object.entries(document.paths ?? {})) {
    for (const [method, operation] of Object.entries(pathItem ?? {})) {
      if (httpMethods.has(method) && operation && typeof operation === 'object') result.push({ path, method, operation });
    }
  }
  return result;
}
function operationById(document, operationId) {
  const matches = operations(document).filter(({ operation }) => operation.operationId === operationId);
  assert(matches.length === 1, `expected one ${operationId}, found ${matches.length}`);
  return matches[0].operation;
}
function operationParameter(document, operationId, name) {
  const operation = operationById(document, operationId);
  const matches = (operation.parameters ?? []).map((value) => resolveRef(document, value)).filter((parameter) => parameter?.name === name);
  assert(matches.length === 1, `${operationId} must expose exactly one ${name} parameter`);
  return matches[0];
}
function assertParameterEnum(document, operationId, name, expected) {
  const parameter = operationParameter(document, operationId, name);
  sameSet(parameter.schema?.enum ?? [], expected, `${operationId}.${name}`);
}
function assertParameterDateTime(document, operationId, name) {
  const parameter = operationParameter(document, operationId, name);
  assert(parameter.schema?.type === 'string' && parameter.schema?.format === 'date-time', `${operationId}.${name} must be date-time`);
}
function assertRequired(schema, names, label) {
  for (const name of names) assert((schema.required ?? []).includes(name), `${label} must require ${name}`);
}
function assertProperties(schema, names, label) {
  for (const name of names) assert(Object.hasOwn(schema.properties ?? {}, name), `${label} must expose ${name}`);
}
function assertForbidden(schema, names, label) {
  for (const name of names) assert(!Object.hasOwn(schema.properties ?? {}, name), `${label} must not expose ${name}`);
}
function validate(document) {
  const ops = operations(document);
  assert(ops.length === 99, `Product operation count changed: ${ops.length}/99`);
  const ordinaryPermissions = normalize(ops.map(({ operation }) => operation['x-mpc-required-permission']).filter((value) => value && value !== 'authenticated'));
  assert(ordinaryPermissions.length === 30, `ordinary Permission count changed: ${ordinaryPermissions.length}/30`);

  const businessOrder = schemaBySuffix(document, 'BusinessOrderIntentListItem');
  assertRequired(businessOrder, ['business_order_intent_id', 'sale', 'target_source_instance_id', 'external_effect_state', 'convergence', 'created_at'], 'BusinessOrderIntentListItem');
  assertProperties(businessOrder, ['convergence'], 'BusinessOrderIntentListItem');
  assertParameterEnum(document, 'ListBusinessOrderIntents', 'external_effect_state', effectStates);
  assertParameterEnum(document, 'ListBusinessOrderIntents', 'convergence', convergenceStates);

  const invoicing = schemaBySuffix(document, 'InvoicingIntent');
  assertRequired(invoicing, ['invoicing_intent_id', 'sale', 'business_order_intent_id', 'external_effect_state', 'convergence', 'created_at'], 'InvoicingIntent');
  const invoicingList = schemaBySuffix(document, 'InvoicingIntentListItem');
  assertRequired(invoicingList, ['invoicing_intent_id', 'sale', 'business_order_intent_id', 'external_effect_state', 'convergence', 'created_at'], 'InvoicingIntentListItem');
  assertProperties(invoicingList, ['fulfillment_execution_id'], 'InvoicingIntentListItem');
  assertParameterEnum(document, 'ListInvoicingIntents', 'external_effect_state', effectStates);
  assertParameterEnum(document, 'ListInvoicingIntents', 'convergence', convergenceStates);

  const fulfillment = schemaBySuffix(document, 'FulfillmentExecutionListItem');
  assertRequired(fulfillment, ['fulfillment_execution_id', 'sale', 'scope', 'separation', 'physical_conference', 'packing', 'dispatch_handoff', 'physical_readiness', 'created_at'], 'FulfillmentExecutionListItem');
  assertProperties(fulfillment, ['provider_dispatch_deadline'], 'FulfillmentExecutionListItem');
  assertParameterEnum(document, 'ListFulfillmentExecutions', 'physical_readiness', ['ready', 'blocked', 'unknown']);
  for (const name of ['separation_state', 'physical_conference_state', 'packing_state', 'dispatch_handoff_state']) {
    assertParameterEnum(document, 'ListFulfillmentExecutions', name, checkpointStates);
  }
  assertParameterDateTime(document, 'ListFulfillmentExecutions', 'provider_dispatch_deadline_before');

  const shipment = schemaBySuffix(document, 'ShipmentListItem');
  assertProperties(shipment, ['sale', 'dispatch_deadline'], 'ShipmentListItem');
  assertParameterEnum(document, 'ListShipments', 'state', shipmentStates);

  const forbidden = ['operational_stage', 'next_action', 'priority', 'urgency_score', 'total_count', 'kanban_column'];
  for (const [name, schema] of [
    ['BusinessOrderIntentListItem', businessOrder],
    ['InvoicingIntentListItem', invoicingList],
    ['FulfillmentExecutionListItem', fulfillment],
    ['ShipmentListItem', shipment],
  ]) assertForbidden(schema, forbidden, name);

  for (const { path, operation } of ops) {
    assert(!path.includes('operational-dashboard'), `screen-shaped operational dashboard path is forbidden: ${path}`);
    assert(operation['x-mpc-semantic-owner'] !== 'OperationalWorkflow', `${operation.operationId} must not create OperationalWorkflow owner`);
    for (const parameterRef of operation.parameters ?? []) {
      const parameter = resolveRef(document, parameterRef);
      assert(!forbidden.includes(parameter?.name), `${operation.operationId} must not expose synthetic ${parameter?.name} filter`);
    }
  }
}
function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `operational-read negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

try {
  mkdirSync(temp, { recursive: true });
  const bundlePath = join(temp, 'product.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundlePath]);
  const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
  validate(document);

  expectFailure('BusinessOrderIntent list loses convergence', document, (candidate) => {
    const schema = schemaBySuffix(candidate, 'BusinessOrderIntentListItem');
    schema.required = (schema.required ?? []).filter((name) => name !== 'convergence');
  });
  expectFailure('Fulfillment list gains synthetic operational stage', document, (candidate) => {
    const schema = schemaBySuffix(candidate, 'FulfillmentExecutionListItem');
    schema.properties.operational_stage = { type: 'string' };
  });

  assert(negativeControls === 2, `operational-read negative-control count mismatch: ${negativeControls}/2`);
  console.log('operational_read_owner_local_projection=PASS');
  console.log('operational_read_filters=PASS');
  console.log(`operational_read_negative_controls=${negativeControls}/2`);
  console.log('operational_read_contract_proof=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
