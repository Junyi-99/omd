## 3. TPU 申请失败分层诊断与双队列交付（候选）

- 触发：配额看似够用但 TPU 创建失败，或希望 spot/on-demand 谁先可用就用谁。
- 可复用步骤：确认有效账号/project且显式传参 → 列出现有节点和 queued resources → 按 zone 核对 accelerator 配额、实际容量、tenant 初始化/权限错误 → 保留原始错误并区分证据与推测 → 用户授权后创建所需模式 → Python 轮询终态并查残留 → 若已有清理授权，赢家就绪后删除另一请求 → 报告实际机器、SSH、计费资源及实验 controller 配置。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/256f3b72-c80d-4dd2-8616-c894f319bfd0.jsonl:9`：“哪个账号…project id”；`:210`：“不要切换…有什么问题？”
  - 同文件 `:249`、`:265`、`:282` 分别追问直接 on-demand、配额与失败原因；`:299`：“同时提交 on-demand 和 spot 两个排队请求”。
  - 同文件 `:261` 纠正先前“region”推断；`:295` 区分容量、零配额、tenant project 创建失败，并将 Google 侧权限原因标为推测。
  - 同文件 `:309` 使用 Python 轮询；`:313` 报告尚未获得机器及两个队列状态。
- 频率/信心：1 session，多个追问服务于同一次申请；7 次 loop 命令记录属于同一监控任务，不计 7 次复现。方法合理，交付闭环未在已分配日志中证实。
- 建议 skill description：当用户需要申请 TPU 或排查排队/创建失败时，核对身份和项目，按 zone 区分配额、容量与初始化错误，并跟踪授权请求到可连接或明确失败。
- 边界及反证：不能硬编码历史项目、可用 zone、配额和 runtime；不能从 tenant 错误直接推出组织权限原因。双申请会额外占用资源，必须依据本次授权，不能作为默认策略。


## 1. 对齐参考实现与执行阶段，再做逐层/逐分片数值定位（优先候选）
**触发**：模型移植后能运行但输出错误，或单算子测试通过而端到端不对。
**步骤**：固定输入与 prefill/decode 阶段；以独立实现为参考，逐层比较中间量；遇并行张量进一步按 shard 核对值与语义归属；先恢复正确性，再测算子和端到端性能，并区分局部加速与整体收益。

## 2. 将带标注 profiling 与真实性能计时分离（专门候选）
**触发**：需要查看 Pallas TPU kernel 子步骤，或 trace 为空、归因消失、文件异常膨胀。
**步骤**：记录 JAX/libtpu/xprof 版本；在 import 前设置 tracing 环境；先用小 workload 验证 trace 非空及 lane；在已知受限版本每个 shape 使用新进程；保留必要 lanes 并压缩导出；用无标注 build 单独计时，避免把 instrumentation 开销当 kernel 性能。

## 3. Spot 实验的环境恢复与资源清理闭环（可合并候选）
**触发**：抢占式加速器实例重建，或临时创建昂贵 GPU 进行验证。
**步骤**：先核实当前项目、机型、配额 API、区域与镜像；环境按稳定依赖/易变源码/模型分层缓存；恢复后校验绝对解释器路径、editable 包解析位置及版本；长恢复任务在远端脱离 SSH 执行并以 sentinel 判断结束；实验完按用户意图停机或删除，再检查磁盘等残留资源。
为抢占式 GPU/TPU 实验建立可恢复环境，断线安全地恢复大文件，并在结束后核查残留计费资源。
配额、价格、区域顺序和硬件可用性是动态事实；原日志 assistant 自行改用另一硬件型号不应成为 skill 默认授权。删除资源需已有用户授权，不能按一律清理规则删除用户数据。


## 1. 基准真实性审计：先证明输入和变量确实在测目标问题

- 触发场景：算子微基准异常快、不同 batch 数据全平、重复查表或链式 GEMM 的峰值用于工程决策。
- 可复用步骤：
  1. 列出被测变量、真实 shape、实际访问量和计时边界，确认循环变量进入输入构造。
  2. 检查重复输入是否命中缓存，链式计算是否退化为零、inf 或 NaN；每轮换冷数据并保持数值有限且非平凡。
  3. 用均匀随机、分布代理和代表性文本作对照，报告代理的偏差及局限。
  4. 扫描工作量与线程数，分开计量 gather、传输、计算，再测重叠；区分实测和瓶颈推断。
  5. 明确撤回失效数字，不把重复测同一工作负载当成 batch 规律。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu/089f8c4e-19b8-4f86-8dde-bd5e876d1aaa.jsonl:984`，用户追问“不同 bs 之间的 lookup 是随机的吗”；`:993` 承认“B 根本没被用到”“同一件事测了六遍”。
  - 同文件 `:978`：“每轮复用同一组 id……全在 L3 里”，冷行重测使单线程结果变化 2.3 倍；`:1016` 比较 uniform、Zipf、real text，并披露词级小词表高估复用。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/memory/mxu-speed-data-dependent.md:11`–15：链式 dot 溢出后约 908 TFLOP/s，有限数据约 847–853；建议保持 finite/non-trivial。其功耗/时钟解释只是推测，不应写成机制事实。
- 频率/信心：2 个项目中的独立证据，日志 1 session + memory 1 份；方法价值高，复现频率中等。
- 建议 skill description：审计加速器与内存微基准的输入真实性、缓存状态、数值退化和工作量定义；通过对照实验修正误导性的性能结论。
- 边界及反证：硬件特定 TFLOP/s、线程阈值和 Zipf 偏差不应固定进 skill；日志中的“最多高估20%”只适用于当次代理，不能推广成保证。

## 2. 大模型特性移植的语义与生命周期对照（候选）

- 触发场景：把已有 GPU 推理特性移植到 TPU，已有 shape 单测通过但尚未加载真实权重。
- 可复用步骤：
  1. 先读论文、真实 config 和两个参考实现，按 hash、dtype、gate、卷积、注入逐项列等价关系。
  2. 追踪 state 从 config、pool、batch 到 backend 的生命周期，覆盖 prefill/decode、chunk、前缀复用、复制与 speculative verify。
  3. 用 checkpoint 自带的小型 buffer 和 header 核实推导；不必先下载整张大表。
  4. 验证数值类型在 host→device 边界的真实语义，不只验证存储位宽；对原始位模式采用位级比较。
  5. 把缺陷、等价但有意的差异、未支持能力分开报告，再运行真实表。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu/089f8c4e-19b8-4f86-8dde-bd5e876d1aaa.jsonl:1566`：“对照你的代码……等肉眼代码通过了之后我们再跑真实的表”。
  - 同文件 `:1563`：真实 safetensors header 揭示 shard 应切 padded 高度；hash buffers 不匹配则抛错，以阻断“合法但错误的行”。
  - 同文件 `:1667`：发现 gather 的 uint16 位模式被当整数上传，以及 mmap 缺 `MADV_RANDOM`；另指出 host 的 `fill_ids` 无法覆盖设备 draft token 的 target verify。
