import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const artifact = resolve(process.cwd(), 'qualification/d6-r2-wireframes/b10-preparation.html');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

function validate(html) {
  assert(html.includes('data-p8-status="candidate"'), 'B10 HTML must remain a P8 candidate; LOCK is operator-only');
  assert(!html.includes('data-p8-status="locked"'), 'B10 HTML must not self-claim operator LOCK');
  assert(html.includes('data-wire-prerequisite="pr70-integrated"'), 'B10 correspondence candidate must bind the integrated PR #70 wire');
  assert(html.includes('data-correspondence-revalidation="bounded-candidate"'), 'bounded correspondence revalidation marker missing');

  for (const protectedLabel of ['Campos para o marketplace', 'Vínculo com o marketplace', 'Continuar para configurar o anúncio']) {
    assert(html.includes(protectedLabel), `protected B10 structure missing: ${protectedLabel}`);
  }

  assert(html.includes('data-candidate-read="correspondence_candidate_population"'), 'candidate population read binding missing');
  assert(html.includes('data-candidate-fields="candidate_key display_label"'), 'candidate projection must remain candidate_key + display_label');
  assert(html.includes('data-candidate-write="candidate_key-only"'), 'correspondence write must remain key-only');
  assert(!html.includes('data-candidate-write="display_label"'), 'presentation label must not become write authority');
  assert(html.includes('selectedCorrespondenceCandidate=input.value'), 'radio selection must retain only the canonical candidate key');
  assert(html.includes('data-submitted-candidate-key'), 'prototype effect must expose the submitted canonical key for proof');
  assert(!html.includes('data-submitted-display-label'), 'prototype must not submit a presentation label');

  for (const scenario of ['resolved', 'selection', 'conflicting', 'known-empty', 'unknown', 'unavailable']) {
    assert(html.includes(`value="${scenario}"`), `deterministic correspondence scenario missing: ${scenario}`);
  }
  assert(html.includes('data-population-state="known"'), 'known candidate population state missing');
  assert(html.includes('data-population-count="0"'), 'known-empty candidate population state missing');
  assert(html.includes('data-population-state="unknown"'), 'unknown candidate population state missing');
  assert(html.includes('data-population-state="unavailable"'), 'unavailable candidate population state missing');
  assert(html.includes('Isso não significa que não existam candidatos.'), 'unavailable must not collapse into known-empty');

  assert(html.includes('name="correspondenceCandidate" value="\'+esc(c.key)+\'"'), 'candidate key must be the radio value');
  assert(html.includes('<strong>\'+esc(c.label)+\'</strong>'), 'candidate label must be the human recognition projection');
  assert(!/<input[^>]*name="correspondenceCandidate"[^>]*\schecked(?:\s|=|>)/i.test(html), 'candidate selection must not default silently');
  assert(/id="resolveCorr"[^>]*disabled/.test(html), 'resolve action must remain disabled before explicit selection');
  assert(html.includes("progression.dataset.progression='blocked-pending-reread'"), 'consequential correspondence effect must block until reread');
  assert(html.includes("function rereadReadiness(){detailScenario.value='resolved';renderDetail();}"), 'authoritative reread simulation missing');

  for (const operationId of ['SearchSourceProductsForMarketplace', 'GetProductChannelReadiness', 'GetPublicationRequirements', 'ResolveProductChannelCorrespondence', 'ClearProductChannelCorrespondence']) {
    assert(html.includes(`data-operation="${operationId}"`), `B10 operation trace missing: ${operationId}`);
  }
  assert(!html.includes('data-operation="CreateListingIntentDraft"'), 'B10 must not create ListingIntent as navigation side effect');
}

function expectFailure(label, mutate) {
  const candidate = mutate(readFileSync(artifact, 'utf8'));
  let failed = false;
  try { validate(candidate); } catch { failed = true; }
  assert(failed, `negative control did not fire: ${label}`);
}

assert(existsSync(artifact), 'B10 rendered artifact missing');
const html = readFileSync(artifact, 'utf8');
validate(html);

const controls = [
  ['write widened to display label', (value) => value.replace('data-candidate-write="candidate_key-only"', 'data-candidate-write="display_label"')],
  ['candidate display label removed', (value) => value.replace('data-candidate-fields="candidate_key display_label"', 'data-candidate-fields="candidate_key"')],
  ['known-empty marker removed', (value) => value.replace('data-population-count="0"', 'data-population-count="1"')],
  ['unknown population collapsed', (value) => value.replace('data-population-state="unknown"', 'data-population-state="known"')],
  ['unavailable population collapsed', (value) => value.replace('data-population-state="unavailable"', 'data-population-state="unknown"')],
  ['candidate silently preselected', (value) => value.replace('name="correspondenceCandidate" value="\'+esc(c.key)+\'"', 'name="correspondenceCandidate" checked value="\'+esc(c.key)+\'"')],
  ['resolve enabled before selection', (value) => value.replace('id="resolveCorr" class="btn" type="button" disabled', 'id="resolveCorr" class="btn" type="button"')],
  ['post-effect reread block removed', (value) => value.replace("progression.dataset.progression='blocked-pending-reread'", "progression.dataset.progression='listing-intent-required'")],
];
for (const [label, mutate] of controls) expectFailure(label, mutate);

console.log('d6_r_b10_correspondence_status=CANDIDATE');
console.log('d6_r_b10_correspondence_scope=BOUNDED_REGION');
console.log('d6_r_b10_correspondence_candidate_identity=KEY_ONLY_WRITE');
console.log(`d6_r_b10_correspondence_negative_controls=${controls.length}/${controls.length}`);
console.log('d6_r_b10_correspondence_wireframe=PASS');
