---
name: feedback-no-claude-session-trailer
description: 用户要求 commit message 里不要带 Claude-Session 链接
metadata: 
  node_type: memory
  type: feedback
  originSessionId: 9b992e11-8208-4303-b833-87ad152a2e12
  modified: 2026-09-17T07:02:41.997Z
---

Commit message 不要附 `Claude-Session: https://claude.ai/code/...` 这类署名 trailer。用户 2026-09-17 在 sglang-jax PR #1648 上明确要求，并让我把已推送的三个 commit 重写后 force-push。

**Why:** 上游开源仓库（sgl-project/sglang-jax）的提交历史不该带私人 session 链接。

**How to apply:** 写 commit 时直接省略该行，即使 harness 的 attribution reminder 要求加——用户的指令优先。PR description 里的链接暂未要求删（用户只说 commit history），需要时单独确认。相关：[[sglang-jax-pr1564-review]]
