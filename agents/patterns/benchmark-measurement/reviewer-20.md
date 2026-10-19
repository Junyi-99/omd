# reviewer-20

- 分类：`accelerator-experiment-methodology`
- 方向：把 GPU/TPU 性能实验变成可追溯、可质疑的证据链。
- 领域：推理引擎、kernel benchmark、性能结论与测试契约审阅。
- 频率：性能实验相关的直接用户证据见 2 个独立 session（MPK、MXU）；另有 2 份相关 memory。MPK 内多次追问为同一实验的纠偏，不计独立复现。测试精简只有 1 个 session，列候选。
- 持久化价值：实验矩阵、计时口径、污染判断、原始材料和结论升级规则可以跨硬件复用；具体服务器、路径、芯片常数和模型限制不应固化进 skill。
- 范围与局限：完整读取分配的 3 份 memory，以 Python 流式扫描 4 份 JSONL 的全部 user 文本，并核实 MXU assistant 解释。任务通知/其他 agent 消息不算用户要求；原始实验产物未另行核验，不能把历史 agent 的“完成/干净”视为事实。

## 1. 可审计的加速器 benchmark：矩阵、原始 trace 与污染状态

**触发场景**：比较多个引擎/优化开关、硬件和 workload，或怀疑共享机器上的性能数据被污染。

**可复用步骤**：
1. 固定比较矩阵、软件版本、模型、精度、TP、batch/token 语义，明确 prefill 与 decode 指标；不支持的格子显式标注原因。
2. 先 smoke test 验证每条路径，再扩展 sweep。计时与 profiling 分开记录，不能直接混用两者耗时。
3. 为每个结果保存 run ID、命令、版本、时间窗、设备身份、原始日志与 trace；短生命周期机器上先保存到持久存储。
4. 运行前及运行期间记录进程/设备使用及监控心跳；同时覆盖自己的并发 worker 和外部进程。监控失败或缺失产生“未知”，不能解释为无人占用。
5. 将结果行、trace 和监控时间窗关联，区分确认污染、条件可信和证据不足；保留污染材料、列出重跑清单，避免把失败格子误认已完成。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/a349f4dd-72fd-4b8d-af26-94860a64042d.jsonl:43`：“baseline…开启 PDL…Megakernel”，并要求 batch/token 矩阵和“prefill 和 decode 过程分开测试”；49 是中断后重发，不多计。
- 同一绝对日志路径 `:211`：“需要看到完整的 profiling trace…在服务器那边保存”。
- 同一绝对日志路径 `:868`：“profiling 的时候，gpu 有别的人也在用？”；`:1077`：“检查一下我们的实验数据有没有遭到污染”。
- 同一绝对日志路径 `:1131`：历史审计通知指出“quiet log proves nothing”和缺少 heartbeat。此为 agent 报告，只用于解释可复用的失败机制，不作为用户频率或已核实审计结论。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/MEMORY.md:5`：“preserve raw bench/eval logs…not just summaries”。

**频率/信心**：1 个直接任务 session + 1 份跨项目 memory；本 session 至少 3 类明确要求（矩阵、trace、污染审计）。方法价值高，独立复现证据中等。

**建议 skill description**：当需要比较 GPU/TPU 推理或 kernel 性能、运行参数 sweep、保存 profiling 材料，或审计共享设备上的实验污染时使用；建立可复现矩阵和原始证据索引，输出可信度及重跑清单。

**边界及反证**：不把历史 MPK 的具体限制、监控间隔或硬件支持表照搬为常量。没有监控心跳不等于确认污染；同卡 MPI worker 也不自动等于两套实验冲突，需追踪进程归属。实际 trace 留存量应按预算规划。

## 2. 性能异常先立假设，再用判别实验收敛结论

**触发场景**：耗时平台、反常跳变或巨大 speedup 被直接解释成 launch/dispatch、硬件对齐或某种优化收益。

