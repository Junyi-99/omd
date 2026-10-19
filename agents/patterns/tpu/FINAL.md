# Kernel 正确性与硬件验证

先做正确性门禁和最小能力探针。

## 1. 正确性门禁

`kernel-correctness-gate` · 优先

- **什么时候用**：优化数值 kernel、改变精度，或移植模型。
- **怎么做**：选独立参考 → 核对真实 checkpoint、dtype 和布局 → 逐层、逐阶段、逐分片找首个偏差 → 检查有限值和误差 → 通过后计时。
- **频率**：1 个原始会话 + 2 个不同来源会话的 memory。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-lower-bound-oracle.md:17)：“不把自己写的代码当 ground truth”。

## 2. 最小硬件能力探针

`accelerator-capability-probe` · 优先

- **什么时候用**：方案依赖尚未确认的 dtype、内存空间或异步通路。
- **怎么做**：列容量账 → 抽最小算子链 → 固定配置 → 在目标硬件编译、执行并对数 → 保存版本、错误和复测命令。
- **频率**：2 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu/089f8c4e-19b8-4f86-8dde-bd5e876d1aaa.jsonl:776)：“先验证显存账和 pinned_host gather”。

## 3. 云实验资源管理

`cloud-accelerator-experiment-lifecycle` · 配方

- **什么时候用**：创建云加速器、恢复 spot 环境，或结束实验。
- **怎么做**：按项目、区域和规格定位资源 → 区分配额、容量和控制面错误 → 恢复缓存并核实解释器及源码 → 跟踪完成状态 → 核对实例、队列和磁盘。
- **频率**：创建与清理见 2 个会话；恢复见 2 份不同来源 memory。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/9a435ff1-bcf3-4a9d-97ec-dc0eb8457bc5.jsonl:146)：要求扫描区域，随后发现瓶颈是容量。

详细记录：[reviewer-18](reviewer-18.md)、[reviewer-19](reviewer-19.md)、[reviewer-24](reviewer-24.md)、[reviewer-26](reviewer-26.md)。

## 4. 新模型接入审查

`inference-runtime-contract-review` · 优先

- **什么时候用**：模型新增共享状态、batch 字段，或接入 warmup、AOT、缓存。
- **怎么做**：列全部构造入口 → 跟踪分配、初始化、快照、复制和释放 → 逐路径检查 shape、dtype、分片 → 提前拒绝不支持组合 → 验证完整模型。
- **频率**：同一项目 3 个会话，其中 2 个讨论同组 PR 评论。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/c510a35c-9034-47b4-b49e-f9ef5f4bb239.jsonl:1473)：“应该报错的。不应该默认喂0”。

## 5. 上游差分验证

`upstream-differential-validation` · 模块

- **什么时候用**：把上游算子或测试移植到新运行时。
- **怎么做**：定位对应实现和用例 → 搬迁空输入、padding、多请求等边界 → 核对数值适用域 → 对照阶段输出 → 验证完整集成。
- **频率**：1 个会话；并入接入审查。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/aa53935d-b777-4bf5-9a88-877f0a7e0b96.jsonl:1123)：“我们不要自己造了，我们按照他们的来模仿参考”。
