import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const artifactPath = resolve(root, 'qualification/d6-r2-wireframes/b10-preparation.html');
const studyPath = resolve(root, 'docs/engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(artifactPath), 'B10 rendered artifact missing');
assert(existsSync(studyPath), 'B10 P6 reference/revalidation study missing');

const html = readFileSync(artifactPath, 'utf8');
const study = readFileSync(studyPath, 'utf8');
const methodPin = '9c7210d1504bef01c0d134a6c3ae8627deebb535';

function verifyHtml(text) {
  assert(text.includes('data-p8-status="candidate"'), 'B10 must identify itself as P8 CANDIDATE');
  assert(!text.includes('data-p8-status="locked"'), 'B10 must not self-claim operator LOCK');
  assert(text.includes('data-surface="R10"'), 'B10 must identify the R10 Preparation surface');
  assert(text.includes('data-responsive-law="search-detail-mobile-stack"'), 'B10 responsive structural law missing');

  for (const operationId of [
    'SearchSourceProductsForMarketplace',
    'GetProductChannelReadiness',
    'GetPublicationRequirements',
    'ResolveProductChannelCorrespondence',
    'ClearProductChannelCorrespondence',
  ]) {
    assert(text.includes(`data-operation="${operationId}"`), `B10 operation trace missing: ${operationId}`);
  }
  assert(text.includes('data-owner="ProductChannelReadiness"'), 'B10 semantic owner trace missing');
  assert(text.includes('data-required-permission="readiness.read"'), 'B10 readiness.read trace missing');
  assert(text.includes('data-write-permission="readiness.manage"'), 'B10 readiness.manage trace missing');

  assert(text.includes('id="organization"'), 'B10 Organization control must be explicit');
  assert(text.includes('id="installation"'), 'B10 exact Marketplace Installation control missing');
  assert(text.includes('function invalidateOrganizationContext'), 'B10 Organization-switch invalidation function missing');
  assert(/invalidateOrganizationContext\([^)]*\)[\s\S]*installation\.value\s*=\s*['"]{2}/u.test(text), 'Organization switch must clear Marketplace Installation context');

  assert(text.includes('function runSearch'), 'B10 material search interaction must be implemented');
  assert(!text.includes("document.getElementById('search').onclick=()=>{}"), 'B10 search button must not be a no-op');
  assert(text.includes('data-search-state="known-populated"'), 'B10 known-populated search state missing');
  assert(text.includes('data-search-state="known-empty"'), 'B10 known-empty search state missing');
  assert(text.includes('Nenhum produto encontrado nesta busca conhecida.'), 'B10 known-empty human state missing');
  assert(text.includes('data-knowledge-state="unavailable"'), 'B10 unavailable knowledge state missing');

  assert(text.includes('SourceInstance'), 'B10 must expose SourceInstance qualification');
  assert(text.includes('native product key'), 'B10 must expose native product key qualification');
  assert(text.includes('Re-leitura de prontidão necessária'), 'B10 must require reread after correspondence effect');
  assert(text.includes('function markCorrespondenceEffect'), 'B10 correspondence effect handler missing');
  assert(text.includes('function rereadReadiness'), 'B10 readiness reread handler missing');

  assert(text.includes('data-boundary="listing-intent"'), 'B10 must terminate at an explicit ListingIntent boundary');
  assert(text.includes('data-next-operation="CreateListingIntentDraft"'), 'B10 ListingIntent boundary must trace the downstream operation without implementing it');
  assert(text.includes('function showListingIntentBoundary'), 'B10 downstream boundary interaction missing');

  assert(text.includes('aria-controls="sidebar"'), 'mobile menu must identify controlled navigation');
  assert(text.includes('aria-expanded="false"'), 'mobile menu must expose collapsed state');
  assert(text.includes('id="sidebar"'), 'mobile sidebar target missing');
  assert(text.includes("event.key==='Escape'"), 'mobile drawer must close on Escape');
  assert(text.includes('menu.focus()'), 'mobile drawer must return focus to menu control');
}

function verifyStudy(text) {
  assert(text.includes('## 9. Post-methodology adoption revalidation'), 'B10 post-methodology revalidation section missing');
  assert(text.includes(methodPin), 'B10 revalidation must cite the accepted methodology pin');
  assert(text.includes('CURRENT STRUCTURE CONFIRMED'), 'B10 revalidation outcome must be explicit');
  assert(text.includes('P7 layout hypotheses remain NOT TRIGGERED'), 'B10 must preserve the no-manufactured-ambiguity disposition');
  for (const label of ['Required fields / summaries', 'Identity sources', 'Pagination / scale', 'Sort / filter', 'Preview / content truth', 'Material writes']) {
    assert(text.includes(label), `B10 P7 feasibility disposition missing: ${label}`);
  }
  assert(text.includes('A01'), 'B10 must carry the materially depended A01 assumption into lock-time disposition');
  assert(text.includes('PENDING OPERATOR'), 'B10 A01 lock-time disposition must remain operator-owned before walkthrough');
  assert(text.includes('ACCEPT_FOR_LOCK_WITH_LATER_PROBE'), 'B10 must expose the accepted lock-time assumption option');
  assert(text.includes('BLOCK_LOCK'), 'B10 must expose the blocking assumption option');
  assert(text.includes('Operator walkthrough: PENDING'), 'B10 operator walkthrough must remain pending until actually operated');
  assert(text.includes('P8 status: CANDIDATE / NOT LOCKED'), 'B10 must remain candidate before operator adjudication');
  assert(text.includes('UPSTREAM FINDING: NONE'), 'B10 revalidation must state whether Product authority was falsified');
}

verifyHtml(html);
verifyStudy(study);

let negativeControls = 0;
function expectFailure(name, body) {
  let failed = false;
  try { body(); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('search no-op', () => verifyHtml(html.replace('function runSearch', 'function missingRunSearch')));
expectFailure('known-empty collapsed', () => verifyHtml(html.replace('data-search-state="known-empty"', 'data-search-state="known-populated"')));
expectFailure('organization context leak', () => verifyHtml(html.replace(/installation\.value\s*=\s*['"]{2}/u, 'installation.value=\'ml-a\'')));
expectFailure('premature lock', () => verifyHtml(html.replace('data-p8-status="candidate"', 'data-p8-status="locked"')));

assert(negativeControls === 4, `B10 negative-control count mismatch: ${negativeControls}/4`);

console.log('d6_r_b10_status=CANDIDATE');
console.log('d6_r_b10_structure=SEARCH_TO_EXACT_SUBJECT_DETAIL');
console.log('d6_r_b10_p7_layout=NOT_TRIGGERED');
console.log('d6_r_b10_upstream_finding=NONE');
console.log('d6_r_b10_assumption_A01=PENDING_OPERATOR');
console.log(`d6_r_b10_negative_controls=${negativeControls}/4`);
console.log('d6_r_b10_wireframe=PASS');
