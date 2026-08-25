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
const pr68Merge = 'ed3d164b0574b7950c2c7467d150c89576bba1ec';

function extractOperatorSurface(text) {
  const start = text.indexOf('<!-- OPERATOR_SURFACE_START -->');
  const end = text.indexOf('<!-- OPERATOR_SURFACE_END -->');
  assert(start >= 0 && end > start, 'B10 operator-surface markers missing');
  return text
    .slice(start, end)
    .replace(/<script[\s\S]*?<\/script>/gu, ' ')
    .replace(/<style[\s\S]*?<\/style>/gu, ' ')
    .replace(/<[^>]+>/gu, ' ')
    .replace(/\s+/gu, ' ');
}

function verifyHtml(text) {
  assert(text.includes('data-p8-status="candidate"'), 'B10 must remain a P8 candidate after bounded reopen');
  assert(!text.includes('data-p8-status="locked"'), 'B10 HTML must never self-author operator LOCK');
  assert(text.includes('data-surface="R10"'), 'B10 must identify R10');
  assert(text.includes('data-wire-prerequisite="pr68-integrated"'), 'PR68 prerequisite must remain integrated');
  assert(text.includes('data-b10-role="requirements-values-handoff"'), 'B10 simplified role marker missing');
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
  assert(text.includes('data-owner="ProductChannelReadiness"'), 'B10 owner trace missing');
  assert(text.includes('data-required-permission="readiness.read"'), 'B10 readiness.read trace missing');
  assert(text.includes('data-write-permission="readiness.manage"'), 'B10 readiness.manage trace missing');

  assert(text.includes('id="organization"'), 'B10 Organization control missing');
  assert(text.includes('id="installation"'), 'B10 exact Marketplace Installation control missing');
  const organizationInvalidation = text.match(/function invalidateOrganizationContext\(\)\s*\{([^}]*)\}/u);
  assert(organizationInvalidation, 'B10 Organization-switch invalidation function missing');
  assert(/installation\.value\s*=\s*['"]{2}/u.test(organizationInvalidation[1]), 'Organization switch must clear Marketplace Installation');

  assert(text.includes('function runSearch'), 'B10 material search interaction missing');
  assert(text.includes('data-search-state="known-populated"'), 'B10 known-populated search state missing');
  assert(text.includes('data-search-state="known-empty"'), 'B10 known-empty search state missing');
  assert(text.includes('data-knowledge-state="unavailable"'), 'B10 unavailable search/read state missing');

  assert(text.includes('data-provider-authority="provider-authoritative"'), 'provider authority marker missing');
  assert(text.includes('data-requirements-census="all-applicable"'), 'complete requirement census marker missing');
  assert(text.includes('data-provider-context="installation-category-product-type"'), 'publication context marker missing');
  assert(text.includes('data-requirements-revision='), 'requirements_revision marker missing');

  const requirementRows = [...text.matchAll(/data-requirement-key="([^"]+)"/gu)];
  assert(requirementRows.length >= 7, `B10 fixture must render a meaningful requirement census; found ${requirementRows.length}`);
  for (const requirementClass of ['required', 'recommended', 'optional', 'conditional']) {
    assert(text.includes(`data-requirement-class="${requirementClass}"`), `requirement_class fixture missing: ${requirementClass}`);
  }
  for (const applicability of ['current', 'draft_dependent']) {
    assert(text.includes(`data-evaluation-applicability="${applicability}"`), `applicability fixture missing: ${applicability}`);
  }
  for (const state of ['known', 'missing', 'conflicting', 'unknown', 'unavailable', 'unsupported']) {
    assert(text.includes(`data-source-knowledge="${state}"`), `source_evidence fixture missing: ${state}`);
  }
  for (const kind of ['text', 'exact_decimal', 'boolean', 'option', 'text_list', 'option_list', 'number_unit']) {
    assert(text.includes(`data-value-spec-kind="${kind}"`), `value_spec fixture missing: ${kind}`);
  }
  assert(text.includes('data-source-candidate-key='), 'opaque source candidate identity missing');
  assert(text.includes('data-source-candidate-count="2"'), 'conflicting source candidates must remain distinct');
  assert(text.includes('data-source-media-candidates='), 'source media candidates must remain separate');
  assert(text.includes('data-source-missing-policy="preserve-source-truth"'), 'missing source truth must remain non-blocking by itself');
  assert(text.includes('data-progression="listing-intent-required"'), 'B10 must hand off to ListingIntent configuration');
  assert(!text.includes('provider_fields'), 'raw provider field bag must not become Product authority');
  assert(!text.includes('source_sufficiency'), 'B10 must not invent a source_sufficiency layer');

  const operatorText = extractOperatorSurface(text);
  for (const phrase of [
    'Campos para o marketplace',
    'Campo do marketplace', 'Exigência', 'Valor encontrado', 'Na configuração do anúncio',
    'Obrigatório', 'Recomendado', 'Opcional', 'Condicional',
    'Preencher na configuração do anúncio', 'Escolher na configuração do anúncio',
    'Continuar para configurar o anúncio', 'Ver detalhes técnicos',
  ]) {
    assert(operatorText.includes(phrase), `B10 simplified operator language missing: ${phrase}`);
  }
  for (const overclaim of ['Atendido', 'requisitos atendidos', '>Situação<', 'Pronto pela fonte']) {
    assert(!text.includes(overclaim), `B10 must not expose per-requirement satisfaction semantics: ${overclaim}`);
  }

  for (const jargon of [
    'SourceInstance', 'native product key', 'Readiness', 'value_spec', 'source_evidence',
    'not_applicable', 'Handoff', 'FOLLOW_SOURCE', 'EXPLICIT_OVERRIDE', 'ListingIntent',
    'candidate key', 'product_type_key', 'requirements_revision', 'category_key',
    'draft_dependent', 'dispatchability', 'source_sufficiency',
  ]) {
    assert(!operatorText.includes(jargon), `B10 primary operator surface leaks technical jargon: ${jargon}`);
  }

  assert(text.includes('function markCorrespondenceEffect'), 'correspondence effect handler missing');
  assert(text.includes('function rereadReadiness'), 'authoritative reread handler missing');
  assert(text.includes('data-boundary="listing-intent"'), 'explicit ListingIntent boundary missing');
  assert(text.includes('data-next-operation="CreateListingIntentDraft"'), 'downstream ListingIntent trace missing');
  assert(text.includes('FOLLOW_SOURCE'), 'technical trace must preserve FOLLOW_SOURCE');
  assert(text.includes('EXPLICIT_OVERRIDE'), 'technical trace must preserve EXPLICIT_OVERRIDE');
  assert(text.includes('function showListingIntentBoundary'), 'downstream boundary interaction missing');

  assert(text.includes('aria-controls="sidebar"'), 'mobile menu target missing');
  assert(text.includes("event.key==='Escape'"), 'mobile drawer Escape behavior missing');
  assert(text.includes('menu.focus()'), 'mobile drawer focus return missing');
}

function verifyStudy(text) {
  assert(text.includes(methodPin), 'B10 study must cite accepted methodology pin');
  assert(text.includes(pr68Merge), 'B10 study must cite integrated PR68 merge');
  assert(text.includes('GLOBAL MAXIMUM REVALIDATED'), 'B10 must record the Global Maximum revalidation');
  assert(text.includes('REJECTED — `source_sufficiency`'), 'B10 must explicitly reject the unnecessary sufficiency layer');
  assert(text.includes('requirements + source values + downstream authoring/provider validation'), 'B10 simplified target model missing');
  assert(text.includes('AnyMarket'), 'B10 research evidence must include AnyMarket');
  assert(text.includes('Hub2b'), 'B10 research evidence must include Hub2b');
  assert(text.includes('Magis5'), 'B10 research evidence must include Magis5');
  assert(text.includes('Mercado Livre'), 'B10 research evidence must include Mercado Livre');
  assert(text.includes('NO NEW UPSTREAM WIRE FIELD'), 'B10 must reject a new Product wire field for per-requirement sufficiency');
  assert(text.includes('P8 REOPENED / CANDIDATE'), 'B10 study must mark the operator-authorized bounded reopen');
  for (const lock of ['B00 | UNAFFECTED', 'B01 | UNAFFECTED', 'B00-R2 | UNAFFECTED', 'B11 | UNAFFECTED', 'B12 | UNAFFECTED', 'B110 | UNAFFECTED']) {
    assert(text.includes(lock), `B10 impact sweep missing: ${lock}`);
  }
  assert(text.includes('A01 = ACCEPT_FOR_LOCK_WITH_LATER_PROBE'), 'accepted A01 debt must survive the bounded reopen');
  assert(text.includes('next gate: operator walkthrough'), 'B10 must stop at the new operator walkthrough gate');
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

expectFailure('satisfaction status reintroduced', () => verifyHtml(html.replace('Valor encontrado', 'Atendido')));
expectFailure('sufficiency layer invented', () => verifyHtml(html.replace('data-b10-role="requirements-values-handoff"', 'data-b10-role="source_sufficiency"')));
expectFailure('source missing becomes blocker', () => verifyHtml(html.replace('data-source-missing-policy="preserve-source-truth"', 'data-source-missing-policy="publication-impossible"')));
expectFailure('requirement census collapsed', () => verifyHtml(html.replaceAll('data-requirements-census="all-applicable"', 'data-requirements-census="summary-only"')));
expectFailure('candidate identity erased', () => verifyHtml(html.replace(/data-source-candidate-key="[^"]+"/gu, 'data-source-candidate="anonymous"')));
expectFailure('global maximum evidence erased', () => verifyStudy(study.replace('GLOBAL MAXIMUM REVALIDATED', 'LOCAL PATCH')));

assert(negativeControls === 6, `B10 negative-control count mismatch: ${negativeControls}/6`);

console.log('d6_r_b10_status=CANDIDATE_REOPENED');
console.log('d6_r_b10_model=REQUIREMENTS_VALUES_HANDOFF');
console.log('d6_r_b10_source_sufficiency=REJECTED_OVERENGINEERING');
console.log('d6_r_b10_provider_validation=DOWNSTREAM');
console.log('d6_r_b10_A01=ACCEPT_FOR_LOCK_WITH_LATER_PROBE');
console.log(`d6_r_b10_negative_controls=${negativeControls}/6`);
console.log('d6_r_b10_wireframe=PASS');
