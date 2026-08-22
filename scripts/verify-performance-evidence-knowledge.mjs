import { mkdirSync, mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const temp = mkdtempSync(join(tmpdir(), 'mpc-performance-knowledge-'));
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
function coverageStates(document, schema) {
  const coverage = resolveRef(document, schema.properties?.coverage);
  return (coverage?.oneOf ?? []).map((branch) => resolveRef(document, branch)?.properties?.state?.const).filter(Boolean);
}
function assertClosedObject(schema, name) {
  assert(schema?.type === 'object' && schema.additionalProperties === false, `${name} must be a closed object`);
}
function assertRequired(schema, names, label) {
  for (const name of names) assert((schema.required ?? []).includes(name), `${label} must require ${name}`);
}
function assertForbiddenProperties(schema, names, label) {
  for (const name of names) assert(!Object.hasOwn(schema.properties ?? {}, name), `${label} must forbid ${name}`);
}
function assertCarrier(document, name, measureNames, options = {}) {
  const wrapper = schemaBySuffix(document, name);
  assert(Array.isArray(wrapper.oneOf) && wrapper.oneOf.length === 2, `${name} must split available/unavailable knowledge with oneOf`);
  const available = schemaBySuffix(document, `${name}Available`);
  const unavailable = schemaBySuffix(document, `${name}Unavailable`);
  assertClosedObject(available, `${name}Available`);
  assertClosedObject(unavailable, `${name}Unavailable`);
  assertRequired(available, options.baseRequired ?? ['coverage', 'provenance'], `${name}Available`);
  assertRequired(unavailable, options.baseRequired ?? ['coverage', 'provenance'], `${name}Unavailable`);
  sameSet(coverageStates(document, available), ['complete', 'partial'], `${name} available coverage`);
  sameSet(coverageStates(document, unavailable), ['unknown', 'unavailable', 'unsupported'], `${name} unavailable coverage`);
  if (options.anyMeasure) {
    const alternatives = (available.anyOf ?? []).map((branch) => normalize(branch.required ?? [])).filter((required) => required.length > 0);
    sameSet(alternatives.map((required) => required.join(',')), measureNames.map((name) => name), `${name} available measure alternatives`);
  } else {
    assertRequired(available, measureNames, `${name}Available`);
  }
  assertForbiddenProperties(unavailable, measureNames, `${name}Unavailable`);
}
function validate(document) {
  assertCarrier(document, 'TrafficPerformance', ['visits']);
  assertCarrier(document, 'TrafficPerformanceComparison', ['visits'], { baseRequired: ['period', 'coverage', 'provenance'] });
  assertCarrier(document, 'SalesActivityPerformance', ['sales_count', 'units_count']);
  assertCarrier(document, 'SalesActivityPerformanceComparison', ['sales_count', 'units_count'], { baseRequired: ['period', 'coverage', 'provenance'] });
  assertCarrier(document, 'RetailMediaPerformance', ['measures'], { baseRequired: ['coverage', 'measurement_basis', 'provenance'] });
  assertCarrier(document, 'RetailMediaPerformanceComparison', ['measures'], { baseRequired: ['period', 'coverage', 'measurement_basis', 'provenance'] });
  assertCarrier(document, 'RetailMediaSummary', ['spend', 'return_on_ad_spend'], { baseRequired: ['coverage', 'measurement_basis', 'provenance'], anyMeasure: true });
  assertCarrier(document, 'RetailMediaSummaryComparison', ['spend', 'return_on_ad_spend'], { baseRequired: ['period', 'coverage', 'measurement_basis', 'provenance'], anyMeasure: true });
}
function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `performance knowledge negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

try {
  mkdirSync(temp, { recursive: true });
  const bundlePath = join(temp, 'product.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundlePath]);
  const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
  validate(document);
  expectFailure('available traffic without visits requirement', document, (candidate) => {
    const schema = schemaBySuffix(candidate, 'TrafficPerformanceAvailable');
    schema.required = (schema.required ?? []).filter((name) => name !== 'visits');
  });
  expectFailure('unavailable retail media carries measures', document, (candidate) => {
    const schema = schemaBySuffix(candidate, 'RetailMediaPerformanceUnavailable');
    schema.properties.measures = {$ref: '#/components/schemas/__invalid__'};
  });
  assert(negativeControls === 2, `performance knowledge negative-control count mismatch: ${negativeControls}/2`);
  console.log('performance_evidence_available_requires_measure=PASS');
  console.log('performance_evidence_unavailable_forbids_measure=PASS');
  console.log(`performance_evidence_negative_controls=${negativeControls}/2`);
  console.log('performance_evidence_knowledge_proof=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
