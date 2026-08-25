function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function clone(value) { return structuredClone(value); }
function refName(node) {
  const ref = node?.$ref;
  return typeof ref === 'string' ? ref.split('/').at(-1) : null;
}
function sameSet(actual, expected, label) {
  const a = [...new Set(actual)].sort();
  const e = [...new Set(expected)].sort();
  assert(JSON.stringify(a) === JSON.stringify(e), `${label} mismatch\nactual=${JSON.stringify(a)}\nexpected=${JSON.stringify(e)}`);
}
function requiredFields(schema, expected, label) {
  sameSet(schema?.required ?? [], expected, `${label} required fields`);
}
function schemaMap(document) {
  const schemas = document?.components?.schemas;
  assert(schemas && typeof schemas === 'object', 'bundled Product OAD components.schemas missing');
  return schemas;
}
function arrayItemRef(schema, label) {
  assert(schema?.type === 'array', `${label} must be an array`);
  return refName(schema.items);
}
function typedCandidateValueRef(schema, label) {
  assert(schema?.type === 'array', `${label} must be an array`);
  assert(schema.items?.type === 'object', `${label} items must constrain a candidate object`);
  return refName(schema.items.properties?.value);
}

const VALUE_FAMILIES = [
  ['text', 'PublicationTextRequirementSpec', 'PublicationTextValue', 'PublicationTextSourceEvidence'],
  ['exact_decimal', 'PublicationExactDecimalRequirementSpec', 'PublicationExactDecimalValue', 'PublicationExactDecimalSourceEvidence'],
  ['boolean', 'PublicationBooleanRequirementSpec', 'PublicationBooleanValue', 'PublicationBooleanSourceEvidence'],
  ['option', 'PublicationOptionRequirementSpec', 'PublicationOptionValueView', 'PublicationOptionSourceEvidence'],
  ['text_list', 'PublicationTextListRequirementSpec', 'PublicationTextListValue', 'PublicationTextListSourceEvidence'],
  ['option_list', 'PublicationOptionListRequirementSpec', 'PublicationOptionListValueView', 'PublicationOptionListSourceEvidence'],
  ['number_unit', 'PublicationNumberUnitRequirementSpec', 'PublicationNumberUnitValueView', 'PublicationNumberUnitSourceEvidence'],
];

function validateDescriptor(schemas, name, key) {
  const schema = schemas[name];
  assert(schema?.type === 'object' && schema.additionalProperties === false, `${name} must be a closed object`);
  requiredFields(schema, [key, 'display_name'], name);
  assert(refName(schema.properties?.[key]) === 'OpaqueKey', `${name}.${key} must remain opaque`);
  assert(schema.properties?.display_name?.type === 'string' && schema.properties.display_name.minLength === 1, `${name}.display_name must be non-empty presentation`);
}

