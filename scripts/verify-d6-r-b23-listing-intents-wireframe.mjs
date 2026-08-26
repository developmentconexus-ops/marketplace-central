import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const artifact = resolve(process.cwd(), 'qualification/d6-r2-wireframes/b23-listing-intents.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

function validate(html) {
  assert(html.includes('data-p8-status="candidate"'), 'B23 HTML must remain a P8 candidate; LOCK is operator-only');
  assert(!html.includes('data-p8-status="locked"'), 'B23 HTML must not self-claim operator LOCK');
  assert(html.includes('data-surface="R22-R23"'), 'B23 surface binding missing');
  assert(html.includes('data-external-effect="simulated"'), 'prototype must declare all external effects simulated');

  // Write authority: Offering only, canonical keys only, idempotent creates.
  assert(html.includes('data-owner="Offering"'), 'Offering write ownership binding missing');
  assert(html.includes('data-write-permission="listing.manage"'), 'listing.manage write permission binding missing');
  assert(html.includes('data-write-principals="H-A"'), 'H/A write principal binding missing');
  assert(html.includes('data-write-identity="canonical-keys-only"'), 'writes must carry canonical keys only');
  assert(html.includes('data-create-idempotency="Idempotency-Key"'), 'create idempotency binding missing');

  // Requirement resolutions: exactly the two admitted kinds.
  assert(html.includes('data-resolution-kinds="follow_source explicit_override"'), 'resolution kinds binding missing');
  assert(html.includes('data-kind="follow_source"'), 'follow_source resolution missing');
  assert(html.includes('data-kind="explicit_override"'), 'explicit_override resolution missing');
  assert(!/data-kind="(auto|inferred|suggested)/.test(html), 'no third resolution kind may be invented');

  // Submit gate is server dispatchability; blockers are listed, never computed away.
  assert(html.includes('data-submit-gate="server-dispatchability"'), 'submit gate binding missing');
  for (const state of ['dispatchable', 'blocked', 'unknown', 'unavailable']) {
    assert(html.includes(`'${state}'`) || html.includes(`"${state}"`), `dispatchability state missing: ${state}`);
  }
  assert(/id="submitIntent"[^>]*disabled/.test(html), 'submit must start disabled until server dispatchability allows it');
  assert(html.includes("submitIntent.disabled=state!=='dispatchable'"), 'submit enablement must follow server dispatchability only');

  // External effect honesty: five states; ambiguous verifies, never blind-retries.
  for (const state of ['pending', 'accepted', 'rejected', 'ambiguous']) {
    assert(html.includes(`'${state}'`), `external effect state missing: ${state}`);
  }
  assert(html.includes('data-ambiguous-retry="forbidden"'), 'ambiguous blind-retry prohibition missing');
  assert(html.includes('data-verify-action="authoritative-reread"'), 'ambiguous outcome must offer authoritative verification');
  assert(html.includes('nunca é reenviado às cegas'), 'ambiguous law copy missing');
  assert(html.includes('data-provider-feedback="verbatim"'), 'rejection must surface provider feedback');
  assert(html.includes('data-convergence'), 'convergence projection missing');

  // Publication context, typed technical sheet and grouped census.
  assert(html.includes('data-publication-context="category product_type"'), 'publication context region missing');
  assert(html.includes('data-value-specs="text exact_decimal boolean option option_list text_list number_unit"'), 'typed value-spec binding missing');
  for (const spec of ['data-value-spec="text"', 'data-value-spec="option"', 'data-value-spec="number_unit"', 'data-value-spec="boolean"', 'data-value-spec="text_list"', 'data-value-spec="exact_decimal"']) {
    assert(html.includes(spec), `typed field rendering missing: ${spec}`);
  }
  assert(html.includes('maxlength="60"'), 'title character limit missing');
  assert(html.includes('data-requirement-class="required"'), 'required group with pendency count missing');
  assert(html.includes('Exibir todos os campos'), 'progressive show-all-fields disclosure missing');
  assert(html.includes('data-resolution="not_applicable"'), 'not-applicable resolution missing');
  assert(html.includes('data-media-role="primary"'), 'primary photo emphasis missing');

  // Variations: provider vocabulary read, coordinate-keyed writes, per-variation scope, excluded owners.
  assert(html.includes('data-variation-axes-read="variation_axes"'), 'variation axes vocabulary read binding missing');
  assert(html.includes('data-variation-write="coordinate-keys-only"'), 'variation write identity must be coordinate keys only');
  assert(html.includes('data-requirement-scope="listing"'), 'listing-scoped requirement marker missing');
  assert(html.includes('data-variation-scoped-fields="per_variation"'), 'per-variation scope binding missing');
  assert(html.includes('data-axis-kind="option"'), 'option-kind axis rendering missing');
  assert(html.includes('data-option-coordinates='), 'variation options must carry coordinate identity');
  assert(html.includes('data-variation-excluded="price quantity"'), 'price/quantity exclusion note missing');
  assert(html.includes('Cor: Inox') && html.includes('sem foto'.toLowerCase() ? 'Sem foto' : 'Sem foto'), 'per-option blocker evidence missing');
  assert(html.includes('data-operation="GetPublicationRequirements"'), 'census read trace missing');

  // Lifecycle and revision honesty.
  for (const lifecycle of ['draft', 'submitted', 'discarded']) {
    assert(html.includes(`${lifecycle}:`) || html.includes(`'${lifecycle}'`), `intent lifecycle missing: ${lifecycle}`);
  }
  assert(html.includes('data-revision-state="stale"'), 'stale requirements-revision warning missing');

  // Collection honesty and grammar.
  assert(html.includes('data-intent-population="known" data-intent-count="0"'), 'known-empty intent population marker missing');
  assert(html.includes('data-intent-population="unknown"'), 'unknown intent population state missing');
  assert(html.includes('data-intent-population="unavailable"'), 'unavailable intent population state missing');
  assert(html.includes('Isso não significa que não existam intenções.'), 'unavailable must not collapse into known-empty');
  assert(html.includes('data-collection-grammar="cursor"'), 'cursor collection grammar note missing');

  // Operation trace: full ListingIntent family, nothing outside it.
  for (const operationId of ['ListListingIntents', 'CreateListingIntentDraft', 'GetListingIntent', 'UpdateListingIntentDraft', 'DiscardListingIntentDraft', 'SubmitListingIntent', 'CreateListingIntentMedia']) {
    assert(html.includes(`data-operation="${operationId}"`), `B23 operation trace missing: ${operationId}`);
  }
  assert(!html.includes('data-operation="ResolveProductChannelCorrespondence"'), 'B23 must not carry correspondence writes');
}

function expectFailure(label, mutate) {
  const candidate = mutate(readFileSync(artifact, 'utf8'));
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `negative control did not fire: ${label}`);
}

assert(existsSync(artifact), 'B23 rendered artifact missing');
const html = readFileSync(artifact, 'utf8');
validate(html);

const controls = [
  ['submit enabled before server dispatchability', (value) => value.replace('id="submitIntent" class="btn primary" type="button" disabled', 'id="submitIntent" class="btn primary" type="button"')],
  ['submit gate rewired to client logic', (value) => value.replace("submitIntent.disabled=state!=='dispatchable'", 'submitIntent.disabled=false')],
  ['ambiguous blind retry admitted', (value) => value.split('data-ambiguous-retry="forbidden"').join('data-ambiguous-retry="allowed"')],
  ['third resolution kind invented', (value) => value.replace('data-kind="follow_source"', 'data-kind="auto"')],
  ['write identity widened past canonical keys', (value) => value.replace('data-write-identity="canonical-keys-only"', 'data-write-identity="labels"')],
  ['create idempotency dropped', (value) => value.split('data-create-idempotency="Idempotency-Key"').join('')],
  ['known-empty marker removed', (value) => value.replace('data-intent-population="known" data-intent-count="0"', 'data-intent-population="known" data-intent-count="1"')],
  ['unavailable population collapsed', (value) => value.split('data-intent-population="unavailable"').join('data-intent-population="unknown"')],
  ['stale revision warning removed', (value) => value.replace('data-revision-state="stale"', 'data-revision-state="current"')],
  ['live external effect claimed', (value) => value.replace('data-external-effect="simulated"', 'data-external-effect="live"')],
  ['variation write widened past coordinate keys', (value) => value.split('data-variation-write="coordinate-keys-only"').join('data-variation-write="labels"')],
  ['per-variation scope collapsed', (value) => value.split('data-variation-scoped-fields="per_variation"').join('data-variation-scoped-fields="listing"')],
  ['price/quantity pulled into variation authoring', (value) => value.split('data-variation-excluded="price quantity"').join('data-variation-excluded="none"')],
  ['title limit removed', (value) => value.replace('maxlength="60"', '')],
];
for (const [label, mutate] of controls) expectFailure(label, mutate);

console.log('d6_r_b23_listing_intents_status=CANDIDATE');
console.log('d6_r_b23_listing_intents_scope=R22_R23_AUTHORING');
console.log('d6_r_b23_listing_intents_write_identity=CANONICAL_KEYS_ONLY');
console.log(`d6_r_b23_listing_intents_negative_controls=${controls.length}/${controls.length}`);
console.log('d6_r_b23_listing_intents_wireframe=PASS');
