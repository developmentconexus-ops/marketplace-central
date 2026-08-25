import { createHash } from 'node:crypto';
import { mkdtempSync, readFileSync, readdirSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { basename, dirname, join, normalize, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';
import { provePublicationRequirementsOad } from './lib/publication-requirements-oad-proof.mjs';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const contractDir = join(root, 'contracts/api/product');
const entrypoint = join(contractDir, 'openapi.yaml');
const redoclyConfig = join(contractDir, 'redocly.yaml');
const allowlistPath = join(contractDir, 'source-orphan-allowlist.json');
const temp = mkdtempSync(join(tmpdir(), 'mpc-oad-source-reachability-'));
const ENFORCED_SECTIONS = new Set(['pathItems', 'schemas']);

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function npxSpec(args) {
  const npxCli = process.env.npm_execpath
    ? resolve(dirname(process.env.npm_execpath), 'npx-cli.js')
    : resolve(dirname(process.execPath), 'node_modules/npm/bin/npx-cli.js');
  return process.platform === 'win32'
    ? { command: process.execPath, args: [npxCli, ...args] }
    : { command: 'npx', args };
}
function run(command, args) {
  const result = spawnSync(command, args, { cwd: root, encoding: 'utf8', shell: false });
  if (result.error) fail(`${command} failed to start: ${result.error.message}`);
  if (result.status !== 0) fail([`${command} ${args.join(' ')} failed with exit ${result.status}`, result.stdout?.trim(), result.stderr?.trim()].filter(Boolean).join('\n'));
  return result;
}
function sha256(text) { return createHash('sha256').update(text).digest('hex'); }
function sortedUnique(values) { return [...new Set(values)].sort(); }

function parseDefinitions(file, text) {
  const lines = text.split(/\r?\n/);
  const definitions = new Map();
  let section = null;
  let current = null;

  const flush = (end) => {
    if (!current) return;
    current.text = lines.slice(current.start, end).join('\n');
    definitions.set(`${file}#/${current.section}/${current.name}`, current);
    current = null;
  };

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const top = line.match(/^([A-Za-z][A-Za-z0-9_-]*):\s*$/);
    if (top) {
      flush(i);
      section = top[1];
      continue;
    }
    if (!section) continue;
    const child = line.match(/^  ([A-Za-z0-9_.-]+):(?:\s|$)/);
    if (child) {
      flush(i);
      current = { file, section, name: child[1], start: i, text: '' };
    }
  }
  flush(lines.length);
  return definitions;
}

function rawRefs(text) {
  const refs = [];
  const regex = /\$ref:\s*['"]?([^'"\s,}\]]+)/g;
  for (const match of text.matchAll(regex)) refs.push(match[1]);
  return refs;
}

function normalizeRef(sourceFile, ref) {
  if (!ref.includes('#/')) return null;
  const [rawFile, pointer] = ref.split('#/');
  let file = sourceFile;
  if (rawFile) {
    if (/^[a-zA-Z][a-zA-Z0-9+.-]*:/.test(rawFile)) return null;
    file = basename(normalize(join(dirname(sourceFile), rawFile)));
  }
  const parts = pointer.split('/');
  if (parts[0] === 'components' && parts.length >= 3) return null;
  if (parts.length < 2) return null;
  return `${file}#/${parts[0]}/${parts[1].replaceAll('~1', '/').replaceAll('~0', '~')}`;
}

function analyze(fileTexts) {
  const defs = new Map();
  for (const [file, text] of fileTexts) {
    if (file === 'openapi.yaml' || file === 'redocly.yaml') continue;
    for (const [id, def] of parseDefinitions(file, text)) defs.set(id, def);
  }

  const reachable = new Set();
  const queue = [];
  const rootText = fileTexts.get('openapi.yaml');
  assert(rootText, 'canonical openapi.yaml missing from source reachability analysis');

  for (const ref of rawRefs(rootText)) {
    const id = normalizeRef('openapi.yaml', ref);
    if (id && defs.has(id) && !reachable.has(id)) { reachable.add(id); queue.push(id); }
  }

  while (queue.length) {
    const id = queue.shift();
    const def = defs.get(id);
    for (const ref of rawRefs(def.text)) {
      const target = normalizeRef(def.file, ref);
      if (target && defs.has(target) && !reachable.has(target)) {
        reachable.add(target);
        queue.push(target);
      }
    }
  }

  const enforced = [...defs.entries()].filter(([, def]) => ENFORCED_SECTIONS.has(def.section));
  const orphans = enforced.filter(([id]) => !reachable.has(id)).map(([id]) => id).sort();
  const orphanPathItems = orphans.filter((id) => id.includes('#/pathItems/'));
  const orphanSchemas = orphans.filter((id) => id.includes('#/schemas/'));
  return { defs, reachable, enforced, orphans, orphanPathItems, orphanSchemas };
}

function orphanSnapshotHash(result) {
  const payload = result.orphans.map((id) => `${id}\n${result.defs.get(id).text.trimEnd()}`).join('\n---\n');
  return sha256(payload);
}

function validatePolicy(result, manifest) {
  assert(manifest.policy === 'exact-frozen-historical-proof-source-debt', 'source orphan manifest policy changed');
  const expected = sortedUnique(manifest.definitions ?? []);
  assert(expected.length === (manifest.definitions ?? []).length, 'source orphan manifest contains duplicate definitions');
  assert(JSON.stringify(result.orphans) === JSON.stringify(expected), [
    'current unreachable source set differs from the exact historical-debt manifest',
    `current=${JSON.stringify(result.orphans)}`,
    `manifest=${JSON.stringify(expected)}`,
  ].join('\n'));
  assert(result.orphanPathItems.length === manifest.path_item_count, `source orphan pathItem count changed: ${result.orphanPathItems.length}/${manifest.path_item_count}`);
  assert(result.orphanSchemas.length === manifest.schema_count, `source orphan schema count changed: ${result.orphanSchemas.length}/${manifest.schema_count}`);
  assert(result.orphans.length === manifest.total_count, `source orphan total count changed: ${result.orphans.length}/${manifest.total_count}`);
  const actualHash = orphanSnapshotHash(result);
  assert(actualHash === manifest.snapshot_sha256, `frozen historical source debt content drifted: ${actualHash}/${manifest.snapshot_sha256}`);
  return actualHash;
}

try {
  const bundlePath = join(temp, 'product.json');
  const spec = npxSpec(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundlePath]);
  run(spec.command, spec.args);
  const bundleText = readFileSync(bundlePath, 'utf8');
  const bundleHash = sha256(bundleText);
  const publicationProof = provePublicationRequirementsOad(JSON.parse(bundleText));

  const files = readdirSync(contractDir).filter((name) => name.endsWith('.yaml')).sort();
  const fileTexts = new Map(files.map((file) => [file, readFileSync(join(contractDir, file), 'utf8')]));
  const manifest = JSON.parse(readFileSync(allowlistPath, 'utf8'));
  const result = analyze(fileTexts);
  const orphanHash = validatePolicy(result, manifest);

  console.log(`oad_source_bundle_sha256=${bundleHash}`);
  console.log(`oad_source_definitions=${result.defs.size}`);
  console.log(`oad_source_enforced_definitions=${result.enforced.length}`);
  console.log(`oad_source_reachable_definitions=${result.reachable.size}`);
  console.log(`oad_source_orphan_pathitems=${result.orphanPathItems.length}`);
  console.log(`oad_source_orphan_schemas=${result.orphanSchemas.length}`);
  console.log(`oad_source_orphan_snapshot_sha256=${orphanHash}`);
  console.log('oad_source_orphan_policy=EXACT_FROZEN_HISTORICAL_ALLOWLIST');
  console.log(`oad_source_allowed_orphans=${result.orphans.length}/${manifest.total_count}`);
  console.log('oad_source_new_orphans=0');
  console.log('publication_requirements_proof=PARSED_BUNDLE');
  console.log('publication_requirements_candidate_identity=TYPED_OPAQUE_KEY_VIEW');
  console.log('publication_requirements_source_value_families=7/7');
  console.log(`publication_requirements_negative_controls=${publicationProof.negativeControls}/13`);

  let negativeControls = 0;

  {
    const candidate = new Map(fileTexts);
    const source = candidate.get('paths-authorization-requests.yaml');
    assert(source.includes('\nparameters:\n'), 'source reachability pathItem control cannot locate pathItems boundary');
    candidate.set('paths-authorization-requests.yaml', source.replace('\nparameters:\n', [
      '',
      '  DefinitelyUnreachablePathItem:',
      '    get:',
      '      operationId: DefinitelyUnreachablePathItem',
      "      responses: {'200': {description: never}}",
      '',
      'parameters:',
      '',
    ].join('\n')));
    let failed = false;
    try { validatePolicy(analyze(candidate), manifest); } catch { failed = true; }
    assert(failed, 'source reachability negative control unexpectedly admitted a new orphan');
    negativeControls++;
  }

  {
    const candidate = new Map(fileTexts);
    const source = candidate.get('components.yaml');
    const needle = '  AboutBlank405Problem:';
    assert(source.includes(needle), 'historical orphan mutation control cannot locate AboutBlank405Problem');
    candidate.set('components.yaml', source.replace(needle, '  AboutBlank405Problem:\n    description: forbidden silent drift'));
    let failed = false;
    try { validatePolicy(analyze(candidate), manifest); } catch { failed = true; }
    assert(failed, 'source reachability negative control unexpectedly admitted allowlisted-orphan content drift');
    negativeControls++;
  }

  assert(negativeControls === 2, `source-reachability negative controls must be 2, found ${negativeControls}`);
  console.log(`oad_source_reachability_negative_controls=${negativeControls}/2`);
  console.log('oad_source_reachability=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
