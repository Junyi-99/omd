---
name: User controls state-change actions (push, close, merge)
description: Prepare the change locally, but don't push, close, or merge until the user gives explicit go-ahead — even in multi-step flows where "writing the fix" implies "shipping the fix"
type: feedback
originSessionId: 684260bd-b5d2-47ef-a332-8f3b3ef464cd
---
Prepare code changes, commits, and branches locally; do NOT execute outward-facing state changes (`git push`, `gh pr close`, `gh pr merge`) until the user says to. The default is: write → show diff → wait.

**Why:** The user wants to verify changes locally before they go anywhere visible. Happened twice in one session:
1. I bundled `gh pr comment` (reject message) with `gh pr close` in one command — user interrupted and clarified they wanted comment only.
2. While iterating on a PR, I offered to push the next fix immediately after writing it — user interrupted with "别急着推，我们先本地测完了，我让你推你在推" ("don't rush to push, let's test locally first, I'll tell you when to push").

The pattern: actions that only touch local state (edit files, commit in worktree) can proceed; actions that affect remote state (push, close, merge, create/close issues) need explicit per-action confirmation, even within an approved multi-step plan.

**How to apply:**
- For `git push`, `gh pr close`, `gh pr merge`, `gh pr create` — always confirm immediately before running, never bundle with the change itself.
- For `gh pr comment` — if it's a direct task instruction ("tell the author X"), execute; but don't pair it with a state change.
- When presenting a diff or committed change, end with "OK 我再 push / OK 我再关 PR" style prompt; do not auto-proceed.
- One approval applies to one action, not to a whole chain. Approving "commit this" does not authorize "push this".
