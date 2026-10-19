# 性能测试

先做测量口径、性能诊断、实验设计三个 skill。

## 1. 统一性能测试口径

`accelerator-benchmark-contract` · 优先

- **什么时候用**：比较设备、后端、batch，或给报告补新硬件。
- **怎么做**：固定输入、版本和工作量 → 分开预热、设备时间与端到端时间 → 等待异步完成 → 检查设备占用 → 保存样本并补齐比较矩阵。
- **频率**：5 个会话，涉及 Laya、互联、MPK。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/eba5c113-e6ae-4eb2-8f61-07795e767711.jsonl:47)：“先检测正确性”“记录每一次的 Latency”。

## 2. 诊断性能异常

`accelerator-performance-diagnosis` · 优先

- **什么时候用**：计时不一致、加速倍数异常，或需要解释优化机制。
- **怎么做**：核对计时代码 → 分开计时与插桩 → 用消融和重复扫描区分假设 → 检查 trace 完整性 → 用 HLO/MLIR 验证编译行为。
- **频率**：6 个会话，涉及 MXU、EPMoE、Laya。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-Downloads-pallas-kernel-targets/8c4c1120-fca7-4bd6-8334-67918f04937f.jsonl:537)：“你没查实的机制应该尽快 xprof 做 profile”。

## 3. 从研究问题设计实验

`performance-experiment-review` · 优先

- **什么时候用**：设计 sweep，或检查实验是否回答了最初的问题。
- **怎么做**：列问题、假设、变量、对照和指标 → 选真实工作量 → 验证控制变量 → 对照局部与端到端结果 → 逐问题整理结论和补测。
- **频率**：3 个会话、3 组研究。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/7dff67a6-e2bc-469e-a770-fffc93de1f63.jsonl:1521)：“解决了 design-2.html 提出的问题吗？”。

## 4. 估算推理成本与容量

`inference-cost-scenarios` · 配方

- **什么时候用**：已有吞吐数据，需要估算 tokens/$、设备数或利用率。
- **怎么做**：统一 token 和时间单位 → 记录报价日期 → 计算 tokens/$ → 按流量、峰均比和目标利用率推算容量。
- **频率**：1 个会话；并入性能测试 skill。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-laya/889c74cd-ee30-4086-b29e-673be378e0ec.jsonl:656)：“加一个 token/$ 的页面”。

## 5. 判断优化空间

`kernel-resource-floor` · 配方

- **什么时候用**：多次调优收益很小，需要判断瓶颈和剩余空间。
- **怎么做**：按指令计数与资源容量估算下界 → 拆启动、稳态和收尾 → 计算优化的资源转移 → 对照实测选择方向。
- **频率**：1 份项目 memory；并入性能诊断。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-tmp-humanize-demo-original-performance-takehome/memory/perf-takehome-floor.md:17)：记录资源下界和多次优化退化。

## 6. 检查 sweep 图表

`benchmark-sweep-plot-review` · 配方

- **什么时候用**：需要从 sweep 标出性能下降点并解释。
- **怎么做**：检查采样缺口 → 找可比样本的下降点 → 对照机制证据 → 渲染检查标签、连线和正文。
- **频率**：1 个会话；并入实验交付。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/5c668753-ed25-444d-b25f-df9808b7c1c1.jsonl:11)：“主要是标出来性能下降的点”。

详细记录：[reviewer-02](reviewer-02.md)、[reviewer-04](reviewer-04.md)、[reviewer-06](reviewer-06.md)、[reviewer-07](reviewer-07.md)、[reviewer-08](reviewer-08.md)、[reviewer-13](reviewer-13.md)、[reviewer-15](reviewer-15.md)、[reviewer-20](reviewer-20.md)、[reviewer-23](reviewer-23.md)、[reviewer-25](reviewer-25.md)。
