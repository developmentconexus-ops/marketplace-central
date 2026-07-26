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
