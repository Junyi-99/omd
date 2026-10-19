# reviewer-30

- 分类：`evidence-driven-research`
- 方向：从原始来源构建可追溯、可更新的研究产物；其次是硬件实验与实现核验。
- 领域：TPU 推理生态、JAX 性能实验、模型实现。
- 频率：明确的方法论需求在 2 个独立 session 出现；另 1 个 session 的实现核验呈现相近证据意识，但主要来自 assistant，不计为用户主动复现。其余候选均只有单 session。
- 持久化价值：来源追踪、结论分级、采样边界与更新差异是可重复执行的方法，适合沉淀；机器名称、研究选题和默认版本偏好不适合独立做 skill。
- 审阅范围：manifest 的 3 份 memory 全文、7 份 JSONL 全部 user turns（Python 逐行解析，排除 tool_result），并补读相关 assistant 文本。未外部核验历史技术结论；不把 assistant 的执行自述直接当作已验证成功。

## 1. 有来源与覆盖边界的增量研究（建议优先沉淀）

**触发场景**：定期追踪多个项目/技术生态，或把研究结果组织成可按方法论重建的 dashboard。

**可复用步骤**：
1. 先读上一版结论、已有分类与原始增量，明确时间窗及“新增”的定义。
2. 记录来源清单、采样/全集边界；标题仅作检索入口，回到原始 PR、issue、代码或权威来源核实。
3. 为每条结论保存来源定位，区分事实、测量、信号与推断；到达可信原始来源后停止无收益深挖。
4. 将结论映射到关注的问题与已有分类，输出新增、实质变化及未发现项。
5. 显式披露未完成搜索、工具阻断和样本缺口；不能把“无法搜索”写成“没有变化”。
6. 如需长期维护，将方法论、实现和二者映射分开保存，检查重建的一致性。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/e58645b9-892d-4b0b-908f-0551ef679fcb.jsonl:5`：用户要求“每一条结论都要带 reference”，并强调“样本不是全集”“区分 fact / measurement / signal / inference”。同一条要求与上一份报告对比。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard-container/b075c613-53b9-46fc-b2aa-76e67496c85d.jsonl:9`：用户要求“提供一个 source of truth”，以及 methodology / mapping / implementation 三部分与“幂等”重建。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/e58645b9-892d-4b0b-908f-0551ef679fcb.jsonl:321`：assistant 披露搜索及写入被阻断；实际只交付聊天内报告，是应保留的失败边界。

**频率/信心**：2 个独立 session 的用户明确请求；中高信心支持共同的“证据研究”核心，低信心支持把自指 genesis 架构整体推广。

**建议 skill description**：当用户需要追踪技术生态、更新研究 dashboard 或比较前后变化时，读取历史结论与原始来源，核实引用、声明覆盖边界并交付可追溯的增量分析。

**边界及反证**：dashboard session 在第 49 行被用户中断，不能声称幂等重建已经验证。genesis 自指结构是单次设计，不应成为所有研究的固定目录。新闻追踪的多仓库清单也应由任务配置。

## 2. 隔离变量、审计计时的硬件微基准（候选）

**触发场景**：检验某个维度是否导致性能阶跃、对齐效应或缓存收益。

**可复用步骤**：定义单一计算与控制变量；设计粗扫与边界密扫；分离编译/预热/测量；按异步执行语义设置同步；记录环境、原始数据与公式；先出可复现曲线，再解释机制。若涉及缓存，则区分冷读/热读与文件大小，并检查本地介质是否成瓶颈。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/6ba43316-7f69-4dd5-b415-37b1bb9837e5.jsonl:9`：用户明确要求“compile once”“exclude compilation time”“block_until_ready() correctly”，并要求原始 CSV、环境信息及“Do not interpret the result yet”。
- 同文件 `:92`：assistant 识别逐调用同步会掩盖小矩阵效应，改为重复执行后同步；同时声明“hasn't been run on a TPU yet”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/8943c2c1-caf3-40a7-886a-19cb9cfcf3e7.jsonl:215` 用户只要求“配一下 file cache”；`:248`、`:253`、`:275` 是 assistant 主动分辨大小文件和本地带宽的测量过程。

