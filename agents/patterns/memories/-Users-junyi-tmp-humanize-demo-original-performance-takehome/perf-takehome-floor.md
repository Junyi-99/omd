---
name: perf-takehome-floor
description: "perf-takehome converged floor — exact slot counts, steady state is LOAD-bound, dead levers"
metadata: 
  node_type: memory
  type: project
  originSessionId: 8a8ccabb-c2df-4dfb-9f0b-defa8b205521
  modified: 2026-08-31T04:27:32.899Z
---

**Turn 120 firsthand recount (same result):** load 2125/1063, valu 6325/1054
(used 1088/1092 bundles = co-binder), alu 12218/1018, flow 555/50.8% (slack but
madd valu-only → unusable). 1092 = 1063 + 29 drain. Level-4 cache still triply
dead (0 free scratch, valu ceiling ~1080, measured 1137). No edit made.


1092 cycles / 135.29x / 9/9 is the converged optimum. Exact compiled slot counts
(built kernel & counted, firsthand): load=2125 (2048 gather + 43 const + 34 vload)
→floor 1063; valu=6325→floor 1055; alu=12218→1019; flow=555.

**Correction (firsthand occupancy profile):** the STEADY STATE is LOAD-bound, not
valu. load saturated in 1062/1092 bundles; load floor 1063 > valu 1055. Profile:
b0-100 ramp (compute/latency-bound), **b100-1000 load pegged 100% = the 2048-gather
wall (~900 cyc)**, b1000-1092 drain (latency-bound, round-15 serial hash). So
1092 = 1024 gather-floor + ~68 latency-bound ramp/drain (serial hash depth 12,
ISA-minimal). flow runs ~36-81% (real headroom) but the L4-mux that would use it
is scratch-dead (below).

**Cleanest kill-argument for any future gather-cut idea (arithmetic, not lore):** load
binder 1063 sits only **8 cycles** above the valu floor 1055. Cutting gathers converts
load→valu AND adds vector-mux valu; valu is already at 1055, so any gather cut caps near
~1055 while pushing valu up → net regression. Measured 3× (L4 round-4 1137, L4 round-15
1228, 25-slot ring 1159). Do not re-attempt a cache-more-levels / ring-resize edit.

Why 2048 gathers is fixed: `problem.py` `indices=[0]*batch` → tree level is static
per round (0..10,0..4) → levels 0-3 muxed, 4-10 gathered = exactly 2048.

**ISA-level why (SLOT_LIMITS in tests/frozen_problem.py:48, so untouchable):** load=2
slots/cyc → 2048 gathers/2 = **1024-cycle hard wall**. Gathers MUST use scalar `load`
(1 elem/slot) because indices diverge per element; `vload` needs a scalar contiguous
base so it structurally cannot vectorize a gather. 1092 = 1024 + 68. Nothing that keeps
the 2048 gathers can go below 1024.

Hash is ISA-minimal: 12 valu/hash. s0/s2/s3/s4 fold into multiply_add (only 3-input
op); s1,s5 = (a^K)^(a>>s) = 3 ops (shift can't fold into a multiply). valu floor 1055
cannot drop → nothing beats ~1055.

DEAD levers (don't re-try): const→flow (already done+balanced 2:1 load:flow,
perf_takehome.py:418-444,737-743); L4 node muxing — even the round-15-only split-mux
(fits flow throughput 480 vselect→1035, turn-102 parity-narrow dodges latency) needs
128 scratch words (16 nodes × 8-wide broadcasts, irreducible for a per-elem vector mux)
and scratch is 0-free (idx/val/t1 cross-round-live, t2 concurrency-bound, top_v live
through round-15 drain). Two measured regressions (round-4 1137, round-15 1228).
No honest one-edit cut exists; a regression violates keep-if-better. See [[perf-takehome-fast-proxy]].
