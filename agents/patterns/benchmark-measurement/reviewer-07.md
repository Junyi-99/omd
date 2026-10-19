# reviewer-07

- 分类：`benchmark-methodology`
- 方向：异构硬件性能实验的口径对齐、补测与结果发布。
- 领域：TPU/GPU 点对点互联、微基准、实验报告。
- 频率：可核实为 **1 个独立日志 session**；4 份 memory 来自另外的项目，但不能据此算作同一方法复现。下述两项是同一长任务的不同阶段，不是两个独立复现。
- 持久化理由：最有价值的是比较前建立测量契约，以及发现证据缺口后补测再决定能否合并曲线。具体硬件数字、机器地址、固定参数和排版偏好应留在项目资料。
- 范围与局限：流式扫描所分配 JSONL 的全部 user 消息，区分真实请求、工具通知、技能注入和压缩摘要；按主题回读 assistant 上下文。memory 全文读取。日志内引用的论文和生成代码未独立复验；以下审阅评价工作流证据，不认可其技术结论为已验证事实。

## 1. 异构互联基准：先对齐测量契约，再引用、补测和合图

**触发场景**：用户希望把自测性能与不同硬件、第三方论文或厂商数据放到同一图表。

**可复用步骤**：
1. 明确比较对象：芯片对还是单链路、scale-up 还是 scale-out、单向还是双向聚合。
2. 为每组数据记录硬件、消息大小、发起方、原语、同步/完成条件、计时域、统计量、并发度及来源。
3. 查原始文字、表格和基准实现；分开记录实测值、拟合参数、模型常数与二次推算，缺失字段显式留空。
4. 按可比程度组织外部证据，不用库级通信、裸 remote-store、NIC loopback 等数据互相替代。
5. 对尺寸或原语不齐的部分补实验；保留重叠尺寸验证一致性。重叠区不一致时保留独立曲线，说明 kernel/同步差异。
6. 在报告中展示口径、来源和限制；审查由数字推导出的倍数、吞吐及归因，修订后保留原始版本。

**原始证据**（以下 `L` 为原始 JSONL 行号）：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/f980cc93-3b6c-4f08-86a3-25be7763f45c.jsonl:L97`，用户：“如果找到 experiment setting 跟我们一样的，就正好画图放进去了”。
- 同一绝对路径 `:L284`，用户：“特别是小包的情况（TPU 这边是单个 RDMA descriptor）”。两轮 survey 是同一需求深化。
- 同一绝对路径 `:L850`、`:L887`，用户追问 TPU/GPU 尺寸起点不同，并要求“让这个图的起点起码同一个起点”。
- 同一绝对路径 `:L1331`，assistant 复核后：“2.5 ns 和 7.8 ns 是拟合量，不是速率”；同时纠正 single-QP、来源归属和 loopback 表述。
- 同一绝对路径 `:L1569`，assistant：“本次比实验 2 高 8–12%……没有和实验 2 的曲线拼起来”。这是保持口径边界的操作证据。

**频率/信心**：1 session 内至少 4 个直接用户请求支撑调查与补测；多轮纠正体现需求强度，不能算跨 session 高频。方法价值高，独立复现信心中等。

**建议 skill description**：当需要比较自测与第三方硬件性能、合并 benchmark 曲线或补齐测量范围时，建立测量契约和来源表，检查指标语义，设计重叠补测，并按一致性决定合并或分列。

**边界及反证**：早期 survey 曾把拟合量当吞吐、过度概括 single-QP 上限；不能把 survey 汇总本身当准确性保证。共同尺寸起点也不等于测量语义一致。此 skill 不应固化某代硬件数值或强行声明跨平台公平性。

## 2. 候选：按现有机器能力分阶段落地微基准 harness

**触发场景**：已有多台共享实验机器，希望先跑能执行的子集，再扩展到目标网络或加速器。

**可复用步骤**：
1. 只读盘点设备、拓扑、peer access、RDMA 链路、驱动/模块、编译工具和当前占用。
2. 形成“可测 / 缺依赖 / 忙碌 / 不具备拓扑”的能力矩阵，先用空闲机器调通最小 harness。
3. 把正确性校验放在计时外；每尺寸校验数据，明确同步作用域，剔除失败结果。
4. 按目标记录独立状态，忙碌设备跳过；校验与上传成功后才标成功，保留失败和旧版本供追踪。
5. 汇总结果时标记版本、fallback、未测项目及干扰；后续只补未完成目标。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/f980cc93-3b6c-4f08-86a3-25be7763f45c.jsonl:L556`，用户提出用自有机器替代部分云实验。
- 同一绝对路径 `:L597`，用户：“先把 harness 调通吧。在我们机器上把能做的实验做了。做不了的先不做。”
- 同一绝对路径 `:L593`，assistant 的盘点区分 NVLink、不可达双机 RoCE、缺模块及 GPU 占用。
- 同一绝对路径 `:L842`，assistant：“skips busy GPUs, resumes from completed pairs”。
- 同一绝对路径 `:L1331`，assistant：“每个尺寸先跑一遍不计时的 payload 校验”；“上传完成后才写 SUCCESS”。

**频率/信心**：1 次 harness 建立及 1 次 review 后重跑，均属 1 session 的同一实验链。候选，信心中等；尚不足以声称稳定高频习惯。

**建议 skill description**：在已有共享 GPU/TPU 机器上搭建可恢复的微基准实验时，先盘点可执行能力，按空闲目标分批运行，执行计时外正确性校验，保留版本与缺失原因，汇总可追溯结果。

**边界及反证**：r1 出现同步和校验问题，r2 的改正才支持这些步骤；不能假定首次 harness 就可靠。无权限或缺物理拓扑时标记未测，不通过修改共享环境强行完成；具体机器清单不是 skill。

## 未提升为独立 pattern 的材料

- `kda-neumann-bounded-gate-only.md:L11–22`：Neumann 仅限有界 gate，是项目算法适用边界；“小规模测试全绿”仍被后续用户纠正。适合项目约束记忆，不足以抽成通用 Neumann 优化技能。
- `perf-takehome-fast-proxy.md:L11–21`：用构建产物计数替代模拟、并核验导入路径确定性，有方法潜力，但这里只有单项目 memory，未见原始日志复验。精确代理与 importlib 陷阱须留在此项目，不推广为通用 Python 结论。`L23–28` 还明确旧瓶颈判断已过时。
- `pallas-kernel-targets/memory/MEMORY.md:L1`：只有符号约定索引；属于术语记忆。
- `medusa/memory/MEMORY.md:L1`：专业语体偏好索引；属于表达偏好。
- 日志 `L457`、`L484`、`L1575`、`L1635` 的图例移出、字节单位、远程发布与 section 排版均为同一报告迭代；不足以独立创建报告设计或发布技能。
- `L1461` 与 `L1529` 是同一实验的重复监控请求，不算两次工作流复现；通知、技能注入、压缩摘要不计用户需求次数。

## 覆盖情况

- 全文读取 `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/memory/MEMORY.md`。
- 全文读取 `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda-aug-12-learn/memory/kda-neumann-bounded-gate-only.md`。
- 全文读取 `/Users/junyi/.claude/projects/-Users-junyi-claude-medusa/memory/MEMORY.md`。
- 全文读取 `/Users/junyi/.claude/projects/-Users-junyi-tmp-humanize-demo-original-performance-takehome/memory/perf-takehome-fast-proxy.md`。
- 全量流式扫描 `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/f980cc93-3b6c-4f08-86a3-25be7763f45c.jsonl` 的 user turns；选读比较口径、实验盘点、harness、复核修正和补测完成的 assistant 上下文。
