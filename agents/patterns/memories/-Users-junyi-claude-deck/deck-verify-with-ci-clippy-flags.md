---
name: deck-verify-with-ci-clippy-flags
description: "In deck, plain `cargo clippy` is not enough — CI uses `-D warnings` on newer stable, so verify with the exact CI commands"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: a7bf189b-7f08-4377-ab53-451ba528dd3b
  modified: 2026-07-29T07:05:48.272Z
---

deck's CI (`.github/workflows/ci.yml`) has two jobs and no fmt job:
`cargo clippy --all-targets --all-features -- -D warnings` and
`cargo test --all-features`, both on `dtolnay/rust-toolchain@stable`.

**Why:** Plain `cargo clippy` on an older local toolchain passed while CI
failed — CI's stable was 1.97 (which added `byte_char_slices`) against a local
1.96, and `-D warnings` turns any lint into an error. Three pushes went out red
because "clippy is clean" locally meant nothing. There is no fmt job, so
`cargo fmt` drift never surfaces in CI and has to be caught by hand.

**How to apply:** Before pushing, run the CI commands verbatim, plus
`cargo fmt --check` (which CI will not catch for you). After pushing, actually
check `gh run list` instead of assuming green. Note that a red CI does not block
a release — the Release workflow has its own `check` job and does not gate on
clippy, so a lint failure can ship a perfectly good binary while main stays red.
Related: [[deck-release-lightweight-tags]], [[deck-push-direct-to-main]].
