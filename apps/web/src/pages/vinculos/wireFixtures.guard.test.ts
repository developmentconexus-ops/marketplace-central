// @vitest-environment node
import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import { KNOWN_IDENTITY_ANCHORS, MERCADO_LIVRE_SUPPLIED_ANCHORS } from "./wireFixtures";

/**
 * The compiler that is missing from this seam.
 *
 * `wireFixtures.ts` decides which candidates are producible by encoding facts
 * that live in Go. Nothing links the two: the FE type system cannot see the
 * generator, and the Go build cannot see these fixtures. That gap is not
 * hypothetical — it is exactly the shape of the defect this whole chip exists
 * to fix (a union enumerated as string literals: type-correct, and therefore
 * silent when the union grows), one layer up. It was filed as an accepted
 * residual risk, which is another way of saying it was left to the next reader
 * to notice.
 *
 * So it is read at test time instead. The tests below fail the moment the Go
 * vocabulary moves, naming what moved, rather than letting every fixture in the
 * screen quietly start describing a backend that no longer exists.
 *
 * These tests read Go SOURCE, not a running backend: they prove the FE's copy
 * matches the declaration in this checkout. They cannot prove anything about a
 * deployed server, and they do not try to.
 */

const here = dirname(fileURLToPath(import.meta.url));
const serverModules = resolve(here, "../../../../server_core/internal/modules");

function readGo(relative: string): string {
  return readFileSync(resolve(serverModules, relative), "utf8");
}

// Read at MODULE scope, not inside the tests. These are the only files in this
// suite outside the vitest root, and the first of these reads once blew the 5s
// per-test budget while the very next test read the same two files in 2ms. The
// cause was not established — three later runs, including a probe built to
// reproduce it, all came back in single-digit milliseconds — so this is a
// mitigation, not a diagnosis: module-load time is not charged to testTimeout,
// so a slow cold read cannot turn a passing guard into a red lane.
const CAPABILITY_PORT = readGo("connectors/ports/marketplace_capability.go");
const MERCADO_LIVRE_CAPABILITY = readGo("connectors/adapters/mercado_livre/capability_adapter.go");

/** `IdentityAnchorSellerSKU IdentityAnchor = "seller_sku"` → { IdentityAnchorSellerSKU: "seller_sku" } */
function identityAnchorConstants(source: string): Map<string, string> {
  const constants = new Map<string, string>();
  for (const match of source.matchAll(/(IdentityAnchor\w+)\s+IdentityAnchor\s*=\s*"([^"]+)"/g)) {
    constants.set(match[1], match[2]);
  }
  return constants;
}

describe("wireFixtures vocabulary vs the Go source it mirrors", () => {
  it("KNOWN_IDENTITY_ANCHORS is exactly knownIdentityAnchors", () => {
    const constants = identityAnchorConstants(CAPABILITY_PORT);
    expect(constants.size).toBeGreaterThan(0);

    const block = CAPABILITY_PORT.match(/var knownIdentityAnchors = \[\]IdentityAnchor\{([^}]*)\}/);
    if (block === null) {
      throw new Error("could not find `var knownIdentityAnchors` — the declaration moved or was renamed");
    }

    const anchors = Array.from(block[1].matchAll(/(IdentityAnchor\w+)/g)).map((entry) => {
      const value = constants.get(entry[1]);
      if (value === undefined) throw new Error(`no constant found for ${entry[1]}`);
      return value;
    });

    // Order matters as little as the Go slice's does, so compare as sets — but
    // compare BOTH directions, because a vocabulary that only grew and a
    // vocabulary that only shrank are different failures and both are real.
    // `refforn` shrank out of this list (D-A) and left one fixture asserting a
    // reason no provider can ever emit.
    expect([...anchors].sort()).toEqual([...KNOWN_IDENTITY_ANCHORS].sort());
  });

  it("MERCADO_LIVRE_SUPPLIED_ANCHORS is exactly what its capability adapter declares", () => {
    const constants = identityAnchorConstants(CAPABILITY_PORT);

    const declaration = MERCADO_LIVRE_CAPABILITY.match(/IdentityAnchors:\s*\[\]ports\.IdentityAnchor\{([^}]*)\}/);
    expect(declaration, "could not find mercado_livre's IdentityAnchors declaration").not.toBeNull();

    const declared = Array.from((declaration as RegExpMatchArray)[1].matchAll(/ports\.(IdentityAnchor\w+)/g)).map(
      (entry) => constants.get(entry[1]) as string,
    );

    expect([...declared].sort()).toEqual([...MERCADO_LIVRE_SUPPLIED_ANCHORS].sort());
    // The consequence the fixtures depend on, asserted rather than assumed: the
    // anchor mercado_livre does NOT supply is the one every one of its
    // candidates carries as UNAVAILABLE.
    expect(declared).not.toContain("marca");
  });
});
