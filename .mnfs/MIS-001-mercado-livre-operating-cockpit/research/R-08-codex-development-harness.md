# Research Note

```yaml
id: R-08
type: research
status: complete
owner: Mission Strategist
parent: MIS-001
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: support
```

## Topic

Codex-native architecture for long-running, token-efficient Marketplace Central
development without local clean-machine simulation.

## Sources Checked

- [Codex subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents) — parallel agent behavior, context pollution, custom agents, model routing, depth, and sandbox inheritance.
- [Codex customization](https://learn.chatgpt.com/docs/customization/overview) — `AGENTS.md`, repo skills, plugins, MCP, and progressive disclosure.
- [Codex worktrees](https://learn.chatgpt.com/docs/environments/git-worktrees) — task/worktree ownership, handoff, ignored files, and lifecycle.
- [Codex app-server API](https://learn.chatgpt.com/docs/app-server#api-overview) — thread start/read/steer/interrupt/goal methods and experimental boundaries.
- [Codex developer commands](https://learn.chatgpt.com/docs/developer-commands#codex-exec) — structured non-interactive sessions, JSONL, output schema, resume, and sandboxing.
- [Lost in the Middle](https://arxiv.org/abs/2307.03172) — retrieval quality degradation in long contexts.
- [Effective context engineering for AI agents](https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents) — bounded context and progressive disclosure.
- [Effective harnesses for long-running agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents) — durable progress artifacts and incremental execution.
- [CodePlan](https://www.microsoft.com/en-us/research/publication/codeplan-repository-level-coding-using-llms-and-planning-2/) — repository-level hierarchical planning evidence.
- [SWE-bench](https://proceedings.iclr.cc/paper_files/paper/2024/hash/edac78c3e300629acfe6cbe9ca88fb84-Abstract-Conference.html) — executable, base-commit-bound evaluation.

The Codex manual helper was attempted first. It reached the official endpoint
but rejected the response because the integrity header `x-content-sha256` was
missing. Official OpenAI documentation MCP pages were used instead.

## Findings

- Subagents reduce main-thread context pollution when work is independent and
  returns a compact summary; they consume additional tokens and are not a
  default justification for parallel writers.
- Current Codex supports visible tasks, subagent threads, task steering,
  worktrees/handoff, goals, repo skills, custom agents, hooks, browser tooling,
  and structured `codex exec` runs.
- The local CLI exposes `app-server` as experimental. A custom orchestration
  client would add an unnecessary runtime dependency to V1.
- Skills provide progressive disclosure; `AGENTS.md` is better suited to short
  durable policy than long workflow prose.
- Worktrees separate Git ownership, not every runtime resource. Separate ports,
  databases, or Compose names are needed only when concurrently active work
  would collide.
- Long-context research and repository-agent evidence favor hierarchical plans,
  fresh bounded workers, durable artifacts, executable graders, and explicit
  risk gates.
- The existing Marketplace Central governance/context work already implements
  much of the knowledge plane. Its main missing pieces are knowledge routing,
  state/lease/resume, native Codex orchestration, risk-selected impact gates,
  and dogfood evals.
- The local cold clone/provisioning experiment adds no product-facing safety
  beyond the accepted unit environment and ephemeral PostgreSQL lanes and has
  repeatedly blocked development on host-specific Git behavior.

## Recommendation

Keep the repository-native knowledge, context, environment, PostgreSQL, and
governance foundations. Retire cold clone/provisioning. Build the remaining V1
as a repo skill plus deterministic state/context/gate modules using native Codex
tasks, depth-one subagents, optional worktrees, and risk-triggered real QA.

## Impact On Mission

- M-08 validation no longer requires two cold runs or clean dependency caches.
- F-04 remains inspectable superseded evidence.
- A cutover feature removes cold-specific runtime/code and introduces the
  current-checkout impact gate.
- F-05 becomes the goal/task/context/lease orchestration control plane.
- F-09 becomes deterministic eval plus fresh-task dogfood and efficiency proof.

## Handoff

- Current status: Research complete and incorporated into M-08 replan.
- Next owner: Milestone Orchestrator after planning artifacts are accepted.
- Next action: Execute pragmatic cutover, then orchestration control plane, then dogfood.
- Required files/evidence: M-08 milestone, validation contract, execution guide, F-10/F-05/F-09 briefs.
- Blockers or open decisions: None; operator explicitly rejected VM-like cold execution.
