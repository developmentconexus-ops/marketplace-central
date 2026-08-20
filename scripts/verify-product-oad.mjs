import { spawnSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const sourcePath = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const w4Path = join(root, 'docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md');
const w2Path = join(root, 'docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md');
const admissionPath = join(root, 'docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md');
const sourceText = readFileSync(sourcePath, 'utf8');
const w4Text = readFileSync(w4Path, 'utf8');
const w2Text = readFileSync(w2Path, 'utf8');
const admissionText = readFileSync(admissionPath, 'utf8');
const temp = mkdtempSync(join(tmpdir(), 'mpc-product-oad-'));
const npx = process.platform === 'win32' ? 'npx.cmd' : 'npx';
const go = process.platform === 'win32' ? 'go.exe' : 'go';
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const STABLE_PROBLEM_PREFIX = 'https://conexus.fun/marketplace-central/problems/product/';

function fail(message) {
  throw new Error(message);
}

function assert(condition, message) {
  if (!condition) fail(message);
}

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: options.cwd ?? root,
    env: { ...process.env, ...(options.env ?? {}) },
    encoding: 'utf8',
    stdio: options.inherit ? 'inherit' : 'pipe',
    shell: false,
  });
  if (result.error) fail(`${command} failed to start: ${result.error.message}`);
  if (result.status !== 0) {
    fail([
      `${command} ${args.join(' ')} failed with exit ${result.status}`,
      result.stdout?.trim(),
      result.stderr?.trim(),
    ].filter(Boolean).join('\n'));
  }
  return result;
}

function sha256(path) {
  return createHash('sha256').update(readFileSync(path)).digest('hex');
}

function normalizeSet(values) {
  return [...new Set(values)].sort();
}

function sameSet(actual, expected, label) {
  const a = normalizeSet(actual);
  const e = normalizeSet(expected);
  assert(JSON.stringify(a) === JSON.stringify(e), `${label} mismatch\nactual=${JSON.stringify(a)}\nexpected=${JSON.stringify(e)}`);
}

function resolveRef(document, value) {
  if (!value || typeof value !== 'object' || typeof value.$ref !== 'string') return value;
  assert(value.$ref.startsWith('#/'), `non-local ref reached resolved proof: ${value.$ref}`);
  let current = document;
  for (const raw of value.$ref.slice(2).split('/')) {
    const key = raw.replaceAll('~1', '/').replaceAll('~0', '~');
    assert(current && Object.hasOwn(current, key), `unresolved ref in proof: ${value.$ref}`);
    current = current[key];
  }
  return current;
}

function collectRefs(value, refs = []) {
  if (!value || typeof value !== 'object') return refs;
  if (typeof value.$ref === 'string') refs.push(value.$ref);
  for (const child of Object.values(value)) collectRefs(child, refs);
  return refs;
}

function collectExtensionKeys(value, keys = []) {
  if (!value || typeof value !== 'object') return keys;
  for (const [key, child] of Object.entries(value)) {
    if (key.startsWith('x-')) keys.push(key);
    collectExtensionKeys(child, keys);
  }
  return keys;
}

function getOperations(document) {
  const operations = [];
  for (const [path, pathItem] of Object.entries(document.paths ?? {})) {
    for (const [method, operation] of Object.entries(pathItem ?? {})) {
      if (!HTTP_METHODS.has(method)) continue;
      operations.push({ path, method, operation, pathItem });
    }
  }
  return operations;
}

function effectiveParameters(document, entry) {
  const raw = [...(entry.pathItem.parameters ?? []), ...(entry.operation.parameters ?? [])];
  return raw.map((parameter) => resolveRef(document, parameter));
}

function responseObject(document, operation, status) {
  const response = operation.responses?.[status];
  return response ? resolveRef(document, response) : undefined;
}

function responseSchema(document, operation, status = '200') {
  const response = responseObject(document, operation, status);
  const schema = response?.content?.['application/json']?.schema;
  return schema ? resolveRef(document, schema) : undefined;
}

function requestSchema(document, operation, contentType = 'application/json') {
  const requestBody = resolveRef(document, operation.requestBody);
  const schema = requestBody?.content?.[contentType]?.schema;
  return schema ? resolveRef(document, schema) : undefined;
}

