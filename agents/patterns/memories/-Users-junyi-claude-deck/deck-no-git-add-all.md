---
name: deck-no-git-add-all
description: "In the deck repo, avoid `git add -A`; stage explicit paths instead"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: 6045b998-3a62-4518-b100-7c15bb4bc6b9
---

In `/Users/junyi/claude/deck`, do not use `git add -A` / `git add .`. The repo's `.gitignore` only covers `/target`, `error.log`, `.DS_Store`, so untracked `.claude/worktrees/` (an embedded git worktree → gitlink, mode 160000) and `.cursor/` get swept into commits. The user chose not to gitignore them, so they persist.

**Why:** `git add -A` once committed `.claude/worktrees/scope+ratatui-vt100` and `.cursor/hooks/state/continual-learning.json` into a cleanup commit; it had to be amended out.

**How to apply:** Stage the explicit files you edited (or `git add <known-source-dir>/`), then verify with `git status --short` before committing.
