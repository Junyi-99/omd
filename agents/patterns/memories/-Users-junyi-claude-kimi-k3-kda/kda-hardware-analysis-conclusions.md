---
name: kda-hardware-analysis-conclusions
description: KDA kernel 在 v6e 上的 hardware utilization 结论（fresh 复测 + HB=16 stage ablation 数据）
metadata: 
  node_type: memory
  type: project
  originSessionId: 911e46b9-0f78-4f3d-8862-5652e40b7599
  modified: 2026-08-03T03:01:23.270Z
---

v6e 单芯片、fp32/highest、B1 T8192 H96 D128 C32（2026-08-01 fresh 复测，
与 07-31 存档逐位一致；artifacts 在 `10-eval/results/v6e-20260801-fresh/`）：

- kda-pallas HB16: 17.32 ms, roofline eff 19.5%, underutilized（不卡任何 roof）
- kda-xla: 22.58 ms, HBM util 0.88, hbm_bandwidth_bound（带宽顶死，无余量）
- HB=1 对照: 46.6 ms（grid-issue 开销 2.7×）
- 校准（单芯片）：821 TFLOP/s (1-pass fp32-data)、1191 GB/s；highest=6 pass → 有效 roof ≈137

HB=16 stage ablation（首测）：Neumann 逆 27.5% 时间 @ 2.3% MXU；s7_output 26.6%；
grid_io 2.1 ms = 959 GB/s = 80% HBM roof（IO 已近满，搬运不是问题）；
算术 stage 利用率 2–19%，病因 = C=32 只占 MXU 1/4 行 + 6-pass fp32 + 串行链。
分 stage FLOPs 手算与 XLA 编译器计数差 0.9%。

优化排序：① inv 降 bf16 单 pass（免 6 pass 税）② C=64（填充×4，需修 gate range NaN）
③ 块三角逆（算术 ÷4.4）④ C=128 不划算（range 要 lb≥-1.375 且逆的无效功 ∝C²）。
坑：analyzer 校准 key 含 used_device_count，多芯片机器上要 TPU_VISIBLE_DEVICES=0；
pkill -f 会自匹配 ssh wrapper 导致断连，模式用 [.] 转义。
相关：[[kda-inv-experiment-pending]]