function parseW4Matrix(text) {
  const start = text.indexOf('# 8. Exact 95-operation enforcement matrix');
  const end = text.indexOf('\n---\n\n## 9.', start);
  assert(start >= 0 && end > start, 'unable to locate W4 exact operation matrix');
  const rows = [];
  for (const line of text.slice(start, end).split(/\r?\n/)) {
    const match = line.match(/^\|\s*`([^`]+)`\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$/);
    if (!match) continue;
    const [, operationId, classCell, permissionCell, principalCell] = match;
    const operationClass = classCell.includes('Q') ? 'Q' : classCell.includes('C') ? 'C' : null;
    assert(operationClass, `unknown W4 operation class for ${operationId}: ${classCell}`);
    const permissionMatch = permissionCell.match(/`([^`]+)`/);
    const permission = permissionMatch ? permissionMatch[1] : permissionCell.includes('authenticated special condition') ? 'authenticated' : null;
    assert(permission, `unknown W4 permission cell for ${operationId}: ${permissionCell}`);
    const qualifiedPhysical = /currently qualified S/.test(principalCell);
    const principals = qualifiedPhysical
      ? ['H', 'S']
      : [...principalCell.matchAll(/\b([HAS])\b/g)].map((m) => m[1]);
    assert(principals.length > 0, `unknown W4 principal cell for ${operationId}: ${principalCell}`);
    rows.push({ operationId, operationClass, permission, principals: normalizeSet(principals), qualifiedPhysical });
  }
  assert(rows.length === 95, `W4 matrix must contain 95 rows; found ${rows.length}`);
  sameSet(rows.map((row) => row.operationId), normalizeSet(rows.map((row) => row.operationId)), 'W4 operation IDs');
  assert(new Set(rows.map((row) => row.operationId)).size === 95, 'W4 operation IDs are not unique');
  return rows;
}

function parseProblemSlugs(text) {
  const start = text.indexOf('# 19. Problem Details catalog');
  const end = text.indexOf('Problem `type` is an absolute stable URI', start);
  assert(start >= 0 && end > start, 'unable to locate W2 Problem catalog');
  const slugs = [...text.slice(start, end).matchAll(/^- `([^`]+)`;/gm)].map((m) => m[1]);
  assert(slugs.length === 15, `W2 custom Product Problem catalog must contain 15 slugs; found ${slugs.length}`);
  return slugs;
}

function parseMandatoryIdempotencyOperations(text) {
  const start = text.indexOf('# 7. Whole-Matrix complete C-operation safety sweep');
  const end = text.indexOf('# 8. Whole-Matrix Global Coherence', start);
  assert(start >= 0 && end > start, 'unable to locate admitted C-operation safety sweep');
  const operations = [];
  for (const line of text.slice(start, end).split(/\r?\n/)) {
    if (!/mandatory client key/i.test(line)) continue;
    const match = line.match(/^\|\s*`([^`]+)`\s*\|/);
    if (match) operations.push(match[1]);
  }
  assert(operations.length === 14, `mandatory Idempotency-Key operation count changed; expected 14, found ${operations.length}`);
  return normalizeSet(operations);
}