- 频率/信心：1 session 的连续移植过程；方法可靠，但没有第二个独立任务证明常用，列候选。
- 建议 skill description：移植推理模型特性时，对照论文、参考实现与真实 checkpoint，审计数学、dtype 和请求状态生命周期，识别静默错误及能力缺口。
- 边界及反证：不能把 SGLang 的全部文件后端、预取和 RSS 管理复制到 TPU；日志也明确拒绝了不适用的复杂度。MTP 当时仍未实现，不能据此声称已验证端到端支持。

## 3. 先做容量与数据通路的最小硬件验证（候选）

- 触发场景：超大只读表无法常驻 HBM，设计依赖 host memory、异步传输或跨设备访问能力。
- 可复用步骤：拆分共享权重、每请求 state、KV 和临时激活的容量账；用最小程序实测 memory kind、索引所在位置和 gather 支持；分别测 CPU gather/H2D/compute 及重叠；按实际调度器验证预处理能否提前；在接完整 runner 前确定可行通路。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu/089f8c4e-19b8-4f86-8dde-bd5e876d1aaa.jsonl:776`、`:824` 为中断后重复同一要求：“先验证显存账和 pinned_host gather”，仅计一次。
  - 同文件 `:945`：pinned_host 不占 HBM，但 gather 两种索引位置均失败；独立测量异步 H2D 和计算重叠，进而改变 batch prep 设计。
  - 同文件 `:978`：纠正“95.4 GiB 全局表”与“180 KiB 每请求 conv state”的混淆。
- 频率/信心：1 session，实测方法清晰，持久化优先级低于模式 1；可作为同一 skill 的容量/传输检查模块。
- 建议 skill description：在加速器推理系统接线前，通过容量账和最小硬件探针验证 host offload、gather 与异步传输方案。
- 边界及反证：当次 pinned_host 错误不代表所有 TPU/JAX 版本永远不支持；DMA 可重叠不等于生产 scheduler 已隐藏传输；单 host 测量不证明多 host 可用。


## 1. 优化前建立独立正确性门禁
- 触发：改写数值 kernel、改变精度或将实验冠军移入已有框架。
- 可复用步骤：先找外部原始实现；确认参考内部 dtype，必要时构造真正高精度参考；用同输入比较参考、候选、上游完整流水线；加入极端输入、有限值断言与 API 契约检查；仅对全部过关候选计时，报告误差包络与速度。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-lower-bound-oracle.md:17–24`：“不把自己写的代码当 ground truth”，交叉验证 FLA、fp64 与完整上游流水线；原实现内部 fp32 cast 会使伪 fp64 参考失效。
  - 同文件 `:56–57`：“必须 isfinite 守卫”；NaN 聚合可能造成阈值比较静默通过。
  - 同文件 `:110–116`：“先精度后速度”，“全局 default…被关卡否决”。
