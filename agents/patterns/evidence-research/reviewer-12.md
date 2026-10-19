# reviewer-12

- 分类：`evidence-pipeline`
- 方向：把持续观察任务变成可追溯、可复跑、带质量闸门的数据与 Agent 协作流程。
- 领域：开源推理引擎能力跟踪、研究工程数据产品；次要领域为历史会话中的项目溯源。
- 频率：分配范围只有 **1 个原始日志 session**；扫描 28 条有文本的 `type=user` 记录，其中含任务通知、skill 注入和 stop-hook 重述，不能算 28 次真实请求。另有 3 份 memory，涉及不同项目，但不能据此认定下面流程跨 session 复现。
- 持久化价值：流程包含证据、程序/判断边界、失败语义和验收链，能减少重复人工侦察；目前适合候选 skill，尚不足证明为高频通用技能。
- 范围与局限：全文读取全部 memory；Python 逐行解析原日志并扫描全部 user 文本，再回读相关 assistant 上下文。没有读取其他 reviewer 报告，也未核查历史声称已生成的仓库文件；实现成果均为会话中的报告，不等同独立代码验证。

以下将长日志路径记为 **L**：`/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl`。证据定位均为此绝对路径的 JSONL 行号。

## 1. 证据驱动的周期性数据产品：程序提取，Agent 判断，验证后发布

**触发场景**：定期追踪多个上游仓库，生成模型/算子支持矩阵、周报或 dashboard，且结果必须解释来源并支持重新运行。

**可复用步骤**：
1. 明确观察对象、数据身份和输出契约，区分具体 checkpoint、架构族和 kernel，避免把不同口径合成一个“支持数量”。
2. 保存原始来源、命令、实际版本和采集状态；固定快照日期与 SHA，完整/部分/失败分别表达。
3. 把抓取、标准化、聚合、证据提取、报告渲染拆为清楚的接口；确定性事实由程序产出。
4. 把“是否收录”“native 或 generic”等解释性交给显式判断表；Agent 读取增量和上一份报告，结论必须带出处并区分事实与推断。
5. 对路径、模型元数据、注册表、schema、链接逐项校验；查不到与查出错误分开，未知状态不能推进 verified。
6. 对齐输入快照后再生成页面，展示实际数据日期；跑完整链及契约检查，最后查看页面。
7. 文档给出“改什么动哪里”，让新增来源、调整分类、改判断分别有稳定入口。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:419`，用户：“再 collect/normalize/derive…最后 landscape + report”；另一条线“按 weekly-job.md 走一轮，verified 等检查退出码”。这是同一 session 中对两套现有流程的实际执行要求。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2269`，用户追问“一个是程序生成，一个是 Agent 自动去寻找…每一步都发生了什么…文件背后的路径”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2273`，assistant 展开 fetch → clean → weekly → html，描述原始来源、判断表、验证闸门和每层产物；属于方法执行证据，不能当成独立用户复现。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2977`，assistant 报告网络异常曾误判死链，改成 unverified 并保持失败退出码；同时承认日期身份与缓存键仍有缺陷。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2989`，用户要求简单使用文档，并“教我怎么迭代它的方法论”。

**频率/信心**：1 session；用户在执行、来源解释、方法论交接三个阶段持续要求同一链路。方法真实性高；跨任务普适性中等，频率仍低。

**建议 skill description**：当用户要建立或更新多来源的能力跟踪报告、支持矩阵或数据 dashboard 时，建立可追溯的数据流水线，分离确定性提取与人工/Agent 判断，以实际验证结果控制数据状态和发布。

**边界及反证**：不要固化四个特定引擎、目录、日期或现有 `osd` 命令；单次静态页面不需要整套流程。会话明确暴露同日覆盖、读取 current 与按 date 判完成不一致、缓存键缺字段等问题，不能把“产物存在即幂等”照搬为正确实践。机械步骤已落代码时，skill 应调用工具、解释边界，而非复制程序逻辑。反复问题进入共享约束/校验可作为此候选的迭代步骤，不另拆成已证实的高频 skill。

## 2. 从会话考古恢复项目及 prompt 的来源链

**触发场景**：用户记得成果或目标，却不记得项目目录；多个近似项目并存，且构建 prompt 可能只留在聊天记录中。

**可复用步骤**：
1. 从成果关键词搜索会话与本地目录，收集候选路径和最近工作，不因首个匹配停止。
2. 用 README、输出文件和 git remote 交叉确认；让用户对候选的排除结果约束后续搜索。
3. 区分执行 prompt、行为规则、一次性迁移规格与普通说明，不因名字含 prompt 就混成同一种来源。
4. 在文件未覆盖时回查历史用户指令；若任务授权归档，再把原文与 session、时间、产物映射一起保存。
5. 输出简短的“路径 → 用途 → 证据 → 当前对应产物”，明确不存在于仓库的部分。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:9`，用户：“可能有不止一个项目…找全…搜索一下 Claude Code 或者 Codex 的聊天记录”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:277`，用户指定 ownership/sept15，并排除其余候选；`:396` 继续要求 GitHub 仓库地址。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2286`、`:2308`，用户从“prompts 有哪些文件”追到“构建整个仓库用的 bootstrap prompt”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:2405`，assistant 报告一条线有归档 schema prompt，另一条线的建仓指令只在历史会话，抽取后以来源 session/时间标注。

**频率/信心**：1 session 内 2 个相关检索目标（定位项目、定位建仓 prompt），不算两次独立 session；用户需求明确，复用性中等，保留为候选。

**建议 skill description**：当用户按模糊成果回忆寻找本地项目、仓库或最初构建指令时，联合搜索代码与会话历史，消歧候选，并恢复可核对的来源与产物映射。

**边界及反证**：首次目录检索与后续 prompt 检索属于同一项目工作，尚无跨项目复现证据。若已有可靠项目索引，应优先查索引。历史记录中个人信息、凭据和工具注入不得作为可转存 prompt；不要默认授权复制整段日志。

## 不建议独立建 skill 的材料

- `...CODON.../memory/user_profile.md:7`：研究主题与身份事实，无方法步骤，属 profile；不复制身份细节。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/feedback-save-scripts-locally.md:11–26`：保存脚本、精确 git add、临时 worktree 执行孤儿分支操作，都是有具体事故理由的工程约束。memory 记载 3 类教训但未分配原始 session；适合项目规则或现有远程实验/安全迁移 skill 的检查项，不能据此宣布“保存脚本技能”高频成立。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard/memory/shell-cp-alias-interactive.md:11–15`：本机 shell alias 导致拷贝挂起的单例，属于环境诊断提示。可融入通用非交互脚本排障；不独立建 skill，也不推广 `cp -f` 为所有环境下可靠的绕过 alias 方法。

## 覆盖情况

- 原日志 **L**：完整流式扫描；全部 28 条 user 文本已读取（含非真实用户注入/通知）；对执行链、来源说明、prompt 考古、验收和交接的 assistant 上下文作定向核实。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-CODON-Collaboration-Optimized-by-DNA-Oriented-Networking-for-Multi-Agent-Team-Composition/memory/user_profile.md`：全文 7 行已读。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/feedback-save-scripts-locally.md`：全文 26 行已读。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard/memory/shell-cp-alias-interactive.md`：全文 15 行已读。
