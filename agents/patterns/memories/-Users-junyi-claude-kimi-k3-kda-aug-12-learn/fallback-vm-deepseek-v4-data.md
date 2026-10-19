---
name: fallback-vm-deepseek-v4-data
description: spot TPU (mamba-v6e-4) 挂掉时的临时备用机 SSH 与挂载点
metadata: 
  node_type: memory
  type: reference
  originSessionId: 05f28c9f-528c-45bd-bb11-d89d5b4dc81c
  modified: 2026-08-18T01:25:46.566Z
---

备用机（2026-08-18 用户提供，spot 机器挂掉时临时用）：

- SSH: `ssh junyi@34.187.66.199`
- 挂载点: `/mnt/disks/deepseek-v4-data`

主力机是 spot TPU `mamba-v6e-4`（us-east1-d，随时可能被抢占）。备用机是否带 TPU、规格未确认，用前先探（`jax.devices()`）。相关：[[kda-neumann-bounded-gate-only]]
