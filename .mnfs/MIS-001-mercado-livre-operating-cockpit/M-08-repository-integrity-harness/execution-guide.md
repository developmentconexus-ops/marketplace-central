# M-08 Execution Guide — Simple Development Sessions

Portfolio starts one visible Milestone session using the packet in the repo
skill. The packet includes Portfolio task ID, objective, base SHA, mission and
milestone files, knowledge routes, constraints, and QA contract. Milestone
acknowledges it and reports compact checkpoints back to Portfolio.

Milestone orders Features and dispatches one bounded Feature worker per eligible
slice. Feature Plan and Feature Execution stay in the same session by default.
After `plan.md`, compile and validate the context file; a new worker receives
the file path, knowledge route IDs, paths, seams, side effects, proof commands,
and stop conditions. Feature returns one commit and impacted evidence.

One writer owns a seam. Parallelize only disjoint Features in worktrees;
investigation may be read-only. Prompts pass paths/selectors, not copied docs or
transcripts. A named route gap permits targeted discovery.

After all Feature commits integrate, Milestone freezes one SHA and requests one
final review plus the QA targets required by the milestone contract. Low-risk
work receives focused checks; real Oracle/provider/database/browser QA runs only
when the contract names it. Only QA passes the milestone.

F-10 and F-05 are the accepted implementation. F-09 remains preserved WIP but
is rejected from the active workflow and is not a completion requirement.
Never restore cold execution or add a custom agent runtime, second CI, hooks,
app server, or token benchmark.

```text
status:
checkpoint_id:
commit:
changed_paths:
evidence:
review:
blockers:
next:
```
