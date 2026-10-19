---
name: deck-push-direct-to-main
description: "For deck, commit and push straight to main — don't open PRs (overrides the repo's CLAUDE.md)"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: 4bc2cdd8-12f1-4ca3-ab50-ee3e88432416
---

For the deck project, the user repeatedly asked to commit and push directly to `main` and NOT open pull requests — e.g. "直接 merge 进 main，不要开 pr" and "直接 merge 到 main，不要提 pr 了" — even though deck's `CLAUDE.md` says to branch + open a PR into main.

**Why:** solo / fast iteration; the user owns the repo and reviews changes inline as they're made.

**How to apply:** after the user approves a change, commit on `main` and `git push origin main` directly. If `origin/main` moved (another worker merged), `git fetch` + `git rebase origin/main` to keep linear history, then push. Only branch/PR if the user explicitly asks. Still stage explicit paths, never `git add -A` ([[deck-no-git-add-all]]).
