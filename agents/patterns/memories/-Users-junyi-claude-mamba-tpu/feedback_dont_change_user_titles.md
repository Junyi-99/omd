---
name: dont-touch-user-opening
description: "When polishing the user's writing, do not change the title, subtitle, or opening narrative paragraphs without explicit permission, even on a \"make it more professional\" pass."
metadata: 
  node_type: memory
  type: feedback
  originSessionId: c34c9208-a02e-48c9-a29f-e4741c64e8e5
---

When polishing the user's writing (blog posts, papers, docs), the title and the opening narrative are off-limits during a general "make it more professional" pass. The body — tables, results, methodology, technical sections, conclusions — is fair game; the framing the user chose for the entrance to their piece is not.

**Why:** Same session, two rejections on the Mamba2 SSD on TPU blog post (Google Doc 1PlqV5a...):
1. A subagent proposed changing the title from "Model Onboarding on TPU with Medusa: A Mamba2 SSD Case Study" to "Mamba2 SSD on TPU: A 4.60× Onboarding Harness". User: "不要改我的 Title".
2. Right after that: "开头的前两个段落都不要改（Why A Dedicated Harness for TPU 和 Mamba V2 SSD as a Concrete Onboarding Case Study）" — don't touch the first two opening sections either.

The pattern: the author's voice lives in the title + opening. They want to own that framing. Polish improves their argument, not their voice.

**How to apply:**
- Title: leave verbatim.
- Subtitle / italic strap below the title: leave verbatim unless the user says otherwise.
- The first 1–2 narrative sections that set up the post (motivation, why this case study): leave verbatim.
- Body sections from methodology / setup / results / closing onward: fair game for restructuring, tightening, etc.
- Do NOT inject new top-matter (TL;DR boxes, abstract callouts, author lines) into the opening zone without asking first.
- If the opening seems genuinely weak and you want to propose a change, ask via [[AskUserQuestion]] BEFORE any subagent runs on it.
