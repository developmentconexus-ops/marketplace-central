// @vitest-environment node
import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import { KNOWN_IDENTITY_ANCHORS, MERCADO_LIVRE_SUPPLIED_ANCHORS } from "./wireFixtures";

/**
 * A divergence detector for the seam that has no compiler.
 *
 * `wireFixtures.ts` decides which candidates are producible by encoding facts
 * that live in Go. Nothing links the two: the FE type system cannot see the
 * generator, and the Go build cannot see these fixtures. That gap is exactly the
 * shape of the defect this whole chip exists to fix (a union enumerated as
 * string literals: type-correct, and therefore silent when the union grows), one
 * layer up.
 *
 * WHAT THIS IS NOT, said first so the pack cannot over-read a green run: it is
 * NOT a compiler, and reading Go with regexes is string-scraping. A rename, a
 * moved file, or a literal replaced by a generated `const` makes the patterns
 * below match NOTHING — and a guard that finds nothing would otherwise PASS,
 * green with the Go intact and green with the Go gone. That is the same stable
 * observable that fails to discriminate as the defects this chip has been
 * fixing all week.
 *
 * So the first test exists to make that failure loud: the sources must be where
 * this file expects them, non-empty, and still carrying the symbols by name. A
 * guard that cannot see its subject FAILS instead of passing vacuously. What
 * that buys is narrower than a compiler and worth stating exactly: this is a
 * detector that knows how to recognize its own blindness. It cannot prove the
 * FE matches a DEPLOYED server, and it does not try to; it proves the FE's copy
 * matches the declaration in THIS checkout, and refuses to answer at all when
 * it can no longer read that declaration.
 */

const here = dirname(fileURLToPath(import.meta.url));
const serverModules = resolve(here, "../../../../server_core/internal/modules");

/**
 * The seam, as paths and the EXTRACTION ITSELF. `extract` is not a description
 * of the pattern used below — it is the pattern used below, and the sentinel
 * asserts with the same object.
 *
 * It was two fields before, a `symbol` string for the sentinel and a separate
 * regex inside each test, and the hub's executor seat measured what that costs:
 * the sentinel asserted `.toContain("var knownIdentityAnchors")`, which
 * `var knownIdentityAnchorsXX` SATISFIES. A rename that appends leaves the
 * substring intact, so the sentinel passed while the extraction found nothing —
 * the guard was wider than the fact it claimed to check, which is the class this
 * chip keeps finding one layer down. Sharing the object is what makes sentinel
 * and extraction unable to disagree; a tighter second pattern would only have
 * moved the disagreement somewhere harder to see.
 */
const GO_SEAM = {
  capabilityPort: {
    path: "connectors/ports/marketplace_capability.go",
    symbol: "var knownIdentityAnchors = []IdentityAnchor{…}",
    extract: /var knownIdentityAnchors = \[\]IdentityAnchor\{([^}]*)\}/,
    what: "the identity-anchor vocabulary itself",
  },
  mercadoLivreCapability: {
    path: "connectors/adapters/mercado_livre/capability_adapter.go",
    symbol: "IdentityAnchors: []ports.IdentityAnchor{…}",
    extract: /IdentityAnchors:\s*\[\]ports\.IdentityAnchor\{([^}]*)\}/,
    what: "which anchors mercado_livre declares supplied",
  },
} as const;

type SeamKey = keyof typeof GO_SEAM;

// Read at MODULE scope, not inside the tests. These are the only files in this
// suite outside the vitest root, and the first of these reads once blew the 5s
// per-test budget while the very next test read the same two files in 2ms. The
// cause was not established — three later runs, including a probe built to
// reproduce it, all came back in single-digit milliseconds — so this is a
// mitigation, not a diagnosis: module-load time is not charged to testTimeout,
// so a slow cold read cannot turn a passing guard into a red lane.
//
// A read FAILURE is captured rather than thrown, so it surfaces as the named
// assertion below instead of an ENOENT that takes the whole file down with a
// message about a path and nothing about what the path was for.
const SOURCES: Record<SeamKey, { text: string; error: string | null }> = {
  capabilityPort: readSeam("capabilityPort"),
  mercadoLivreCapability: readSeam("mercadoLivreCapability"),
};

function readSeam(key: SeamKey): { text: string; error: string | null } {
  try {
    return { text: readFileSync(resolve(serverModules, GO_SEAM[key].path), "utf8"), error: null };
  } catch (cause) {
    return { text: "", error: String(cause) };
  }
}

/**
 * Every later test goes through this. An unreadable source must never reach an
 * extraction as an empty string, because "extracted nothing" and "there is
 * nothing to extract" are the two states this guard exists to tell apart.
 */
