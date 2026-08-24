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

  const expectedNavLabels = [
    'Visão geral', 'Preparação', 'Anúncios', 'Preços', 'Disponibilidade',
    'Visão operacional', 'Vendas', 'Expedição', 'Pós-venda', 'Performance',
    'Mercado', 'Economia', 'Trabalho', 'Aprovações', 'Configurações',
  ];
  for (const label of expectedNavLabels) {
    assert(text.includes(`>${label}</button>`) || text.includes(`>${label}</a>`), `locked B00 IA label missing from B10 shell: ${label}`);
  }
  assert(text.includes('>Oferta<') || text.includes('>OFERTA<') || text.includes('>Oferta</'), 'B10 must inherit the locked OFERTA navigation mass');
  assert(text.includes('>Operação<') || text.includes('>OPERAÇÃO<') || text.includes('>Operação</'), 'B10 must inherit the locked OPERAÇÃO navigation mass');

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
  const organizationInvalidation = text.match(/function invalidateOrganizationContext\(\)\{([^}]*)\}/u);
  assert(organizationInvalidation, 'B10 Organization-switch invalidation function missing');
  assert(/installation\.value\s*=\s*['"]{2}/u.test(organizationInvalidation[1]), 'Organization switch must clear Marketplace Installation context');

  assert(text.includes('function runSearch'), 'B10 material search interaction must be implemented');
  assert(!text.includes("document.getElementById('search').onclick=()=>{}"), 'B10 search button must not be a no-op');
  assert(text.includes('data-search-state="known-populated"'), 'B10 known-populated search state missing');
  assert(text.includes('data-search-state="known-empty"'), 'B10 known-empty search state missing');
  assert(text.includes('Nenhum produto encontrado nesta busca conhecida.'), 'B10 known-empty human state missing');
  assert(text.includes('data-knowledge-state="unavailable"'), 'B10 unavailable knowledge state missing');

  assert(text.includes('SourceInstance'), 'B10 must expose SourceInstance qualification');
  assert(text.includes('native product key'), 'B10 must expose native product key qualification');

  assert(text.includes('data-provider-authority="provider-authoritative"'), 'B10 must state provider-authoritative publication requirements');
  assert(text.includes('data-requirements-census="all-applicable"'), 'B10 must expose the complete applicable publication-requirement census');
  assert(text.includes('data-provider-context="installation-category-product-type"'), 'B10 requirements must remain qualified by Installation/category/product-type context');
  assert(text.includes('data-requirements-revision='), 'B10 must expose the provider/readiness requirements revision');
  assert(text.includes('Todos os requisitos aplicáveis'), 'B10 must label the complete applicable requirement set for the operator');
  assert(text.includes('Os requisitos variam conforme marketplace, categoria e product type.'), 'B10 must explain provider/context-specific requirement variation');

  const requirementRows = [...text.matchAll(/data-requirement-key="([^"]+)"/gu)];
  assert(requirementRows.length >= 5, `B10 fixture must render a meaningful requirement census; found ${requirementRows.length}`);
  for (const applicability of ['required', 'recommended', 'conditional']) {
    assert(text.includes(`data-applicability="${applicability}"`), `B10 requirement applicability fixture missing: ${applicability}`);
  }
  for (const sourceState of ['met', 'missing', 'unavailable']) {
    assert(text.includes(`data-source-state="${sourceState}"`), `B10 source-evidence state fixture missing: ${sourceState}`);
  }
  assert(text.includes('data-resolution-mode="follow-source"'), 'B10 must show FOLLOW_SOURCE as source-backed resolution meaning');
  assert(text.includes('data-listing-intent-resolution="explicit-override-eligible"'), 'B10 must show when a missing source value may be resolved later by ListingIntent EXPLICIT_OVERRIDE');
  assert(text.includes('data-source-missing-policy="listing-intent-override-eligible"'), 'B10 must not equate missing source evidence with publication impossibility');
  assert(text.includes('data-progression="listing-intent-required"'), 'B10 must allow progression to ListingIntent when source debt is override-eligible');
  assert(text.includes('data-progression="blocked-by-correspondence"'), 'B10 must distinguish true pre-ListingIntent correspondence blockers');
  assert(!text.includes('provider_fields'), 'B10 must not expose a raw provider field bag as Product authority');

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
  assert(text.includes('## 10. Operator REVISE — requirement census and provider-specific readiness'), 'B10 operator REVISE disposition missing');
  assert(text.includes('provider-authoritative'), 'B10 study must preserve provider authority for publication requirements');
  assert(text.includes('missing source != publication impossible'), 'B10 study must distinguish source insufficiency from publication impossibility');
  assert(text.includes('FOLLOW_SOURCE'), 'B10 study must preserve FOLLOW_SOURCE boundary meaning');
  assert(text.includes('EXPLICIT_OVERRIDE'), 'B10 study must preserve EXPLICIT_OVERRIDE boundary meaning');
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
expectFailure('organization context leak', () => verifyHtml(html.replace("function invalidateOrganizationContext(){installation.value='';", "function invalidateOrganizationContext(){installation.value='ml-a';")));
expectFailure('premature lock', () => verifyHtml(html.replace('data-p8-status="candidate"', 'data-p8-status="locked"')));
expectFailure('requirement census collapsed', () => verifyHtml(html.replace('data-requirements-census="all-applicable"', 'data-requirements-census="summary-only"')));
expectFailure('provider authority erased', () => verifyHtml(html.replace('data-provider-authority="provider-authoritative"', 'data-provider-authority="frontend-normalized"')));
expectFailure('missing source treated as blocker', () => verifyHtml(html.replace('data-source-missing-policy="listing-intent-override-eligible"', 'data-source-missing-policy="blocked"')));

assert(negativeControls === 7, `B10 negative-control count mismatch: ${negativeControls}/7`);

console.log('d6_r_b10_status=CANDIDATE');
console.log('d6_r_b10_structure=SEARCH_TO_EXACT_SUBJECT_DETAIL');
console.log('d6_r_b10_requirements=ALL_APPLICABLE_PROVIDER_CONTEXT');
console.log('d6_r_b10_source_missing=LISTING_INTENT_OVERRIDE_ELIGIBLE');
console.log('d6_r_b10_provider_authority=PROVIDER_AUTHORITATIVE_VIA_READINESS');
console.log('d6_r_b10_p7_layout=NOT_TRIGGERED');
console.log('d6_r_b10_upstream_finding=NONE');
console.log('d6_r_b10_assumption_A01=PENDING_OPERATOR');
console.log(`d6_r_b10_negative_controls=${negativeControls}/7`);
console.log('d6_r_b10_wireframe=PASS');