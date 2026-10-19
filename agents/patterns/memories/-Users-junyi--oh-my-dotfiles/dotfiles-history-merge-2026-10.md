---
name: dotfiles-history-merge-2026-10
description: On 2026-10-09 local main was merged with the rewritten origin/main (unrelated histories); old history kept in branch backup/local-main; other Claude sessions commit to this repo concurrently
metadata:
  node_type: memory
  type: project
  originSessionId: db587ed6-88e2-4e98-a0ee-3c3757607140
  modified: 2026-10-09T10:39:53.861Z
---

On 2026-10-09 `origin/main` (github.com/Junyi-99/dotfiles) had a rewritten history
(squash commit `7915e25 oh-my-dotfiles` + Go `cmd/omd` rewrite). Local `main` was merged
with `--allow-unrelated-histories` (merge commit `49bb84b`); the pre-merge local history
is saved as branch `backup/local-main`. Resolution rule used: remote wins for omd/nvim/
README/.gitignore/.funcs.sh, local wins for etc/tmux, zshrc hand-merged.

**Why:** the user chose the merge over a reset so nothing is rewritten; the backup branch
is the escape hatch if something from the old history is missed.

**How to apply:** the user runs several Claude sessions on this repo at once — commits can
land mid-task. Re-check `git log`/`git status` immediately before any reset, rebase or
stash, and never assume the session-start snapshot is current. `etc/omd/` is gitignored
(packages.conf, settings.conf, env.sh live there). A Docker harness that follows the
README across distros lives only in the session scratchpad; recreate it rather than look
for it.
