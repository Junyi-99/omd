---
name: kda-inv-experiment-pending
description: KDA Neumann 逆 bf16 + C=64 range 问题——已被 sub_chunk=16 解决（见 kda-lower-bound-oracle），本文留历史数据
metadata: 
  node_type: memory
  type: project
  originSessionId: 911e46b9-0f78-4f3d-8862-5652e40b7599
  modified: 2026-08-10T13:31:08.821Z
---

实验完成（2026-08-03，mamba-v6e-4，T=8192 H=96 生产形状，20 次中位数）：

| C | inv=highest(6pass) | inv=bf16(1pass) | err_ratio(bf16 inv) |
|---|---|---|---|
| 32 (HB16) | 17.61 ms (基线) | 14.62 ms **1.20×** | 9.2e-06（lb=-5 生产 gate）|
| 64 (HB16) | 15.05 ms 1.17× | 10.90 ms **1.62×** | 9.5e-05（lb=-2）|
| 128 (HB8) | 15.21 ms 1.16× | 10.19 ms **1.73×** | 3.6e-04（lb=-1）|

- fp32 基线 err ~1.3e-07；bf16 逆误差随 C 增长但都远低于 bf16 I/O 舍入（4e-03）
- C=64/128 在生产 gate lb=-5 下 NaN（实测确认），C=64 需两级 rescale 修复才能部署
- C=128 需 HB=8（VMEM 33.3M>32M @HB16）且要求 lb≥-1.375，不可部署，仅作数据点
- Mosaic 不支持 Precision.HIGH（3-pass），只有 default/highest 两档
- 立即可部署：C=32 + inv_precision="default" → 1.20×（kernel 补丁在 tpu/kda_pallas.py）

**2026-08-10 更新：C=64 range 问题已解决。** `sub_chunk=16` 多参考点因子化
（`_chunk_body(sub=16)`，tpu/kda_pallas.py + kda-lower-bound/kda_pallas.py）在
lb=-5 生产 gate 下 C=64 全正确（含 g≡-5 worst case）。v6e T=8192 H=24：
c64-s16-hb8-invbf16 = 3.56 ms，3.35× vs 旧 pallas-c32，err ~1e-5。
详见 [[kda-lower-bound-oracle]]。
相关：[[kda-hardware-analysis-conclusions]]
