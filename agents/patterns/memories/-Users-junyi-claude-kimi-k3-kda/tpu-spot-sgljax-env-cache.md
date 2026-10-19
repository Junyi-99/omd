---
name: tpu-spot-sgljax-env-cache
description: GCS 里 sglang-jax 的环境缓存是 sgljax-py3.12-jax0.10.2.tar.gz（解到 ~/venv-sgl），文件名不符合 env_cache.py 的 <key>.tar 约定；model-cache 里没有 gemma；unsloth/gemma-2-2b-it 是非 gated 镜像
metadata: 
  node_type: memory
  type: project
  originSessionId: 18ea2898-6b8c-4aa4-a18d-e9fa279504bd
  modified: 2026-09-07T03:30:16.174Z
---

- `gs://medusa-experiments/env-cache/sgljax-py3.12-jax0.10.2.tar.gz`（617MB，2026-08-17）：tar 顶层是 `venv-sgl/`（py3.12，jax 0.10.2）。
  它是 `.tar.gz`，而 `infra/vm/env_cache.py` 的 `has/restore` 只认 `<key>.tar`，所以要手动 `gcloud storage cp` + `tar xzf -C $HOME`。
  sglang-jax 仓库现在 pin `jax[tpu]==0.11.1`（`python/pyproject.toml`），还原后需 `uv pip install -e "python[tpu]"` 升级，再另存新 key（如 `sgljax-py3.12-jax0.11.1`）。
- `gs://medusa-experiments/model-cache/` 2026-09-07 只有 `kimi-linear-48b/`；没有任何 hybrid SWA 模型。
- 本机没有 HF token；`google/gemma-2-2b-it` gated，`unsloth/gemma-2-2b-it` 非 gated 且含 safetensors，可作为 sglang-jax gemma2 路径的替代权重（Gemma-2 sliding_window=4096，要触发 SWA 逐出序列必须 >4096，`--context-length 8192`）。
- sglang-jax 支持 sliding window 的模型只有 gemma2 / gemma4 / mimo_v2_flash（`grep sliding_window python/sgl_jax/srt/models/`）。
- 2026-09-07 上午 us-east1-d 的 spot v6e-4 连续 3 次 "no more capacity"；europe-west4-a 连续 7 次 "Failed to perform tenant project creation"（preflight 配额 64 chips、v6e-4 offered 都 PASS，约 10 秒即失败，疑似该区域项目侧问题，非容量）。
  退路：单芯 `TPU_SPOT_NAME=pr1564-v6e-1 TPU_SPOT_ACCEL=v6e-1` 在 us-east1-d；gemma-2-2b tp=1 够用。`tpu_vm.py up` 要直接 `2>&1` 落盘，别接 `| tail`，否则中间报错看不到。
- 作者跑 A/B 用的是仓库自带 `benchmark/hicache/bench_unified_radix_ab.py`（configs no-cache/radix/unified，workloads random/gsp），但它不透传 `--context-length/--swa-full-tokens-ratio`。

**Why:** 每次上 TPU 都先查缓存再装（CLAUDE.md 规则），这几个命名/gated 的坑每次都会浪费时间。
**How to apply:** 用 `exp-2026-09-07-pr1564-review/scripts/tpu_setup_pr1564.sh`（restore-first → 安装 → 回写缓存）和 `tpu_e2e_swa_ab.sh`（模型 restore-first → legacy/unified 各起一次 server）。相关：[[tpu-spot-ssh-host-block]] [[sglang-jax-pr1564-review]]
