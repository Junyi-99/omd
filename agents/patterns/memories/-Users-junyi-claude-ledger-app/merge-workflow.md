---
name: merge-workflow
description: How the user wants commits landed in the ledger-app repo
metadata: 
  node_type: memory
  type: feedback
  originSessionId: b9a4b8b0-65f2-4603-b06f-0d90bb4f4356
---

For the ledger-app repo, the user wants changes merged **directly into `main`, no PR**.

**Why:** It's their personal repo (Junyi-99/ledger-app); they explicitly said "没必要开 pr。直接合入 main" and pushed back on PR ceremony.

**How to apply:** When asked to commit & push, land it on `main` (a quick branch → `--ff-only` merge → push → delete-branch flow is fine, but don't open a PR or ask). Still use the required commit trailers.
