---
name: shared-gpu-benchmark-hygiene
description: "Latency benchmarks on shared xtra GPU servers need exclusivity checks with heartbeat logs, serial runs, and per-run nvidia-smi/ps snapshots"
metadata:
  node_type: memory
  type: feedback
  originSessionId: a349f4dd-72fd-4b8d-af26-94860a64042d
  modified: 2026-10-15T03:49:08.506Z
---

On the shared xtra* GPU servers, other users (vLLM servers etc.) grab GPUs at any time. In the 2026-10-15 MPK/PDL benchmark, an audit could not prove most rows clean.

**Why:**
- The co-tenant monitor only logged when foreign processes were present, so silence proved nothing.
- Checks ran only at job start.
- Our own jobs ran in parallel on other GPUs of the same host and once collided on the same GPU.
- Logs were overwritten on restart.
- The user audits closely: they spotted mpi4py/TRT processes in nvidia-smi themselves.

**How to apply** (when benchmarking on these hosts):
- Run one job at a time per host.
- Poll every GPU continuously, writing a heartbeat line on every poll, free or not, and covering our own processes too.
- Record `nvidia-smi --query-compute-apps` + `ps -o pid,user,args` at the start and end of each run.
- Wait for our own killed processes to release memory before the next run.
- Never overwrite per-run logs; timestamp every log line.
- Tag each CSV row with run_id/start/end.
- Don't state "clean" or attribute a PID without evidence; nsys sqlite PROCESSES tables list all host processes and can identify foreign PIDs.

See [[xtra-disk-placement]].