function validatePublicationRequirementsOad(document) {
  const schemas = schemaMap(document);

  const contextRef = schemas.PublicationContextRef;
  assert(contextRef?.type === 'object' && contextRef.additionalProperties === false, 'PublicationContextRef must be a closed object');
  assert(contextRef.properties?.category_key, 'PublicationContextRef.category_key missing');
  assert(contextRef.properties?.product_type_key, 'PublicationContextRef.product_type_key missing');
  assert(refName(contextRef.properties.category_key) === 'OpaqueKey', 'publication category_key must remain opaque');
  assert(refName(contextRef.properties.product_type_key) === 'OpaqueKey', 'publication product_type_key must remain opaque');

  validateDescriptor(schemas, 'PublicationCategoryDescriptor', 'category_key');
  validateDescriptor(schemas, 'PublicationProductTypeDescriptor', 'product_type_key');
  validateDescriptor(schemas, 'PublicationOptionDescriptor', 'option_key');
  validateDescriptor(schemas, 'PublicationUnitDescriptor', 'unit_key');

  const contextView = schemas.PublicationContextView;
  assert(contextView?.type === 'object' && contextView.additionalProperties === false, 'PublicationContextView must be a closed object');
  assert(refName(contextView.properties?.category) === 'PublicationCategoryDescriptor', 'PublicationContextView.category descriptor missing');
  assert(refName(contextView.properties?.product_type) === 'PublicationProductTypeDescriptor', 'PublicationContextView.product_type descriptor missing');

  sameSet((schemas.PublicationValue?.oneOf ?? []).map(refName), [
    'PublicationTextValue', 'PublicationExactDecimalValue', 'PublicationBooleanValue', 'PublicationOptionValue',
    'PublicationTextListValue', 'PublicationOptionListValue', 'PublicationNumberUnitValue', 'PublicationNotApplicableValue',
  ], 'canonical PublicationValue variants');
  sameSet((schemas.PublicationValueView?.oneOf ?? []).map(refName), [
    'PublicationTextValue', 'PublicationExactDecimalValue', 'PublicationBooleanValue', 'PublicationOptionValueView',
    'PublicationTextListValue', 'PublicationOptionListValueView', 'PublicationNumberUnitValueView', 'PublicationNotApplicableValue',
  ], 'PublicationValueView variants');

  const candidate = schemas.PublicationSourceCandidateView;
  assert(candidate?.type === 'object' && candidate.additionalProperties === false, 'PublicationSourceCandidateView must be a closed object');
  requiredFields(candidate, ['source_candidate_key', 'display_label', 'value'], 'PublicationSourceCandidateView');
  assert(refName(candidate.properties?.source_candidate_key) === 'OpaqueKey', 'source candidate identity must remain opaque');
  assert(candidate.properties?.display_label?.minLength === 1, 'source candidate display_label must be non-empty');
  assert(refName(candidate.properties?.value) === 'PublicationValueView', 'source candidate value must use PublicationValueView');

  const known = schemas.PublicationSourceEvidenceKnown;
  requiredFields(known, ['state', 'candidates'], 'PublicationSourceEvidenceKnown');
  assert(known.properties?.state?.const === 'known', 'known source evidence discriminant changed');
  assert(known.properties?.candidates?.minItems === 1, 'known source evidence requires at least one candidate');
  assert(arrayItemRef(known.properties?.candidates, 'PublicationSourceEvidenceKnown.candidates') === 'PublicationSourceCandidateView', 'known source evidence must use candidate views');

  const conflicting = schemas.PublicationSourceEvidenceConflicting;
  requiredFields(conflicting, ['state', 'candidates'], 'PublicationSourceEvidenceConflicting');
  assert(conflicting.properties?.state?.const === 'conflicting', 'conflicting source evidence discriminant changed');
  assert(conflicting.properties?.candidates?.minItems === 2, 'conflicting source evidence requires at least two candidates');
  assert(arrayItemRef(conflicting.properties?.candidates, 'PublicationSourceEvidenceConflicting.candidates') === 'PublicationSourceCandidateView', 'conflicting source evidence must use candidate views');

  const evidenceStates = {
    PublicationSourceEvidenceMissing: 'missing',
    PublicationSourceEvidenceUnknown: 'unknown',
    PublicationSourceEvidenceUnavailable: 'unavailable',
    PublicationSourceEvidenceUnsupported: 'unsupported',
  };
  for (const [name, state] of Object.entries(evidenceStates)) {
    const schema = schemas[name];
    requiredFields(schema, ['state'], name);
    assert(schema?.properties?.state?.const === state, `${name} must preserve state=${state}`);
  }

  sameSet(
    (schemas.PublicationSourceEvidence?.oneOf ?? []).map(refName),
    ['PublicationSourceEvidenceKnown', 'PublicationSourceEvidenceMissing', 'PublicationSourceEvidenceConflicting', 'PublicationSourceEvidenceUnknown', 'PublicationSourceEvidenceUnavailable', 'PublicationSourceEvidenceUnsupported'],
    'PublicationSourceEvidence variants',
  );

  sameSet((schemas.PublicationRequirementValueSpec?.oneOf ?? []).map(refName), VALUE_FAMILIES.map(([, spec]) => spec), 'PublicationRequirementValueSpec variants');
  assert(arrayItemRef(schemas.PublicationOptionRequirementSpec?.properties?.options, 'PublicationOptionRequirementSpec.options') === 'PublicationOptionDescriptor', 'option requirement must expose descriptors');
  assert(arrayItemRef(schemas.PublicationOptionListRequirementSpec?.properties?.options, 'PublicationOptionListRequirementSpec.options') === 'PublicationOptionDescriptor', 'option-list requirement must expose descriptors');
  assert(arrayItemRef(schemas.PublicationNumberUnitRequirementSpec?.properties?.units, 'PublicationNumberUnitRequirementSpec.units') === 'PublicationUnitDescriptor', 'number-unit requirement must expose descriptors');
  assert(refName(schemas.PublicationNumberUnitRequirementSpec?.properties?.default_unit_key) === 'OpaqueKey', 'default_unit_key must remain canonical');

  for (const [kind, specName, valueName, evidenceName] of VALUE_FAMILIES) {
    const spec = schemas[specName];
    assert(spec?.properties?.kind?.const === kind, `${specName} must discriminate kind=${kind}`);
    const typedEvidence = schemas[evidenceName];
    const allOf = typedEvidence?.allOf ?? [];
    assert(allOf.some((entry) => refName(entry) === 'PublicationSourceEvidence'), `${evidenceName} must preserve owner-local source knowledge union`);
    const conditional = allOf.find((entry) => entry?.if && entry?.then);
    assert(conditional, `${evidenceName} must constrain known/conflicting candidate values`);
    sameSet(conditional.if?.properties?.state?.enum ?? [], ['known', 'conflicting'], `${evidenceName} constrained states`);
    assert((conditional.if?.required ?? []).includes('state'), `${evidenceName} type constraint must require state discriminant`);
    assert(typedCandidateValueRef(conditional.then?.properties?.candidates, `${evidenceName}.candidates`) === valueName, `${evidenceName} must bind candidates to ${valueName}`);
    assert(JSON.stringify(typedEvidence).includes('PublicationNotApplicableValue') === false, `${evidenceName} must not admit not_applicable source evidence`);
  }

  const requirement = schemas.PublicationRequirement;
  requiredFields(requirement, ['requirement_key', 'display_name', 'requirement_class', 'applicability', 'value_spec', 'not_applicable_allowed', 'source_evidence'], 'PublicationRequirement');
  assert(requirement.properties?.display_name?.minLength === 1, 'PublicationRequirement.display_name must be non-empty');
  sameSet(requirement.properties?.requirement_class?.enum ?? [], ['required', 'recommended', 'optional', 'conditional'], 'PublicationRequirement requirement_class');
  sameSet(requirement.properties?.applicability?.enum ?? [], ['current', 'draft_dependent'], 'PublicationRequirement applicability');
  assert(refName(requirement.properties?.value_spec) === 'PublicationRequirementValueSpec', 'PublicationRequirement.value_spec must use bounded value spec union');
  assert(refName(requirement.properties?.source_evidence) === 'PublicationSourceEvidence', 'PublicationRequirement.source_evidence must preserve owner-local knowledge union');
  const coupling = (requirement.allOf ?? []).find((entry) => Array.isArray(entry?.oneOf))?.oneOf ?? [];
  assert(coupling.length === VALUE_FAMILIES.length, `PublicationRequirement value/source coupling must cover ${VALUE_FAMILIES.length} families`);
  for (const [, specName, , evidenceName] of VALUE_FAMILIES) {
    assert(coupling.some((branch) => refName(branch?.properties?.value_spec) === specName && refName(branch?.properties?.source_evidence) === evidenceName), `PublicationRequirement coupling missing ${specName} -> ${evidenceName}`);
  }
  assert(coupling.every((branch) => refName(branch?.properties?.value_spec) !== 'PublicationNotApplicableValue'), 'source requirement coupling must not create not_applicable source family');

  const requirements = schemas.PublicationRequirements;
  requiredFields(requirements, ['subject', 'subject_presentation', 'publication_context', 'requirements_revision', 'requirements', 'source_media_candidates', 'evaluated_at'], 'PublicationRequirements');
  assert(refName(requirements.properties?.subject_presentation) === 'SourceProductPresentation', 'PublicationRequirements subject presentation missing');
  assert(refName(requirements.properties?.publication_context) === 'PublicationContextView', 'PublicationRequirements publication context view missing');
}

