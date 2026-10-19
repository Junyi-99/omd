---
name: paper-kb-for-kernel-gen
description: User is writing a paper on how to better build a knowledge base for kernel generation
metadata: 
  node_type: memory
  type: project
  originSessionId: 26f2cdee-7374-4cca-bd3f-bf8eb66c7eae
---

User is writing a research paper on **how to better build a knowledge base (KB) for kernel generation** (GPU/ASIC kernel design & optimization). They want related-work coverage across four axes, all of interest:
1. KB content / representation — what to store (slow-fast kernel pairs, reusable skills, insights, structured wiki).
2. Data augmentation / synthetic data — generating & verifying kernel training data, curricula, anti reward-hacking.
3. Retrieval / reuse mechanism — how the KB feeds back into generation (in-context, textual gradient, skill invocation).
4. Cross-hardware transfer — reusing KB/data across GPU↔ASIC/NPU.

**How to apply:** In the `research/` wiki, the relevant clusters are the `memory-and-knowledge-reuse` concept (10 papers) and the RL-for-kernels synthetic-data papers. Closest-fit papers: [[kda-kernel-design-agents]] (KernelWiki), accelopt (optimization memory), kernelskill (dual-level skill memory). Deliverable agreed: a dedicated synthesis page `knowledge-base-for-kernel-generation.md` organizing papers by the four axes + a gap analysis that motivates the new paper. Build it after the source-page deepening run completes.
