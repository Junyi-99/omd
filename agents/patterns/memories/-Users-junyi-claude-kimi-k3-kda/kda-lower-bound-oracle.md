---
name: kda-lower-bound-oracle
description: kda-lower-bound/ release 目录：K3 oracle + sub_chunk=16 kernel，v6e 实测 3.35×，两个 milestone 均完成
metadata: 
  node_type: memory
  type: project
  originSessionId: e645184a-6feb-471a-ae36-e740b711b0e0
  modified: 2026-08-10T13:31:31.217Z
---

`kda-lower-bound/` 是用户计划 release 的独立目录（2026-08-10 定），kernel 正式版本放里面
（`kda_pallas.py`，与 tpu/kda_pallas.py 同步的 canonical copy）。

**Milestone 1（完成）**：K3 correctness oracle（`test_k3_oracle.py`）：
g = -5·sigmoid(exp(A_log)·(raw_g+dt_bias)) → `kda_recurrent_ref` 为 oracle。

**独立 ground truth 锚定（2026-08-11，crossval/RESULTS.md）**——用户明确要求不把自己写的
代码当 ground truth：① 上游 FLA torch naive（原文）与 oracle 差 4.5e-8，双方对真 fp64
金标准都在 fp32 舍入地板 ~9e-8（注意：FLA naive 内部强制 cast fp32，直接喂 fp64 输入
得不到 fp64 参考，必须改源码里的 cast）；② FLA naive_kda_lowerbound_gate 与 kda_gate
差 7.2e-7；③ H200 上全 FLA 流水线（他们的 gate+naive，零行我方代码）与 oracle 差 1.0e-7、
与 chunk kernel 差 1.9e-6；④ FLA 官方 chunk_kda(safe_gate) 离它自家 naive 有 5.6e-4
（triton tl.dot TF32），我们的 kernel 比 FLA 官方 kernel 更贴近共同 naive 基准 ~300×；
⑤ sglang-jax vendor naive 3e-8。

**Milestone 2（完成 2026-08-10）**：`sub_chunk=16` 多参考点因子化解锁 C=64@lb=-5。
原理：decay 因子化 `exp(gc_i-R)·exp(R-gc_j)`，ref 精确抵消；单 midpoint ref 时指数
span C/2·|g|（C=64→160 溢出），每 16-token 子块一个 ref 后 kept 区域 ≤40。L/M 变成
C/16 个 [C,128]@[128,16] 列条带 GEMM，MAC 总数不变。实现在 `_chunk_body(sub=)`，
XLA/Pallas 共享。v6e-4 实测（B=1 T=8192 H=24 D=128，bounded gate）：
c64-s16=1.38×，+hb8=2.67×，+invbf16=**3.56ms, 3.35×** vs 旧 pallas-c32（11.91ms），
1.60× vs XLA baseline。最快配置 err ~1e-5，g≡-5 worst case 8.9e-8。

**xprof profile（2026-08-10，prof/PROFILE.md）**：Pallas op 本体 10.79→2.46ms（4.38×，
FLOPs +14% → 实际 FLOP/s ≈3.8×）。XLA 同体代理：GEMM 6.9→14.0 TFLOP/s，Neumann 逆
占比 12.3%→4.1%。⚠️ xprof 对 fused Pallas op 是黑盒（mxu_utilization 报 0，无硬件
counter lane；libtpu tensorcore_util 是整核口径）——单元级拆分只能走 XLA 代理或 stage ablation。
`--xla_xprof_enable_custom_call_tracing`（2026-08-10 实测）：libtpu 0.0.41 不认识该 flag；
0.0.43.1/0.0.44.1/0.0.45（配 jaxlib 0.10.2）接受 flag、Pallas 正常执行，但 **profiler 抓取
Pallas custom call 时必 segfault**；"nightly 无更新版本"是误判：jax 0.11.0（2026-07-16 发布）requires_python>=3.12，
py3.10/3.11 的 venv 会静默回退 0.10.2——建 venv 用 uv --python 3.12。
jax 0.11.0 + libtpu 0.0.44.1 下该 flag 小抓取（T=1024×3 iter）能过，但正式抓取
（T=8192×20×3 配置）**连续两次把整台 VM 搞到失联**（gcloud 仍报 READY，ssh 死）——
不只是进程 segfault，是 host hang。别再用这个 flag 跑大 capture，等更稳的 libtpu
（探测脚本 kda-lower-bound/vm_probe_tracing_flag.py）。实验在独立 ~/venv-prof 做的，
主 venv 已回滚 libtpu 0.0.41，未回写 GCS 缓存。

