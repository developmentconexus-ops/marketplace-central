import { readFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const d7rPath = join(root, 'docs/engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md');
const repairPath = join(root, 'docs/engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md');
const doc = readFileSync(d7rPath, 'utf8');
const repair = readFileSync(repairPath, 'utf8');
const VALIDITY_UNAVAILABLE_TYPE = 'https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable';
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function requireText(text, fragment, label) {
  assert(text.includes(fragment), `${label} missing required authority fragment: ${fragment}`);
}
function forbidText(text, fragment, label) {
  assert(!text.includes(fragment), `${label} contains superseded authority fragment: ${fragment}`);
}

function assertRatifiedAuthority(text, currentRepair) {
  requireText(text, '> **Status:** OPERATOR-RATIFIED / ACCEPTED', 'D7-R status');
  requireText(text, 'D7R_GLOBAL_MAXIMUM:ACCEPTED', 'D7-R Global Maximum');
  requireText(text, 'D7R_IDEMPOTENCY_SCOPE:ORGANIZATION_PRINCIPAL_OPERATION_KEY', 'D7-R idempotency scope');
  requireText(text, 'D7R_IDEMPOTENCY_FINGERPRINT:REQUEST_IFMATCH_OUTCOME', 'historical D7-R fingerprint snapshot');
  forbidText(text, 'D7R_IDEMPOTENCY_FINGERPRINT:PRINCIPAL_REQUEST_IFMATCH_OUTCOME', 'D7-R idempotency fingerprint');
  requireText(text, 'same raw Idempotency-Key under a different effective Principal is a different idempotency namespace', 'D7-R cross-Principal namespace law');
  requireText(text, 'D7R_REPLAY_ORDER:IDEMPOTENCY_BEFORE_IF_MATCH', 'historical D7-R replay-order snapshot');
  requireText(text, 'D7R_ELIGIBILITY:AUTHORITATIVE_Q_NOT_EVENT_CACHE', 'D7-R eligibility authority');
  requireText(text, 'D7R_VALIDITY_Q:IN_PROCESS_NO_EXTERNAL_NETWORK', 'D7-R material-validity Q');
  requireText(text, 'D7R_503:EXACT_TYPED_KNOWN_NO_EFFECT', 'D7-R semantic 503');
  requireText(text, VALIDITY_UNAVAILABLE_TYPE, 'D7-R semantic 503 Problem type');
  requireText(text, 'D7R_TERMINAL_REPLAY_RETENTION:AUTHORIZATION_HISTORY_LIFETIME', 'D7-R replay retention');
  requireText(text, 'D7R_F13:REVALIDATE_PENDING_AND_ELIGIBLE', 'D7-R F13');
  requireText(text, 'D7R_F14:DECISION_OCCURRENCE_IDEMPOTENT', 'D7-R F14');
  requireText(text, 'D7R_ZERO_DECIDER:WORK_MATERIALIZE_AND_RECONCILE', 'D7-R zero-decider Work');
  requireText(text, 'D7R_RECOVERY:EVENT_PLUS_DURABLE_PENDING_SWEEP', 'D7-R recovery sweep');
  requireText(text, 'D7R_INVALIDATION:OWNER_EVENT_NOT_SOLE_AUTHORITY', 'D7-R invalidation');
  requireText(text, 'D7R_PROOF_CLASS:STRUCTURAL_AUTHORITY_ONLY_NOT_RUNTIME_PROOF', 'D7-R proof-class fence');
  requireText(text, 'Real PostgreSQL/River/runtime execution is required when implementation proof becomes authorized', 'D7-R real-dependency proof obligation');
  forbidText(text, '**SELECTED — Global-Maximum candidate.**', 'D7-R accepted alternative');
  forbidText(text, '## 16. Candidate result', 'D7-R accepted result');

  requireText(currentRepair, '> **Status:** OPERATOR-RATIFIED / ACCEPTED', 'D5-R7 current repair');
  requireText(currentRepair, 'D5R7_W1_CARRIER:TYPED_REQUEST_ETAG', 'D5-R7 current carrier');
  requireText(currentRepair, 'D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION', 'D5-R7 current replay ordering');
  requireText(currentRepair, 'authorization_request_id\n+ typed request etag\n+ outcome', 'D5-R7 current fingerprint');
  requireText(currentRepair, 'D5R7_SUPERSEDES:D5R6_P9_D7R_D8R_CARRIER_ONLY', 'D5-R7 bounded supersession');
}

function expectFailure(name, mutateDoc, mutateRepair = (value) => value) {
  let failed = false;
  try { assertRatifiedAuthority(mutateDoc(doc), mutateRepair(repair)); } catch { failed = true; }
  assert(failed, `D7-R negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

assertRatifiedAuthority(doc, repair);

expectFailure('Principal idempotency namespace can be removed', (text) =>
  text.replace('D7R_IDEMPOTENCY_SCOPE:ORGANIZATION_PRINCIPAL_OPERATION_KEY', 'D7R_IDEMPOTENCY_SCOPE:ORGANIZATION_OPERATION_KEY'));
expectFailure('current replay-before-revision law can be removed', (text) => text, (text) =>
  text.replace('D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION', 'D5R7_REPLAY_ORDER:REVISION_BEFORE_IDEMPOTENCY'));
expectFailure('semantic 503 can lose its exact discriminator', (text) =>
  text.replace('D7R_503:EXACT_TYPED_KNOWN_NO_EFFECT', 'D7R_503:STATUS_ONLY'));
expectFailure('event/cache can become eligibility authority', (text) =>
  text.replace('D7R_ELIGIBILITY:AUTHORITATIVE_Q_NOT_EVENT_CACHE', 'D7R_ELIGIBILITY:EVENT_CACHE_AUTHORITY'));
expectFailure('architecture proof can masquerade as runtime proof', (text) =>
  text.replace('D7R_PROOF_CLASS:STRUCTURAL_AUTHORITY_ONLY_NOT_RUNTIME_PROOF', 'D7R_PROOF_CLASS:RUNTIME_PROVED'));

assert(negativeControls === 5, `D7-R negative-control count must be 5, found ${negativeControls}`);
console.log('authorization_request_d7r_ratification=ACCEPTED');
console.log('authorization_request_d7r_idempotency_scope=ORGANIZATION_PRINCIPAL_OPERATION_KEY');
console.log('authorization_request_d7r_fingerprint=REQUEST_ETAG_OUTCOME_VIA_D5R7');
console.log('authorization_request_d7r_replay_order=IDEMPOTENCY_BEFORE_REVISION_PRECONDITION');
console.log('authorization_request_d7r_currentness=AUTHORITATIVE_Q');
console.log('authorization_request_d7r_recovery=F13_F14_WORK_INVALIDATION');
console.log('authorization_request_d7r_proof_class=STRUCTURAL_AUTHORITY_ONLY');
console.log(`authorization_request_d7r_negative_controls=${negativeControls}/5`);
console.log('authorization_request_d7r=PASS');
