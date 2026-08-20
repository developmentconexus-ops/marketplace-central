# AI Dialog — review-branch gate proof

> **TEMPORARY / NON-AUTHORITATIVE.** This file exists only to prove that the accepted review-branch lifecycle can run repository gates while the merge candidate remains free of review dialogue.

```text
candidate head   2ec3b8482d3dfba77869ddb1705168bc01976c10
review branch    review/authority-surface-cleanup-gate-proof
candidate PR     #45
```

The review branch must differ from the candidate only by this file. The documentation-authority guard must inspect the declared candidate tree for merge contamination, while still requiring this ambient branch to be named `review/*`.

Delete or retire this temporary branch after the proof is recorded. No architecture decision or merge authorization is carried here.
