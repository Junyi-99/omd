分类：`accelerator-experiment-workflow`
方向：把研究问题转成可审计的加速器实验，并控制运行、证据和结论之间的落差。
领域：TPU ICI 微基准、云端 TPU/GPU 资源、实验科学性复核。
频率：1 个原始日志 session，包含 2 轮实验、3 轮不同目的的 review；另有 3 份来自不同项目的 memory。不能把轮询唤醒、任务通知、上下文摘要计为独立复现。
为什么值得持久化：实验设计、运行状态机、结论覆盖审计具有可复用输入/步骤/输出；具体硬件常数、云项目标识和临时资源状态应留在项目 memory。
范围与局限：已流式扫描全部日志 user 记录，读取全部 3 份 memory；重点核对 review 汇报与用户约束。未打开实际实验源码/CSV，日志中的测量数字和云状态不视作已独立验证。

## 1. 研究问题 → 有来源的实验矩阵 → 问题覆盖审计

- 触发：用户想知道某种硬件行为是否值得优化，或问实验结果是否回答原始问题。
- 可复用步骤：
  1. 记录研究问题、可用硬件/拓扑、workload 和估计量；先查现有微基准及模型配置。
  2. 给每个消息大小、设备映射、并发度、对照实验写 why 与证据来源；将其他平台的术语映射到当前硬件可执行原语。
  3. 为假设建立“问题—实验—指标—所需对照”表，保留用户原始问题，区分必做与超出当前范围。
  4. 运行后逐项标记 answered / partial / not answered / framing wrong；指出最小缺失测量，而非用已有曲线硬答所有问题。
  5. 当用户禁止新实验时，明确证据缺口，修正文案和判定，保留未回答项。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/7dff67a6-e2bc-469e-a770-fffc93de1f63.jsonl:190`：“先搜索一下 paper…怎么做这个通信性能实验”。
  - 同一绝对日志路径 `:399`：“写清楚 why…证据来源在哪”；`:750` 给出第二轮粒度、并发、容量与拓扑矩阵；`:806` 要求将十个问题追加到设计。
  - 同一绝对日志路径 `:1521`：“现在的实验内容和结果，解决了 design-2.html 提出的问题吗？”；`:1556` 的任务结果产出逐项覆盖矩阵，并指出缺少匹配 wait-all 基线等证据。
- 频率/信心：1 session 内 2 轮设计、1 次专门覆盖审计；方法信心高，跨 session 常见程度未证实。
- 建议 skill description：当用户设计或审查加速器性能实验时，从研究问题、文献和实际硬件能力建立带来源的实验矩阵，并逐项核对结果是否足以回答问题。
- 边界及反证：不把 GPU threads 机械等同 TPU DMA 并发；部分问题需要新 workload，更多分析不能替代缺失对照。不能将该 session 的 TPU 数值泛化成硬件定律。

## 2. 独立 review 驱动的实验结论去过度声称

- 触发：实验已有设计、代码与报告，需要审查科学性、回归修复或结果文案。
- 可复用步骤：
  1. 给独立 reviewer 划定只读范围、原始设计、实现、原始数据/分析产物与已知缺口。
  2. 分开审实现错误、计时范围、拟合可识别性、对照完整性、覆盖率及结论措辞。
  3. 输出可复核 file:line、solid 与 over-claimed 清单、修复回归状态；主审汇总并抽查引用。
  4. 将测量事实、模型推算、机制猜测分开；对不可比估计量、独立小块延迟与流内首到达等概念明确标注。
  5. 在授权范围内同步修分析生成器、报告、答案和 memory，避免仅改最终 HTML；保留未验证项。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/7dff67a6-e2bc-469e-a770-fffc93de1f63.jsonl:1033`：“让 codex…开一些…subagent…master agent 汇报”；`:1393` 再次请求全面 review；`:1521` 第三轮改审问题覆盖。
  - 同一绝对日志路径 `:1401` 将第二轮拆成回归、科学性、答案覆盖、代码/运维；`:1500` 明确“固定 N 截距不是通用常数”“独立内核不是流内到达时间”等错误。
  - 同一绝对日志路径 `:1503`：“修复 over-claimed 结论。不用补更多实验了”；`:1517` 汇报修分析器、答案、报表及 memory 的限定表述。