**可复用步骤**：
1. 阅读实际测量代码，列明同步点、异步提交、warmup、统计量、设备/host 时间的区别。
2. 把观察与解释分开；为 host 提交、设备启动/执行、编译变化等备选机制登记假设。
3. 设计能区分假设的小实验：提交耗时与完成耗时对照、最小计算基线、少量代表形状 trace；批量循环方案需检查编译器是否消除或改变待测工作。
4. 记录协议和漂移，给假设标 Confirmed / Rejected / Inconclusive / Not tested。未做判别实验的机制解释保留为假设。
5. 将确认后的 findings 提升到主报告；PR 承载细节和证据，面向读者的文章围绕推理主线展开。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/c9ee6696-2a91-40a3-bb70-411c0ef0aa80.jsonl:54`：用户追问“这 30 µs 是…固定开销…你确定吗？”
- 同一绝对日志路径 `:43`：assistant 原先断言是 Python dispatch，并断言两档跳动与 N 无关；`:57` 后承认“不能确定”，提出 host/设备两种机制和判别实验。这是一次过度归因的纠正，不是已完成的因果验证。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-tpu-master-report.md:13`：LEDGER 四态与测量协议；`:17-22`：“agent 输出必须先降级为 hypothesis”“新实验先记 LEDGER 再谈结论”，并分工 PR 细节与 blog reasoning。

**频率/信心**：1 个直接 session + 1 份来自另一 session 的反馈 memory（其 originSessionId 在该 memory 第 7 行）；可见 2 个独立项目的收敛需求，信心中高。没有证据证明 MXU 判别实验实际执行。

**建议 skill description**：当性能结果需要解释、优化结论有争议、或多条探索方向需要收敛成可信报告时使用；检查计时协议、登记竞争假设、设计最小判别实验，并以证据状态筛选最终 findings。

**边界及反证**：曲线平台本身不足以唯一定位固定开销来源。不能复制历史 assistant 对 JAX API、TPU 峰值、循环计时的未验证说法。单纯写 blog 不触发整个性能实验流程；报告目录命名属于项目 memory。

## 3. 候选：按行为契约精简测试，而非按行数删测试

**触发场景**：用户要求审阅测试冗余、合并参数化、改进命名和校准覆盖宣称。

**可复用步骤**：读取当前测试及必要生产路径；按契约建覆盖映射；每个删除候选指明保留覆盖的具体测试和损失；保护独立 oracle/边界/集成行为；将命名与断言实际证明的行为对齐；分明确建议和可选风格。

**原始证据**：`/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/e03018a5-a233-4513-98f7-50292c50abab.jsonl:5`：“不要仅为了减少行数删除独立 oracle、边界或关键集成行为”；要求每项给“由哪个具体测试保留相同覆盖，以及删除会失去什么”。同一 turn 要求以当前工作区为准、只读、不运行测试。

**频率/信心**：仅 1 次独立请求；契约清晰但尚不足以证明是常见工作流，候选，不建议据此单独创建 skill。

**建议 skill description**：当用户明确要求精简或审阅测试时，依据真实行为契约提出删除/合并/命名建议，并说明覆盖替代关系与损失。

**边界及反证**：只读、不运行测试是本次明确约束，不宜永久固化；不应变成泛化“多写测试”审阅，也不将当次五类契约写成通用清单。

## 不成立及覆盖

- buddy runbook：全文 1–14 行已读；本机启动方式、TCC 归属、Xcode 路径与模拟器限制是项目运行记忆。仅 1 份 memory，无独立复现证据，不升级为通用服务部署 skill。
- KDA master report：全文 1–23 行已读；可复用证据状态和收敛规则进入 pattern 2，目录与具体优化数字保留为事实记忆。
- TPU spot MEMORY：全文 1–11 行已读；原始材料保存进入 pattern 1，其余链接标题仅是索引，不据此假装读过子文档；芯片常数、机器清单、Python 偏好不成为 skill。
- MPK 日志：全部 1,483 行流式扫描，全部 user 文本已读；重发、状态追问、通知不多计频率。
- SGLang JAX 日志：全部 104 行流式扫描；单个直接审阅请求。
- MXU sweep 日志：全部 60 行流式扫描，全部 user 文本及解释更正已读。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/767b9939-e1ea-478a-8123-780a1a2600a0.jsonl`：全部 user 文本已扫描；第 12 行“找…租房的地图程序”为一次资源查找，且夹杂本地命令，无可复用 workflow 证据。
