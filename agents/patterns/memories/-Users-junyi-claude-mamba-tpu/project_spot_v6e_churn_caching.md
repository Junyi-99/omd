---
name: project_spot_v6e_churn_caching
description: "Operating the medusa-compute spot v6e-4 box — zone quota, GCS cache for fast recovery, tmux launch, mid-recover-kill hazard"
metadata: 
  node_type: memory
  type: project
  originSessionId: 2c9685e2-d246-42ef-a65f-42c042fb5631
---

Running E2E work on the `mamba-v6e-4` **spot** TPU (project `medusa-compute`, hard
constraint: never switch project). Operational facts learned 2026-06-01/02:

- **Spot quota is only in TWO zones: `us-east1-d` and `europe-west4-a`** (64 v6e
  chips each). When one is capacity-exhausted (`create` → "no more capacity in the
  zone"), fall back to the other via `TPU_SPOT_ZONE=europe-west4-a ./tpu-spot.sh recover`.
  Both `tpu-spot.sh` and `tpu-spot/gdn-v3-fused-prefill/auto_retry_v3.sh` read `TPU_SPOT_ZONE`.
- **Churn is severe** — boxes preempt every ~5–13 min, often mid-run. A fresh
  install (~8min) + model download (~3min) + eval rarely fits one window.
- **Fix: GCS cache** at `gs://medusa-experiments/v3-cache/` (`vllm_env_uv.tar` =
  venv + `.local/share/uv` + uv; `qwen35_4b_hf.tar` = the model). `box_deploy_v3.sh`
  restores+validates them (falls back to fresh install on bad cache) → on-box prep
  drops to ~1–2 min. Same v6e image across zones, so the cache works in either zone.
- **Launch on-box work in tmux**, not `nohup`/`setsid` — over ssh the latter two got
  SIGHUP'd and died with an empty log. `tmux new-session -d` survives disconnect.
  Do NOT `pkill -f box_deploy_v3.sh` in the launching ssh command — the command line
  itself contains that string, so pkill self-matches and kills the launcher before
  tmux starts. `tmux kill-session -t v3` alone is enough.
- **Never kill the auto-retry loop while `tpu-spot.sh recover` is mid-DELETE** — it
  leaves the node DELETING→deleted and not recreated. Stop the loop only between
  attempts, or let recover finish first.
- To fit short windows, lower `LIMIT` (lm_eval `--limit`); `LIMIT=3` keeps the
  0-vs-0.33 signal but cuts the v3 eval to ~2min. Drive/loop pass it through.

Related: [[project_gdn_v3_chunk128_overflow]].
