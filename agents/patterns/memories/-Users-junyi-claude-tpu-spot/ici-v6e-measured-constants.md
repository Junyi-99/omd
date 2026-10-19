---
name: ici-v6e-measured-constants
description: "Measured TPU v6e-4 ICI ppermute constants (alpha, beta, per-hop, host dispatch floor) from experiment ici-v6e4-2x2-r1, 2026-10-11"
metadata:
  node_type: memory
  type: project
  originSessionId: 7dff67a6-e2bc-469e-a770-fffc93de1f63
  modified: 2026-10-11T17:37:46.880Z
---

Experiment `ici-v6e4-2x2-r1` (2026-10-11, v6e-4 2×2, jax 0.11.2 / libtpu 0.0.50, 25/25 shards, workload in `~/claude/ici-bench/`, results `gs://medusa-experiments/experiments/ici-v6e4-2x2-r1`):

- Device-side (XProf) ppermute: α ≈ 2.95 µs for 1 hop, 4.13 µs for the 2-hop diagonal (+1.18 µs/hop, +40%); β = 93.3 GB/s for every pair (spec 90). m* = α·β ≈ 275 KB (1 hop) / 386 KB (2 hop).
- Host-side: single jit call floor ≈ 185 µs (dispatch + block_until_ready, 60× the device time); inside a fori_loop of 100 ppermutes it drops to ≈ 4.7 µs per call.
- Zero points: 1-element ppermute 2.78 µs, 1-element psum over 4 chips 5.97 µs. The jitted-identity device time is unstable (p50 15 µs, min 2.3) and is not a usable zero.
- Anomalies: every transfer INTO device 0 (host/profiler chip) is ~0.45 µs slower across all sizes (p50 and min); looped mode at 64 MiB drops 4% (gap between collectives, not transport; β rises 87.9→92.3 GB/s without that point).
- Model points (5–20 KB decode messages) are 93–98% fixed overhead on the device side.

Experiment 2 `ici-v6e4-2x2-r2` (same machine, 2026-10-11, Pallas in-kernel RDMA = path B vs XLA ppermute = path A; results `gs://medusa-experiments/experiments/ici-v6e4-2x2-r2`, analysis `~/claude/ici-bench/out-r2/`):
- Path B α = 2.08 µs (1 hop) / 3.27 µs (diag) vs path A 2.95 / 4.13 (different estimands: B's op includes the barrier handshake and excludes XLA's output copy); β 93.33 GB/s both; per-hop +1.19 µs both paths.
- Per-descriptor cost: fixed-N (4096) mixed-layout fit intercept 7.8 ns, matched tiled control 2.5 ns (no N sweep, so not a universal constant); tiled refs of 1–128 B take near-identical time (34.98–34.99 µs for 4096 descriptors), consistent with tile padding but not proven (semaphore byte reads are 0).
- In-flight DMA count T gives no gain over T=1 at any size; per-copy semaphores cost ≈19.5 ns each (T=128 at 1 MB is −16%). per-copy semaphores verified up to 128 copies × 2 (256 s32 = 1 KiB, a conservative allowance; exact sflag capacity between T=128 and 256 unmeasured); 16 GiB single buffers do not compile (HBM temporaries > 31.24 GiB); 8 GiB is the largest size tested, not a proven hardware max.
- Aggregate: 2bi 186.6, 4pairs-x 186.6, 4pairs-bi 373.3, ring 373.3 GB/s (= links × 93.3 within 0.1%); 4all = 535 GB/s = 1.43× ring, 95.5% of the 559.8 GB/s link-occupancy bound (8×93.3×12/16; diagonals relay through a neighbour), per-flow time 2.09× ring.
- Overlap: in-kernel RDMA + matmul overlap efficiency 1.00; XLA ppermute + matmul 0.97. Chunked consume (64 × 1 MiB, per-chunk wait + HBM→VMEM + accumulate) is +25% vs comm-only, but there is no matched wait-all-then-consume baseline so the per-chunk cost is not isolated; 'first chunk at 2.4%' is a separate 1 MiB T=1 kernel (17 µs), not an in-stream arrival time.
- Device-0-as-destination penalty is +0.72 µs for path B (vs +0.45 for A single, ~0 in looped mode); cause unknown.

**Why:** these constants are needed whenever estimating TPU v6e communication cost; the scaling-book 1 µs/hop assumption is close but the real α is ~3 µs and host dispatch dominates unfused calls.
**How to apply:** use T_us ≈ 2.94 + (hops−1)×1.19 + bytes/93328 (bytes/93.3 GB/s in µs) for v6e ICI point-to-point; α already includes the first hop; prefer fused loops to avoid the 185 µs host floor. Experiment 2 (`ici-v6e4-2x2-r2`, Pallas in-kernel RDMA vs XLA ppermute) follows [[ici-experiment-2-design]]. Machine: on-demand `tpu-v6e-4-od-us-east5-a` in project Medusa Weihao (`project-68023df7-96eb-4779-857`), its service account has objectAdmin on gs://medusa-experiments.
