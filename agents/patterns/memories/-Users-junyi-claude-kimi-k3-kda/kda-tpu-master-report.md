---
name: kda-tpu-master-report
description: aug-18-shipping/ 是全项目 source of truth；分工=report 全量 / PR 细节 / blog 只讲 reasoning 主线
metadata: 
  node_type: memory
  type: feedback
  originSessionId: b8b46a0a-02da-4981-8eb6-820cddada1e0
  modified: 2026-08-18T06:00:49.164Z
---

2026-08-18 用户定下的整理方法（源自其六层收敛框架）：agent 横向展开、人纵向收敛。
`aug-18-shipping/` 目录（原名 kda-tpu，用户改名）= source of truth：`report.md`（12 节 master report）、
`experiments/LEDGER.md`（~25 条 hypothesis，四态：Confirmed/Rejected/Inconclusive/Not tested）、
`benchmarks/INDEX.md`（三套测量协议 A/B/C 及漂移记录，canonical=B: w20/r100 avg）、
`notes/pr4-data-audit.md`（PR#4 数据审计）、`notes/blog-outline.md`。

**Why:** 之前几十个探索方向分不清哪些是真结论；agent 输出必须先降级为 hypothesis，
只有进入 Findings（现有 F1-F4）的才能进 PR 和 blog。

**How to apply:** 写 blog 只用 report §7 Findings + notes/blog-outline.md 的六幕主线，
少放数字表格（用户明确：detail 归 PR description，blog 归 reasoning）。新实验先记
LEDGER 再谈结论。未 ship 的第六级（hb+交错 2.64ms/18.4×）是下一个 PR 候选。
相关：[[kda-fused-neumann-champion]]、[[kda-lower-bound-oracle]]、[[feedback-experiment-folders]]
