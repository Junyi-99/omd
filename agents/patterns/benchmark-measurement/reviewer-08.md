分类：`performance-analysis`
方向：性能实验的可信归因与结果可视化。
领域：TPU / JAX / Pallas kernel、吞吐量图表、指令级优化。
频率：1 个原始日志 session；另有 4 份 memory。图表任务有 6 条实质用户请求（初始请求及 5 次修订），属于同一任务，不能计为 6 次独立复现。性能方法在 3 份不同项目 memory 中出现，但来源 session 未分配，无法核实历史执行次数。
为什么值得持久化：性能数字容易受输入布局、插桩与硬件约束影响；图表也容易通过阈值和标注布局漏掉真正的下降点。这些流程具有复用价值，但本批证据只支持 1 个图表候选及 2 个性能方法候选，尚不足以认定高频成熟 skill。
审阅范围和局限：全文读取全部 4 份 memory；Python 逐行解析唯一 JSONL，扫描所有 user 文本，再核实 assistant 文本上下文。图像块和工具结果不计为用户请求；未审查其他 reviewer 结果，也未访问 memory 引用的原 session。以下历史技术参数均为来源当时记录，不能直接视为当前通用事实。

## 1. 性能下降点的检出、解释与图表验收（候选，证据最完整）

**触发场景**：已有 sweep CSV 与绘图脚本，需要标出吞吐量突降点并生成可读的性能图。

**可复用步骤**：
1. 读取完整数据及采样方式，区分粗扫、局部密扫和不连续窗口，避免把断开的窗口连接成趋势。
2. 对相邻可比样本计算相对下降，列出候选点；阈值作为可解释参数，检查用户指出的漏检位置。
3. 在 lowering 或硬件边界有证据时结合边界解释；对于 dispatch floor 附近的小幅波动，核实后区分噪声与真实下降，不能只靠不断降低阈值。
4. 标记峰值位置、下降幅度与对应序列。连接线采用与下降曲线不同的方向，减少重合；控制标签密度和遮挡。
5. 使用双轴方向的浅灰网格，极限/边界参考线使用可区分的线型；重新渲染并查看完整图。
6. 核对图例、正文 Findings 与新数据是否一致；有待确认的解释应标明，不能让旧结论继续充当已验证事实。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/5c668753-ed25-444d-b25f-df9808b7c1c1.jsonl:11`：“重画 gemm_square.png……主要是标出来性能下降的点”。附带完整脚本，已有粗扫/密扫和不连续窗口处理。
- 同文件 `:94`：“3584 后面、4608 前面，也有一次 drop”；`:128`：“7680 后面还有一次”。assistant 在 `:124`、`:168` 记录阈值从 15% 到 12% 再到 8%，并报告 dispatch floor 噪声，说明阈值漏检是实际问题。
- 同文件 `:172`：“标志线应该有倾斜……那个线和性能曲线有点重合了”；`:190`：“垂直方向应该有灰色线条”；`:207`：“Pallas Limit 的那个地方用虚线表示”。
- 同文件 `:219`：最终仍报告参考线穿过标签、Findings 周期描述与标记矛盾。说明本 session 的流程未完全闭环，不能把“验收通过”当作已发生事实。

**频率 / 信心**：1 session，6 条任务/修订请求；任务内持续出现。可复用性中高，跨任务复现证据低。

**建议 skill description**：基于原始 sweep 数据重绘性能图，检查下降点漏检与测量噪声，优化标注、网格和参考线，并核对图文结论的一致性。

**边界及反证**：不要把 8%、256 对齐或 N≈2048 固化为通用规则；这些是该数据集的选择。不要把 Apple 配色、倾斜量 110 或固定标签位置当作方法本身。一个 session 的审美修订尚不能独立证明需要专用 skill；可先作为性能分析 skill 的绘图章节。

## 2. 保持输入契约并分离归因与计时（候选，可并入性能实验 skill）

**触发场景**：比较 fused/custom kernel 与 baseline，或使用插桩分析 VPU/MXU 活动。

**可复用步骤**：
1. 固定 baseline 的函数签名、输入 shape、dtype 与布局，明确计时边界。
2. 候选所需的 reshape/relayout 放入被计时函数；依赖上游改变布局的优化作为另一项实验。
3. 使用独立的插桩运行定位内部活动；以未插桩运行提供性能数字。
4. 检查工具/运行库版本与 flag 支持，保留软件栈与实验配置；不可因默认视图缺失就断言无法观测。
5. 重新计算速度比，记录曾经因测量边界不公平导致的反转。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/memory/epmoe-notation.md:32`：“Any reshape/relayout the candidate needs is part of its cost”。`:36–44` 记录计入 relayout 后速度比由 1.9–2.5x 改为 1.2–1.53x，且一个 shape 从赢变输。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda-aug-12-learn/memory/xprof-custom-call-tracing.md:18`：“归因用插桩 trace，计时数字用未插桩运行”。`:20–26` 记录错误断言被修正、不同 libtpu 版本对 flag 的支持不同。

**频率 / 信心**：2 份不同项目 memory 分别提供 1 个相关案例；不是同一完整工作流的两次复现。方法价值高；独立 session 次数不可核实，skill 成熟度中低。

**建议 skill description**：为 JAX/Pallas kernel 建立公平性能对比，校验输入布局和计时边界，并分离插桩归因与正式计时。

**边界及反证**：memory 中的具体 flag 与版本属于参考资料，应在使用时验证；单独的开关速查更适合作为 reference。EPMoE 符号 T/E/K/H/F 和实现名称是项目约定，不应扩展成独立 skill。

## 3. 用资源下界判断性能优化是否值得继续（低频候选）

**触发场景**：优化接近平台吞吐上限，反复尝试缓存、改调度或指令替换而收益有限。

**可复用步骤**：
1. 从编译后指令统计各资源使用量，除以每周期容量，建立各类资源下界。
2. 分段查看 ramp / steady state / drain 占用，区分稳态瓶颈和启动、排空延迟。
3. 对候选优化计算资源转移：减少 load 是否增加已接近上限的向量计算，是否超出 scratch 容量。
4. 记录实测回退与约束推导，以此筛掉等价的重复尝试；只有约束改变或提出新机制时重开实验。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-tmp-humanize-demo-original-performance-takehome/memory/perf-takehome-floor.md:17–26`：记录编译指令数、每类 floor，并修正稳态为 LOAD-bound。
- 同文件 `:29–33`：“load binder 1063 sits only 8 cycles above the valu floor 1055”，减少 gather 转而增加 valu，随后列出 3 次实测回退。
- 同文件 `:38–54`：基于固定 slot 限制、非连续 gather 和 scratch 容量说明候选方案为何无效。

