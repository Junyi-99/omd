## 1. 候选：跨执行路径核对运行时数据契约

**触发场景**：推理框架新增 batch 字段、状态 buffer、模型特性；或 rebase 后正常路径可用但预编译、AOT、缓存恢复异常。

**可复用步骤**：
1. 列出数据的所有生产者和消费者：真实请求、dummy/warmup、AOT 直接构造、abstract pool、prefill/decode、prefix hit。
2. 对每条路径核对 shape、dtype、sharding、pytree 结构、状态初始化与所有权。
3. 把兼容性校验放到资源分配前；同时检查下游读的是派生配置还是原始开关。
4. 用相邻已支持模块的状态写回方式修复；检查转置和 snapshot slot 的布局契约。
5. 区分“编译用占位数据”和“缺必要资源”；不能让补零掩盖真实配置错误。
6. 报告验证覆盖与未运行路径，不把 lint/import 通过当作功能正确。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:9`：用户提供“conv state is never written to the track slot”，涉及 prefix hit 静默错误。
- 同文件绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:305`：“模仿 gdn 的 track 写法”。这是局部修复策略，不应泛化成禁止抽象的永久偏好。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:541`：“Overlap and spec are rejected at the first batch, not at startup.”
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:703`：“Precompile batches don't carry ple_embeddings”。助手在 :744 核实 AOT 直接构造不走 `init_new`，反驳了 reviewer 的覆盖假设。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:837`：“main's abstract recurrent pool assumes one conv buffer”；:880 指出不冲突的调用方也会漏传新参数。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:1473`：“该装 table 却没装…应该报错…不应该默认喂0”。

**频率/信心**：1 session、4 个独立技术条目（track、启动校验、batch 构造、abstract pool），另 1 次缺表语义纠正；方法价值高，独立复现信心中等。

**建议 skill description**：审阅或修改 JAX 推理运行时字段与状态时，沿 serving、warmup、AOT、缓存恢复等路径核对数据结构、资源和配置契约，给出最小修复及验证边界。

**边界及反证**：多条问题来自同一 reviewer/PR，可能共享根因；补零并非普适修复。用户要求先解释再改、精简代码等属于协作偏好，不构成这个 skill 的核心。

## 2. 候选：以有效增量和外部格式证据审阅堆叠模型 PR

**触发场景**：模型 PR 依赖另一个 PR、squash/rebase 后声称已覆盖旧修复，需要确认接线与权重加载。

**可复用步骤**：
1. 拉取真实分支，比较基线 tree，剥离已审阅依赖，确定有效增量。
2. 顺调用链检查 metadata、state 合并顺序、权重名/shape/sharding；不能只看函数参数名字。
3. 用发布 config 和 checkpoint index 核实层号、字段位置与权重映射。
4. 检查 dummy 模式、远程 WeightSource 等异常路径，以及已有测量是否发生在当前 rebase 之前。
5. 将发现按运行错误、待确认支持范围、可读性问题分级；明确静态审阅与实际执行证据。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:2200`：“把 #1781 拉下来，帮我过一遍他的代码”。:2017、:2037 是同一请求的重复，不能算三次独立审阅。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:2152`：助手通过 tree 一致排除依赖代码，并用 checkpoint index 证实层号。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:2298`：“测试没跑成…结论都来自读代码”，同时发现 CI 测试 skipped 及 rebase 前数据不能覆盖新行为。

**频率/信心**：1 次独立 PR 审阅、同 session 两轮核对；中等信心。适合作为第 1 项 skill 的审阅模式，暂不建议独立安装。

**建议 skill description**：审阅依赖其他分支的模型接入 PR，先核实实际代码增量，再用 config/checkpoint 与调用链验证接线，并标明测试和 CI 的证据缺口。

**边界及反证**：历史结论同时说 dummy 路径失败且“没有阻塞合并的 bug”，分级存在张力；skill 应要求给运行失败明确严重性，不能照搬该结论。发布配置约束也可能让泛化的层位置警告失效。

## 3. 候选：将实验过程压缩为可运行的维护文档

**触发场景**：PR 文档、注释混入大量试验历史，且补充文档修正主 README。

**可复用步骤**：先收集新增注释/主文档及纠正文档；保留代码看不出的原因和硬件限制；将旧结论替换成当前结论；删除过程性日志、临时环境命令与重复代码叙述；逐个核对命令参数，说明数字是否重测。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:1141`：“review…新增的 comments，精简一下”；:1254 将保留依据区分为“代码本身看不出来的原因和约束”。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:1632`：“README 写的太臃肿太长”；:1684 记录合并修正文档、291 行压到约 70 行、核对 CLI 参数，并披露数字未重跑。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:1705`：继续要求检查 PR 其他文档、comment、代码的臃肿。

**频率/信心**：1 session 内 2 次实际清理（代码注释、benchmark 文档），随后 1 次扩展检查；中等信心，暂作通用 review skill 的子流程。

**建议 skill description**：清理实验型 PR 的说明文档与新增注释，保留维护所需约束、当前结论和可运行命令，移除失效过程记录与重复叙述。

**边界及反证**：精简本身是偏好；不可按删行数评估质量，也不应默认删除原理、出处或复现实验所需信息。

## 不建议单独建 skill 的 memory

- `/Users/junyi/.claude/projects/-Users-junyi-claude-ChatGPT-API-Scanner/memory/reference_github_code_search.md:7-15`：代码搜索日期限定符与 repo 预筛选，是 1 次经验型工具限制；适合搜索 skill 的参考条目，不能据此声称常见工作流。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-fused-neumann-champion.md:11-23`：记录精度门、端到端性能、BT 改变导致下游成本反转，具有优化方法启发；但只有 1 份项目终态摘要，具体冠军、参数和未提交状态都是易过期项目事实，不宜固化。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/memory/MEMORY.md:1-3`：仅索引。空行是偏好，VMEM→MXU 是研究目标，数据相关峰值是值得查证的基准陷阱；未读取未分配的链接正文，不能从标题推导完整 benchmark skill。

## 覆盖

- 3/3 分配 memory 全文读取。
- 1/1 分配 JSONL 用 Python 逐行解析，扫描全部 2301 行中所有 user 文本，排除 tool_result；按 track、启动检查、precompile/AOT、文档精简与模型 PR review 回读助手上下文。
- :1965 的压缩摘要仅用于定位，不计为新的任务或频次；task-notification 不计用户请求。未读取其他 reviewer 结果；未创建 skill。
