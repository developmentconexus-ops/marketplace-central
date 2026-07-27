# HUB GATE — CHIP-ANCHORS (`chip/anchors`)

**Hub-owned pack.** The chip did not self-stamp the marker below and was right not to: the
`P6-DUAL-GATE:` line is the hub's authority, not the chip's. The chip's own closing sentence reads
`AGREEMENT — P6 discharged` and carries the discharge ledger beside it
(`_chip-anchors/EVIDENCE.md`, section *P6 DISCHARGE*). This file is the hub's independent record and
the one the merge hook reads, because the hook greps the working tree of `main` before the chip's
branch arrives.

- Branch: `chip/anchors`
- Final tip: `cbe045bb7dcf10303714e6c8014041aa6e589529`
- Frozen CODE tip: `85b6c367` — every commit after it is pack-only, verified below
- Scope: F-01 anchor-capability-as-data · F-02 dimension units `big.Rat` · F-03 independent limits
- Rulings R-1…R-26: `_chip-anchors/hub-rulings.md`

---

P6-DUAL-GATE: AGREEMENT

LIVE-VERIFIED: U1-U3 driven by the HUB against the connected ML account (R-8) — evidence
`.mnfs/MIS-006-integracao-fundacao/_chip-anchors/hub-live-drive-u1-u3.md`, committed to `main`
@ `40623b57`, U2 inversion corrected by a NEW commit @ `874d00e5` (a pin names its SHA; the answer
to a wrong pinned artifact is a new commit, never an edit that makes the pin lie).

---

## What AGREEMENT means here, and what it does not

The two gates agreed on the CODE, at `85b6c367`. The cold Opus gate returned PASS on semantics
(C1–C4, C6–C8, read from source). The GPT-5.6 Sol gate returned PASS on its own scope
(C2/C3/C5/C9/C10/C12, re-derived from git). Neither verdict on the code depended on the chip's
`cite-table.py` — the observation that ended round 6 and produced R-24.

Rounds 4–6 were about the PACK, not the code. Their five DERIVATION blockers are discharged:

| # | Finding | Closed by | Means |
|---|---|---|---|
| 1 | `connectorsports` sentence claimed an empty result without the exclusion that makes it empty | GPT gate, at the corrected tip | Both grep forms re-derived at the frozen tip and at the control |
| 2 | S1+S2 write set ran two commits past S2 | GPT gate, at the corrected tip | Range re-pinned, `SETS_EQUAL` confirmed, coincidence stated rather than the SHA quietly swapped |
| 3 | Tool silently discarded unresolvable anchors | GPT gate, which had filed it BLOCKING | DISSOLVED-BY-R-24, accepted as closure by its own filer |
| 4 | Mutable-axis pass called itself a ban | **HUB OBSERVATION, not a gate** | Sentence re-aimed at what the report PRINTS — a gate-certified target — absence then verified by string |
| 5 | Prose scan called itself a SUPERSET | **HUB OBSERVATION, not a gate** | Comment DELETED, not reworded; absence verified by string |

**Honest limit, stated and not softened.** GPT did not re-run and does not certify the final tip. The
hub holds (a) its verdict that those two lines were the ENTIRE remaining defect and that the corrective
introduced none, and (b) the hub's own observation of absence. **The hub does NOT hold a model-side
verdict on the final tip.** The gap is accepted deliberately under R-25(d): when a remedy is fully
determined by its finding — two files, exact strings, a replacement target already certified — the
verification is a check of absence, one bit, and that wants a shell rather than a second model. A
round whose entire content is confirming three strings are gone would itself need another round if it
found a typo, and that regress is what R-24 ended.

A closing sentence that omits HOW the close was reached is round 4's silent-omission class committed
in the last line of the document. Hence the ledger above the marker rather than the marker alone.

## Merge conditions — all three run by the HUB at the final tip

Verified by STRING, never by line: a gate's coordinate is a claim like any other, and the round-6
gate cited `EVIDENCE.md:1028` where the text actually lived at `:1030` and `:1057`.

```
$ git diff --name-only 85b6c367 cbe045bb -- apps contracts packages
(no output)                                                          COUNT=0

$ git grep -n -E "Banned outright|lane check now|preceding character class is NOT" \
      cbe045bb -- _chip-anchors/EVIDENCE.md _chip-anchors/cite-table.py
(no output)                                                          EXIT=1  (absent, expected)

$ python cite-table.py --strict          # run in the chip worktree at cbe045bb
...
REPORT — pack hygiene. NOT A LADDER RUNG (R-24).                     EXIT=0
```

**One string class survives the widened grep on purpose.** The same false sentences still appear at
`_chip-anchors/verdict-r6-recheck.md:48,49`, inside the ` ```quoted-verdict ` fence, because that is
the GPT gate's own text transcribed under R-23 and no party may edit it. A false claim quoted as the
report of the gate that found it is not the pack asserting it. Scope condition 2 to the two
claim-holders or expect the count of 2.

## Hub-side backlog carried out of this chip (none blocking)

- FOREIGN-fence provenance hole: a pack author can still exempt authored text by tagging it
  `quoted-verdict`; the tool checks the tag, never the provenance. SHOULD-FIX in an artifact that no
  longer gates anything — the same reason three other blockers dissolved. **Carried with its
  provenance, because the provenance is the argument:** the hole landed as a real defect in the HUB's
  own copies on the same day it sat backlogged as a theoretical defect in the chip's tool. It cost an
  honest party who was not trying to exploit it, which is the objection *"nobody would actually do
  this"* answered by observation rather than by argument.
- `refforn` does not belong in the cross-side anchor vocabulary.
- Declared-but-incomparable anchors vanish silently (`generation_service.go` ~`:639`).
- `slices.SortFunc` + `CompactFunc` display nondeterminism at `:764-769`.
- 53 standing governance violations on `main`, pre-existing and unrelated to this chip.
