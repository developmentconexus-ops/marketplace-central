import { createHash } from 'node:crypto';
import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const entrypoint = join(root, 'contracts/api/product/openapi.yaml');
const redoclyConfig = join(root, 'contracts/api/product/redocly.yaml');
const preAuthVerifier = join(root, 'scripts/verify-product-oad-pre-auth.mjs');
const temp = mkdtempSync(join(tmpdir(), 'mpc-product-auth-proof-'));
const HTTP_METHODS = new Set(['get', 'put', 'post', 'delete', 'options', 'head', 'patch', 'trace']);
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
function operations(document) {
  const result = [];
  for (const [path, pathItem] of Object.entries(document.paths ?? {})) {
    for (const [method, operation] of Object.entries(pathItem ?? {})) {
      if (HTTP_METHODS.has(method)) result.push({ path, method, operation });
    }
  }
  return result;
}

function validateAuth(document) {
  assert(JSON.stringify(document.security) === JSON.stringify([
    { MpcHumanSessionAuth: [] },
    { MpcMachineBearerAuth: [] },
  ]), 'root Product security must allow exactly human session OR machine bearer');

  const schemes = document.components?.securitySchemes ?? {};
  sameSet(Object.keys(schemes), ['MpcHumanCsrf', 'MpcHumanSessionAuth', 'MpcMachineBearerAuth'], 'Product security scheme names');
  assert(!Object.hasOwn(schemes, 'MpcBearerAuth'), 'legacy all-client MpcBearerAuth must be absent');

  const session = schemes.MpcHumanSessionAuth;
  assert(session?.type === 'apiKey' && session?.in === 'cookie' && session?.name === '__Host-mpc_session', 'human session must be __Host-mpc_session cookie auth');

  const csrf = schemes.MpcHumanCsrf;
  assert(csrf?.type === 'apiKey' && csrf?.in === 'header' && csrf?.name === 'X-CSRF-Token', 'human CSRF carrier must be X-CSRF-Token header');

  const machine = schemes.MpcMachineBearerAuth;
  assert(machine?.type === 'http' && machine?.scheme === 'bearer', 'machine auth must remain HTTP bearer');
  assert(!Object.hasOwn(machine ?? {}, 'bearerFormat'), 'machine bearer format remains a D7 mechanism');

  const profile = document['x-mpc-authentication-profile'];
  assert(profile?.human?.principal_kind === 'H', 'human auth profile must be H only');
  assert(profile?.human?.security_scheme === 'MpcHumanSessionAuth', 'human auth profile must use MpcHumanSessionAuth');
  assert(profile?.human?.csrf?.security_scheme === 'MpcHumanCsrf', 'human unsafe-request profile must use MpcHumanCsrf');
  sameSet(profile?.human?.csrf?.required_on_methods ?? [], ['DELETE', 'PATCH', 'POST', 'PUT'], 'human CSRF unsafe methods');
  assert(profile?.machine?.security_scheme === 'MpcMachineBearerAuth', 'machine auth profile must use MpcMachineBearerAuth');
  sameSet(profile?.machine?.principal_kinds ?? [], ['A', 'S'], 'machine bearer principal kinds');
  assert(profile?.machine?.grant === 'client_credentials', 'machine bearer must be issued through Client Credentials baseline');

  const all = operations(document);
  assert(all.length === 99, `auth repair must preserve 99 Product operations, found ${all.length}`);
  for (const entry of all) {
    assert(!Object.hasOwn(entry.operation, 'security'), `${entry.operation.operationId} must not shadow the canonical split auth profile`);
    const kinds = entry.operation['x-mpc-principal-kinds'] ?? [];
    assert(kinds.length > 0 && kinds.every((kind) => ['H', 'A', 'S'].includes(kind)), `${entry.operation.operationId} principal kinds changed during auth repair`);
  }
  return all;
}

function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validateAuth(candidate); } catch { failed = true; }
  assert(failed, `auth negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

function negativeAuthControls(document) {
  expectFailure('legacy universal bearer', document, (candidate) => {
    candidate.security = [{ MpcBearerAuth: [] }];
  });
  expectFailure('browser-visible session cookie', document, (candidate) => {
    candidate.components.securitySchemes.MpcHumanSessionAuth.name = 'mpc_session';
  });
  expectFailure('missing unsafe CSRF method', document, (candidate) => {
    candidate['x-mpc-authentication-profile'].human.csrf.required_on_methods = ['POST', 'PUT', 'PATCH'];
  });
  expectFailure('human admitted through machine bearer profile', document, (candidate) => {
    candidate['x-mpc-authentication-profile'].machine.principal_kinds.push('H');
  });
  expectFailure('machine bearer becomes cookie', document, (candidate) => {
    candidate.components.securitySchemes.MpcMachineBearerAuth = { type: 'apiKey', in: 'cookie', name: 'machine' };
  });
  assert(negativeControls === 5, `auth negative-control count must be 5, found ${negativeControls}`);
}

function currentAuthProof() {
  runNpx(['--yes', '@redocly/cli@2.45.0', 'lint', entrypoint, '--config', redoclyConfig]);
  const bundleA = join(temp, 'auth-a.json');
  const bundleB = join(temp, 'auth-b.json');
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleA]);
  runNpx(['--yes', '@redocly/cli@2.45.0', 'bundle', entrypoint, '--config', redoclyConfig, '-o', bundleB]);
  assert(sha256(bundleA) === sha256(bundleB), 'auth-repaired Product bundle is not deterministic');
  const document = JSON.parse(readFileSync(bundleA, 'utf8'));
  const all = validateAuth(document);
  negativeAuthControls(document);
  console.log(`product_oad_operations=${all.length}/99`);
  console.log('product_oad_auth_human=OIDC_SERVER_SESSION_COOKIE_CSRF');
  console.log('product_oad_auth_machine=CLIENT_CREDENTIALS_BEARER_A_S');
  console.log(`product_oad_auth_negative_controls=${negativeControls}/5`);
  console.log('product_oad_auth_profile=PASS');
}

try {
  const legacy = run(process.execPath, [preAuthVerifier]);
  if (legacy.stdout?.trim()) console.log(legacy.stdout.trim());
  console.log('product_oad_pre_auth_surface_non_regression=PASS');
  currentAuthProof();
} finally {
  rmSync(temp, { recursive: true, force: true });
}
