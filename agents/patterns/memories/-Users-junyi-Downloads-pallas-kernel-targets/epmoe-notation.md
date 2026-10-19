---
name: epmoe-notation
description: Fixed symbol set the user requires for all EPMoE / MoE kernel discussion
metadata:
  type: feedback
---

When discussing EPMoE (sglang-jax MoE), always use this symbol set:
`T` = token count, `E` = expert count, `K` = experts selected per token (top_k),
`H` = model hidden size, `F` = expert FFN intermediate size.

**Why:** the user set this on 2026-09-03 after I had drifted between M/D/N/L across
many messages, which made cross-message tables hard to compare.

**How to apply:** map old usage — my `M` -> `T`, my `D` (unpermute row width) -> `H`,
my `N = M*K` -> `T*K` rows. Never reuse `T` for anything else; the Pallas
tokens-per-program tuning parameter is now `BT`. Note `F` does not appear in
`_unpermute` at all (it lives in the GMM), and for latent-MoE models like Kimi K3
the unpermute width is `routed_expert_hidden_size`, not `H`.

Name the four unpermute implementations in full, never by letter:
**baseline** (sglang-jax main), **#1058** (that PR's idea rebased), **unroll**
(the ~10-line JAX change), **pallas** (the fused Pallas-TPU kernel).

**Why:** the user pushed back on 2026-09-03 — "A C P 都是啥？不要这样。是 baseline
就直接写 baseline". Single letters and aliases (A_main/C_unroll/P3d/fuse/fused3d)
made tables unreadable across messages.

**How to apply:** use these names in prose, tables, and script output alike;
`scratchpad/*.py` already emit them.

Benchmark a candidate kernel against baseline's own function signature and input
shapes. Any reshape/relayout the candidate needs is part of its cost and must be
inside the timed function.

**Why:** the user corrected this on 2026-09-03. I had been passing the Pallas
kernel a pre-reshaped 3-D source and quoting 1.9-2.5x; with the relayout counted
the real figure is 1.2-1.53x, and one shape flips from win to loss. Never assume
the producer (GMM) emits a convenient layout.

**How to apply:** `scratchpad/unpermute_pallas.py:unpermute()` now takes
`bf16[T*K, H]` and reshapes internally. If a layout change at the producer would
help, measure and propose it as a separate change — do not fold it into the
kernel's numbers as an assumption.
