# reviewer-28

- 分类：`evidence-preserving-workflows`
- 方向：让研究/实验产物可追溯、可重建；其次是把文件约定落实为可验证契约。
- 领域：GitHub 项目观察、TPU 实验、prompt 生成系统。
- 频率：核心证据保全方法见 1 个日志 session + 1 份来自另一 session 的 memory；其余两项各仅 1 个独立 session，属于候选，不能将多轮修订或并行生成次数当作独立复现。
- 持久化价值：容易反复遗漏的是原始证据、覆盖缺口、生成输入与验收口径；这些有明确步骤及失败条件，适合成为窄范围 skill。工具版本、仓库落地偏好不应混入。
- 范围与局限：全文读取 3 份 memory；Python 流式扫描 5 份 JSONL 的全部 user 消息，并核对相关 assistant 上下文。工具结果、通知、goal 自动重述不计用户独立需求。未查分配范围外的原始 memory 来源 session，也未复跑历史代码；完成情况以历史记录为限。

## 1. 保留原始证据，再生成标准化与派生结果（建议沉淀）

**触发场景**：采集开源项目资料或在临时机器上运行实验，需要让结论可核查、可重算。

**可复用步骤**：
1. 按运行/快照建立独立目录，记录日期、对象全名、版本或提交、请求时间范围。
2. 将原始输出原样保存；与标准化数据、分析结果分层，保存实际运行命令。
3. manifest 分别记录请求范围、实际覆盖、各来源状态、分页/采集错误、已知缺口；未采集不得写成 complete。
4. 在临时机器销毁前上传原始日志到持久存储；核对产物清单，分析结束保留本地副本与引用。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/c8e63cde-9fb2-4fa0-8a2c-0fce0226703d.jsonl:9`：要求快照目录与 `manifest、raw、json、derived`；`:112` manifest 明确给出 `actual_start`、`pages_fetched`、`known_gaps`、`collection_errors`。
- 同一绝对文件 `:331`：用户指定 `gh api .../pulls/1600` 获取 raw；`:335` assistant 核对 REST 原文与 `gh pr view` 字段输出不同；`:345` 落地实际命令 `_cmd.txt` 及 `/issues` 需剔除 PR 的规则。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/save-raw-experiment-logs.md:10-16`：`always save the raw logs ... not just the parsed summary tables`，解释临时 VM 回收会丢证据；`:18-21` 要求持久化后拉本地副本。

**频率/信心**：2 个独立来源 session（其中 1 个仅 memory 间接可见）；跨 GitHub 采集与实验日志，高信心存在共通需求，中等信心适合合为一个入口。

**建议 skill description**：当采集外部资料或运行临时环境实验并需可复核结果时，建立原始、标准化、派生分层，记录覆盖缺口与采集命令，及时持久化原始证据并验证归档。

**边界及反证**：GitHub schema 与 TPU 日志类型应是不同配方；具体 bucket、目录根和 YAML 偏好保留项目配置。ownership 的 `:356` 明确 raw 子目录仍为空，因此只证明骨架与规范落地，不能声称采集器已实现或真实覆盖完整。

## 2. 将目录与 YAML 约定变为单一类型契约（候选）

**触发场景**：脚本逐渐变成采集工具，目录名、manifest、状态枚举、注释多处重复，人工填写容易漂移。

