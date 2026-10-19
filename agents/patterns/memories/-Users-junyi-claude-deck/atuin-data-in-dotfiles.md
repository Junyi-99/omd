---
name: atuin-data-in-dotfiles
description: atuin 数据/配置在 ~/.oh-my-dotfiles 下，~/.atuin 只是二进制；别把 ~/.local/share/atuin 不存在当成数据丢失
metadata: 
  node_type: memory
  type: reference
  originSessionId: 73e77a24-372b-4e45-8c06-4e10fec311da
  modified: 2026-08-11T08:55:12.544Z
---

atuin 的真实路径（dotfiles 重定向，非默认 XDG）：

- 数据库：`~/.oh-my-dotfiles/data/atuin/history.db`（含 key、records.db）
- 配置：`~/.oh-my-dotfiles/etc/atuin/config.toml`
- `~/.atuin/` 只放二进制（官方 installer 路径）；zshrc 以 `[ -f ~/.atuin/bin/env ]` 为开关，二进制丢了 hook 静默失效但数据无损
- 用 `atuin info` 查实际路径，别猜 `~/.local/share/atuin`
- 未登录 sync（data 目录下无 session 文件）
- 2026-08-11：`~/.atuin` 被不明 session 清空过一次，重装即恢复；shell 集成在 `~/.oh-my-dotfiles/etc/zsh/zshrc:97`（带 `--disable-up-arrow`），installer 塞进 `~/.zshrc` 的裸 `atuin init zsh` 要删
