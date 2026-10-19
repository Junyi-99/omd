---
name: ici-external-interconnect-refs
description: Published RoCE/InfiniBand/NVLink p2p numbers collected 2026-10-13 for comparison with our v6e ICI results; none match our measurement setting exactly
metadata:
  type: project
---

2026-10-13: collected external p2p numbers into `~/claude/ici-bench/external/` (sweeps.csv, points.csv, sources.md, plot_external.py -> charts/ext_01_latency, ext_02_bandwidth, compare_table.md). Run with `../.venv/bin/python plot_external.py`. Not wired into report2.py.

2026-10-13 (later): expanded with a 4-agent survey -> `survey.md` (small-message NIC latency sweeps, per-descriptor/message-rate, GPU-initiated puts, scale-up fabrics), `per_descriptor.csv`, charts ext_03_small_latency / ext_04_per_descriptor. Closest-match source: arXiv 2610.01380 (single 8 B WQE, device globaltimer, p50; RTT 6.85 us on H200+CX-7). No published device-timed TPU ICI p2p latency exists.

No published source matches our setting (one unloaded transfer, device-side XProf time, p50, single link). Only full per-size sweeps found in text: OCI RoCEv2 H100 OSU (host + GPU buffers, mean, multi-rail UCX) and an H100 + ConnectX-6 IB osu_bibw. InfiniBand NDR per-size curves exist only as figure images (arxiv 2511.15076, 2606.05951).

**Why:** the user wants ICI vs RoCE vs InfiniBand in experiment reports; searching again wastes time.
**How to apply:** reuse these files; for a matched comparison, a GPU-cluster run with the same sizes and statistic is required. Related: [[ici-v6e-measured-constants]].
