# reviewer-05

- 分类：`compiler-validation`
- 方向：从真实编译产物建立可复现的可视化语义验证；次要候选为面向使用者的项目 README。
- 领域：JAX / TPU / XLA HLO 工具开发。
- 频率：仅 1 个独立日志 session；语义验证有 4 次相关用户请求，README 文案有 5 次相关请求，均是同一项目内迭代，不能计为跨项目复现。
- 持久化价值：编译器输出充当外部判据、逐片段覆盖检查、反向验证及多视图共享语义，是可迁移的方法；TPU 编号、部署地址和个人偏好本身不应成为 skill。
- 范围与局限：全文读取分配的 4 个 memory；Python 逐行解析完整 2021 行 JSONL，抽取全部 27 条非空 user 文本，再按主题核实 assistant 上下文。两条超长文本是注入的既有 skill 内容，不计为用户工作流证据。以下测试结果是历史 assistant 的报告，未重跑，也未独立验证硬件结论。

## 1. 用真实编译语料与外部判据验证解释器（优先候选）

**触发场景**：开发 HLO 等编译产物的解析器、解释器或图形查看器，需要证明每段解释与实际执行语义一致。

**可复用步骤**：
1. 写出覆盖控制流、调用、tuple、布局、特殊算子等边界的最小源程序，在目标后端生成编译前后产物。
2. 同次编译保存 HLO、调试 metadata 和 buffer assignment；避免跨编译进程用名称拼接“标准答案”。
3. 检查文本重建与原文一致、片段无遗漏、片段与说明双向对应；用已知源码调用链校验 metadata 解析。
4. 将解析结论与编译器分配结果双向比对：显示的结论必须有证据，应显示的项目不能漏掉。
5. 遍历示例中的实际 UI 节点；刻意恢复一个旧错误，确认测试能抓到它；记录后端、版本及未覆盖情形。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:428`，用户：“构建一些 jax 程序……拉取 HLO……确保所有内容都正确（HLO instruction 都能逐段解读）”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:986`，assistant 报告：“渲染后的文字和原文一致”“没有未被标出的文字”，并报告 21 个程序、42 份前后产物及逐节点遍历。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:995`，用户：“你怎么保证 ‘展示内存位置’ 展示的位置不出错？”
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1257`，assistant：“必须和例子来自同一次编译”；报告双向比对与恢复旧规则后的失败测试。

**频率 / 信心**：2 次直接请求、1 session；另有下述 2 次追问共同验证此主线。方法抽取信心高，常见频率证据不足。

**建议 skill description**：当开发编译产物解析、解释或可视化工具，并需验证语义正确性时，用目标后端真实样例、同次编译的外部判据、片段覆盖与反向验证建立回归语料。

**边界及反证**：片段都有说明不等于说明正确，必须补独立语义判据。历史“84 个测试通过”之后仍被用户发现另一视图有错误；不能把测试数量当正确性证明。`is_scheduled` 等具体判断不应无条件推广到所有 XLA 版本。

## 2. 把语义决策集中到共享模型，并用实验消除未知（建议并入 1）

**触发场景**：同一硬件语义在节点标签、解释、提示等多处显示，出现互相矛盾的结论；或遇到未公开编号。

**可复用步骤**：追溯具体指令及调用上下文 → 检查独立缓冲区是否存在 → 用统一状态模型供所有视图调用 → 对各视图同时做外部判据比对。未知编号先保留原值；收集使用位置、实现断言，再写显式指定资源的最小实验交叉验证；注明硬件与版本。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1260`，用户给出 bitcast 指令并问：“为什么 bitcast 那个 ROOT 的内存空间被标注成 HBM？”
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1346`，assistant：“上一轮只修了图上节点的内存标签”，解释面板仍用另一逻辑；改为共享 `bufferStatus`。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1356`，用户：“S(2) 是什么位置”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1454`，assistant 报告从 libtpu 断言与指定 SMEM 的 Pallas 实验取证，将散在四处的名称改为统一表，并限定 v6e / JAX 0.11.2。

**频率 / 信心**：2 次追问、1 session；共享模型与实验取证各 1 次完整实例。作为主模式的排错步骤信心高，单独建 skill 的重复证据弱。

**建议 skill description**：当编译器可视化多个视图对同一事实解释不一致，或后端字段含义未知时，追溯编译证据、集中语义状态，并以最小可复现实验验证未知映射。

**边界及反证**：不能仅凭标量形状推断存储位置，也不能把某代 TPU 的编号当通用常量；二进制断言与实验结论需明确版本。此处和模式 1 同属一次修复链，建议合并而非拆出多个技能。

## 3. 面向在线工具使用者压缩 README（次级候选）

**触发场景**：项目已提供可用网站，README 却围绕开发环境、技术栈或冗长 UI 操作解释。

**可复用步骤**：确定读者要完成的任务 → 放置直接使用入口 → 以领域收益组织简短功能条目 → 删去与使用决策无关的实现细节 → 用真实截图或基准对照证明能力；资源入库后检查相对链接。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1557`：“make it simple. don't include npm or any thing to run, let the user use website”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1602`：“don't mention any tech stack……just focus on JAX, TPU, HLO”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1682`：“What you can see 可以写短一点，只放大家想要的”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1740`：“这个图放到仓库里做对比。左边放这个图，右边放我们的”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1973`：要求把三张功能截图放入 README。

**频率 / 信心**：5 次文案请求（1550、1557、1578、1602、1682），另 3 次图片整合请求（1658、1740、1973）；1 session、同一 README 的连续打磨，不能算 8 次独立复现。候选信心中。

**建议 skill description**：当为已上线的交互工具编写用户向 README 时，围绕直接使用入口、领域能力和真实截图压缩文案。

**边界及反证**：不应变成“所有 README 禁写 npm / 技术栈”；开发库与自托管工具仍需安装步骤。若无其他 session 佐证，更适合作为文档技能中的一种受众分支。

## 未成立与覆盖清单

- `/Users/junyi/.claude/projects/-Users-junyi-Desktop-phd-qe-qe-slides/memory/MEMORY.md:1`：全文已读；脚本落文件供审阅是个人协作偏好，只有索引句，无完整方法链，不建议单独 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda-aug-12-learn/memory/MEMORY.md:1`：全文 3 行已读；有界 gate 验证边界、备用环境与插桩入口均为技术事实或环境记忆，缺乏过程证据。不展开未分配的链接文档。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-medusa-reborn/memory/MEMORY.md:1`：全文已读；模糊删除范围需澄清是一次反馈与操作约束，不足以成立专门 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude/memory/rental-search-near-mrt.md:1`：全文 17 行已读；存在“先筛地铁距离，再租金、双地点通勤”的顺序，但主要是个人找房条件，未见多 session 复现或可检验路线算法，不成立 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1`：完整流式扫描 2021 行、27 条非空 user 文本；按主题核实 986、1257、1346、1454 等 assistant 文本。UI 主题按钮、N-hop 点击时机、部署与提交请求均是单次功能/交付动作，不单独提炼。
