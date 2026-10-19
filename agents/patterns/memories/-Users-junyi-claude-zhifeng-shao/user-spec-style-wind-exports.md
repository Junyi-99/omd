---
name: user-spec-style-wind-exports
description: 用户给需求的方式：编号提示词 + Wind SQL + 导出 CSV 附件放 assets/；本地 Wind 导出多为 GB18030
metadata:
  type: user
---

用户（Junyi）给 zhifeng-shao 项目需求时，习惯写编号提示词，每条附上在本地 Wind 库执行的 SQL
和导出的结果 CSV（放在 `assets/附件N....csv`，用 @ 引用）。附件即数据格式的权威样本：
列名、编码、区间字段语义都以附件核对为准，不要凭记忆假设 Wind 表结构。

**Why:** 2026-08-25 的中信行业附件是 GB18030 且为联表导出（重复列名、名称内联），
和"每表一文件"的默认约定不同；先核对附件再写代码避免了返工。

**How to apply:** 收到带附件的需求先用 pandas 读附件（试 utf-8 再 gb18030）检查列、行数、
区间端点语义，再改代码；把核对结论写进 README 的数据约定部分。
