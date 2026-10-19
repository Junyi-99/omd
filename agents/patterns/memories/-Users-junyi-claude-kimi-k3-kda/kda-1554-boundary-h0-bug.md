---
name: kda-1554-boundary-h0-bug
description: "#1554 mega kernel 的 bounded gate BOUNDARY tile 在段 B 带非零 h0 时输出错（92–99% 相对误差），已 CPU interpret 复现并定位修补"
metadata: 
  node_type: memory
  type: project
  originSessionId: f57bdb37-d0ec-4a7d-b795-5c90684723c3
  modified: 2026-08-27T19:27:25.127Z
---

2026-08-28 复现：sglang-jax #1554 (cjx0709, mega KDA) `_boundary()` 里 unbounded 分支在建 Aqk/L 前把 B 段
g_cumsum 减去 A 段总和，bounded 分支没做 → B 行 k_eg_beta / b_qg 多乘 2^(g_A_total)，w@h0_B、q@h0_B 被压到 0。
h0_B=0（新请求首 chunk）看不出来；chunked prefill 续算 / 前缀命中带状态会错。B@tile0 相对误差 0.925（A=32）/ 0.995（A=37），
其它组合 ~2e-3。修补（建完 Aqk/L 后对 B 行减 first_segment_total，g_cumsum_B_local 统一为 g_cumsum）后回到 2e-3。

验证：1554 自带 15 测试原版/修补版都过；6 几何×2 gate 慢衰减+O(1) 初态，原版 bounded max|Δo| 0.15–0.26、max|Δht| 到 0.5，修补版 ≤9.8e-4 / 2.9e-6。
1554 自己的 (63,65) 用例盖住了这个几何但 initial_state*0.1 + atol 5e-2 测不出。
TPU v6e-1 真机（kda-1554-boundary-v1）：原版数字与 CPU 逐位一致，修补版 max|Δo| 1.95e-3 / max|Δht| ≤1.3e-2 = 该 kernel 未受影响行的水平。
gist（公开，已发）：https://gist.github.com/Junyi-99/a44b942146a1a66ec851094aba6af4d8
脚本：exp-2026-08-28-1554-vs-1550/scripts/probe_1554_boundary_bounded_h0.py（gist 版）、test_1554_boundary_extra.py、report_1554_boundary_extra.py、patch_boundary_restart.py；
报告：exp-2026-08-28-1554-vs-1550/code-comparison-1554-vs-1550.md（含两 PR 方法论对比表 + 优化空间清单）。

**Why:** 用户是 #1550 作者，两 PR 竞争同一功能；这是 1554 的确定性正确性缺陷，PR 评审可直接引用。
**How to apply:** 讨论 1554 时引用该报告 §4；1554 若修了，重跑 probe 确认。相关：[[kda-lower-bound-oracle]]、[[kda-fused-neumann-champion]]。
