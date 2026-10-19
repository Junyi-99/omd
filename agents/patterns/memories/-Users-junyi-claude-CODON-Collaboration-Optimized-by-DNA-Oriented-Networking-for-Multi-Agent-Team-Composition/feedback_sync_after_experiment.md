---
name: Sync results from remote after experiments finish
description: Whenever an experiment run on xtrah200 completes, promptly rsync ~/codon/results/ and ~/codon/logs/ back to the local repo
type: feedback
originSessionId: 2e1dfa5c-4d8e-49f0-8f70-ddc536395689
---
用户指令：跑完实验，及时同步。

**Why:** 原文 "你记住，跑完实验，及时同步"。实验数据（JSONL / logs）只存在 xtrah200 上，本地仓库没有副本，分析和写作都在本地做。不同步就可能漏追最新结果。

**How to apply:** 每当 xtrah200 上启动的实验（phase1_self_refine / phase2_* 等）完成时，立即执行：

```bash
rsync -az xtrah200:codon/results/ /Users/junyi/claude/CODON-Collaboration-Optimized-by-DNA-Oriented-Networking-for-Multi-Agent-Team-Composition/results/
rsync -az xtrah200:codon/logs/    /Users/junyi/claude/CODON-Collaboration-Optimized-by-DNA-Oriented-Networking-for-Multi-Agent-Team-Composition/logs/
```

不用等用户要，实验刚结束就主动拉。长跑任务中途收到进度查询时不必 sync（只在完成/失败时做）。
