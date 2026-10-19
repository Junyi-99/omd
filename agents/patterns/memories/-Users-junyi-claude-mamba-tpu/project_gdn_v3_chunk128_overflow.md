---
name: project_gdn_v3_chunk128_overflow
description: v3 fused GDN prefill gives mmlu_pro=0 in vllm — root cause is CHUNK=128 inverse overflow; fix is in-tile 2×64-subchunk recurrence
metadata: 
  node_type: memory
  type: project
  originSessionId: 2c9685e2-d246-42ef-a65f-42c042fb5631
---

The v3 fused GDN prefill kernel (`tpu-inference/.../kernels/gdn/v3_fused_prefill/`,
mirrored in `tpu-spot/gdn-v3-fused-prefill/`) was microbench-fast (3.3–7.3×) and
matched production standalone, but in the full vllm lm-eval (Qwen3.5-4B, mmlu_pro)
scored **exactly 0**.

**Root cause (confirmed 2026-06-01 via in-kernel `GDN_V3_DIFF=1` diff vs production
on real inputs):** v3 is per-call correct on most chunks (out_d~1e-4), but on
pathological 128-token chunks (correlated keys) its intra-chunk `(I+S)⁻¹` overflows
to ~1e35 → ±inf → NaN cascade → garbage logits → 0. Production never hits this because
it runs at **chunk_size=64** (a 64×64 unit-lower inverse maxes ~1e18, inside f32;
128×128 reaches ~1e38). Newton–Schulz vs Gaussian is NOT the issue — the chunk SIZE is.
`test_dropin.py` passed only because it compared v3-vs-prod **both at 128** with random
inputs that never produced a pathological chunk.

**Fix chosen by user (2026-06-01): process the 128-tile as 2×64 sub-chunks with an
in-kernel state carry between them** — keeps full 128-tile utilization (no wasted pad)
and the ~3–7× speedup, while every inverse is 64×64 (bounded). Sequences are packed to
multiples of 128, so both 64-halves of a tile are always the same sequence → reset/
init_h logic stays per-128-tile (sub-chunk a only); only `gcs` becomes per-64-block
cumsum, and S / attn_i use block-diagonal masks. See [[project_spot_v6e_churn_caching]].

**Why:** the kernel's core CHUNK=128 design is numerically unstable for real
activations; this was the blocker to E2E correctness.
**How to apply:** when touching the v3 kernel, never form a 128×128 inverse; keep
inverses at 64. Validate math with a numpy reference before burning a box cycle.