**频率/信心**：1 个明确微基准需求 session，另 1 个配置任务内出现 assistant 主导的实验；中等方法价值、较低重复需求证据。

**建议 skill description**：当用户要设计或审计 JAX/TPU 微基准时，检查变量、异步计时、编译与预热隔离、扫参覆盖和原始数据交付，并明确实际运行环境。

**边界及反证**：不能把缓存调参自动算作用户第二次要求微基准。具体试验时长、重复数、矩阵范围不是通用默认；同步策略须核实实际执行/排队语义。CPU 检查不能替代 TPU 实测，日志中的具体性能数字也不应固化。

## 3. 用 checkpoint 小数据与独立参考核验模型实现（候选）

**触发场景**：用户想下载巨大权重来验证实现，但正确性可能只依赖少量元数据或参数。

**可复用步骤**：拆开待验证命题；定位 checkpoint 的配置、索引和必要 buffer；用文件头/范围读取缩小下载；将派生参数与 checkpoint 真值对账；找到外部参考实现作差分；对共同假设和缺失覆盖单独列风险。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/8943c2c1-caf3-40a7-886a-19cb9cfcf3e7.jsonl:280`：用户要求拉 embedding table，“测一下现在的代码，ngram 实现的是否正确”。
- 同文件 `:331`、`:359`、`:377`：assistant 将大权重需求缩成小 buffer，并通过 safetensors offset 和 HTTP range 取得对账参数。
- 同文件 `:414`：assistant 指出独立标量测试仍可能和实现“共享同一个假设”；`:424`、`:443` 转向 transformers 参考源码和差分测试。

**频率/信心**：仅 1 个独立 session，方法由 assistant 提出；中等技术启发、低重复需求信心，暂不独立建 skill。

**建议 skill description**：当用户需要核验大模型算子或 checkpoint 参数推导时，提取最小必要真值，结合权威参考作差分测试，识别自写测试共享假设造成的盲点。

**边界及反证**：只核验 hash 或布局不能证明 embedding lookup、设备执行及端到端模型正确。所读日志尚未给出差分最终结论，也不能把它描述为已通过。范围读取取决于存储格式和服务能力。

## 不建议独立沉淀的内容

- deck 发布 memory（`/Users/junyi/.claude/projects/-Users-junyi-claude-deck/memory/deck-release-lightweight-tags.md:10`、`:16`、`:18`）提供可执行发布步骤，但只有一份项目记忆。轻量 tag、直推 main 与签名配置是项目约束；可并入已有发布流程，不能泛化。
- vibe-usage memory（`/Users/junyi/.claude/projects/-Users-junyi-claude-vibe-usage/memory/MEMORY.md:1`）仅是默认 patch bump 的索引。引用文件未分配且未读取，不能据此声称有第二次同类发布流程复现。
- kernel KB memory（`/Users/junyi/.claude/projects/-Users-junyi-claude-llm-wiki/memory/paper-kb-for-kernel-gen.md:10`–`:16`）是论文选题、四轴分类及预定交付，缺实际操作证据；保留项目记忆即可。
- TPU 建机、磁盘切换、个人主页措辞和 Zed 设置均属单次操作或偏好，不足以单独建 skill。

## 覆盖情况

- Memory：deck-release-lightweight-tags.md、paper-kb-for-kernel-gen.md、vibe-usage/MEMORY.md，3/3 全文读取。
- Logs：e58645b9（追踪）、8943c2c1（TPU/GCS/核验）、6ba43316（微基准）、99f11088（主页）、b075c613（dashboard）、ecb4a99e（sglang-jax）、31b0962b（skim-cc），7/7 全量扫描 user turns。
- ecb4a99e 与 31b0962b 只有 clear/本地命令提示，无实质用户请求；不计作方法复现。主页日志仅 1 个用户请求。各候选频率不包含 memory 的未知来源次数，也不把同 session 的追问累加为独立复现。
