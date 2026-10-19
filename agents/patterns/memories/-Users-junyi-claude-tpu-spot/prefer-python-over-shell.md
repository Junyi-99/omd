---
name: prefer-python-over-shell
description: 用户要求用 Python 写多步/循环逻辑，不要写 shell（zsh）脚本
metadata:
  type: feedback
---

写重试循环、批量探测这类多步逻辑时，用 Python 脚本，不要在 Bash 工具里写 zsh 循环。

**Why:** 用户明确说过"你别 jb 写 zsh shell 了，写 python"。zsh 还有默认不做变量分词等坑，循环容易静默跑错。

**How to apply:** 需要循环/重试/解析输出时，先把脚本写到 scratchpad 的 .py 文件再 `python3` 跑。单条只读命令（cat/grep/gcloud describe）仍直接用 Bash。
