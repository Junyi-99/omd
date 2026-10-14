#!/bin/bash
# 左侧竖向侧边栏：在每个窗口最左边开一个窄 pane，竖着列出所有 session 及其窗口，替代底部状态栏
#
# 用法：
#   sidebar.sh            开/关侧边栏（全局生效；开启时隐藏底部状态栏）
#   sidebar.sh <width>    以指定宽度开启；已开启时调整宽度
# 内部（由 tmux.conf 的 hook / 鼠标绑定调用）：
#   hook <window_id>      若已开启，给该窗口补上侧边栏并刷新
#   refresh               唤醒所有侧边栏立即重绘
#   fix-width [window_id] 把侧边栏宽度拉回 @sidebar_width（不给窗口则处理全部）
#   adopt-width <window_id> 鼠标拖完边框后，以该窗口侧边栏的新宽度作为全局宽度
#   click <y> <client>
#   render                侧边栏 pane 里运行的绘制循环

SELF="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/$(basename "${BASH_SOURCE[0]}")"
CONFIG_DIR="$(dirname "$SELF")"
DEFAULT_WIDTH=24
SEP=$'\x1f'    # 非空白分隔符，空字段不会被 read 合并

export LC_CTYPE="${LC_CTYPE:-en_US.UTF-8}"

enabled() { [ "$(tmux show -gqv @sidebar_enabled)" = 1 ]; }

all_sidebars() { tmux list-panes -a -f '#{@is_sidebar}' -F '#{pane_id}' 2>/dev/null; }

sidebar_of() { tmux list-panes -t "$1" -f '#{@is_sidebar}' -F '#{pane_id}' 2>/dev/null | head -n1; }

# 新建窗口会同时触发多个 hook，加锁避免同一窗口被补上两个侧边栏
add_to_window() {
  tmux wait-for -L sidebar-lock
  _add_to_window "$1"
  tmux wait-for -U sidebar-lock
}

_add_to_window() {
  [ -n "$(sidebar_of "$1")" ] && return
  local width pane
  width="$(tmux show -gqv @sidebar_width)"
  pane="$(tmux split-window -hbfd -l "${width:-$DEFAULT_WIDTH}" -t "$1" -P -F '#{pane_id}' "exec '$SELF' render")" || return
  tmux set -p -t "$pane" @is_sidebar 1 \; set -w -t "$1" @sidebar_pane "$pane"
}

# tmux 在终端缩放、关闭相邻 pane、布局变化时会按比例调整 pane 宽度，
# 每个窗口的侧边栏各自漂移，所以由 hook 统一拉回固定宽度。
# 放大（zoom）中的窗口跳过：resize-pane 会取消 zoom
fix_width() {
  # 正在用鼠标拖边框：别跟用户抢，松手后由 adopt_width 统一
  [ "$(tmux show -gqv @sidebar_dragging)" = 1 ] && return 0
  local width pane w zoomed scope=(-a)
  [ -n "$1" ] && scope=(-t "$1")
  width="$(tmux show -gqv @sidebar_width)"
  width="${width:-$DEFAULT_WIDTH}"
  tmux list-panes "${scope[@]}" -f '#{@is_sidebar}' -F "#{pane_id}$SEP#{pane_width}$SEP#{window_zoomed_flag}" 2>/dev/null |
    while IFS="$SEP" read -r pane w zoomed; do
      [ "$w" != "$width" ] && [ "$zoomed" != 1 ] && tmux resize-pane -t "$pane" -x "$width"
    done
  return 0
}

adopt_width() {
  local w
  w="$(tmux list-panes -t "$1" -f '#{@is_sidebar}' -F '#{pane_width}' 2>/dev/null | head -n1)"
  [ -n "$w" ] && tmux set -g @sidebar_width "$w"
  fix_width
}

