"""Coordinate table generator and pack REPORT for EVIDENCE.md — hub rulings R-22 and R-24.

WHAT R-24 CHANGED, AND READ IT BEFORE THE R-22 HISTORY BELOW
===========================================================
Round 6 failed here, not in the prose: the three blocking findings were all claims THIS TOOL MADE
ABOUT ITS OWN COVERAGE. The prose scan called itself a SUPERSET and a negative lookbehind made
`123:456` invisible. The mutable-axis ban called itself a ban and did not see bare `HEAD`,
`git rev-parse`, or a mutable named ref. The table claimed to resolve every anchor the prose names
and silently dropped the ones it could not find. Each claim was false in the same shape: TOTAL in
the wording, PARTIAL in the code.

The hub's ruling did not widen any recognizer. It removed the totality:

    A VERIFICATION ARTIFACT KEEPS ONLY THE CLAIMS IT CAN MAKE TOTALLY.
    A CLAIM IT CANNOT MAKE TOTALLY BECOMES A REPORT, NOT A GATE.

A report saying "this I resolved, this I did not" has no totality claim left to falsify. So the
three findings cease to exist BY CONSTRUCTION rather than by pardon — the false claim is removed,
not excused. This is ADR-17's honest-unknown rule turned on our own evidence apparatus: the same
discipline the product uses for data it does not have.

What follows from that, and what this file must NOT drift back into:

  * THE PACK RUNG IS OFF THE LADDER. Nothing merges or fails on this tool's exit code. It is a
    report the author runs and the reader may re-run.
  * NO WIDENING. TOPLEVEL is not taught grouped `var (…)`. The git-verb list is not extended. The
    lookbehind is not touched. Widening a recognizer whose output no longer gates anything buys
    nothing and is precisely the move that produced rounds 3, 4, 5 and 6.
  * The table is NAVIGATION, not evidence. A wrong row costs a reader a wrong jump, not a false
    proof; R-11's mechanism is retired with its object, its insight kept.
  * Every pass says what it does NOT see, in the same breath as what it does.

Coordinate table generator and pack-hygiene lane for EVIDENCE.md — hub ruling R-22.

WHAT R-22 CHANGED, AND WHY THIS TOOL IS SMALL
=============================================
Five P6 rounds on this chip failed on the pack's account of itself, never on the code. Three of those
rounds failed INSIDE the remedy for the previous one, each time by widening a recognizer:

  round 3  the audit recognized `path.go:NNN` only; 89 bare `:NNN` citations were invisible while the
           pack claimed every citation had been re-derived.
  round 4  a bare pass was added, but it matched a WHOLE backtick span, so `:527,:558,:582,:653`
           matched nothing. 61 invisible; the run printed ZERO.
  round 5  both passes require BACKTICKS. An unbackticked citation in a header, a comma-list, and a
           citation inside a fence were invisible — and one was stale by the same +27 shift, on the
           pack's own EXEMPLO-IO marker line.

The hub's own finding closed it, and no gate could have made it because each gate sees one round: the
pack had grown an audit of citations inside the document containing the audit, exemptions declared as
line ranges of that same document, and a count of citations in the text that stated the count.
**Verifying the pack had become more expensive than verifying the code**, which inverts what a pack is
for. The ruling that follows outlives this chip: an evidence apparatus that costs more to verify than
the thing it evidences has inverted its purpose, and the answer to that is not more apparatus.

So coordinates LEAVE THE PROSE. Prose claims become behavioural and anchored to test and symbol
names, which are stable under line shifts and checkable with `go test -run`. The coordinate table
stays — R-22 explicitly does not overturn the §11 judge, who rejected replacing the coordinate with
the symbol — but it lives HERE, generated from the code, regenerated and diffed by the lane. What is
removed is the DUPLICATION of coordinates into prose: R-14 applied to the last place a copy still
lived.

THIS TOOL THEREFORE DOES FOUR THINGS AND NOTHING ELSE
=====================================================
1. GENERATE the coordinate table, from the CODE, keyed by the anchors the prose names. If the prose
   says `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm`, this resolves that symbol in the
   tree and emits its file and line. Drift is then a diff in a generated artifact, not a defect a
   reviewer has to notice.
2. LOOK FOR COORDINATES LEFT IN PROSE. What it sees: `:<digits>` not preceded by another digit,
   anywhere on a line, inside or outside backticks and fences. What it therefore does NOT see:
   a colon-digits pair whose left neighbour is a digit, so `123:456` is invisible to it. That is a
   real hole and it is written here rather than papered over, because the claim that this scan was
   a SUPERSET is the exact sentence R-24 struck.
3. LOOK FOR MUTABLE-AXIS COMMANDS. Freezing a SHA in prose does not freeze the command. `..HEAD`
   numstat rows in this pack decayed from 13 files to 20, and a listing was called a frozen-SHA fact
   while its command was `git ls-files`, which reads the index and accepts no tree-ish. What this
   pass sees: `..HEAD`, `HEAD~`, `git ls-files`, `git status`, and `git diff|grep|ls-tree|show|log`
   on a line carrying no hex-ish token. What it does NOT see: bare `HEAD`, `git rev-parse`, a
   mutable named ref such as `main`, an ordinary worktree `grep`, or a line where an unrelated SHA
   elsewhere satisfies the hex-ish test. The principle stands whole — every evidence command names
   its own tree-ish or it is not evidence — and this pass is a partial reading of it, not the thing
   itself.
4. REPORT WHAT IT COULD NOT RESOLVE. Anchors recognized in the prose with no matching declaration,
   and backticked spans carrying a call suffix that were never treated as anchors at all. R-17 is
   ratified a fourth time by R-24: the requirement to say what could not be resolved stands. What
   R-24 changed is the CONSEQUENCE — a report line, not an exit code. And this report does not
   claim to be complete either; see its own footnote where it prints.

NO DECLARATION MAY BE A PACK LINE RANGE
=======================================
The first draft of the round-5 tool declared exemptions as pack line ranges. That is R-14 one level
down, inside R-14's own fix: a line-range exemption decays the instant a corrective inserts a
paragraph above it, and then it silently excuses the wrong lines. Ratified by the hub: every
declaration TRAVELS WITH THE CONTENT it describes.

  historical citation   carries its SHA inline           `generation_service.go:88@2921d563`
  text the chip quotes  ANY ```quoted-* fence tag        (`file.go:12:` there is the QUOTATION)
  literal non-citation  listed in cite-regions.txt       (a literal is not a coordinate)

The middle row is written as the code behaves, not as an earlier draft described it. Round 6 found
the prose naming `quoted-output` specifically while the check honoured every `quoted-*` tag; the
sentence was corrected rather than the check narrowed, because a tag is a declaration the author
makes in the open and an unaudited tag is visible in the diff that adds it.

WHAT THIS TOOL DOES NOT DECIDE
==============================
Whether the prose is true. A behavioural claim naming a real test can still misdescribe what that
test asserts. That is the cold reviewer's seat, and moving prose onto test names is what makes it
checkable — `go test -run <name>` either exists and passes or does not.

    python cite-table.py                 # GENERATE: rewrite coordinates.txt from the code
    python cite-table.py --strict        # CHECK: write nothing, exit 1 on a finding OR on drift
    python cite-table.py --show          # also dump the table to stdout
    python cite-table.py --suggest       # for each coordinate still in prose, its enclosing symbol

The two modes are separate on purpose. A checker that regenerated the artifact could never detect
drift in it — it would overwrite the evidence and report success. So `--strict` cannot write. The
exit code survives R-24 because it is useful to the author; what does not survive is anything
depending on it. The rung is off the ladder, and the unresolved report never touches the exit code.
"""
import io
import os
import re
import sys

