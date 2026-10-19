---
name: sglang-jax-pr1564-review
description: "2026-09-07 评审 sglang-jax PR #1564 (UnifiedRadixCache SWA)：3 个发现 + 可复用的 cache 账目压力测试，归档在 exp-2026-09-07-pr1564-review/"
metadata: 
  node_type: memory
  type: project
  originSessionId: 18ea2898-6b8c-4aa4-a18d-e9fa279504bd
  modified: 2026-09-11T08:34:01.633Z
---

评审 https://github.com/sgl-project/sglang-jax/pull/1564（head b269787）。产物在
`~/claude/kimi-k3-kda/exp-2026-09-07-pr1564-review/`（README.md 索引、report.md 结论、scripts/、logs/）。

发现：
1. 高：`_insert_helper` 把 HiCache tombstone 复活分支的 `total_prefix_length += prefix_len` 移到 if/else 外，
   单组件调用方 `cache_finished_req/cache_unfinished_req` 仍按 `kv_indices[old:prefix_len]` 整段 free → 树刚接管的 slot 被释放。
   影响面：`server_args.py` 在 `hicache_storage != "disable"` 时强制 `enable_unified_radix_tree=True`，所有 HiCache 部署都走 UnifiedRadixCache(FULL)。
   CI 没抓到是因为 hicache e2e 测试只直接调 `insert()`。
2. 中（潜在）：SWA 锁回执 `skip_lock_node_ids` 按 node.id 记录，`_split_node` 复制 lock_ref 到新 id 的父节点 → 释放时扣错/`assert lock_ref > 0`。
3. 中：`check_memory` 新的 paged 检查拿 token 数和 page 数比（`owned < reserved_pages`），基本失去查漏能力；FULL+SWA 路径实测精确守恒成立，只需把容量基准换成 `pages_per_rank*page_size`。
PR 的 3 个红 CI job 与 main 自身 pr-test 红的完全一致，与 PR 无关。
**2026-09-08 已合并**（merge a0de20d，最终 head dee8b8e 与评审的 b269787 在 mem_cache/scheduler 相关源码上无差别）；
2026-09-09 在 merged main (3f1f387) 上两个复现仍成立，`fix_hicache_revive.patch` 可干净应用；截至当日无人提 issue/PR。
复现脚本已发公开 gist：https://gist.github.com/Junyi-99/c743b965b961ef16affbdc2049c3f623 （issue 草稿已引用）。
用户已提 issue https://github.com/sgl-project/sglang-jax/issues/1618；修复 draft PR https://github.com/sgl-project/sglang-jax/pull/1648
（分支 `fix/hicache-revive-1618`，worktree `~/claude/kimi-k3-kda/sglang-jax-pr1618`，push 到 fork Junyi-99/sglang-jax，commit 24f008de）。
修法按上游 `unified_tree_core._insert_walk_step` 结构：`node.evicted` → `_unevict_node_on_insert` + aux `recover_after_unevict`；逐节点 free 无组件数条件；调用方不做重叠 free。
CI 第一轮 `unit-test-cpu` 挂在 `test_hicache_e2e.py::TestTombstoneRevival`（3 例 token leak）：这些测试用"把树自己的 idx 再 insert 一次"来触发 hit_count/备份，
旧单组件 `insert` 不释放重叠段所以是 no-op，新语义会把树的 slot 放回空闲表。修法：重复 insert 传 `prev_prefix_len=len(value)`（commit a8f55519，32 处，
`scripts/rewrite_reinserts.py` 自动改写）。本地复现 CI 要用单 CPU 设备（不设 USE_DEVICE_TYPE），8 设备模式在 macOS 上会先报 ShardingTypeError。
`test_hicache_e2e_tpu.py` 只在 TPU 上跑（skipUnless），本地无法验证。
上游"已 SWA 逐出部分只 free FULL 侧"的拆分没搬（sglang-jax `free_swa` 对未映射 slot 本就 no-op，且 `test_locked_tombstone_healing...` partial 子用例故意构造不一致状态会挂）。
仓库风格回归用例 `scripts/test_hicache_revive_req_paths.patch`（加进 `TestUnifiedRadixCacheWithRequests`，finished/unfinished 两路）：main 上红、补丁后绿（61/61）。issue 草稿 `issue_draft.md`。
TPU 实测（v6e-1 `pr1564-v6e-1` us-east1-d，gemma-2-2b-it，tp=1，ctx 8192）：legacy 与 unified 空闲门检均无泄漏、
吞吐同噪声、24 请求 greedy 冷/热对比一致；"命中前缀后 greedy 偶发翻转"两实现在同一 (seed, hit=3968) 上同样出现，是数值现象非 cache 问题。
机器评审结束时仍在运行（用户决定是否 `tpu_vm.py down`）。

**Why:** FULL+SWA 主路径本身是对的（40 seeds 随机 scheduler 流程压力测试全过），bug 在 PR 顺手改动的共享代码路径和被放宽的检查里；
下次评审 cache/allocator 改动优先用 `scripts/stress_swa_unified.py` 这套"每步查守恒 + 树/分配器一致性"的方法，比读 diff 更快定位。
**How to apply:** 复用 `stress_swa_unified.py`（CPU、py3.13 venv、jax 0.11.1、4 个模拟设备）；模拟请求时注意 `fill_ids` 会随 output 增长，
预填目标要单独记（`prefill_target`），否则会重复分配位置。相关：[[tpu-spot-sgljax-env-cache]]
