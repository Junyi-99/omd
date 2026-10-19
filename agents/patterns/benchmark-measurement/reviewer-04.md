# reviewer-04

- 分类：`ml-performance-benchmarking`
- 方向：可审计的跨机器推理基准、性能归因与结果复核。
- 领域：GPU 推理、CUDA Graph、HTTP 延迟、可抢占计算环境。
- 频率：1 个分配日志 session，含同一 Laya 实验的多轮扩展和纠错；不能把五台机器、四个变体或反复问进度计作独立复现。另有 4 份 memory，只有 1 份包含另一项目的方法性经验。
- 持久化价值：测量口径、消融控制、原始 trace 和代码来源记录可复用于新模型；具体硬件排名、项目路径和写作偏好不应变成 skill。
- 审阅范围和局限：全文读取 4 份 memory；用 Python 流式解析完整 JSONL，扫描全部 user 文本，再剔除自动任务通知、命令包装和注入的历史 skill 文本，核实关键 assistant 上下文。以下是高价值候选，尚不足证明跨 session 高频使用。

## 1. 候选：正确性先行的跨 GPU 推理基准

**触发场景**：需要比较同一模型在多种 GPU 上的函数调用和 HTTP 服务延迟。

**可复用步骤**：
1. 固定模型 revision、上游代码、依赖和统一 harness，在每份结果中写入实际运行代码的 hash；先提交再分发，避免 rsync 未提交代码失去来源。
2. 先验证 checkpoint、输出标签与数值一致性，再测性能；记录 GPU 选择、CPU 负载和其他任务占用。
3. 区分进程内调用与客户端请求到响应的全链路延迟；固定预热、调用次数、输入形状和真实测试输入，保存每次样本，核查缓存路径。
4. 输出按变体 × 主机的 p50/p95 表格、环境信息与图；争用或异常数据标注条件，重跑后保留替换依据。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/eba5c113-e6ae-4eb2-8f61-07795e767711.jsonl:47`：“先检测正确性”“warmup 20 次”“记录每一次的 Latency”“五台机器的测试代码应该保持一致”。
- 同绝对文件 `:731`：“每行是跑的变体，每列是每台机器……HTTP p50/p95”。
- 同绝对文件 `:1347`：“51.4 ms 就不对劲。重跑”。
- 同绝对文件 `:1464`：历史审查发现 fast/profile 结果未记录 harness sha、运行过未提交代码，以及空闲 GPU 不等于整机无争用。这是执行反证，不是用户新增需求。

**频率/信心**：1 session；至少 4 次明确相关要求，属于同一实验推进。方法价值高，独立复现信心低。

**建议 skill description**：当用户要求跨 GPU 或跨机器比较模型推理速度时，建立统一且可追溯的正确性与延迟基准，保存逐次样本，检查缓存和资源争用，并生成可复核报告。

**边界及反证**：预热 20 次、测试 100 次是本次规格，不能成为所有模型的默认充分样本量。历史结果存在来源和争用问题，不能把当时脚本视为已验证模板。HTML 浅色和 SVG 是交付偏好，不是基准方法的核心。

## 2. 候选：用消融与 trace 验证优化收益

**触发场景**：GPU 性能排名反常，或需要解释融合 kernel、编译和 CUDA Graph 各自带来的收益。

**可复用步骤**：
1. 先列出变体实际改变的维度，分开 kernel 后端、权重精度、图捕获和输入批量，增加单因素消融。
2. 对每个变体验证实际启用的 backend、graph 捕获/回放和输出 parity；分开冷启动与稳态成本。
3. 同时保存未 profile 的 wall 延迟、host/device 时间、launch 数和完整 trace，避免把 profiler 开销当成真实请求延迟。
4. 覆盖单问和多问，解释 host 开销与 device 计算的相对变化；将机制推测和实测结论明确区分。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/eba5c113-e6ae-4eb2-8f61-07795e767711.jsonl:391`：“用 profiler 去测 GPU 那边的耗时”。
- 同绝对文件 `:548`：“你跑了吗就说有效？”；`:550` assistant 承认只是从 profile 推断，收益和数值一致性待测。
- 同绝对文件 `:702`：“cudagraph 在哪”；`:708`：“tilelang 不录 graph 的也跑一下”。
- 同绝对文件 `:812`：追加 multilingual 50 问 p50；`:932`：“device time 和 host time 的 decompose”“人工review trace”。

**频率/信心**：1 session；6 处明确要求反复完善同一归因链。单 session 内需求强，跨项目频率未验证。

**建议 skill description**：当用户需要解释 GPU 推理优化为何有效或出现性能反常时，通过后端与 CUDA Graph 消融、输出一致性验证、host/device 分解和完整 trace，形成有实测支撑的性能结论。

**边界及反证**：kernel 时间求和除以 wall 不自动等于 GPU 利用率；CPU self time 可能大于单次 wall。历史 assistant 曾把机制推测提前写成优化有效，skill 应防止复述这类结论。此候选可并入上一项的诊断阶段，不必另建高度重叠的 skill。

## 3. 弱候选：短可抢占窗口中的实验恢复

**触发场景**：spot 机器频繁回收，环境安装和模型下载挤占实验窗口。

**可复用步骤**：缓存环境和模型至对象存储，恢复后校验并对坏缓存回退安装；在已验证配额区域间重试；用可靠的会话托管远端任务；把恢复过程视为状态机，避免中途打断删除再创建；先用小规模评估检查信号。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/project_spot_v6e_churn_caching.md:17`：“preempt every ~5–13 min”；`:19-22`：缓存恢复校验，将准备时间降至约 1–2 分钟。
- 同绝对文件 `:23-30`：tmux、进程字符串自匹配导致自杀、recover 删除阶段不可直接终止等经验。

**频率/信心**：仅 1 份 memory，记载 2026-06-01/02 的经验；不能据此断言 2 个 session。中等方法信心，低频率信心。

**建议 skill description**：当可抢占实例的生命周期短于安装和实验准备时，建立可校验缓存、区域回退和不中断恢复状态的实验启动流程。

**边界及反证**：具体区域配额、bucket 和项目限制应留在项目配置；不能推广为所有 SSH 环境都必须用 tmux，也不能把小样本评估当作最终结论。当前仅列候选，宜等待第二个独立项目复现。

## 不建议单独沉淀

- `feedback-plain-language.md:10-14` 与 `feedback_writing_style.md:7-12` 跨两个来源共同支持“朴素直接”的表达偏好，但没有独立操作流程，保留为 feedback 即可。
- `MEMORY.md:1` 是找房条件索引，属于个人事实，不扩读未分配链接，不建 skill。
- 日志 `:583` 的 SSH socket 复用只有一次指令，适合作为远端执行技能的一条措施。
- 日志 `:1404` 将“做完”与“做对”交给独立复核有价值，但仅一次显式请求，宜纳入候选 1 的结束检查，不单独建立多代理 skill。

## 覆盖情况

- memory：`/Users/junyi/.claude/projects/-Users-junyi-Desktop-Notes/memory/feedback-plain-language.md`，全文 14 行。
- memory：`/Users/junyi/.claude/projects/-Users-junyi-claude-junyi-blog/memory/feedback_writing_style.md`，全文 12 行。
- memory：`/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/project_spot_v6e_churn_caching.md`，全文 34 行。
- memory：`/Users/junyi/.claude/projects/-Users-junyi-claude/memory/MEMORY.md`，全文 1 行。
- logs：`/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/eba5c113-e6ae-4eb2-8f61-07795e767711.jsonl`，完整流式扫描所有 user turns，关键区间回读上下文；自动通知不计用户需求或独立 session。
