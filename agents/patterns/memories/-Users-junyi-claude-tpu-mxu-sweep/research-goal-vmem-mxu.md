---
name: research-goal-vmem-mxu
description: Core research question of tpu-mxu-sweep — how to allocate VMEM so the MXU runs at peak
metadata:
  node_type: memory
  type: project
  originSessionId: e0390c91-ec12-43eb-93e6-71bbfe6a0839
  modified: 2026-10-09T12:22:29.835Z
---

The core question (user, 2026-10-09): 如何分配 VMEM 从而让 MXU 发挥最大效能 — how VMEM allocation (tile shapes under a VMEM budget) drives MXU utilization on TPU v6e.

**Why:** all sweeps (square N, tiled bm/bn/bk, VMEM budget levels 128→4 MiB, small-matrix packing, M/K/N alignment) are means to this end.

**How to apply:** frame experiment design and result write-ups around VMEM footprint vs TFLOP/s; record each config's VMEM footprint alongside throughput.
