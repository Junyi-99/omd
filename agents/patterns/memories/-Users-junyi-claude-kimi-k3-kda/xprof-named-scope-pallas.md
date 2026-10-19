---
name: xprof-named-scope-pallas
description: 用 jax.named_scope 给 Pallas TPU kernel 打 xprof 子步骤标注的可用配方 + 三个会静默坑人的限制
metadata: 
  node_type: memory
  type: reference
  originSessionId: f75e7da8-0e4c-4ca4-ac32-3772a62abd60
  modified: 2026-08-31T11:33:41.337Z
---

`jax.named_scope` 写在 Pallas TPU kernel body 里会被 lower 成 `tpu.trace_start` /
`tpu.trace_stop`（`jax/_src/pallas/mosaic/lowering.py`，注释 "Insert trace_start and
trace_stop ops on named_scope boundaries"）。开 `--xla_xprof_enable_custom_call_tracing=true`
后，这些 region 以 **`XLA TraceMe`** lane 的 X 事件出现，每个 scope 每个 grid step 一条。
文档：<https://openxla.org/xprof/custom_call_profiling>；raw md 在
`https://raw.githubusercontent.com/openxla/xprof/master/docs/custom_call_profiling.md`
（网页版是 JS 渲染的，curl 抓不到正文）。

三个会静默坑人的限制（2026-08-31 在 jax 0.11.1 / libtpu 0.0.47 / xprof-nightly 2.24.1a20260830 上实测）：

1. **一个进程只有第一次 `jax.profiler.trace` 出数据**。开了 custom-call tracing 之后，同进程
   第二次及以后的 capture 写出 400 字节的空 trace，不报错。扫参数时必须 **一个 shape 一个进程**。
   不开这个 flag 的话同进程连续 capture 正常。
2. **开了 tracing 后 custom call 自己的 `XLA Ops` 事件会消失**（被 trace region 取代）。
   算 kernel 时间要用 `XLA Modules` 的时长减掉非 kernel 的 XLA Ops，不能直接找 kernel op。
3. **trace 体积爆炸**：这个 flag 顺带打开 instruction bundle instrumentation，4k token 的
   一次 capture ≈ 70 MB `trace.json.gz` + 230 MB `xplane.pb`。只留
   `XLA Ops / XLA TraceMe / XLA Modules / Steps` 四条 lane 并删掉 xplane 后降到 ~100 KB
   （`exp-2026-08-31-kda-xprof-scope/scripts/trace_tools.py`）。Perfetto 只需要 json.gz。

另外：`LIBTPU_INIT_ARGS` 必须在 `import jax` 之前设好（容器 env 或 shell export，不能在 py 里 import 后设）；
trace region 本身有开销（紧的 kernel 受影响更大：#1554 慢 47%，#1550 只慢 12%），所以
**计时和标注要分成两个 build**——把 named_scope 包在 `KDA_XPROF_SCOPES` 环境变量后面即可。

**Why:** 这三条都不会报错，只会给出空数据或错误归因，第一次踩要花一两个小时。
**How to apply:** 复用 `exp-2026-08-31-kda-xprof-scope/scripts/{bench_kda,trace_tools,run_phase}.py|sh`。
相关：[[kda-1550-vs-1554-xprof-substep]]。