try:
    sys.stdout.reconfigure(encoding="utf-8")
except AttributeError:  # pragma: no cover - Python < 3.7
    pass

HERE = os.path.dirname(os.path.abspath(__file__))
if "--pack-dir" in sys.argv:
    HERE = os.path.abspath(sys.argv[sys.argv.index("--pack-dir") + 1])
REPO = os.path.abspath(os.path.join(HERE, "..", "..", ".."))
PACK = os.path.join(HERE, "EVIDENCE.md")
REGIONS = os.path.join(HERE, "cite-regions.txt")

ROOTS = [
    os.path.join(REPO, "apps", "server_core", "internal"),
    os.path.join(REPO, "contracts"),
]

# ---------------------------------------------------------------------------
# The symbol map. For every top-level Go declaration: name -> [(path, line)]. Also every line's
# ENCLOSING declaration, which is what makes --suggest able to turn a coordinate into an anchor.
TOPLEVEL = re.compile(
    r"^(?:func\s+\((?:[^)]*)\)\s*(?P<method>\w+)"      # method
    r"|func\s+(?P<func>\w+)"                            # function
    r"|type\s+(?P<type>\w+)"                            # type
    r"|(?:var|const)\s+(?P<vc>\w+))"                    # single var/const
)

symbols = {}          # name -> [(relpath, lineno)]
enclosing = {}        # relpath -> [(lineno, name)] ascending


