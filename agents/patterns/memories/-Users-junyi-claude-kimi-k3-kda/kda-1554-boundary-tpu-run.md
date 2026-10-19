---
name: kda-1554-boundary-tpu-run
description: "kda-1554-boundary-v1 已完成：v6e-1 真机复现 #1554 边界 bug（与 CPU 逐位一致）并验证 3 行修补；VM 仍保留计费"
metadata:
  type: project
---

2026-08-28 启动 `kda-1554-boundary-v1`（tpu-spot infra，v6e-1 spot，us-east1-d，VM `kda-1554-bnd-v6e1`），
目录 `exp-2026-08-28-1554-boundary-tpu/`，单 shard `boundary`，产出 `boundary_original|patched.{json,md}` 到
`gs://medusa-experiments/experiments/kda-1554-boundary-v1/`。镜像 `ghcr.io/junyi-99/kda-tpu-bench:latest`。

约束：infra preflight 只放行 us-east1-d / europe-west4-a；europe 报 "Failed to perform tenant project creation"（非容量）。
us-east1-d 2026-08-27 18:00–20:30Z 一直 "no more capacity" 拿不到 v6e-1；同期 v6e-4 能建（kda-1554-vs-1550-v6e4 CREATING）。
同一时刻 `kda-1554-vs-1550-v1` 也在 us-east1-d 抢 v6e-1，两者会竞争。

结果（2026-08-28 01:39Z SUCCEEDED，data/ 已回收）：原版 bounded max|Δo| 0.15–0.26、max|Δht| 到 0.50，与 CPU 逐位一致；
修补版 1.95e-3 / ≤1.3e-2，probe B@tile0 0.93→4.0e-3（= 未受影响的 A 行）。bounded 行比 unbounded 行噪声大是 K3 路径 DEFAULT 单遍 bf16 所致。
卡了 5 小时的真因是 ~/.ssh/config 缺 Host 块（见 [[tpu-spot-ssh-host-block]]），不是容量。VM kda-1554-bnd-v6e1 (34.74.99.121) 仍 READY 在计费。

**Why:** 评审 #1554 需要 TPU 上（bf16 MXU）的误差数字，CPU interpret 的 5e-4/2e-6 不能直接引用。
**How to apply:** `./infra/tpu_experiment.py status kda-1554-boundary-v1` 轮询；拿到结果放进 `data/`，更新
[[kda-1554-boundary-h0-bug]] 与 code-comparison §4。若 us-east1-d 长期无容量，可改 v6e-4 spec（compute 仍单 chip）。
