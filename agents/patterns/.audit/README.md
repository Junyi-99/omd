# 审计材料

- inventory.json：原始来源文件清单。
- reviewer-*.json：32份分片清单，每个源文件恰好分配一次。
- REVIEW.md / CLASSIFY.md / VERIFY.md：各阶段执行说明。
- categories.json：7个实际分类目录。
- validation.json：最终文件完整性和分配核验。

审阅期间使用的 user-index.json 与 user-turns/ 是从原始 JSONL 提取的临时全文索引，完成后已移除；历史执行说明中的相关路径不再存在。证据仍可从各报告引用的原始路径和行号回查。
