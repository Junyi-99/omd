# reviewer-02：性能实验的设计、校验与交付

- **分类**：`performance-experimentation`
- **方向**：从硬件性能疑问形成可复现、可比较、能被解释的实验；优先沉淀测量可信度检查。
- **领域**：TPU / JAX / Pallas GEMM、VMEM、MXU、流水线与数据可视化。
- **频率**：主日志仅 **1 个独立 session**；4 个 memory 来自其他项目，但未提供相应完整日志，不用于夸大工作流频次。其中多次实验迭代，不能算跨 session 复现。下文次数均指该 session 内可定位的独立问题或交付轮次，排除进度询问、截断重发、工具回报。
- **持久化价值**：固定开销、计时边界、编译器消除、实验模式混淆在同一研究中反复造成错误解释。这些有明确的纠错步骤，可迁移到其他加速器实验；具体峰值和 Apple 风格偏好本身不应成为 skill。
- **范围与局限**：完整流式扫描 manifest 唯一日志的 3,411 行、491 条 `type=user` 记录；138 条含文本（包含通知、技能注入及图片占位），逐条主题审阅其中人工文本，并按需读取 assistant 上下文。另全文读取更新 manifest 中的 4 个 memory 文件。以下是历史证据审阅，未重新运行实验；历史助手报告的测量结果不视作独立外部验证。

证据文件统一为：
`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl`。
下列每项原始证据仍给完整路径，方便独立定位。

## 1. 用计时边界与重复次数扫描解释反常性能

**触发场景**：理论上更快的实现反而更慢；同一矩阵形状跨实验吞吐相差很大；宣称数据常驻但结果包含搬运成本。

**可复用步骤**：
1. 对齐 shape、dtype、计算语义、FLOP 计数、每次调用的工作量、输入输出所在存储层；明确包含编译、预热、dispatch、DMA 中的哪些阶段。
2. 构造只搬运不计算与搬一次重复计算 R 次的对照；扫描 R，拟合 `latency = fixed + R × compute`，比较边际吞吐与端到端吞吐。
3. 检查循环预算量纲、重复次数上限与 JIT cache key，确认每个 R 真正重编译，并验证结果未被消除。
4. 区分被固定成本污染的绝对值、排名和仍成立的定性结论；重测受影响数据，再更新图例、正文和汇总页。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:2587`：“double buffer 比 resident 的速度还要快……确定计时边界没有错吗？”
- 同文件绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:2608`，助手报告采用“只拷不算”“拷一次、算 R 遍”分离固定成本；行 2612 用户要求图例明确 HBM→VMEM 与数据已在 VMEM 两条线。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:3157`：“同样是 1024³，这里只有约 380 TFLOP/s……777……856”，要求独立调查；行 3176 的子任务回报定位 FLOP budget 量纲错误，用 R 扫描展示排名会改变，并指出修改 R 不改变 JIT key 会误用缓存。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:3317`：“best TFLOPs 只有 752？但是 pallas 那边有 912？”行 3329 再追问 VMEM 常驻情形；行 3325 明确将部分解释标为未验证推测。

**频率 / 信心**：1 session 内至少 3 组异常比较（resident 对 double buffer、小形状跨脚本、VMEM 最优对方阵峰值）；前两组有探针结果回报，第三组尚有未验证解释。方法价值高，跨项目频率未知。

**建议 skill description**：当加速器 benchmark 出现反常排名或跨脚本吞吐不一致时，审计计时边界、工作量与 JIT 缓存，通过重复次数扫描分离固定开销，并标记需要重测的结论。

**边界及反证**：不能把所有差距归因于计时错误；搬运与计算重叠、tile 结构也会真实改变性能。不固定 TPU 型号或峰值，不用近峰值作为正确性的充分条件。独立盲审仅被明确请求一次，宜为可选步骤，不足以单列高频 skill。

## 2. 将硬件机制假说变成受控 sweep，并验证控制变量确实生效

**触发场景**：想知道容量、tile、对齐、buffer 深度如何影响吞吐，而不是只得到一张规模曲线。

**可复用步骤**：
1. 先写问题与可证伪假说，列出实验变量、固定变量和对照；分别改变存储预算、真实占用、矩阵形状或实现方式。
2. 先 coarse 探索边界，再针对机制补点；预先筛除装不下的配置并记录失败、跳过及原因，不把缺测当零性能。
3. 区分配置额度与真实存活分配；内存占用探针要完整写入并让读回结果影响可观察输出，避免优化器删掉占用。
4. 用编译成功/失败、结果正确性及耗时共同确认变量生效，再解释机制；遇反例修正模型和已发布说明。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:450`：“VMEM 从 128 每次缩小 4MB……排除掉……跑不了”的 tile；行 427 要求 coarse 足够、无需 dense。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:511`：同时 sweep M/K/N，关注对齐及拼接与小矩阵连续发送的差异。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:959`：“allocate 一个 placeholder 的 VMEM 占用……编译器会根据实际 VMEM 占用调整 prefetch 和 overlap”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:2513`：“48 MiB 的 scratch，vmem_limit_bytes 设成 16 MiB，照样编译成功”；行 2551 助手纠正：原探针仅写首尾，可能被优化掉；完整写读后才得到符合额度约束的结果。

