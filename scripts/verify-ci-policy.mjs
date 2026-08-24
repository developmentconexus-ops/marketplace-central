import { readFileSync } from 'node:fs';

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

const ci = readFileSync('.github/workflows/ci.yml', 'utf8');
const prTitle = readFileSync('.github/workflows/pr-title.yml', 'utf8');
const pkg = JSON.parse(readFileSync('package.json', 'utf8'));
const gate = readFileSync('scripts/gate.ps1', 'utf8');

const draftCondition = "if: ${{ github.event_name == 'pull_request' && github.event.pull_request.draft == true }}";
const fullCondition = "if: ${{ github.event_name == 'push' || github.event.pull_request.draft == false }}";
const dynamicCheckName = "name: ${{ github.event_name == 'pull_request' && github.event.pull_request.draft == true && 'quick' || 'required' }}";

function stepBlock(text, anchor) {
  const lines = text.split('\n');
  const start = lines.findIndex((line) => line.includes(anchor));
  assert(start >= 0, `CI step missing: ${anchor}`);
  const block = [lines[start]];
  for (let index = start + 1; index < lines.length; index += 1) {
    if (/^      - (?:name:|uses:)/u.test(lines[index])) break;
    block.push(lines[index]);
  }
  return block.join('\n');
}

function verifyWorkflowPolicy(text) {
  assert(text.includes(dynamicCheckName), 'Draft runs must report quick while Ready/main runs report ruleset-required context required');

  const quick = stepBlock(text, '- name: quick gate');
  assert(quick.includes(draftCondition), 'quick gate must be bound to Draft pull requests');
  assert(quick.includes('run: npm run gate\n') || quick.endsWith('run: npm run gate'), 'quick gate must run npm run gate');
  assert(!quick.includes('gate:full'), 'quick gate must never run the full lane');

  const full = stepBlock(text, '- name: full gate');
  assert(full.includes(fullCondition), 'full gate must be bound to Ready/non-Draft candidates and main pushes');
  assert(full.includes('run: npm run gate:full'), 'full gate must run npm run gate:full');

  const go = stepBlock(text, 'uses: actions/setup-go@v6');
  assert(go.includes(fullCondition), 'Go setup must be bound to the full verification condition');
}

verifyWorkflowPolicy(ci);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('quick lane inverted to full', () => verifyWorkflowPolicy(
  ci.replace('run: npm run gate\n', 'run: npm run gate:full\n'),
));
expectFailure('full lane inverted to quick', () => verifyWorkflowPolicy(
  ci.replace('run: npm run gate:full\n', 'run: npm run gate\n'),
));
expectFailure('required context leaked to Draft', () => verifyWorkflowPolicy(
  ci.replace(dynamicCheckName, 'name: required'),
));
assert(negativeControls === 3, `CI policy negative-control count mismatch: ${negativeControls}/3`);

assert(!prTitle.includes('synchronize'), 'PR title check must not rerun on synchronize-only events');
assert(prTitle.includes('opened') && prTitle.includes('edited') && prTitle.includes('reopened'), 'PR title check must still cover title-changing lifecycle events');

assert(pkg.scripts?.gate === 'node scripts/verify-ci-policy.mjs && pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate.ps1 -Lane quick', 'npm run gate must be a genuinely quick lane');
assert(typeof pkg.scripts?.['gate:full'] === 'string' && pkg.scripts['gate:full'].includes('scripts/gate.ps1 -Lane full'), 'npm run gate:full must retain the aggregate lane');
assert(pkg.scripts.gate !== pkg.scripts['gate:full'], 'quick and full scripts must not be equivalent');

assert(gate.includes("if ($Lane -eq 'full')"), 'gate.ps1 must make full-only work conditional on Lane');
assert(gate.includes('quick_verifiers:'), 'gate.ps1 must report the change-aware verifier set');
assert(gate.includes('product_proof:'), 'gate.ps1 must report whether Product OAD proof ran');

console.log('ci_policy_draft_context=QUICK_ADVISORY');
console.log('ci_policy_candidate_context=REQUIRED_FULL');
console.log('ci_policy_negative_controls=3/3');
console.log('ci_policy_pr_title_synchronize=SKIPPED');
console.log('ci_policy=PASS');
