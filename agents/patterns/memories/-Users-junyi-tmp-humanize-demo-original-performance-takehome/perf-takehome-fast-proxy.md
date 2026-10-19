---
name: perf-takehome-fast-proxy
description: perf_takehome cycle count has an exact fast deterministic proxy; importlib copies are unreliable
metadata: 
  node_type: memory
  type: project
  originSessionId: e5c7c878-d875-4e56-9c33-394535b15c66
  modified: 2026-08-30T22:41:25.044Z
---

For the perf_takehome kernel: every bundle has a non-debug slot, so
`len(kb.instrs) == machine.cycle` exactly (1118 at the current optimum). Building
via a **normal** `from perf_takehome import KernelBuilder` is fully deterministic
(stable across PYTHONHASHSEED values), so knob sweeps need only `build_kernel` +
`len(kb.instrs)` — milliseconds each, no sim run.

**Trap:** loading a *copied* file via `importlib.util.spec_from_file_location(name, path)`
gives NONDETERMINISTIC WRONG counts (saw 1180 and 1258 for the byte-identical
kernel, varying by module-name string). Sweep by monkeypatching the
normally-imported class (e.g. `KB._emission_order = ...`) or edit the real file
in place and re-import in a fresh process — never via importlib copies.

Current optimum is **1098 cycles** (Turn 81). As of Turn 82 the binder is **load,
not valu**: turn 80's flex cut dropped valu to floor 1059.8, so load (2128 slots
-> 1064) is now the sole wall = 40-cyc rounds-0-3 ramp + 1024 gathers, + 34 drain
(last vector's serial hash). valu is 4 cyc behind, so load->cache trades over-bind
valu. The old "valu 1073.5 triple co-bound at 1118" framing (turns 54-73) is
STALE. See [[perf-takehome-uncommitted-kernel]].
