# CHIP-SAT dispatch ledger — MIS-003 W1 (M-05 F-01 + M-06 F-02)

Branch `chip/sat-m05f01-m06f02` · base `a49168e641ffd6f61932ca57c29b1d1bdcde2fb0` · worktree `../marketplace-central-chip-sat`

| # | Role | Model/effort | Path | Prompt | Log / output | Status | Verdict/artifact |
|---|---|---|---|---|---|---|---|
| 1 | Feature planner (P2 batch, both features) | gpt-5.6-sol / medium | OS-process | scratchpad/prompt-plan-sat.md | scratchpad/agent__plan-sat.log · agent__plan-sat.last.md | done | CHIP-SAT-P2-PLAN.md (23 slices; 1 blocking REQUEST: orders canonical fields) |
| 2 | Implement F01-O1 orders query/cursor grammar | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-o1.md | scratchpad/agent__f01-o1.log · agent__f01-o1.last.md | dispatched | pending |
| 3 | Implement F01-S1 sync-runs cursor+repo (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-f01-s1.md | scratchpad/agent__f01-s1.log · agent__f01-s1.last.md | done | quick_validation_passed; 6 files, write-set respected |
| 5 | Slice review F01-S1 | claude sonnet subagent | Agent SYNC | scratchpad/review-pack-f01-s1.md | scratchpad/review-verdict-f01-s1.md | done | REJECT — 1 important (StartedAt nullable), fix dispatched row 8 |
| 6 | Implement F01-O1 orders grammar | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-o1.md | scratchpad/agent__f01-o1.log · agent__f01-o1.last.md | done | green in write-set; ruling delta (fulfillment) folded into row 8 |
| 7 | Implement F01-D2 link counters | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d2.md | scratchpad/agent__f01-d2.log · agent__f01-d2.last.md | done | green in write-set |
| 8 | Corrective: S1 StartedAt + O1 unsupported_filter | gpt-5.6-luna / high | OS-process | scratchpad/prompt-fix-s1-o1.md | scratchpad/agent__fix-s1-o1.log · agent__fix-s1-o1.last.md | dispatched | pending |
| 9 | Slice review F01-D2 | claude sonnet subagent | Agent SYNC | scratchpad/review-pack-f01-d2.md | scratchpad/review-verdict-f01-d2.md | done | ACCEPT — committed bee6df3 |
| 10 | Delta re-review F01-S1 | claude sonnet subagent (resumed row 5) | SendMessage | (delta instruction) | scratchpad/review-verdict-f01-s1.md (appended) | done | DELTA ACCEPT — committed 1fc1253 |
| 11 | Slice review F01-O1 (incl. corrective) | claude sonnet subagent | Agent SYNC | scratchpad/review-pack-f01-o1.md | scratchpad/review-verdict-f01-o1.md | done | ACCEPT-WITH-CONDITIONS (3 test gaps + 2 nits → folded into O2) — committed 1fc1253 |
| 12 | Implement F01-O2 orders keyset repo (complex) + O1 conditions | gpt-5.6-sol / low | OS-process | scratchpad/prompt-f01-o2.md | scratchpad/agent__f01-o2.log · agent__f01-o2.last.md | done | quick_validation_passed; 5 files in write-set; repo-wide vet green |
| 13 | Implement F01-S2 sync-runs app+transport | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-s2.md | scratchpad/agent__f01-s2.log · agent__f01-s2.last.md | done | quick_validation_passed; 6 files in write-set; vet blocker was O2 in-flight, resolved |
| 14 | Slice review F01-O2 | claude sonnet subagent | Agent | scratchpad/review-pack-f01-o2.md | scratchpad/review-verdict-f01-o2.md | done | REJECT — 1 blocking (nf_state missing evidence_state='exact' gate), fix row 16 |
| 16 | Corrective: O2 nf_state exact-evidence gate | gpt-5.6-luna / high | OS-process | scratchpad/prompt-fix-o2.md | scratchpad/agent__fix-o2.log · agent__fix-o2.last.md | done | evidence_state='exact' on both subqueries + unknown→nil regression test; tests PASS |
| 18 | Delta re-review F01-O2 | claude sonnet subagent (resumed row 14) | SendMessage | (delta instruction) | scratchpad/review-verdict-f01-o2.md (appended) | done | DELTA ACCEPT (1 retry — server error mid-append); suggestion+nit carried non-blocking |
| 20 | Implement F01-D1 orders dashboard counters | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d1.md | scratchpad/agent__f01-d1.log · agent__f01-d1.last.md | done | quick_validation_passed; 5 files in write-set; tagged repo-vet noise = D3 in-flight |
| 21 | Slice review F01-D1 | claude sonnet subagent | Agent | scratchpad/review-pack-f01-d1.md | scratchpad/review-verdict-f01-d1.md | done | ACCEPT (1 suggestion + 1 nit, non-blocking) |
| 23 | Implement F01-O3 orders app+HTTP evolution | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-o3.md | scratchpad/agent__f01-o3.log · agent__f01-o3.last.md | dispatched | pending |
| 19 | Implement F01-D3 last-sync projection | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d3.md | scratchpad/agent__f01-d3.log · agent__f01-d3.last.md | done | quick_validation_passed; write-set respected; vet green both modes |
| 22 | Slice review F01-D3 | claude sonnet subagent | Agent | scratchpad/review-pack-f01-d3.md | scratchpad/review-verdict-f01-d3.md | done | ACCEPT (2 nits, non-blocking) |
| 24 | Implement F01-D4 dashboard composition | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d4.md | scratchpad/agent__f01-d4.log · agent__f01-d4.last.md | done | quick_validation_passed; 4 new files, dashboard module only |
| 25 | Slice review F01-D4 | claude sonnet subagent | Agent | scratchpad/review-pack-f01-d4.md | scratchpad/review-verdict-f01-d4.md | done | ACCEPT-WITH-CONDITIONS (1 important test gap, 1 nit) — fix row 26 |
| 26 | Corrective: D4 partial-failure sibling assertions | claude sonnet subagent (fallback, core §1) | Agent | (inline instruction) | service_test.go diff | done | sibling survival asserted + test renamed; tests+vet green |
| 26b | Delta re-review F01-D4 | claude sonnet subagent (resumed row 25) | SendMessage | (delta instruction) | scratchpad/review-verdict-f01-d4.md (appended) | done | DELTA ACCEPT — production code unchanged |
| 27 | Implement F01-D5 dashboard HTTP transport | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d5.md | scratchpad/agent__f01-d5.log · agent__f01-d5.last.md | dispatched | pending |
| 15 | Slice review F01-S2 | claude sonnet subagent | Agent | scratchpad/review-pack-f01-s2.md | scratchpad/review-verdict-f01-s2.md | done | REJECT — 1 blocking (4 dead exported error lines), fix row 17 |
| 17 | Corrective: S2 delete dead exports (trivial) | claude sonnet subagent (fallback, core §1) | Agent | (inline instruction) | run_query.go 149→143 lines | done | 4 dead exports deleted; vet+tests green |
| 17b | Delta re-review F01-S2 | claude sonnet subagent (resumed row 15) | SendMessage | (delta instruction) | scratchpad/review-verdict-f01-s2.md (appended) | done | DELTA ACCEPT — no open blocking/important |
| 4 | Implement F01-D2 link counters | gpt-5.6-luna / high | OS-process | scratchpad/prompt-f01-d2.md | scratchpad/agent__f01-d2.log · agent__f01-d2.last.md | dispatched | pending |