**可复用步骤**：
1. 识别真正重复的目录/状态定义，指定唯一代码来源。
2. 用 typed models、enum 及跨字段校验约束持久化格式；保留未知值与失败原因，拒绝拼错字段。
3. 序列化文档及目录从模型派生；无法消除的契约重复用一致性检查守住。
4. 用临时输出生成一次实际树，运行严格类型检查与少量校验测试，再接入 pre-commit。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/c8e63cde-9fb2-4fa0-8a2c-0fce0226703d.jsonl:239`：`folder structure 能用代码定义吗？`；`:273` 追问三个状态是否有 enum；`:282` 要求所有 YAML 成为 typed models；`:360` 要求 precommit 与严格 mypy。
- 同一绝对文件 `:229`：记录 `extra="forbid"`、失败状态需 reason、显式 null；`:252` 目录集中到 `layout.py` 并用检查防 manifest 漂移；`:327` 从 `Sources.model_fields` 派生 raw 子目录；`:403` 报告严格 mypy 与 pre-commit 通过。

**频率/信心**：1 session 的至少 4 个契约收敛请求；方法证据充分，独立复现不足，中信心候选。

**建议 skill description**：当 Python 数据采集脚本的目录、YAML schema 与状态定义开始重复时，将其收敛为类型契约，派生输出并验证目录、序列化及覆盖状态一致。

**边界及反证**：不能泛化为所有 Python 项目都要固定版本、十个 schema 文件或双类型检查器；`:403` 实际删掉 Pyright 改用 mypy，说明依赖栈是任务选择。只生成少量静态文件时不应预先搭完整框架。

## 3. 用隔离重建与逐轮差异校准生成 prompt（实验候选）

**触发场景**：用户要求自然语言规范能重建同一文件树，并明确要验证结构/内容一致性。

**可复用步骤**：
1. 明确一致性口径：语义、结构或逐字节；保存输入与每轮输出。
2. 新执行上下文只读本轮规范，避免复制前一轮产物伪装重建。
3. 比较文件树及内容，将差异追溯到自由决定，修改规范再生成下一版。
4. 保持规范不变做确认轮，保存 diff、hash、检查脚本与未满足项。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard-genesis/ece58979-6af1-4700-94ff-a79c3fbb7ac8.jsonl:13`：要求 genesis-0 到 genesis-10，比较相邻文件夹，修改 prompt，并保留每版 review。
- 同一绝对文件 `:259`：不变 prompt 检验固定点；`:273` 报告首次逐字节一致；`:496` 总结结构约束后仍有文字差异，最终把文件内容写入 prompt，规范从 1178 字节增到 22386 字节。

**频率/信心**：1 session、11 个生成版本；仍仅 1 次独立实验。中等信心作为验证配方，低信心证明普遍 skill 需求。

**建议 skill description**：当用户明确要求 prompt 生成产物可重建时，隔离生成上下文，保留迭代版本，以结构和内容差异收敛规范，并验证不变输入下的重复结果。

**边界及反证**：把全文嵌入 prompt 证明的是抽取/复制确定性，不能证明独立推导能力；不要默认所有写作都需字节一致。`:496` 还说明早期 prompt 在生成后被修改，导致自身 hash 检查失败；实际 skill 应分开保存输入快照与修订版，避免破坏审计链。

## 不建议独立成 skill

- `/Users/junyi/.claude/projects/-Users-junyi-claude-deck/memory/deck-no-git-add-all.md:10-14`：特定仓库误 stage worktree 的一次反馈，适合 repo guardrail；可提供显式 stage + status 的通用检查，但当前证据不足独立建 skill。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-ledger-app/memory/merge-workflow.md:10-14`：直接 main、不开 PR 是仓库偏好，不是跨场景方法。
- `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/4109ca0a-4db2-407a-be71-a61f42b27801.jsonl:90,129,157`：表达自建/统一 skills 的意愿及迁移授权边界，没有已验证的可复用迁移流程。
- `/Users/junyi/.claude/projects/-Users-junyi-Downloads-lectures/1ded4a82-e6c2-412a-b97b-e818cf6b6327.jsonl:115,141`：两次 trace 生成请求在同一 session；`:155` 有批处理成功/失败记录，但课程命令、GPU 限制是单项目操作，暂作候选配方素材。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/d694c031-702e-4cc4-8495-da1469d58ad3.jsonl:9`：新私有仓库 commit/push 单次操作，不足以沉淀独立 skill。

## 覆盖

- memory 3/3 全文：deck-no-git-add-all、merge-workflow、save-raw-experiment-logs。
- JSONL 5/5 全文件流式扫描：genesis（65 条 type=user）、ownership（58）、dotfiles（14）、lectures（22）、hlo-visualizer（4）。这些为记录条数，包含工具回传/自动消息，不能当用户需求次数。
- 已核对 genesis、ownership、lectures、hlo-visualizer 的 assistant 文本；dotfiles 以全部 user 文本判断，不推断未验证迁移成果。