function sourceOf(key: SeamKey): string {
  const source = SOURCES[key];
  if (source.error !== null || source.text.length === 0) {
    throw new Error(
      `the Go source for ${GO_SEAM[key].what} could not be read at ${GO_SEAM[key].path} — ` +
        "this guard cannot answer, and an unreadable source is a FAILURE here rather than a skip. " +
        `Underlying: ${source.error ?? "file was empty"}`,
    );
  }
  return source.text;
}

/** `IdentityAnchorSellerSKU IdentityAnchor = "seller_sku"` → { IdentityAnchorSellerSKU: "seller_sku" } */
function identityAnchorConstants(source: string): Map<string, string> {
  const constants = new Map<string, string>();
  for (const match of source.matchAll(/(IdentityAnchor\w+)\s+IdentityAnchor\s*=\s*"([^"]+)"/g)) {
    constants.set(match[1], match[2]);
  }
  return constants;
}

describe("wireFixtures vocabulary vs the Go source it mirrors", () => {
  // ARM (b) of the must-fail. Without this test, every assertion below can be
  // satisfied by a vocabulary that moved out from under them: the patterns stop
  // matching, the extraction comes back empty, and an empty guard reports no
  // divergence. This is the test that turns "I found nothing" into red.
  it("can still SEE the Go declarations it claims to guard", () => {
    for (const key of Object.keys(GO_SEAM) as SeamKey[]) {
      const { path, symbol, what } = GO_SEAM[key];

      expect(SOURCES[key].error, `${path} — the file this guard reads for ${what} is gone or moved`).toBeNull();
      expect(SOURCES[key].text.length, `${path} is empty`).toBeGreaterThan(0);

      // The path existing is not enough: the FACT has to still be in it, in the
      // form the extraction reads. Asserted with `extract` — the same object the
      // tests below run — because a sentinel with its own looser pattern is a
      // sentinel that can be satisfied while the extraction comes back empty,
      // which is precisely how a suffix rename went undetected here.
      expect(
        GO_SEAM[key].extract.test(SOURCES[key].text),
        `${path} no longer matches \`${symbol}\` — the seam moved or was renamed, so this guard is ` +
          `reading a file that no longer declares ${what}, and every assertion after this one would ` +
          "pass vacuously",
      ).toBe(true);
    }
  });

  it("KNOWN_IDENTITY_ANCHORS is exactly knownIdentityAnchors", () => {
    const capabilityPort = sourceOf("capabilityPort");
    const constants = identityAnchorConstants(capabilityPort);
    expect(constants.size, "no `IdentityAnchor… = \"…\"` constants found — the declaration form changed").toBeGreaterThan(0);

    const block = capabilityPort.match(GO_SEAM.capabilityPort.extract);
    if (block === null) {
      throw new Error("could not find `var knownIdentityAnchors` — the declaration moved or was renamed");
    }

    const anchors = Array.from(block[1].matchAll(/(IdentityAnchor\w+)/g)).map((entry) => {
      const value = constants.get(entry[1]);
      if (value === undefined) throw new Error(`no constant found for ${entry[1]}`);
      return value;
    });
    expect(anchors.length, "the knownIdentityAnchors block parsed to zero anchors").toBeGreaterThan(0);

    // Order matters as little as the Go slice's does, so compare as sets — but
    // compare BOTH directions, because a vocabulary that only grew and a
    // vocabulary that only shrank are different failures and both are real.
    // `refforn` shrank out of this list (D-A) and left one fixture asserting a
    // reason no provider can ever emit.
    expect([...anchors].sort()).toEqual([...KNOWN_IDENTITY_ANCHORS].sort());
  });

  it("MERCADO_LIVRE_SUPPLIED_ANCHORS is exactly what its capability adapter declares", () => {
    const constants = identityAnchorConstants(sourceOf("capabilityPort"));
    expect(constants.size, "no IdentityAnchor constants to resolve the declaration against").toBeGreaterThan(0);

    const declaration = sourceOf("mercadoLivreCapability").match(GO_SEAM.mercadoLivreCapability.extract);
    expect(declaration, "could not find mercado_livre's IdentityAnchors declaration").not.toBeNull();

    const declared = Array.from((declaration as RegExpMatchArray)[1].matchAll(/ports\.(IdentityAnchor\w+)/g)).map(
      (entry) => {
        const value = constants.get(entry[1]);
        if (value === undefined) throw new Error(`no constant found for ${entry[1]}`);
        return value;
      },
    );
    expect(declared.length, "the IdentityAnchors declaration parsed to zero anchors").toBeGreaterThan(0);

    expect([...declared].sort()).toEqual([...MERCADO_LIVRE_SUPPLIED_ANCHORS].sort());
    // The consequence the fixtures depend on, asserted rather than assumed: the
    // anchor mercado_livre does NOT supply is the one every one of its
    // candidates carries as UNAVAILABLE.
    expect(declared).not.toContain("marca");
  });
});
