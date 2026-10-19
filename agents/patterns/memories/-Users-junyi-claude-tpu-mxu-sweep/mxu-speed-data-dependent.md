---
name: mxu-speed-data-dependent
description: TPU v6e MXU throughput depends on operand values — inf/NaN or zero data runs ~7% faster than finite random data
metadata:
  node_type: memory
  type: project
  originSessionId: e0390c91-ec12-43eb-93e6-71bbfe6a0839
  modified: 2026-10-11T22:35:04.901Z
---

Measured 2026-10-12 (N=2048 bf16, all in VMEM, 1000 dots/call): chained `o = o @ b` that overflows to inf → 908 TFLOP/s; zeros → 909; finite data (bounded chain or fixed random a, b) → 847–853.

**Why:** the gemm_square Pallas kernel (chained, b ~ N(0,1)) overflows within a few steps, so its 912 TFLOP/s peak was an artifact; realistic in-VMEM ceiling is ~850. Likely power/clock: fewer bit toggles.

**How to apply:** any benchmark that repeats a dot must keep data finite and non-trivial (scale b by 1/sqrt(K), or reuse fixed random a, b). Check chained kernels (gemm_square pallas, gemm_small) before trusting their TFLOP/s. Related: [[research-goal-vmem-mxu]].