**频率 / 信心**：1 份项目 memory；其中记载 3 次相关失败实验，不能计为 3 个独立 session。推导结构可复用，跨项目复现证据不足。

**建议 skill description**：结合编译指令计数、资源容量与分阶段 profile 推导性能下界，评估优化的资源转移和容量约束，避免重复无效实验。

**边界及反证**：1092 cycles、2048 gathers、slot 限制和“不要重试”仅属于原项目。下界依赖当前算法、ISA 和工作量；改变算法或数据表示后必须重新推导。不要把项目收敛结论复制为通用禁止规则。

## 不应独立沉淀的内容

- `epmoe-notation.md:8–30`：固定符号及 baseline/#1058/unroll/pallas 命名，属于项目术语与沟通约定。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-medusa/memory/feedback_technical_register.md:7–30`：书面技术语体、禁用口语与第二人称，属于个人表达偏好；准确术语校对可放文档规范，但本批没有独立任务频率证据。
- XProf flag 表、takehome 精确周期与已失败参数属于参考/项目记忆，应作为上述方法的附属材料而非 skill 主体。

## 覆盖情况

- 全文读取 memory 4/4：`epmoe-notation.md`、`xprof-custom-call-tracing.md`、`feedback_technical_register.md`、`perf-takehome-floor.md`。
- 流式扫描日志 1/1：`5c668753-ed25-444d-b25f-df9808b7c1c1.jsonl` 全部行，读取全部 user 文本，并复核 assistant 中的阈值变化和最终未解决项。
- 无其他文件修改；未创建 skill。
