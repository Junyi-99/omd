---
name: blank-lines-between-blocks
description: User wants blank lines separating logical blocks inside functions
metadata:
  node_type: memory
  type: feedback
  originSessionId: 6bdc8a1a-301f-4d3d-9bb6-2bdb1171d0b7
  modified: 2026-10-08T12:33:15.488Z
---

Inside function bodies, separate logical steps (setup / loop / compute / output) with blank lines.

**Why:** 用户说"我读代码都是一个 block 一个 block 读的"，空行分区方便阅读。

**How to apply:** Any code written or edited for this user — group related lines, put one blank line between groups; don't write long unbroken function bodies.
