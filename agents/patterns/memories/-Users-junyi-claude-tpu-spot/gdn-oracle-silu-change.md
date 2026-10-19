---
name: gdn-oracle-silu-change
description: "GDN prefill silus q/k/v internally; the 0.20.0 oracle needs pre-silu'd input"
metadata: 
  node_type: memory
  type: project
  originSessionId: 5d50b5b8-1c61-438f-9787-d3a9dd744758
---

GDN prefill 的正确语义**包含 silu**(post-conv 对 q/k/v 做 silu)。具体约定(2026-06-05 读 0.20.0 源确认):
- **prefill impl**(`kernels/gdn/recurrent_scan_v2.recurrent_scan`,v3 要替换的 baseline)收 **RAW** post-conv mixed_qkv,**内部 silu**(recurrent_scan_v2.py 行 207/413/591)。wrapper 的 'jax'/'fused' 分支则在调用前 silu。
- **oracle** `ragged_gated_delta_rule_ref`(0.20.0)**内部不 silu** → 校验时必须对它喂 `silu(mixed_qkv)`。

**How to apply:** v3 必须**内部 silu**(同 v2),否则当 prefill impl 接 RAW 输入会差一个 silu(diff~0.06–0.12)。测试/对 oracle 比较时给 oracle 喂 silu'd 输入。

**Why / 纠错:** 本会话早期一度「去掉 v3 的 silu」来匹配「对 oracle 喂 RAW」的写法——那是**错的**(把 oracle 调用约定也搞错了)。real pipeline 是 silu 的。已在 commit 346b950 恢复 v3 内部 silu 并修正测试。详见 gdn-kernel-opt/v3/integration_test.py、state_test.py。

旧 `kernels/gdn/v2/` 路径在 0.20.0 移到 `kernels/gdn/recurrent_scan_v2`(仍在);prefill baseline 用它,decode 用 `fused_gdn`。集成 seam:`ragged_gated_delta_rule_wrapper`(prefill_impl/decode_impl 可分别配,is_decode_only 路由),monkeypatch 其 `recurrent_scan` 符号即可换成 v3。
