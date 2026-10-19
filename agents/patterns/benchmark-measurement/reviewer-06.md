分类：`performance-benchmarking`
方向：从异构推理实测生成可追溯的性能、成本与容量报告。
领域：GPU/TPU 推理基准、profiler、技术报告。
频率：仅 1 个独立日志 session；其中出现多轮报告补全、成本口径追问和利用率诊断，不能算跨 session 复现。
持久化价值：有明确的多步骤工作流和可检查产物；最有价值的是避免“新增硬件只补一节”和“忙碌率等于有效算力利用率”两类重复失误。当前适合候选 skill，尚不足证明是用户跨项目常用方法。
范围与局限：完整读取 4 份 memory；流式扫描 1 份 JSONL 的全部 user 文本，过滤工具结果，并核对相关 assistant 上下文。日志里的 task-notification 是系统回传，未计为用户需求。memory 所链接的其他文件未分配，未扩展读取；未验证历史代码与数字正确性。

## 1. 新增硬件时按指标矩阵补全异构性能报告（优先候选）

**触发场景**：已有 GPU 基准报告加入 TPU、不同后端或新硬件。

**可复用步骤**：
1. 先盘点报告全部指标与图表，建立“硬件 × 模型 × batch × 实现 × 指标”覆盖矩阵。
2. 固定输入、计时边界、预热/重复次数；记录 Python、框架、驱动、checkpoint、harness 版本与 dirty 状态。
3. 同时采集端到端、forward、设备时间、HTTP 与首次编译时间；后端不可比或缺测的项明确标注，不补造数据。
4. 新硬件同时进入 setup、基线、吞吐缩放、设备/非设备时间等公共图，不只增加独立章节。区分 GPU eager 与 TPU JAX jit。
5. 从原始结果重新生成报告，检查图表覆盖和无新硬件数据时的退化路径，并明确权重差异与未做 parity 的限制。

**原始证据**（以下 L 均为绝对路径对应 JSONL 行）：
- `/Users/junyi/.claude/projects/-Users-junyi-laya/889c74cd-ee30-4086-b29e-673be378e0ec.jsonl:187`：“把 TPU 这个加到报告里。”
- 同文件 `:318`：“GPU 和 TPU 的 latency 啥的画到一起”。
- 同文件 `:390`：“Setup section 里为什么没有 tpu 的？”
- 同文件 `:465`：“Device vs non-device time, HTTP latency per variant，也都没有 tpu 的图呢？”
- 同文件 `:606`：“Baseline in-process, eager 里面也要放 tpu 的（tpu 对应 jax）”。
- 同文件 `:184`（assistant）：记录 checkpoint 不同、未逐个 parity、harness 版本不准；是应改进的测量边界，不是已保证公平比较的证据。

**频率/信心**：1 session 中 1 次集成需求、4 次关联补全/纠正；主题强，跨项目频率未知。

**建议 skill description**：当向现有推理基准加入新硬件或后端时，按指标覆盖矩阵补齐测量、版本元数据和公共图表，并标注可比性限制，生成可复查报告。

**边界及反证**：这是一次连续修订，不是 5 次独立工作流。画图优先、表格折叠是本 session 的呈现要求（`:878`），不能默认推广为全局偏好。没有 HTTP 测量就应显示缺测，不能硬把 TPU 数据塞进不适用指标。

## 2. 从实测吞吐推导成本与容量，并公开假设（候选）

**触发场景**：用户从 latency 比较追问 tokens/$、每日请求规模、设备数和平均利用率。

