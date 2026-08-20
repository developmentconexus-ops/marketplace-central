import { createHash } from 'node:crypto';
import { cpSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const contractDir = join(root, 'contracts/api/product');
const entrypoint = join(contractDir, 'openapi.yaml');
const redoclyConfig = join(contractDir, 'redocly.yaml');
const sourceNames = ['openapi.yaml','components.yaml','paths-identity-portfolio-readiness.yaml','paths-offering-availability-market.yaml','paths-economics-governance-sales-materialization.yaml','paths-fulfillment-postsale-work.yaml'];
const sourceTreeText = sourceNames.map((name) => readFileSync(join(contractDir, name), 'utf8')).join('\n');
const componentsText = readFileSync(join(contractDir, 'components.yaml'), 'utf8');
const w4Text = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md'), 'utf8');
const w2Text = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md'), 'utf8');
const admissionText = readFileSync(join(root, 'docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md'), 'utf8');
const temp = mkdtempSync(join(tmpdir(), 'mpc-product-oad-'));
const npx = process.platform === 'win32' ? 'npx.cmd' : 'npx';
const go = process.platform === 'win32' ? 'go.exe' : 'go';
const HTTP_METHODS = new Set(['get','put','post','delete','options','head','patch','trace']);
const PROBLEM_PREFIX = 'https://conexus.fun/marketplace-central/problems/product/';

const fail = (message) => { throw new Error(message); };
const assert = (condition, message) => { if (!condition) fail(message); };
const normalize = (values) => [...new Set(values)].sort();
function sameSet(actual, expected, label) {
  const a = normalize(actual), e = normalize(expected);
  assert(JSON.stringify(a) === JSON.stringify(e), `${label} mismatch\nactual=${JSON.stringify(a)}\nexpected=${JSON.stringify(e)}`);
}
function run(command, args, options = {}) {
  const result = spawnSync(command, args, { cwd: options.cwd ?? root, env: {...process.env, ...(options.env ?? {})}, encoding: 'utf8', shell: false });
  if (result.error) fail(`${command} failed to start: ${result.error.message}`);
  if (result.status !== 0) fail([`${command} ${args.join(' ')} failed with exit ${result.status}`, result.stdout?.trim(), result.stderr?.trim()].filter(Boolean).join('\n'));
  return result;
}
const sha256 = (path) => createHash('sha256').update(readFileSync(path)).digest('hex');
const clone = (value) => structuredClone(value);
function resolveRef(document, value, seen = new Set()) {
  if (!value || typeof value !== 'object' || typeof value.$ref !== 'string') return value;
  const ref = value.$ref;
  assert(ref.startsWith('#/'), `resolved bundle contains non-local ref: ${ref}`);
  assert(!seen.has(ref), `cyclic ref reached proof resolver: ${ref}`);
  let current = document;
  for (const raw of ref.slice(2).split('/')) {
    const key = raw.replaceAll('~1','/').replaceAll('~0','~');
    assert(current && Object.hasOwn(current, key), `unresolved bundled ref: ${ref}`);
    current = current[key];
  }
  return resolveRef(document, current, new Set(seen).add(ref));
}
function walk(value, visitor) { if (!value || typeof value !== 'object') return; visitor(value); for (const child of Object.values(value)) walk(child, visitor); }
function operations(document) {
  const result = [];
  for (const [path, raw] of Object.entries(document.paths ?? {})) {
    const pathItem = resolveRef(document, raw);
    for (const [method, operation] of Object.entries(pathItem ?? {})) if (HTTP_METHODS.has(method)) result.push({path, method, operation, pathItem});
  }
  return result;
}
const parameters = (document, entry) => [...(entry.pathItem.parameters ?? []), ...(entry.operation.parameters ?? [])].map((v) => resolveRef(document, v));
const response = (document, operation, status) => operation.responses?.[status] ? resolveRef(document, operation.responses[status]) : undefined;
const responseSchema = (document, operation, status='200') => { const raw = response(document, operation, status)?.content?.['application/json']?.schema; return raw ? resolveRef(document, raw) : undefined; };
const requestSchema = (document, operation, media='application/json') => { const raw = resolveRef(document, operation.requestBody)?.content?.[media]?.schema; return raw ? resolveRef(document, raw) : undefined; };
function byId(all, id) { const found = all.find((entry) => entry.operation.operationId === id); assert(found, `operation not found: ${id}`); return found; }
function component(document, kind, name) {
  const direct = document.components?.[kind]?.[name];
  if (direct) return resolveRef(document, direct);
  const found = Object.entries(document.components?.[kind] ?? {}).find(([candidate]) => candidate === name || candidate.endsWith(`_${name}`) || candidate.endsWith(name));
  assert(found, `bundled component not found: ${kind}.${name}`);
  return resolveRef(document, found[1]);
}

function parseW4(text) {
  const start = text.indexOf('# 8. Exact 95-operation enforcement matrix'), end = text.indexOf('\n---\n\n## 9.', start);
  assert(start >= 0 && end > start, 'unable to locate W4 exact matrix');
  const rows = [];
  for (const line of text.slice(start,end).split(/\r?\n/)) {
    const m = line.match(/^\|\s*`([^`]+)`\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$/); if (!m) continue;
    const [,operationId,classCell,permissionCell,principalCell] = m;
    const operationClass = classCell.includes('Q') ? 'Q' : classCell.includes('C') ? 'C' : null;
    const permission = permissionCell.match(/`([^`]+)`/)?.[1] ?? (permissionCell.includes('authenticated special condition') ? 'authenticated' : null);
    const principalKinds = [...principalCell.matchAll(/\b([HAS])\b/g)].map((match) => match[1]);
    assert(operationClass && permission && principalKinds.length, `unparseable W4 row: ${line}`);
    rows.push({operationId, operationClass, permission, principalKinds: normalize(principalKinds), qualifiedPhysical: /currently qualified S/.test(principalCell)});
  }
  assert(rows.length === 95, `W4 row count changed: ${rows.length}`);
  assert(new Set(rows.map((row) => row.operationId)).size === 95, 'W4 operation IDs are not unique');
  return rows;
}
function parseProblemSlugs(text) {
  const start = text.indexOf('# 19. Problem Details catalog'), end = text.indexOf('Problem `type` is an absolute stable URI', start);
  assert(start >= 0 && end > start, 'unable to locate W2 Problem catalog');
  const slugs = [...text.slice(start,end).matchAll(/^- `([^`]+)`[;.]/gm)].map((m) => m[1]);
  assert(slugs.length === 15, `W2 Product Problem count changed: ${slugs.length}`);
  return slugs;
}
function admissionSafetyRows(text) {
  const start = text.indexOf('# 7. Whole-Matrix complete C-operation safety sweep'), end = text.indexOf('### 7.6 Safety-sweep result', start);
  assert(start >= 0 && end > start, 'unable to locate admission safety sweep');
  const rows = [];
  for (const line of text.slice(start,end).split(/\r?\n/)) {
    const m = line.match(/^\|\s*`([^`]+)`\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|\s*$/);
    if (m) rows.push({operationId:m[1], consequence:m[2], idempotency:m[3], precondition:m[4]});
  }
  return rows;
}

const w4 = parseW4(w4Text);
const w4ById = new Map(w4.map((row) => [row.operationId,row]));
const expectedPermissions = normalize(w4.map((row) => row.permission).filter((permission) => permission !== 'authenticated'));
const expectedProblems = parseProblemSlugs(w2Text).map((slug) => `${PROBLEM_PREFIX}${slug}`);
const safetyRows = admissionSafetyRows(admissionText);
const expectedIdempotency = normalize(safetyRows.filter((row) => /mandatory client key/i.test(row.idempotency)).map((row) => row.operationId));
const expectedCOperations = normalize(w4.filter((row) => row.operationClass === 'C').map((row) => row.operationId));
assert(expectedPermissions.length === 29, `ordinary Permission count changed: ${expectedPermissions.length}`);
assert(expectedIdempotency.length === 14, `mandatory Idempotency-Key disposition count changed: ${expectedIdempotency.length}`);
sameSet(safetyRows.map((row) => row.operationId), expectedCOperations, 'W4/admission C-operation coverage');

function validate(document) {
  assert(document.openapi === '3.1.2', `OpenAPI must be 3.1.2, found ${document.openapi}`);
  assert(document.jsonSchemaDialect === 'https://spec.openapis.org/oas/3.1/dialect/base', 'unexpected JSON Schema dialect');
  assert(!Object.hasOwn(document,'servers'), 'canonical/resolved Product OAD must omit servers');
  assert(JSON.stringify(document.security) === JSON.stringify([{MpcBearerAuth:[]}]), 'root security must be exactly MpcBearerAuth');
  const bearer = document.components?.securitySchemes?.MpcBearerAuth;
  assert(bearer?.type === 'http' && bearer?.scheme === 'bearer', 'MpcBearerAuth must be HTTP bearer');

  const all = operations(document), ids = all.map((entry) => entry.operation.operationId);
  assert(all.length === 95, `Product operation count must be 95, found ${all.length}`);
  assert(new Set(ids).size === 95, 'operationId values are not unique');
  sameSet(ids, [...w4ById.keys()], 'OAD/W4 operation IDs');
  for (const path of Object.keys(document.paths ?? {})) {
    assert(path === '/access-context' || path.startsWith('/organizations/{organization_id}/'), `path outside canonical Product roots: ${path}`);
    assert(!/(^|\/)(providers?|integrations?|webhooks?|oauth|callbacks?|acquisition)(\/|$|:)/i.test(path), `technical/provider path leaked into Product API: ${path}`);
    assert(!path.startsWith('/v1'), `version prefix leaked into Product path: ${path}`);
  }

  const allowedExtensions = new Set(['x-mpc-operation-class','x-mpc-required-permission','x-mpc-principal-kinds','x-mpc-semantic-owner','x-mpc-required-physical-qualification']);
  for (const entry of all) {
    const op = entry.operation, expected = w4ById.get(op.operationId);
    assert(/^[A-Z][A-Za-z0-9]*$/.test(op.operationId), `operationId is not stable PascalCase: ${op.operationId}`);
    assert(!op.operationId.startsWith('AcquireMarketplace'), `internal acquisition operation leaked: ${op.operationId}`);
    for (const key of ['x-mpc-operation-class','x-mpc-required-permission','x-mpc-principal-kinds','x-mpc-semantic-owner']) assert(Object.hasOwn(op,key), `${op.operationId} missing ${key}`);
    for (const key of Object.keys(op).filter((key) => key.startsWith('x-'))) assert(allowedExtensions.has(key), `${op.operationId} contains non-allowlisted extension ${key}`);
    assert(!Object.hasOwn(op,'security'), `${op.operationId} overrides global Product security`);
    assert(!Object.hasOwn(op.responses ?? {},'202'), `${op.operationId} introduces generic 202`);
    assert(op['x-mpc-operation-class'] === expected.operationClass, `${op.operationId} class mismatch`);
    assert(op['x-mpc-required-permission'] === expected.permission, `${op.operationId} permission mismatch`);
    sameSet(op['x-mpc-principal-kinds'], expected.principalKinds, `${op.operationId} principal kinds`);
    assert(op['x-mpc-principal-kinds'].every((kind) => ['H','A','S'].includes(kind)), `${op.operationId} introduces non-H/A/S principal kind`);
    assert((op['x-mpc-required-physical-qualification'] === true) === expected.qualifiedPhysical, `${op.operationId} physical qualification mismatch`);
  }

  sameSet(component(document,'schemas','Permission').enum ?? [], expectedPermissions, 'Permission enum');
  sameSet(all.map((entry) => entry.operation['x-mpc-required-permission']).filter((p) => p !== 'authenticated'), expectedPermissions, 'operation Permission vocabulary');
  const listSearch = all.filter((entry) => /^(List|Search)/.test(entry.operation.operationId));
  assert(listSearch.length === 26, `List/Search operation count must be 26, found ${listSearch.length}`);
  assert(listSearch.filter((entry) => entry.operation.operationId.startsWith('Search')).length === 1, 'exactly one Search operation is required');
  for (const entry of listSearch) {
    const names = parameters(document,entry).filter((p) => p.in === 'query').map((p) => p.name);
    assert(names.includes('limit') && names.includes('cursor'), `${entry.operation.operationId} lacks limit/cursor`);
    assert(!names.some((name) => ['page','offset','sort','total','total_count'].includes(name)), `${entry.operation.operationId} contains offset/sort/total grammar`);
    const schema = responseSchema(document,entry.operation);
    assert(schema?.type === 'object' && Object.hasOwn(schema.properties ?? {},'next_cursor') && !(schema.required ?? []).includes('next_cursor'), `${entry.operation.operationId} collection cursor response is non-canonical`);
  }
  const limit = component(document,'parameters','Limit');
  assert(limit.schema?.type === 'integer' && limit.schema.minimum === 1 && Number.isInteger(limit.schema.default) && Number.isInteger(limit.schema.maximum) && limit.schema.default > 0 && limit.schema.maximum >= limit.schema.default, 'Limit bounds are incoherent');

  for (const [name,raw] of Object.entries(document.components?.schemas ?? {})) {
    const schema = resolveRef(document,raw);
    if (schema?.type === 'object') assert(schema.additionalProperties === false, `${name} semantic object is not closed`);
    if (Array.isArray(schema?.oneOf)) for (const branchRef of schema.oneOf) {
      const branch = resolveRef(document,branchRef);
      assert(Object.values(branch?.properties ?? {}).some((property) => property && typeof property === 'object' && Object.hasOwn(property,'const')), `${name} oneOf branch lacks const discriminator`);
    }
  }
  const decimal = component(document,'schemas','ExactDecimalString');
  assert(decimal.type === 'string' && typeof decimal.pattern === 'string' && !decimal.pattern.toLowerCase().includes('e'), 'exact decimal grammar is not a non-exponent string');

  const customProblemTypes = [];
  walk(document,(node) => { const value = node?.properties?.type?.const; if (typeof value === 'string' && value.startsWith(PROBLEM_PREFIX)) customProblemTypes.push(value); });
  sameSet(customProblemTypes, expectedProblems, 'Product Problem catalog/origin');
  assert(/MethodNotAllowed:[\s\S]{0,1200}?Allow:/.test(componentsText), '405 MethodNotAllowed source response lacks Allow');
  for (const [name,status] of [['AboutBlank413Problem',413],['AboutBlank415Problem',415]]) {
    const schema = component(document,'schemas',name);
    assert(schema.properties?.type?.const === 'about:blank' && schema.properties?.status?.const === status, `${name} must be about:blank/${status}`);
  }

  const idempotencyOps = all.filter((entry) => parameters(document,entry).some((p) => p.in === 'header' && p.name.toLowerCase() === 'idempotency-key')).map((entry) => entry.operation.operationId);
  sameSet(idempotencyOps, expectedIdempotency, 'Idempotency-Key carrier set');
  const ifMatchExpected = ['UpdateMarketplaceInstallationConfiguration','UpdateListingIntentDraft','UpdateInventorySource','UpdateAvailabilityAllocationScopePolicy','UpdateCommercialPolicy','UpdateAuthorizationDelegation','UpdateFulfillmentNode','UpdateFulfillmentOperatingTargets'];
  const ifMatchActual = all.filter((entry) => parameters(document,entry).some((p) => p.in === 'header' && p.name.toLowerCase() === 'if-match')).map((entry) => entry.operation.operationId);
  sameSet(ifMatchActual, ifMatchExpected, 'If-Match carrier set');

  for (const id of ['DeactivateMarketplaceInstallation','DiscardListingIntentDraft','SubmitListingIntent','DeactivateInventorySource','ResolveEconomicAttribution','ResolveSaleSellingEntityAttribution','ResolveBusinessSystemPartyResolution','DeactivateFulfillmentNode','RecordSeparation','RecordPhysicalConference','RecordPacking','RecordDispatchHandoff','AssignWork','ClearWorkAssignment','HoldWork','ResumeWork','EscalateWork']) {
    const schema = requestSchema(document,byId(all,id).operation); assert(schema?.properties?.etag && (schema.required ?? []).includes('etag'), `${id} lacks required typed etag`);
  }
  for (const id of ['ResolveProductChannelCorrespondence','ClearProductChannelCorrespondence']) {
    const schema = requestSchema(document,byId(all,id).operation); assert(schema?.properties?.correspondence_etag && (schema.required ?? []).includes('correspondence_etag'), `${id} lacks required correspondence_etag`);
  }
  for (const branch of component(document,'schemas','AuthorizationTarget').oneOf ?? []) { const schema = resolveRef(document,branch); assert(schema.properties?.etag && (schema.required ?? []).includes('etag'), 'AuthorizationDecision target branch lacks exact ETag'); }
  const superseded = component(document,'schemas','SupersededPriceIntentRef'); assert(superseded.properties?.etag && (superseded.required ?? []).includes('etag'), 'superseded PriceIntent ref lacks exact ETag');

  const durableCreates = ['CreateMarketplaceInstallation','CreateListingIntentDraft','CreatePriceIntent','CreateInventorySource','CreateAuthorizationDecision','EstablishAuthorizationDelegation','CreateFulfillmentNode','CreatePostSaleResolution'];
  for (const id of durableCreates) assert(response(document,byId(all,id).operation,'201')?.headers?.Location, `${id} must return 201 + Location`);
  for (const entry of all.filter((entry) => entry.method === 'post' && !durableCreates.includes(entry.operation.operationId))) assert(response(document,entry.operation,'200'), `${entry.operation.operationId} custom POST must return 200`);

  const media = byId(all,'CreateListingIntentMedia').operation, mediaContent = resolveRef(document,media.requestBody)?.content?.['multipart/form-data'], mediaSchema = resolveRef(document,mediaContent?.schema);
  sameSet(Object.keys(mediaSchema?.properties ?? {}), ['file','etag'], 'CreateListingIntentMedia multipart parts');
  sameSet(mediaSchema?.required ?? [], ['file','etag'], 'CreateListingIntentMedia required multipart parts');
  assert(mediaContent?.encoding?.etag?.contentType === 'text/plain', 'CreateListingIntentMedia etag part must be text/plain');

  const workOriginKinds = (component(document,'schemas','WorkOrigin').oneOf ?? []).map((branch) => { const schema = resolveRef(document,branch); return Object.values(schema.properties ?? {}).find((p) => p && typeof p === 'object' && Object.hasOwn(p,'const'))?.const; }).filter(Boolean);
  sameSet(workOriginKinds, ['readiness','listing_intent','price_intent','marketplace_sale','business_order_intent','fulfillment_execution','shipment','post_sale_resolution','economic_attribution','authorization_decision'], 'owner-specific Work origins');
  return {all,idempotencyOps,listSearch};
}

function expectFailure(name,factory) { let failed=false; try { validate(factory()); } catch { failed=true; } assert(failed, `negative control unexpectedly passed: ${name}`); }
function negativeSemanticControls(document) {
  expectFailure('missing operation',()=>{const c=clone(document);delete c.paths['/organizations/{organization_id}/work/{work_id}:escalate'];return c;});
  expectFailure('wrong W4 permission',()=>{const c=clone(document);byId(operations(c),'CreatePriceIntent').operation['x-mpc-required-permission']='listing.manage';return c;});
  expectFailure('technical path leakage',()=>{const c=clone(document);c.paths['/organizations/{organization_id}/webhooks']=clone(c.paths['/organizations/{organization_id}/work']);return c;});
  expectFailure('fourth principal kind',()=>{const c=clone(document);byId(operations(c),'GetCurrentAccessContext').operation['x-mpc-principal-kinds'].push('X');return c;});
  expectFailure('missing idempotency carrier',()=>{const c=clone(document),e=byId(operations(c),'CreateInventorySource');e.operation.parameters=(e.operation.parameters??[]).filter((raw)=>resolveRef(c,raw)?.name!=='Idempotency-Key');return c;});
  expectFailure('wrong Product Problem origin',()=>{const c=clone(document);component(c,'schemas','AccessDeniedProblem').properties.type.const='https://example.invalid/problem';return c;});
}
function redoclyProof() {
  run(npx,['--yes','@redocly/cli@2.45.0','lint',entrypoint,'--config',redoclyConfig]);
  const a=join(temp,'product-a.json'),b=join(temp,'product-b.json');
  run(npx,['--yes','@redocly/cli@2.45.0','bundle',entrypoint,'--config',redoclyConfig,'-o',a]); run(npx,['--yes','@redocly/cli@2.45.0','bundle',entrypoint,'--config',redoclyConfig,'-o',b]);
  assert(sha256(a)===sha256(b),'resolved Product bundle is not deterministic');
  const fixture=join(temp,'broken-source');mkdirSync(fixture,{recursive:true});for(const name of sourceNames)cpSync(join(contractDir,name),join(fixture,name));cpSync(redoclyConfig,join(fixture,'redocly.yaml'));
  const brokenEntry=join(fixture,'openapi.yaml'),original=readFileSync(brokenEntry,'utf8'),broken=original.replace('#/pathItems/AccessContext','#/pathItems/__missing__');assert(broken!==original,'unable to construct unresolved-ref negative fixture');writeFileSync(brokenEntry,broken,'utf8');
  assert(spawnSync(npx,['--yes','@redocly/cli@2.45.0','lint',brokenEntry,'--config',join(fixture,'redocly.yaml')],{cwd:root,encoding:'utf8',shell:false}).status!==0,'Redocly unresolved-ref negative fixture unexpectedly passed');
  return {document:JSON.parse(readFileSync(a,'utf8')),bundle:a};
}
function typescriptProof() {
  const a=join(temp,'product-a.d.ts'),b=join(temp,'product-b.d.ts');run(npx,['--yes','openapi-typescript@7.13.0',entrypoint,'-o',a,'--redocly',redoclyConfig]);run(npx,['--yes','openapi-typescript@7.13.0',entrypoint,'-o',b,'--redocly',redoclyConfig]);
  assert(sha256(a)===sha256(b),'TypeScript projection is not deterministic');const generated=readFileSync(a,'utf8');assert(generated.includes('This file was auto-generated by openapi-typescript'),'TypeScript projection lacks generated banner');assert(generated.includes('GetCurrentAccessContext')&&generated.includes('EscalateWork'),'TypeScript projection lacks boundary operation IDs');run(npx,['--yes','-p','typescript@5.9.3','tsc','--noEmit','--strict','--skipLibCheck',a]);
}
function goProof(bundle) {
  const dir=join(temp,'go-proof');mkdirSync(dir,{recursive:true});writeFileSync(join(dir,'oapi-codegen.yaml'),['package: productapi','output: product.gen.go','generate:','  models: true','  std-http-server: true','  strict-server: true','output-options:','  skip-prune: true',''].join('\n'));
  const generate=()=>run(go,['run','github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0','-config','oapi-codegen.yaml',bundle],{cwd:dir});generate();const generatedPath=join(dir,'product.gen.go'),first=sha256(generatedPath);generate();assert(sha256(generatedPath)===first,'Go projection is not deterministic');
  const generated=readFileSync(generatedPath,'utf8');assert(generated.includes('Code generated by github.com/oapi-codegen/oapi-codegen'),'Go projection lacks generated banner');assert(generated.includes('type StrictServerInterface interface'),'Go projection lacks StrictServerInterface');assert(generated.includes('GetCurrentAccessContext')&&generated.includes('EscalateWork'),'Go projection lacks boundary operation IDs');assert(!/\"(os\/exec|unsafe|syscall)\"/.test(generated),'Go projection contains forbidden high-risk import');
  writeFileSync(join(dir,'go.mod'),['module productoadproof','','go 1.25.1','','require github.com/oapi-codegen/runtime v1.7.0',''].join('\n'));writeFileSync(join(dir,'colon_suffix_test.go'),`package productapi\nimport(\"net/http\";\"net/http/httptest\";\"testing\")\ntype exactSuffixMux struct{method,path string;handler http.Handler}\nfunc(m exactSuffixMux)ServeHTTP(w http.ResponseWriter,r *http.Request){if r.Method==m.method&&r.URL.Path==m.path{m.handler.ServeHTTP(w,r);return};http.NotFound(w,r)}\nfunc TestCanonicalColonSuffixDispatch(t *testing.T){called:=false;mux:=exactSuffixMux{http.MethodPost,\"/organizations/org/listing-intents/li:submit\",http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request){called=true;w.WriteHeader(http.StatusNoContent)})};res:=httptest.NewRecorder();mux.ServeHTTP(res,httptest.NewRequest(http.MethodPost,\"/organizations/org/listing-intents/li:submit\",nil));if !called||res.Code!=http.StatusNoContent{t.Fatalf(\"colon-suffix dispatch failed\")}}\n`);
  run(go,['mod','tidy'],{cwd:dir});assert(run(go,['list','-m','-f','{{.Version}}','github.com/oapi-codegen/runtime'],{cwd:dir}).stdout.trim()==='v1.7.0','unexpected oapi-codegen runtime');run(go,['test','./...'],{cwd:dir,env:{GOTOOLCHAIN:'go1.25.1'}});const minimum=run(go,['version'],{cwd:dir,env:{GOTOOLCHAIN:'go1.25.1'}}).stdout.trim();assert(/go1\.25\.1\b/.test(minimum),`minimum toolchain proof did not execute go1.25.1: ${minimum}`);run(go,['test','./...'],{cwd:dir});console.log(`go_minimum=${minimum}`);console.log(`go_current=${run(go,['version'],{cwd:dir}).stdout.trim()}`);
}
function sourceTreeProof(){assert(!sourceTreeText.includes('ngrok-free.dev'),'preview ngrok origin leaked into Product source tree');assert(!sourceTreeText.includes('x-go-'),'forbidden x-go-* override leaked into Product source tree');assert(!sourceTreeText.includes('problems/technical/'),'technical Problem namespace leaked into Product source tree');assert(!/(^|\n)\s*servers\s*:/m.test(sourceTreeText),'servers leaked into Product source tree');const refs=[...sourceTreeText.matchAll(/\$ref:\s*['"]?([^'"}\s]+)['"]?/g)].map((m)=>m[1]);assert(refs.length>0,'Product source tree contains no refs');assert(refs.every((ref)=>ref.startsWith('./')||ref.startsWith('#/')),`remote/non-local Product source ref found: ${refs.find((ref)=>!ref.startsWith('./')&&!ref.startsWith('#/'))}`);}
function legacyRuntimeProof(){const tracked=run('git',['-C',root,'ls-files']).stdout.split(/\r?\n/).filter(Boolean);for(const prefix of ['apps/','cmd/','internal/','server/','backend/','frontend/'])assert(!tracked.some((path)=>path.startsWith(prefix)),`active legacy/runtime population found under ${prefix}`);sameSet(tracked.filter((path)=>/(^|\/)openapi[^/]*\.ya?ml$/i.test(path)),['contracts/api/product/openapi.yaml'],'tracked OpenAPI entrypoint authority');}

try {
  const before=run('git',['-C',root,'status','--porcelain=v1']).stdout;sourceTreeProof();legacyRuntimeProof();const {document,bundle}=redoclyProof();const result=validate(document);negativeSemanticControls(document);typescriptProof();goProof(bundle);const after=run('git',['-C',root,'status','--porcelain=v1']).stdout;assert(after===before,'Product OAD proof dirtied the repository');
  console.log('product_oad_openapi=3.1.2');console.log(`product_oad_operations=${result.all.length}/95`);console.log(`product_oad_permissions=${expectedPermissions.length}/29`);console.log('product_oad_principal_kinds=H/A/S');console.log('product_oad_stable_origin=https://conexus.fun');console.log(`product_oad_idempotency_carriers=${result.idempotencyOps.length}/14`);console.log(`product_oad_collection_operations=${result.listSearch.length}/26`);console.log('product_oad_negative_controls=7/7');console.log('product_oad_legacy_runtime_population=0');console.log('product_oad_runtime_schema_enforcement=NOT_CLAIMED_D7');console.log('product_oad_router_selection=NONE_D7');console.log('product_oad_proof=PASS');
} finally { rmSync(temp,{recursive:true,force:true}); }
