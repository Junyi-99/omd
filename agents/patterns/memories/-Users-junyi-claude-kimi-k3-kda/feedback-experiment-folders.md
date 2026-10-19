---
name: feedback-experiment-folders
description: 每轮大实验的产物必须单独建目录（exp-<日期>-<主题>/），不与历史实验混在 prof/
metadata: 
  node_type: memory
  type: feedback
  originSessionId: e645184a-6feb-471a-ae36-e740b711b0e0
  modified: 2026-08-17T05:29:39.697Z
---

用户要求（2026-08-17）：一轮实验的所有数据、日志、图表、脚本、原始 capture 单独放一个目录，
不要和之前的实验混在 kda-lower-bound/prof/ 里。

**Why:** prof/ 积累了几十个实验的产物后无法分辨哪些属于哪轮；实验包要能整体交付/归档。

**How to apply:** 目录名 `exp-<YYYY-MM-DD>-<主题>/`，内含 README.md（文件索引）+ 主报告 md
（从原始 JSON 程序化生成，数字零手抄）+ 数据 json + 日志 + fig/ + 采集与画图脚本 + 原始 capture。
报告内部用相对路径。首例：kda-lower-bound/exp-2026-08-17-pr2-ablation/。
相关：[[feedback-save-scripts-locally]]