- 频率/信心：3 轮 review，均属同一研究 session；`:1030` 与 `:1033` 是首次需求澄清，不计两轮。方法信心高。
- 建议 skill description：当用户要求独立审查实验或修复过度声称时，核对实现、测量口径、对照与问题覆盖，给出证据定位并同步收紧派生产物中的结论。
- 边界及反证：指定模型与多 agent 是本次授权，不应固化为所有 review 的必需条件；review 本身也可能误读数字，不能不核实就批量改结论。第二轮后仍有残余过度声称，说明“review 完成”不等于科学结论全部成立。

## 3. 云端分片实验的终态驱动监控与补跑（候选）

- 触发：长时间运行的 TPU/GPU 分片实验，需要监控、回收部分结果、失败恢复与共享设备交接。
- 可复用步骤：
  1. 用项目控制器查询一次状态，记录实验 id、已完成/总 shard、资源占用。
  2. RUNNING 时只报告进度并安排下一次检查；允许获取 partial，但标明不完整。
  3. FAILED 时读取控制器日志和 stdout，定位失败 shard，只补跑失败部分，记录补跑来源。
  4. SUCCEEDED 且覆盖完整后运行分析和报表；汇报关键指标与覆盖率，再释放共享设备给下一轮。
  5. 把资源连接、配额/容量/权限、实际实验失败分层诊断；记录仍存活计费资源，不把历史状态当当前事实。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/7dff67a6-e2bc-469e-a770-fffc93de1f63.jsonl:854`：“partial 的给我拉下来看看”；`:875` 明确 SUCCEEDED/FAILED/RUNNING 分支及第一轮结束才能交接 TPU。
  - 同一绝对日志路径 `:1203` 为第二轮定义 52/52 终态分析与 failed-shard fix-up；`:1233`、`:1253`、`:1268` 是同一监控任务重复唤醒，不增加复现次数。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-1554-boundary-tpu-run.md:12` 区分 tenant 创建错误与容量；`:18` 记录等待真因是 SSH Host 配置，并指出 VM 仍计费。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/memory/gcp-b200-spot-provisioning.md:13`–`:20` 记录配额入口、容量、建机与停机陷阱，提供跨项目运维佐证。
- 频率/信心：日志内 2 个实验监控实例（1 session），另 2 份不同项目 memory 提供部分基础设施佐证；候选强度中等，不能声称完整状态机跨 3 session 复现。
- 建议 skill description：当云端分片实验需要持续监控、失败补跑或共享设备交接时，按控制器终态推进结果回收、覆盖校验、分析报告和下一轮启动。
- 边界及反证：各硬件的创建/停机参数会变，不能固化 zone 优先级、镜像版本或云项目；也不能自动删除计费实例。若已有项目控制器，skill 应编排它而非另造轮询系统。

## 不建议提炼成独立 skill

- SeleniumBase 迁移：`/Users/junyi/.claude/projects/-Users-junyi-claude-ChatGPT-API-Scanner/memory/project_seleniumbase_migration.md:7`–`:19` 仅 1 份计划 memory，无实施复现。保留为项目待办，不升级为通用浏览器绕过流程。
- 具体 B200 机型、配额/zone、TPU 数值常数及 bug 修补：事实记忆与环境快照，应即时核实；不是方法。
- 图轴改成 B/MiB/GiB（日志 `:1596`）与 TLDR/折叠详情（`:1767`）：各一次报告反馈，适合作为既有科学报告 skill 的展示建议，无独立频率支持。

## 覆盖情况

- memory 3/3 全文读取：`project_seleniumbase_migration.md`（19 行）、`kda-1554-boundary-tpu-run.md`（22 行）、`gcp-b200-spot-provisioning.md`（22 行），绝对路径见上。
- log 1/1：`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/7dff67a6-e2bc-469e-a770-fffc93de1f63.jsonl` 全文件流式扫描，244 条 type=user 记录；排除工具结果、系统注入的 Artifact 规则、任务通知、续接摘要与重复唤醒后的频率判断见各 pattern。
- 无其他 reviewer 结果读取；仅写本报告。优先沉淀 1、2，3 保留为候选或现有云实验 skill 的补充。
