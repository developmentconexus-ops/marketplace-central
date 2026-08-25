import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const contractPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md');
const ratificationPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md');
const htmlPath = resolve(root, 'qualification/d6-r2-wireframes/b10-preparation.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

for (const path of [contractPath, ratificationPath, htmlPath]) assert(existsSync(path), `B10 P9 input missing: ${path}`);

const contract = readFileSync(contractPath, 'utf8');
const ratification = readFileSync(ratificationPath, 'utf8');
const html = readFileSync(htmlPath, 'utf8');

function verify(contractText, ratificationText, htmlText) {
  assert(ratificationText.includes('CURRENT P8: REOPENED / CANDIDATE'), 'P9 must see the reopened P8 candidate');
  assert(contractText.includes('PAUSED — P8 REOPENED'), 'P9 must be paused while P8 is reopened');
  assert(contractText.includes('F-P9-B10-01'), 'P9 must preserve the finding that triggered revalidation');
  assert(contractText.includes('GLOBAL MAXIMUM REVALIDATED'), 'P9 must record the Global Maximum decision');
  assert(contractText.includes('REJECTED — `source_sufficiency`'), 'P9 must reject the unnecessary sufficiency layer');
  assert(contractText.includes('NO NEW UPSTREAM WIRE FIELD'), 'P9 must not invent a new Product field');
  assert(contractText.includes('requirements + source values + downstream authoring/provider validation'), 'P9 simplified model missing');
  assert(contractText.includes('rerun P9 after operator re-LOCK'), 'P9 must not close against an unratified candidate');

  assert(contractText.includes('`/preparacao`'), 'B10 route must remain /preparacao');
  for (const operation of [
    'SearchSourceProductsForMarketplace',
    'GetProductChannelReadiness',
    'GetPublicationRequirements',
    'ResolveProductChannelCorrespondence',
    'ClearProductChannelCorrespondence',
    'CreateListingIntentDraft',
  ]) assert(contractText.includes(operation), `B10 operation trace missing: ${operation}`);

  assert(contractText.includes('B10 does not call `CreateListingIntentDraft`'), 'B10 must preserve downstream ListingIntent boundary');
  assert(contractText.includes('frontend → backend'), 'frontend-to-backend trace marker missing');
  assert(contractText.includes('backend → frontend'), 'backend-to-frontend trace marker missing');
  assert(htmlText.includes('data-b10-role="requirements-values-handoff"'), 'P9 must target the simplified B10 candidate');
  assert(!htmlText.includes('Atendido'), 'P9 target must not carry satisfaction overclaim');
  assert(!htmlText.includes('source_sufficiency'), 'P9 target must not invent sufficiency state');
}

verify(contract, ratification, html);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('P9 closed before relock', () => verify(contract.replaceAll('PAUSED — P8 REOPENED', 'CLOSED'), ratification, html));
expectFailure('sufficiency layer returns', () => verify(contract.replaceAll('REJECTED — `source_sufficiency`', 'ACCEPTED — `source_sufficiency`'), ratification, html));
expectFailure('wire field invented', () => verify(contract.replaceAll('NO NEW UPSTREAM WIRE FIELD', 'ADD NEW UPSTREAM WIRE FIELD'), ratification, html));
expectFailure('ListingIntent boundary collapses', () => verify(contract.replaceAll('B10 does not call `CreateListingIntentDraft`', 'B10 calls `CreateListingIntentDraft`'), ratification, html));

assert(negativeControls === 4, `B10 P9 paused negative-control count mismatch: ${negativeControls}/4`);

console.log('d6_r_b10_p9=PAUSED_P8_REOPENED');
console.log('d6_r_b10_p9_global_maximum=REQUIREMENTS_VALUES_HANDOFF');
console.log('d6_r_b10_p9_source_sufficiency=REJECTED');
console.log('d6_r_b10_p9_new_wire_field=NONE');
console.log(`d6_r_b10_p9_negative_controls=${negativeControls}/4`);
console.log('d6_r_b10_p9_contract=PASS');
