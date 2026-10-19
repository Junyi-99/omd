# Memory index

- [KDA hardware utilization 结论](kda-hardware-analysis-conclusions.md) — v6e 实测：pallas 1.30× 但 roofline eff 20%，Neumann 逆 27.5%@2.3% MXU
- [inv_precision + C=64/128 实验结果](kda-inv-experiment-pending.md) — bf16 逆历史数据；C=64 range 问题已由 sub16 解决
- [kda-lower-bound release](kda-lower-bound-oracle.md) — K3 oracle + sub_chunk=16 kernel，v6e 3.35×，用户要 release 此目录
- [脚本必须落盘本地](feedback-save-scripts-locally.md) — 禁止一次性 heredoc，写完存 kda-lower-bound/ 再 scp
- [fused-Neumann 新冠军](kda-fused-neumann-champion.md) — fused-wide bf16 @BT=64 = 80ms/2.5TF；fused 才 bf16 稳；BT=16 管线在 v6e 不可行
- [实验目录规范](feedback-experiment-folders.md) — 每轮实验单独 exp-<日期>-<主题>/ 目录，README+报告+数据+图+脚本+capture 齐套
- [aug-18-shipping master report + 分工规则](kda-tpu-master-report.md) — report全量/PR细节/blog主线；LEDGER四态；第六级未ship
- [1554 bounded BOUNDARY h0_B bug](kda-1554-boundary-h0-bug.md) — 已复现+修补验证；报告在 exp-2026-08-28-1554-vs-1550/code-comparison-1554-vs-1550.md
- [1554 边界 bug 真机验证实验](kda-1554-boundary-tpu-run.md) — 已完成：TPU 复现与 CPU 逐位一致，修补版过关；VM 仍保留
- [tpu-spot 新 VM 必须先加 ssh Host 块](tpu-spot-ssh-host-block.md) — 缺块则 READY 后永远 SSH 不通、空转到被回收；check 不查
- [#1550 vs #1554 xprof 子步骤归因](kda-1550-vs-1554-xprof-substep.md) — v6e-1 实测 #1554 快 3.6-6.1×；瓶颈在 intra_build/state_read
- [Pallas named_scope + xprof 配方](xprof-named-scope-pallas.md) — 一进程只出第一条 trace；custom call 的 XLA Ops 事件会消失；trace 要 trim
- [Mosaic 布局传播](mosaic-layout-propagation.md) — 转置不物化，改生产端布局的账由所有 VPU 侧消费者付
- [TPU kernel 优化第一性原理 (knowledge/)](tpu-kernel-first-principles.md) — MXU 指令只差 1.10× 而耗时差 7.6×；布局是全局属性，局部优化慢 5×
- [sglang-jax PR #1564 评审](sglang-jax-pr1564-review.md) — HiCache 复活 double-free（高）、SWA 锁回执跨分裂失效、check_memory 放宽；FULL+SWA 主路径 40 seeds 守恒
- [tpu-spot 上的 sglang-jax 环境/模型缓存坑](tpu-spot-sgljax-env-cache.md) — sgljax tar.gz 不符 env_cache 命名；无 gemma 权重；unsloth 镜像非 gated；us-east1-d 无容量
- [commit 不带 Claude-Session](feedback-no-claude-session-trailer.md) — 上游仓库提交历史不放私人 session 链接
