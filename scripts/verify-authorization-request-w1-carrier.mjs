import { readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const read = (path) => readFileSync(resolve(root, path), 'utf8');
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function section(text, start, end) {
  const startIndex = text.indexOf(start);
  assert(startIndex >= 0, `missing section start: ${start}`);
  const endIndex = end ? text.indexOf(end, startIndex + start.length) : -1;
  return text.slice(startIndex, endIndex >= 0 ? endIndex : undefined);
}

function validate({ pathYaml, w1, repair }) {
  const decide = section(pathYaml, '  DecideAuthorizationRequest:', '  AuthorizationDecisions:');
  const request = section(pathYaml, '  CreateAuthorizationDecisionRequest:', '  AuthorizationDecision:');
  const w1Custom = section(w1, '### 14.2 Owner custom method', '### 14.3 Exact revision of another referenced resource');

  assert(w1Custom.includes('does **not** use the base resource\'s ETag as `If-Match`'), 'W1 custom-method carrier law missing');
  assert(w1Custom.includes('typed technical request data'), 'W1 typed technical carrier law missing');
  assert(w1Custom.includes('`422 validation-error`'), 'W1 missing/invalid typed ETag failure grammar missing');
  assert(w1Custom.includes('`409 resource-revision-conflict`'), 'W1 stale typed ETag failure grammar missing');

  assert(decide.includes('operationId: CreateAuthorizationDecision'), 'CreateAuthorizationDecision operation missing');
  assert(decide.includes("./components.yaml#/parameters/IdempotencyKey"), 'CreateAuthorizationDecision must retain Idempotency-Key');
  assert(!decide.includes("./components.yaml#/parameters/IfMatch"), 'CreateAuthorizationDecision custom method must not use If-Match');
  assert(!decide.includes("'412':"), 'CreateAuthorizationDecision custom method must not expose 412');
  assert(!decide.includes("'428':"), 'CreateAuthorizationDecision custom method must not expose 428');
  assert(decide.includes("'409':"), 'CreateAuthorizationDecision must expose stale/state conflict as 409');
  assert(decide.includes("'422':"), 'CreateAuthorizationDecision must expose missing/invalid typed etag as 422');

  assert(request.includes('required: [etag, outcome]'), 'CreateAuthorizationDecisionRequest must require etag + outcome');
  assert(request.includes("etag: {$ref: './components.yaml#/schemas/StrongETag'}"), 'CreateAuthorizationDecisionRequest etag must use StrongETag');
  assert(request.includes('outcome: {type: string, enum: [authorize, reject]}'), 'CreateAuthorizationDecisionRequest outcome missing');
  assert(!request.includes('authorization_request_id:'), 'decision body must not duplicate request identity');
  assert(!request.includes('target:'), 'decision body must not carry governed target');

  for (const marker of [
    '> **Status:** OPERATOR-RATIFIED / ACCEPTED',
    'D5R7_W1_CARRIER:TYPED_REQUEST_ETAG',
    'D5R7_FAILURES:MISSING_INVALID_422_STALE_409',
    'D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION',
    'D5R7_B110:REVALIDATED_STRUCTURE_UNAFFECTED',
    'D5R7_SUPERSEDES:D5R6_P9_D7R_D8R_CARRIER_ONLY',
    'D5R7_NEW_PRODUCT_OPERATIONS:0',
    'D5R7_NEW_PERMISSIONS:0',
    'D5R7_NEW_BUSINESS_MEANING:0',
  ]) assert(repair.includes(marker), `D5-R7 repair authority missing marker: ${marker}`);

  assert(repair.includes('B110 LOCK disposition = REVALIDATE'), 'D5-R7 must explicitly revalidate B110 LOCK');
  assert(repair.includes('P8 reopen             = NO'), 'D5-R7 must preserve B110 structure without unnecessary P8 reopen');
  assert(repair.includes('Historical proof text is intentionally not rewritten'), 'D5-R7 must preserve prior artifacts as historical snapshots');
}

function expectFailure(name, base, mutate) {
  const candidate = structuredClone(base);
  mutate(candidate);
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `W1 carrier negative control unexpectedly passed: ${name}`);
  negativeControls++;
}

const base = {
  pathYaml: read('contracts/api/product/paths-authorization-requests.yaml'),
  w1: read('docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md'),
  repair: read('docs/engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md'),
};

validate(base);
expectFailure('custom decide regresses to If-Match header', base, (c) => {
  c.pathYaml = c.pathYaml.replace("parameters: [{$ref: './components.yaml#/parameters/IdempotencyKey'}]", "parameters: [{$ref: './components.yaml#/parameters/IfMatch'}, {$ref: './components.yaml#/parameters/IdempotencyKey'}]");
});
expectFailure('typed etag becomes optional', base, (c) => {
  c.pathYaml = c.pathYaml.replace('required: [etag, outcome]', 'required: [outcome]');
});
expectFailure('stale typed etag maps to 412', base, (c) => {
  const s = section(c.pathYaml, '  DecideAuthorizationRequest:', '  AuthorizationDecisions:');
  c.pathYaml = c.pathYaml.replace(s, s.replace("'409':", "'412':"));
});
expectFailure('bounded repair loses carrier-neutral replay law', base, (c) => {
  c.repair = c.repair.replace('D5R7_REPLAY_ORDER:IDEMPOTENCY_BEFORE_REVISION_PRECONDITION', 'D5R7_REPLAY_ORDER:REVISION_BEFORE_IDEMPOTENCY');
});

console.log('authorization_request_w1_carrier=TYPED_REQUEST_ETAG');
console.log('authorization_request_w1_failures=MISSING_INVALID_422_STALE_409');
console.log('authorization_request_w1_replay_order=IDEMPOTENCY_BEFORE_REVISION_PRECONDITION');
console.log('authorization_request_w1_b110=REVALIDATED_STRUCTURE_UNAFFECTED');
console.log(`authorization_request_w1_negative_controls=${negativeControls}/4`);
console.log('authorization_request_w1_carrier=PASS');