**附带发现**：旧 c32 midpoint 路径在 g≡-5 时 o_err~7e-4（e^±80 fp32 精度损失），
sub16 一并修复（因子 ≤e^40）。

**pkill 自匹配坑（一天踩五次）**：`ssh 'pkill -f PATTERN; ...'` 里 PATTERN 会匹配这条
ssh 自己的 bash -c 命令行，把自己杀掉、后续命令全不执行且无输出。`[x]` 括号技巧只在
同一命令行**不含**对应明文时有效（echo 文本/路径/参数里出现明文照样中招）。
铁律：kill 用独立 ssh 会话、按 PID kill、会话里不得出现任何与模式匹配的明文。

**坑**（复用于一切比较测试）：NaN diff 经 jnp.max 打印 -inf 且 `-inf<tol` 静默通过
→ 必须 isfinite 守卫；q/k 必须 L2 归一化否则 delta rule 发散；mask 用 where 不用乘法。

**4 卡 + B sweep（2026-08-10）**：H=96 shard_map 头并行。B=1 prefill：新 3.65ms vs
旧强基线 11.01ms = 3.0×（旧 hb16 因 24%16≠0 静默回退 hb1——head-block 必须整除 B·H）。
B≥2 后两者都 ~100% 扩展、饱和在 HBM roof，新 vs 旧收敛到 1.10×（2.39M vs 2.16M tok/s）。
结论：大收益在低并行度延迟敏感场景（单请求长上下文 prefill / 每卡少 head），吞吐饱和区 ~10%。
坑：bench 输入必须先 device_put+NamedSharding 预分片，否则每次调用的 resharding 传输
（~1.6GB）淹没计时；`tpu_vm.py status` 会在实际 PREEMPTED 时误报 READY，判抢占以
`gcloud compute tpus tpu-vm describe` 为准；"READY 但 ssh 死"的另一原因：`describe --format="value(state,health)"` 可能给出
READY + UNHEALTHY_MAINTENANCE——判活要看 health 字段，不健康直接 down/up；
spot 重建可能复用旧 IP → known_hosts 指纹冲突
使 scp "Connection closed"（accept-new 不覆盖冲突条目），修复：
`ssh-keygen -R <ip> -f ~/.ssh/google_compute_known_hosts`。

**MXU 利用率 & MXU-vs-VPU ablation（2026-08-11，prof/PROFILE.md + prof/lowering.txt）**：
HLO 证实 XLA/Mosaic 都不聚合 4 条 strip（各 strip LHS 行缩放不同，数学上不可 concat），
8 个独立 [*,64,16] dot。⚠️ **v6e MXU 是 256×256**（v2-v5 才是 128×128，用户纠正）：
[64,128]@[128,16] 的结构上限 = K·N/256² = 3.1%，实测 2.6% ≈ 上限的 84%——strip GEMM 是
shape-bound 不是 overhead-bound。同 kernel 只换 L/M 构造：VPU [C,C,D] elementwise vs MXU
strip：hb8 下 6.08→3.56ms = 1.71×（stage 级 ≈2.5-3×）。纯 JAX 写法不缓解。
下一杠杆（256×256 修正后 fat-K 结论翻转）：fat-K [64,512]@[512,64] 上限 25%（8× strips，
净省 ~2× MXU cycles，值得做）；s4 precision highest→default（6→1 pass）可叠加。

**hb sweep + FLOPS（2026-08-11，vm_bench_hb_sweep.py，H=96 T=8192）**：
hb 边际收益逐倍减半（1/2/4/8/16 → 28.9/19.8/15.0/13.2/12.3ms 单卡），hb32 VMEM 爆；
四卡最优 hb8=3.62ms（聚合 55.1 TFLOP/s）；单卡最优 hb16=12.33ms（16.2 TFLOP/s ≈
fp32-6pass roof 的 10.6%）。sglang chunk_kda_fwd 同形状 180.3ms 单卡（1.11 TFLOP/s）、
49.4ms 四卡——比我们慢 13.7×（其中 intra 构造 ~1.7×，其余是逐行 solve、无 hb、
四阶段间 HBM 往返）。VPU 口径（解析）：我们 VPU ~2G ops/iter（1%），sglang intra
构造 ~51.5G VPU ops/iter，实测推断 ~3.1 Top/s ≈ VPU 峰值 80%，VPU-bound。

