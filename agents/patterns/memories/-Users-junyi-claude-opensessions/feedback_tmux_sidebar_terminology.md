---
name: tmux-sidebar terminology mapping
description: User wants tmux concepts displayed with custom names — Session→Project, Window→Agent Session, Pane→auxiliary agent/tool
type: feedback
originSessionId: 51db8f58-6bc2-4f4a-bdaf-e1ee9816c19f
---
tmux-sidebar 项目中，用户要求 UI 展示层使用自定义术语，代码内部保留 tmux 原生术语。

映射：
- tmux Session → Project
- tmux Window → Agent Session
- tmux Pane → 辅助 agent / 工具（sidebar 当前不展示 pane 级别信息）

**Why:** 用户把 tmux 当作多 agent 工作空间来用，Session=项目，Window=该项目下的 agent 会话。
**How to apply:** 所有面向用户的文本（header、footer、提示文案）用 Project/Agent Session，代码变量名和 tmux API 调用保持 session/window/pane。
