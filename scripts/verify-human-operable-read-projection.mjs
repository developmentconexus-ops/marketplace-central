import { readFileSync } from 'node:fs';

const bundlePath = process.argv[2];
if (!bundlePath) throw new Error('usage: node scripts/verify-human-operable-read-projection.mjs <resolved-bundle.json>');
const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function schemas(doc) { return doc.components?.schemas ?? {}; }
function refName(node) {
  const ref = node?.$ref;
  return typeof ref === 'string' ? ref.split('/').at(-1) : null;
}
function requireFieldsFrom(s, name, fields) {
  assert(s[name], `missing schema ${name}`);
  const required = new Set(s[name].required ?? []);
  for (const field of fields) assert(required.has(field), `${name} must require ${field}`);
}
function requirePropertyRef(s, name, property, expected) {
  assert(refName(s[name]?.properties?.[property]) === expected, `${name}.${property} must reference ${expected}`);
}
function requireClosedDiscriminant(s, name, property, value) {
  const schema = s[name];
  assert(schema?.type === 'object' && schema.additionalProperties === false, `${name} must be a closed object`);
  requireFieldsFrom(s, name, [property]);
  assert(schema.properties?.[property]?.const === value, `${name}.${property} must be ${value}`);
}
function requireUnionRefs(s, name, expected) {
  const actual = (s[name]?.oneOf ?? []).map(refName).sort();
  const wanted = [...expected].sort();
  assert(JSON.stringify(actual) === JSON.stringify(wanted), `${name} variants mismatch`);
}
function requirePairedRefs(s, name, canonicalProperty, presentationProperty, pairs) {
  const branches = (s[name]?.allOf ?? []).flatMap((entry) => entry?.oneOf ?? []);
  assert(branches.length === pairs.length, `${name} must correlate ${canonicalProperty} with ${presentationProperty}`);
  for (const [canonicalRef, presentationRef] of pairs) {
    assert(branches.some((branch) => refName(branch?.properties?.[canonicalProperty]) === canonicalRef
      && refName(branch?.properties?.[presentationProperty]) === presentationRef), `${name} missing ${canonicalRef} -> ${presentationRef} correlation`);
  }
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
  requirePropertyRef(s, 'SourceProductSearchHit', 'presentation', 'SourceProductPresentationKnown');
  requireClosedDiscriminant(s, 'SourceProductPresentationKnown', 'state', 'known');
  requireClosedDiscriminant(s, 'SourceProductPresentationUnknown', 'state', 'unknown');
  requireClosedDiscriminant(s, 'SourceProductPresentationUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'SourceProductPresentation', ['SourceProductPresentationKnown', 'SourceProductPresentationUnknown', 'SourceProductPresentationUnavailable']);
  requireFieldsFrom(s, 'ProductChannelReadiness', ['subject', 'subject_presentation', 'correspondence', 'correspondence_candidate_population', 'correspondence_etag', 'readiness', 'blockers', 'evaluated_at']);
  requirePropertyRef(s, 'ProductChannelReadiness', 'subject_presentation', 'SourceProductPresentation');
  requirePropertyRef(s, 'ProductChannelReadiness', 'correspondence_candidate_population', 'CorrespondenceCandidatePopulation');
  requireClosedDiscriminant(s, 'CorrespondenceCandidatePopulationKnown', 'state', 'known');
  const correspondenceCandidates = s.CorrespondenceCandidatePopulationKnown.properties?.candidates;
  assert(correspondenceCandidates?.type === 'array' && refName(correspondenceCandidates.items) === 'ProductChannelCorrespondenceCandidate', 'known correspondence candidates must be typed candidate views');
  requireClosedDiscriminant(s, 'CorrespondenceCandidatePopulationUnknown', 'state', 'unknown');
  requireClosedDiscriminant(s, 'CorrespondenceCandidatePopulationUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'CorrespondenceCandidatePopulation', ['CorrespondenceCandidatePopulationKnown', 'CorrespondenceCandidatePopulationUnknown', 'CorrespondenceCandidatePopulationUnavailable']);
  requireFieldsFrom(s, 'PublicationRequirements', ['subject', 'subject_presentation', 'publication_context', 'requirements_revision', 'requirements', 'source_media_candidates', 'evaluated_at']);
  requirePropertyRef(s, 'PublicationRequirements', 'subject_presentation', 'SourceProductPresentation');
  requirePropertyRef(s, 'PublicationRequirements', 'publication_context', 'PublicationContextView');
  requireFieldsFrom(s, 'PublicationRequirement', ['requirement_key', 'display_name', 'requirement_class', 'applicability', 'value_spec', 'not_applicable_allowed', 'source_evidence']);
  assert(JSON.stringify(s.PublicationOptionRequirementSpec).includes('PublicationOptionDescriptor'), 'option requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationOptionListRequirementSpec).includes('PublicationOptionDescriptor'), 'option-list requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationNumberUnitRequirementSpec).includes('PublicationUnitDescriptor'), 'number-unit requirement must expose unit descriptors');
  assert(JSON.stringify(s.PublicationSourceEvidenceKnown).includes('PublicationSourceCandidateView'), 'known source evidence must expose candidate views');
  assert(JSON.stringify(s.PublicationSourceEvidenceConflicting).includes('PublicationSourceCandidateView'), 'conflicting source evidence must expose candidate views');
  assert(JSON.stringify(s.SourceMediaCandidate).includes('SourceMediaPresentation'), 'source media must use its own presentation trust type');
  requireClosedDiscriminant(s, 'SourceMediaPresentationKnown', 'state', 'known');
  requireClosedDiscriminant(s, 'SourceMediaPresentationUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'SourceMediaPresentation', ['SourceMediaPresentationKnown', 'SourceMediaPresentationUnavailable']);
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
  requireFieldsFrom(s, 'MarketplaceListingListItem', ['listing', 'presentation', 'lifecycle', 'source_product_link', 'observed_at']);
  requirePropertyRef(s, 'MarketplaceListingListItem', 'presentation', 'MarketplaceListingPresentation');
  requirePropertyRef(s, 'MarketplaceListingListItem', 'source_product_link', 'ListingSourceProductLink');
  requireFieldsFrom(s, 'MarketplaceListing', ['listing', 'presentation', 'lifecycle', 'source_product_link', 'publication_context', 'observed_fields', 'observed_media', 'observed_at', 'provenance']);
  requirePropertyRef(s, 'MarketplaceListing', 'source_product_link', 'ListingSourceProductLink');
  requireClosedDiscriminant(s, 'ListingSourceProductLinkResolved', 'state', 'resolved');
  requireClosedDiscriminant(s, 'ListingSourceProductLinkUnresolved', 'state', 'unresolved');
  requireClosedDiscriminant(s, 'ListingSourceProductLinkUnknown', 'state', 'unknown');
  requireClosedDiscriminant(s, 'ListingSourceProductLinkUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'ListingSourceProductLink', ['ListingSourceProductLinkResolved', 'ListingSourceProductLinkUnresolved', 'ListingSourceProductLinkUnknown', 'ListingSourceProductLinkUnavailable']);
  requirePropertyRef(s, 'ListingSourceProductLinkResolved', 'presentation', 'SourceProductPresentation');
  requirePropertyRef(s, 'MarketplaceListing', 'presentation', 'MarketplaceListingPresentation');
  requirePropertyRef(s, 'MarketplaceListing', 'publication_context', 'PublicationContextView');
  requireClosedDiscriminant(s, 'MarketplaceListingPresentationKnown', 'state', 'known');
  requireClosedDiscriminant(s, 'MarketplaceListingPresentationUnknown', 'state', 'unknown');
  requireClosedDiscriminant(s, 'MarketplaceListingPresentationUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'MarketplaceListingPresentation', ['MarketplaceListingPresentationKnown', 'MarketplaceListingPresentationUnknown', 'MarketplaceListingPresentationUnavailable']);
  requireUnionRefs(s, 'ListingObservedField', ['ListingObservedFieldKnown', 'ListingObservedFieldUnknown', 'ListingObservedFieldUnavailable', 'ListingObservedFieldNotApplicable']);
  assert(JSON.stringify(s.ListingObservedFieldKnown).includes('PublicationValueView'), 'known Listing field must use PublicationValueView');
  assert(JSON.stringify(s.MarketplaceListingPerformanceListItem).includes('MarketplaceListingPresentation'), 'Performance Listing item must reuse MarketplaceListingPresentation');
  requirePropertyRef(s, 'MarketplaceListingPerformanceListItem', 'presentation', 'MarketplaceListingPresentation');
  assert(!Object.hasOwn(s.MarketplaceListingPerformanceListItem?.properties ?? {}, 'display_name'), 'Performance Listing item must not keep parallel display_name');
}

function validateConsumers(doc) {
  const s = schemas(doc);
  requireFieldsFrom(s, 'ListingIntentListItem', ['source_product_presentation']);
  requirePropertyRef(s, 'ListingIntentListItem', 'source_product_presentation', 'SourceProductPresentation');
  requireFieldsFrom(s, 'PriceIntent', ['target_presentation']);
  requirePropertyRef(s, 'PriceIntent', 'target_presentation', 'PriceIntentTargetPresentation');
  requireFieldsFrom(s, 'PriceIntentListItem', ['target_presentation']);
  requirePropertyRef(s, 'PriceIntentListItem', 'target_presentation', 'PriceIntentTargetPresentation');
  requireFieldsFrom(s, 'SellableAvailability', ['target_presentation']);
  requirePropertyRef(s, 'SellableAvailability', 'target_presentation', 'AvailabilityTargetPresentation');
  requireFieldsFrom(s, 'CompetitivePosition', ['subject_presentation']);
  requirePropertyRef(s, 'CompetitivePosition', 'subject_presentation', 'MarketSubjectPresentation');
  requireFieldsFrom(s, 'CompetitivePositionListItem', ['subject_presentation']);
  requirePropertyRef(s, 'CompetitivePositionListItem', 'subject_presentation', 'MarketSubjectPresentation');
  requireFieldsFrom(s, 'ExpectedEconomics', ['subject_presentation']);
  requirePropertyRef(s, 'ExpectedEconomics', 'subject_presentation', 'EconomicsSubjectPresentation');
  requireFieldsFrom(s, 'ExpectedEconomicsListItem', ['subject_presentation']);
  requirePropertyRef(s, 'ExpectedEconomicsListItem', 'subject_presentation', 'EconomicsSubjectPresentation');
  requireFieldsFrom(s, 'PriceScenarioEvaluation', ['subject_presentation']);
  requirePropertyRef(s, 'PriceScenarioEvaluation', 'subject_presentation', 'EconomicsSubjectPresentation');
  for (const name of ['PriceIntentTargetPresentation', 'AvailabilityTargetPresentation', 'MarketSubjectPresentation', 'EconomicsSubjectPresentation']) assert(s[name], `missing schema ${name}`);
  requireUnionRefs(s, 'PriceIntentTargetPresentation', ['PriceIntentExistingTargetPresentation', 'PriceIntentPreCreationTargetPresentation']);
  requireUnionRefs(s, 'AvailabilityTargetPresentation', ['AvailabilityExistingTargetPresentation', 'AvailabilityPreCreationTargetPresentation']);
  requireUnionRefs(s, 'MarketSubjectPresentation', ['MarketExistingListingPresentation', 'MarketSourceProductPresentation']);
  requireUnionRefs(s, 'EconomicsSubjectPresentation', ['EconomicsExistingListingPresentation', 'EconomicsSourceProductPresentation']);
  const pricePairs = [
    ['PriceIntentExistingTarget', 'PriceIntentExistingTargetPresentation'],
    ['PriceIntentPreCreationTarget', 'PriceIntentPreCreationTargetPresentation'],
  ];
  for (const name of ['PriceIntent', 'PriceIntentListItem']) requirePairedRefs(s, name, 'target', 'target_presentation', pricePairs);
  requirePairedRefs(s, 'SellableAvailability', 'target', 'target_presentation', [
    ['AvailabilityExistingTarget', 'AvailabilityExistingTargetPresentation'],
    ['AvailabilityPreCreationTarget', 'AvailabilityPreCreationTargetPresentation'],
  ]);
  const marketPairs = [
    ['MarketExistingListingSubject', 'MarketExistingListingPresentation'],
    ['MarketSourceProductSubject', 'MarketSourceProductPresentation'],
  ];
  for (const name of ['CompetitivePosition', 'CompetitivePositionListItem']) requirePairedRefs(s, name, 'subject', 'subject_presentation', marketPairs);
  const economicsPairs = [
    ['EconomicsExistingListingSubject', 'EconomicsExistingListingPresentation'],
    ['EconomicsSourceProductSubject', 'EconomicsSourceProductPresentation'],
  ];
  for (const name of ['ExpectedEconomics', 'ExpectedEconomicsListItem', 'PriceScenarioEvaluation']) requirePairedRefs(s, name, 'subject', 'subject_presentation', economicsPairs);
}

function validateListingIntent(doc) {
  const s = schemas(doc);
  for (const name of [
    'ListingIntentResolvedValueKnown', 'ListingIntentResolvedValueMissing', 'ListingIntentResolvedValueUnknown', 'ListingIntentResolvedValueUnavailable', 'ListingIntentResolvedValueUnsupported', 'ListingIntentResolvedValue',
    'ListingIntentFollowSourceResolutionView', 'ListingIntentExplicitOverrideResolutionView', 'ListingIntentRequirementResolutionView',
    'ListingIntentMediaPresentationKnown', 'ListingIntentMediaPresentationUnavailable', 'ListingIntentMediaPresentation', 'ListingIntentMediaPresentationDescriptor',
    'ListingIntentAttemptRequirementResolution', 'ListingIntentAttemptMediaBasis', 'ListingIntentAttemptAvailabilityInput', 'ListingIntentEffectAttempt',
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'ListingIntent', ['source_product_presentation', 'resolved_requirements', 'authored_media_presentations', 'dispatch_blockers', 'created_by_principal_id', 'updated_by_principal_id', 'effect_history']);
  requirePropertyRef(s, 'ListingIntent', 'source_product_presentation', 'SourceProductPresentation');
  assert(JSON.stringify(s.ListingIntentDesired).includes('PublicationContextRef'), 'ListingIntent desired must carry key-based publication context');
  requirePropertyRef(s, 'ListingIntentMediaPresentationDescriptor', 'presentation', 'ListingIntentMediaPresentation');
  requireClosedDiscriminant(s, 'ListingIntentMediaPresentationKnown', 'state', 'known');
  requireClosedDiscriminant(s, 'ListingIntentMediaPresentationUnavailable', 'state', 'unavailable');
  requireUnionRefs(s, 'ListingIntentMediaPresentation', ['ListingIntentMediaPresentationKnown', 'ListingIntentMediaPresentationUnavailable']);
  assert(!JSON.stringify(s.ListingIntentMediaPresentationDescriptor).includes('SourceMediaPresentation'), 'authored media must not reuse source-media presentation trust type');
  for (const name of ['ListingIntentDesired', 'RequirementResolution', 'ExplicitOverrideResolution', 'PublicationValue', 'MediaSelection', 'CreateListingIntentMediaMultipart']) {
    const schemaText = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name', 'display_label', 'access_ref', 'authored_by_principal_id', 'subject_presentation']) assert(!schemaText.includes(`"${forbidden}"`), `${name} must remain presentation-free`);
  }
}

function validateVariations(doc) {
  const s = schemas(doc);
  requireFieldsFrom(s, 'PublicationRequirement', ['scope']);
  requireFieldsFrom(s, 'PublicationRequirements', ['variation_axes']);
  requireClosedDiscriminant(s, 'VariationAxisSpecOptionKind', 'kind', 'option');
  requireClosedDiscriminant(s, 'VariationAxisSpecTextKind', 'kind', 'text');
  requireUnionRefs(s, 'VariationAxisSpec', ['VariationAxisSpecOptionKind', 'VariationAxisSpecTextKind']);
  requireClosedDiscriminant(s, 'VariationCoordinateOption', 'kind', 'option');
  requireClosedDiscriminant(s, 'VariationCoordinateText', 'kind', 'text');
  requireUnionRefs(s, 'VariationCoordinate', ['VariationCoordinateOption', 'VariationCoordinateText']);
  requireFieldsFrom(s, 'ListingIntentVariations', ['axes', 'options']);
  requireFieldsFrom(s, 'ListingIntentVariationOption', ['option_coordinates', 'requirement_resolutions', 'media_selection']);
  assert(JSON.stringify(s.ListingIntentVariationOption.properties.option_coordinates).includes('VariationCoordinate'), 'variation option identity must be coordinate-typed');
  for (const name of ['ListingIntentVariations', 'ListingIntentVariationOption']) {
    const schemaText = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['price', 'quantity', 'availability', 'display_name', 'display_label']) {
      assert(!schemaText.includes(`"${forbidden}"`), `${name} must not carry ${forbidden}; price/quantity stay with their owners and labels never enter writes`);
    }
  }
  requireFieldsFrom(s, 'MarketplaceListingObservedVariation', ['option_coordinates', 'presentation']);
  requirePropertyRef(s, 'MarketplaceListingObservedVariation', 'presentation', 'MarketplaceListingPresentation');
}

function validateAll(doc) {
  validateReadiness(doc);
  if (typeof validateMarketplaceListing === 'function') validateMarketplaceListing(doc);
  if (typeof validateConsumers === 'function') validateConsumers(doc);
  if (typeof validateListingIntent === 'function') validateListingIntent(doc);
  validateVariations(doc);
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
expectMutationFailure('Availability target presentation removed', (d) => { d.components.schemas.SellableAvailability.required = d.components.schemas.SellableAvailability.required.filter((x) => x !== 'target_presentation'); });
expectMutationFailure('override write accepts display label', (d) => { d.components.schemas.ExplicitOverrideResolution.properties.display_label = { type: 'string' }; });
expectMutationFailure('ListingIntent effect history removed', (d) => { d.components.schemas.ListingIntent.required = d.components.schemas.ListingIntent.required.filter((x) => x !== 'effect_history'); });
expectMutationFailure('authored media reuses source trust type', (d) => { d.components.schemas.ListingIntentMediaPresentationDescriptor.properties.presentation = { $ref: '#/components/schemas/SourceMediaPresentation' }; });
expectMutationFailure('correspondence candidate items weakened', (d) => { d.components.schemas.CorrespondenceCandidatePopulationKnown.properties.candidates.items = { type: 'string' }; });
expectMutationFailure('PriceIntent presentation ref weakened', (d) => { d.components.schemas.PriceIntent.properties.target_presentation = { type: 'string' }; });
expectMutationFailure('PriceIntent target correlation removed', (d) => { d.components.schemas.PriceIntent.allOf = []; });
expectMutationFailure('requirement scope removed', (d) => { d.components.schemas.PublicationRequirement.required = d.components.schemas.PublicationRequirement.required.filter((x) => x !== 'scope'); });
expectMutationFailure('variation option gains price', (d) => { d.components.schemas.ListingIntentVariationOption.properties.price = { type: 'string' }; });
expectMutationFailure('variation coordinates weakened to strings', (d) => { d.components.schemas.ListingIntentVariationOption.properties.option_coordinates.items = { type: 'string' }; });
expectMutationFailure('observed variation presentation weakened', (d) => { d.components.schemas.MarketplaceListingObservedVariation.properties.presentation = { type: 'string' }; });
assert(negativeControls === 16, `negative-control count must be 16, found ${negativeControls}`);
console.log('human_operable_read_projection=PASS');
console.log(`human_operable_read_projection_negative_controls=${negativeControls}/16`);