**VPU-vs-MXU microbench（2026-08-11，vm_bench_dot_vpu_mxu2.py，VMEM 驻留 ×16 复用）**：
[16,128]@[128,16]：MXU fp32-hi 1.01 / bf16 3.47 / VPU 0.85 TFLOP/s（fp32 下仅险胜 1.19×）；
[64,128]@[128,16]（strip 实形）：MXU fp32-hi 3.40 / bf16 12.60 / VPU 1.15（3.0×/11×）。
3.40 与 kernel 内 s4 实测 3.95 交叉吻合。bf16 单 pass = s4 最大杠杆（~3.7×）。
坑：v1 bench 算术强度 4 FLOP/B 全员同速（内存受限），bf16==fp32 同时间即为铁证——
microbench 必须先验证 bf16/fp32 分离再谈计算结论。

**直接 A/B vs sglang 官方 kernel（2026-08-11，vm_ab_vs_sglang.py）**：sglang chunk_kda_fwd
原样（varlen API，K3 gate 在其 kernel 内激活）vs 我们 c64s16，同 tensor：o=2.9e-4/s=2.2e-4，
恰好等于 sglang kernel 自己离 fp32 oracle 的距离（其内部 bf16 中间量）；我们离 oracle 1.7e-7。
结论：一致，差异全部来自 sglang 侧自身精度。sgl_kda.py 单文件只依赖 jax，可独立跑。

**s4-bf16 + A/B + LLO 硬件追踪（2026-08-11，prof/PROFILE.md + prof/llo/）**：
s4_precision="default"（旋钮已加进 kda_pallas.py）：err 1.2e-5→2.1e-4（仍 <5e-4，与 sglang
包络同级），整 kernel 1.16×；新最优 c64s16+hb8+invbf16+s4bf16=3.07ms@H24（累计 3.88×）。
sglang kernel 只变 safe_gate：187.1→158.8ms=1.18×，gate 公式零成本，ValueError 契约 ✓。
**custom-call profiling 在 jax 0.11.0 + 小抓取（T=2048×3iter 单配置单进程）下可用**：
trace 出现 VALU/MXU0/MXU1/EUP/XLU/VLD/VST 指令 lane + _counters_ 利用率线 + 源码行 lane。
generic vs strip 指令数：MXU 21.6K→140.7K（+6.5×），XLU −74%，VLD/VST −74%，VALU −20%，
EUP(vpow2.f32=exp2) −47%——VPU→MXU 的硬件级实锤。坑：指令时长是机会采样（看 count 别看
时长和）；_counters_ 小抓取采到 0；大抓取仍会挂主机。原始抓取 1.8GB 在 prof/llo/。

**冠军合体（2026-08-11，规则=先精度后速度）**：发现 repo 已有 candidate_v1..v11 战役
（gorgon/examples/kda_chunk_pallas/RESULTS.md，8-01 Medusa run；v1 早已做 sub-block refs——
本轮部分重复发明，教训：动手前读战役存档）。v9=原生布局+块三角逆+log-step cumsum。
同台：in-spec v9-nb2 11.51ms（但 g≡-5 err 6.9e-4）/v9-nb4 12.86ms（err 1e-7，60/60 关卡+
GateOutOfRange 守卫验证）/ours 16.41ms。合体 candidate_v9s.py（v9+nb4+选择性 s4/inv bf16）：
**7.893ms @ worst 2.1e-4，全关卡通过——现任冠军**（累计 5.95×；vs 前代部署 1.85×）。
全局 default（3.99ms@1.8e-3）被关卡否决。vm_gauntlet.py=通用准入流程（先过关再计时）。
你记的"砍 exp"kernel=candidate_v5，当时判 no-op（Mosaic 已 CSE）。

**冠军细扫（2026-08-11，prof/champion_sweep.txt）**：T{2K..32K}×H{8,24,48,96}+B{1,2,4} 全 22 点
无失败，champion 全域赢 1.06-1.44×（随 B·T·H 单调放大）；T 严格线性（linear attention 兑现）；
饱和 25.9 TFLOP/s；B 与 T 等效；四卡分片格（H=24,T=8192）=2.47ms。

