# Claude Code → Codex Portfolio bridge

This local MCP server lets Claude Code send one advisory message to an existing
Codex Portfolio task and receive its final response. It does not create a task,
open a network listener, or allow Claude to select arbitrary task IDs. The
resumed Codex turn runs read-only with approvals disabled.

## Register once

From the repository root in PowerShell:

```powershell
node scripts/claude-codex-bridge.mjs --task-id 019f5cf6-8c9f-7321-ba07-f5b5b5e6bc77 --check
claude mcp add-json --scope local codex-portfolio '{"type":"stdio","command":"node","args":["C:\\Users\\leandro.theodoro\\Documents\\marketplace-central\\scripts\\claude-codex-bridge.mjs"],"env":{"CODEX_PORTFOLIO_TASK_ID":"019f5cf6-8c9f-7321-ba07-f5b5b5e6bc77"}}'
claude mcp get codex-portfolio
```

The local scope is private to this checkout and requires no repository config.

## Use

In Claude Code, ask it explicitly to call the tool:

```text
Use ask_codex_portfolio to send this proposal to the Portfolio Hub and return
its response before continuing: <proposal>
```

The tool waits for the Codex turn to finish. Do not call it while the same Codex
task is already running, and keep only one bridge call active at a time.

## Remove or retarget

```powershell
claude mcp remove codex-portfolio --scope local
```

Then run the registration command again with the new task ID.
