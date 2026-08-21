import { createHash } from 'node:crypto';
import { cpSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const contractDir = join(root, 'contracts/api/product');
const entrypoint = join(contractDir, 'openapi.yaml');
const redoclyConfig = join(contractDir, 'redocly.yaml');
const sourceNames = [
  'openapi.yaml',
  'components.yaml',
  'paths-identity-portfolio-readiness.yaml',
  'paths-offering-availability-market.yaml',
  'paths-economics-governance-sales-materialization.yaml',
  'paths-fulfillment-postsale-work.yaml',
];
const sourceTreeText = sourceNames.map((name) => readFileSync(join(contractDir, name), 'utf8')).join('\n');
const w4Text = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md'), 'utf8');
const w2Text = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md'), 'utf8');
const admissionText = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md'), 'utf8');
const temp = mkdtempSync(join(tmpdir(), 'mpc-product-oad-'));
const go = process.platform === 'win32' ? 'go.exe' : 'go';
const npxCli = process.env.npm_execpath
  ? resolve(dirname(process.env.npm_execpath), 'npx-cli.js')
  : resolve(dirname(process.execPath), 'node_modules/npm/bin/npx-cli.js');
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const PROBLEM_PREFIX = 'https://conexus.fun/marketplace-central/problems/product/';
const ALLOWED_EXTENSIONS = new Set(['x-mpc-operation-class', 'x-mpc-required-permission', 'x-mpc-principal-kinds', 'x-mpc-semantic-owner', 'x-mpc-required-physical-qualification']);
const EXPECTED_NEGATIVE_CONTROLS = 12;
let negativeControlCount = 0;

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
  return process.platform === 'win32'
    ? { command: process.execPath, args: [npxCli, ...args] }
    : { command: 'npx', args };
}
function runNpx(args, options = {}) {
  const spec = npxSpec(args);
  return run(spec.command, spec.args, options);
}
function spawnNpx(args, options = {}) {
  const spec = npxSpec(args);
  return spawnSync(spec.command, spec.args, {
    cwd: options.cwd ?? root,
    env: { ...process.env, ...(options.env ?? {}) },
    encoding: 'utf8',
    shell: false,
  });
}
function sha256(path) { return createHash('sha256').update(readFileSync(path)).digest('hex'); }
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
function walk(value, visitor) {
  if (!value || typeof value !== 'object') return;
  visitor(value);
  for (const child of Object.values(value)) walk(child, visitor);
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
function parameters(document, entry) {
  return [...(entry.pathItem.parameters ?? []), ...(entry.operation.parameters ?? [])].map((value) => resolveRef(document, value));
}
function response(document, operation, status) {
  const raw = operation.responses?.[status];
  return raw ? resolveRef(document, raw) : undefined;
}
function responseSchema(document, operation, status = '200') {
  const raw = response(document, operation, status)?.content?.['application/json']?.schema;
  return raw ? resolveRef(document, raw) : undefined;
}
function requestSchema(document, operation, mediaType = 'application/json') {
  const raw = resolveRef(document, operation.requestBody)?.content?.[mediaType]?.schema;
  return raw ? resolveRef(document, raw) : undefined;
}
function byId(all, id) {
  const found = all.find((entry) => entry.operation.operationId === id);
  assert(found, `operation not found: ${id}`);
  return found;
}
function component(document, kind, name) {
  const direct = document.components?.[kind]?.[name];
  assert(direct, `bundled component not found: ${kind}.${name}`);
  return resolveRef(document, direct);
}
function generatedSchemaExpression(text, name, lane) {
  const match = text.match(new RegExp(`${name}:\\s*([^;]+);`));
  assert(match, `${lane} projection lacks ${name} schema expression`);
  return match[1];
}
function generatedObjectBlock(text, name, lane) {
  const match = text.match(new RegExp(`${name}:\\s*\\{([\\s\\S]*?)\\n\\s*\\};`));
  assert(match, `${lane} projection lacks ${name} object`);
  return match[1];
}
function generatedGoStruct(text, name) {
  const match = text.match(new RegExp(`type ${name} struct \\{([\\s\\S]*?)\\n\\}`));
  assert(match, `Go projection lacks ${name} struct`);
  return match[1];
}

function parseW4(text) {
  const start = text.indexOf('# 8. Exact 95-operation enforcement matrix');
  const end = text.indexOf('\n---\n\n## 9.', start);
  assert(start >= 0 && end > start, 'unable to locate W4 exact matrix');
  const rows = [];
  for (const line of text.slice(start, end).split(/\r?\n/)) {
    const match = line.match(/^\|\s*`([^`]+)`\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$/);
    if (!match) continue;
    const [, operationId, classCell, permissionCell, principalCell] = match;
    const operationClass = classCell.includes('Q') ? 'Q' : classCell.includes('C') ? 'C' : null;
    const permission = permissionCell.match(/`([^`]+)`/)?.[1] ?? (permissionCell.includes('authenticated special condition') ? 'authenticated' : null);
    const principalKinds = [...principalCell.matchAll(/\b([HAS])\b/g)].map((m) => m[1]);
    assert(operationClass && permission && principalKinds.length, `unparseable W4 row: ${line}`);
    rows.push({ operationId, operationClass, permission, principalKinds: normalize(principalKinds), qualifiedPhysical: /currently qualified S/.test(principalCell) });
  }
  assert(rows.length === 95, `W4 row count changed: ${rows.length}`);
  assert(new Set(rows.map((row) => row.operationId)).size === 95, 'W4 operation IDs are not unique');
  return rows;
}
function parseProblemSlugs(text) {
  const start = text.indexOf('# 19. Problem Details catalog');
  const end = text.indexOf('Problem `type` is an absolute stable URI', start);
  assert(start >= 0 && end > start, 'unable to locate W2 Problem catalog');
  const slugs = [...text.slice(start, end).matchAll(/^- `([^`]+)`[;.]/gm)].map((m) => m[1]);
  assert(slugs.length === 15, `W2 Product Problem count changed: ${slugs.length}`);
  return slugs;
}
function parseMandatoryIdempotency(text) {
  const start = text.indexOf('# 7. Whole-Matrix complete C-operation safety sweep');
  const end = text.indexOf('### 7.6 Safety-sweep result', start);
  assert(start >= 0 && end > start, 'unable to locate admission safety sweep');
  const result = [];
  for (const line of text.slice(start, end).split(/\r?\n/)) {
    if (!line.startsWith('|')) continue;
    const cells = line.split('|').slice(1, -1).map((cell) => cell.trim());
    if (cells.length !== 4 || cells[0] === 'Operation' || /^-+$/.test(cells[0])) continue;
    const operation = cells[0].replace(/^`|`$/g, '');
    if (/mandatory client key/i.test(cells[2])) result.push(operation);
  }
  assert(result.length === 14, `mandatory Idempotency-Key disposition count changed: ${result.length}`);
  return normalize(result);
}

const w4 = parseW4(w4Text);
const w4ById = new Map(w4.map((row) => [row.operationId, row]));
const expectedPermissions = normalize(w4.map((row) => row.permission).filter((permission) => permission !== 'authenticated'));
const expectedProblems = parseProblemSlugs(w2Text).map((slug) => `${PROBLEM_PREFIX}${slug}`);
const expectedIdempotency = parseMandatoryIdempotency(admissionText);
assert(expectedPermissions.length === 29, `ordinary Permission count changed: ${expectedPermissions.length}`);
assert(w4ById.has('UpdateAvailabilityAllocationScopePolicy'), 'W4 final Availability policy operation name is missing');
assert(admissionText.includes('| update allocation/scope policy |'), 'accepted pre-wire Availability policy safety row is missing');

function validate(document) {
  assert(document.openapi === '3.1.2', `OpenAPI must be 3.1.2, found ${document.openapi}`);
  assert(document.jsonSchemaDialect === 'https://spec.openapis.org/oas/3.1/dialect/base', 'unexpected JSON Schema dialect');
  assert(!Object.hasOwn(document, 'servers'), 'canonical/resolved Product OAD must omit servers');
  assert(JSON.stringify(document.security) === JSON.stringify([{ MpcBearerAuth: [] }]), 'root security must be exactly MpcBearerAuth');
  const bearer = document.components?.securitySchemes?.MpcBearerAuth;
  assert(bearer?.type === 'http' && bearer?.scheme === 'bearer', 'MpcBearerAuth must be HTTP bearer');
  assert(!Object.hasOwn(bearer ?? {}, 'bearerFormat'), 'Product bearer credential format is a D7 mechanism and must remain unspecified');

  walk(document, (node) => {
    for (const key of Object.keys(node)) {
      if (key.startsWith('x-')) assert(ALLOWED_EXTENSIONS.has(key), `resolved Product OAD contains non-allowlisted extension ${key}`);
    }
  });

  const all = operations(document);
  const ids = all.map((entry) => entry.operation.operationId);
  assert(all.length === 95, `Product operation count must be 95, found ${all.length}`);
  assert(new Set(ids).size === 95, 'operationId values are not unique');
  sameSet(ids, [...w4ById.keys()], 'OAD/W4 operation IDs');

  for (const path of Object.keys(document.paths ?? {})) {
    assert(path === '/access-context' || path.startsWith('/organizations/{organization_id}/'), `path outside canonical Product roots: ${path}`);
    assert(!/(^|\/)(providers?|integrations?|webhooks?|oauth|callbacks?|acquisition)(\/|$|:)/i.test(path), `technical/provider path leaked into Product API: ${path}`);
    assert(!path.startsWith('/v1'), `version prefix leaked into Product path: ${path}`);
  }

  for (const entry of all) {
    const op = entry.operation;
    const expected = w4ById.get(op.operationId);
    assert(/^[A-Z][A-Za-z0-9]*$/.test(op.operationId), `operationId is not stable PascalCase: ${op.operationId}`);
    assert(!op.operationId.startsWith('AcquireMarketplace'), `internal acquisition operation leaked: ${op.operationId}`);
    for (const key of ['x-mpc-operation-class', 'x-mpc-required-permission', 'x-mpc-principal-kinds', 'x-mpc-semantic-owner']) assert(Object.hasOwn(op, key), `${op.operationId} missing ${key}`);
    for (const key of Object.keys(op).filter((key) => key.startsWith('x-'))) assert(ALLOWED_EXTENSIONS.has(key), `${op.operationId} contains non-allowlisted extension ${key}`);
    assert(!Object.hasOwn(op, 'security'), `${op.operationId} overrides global Product security`);
    assert(!Object.hasOwn(op.responses ?? {}, '202'), `${op.operationId} introduces generic 202`);
    assert(op['x-mpc-operation-class'] === expected.operationClass, `${op.operationId} class mismatch`);
    assert(op['x-mpc-required-permission'] === expected.permission, `${op.operationId} permission mismatch`);
    sameSet(op['x-mpc-principal-kinds'], expected.principalKinds, `${op.operationId} principal kinds`);
    assert(op['x-mpc-principal-kinds'].every((kind) => ['H', 'A', 'S'].includes(kind)), `${op.operationId} introduces non-H/A/S principal kind`);
    assert((op['x-mpc-required-physical-qualification'] === true) === expected.qualifiedPhysical, `${op.operationId} physical qualification mismatch`);
    if (entry.path.startsWith('/organizations/{organization_id}/')) {
      assert(response(document, op, '404'), `${op.operationId} must preserve privacy-protecting 404 for path-Organization non-membership/not-found`);
    }
  }

  sameSet(component(document, 'schemas', 'Permission').enum ?? [], expectedPermissions, 'Permission enum');
  sameSet(all.map((entry) => entry.operation['x-mpc-required-permission']).filter((permission) => permission !== 'authenticated'), expectedPermissions, 'operation Permission vocabulary');

  const listSearch = all.filter((entry) => /^(List|Search)/.test(entry.operation.operationId));
  assert(listSearch.length === 26, `List/Search operation count must be 26, found ${listSearch.length}`);
  assert(listSearch.filter((entry) => entry.operation.operationId.startsWith('Search')).length === 1, 'exactly one Search operation is required');
  for (const entry of listSearch) {
    const names = parameters(document, entry).filter((parameter) => parameter.in === 'query').map((parameter) => parameter.name);
    assert(names.includes('limit') && names.includes('cursor'), `${entry.operation.operationId} lacks limit/cursor`);
    assert(!names.some((name) => ['page', 'offset', 'sort', 'total', 'total_count'].includes(name)), `${entry.operation.operationId} contains offset/sort/total grammar`);
    const schema = responseSchema(document, entry.operation);
    assert(schema?.type === 'object', `${entry.operation.operationId} collection response is not an object`);
    assert(Object.hasOwn(schema.properties ?? {}, 'next_cursor'), `${entry.operation.operationId} lacks next_cursor`);
    assert(!(schema.required ?? []).includes('next_cursor'), `${entry.operation.operationId} requires next_cursor at exhaustion`);
  }
  const limit = component(document, 'parameters', 'Limit');
  assert(limit.schema?.type === 'integer' && limit.schema.minimum === 1 && Number.isInteger(limit.schema.default) && Number.isInteger(limit.schema.maximum) && limit.schema.default > 0 && limit.schema.maximum >= limit.schema.default, 'Limit bounds are incoherent');

  for (const [name, raw] of Object.entries(document.components?.schemas ?? {})) {
    const schema = resolveRef(document, raw);
    if (schema?.type === 'object') assert(schema.additionalProperties === false, `${name} semantic object is not closed`);
    if (Array.isArray(schema?.oneOf)) {
      for (const branchRef of schema.oneOf) {
        const branch = resolveRef(document, branchRef);
        const hasConst = Object.values(branch?.properties ?? {}).some((property) => property && typeof property === 'object' && Object.hasOwn(property, 'const'));
        assert(hasConst, `${name} oneOf branch lacks const discriminator`);
      }
    }
  }
  const decimal = component(document, 'schemas', 'ExactDecimalString');
  assert(decimal.type === 'string' && typeof decimal.pattern === 'string' && !decimal.pattern.toLowerCase().includes('e'), 'exact decimal grammar is not a non-exponent string');
  const etag = component(document, 'schemas', 'StrongETag');
  assert(etag.type === 'string' && typeof etag.pattern === 'string', 'StrongETag must be a patterned string');
  const etagPattern = new RegExp(etag.pattern);
  assert(etagPattern.test('"opaque-validator"') && !etagPattern.test('W/"weak"') && !etagPattern.test('opaque-validator'), 'StrongETag must admit only quoted non-weak entity tags');
  const desiredKnown = component(document, 'schemas', 'DesiredQuantityKnown');
  const desiredQuantity = resolveRef(document, desiredKnown.properties?.quantity);
  const desiredQuantityPattern = new RegExp(desiredQuantity?.pattern ?? '');
  assert(desiredQuantity?.type === 'string' && desiredQuantityPattern.test('0') && desiredQuantityPattern.test('1.25') && !desiredQuantityPattern.test('-1'), 'Availability desired quantity must remain an exact non-negative decimal string');
  const assignment = component(document, 'schemas', 'WorkAssignment');
  const assignmentBranches = (assignment.oneOf ?? []).map((branch) => resolveRef(document, branch));
  const assignedBranch = assignmentBranches.find((branch) => branch.properties?.state?.const === 'assigned');
  const unassignedBranch = assignmentBranches.find((branch) => branch.properties?.state?.const === 'unassigned');
  assert(assignedBranch && (assignedBranch.required ?? []).includes('principal_id'), 'assigned WorkAssignment must require principal_id');
  assert(unassignedBranch && !Object.hasOwn(unassignedBranch.properties ?? {}, 'principal_id'), 'unassigned WorkAssignment must not carry principal_id');
  const fulfillmentArtifact = component(document, 'schemas', 'FulfillmentArtifact');
  assert((fulfillmentArtifact.required ?? []).includes('sensitivity'), 'FulfillmentArtifact must carry explicit sensitivity classification');

  const customProblemTypes = [];
  walk(document, (node) => {
    const value = node?.properties?.type?.const;
    if (typeof value === 'string' && value.startsWith(PROBLEM_PREFIX)) customProblemTypes.push(value);
  });
  sameSet(customProblemTypes, expectedProblems, 'Product Problem catalog/origin');
  for (const [name, status] of [['AboutBlank413Problem', 413], ['AboutBlank415Problem', 415]]) {
    const schema = component(document, 'schemas', name);
    assert(schema.properties?.type?.const === 'about:blank' && schema.properties?.status?.const === status, `${name} must be about:blank/${status}`);
  }

  const idempotencyOps = all.flatMap((entry) => {
    const carriers = parameters(document, entry).filter((parameter) => parameter.in === 'header' && parameter.name.toLowerCase() === 'idempotency-key');
    if (carriers.length === 0) return [];
    assert(carriers.length === 1 && carriers[0].required === true, `${entry.operation.operationId} must have exactly one required Idempotency-Key header`);
    return [entry.operation.operationId];
  });
  sameSet(idempotencyOps, expectedIdempotency, 'Idempotency-Key carrier set');

  const ifMatchExpected = ['UpdateMarketplaceInstallationConfiguration', 'UpdateListingIntentDraft', 'UpdateInventorySource', 'UpdateAvailabilityAllocationScopePolicy', 'UpdateCommercialPolicy', 'UpdateAuthorizationDelegation', 'UpdateFulfillmentNode', 'UpdateFulfillmentOperatingTargets'];
  const ifMatchActual = all.filter((entry) => parameters(document, entry).some((parameter) => parameter.in === 'header' && parameter.name.toLowerCase() === 'if-match')).map((entry) => entry.operation.operationId);
  sameSet(ifMatchActual, ifMatchExpected, 'If-Match carrier set');

  for (const id of ['DeactivateMarketplaceInstallation', 'DiscardListingIntentDraft', 'SubmitListingIntent', 'DeactivateInventorySource', 'ResolveEconomicAttribution', 'ResolveSaleSellingEntityAttribution', 'ResolveBusinessSystemPartyResolution', 'DeactivateFulfillmentNode', 'RecordSeparation', 'RecordPhysicalConference', 'RecordPacking', 'RecordDispatchHandoff', 'AssignWork', 'ClearWorkAssignment', 'HoldWork', 'ResumeWork', 'EscalateWork']) {
    const schema = requestSchema(document, byId(all, id).operation);
    assert(schema?.properties?.etag && (schema.required ?? []).includes('etag'), `${id} lacks required typed etag`);
  }
  for (const id of ['ResolveProductChannelCorrespondence', 'ClearProductChannelCorrespondence']) {
    const schema = requestSchema(document, byId(all, id).operation);
    assert(schema?.properties?.correspondence_etag && (schema.required ?? []).includes('correspondence_etag'), `${id} lacks required correspondence_etag`);
  }
  for (const branch of component(document, 'schemas', 'AuthorizationTarget').oneOf ?? []) {
    const schema = resolveRef(document, branch);
    assert(schema.properties?.etag && (schema.required ?? []).includes('etag'), 'AuthorizationDecision target branch lacks exact ETag');
  }
  const superseded = component(document, 'schemas', 'SupersededPriceIntentRef');
  assert(superseded.properties?.etag && (superseded.required ?? []).includes('etag'), 'superseded PriceIntent ref lacks exact ETag');

  const durableCreates = ['CreateMarketplaceInstallation', 'CreateListingIntentDraft', 'CreatePriceIntent', 'CreateInventorySource', 'CreateAuthorizationDecision', 'EstablishAuthorizationDelegation', 'CreateFulfillmentNode', 'CreatePostSaleResolution'];
  for (const id of durableCreates) assert(response(document, byId(all, id).operation, '201')?.headers?.Location, `${id} must return 201 + Location`);
  for (const entry of all.filter((entry) => entry.method === 'post' && !durableCreates.includes(entry.operation.operationId))) assert(response(document, entry.operation, '200'), `${entry.operation.operationId} custom POST must return 200`);

  const media = byId(all, 'CreateListingIntentMedia').operation;
  const mediaContent = resolveRef(document, media.requestBody)?.content?.['multipart/form-data'];
  const mediaSchema = resolveRef(document, mediaContent?.schema);
  sameSet(Object.keys(mediaSchema?.properties ?? {}), ['file', 'etag'], 'CreateListingIntentMedia multipart parts');
  sameSet(mediaSchema?.required ?? [], ['file', 'etag'], 'CreateListingIntentMedia required multipart parts');
  assert(mediaContent?.encoding?.etag?.contentType === 'text/plain', 'CreateListingIntentMedia etag part must be text/plain');
  const mediaResult = responseSchema(document, media);
  assert(mediaResult?.properties?.listing_intent_etag && (mediaResult.required ?? []).includes('listing_intent_etag'), 'CreateListingIntentMedia result lacks current parent ListingIntent validator');

  const workOriginKinds = (component(document, 'schemas', 'WorkOrigin').oneOf ?? []).map((branch) => {
    const schema = resolveRef(document, branch);
    return Object.values(schema.properties ?? {}).find((property) => property && typeof property === 'object' && Object.hasOwn(property, 'const'))?.const;
  }).filter(Boolean);
  sameSet(workOriginKinds, ['readiness', 'listing_intent', 'price_intent', 'marketplace_sale', 'business_order_intent', 'fulfillment_execution', 'shipment', 'post_sale_resolution', 'economic_attribution', 'authorization_decision'], 'owner-specific Work origins');

  return { all, idempotencyOps, listSearch };
}

function expectFailure(name, factory) {
  let failed = false;
  try { validate(factory()); } catch { failed = true; }
  assert(failed, `negative control unexpectedly passed: ${name}`);
  negativeControlCount++;
}
function negativeSemanticControls(document) {
  expectFailure('missing operation', () => { const candidate = clone(document); delete candidate.paths['/organizations/{organization_id}/work/{work_id}:escalate']; return candidate; });
  expectFailure('wrong W4 permission', () => { const candidate = clone(document); byId(operations(candidate), 'CreatePriceIntent').operation['x-mpc-required-permission'] = 'listing.manage'; return candidate; });
  expectFailure('technical path leakage', () => { const candidate = clone(document); candidate.paths['/organizations/{organization_id}/webhooks'] = clone(candidate.paths['/organizations/{organization_id}/work']); return candidate; });
  expectFailure('fourth principal kind', () => { const candidate = clone(document); byId(operations(candidate), 'GetCurrentAccessContext').operation['x-mpc-principal-kinds'].push('X'); return candidate; });
  expectFailure('missing idempotency carrier', () => { const candidate = clone(document); const entry = byId(operations(candidate), 'CreateInventorySource'); entry.operation.parameters = (entry.operation.parameters ?? []).filter((raw) => resolveRef(candidate, raw)?.name !== 'Idempotency-Key'); return candidate; });
  expectFailure('optional idempotency carrier', () => { const candidate = clone(document); component(candidate, 'parameters', 'IdempotencyKey').required = false; return candidate; });
  expectFailure('wrong Product Problem origin', () => { const candidate = clone(document); component(candidate, 'schemas', 'AccessDeniedProblem').properties.type.const = 'https://example.invalid/problem'; return candidate; });
  expectFailure('schema extension leakage', () => { const candidate = clone(document); component(candidate, 'schemas', 'Money')['x-enum-varnames'] = ['Money']; return candidate; });
  expectFailure('missing Organization privacy 404', () => { const candidate = clone(document); delete byId(operations(candidate), 'GetWork').operation.responses['404']; return candidate; });
  expectFailure('missing fulfillment artifact sensitivity', () => { const candidate = clone(document); const schema = component(candidate, 'schemas', 'FulfillmentArtifact'); schema.required = (schema.required ?? []).filter((name) => name !== 'sensitivity'); return candidate; });
}

function redoclyProof() {
  runNpx(['--yes', '@redocly/cli@2.45.0', 'lint', entrypoint, '--config', redoclyConfig]);
  const bundleA = join(temp, 'product-a.json');
  const bundleB = join(temp, 'product-b.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleA]);
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleB]);
  assert(sha256(bundleA) === sha256(bundleB), 'resolved Product bundle is not deterministic within one proof run');

  const fixtureDir = join(temp, 'broken-source');
  mkdirSync(fixtureDir, { recursive: true });
  for (const name of sourceNames) cpSync(join(contractDir, name), join(fixtureDir, name));
  cpSync(redoclyConfig, join(fixtureDir, 'redocly.yaml'));
  const brokenEntry = join(fixtureDir, 'openapi.yaml');
  const original = readFileSync(brokenEntry, 'utf8');
  const broken = original.replace('#/pathItems/AccessContext', '#/pathItems/__missing__');
  assert(broken !== original, 'unable to construct unresolved-ref negative fixture');
  writeFileSync(brokenEntry, broken, 'utf8');
  const negative = spawnNpx(['--yes', '@redocly/cli@2.45.0', 'lint', brokenEntry, '--config', join(fixtureDir, 'redocly.yaml')]);
  assert(negative.status !== 0, 'Redocly unresolved-ref negative fixture unexpectedly passed');
  negativeControlCount++;

  const methodProbeDir = join(temp, 'method-not-allowed-source-probe');
  mkdirSync(methodProbeDir, { recursive: true });
  for (const name of sourceNames) cpSync(join(contractDir, name), join(methodProbeDir, name));
  cpSync(redoclyConfig, join(methodProbeDir, 'redocly.yaml'));
  const methodProbeEntry = join(methodProbeDir, 'openapi.yaml');
  const methodProbeOriginal = readFileSync(methodProbeEntry, 'utf8');
  const componentMarker = '\ncomponents:\n';
  assert(methodProbeOriginal.includes(componentMarker), 'unable to construct 405 source expressibility probe');
  const methodProbe = methodProbeOriginal.replace(componentMarker, `\n  '/__proof-method-not-allowed':\n    get:\n      operationId: ProofMethodNotAllowed\n      responses:\n        '405': {$ref: './components.yaml#/responses/MethodNotAllowed'}\n${componentMarker}`);
  writeFileSync(methodProbeEntry, methodProbe, 'utf8');
  const methodProbeBundle = join(temp, 'method-not-allowed-probe.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', methodProbeEntry, '--config', join(methodProbeDir, 'redocly.yaml'), '-o', methodProbeBundle]);
  const methodProbeDocument = JSON.parse(readFileSync(methodProbeBundle, 'utf8'));
  const methodResponse = resolveRef(methodProbeDocument, methodProbeDocument.paths?.['/__proof-method-not-allowed']?.get?.responses?.['405']);
  assert(methodResponse?.headers?.Allow, 'source OAD 405 grammar does not resolve an Allow header');
  const methodProblem = resolveRef(methodProbeDocument, methodResponse?.content?.['application/problem+json']?.schema);
  assert(methodProblem?.properties?.type?.const === 'about:blank' && methodProblem?.properties?.status?.const === 405, 'source OAD 405 grammar must resolve about:blank/405');

  return { document: JSON.parse(readFileSync(bundleA, 'utf8')), bundlePath: bundleA };
}

function typescriptProof() {
  const first = join(temp, 'product-a.d.ts');
  const second = join(temp, 'product-b.d.ts');
  runNpx(['--yes', 'openapi-typescript@7.13.0', entrypoint, '-o', first, '--redocly', redoclyConfig]);
  runNpx(['--yes', 'openapi-typescript@7.13.0', entrypoint, '-o', second, '--redocly', redoclyConfig]);
  assert(sha256(first) === sha256(second), 'TypeScript projection is not deterministic within one proof run');
  const generated = readFileSync(first, 'utf8');
  assert(generated.includes('This file was auto-generated by openapi-typescript'), 'TypeScript projection lacks generated banner');
  assert(generated.includes('GetCurrentAccessContext') && generated.includes('EscalateWork'), 'TypeScript projection lacks boundary operation IDs');
  assert(/ExactDecimalString:\s*string;/.test(generated), 'TypeScript projection does not preserve ExactDecimalString as string');
  const publication = generatedSchemaExpression(generated, 'PublicationValue', 'TypeScript');
  for (const branch of ['PublicationTextValue', 'PublicationExactDecimalValue', 'PublicationBooleanValue', 'PublicationOptionValue', 'PublicationTextListValue', 'PublicationOptionListValue', 'PublicationNumberUnitValue', 'PublicationNotApplicableValue']) assert(publication.includes(`components["schemas"]["${branch}"]`), `TypeScript PublicationValue lost ${branch}`);
  const desiredQuantity = generatedSchemaExpression(generated, 'DesiredQuantity', 'TypeScript');
  for (const branch of ['DesiredQuantityKnown', 'DesiredQuantityUnknown', 'DesiredQuantityUnavailable']) assert(desiredQuantity.includes(`components["schemas"]["${branch}"]`), `TypeScript DesiredQuantity lost ${branch}`);
  const mediaMultipart = generatedObjectBlock(generated, 'CreateListingIntentMediaMultipart', 'TypeScript');
  assert(/\bfile:\s*[^;]+;/.test(mediaMultipart), 'TypeScript multipart projection lost file part');
  assert(/etag:\s*components\["schemas"\]\["StrongETag"\];/.test(mediaMultipart), 'TypeScript multipart projection lost typed etag part');
  assert(generated.includes('"multipart/form-data": components["schemas"]["CreateListingIntentMediaMultipart"];'), 'TypeScript operation projection does not reach CreateListingIntentMedia multipart schema');
  const mediaResult = generatedObjectBlock(generated, 'CreateListingIntentMediaResult', 'TypeScript');
  assert(/listing_intent_etag:\s*components\["schemas"\]\["StrongETag"\];/.test(mediaResult), 'TypeScript media result lost typed parent ListingIntent validator');
  assert(generated.includes('"application/json": components["schemas"]["CreateListingIntentMediaResult"];'), 'TypeScript operation projection does not reach CreateListingIntentMedia result schema');
  runNpx(['--yes', '-p', 'typescript@5.9.3', 'tsc', '--noEmit', '--strict', '--skipLibCheck', first]);
}

function goProof(bundlePath, document) {
  const dir = join(temp, 'go-proof');
  mkdirSync(dir, { recursive: true });
  writeFileSync(join(dir, 'oapi-codegen.yaml'), [
    'package: productapi',
    'output: product.gen.go',
    'generate:',
    '  models: true',
    '  std-http-server: true',
    '  strict-server: true',
    'output-options:',
    '  skip-prune: true',
    '',
  ].join('\n'));
  const generate = () => run(go, ['run', 'github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0', '-config', 'oapi-codegen.yaml', bundlePath], { cwd: dir });
  generate();
  const generatedPath = join(dir, 'product.gen.go');
  const firstHash = sha256(generatedPath);
  generate();
  assert(sha256(generatedPath) === firstHash, 'Go projection is not deterministic within one proof run');
  const generated = readFileSync(generatedPath, 'utf8');
  assert(generated.includes('Code generated by github.com/oapi-codegen/oapi-codegen'), 'Go projection lacks generated banner');
  assert(generated.includes('type StrictServerInterface interface'), 'Go projection lacks StrictServerInterface');
  assert(generated.includes('GetCurrentAccessContext') && generated.includes('EscalateWork'), 'Go projection lacks boundary operation IDs');
  assert(!/\"(os\/exec|unsafe|syscall)\"/.test(generated), 'Go projection contains forbidden high-risk import');
  assert(/type ExactDecimalString\s+(?:=\s*)?string\b/.test(generated), 'Go projection does not preserve ExactDecimalString as string');
  for (const branch of ['PublicationTextValue', 'PublicationExactDecimalValue', 'PublicationBooleanValue', 'PublicationOptionValue', 'PublicationTextListValue', 'PublicationOptionListValue', 'PublicationNumberUnitValue', 'PublicationNotApplicableValue']) assert(generated.includes(`As${branch}()`), `Go PublicationValue lost ${branch} union projection`);
  for (const branch of ['DesiredQuantityKnown', 'DesiredQuantityUnknown', 'DesiredQuantityUnavailable']) assert(generated.includes(`As${branch}()`), `Go DesiredQuantity lost ${branch} knowledge-state projection`);
  const mediaRequest = generatedGoStruct(generated, 'CreateListingIntentMediaRequestObject');
  assert(/\bBody\s+\*multipart\.Reader\b/.test(mediaRequest), 'Go strict request projection must expose CreateListingIntentMedia multipart.Reader');
  const goMediaMultipart = generatedGoStruct(generated, 'CreateListingIntentMediaMultipart');
  assert(/`json:"file"`/.test(goMediaMultipart) && /`json:"etag"`/.test(goMediaMultipart), 'Go multipart model lost file or etag part');
  assert(/type CreateListingIntentMedia200JSONResponse(?:\s*=\s*|\s+)CreateListingIntentMediaResult\b/.test(generated), 'Go operation response does not preserve CreateListingIntentMediaResult');
  const goMediaResult = generatedGoStruct(generated, 'CreateListingIntentMediaResult');
  assert(/\bStrongETag\s+`json:"listing_intent_etag"`/.test(goMediaResult), 'Go media result lost typed parent ListingIntent validator');

  const expectedPatterns = operations(document)
    .map((entry) => `${entry.method.toUpperCase()} ${entry.path}`)
    .sort();
  const goExpectedPatterns = expectedPatterns.map((pattern) => `\t${JSON.stringify(pattern)},`).join('\n');

  writeFileSync(join(dir, 'go.mod'), ['module productoadproof', '', 'go 1.25.1', '', 'require github.com/oapi-codegen/runtime v1.7.0', ''].join('\n'));
  writeFileSync(join(dir, 'colon_suffix_test.go'), `package productapi\n\nimport (\n  \"net/http\"\n  \"net/http/httptest\"\n  \"sort\"\n  \"testing\"\n)\n\ntype captureMux struct {\n  handlers map[string]func(http.ResponseWriter, *http.Request)\n  patterns []string\n}\n\nfunc newCaptureMux() *captureMux { return &captureMux{handlers: map[string]func(http.ResponseWriter, *http.Request){}} }\nfunc (m *captureMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {\n  m.patterns = append(m.patterns, pattern)\n  m.handlers[pattern] = handler\n}\nfunc (m *captureMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n  const pattern = \"POST /organizations/{organization_id}/listing-intents/{listing_intent_id}:submit\"\n  handler, ok := m.handlers[pattern]\n  if !ok { http.NotFound(w, r); return }\n  r.SetPathValue(\"organization_id\", \"org\")\n  r.SetPathValue(\"listing_intent_id\", \"li\")\n  handler(w, r)\n}\n\ntype probeServer struct { ServerInterface }\n\nfunc TestGeneratedCanonicalRouteRegistrationAndColonSuffixDispatch(t *testing.T) {\n  mux := newCaptureMux()\n  called := false\n  middleware := func(_ http.Handler) http.Handler {\n    return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { called = true; w.WriteHeader(http.StatusNoContent) })\n  }\n  handler := HandlerWithOptions(probeServer{}, StdHTTPServerOptions{BaseRouter: mux, Middlewares: []MiddlewareFunc{middleware}})\n\n  got := append([]string(nil), mux.patterns...)\n  want := []string{\n${goExpectedPatterns}\n  }\n  sort.Strings(got)\n  sort.Strings(want)\n  if len(got) != len(want) { t.Fatalf(\"generated route count mismatch: got=%d want=%d\", len(got), len(want)) }\n  for i := range want {\n    if got[i] != want[i] { t.Fatalf(\"generated route mismatch at %d: got=%q want=%q\", i, got[i], want[i]) }\n  }\n\n  res := httptest.NewRecorder()\n  handler.ServeHTTP(res, httptest.NewRequest(http.MethodPost, \"/organizations/org/listing-intents/li:submit\", nil))\n  if !called || res.Code != http.StatusNoContent { t.Fatalf(\"generated colon-suffix dispatch failed: called=%v status=%d\", called, res.Code) }\n}\n`);
  run(go, ['mod', 'tidy'], { cwd: dir });
  const runtimeVersion = run(go, ['list', '-m', '-f', '{{.Version}}', 'github.com/oapi-codegen/runtime'], { cwd: dir }).stdout.trim();
  assert(runtimeVersion === 'v1.7.0', `unexpected oapi-codegen runtime: ${runtimeVersion}`);
  run(go, ['test', './...'], { cwd: dir, env: { GOTOOLCHAIN: 'go1.25.1' } });
  const minimumVersion = run(go, ['version'], { cwd: dir, env: { GOTOOLCHAIN: 'go1.25.1' } }).stdout.trim();
  assert(/go1\.25\.1\b/.test(minimumVersion), `minimum toolchain proof did not execute go1.25.1: ${minimumVersion}`);
  run(go, ['test', './...'], { cwd: dir });
  console.log(`go_minimum=${minimumVersion}`);
  console.log(`go_current=${run(go, ['version'], { cwd: dir }).stdout.trim()}`);
}

function sourceRefs(text) {
  return [...text.matchAll(/(?:["']?\$ref["']?)\s*:\s*['"]?([^'"}\s]+)['"]?/g)].map((match) => match[1]);
}
function sourceTreeProof(text = sourceTreeText) {
  assert(!text.includes('ngrok-free.dev'), 'preview ngrok origin leaked into Product source tree');
  assert(!text.includes('x-go-'), 'forbidden x-go-* override leaked into Product source tree');
  assert(!text.includes('problems/technical/'), 'technical Problem namespace leaked into Product source tree');
  assert(!/(^|\n)\s*servers\s*:/m.test(text), 'servers leaked into Product source tree');
  const sourceExtensions = [...text.matchAll(/(?:^|[\s,{])(x-[A-Za-z0-9_-]+)\s*:/gm)].map((match) => match[1]);
  for (const key of sourceExtensions) assert(ALLOWED_EXTENSIONS.has(key), `Product source tree contains non-allowlisted extension ${key}`);
  const refs = sourceRefs(text);
  assert(refs.length > 0, 'Product source tree contains no refs');
  assert(refs.every((ref) => ref.startsWith('./') || ref.startsWith('#/')), `remote/non-local Product source ref found: ${refs.find((ref) => !ref.startsWith('./') && !ref.startsWith('#/'))}`);
}
function sourceTreeNegativeControls() {
  const quotedRemote = sourceTreeText.replace("$ref: './paths-identity-portfolio-readiness.yaml#/pathItems/AccessContext'", '"$ref": "https://example.invalid/product.yaml"');
  assert(quotedRemote !== sourceTreeText, 'unable to construct quoted remote-ref source fixture');
  let failed = false;
  try { sourceTreeProof(quotedRemote); } catch { failed = true; }
  assert(failed, 'quoted remote-ref source fixture unexpectedly passed');
  negativeControlCount++;
}
function legacyRuntimeProof() {
  const tracked = run('git', ['-C', root, 'ls-files']).stdout.split(/\r?\n/).filter(Boolean);
  for (const prefix of ['apps/', 'cmd/', 'internal/', 'server/', 'backend/', 'frontend/']) assert(!tracked.some((path) => path.startsWith(prefix)), `active legacy/runtime population found under ${prefix}`);
  sameSet(tracked.filter((path) => /(^|\/)openapi[^/]*\.ya?ml$/i.test(path)), ['contracts/api/product/openapi.yaml'], 'tracked OpenAPI entrypoint authority');
}

try {
  const before = run('git', ['-C', root, 'status', '--porcelain=v1']).stdout;
  sourceTreeProof();
  sourceTreeNegativeControls();
  legacyRuntimeProof();
  const { document, bundlePath } = redoclyProof();
  const result = validate(document);
  negativeSemanticControls(document);
  typescriptProof();
  goProof(bundlePath, document);
  const after = run('git', ['-C', root, 'status', '--porcelain=v1']).stdout;
  assert(after === before, 'Product OAD proof dirtied the repository');
  assert(negativeControlCount === EXPECTED_NEGATIVE_CONTROLS, `negative-control execution count mismatch: ${negativeControlCount}/${EXPECTED_NEGATIVE_CONTROLS}`);
  console.log('product_oad_openapi=3.1.2');
  console.log(`product_oad_operations=${result.all.length}/95`);
  console.log(`product_oad_permissions=${expectedPermissions.length}/29`);
  console.log('product_oad_principal_kinds=H/A/S');
  console.log('product_oad_stable_origin=https://conexus.fun');
  console.log(`product_oad_idempotency_carriers=${result.idempotencyOps.length}/14`);
  console.log(`product_oad_collection_operations=${result.listSearch.length}/26`);
  console.log(`product_oad_negative_controls=${negativeControlCount}/${EXPECTED_NEGATIVE_CONTROLS}`);
  console.log('product_oad_generated_projection_semantics=PASS');
  console.log('product_oad_405_allow_source_expressibility=PASS');
  console.log('product_oad_generated_colon_suffix_dispatch=PASS');
  console.log('product_oad_organization_privacy_404=PASS');
  console.log('product_oad_legacy_runtime_population=0');
  console.log('product_oad_runtime_schema_enforcement=NOT_CLAIMED_D7');
  console.log('product_oad_router_selection=NONE_D7');
  console.log('product_oad_proof=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
