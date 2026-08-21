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
const baselineVerifier = join(root, 'scripts/verify-product-oad-baseline.mjs');
const repairDoc = join(root, 'docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md');
const w4Doc = join(root, 'docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md');
const baselineMatrixText = readFileSync(w4Doc, 'utf8');
const repairText = readFileSync(repairDoc, 'utf8');
const temp = mkdtempSync(join(tmpdir(), 'mpc-product-oad-d6r1-'));
const go = process.platform === 'win32' ? 'go.exe' : 'go';
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
const PERFORMANCE_OWNER = 'MarketplacePerformanceIntelligence';
const PERFORMANCE_PERMISSION = 'performance.read';
const PERFORMANCE_START = '  # D6-R1 performance paths start';
const PERFORMANCE_END = '  # D6-R1 performance paths end';
const FORBIDDEN_PERFORMANCE_QUERY = new Set(['metrics', 'metric', 'group_by', 'dimensions', 'dimension', 'granularity', 'sort', 'order_by', 'top']);
let repairNegativeControls = 0;

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
function byId(all, id) {
  const found = all.find((entry) => entry.operation.operationId === id);
  assert(found, `operation not found: ${id}`);
  return found;
}
function schemaBySuffix(document, suffix) {
  const matches = Object.entries(document.components?.schemas ?? {}).filter(([name]) => name === suffix || name.endsWith(`_${suffix}`) || name.endsWith(`-${suffix}`));
  assert(matches.length === 1, `expected exactly one bundled schema ending ${suffix}, found ${matches.map(([name]) => name).join(', ')}`);
  return resolveRef(document, matches[0][1]);
}
function parseMatrixRows(text, start, end) {
  const startAt = text.indexOf(start);
  const endAt = text.indexOf(end, startAt + start.length);
  assert(startAt >= 0 && endAt > startAt, `unable to locate matrix ${start}`);
  const rows = [];
  for (const line of text.slice(startAt, endAt).split(/\r?\n/)) {
    const match = line.match(/^\|\s*`([^`]+)`\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$/);
    if (!match) continue;
    const [, operationId, classCell, permissionCell, principalCell] = match;
    const operationClass = classCell.includes('Q') ? 'Q' : classCell.includes('C') ? 'C' : null;
    const permission = permissionCell.match(/`([^`]+)`/)?.[1] ?? (permissionCell.includes('authenticated special condition') ? 'authenticated' : null);
    const principalKinds = [...principalCell.matchAll(/\b([HAS])\b/g)].map((m) => m[1]);
    assert(operationClass && permission && principalKinds.length, `unparseable matrix row: ${line}`);
    rows.push({ operationId, operationClass, permission, principalKinds: normalize(principalKinds) });
  }
  return rows;
}
function baselineRows() {
  const rows = parseMatrixRows(baselineMatrixText, '# 8. Exact 95-operation enforcement matrix', '\n---\n\n## 9.');
  assert(rows.length === 95, `baseline W4 row count changed: ${rows.length}`);
  return rows;
}
function repairRows() {
  const rows = parseMatrixRows(repairText, '<!-- performance-operation-matrix:start -->', '<!-- performance-operation-matrix:end -->');
  assert(rows.length === 4, `performance repair row count must be 4, found ${rows.length}`);
  return rows;
}

const acceptedBaselineRows = baselineRows();
const acceptedRepairRows = repairRows();
const acceptedRows = [...acceptedBaselineRows, ...acceptedRepairRows];
const acceptedById = new Map(acceptedRows.map((row) => [row.operationId, row]));
const expectedPermissions = normalize(acceptedRows.map((row) => row.permission).filter((permission) => permission !== 'authenticated'));
const expectedPerformanceIds = normalize(acceptedRepairRows.map((row) => row.operationId));
assert(acceptedById.size === 99, `combined operation authority must be 99 unique rows, found ${acceptedById.size}`);
assert(expectedPermissions.length === 30, `combined ordinary Permission authority must be 30, found ${expectedPermissions.length}`);
sameSet(expectedPerformanceIds, ['GetMarketplaceListingPerformance', 'GetMarketplacePerformanceSummary', 'ListMarketplaceListingPerformance', 'ListRetailMediaPerformance'], 'performance operation authority');

