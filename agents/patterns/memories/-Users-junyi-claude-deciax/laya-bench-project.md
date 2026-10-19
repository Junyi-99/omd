---
name: laya-bench-project
description: "Laya (convaiinnovations/laya) benchmark harness repo, the 5 NUS GPU hosts used, and host quirks hit on 2026-10-12"
metadata:
  node_type: memory
  type: project
  originSessionId: eba5c113-e6ae-4eb2-8f61-07795e767711
  modified: 2026-10-12T12:07:30.537Z
---

Benchmark of https://huggingface.co/convaiinnovations/laya (upstream NandhaKishorM/laya @ c73908d, 2026-10-12) lives in the private repo `Junyi-99/laya-bench`, cloned at `~/laya` on the Mac and on each host. Results + chrome traces committed under `results/<host>/`; `report.py` builds `report.html` (gitignored).

Hosts (ssh `***REMOVED***`, user junyi, no sudo, all shared with other users): xtra3090 (8x3090), xtraa100 (8xA100, driver 570 → torch must be cu126, not cu130), xtraa6k03 (2xA6000, usually idle), xtrah100 (2xH100, GPU 0 permanently taken), xtrah200 (4xH200, heavily shared, load ~12-30).

**Why:** hosts have no GitHub creds, so code is rsynced from the Mac; PyPI/NVIDIA wheel downloads time out → run setup with `UV_HTTP_TIMEOUT=600`; default g++ on xtra3090/xtraa6k03 is 9.5 (tilelang/Inductor need g++-11+); installing tilelang flips laya's `backend="auto"` to tilelang, so eager baselines must run before installing it.

**How to apply:** reuse `~/laya/setup.sh`, `run_all.sh`, `run_fast.sh`, `bench_profile_all.py`; ssh ControlMaster block for these 5 hosts is in `~/.ssh/config`; always check `nvidia-smi` for an idle GPU before a latency run and record the load average. See [[laya-bench-findings]].