def scan_code():
    for root in ROOTS:
        for dirpath, _dirnames, filenames in os.walk(root):
            for fn in filenames:
                if not fn.endswith(".go"):
                    continue
                path = os.path.join(dirpath, fn)
                rel = os.path.relpath(path, REPO).replace("\\", "/")
                with io.open(path, encoding="utf-8", newline="") as fh:
                    lines = fh.read().replace("\r\n", "\n").split("\n")
                marks = []
                for no, line in enumerate(lines, 1):
                    m = TOPLEVEL.match(line)
                    if not m:
                        continue
                    name = m.group("method") or m.group("func") or m.group("type") or m.group("vc")
                    if not name:
                        continue
                    symbols.setdefault(name, []).append((rel, no))
                    marks.append((no, name))
                enclosing[rel] = marks


scan_code()

files_by_base = {}
for name, locs in symbols.items():
    for rel, _no in locs:
        files_by_base.setdefault(os.path.basename(rel), set()).add(rel)
files_by_base = {k: sorted(v) for k, v in files_by_base.items()}


def body_at(rel, lineno):
    path = os.path.join(REPO, rel)
    with io.open(path, encoding="utf-8", newline="") as fh:
        lines = fh.read().replace("\r\n", "\n").split("\n")
    if lineno > len(lines):
        return None
    return lines[lineno - 1].strip()


def enclosing_symbol(rel, lineno):
    marks = enclosing.get(rel)
    if not marks:
        return None
    found = None
    for no, name in marks:
        if no <= lineno:
            found = name
        else:
            break
    return found


# ---------------------------------------------------------------------------
def load_regions():
    """Declared literal non-citations. Data, not judgement in the tool. `NOTACITE <literal> <why>`.

    Note what is NOT here: no line ranges, and no historical regions. Historical citations declare
    themselves inline with `@<sha>`, and quoted output declares itself with a fence tag, so neither
    needs an entry that could drift away from what it describes.
    """
    notacite, problems = {}, []
    if not os.path.exists(REGIONS):
        return notacite, ["cite-regions.txt ABSENT — no literal declarations loaded"]
    with io.open(REGIONS, encoding="utf-8") as fh:
        for no, raw in enumerate(fh, 1):
            line = raw.strip()
            if not line or line.startswith("#"):
                continue
            parts = line.split(None, 2)
            if parts[0] == "NOTACITE" and len(parts) >= 2:
                notacite[parts[1]] = parts[2] if len(parts) > 2 else ""
            else:
                problems.append("cite-regions.txt:%d unparsable: %s" % (no, line))
    return notacite, problems


NOTACITE, REGION_PROBLEMS = load_regions()

with io.open(PACK, encoding="utf-8", newline="") as fh:
    pack_lines = fh.read().split("\n")

# ---------------------------------------------------------------------------
# THE SUPERSET. Any colon followed by digits, anywhere. The preceding character class is NOT
# restricted, because restricting it is exactly how the last three blind spots were built.
CANDIDATE = re.compile(r"(?<![0-9]):(\d+)(?:\s*[-–]\s*(\d+))?")
# A historical citation declares its own SHA inline and is exempt BY CARRYING IT.
HISTORICAL = re.compile(r"`[A-Za-z0-9_./-]+\.[A-Za-z]+:\d+(?:\s*[-–]\s*\d+)?@[0-9a-f]{7,40}`")
FENCE = re.compile(r"^\s*```\s*(\S+)?")
# Anchors the prose names: backticked Go-ish identifiers. Tests are the load-bearing ones.
ANCHOR = re.compile(r"`([A-Z][A-Za-z0-9_]{3,}|[a-z][A-Za-z0-9_]{5,})`")
# REPORT ONLY, and it gates nothing. A backticked identifier carrying a call suffix — `Foo()` — is
# not an anchor to ANCHOR, which requires the whole span to be a bare identifier. The cold gate found
# `ProviderCapabilitySet()` sitting in the prose, real and invisible. R-24 forbids widening ANCHOR to
# swallow it, because a wider recognizer feeding a table that gates nothing buys nothing. It does
# require saying the span was seen and not resolved, which is what this pattern is for.
CALLSPAN = re.compile(r"`([A-Za-z_][A-Za-z0-9_]{2,})\(\)`")

