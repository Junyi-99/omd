分类：`tpu-performance-validation`
方向：把 TPU/Pallas 优化猜想转为可复核的测量、归因和适用范围。
领域：JAX/XLA、Pallas、MoE kernel、XProf。
频率：1 个原始 JSONL session（含延续/压缩记录）；另有 1 份不同 originSessionId 的 KDA 实验 memory 提供跨任务佐证。不能把追问、监控通知或压缩摘要算作独立复现。
为什么值得持久化：测量口径、真实形状覆盖和布局成本验证能跨算子复用；具体 PR、倍数和模型配置不应写死进 skill。
审阅范围与局限：全文读取 3 份指定 memory，Python 流式扫描日志全部 4,467 行中的 137 条含文本 user 记录，复核关键 assistant 上下文。通知和压缩摘要不作为新的用户需求；未访问链接或运行实验，历史性能数字仅作方法证据。

## 1. 先验证测量工具，再用 profile 解释优化

- **触发场景**：声称 fusion、layout、block 参数减少访存或等待，却缺少硬件证据；scope 指标与墙钟矛盾。
- **可复用步骤**：
  1. 固定基线版本、设备、shape、精度和输入；先测数值正确性。
  2. 将真实被测入口 JIT，预热并同步设备后计时；保留原始墙钟，明确 dispatch 与 kernel 的边界。
  3. 用 HBM bytes、HLO/LLO 和 XProf 检查具体机制；scope 按 lane 建嵌套关系，区分 inclusive/exclusive，不把跨 core 累加时间当墙钟。
  4. 对 trace 派生指标做规模敏感性校验：增加工作量，观察时间、事件数和捕获完整性；事件数封顶时停止用其比较总时长。
  5. 要求独立计时口径交叉验证；一旦矛盾推翻结论，同步更正文档、PR、图表。
- **原始证据**：
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:105`：“抓一次 HBM byte，看 gather 的往返还在不在”。
  - 同文件绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:537`：“你没查实的机制应该尽快 xprof 做 profile。”
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:2755`：“补上 exclusive self-time 再跑一轮”。
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:3622`：“加上 jit 重跑”。
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:4447`（assistant 事后解释）：“trace 的事件缓冲区会满”；对 bts=64 的收益论断已撤回（4440）。这是失败案例，不是可照抄的成功优化。
- **频率/信心**：上述 4 个明确测量请求均在同一 session；高信心支持流程，跨 session 的 trace 完整性校验仅有这一失败案例。
- **建议 skill description**：当优化 JAX/Pallas TPU kernel、比较实现或解释 XProf 瓶颈时，先验证 JIT 墙钟和 trace 完整性，再以 HBM/HLO/LLO 证据归因，校验规模敏感性并更正被推翻的结论。
- **边界及反证**：named scope 的嵌套减法正确不代表输入完整；源码确有重复加载，也不能推出特定参数一定更快。历史日志先把 self-time 称为可信，最后被墙钟和事件饱和证据推翻，因此不得沉淀“self-time 优先”的规则。具体 profiler flag 随版本变化，需现场探测。

## 2. 用真实模型形状和完整调用链约束优化收益

- **触发场景**：微基准显示收益，但真实模型、padding bucket、拓扑、布局转换或端到端收益尚未覆盖。
- **可复用步骤**：
  1. 先搜索仓库现有实现和 PR 合入情况，列清 baseline、unroll、fused 的代码来源，避免重复实现或错比版本。
  2. 从实际支持模型提取 K、hidden、FFN 维度，结合真实 token padding bucket 建覆盖表；显式列出未支持条件。
  3. 固定同一输入布局，把必要 reshape/relayout 和外围处理算入候选成本；比较数值误差与性能退化。
  4. 按实际硬件逐步扩展 token、ep 和合法 block 参数，保留失败点，不把非法配置或 OOM 当成无结果省略。
  5. 同时报 isolated kernel、完整模块占比、compiled graph 和 throughput；按形状说明收益，不报单一最好倍数。
- **原始证据**：
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:834`：“按实际的 padding bucket 重扫一轮。”
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:1438`：“找现有的用到 EPMOE 的模型，拿到他们的 K 和 hidden”；1480 扩展 K 范围，属同一次 sweep 扩充。
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:1731`：“测一下我们的 pr 的 compiled graph 和 throughput”；1785：“测一下 unpermute 在完整 EPMoE.__call__ 里的占比”。
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:2160`：“先搜索一下……有没有实现 fused unpermute”。
  - `/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:3159`：“把扫描扩到 token 桶和 ep 这两维。”
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-1550-vs-1554-xprof-substep.md:12`：60 个 shape、warmup/迭代口径；14 行先交代数值一致；35–41 行报告 kernel 外开销及未测 ragged/unbounded 的限制。
- **频率/信心**：日志中可区分 unpermute 和 fused-v2 block config 两组实验，但仍只算 1 个 session；KDA memory 增加 1 个独立实验来源。高信心，适合成为统一性能验证 skill 的覆盖模块。
- **建议 skill description**：当评测 TPU kernel 优化或准备性能 PR 时，从真实模型与 padding bucket 建 shape 矩阵，计入布局和外围成本，验证正确性及退化，并连接微基准、模块占比与端到端吞吐。
- **边界及反证**：不要固化这次 token 集、K 值或“忽略某模型”的临时选择；全笛卡尔积可能很贵，应先代表点再扩充。KDA memory 的 3.6–6.1× 是该实验结果，不能外推至其他算子。此 pattern 与上一项可合为一份 skill，避免拆成多个窄技能。

## 候选与不成立

- **机制图辅助解释**：日志 338、3414、3484 要求解释消掉什么维度、before/after 气泡，以及 HBM/VMEM 的反复读写。可作为性能报告中的可选步骤；仅同一 session 的迭代修图，尚不足单独确立绘图 skill。
- **Spot 实验恢复/状态巡检**：2374、2485 有抢占和每 5 分钟检查请求；另一份 sglang memory 只有 B200 建机索引，没有实际恢复步骤。保留候选，不从索引标题推演通用云资源编排方法。
- **外部写入控制、PR 文风、固定符号**：scanner memory 第 1 行和日志 918、1013、1438、1714 属权限/表达偏好，不能各自转成 skill。具体 SeleniumBase 迁移、GitHub 搜索限制和 B200 配额也是项目状态/参考事实。

## 覆盖情况

- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-ChatGPT-API-Scanner/memory/MEMORY.md`，3 行；未追读未分配链接。
- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-1550-vs-1554-xprof-substep.md`，46 行。
- 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/memory/MEMORY.md`，2 行；仅索引，不冒充读取底层建机记录。
- 全量流式扫描：`/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl`；全部 user 文本已扫描，选取 530–545、2748–2760、3600–3655、4400–4470 区间核实 assistant 上下文。未读取其他 reviewer 结果，未创建 skill。
