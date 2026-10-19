---
name: tpu-spot-ssh-host-block
description: tpu-spot 曾因 ~/.ssh/config 缺 alias 的 Host 块让 controller 永远 SSH 不通；2026-08-28 已修：tpu_vm.py up/_refresh_ip 自动补块 + ssh-setup 子命令
metadata: 
  node_type: memory
  type: project
  originSessionId: f57bdb37-d0ec-4a7d-b795-5c90684723c3
  modified: 2026-08-28T01:43:10.668Z
---

`infra/tpu_vm.py::_refresh_ip` 只改写已存在的 `Host <alias>` 块里的 HostName，没有块就
`die("no Host block aliasing '<alias>' in ~/.ssh/config (run ssh setup first)")`；experiment controller 的
SSH probe 走 `ssh <alias>`，于是 VM READY 后永远探不通，几分钟后被抢占/删除，日志表现像"spot 不稳定"。
2026-08-27/28 三个 kda-1554-* 实验因此卡了一整天；补块后 7 分钟跑完。`tpu_experiment.py check --online` 不检查这一项。

模板（= `infra/fleet.sh::append_host_block`）：
```
Host <alias>
    HostName 0.0.0.0
    User junyi
    IdentityFile ~/.ssh/google_compute_engine
    IdentitiesOnly yes
    CheckHostIP no
    StrictHostKeyChecking no
    UserKnownHostsFile ~/.ssh/google_compute_known_hosts
    HashKnownHosts no
```

2026-08-28 修复（未 commit，在用户 WIP 工作树里）：`tpu_vm.py` 新增 `ensure_ssh_host_block()`（模板搬自 fleet.sh），`_refresh_ip` 先补块再改 IP，
新增 `ssh-setup` 子命令；`fleet.sh ensure_host_blocks` 改为串行调 `tpu_vm.py ssh-setup`；测试 `infra/tests/test_tpu_vm.py`（4 个）。

**Why:** 起实验前漏这一步的代价是整天拿不到结果，而且日志只在 READY 那一行后面才出现这条 ERROR，容易被容量错误淹没。
**How to apply:** 修复后 `up` 自动处理；老的 controller 子进程仍跑旧代码时，手动 `TPU_SPOT_NAME=<alias> TPU_SPOT_ZONE=<zone> ./infra/tpu_vm.py ssh-setup` 即可；
排查 REACQUIRING 反复时先 grep 日志 "no Host block"。相关：[[kda-1554-boundary-tpu-run]]。
