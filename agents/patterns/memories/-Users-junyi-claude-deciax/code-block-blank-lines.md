---
name: code-block-blank-lines
description: Separate logical blocks inside function bodies with blank lines; the user reads code block by block
metadata:
  type: feedback
---

Inside function bodies, put a blank line between logical blocks: validation, setup, main computation, loops, return. Tests separate arrange / act / assert. Short functions stay as is. Check with `git diff --ignore-blank-lines` (should be empty when only spacing changed) and `ruff format --check`.

**Why:** User (2026-10-15): "代码应该有故意的空行，因为我读代码的时候是按 block 去读的。" Repo commit a9af4a2 applied the same style to the original code.
**How to apply:** Any code written or reviewed in this repo, including code from other branches/agents before it goes into a PR.
