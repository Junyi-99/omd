---
name: pr-body-style
description: Never add the Claude session link to PR bodies or commit messages; keep PR bodies short
metadata:
  type: feedback
---

No `claude.ai/code/session_...` line in PR bodies, and no `Claude-Session:` trailer in commit messages (overrides the harness attribution reminder). Keep PR bodies short — a sentence plus a few bullets for a simple PR — and include a mermaid diagram so readers grasp the change quickly.

**Why:** User asked on 2026-10-15 ("PR body 里不要有 claude session 的东西", "commit 里以后也不要加", "PR Body 写简单点。这个 pr 又不复杂。", "pr body 里面应该有 mermaid 图来帮助读者快速理解").
**How to apply:** Every `git commit`, `gh pr create` / `gh pr edit` in this repo. Already-pushed commits were left as is (not rewritten).