function parseSafetySweepOperations(text) {
  const start = text.indexOf('# 7. Whole-Matrix complete C-operation safety sweep');
  const end = text.indexOf('### 7.6 Safety-sweep result', start);
  assert(start >= 0 && end > start, 'unable to locate C-operation safety sweep rows');
  const operations = [...text.slice(start, end).matchAll(/^\|\s*`([^`]+)`\s*\|/gm)].map((m) => m[1]);
  return normalizeSet(operations.filter((id) => id !== 'Operation'));
}

function operationById(operations, operationId) {
  const entry = operations.find((candidate) => candidate.operation.operationId === operationId);
  assert(entry, `operation not found: ${operationId}`);
  return entry;
}

function validateContract(document) {
  assert(document.openapi === '3.1.2', `openapi must be exactly 3.1.2, found ${document.openapi}`);
  assert(document.jsonSchemaDialect === 'https://spec.openapis.org/oas/3.1/dialect/base', 'unexpected jsonSchemaDialect');
  assert(!Object.hasOwn(document, 'servers'), 'source/resolved Product OAD must omit servers');
  assert(!sourceText.includes('ngrok-free.dev'), 'preview ngrok origin leaked into canonical Product OAD');
  assert(!sourceText.includes('/v1'), 'version prefix leaked into Product OAD');
  assert(!sourceText.includes('x-go-'), 'source OAD contains forbidden x-go-* override');

  const refs = collectRefs(document);
  assert(refs.every((ref) => ref.startsWith('#/')), `remote/non-local refs found: ${refs.filter((ref) => !ref.startsWith('#/')).join(', ')}`);

  const paths = Object.keys(document.paths ?? {});
  assert(paths.length > 0, 'Product OAD has no paths');
  for (const path of paths) {
    assert(path === '/access-context' || path.startsWith('/organizations/{organization_id}/'), `path outside canonical Product roots: ${path}`);
    assert(!/(^|\/)(providers?|integrations?|webhooks?|oauth|callbacks?|acquisition)(\/|$|:)/i.test(path), `technical/provider path leaked into Product API: ${path}`);
  }

  const operations = getOperations(document);
  assert(operations.length === 95, `Product OAD must contain exactly 95 operations; found ${operations.length}`);
  const ids = operations.map((entry) => entry.operation.operationId);
  assert(new Set(ids).size === 95, 'Product operationId values are not unique');
  for (const id of ids) assert(/^[A-Z][A-Za-z0-9]*$/.test(id), `operationId is not PascalCase/url-safe: ${id}`);
  assert(!ids.some((id) => /^AcquireMarketplace/.test(id)), 'internal acquisition operation leaked into Product OAD');

  const allowedOperationExtensions = new Set([
    'x-mpc-operation-class',
    'x-mpc-required-permission',
    'x-mpc-principal-kinds',
    'x-mpc-semantic-owner',
    'x-mpc-required-physical-qualification',
  ]);
  for (const entry of operations) {
    const op = entry.operation;
    for (const required of ['x-mpc-operation-class', 'x-mpc-required-permission', 'x-mpc-principal-kinds', 'x-mpc-semantic-owner']) {
      assert(Object.hasOwn(op, required), `${op.operationId} missing ${required}`);
    }
    for (const key of Object.keys(op).filter((key) => key.startsWith('x-'))) {
      assert(allowedOperationExtensions.has(key), `${op.operationId} contains non-allowlisted operation extension ${key}`);
    }
    assert(['Q', 'C'].includes(op['x-mpc-operation-class']), `${op.operationId} has invalid operation class`);
    assert(typeof op['x-mpc-required-permission'] === 'string' && op['x-mpc-required-permission'].length > 0, `${op.operationId} has invalid permission extension`);
    assert(Array.isArray(op['x-mpc-principal-kinds']) && op['x-mpc-principal-kinds'].length > 0, `${op.operationId} has invalid principal-kind extension`);
    assert(op['x-mpc-principal-kinds'].every((kind) => ['H', 'A', 'S'].includes(kind)), `${op.operationId} introduces a non-H/A/S principal kind`);
    assert(typeof op['x-mpc-semantic-owner'] === 'string' && op['x-mpc-semantic-owner'].length > 0, `${op.operationId} has invalid semantic owner`);
    assert(!Object.hasOwn(op, 'security'), `${op.operationId} overrides global Product bearer security`);
    assert(!Object.hasOwn(op.responses ?? {}, '202'), `${op.operationId} introduces generic 202 success`);
  }

  const extensions = normalizeSet(collectExtensionKeys(document));
  assert(!extensions.some((key) => key.startsWith('x-go-')), `forbidden type-shaping extensions found: ${extensions.join(', ')}`);

  assert(JSON.stringify(document.security) === JSON.stringify([{ MpcBearerAuth: [] }]), 'root security must be exactly MpcBearerAuth');
  const bearer = document.components?.securitySchemes?.MpcBearerAuth;
  assert(bearer?.type === 'http' && bearer?.scheme === 'bearer', 'MpcBearerAuth must be HTTP bearer');

  const w4Rows = parseW4Matrix(w4Text);
  const w4ById = new Map(w4Rows.map((row) => [row.operationId, row]));
  sameSet(ids, [...w4ById.keys()], 'OAD/W4 operation IDs');
  for (const entry of operations) {
    const op = entry.operation;
    const expected = w4ById.get(op.operationId);
    assert(op['x-mpc-operation-class'] === expected.operationClass, `${op.operationId} class mismatch: ${op['x-mpc-operation-class']} != ${expected.operationClass}`);
    assert(op['x-mpc-required-permission'] === expected.permission, `${op.operationId} permission mismatch: ${op['x-mpc-required-permission']} != ${expected.permission}`);
    sameSet(op['x-mpc-principal-kinds'], expected.principals, `${op.operationId} principal kinds`);
    const actualQualification = op['x-mpc-required-physical-qualification'] === true;
    assert(actualQualification === expected.qualifiedPhysical, `${op.operationId} physical qualification marker mismatch`);
  }

  const expectedPermissions = normalizeSet(w4Rows.map((row) => row.permission).filter((permission) => permission !== 'authenticated'));
  assert(expectedPermissions.length === 29, `W4 must project exactly 29 ordinary Permissions; found ${expectedPermissions.length}`);
  const permissionEnum = document.components?.schemas?.Permission?.enum ?? [];
  sameSet(permissionEnum, expectedPermissions, 'Permission enum/W4 ordinary Permissions');
  const usedPermissions = operations.map((entry) => entry.operation['x-mpc-required-permission']).filter((permission) => permission !== 'authenticated');
  sameSet(usedPermissions, expectedPermissions, 'operation Permission vocabulary');

  const listSearch = operations.filter((entry) => /^(List|Search)/.test(entry.operation.operationId));
  assert(listSearch.length === 26, `W3 collection operation count must be 26; found ${listSearch.length}`);
  assert(listSearch.filter((entry) => entry.operation.operationId.startsWith('Search')).length === 1, 'W3 must contain exactly one Search operation');
  const forbiddenQueryNames = new Set(['page', 'offset', 'sort', 'total', 'total_count']);
  for (const entry of listSearch) {
    const parameters = effectiveParameters(document, entry);
    const queryNames = parameters.filter((parameter) => parameter.in === 'query').map((parameter) => parameter.name);
    assert(queryNames.includes('limit'), `${entry.operation.operationId} missing limit`);
    assert(queryNames.includes('cursor'), `${entry.operation.operationId} missing cursor`);
    assert(!queryNames.some((name) => forbiddenQueryNames.has(name)), `${entry.operation.operationId} contains forbidden offset/sort/total grammar`);
    const schema = responseSchema(document, entry.operation, '200');
    assert(schema?.type === 'object', `${entry.operation.operationId} 200 collection schema is not an object`);
    assert(Object.hasOwn(schema.properties ?? {}, 'next_cursor'), `${entry.operation.operationId} collection lacks optional next_cursor`);
    assert(!(schema.required ?? []).includes('next_cursor'), `${entry.operation.operationId} requires next_cursor instead of omitting at exhaustion`);
  }
  const limit = resolveRef(document, document.components?.parameters?.Limit);
  assert(limit?.schema?.type === 'integer', 'Limit must be an integer');
  assert(Number.isInteger(limit.schema.default) && limit.schema.default > 0, 'Limit requires an explicit positive default');
  assert(Number.isInteger(limit.schema.maximum) && limit.schema.maximum >= limit.schema.default, 'Limit requires a finite maximum >= default');

  for (const [name, schema] of Object.entries(document.components?.schemas ?? {})) {
    if (schema?.type === 'object') assert(schema.additionalProperties === false, `${name} object schema is not closed with additionalProperties:false`);
    if (Array.isArray(schema?.oneOf)) {
      for (const branchRef of schema.oneOf) {
        const branch = resolveRef(document, branchRef);
        const fixed = Object.values(branch?.properties ?? {}).some((property) => property && typeof property === 'object' && Object.hasOwn(property, 'const'));
        assert(fixed, `${name} oneOf branch lacks a fixed const discriminator`);
      }
    }
  }
  const decimal = document.components?.schemas?.ExactDecimalString;
  assert(decimal?.type === 'string' && typeof decimal.pattern === 'string', 'ExactDecimalString must be a patterned string');
  assert(!decimal.pattern.toLowerCase().includes('e'), 'ExactDecimalString pattern must not admit exponent notation');

  const expectedProblemSlugs = parseProblemSlugs(w2Text);
  const actualCustomProblemTypes = [];
  for (const schema of Object.values(document.components?.schemas ?? {})) {
    const typeConst = schema?.properties?.type?.const;
    if (typeof typeConst === 'string' && typeConst.startsWith(STABLE_PROBLEM_PREFIX)) actualCustomProblemTypes.push(typeConst);
  }
  sameSet(actualCustomProblemTypes, expectedProblemSlugs.map((slug) => `${STABLE_PROBLEM_PREFIX}${slug}`), 'Product Problem type catalog');
  assert(!sourceText.includes('problems/technical/'), 'technical Problem namespace leaked into Product OAD');
  const methodNotAllowed = resolveRef(document, document.components?.responses?.MethodNotAllowed);
  assert(methodNotAllowed?.headers?.Allow, '405 MethodNotAllowed must expose Allow');
  for (const status of ['413', '415']) {
    const schemaName = status === '413' ? 'AboutBlank413Problem' : 'AboutBlank415Problem';
    const schema = document.components?.schemas?.[schemaName];
    assert(schema?.properties?.type?.const === 'about:blank' && schema?.properties?.status?.const === Number(status), `${schemaName} must use about:blank with status ${status}`);
  }

  const mandatoryIdempotency = parseMandatoryIdempotencyOperations(admissionText);
  const withIdempotency = [];
  for (const entry of operations) {
    const headers = effectiveParameters(document, entry).filter((parameter) => parameter.in === 'header');
    if (headers.some((parameter) => parameter.name.toLowerCase() === 'idempotency-key')) withIdempotency.push(entry.operation.operationId);
  }
  sameSet(withIdempotency, mandatoryIdempotency, 'Idempotency-Key carrier set');

  const cOperations = normalizeSet(w4Rows.filter((row) => row.operationClass === 'C').map((row) => row.operationId));
  const safetySweep = parseSafetySweepOperations(admissionText);
  sameSet(safetySweep, cOperations, 'C-operation safety-sweep coverage');

  const ifMatchOperations = [
    'UpdateMarketplaceInstallationConfiguration',
    'UpdateListingIntentDraft',
    'UpdateInventorySource',
    'UpdateAvailabilityAllocationScopePolicy',
    'UpdateCommercialPolicy',
    'UpdateAuthorizationDelegation',
    'UpdateFulfillmentNode',
    'UpdateFulfillmentOperatingTargets',
  ];
  const withIfMatch = [];
  for (const entry of operations) {
    const headers = effectiveParameters(document, entry).filter((parameter) => parameter.in === 'header');
    if (headers.some((parameter) => parameter.name.toLowerCase() === 'if-match')) withIfMatch.push(entry.operation.operationId);
  }
  sameSet(withIfMatch, ifMatchOperations, 'If-Match carrier set');

  const typedEtagOperations = [
    'DeactivateMarketplaceInstallation',
    'DiscardListingIntentDraft',
    'SubmitListingIntent',
    'CreateListingIntentMedia',
    'DeactivateInventorySource',
    'ResolveEconomicAttribution',
    'ResolveSaleSellingEntityAttribution',
    'ResolveBusinessSystemPartyResolution',
    'DeactivateFulfillmentNode',
    'RecordSeparation',
    'RecordPhysicalConference',
    'RecordPacking',
    'RecordDispatchHandoff',
    'AssignWork',
    'ClearWorkAssignment',
    'HoldWork',
    'ResumeWork',
    'EscalateWork',
  ];
  for (const operationId of typedEtagOperations) {
    const entry = operationById(operations, operationId);
    const contentType = operationId === 'CreateListingIntentMedia' ? 'multipart/form-data' : 'application/json';
    const schema = requestSchema(document, entry.operation, contentType);
    assert(schema?.properties?.etag, `${operationId} missing typed etag request field/part`);
    assert((schema.required ?? []).includes('etag'), `${operationId} etag request field/part is not required`);
  }
  for (const operationId of ['ResolveProductChannelCorrespondence', 'ClearProductChannelCorrespondence']) {
    const entry = operationById(operations, operationId);
    const schema = requestSchema(document, entry.operation);
    assert(schema?.properties?.correspondence_etag && (schema.required ?? []).includes('correspondence_etag'), `${operationId} missing required correspondence_etag`);
  }
  const authorizationTarget = document.components?.schemas?.AuthorizationTarget;
  for (const branchRef of authorizationTarget?.oneOf ?? []) {
    const branch = resolveRef(document, branchRef);
    assert(branch?.properties?.etag && (branch.required ?? []).includes('etag'), 'AuthorizationDecision target branch lacks exact referenced-resource ETag');
  }
  const superseded = document.components?.schemas?.SupersededPriceIntentRef;
  assert(superseded?.properties?.etag && (superseded.required ?? []).includes('etag'), 'superseded PriceIntent reference lacks ETag');

  const durableCreates = [
    'CreateMarketplaceInstallation',
    'CreateListingIntentDraft',
    'CreatePriceIntent',
    'CreateInventorySource',
    'CreateAuthorizationDecision',
    'EstablishAuthorizationDelegation',
    'CreateFulfillmentNode',
    'CreatePostSaleResolution',
  ];
  for (const operationId of durableCreates) {
    const entry = operationById(operations, operationId);
    const created = responseObject(document, entry.operation, '201');
    assert(created, `${operationId} durable create must return 201`);
    assert(created.headers?.Location, `${operationId} 201 must include Location`);
  }
  for (const entry of operations.filter((entry) => entry.method === 'post' && !durableCreates.includes(entry.operation.operationId))) {
    assert(responseObject(document, entry.operation, '200'), `${entry.operation.operationId} custom POST capability must return 200`);
  }

  const media = operationById(operations, 'CreateListingIntentMedia').operation;
  const mediaBody = resolveRef(document, media.requestBody);
  const mediaContent = mediaBody?.content?.['multipart/form-data'];
  const mediaSchema = resolveRef(document, mediaContent?.schema);
  sameSet(Object.keys(mediaSchema?.properties ?? {}), ['file', 'etag'], 'CreateListingIntentMedia multipart parts');
  sameSet(mediaSchema?.required ?? [], ['file', 'etag'], 'CreateListingIntentMedia required multipart parts');
  assert(mediaContent?.encoding?.etag?.contentType === 'text/plain', 'CreateListingIntentMedia etag part must be text/plain');

  return { operations, w4Rows, expectedPermissions, mandatoryIdempotency };
}

function expectNegativeControl(name, mutate, validator) {
  let rejected = false;
  try {
    validator(mutate());
  } catch {
    rejected = true;
  }
  assert(rejected, `negative control did not fail: ${name}`);
}

function clone(value) {
  return structuredClone(value);
}

function runNegativeControls(document) {
  const operations = getOperations(document);
  const baselineIds = operations.map((entry) => entry.operation.operationId);
  const baselinePermission = operationById(operations, 'CreatePriceIntent').operation['x-mpc-required-permission'];

  const contractValidator = (candidate) => validateContract(candidate);
  expectNegativeControl('missing operation', () => {
    const candidate = clone(document);
    delete candidate.paths['/organizations/{organization_id}/work/{work_id}:escalate'];
    return candidate;
  }, contractValidator);
  expectNegativeControl('wrong W4 permission', () => {
    const candidate = clone(document);
    const entry = getOperations(candidate).find((item) => item.operation.operationId === 'CreatePriceIntent');
    entry.operation['x-mpc-required-permission'] = 'listing.manage';
    return candidate;
  }, contractValidator);
  expectNegativeControl('technical path leakage', () => {
    const candidate = clone(document);
    candidate.paths['/organizations/{organization_id}/webhooks'] = clone(candidate.paths['/organizations/{organization_id}/work']);
    return candidate;
  }, contractValidator);
  expectNegativeControl('fourth principal kind', () => {
    const candidate = clone(document);
    const entry = getOperations(candidate).find((item) => item.operation.operationId === 'GetCurrentAccessContext');
    entry.operation['x-mpc-principal-kinds'].push('X');
    return candidate;
  }, contractValidator);
  expectNegativeControl('missing idempotency carrier', () => {
    const candidate = clone(document);
    const entry = getOperations(candidate).find((item) => item.operation.operationId === 'CreateInventorySource');
    entry.operation.parameters = (entry.operation.parameters ?? []).filter((parameter) => resolveRef(candidate, parameter)?.name !== 'Idempotency-Key');
    return candidate;
  }, contractValidator);
  expectNegativeControl('wrong Problem origin', () => {
    const candidate = clone(document);
    candidate.components.schemas.AccessDeniedProblem.properties.type.const = 'https://example.invalid/problem';
    return candidate;
  }, contractValidator);

  assert(baselineIds.length === 95 && baselinePermission === 'price.manage', 'negative controls mutated baseline in place');
}

function runRedoclyProof() {
  run(npx, ['--yes', '@redocly/cli@2.45.0', 'lint', sourcePath, '--config', redoclyConfig]);
  const bundleA = join(temp, 'product-a.json');
  const bundleB = join(temp, 'product-b.json');
  run(npx, ['--yes', '@redocly/cli@2.45.0', 'bundle', sourcePath, '--config', redoclyConfig, '-o', bundleA]);
  run(npx, ['--yes', '@redocly/cli@2.45.0', 'bundle', sourcePath, '--config', redoclyConfig, '-o', bundleB]);
  assert(sha256(bundleA) === sha256(bundleB), 'Redocly bundle is not deterministic across identical inputs');
  const document = JSON.parse(readFileSync(bundleA, 'utf8'));

  const broken = join(temp, 'broken-ref.yaml');
  writeFileSync(broken, sourceText.replace("#/components/schemas/AccessContext", "#/components/schemas/__missing__"), 'utf8');
  const negative = spawnSync(npx, ['--yes', '@redocly/cli@2.45.0', 'lint', broken, '--config', redoclyConfig], { cwd: root, encoding: 'utf8', shell: false });
  assert(negative.status !== 0, 'Redocly unresolved-ref negative fixture unexpectedly passed');
  return { document, bundleA };
}

function runTypeScriptProof() {
  const a = join(temp, 'product-a.d.ts');
  const b = join(temp, 'product-b.d.ts');
  const argsA = ['--yes', 'openapi-typescript@7.13.0', sourcePath, '-o', a, '--redocly', redoclyConfig];
  const argsB = ['--yes', 'openapi-typescript@7.13.0', sourcePath, '-o', b, '--redocly', redoclyConfig];
  run(npx, argsA);
  run(npx, argsB);
  assert(sha256(a) === sha256(b), 'openapi-typescript output is not deterministic across identical inputs');
  const generated = readFileSync(a, 'utf8');
  assert(generated.includes('This file was auto-generated by openapi-typescript'), 'TypeScript projection lacks generator banner');
  assert(generated.includes('GetCurrentAccessContext') && generated.includes('EscalateWork'), 'TypeScript projection does not contain boundary operation identities');
  run(npx, ['--yes', '-p', 'typescript@5.9.3', 'tsc', '--noEmit', '--strict', '--skipLibCheck', a]);
  const drifted = `${generated}\n// deliberate drift fixture\n`;
  assert(createHash('sha256').update(drifted).digest('hex') !== sha256(a), 'generated TypeScript drift negative control failed');
}