**blocked-layout 入口——1.36×，追平冠军（2026-08-12，vm_bench_blocked.py）**：
kda_chunk_pallas_blocked 接受预分块 [B·H,NT,C,Dp]（上游 QKV epilogue 可直接产出），
免全部 HBM layout：10.442→7.685ms（layout 税实测 2.76ms=26.4%，与 xprof 推算吻合），
与 v9s 冠军 7.893 打平（~噪声）。两线在去 layout 后第二次收敛——layout 是最后的大石头。
契约差异：v9 原生布局无上游要求；blocked API 需跨算子布局契约。wrapper 已改为路由过
blocked core，输出逐位同。sglang follow-up：其四阶段间 HBM 往返使 blocked API 收益更大。

**mfold 移植 kernel：精度无损、性能中性（2026-08-12）**：s4 四 strip 的 M-folding 打包
进 kernel（sub_chunk="mfold"，行向 reshape 摞叠免 lane concat）：精度与 sub-chunk 逐位同，
perf 仅 1.00-1.03×。原因：① kernel strip 是 M=64（12.6 TFLOP/s，轻度发射受限），不是
microbench 的 M=16（7.4，重度）——上限 ~2× 被 4× 交叉 MAC 吃掉；② s4 占比 ~15%，Amdahl
封顶 ~1.08×；③ hb 已摊薄发射。**适用边界definite：all-pairs 打包只救 M≤16 且操作数预排布
的发射饥饿 dot；M≥64+hb 已不在此区**。默认路径仍 sub-chunk，冠军不变（v9s 7.89ms）。

**[16,128]@[128,16] 打包实测成功——5.03×（2026-08-12，vm_bench_pack_tile.py，llo_tile/）**：
KDA 真实 tile 形状下 all-pairs 打包（A 沿 M 摞 [256,128]、B 沿 N 拼 [128,256]、单 GEMM 取对角）
= 0.92ms/37.2 TFLOP/s vs batched 4.65ms/7.4 = **5.03×**，数值逐位一致。LLO 归因：batched 是
**MXU 指令发射受限**（442K vs 194K 指令，后者 16× MAC——每指令有效 MAC ~36×）。致命细节：
kernel 内 lane 方向 concat 一项就吃掉全部收益（allpairs in-kernel=1.01×）——**打包只在操作数
预排布进 kernel 时成立**。与 16×16 实验（1.08×）对照：K=16 打包是零填充、K=128 打包是真胖
GEMM——形状决定成败。KDA 启示（未测）：s4 四条 strip 可 M-折叠打包（sublane 摞 q、dot_general
直接收缩、免 lane concat），MAC×4 换指令 ÷4。fp32-hi batched 在 R=16 下 scoped-VMEM OOM。

**16×16 block-diag 打包实测否决（2026-08-12，vm_bench_pack16.py）**：用户假设把 16 个
16×16 塞满 256×256 阵列可得 ~6×；实测 prepack 上限 1.08×、含构造 0.97×（数值逐位一致）。
机理：Mosaic 小 dot 调度远好于朴素脉动账（seq≈prepack）；打包大 GEMM 自身只跑峰值 ~6%，
发射/流水开销随总 MAC 放大，16× 零乘通胀抵消填充收益。bf16 vs fp32 仅差 13% → 该尺度
瓶颈不在 MXU pass。与 fat-K 同律：**padded-MAC 脉动账系统性高估小形状浪费，动手前先跑
prepack 式上限探针**。

**fat-K 实测否决（2026-08-11，vm_bench_fatk.py，实现留在 kda_pallas.py sub_chunk="fat"）**：
正确性与 sub-chunk 逐位相同，但性能 0.80-0.84×（全 hb×精度配置皆负）。预测错在假设窄 N
dot 浪费全部 256-lane padding；实测 Mosaic 对窄 dot 的处理好于最坏情形（microbench 44%
of ceiling 早有暗示），4× 真实 MAC 补不回来。教训：per-pass ceiling 推算 ≠ 实测，凡
"padded-MAC 账"必须上机验证。实现要点：fat-K 必须 clamp 行因子（inf 在 dot 内撞结构零
=NaN，where 救不了）+ float 块 mask（Mosaic 拒 3D bool）。附带新最优：sub-chunk+hb16+
s4bf16=10.46ms@H96 单卡（累计 4.5×；4 卡每卡 24 头仍用 hb8）。

