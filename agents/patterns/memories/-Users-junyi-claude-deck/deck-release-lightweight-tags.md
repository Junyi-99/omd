---
name: deck-release-lightweight-tags
description: "How to cut a deck release — lightweight tags, the gpgSign gotcha, and writing notes"
metadata: 
  node_type: memory
  type: reference
  originSessionId: 4bc2cdd8-12f1-4ca3-ab50-ee3e88432416
---

deck releases tag `main` directly with `vX.Y.Z` (no branch/PR); pushing the tag triggers a GitHub Actions **Release** run that builds binaries and updates the `cross-entropy-ai/homebrew-tap`.

This machine's git has `tag.gpgSign=true` and `tag.forceSignAnnotated=true`, so a plain `git tag vX.Y.Z` fails with `fatal: no tag message?` (it tries to make a signed *annotated* tag). All prior release tags (v0.6.0, v0.7.0, …) are **lightweight**. Create them with:

`git -c tag.gpgSign=false tag vX.Y.Z && git push origin vX.Y.Z`

Version bump first: edit `Cargo.toml` version, run `cargo check` to sync `Cargo.lock`, commit as `chore: bump version to X.Y.Z`, push `main`, then tag.

After the Release run finishes, write user-facing notes with `gh release edit vX.Y.Z --notes-file <file>` (a quoted heredoc → file, so backticks aren't shell-expanded). The auto-generated notes only list merged PRs and miss direct-to-`main` commits, so diff `vPREV..vX.Y.Z` and group into New Features / Enhancements / Bug Fixes / Under the hood. Wrap literal `@`/`#` (e.g. `@local`, `@deck_order`) in backticks; bare `#NN` only for real PR links. See [[deck-push-direct-to-main]].
