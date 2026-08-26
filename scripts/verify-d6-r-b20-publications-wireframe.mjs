import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const artifact = resolve(process.cwd(), 'qualification/d6-r2-wireframes/b20-publications.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

function validate(html) {
  assert(html.includes('data-p8-status="candidate"'), 'B20 HTML must remain a P8 candidate; LOCK is operator-only');
  assert(!html.includes('data-p8-status="locked"'), 'B20 HTML must not self-claim operator LOCK');
  assert(html.includes('data-surface="R20-R21"'), 'B20 surface binding missing');
  assert(html.includes('data-wire-prerequisite="pr70-integrated"'), 'B20 must bind the integrated PR #70 read projections');

  // Presentation is the human projection; canonical ref never is, and never becomes write authority.
  assert(html.includes('data-read-projection="typed-presentation"'), 'typed presentation read binding missing');
  assert(html.includes('data-presentation-fields="display_name"'), 'presentation projection must be display_name');
  assert(html.includes('data-presentation-write="never"'), 'presentation must never become write authority');
  assert(html.includes("presentation.state==='known'?l.presentation.display_name"), 'known presentation must render display_name, not the canonical key');
  assert(html.includes('data-presentation-state="unknown"'), 'unknown presentation state missing');
  assert(html.includes('data-presentation-state="unavailable"'), 'unavailable presentation state missing');

  // List population honesty: known-empty, unknown and unavailable stay distinct.
  assert(html.includes('data-listing-population="known" data-listing-count="0"'), 'known-empty listing population marker missing');
  assert(html.includes('data-listing-population="unknown"'), 'unknown listing population state missing');
  assert(html.includes('data-listing-population="unavailable"'), 'unavailable listing population state missing');
  assert(html.includes('Isso não significa que não existam anúncios.'), 'unavailable must not collapse into known-empty');
  for (const scenario of ['normal', 'known-empty', 'unknown', 'unavailable']) {
    assert(html.includes(`value="${scenario}"`), `deterministic list scenario missing: ${scenario}`);
  }

  // Lifecycle is observed fact, including unknown; no synthetic status/score.
  for (const lifecycle of ['active', 'paused', 'closed', 'unknown']) {
    assert(html.includes(`${lifecycle}:`) || html.includes(`'${lifecycle}'`), `observed lifecycle missing: ${lifecycle}`);
  }
  assert(!/data-readiness-score|data-listing-status="satisfied"/.test(html), 'no synthetic readiness/status projection is admitted');

  // Detail composes owner-separated read-only regions; no mutation home in R20/R21.
  assert(html.includes('data-mutation-home="none"'), 'R20/R21 must declare no mutation home');
  for (const owner of ['Offering', 'Availability', 'Performance', 'Market', 'Economics', 'Work']) {
    assert(html.includes(`data-region-owner="${owner}"`), `owner-separated detail region missing: ${owner}`);
  }
  assert(html.includes('data-write-controls="none"'), 'Offering observation region must carry no write controls');
  assert(html.includes('data-continuation-kind="navigation-only"'), 'ListingIntent continuation must be navigation-only');
  assert(html.includes('nenhum rascunho foi criado por esta navegação'), 'boundary must state that navigation creates nothing');

  // Source-product link: typed states, SKU shown only from resolved presentation, never a write carrier.
  assert(html.includes('data-source-link-read="source_product_link"'), 'source_product_link read binding missing');
  for (const state of ['resolved', 'unresolved', 'unknown', 'unavailable']) {
    assert(html.includes(`data-source-link-state="${state}"`), `source link state missing: ${state}`);
  }
  assert(html.includes('Seu produto'), 'human source-product column missing');
  assert(html.includes('SKU '), 'resolved source link must show the SKU');
  assert(html.includes('Defina o vínculo na Preparação.'), 'unresolved link must route to Preparação, not invent a link');
  assert(!html.includes('data-source-link-write'), 'source link must never become write authority');

  // Peek panel: read-only summary, never a second mutation surface.
  assert(html.includes('data-peek-kind="read-only-summary"'), 'peek panel must be a read-only summary');
  assert(html.includes('data-peek-writes="none"'), 'peek panel must declare no writes');
  assert(html.includes('Espiada somente-leitura'), 'peek read-only law copy missing');

  // Collection grammar and boundaries.
  assert(html.includes('data-collection-grammar="cursor"'), 'cursor collection grammar note missing');
  assert(html.includes('sem seleção em massa'), 'bulk-selection rejection missing');
  for (const operationId of ['ListMarketplaceListings', 'GetMarketplaceListing', 'GetSellableAvailability', 'GetMarketplaceListingPerformance', 'GetCompetitivePosition', 'GetExpectedEconomics']) {
    assert(html.includes(`data-operation="${operationId}"`), `B20 operation trace missing: ${operationId}`);
  }
  for (const forbidden of ['CreateListingIntentDraft', 'UpdateListingIntentDraft', 'SubmitListingIntent']) {
    assert(!html.includes(`data-operation="${forbidden}"`), `B20 must not carry a mutation operation: ${forbidden}`);
  }
}

function expectFailure(label, mutate) {
  const candidate = mutate(readFileSync(artifact, 'utf8'));
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `negative control did not fire: ${label}`);
}

assert(existsSync(artifact), 'B20 rendered artifact missing');
const html = readFileSync(artifact, 'utf8');
validate(html);

const controls = [
  ['presentation widened to write authority', (value) => value.replace('data-presentation-write="never"', 'data-presentation-write="resolve"')],
  ['canonical key used as human projection', (value) => value.split("presentation.state==='known'?l.presentation.display_name").join("presentation.state==='known'?l.key")],
  ['known-empty marker removed', (value) => value.replace('data-listing-population="known" data-listing-count="0"', 'data-listing-population="known" data-listing-count="1"')],
  ['unknown population collapsed', (value) => value.replace('data-listing-population="unknown"', 'data-listing-population="known"')],
  ['unavailable population collapsed', (value) => value.replace('data-listing-population="unavailable"', 'data-listing-population="unknown"')],
  ['mutation home introduced', (value) => value.replace('data-mutation-home="none"', 'data-mutation-home="inline"')],
  ['owner separation dropped', (value) => value.replace('data-region-owner="Availability"', 'data-region-owner="Offering"')],
  ['continuation became a draft-creating call', (value) => value.replace('data-continuation-kind="navigation-only"', 'data-continuation-kind="create-draft"') .replace('<span data-operation="GetExpectedEconomics">', '<span data-operation="CreateListingIntentDraft"></span><span data-operation="GetExpectedEconomics">')],
  ['unresolved source link collapsed into resolved', (value) => value.split('data-source-link-state="unresolved"').join('data-source-link-state="resolved"')],
  ['peek widened into a write surface', (value) => value.split('data-peek-writes="none"').join('data-peek-writes="inline"')],
];
for (const [label, mutate] of controls) expectFailure(label, mutate);

console.log('d6_r_b20_publications_status=CANDIDATE');
console.log('d6_r_b20_publications_scope=R20_R21_READ_ONLY');
console.log('d6_r_b20_publications_projection=TYPED_PRESENTATION');
console.log(`d6_r_b20_publications_negative_controls=${controls.length}/${controls.length}`);
console.log('d6_r_b20_publications_wireframe=PASS');
