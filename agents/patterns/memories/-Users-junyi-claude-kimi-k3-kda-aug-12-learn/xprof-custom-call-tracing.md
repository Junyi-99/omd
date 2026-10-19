---
name: xprof-custom-call-tracing
description: XProf 看 pallas/custom-call 内部 VPU/MXU 活动的开关（LIBTPU_INIT_ARGS）
metadata: 
  node_type: memory
  type: reference
  originSessionId: 05f28c9f-528c-45bd-bb11-d89d5b4dc81c
  modified: 2026-08-19T07:22:08.207Z
---

XProf 默认看不到 pallas kernel（custom call）内部的 VPU/MXU 单元级活动，但可以开插桩：

```bash
LIBTPU_INIT_ARGS="--xla_xprof_enable_custom_call_tracing=true --xla_xprof_register_llo_debug_info=true" \
  python your_jax_workload.py
```

- 会触发带 instrumentation 的重编译，trace 中出现 LLO 级 debug 信息（VPU/MXU 归属）
- 插桩本身有开销：**归因用插桩 trace，计时数字用未插桩运行**，两者分开
- 2026-08-19 由用户提供；此前我错误断言"xprof 无法暴露 pallas 内部单元分类"，
  相关报告表述（tpu-kda report §3/§12 的方法学声明）需要修订

**版本要求（2026-08-19 实测）**：这两个新 flag 名需要 **libtpu ≥ 0.0.46**（jax 0.11.1）。
在 libtpu 0.0.42.1（jax 0.10.2）上不存在，传入直接 "Unknown command line flag" 崩溃；
那一版只能用 legacy 名 `xla_enable_custom_call_region_trace` / `xla_enable_mxu_trace` /
`xla_enable_local_dma_trace`（strings libtpu.so 探测得到，0.0.46 里两套都在）。

新 flag 生效后（KDA S=1024×B=8，30 iter，v6e）：trace 13MB→753MB、9→17 泳道，
新增 MXU0/MXU1/VALU/VLD/VST/XLU/EUP/SALU Instructions + Pallas Primitives + _counters_，
且 Source code 泳道按 pallas 源码行归因（如 kda/v2/head_block.py:109）。
镜像 [[kda-tpu-bench]]（ghcr.io/junyi-99/kda-tpu-bench:latest）已钉到这套新栈。
