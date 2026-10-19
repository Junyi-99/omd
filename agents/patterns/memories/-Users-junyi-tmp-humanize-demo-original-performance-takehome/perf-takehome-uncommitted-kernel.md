---
name: perf-takehome-uncommitted-kernel
description: The optimized perf_takehome.py kernel is uncommitted; git checkout reverts it to baseline
metadata: 
  node_type: memory
  type: project
  originSessionId: 80cb2660-2d30-41fc-8173-92a192b0448d
  modified: 2026-08-30T17:27:03.990Z
---

In `~/tmp/humanize-demo/original_performance_takehome`, the optimized kernel
(1123 cycles, 131.6x as of Turn 46) lives ONLY as an **uncommitted working-tree change**
(`git status` shows `M perf_takehome.py`). The committed version at `e7b25a7`
("the task") is the 147734-cycle scalar baseline.

**Hazard:** `git checkout perf_takehome.py`, `git restore`, or `git stash`
silently reverts to the baseline and destroys the optimized kernel. This
happened once (Turn 42) and required recovery from conversation context.

**Why:** the takehome ships the optimization as a dirty working tree, not a
commit — intentional.

**How to apply:** to revert an *experiment*, first `cp perf_takehome.py
/tmp/keep.py` and restore from that copy, NEVER from git. Before any risky git
op, back up the file. `git fsck --lost-found` dangling WIP commits may exist but
can be stale/broken — trust a manual file copy instead.