function expectFailure(name, document, mutate) {
  const candidate = clone(document);
  mutate(candidate);
  let failed = false;
  try { validatePublicationRequirementsOad(candidate); } catch { failed = true; }
  assert(failed, `publication requirements negative control unexpectedly passed: ${name}`);
}

function publicationRequirementsNegativeControls(document) {
  const controls = [
    ['description decoy cannot replace context property', (candidate) => {
      const context = candidate.components.schemas.PublicationContextRef;
      delete context.properties.category_key;
      context.description = 'decoy category_key: text must not satisfy structural proof';
    }],
    ['context view cannot fall back to opaque key', (candidate) => {
      candidate.components.schemas.PublicationContextView.properties.category = { $ref: '#/components/schemas/OpaqueKey' };
    }],
    ['candidate identity weakened', (candidate) => {
      candidate.components.schemas.PublicationSourceCandidateView.properties.source_candidate_key = { type: 'string' };
    }],
    ['candidate view opened', (candidate) => {
      candidate.components.schemas.PublicationSourceCandidateView.additionalProperties = true;
    }],
    ['candidate value contract erased', (candidate) => {
      candidate.components.schemas.PublicationSourceCandidateView.properties.value = { type: 'string' };
    }],
    ['conflict distinct candidate bound weakened', (candidate) => {
      candidate.components.schemas.PublicationSourceEvidenceConflicting.properties.candidates.minItems = 1;
    }],
    ['unsupported source state collapsed', (candidate) => {
      candidate.components.schemas.PublicationSourceEvidenceUnsupported.properties.state.const = 'unknown';
    }],
    ['text source candidate type widened', (candidate) => {
      const conditional = candidate.components.schemas.PublicationTextSourceEvidence.allOf.find((entry) => entry.if && entry.then);
      conditional.then.properties.candidates.items.properties.value = { $ref: '#/components/schemas/PublicationValueView' };
    }],
    ['not-applicable leaked into text source candidates', (candidate) => {
      const conditional = candidate.components.schemas.PublicationTextSourceEvidence.allOf.find((entry) => entry.if && entry.then);
      conditional.then.properties.candidates.items.properties.value = { $ref: '#/components/schemas/PublicationNotApplicableValue' };
    }],
    ['requirement value/source coupling removed', (candidate) => {
      candidate.components.schemas.PublicationRequirement.allOf = [];
    }],
    ['requirement display name omitted', (candidate) => {
      candidate.components.schemas.PublicationRequirement.required = candidate.components.schemas.PublicationRequirement.required.filter((field) => field !== 'display_name');
    }],
    ['response publication context omitted', (candidate) => {
      candidate.components.schemas.PublicationRequirements.required = candidate.components.schemas.PublicationRequirements.required.filter((field) => field !== 'publication_context');
    }],
    ['response subject presentation omitted', (candidate) => {
      candidate.components.schemas.PublicationRequirements.required = candidate.components.schemas.PublicationRequirements.required.filter((field) => field !== 'subject_presentation');
    }],
  ];
  for (const [name, mutate] of controls) expectFailure(name, document, mutate);
  return controls.length;
}

export function provePublicationRequirementsOad(document) {
  validatePublicationRequirementsOad(document);
  const negativeControls = publicationRequirementsNegativeControls(document);
  return { negativeControls };
}
