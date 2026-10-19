---
name: laya-bench-findings
description: Headline results of the laya single-GPU latency study (eager vs TileLang vs torch.compile CUDA graphs) on 2026-10-12
metadata:
  node_type: memory
  type: project
  originSessionId: eba5c113-e6ae-4eb2-8f61-07795e767711
  modified: 2026-10-12T12:07:38.528Z
---

Eager 1-question `router.predict` is 13-26 ms on 3090/A100/A6000/H100/H200 (H200 re-run 2026-10-13 on an idle GPU, CPU load ~146) (upstream T4 claim 39.5/32.8 ms) and is CPU/launch-bound: ~1200 kernel launches, GPU busy 22-58%, 42-46% of GPU time is fp32→bf16 weight re-casting under autocast. A6000 host beats A100 at 1 question because of faster single-core CPU and zero load; A100/H100 win at 50 questions.

Fast paths (parity 4/4, max |Δp| ≤ 0.013): TileLang kernels + bf16-resident weights cut 1-question latency ~2.5-3x (H100 english 5.7 ms, multilingual 3.6 ms); CUDA-graph replay adds only 0.2-0.5 ms on top because upstream's graph covers only encoder+head (48 launches remain). `compile=True, compile_mode="reduce-overhead"` is best at 50 questions (both paths pad rows to 64; TileLang's fixed-tile GEMMs lose at large M) but costs 1-2 min warmup per model and still re-casts weights.

**How to apply:** when asked about laya speed, quote these as measured on shared NUS hosts with ~20% run-to-run noise; HTTP adds ~1 ms (3-4 ms on loaded hosts). Details in [[laya-bench-project]] repo `report.html`.
