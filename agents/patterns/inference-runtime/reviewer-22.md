## 1. 新模型接入时审计共享状态契约，同时限制改动范围
- 触发：新模型向现有 request pool/batch 添加状态；用户要求不碰已有 backend，或问“这与我们的模型有什么关系”。
- 可复用步骤：
  1. 画出 host/device、scheduler/worker、batch/pool 的所有权与转换点。
  2. 从新字段沿实际调用链追到预算、分配、读写、快照、复制、释放；区分类型定义和模型实例配置。
  3. 查已有 backend 依赖的布局/位置约定；尽量保持其代码不变，只补新模型承担的契约。
  4. 对未支持的 overlap/spec/cache 组合明确拒绝；不要让合法索引掩盖错误状态。
  5. 输出修改归属和剩余集成限制，避免覆盖用户并行改动。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/aa53935d-b777-4bf5-9a88-877f0a7e0b96.jsonl:691`：“我们这次只做 qwen，我不想碰 gdn 和 kda”。同文件 `:740` 允许短 TODO；`:1318` 记录 backend 回退、新旧修改归属及未完成真表验证。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/b2b4854a-f7b2-48af-a1e3-9a5f0144688b.jsonl:137`：“我们不去修改 GDN 的 cache 机制。和我们有什么关系？”；`:150` 解释 pool 整体复制所有 conv buffer，PLE 未写 track 快照会复制过期状态，并给出支持或启动拒绝两条路径。
- 频率/信心：2 独立 session、同一项目的重复边界问题；方法可信度高，跨项目普遍性中等。
- 建议 skill description：当为推理框架接入新模型状态或 batch 字段时，追踪共享资源的完整生命周期与布局契约，在不重构现有 backend 的前提下识别集成缺口，并验证或拒绝不支持的运行组合。
- 边界及反证：“不碰 GDN”是本任务范围，不能固化成所有任务的禁令。后一个 session 是理解评审意见，用户明确要求先别改（`:70`）；技能应允许只读解释模式。独立审阅代理只出现一次，不足以强制每次委派。

## 2. 迁移上游用例并核对数值适用域，而非只造本地测试
- 触发：实现与 SGLang/vLLM/Transformers 等已有实现应语义等价，尤其 hash、卷积状态和 embedding gather。
- 可复用步骤：
  1. 定位真正对应的上游函数和现有 testcase，记录 dtype、参数范围与调用模式。
  2. 迁移命名边界例：空请求、padding、多请求、奇数维、输出缓冲复用等；再核对 prefill/decode/chunked 的同位置一致性。
  3. 不机械照搬数值边界：明确 signed/unsigned overflow 等差异，将“不可能进入该输入域”的论证写成断言。
  4. 先用轻量环境检查纯逻辑，再在目标设备验证精度和分片；报告 skipped 与尚未覆盖的真实资源端到端路径。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/aa53935d-b777-4bf5-9a88-877f0a7e0b96.jsonl:861`：“和 vllm/sglang 一起核对一下 ngram 实现是否正确”。
  - 同绝对路径 `:1123`：“我们不要自己造了，我们按照他们的来模仿参考”。
  - 同绝对路径 `:1293` 记录 vLLM case 表、overflow 适用域断言、SGLang gather 用例，以及 TPU 的“63 passed, 1 skipped, 14 subtests passed”；`:1318` 同时声明真表端到端未验证。
- 频率/信心：1 session、2 次推进同一验证任务，不算独立复现；候选，方法清晰但频率证据偏弱。
- 建议 skill description：验证推理算子移植时，优先迁移上游边界用例并进行参考实现差分，显式检查数值适用域，再分层执行 CPU 与目标设备测试并说明覆盖限制。
- 边界及反证：上游用例不代替集成测试；后续 PR 审阅日志 `b2b4854a-f7b2-48af-a1e3-9a5f0144688b.jsonl:132` 指出单层测试不构造完整模型、不安装表，因此未暴露集成问题。建议作为 runtime 验证技能的一部分，尚不单独建 skill。

## 3. Cache/allocator 评审采用逐步守恒压力测试（高价值候选）
- 触发：cache 所有权、tombstone revival、分页容量或锁引用逻辑发生修改。
- 可复用步骤：建最小请求生命周期模拟；每一步检查树拥有的 slot 与 allocator 空闲/占用是否互斥且守恒；统一 token/page 单位；覆盖调用方 finished/unfinished 路径及拆分/复活；固定随机 seed；用主分支红、修复绿证明回归。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/sglang-jax-pr1564-review.md:15`–20：共享复活路径重复 free、锁回执 split 后失配、token/page 单位混用。
  - 同绝对路径 `:33`：“main 上红、补丁后绿（61/61）”；`:38`–41：“每步查守恒 + 树/分配器一致性”，以及模拟时单独保留 `prefill_target`。
- 频率/信心：1 memory 描述 1 次 PR 审阅及后续修复，不把 40 seeds 或多天追踪算作多次需求。技术价值高，复现频率不足。
- 建议 skill description：审阅 cache/allocator 修改时，以请求生命周期模拟逐步核对资源所有权、容量单位和树/分配器一致性，构造可重复的故障种子与修复前后回归证据。
- 边界及反证：memory 是历史摘要，原始执行日志不在分配范围；不可直接相信或固化其中依赖版本、脚本位置。与 pattern 1 有相关性，但尚无第二次独立执行守恒测试的证据。

## 不应技能化
- `code-block-blank-lines.md:8`–11：按逻辑块留空行是代码阅读/格式偏好，适合规则或 memory；`git diff --ignore-blank-lines` 可作轻量检查，单独不构成 skill。
- `ici-experiment-3-small.md:8`–9：机器容量、资源占用、实测常数是项目事实。“不同 grain/flow kernel 曲线不拼接”有方法价值，但只有一次摘要；“不改全局 gcloud project”是共享环境约束。
- kitty/tmux URL 排查仅 1 次，且历史回复对 tmux 命令参数自承未验证；不据此提出成熟终端修复技能。
- HUD configure 日志主体是插件注入的现成操作说明，没有用户决策/重复工作证据；不将其当新技能需求。
- 查 Claude 历史目录只有 1 个用户问题，适合事实回答。

## 覆盖情况
- memory 全文：`/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/memory/code-block-blank-lines.md`（11 行）；`/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/sglang-jax-pr1564-review.md`（41 行）；`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/ici-experiment-3-small.md`（10 行）。
- 日志全部 user turns：`/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/aa53935d-b777-4bf5-9a88-877f0a7e0b96.jsonl`；`/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/b2b4854a-f7b2-48af-a1e3-9a5f0144688b.jsonl`；`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-architecture/4f0f595d-15fa-4bc0-9bb6-c76bbc5cc13d.jsonl`；`/Users/junyi/.claude/projects/-Users-junyi-claude/673e1f40-8170-4ea1-86df-d18c8956ddb5.jsonl`；`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/d16f724f-844c-4546-b83b-cb70142b9957.jsonl`。
