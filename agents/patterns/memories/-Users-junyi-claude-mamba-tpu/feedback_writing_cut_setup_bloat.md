---
name: writing-cut-setup-bloat
description: "When polishing the user's writing, aggressively merge sections that say the same thing. Don't keep \"What we compare\" and \"Experimental setup\" as separate sections — they are one section about the experiment, and the prose between Table 1 and the workload command is usually filler."
metadata: 
  node_type: memory
  type: feedback
  originSessionId: c34c9208-a02e-48c9-a29f-e4741c64e8e5
---

When polishing the user's technical writing, look for adjacent sections that share a topic and merge them. In the Mamba2 SSD on TPU blog post, the first rewrite kept "What We Compare (and What We Do Not)" and "Experiment Settings" as separate top-level sections, with subsection breakouts inside each. The user's verdict: "What We Compare 和 Experiment Settings 说的是同一回事" — they are saying the same thing — and the prose between Table 1 and the workload command was "废话太多" (too much filler).

The fix that worked:
- Merge to one section ("Setup").
- Put the framing table (Table 1) right after a one-sentence intro.
- One short paragraph for hardware + shape constants.
- One short paragraph that folds workload generator command + uniform-vs-jagged + mamba2-jax pad-adapter + GDN-port harness into a single tight block.
- One short paragraph + Table 2 for version roles.
- One trailing pointer line: "Methodology notes — X, Y, Z — are in the appendix."
- Move kernel-time-vs-wall-clock methodology, raw-Bonsai-vs-reproduced clarification, and `reference.py`-is-not-a-baseline notes ALL into the appendix.

Net result: the merged §3+§4 was ~50% the original length, no information lost.

**Why:** This is the author's signature, a technical blog, not a conference paper — every "we will explain X below" sentence, every "Note that..." clarification, every redundant Figure-N citation is a tax the reader pays. Cut them.

**How to apply (general):**
- When dispatching section-drafting agents, give them adjacent-section merging power, not just "rewrite this section in place".
- Look for telltale duplication: same Figure cited in adjacent paragraphs, same fact stated twice with different framings, "Three slices answer three different questions" prose that just lists the upcoming subsection titles.
- Methodology footnotes belong in the appendix unless they actually load-bear in the main flow. Default: appendix.
- After the agent returns, do a redundancy sweep — read every paragraph against its neighbors and ask "does this say something the neighbor doesn't already say." If not, cut it.
