import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const contractPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md');
const ratificationPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md');
const htmlPath = resolve(root, 'qualification/d6-r2-wireframes/b10-preparation.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(contractPath), 'B10 P9 screen contract missing');
assert(existsSync(ratificationPath), 'B10 P8 ratification missing');
assert(existsSync(htmlPath), 'B10 locked candidate evidence missing');

const contract = readFileSync(contractPath, 'utf8');
const ratification = readFileSync(ratificationPath, 'utf8');
const html = readFileSync(htmlPath, 'utf8');

function verify(contractText, ratificationText, htmlText) {
  assert(ratificationText.includes('OPERATOR-RATIFIED / LOCKED'), 'P9 requires a ratified P8 input');
  assert(contractText.includes('DERIVED / BLOCKED — P8 REOPEN REQUIRED'), 'P9 must expose the current blocked disposition');
  assert(contractText.includes('F-P9-B10-01'), 'P9 frontend finding identity missing');
  assert(contractText.includes('UPSTREAM FINDING: NONE'), 'P9 must distinguish the frontend mismatch from an upstream Product finding');
  assert(contractText.includes('`known` source evidence does not equal requirement satisfied'), 'P9 must state the exact overclaim finding');
  assert(htmlText.includes('Atendido'), 'P9 finding proof expects the locked candidate wording that was falsified');

  assert(contractText.includes('`/preparacao`'), 'B10 canonical route missing');
  for (const stateClass of ['GLOBAL_WORKSPACE_CONTEXT', 'URL_NAVIGATION_STATE', 'SERVER_STATE', 'LOCAL_EPHEMERAL']) {
    assert(contractText.includes(stateClass), `B10 P9 state class missing: ${stateClass}`);
  }
  for (const key of ['marketplace_installation_id', 'q', 'source_instance_id', 'selected_source_instance_id', 'selected_native_product_key']) {
    assert(contractText.includes(`\`${key}\``), `B10 P9 URL/navigation key missing: ${key}`);
  }

  for (const operation of [
    'SearchSourceProductsForMarketplace',
    'GetProductChannelReadiness',
    'GetPublicationRequirements',
    'ResolveProductChannelCorrespondence',
    'ClearProductChannelCorrespondence',
    'CreateListingIntentDraft',
  ]) {
    assert(contractText.includes(operation), `B10 P9 operation trace missing: ${operation}`);
  }

  for (const marker of [
    'readiness.read', 'readiness.manage', 'listing.manage',
    'ProductChannelReadiness', 'Offering',
    'correspondence_etag', 'Idempotency-Key', 'requirements_revision',
    'no blind retry', 'candidate evidence', 'Informação disponível',
  ]) {
    assert(contractText.includes(marker), `B10 P9 contract marker missing: ${marker}`);
  }

  assert(contractText.includes('B10 does not call `CreateListingIntentDraft`'), 'P9 must preserve the unopened ListingIntent boundary');
  assert(contractText.includes('frontend → backend'), 'P9 frontend-to-backend trace missing');
  assert(contractText.includes('backend → frontend'), 'P9 backend-to-frontend trace missing');
}

verify(contract, ratification, html);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('finding erased', () => verify(contract.replaceAll('F-P9-B10-01', 'F-P9-B10-XX'), ratification, html));
expectFailure('route drift', () => verify(contract.replace('`/preparacao`', '`/produto`'), ratification, html));
expectFailure('write authority drift', () => verify(contract.replaceAll('readiness.manage', 'readiness.read'), ratification, html));
expectFailure('correspondence validator erased', () => verify(contract.replaceAll('correspondence_etag', 'etag_missing'), ratification, html));
expectFailure('downstream boundary becomes B10 write', () => verify(contract.replace('B10 does not call `CreateListingIntentDraft`', 'B10 calls `CreateListingIntentDraft`'), ratification, html));

assert(negativeControls === 5, `B10 P9 negative-control count mismatch: ${negativeControls}/5`);

console.log('d6_r_b10_p9=BLOCKED_P8_REOPEN_REQUIRED');
console.log('d6_r_b10_p9_finding=F-P9-B10-01');
console.log('d6_r_b10_p9_upstream_finding=NONE');
console.log(`d6_r_b10_p9_negative_controls=${negativeControls}/5`);
console.log('d6_r_b10_p9_contract=PASS');
