import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';

const repoRoot = path.resolve(path.dirname(new URL(import.meta.url).pathname), '../..');
const cockpitPath = path.join(repoRoot, 'docs/engineering/rebaseline/cockpit.html');

assert.ok(fs.existsSync(cockpitPath), 'cockpit.html must exist');
const html = fs.readFileSync(cockpitPath, 'utf8');

const required = [
  'VISUAL PROJECTION — NOT ARCHITECTURE AUTHORITY',
  'ca6cb9549e168500061fa6cb6a310d8a4a38fee8',
  'W1', 'W2', 'W3', 'W4',
  'ACCEPTED / CANONICAL',
  'Technical non-Product ingress classification',
  'Principal access eligibility',
  '95 / 95',
  '29',
  'Implementation blocked until D9',
  'D0', 'D1', 'D2', 'D3', 'D4', 'D5', 'D6', 'D7', 'D8', 'D9',
  'One business meaning → one authority',
  'Mechanism ≠ Authority',
  'Unknown ≠ zero',
  'Copy New Session Prompt',
  'docs/engineering/rebaseline/README.md',
  'D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md',
  'Marketplace Operations Control Plane',
  'ListingIntent', 'PriceIntent', 'Availability',
  'BusinessOrderIntent', 'FulfillmentExecution', 'PostSaleResolution', 'Operational Work',
];

for (const token of required) assert.ok(html.includes(token), `cockpit.html missing required token: ${token}`);
assert.match(html, /id="copy-handoff"/);
assert.match(html, /navigator\.clipboard\.writeText/);
assert.match(html, /data-status="accepted"/);
assert.match(html, /data-status="active"/);
assert.match(html, /data-status="blocked"/);
assert.match(html, /data-status="deferred"/);
assert.match(html, /<svg[\s>]/);
assert.match(html, /@media\s*\(/);
console.log('rebaseline cockpit contract: PASS');
