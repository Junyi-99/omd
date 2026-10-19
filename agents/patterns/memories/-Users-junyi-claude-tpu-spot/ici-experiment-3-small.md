---
name: ici-experiment-3-small
description: TPU ICI experiment 3 (E6 small single transfers, 1-4 B..4 KiB) — specs, where it ran, why the on-demand box was not used
metadata:
  type: project
---

2026-10-13: added shards e6-p{s}-{d}-{A|B} to ~/claude/ici-bench/bench2.py (path A ppermute 2 B..4 KiB; path B one descriptor (1,1)..(32,128) int8 = 1 B..4 KiB) so TPU curves start where GPU/RoCE curves start. Analysis analyze3.py -> out-r3/. The on-demand tpu-v6e-4-od-us-east5-a was busy with the user's deciax.serve clef eval (not touched); spot attempts: ici-v6e4-2x2-r3-small-spot (us-east5-a, no capacity at first) and ici-v6e4-2x2-r3-small-use1d (us-east1-d), both Medusa project, cleanup_on_success. Stop the slower one once one runs. medusa-compute has no spot v6e quota now; Medusa project has 1536 chips in us-east1-d / europe-west4-a / us-east5-a. Don't change the global gcloud project (other sessions use it).
Result 2026-10-13 (ici-v6e4-2x2-r3-small-od, on-demand VM after the clef server exited, 24/24 shards, 300/300 configs ok): x-axis p50 path A 2.90 us flat 4 B-1 KiB (2 B 2.75), matches exp1 (0.985-1.008); path B single descriptor 2.27 us flat 2 B-1 KiB (1 B 2.22) but 8-12% above E2 at 4 KiB because the grain kernel differs from the flow kernel -> plotted separately, not spliced. Spot runs were stopped without ever getting a VM.
Related: [[ici-v6e-measured-constants]], [[ici-external-interconnect-refs]].
