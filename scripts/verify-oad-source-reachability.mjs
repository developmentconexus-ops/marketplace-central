import { createHash } from 'node:crypto';
import { mkdtempSync, readFileSync, readdirSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { basename, dirname, join, normalize, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const contractDir = join(root, 'contracts/api/product');
const entrypoint = join(contractDir, 'openapi.yaml');
const redoclyConfig = join(contractDir, 'redocly.yaml');
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
  if (parts[0] === 'components' && parts.length >= 3) return null; // root-local aliases are not source-definition ids
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

try {
  const bundlePath = join(temp, 'product.json');
  const spec = npxSpec(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundlePath]);
  run(spec.command, spec.args);
  const bundleText = readFileSync(bundlePath, 'utf8');
  const bundleHash = sha256(bundleText);

  const files = readdirSync(contractDir).filter((name) => name.endsWith('.yaml')).sort();
  const fileTexts = new Map(files.map((file) => [file, readFileSync(join(contractDir, file), 'utf8')]));
  const result = analyze(fileTexts);

  console.log(`oad_source_bundle_sha256=${bundleHash}`);
  console.log(`oad_source_definitions=${result.defs.size}`);
  console.log(`oad_source_enforced_definitions=${result.enforced.length}`);
  console.log(`oad_source_reachable_definitions=${result.reachable.size}`);
  console.log(`oad_source_orphan_pathitems=${result.orphanPathItems.length}`);
  console.log(`oad_source_orphan_schemas=${result.orphanSchemas.length}`);
  if (result.orphans.length) {
    fail(`unreachable OAD source definitions (${result.orphans.length}):\n${result.orphans.join('\n')}`);
  }

  // GREEN-only falsifiers: once the live graph is clean, prove the detector itself fires.
  let negativeControls = 0;
  for (const [name, injection] of [
    ['orphan pathItem', '\n  DefinitelyUnreachablePathItem:\n    get:\n      operationId: DefinitelyUnreachablePathItem\n      responses: {\'200\': {description: never}}\n'],
    ['orphan schema', '\n  DefinitelyUnreachableSchema:\n    type: object\n    additionalProperties: false\n'],
  ]) {
    const candidate = new Map(fileTexts);
    const targetFile = name.includes('pathItem') ? 'paths-authorization-requests.yaml' : 'components.yaml';
    const marker = name.includes('pathItem') ? '\nschemas:\n' : null;
    if (marker) {
      const source = candidate.get(targetFile);
      candidate.set(targetFile, source.replace(marker, `${injection}${marker}`));
    } else {
      candidate.set(targetFile, `${candidate.get(targetFile).trimEnd()}${injection}\n`);
    }
    const mutated = analyze(candidate);
    assert(mutated.orphans.length > 0, `source-reachability negative control unexpectedly passed: ${name}`);
    negativeControls++;
  }
  assert(negativeControls === 2, `source-reachability negative controls must be 2, found ${negativeControls}`);
  console.log(`oad_source_reachability_negative_controls=${negativeControls}/2`);
  console.log('oad_source_reachability=PASS');
} finally {
  rmSync(temp, { recursive: true, force: true });
}
