---
name: mosaic-layout-propagation
description: Pallas/Mosaic 不物化转置，而是给整个值选一种物理布局；改生产端布局的代价会摊到所有 VPU 侧消费者
metadata:
  type: reference
---

2026-09-01 在 sglang-jax #1554 mega KDA 上实测（v6e-1）。

把 `[MB,BT,K]` 上的 Hillis-Steele 扫描改成在 `[BT,MB,K]`（BT 落主轴，移位免费）上做、
结尾 `.transpose(1,0,2)` 转回来 —— 预期省 ~460 ns/chunk，**实测整个 kernel 慢 2.7–6.2×**。

原因不是转置贵：
- XLU 指令只从 48 涨到 50，**Mosaic 根本没物化那个转置**，它把转置退化成 relabel，
  让这个值继续留在 `[BT,MB,K]` 推出的 unpacked 布局里（每 vreg 1 个 sublane 行，不是 8 个）。
- 于是代价摊到**所有逐元素消费这个值的 scope**，不只是生产端：
  s1 扫描 6.7×、s8 `q*exp2(g)` 6.4×、s3 `k*exp2(g)*β` 5.2×、s2（一半 MXU 一半 exp2）3.1×；
  纯 MXU 的 s9/s10/s11 完全没变（1.0×），因为喂 MXU 前本来就要重新打包。
- 指纹：sigmoid 的 EUP 指令数正好 ×8（384→3072）而元素数没变 = 同样的活以 1/8 向量宽度跑。
- 另外 `g_ref[0].transpose(1,0,2)` 这种**从 ref 读出来时的转置是折进 block load 的，免费**；
  改动等于拿一个免费的 load-time relayout 换了一个全局布局决策。

**Why:** 「把转置挪个位置」这一整类微优化在 Pallas 上不成立——一个值全生命周期只有一种物理布局，
生产端和 MXU 消费端布局偏好一旦冲突，账单由全部 VPU 侧消费者付。
**How to apply:** 改 Pallas 里某个共享中间值的布局前，先列出它的所有消费者；用
`llo_scope_mix.py` 看 EUP 指令数是否与元素数成比例（不成比例 = 布局被拆包）。
分支 `wt-1554/s1-major-axis`（KDA_S1_MAJOR_AXIS，默认关）保留为负面结果记录。
相关：[[kda-1550-vs-1554-xprof-substep]]、[[xprof-named-scope-pallas]]。
