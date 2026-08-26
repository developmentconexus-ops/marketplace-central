import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const artifact = resolve(process.cwd(), 'qualification/d6-r2-wireframes/b24-price-intents.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

function validate(html) {
  assert(html.includes('data-p8-status="candidate"'), 'B24 HTML must remain a P8 candidate; LOCK is operator-only');
  assert(!html.includes('data-p8-status="locked"'), 'B24 HTML must not self-claim operator LOCK');
  assert(html.includes('data-surface="R24"'), 'B24 surface binding missing');
  assert(html.includes('data-external-effect="simulated"'), 'prototype must declare all external effects simulated');

  // Write authority and price-change law.
  assert(html.includes('data-owner="Offering"'), 'Offering write ownership binding missing');
  assert(html.includes('data-write-permission="price.manage"'), 'price.manage write permission binding missing');
  assert(html.includes('data-write-principals="H-A"'), 'H/A write principal binding missing');
  assert(html.includes('data-create-idempotency="Idempotency-Key"'), 'create idempotency binding missing');
  assert(html.includes('data-price-change="supersede-only"'), 'price change must be supersede-only');
  assert(html.includes('Preço não é editado no lugar.'), 'supersede-instead-of-edit law copy missing');
  assert(html.includes('data-repricing-engine="none"'), 'automatic repricing must be excluded');
  assert(/id="createIntent"[^>]*disabled/.test(html), 'create must start disabled until target and price are explicit');

  // Honest states: full PriceIntent state and convergence vocabularies.
  for (const state of ['pending', 'applied', 'rejected', 'ambiguous', 'superseded']) {
    assert(html.includes(`${state}:`) || html.includes(`'${state}'`), `price intent state missing: ${state}`);
  }
  for (const conv of ['pending', 'converged', 'divergent', 'unknown', 'not_applicable']) {
    assert(html.includes(`${conv}:`) || html.includes(`'${conv}'`), `convergence state missing: ${conv}`);
  }
  assert(html.includes('data-intent-population="known" data-intent-count="0"'), 'known-empty population marker missing');
  assert(html.includes('data-intent-population="unknown"'), 'unknown population state missing');
  assert(html.includes('data-intent-population="unavailable"'), 'unavailable population state missing');
  assert(html.includes('Isso não significa que não existam intenções.'), 'unavailable must not collapse into known-empty');

  // Ambiguous law and provider feedback.
  assert(html.includes('data-ambiguous-retry="forbidden"'), 'ambiguous blind-retry prohibition missing');
  assert(html.includes('data-verify-action="authoritative-reread"'), 'ambiguous outcome must offer authoritative verification');
  assert(html.includes('nunca é reenviado às cegas'), 'ambiguous law copy missing');
  assert(html.includes('data-provider-feedback="verbatim"'), 'rejection must surface provider feedback');

  // Pricing workbench: listing-centric grid, per-row explicit writes, fixed views only.
  assert(html.includes('data-views="decide intents"'), 'the two fixed views (decide/intents) are missing');
  assert(html.includes('data-view="decide"') && html.includes('data-view="intents"'), 'view tabs missing');
  assert(html.includes('data-row-write="one-intent-per-row"'), 'per-row explicit write law missing');
  assert(html.includes('data-workbench-composition="page-level-owner-collections"'), 'workbench must compose owner collections page-level');
  assert(html.includes('data-fact-owner="Economics"') && html.includes('data-fact-owner="Market"'), 'owner-attributed row facts missing');
  assert(html.includes('data-filter-basis="server-facts"'), 'attention filter must project server facts, not client scoring');
  assert(/class="btn row-confirm"[^>]*disabled/.test(html), 'row confirm must stay disabled until an explicit price is typed');
  assert(html.includes('data-workbench-population="known" data-workbench-count="0"'), 'workbench known-empty marker missing');
  assert(html.includes('data-workbench-population="unknown"'), 'workbench unknown population missing');
  assert(html.includes('data-workbench-population="unavailable"'), 'workbench unavailable population missing');

  // Read-only context facts; targets typed.
  assert(html.includes('data-context-facts-kind="inline-read-only"'), 'inline read-only context facts binding missing');
  assert(html.includes('nenhum preço é calculado ou aplicado automaticamente'), 'no-auto-pricing law copy missing');
  assert(html.includes('anúncio ainda não criado'), 'pre-creation target presentation missing');
  assert(html.includes('data-collection-grammar="cursor"'), 'cursor collection grammar note missing');
  assert(html.includes('sem reprecificação em massa'), 'bulk repricing rejection missing');

  for (const operationId of ['ListMarketplaceListings', 'ListPriceIntents', 'ListExpectedEconomics', 'ListCompetitivePositions', 'CreatePriceIntent', 'GetPriceIntent']) {
    assert(html.includes(`data-operation="${operationId}"`), `B24 operation trace missing: ${operationId}`);
  }
  assert(!html.includes('data-operation="SubmitListingIntent"'), 'B24 must not carry listing mutations');
}

function expectFailure(label, mutate) {
  const candidate = mutate(readFileSync(artifact, 'utf8'));
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `negative control did not fire: ${label}`);
}