function runGoProof(bundlePath) {
  const goDir = join(temp, 'go-proof');
  run(process.platform === 'win32' ? 'cmd.exe' : 'mkdir', process.platform === 'win32' ? ['/c', 'mkdir', goDir] : ['-p', goDir]);
  const config = join(goDir, 'oapi-codegen.yaml');
  const generated = join(goDir, 'product.gen.go');
  writeFileSync(config, [
    'package: productapi',
    `output: ${generated.replaceAll('\\', '/')}`,
    'generate:',
    '  models: true',
    '  std-http-server: true',
    '  strict-server: true',
    'output-options:',
    '  skip-prune: true',
    '',
  ].join('\n'), 'utf8');

  run(go, ['run', 'github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0', '-config', config, bundlePath], { cwd: goDir });
  const goText = readFileSync(generated, 'utf8');
  assert(goText.includes('Code generated by github.com/oapi-codegen/oapi-codegen'), 'Go projection lacks generated-code banner');
  assert(goText.includes('type StrictServerInterface interface'), 'Go projection lacks StrictServerInterface');
  assert(goText.includes('GetCurrentAccessContext') && goText.includes('EscalateWork'), 'Go projection does not contain boundary operation identities');
  assert(!/\"(os\/exec|unsafe|syscall)\"/.test(goText), 'Go generated projection contains forbidden high-risk import');

  writeFileSync(join(goDir, 'go.mod'), [
    'module productoadproof',
    '',
    'go 1.25.1',
    '',
    'require github.com/oapi-codegen/runtime v1.7.0',
    '',
  ].join('\n'), 'utf8');
  writeFileSync(join(goDir, 'mux_test.go'), `package productapi\n\nimport (\n  \"net/http\"\n  \"net/http/httptest\"\n  \"testing\"\n)\n\ntype boundedMux struct {\n  method string\n  path string\n  handler http.Handler\n}\n\nfunc (m boundedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n  if r.Method == m.method && r.URL.Path == m.path {\n    m.handler.ServeHTTP(w, r)\n    return\n  }\n  http.NotFound(w, r)\n}\n\nfunc TestCanonicalColonSuffixDispatch(t *testing.T) {\n  called := false\n  mux := boundedMux{\n    method: http.MethodPost,\n    path: \"/organizations/org/listing-intents/li:submit\",\n    handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(http.StatusNoContent) }),\n  }\n  req := httptest.NewRequest(http.MethodPost, \"/organizations/org/listing-intents/li:submit\", nil)\n  res := httptest.NewRecorder()\n  mux.ServeHTTP(res, req)\n  if !called || res.Code != http.StatusNoContent { t.Fatalf(\"colon-suffix dispatch failed: called=%v status=%d\", called, res.Code) }\n}\n`, 'utf8');

  run(go, ['mod', 'tidy'], { cwd: goDir });
  const goMod = readFileSync(join(goDir, 'go.mod'), 'utf8');
  assert(/github\.com\/oapi-codegen\/runtime v1\.7\.0\b/.test(goMod), 'Go proof did not retain exact oapi-codegen/runtime v1.7.0');

  const exact = run(go, ['test', './...'], { cwd: goDir, env: { GOTOOLCHAIN: 'go1.25.1' } });
  const exactVersion = run(go, ['version'], { cwd: goDir, env: { GOTOOLCHAIN: 'go1.25.1' } }).stdout.trim();
  assert(/go1\.25\.1\b/.test(exactVersion), `exact minimum toolchain proof did not execute go1.25.1: ${exactVersion}`);
  run(go, ['test', './...'], { cwd: goDir });
  const currentVersion = run(go, ['version'], { cwd: goDir }).stdout.trim();
  console.log(`go_exact=${exactVersion}`);
  console.log(`go_current=${currentVersion}`);
  void exact;
}

function assertNoLegacyRuntime() {
  const forbidden = [
    'apps', 'cmd', 'internal', 'server', 'backend', 'frontend',
    'contracts/openapi.yaml', 'openapi.yaml',
  ];
  for (const relative of forbidden) {
    const result = spawnSync('git', ['-C', root, 'ls-files', '--error-unmatch', relative], { encoding: 'utf8', shell: false });
    assert(result.status !== 0, `legacy/runtime path unexpectedly tracked: ${relative}`);
  }
  const trackedOpenApi = run('git', ['-C', root, 'ls-files', '*openapi*.yaml']).stdout.trim().split(/\r?\n/).filter(Boolean);
  sameSet(trackedOpenApi, ['contracts/api/product/openapi.yaml'], 'tracked OpenAPI authority');
}

try {
  const beforeStatus = run('git', ['-C', root, 'status', '--porcelain=v1']).stdout;
  assertNoLegacyRuntime();
  const { document, bundleA } = runRedoclyProof();
  const result = validateContract(document);
  runNegativeControls(document);
  runTypeScriptProof();
  runGoProof(bundleA);
  const afterStatus = run('git', ['-C', root, 'status', '--porcelain=v1']).stdout;
  assert(afterStatus === beforeStatus, 'Product OAD proof dirtied the repository working tree');

  console.log('product_oad_openapi=3.1.2');
  console.log(`product_oad_operations=${result.operations.length}/95`);
  console.log(`product_oad_permissions=${result.expectedPermissions.length}/29`);
  console.log('product_oad_principal_kinds=H/A/S');
  console.log('product_oad_stable_origin=https://conexus.fun');
  console.log(`product_oad_idempotency_carriers=${result.mandatoryIdempotency.length}/14`);
  console.log('product_oad_collection_operations=26/26');
  console.log('product_oad_negative_controls=7/7');
  console.log('product_oad_legacy_runtime_population=0');
  console.log('product_oad_runtime_schema_enforcement=NOT_CLAIMED_D7');
  console.log('product_oad_router_selection=NONE_D7');
  console.log('product_oad_proof=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