MUTABLE = [
    (re.compile(r"\.\.HEAD\b"), "`..HEAD` — the range moves with every commit"),
    (re.compile(r"\bHEAD~\d*"), "`HEAD~` — relative to a moving ref"),
    (re.compile(r"\bgit\s+ls-files\b"), "`git ls-files` reads the INDEX and accepts no tree-ish"),
    (re.compile(r"\bgit\s+status\b"), "`git status` reads the working tree"),
]
SHAISH = re.compile(r"\b[0-9a-f]{7,40}\b")
GITREAD = re.compile(r"\bgit\s+(diff|grep|ls-tree|show|log)\b")

# Fence tags holding text the chip did not author and may not edit. Exempt from the mutable-axis
# pass only — a coordinate inside them is still a coordinate, and still reported.
FOREIGN = ("quoted-contract", "quoted-verdict")

violations = []
mutable_hits = []
anchors_seen = {}
callspans_seen = {}
exempted = {}

in_fence = None
for pack_no, line in enumerate(pack_lines, 1):
    fence = FENCE.match(line)
    if fence:
        in_fence = None if in_fence is not None else (fence.group(1) or "")
        continue

    # A `quoted-*` fence tag declares its contents as text this pack did not author: program output,
    # or a criterion quoted verbatim from the hub's contract. Neither is the pack making a claim in
    # its own voice, and neither is the chip's to edit. The tag travels WITH the content, which is the
    # whole point — the round-5 draft declared its exemptions as pack line ranges, and a line range
    # decays the moment a corrective inserts a paragraph above it, silently excusing the wrong lines.
    quoted = in_fence is not None and in_fence.startswith("quoted-")

    # --- job 3: mutable-axis commands. Checked inside quoted fences too, EXCEPT for FOREIGN TEXT —
    # words the chip neither wrote nor may edit. A criterion the hub authored is evidence of what the
    # hub required; a gate's verdict is evidence of what the gate said. When either one names a
    # command, the pack is not offering that command as its own evidence, it is recording that
    # someone else used or criticised it. Program output is NOT foreign in this sense and stays in
    # scope, because a pasted command IS the claim about how the pasted result was obtained.
    #
    # Widening an EXEMPTION is the dangerous direction — widening the violation recognizer is what
    # failed in rounds 3, 4 and 5, and an exemption fails the same way but quieter, since it removes
    # findings instead of adding them. So every exemption is counted and printed below. An exemption
    # that excuses nothing is inert; one that excuses a surprising number is visible before it is
    # load-bearing. The defect was never narrowness, it was silence.
    if in_fence in FOREIGN:
        exempted[in_fence] = exempted.get(in_fence, 0) + 1
    else:
        for rx, why in MUTABLE:
            if rx.search(line):
                mutable_hits.append((pack_no, why, line.strip()[:96]))
        if GITREAD.search(line) and not SHAISH.search(line) and "--cached" not in line:
            mutable_hits.append((pack_no, "a read command with no tree-ish on the line", line.strip()[:96]))

    # --- job 1 input: anchors the prose names
    if not quoted:
        for m in ANCHOR.finditer(line):
            anchors_seen.setdefault(m.group(1), []).append(pack_no)
        for m in CALLSPAN.finditer(line):
            callspans_seen.setdefault(m.group(1), []).append(pack_no)

    # --- job 2: no coordinate in prose
    if quoted:
        continue
    # Exemptions are SPANS, never bare tokens. A declared literal like `mm:500/1` must exempt the
    # `:500` INSIDE it and nothing else — comparing the match text (`:500`) against the declaration
    # would never fire, and comparing loosely would excuse every `:500` in the pack. Both failures
    # are the same one the recognizer rounds kept making: the thing checked was not the thing meant.
    exempt = [(m.start(), m.end()) for m in HISTORICAL.finditer(line)]
    for literal in NOTACITE:
        start = line.find(literal)
        while start != -1:
            exempt.append((start, start + len(literal)))
            start = line.find(literal, start + 1)
    for m in CANDIDATE.finditer(line):
        if any(lo <= m.start() < hi for lo, hi in exempt):
            continue
        violations.append((pack_no, m.group(0), line.strip()[:96]))