- 频率/信心：同一项目 memory 中 3 组明确证据；独立原始 session 数不可核实。方法信心高，跨项目频率证据不足。
- 建议 skill description：当优化数值 kernel 或降低精度时，构造独立参考、极端输入与有限值门禁，先验证正确性和接口契约，再筛选性能冠军。
- 边界及反证：并非每次代码调整都需要新建高精度 oracle；容差由应用决定，不能把此项目的固定阈值泛化。上游官方 kernel 也可能低于 naive 参考精度，不能仅凭“官方”认定真值。

## 2. 用上限探针、消融与真实形状验证硬件猜想
- 触发：怀疑小 GEMM、VMEM、布局或硬件利用率限制性能，计划打包、融合、改 tile。
- 可复用步骤：先阅读既有实验档案；记录真实 shape、dtype、VMEM 足迹与吞吐；建立预排布/驻留数据的理想上限探针，再测包含构造与搬运的版本；通过精度对照确认测试确为计算受限；做阶段消融，并以实际阶段占比估算整体上限；负结果也记录适用边界，避免把 microbench 收益直接外推到 kernel。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/memory/research-goal-vmem-mxu.md:11–15`：“record each config's VMEM footprint alongside throughput”。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/kda-lower-bound-oracle.md:88–93`：首版 microbench 内存受限，“必须先验证 bf16/fp32 分离再谈计算结论”。
  - 同文件 `:130–151`：小 tile 预排布打包 5.03×，移入真实 kernel 仅 1.00–1.03×；concat 吞噬收益，阶段占比给出 Amdahl 上限；16×16 打包被实测否决。
  - 同文件 `:153–158`：“per-pass ceiling 推算 ≠ 实测”；fat-K 正确却变慢。
  - 同文件 `:110–112`：“动手前读战役存档”，已有实现导致重复发明。
- 频率/信心：2 个不同项目 memory 支撑实验框架，其中详细负结果集中在 KDA；至少 3 类独立技术实验（fat-K、两类打包、microbench 校验），不等于 3 个 session。方法信心高。
- 建议 skill description：当研究 TPU kernel 的 tile、VMEM、布局或矩阵打包优化时，用真实形状上限探针、计算/内存判别和阶段消融验证猜想，记录失败边界再决定是否移植。
- 边界及反证：具体 MXU 尺寸、版本 flag 和最佳 hb 是设备事实，不能写成普遍规则；理论 ceiling 能指导假设，不能代替实测。大 profiler capture 曾挂死主机，小抓取成功也不能推定大抓取安全（同文件 :38–46、:104–108）。


