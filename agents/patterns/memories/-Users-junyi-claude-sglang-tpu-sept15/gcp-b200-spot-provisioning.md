---
name: gcp-b200-spot-provisioning
description: How to rent 8xB200 spot GPUs in the "Medusa Weihao" GCP project — quota shape, working image, zone order, stop/delete gotchas
metadata:
  type: project
---

GCP 项目 **Medusa Weihao**：`project-68023df7-96eb-4779-857`（number 556758793357），org 615654616044。

截至 2026-09-20 实测：

- **B300 在 GCP 尚未开放**。该项目 accelerator-types 只有 `nvidia-b200` / `nvidia-gb200`（外加 H100/H200/A100/L4/RTX-PRO-6000 和各代 TPU）。
- **B200 只有 spot 配额**：Cloud Quotas 里只存在 `PREEMPTIBLE-NVIDIA-B200-GPUS-per-project-region`（= 64，usage 0）和 `COMMITTED-...`，**没有按需 on-demand 配额**。想要非 spot 得先申请配额或买预留。
- B200 配额**不在** `gcloud compute regions describe` 的老 quota 列表里，要用 `gcloud alpha quotas info list --service=compute.googleapis.com`。
- 机型 `a4-highgpu-8g` = 8×B200 + 224 vCPU + 3968 GB RAM，GPU 隐含在机型里，不用另加 `--accelerator`。
- 镜像用 `ubuntu-accelerator-2404-amd64-with-nvidia-580`（project `ubuntu-os-accelerator-images`），开机即有驱动，nvidia-smi 直接认到 8 卡。
- 启动盘必须 `hyperdisk-balanced`（A4 不支持 pd-*）。
- **spot 容量是唯一瓶颈**，按区顺序试：2026-09-20 当天 us-central1-b ✗、us-east1-d ✗、**us-east1-b ✓**。写个 for 循环轮询建机，第一个成功就停。
- **停机要加 `--discard-local-ssd=true`**：A4 自带 Local SSD，不给这个 flag 会 HTTP 400 直接失败。
- spot 实例 stop 之后会被回收消失（连启动盘一起），不用再手动 delete。

可用命令见 [[gcp-b200-create-command]]。TPU 侧的环境/权重缓存约定见 CLAUDE.md（env-cache / model-cache），与这里的 GPU 机器无关。