function baselineProof() {
  const baselineRoot = join(temp, 'baseline-root');
  const baselineContracts = join(baselineRoot, 'contracts/api/product');
  const baselineDocs = join(baselineRoot, 'docs/engineering/rebaseline');
  const baselineScripts = join(baselineRoot, 'scripts');
  mkdirSync(baselineContracts, { recursive: true });
  mkdirSync(baselineDocs, { recursive: true });
  mkdirSync(baselineScripts, { recursive: true });

  cpSync(baselineVerifier, join(baselineScripts, 'verify-product-oad.mjs'));
  for (const name of ['components.yaml', 'paths-identity-portfolio-readiness.yaml', 'paths-offering-availability-market.yaml', 'paths-economics-governance-sales-materialization.yaml', 'paths-fulfillment-postsale-work.yaml', 'redocly.yaml']) {
    cpSync(join(contractDir, name), join(baselineContracts, name));
  }
  for (const name of ['D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md', 'D5-B2-W2-SCHEMA-GRAMMAR.md', 'D5-B2-OPERATION-ADMISSION-MATRIX.md']) {
    cpSync(join(root, 'docs/engineering/rebaseline', name), join(baselineDocs, name));
  }

  let source = readFileSync(entrypoint, 'utf8');
  assert(source.includes(PERFORMANCE_START) && source.includes(PERFORMANCE_END), 'performance OAD path markers are missing');
  const startAt = source.indexOf(PERFORMANCE_START);
  const endAt = source.indexOf(PERFORMANCE_END, startAt);
  source = `${source.slice(0, startAt)}${source.slice(endAt + PERFORMANCE_END.length)}`;
  source = source.replaceAll('./paths-access-performance.yaml#/pathItems/', './paths-identity-portfolio-readiness.yaml#/pathItems/');
  writeFileSync(join(baselineContracts, 'openapi.yaml'), source, 'utf8');

  run('git', ['init'], { cwd: baselineRoot });
  run('git', ['add', '.'], { cwd: baselineRoot });
  run('git', ['-c', 'user.name=Marketplace Central Baseline Proof', '-c', 'user.email=baseline-proof@local.invalid', 'commit', '-m', 'baseline proof fixture'], { cwd: baselineRoot });

  const result = run(process.execPath, [join(baselineScripts, 'verify-product-oad.mjs')], { cwd: baselineRoot });
  console.log('product_oad_baseline_non_regression=PASS');
  if (result.stdout?.trim()) console.log(result.stdout.trim());
}

function effectiveAccessPermissionEnum(document) {
  const all = operations(document);
  const schema = responseSchema(document, byId(all, 'GetCurrentAccessContext').operation);
  const organization = resolveRef(document, schema?.properties?.organizations?.items);
  const permission = resolveRef(document, organization?.properties?.permissions?.items);
  assert(Array.isArray(permission?.enum), 'effective AccessContext Permission enum is missing');
  return permission.enum;
}
function effectiveRolePermissionEnum(document) {
  const all = operations(document);
  const schema = responseSchema(document, byId(all, 'ListAccessRoles').operation);
  const role = resolveRef(document, schema?.properties?.access_roles?.items);
  const permission = resolveRef(document, role?.properties?.permissions?.items);
  assert(Array.isArray(permission?.enum), 'effective AccessRole Permission enum is missing');
  return permission.enum;
}

