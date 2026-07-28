"""Measure the pack at a NAMED COMMIT, off the blobs, never off the working tree.

Why the blob: a figure measured on disk and then written to disk describes a file that was
never committed. Measuring an immutable commit removes the self-reference entirely.

Units are named on every row:
  lines      = LF-terminated lines (a trailing newline does NOT create a final empty line)
  sentences  = split on [.!?:;] + whitespace, and on newline; empty fragments dropped
  chars      = len() of the UTF-8-decoded text
  bytes      = len() of the raw blob
"""
import re
import subprocess
import sys

COMMIT = sys.argv[1]
REPO = r"C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/happy-montalcini-b010c0"
PREFIX = ".mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/"
FILES = ["EVIDENCE.md", "dispatches/p6-gate-brief.md", "dispatches/p6-gate-brief-r2.md",
         "dispatches/p6-gate-brief-r3.md", "dispatches/a2-assertion-before-after.md"]

P1 = ("nunca|sempre|apenas|inalcanç|não pode|impossív|todo |toda |todos|todas|"
      "nenhum|qualquer|única|único|só o|só a|never|always|only|unreachable|cannot|every|no longer")
cs = re.compile("(" + P1 + ")")
ci = re.compile("(" + P1 + ")", re.I)
SENT = re.compile(r"(?<=[.!?:;])\s+|\n")


def blob(path):
    return subprocess.run(["git", "-C", REPO, "cat-file", "-p", f"{COMMIT}:{PREFIX}{path}"],
                          capture_output=True, check=True).stdout


rows = []
tot = [0] * 7
for rel in FILES:
    raw = blob(rel)
    txt = raw.decode("utf-8")
    lines = txt.split("\n")
    if lines and lines[-1] == "":
        lines = lines[:-1]
    sents = [s for s in SENT.split(txt) if s.strip()]
    row = [rel, len(lines),
           sum(1 for l in lines if cs.search(l)), sum(1 for l in lines if ci.search(l)),
           len(sents),
           sum(1 for s in sents if cs.search(s)), sum(1 for s in sents if ci.search(s)),
           len(cs.findall(txt)), len(ci.findall(txt)), len(txt), len(raw)]
    rows.append(row)

print(f"commit {COMMIT}")
print("file | lines | cs-lines | ci-lines | sentences | cs-sent | ci-sent | cs-occ | ci-occ | chars | bytes")
agg = [0] * 10
for r in rows:
    print(" | ".join(str(x) for x in r))
    for i in range(1, 11):
        agg[i - 1] += r[i]
print("TOTAL | " + " | ".join(str(x) for x in agg))
