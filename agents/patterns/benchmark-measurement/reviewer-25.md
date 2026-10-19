# reviewer-25

- 分类：`accelerator-benchmarking`
- 方向：把加速器微基准的计时、trace 与编译结果核对做成可复用诊断流程。
- 领域：JAX / TPU 性能实验；次要候选为基础设施预检与仓库侦查清单。
- 频率：主模式在 1 个独立 session 中反复出现，至少 4 组计时纠偏主题；不能当作 4 次跨会话复现。其余 4 个日志各 1 个 session，3 份 memory 独立全文审阅。
- 持久化价值：测量边界、异步派发与编译优化会改变实验解释，流程可跨 shape 和 kernel 复用；具体测量数字和机器地址不应写进 skill。
- 范围与局限：已流式扫描全部 5 个 JSONL 的全部 user 文本，并按主题核对 assistant 上下文；部分日志只有继续任务的一句话、图片或不可访问链接，没有完整起因。以下只有第 1 项有较充分的流程证据，另外两项保留候选。

## 1. 用 trace 与 HLO 审核 JAX 微基准的测量语义

**触发场景**：用户要求测 device time、wall-clock 和设备耗时相差很大，或希望用设备内循环摊薄 host 开销。

**可复用步骤**：
1. 明确目标是独立 kernel 调用、包含内存搬运的程序，还是片上重复计算。
2. 输入先就绪；编译和热身移到计时区间外；在正确的异步边界同步，并分别报告 wall-clock 和 device time。
3. 为 shape / loop 参数加 trace 分组和 op 名；抓 trace 时确保设备执行在 capture 结束前完成；导出可供用户打开的原始 trace。
4. 用 trace 的 host dispatch 次数、设备程序次数和间隔解释差异；不要把 host 循环误认为一次设备派发。
5. 改用 `fori_loop` 时建立真实的数据依赖，检查 HLO 确认运算未被提升；核对 VMEM/HBM 驻留、拷贝和循环边界，防止把内存场景变化误称为单纯消除了计时开销。
6. 复用公共 benchmark/profile 驱动；尊重运行范围，若用户限定本地 CPU 验证，不继续远程实验。

