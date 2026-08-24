import { existsSync, readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const d8rPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-NOTIF-01-D8-R-AUTHORIZATION-REQUEST-GOLDEN-FLOW-REVALIDATION.md');
const d7rPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md');
const repairPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md');
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function requireText(text, fragment, label) {
  assert(text.includes(fragment), `${label} missing required authority fragment: ${fragment}`);
}
function assertD7Prerequisite(text, repair) {
  requireText(text, '> **Status:** OPERATOR-RATIFIED / ACCEPTED', 'D7-R prerequisite');
  requireText(text, 'D7R_IDEMPOTENCY_SCOPE:ORGANIZATION_PRINCIPAL_OPERATION_KEY', 'D7-R prerequisite');
  requireText(text, 'D7R_REPLAY_ORDER:IDEMPOTENCY_BEFORE_IF_MATCH', 'historical D7-R replay snapshot');
  requireText(text, 'D7R_PROOF_CLASS:STRUCTURAL_AUTHORITY_ONLY_NOT_RUNTIME_PROOF', 'D7-R prerequisite');
  requireText(repair, 'D5R7_W1_CARRIER:TYPED_REQUEST_ETAG', 'D5-R7 current carrier');
  requireText(repair, 'D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION', 'D5-R7 current replay law');
  requireText(repair, 'D5R7_SUPERSEDES:D5R6_P9_D7R_D8R_CARRIER_ONLY', 'D5-R7 bounded supersession');
}
function assertD8Authority(text) {
  const required = [
    '> **Status:** OPERATOR-RATIFIED / ACCEPTED',
    'D8R_GLOBAL_MAXIMUM:ACCEPTED',
    'D8R_GOLDEN_FLOW_SET:THREE_BUSINESS_PLUS_SR01_UNCHANGED',
    'D8R_PRODUCT_SURFACE:106_OPERATIONS_31_PERMISSIONS_HAS',
    'D8R_GF01:REVALIDATED',
    'D8R_GF02:REVALIDATED',
    'D8R_GF03:NOT_MATERIALLY_AFFECTED',
    'D8R_SR01:REVALIDATED_CONTINUOUS_TIMELINE_VS_PITR',
    'D8R_REVIEW_BASIS:LISTING_INTENT_PRICE_INTENT_BUSINESS_ORDER_INTENT_INVOICING_INTENT',
    'D8R_TYPED_503:KNOWN_NO_EFFECT_ONLY_FOR_EXACT_PROBLEM_TYPE',
    'D8R_NOTIFICATIONS:F13_CURRENTNESS_F14_DECISION_OCCURRENCE',
    'D8R_ZERO_DECIDER:WORK_NOT_AUTHORITY',
    'D8R_PITR:IDEMPOTENCY_IS_NOT_CONTINUITY_ORACLE',
    'D8R_LIVE_PROBES:NOT_REOPENED',
    'D8R_P2:PRESERVED_FUTURE_PROOF',
    'D8R_NEW_PRODUCT_OPERATIONS:0',
    'D8R_NEW_PERMISSIONS:0',
    'D8R_NEW_GOLDEN_FLOWS:0',
    'D8R_NEW_INFRASTRUCTURE:0',
    'D8R_RUNTIME_PROOF:POST_D9_IMPLEMENTATION',
  ];
  for (const fragment of required) requireText(text, fragment, 'D8-R authority');

  requireText(text, 'listing_intent', 'D8-R GF-01 mapping');
  requireText(text, 'price_intent', 'D8-R GF-01 mapping');
  requireText(text, 'business_order_intent', 'D8-R GF-02 mapping');
  requireText(text, 'invoicing_intent', 'D8-R GF-02 mapping');
  requireText(text, 'authorization-validity-unavailable', 'D8-R semantic 503');
  requireText(text, 'same raw Idempotency-Key under a different effective Principal', 'D8-R Principal isolation');
  requireText(text, 'recovery fence', 'D8-R PITR composition');
}
function expectFailure(name, fn) {
  let failed = false;
  try { fn(); } catch { failed = true; }
  assert(failed, `D8-R negative control unexpectedly passed: ${name}`);
  negativeControls++;
}
function negativeProof(text) {
  expectFailure('GF-04 can replace the unchanged golden-flow set', () => {
    assertD8Authority(text.replace('D8R_NEW_GOLDEN_FLOWS:0', 'D8R_NEW_GOLDEN_FLOWS:1'));
  });
  expectFailure('GF-03 can be silently reopened', () => {
    assertD8Authority(text.replace('D8R_GF03:NOT_MATERIALLY_AFFECTED', 'D8R_GF03:REVALIDATED'));
  });
  expectFailure('live probes can be reopened by ceremony', () => {
    assertD8Authority(text.replace('D8R_LIVE_PROBES:NOT_REOPENED', 'D8R_LIVE_PROBES:REOPENED'));
  });
  expectFailure('PITR can treat idempotency as continuity proof', () => {
    assertD8Authority(text.replace('D8R_PITR:IDEMPOTENCY_IS_NOT_CONTINUITY_ORACLE', 'D8R_PITR:IDEMPOTENCY_IS_CONTINUITY_ORACLE'));
  });
  expectFailure('current Product census can regress to the old D8 baseline', () => {
    assertD8Authority(text.replace('D8R_PRODUCT_SURFACE:106_OPERATIONS_31_PERMISSIONS_HAS', 'D8R_PRODUCT_SURFACE:99_OPERATIONS_30_PERMISSIONS_HAS'));
  });
  assert(negativeControls === 5, `D8-R negative-control count must be 5, found ${negativeControls}`);
}

assert(existsSync(d7rPath), 'accepted D7-R authority document missing');
assert(existsSync(d8rPath), 'D8-R authority document missing; operator-ratified D8-R has not been recorded');
assert(existsSync(repairPath), 'D5-R7 W1 carrier repair missing');
const d7r = readFileSync(d7rPath, 'utf8');
const d8r = readFileSync(d8rPath, 'utf8');
const repair = readFileSync(repairPath, 'utf8');
assertD7Prerequisite(d7r, repair);
assertD8Authority(d8r);
negativeProof(d8r);

console.log('authorization_request_d8r_ratification=ACCEPTED');
console.log('authorization_request_d8r_golden_flow_set=3_BUSINESS_PLUS_SR01_UNCHANGED');
console.log('authorization_request_d8r_affected=GF01_GF02_SR01');
console.log('authorization_request_d8r_gf03=NOT_MATERIALLY_AFFECTED');
console.log('authorization_request_d8r_product_surface=106_31_HAS');
console.log('authorization_request_d8r_replay_order=IDEMPOTENCY_BEFORE_REVISION_PRECONDITION_VIA_D5R7');
console.log('authorization_request_d8r_live_probes=NOT_REOPENED');
console.log('authorization_request_d8r_pitr=IDEMPOTENCY_NOT_CONTINUITY_ORACLE');
console.log('authorization_request_d8r_proof_class=STRUCTURAL_AUTHORITY_ONLY');
console.log(`authorization_request_d8r_negative_controls=${negativeControls}/5`);
console.log('authorization_request_d8r=PASS');