# ---------------------------------------------------------------------------
# job 1: the table
rows = []
unresolved = []      # recognized as an anchor, no declaration found in the scanned roots
for name in sorted(anchors_seen):
    locs = symbols.get(name)
    if not locs:
        # R-24 / R-17: the drop is recorded, not swallowed. `if not locs: continue` was the whole
        # defect — the name left the pipeline and the footer still printed a resolved count as if
        # the table covered everything the prose named.
        unresolved.append((name, anchors_seen[name]))
        continue
    if len(locs) > 1:
        rows.append((name, "AMBIGUOUS", "; ".join("%s:%d" % l for l in locs), ""))
        continue
    rel, no = locs[0]
    rows.append((name, rel, str(no), (body_at(rel, no) or "")[:78]))

table_lines = []
out = table_lines.append
out("=" * 112)
out("COORDINATE TABLE — generated from the CODE, keyed by the anchors EVIDENCE.md names (R-22).")
out("=" * 112)
out("")
out("GENERATED FILE. Do not hand-edit: the lane regenerates it and fails on the diff, so an edit here")
out("is reverted with a red lane rather than believed. Regenerate with `python cite-table.py`.")
out("")
out("  pack scanned    %s" % os.path.relpath(PACK, REPO).replace(os.sep, "/"))
out("  declarations    %s" % (os.path.relpath(REGIONS, REPO).replace(os.sep, "/")
                              if os.path.exists(REGIONS) else "cite-regions.txt  (ABSENT)"))
out("  code roots      %s" % "  ".join(os.path.relpath(r, REPO).replace(os.sep, "/") for r in ROOTS))
out("")
out("Coordinates live HERE and nowhere else. The prose names behaviour and test names; this table")
out("resolves them to positions. R-22 does not overturn the §11 judge — the anchor's TYPE is still a")
out("coordinate, and it is still mechanical; what is gone is the COPY of it into prose, which is the")
out("last place a decaying duplicate still lived.")
out("")
out("THIS TABLE IS NAVIGATION, NOT EVIDENCE, and that is R-24. It resolves THE ANCHORS LISTED BELOW.")
out("It does not resolve every anchor the prose names, and the ones it could not resolve are printed")
out("under UNRESOLVED rather than dropped. A wrong or missing row costs a reader a wrong jump; it")
out("cannot cost anyone a false proof, because nothing gates on this file. The claim it used to make")
out("— every anchor, resolved — was total in the wording and partial in the code, and R-24's remedy")
out("was to delete the claim rather than to widen the recognizer underneath it.")
out("")
out("AMBIGUOUS is a resolution, not a failure. A name carried by more than one declaration in the")
out("scanned roots resolves to ALL of them, printed. Narrowing it — by guessing the package the prose")
out("meant — would be an inference dressed as a lookup, which is the one thing this tool must never")
out("do (R-20: an inference must be visible as an inference).")
out("")
out("%-64s %-52s %-6s" % ("ANCHOR", "FILE", "LINE"))
out("-" * 112)
for name, rel, no, body in rows:
    out("%-64s %-52s %-6s" % (name, rel, no))
    if body:
        out("%-64s   %s" % ("", body))
