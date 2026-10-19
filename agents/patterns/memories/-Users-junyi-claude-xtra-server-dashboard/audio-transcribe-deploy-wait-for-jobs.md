---
name: audio-transcribe-deploy-wait-for-jobs
description: "Before restarting/deploying audio.xtra.science, wait until no transcription job is running at all"
metadata:
  node_type: memory
  type: feedback
  originSessionId: c9cb212f-e317-436c-85fc-fa399b9ed2a2
  modified: 2026-10-08T11:27:41.145Z
---

Before any deploy/restart of the audio-transcribe app (audio.xtra.science, xtras3, gitops repo xtra-gitops), wait until the `jobs` table in `/data/transcribe.db` has zero jobs outside `done`/`error` — not just zero jobs in the upload stage.

**Why:** The user (2026-10-08) was mid-transcription when I deployed v0.7.0 because I judged their already-submitted job "resumable". They had explicitly asked not to drop their job. Resume works in local tests but a restart still risks it, and users notice.

**How to apply:** Query the live DB via `kubectl exec`; if anything is active, prepare the gitops change locally and run a terminating wait loop (background Bash, 30s checks) before pushing. Re-check right before `git push`.
