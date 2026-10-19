# 研究与证据整理

先做实验产物打包、增量研究、模型支持审计。

## 1. 实验产物打包

`research-artifact-bundle` · 优先

- **什么时候用**：实验结束、远程环境将释放，或需要分享和复跑结果。
- **怎么做**：每轮建目录 → 分开原始、标准化、派生数据 → 保存脚本、命令、版本和日志 → 同步产物 → 提供重建命令。
- **频率**：2 个原始会话 + 3 个不同来源会话的 memory。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2770)：“跑的脚本为什么不落盘”。

## 2. 增量研究流水线

`evidence-research-pipeline` · 优先

- **什么时候用**：周期更新多来源研究、报告或 dashboard。
- **怎么做**：读取上一版 → 确定观察窗口 → 程序提取事实 → Agent 判断变化 → 检查来源和状态 → 输出新增、变化及缺口。
- **频率**：3 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/83326f11-ca34-495a-b889-48cbb628584b.jsonl:729)：“程序是给 AI Agent 提供 原语、事实”。

## 3. 模型支持审计

`model-support-recon` · 优先

- **什么时候用**：从仓库、配置和 PR 构建模型或算子支持矩阵。
- **怎么做**：按真实 config 和推理图识别模型 → 核对 ID 与别名 → 分别记录实现、加载、验证和优化证据 → 检查引用和展示状态。
- **频率**：2 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:2294)：“为什么会出现拼出来的假身份”。

## 4. 资料驱动架构展示

`source-backed-architecture-explorer` · 配方

- **什么时候用**：把论文、照片和规格资料做成架构学习工具。
- **怎么做**：统一实体层级和单位 → 建带来源的规格与术语模型 → 同一数据源驱动表图 → 分开物理与逻辑视图 → 逐层缩放验证。
- **频率**：1 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-architecture/d88da7fc-99bc-415f-945c-438a4a6ecb51.jsonl:596)：“一个 node 应该画成 chip”。

## 5. 审计 kernel 优化机会

`kernel-opportunity-audit` · 配方

- **什么时候用**：自动生成的优化列表与认可样本差距很大。
- **怎么做**：区分功能、实现和融合缺口 → 对照样本找漏检 → 查后端及调用链 → 用反证筛候选 → 列待 profile 项。
- **频率**：1 个会话；作为支持审计配方。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/5c3d2ba3-02d4-43b3-9c8d-7caedf662eb2.jsonl:4443)：“三类候选并排”。

## 6. 找回项目与构建指令

`project-provenance-recovery` · 配方

- **什么时候用**：记得成果，但忘了项目目录或最初的 bootstrap prompt。
- **怎么做**：联合搜索项目与会话 → 用 README、remote 和产物消歧 → 找到构建指令 → 整理目录、用途和来源。
- **频率**：1 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude/66184f65-6f18-4649-91b1-7058cdad7e64.jsonl:9)：“可能有不止一个项目…找全…聊天记录”。

## 7. 统一研究数据契约

`typed-research-contracts` · 模块

- **什么时候用**：目录、状态和持久化 schema 出现重复定义。
- **怎么做**：找出重复 → 选唯一类型与枚举来源 → 派生目录和序列化 → 检查跨字段约束 → 运行类型检查。
- **频率**：1 个会话；并入研究流水线。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-ownership/c8e63cde-9fb2-4fa0-8a2c-0fce0226703d.jsonl:282)：要求 YAML typed models。

## 8. Prompt 重建实验

`prompt-rebuild-audit` · 观察

- **什么时候用**：要求用自然语言规范重复生成相同文件树。
- **怎么做**：约定结构、语义或字节一致目标 → 保存输入版本 → 隔离生成 → 比较产物 → 根据 diff 修订并复测。
- **频率**：1 个会话；曾通过嵌入全文实现复制式收敛。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-opensource-dashboard-genesis/ece58979-6af1-4700-94ff-a79c3fbb7ac8.jsonl:13)：要求 genesis 多版本、比较结构和内容并保留 review。

详细记录：[reviewer-01](reviewer-01.md)、[reviewer-10](reviewer-10.md)、[reviewer-12](reviewer-12.md)、[reviewer-21](reviewer-21.md)、[reviewer-28](reviewer-28.md)、[reviewer-30](reviewer-30.md)。