out("")
out("  anchors resolved  %d" % len([r for r in rows if r[1] != "AMBIGUOUS"]))
out("  ambiguous         %d" % len([r for r in rows if r[1] == "AMBIGUOUS"]))
out("  UNRESOLVED        %d" % len(unresolved))
out("")
if unresolved:
    out("-" * 112)
    out("UNRESOLVED — recognized as an anchor in the prose, no declaration under the code roots above.")
    out("-" * 112)
    out("")
    out("Not a failure and not a pardon. A name lands here for ordinary reasons: it is declared inside")
    out("a grouped `var (…)` or `const (…)` block, which the column-0 matcher does not see; it lives")
    out("outside the scanned roots; or it is prose that merely looks like an identifier. R-24 forbids")
    out("widening the matcher to chase the first case, because a wider matcher feeding a table that")
    out("gates nothing buys nothing. What it requires is that the name appear here instead of leaving")
    out("the pipeline in silence. Resolve one by hand with `git grep -n '<name>'`.")
    out("")
    for name, at in unresolved:
        out("%-64s named on %d line%s" % (name, len(at), "" if len(at) == 1 else "s"))
    out("")

if callspans_seen:
    out("-" * 112)
    out("NOT TREATED AS ANCHORS — backticked spans carrying a call suffix, `Foo()`.")
    out("-" * 112)
    out("")
    out("The anchor pattern requires the whole backticked span to be a bare identifier, so these were")
    out("never anchors and never reached the table. The cold gate found one such span in the prose,")
    out("real in the code and invisible here. They are LISTED rather than RECOGNIZED, because R-24")
    out("forbids widening. Where the bare name does resolve, its position is shown — as a courtesy to")
    out("a reader, not as a row of the table above.")
    out("")
    for name in sorted(callspans_seen):
        locs = symbols.get(name) or []
        where = "; ".join("%s:%d" % l for l in locs) if locs else "(no declaration under the code roots)"
        out("%-40s %s" % (name + "()", where[:70]))
    out("")

out("-" * 112)
out("WHAT THIS TABLE DOES NOT COVER — R-24 requires the limit to travel with the artifact")
out("-" * 112)
out("")
out("Only backticked spans are considered at all, so an anchor named in bare prose is not seen. Only")
out("column-0 `func`/`type`/`var`/`const` declarations are indexed, so grouped-block declarations are")
out("absent. Only the code roots printed at the top are walked. This list of gaps is itself written")
out("by hand and is not guaranteed complete. None of that is a defect to be fixed by widening: it is")
out("the honest extent of a navigation aid, stated so that no reader mistakes its silence for a")
out("finding of absence.")
out("")

# Two modes, and the split is what makes drift a real failure rather than a sentence.
#
# Default = GENERATE: rewrite coordinates.txt from the code. That is the author's step.
# --strict  = CHECK  : write NOTHING, compare, and fail on any difference. That is the lane's step.
#
# A tool that regenerated in lane mode could never detect drift — it would overwrite the evidence of
# it and report success, which is the shape of every defect this chip has produced. Writing in the
# checker is the bug; the fix is that the checker cannot write.
TABLE = os.path.join(HERE, "coordinates.txt")
generated = "\n".join(table_lines)
committed = None
if os.path.exists(TABLE):
    with io.open(TABLE, encoding="utf-8", newline="") as fh:
        # Line endings are normalized before comparing, and this is not laxity. Git is configured to
        # convert on checkout here, so a fresh clone gets CRLF while this tool writes LF — the guard
        # would then be red on every machine but the author's, for a reason having nothing to do with
        # what it guards. A guard that cries wolf is a guard someone turns off, and this pack cannot
        # afford another check that is technically present and practically ignored. Round 1 of this
        # chip lost a real gofmt violation inside exactly this noise.
        committed = fh.read().replace("\r\n", "\n")
drift = committed is not None and committed != generated
missing = committed is None

if "--strict" in sys.argv:
    if missing:
        print("coordinate table MISSING: %s" % os.path.relpath(TABLE, REPO).replace(os.sep, "/"))
    elif drift:
        print("coordinate table DRIFTED: %s" % os.path.relpath(TABLE, REPO).replace(os.sep, "/"))
        print("  the committed table is not what the code generates. Regenerate with")
        print("  `python cite-table.py` and commit the result; do not hand-edit it.")
        import difflib
        d = list(difflib.unified_diff(committed.split("\n"), generated.split("\n"),
                                      "committed", "generated", lineterm="", n=1))
        for line in d[:40]:
            print("    %s" % line)
        if len(d) > 40:
            print("    … %d more diff lines" % (len(d) - 40))
    else:
        print("coordinate table CURRENT: %s" % os.path.relpath(TABLE, REPO).replace(os.sep, "/"))