assert(existsSync(artifact), 'B24 rendered artifact missing');
const html = readFileSync(artifact, 'utf8');
validate(html);

const controls = [
  ['price edit-in-place admitted', (value) => value.split('data-price-change="supersede-only"').join('data-price-change="edit-in-place"')],
  ['repricing engine introduced', (value) => value.replace('data-repricing-engine="none"', 'data-repricing-engine="rules"')],
  ['create enabled before explicit input', (value) => value.replace('id="createIntent" class="btn primary" type="button" disabled', 'id="createIntent" class="btn primary" type="button"')],
  ['ambiguous blind retry admitted', (value) => value.split('data-ambiguous-retry="forbidden"').join('data-ambiguous-retry="allowed"')],
  ['context facts widened into writes', (value) => value.replace('data-context-facts-kind="inline-read-only"', 'data-context-facts-kind="inline-actions"')],
  ['known-empty marker removed', (value) => value.replace('data-intent-population="known" data-intent-count="0"', 'data-intent-population="known" data-intent-count="1"')],
  ['unavailable population collapsed', (value) => value.split('data-intent-population="unavailable"').join('data-intent-population="unknown"')],
  ['create idempotency dropped', (value) => value.split('data-create-idempotency="Idempotency-Key"').join('')],
  ['provider feedback hidden', (value) => value.split('data-provider-feedback="verbatim"').join('')],
  ['listing mutation smuggled in', (value) => value.replace('<span data-operation="ListExpectedEconomics">', '<span data-operation="SubmitListingIntent"></span><span data-operation="ListExpectedEconomics">')],
  ['bulk apply introduced', (value) => value.split('data-row-write="one-intent-per-row"').join('data-row-write="bulk-apply"')],
  ['row confirm enabled without input', (value) => value.replace('class="btn row-confirm" type="button" data-row-target="\'+esc(l.key)+\'" disabled', 'class="btn row-confirm" type="button" data-row-target="\'+esc(l.key)+\'"')],
  ['attention filter became client scoring', (value) => value.replace('data-filter-basis="server-facts"', 'data-filter-basis="client-score"')],
];
for (const [label, mutate] of controls) expectFailure(label, mutate);

console.log('d6_r_b24_price_intents_status=CANDIDATE');
console.log('d6_r_b24_price_intents_scope=R24_PRICE_AUTHORING');
console.log('d6_r_b24_price_intents_change_law=SUPERSEDE_ONLY');
console.log(`d6_r_b24_price_intents_negative_controls=${controls.length}/${controls.length}`);
console.log('d6_r_b24_price_intents_wireframe=PASS');