**频率 / 信心**：至少 3 类实验问题（容量×tile、shape/对齐×实现、流水线深度），均在同一 session；“控制变量被优化消除”的明确纠错只有 1 例。建议合为受控性能实验方法，不单独泛化成编译器知识 skill。

**建议 skill description**：将硬件吞吐问题转成有对照的参数扫描，验证容量、真实占用与编译器行为的控制变量确实生效，记录不可行点并据反例修正解释。

**边界及反证**：配置额度不一定直接改变速度；编译成功不能证明声称占用真实存在。不要把日志中某版本的 VMEM 默认值、buffer 份数或特定 tile 最优值固化成通用规则。与模式 1 共用测量框架，但本项关注实验因果设计。

## 3. 把实验代码、绘图与可分享说明组织成可追溯交付

**触发场景**：已有多组实验，读者分不清做了什么、方法语义或图上曲线对应什么，需要可复现研究汇总。

**可复用步骤**：
1. 每个实验脚本只产出数据，配对应 `_plot.py`；README 简述问题、验证目标和初步结论。
2. 先用具体张量形状说明各模式的 A/B 是否共享、是否含 batch、拼接方式和计时存储层，避免只用实现术语命名。
3. 按语义分面，重叠曲线分图；图中交代 Experiment / Question / Findings、计时边界和缺测情况；数字从数据计算。
4. 渲染检查标签冲突与裁切，随后把确认结果、局限与交互图链接同步到分享页面；解释变化时同步修改旧说明。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:831`：“图里面标注清楚做的是什么实验，为了探究什么内容，结论如何”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:1373`：“一个实验一个脚本……对应一个 _plot.py……不希望画图和实验混在一起”，并要求简短 README。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:1762`：“shared weights 和 independent weights 分开画……有些重叠”；行 1884、1916 用户持续纠正 batch 维度与图例语义。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl:2254`：“给我一个完整描述，说清楚我们做了什么实验”；行 2289 要求分享页面，2468、2562、3125 要求更新与补链接。

**频率 / 信心**：1 次明确代码组织要求，至少 2 轮图表组织要求，1 次创建分享页与 3 次后续更新。均属同一研究交付链，不算独立复现。建议候选，优先级低于前两项。

**建议 skill description**：将多组性能实验整理为实验脚本、独立绘图脚本和简短研究索引，用张量语义、计时边界及可追溯数据制作可分享的结果说明。

**边界及反证**：Apple 风格、删某子图、改特定图例是局部偏好或单次编辑，不足以单列 skill。已有通用数据可视化 skill 的长文本在日志中属于注入内容，不能当用户历史需求或独立复现。

## 补充 memory 的判定

4 个 memory 全文读取后，没有新增足够强的独立 skill：

- `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/memory/dotfiles-history-merge-2026-10.md:11`–25：一次历史合并事实和并发会话下操作 Git 前复查状态的项目规则；保留为项目记忆。“远端某目录胜出”的解决规则不适用于其他合并。不能把一次成功处置计为经常复现。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/memory/xtra-disk-placement.md:11`–14：安装前检查挂载空间有复用价值，但指定服务器与目录主要是环境事实；只记录一次容量事件，宜留作环境约束，不独立建 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/feedback_writing_cut_setup_bloat.md:10`–29：“相邻章节同题合并—表格前置—方法细节移附录—逐段查重”是明确可执行的编辑方法，可作**候选技术写作去冗余模块**。来源只有一次文章重写反馈，未达到独立高频模式；它对模式 3 的“简短说明”提供跨项目方向性支持，不能算模式 3 全流程复现。用户要求简短与章节安排仍应优先作为偏好保存，不把所有技术文章强制改成同一结构。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-zhifeng-shao/memory/MEMORY.md:1`：只有索引摘要“附件是格式权威样本，先核对再写”；未分配目标文件，不能据索引推断已执行的步骤或次数，证据不足。

## 不独立成立的候选

- **远程运行、轮询、拉取结果**：有实际执行链及大量“进度如何”，但只见一台服务器上的连续研究；许多重复消息来自后台通知。可作为上述实验 skill 的执行模块，证据不足以证明值得另建通用 SSH 调度 skill。
- **术语问答**：多次询问 batched / concat / scratch 等，暴露命名与模型不清；宜纳入实验协议和解释步骤，不能把具体答案变成方法 skill。

## 覆盖情况

| manifest 文件 | 覆盖 |
|---|---|
| `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/e0390c91-ec12-43eb-93e6-71bbfe6a0839.jsonl` | 全部 3,411 行流式解析；全部 user 文本扫描；按主题抽读 assistant 与子任务结果，不将通知计为人工需求。 |
| `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/memory/dotfiles-history-merge-2026-10.md` | 全文 25 行。 |
| `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/memory/xtra-disk-placement.md` | 全文 14 行。 |
| `/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/feedback_writing_cut_setup_bloat.md` | 全文 29 行。 |
| `/Users/junyi/.claude/projects/-Users-junyi-claude-zhifeng-shao/memory/MEMORY.md` | 全文 1 行；索引指向的其他文件不在分配范围。 |

只写本报告，未创建 skill，未读取其他 reviewer 结果，未修改实验文件。
