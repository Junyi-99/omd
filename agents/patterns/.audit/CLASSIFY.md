你是 master 分类 agent。32名独立reviewer已完成，报告位于 agents/patterns/reviewer-01.md 至 reviewer-32.md。请阅读全文（可分批）并依据可复用方法归类。
执行操作：
1. 只在 patterns/ 下创建有意义的 category slug 文件夹，通常5–9个；按实际内容决定，不强制数目。每份报告按主要方法选择唯一主类，移动到 patterns/<category>/reviewer-NN.md，保留原始报告内容。不要删减、拆散、复制报告。跨主题在分类索引列次级标签。不要把 .audit 当category。
2. 写 patterns/CLASSIFICATION.md：解释分类标准、各目录范围、全部32文件的分配表、每份次级主题、跨类候选归属建议，避免重复创建skill。特别注意偏好/项目事实不等于skill。
3. 写 .audit/categories.json，以JSON数组列出实际category目录名，供下一轮每类1个agent复核。
4. 不创建FINAL.md，不替代下一轮原文复核，不创建skill。禁止改动源memory/聊天记录，不动其他目录。所有路径相对 /Users/junyi/.oh-my-dotfiles。
5. 完成时核验32份报告恰好各有一份，输出分类数N、目录名及每类数量。
