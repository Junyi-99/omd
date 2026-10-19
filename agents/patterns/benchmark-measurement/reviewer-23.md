分类：`accelerator-performance-validation`
方向：用可比的测量与编译产物解释加速器性能，辅以有限范围的代码审计候选。
领域：JAX / Pallas / TPU、GPU 推理基准、Python API 类型。
频率：7 份日志中只有 4 份有真实用户请求；性能相关请求出现在其中 2 个独立 session（deciax、MXU sweep）。另有 3 份 memory 提供不同项目的测量方法，但不把 memory 条数换算成独立 session 次数。
为什么值得持久化：对齐工作量、计时边界与编译层次能避免把吞吐当延迟、把 FLOPs 当瓶颈；这些方法跨项目复用，具体机器数字应继续留在项目记忆。
审阅范围和局限：全文读取全部 3 份 memory；Python 流式扫描全部 7 份 JSONL 的 user 文本，再核实相关 assistant 上下文。工具结果、任务通知、续接摘要不算用户需求；没有复跑历史实验，历史技术判断只作为方法证据。

## 1. 对齐测量口径，再定位 host / device 性能瓶颈（建议沉淀）
- **触发场景**：启用优化后端仍没有预期收益；比较单请求与 batch；引用不同设备或互连的公开数字。
- **可复用步骤**：
  1. 记录实际后端、设备、输入 shape / padding、dtype、host 负载及编译预热条件。
  2. 为每种 shape 单独预热；异步设备计算显式等待完成，分开测 device forward 与端到端。
  3. 同时报告请求延迟、每项分摊时间、吞吐；分母写实际问题数，不能拿 batch 分摊时间冒充单题响应延迟。
  4. 扫工作量并对照参考实现；用差值定位 host 固定开销、launch、拷贝或设备计算，不先归因 kernel。
  5. 跨来源比较建立设置清单：计时位置、p50/mean、消息大小、并发、链路数；不匹配的标成参考，设计同口径补测。
- **原始证据**：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/b8719eb3-b653-4b28-a0e3-6d7c997c4778.jsonl:130`：“在 v5e 上用新的包跑还是 4ms 多？”；`:165` 给出同步 device forward / e2e 对照，并明确 v5e 小尺寸归因尚未实测。
  - 同文件 `:169`：“给我个 50问题的 batch”；`:182`：“per question 的计算呢？”；`:189`：“和 single request 的能放一起比吗？”；`:192` 纠正原请求其实包含 3 个问题。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/memory/laya-bench-findings.md:11`：“CPU/launch-bound”；`:13` 比较 parity、warmup 与不同 batch 的优劣；`:15` 交代共享主机噪声。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/ici-external-interconnect-refs.md:12`：“No published source matches our setting”；`:15` 要求同 sizes / statistic 的 GPU-cluster 补测。
- **频率 / 信心**：直接请求为 1 个 session 内 4 次相关提问，另 2 份项目 memory 佐证；不等于 6 次独立复现。方法价值高，独立重复需求证据中等。
- **建议 skill description**：当加速器优化收益异常，或比较 batch、设备、互连性能时，先统一输入与计时口径，拆分 host/device 开销，输出可复现基准与明确的比较限制。
- **边界及反证**：不固定使用 50 题、某 tile 或某硬件阈值；memory 中测得的延迟不是可泛化定律。没有 TPU 的本地脚本或推测不能标成已验证性能结论。

## 2. 用降低后的 IR 与指令计数检验 kernel 性能假说（建议沉淀，可与上项组成同一 skill）
- **触发场景**：Python 写法相似但速度悬殊；batched、concat、循环或布局变体的效果难解释。
- **可复用步骤**：
  1. 逐变体列出输入形状、数据依赖、循环边界与矩阵乘形状。
  2. 导出实际 lowering 产物，必要时开启文本 debug；定位循环是否保留、batch 是否展开、索引是否动态。
  3. 对照相应编译器版本源码，区别 Python 函数嵌套与设备 IR 结构。
  4. 用 profile、低层指令数与布局搬运解释耗时；检查 scope 分母与归一化，避免仅凭 FLOPs 判断。
  5. 设计一个可证伪的最小对照（如 unroll）；计时与带 trace 标注的构建分离；保留负面结果。
- **原始证据**：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/0bf60b63-b45f-4292-a817-fd921134d390.jsonl:9`：“编译后有区别吗？不都是一个 fori_loop 吗？”
  - 同文件 `:65`：“改用 `debug=True` 打印文本 MLIR”；`:117` 基于 5 个 variant 的 Mosaic MLIR 对比循环 / 形状，提出 unroll 验证，但没有报告验证结果。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/tpu-kernel-first-principles.md:10`–`:20`：MXU 指令数只差 1.10× 而耗时差 7.6×；布局局部优化反而变慢；同名 scope 合并和 trace region 引入测量偏差。
  - 同 memory `:22`：“252 个 shape 的 sweep + LLO 指令计数”，并保留负面结果。
- **频率 / 信心**：1 个直接用户 session + 1 份另项目 memory；方法交叉佐证，中高信心。没有足够证据称为多 session 高频显式请求。
- **建议 skill description**：当 JAX/Pallas kernel 的循环、batch 或布局变体出现难解性能差异时，用实际 lowering、编译器源码、指令计数与最小对照验证原因。
- **边界及反证**：历史源码版本、静态展开策略与硬件行为必须重新确认；“加 unroll 应靠拢”只是待测假说。不可把“不要从 FLOPs 起手”写成所有 kernel 的绝对规则。

## 3. 沿调用链收紧闭集类型并校验外部入口（候选）
- **触发场景**：实现选择、dtype、模型名等有限值被标成通用 `str`，调用链丢失约束。
- **可复用步骤**：从公开 API 追到配置、服务入口、dispatch 和模型；定义共享 Literal / 别名并传播；外部字符串先校验再转入闭集；运行严格类型检查与相关现有测试；核实自由文本仍保持开放。
- **原始证据**：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/b8719eb3-b653-4b28-a0e3-6d7c997c4778.jsonl:197`：“把链路上的 type 都检查一下”；`:256`：“至少是个 Literal。load 和 serve 的 dtype 也改了”。
  - 同文件 `:389`：列出类型沿调用链传播，并发现 argparse 不校验环境变量默认值；保留 revision / device 等自由文本为 `str`。