**hb 全档 ablation（2026-08-11，H=96）**：MXU-vs-VPU speedup 随 hb 单调饱和：
1.15/1.29/1.63/1.76×@hb=1/2/4/8；hb=16 VPU 变体 VMEM OOM（[64,64,128]×16 头超限）而
GEMM 变体正常（12.33ms）——快路径每头省 ~16× VMEM，这本身解锁更高 hb。机理：固定步开销
O/hb 稀释（(V+O/hb)/(M+O/hb)→V/M）+ batching 对 MXU 不对称增益（VPU 变体 hb≥4 后平台 ~23ms）。
图 fig/kda_hb_ablation.png，已入 PR。

**命名定稿 + roofline（2026-08-11）**：strip 全量改为 **safe_gate path**（路径名，可 grep）
+ **sub-chunk GEMM**（技法名，对齐 FLA 的 sub_chunk）；lower_bound 不当路径名（A/B 第 2/3 行
同 lb 不同路径证明会混淆）。commit 3784449（rename，4 tests 重验过）。roofline 结论：
generic 在 VPU roof（3.85T 估计）跑到 ~80% 效率，safe_gate 在形状 ceiling（28.7T）跑 44%
（裸 918T 的 1.4%）——加速来自换 roof 而非提效率；44%→ceiling 是调度、ceiling→peak 23× 是
形状（fat-K）。图 prof/fig/kda_safe_gate_roofline.{pdf,png}（fig_roofline.py，dataviz 调色板
已验证），fork pr-assets 分支托管、嵌入 PR 正文。PR 现 5 commits / 2 files。

**PR 最终形态（2026-08-11 收官）**：3 个 commit——782ccd3 strip 路径、307c33a safe_gate 默认
False（对齐 FLA）、1f75355 where-mask 替代 clamp+乘法 mask、删除 safe_gate_cap 整条参数链、
FLA 式入口校验（safe_gate∧use_gate_in_kernel 要求 lb∈[-5,0) 否则 ValueError；预激活 g 信任调用方）。
9 个 CI check 等价项全绿（2 个环境失败已定性并修复：wan 资产下载 27G、part3/4 单芯重跑
TPU_VISIBLE_CHIPS=0）。评论 #5250042043（首轮汇总）+ #5250435959（全绿收官）。
设计演化由用户逐层质询驱动：cap 为何存在→FLA/FlashKDA 如何做→保留 safe_gate（预激活 g 场景）
+ where + 写死范围校验。FLA 机制：lower_bound 只选激活公式（值），safe_gate 只选 intra kernel
（bit），预激活模式下 flag 是请求快路径的唯一手段；FlashKDA 单模式无 flag。

**sglang-jax PR 已发（2026-08-11，仅发在自己 fork，官方 repo 零接触——用户明确指示）**：
https://github.com/Junyi-99/sglang-jax/pull/2（base=main head=kda-safe-gate-mxu，
cross-repo=false，官方无任何 PR/分支）。patch=782ccd3（本地 clone ~/claude/kimi-k3-kda/sglang-jax）：
kernel strip 分支 + safe_gate_cap 传参 + gating（safe_gate∧lower_bound≠None∧2cap+log2K<126）
+ 2 个 bounded 32k varlen 测试节点。TPU 实测：4 tests passed（654s）；patched chunk_kda_fwd
183.4→158.8ms=1.16×（比 standalone ablation 1.71× 低是因为 intra 只占其 pipeline 一部分），
correctness 2.4e-4（同 unpatched 精度包络）。正文含全部证据链。
方案：在他们的 `_kda_fwd_intra_kernel` 加 safe_gate 分支（sub-chunk-16 strip GEMM 替换
[64,64,128] decay tensor）。关键适配：他们是 log2 空间（lb=-5 → 步长 7.213，16-sub 界 57.7），
且 mask 用乘法而非 where——因子化路径 anti-causal 行会 inf，必须 clamp(+60) 防 inf×0=NaN。

sglang-jax 对照（commit 13fc48de4e）：其 safe_gate/lower_bound 只是 gate plumbing，
Pallas intra kernel 未用（[64,64,128] decay tensor + 逐行 solve，VPU-heavy）——
upstream 候选就是本目录的 strip-GEMM 方案。相关：[[kda-inv-experiment-pending]]
