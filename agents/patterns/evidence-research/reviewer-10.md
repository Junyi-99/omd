分类：`evidence-driven-research`
方向：把仓库研究变成可复跑、可核查的证据管线。
领域：开源推理框架观察、TPU kernel 机会分析、远程实验产物管理。
频率：1 个完整日志 session + 3 份 memory（各指向不同 originSessionId，共 4 个来源 session）；跨项目重复最明确的是产物持久化，其余是单 session 多轮纠错，不算独立复现。
值得持久化的原因：可抽出数据覆盖验证、身份与证据审计、产物交付等稳定步骤，能减少“程序生成所以可信”“有同名模块所以已实现”的系统性误判。
审阅范围和局限：流式扫描指定 JSONL 全文件，提取全部 user text；排除工具结果、通知、注入的 skill 正文和压缩摘要作为用户需求计数，并按需核查 assistant 上下文。memory 三份全文读取。未访问研究项目代码或外部 ground truth；日志中 assistant 声称的修复和统计仅作上下文，不当作独立验证结果。

## 1. 可复跑的研究／实验产物包（推荐）

**触发场景**：远端实验完成、生成研究报告／artifact、交付可定期重跑的观察任务。

**可复用步骤**：
1. 先设计每轮独立的目录，区分原始数据、日志、脚本、报告、图表和索引；报告使用相对链接。
2. 把采集、分析、渲染脚本保存进仓库，给出明确的一条重跑命令，避免只有临时脚本或手工写表。
3. 数字从原始数据程序化生成；自动化可复现与结论真实性分别验收。
4. 远端任务完成或失败后同步数据与日志；保留本地可分析副本。
5. 周期性 agent 研究保存输入 prompt、带引用的报告和日期快照，读取先前报告后处理新增内容。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-CODON-Collaboration-Optimized-by-DNA-Oriented-Networking-for-Multi-Agent-Team-Composition/memory/feedback_sync_after_experiment.md:7–18`：“跑完实验，及时同步”；明确完成／失败才同步，进度查询不必同步。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/feedback-experiment-folders.md:11–18`：“一轮实验的所有数据、日志、图表、脚本、原始 capture 单独放一个目录”；附 README、程序生成报告和相对路径约定。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2572`：“下次我运行什么命令可以得到这个 artifact？”
- 同日志绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2770`：“跑的脚本为什么不落盘呢？…放到 repo 里”；`:3299`：“下次我重跑要能跑出同样 artifact”；`:3767`：每日 snapshot 保存 prompt、markdown report，先写 folder layout。

**频率／信心**：3 个来源 session（2 份 memory + 1 日志）。日志内至少 4 个明确相关用户回合（2572、2770、3299、3767），属于同项目细化；跨项目可复用性高。每日 agent 快照仅 1 次明确需求，不单独包装成高频 skill。

**建议 skill description**：当研究或远程实验需要交付可重跑结果时，建立独立产物包，保存数据、脚本、日志和引用，完成后同步远端结果，并提供重建报告的命令。

**边界及反证**：远端主机、绝对目录与 rsync 命令是项目配置；不固化为通用默认。不是所有临时问题都要建实验包；缓存／增量是用户在同日志 3258 要求的效率约束，不能把“可复跑”解释为每次全量重采。

## 2. 仓库数据到支持矩阵的证据审计（推荐候选）

**触发场景**：采集 GitHub 数据、汇总多仓模型／kernel 支持情况、生成首次落地时间线或跨仓对比。

**可复用步骤**：
1. 声明时间窗、全量分母、分页及窗外仍开放条目的口径，实际验证 manifest 承诺；缺项明确列出。
2. 保留 body、讨论及事件等论证所需原始资料；已有 git 历史可取的 patch 不必重复物化。
3. 用真实架构注册、准确标识符和经核实别名连接实体；不可把不同来源的真实片段拼成不存在的身份。
4. 检查子串碰撞、共文件多架构、平台不相容、迁移改名；保留正反例，兼顾误报和漏报。
5. 分开记录 Fact、Measurement、Signal、Inference、Assumption；日期代理量要标明来源，标题提及不等于支持。
6. 数据到 UI 的最后一层也审计：unknown、仅提及、已落地但缺路径证据不能都显示为勾。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:261`：“收集了 PR 600 条，但一共多少呢？”；`:322` 要求补全；`:392` 追问缺失 body、conversation/action history；`:526` 撤回补 patches，因为可从 git 获取。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:1165`：“所有重要事实必须能够回到 source”，并明确区分五种证据类型、收集 negative evidence。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2239` 质疑 Mamba 身份；`:2294`：“模型不是用 python 生成出来的吗？为什么会出现拼出来的假身份？”；assistant `:2297` 给出拼接假标识符、错误连接键和日期代理量三类根因。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2918` 质疑 hca 搜索出现 attention；`:2946` 质疑 TPU 支持 CUTLASS；`:2974`：“检查一下还有哪些词表误伤”；`:3016` 要求区分真落地缺路径证据与 CUTLASS 误判。assistant `:2953` 明确渲染层把 unknown 画成支持。

