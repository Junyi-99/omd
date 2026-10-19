---
name: kda-1550-vs-1554-xprof-substep
description: "v6e-1 实测：#1554 mega kernel 比 #1550 快 3.6-6.1×（kernel time），子步骤归因指向 intra_build 和 state_read"
metadata: 
  node_type: memory
  type: project
  originSessionId: f75e7da8-0e4c-4ca4-ac32-3772a62abd60
  modified: 2026-08-31T11:34:04.170Z
---

2026-08-31，目录 `exp-2026-08-31-kda-xprof-scope/`（README + REPORT + RESULTS + 240 条 trace）。
v6e-1 单芯片，jax 0.11.1 / libtpu 0.0.47 / xprof-nightly，60 个 shape（T∈{4k,8k,16k,32k} ×
N∈{1,2,4,8,16} × H∈{8,16,24}），warmup 20 + 100 iter，bounded gate lower_bound=-5。
两边数值一致（rel 4.9e-3 / 6.4e-3 = bf16 噪声），所以比的是两个都对的实现。

**结果：#1554 kernel time 快 3.6-6.1×（均值 4.8×），wall 快 2.3-5.7×。H 越大差距越大。**
T=32768/N=1/H=24：#1550 14.80 ms vs #1554 2.71 ms。

归因（named scope，T=32768/N=1/H=24）：

| 子步骤 | #1550 | #1554 | 比 |
|---|---|---|---|
| intra_build（建 Aqk/L） | 1700 us | 202 us | **8.4×** |
| state_read（w@S, qg@S） | 947 us | 98 us | **9.7×** |
| K1 spill 到 HBM | 525 us | 无 | — |
| solve (I+L)⁻¹ | 719 us | 631 us | 1.1×（打平） |
| gate cumsum | 99 us | 228 us | **0.4×（#1550 赢）** |
| load/cast refs | 128 us | 217 us | **0.6×（#1550 赢）** |

- intra_build 差 8-11× 的原因：#1550 `for h in range(H)` 发 2·(BT/SB)·H=192 个
  `[64,128]×[16,128]→[64,16]` 的 dot，输出只占 128×128 MXU tile 的 1/16；#1554 把 head 放进
  dot_general 的 batch 维（一条指令覆盖 MB=min(H,32) 个 head）、q/k 叠成 32 行单个 LHS、只算下三角块。
- state_read 差 8-10× 的原因：#1550 用 `Precision.HIGHEST`（TPU 上是多遍 bf16）且每 head 分开发 dot；
  #1554 是 `[w;qg]@h` 单次 DEFAULT 单遍 bf16、batched over heads。
- #1550 还有 **22-37% 的时间根本不在 Pallas kernel 里**（`_align_seqs`/`_unalign_output` gather +
  padding + XLA 里的 q/k L2 norm），随 N 增长；#1554 只有 1-3%。
- #1550 赢的两项确认了之前 code review 的预测：#1554 的 `[BT,MB,K]↔[MB,BT,K]` relayout 是真实成本，
  只是比它赢的地方小一个量级。

**口径**：S=T/N 全部 64 对齐，所以 #1550 的 align 不产生 padding、#1554 只走 FULL tile
（PARTIAL/BOUNDARY 分支 count 为 0）——两边都是 varlen 的最好情况。ragged 长度和 unbounded gate 没测。

**Why:** 用户是 #1550 作者，这是两 PR 竞争同一功能的定量证据，也给出了 #1550 的改造顺序
（batch head 轴 → 去 HIGHEST/合并 dot → 融合 K1/K2 → 换 native segment ids）。
**How to apply:** 引用 `exp-2026-08-31-kda-xprof-scope/REPORT.md` §3/§6。
相关：[[kda-1554-boundary-h0-bug]]、[[xprof-named-scope-pallas]]、[[kda-tpu-master-report]]。
