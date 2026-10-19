---
name: gcp-b200-create-command
description: Working gcloud command to create an 8xB200 spot VM in the Medusa Weihao project
metadata:
  type: reference
---

2026-09-20 实测可用的建机命令（背景见 [[gcp-b200-spot-provisioning]]）：

```bash
gcloud compute instances create b200-8g-spot \
  --project=project-68023df7-96eb-4779-857 --zone=us-east1-b \
  --machine-type=a4-highgpu-8g \
  --provisioning-model=SPOT --instance-termination-action=STOP \
  --maintenance-policy=TERMINATE \
  --image-family=ubuntu-accelerator-2404-amd64-with-nvidia-580 \
  --image-project=ubuntu-os-accelerator-images \
  --boot-disk-size=1000GB --boot-disk-type=hyperdisk-balanced \
  --network-interface=nic-type=GVNIC,network=default,subnet=default \
  --scopes=cloud-platform
```

没容量报 `The zone ... does not have enough resources`，换下一个区重试。
区序：us-east1-b → us-east4-b → us-west2-c → us-south1-b → europe-west4-b → asia-northeast1-b。

停机：

```bash
gcloud compute instances stop b200-8g-spot --zone=us-east1-b \
  --project=project-68023df7-96eb-4779-857 --discard-local-ssd=true
```
