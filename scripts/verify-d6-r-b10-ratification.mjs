import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const ratificationPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md');
const htmlPath = resolve(root, 'qualification/d6-r2-wireframes/b10-preparation.html');
const roadmapPath = resolve(root, 'docs/roadmap.md');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

for (const path of [ratificationPath, htmlPath, roadmapPath]) {
  assert(existsSync(path), `B10 ratification proof input missing: ${path}`);
}

const ratification = readFileSync(ratificationPath, 'utf8');
const html = readFileSync(htmlPath, 'utf8');
const roadmap = readFileSync(roadmapPath, 'utf8');

function verify(ratificationText, htmlText, roadmapText) {
  assert(ratificationText.includes('OPERATOR-RATIFIED / LOCKED'), 'B10 ratification must record operator LOCK');
  assert(ratificationText.includes('final disposition: LOCK'), 'B10 walkthrough must end in LOCK');
  assert(ratificationText.includes('A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE'), 'B10 A01 lock-time disposition missing');
  assert(ratificationText.includes('No blocking frontend or upstream Finding remains'), 'B10 P8 lock must state blocking-finding disposition');
  assert(ratificationText.includes('P9 — B10 Screen Contract + bidirectional backend trace'), 'B10 LOCK must route to P9');

  // The operated low-fi remains candidate evidence; ratification, not the HTML, is the operator LOCK carrier.
  assert(htmlText.includes('data-p8-status="candidate"'), 'B10 operated HTML evidence must remain the exact candidate form');
  assert(!htmlText.includes('data-p8-status="locked"'), 'B10 HTML must not self-author operator LOCK');
  assert(htmlText.includes('Resumo da preparação'), 'B10 locked evidence lost the human-first summary');
  assert(htmlText.includes('Requisito do marketplace'), 'B10 locked evidence lost the human-first requirement language');

  assert(/B10[^|\n]{0,100}P8 OPERATOR-RATIFIED \/ LOCKED/u.test(roadmapText), 'roadmap must expose current B10 P8 LOCK');
  assert(roadmapText.includes('A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`'), 'roadmap must expose accepted A01 debt');
  assert(roadmapText.includes('P9 — B10 Screen Contract + bidirectional backend trace'), 'roadmap next action must be P9');
}

verify(ratification, html, roadmap);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('operator lock erased', () => verify(ratification.replace('final disposition: LOCK', 'final disposition: REVISE'), html, roadmap));
expectFailure('A01 silently reopened', () => verify(ratification.replace('A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE', 'A01 = PENDING'), html, roadmap));
expectFailure('roadmap regresses to pending', () => verify(ratification, html, roadmap.replace('P8 OPERATOR-RATIFIED / LOCKED', 'P8 NOT LOCKED')));

assert(negativeControls === 3, `B10 ratification negative-control count mismatch: ${negativeControls}/3`);

console.log('d6_r_b10_p8=LOCKED');
console.log('d6_r_b10_operator=RATIFIED');
console.log('d6_r_b10_A01=ACCEPT_FOR_LOCK_WITH_LATER_PROBE');
console.log(`d6_r_b10_ratification_negative_controls=${negativeControls}/3`);
console.log('d6_r_b10_ratification=PASS');
