"""SUPERSEDED — round-4 apparatus, kept as the record of what round 4 found.

Do not run this to check the pack. The live check is `cite-table.py --strict` (ladder rung PACK), and
the coordinate table it generates is `coordinates.txt`. This script's coverage claim was itself found
false by the round-5 derivation gate: its total counted unresolved rows as resolved, and whole
citation forms were outside its regex. That finding is why R-22 exists, so the file is evidence and
not a tool.
"""

"""Mechanical citation audit for EVIDENCE.md — committed pack artifact (hub ruling R-11).

Why this exists. Three consecutive P6 rounds on CHIP-ANCHORS failed on the same class: a `file:line`
citation that was correct when typed and wrong by the time a reviewer read it, because a corrective
inserted lines above it. Hand-fixing the instances a gate happened to name is what produced rounds 3
and 4. The D21 independent judgement REJECTED the chip's own proposal (cite by symbol) and prescribed
mechanical regeneration at the final tip; R-11 then made this a committed artifact, re-run at the
frozen tip, with the gate checking it ran clean instead of sampling citations by hand — the cold gate
has no Bash and is structurally blind to this class.

What it proves and what it does NOT. It resolves every `file:line` / `file:line-line` citation in the
pack against the working tree and prints what the cited line ACTUALLY contains, so a human decides
only the ambiguous ones. It does not prove the PROSE about that line is true: a claim that "this
branch has no test" resolves perfectly while being false. That half is the cold reviewer's job.

Read-only. Writes nothing into the worktree. Run from anywhere; paths resolve from this file.

    python cite-audit.py > cite-audit.txt
"""
import io
import os
import re
import sys

# Windows stdout defaults to cp1252, which turns every accented character in the cited Go source
# into a replacement char in the committed artifact. The artifact is evidence; it has to be faithful.
try:
    sys.stdout.reconfigure(encoding="utf-8")
except AttributeError:  # pragma: no cover - Python < 3.7
    pass

HERE = os.path.dirname(os.path.abspath(__file__))
# .mnfs/MIS-006-integracao-fundacao/_chip-anchors/ -> repo root
REPO = os.path.abspath(os.path.join(HERE, "..", "..", ".."))
PACK = os.path.join(HERE, "EVIDENCE.md")

# Where a bare `foo.go:123` citation lives, by basename. Ambiguous basenames resolve to a list.
ROOTS = [
    os.path.join(REPO, "apps", "server_core", "internal"),
]

index = {}
for root in ROOTS:
    for dirpath, _dirnames, filenames in os.walk(root):
        for fn in filenames:
            if fn.endswith(".go"):
                index.setdefault(fn, []).append(os.path.join(dirpath, fn))

CITE = re.compile(r"`([A-Za-z0-9_./-]+\.(?:go|yaml|md|ts))(?::(\d+)(?:\s*[-–]\s*(\d+))?)?`")

with io.open(PACK, encoding="utf-8", newline="") as fh:
    pack_lines = fh.read().split("\n")

cache = {}


def read_file(path):
    if path not in cache:
        with io.open(path, encoding="utf-8", newline="") as fh:
            cache[path] = fh.read().replace("\r\n", "\n").split("\n")
    return cache[path]


rows = []
for pack_no, line in enumerate(pack_lines, 1):
    for m in CITE.finditer(line):
        target, start, end = m.group(1), m.group(2), m.group(3)
        if start is None:
            continue
        base = os.path.basename(target)
        candidates = index.get(base, [])
        if len(candidates) != 1:
            rows.append((pack_no, m.group(0), "AMBIGUOUS-OR-UNKNOWN", "%d candidates" % len(candidates)))
            continue
        src = read_file(candidates[0])
        n = int(start)
        if n > len(src):
            rows.append((pack_no, m.group(0), "OUT-OF-RANGE", "file has %d lines" % len(src)))
            continue
        body = src[n - 1].strip()
        tail = ""
        if end:
            e = int(end)
            if e <= len(src):
                tail = "  ..  ->%d: %s" % (e, src[e - 1].strip()[:60])
        rows.append((pack_no, m.group(0), body[:90], tail))

print("EVIDENCE.md citations resolved against the code tip: %d" % len(rows))
print("")
for pack_no, cite, body, tail in rows:
    print("pack:%-5d %-58s | %s%s" % (pack_no, cite, body, tail))

# ---------------------------------------------------------------------------
# SECOND PASS — bare `:NNN` citations.
#
# The round-4 cold gate found the first pass structurally blind to these: CITE requires a
# filename with an extension, so a bare `:640` is invisible to it, and the pack's closure
# sentence nonetheless claimed EVERY file:line-shaped citation had been re-derived. Two
# survivors were demonstrated by hand. This pass exists so the claim can be made honestly.
#
# A bare citation carries no file — it inherits one from the prose above it. That inheritance
# is an INFERENCE, and a tool that resolves it silently would manufacture exactly the false
# confidence this artifact was built to remove. So the inferred file is PRINTED on every row,
# together with the distance in pack lines back to where it was inherited from. A wrong
# inference is then visible: the resolved line will not match the sentence around it.
#
# Rows are flagged, never silently dropped: NO-CONTEXT when no file-qualified citation
# precedes it within the window, OUT-OF-RANGE when the line does not exist in the inferred
# file (this is also what catches non-citations that merely look like one, e.g. a port `:8080`).
#
# COMMA LISTS. The first version of this pass matched a whole backtick span (`` `:640` ``), so a span
# holding SEVERAL citations — `` `:527,:558,:582,:653` `` — matched nothing and passed in silence.
# Four such spans sat in the decay table and survived a run that reported zero bare citations. That is
# the same defect as the original narrow regex, one level down, and it is why this pass now scans
# INSIDE each span for every `:NNN` rather than requiring the span to be exactly one citation.
INHERIT_WINDOW = 20
SPAN = re.compile(r"`[^`]*`")
BARE_IN_SPAN = re.compile(r"(?<![A-Za-z0-9_])[:](\d+)(?:\s*[-–]\s*(\d+))?")


