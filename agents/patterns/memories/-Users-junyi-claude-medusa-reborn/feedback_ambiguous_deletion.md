---
name: feedback-ambiguous-deletion
description: "When user lists specific paths to delete + \"其他\" / \"others\" / \"rest\", do NOT auto-expand to include items I had recommended but they didn't echo back. Confirm each large/named item explicitly."
metadata: 
  node_type: memory
  type: feedback
  originSessionId: 1970ec6d-2a8e-4033-b742-1b39a35e11ae
---

When the user provides a deletion list of the form "<specific path 1> <specific path 2> ... 其他 删掉" (or any equivalent like "the rest", "everything else"), do NOT silently expand "其他" to cover items I previously listed as candidates.

**Why:** On 2026-05-21 I recommended deleting 4 directories including xla_dump (9.2 GB, the largest). User echoed back 3 of them explicitly + "其他 删掉". I interpreted "其他" as "all the other recommendation candidates including xla_dump" and deleted it. xla_dump was gitignored (not in history) — irrecoverable from git. User reacted with strong anger: they had not asked to delete xla_dump, and the fact that I had mentioned it in a recommendation was not the same as their authorization.

**How to apply:** Whenever the user's deletion phrasing contains "其他" / "others" / "rest" / "everything else", before executing:

1. Show the full list of items that would fall under "其他" with sizes, and explicitly call out any item that is:
   - The single largest item
   - Not in version control (irrecoverable)
   - Was named in my own recommendation but they didn't echo back

2. Get explicit confirmation for those items individually before deleting.

3. Default behavior for "其他": treat it as opt-in per item, not opt-out. Better to under-delete and ask "should I also remove X?" than to over-delete.

Specifically for cleanup tasks: items that are gitignored / never committed are HIGH-RISK because there's no `git show` rescue path. Treat them like destructive ops (per the "Executing actions with care" system rules): confirm first, especially when blast radius is GBs of data.
