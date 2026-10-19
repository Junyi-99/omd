#!/bin/bash
# Auto-detect macOS light/dark mode and apply tmux theme

mode=$(defaults read -g AppleInterfaceStyle 2>/dev/null || echo "Light")

if [ "$mode" = "Dark" ]; then
    # --- Tokyo Night Dark ---
    bg="#1a1b26"
    fg="#a9b1d6"
    accent="#7aa2f7"
    secondary="#bb9af7"
    dim="#7982a9"
    border="#3b4261"
    alert="#f7768e"
else
    # --- Tokyo Night Light ---
    bg="#f0f1f5"
    fg="#343b58"
    accent="#34548a"
    secondary="#5a4a78"
    dim="#646b82"
    border="#b4b5b9"
    alert="#c0392b"
fi

# --- 状态栏 ---

tmux set -g status-style "bg=$bg,fg=$fg"
tmux set -g status-left "#[fg=$bg,bg=$accent,bold]  #S #[fg=$accent,bg=$bg] "
tmux set -g status-right "#{?pane_synchronized,#[bg=$alert]#[fg=#ffffff]#[bold] SYNC #[default],}#{?window_zoomed_flag,#[bg=#2ecc71]#[fg=#ffffff]#[bold] ZOOM #[default],}#[fg=$accent,bg=$bg]#[fg=$bg,bg=$accent,bold] #h "
tmux setw -g window-status-format "#{?window_bell_flag,#[fg=#ffffff bg=$alert bold] 🔔 #I:#W #[default bg=$bg],#[fg=$dim] #I:#W }"
tmux setw -g window-status-current-format "#[fg=$bg,bg=$secondary,bold] #I:#W #[fg=$secondary,bg=$bg]"
tmux setw -g window-status-activity-style "fg=$accent,bg=$bg"

# 在窗口列表末尾（紧贴最后一个 tab）插入可点击的 + 按钮（tmux.conf 的 MouseDown1Status 按 range=user|newwin 识别）
# _sf0 是 tmux 内置默认 status-format[0]，右对齐段前已插入 + 按钮；__PC__ 为配色占位
_sf0='#[align=left range=left #{E:status-left-style}]#[push-default]#{T;=/#{status-left-length}:status-left}#[pop-default]#[norange default]#[list=on align=#{status-justify}]#[list=left-marker]<#[list=right-marker]>#[list=on]#{W:#[range=window|#{window_index} #{E:window-status-style}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-format}#[pop-default]#[norange default]#{?loop_last_flag,,#{window-status-separator}},#[range=window|#{window_index} list=focus #{?#{!=:#{E:window-status-current-style},default},#{E:window-status-current-style},#{E:window-status-style}}#{?#{&&:#{window_last_flag},#{!=:#{E:window-status-last-style},default}}, #{E:window-status-last-style},}#{?#{&&:#{window_bell_flag},#{!=:#{E:window-status-bell-style},default}}, #{E:window-status-bell-style},#{?#{&&:#{||:#{window_activity_flag},#{window_silence_flag}},#{!=:#{E:window-status-activity-style},default}}, #{E:window-status-activity-style},}}]#[push-default]#{T:window-status-current-format}#[pop-default]#[norange list=on default]#{?loop_last_flag,,#{window-status-separator}}}#[nolist]#[range=user|newwin]#[__PC__] + #[norange]#[nolist align=right range=right #{E:status-right-style}]#[push-default]#{T;=/#{status-right-length}:status-right}#[pop-default]#[norange default]'
tmux set -g 'status-format[0]' "${_sf0/__PC__/fg=$secondary,bg=$bg,bold}"

# --- 面板、消息、菜单 ---

tmux set -g pane-border-style "fg=$border"
tmux set -g pane-active-border-style "fg=$accent"
tmux set -g display-panes-colour "$dim"
tmux set -g display-panes-active-colour "$accent"
tmux set -g message-style "fg=$accent,bg=$bg,bold"
tmux set -g message-command-style "fg=$secondary,bg=$bg"
tmux setw -g mode-style "fg=$bg,bg=$secondary,bold"
tmux set -g menu-style "fg=$fg,bg=$bg"
tmux set -g menu-selected-style "bg=$accent,fg=$bg,bold"
tmux set -g menu-border-style "fg=$border,bg=$bg"
tmux set -g menu-border-lines rounded

# --- 浮层与脚本共用的主题色 ---

# 浮层（popup）统一风格
tmux set -g popup-border-lines rounded
tmux set -g popup-border-style "fg=$accent,bg=$bg"
tmux set -g popup-style "bg=$bg,fg=$fg"

# 把主题色暴露给 new-item-prompt.py（弹窗配色）和 sidebar.sh（高亮/提醒）
tmux set -g @popup_accent "$accent"
tmux set -g @popup_dim "$dim"
tmux set -g @popup_fg "$fg"
tmux set -g @popup_alert "$alert"

# 侧边栏（sidebar.sh）选中项配色
tmux set -g @theme_bg "$bg"
tmux set -g @theme_secondary "$secondary"