**原始证据**（同一日志，不计独立复现）：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/6bdc8a1a-301f-4d3d-9bb6-2bdb1171d0b7.jsonl:46`：用户问数据是否先传到设备；`:67` 提供含 `device_put` / `block_until_ready` 的计时例子；`:74` 要求同步后的 profiler capture。
- 同一绝对路径 `:175`：要求“记录下 device time”“trace 拉下来”；`:284`：追问小 N 与大 N 的 wall-clock/device-time 差距；`:291`：“每次调用都是 host 在发吗？”
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/6bdc8a1a-301f-4d3d-9bb6-2bdb1171d0b7.jsonl:303`：“用 fori_loop 改一下”；`:325` assistant 报告检查 HLO，指出循环版改变内存驻留；`:390` 用户再次追问 10 次循环为何快，`:398` assistant 给出拷入/循环/写回时间线。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/6bdc8a1a-301f-4d3d-9bb6-2bdb1171d0b7.jsonl:429`：要求公共函数提取到 `tools/`；`:747` 再次指出重复 profile/benchmark，且明确“先只做代码的更改和本地用 cpu 进行验证”。

**频率/信心**：1 session，输入/同步、trace 采集、派发差异、循环改变实验语义至少 4 组连续问题；流程价值高，跨会话常见性中低。建议作为领域型 skill 候选，待其它 reviewer 的独立证据验证频率。

**建议 skill description**：当审核或实现 JAX/TPU 微基准、解释 host 与 device 时间差异、采集 XProf trace 或将循环移入 JIT 时，先确定测量边界，再用同步、trace 和编译 IR 验证实验语义。

**边界及反证**：日志里的具体耗时、性能阈值和 profiler 选项是特定硬件/版本的结果，不能作为通用常量。CPU smoke test 不能证明 TPU 性能；循环内数值溢出在该实验中未单独验证。`:751` assistant 在“不要去机器上跑”后仍提出远程停止任务，提示应把用户运行边界作为检查点，而非复制历史全部行为。

## 2. 起实验前检查控制面可达性与共享资源（候选）

**触发场景**：spot VM 显示 READY 但 controller 反复重试，或希望使用共享 GPU 机器开展实验。

**可复用步骤**：先区分容量/抢占与 SSH 配置故障；检查 controller 实际使用的 alias 及 Host 块；将缺失配置补齐做成幂等入口；再只读检查设备拓扑、互联链路、负载与空闲情况，记录可复现环境，避免忙碌节点直接开跑。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/tpu-spot-ssh-host-block.md:11`–14：缺 Host 块使 SSH probe 永远失败，表现似 spot 不稳定；`:29`–34 记载幂等补块和排查日志的方法。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/lab-gpu-machines.md:8`：“shared multi-user boxes, check load first”；`:17`–19 记录互联不通与忙碌设备跳过后再跑。

**频率/信心**：2 份 memory、2 类基础设施场景；前者记载 3 个实验受同一故障影响，不算 3 个独立 session。原始修复/探测 session 未分配，跨场景统一流程的信心中低。

**建议 skill description**：在远程 spot 或共享 GPU 环境启动实验前，验证 SSH alias、控制器可达性、硬件互联与负载，区分配置故障和资源不可用，并输出可执行的预检结论。

**边界及反证**：两份 memory 主要是已解决缺陷与机器清单，适合保留为项目记忆；当前证据不足以单独创建通用 skill。SSH 模板与特定内网状态不能泛化，更不能照搬关闭主机校验的配置。

## 3. 建立有核验依据的仓库侦查清单（候选）

**触发场景**：为一个技术领域建立后续监测的目标仓库 YAML。

**可复用步骤**：以用户指定仓库为种子扩展相关基础设施；核验 URL、默认分支与关注路径；区分迁移、归档和镜像；把抓取所需字段与解释性 metadata 分开；解析 YAML 并校验关联引用，尚未启动巡检时不提前堆积运行状态字段。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/0ed432cf-36cf-4184-a9a5-04e157d887e2.jsonl:9`：要求专业的 `repos.yaml` 并扩展 TPU 相关库。
- 同一绝对路径 `:103`：assistant 报告逐一 HTTP 核验路径、处理迁移/归档；`:107` 用户删除多余字段；`:120` 将 `tags, why, related` 收入 `metadata`；`:132` 校验关联无悬挂引用。

**频率/信心**：1 session，多个 schema 修改属于同次设计迭代；中等流程信心、低重复性信心。

**建议 skill description**：当创建技术生态的仓库侦查清单时，核验仓库现状与关注路径，设计精简的 YAML schema，分离说明性 metadata 并校验交叉引用。

**边界及反证**：具体字段和全英文要求仍是此任务选择，不宜提升为所有项目规范；没有定期巡检的执行证据，不能声称已形成持续监测工作流。

## 未升格的内容

- PR body 简短、不含 session 链接、希望 Mermaid：`/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/memory/pr-body-style.md:8`–11 是明确风格偏好，保留反馈 memory 即可，不是独立方法 skill。
- 会议网站互换两个报告：`/Users/junyi/.claude/projects/-Users-junyi-claude-icpp2026-github-io/e7aefb0f-4e42-4373-bf90-fb4cfb1a38c4.jsonl:9,134,146` 是单次编辑、看 diff、再发布；没有足够反复证据支持会议排程 skill。
- 继续爬数据：`/Users/junyi/.claude/projects/-Users-junyi-claude-chagee/36c5016c-b86d-4537-9ed1-ecb888373118.jsonl:9,30,84,101` 显示先查既有进程、确认入库并处理重复进程的有用片段，但只有一次含糊请求，缺少独立复现，暂不升格。
- Artifact finding 请求：`/Users/junyi/.claude/projects/-Users-junyi/6302853c-8676-49ac-ba87-1846c0091ad6.jsonl:9,27` 因会话额度终止，无分析产出，不能提炼方法；未复制链接查询参数。

## 覆盖情况

- Memory：`pr-body-style.md`、`tpu-spot-ssh-host-block.md`、`lab-gpu-machines.md` 全文读取，3/3。
- Logs：manifest 所列 `6bdc8a1a…`、`e7aefb0f…`、`0ed432cf…`、`36c5016c…`、`6302853c…` 全部逐行 JSON 解析并扫描全部 user 文本，5/5；排除工具结果、命令提示与后台通知，不将其算作用户需求次数；按上述主题补读 assistant 文本核对上下文。
