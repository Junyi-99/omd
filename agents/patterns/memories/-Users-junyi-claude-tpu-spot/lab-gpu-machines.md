---
name: lab-gpu-machines
description: NUS xtra* lab GPU machines (ssh junyi@***REMOVED***) — GPUs, NVLink topology, RDMA NICs, probed 2026-10-13
metadata:
  type: reference
---

Probed 2026-10-13 (read-only nvidia-smi / sysfs). ssh as junyi via ~/.ssh/config (ControlMaster); shared multi-user boxes, check load first.

- xtra3090: 8x RTX 3090, PCIe only (PIX/PXB, no NVLink), no RDMA NIC. CUDA 13.0.
- xtraa100: 8x A100-SXM4-80GB, NVSwitch (NV12, 12x25 GB/s), ConnectX-6 Dx mlx5_0 100GbE RoCE ACTIVE at 192.168.1.116/24 (peer seen: 192.168.1.2, Mellanox MAC, unknown host), active_mtu 1024, no nvidia_peermem. perftest installed. Very busy (load ~130, GPUs 0,1,4,5 at 100%).
- xtrah100: 2x H100 PCIe with NVLink bridge (NV12), ConnectX-5 Ex x2 both DOWN. GPUs busy.
- xtrah200: 4x H200 NVL, NV6 all-to-all (18 links x 26.5 GB/s per GPU), no RDMA NIC (X710 10G). GPUs busy.
- xtraa6k02: 2x RTX A6000, NV4 bridge, no RDMA NIC.
- xtraa6k03: 2x RTX A6000, NV4 bridge, ConnectX-6 mlx5_0 100GbE RoCE ACTIVE at 192.168.4.3/24 (no neighbours), port 2 down. Idle. Note: ssh config names it a6k03; user may say a6k02.

No two machines currently share a 100G subnet (a100 192.168.1.x vs a6k03 192.168.4.x, ping fails). No InfiniBand link layer anywhere.
2026-10-13 run gpu-p2p-lab-r1 (harness ~/claude/ici-bench/gpu-p2p, results gs://medusa-experiments/experiments/gpu-p2p-lab-r1): done a6k03, a6k02, 3090 (0-1, 0-4), a100 (2-3, 6-7); skipped h100 and h200 (GPUs busy) -> rerun with `python3 run_lab.py` when idle. 3090 has no P2P (p2p=0). A6000 globaltimer is 1.024 us quantised -> use calibrated clock64.
2026-10-13 r2 (gpu-p2p-lab-r2, after Codex review: system-scope atomic flags via fence + relaxed, ack-boundary validation, per-size copy checks, bidir bw): a6k02/a6k03/a100 pairs done; 3090 busy in r2 so its r1 rows are reused (no ping-pong there). Key: A6000 put RTT/2 1.09 us, A100 2.2 us; pair bw uni/bidir A100 279/534, A6000 52.9/102 GB/s vs TPU pair 93.3/186.7.
Related: [[ici-external-interconnect-refs]].