function validatePerformance(document) {
  const all = operations(document);
  const ids = all.map((entry) => entry.operation.operationId);

  for (const path of Object.keys(document.paths ?? {})) {
    assert(!/(^|\/)(analytics|metrics|strategy)(\/|$|:)/i.test(path), `generic analytics/metrics/strategy Product path leaked: ${path}`);
  }

  assert(all.length === 99, `Product operation count must be 99, found ${all.length}`);
  assert(new Set(ids).size === 99, 'operationId values are not unique in repaired Product OAD');
  sameSet(ids, [...acceptedById.keys()], 'repaired OAD/accepted operation IDs');

  for (const entry of all) {
    const expected = acceptedById.get(entry.operation.operationId);
    assert(expected, `unmapped repaired operation ${entry.operation.operationId}`);
    assert(entry.operation['x-mpc-operation-class'] === expected.operationClass, `${entry.operation.operationId} class mismatch`);
    assert(entry.operation['x-mpc-required-permission'] === expected.permission, `${entry.operation.operationId} permission mismatch`);
    sameSet(entry.operation['x-mpc-principal-kinds'] ?? [], expected.principalKinds, `${entry.operation.operationId} principal kinds`);
  }

  const effectivePermissions = effectiveAccessPermissionEnum(document);
  sameSet(effectivePermissions, expectedPermissions, 'effective AccessContext Permission vocabulary');
  sameSet(effectiveRolePermissionEnum(document), expectedPermissions, 'effective AccessRole Permission vocabulary');
  sameSet(all.map((entry) => entry.operation['x-mpc-required-permission']).filter((permission) => permission !== 'authenticated'), expectedPermissions, 'operation Permission vocabulary');

  const performanceEntries = expectedPerformanceIds.map((id) => byId(all, id));
  for (const entry of performanceEntries) {
    assert(entry.operation['x-mpc-operation-class'] === 'Q', `${entry.operation.operationId} must remain Q`);
    assert(entry.operation['x-mpc-required-permission'] === PERFORMANCE_PERMISSION, `${entry.operation.operationId} must require ${PERFORMANCE_PERMISSION}`);
    assert(entry.operation['x-mpc-semantic-owner'] === PERFORMANCE_OWNER, `${entry.operation.operationId} semantic owner mismatch`);
    sameSet(entry.operation['x-mpc-principal-kinds'], ['H', 'A', 'S'], `${entry.operation.operationId} principal classes`);
    assert(entry.path.startsWith('/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/'), `${entry.operation.operationId} must be exact Installation scoped`);
    assert(response(document, entry.operation, '404'), `${entry.operation.operationId} must preserve Organization/resource privacy 404`);
    assert(response(document, entry.operation, '422'), `${entry.operation.operationId} must expose validation failure for period semantics`);

    const query = parameters(document, entry).filter((parameter) => parameter.in === 'query');
    const byName = new Map(query.map((parameter) => [parameter.name, parameter]));
    assert(byName.get('period_from')?.required === true, `${entry.operation.operationId} must require period_from`);
    assert(byName.get('period_before')?.required === true, `${entry.operation.operationId} must require period_before`);
    assert(byName.get('comparison_from')?.required === false, `${entry.operation.operationId} comparison_from must be optional`);
    assert(byName.get('comparison_before')?.required === false, `${entry.operation.operationId} comparison_before must be optional`);
    for (const name of byName.keys()) assert(!FORBIDDEN_PERFORMANCE_QUERY.has(name), `${entry.operation.operationId} leaked generic query vocabulary ${name}`);
  }

  const listingList = byId(all, 'ListMarketplaceListingPerformance');
  const mediaList = byId(all, 'ListRetailMediaPerformance');
  for (const entry of [listingList, mediaList]) {
    const query = new Map(parameters(document, entry).filter((parameter) => parameter.in === 'query').map((parameter) => [parameter.name, parameter]));
    assert(query.get('limit')?.required === false && query.get('cursor')?.required === false, `${entry.operation.operationId} must use canonical optional limit/cursor`);
  }

  const listSearch = all.filter((entry) => /^(List|Search)/.test(entry.operation.operationId));
  assert(listSearch.length === 28, `List/Search operation count must be 28, found ${listSearch.length}`);
  assert(listSearch.filter((entry) => entry.operation.operationId.startsWith('Search')).length === 1, 'repair must not add another Search operation');

  const coverage = schemaBySuffix(document, 'PerformanceEvidenceCoverage');
  const coverageStates = (coverage.oneOf ?? []).map((branch) => resolveRef(document, branch)?.properties?.state?.const).filter(Boolean);
  sameSet(coverageStates, ['complete', 'partial', 'unknown', 'unavailable', 'unsupported'], 'Performance evidence coverage states');

  const scope = schemaBySuffix(document, 'RetailMediaScope');
  const scopeKinds = (scope.oneOf ?? []).map((branch) => resolveRef(document, branch)?.properties?.kind?.const).filter(Boolean);
  sameSet(scopeKinds, ['campaign', 'listing', 'marketplace_catalog_group', 'marketplace_family_group'], 'Retail Media scope kinds');

  const count = schemaBySuffix(document, 'NonNegativeCountString');
  const countPattern = new RegExp(count.pattern ?? '');
  assert(count.type === 'string' && countPattern.test('0') && countPattern.test('18446744073709551615') && !countPattern.test('-1') && !countPattern.test('1.5'), 'NonNegativeCountString must be an integer decimal string');

  const percentage = schemaBySuffix(document, 'PercentageValue');
  assert(resolveRef(document, percentage.properties?.unit)?.const === 'percent', 'PercentageValue must carry unit=percent');
  const multiple = schemaBySuffix(document, 'MultipleValue');
  assert(resolveRef(document, multiple.properties?.unit)?.const === 'multiple', 'MultipleValue must carry unit=multiple');

  const basis = schemaBySuffix(document, 'ProviderMeasurementBasis');
  assert(resolveRef(document, basis.properties?.origin)?.const === 'provider_reported', 'ProviderMeasurementBasis must preserve provider_reported origin');

  const measures = schemaBySuffix(document, 'RetailMediaMeasures');
  assert(measures.type === 'object' && measures.additionalProperties === false && measures.minProperties === 1, 'RetailMediaMeasures must be closed and non-empty');
  for (const forbidden of ['metrics', 'metric_name', 'dimensions', 'provider_fields', 'raw_payload']) assert(!Object.hasOwn(measures.properties ?? {}, forbidden), `RetailMediaMeasures leaked generic/raw property ${forbidden}`);

  for (const forbiddenName of ['Metric', 'Analytics', 'StrategyWorkspace', 'GenericMetric']) {
    const matches = Object.keys(document.components?.schemas ?? {}).filter((name) => name === forbiddenName || name.endsWith(`_${forbiddenName}`));
    assert(matches.length === 0, `generic schema authority leaked: ${matches.join(', ')}`);
  }

  return { all, listSearch };
}

function expectPerformanceFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validatePerformance(candidate); } catch { failed = true; }
  assert(failed, `performance negative control unexpectedly passed: ${name}`);
  repairNegativeControls++;
}

function negativePerformanceControls(document) {
  expectPerformanceFailure('wrong performance Permission', document, (candidate) => {
    byId(operations(candidate), 'GetMarketplacePerformanceSummary').operation['x-mpc-required-permission'] = 'market.read';
  });
  expectPerformanceFailure('effective access omits performance.read', document, (candidate) => {
    const values = effectiveAccessPermissionEnum(candidate);
    values.splice(values.indexOf(PERFORMANCE_PERMISSION), 1);
  });
  expectPerformanceFailure('generic analytics Product path', document, (candidate) => {
    candidate.paths['/organizations/{organization_id}/analytics'] = clone(candidate.paths['/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/performance']);
  });
  expectPerformanceFailure('generic metrics selector', document, (candidate) => {
    byId(operations(candidate), 'ListRetailMediaPerformance').operation.parameters.push({ name: 'metrics', in: 'query', required: false, schema: { type: 'string' } });
  });
  expectPerformanceFailure('optional primary period', document, (candidate) => {
    const entry = byId(operations(candidate), 'GetMarketplacePerformanceSummary');
    parameters(candidate, entry).find((parameter) => parameter.name === 'period_from').required = false;
  });
  expectPerformanceFailure('performance mutation class', document, (candidate) => {
    byId(operations(candidate), 'GetMarketplaceListingPerformance').operation['x-mpc-operation-class'] = 'C';
  });
  expectPerformanceFailure('generic retail-media scope', document, (candidate) => {
    const scope = schemaBySuffix(candidate, 'RetailMediaScope');
    resolveRef(candidate, scope.oneOf[0]).properties.kind.const = 'generic';
  });
  assert(repairNegativeControls === 7, `performance negative-control count must be 7, found ${repairNegativeControls}`);
}

