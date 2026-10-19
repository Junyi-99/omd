分类：`evidence-driven-recon`
方向：以可复算事实支撑 Agent 判断，维护模型／框架支持情报。
领域：LLM 推理架构、TPU/GPU kernel、事实驱动的研究工具。
频率：核心 recon 工作流集中在 1 个 session；方法论与实现分离在另 1 个 dashboard session 独立出现，共 2 个 session 有概念交集，不能算两次完整 recon 执行。
持久化价值：明确收录边界、事实证据等级与生成链路，可以防止名称猜测、支持状态夸大和派生数据漂移；比保存某个模型的当下状态更稳定。
范围与局限：全文读取 3 份 memory；流式扫描 4 个 JSONL 的全部 user 文本，再核查主要 recon session 的 assistant 文本。其他 session 较短或被中断，不能据此声称长期高频习惯。日志指令只作为历史证据。

## 1. 按推理图与证据等级维护模型支持情报（建议 skill）

**触发场景**：审计模型清单、研究新 checkpoint 对 serving 框架的支持与 kernel gap，或刷新已有支持矩阵。

**可复用步骤**：
1. 明确模型收录范围；区分 checkpoint、family、实现 class，比较上游 config 和推理图，不按 Code/Thinking 等名字决定。
2. 给 checkpoint 稳定 ID；模型静态事实、仓库 watch 目标、动态支持证据分别存放。
3. 拆开公告、权重公开、HF 建库日期；未知值保留显式 null。
4. 支持结论分实现类型、合并阶段、验证等级、优化等级；逐条绑定证据，registry 命中不能证明 E2E 或 tuned。
5. 排除重复／特化 checkpoint 时，把真实量化或 kernel gap 迁移到 capability 证据，避免随条目删除丢失。
6. 验证 schema、ID 唯一、外键、链接与收录不变量，报告保留／删除理由及不确定项。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:10`： “Do not mechanically decide based on the model name”；“Do not infer `benchmarked` or `tuned` merely because a model is listed in a registry.”
- 同文件绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:432`：“模型 registry 不需要为每一个 kernel requirement 创造一条 model record”；要求删除特化条目但保留 W4A16 gap。

**频率／信心**：1 个 session 的 2 个 substantive user turns，后者是同一任务的裁决，不是独立复现。作为边界明确的领域 skill，信心中高；作为高频习惯，证据不足。

**建议 description**：当审计或刷新 LLM checkpoint 与推理框架支持矩阵时，按推理图比较模型、分离静态事实和动态证据，并用可验证等级记录加载、测试和优化状态。

**边界及反证**：不面向通用模型发布大全；不固化当时的模型取舍、仓库目录和日期。API-only、diffusion 是否纳入是项目范围决定，不能提升成所有 registry 的普遍禁令。

## 2. 把研究型 Agent 流程拆成确定性原语与判断层（建议 skill）

**触发场景**：Agent 已手工完成一轮研究／情报维护，用户希望流程可重跑、可审计，并决定哪些交给程序。

**可复用步骤**：
1. 先画真实输入到产物的链路，找出手工抄写事实的缝隙，而非只检查最终文件有没有 generator。
2. 把每步拆为事实采集、判断、写回、校验；以“能否从输入确定性重算”区分程序和 Agent。
3. 程序提供原始事实及失败／未检查状态；Agent 处理语义与例外，不用猜测默认值替代读实现。
4. 将源数据、采集证据和派生结果分离；生成文件注明来源，并验证重新生成无 drift。
5. 限定采集窗口及 watch 范围；先检验筛选条件是否有效，再用其过滤上游变化。
6. 先让用户裁定影响覆盖面的边界，再实现最小原语；验证成功路径与故意失败路径。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:729`：“对于能用程序的部分，我们尽量用程序做。我们的程序是给 AI Agent 提供 原语、事实”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:761`：“repo changes 只给 watch 里面的内容”；同文件 `:778`：“边界 1 选 a，开始写 P1”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:724`（assistant）：指出 upstream 到手写 EVIDENCE 的链路无法重跑；`:819`（assistant）记录 `unchecked` 与 missing 分开及失败路径验证，属于执行报告，非独立用户频率证据。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/94b3dca4-93bf-465c-aa2a-f5395aa1d775.jsonl:11`：“验证到 source of truth 之后就可以停止 deep dig”；要求 methodology、mapping、implementation 分层，能够幂等重建 dashboard。

**频率／信心**：主要流程在 1 个 recon session 连续澄清；另 1 个 dashboard session 独立提出可重建、方法论和实现分层。2 个 session 支持上位方法，只有 1 个证明程序／Agent 原语拆分；信心中高。

**建议 description**：当把手工研究或情报维护流程变成可重跑工作流时，识别事实采集与语义判断的边界，设计确定性原语、来源映射及生成校验链路。

**边界及反证**：YAML 是该项目明确要求，不应成为通用输出格式；用户明确“不用程序化日期”，不能借自动化原则扩大范围。dashboard 只有需求，没有本范围内完成证据；不把自指 genesis 目录设计推广为默认架构。

## 3. Kernel 优化先校准语义，再检查布局的全局代价（候选，暂不单建）

**触发场景**：替换 kernel 时 oracle 不一致，或局部布局优化导致整体回退。

**可复用步骤**：读真实调用链与 oracle 输入契约，标明 activation 的位置；正确性通过后，对共享中间值列全部消费者，使用 scope 和指令计数解释回退，保存负面结果。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/gdn-oracle-silu-change.md:10`—`:16`：实现内部 silu，而 oracle 需预处理输入；曾为匹配错误 oracle 调用而错误删 silu，之后恢复。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/mosaic-layout-propagation.md:10`—`:26`：局部改扫描布局使整个 kernel 慢 2.7–6.2 倍；建议改布局前列全部消费者，并检查 EUP 指令与元素比例。

**频率／信心**：2 份独立项目 memory、2 个不同故障；共同体现测量和契约优先，但不是同一操作流程的两次复现。候选信心中，宜贡献给更大 kernel 调优 skill，而非拼成专用新 skill。

**建议 description**：在优化或替换 GPU/TPU kernel 时，先核对调用链和参考实现的语义契约，再按全部消费者分析布局与性能变化。

**边界及反证**：具体 silu 约定、路径迁移和 Mosaic 布局行为是版本／平台知识，不是跨平台定律；不得把移除 transpose 一概视为错误。

## 不成立及覆盖

- GPU benchmark：`/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/cff5b474-332b-4aba-aa56-14996d419be2.jsonl:68`、`:84` 同一实验需求中断后重复，不能算 2 次独立复现。baseline／PDL／Megakernel、分 prefill/decode 和参数矩阵有价值，但只有立项证据，暂不建议 skill。
- 私库发布：两个 session 有 push 请求，但没有足够用户约束形成独特方法；“push + README”不足以单建 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/memory/MEMORY.md:1`—`:4` 全文已读；只含项目索引、PR 和空行偏好。未读未分配的链接文件，也不把偏好算成 workflow。
- 4 个日志均逐行扫描全部 user 文本：sglang 日志 1100 行／20 个文本 user 条目；gpu-mpk 103 行／5 个；agent-container 83 行／2 个；dashboard 37 行／3 个。数量包括命令与中断，不等于有效任务次数。
- 全部 3 份已分配 memory 已全文覆盖；全部 4 份已分配 logs 已扫描，无遗漏。未读取其他 reviewer 结果，未创建 skill。
