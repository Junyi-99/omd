# reviewer-11

- 分类：`developer-tooling`
- 方向：从真实编译数据出发，改进开发者可视化工具的可靠性与一致性。
- 领域：React/TypeScript UI、编译器 dump、HLO/LLO 调试。
- 频率：分配到 1 个独立 session；代码清理有 2 次明确用户请求；跨视图复用与数据溯源各只有本 session 内的连续演进，不能算跨 session 复现。
- 持久化价值：可沉淀“证据检查→范围受控的实现→行为验证”的操作流程；不宜把 TPU 事实、发布授权或用户提交署名偏好直接做成 skill。以下均为候选，跨项目复用证据仍不足。
- 审阅范围与局限：全文读取 3 个 memory；以 Python 流式解析唯一 JSONL，扫描全部 user 文本，再核对相关 assistant 上下文。过滤 skill 注入、任务通知和压缩摘要，不把它们算新的用户需求。未运行历史项目、未核验历史代码当前状态；下述完成效果是日志中 assistant 的报告。

## 1. 证据驱动的 dead-code 清理

**触发场景**：用户要求删除冗余代码、无用注释，或重构后做清理。

**可复用步骤**：
1. 按组件、库、脚本和样式分区检查，并先识别项目已有格式与静态检查命令。
2. 删除 export 前检查全仓调用；收紧类型前证明调用方满足约束。
3. 对 CSS 搜索结果复查动态类名、模板拼接和共享消费者，不能仅凭没有字面引用删除。
4. 清理写入但从未读取的字段，保留有实际消费者的类型和接口。
5. 汇总改动，执行格式、lint、类型和构建检查；不把机械扫描结果直接当结论。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-visualizer/8c7cd77a-03bc-4211-9bb0-e396333ae734.jsonl:11`：“review the code…pass all formatter and linter…clean the redudant codes or comments”。
- 同一绝对文件 `:427`：“清理dead code”。
- 同一绝对文件 `:484`：报告扫描命中的 `memory-space-*`、`part-*` 等“built at runtime…live”，并删除 parser 中只写不读的字段。

**频率/信心**：2 次明确请求 / 1 session；方法信心中高，常见性证据中等。

**建议 skill description**：当用户要求清理冗余代码或重构后清扫时，结合全仓引用、动态类名与数据消费关系识别可安全删除内容，并用项目静态检查验证。

**边界及反证**：第二次清理紧接同一次侧栏重构，并非独立项目复现。不要扩成自动重构或默认提交推送；动态 CSS 是本例最明确的误删反证。

## 2. 把跨视图一致性修复落实到共享组件，并验证实际布局

**触发场景**：两个关联视图出现相同交互缺陷，或同一功能在桌面/手机入口不一致。

**可复用步骤**：
1. 查清两个视图是否复制了 shell、toggle、状态样式；将共同交互放到共享组件，内容以插槽或属性承载。
2. 拆开公共规则与视图专属断点，避免 CSS 互相泄漏。
3. 对折叠动画采集中间帧，检查宽度、文本换行和最终布局，而不只看动画结束截图。
4. 在两个视图及窄屏分别验证导航入口与当前上下文传递；发现 overflow 后实际复测。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-visualizer/8c7cd77a-03bc-4211-9bb0-e396333ae734.jsonl:193`：“侧边栏 展开收起 的 动画 不对劲”。
- 同一绝对文件 `:423`：报告“Both sidebars now share one component”；固定宽度后用中间帧确认标题不再逐帧换行，390px 首轮溢出后重新测量。
- 同一绝对文件 `:1080`：“LLO…explain EXACT INSTRUCTION，就像 HLO…一样”；`:1302` 报告复用 `GuideRow` 和解释样式。

**频率/信心**：1 个侧栏修复链 + 1 个解释 UI 对齐需求 / 1 session。共享实现方向信心中高，通用 skill 必要性中等。

**建议 skill description**：当关联视图需要一致的组件行为、响应式入口或解释交互时，提取共同结构并用多视图、多宽度和动画中间状态验证。

**边界及反证**：解释 UI 对齐不是第二次侧栏修复。用户关于“共享同一份组件”的细节仅在压缩摘要中保存，不能冒充独立原始 user turn；结论主要以直接请求和实现报告支撑。不要强迫业务不同的视图统一断点，也不要替代专门的动画 skill。

## 3. 编译产物互链先证明来源和粒度，再生成索引

**触发场景**：把上层 IR、底层程序和示例互相链接，或给非公开指令添加解释。

**可复用步骤**：
1. 检查真实产物的文件名、入口注释、操作数与元数据，明确能确定对应到哪一级。
2. 使用同次编译的上层产物校验名称、opcode、操作数，拒绝只凭名字或文件序号关联。
3. 给索引加入拓扑/示例命名空间；无法证明匹配的产物允许展示，但不提供伪精确链接。
4. 为上层节点到后端程序、后端程序回选上层节点设计双向路径，并保留已有 URL 的兼容。
5. 明确区分确定来源与推断语义；用多个实际指令类别验证 parser 和解释界面。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-visualizer/8c7cd77a-03bc-4211-9bb0-e396333ae734.jsonl:505`：“有算法能把 HLO 和 LLO 对应起来吗？”；`:554`：“先把第一层做进 visualizer”。
- 同一绝对文件 `:543`：区分“HLO 指令↔LLO 程序”的确定对应与 fusion 内部无法精确对应的粒度。
- 同一绝对文件 `:761`：报告要求名字、opcode、操作数三者匹配；发现 118 个其他 module 的程序，文件名前缀不能可靠分组。
- 同一绝对文件 `:827`：“把 2x2 的 HLO 也加进 examples”；`:1302`：解释面板明确语义为推断，并报告修复 `>`/`>>` bundle 被当续行吞掉的问题。

**频率/信心**：1 个互链实现链，后续 2 次相邻扩展 / 1 session；技术方法信心高，独立复现频率低。

**建议 skill description**：当开发者工具需要关联多层编译产物时，先核对真实 provenance 和可证明的映射粒度，再构建带命名空间的索引、双向导航和无法匹配时的降级。

**边界及反证**：TPU opcode、具体拓扑及数值是项目知识，不应硬编码成通用 skill。缺少血缘信息时不能把启发式解释包装成精确映射；仅凭本次项目不足以断言高频。

## 不建议提升为 skill 的内容

- `project_submission_plan.md:7–11`：EMNLP 优先、KDD 后备及模板要求，是一次项目投稿计划。
- `feedback-no-claude-session-trailer.md:11–15`：不附私人 session trailer，是提交偏好；适合保留 memory，不另建方法 skill，也不推广其中对其他指令优先级的说法。
- `MEMORY.md:3`：交互式 `cp` alias 的规避，是环境注意事项；仅索引一句话，未分配其链接正文，无足够证据扩写。
- 多次“push”是同 session 的发布指令，不代表可跨 session 自动发布的授权。

## 覆盖情况

- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-CODON-Collaboration-Optimized-by-DNA-Oriented-Networking-for-Multi-Agent-Team-Composition/memory/project_submission_plan.md`（11 行）。
- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/feedback-no-claude-session-trailer.md`（15 行）。
- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard/memory/MEMORY.md`（3 行）。
- 全 JSONL 流式扫描：`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-visualizer/8c7cd77a-03bc-4211-9bb0-e396333ae734.jsonl`；扫描全部 user 文本，重点回读清理、侧栏、互链和指令解释的 assistant 上下文。
