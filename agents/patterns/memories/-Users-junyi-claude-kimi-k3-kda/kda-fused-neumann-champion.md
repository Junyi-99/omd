---
name: kda-fused-neumann-champion
description: sglang KDA solve 的演化终点：fused-wide Neumann bf16 @ BT=64 = 80.05ms 新冠军；fused 形式 bf16 稳定、直接式不稳；BT=16 管线在 v6e 不可行
metadata: 
  node_type: memory
  type: project
  originSessionId: e645184a-6feb-471a-ae36-e740b711b0e0
  modified: 2026-08-13T11:45:25.469Z
---

sglang-jax KDA intra kernel 的 solve 优化链（2026-08-13，B=1 T=8192 H=96 D=128 单 v6e）：

- 串行 solve（原 PR 路径）158.7ms → LLO 证实是延迟瓶颈（所有单元 <8% 忙）
- Neumann 直接式 HIGHEST 85.2ms（bf16 直接式 FAIL 6.9e-04——显式逆平方级放大误差）
- **fused-wide bf16 = 80.05ms 冠军（2.50 TFLOP/s，2.34× vs baseline 187ms）**：Z=[I|Vβ|Kβ] (64×320) 逐因子传播 (I−L)(I+L²)…(I+L³²)，不显式造逆 → bf16 稳定（2.77e-04 过 5e-4 门），零 6-pass 税
- fused-HIGHEST 87.9ms 比直接式 HIGHEST 还慢——fused 只有解锁 bf16 时才划算
- BT=16 全管线（FlashKDA CHUNK=16 绑定）在 v6e 不成立：最好 144ms，4× chunk 数的 fwd_h 流量/grid 步进吃掉 solve 的 10× 每 token 节省

**最终结构（用户拍板，2026-08-13）**：general 路径保留上游原 `_solve_unit_lower_triangular`；
K3 作为特殊分支，以已有 `safe_gate` 参数为键（绑定 bounded gate + [-5,0) 校验），solve 用
`_neumann_fused_wide` bf16。所有实验参数（fast_solve/intra_pack/assume_aligned/k3*）已全部拆除
——用户原则：不加新参数。终态验证：general 187.0ms 不变、K3 80.03ms（2.34×）、精度门 + 上游
kda_test 7/7 全过。改动未提交（sglang-jax 本地树，branch kda-safe-gate-mxu）。
相关：[[kda-lower-bound-oracle]]、[[kda-hardware-analysis-conclusions]]
