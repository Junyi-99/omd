---
name: feedback-save-scripts-locally
description: 用户要求：写过的脚本必须存到本地仓库，不要用一次性 heredoc、不要重写
metadata: 
  node_type: memory
  type: feedback
  originSessionId: e645184a-6feb-471a-ae36-e740b711b0e0
  modified: 2026-08-10T14:46:30.841Z
---

写过的脚本一律先存到本地项目目录（如 `kda-lower-bound/`），再 scp/执行；禁止一次性
heredoc 内联脚本（探测类小脚本也一样）。

**Why:** 一次性脚本无法复跑、无法审计；spot 机器随时重建，本地是唯一可靠存储；用户要
release 项目目录，脚本本身是交付物的一部分。
**另一条同源教训（2026-08-11）**：共享工作区里禁用 `git add -A`——用户放在树里的审阅
副本（kda-old.py/kda-v1.py）被扫进 commit 并 push，需要追加 revert commit 清理。
只 add 明确的目标文件。

**孤儿分支 + rm -rf 坑（2026-08-11）**：`git checkout --orphan` 后 `rm -rf ./*` 会删掉
工作树里未提交的改动（本例 kda.py 的 rename 险些丢失，靠 VM 副本找回）；切回原分支不会
自动恢复被删文件，需 `git restore .`。要做 assets 孤儿分支，应在临时目录 clone/worktree
里做，不要在有未提交改动的主工作区做。

**How to apply:** 任何 `python - <<EOF` / 远端临时文件的冲动 → 改为 Write 到项目目录
+ scp；探测型脚本可集中放一个 tools/ 或 probe 文件里反复扩充。相关：[[kda-lower-bound-oracle]]