**可复用步骤**：
1. 确认实测输入 token 数及吞吐单位；写清采用 eager、固定实现还是每点最快实现，保留所选 variant。
2. 为小时价格记录来源、采集时间、地区/型号和按需/抢占式口径，分别计算各价格情景。
3. 用 `tokens/call ÷ seconds/call × 3600 ÷ hourly_price` 算 tokens/$；注明串行实测不等于已测并发服务极限。
4. 将每日请求量和每请求 token 数换成平均 token/s；由单卡实测吞吐推忙碌卡时与台数。
5. 固定集群按峰均比和目标利用率配备，另列理想弹性成本；展示敏感性分析，而非把假定流量曲线当观测。
6. 图表同时暴露吞吐实现、价格情景、batching 假设与尚未实测的并发/排队条件。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-laya/889c74cd-ee30-4086-b29e-673be378e0ec.jsonl:656`：“加一个 token/$ 的页面”，并要求区分按需和抢占式。
- 同文件 `:821`：“你用什么算的？是 eager 模式，还是 cudagraph 模式？”
- 同文件 `:828`：“如果 workload 跑不满硬件呢？……多少硬件 surve 多久，平均利用率多少？”
- 同文件 `:824`（assistant）：逐硬件/模型/问题数取最快实现，说明必须公开口径。
- 同文件 `:831`、`:874`（assistant）：采用假设峰均比 2、目标利用率 80%，给固定与弹性方案；均是推算而非线上容量验证。

**频率/信心**：1 session；1 次成本需求、1 次实现口径追问、1 次容量扩展。可复用结构中等偏强，跨 session 证据不足。

**建议 skill description**：当已有推理基准需要估算 tokens/$ 和生产容量时，统一 token、延迟和价格口径，计算固定/弹性集群规模并公开峰值、batching 与利用率假设。

**边界及反证**：不持久化硬件报价、请求规模或历史排名。历史回答曾把设备忙碌当“基本满载”，而后续 MFU 分析显示有效利用率很低；因此该容量模型只是情景估算，不能直接承诺生产 SLA。

## 3. 从设备忙碌追到真实 token 的有效算力（候选）

**触发场景**：用户问硬件是否打满、HBM 读写是否频繁，或准备优化 TPU kernel。

**可复用步骤**：
1. 分开看主机等待、设备空闲、算子时间与矩阵乘效率；不以低 idle 认定计算峰值已打满。
2. 从 tokenizer 实测逐行长度，区分 state、问题、sequence padding 与 batch padding。
3. 报告 padded 与有效 token 的计算利用率，说明 FLOP 估算、峰值规格及 profiler 扰动。
4. 将热点映射回源码与 dtype，区分已测瓶颈和待核实猜测，再选择融合、flash attention、batch 扫描等实验。
5. 在可 review 的隔离 worktree 中保留 kernel 与实验脚本；正确性和性能验证完成前不宣称优化成功。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-laya/889c74cd-ee30-4086-b29e-673be378e0ec.jsonl:1185`：“v5e 和 v6e 性能打满了吗”。
- 同文件 `:1204`：“50个 question 的 token 数量分布，你画在 report 上”。
- 同文件 `:1256`：要求检查 HBM/VMEM、算子融合与 flash attention，并要求“开个 worktree……我要 review 代码”。
- 同文件 `:1195`（assistant）：区分含 padding 与真实 token 的 MFU，并承认 fp32 搬运解释未逐算子验证。
- 同文件 `:1247`（assistant）：逐行分解 state、问题及 padding，指出基准只有两种长度。

**频率/信心**：1 session 的 3 个连续诊断/优化需求。诊断流程中等信心；优化落地证据不足，日志以“无需响应”结束，不能声称 kernel 已实现。

**建议 skill description**：当诊断 TPU 推理是否打满或准备 kernel 优化时，结合 profiler 与真实 token 分布区分 idle、padding、HBM 和 compute 瓶颈，输出带证据的实验计划。

**边界及反证**：用户在本次明确排除减少 padding 和 encode/decode 流水化（`:1256`）；不能把通用建议凌驾于实验范围。两个题型交替的合成负载不代表生产长度分布。

## 不建议单独沉淀的 memory

- `/Users/junyi/.claude/projects/-Users-junyi-Desktop-phd-qe-qe-slides/memory/feedback_scripts_over_inline.md:7-11`：脚本先写文件、由用户 review/执行，是 1 份偏好反馈；可作项目约束，独立 skill 证据不足。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-medusa-reborn/memory/feedback_ambiguous_deletion.md:10-25`：删除“其他”不得扩展授权、先列大小与恢复能力，是 1 次严重事故后的安全规则；宜保留通用行为约束，不伪装成频繁工作流。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda-aug-12-learn/memory/fallback-vm-deepseek-v4-data.md:11-16`：临时备用机与挂载点是时效性基础设施事实；不复制地址，不作为 skill。设备能力先探测可并入基准准备步骤。
- `/Users/junyi/.claude/projects/-Users-junyi-tmp-humanize-demo-original-performance-takehome/memory/MEMORY.md:1-3`：dirty kernel、代理指标和性能下限均为项目索引；未读链接正文，不能据此推广优化方法。

## 覆盖

- memory：manifest 全部 4/4，全文读取（含元数据）。
- logs：manifest 全部 1/1；流式扫描所有 user turns，核实上述相关 assistant 文本；未逐条读取无关工具输出。
- 未创建 skill，未读取其他 reviewer 报告。