def bare_matches(line):
    """Yield (abs_start, abs_end, lineno, endno) for each bare citation inside a backtick span.

    A citation is bare when nothing filename-shaped precedes the colon inside its own span. The
    span is the unit because that is how a reader sees it: `` `foo.go:12` `` carries its file,
    `` `:12,:20` `` carries none, and both live inside one pair of backticks.
    """
    for span in SPAN.finditer(line):
        body = span.group(0)
        for m in BARE_IN_SPAN.finditer(body):
            before = body[: m.start()]
            # File-qualified iff a filename with an extension ends right where the colon starts.
            if re.search(r"[A-Za-z0-9_./-]+\.(?:go|yaml|md|ts)$", before):
                continue
            yield (span.start() + m.start(), span.start() + m.end(), m.group(1), m.group(2))

bare_rows = []
last_target = None
last_no = None
for pack_no, line in enumerate(pack_lines, 1):
    events = []
    for m in CITE.finditer(line):
        events.append((m.start(), "file", m.group(1)))
    for start, _stop, start_no, end_no in bare_matches(line):
        events.append((start, "bare", (start_no, end_no)))
    for _pos, kind, payload in sorted(events, key=lambda e: e[0]):
        if kind == "file":
            last_target, last_no = payload, pack_no
            continue
        start_no, end_no = payload
        cite = ":%s" % start_no if not end_no else ":%s-%s" % (start_no, end_no)
        distance = None if last_no is None else pack_no - last_no
        if last_target is None or distance > INHERIT_WINDOW:
            bare_rows.append((pack_no, cite, "-", "NO-CONTEXT",
                              "no file-qualified citation within %d lines" % INHERIT_WINDOW))
            continue
        base = os.path.basename(last_target)
        candidates = index.get(base, [])
        if len(candidates) != 1:
            bare_rows.append((pack_no, cite, base, "AMBIGUOUS-OR-UNKNOWN",
                              "%d candidates" % len(candidates)))
            continue
        src = read_file(candidates[0])
        n = int(start_no)
        if n > len(src):
            bare_rows.append((pack_no, cite, base, "OUT-OF-RANGE",
                              "file has %d lines" % len(src)))
            continue
        body = src[n - 1].strip()[:70]
        if end_no and int(end_no) <= len(src):
            body += "  ..  ->%s: %s" % (end_no, src[int(end_no) - 1].strip()[:40])
        bare_rows.append((pack_no, cite, base, "(-%d)" % distance, body))

print("")
print("=" * 100)
print("COVERAGE — this artifact declares its own reach (hub ruling R-17).")
print("")
print("  resolved (file-qualified citations, checked against the code tip) : %d" % len(rows))
print("  unresolvable (bare citation, no file)                             : %d" % len(bare_rows))
print("")
print("R-17 BANS bare `:NNN` citations from the pack. Every citation carries its file, so coverage")
print("is total BY CONSTRUCTION rather than by a better regex. The count above is therefore the")
print("enforcement point, not a statistic: a NON-ZERO unresolvable count IS the violation (R-20).")
print("")
print("Why this block exists at all. Until round 4 the regex silently skipped bare citations, so the")
print("pack's closure sentence claimed 'every file:line citation re-derived' while 89 were invisible")
print("to the tool that was supposed to have checked them. Two real survivors were found by hand.")
print("A tool that skips in silence makes its own closure sentence overclaim by construction — the")
print("defect was never the narrowness, it was the silence.")
print("")
print("And the SAME defect recurred one level down, inside the fix for it. The first bare pass matched")
print("a whole backtick span, so a span holding several citations — `:527,:558,:582,:653` — matched")
print("nothing and the run reported ZERO bare citations while 61 were still there, four of them in the")
print("decay table. The pass now scans INSIDE each span. The lesson is the reason this block is printed")
print("at all: a coverage number is only worth reading next to the rule that produced it, because the")
print("number cannot tell you what the rule cannot see.")
print("")
print("What this artifact still does NOT decide: whether the PROSE around a citation is true. A")
print("citation resolving to plausible-but-wrong content passes here (R-21), which is why every")
print("resolved line is printed above — reading compares a sentence to an expectation, never to the")
print("actual line, unless the actual line is on the page.")
print("=" * 100)
if bare_rows:
    print("")
    print("VIOLATIONS — bare citations remaining. File INFERRED from the nearest preceding")
    print("file-qualified citation within %d pack lines. The inference is PRINTED, not assumed:" % INHERIT_WINDOW)
    print("an inference must be visible as an inference (R-20). Qualify each one at its site.")
    print("")
    for pack_no, cite, base, note, body in bare_rows:
        print("pack:%-5d %-14s %-38s %-22s | %s" % (pack_no, cite, base, note, body))