refresh() {
  local args=() p
  for p in $(all_sidebars); do args+=(send-keys -t "$p" r ';'); done
  [ ${#args[@]} -eq 0 ] || tmux "${args[@]}"
}

enable() {
  local width="${1:-$(tmux show -gqv @sidebar_width)}" win p
  tmux set -g @sidebar_width "${width:-$DEFAULT_WIDTH}" \; set -g @sidebar_enabled 1 \; set -g status off
  for win in $(tmux list-windows -a -F '#{window_id}' | sort -u); do add_to_window "$win"; done
  fix_width
}

disable() {
  local p
  tmux set -gu @sidebar_enabled \; set -g status on
  for p in $(all_sidebars); do tmux kill-pane -t "$p"; done
}

# 侧边栏的行布局：每个 session 一行标题，下面是它的窗口，最后一行“+ new window”。
# list-windows -a 按 session 再按窗口序号排列；draw 和 click 共用这个顺序
#   $1 = 每个窗口额外输出的字段（format），输出：session_id SEP window_id SEP 额外字段
windows_tree() { tmux list-windows -a -F "#{session_id}$SEP#{window_id}$SEP$1"; }

click() { # $1=pane 内行号 $2=client
  local y=0 sid wid prev="" dir
  while IFS="$SEP" read -r sid wid _; do
    if [ "$sid" != "$prev" ]; then
      [ "$y" -eq "$1" ] && { tmux switch-client -c "$2" -t "$sid"; return; }
      y=$((y + 1))
      prev="$sid"
    fi
    [ "$y" -eq "$1" ] && { tmux switch-client -c "$2" -t "$sid:$wid"; return; }
    y=$((y + 1))
  done < <(windows_tree '')
  if [ "$y" -eq "$1" ]; then
    dir="$(tmux display -p -c "$2" '#{pane_current_path}')"
    tmux display-popup -c "$2" -E -w 64 -h 14 -x C -y C -T ' New window ' -d "$dir" "'$CONFIG_DIR/new-window-prompt.sh' window"
  fi
}

# --- 绘制 ---

fg() { [ -n "$1" ] && printf '\e[38;2;%d;%d;%dm' "0x${1:1:2}" "0x${1:3:2}" "0x${1:5:2}"; }
bg() { [ -n "$1" ] && printf '\e[48;2;%d;%d;%dm' "0x${1:1:2}" "0x${1:3:2}" "0x${1:5:2}"; }

# 补齐/截断到 n 列。名字可能含中文等双宽字符，交给 tmux 按显示宽度处理
fitw() { printf '#{p%d:#{=%d:%s}}' "$1" "$1" "$2"; }

# 仅用于纯 ASCII 文本
fit() { printf '%-*.*s' "$W" "$W" "$1"; }

# 输出整帧；返回 1 表示该退出，返回 2 表示侧边栏当前不可见、无需重绘
draw() {
  local mysid npanes visible sync zoom host accent second dim tbg alert
  IFS="$SEP" read -r mysid npanes visible W H < <(
    tmux display -p -t "$TMUX_PANE" "#{session_id}$SEP#{window_panes}$SEP#{&&:#{window_active},#{session_attached}}$SEP#{pane_width}$SEP#{pane_height}"
  )
  # 窗口里只剩侧边栏（主 pane 已退出）→ 退出，让窗口一起关闭
  [ -n "$npanes" ] && [ "$npanes" -gt 1 ] || return 1
  # 后台窗口里的输出会被 tmux 记为 activity 并向终端发 BEL（整屏闪烁），所以只画看得见的侧边栏
  [ "$visible" = 1 ] || return 2
  IFS="$SEP" read -r sync zoom host accent second dim tbg alert < <(
    tmux display -p -t "$TMUX_PANE" "#{synchronize-panes}$SEP#{window_zoomed_flag}$SEP$(fitw $((W - 1)) host_short)$SEP#{@popup_accent}$SEP#{@popup_dim}$SEP#{@theme_secondary}$SEP#{@theme_bg}$SEP#{@popup_alert}"
  )

  # session 行：" " + 图标 + " " + 名字（当前 session 用打开的文件夹，其他用合上的，Nerd Font 图标）；窗口行："   " + 序号（2 列）+ " " + 名字，共 W 列；有提醒时名字让出 2 列给标记
  local R=$'\e[0m' B=$'\e[1m' lines=() sid wid sname idx name name2 active bell activity prev=""
  local dimc="${dim:-#7982a9}"
  while IFS="$SEP" read -r sid wid sname idx name name2 active bell activity; do
    if [ "$sid" != "$prev" ]; then
      if [ "$sid" = "$mysid" ]; then
        lines+=("$B$(fg "$accent")  $sname$R")
      else
        lines+=("$B$(fg "$dimc")  $sname$R")
      fi
      prev="$sid"
    fi
    if [ "$active" = 1 ] && [ "$sid" = "$mysid" ]; then
      lines+=("$B$(bg "$second")$(fg "$tbg")   $idx $name$R")
    elif [ "$bell" = 1 ]; then
      lines+=("$B$(fg "$alert")   $idx $name2 !$R")
    elif [ "$activity" = 1 ]; then
      lines+=("$(fg "$accent")   $idx $name2 •$R")
    elif [ "$active" = 1 ]; then
      lines+=("$(fg "$second")   $idx $name$R")
    else
      lines+=("$(fg "$dimc")   $idx $name$R")
    fi
  done < <(windows_tree "$(fitw $((W - 3)) session_name)$SEP#{p2:window_index}$SEP$(fitw $((W - 6)) window_name)$SEP$(fitw $((W - 8)) window_name)$SEP#{window_active}$SEP#{window_bell_flag}$SEP#{window_activity_flag}")
  lines+=("$(fg "$second")$(fit " +  new window")$R")

  # 底部：SYNC / ZOOM 标记 + 主机名；列表过长时截断列表，保留底部
  local footer=() i
  [ "$sync" = 1 ] && footer+=("$B$(bg "$alert")$(fg '#ffffff') SYNC $R")
  [ "$zoom" = 1 ] && footer+=("$B$(bg '#2ecc71')$(fg '#ffffff') ZOOM $R")
  footer+=("$B$(fg "$accent") $host$R")
  while [ $((${#lines[@]} + ${#footer[@]})) -gt "$H" ] && [ ${#lines[@]} -gt 0 ]; do unset "lines[${#lines[@]}-1]"; done
  while [ $((${#lines[@]} + ${#footer[@]})) -lt "$H" ]; do lines+=(""); done
  lines=("${lines[@]}" "${footer[@]}")

  for ((i = 0; i < H && i < ${#lines[@]}; i++)); do
    [ "$i" -gt 0 ] && printf '\n'
    printf '%s\e[K' "${lines[$i]}"
  done
}

render() {
  local prev="" frame
  trap 'tmux set -wu -t "$TMUX_PANE" @sidebar_pane 2>/dev/null' EXIT
  trap 'exit' HUP TERM
  # 备用屏幕：鼠标滚轮不会把侧边栏带进 copy-mode；隐藏光标、关闭回显
  printf '\e[?1049h\e[?25l\e]2;sidebar\e\\'
  stty -echo -icanon 2>/dev/null
  while :; do
    frame="$(draw)"
    case $? in
      1) break ;;
      # draw 期间可能刚切走窗口：输出前再确认一次仍可见，避免在后台窗口产生 activity
      0) if [ "$frame" != "$prev" ] && [ "$(tmux display -p -t "$TMUX_PANE" '#{&&:#{window_active},#{session_attached}}')" = 1 ]; then
           printf '\e[H%s\e[J' "$frame"
           prev="$frame"
         fi ;;
    esac
    # 每秒兜底重绘一次；hook 通过 send-keys 唤醒以立即刷新
    read -rsn1 -t 1 _
  done
}

case "$1" in
  render) render ;;
  refresh) refresh ;;
  click) click "$2" "$3"; exit 0 ;;
  fix-width) enabled && fix_width "$2"; exit 0 ;;
  adopt-width) enabled && adopt_width "$2"; exit 0 ;;
  hook)
    enabled || exit 0
    add_to_window "$2"
    fix_width "$2"
    refresh
    ;;
  '')
    if enabled; then disable; else enable; fi
    ;;
  *[!0-9]*)
    echo "usage: $0 [width]" >&2
    exit 1
    ;;
  *) enable "$1" ;;
esac
