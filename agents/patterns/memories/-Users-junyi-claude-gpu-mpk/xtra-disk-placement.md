---
name: xtra-disk-placement
description: "On xtra* GPU servers, put large data under /disk*/junyi or /ssd*/junyi when /home or / is low"
metadata:
  node_type: memory
  type: feedback
  originSessionId: a349f4dd-72fd-4b8d-af26-94860a64042d
  modified: 2026-10-14T20:32:39.550Z
---

On the NUS xtra GPU servers (xtrah100, xtrah200, xtraa100 ***REMOVED***), when /home or / is short on space, check /disk* and /ssd* mounts, create /disk*/junyi or /ssd*/junyi, and put envs/models/data there.

**Why:** xtrah200 root was 99% full (102G free) on 2026-10-15; user pointed to these extra disks.
**How to apply:** Before installing envs or downloading weights on these hosts, run `df -h /disk* /ssd*` and pick the emptiest; tell subagents to do the same.
