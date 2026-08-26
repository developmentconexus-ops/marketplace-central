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

  // Owner-evaluated pricing indicators: waterfall, contribution+margin, policy judgment, delivered-price position.
  assert(html.includes('data-evaluation-source="EvaluatePriceScenario"'), 'evaluation must come from the owner operation');
  assert(html.includes('data-evaluation-trigger="debounced-owner-call"'), 'evaluation must be a debounced owner call');
  assert(html.includes('data-evaluation-debounce-ms="400"'), 'debounce window binding missing');
  assert(html.includes('data-client-computation="none"'), 'the screen must declare it computes nothing');
  assert(html.includes('setTimeout(()=>requestEvaluation(key),DEBOUNCE_MS)'), 'typing must schedule a debounced owner evaluation');
  assert(html.includes('data-waterfall="owner-components"'), 'price waterfall from owner components missing');
  for (const line of ['Tarifa do marketplace', 'Frete que você paga', 'Imposto', 'Promoção', 'Custo do produto', 'Contribuição']) {
    assert(html.includes(line), `waterfall line missing: ${line}`);
  }
  assert(html.includes('data-policy-judgment="owner"'), 'policy judgment must be owner-issued');
  assert(html.includes('data-profitability="'), 'acceptable/below_policy projection missing');
  assert(html.includes('Novo preço e margem nova'), 'new-margin column heading missing');
  assert(html.includes('preço entregue '), 'delivered-price position indicator missing');
  assert(html.includes('data-economics-conclusion='), 'honest economics conclusion states missing');
  assert(html.includes('data-market-evidence="insufficient"') && html.includes('data-market-evidence="unavailable"'), 'honest market evidence states missing');
  assert(html.includes('data-anchor-apply="fills-input-only"'), 'anchors must only fill the input');
  assert(html.includes('data-market-anchor-gate="evidence_sufficiency"') && html.includes('data-anchor-gate="evidence_sufficiency"'), 'market anchor must be gated by evidence sufficiency');
  assert(html.includes('data-anchor-missing="market"'), 'missing-market-anchor explanation required');
  assert(html.includes('Piso da política'), 'policy floor anchor missing');
  // Live indicators while typing, with the current margin preserved.
  assert(html.includes('data-live-evaluation="'), 'live evaluation region missing');
  assert(html.includes('data-live-contribution'), 'live contribution/margin indicator missing');
  assert(html.includes('vs atual '), 'the current margin must stay visible next to the new one');
  assert(html.includes('data-current-margin-preserved="true"') && html.includes('data-current-margin="known"'), 'current margin column must be preserved');
  assert(html.includes("' p.p.'") || html.includes('p.p.'), 'percentage-point delta missing');
  // Market range restored in the column, plus the positional bar on the typed price.
  assert(html.includes('data-market-range="low-high"'), 'market range must be shown in the Mercado column');
  assert(html.includes('Faixa entregue: '), 'delivered-price range copy missing');
  assert(html.includes('class="rangebar"'), 'positional range bar missing');
  // Waterfall behind an expand/collapse disclosure.
  assert(html.includes('data-waterfall-disclosure="collapsible"'), 'waterfall disclosure binding missing');
  assert(/class="disclosure hidden" type="button" data-disclosure-for="[^"]*" aria-expanded="false"/.test(html), 'waterfall disclosure must start collapsed with aria-expanded');
  assert(html.includes('▸') && html.includes('▾'), 'expand/collapse affordance missing');

  // Single create home: the ledger view never creates.
  assert(html.includes('data-create-home="workbench-only"'), 'single create home binding missing');
  assert(html.includes('data-intents-view="ledger-only"'), 'intents view must be ledger-only');
  assert(html.includes('Esta visão não cria nem altera intenções.'), 'ledger-only law copy missing');
  assert(!html.includes('id="createIntent"'), 'the ledger view must not carry a create form');
  assert(html.includes("kind:'pre_creation'") && html.includes('Alvo pré-criação'), 'pre-creation rows must live in the workbench');
  assert(html.includes('data-target-kind="'), 'workbench rows must carry their target kind');

  // Pricing workbench: listing-centric grid, per-row explicit writes, fixed views only.
  assert(html.includes('data-views="decide intents"'), 'the two fixed views (decide/intents) are missing');
  assert(html.includes('data-view="decide"') && html.includes('data-view="intents"'), 'view tabs missing');
  assert(html.includes('data-row-write="one-intent-per-row"'), 'per-row explicit write law missing');
  assert(html.includes('data-workbench-composition="page-level-owner-collections"'), 'workbench must compose owner collections page-level');
  assert(html.includes('data-fact-owner="Economics"') && html.includes('data-fact-owner="Market"'), 'owner-attributed row facts missing');
  assert(html.includes('data-filter-basis="server-facts"'), 'attention filter must project server facts, not client scoring');
  assert(/class="btn small row-confirm"[^>]*disabled/.test(html), 'row confirm must stay disabled until an explicit price is typed');
  assert(html.includes('data-workbench-population="known" data-workbench-count="0"'), 'workbench known-empty marker missing');
  assert(html.includes('data-workbench-population="unknown"'), 'workbench unknown population missing');
  assert(html.includes('data-workbench-population="unavailable"'), 'workbench unavailable population missing');

  // Read-only context facts; targets typed.
  assert(html.includes('data-context-facts-kind="inline-read-only"'), 'inline read-only context facts binding missing');
  assert(html.includes('nenhum preço é calculado nesta tela, sugerido como recomendado, nem aplicado automaticamente'), 'no-auto-pricing law copy missing');
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
  ['ambiguous blind retry admitted', (value) => value.split('data-ambiguous-retry="forbidden"').join('data-ambiguous-retry="allowed"')],
  ['context facts widened into writes', (value) => value.replace('data-context-facts-kind="inline-read-only"', 'data-context-facts-kind="inline-actions"')],
  ['known-empty marker removed', (value) => value.replace('data-intent-population="known" data-intent-count="0"', 'data-intent-population="known" data-intent-count="1"')],
  ['unavailable population collapsed', (value) => value.split('data-intent-population="unavailable"').join('data-intent-population="unknown"')],
  ['create idempotency dropped', (value) => value.split('data-create-idempotency="Idempotency-Key"').join('')],
  ['provider feedback hidden', (value) => value.split('data-provider-feedback="verbatim"').join('')],
  ['listing mutation smuggled in', (value) => value.replace('<span data-operation="ListExpectedEconomics">', '<span data-operation="SubmitListingIntent"></span><span data-operation="ListExpectedEconomics">')],
  ['bulk apply introduced', (value) => value.split('data-row-write="one-intent-per-row"').join('data-row-write="bulk-apply"')],
  ['row confirm enabled without input', (value) => value.replace('<button class="btn small row-confirm" type="button" data-row-target="\'+esc(r.key)+\'" disabled>', '<button class="btn small row-confirm" type="button" data-row-target="\'+esc(r.key)+\'">')],
  ['attention filter became client scoring', (value) => value.replace('data-filter-basis="server-facts"', 'data-filter-basis="client-score"')],
  ['evaluation detached from the owner operation', (value) => value.split('data-evaluation-source="EvaluatePriceScenario"').join('data-evaluation-source="client-formula"')],
  ['evaluation fired per keystroke', (value) => value.split('data-evaluation-trigger="debounced-owner-call"').join('data-evaluation-trigger="keystroke"')],
  ['debounce removed from typing', (value) => value.replace('setTimeout(()=>requestEvaluation(key),DEBOUNCE_MS)', 'requestEvaluation(key)')],
  ['screen started computing locally', (value) => value.replace('data-client-computation="none"', 'data-client-computation="local-formula"')],
  ['current margin dropped in favour of the new one', (value) => value.split('data-current-margin-preserved="true"').join('data-current-margin-preserved="false"')],
  ['market range removed from the column', (value) => value.split('data-market-range="low-high"').join('data-market-range="none"')],
  ['waterfall forced always-open', (value) => value.split('data-waterfall-disclosure="collapsible"').join('data-waterfall-disclosure="always-open"')],
  ['policy judgment taken over by the screen', (value) => value.split('data-policy-judgment="owner"').join('data-policy-judgment="client"')],
  ['market anchor ungated from evidence', (value) => value.split('data-anchor-gate="evidence_sufficiency"').join('')],
  ['anchor started auto-applying', (value) => value.split('data-anchor-apply="fills-input-only"').join('data-anchor-apply="auto-submit"')],
  ['waterfall stripped from the evaluation', (value) => value.split('data-waterfall="owner-components"').join('data-waterfall="none"')],
  ['ledger view regained a create form', (value) => value.split('data-intents-view="ledger-only"').join('data-intents-view="create-and-ledger"')],
];
for (const [label, mutate] of controls) expectFailure(label, mutate);

console.log('d6_r_b24_price_intents_status=CANDIDATE');
console.log('d6_r_b24_price_intents_scope=R24_PRICE_AUTHORING');
console.log('d6_r_b24_price_intents_change_law=SUPERSEDE_ONLY');
console.log(`d6_r_b24_price_intents_negative_controls=${controls.length}/${controls.length}`);
console.log('d6_r_b24_price_intents_wireframe=PASS');
