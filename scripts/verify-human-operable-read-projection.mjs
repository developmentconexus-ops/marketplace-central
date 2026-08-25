import { readFileSync } from 'node:fs';

const bundlePath = process.argv[2];
if (!bundlePath) throw new Error('usage: node scripts/verify-human-operable-read-projection.mjs <resolved-bundle.json>');
const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function schemas(doc) { return doc.components?.schemas ?? {}; }
function requireFieldsFrom(s, name, fields) {
  assert(s[name], `missing schema ${name}`);
  const required = new Set(s[name].required ?? []);
  for (const field of fields) assert(required.has(field), `${name} must require ${field}`);
}

function validateReadiness(doc) {
  const s = schemas(doc);
  for (const name of [
    'SourceProductPresentationKnown', 'SourceProductPresentationUnknown', 'SourceProductPresentationUnavailable', 'SourceProductPresentation',
    'PublicationCategoryDescriptor', 'PublicationProductTypeDescriptor', 'PublicationContextRef', 'PublicationContextView',
    'PublicationOptionDescriptor', 'PublicationUnitDescriptor', 'PublicationValueView', 'PublicationSourceCandidateView',
    'ProductChannelCorrespondenceCandidate', 'CorrespondenceCandidatePopulationKnown', 'CorrespondenceCandidatePopulationUnknown', 'CorrespondenceCandidatePopulationUnavailable', 'CorrespondenceCandidatePopulation',
    'SourceMediaPresentationKnown', 'SourceMediaPresentationUnavailable', 'SourceMediaPresentation',
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'SourceProductSearchHit', ['source_product', 'presentation']);
  requireFieldsFrom(s, 'ProductChannelReadiness', ['subject', 'subject_presentation', 'correspondence', 'correspondence_candidate_population', 'correspondence_etag', 'readiness', 'blockers', 'evaluated_at']);
  requireFieldsFrom(s, 'PublicationRequirements', ['subject', 'subject_presentation', 'publication_context', 'requirements_revision', 'requirements', 'source_media_candidates', 'evaluated_at']);
  requireFieldsFrom(s, 'PublicationRequirement', ['requirement_key', 'display_name', 'requirement_class', 'applicability', 'value_spec', 'not_applicable_allowed', 'source_evidence']);
  assert(JSON.stringify(s.PublicationOptionRequirementSpec).includes('PublicationOptionDescriptor'), 'option requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationOptionListRequirementSpec).includes('PublicationOptionDescriptor'), 'option-list requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationNumberUnitRequirementSpec).includes('PublicationUnitDescriptor'), 'number-unit requirement must expose unit descriptors');
  assert(JSON.stringify(s.PublicationSourceEvidenceKnown).includes('PublicationSourceCandidateView'), 'known source evidence must expose candidate views');
  assert(JSON.stringify(s.PublicationSourceEvidenceConflicting).includes('PublicationSourceCandidateView'), 'conflicting source evidence must expose candidate views');
  assert(JSON.stringify(s.SourceMediaCandidate).includes('SourceMediaPresentation'), 'source media must use its own presentation trust type');
  for (const name of ['ResolveCorrespondenceRequest', 'ClearCorrespondenceRequest', 'ListingIntentDesired', 'RequirementResolution', 'ExplicitOverrideResolution', 'PublicationValue', 'MediaSelection', 'CreatePriceIntentRequest']) {
    const schemaText = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name', 'display_label', 'access_ref', 'subject_presentation']) assert(!schemaText.includes(`"${forbidden}"`), `${name} must not author ${forbidden}`);
  }
}

function validateMarketplaceListing(doc) {
  const s = schemas(doc);
  for (const name of [
    'MarketplaceListingPresentationKnown', 'MarketplaceListingPresentationUnknown', 'MarketplaceListingPresentationUnavailable', 'MarketplaceListingPresentation',
    'ListingObservedFieldKnown', 'ListingObservedFieldUnknown', 'ListingObservedFieldUnavailable', 'ListingObservedFieldNotApplicable', 'ListingObservedField',
    'MarketplaceListingMediaPresentationKnown', 'MarketplaceListingMediaPresentationUnavailable', 'MarketplaceListingMediaPresentation', 'MarketplaceListingObservedMedia', 'MarketplaceListingObservationProvenance',
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'MarketplaceListingListItem', ['listing', 'presentation', 'lifecycle', 'observed_at']);
  requireFieldsFrom(s, 'MarketplaceListing', ['listing', 'presentation', 'lifecycle', 'publication_context', 'observed_fields', 'observed_media', 'observed_at', 'provenance']);
  assert(JSON.stringify(s.ListingObservedFieldKnown).includes('PublicationValueView'), 'known Listing field must use PublicationValueView');
  assert(JSON.stringify(s.MarketplaceListingPerformanceListItem).includes('MarketplaceListingPresentation'), 'Performance Listing item must reuse MarketplaceListingPresentation');
  assert(!Object.hasOwn(s.MarketplaceListingPerformanceListItem?.properties ?? {}, 'display_name'), 'Performance Listing item must not keep parallel display_name');
}

function validateAll(doc) {
  validateReadiness(doc);
  if (typeof validateMarketplaceListing === 'function') validateMarketplaceListing(doc);
  if (typeof validateConsumers === 'function') validateConsumers(doc);
  if (typeof validateListingIntent === 'function') validateListingIntent(doc);
}

function expectMutationFailure(label, mutate) {
  const candidate = structuredClone(document);
  let failed = false;
  try { mutate(candidate); validateAll(candidate); } catch { failed = true; }
  assert(failed, `negative control unexpectedly passed: ${label}`);
  negativeControls++;
}

validateAll(document);
expectMutationFailure('requirement label removed', (d) => { d.components.schemas.PublicationRequirement.required = d.components.schemas.PublicationRequirement.required.filter((x) => x !== 'display_name'); });
expectMutationFailure('correspondence candidate population removed', (d) => { d.components.schemas.ProductChannelReadiness.required = d.components.schemas.ProductChannelReadiness.required.filter((x) => x !== 'correspondence_candidate_population'); });
expectMutationFailure('resolve write accepts label', (d) => { d.components.schemas.ResolveCorrespondenceRequest.properties.display_label = { type: 'string' }; });
expectMutationFailure('Listing collection presentation removed', (d) => { d.components.schemas.MarketplaceListingListItem.required = d.components.schemas.MarketplaceListingListItem.required.filter((x) => x !== 'presentation'); });
expectMutationFailure('Performance parallel display name restored', (d) => { d.components.schemas.MarketplaceListingPerformanceListItem.properties.display_name = { type: 'string' }; });
assert(negativeControls === 5, `negative-control count must be 5, found ${negativeControls}`);
console.log('human_operable_read_projection=PASS');
console.log(`human_operable_read_projection_negative_controls=${negativeControls}/5`);