**频率／信心**：仅 1 独立 session；有覆盖度、实体身份、词表／展示三组可核实纠错链，不能算三次独立复现。方法明确，泛化价值中高，长期频率证据中等。

**建议 skill description**：当根据多仓代码和 GitHub 历史生成支持矩阵或时间线时，审计采集覆盖度、实体连接、时间代理量和证据等级，追溯每个重要声明并验证展示语义。

**边界及反证**：不把 PR 标题零命中解释为不支持；不把程序确定性当真实性。526 的用户撤回说明“完整”按可追溯性验收，不能机械要求所有 patch 重复保存。

## 3. 对照 ground truth 校准机会发现与证伪（窄域候选）

**触发场景**：生成 kernel 优先级列表，或已有候选列表与人工认可样本差距很大。

**可复用步骤**：
1. 明确候选回答的是缺功能、现有实现优化，还是多 op 融合；先分类型再设“已经有了”的反证条件。
2. 将候选和认可样本逐项对齐，分类统计未命中、被错误过滤、实体类型缺失、证据无法取得。
3. 功能缺口检查等价实现与别名；优化候选检查具体后端／实现方式；融合候选看调用链，而非只枚举目录实体。
4. 执行反证时保留剔除理由与证据，词边界和符号搜索都要审查。
5. 显式 roadmap 与推断分层排序；需要 profile 的候选保留待测状态，不能把静态发现直接当性能收益。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:3525` 要求 P0–P3 top 10；`:3665` 明确本项目 roadmap 优先；`:3695`：“根据‘什么事实会推翻’来检查…去掉已被推翻的 target”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:4330`：以给定 target list 为 ground truth，问“有哪些是我们没获取到的？为什么…怎么才能…一致”；`:4400`：“实体层里压根不存在的内容，怎么知道要做呢？”；`:4443`：“合进 targets 做成三类候选并排”。
- 同日志绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:4396`（assistant 上下文）承认把“有模块”误当“有 Pallas kernel”，因此删掉认可候选；`:4438` 讨论调用链实体和必须 profile 的缺口。

**频率／信心**：1 session 中 1 条连续设计迭代，不能声称跨项目复现。候选级；方法有价值，但实现诊断和 ground truth 覆盖数未在本审阅中独立验证。

**建议 skill description**：当筛选 TPU kernel 机会或校准候选列表时，按功能缺口、实现优化和融合机会分类，用认可样本定位漏检原因，执行分类反证并列明 profiling 缺口。

**边界及反证**：3594 的“只做没有的”被后续 ground truth 讨论细化；不固化成永久排除已有模块的规则。AST 中出现连续调用只是候选线索，不能证明融合合法或更快；优先级规则属于具体项目。

## 不成立／不单独建 skill

- tmux 自定义术语：`/Users/junyi/.claude/projects/-Users-junyi-claude-opensessions/memory/feedback_tmux_sidebar_terminology.md:7–15` 是单项目 UI 术语偏好，适合项目 memory／词汇表；无独立工作流复现证据。
- 颜色互换、页面排序、文案删改、push、状态追问均不构成独立方法。
- 570 的“代码已定义的内容不在文档重复”有维护价值，但仅一次，未独立列为 skill。

## 覆盖情况

- memory：`feedback_sync_after_experiment.md`、`feedback-experiment-folders.md`、`feedback_tmux_sidebar_terminology.md`：3/3 全文读取。
- logs：`5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl`：1/1 全文件流式解析，扫描至末尾 user turn 4704；关键纠错回合补读 assistant 原始上下文。
- 未读其他 reviewer 结果；只写本文件，未创建 skill。