- **频率 / 信心**：1 个 session 中 2 次范围递进请求，属于同一次修复；方法清楚，但不足独立新建 skill。更适合作为 Python API 审计 skill 的检查项。
- **建议 skill description**：审计 Python API 的有限值参数，沿调用链保持类型约束，并在 CLI、环境变量等外部入口提供运行时校验。
- **边界及反证**：不能把所有 `str` 改 Literal；同样不能只靠类型注解替代输入校验。

## 不单独沉淀的内容
- 清理 PR：deciax `:475` 要审查必要测试，`:840` / `:843` 是同一轮五 agent 清理请求的重复表达，不能计为两次独立复现。`:1145` 核实建议后保留不同语义 helper、避免循环导入、说明未跑 TPU 测试，值得作为代码审计案例，尚不足另立“并行清理”skill。
- PR Mermaid、去会话标识、逻辑块空行与不提交 benchmarks 属于个人 / 仓库约定；发 dev3 和加 badge 是单次任务。
- SSH jump / SOCKS / 自动认证仅 1 个 session；历史含凭据，不复制，不据此建立自动保存密码流程。
- 查询 Claude 历史目录仅一次，不构成工作流；互连论文列表和 laya 延迟数值本身属于 reference memory。

## 覆盖情况
全部路径均为 `/Users/junyi/.claude/projects/` 下，以下为已扫描文件的相对路径；上文关键证据保留绝对路径。
- `-Users-junyi-claude-deciax/memory/laya-bench-findings.md`：全文 15 行。
- `-Users-junyi-claude-kimi-k3-kda/memory/tpu-kernel-first-principles.md`：全文 25 行。
- `-Users-junyi-claude-tpu-spot/memory/ici-external-interconnect-refs.md`：全文 15 行。
- `-Users-junyi-claude-deciax/b8719eb3-b653-4b28-a0e3-6d7c997c4778.jsonl`：全量 user 扫描，筛除工具结果 / 通知 / 续接摘要，按主题回读 assistant。
- `-Users-junyi/e3474b52-1ff4-4b53-8978-e66a322d6c59.jsonl`：全量 user 扫描，SSH 单次需求，不输出隐私。
- `-Users-junyi-claude-tpu-mxu-sweep/0bf60b63-b45f-4292-a817-fd921134d390.jsonl`：全量 user 扫描，回读 IR 解释上下文。
- `-Users-junyi-claude-tpu-mxu-sweep/904e9c80-c58e-4ba3-aebc-813ac32f6564.jsonl`：全量 user 扫描；仅历史目录查询与本地登录记录。
- `-Users-junyi-claude-agent-container/31923175-a061-4173-a3f4-0fbb810b3ec4.jsonl`：全量扫描，仅 1 条 bridge-session，无 user turn。
- `-Users-junyi-claude-sglang-tpu-sept15/38b4fbb8-30ae-433d-835f-b36fd32db4af.jsonl`：全量扫描，仅 1 条 bridge-session，无 user turn。
- `-Users-junyi-claude-tpu-mxu-sweep/d7258cf3-dc08-48a7-ba4e-522c3632bb85.jsonl`：全量扫描，仅 1 条 bridge-session，无 user turn。
