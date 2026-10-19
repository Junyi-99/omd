---
name: kda-neumann-bounded-gate-only
description: Neumann 求解仅限有界 gate（K3 lower_bound），无界 softplus 数值不稳定，不要解耦 safe_gate
metadata: 
  node_type: memory
  type: project
  originSessionId: 05f28c9f-528c-45bd-bb11-d89d5b4dc81c
  modified: 2026-08-17T16:24:23.373Z
---

KDA kernel 的 Neumann 求解（fused-wide 有限级数，`_neumann_fused_wide`）**仅在
有界 gate（Kimi-K3，lower_bound ∈ [-5, 0)）下数值稳定**。普通 KDA（无界
softplus gate）不能用 Neumann 版本——数值不稳定（用户 2026-08-18 明确指示，
覆盖我早期"小规模测试全绿"的观察）。

**Why:** PR#2 把 strip-GEMM 与 Neumann 捆绑在 `safe_gate` 开关后面是刻意的
验证边界，不是偶然的耦合。

**How to apply:** 不要提议把求解方法从 `safe_gate` 解耦成独立开关；generic
路径（safe_gate=False）保持 elementwise 衰减 + 逐行前向消元。generic KDA 的
优化空间只在与 gate 无关的搬运/调度层（fuse / unified_layout / flat_grid /
head_block）。树内 kda.py 的 chunk_kda_fwd docstring 已加警示。
