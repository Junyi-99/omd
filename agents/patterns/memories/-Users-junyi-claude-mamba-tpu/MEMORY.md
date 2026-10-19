# Memory index

- [Don't touch user's title or opening](feedback_dont_change_user_titles.md) — title + first 1–2 setup sections stay verbatim during any "make it professional" rewrite; ask before injecting TL;DR/abstract boxes.
- [Cut setup-section bloat aggressively](feedback_writing_cut_setup_bloat.md) — merge adjacent same-topic sections, push methodology footnotes to appendix, sweep for duplicate Figure cites and "we'll explain X below" sentences.
- [Prefill-scope terminology precision](feedback_prefill_terminology.md) — never use "serving" in prefill-only kernel docs (conflates prefill+decode). Use "prefill" / "jagged-prefill" / "packed-jagged input" instead. Overrides the protected-opening rule for surgical word swaps.
- [v3 GDN CHUNK=128 overflow](project_gdn_v3_chunk128_overflow.md) — v3 fused prefill gives mmlu_pro=0 in vllm because the 128-chunk (I+S)⁻¹ overflows to NaN; fix = in-tile 2×64-subchunk recurrence (only 64×64 inverses). Math validated offline.
- [Spot v6e-4 churn + caching ops](project_spot_v6e_churn_caching.md) — spot quota only in us-east1-d / europe-west4-a; GCS cache for fast recovery; launch via tmux (not nohup); don't kill the loop mid-recover-DELETE.
