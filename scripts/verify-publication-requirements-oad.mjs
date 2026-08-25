import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';

const root = process.cwd();
const componentsPath = resolve(root, 'contracts/api/product/components.yaml');

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }

assert(existsSync(componentsPath), 'Product OAD components.yaml missing');
const source = readFileSync(componentsPath, 'utf8');

function schemaBlock(text, name) {
  const marker = `  ${name}:`;
  const start = text.indexOf(marker);
  assert(start >= 0, `publication-requirements schema missing: ${name}`);
  const rest = text.slice(start + marker.length);
  const next = rest.search(/\n  [A-Za-z][A-Za-z0-9]*:/u);
  return marker + (next >= 0 ? rest.slice(0, next) : rest);
}

function verify(text) {
  const context = schemaBlock(text, 'PublicationRequirementsContext');
  assert(context.includes('additionalProperties: false'), 'publication requirements context must be closed');
  assert(context.includes('category_key:'), 'publication requirements context must expose category_key');
  assert(context.includes('product_type_key:'), 'publication requirements context must expose product_type_key');

  const candidateValues = schemaBlock(text, 'PublicationSourceCandidateValues');
  assert(candidateValues.includes('type: object'), 'source candidates must use a key-addressed object so candidate identity is structurally unique');
  assert(candidateValues.includes('minProperties: 1'), 'source candidate values require at least one candidate');
  assert(candidateValues.includes("propertyNames: {$ref: '#/schemas/OpaqueKey'}"), 'source candidate identities must remain opaque Readiness keys');
  assert(candidateValues.includes("additionalProperties: {$ref: '#/schemas/PublicationValue'}"), 'each source candidate key must resolve to one PublicationValue');
  assert(!candidateValues.includes('x-mpc-unique-by'), 'candidate identity uniqueness must be structural, not an unenforced vendor annotation');

  const evidenceVariants = [
    ['PublicationSourceEvidenceKnown', 'known'],
    ['PublicationSourceEvidenceMissing', 'missing'],
    ['PublicationSourceEvidenceConflicting', 'conflicting'],
    ['PublicationSourceEvidenceUnknown', 'unknown'],
    ['PublicationSourceEvidenceUnavailable', 'unavailable'],
    ['PublicationSourceEvidenceUnsupported', 'unsupported'],
  ];
  for (const [schema, state] of evidenceVariants) {
    const block = schemaBlock(text, schema);
    assert(block.includes(`state: {const: ${state}}`), `${schema} must preserve state=${state}`);
  }
  const known = schemaBlock(text, 'PublicationSourceEvidenceKnown');
  assert(known.includes('required: [state, candidates]'), 'known source evidence must expose candidates');
  assert(known.includes("candidates: {$ref: '#/schemas/PublicationSourceCandidateValues'}"), 'known source evidence must use the structurally unique candidate-key map');
  const conflicting = schemaBlock(text, 'PublicationSourceEvidenceConflicting');
  assert(conflicting.includes('required: [state, candidates]'), 'conflicting source evidence must expose conflicting candidates');
  assert(conflicting.includes("{$ref: '#/schemas/PublicationSourceCandidateValues'}"), 'conflicting source evidence must use the structurally unique candidate-key map');
  assert(conflicting.includes('minProperties: 2'), 'conflicting source evidence requires at least two distinct candidate identities');

  const evidence = schemaBlock(text, 'PublicationSourceEvidence');
  for (const [schema] of evidenceVariants) {
    assert(evidence.includes(`#/schemas/${schema}`), `PublicationSourceEvidence union missing ${schema}`);
  }

  const textSpec = schemaBlock(text, 'PublicationTextRequirementSpec');
  assert(textSpec.includes('kind: {const: text}'), 'text requirement spec kind missing');
  assert(textSpec.includes('min_length:'), 'text requirement spec must expose minimum length when material');
  assert(textSpec.includes('max_length:'), 'text requirement spec must expose maximum length when material');

  const decimalSpec = schemaBlock(text, 'PublicationExactDecimalRequirementSpec');
  assert(decimalSpec.includes('minimum:'), 'exact-decimal requirement spec must expose minimum when material');
  assert(decimalSpec.includes('maximum:'), 'exact-decimal requirement spec must expose maximum when material');

  schemaBlock(text, 'PublicationBooleanRequirementSpec');

  const optionSpec = schemaBlock(text, 'PublicationOptionRequirementSpec');
  assert(optionSpec.includes('required: [kind, option_keys]'), 'option requirement spec must expose provider-authoritative options');
  assert(/option_keys: \{type: array, minItems: 1, uniqueItems: true,/u.test(optionSpec), 'option requirement spec requires a non-empty unique option set');

  const textListSpec = schemaBlock(text, 'PublicationTextListRequirementSpec');
  assert(textListSpec.includes('min_items:'), 'text-list requirement spec must expose minimum cardinality when material');
  assert(textListSpec.includes('max_items:'), 'text-list requirement spec must expose maximum cardinality when material');
  assert(textListSpec.includes('item_min_length:'), 'text-list requirement spec must expose item minimum length when material');
  assert(textListSpec.includes('item_max_length:'), 'text-list requirement spec must expose item maximum length when material');

  const optionListSpec = schemaBlock(text, 'PublicationOptionListRequirementSpec');
  assert(optionListSpec.includes('required: [kind, option_keys]'), 'option-list requirement spec must expose provider-authoritative options');
  assert(optionListSpec.includes('min_items:'), 'option-list requirement spec must expose minimum cardinality when material');
  assert(optionListSpec.includes('max_items:'), 'option-list requirement spec must expose maximum cardinality when material');

  const numberUnitSpec = schemaBlock(text, 'PublicationNumberUnitRequirementSpec');
  assert(numberUnitSpec.includes('required: [kind, unit_keys]'), 'number-unit requirement spec must expose allowed units');
  assert(/unit_keys: \{type: array, minItems: 1, uniqueItems: true,/u.test(numberUnitSpec), 'number-unit requirement spec requires a non-empty unique unit set');
  assert(numberUnitSpec.includes('default_unit_key:'), 'number-unit requirement spec must expose default unit when provider meaning supplies one');
  assert(numberUnitSpec.includes('minimum:'), 'number-unit requirement spec must expose minimum when material');
  assert(numberUnitSpec.includes('maximum:'), 'number-unit requirement spec must expose maximum when material');

  const valueSpec = schemaBlock(text, 'PublicationRequirementValueSpec');
  for (const schema of [
    'PublicationTextRequirementSpec',
    'PublicationExactDecimalRequirementSpec',
    'PublicationBooleanRequirementSpec',
    'PublicationOptionRequirementSpec',
    'PublicationTextListRequirementSpec',
    'PublicationOptionListRequirementSpec',
    'PublicationNumberUnitRequirementSpec',
  ]) assert(valueSpec.includes(`#/schemas/${schema}`), `PublicationRequirementValueSpec union missing ${schema}`);

  const requirement = schemaBlock(text, 'PublicationRequirement');
  assert(requirement.includes('required: [requirement_key, requirement_class, applicability, value_spec, not_applicable_allowed, source_evidence]'), 'PublicationRequirement required fields must preserve provider class, value spec and source evidence');
  assert(requirement.includes('requirement_class: {type: string, enum: [required, recommended, optional, conditional]}'), 'PublicationRequirement must preserve required/recommended/optional/conditional class');
  assert(requirement.includes('applicability: {type: string, enum: [current, draft_dependent]}'), 'PublicationRequirement must preserve current vs draft-dependent evaluation class');
  assert(requirement.includes("value_spec: {$ref: '#/schemas/PublicationRequirementValueSpec'}"), 'PublicationRequirement must use bounded value spec');
  assert(requirement.includes("source_evidence: {$ref: '#/schemas/PublicationSourceEvidence'}"), 'PublicationRequirement must use owner-local source evidence union');
  assert(!requirement.includes('required: {type: boolean}'), 'PublicationRequirement must not collapse provider requirement class to a required boolean');
  assert(!requirement.includes('source_candidates:'), 'PublicationRequirement must not leave source knowledge as an unqualified candidate array');
  assert(!requirement.includes('provider_fields'), 'PublicationRequirement must not expose provider field bags');

  const requirements = schemaBlock(text, 'PublicationRequirements');
  assert(requirements.includes('required: [subject, publication_context, requirements_revision, requirements, source_media_candidates, evaluated_at]'), 'PublicationRequirements must preserve exact publication context with its revision');
  assert(requirements.includes("publication_context: {$ref: '#/schemas/PublicationRequirementsContext'}"), 'PublicationRequirements publication_context ref missing');
}

verify(source);

let negativeControls = 0;
function expectFailure(name, mutate) {
  let failed = false;
  try { verify(mutate(source)); } catch { failed = true; }
  if (!failed) fail(`negative control unexpectedly passed: ${name}`);
  negativeControls += 1;
}

expectFailure('publication context erased', (text) => text.replace('  PublicationRequirementsContext:', '  BrokenPublicationRequirementsContext:'));
expectFailure('unsupported source state collapsed', (text) => text.replace('state: {const: unsupported}', 'state: {const: unknown}'));
expectFailure('provider requirement class collapsed', (text) => text.replace('[required, recommended, optional, conditional]', '[required, optional]'));
expectFailure('text max constraint erased', (text) => text.replace('      max_length:', '      erased_length:'));
expectFailure('candidate key-addressed representation erased', (text) => text.replace("propertyNames: {$ref: '#/schemas/OpaqueKey'}", 'propertyNames: {type: string}'));
expectFailure('candidate value contract erased', (text) => text.replace("additionalProperties: {$ref: '#/schemas/PublicationValue'}", 'additionalProperties: true'));
expectFailure('conflict distinct-identity cardinality weakened', (text) => text.replace('        - {minProperties: 2}', '        - {minProperties: 1}'));
expectFailure('response publication context omitted', (text) => text.replace('[subject, publication_context, requirements_revision, requirements, source_media_candidates, evaluated_at]', '[subject, requirements_revision, requirements, source_media_candidates, evaluated_at]'));

assert(negativeControls === 8, `publication requirements negative-control count mismatch: ${negativeControls}/8`);
console.log('publication_requirements_context=EXPLICIT');
console.log('publication_requirements_class=FOUR_WAY');
console.log('publication_requirements_value_spec=BOUNDED_BY_PUBLICATION_VALUE');
console.log('publication_requirements_candidate_identity=STRUCTURALLY_UNIQUE_BY_OBJECT_KEY');
console.log('publication_requirements_source_evidence=KNOWN_MISSING_CONFLICTING_UNKNOWN_UNAVAILABLE_UNSUPPORTED');
console.log(`publication_requirements_negative_controls=${negativeControls}/8`);
console.log('publication_requirements_oad=PASS');
