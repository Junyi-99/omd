---
name: save-raw-experiment-logs
description: "User wants raw experiment/benchmark logs preserved durably (GCS), not just summaries"
metadata: 
  node_type: memory
  type: feedback
  originSessionId: b7dcdd8d-8b34-426a-801c-f66530b17003
---

When running benchmarks / experiments on the spot TPU, always save the raw logs
(driver log, server.log, per-point bench logs, correctness/eval stdout) to GCS —
not just the parsed summary tables.

**Why:** spot VMs are ephemeral (preempted → wiped, /tmp lost). Raw logs hold the
evidence (e.g. the jagged-kernel `CompileTimeScopedVmemOom` trace, "Bad Request"
maxlen edge, sanity-prompt outputs) needed to explain results and file PR findings.

**How to apply:** push raw logs to `gs://medusa-experiments/tco-bench/<tag>/raw-logs/`
(or the run's GCS_OUT). `stage_bench.sh` already pushes its result dir per-point;
additionally copy any `/tmp/*.log` (correctness runs, ad-hoc evals) before they can
be lost. Pull a copy into the local repo (tco-analysis/) once the run is done.
