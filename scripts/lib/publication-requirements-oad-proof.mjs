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

const VALUE_FAMILIES = [
  ['text', 'PublicationTextRequirementSpec', 'PublicationTextValue', 'PublicationTextSourceEvidence'],
  ['exact_decimal', 'PublicationExactDecimalRequirementSpec', 'PublicationExactDecimalValue', 'PublicationExactDecimalSourceEvidence'],
  ['boolean', 'PublicationBooleanRequirementSpec', 'PublicationBooleanValue', 'PublicationBooleanSourceEvidence'],
  ['option', 'PublicationOptionRequirementSpec', 'PublicationOptionValue', 'PublicationOptionSourceEvidence'],
  ['text_list', 'PublicationTextListRequirementSpec', 'PublicationTextListValue', 'PublicationTextListSourceEvidence'],
  ['option_list', 'PublicationOptionListRequirementSpec', 'PublicationOptionListValue', 'PublicationOptionListSourceEvidence'],
  ['number_unit', 'PublicationNumberUnitRequirementSpec', 'PublicationNumberUnitValue', 'PublicationNumberUnitSourceEvidence'],
];

function validatePublicationRequirementsOad(document) {
  const schemas = schemaMap(document);

  const context = schemas.PublicationRequirementsContext;
  assert(context?.type === 'object' && context.additionalProperties === false, 'PublicationRequirementsContext must be a closed object');
  assert(context.properties?.category_key, 'PublicationRequirementsContext.category_key missing');
  assert(context.properties?.product_type_key, 'PublicationRequirementsContext.product_type_key missing');
  assert(refName(context.properties.category_key) === 'OpaqueKey', 'publication category_key must remain opaque');
  assert(refName(context.properties.product_type_key) === 'OpaqueKey', 'publication product_type_key must remain opaque');

  assert(!schemas.PublicationSourceCandidate, 'legacy array-item PublicationSourceCandidate shadow representation must be absent');
  const candidates = schemas.PublicationSourceCandidateValues;
  assert(candidates?.type === 'object', 'PublicationSourceCandidateValues must be an object keyed by candidate identity');
  assert(candidates.minProperties === 1, 'PublicationSourceCandidateValues must require at least one candidate');
  assert(refName(candidates.propertyNames) === 'OpaqueKey', 'candidate map property names must be opaque Readiness keys');
  assert(refName(candidates.additionalProperties) === 'PublicationValue', 'candidate map values must use canonical PublicationValue');

  const known = schemas.PublicationSourceEvidenceKnown;
  requiredFields(known, ['state', 'candidates'], 'PublicationSourceEvidenceKnown');
  assert(known.properties?.state?.const === 'known', 'known source evidence discriminant changed');
  assert(refName(known.properties?.candidates) === 'PublicationSourceCandidateValues', 'known source evidence must use candidate-key map');

  const conflicting = schemas.PublicationSourceEvidenceConflicting;
  requiredFields(conflicting, ['state', 'candidates'], 'PublicationSourceEvidenceConflicting');
  assert(conflicting.properties?.state?.const === 'conflicting', 'conflicting source evidence discriminant changed');
  const conflictCandidateAllOf = conflicting.properties?.candidates?.allOf ?? [];
  assert(conflictCandidateAllOf.some((entry) => refName(entry) === 'PublicationSourceCandidateValues'), 'conflicting source evidence must reuse candidate-key map');
  assert(conflictCandidateAllOf.some((entry) => entry?.minProperties === 2), 'conflicting source evidence requires at least two distinct candidate identities');

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

  const sourceEvidence = schemas.PublicationSourceEvidence;
  sameSet(
    (sourceEvidence?.oneOf ?? []).map(refName),
    ['PublicationSourceEvidenceKnown', 'PublicationSourceEvidenceMissing', 'PublicationSourceEvidenceConflicting', 'PublicationSourceEvidenceUnknown', 'PublicationSourceEvidenceUnavailable', 'PublicationSourceEvidenceUnsupported'],
    'PublicationSourceEvidence variants',
  );

  const valueSpec = schemas.PublicationRequirementValueSpec;
  sameSet((valueSpec?.oneOf ?? []).map(refName), VALUE_FAMILIES.map(([, spec]) => spec), 'PublicationRequirementValueSpec variants');

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
    const typedCandidates = conditional.then?.properties?.candidates;
    assert(typedCandidates?.type === 'object', `${evidenceName} candidate overlay must remain an object map`);
    assert(refName(typedCandidates.additionalProperties) === valueName, `${evidenceName} must bind candidates to ${valueName}`);
    assert(JSON.stringify(typedEvidence).includes('PublicationNotApplicableValue') === false, `${evidenceName} must not admit not_applicable source evidence`);
  }

  const requirement = schemas.PublicationRequirement;
  requiredFields(requirement, ['requirement_key', 'requirement_class', 'applicability', 'value_spec', 'not_applicable_allowed', 'source_evidence'], 'PublicationRequirement');
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
  requiredFields(requirements, ['subject', 'publication_context', 'requirements_revision', 'requirements', 'source_media_candidates', 'evaluated_at'], 'PublicationRequirements');
  assert(refName(requirements.properties?.publication_context) === 'PublicationRequirementsContext', 'PublicationRequirements publication_context ref missing');
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
      const context = candidate.components.schemas.PublicationRequirementsContext;
      delete context.properties.category_key;
      context.description = 'decoy category_key: text must not satisfy structural proof';
    }],
    ['candidate identity map weakened', (candidate) => {
      candidate.components.schemas.PublicationSourceCandidateValues.propertyNames = { type: 'string' };
    }],
    ['conflict distinct identity weakened', (candidate) => {
      const allOf = candidate.components.schemas.PublicationSourceEvidenceConflicting.properties.candidates.allOf;
      const bound = allOf.find((entry) => entry.minProperties === 2);
      bound.minProperties = 1;
    }],
    ['unsupported source state collapsed', (candidate) => {
      candidate.components.schemas.PublicationSourceEvidenceUnsupported.properties.state.const = 'unknown';
    }],
    ['text source candidate type widened', (candidate) => {
      const allOf = candidate.components.schemas.PublicationTextSourceEvidence.allOf;
      const conditional = allOf.find((entry) => entry.if && entry.then);
      conditional.then.properties.candidates.additionalProperties = { $ref: '#/components/schemas/PublicationValue' };
    }],
    ['not-applicable leaked into text source candidates', (candidate) => {
      const allOf = candidate.components.schemas.PublicationTextSourceEvidence.allOf;
      const conditional = allOf.find((entry) => entry.if && entry.then);
      conditional.then.properties.candidates.additionalProperties = { $ref: '#/components/schemas/PublicationNotApplicableValue' };
    }],
    ['requirement value/source coupling removed', (candidate) => {
      candidate.components.schemas.PublicationRequirement.allOf = [];
    }],
    ['response publication context omitted', (candidate) => {
      candidate.components.schemas.PublicationRequirements.required = candidate.components.schemas.PublicationRequirements.required.filter((field) => field !== 'publication_context');
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