else:
    with io.open(TABLE, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(generated)
    print("coordinate table %s -> %s   (%d resolved, %d ambiguous)"
          % ("REGENERATED" if drift or missing else "unchanged",
             os.path.relpath(TABLE, REPO).replace(os.sep, "/"),
             len([r for r in rows if r[1] != "AMBIGUOUS"]),
             len([r for r in rows if r[1] == "AMBIGUOUS"])))
if "--show" in sys.argv:
    print(generated)
print("")

print("=" * 112)
print("REPORT — pack hygiene. NOT A LADDER RUNG (R-24). --strict still exits 1 for the author.")
print("=" * 112)
print("")
print("  coordinates found in prose         %d   (target 0 — R-22)" % len(violations))
print("  mutable-axis commands found        %d   (target 0)" % len(mutable_hits))
print("  anchors this tool could not resolve %d  (a report, never an exit code — R-24)" % len(unresolved))
print("  backticked call spans not anchored %d   (listed in the table, not recognized)"
      % len(callspans_seen))
if REGION_PROBLEMS:
    for p in REGION_PROBLEMS:
        print("  DECLARATION: %s" % p)
print("")
print("  EXEMPTIONS — lines the mutable-axis pass did not check, because they are foreign text the")
print("  chip may not edit. Printed so no exemption can be silent; an unused one is inert.")
for tag in FOREIGN:
    print("    ```%-18s %4d lines" % (tag, exempted.get(tag, 0)))
print("")
print("WHAT THESE TWO PASSES SEE, AND WHAT THEY DO NOT. R-24 struck the sentence that used to stand")
print("here, which called the prose scan a SUPERSET; it was total in the wording and partial in the")
print("code, and that gap was a round-6 blocking finding.")
print("")
print("  coordinates   sees `:<digits>` NOT preceded by another digit, anywhere on a line, inside or")
print("                outside backticks and fences. Does NOT see `123:456`, where the left neighbour")
print("                is a digit. Exempt: an inline SHA, a declared literal, or ANY ```quoted-* fence")
print("                — every such tag, not only ```quoted-output. Each exemption travels with the")
print("                content it describes, so none can drift away from what it excuses.")
print("  mutable axis  sees `..HEAD`, `HEAD~`, `git ls-files`, `git status`, and a bare")
print("                `git diff|grep|ls-tree|show|log` on a line carrying no hex-ish token. Does NOT")
print("                see bare `HEAD`, `git rev-parse`, a mutable named ref such as `main`, a plain")
print("                worktree `grep`, or a line where an unrelated SHA satisfies the hex-ish test.")
print("")
print("The RULE is whole and is not what narrowed: every evidence command names its own tree-ish or")
print("it is not evidence. These passes are a partial reading of that rule, and the reading is now")
print("written down instead of being implied to be complete.")
print("")

if violations:
    print("-" * 112)
    print("COORDINATES STILL IN PROSE")
    print("-" * 112)
    for pack_no, text, ctx in violations:
        print("pack:%-5d %-12s | %s" % (pack_no, text, ctx))
    print("")

if mutable_hits:
    print("-" * 112)
    print("MUTABLE-AXIS COMMANDS — every evidence command names its own tree-ish or it is not evidence")
    print("-" * 112)
    for pack_no, why, ctx in mutable_hits:
        print("pack:%-5d %s" % (pack_no, why))
        print("            %s" % ctx)
    print("")

if "--suggest" in sys.argv:
    print("-" * 112)
    print("SUGGESTIONS — the enclosing declaration at each coordinate still in prose. Replace the")
    print("coordinate with a behavioural claim naming this anchor; the table above resolves it back.")
    print("-" * 112)
    CO = re.compile(r"([A-Za-z0-9_./-]+\.go):(\d+)")
    for pack_no, _text, ctx in violations:
        for m in CO.finditer(pack_lines[pack_no - 1]):
            base, no = os.path.basename(m.group(1)), int(m.group(2))
            for rel in files_by_base.get(base, []):
                sym = enclosing_symbol(rel, no)
                print("pack:%-5d %s:%d  ->  %s" % (pack_no, rel, no, sym or "(no enclosing decl)"))

if "--strict" in sys.argv and (violations or mutable_hits or drift or missing):
    sys.exit(1)
