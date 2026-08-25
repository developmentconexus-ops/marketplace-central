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
  const start = '<!-- OPERATOR_SURFACE_START -->';
  const end = '<!-- OPERATOR_SURFACE_END -->';
  const startIndex = text.indexOf(start);
  const endIndex = text.indexOf(end);
  assert(startIndex >= 0 && endIndex > startIndex, 'B10 operator-surface markers missing');
  return text.slice(startIndex + start.length, endIndex);
}

function visibleText(fragment) {
  return fragment
    .replace(/<script\b[^>]*>[\s\S]*?<\/script>/giu, ' ')
    .replace(/<style\b[^>]*>[\s\S]*?<\/style>/giu, ' ')
    .replace(/<[^>]+>/gu, ' ')
    .replace(/&nbsp;/gu, ' ')
    .replace(/&amp;/gu, '&')
    .replace(/\s+/gu, ' ')
    .trim();
}

function verifyHtml(text) {
  assert(text.includes('data-p8-status="candidate"'), 'B10 must identify itself as P8 CANDIDATE');
  assert(!text.includes('data-p8-status="locked"'), 'B10 must not self-claim operator LOCK');
  assert(text.includes('data-surface="R10"'), 'B10 must identify the R10 Preparation surface');
  assert(text.includes('data-wire-prerequisite="pr68-integrated"'), 'B10 must state that the PR68 wire prerequisite is integrated');
  assert(text.includes('data-responsive-law="search-detail-mobile-stack"'), 'B10 responsive structural law missing');

  const operatorSurface = extractOperatorSurface(text);
  const operatorText = visibleText(operatorSurface);

  const expectedNavLabels = [
    'Visão geral', 'Preparação', 'Anúncios', 'Preços', 'Disponibilidade',
    'Visão operacional', 'Vendas', 'Expedição', 'Pós-venda', 'Performance',
    'Mercado', 'Economia', 'Trabalho', 'Aprovações', 'Configurações',
  ];
  for (const label of expectedNavLabels) {
    assert(text.includes(`>${label}</button>`) || text.includes(`>${label}</a>`), `locked B00 IA label missing from B10 shell: ${label}`);
  }
  assert(text.includes('>OFERTA<'), 'B10 must inherit the locked OFERTA navigation mass');
  assert(text.includes('>OPERAÇÃO<'), 'B10 must inherit the locked OPERAÇÃO navigation mass');

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
  const organizationInvalidation = text.match(/function invalidateOrganizationContext\(\)\s*\{([^}]*)\}/u);
  assert(organizationInvalidation, 'B10 Organization-switch invalidation function missing');
  assert(/installation\.value\s*=\s*['"]{2}/u.test(organizationInvalidation[1]), 'Organization switch must clear Marketplace Installation context');

  assert(text.includes('function runSearch'), 'B10 material search interaction must be implemented');
  assert(text.includes('data-search-state="known-populated"'), 'B10 known-populated search state missing');
  assert(text.includes('data-search-state="known-empty"'), 'B10 known-empty search state missing');
  assert(text.includes('Nenhum produto encontrado nesta busca conhecida.'), 'B10 known-empty human state missing');
  assert(text.includes('data-knowledge-state="unavailable"'), 'B10 unavailable search/read state missing');
  assert(text.includes('SourceInstance'), 'B10 must retain SourceInstance qualification in technical detail');
  assert(text.includes('native product key'), 'B10 must retain native product key qualification in technical detail');

  assert(text.includes('data-provider-authority="provider-authoritative"'), 'B10 must preserve provider authority for publication requirements');
  assert(text.includes('data-requirements-census="all-applicable"'), 'B10 must expose the complete applicable requirement census');
  assert(text.includes('data-provider-context="installation-category-product-type"'), 'B10 must preserve exact publication context');
  assert(text.includes('data-requirements-revision='), 'B10 must retain requirements_revision');
  assert(operatorText.includes('Requisitos do marketplace'), 'B10 must label requirements in operator language');

  const requirementRows = [...text.matchAll(/data-requirement-key="([^"]+)"/gu)];
  assert(requirementRows.length >= 7, `B10 fixture must render a meaningful complete requirement census; found ${requirementRows.length}`);

  for (const requirementClass of ['required', 'recommended', 'optional', 'conditional']) {
    assert(text.includes(`data-requirement-class="${requirementClass}"`), `B10 requirement_class fixture missing: ${requirementClass}`);
  }
  for (const applicability of ['current', 'draft_dependent']) {
    assert(text.includes(`data-evaluation-applicability="${applicability}"`), `B10 evaluation applicability fixture missing: ${applicability}`);
  }
  assert(!text.includes('data-applicability="required"'), 'B10 must not collapse requirement_class into applicability');

  for (const state of ['known', 'missing', 'conflicting', 'unknown', 'unavailable', 'unsupported']) {
    assert(text.includes(`data-source-knowledge="${state}"`), `B10 source_evidence state fixture missing: ${state}`);
  }
  assert(!text.includes('data-source-state="met"'), 'B10 must not synthesize met as source_evidence wire truth');

  for (const kind of ['text', 'exact_decimal', 'boolean', 'option', 'text_list', 'option_list', 'number_unit']) {
    assert(text.includes(`data-value-spec-kind="${kind}"`), `B10 value_spec kind fixture missing: ${kind}`);
  }
  assert(text.includes('data-not-applicable-allowed="true"'), 'B10 must expose not_applicable_allowed=true where applicable');
  assert(text.includes('data-not-applicable-allowed="false"'), 'B10 must expose not_applicable_allowed=false where applicable');

  assert(text.includes('data-source-candidate-key='), 'B10 must expose opaque source candidate identity for known/conflicting evidence');
  assert(text.includes('data-source-candidate-count="2"'), 'B10 conflicting evidence must preserve multiple candidate identities');
  assert(text.includes('data-source-media-candidates='), 'B10 must expose source_media_candidates separately from requirement rows');

  assert(text.includes('data-source-missing-policy="preserve-source-truth"'), 'B10 must preserve missing source truth without declaring publication impossible');
  assert(operatorText.includes('Falta informação'), 'B10 must translate missing source evidence into operator language');
  assert(operatorText.includes('Há informações diferentes'), 'B10 must translate conflicting source evidence into operator language');
  assert(operatorText.includes('Não foi possível verificar'), 'B10 must translate unknown/unavailable evidence into operator language');
  assert(operatorText.includes('Não disponível na fonte'), 'B10 must translate unsupported evidence into operator language');
  assert(text.includes('data-progression="listing-intent-required"'), 'B10 must admit progression to ListingIntent when no pre-ListingIntent blocker exists');
  assert(text.includes('data-progression="blocked-by-correspondence"'), 'B10 must distinguish correspondence blockers');
  assert(text.includes('data-progression="blocked-by-unknown-authority"'), 'B10 must distinguish unknown/unavailable authority blockers');
  assert(!text.includes('provider_fields'), 'B10 must not expose a raw provider field bag as Product authority');

  assert(operatorText.includes('Resumo da preparação'), 'B10 must lead with an operator-oriented preparation summary');
  for (const column of ['Requisito do marketplace', 'Exigência', 'Situação', 'Informação atual', 'O que fazer']) {
    assert(operatorText.includes(column), `B10 operator table column missing: ${column}`);
  }
  assert(operatorText.includes('Atendido'), 'B10 must show human-readable satisfied state');
  assert(operatorText.includes('Preencher ao configurar o anúncio'), 'B10 must make the next operator action explicit for missing source data');
  assert(operatorText.includes('Ver detalhes técnicos'), 'B10 must offer technical detail through secondary progressive disclosure');
  assert(text.includes('id="technicalDetails"'), 'B10 technical-detail disclosure missing');

  for (const jargon of [
    'SourceInstance', 'native product key', 'Readiness', 'value_spec', 'source_evidence',
    'not_applicable', 'Handoff', 'FOLLOW_SOURCE', 'EXPLICIT_OVERRIDE', 'ListingIntent',
    'Boundary:', 'candidate key', 'product_type_key', 'requirements_revision', 'category_key',
    'draft_dependent', 'Product truth', 'source-qualified', 'dispatchability', 'blocker',
  ]) {
    assert(!operatorText.includes(jargon), `B10 primary operator surface leaks technical jargon: ${jargon}`);
  }

  assert(text.includes('Re-leitura de prontidão necessária'), 'B10 must require reread after correspondence effect');
  assert(text.includes('function markCorrespondenceEffect'), 'B10 correspondence effect handler missing');
  assert(text.includes('function rereadReadiness'), 'B10 readiness reread handler missing');

  assert(text.includes('data-boundary="listing-intent"'), 'B10 must terminate at an explicit ListingIntent boundary');
  assert(text.includes('data-next-operation="CreateListingIntentDraft"'), 'B10 must trace but not implement the downstream ListingIntent operation');
  assert(text.includes('FOLLOW_SOURCE'), 'B10 technical trace must preserve FOLLOW_SOURCE meaning');
  assert(text.includes('EXPLICIT_OVERRIDE'), 'B10 technical trace must preserve EXPLICIT_OVERRIDE meaning');
  assert(text.includes('function showListingIntentBoundary'), 'B10 downstream boundary interaction missing');

  assert(text.includes('aria-controls="sidebar"'), 'mobile menu must identify controlled navigation');
  assert(text.includes('aria-expanded="false"'), 'mobile menu must expose collapsed state');
  assert(text.includes('id="sidebar"'), 'mobile sidebar target missing');
  assert(text.includes("event.key==='Escape'"), 'mobile drawer must close on Escape');
  assert(text.includes('menu.focus()'), 'mobile drawer must return focus to menu control');
}

function verifyStudy(text) {
  assert(text.includes(methodPin), 'B10 revalidation must cite the accepted methodology pin');
  assert(text.includes('CURRENT STRUCTURE CONFIRMED'), 'B10 structural outcome must remain explicit');
  assert(text.includes('P7 layout hypotheses remain NOT TRIGGERED'), 'B10 must preserve the no-manufactured-ambiguity disposition');
  assert(text.includes('## 11. PR #68 integrated — bounded B10 rebaseline'), 'B10 must record the post-PR68 bounded rebaseline');
  assert(text.includes(pr68Merge), 'B10 rebaseline must cite the exact integrated PR68 merge commit');
  assert(text.includes('UPSTREAM FINDING: RESOLVED'), 'B10 must supersede the earlier incorrect no-finding disposition');
  assert(text.includes('requirement_class != applicability'), 'B10 study must preserve the two independent wire dimensions');
  assert(text.includes('known / missing / conflicting / unknown / unavailable / unsupported'), 'B10 study must preserve the six source_evidence states');
  assert(text.includes('value_spec'), 'B10 study must preserve bounded safe-authoring constraints');
  assert(text.includes('not_applicable_allowed'), 'B10 study must preserve override-only not_applicable meaning');
  assert(text.includes('## 12. Operator REVISE — human language projection'), 'B10 must record the operator-language revision');
  assert(text.includes('human needs before screens'), 'B10 must cite the frontend-method principle exposed by operator review');
  assert(text.includes('technical contract remains encoded and testable'), 'B10 must preserve wire truth beneath human language');
  assert(text.includes('B00 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B00 LOCK');
  assert(text.includes('B01 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B01 LOCK');
  assert(text.includes('B00-R2 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B00-R2 LOCK');
  assert(text.includes('B11 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B11 LOCK');
  assert(text.includes('B12 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B12 LOCK');
  assert(text.includes('B110 | UNAFFECTED'), 'B10 impact sweep must explicitly preserve B110 LOCK');
  assert(text.includes('A01'), 'B10 must carry the materially depended A01 assumption into lock-time disposition');
  assert(text.includes('PENDING OPERATOR'), 'B10 A01 disposition must remain operator-owned before walkthrough');
  assert(text.includes('ACCEPT_FOR_LOCK_WITH_LATER_PROBE'), 'B10 must expose the accepted lock-time assumption option');
  assert(text.includes('BLOCK_LOCK'), 'B10 must expose the blocking assumption option');
  assert(text.includes('Operator walkthrough: PENDING'), 'B10 operator walkthrough must remain pending until actually operated');
  assert(text.includes('P8 status: CANDIDATE / NOT LOCKED'), 'B10 must remain candidate before operator adjudication');
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
expectFailure('organization context leak', () => verifyHtml(html.replace(/installation\.value\s*=\s*['"]{2}/u, "installation.value='ml-a'")));
expectFailure('premature lock', () => verifyHtml(html.replace('data-p8-status="candidate"', 'data-p8-status="locked"')));
expectFailure('requirement census collapsed', () => verifyHtml(html.replaceAll('data-requirements-census="all-applicable"', 'data-requirements-census="summary-only"')));
expectFailure('requirement class collapsed into applicability', () => verifyHtml(html.replaceAll('data-requirement-class="required"', 'data-applicability="required"')));
expectFailure('source knowledge synthesized', () => verifyHtml(html.replaceAll('data-source-knowledge="known"', 'data-source-state="met"')));
expectFailure('candidate identity erased', () => verifyHtml(html.replace(/data-source-candidate-key="[^"]+"/gu, 'data-source-candidate="anonymous"')));
expectFailure('source missing becomes blocker', () => verifyHtml(html.replace('data-source-missing-policy="preserve-source-truth"', 'data-source-missing-policy="publication-impossible"')));
expectFailure('PR68 integration forgotten', () => verifyHtml(html.replace('data-wire-prerequisite="pr68-integrated"', 'data-wire-prerequisite="missing"')));
expectFailure('operator columns regress to contract language', () => verifyHtml(html.replace('O que fazer', 'Handoff')));
expectFailure('technical jargon leaks into operator copy', () => verifyHtml(html.replace('Resumo da preparação', 'Resumo da preparação SourceInstance')));

assert(negativeControls === 12, `B10 negative-control count mismatch: ${negativeControls}/12`);

console.log('d6_r_b10_status=CANDIDATE');
console.log('d6_r_b10_structure=SEARCH_TO_EXACT_SUBJECT_DETAIL');
console.log('d6_r_b10_pr68_wire=INTEGRATED');
console.log('d6_r_b10_requirements=ALL_APPLICABLE_PROVIDER_CONTEXT');
console.log('d6_r_b10_requirement_dimensions=CLASS_PLUS_APPLICABILITY');
console.log('d6_r_b10_source_knowledge=6_STATES');
console.log('d6_r_b10_value_spec_families=7/7');
console.log('d6_r_b10_operator_language=HUMAN_FIRST');
console.log('d6_r_b10_source_missing=PRESERVE_TRUTH_CONTINUE_TO_LISTING_INTENT_WHEN_SAFE');
console.log('d6_r_b10_p7_layout=NOT_TRIGGERED');
console.log('d6_r_b10_upstream_finding=RESOLVED_BY_PR68');
console.log('d6_r_b10_assumption_A01=PENDING_OPERATOR');
console.log(`d6_r_b10_negative_controls=${negativeControls}/12`);
console.log('d6_r_b10_wireframe=PASS');
