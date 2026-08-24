import { readFileSync } from 'node:fs';

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

const ci = readFileSync('.github/workflows/ci.yml', 'utf8');
const prTitle = readFileSync('.github/workflows/pr-title.yml', 'utf8');
const pkg = JSON.parse(readFileSync('package.json', 'utf8'));
const gate = readFileSync('scripts/gate.ps1', 'utf8');

assert(ci.includes("github.event.pull_request.draft == true"), 'CI must select quick verification for Draft PRs');
assert(ci.includes("github.event_name == 'push' || github.event.pull_request.draft == false"), 'CI must select full verification for non-Draft candidates and main pushes');
assert(ci.includes('run: npm run gate\n'), 'CI quick path must run npm run gate');
assert(ci.includes('run: npm run gate:full'), 'CI full path must run npm run gate:full');
assert(ci.includes('setup-go@v6') && ci.includes("github.event_name == 'push' || github.event.pull_request.draft == false"), 'Go setup must be reserved for full verification');

assert(!prTitle.includes('synchronize'), 'PR title check must not rerun on synchronize-only events');
assert(prTitle.includes('opened') && prTitle.includes('edited') && prTitle.includes('reopened'), 'PR title check must still cover title-changing lifecycle events');

assert(pkg.scripts?.gate === 'node scripts/verify-ci-policy.mjs && pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate.ps1 -Lane quick', 'npm run gate must be a genuinely quick lane');
assert(typeof pkg.scripts?.['gate:full'] === 'string' && pkg.scripts['gate:full'].includes('scripts/gate.ps1 -Lane full'), 'npm run gate:full must retain the aggregate lane');
assert(pkg.scripts.gate !== pkg.scripts['gate:full'], 'quick and full scripts must not be equivalent');

assert(gate.includes("if ($Lane -eq 'full')"), 'gate.ps1 must make full-only work conditional on Lane');
assert(gate.includes('quick_verifiers:'), 'gate.ps1 must report the change-aware verifier set');
assert(gate.includes('product_proof:'), 'gate.ps1 must report whether Product OAD proof ran');

console.log('ci_policy_draft_lane=QUICK');
console.log('ci_policy_candidate_lane=FULL');
console.log('ci_policy_pr_title_synchronize=SKIPPED');
console.log('ci_policy=PASS');
