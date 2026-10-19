---
name: tpu-kernel-first-principles
description: knowledge/ 目录里沉淀了 TPU Pallas kernel 优化的第一性原理 + 测量配方，从 KDA #1554 vs #1550 实测反推
metadata:
  type: reference
---

`~/claude/kimi-k3-kda/knowledge/`（2026-09-05 建）：

- `tpu-kernel-first-principles.md` — 八条原理。起点是一条实测：差距最大的子步骤里
  **MXU 指令数只差 1.10× 而耗时差 7.6×**，所以首要度量是「每条 MXU 指令周围裹了多少条
  非 MXU 指令」，不是 FLOPs。要点：迭代轴放主轴（sublane 上做 shift 要过 VMEM）；批量轴
  放 `dot_general` batch 维；**布局是整条链的全局属性，局部优化实测 5.2–5.4× 慢**
  （Mosaic 反传 relayout，EUP 正好 ×8 = sublane 数）；`HIGHEST` 是多遍开销要量误差再定；
  静态 grid 要恰好等于真实工作量；融合必须配 tiling 旋钮（否则 H≥48 VMEM OOM）；
  合并 dot 和减 FLOPs 放最后且先确认单元饱和。附全部实测为负的反模式表。
- `tpu-profiling-recipe.md` — 怎么拿数据。见 [[xprof-named-scope-pallas]] 的坑，另加：
  同名 scope 会合并计数（分母要用只出现一次的 scope，我一度把 7.6× 报成 3.8×）；
  计时/标注要分两个 build（trace region 让紧的 kernel 慢 47%、松的慢 12%）；
  per-chunk 与 per-token 两种归一化的差异本身就是信息。

**Why:** 这批原理是从 252 个 shape 的 sweep + LLO 指令计数反推的，含三个负面结果
（都在 worktree 里留了默认关闭的实现），比正面结论更难重新获得。
**How to apply:** 下次做 TPU kernel 优化先读 `knowledge/`，别从 FLOPs 或 roofline 起手。
相关：[[kda-1550-vs-1554-xprof-substep]]、[[xprof-named-scope-pallas]]、[[feedback-experiment-folders]]。
