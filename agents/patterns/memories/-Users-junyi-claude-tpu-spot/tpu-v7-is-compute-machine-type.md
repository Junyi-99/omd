---
name: tpu-v7-is-compute-machine-type
description: TPU v7 (tpu7x) 只作为 Compute Engine machine type 存在，不在 gcloud compute tpus accelerator-types 里
metadata:
  type: reference
---

TPU v7 = `tpu7x-standard-1t/4t`、`tpu7x-ultranet-4t`，是 **Compute Engine machine type**，
只能用 `gcloud compute machine-types list --filter="name~tpu7x"` 看到；老的
`gcloud compute tpus accelerator-types list`（Cloud TPU API）永远看不到 v7，
只有 v3/v5p/v6e。2026-05-11 那次"看不到 v7、以为是权限问题"就是查错了 API surface。

2026-09-10 实测（medusa-compute 和 project-68023df7 两个项目都一样）：
tpu7x 出现在 us-central1-b、us-central1-c、us-east1-b、europe-west2-a；
quota 只有 `compute.googleapis.com/preemptible_tpu7x`，没有 on-demand metric。
**放量只在 `project-68023df7-96eb-4779-857`（"Medusa"，556758793357）：regional
bucket us-central1 = 4 芯。** `medusa-compute`（仓库默认项目）没有放量，拿不到 v7。
2026-09-10 在 medusa-compute 排 1 芯 flex-start 等了 62 分钟，最后
ZONE_RESOURCE_POOL_EXHAUSTED 失败。

排队走 GKE node pool：`--machine-type=tpu7x-standard-4t --tpu-topology=...`
配 `--flex-start`（DWS 排队）或 `--enable-queued-provisioning` + Kueue ProvisioningRequest。
相关：[[save-raw-experiment-logs]]