function fullProjectionProof() {
  runNpx(['--yes', '@redocly/cli@2.45.0', 'lint', entrypoint, '--config', redoclyConfig]);
  const bundleA = join(temp, 'full-product-a.json');
  const bundleB = join(temp, 'full-product-b.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleA]);
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleB]);
  assert(sha256(bundleA) === sha256(bundleB), 'repaired full Product bundle is not deterministic');
  const document = JSON.parse(readFileSync(bundleA, 'utf8'));
  const validated = validatePerformance(document);
  negativePerformanceControls(document);

  const tsA = join(temp, 'full-product-a.d.ts');
  const tsB = join(temp, 'full-product-b.d.ts');
  runNpx(['--yes', 'openapi-typescript@7.13.0', entrypoint, '-o', tsA, '--redocly', redoclyConfig]);
  runNpx(['--yes', 'openapi-typescript@7.13.0', entrypoint, '-o', tsB, '--redocly', redoclyConfig]);
  assert(sha256(tsA) === sha256(tsB), 'repaired TypeScript projection is not deterministic');
  const ts = readFileSync(tsA, 'utf8');
  for (const id of expectedPerformanceIds) assert(ts.includes(id), `TypeScript projection lacks ${id}`);
  assert(ts.includes(PERFORMANCE_PERMISSION), 'TypeScript projection lacks performance.read');

  const goDir = join(temp, 'full-go-proof');
  mkdirSync(goDir, { recursive: true });
  writeFileSync(join(goDir, 'oapi-codegen.yaml'), [
    'package: productapi',
    'output: product.gen.go',
    'generate:',
    '  models: true',
    '  std-http-server: true',
    'output-options:',
    '  skip-prune: true',
    '',
  ].join('\n'));
  const generate = () => run(go, ['run', 'github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0', '-config', 'oapi-codegen.yaml', bundleA], { cwd: goDir });
  generate();
  const generatedPath = join(goDir, 'product.gen.go');
  const firstHash = sha256(generatedPath);
  generate();
  assert(sha256(generatedPath) === firstHash, 'repaired Go projection is not deterministic');
  const generated = readFileSync(generatedPath, 'utf8');
  for (const id of expectedPerformanceIds) assert(generated.includes(id), `Go projection lacks ${id}`);
  writeFileSync(join(goDir, 'go.mod'), ['module productoadrepairproof', '', 'go 1.25.1', '', 'require github.com/oapi-codegen/runtime v1.7.0', ''].join('\n'));
  run(go, ['mod', 'tidy'], { cwd: goDir });
  run(go, ['test', './...'], { cwd: goDir });

  console.log('product_oad_openapi=3.1.2');
  console.log(`product_oad_operations=${validated.all.length}/99`);
  console.log(`product_oad_permissions=${expectedPermissions.length}/30`);
  console.log('product_oad_principal_kinds=H/A/S');
  console.log(`product_oad_collection_operations=${validated.listSearch.length}/28`);
  console.log(`product_oad_performance_negative_controls=${repairNegativeControls}/7`);
  console.log('product_oad_full_generated_projection_semantics=PASS');
}

try {
  baselineProof();
  fullProjectionProof();
  console.log('product_oad_performance_repair=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
