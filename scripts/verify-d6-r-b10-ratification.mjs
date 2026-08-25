import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const ratificationPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md');
const htmlPath = resolve(root, 'qualification/d6-r2-wireframes/b10-preparation.html');
const roadmapPath = resolve(root, 'docs/roadmap.md');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

for (const path of [ratificationPath, htmlPath, roadmapPath]) {
  assert(existsSync(path), `B10 reopen proof input missing: ${path}`);
}

const ratification = readFileSync(ratificationPath, 'utf8');
const html = readFileSync(htmlPath, 'utf8');
const roadmap = readFileSync(roadmapPath, 'utf8');

function verify(ratificationText, htmlText, roadmapText) {
  assert(ratificationText.includes('PRIOR P8: OPERATOR-RATIFIED / LOCKED'), 'B10 must preserve the historical operator LOCK');
  assert(ratificationText.includes('CURRENT P8: REOPENED / CANDIDATE'), 'B10 current P8 must be reopened candidate');
  assert(ratificationText.includes('operator-authorized bounded rebaseline'), 'B10 reopen must have operator authority');
  assert(ratificationText.includes('A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE'), 'accepted A01 debt must survive reopen');
  assert(ratificationText.includes('operator walkthrough required'), 'reopened P8 must stop for a fresh operator walkthrough');
  assert(!ratificationText.includes('CURRENT P8: OPERATOR-RATIFIED / LOCKED'), 'reopened P8 must not retain a current LOCK claim');

  assert(htmlText.includes('data-p8-status="candidate"'), 'reopened HTML must remain candidate evidence');
  assert(htmlText.includes('data-b10-role="requirements-values-handoff"'), 'reopened HTML must use simplified B10 role');
  assert(htmlText.includes('Campos para o marketplace'), 'reopened HTML must expose simplified operator model');
  assert(!htmlText.includes('Atendido'), 'reopened HTML must remove per-requirement satisfaction wording');

  assert(roadmapText.includes('P8 REOPENED / CANDIDATE'), 'roadmap must expose current reopened P8 gate');
  assert(roadmapText.includes('A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`'), 'roadmap must preserve accepted A01 debt');
  assert(roadmapText.includes('operator walkthrough'), 'roadmap next action must be fresh operator walkthrough');
}

verify(ratification, html, roadmap);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('historical lock erased', () => verify(ratification.replace('PRIOR P8: OPERATOR-RATIFIED / LOCKED', 'PRIOR P8: UNKNOWN'), html, roadmap));
expectFailure('reopen silently relocked', () => verify(ratification.replace('CURRENT P8: REOPENED / CANDIDATE', 'CURRENT P8: OPERATOR-RATIFIED / LOCKED'), html, roadmap));
expectFailure('A01 silently reopened', () => verify(ratification.replace('A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE', 'A01 = PENDING'), html, roadmap));
expectFailure('walkthrough bypassed', () => verify(ratification.replace('operator walkthrough required', 'operator walkthrough optional'), html, roadmap));

assert(negativeControls === 4, `B10 reopen negative-control count mismatch: ${negativeControls}/4`);

console.log('d6_r_b10_prior_p8=LOCKED');
console.log('d6_r_b10_current_p8=REOPENED_CANDIDATE');
console.log('d6_r_b10_A01=ACCEPT_FOR_LOCK_WITH_LATER_PROBE');
console.log('d6_r_b10_next_gate=OPERATOR_WALKTHROUGH');
console.log(`d6_r_b10_reopen_negative_controls=${negativeControls}/4`);
console.log('d6_r_b10_reopen=PASS');
