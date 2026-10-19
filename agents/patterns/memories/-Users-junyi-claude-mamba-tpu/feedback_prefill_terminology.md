---
name: prefill-terminology-precision
description: "When the user's writing is scoped to prefill kernel work, do not use words that conflate prefill and decode — \"serving\" especially. Use \"prefill\" / \"jagged-prefill\" / \"packed-jagged input\" / \"production prefill\" instead."
metadata: 
  node_type: memory
  type: feedback
  originSessionId: c34c9208-a02e-48c9-a29f-e4741c64e8e5
---

When the doc is about prefill kernel optimization, the word **"serving"** is wrong. The user's exact framing: *"我很讨厌你说 'Serving' 这个词，因为我们在做 Kernel 的 Optimization，这个文档针对 Prefill 这个阶段。你说 Serving 这个词会包括 both Prefill and Decode, which is not Precise."* Serving in LLM/inference context implies both prefill and decode; using it leaks decode connotations into a prefill-only doc.

Phrases to avoid in prefill-scope writing, and what to use instead:
- "serving shape" / "serving layout" / "our packed serving layout" → "packed-jagged input" (per the user's earlier explicit rename)
- "Jagged serving" / "Jagged B=64 serving" → "Jagged prefill" / "Jagged B=64 prefill"
- "real-world serving" → "real-world prefill" or "real-world jagged prefill"
- "serving headline" → "jagged-prefill headline"
- "serving pattern" / "serving result" / "serving APIs" → "prefill workload" / "prefill numbers" / "packed-jagged input APIs"
- "serving path" (if it pops back up — was the v12 persona we already deleted) — do not reintroduce

Acceptable, NOT to rewrite:
- "deployment" — generic word for putting into production, doesn't conflate prefill/decode.
- "production" — same.
- "inference" — technically includes both, but is less loaded than "serving"; check context before changing.
- "real-world" as a modifier of "jagged prefill" / "workload distribution" — fine.

**Protected-opening rule interaction:** The earlier [[dont-touch-user-opening]] rule says don't change the first two sections. This new terminology-precision rule **overrides** that rule for surgical word swaps (not voice/framing rewrites). If a "serving" mention sits in the protected opening, evaluate whether it's a methodological claim (fix it) or an external citation (leave it; flag to user).

In the Mamba2 SSD on TPU blog (Google Doc 1PlqV5a...), §1's "Anthropic is expanding Claude training and serving on Google Cloud TPUs" was left alone because it's quoting Anthropic's industry posture; §2's "our target serving pattern" → "our target jagged-prefill workload" was changed because it's our claim.
