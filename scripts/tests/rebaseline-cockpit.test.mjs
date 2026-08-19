import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';

const repoRoot = path.resolve(path.dirname(new URL(import.meta.url).pathname), '../..');
const cockpitPath = path.join(repoRoot, 'docs/engineering/rebaseline/cockpit.html');

assert.ok(fs.existsSync(cockpitPath), 'cockpit.html must exist');
const html = fs.readFileSync(cockpitPath, 'utf8');

const required = [
  'VISUAL PROJECTION — NOT ARCHITECTURE AUTHORITY',
  '8c685eabf9594459ded4206dfa5092348fffee22',
  'W3',
  'ACCEPTED / CANONICAL',
  'W4',
  'Permission → Operation / Client-Class Enforcement',
  'Implementation blocked until D9',
  'D0', 'D1', 'D2', 'D3', 'D4', 'D5', 'D6', 'D7', 'D8', 'D9',
  'One business meaning → one authority',
  'Mechanism ≠ Authority',
  'Unknown ≠ zero',
  'Copy New Session Prompt',
  'docs/engineering/rebaseline/README.md',
  'D5-B2-W3-COLLECTION-GRAMMAR.md',
  'Marketplace Operations Control Plane',
  'ListingIntent',
  'PriceIntent',
  'Availability',
  'BusinessOrderIntent',
  'FulfillmentExecution',
  'PostSaleResolution',
  'Operational Work',
];

for (const token of required) {
  assert.ok(html.includes(token), `cockpit.html missing required token: ${token}`);
}

assert.match(html, /id="copy-handoff"/);
assert.match(html, /navigator\.clipboard\.writeText/);
assert.match(html, /data-status="accepted"/);
assert.match(html, /data-status="active"/);
assert.match(html, /data-status="blocked"/);
assert.match(html, /data-status="deferred"/);
assert.match(html, /<svg[\s>]/);
assert.match(html, /@media\s*\(/);

console.log('rebaseline cockpit contract: PASS');