## 1. 在真实目标硬件上做基线、消融与瓶颈归因（推荐主 skill）
- 触发：要求融合算子、移植其他框架优化、比较优化前后性能。
- 可复用步骤：明确目标设备及 host/device 执行位置 → 固定环境、输入形状、内存放置和 donation 等运行条件 → 保存旧版本基线并校验数值 → 在目标机测量 → 分项消融和形状扫描 → 检查编译产物/trace → 记录重复运行、复现命令及适用边界；收益与风险不清时先交付实验，不顺手删除共享正确性保护。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/2f92f0d6-1c91-4c6d-be19-6471a4a25a6b.jsonl:74` 用户：“记录一下 fuse 前后的性能差异，完整记录到 markdown 文件里。”
  - 同文件 `:375` 用户：“你在 tpu 上跑啊。怎么在我 m1pro 上跑？”；`:652` assistant 更正目标机与 M1 的收益差异。
  - 同文件 `:728` assistant：“conv 的开销跟 num_slots 走，不跟 batch 走”，通过移除 barrier、keep-mask、write-back 分别定位开销；同时注明 barrier 涉及多机竞态，未修改。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/memory/laya-bench-project.md:15`：“eager baselines must run before installing it”；`:17` 要求延迟测试前检查空闲 GPU 并记录 load。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/ici-v6e-measured-constants.md:14` 区分 host dispatch 与 device 时间；`:21-24` 明确无 N sweep、缺 matched baseline 等推断边界。
- 频率/信心：1 个完整任务 session + 2 个不同 origin session 的 memory 摘要；跨 TPU/GPU 复现，方法信心高，独立用户请求频率仅能确认 1 次。
- 建议 skill description：当用户要求加速器算子优化或性能对比时，在目标硬件构建可复现基线、正确性校验、消融与 trace 证据，交付带环境和因果边界的性能报告。
- 边界及反证：不是“看到 GPU fused kernel 就照搬”。该日志中 XLA 已融合 elementwise，实际瓶颈在状态写回；用户真正要 Pallas kernel，而先交付 numpy 算法融合，说明开始前需澄清优化层次。固定收益数字和主机清单不是 skill。

## 2. 通过最小能力探针验证编译器支持（适合作为主 skill 子流程）
- 触发：文档与实机行为冲突，涉及 dtype、memory space、后端 lowering 或升级能否解除限制。
- 可复用步骤：从真实算法抽最小算子链 → 区分输入/输出和中间 dtype → 显式切换相关配置 → 与可信参考实现数值比对 → 隔离会 abort 的 probe → 在独立环境复测候选版本 → 记录版本、错误方式、升级后的重测判据；通过正确性后再评估端到端传输收益。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-jax/2f92f0d6-1c91-4c6d-be19-6471a4a25a6b.jsonl:949` 用户粘贴支持所有整数精度的文档，`:987` 问“你要跑哪个 probe？”
  - 同文件 `:983` assistant：“x64=off 时它编译通过但算错”；`:1036` 给出递进五个 case，要求与 numpy int64 对数；`:1072` 两个版本均失败，保留明确重测条件。
  - 同文件 `:848` 区分 device/device、pinned_host/device、host/host gather，最后一种 abort，probe 不纳入普通测试套件。
- 频率/信心：1 session 中 2 类能力探针（整数精度、memory space），并含版本复测；不是多 session 习惯。方法可信，独立成 skill 的频率证据中等偏弱。
- 建议 skill description：当加速器后端的文档支持与实际编译或数值行为不一致时，构建隔离的最小能力探针，跨配置和版本校验数值，并给出可重测的结论。
- 边界及反证：不能把“当前 Pallas 不支持 int64”硬编码成永久规则；历史中“能编译”恰是错误结论来源。硬件、库版本和内存规模变化都需要重测。

## 3. TPU 创建失败的跨 zone 分类诊断（候选）
- 触发：节点创建失败、用户质疑配额/区域判断，或需在指定规格和计费模式下找可用资源。
- 可复用步骤：确认项目、规格、on-demand/spot 与接口支持 → 查询可用 zone → 按 zone 保存完整创建/排队错误 → 分开记录配额、即时容量、接口不支持、项目侧错误 → 必要时用第二接口交叉验证 → 核验实际 schedulingConfig 与节点状态 → 停止时清理 queued resource 并确认无残留。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/9a435ff1-bcf3-4a9d-97ec-dc0eb8457bc5.jsonl:13` 用户要求创建 on-demand TPU；`:146`：“你扫描一下看看我们哪些区域有？”
  - 同文件 `:142` assistant 过早断言所有区域配额为零；`:431` 扫描后发现 5 个 zone 有配额；`:495` 更正“卡的是容量”，并确认清理残留。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/memory/tpu-spot-sgljax-env-cache.md:17-18` 区分 no capacity 与 tenant project creation，并强调完整落盘日志而非接 tail。
- 频率/信心：1 个完整 provisioning session + 1 个不同 origin session 的 memory 摘要；分类诊断有复现，完整流程仅单次，候选信心中等。
- 建议 skill description：当云 TPU 创建或排队失败时，按 zone 和计费模式收集原始错误，区分配额、容量与项目/接口问题，核验实际资源状态并清理未完成请求。
- 边界及反证：不能从少数区域失败推导项目全局无配额；单一区域 tenant 错误也不足以证明是容量。区域、配额与可用性是动态事实，需查询。
