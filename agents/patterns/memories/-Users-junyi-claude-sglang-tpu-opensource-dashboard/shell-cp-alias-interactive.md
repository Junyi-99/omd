---
name: shell-cp-alias-interactive
description: "In this user's shell, cp is aliased to `cp -i -v`; overwriting with plain cp blocks on a y/n prompt and hangs agents/scripts"
metadata: 
  node_type: memory
  type: project
  originSessionId: ece58979-6af1-4700-94ff-a79c3fbb7ac8
  modified: 2026-09-23T06:38:30.729Z
---

The login shell aliases `cp` to `cp -i -v` (interactive + verbose). Any `cp` that would overwrite an existing file stops at a y/n prompt; in a non-interactive agent shell that looks like a hung command.

**Why:** A subagent building `genesis-8` hung on exactly this (overwriting `mapping/00-meta-prompt.md`) and was then killed by a rate limit, leaving a half-finished directory. The `-v` part is also why every `cp` prints a `src -> dst` line.

**How to apply:** When overwriting, use `\cp`, `command cp`, `cp -f`, or write via redirection (`cat src > dst`). Same caution for `mv`/`rm` if they turn out to be aliased too. Related: [[genesis-self-reference-iteration]]
