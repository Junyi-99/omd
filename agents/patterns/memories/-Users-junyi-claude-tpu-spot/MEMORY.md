# Memory index

- [GDN oracle silu change](gdn-oracle-silu-change.md) — GDN prefill silus q/k/v internally; 0.20.0 oracle needs pre-silu'd input; v3 must silu internally
- [Nemotron TPU Mamba2](nemotron-tpu-mamba2.md) — running Nemotron-3-Nano (Mamba2+ReLU²-MoE) on v6e via vLLM; runs end-to-end, output-correctness bug open
- [Save raw experiment logs](save-raw-experiment-logs.md) — preserve raw bench/eval logs to GCS (spot VMs are ephemeral), not just summaries
- [TPU v7 是 Compute machine type](tpu-v7-is-compute-machine-type.md) — tpu7x 只在 compute machine-types 下可见，tpus accelerator-types 永远查不到
- [Prefer Python over shell](prefer-python-over-shell.md) — 多步/循环逻辑写 Python 脚本，不要写 zsh 循环
- [ICI v6e measured constants](ici-v6e-measured-constants.md) — ppermute α≈2.95µs (Pallas RDMA 2.08µs), +1.19µs/hop, β=93.3GB/s, host floor 185µs, 8ns/descriptor, T≥1 no gain, overlap eff 1.0; exps ici-v6e4-2x2-r1/r2
- [ICI external interconnect refs](ici-external-interconnect-refs.md) — survey.md + data + charts in ici-bench/external; closest match arXiv 2610.01380; no exact setting match
- [Lab GPU machines](lab-gpu-machines.md) — xtra3090/a100/h100/h200/a6k02/a6k03: GPUs, NVLink, 100G RoCE ports, busy status
- [ICI experiment 3 small](ici-experiment-3-small.md) — E6 small single-transfer shards; on-demand box busy (deciax eval); spot runs in Medusa project
