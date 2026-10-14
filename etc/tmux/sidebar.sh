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
#   menu <mouse_x> <mouse_y> <pane_left> <pane_top> <client>   右键菜单
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

# pane 内第 $1 行对应什么：session <sid> / window <sid> <wid> / new（“+ new window”行）/ 空
row_at() {
  local y=0 sid wid prev=""
  while IFS="$SEP" read -r sid wid _; do
    if [ "$sid" != "$prev" ]; then
      [ "$y" -eq "$1" ] && { echo "session $sid"; return; }
      y=$((y + 1))
      prev="$sid"
    fi
    [ "$y" -eq "$1" ] && { echo "window $sid $wid"; return; }
    y=$((y + 1))
  done < <(windows_tree '')
  [ "$y" -eq "$1" ] && echo new
}

new_item_popup() { # $1=client $2=window|session
  local dir
  dir="$(tmux display -p -c "$1" '#{pane_current_path}')"
  tmux display-popup -c "$1" -E -w 64 -h 14 -x C -y C -T " New $2 " -d "$dir" "'$CONFIG_DIR/new-window-prompt.sh' $2"
}

click() { # $1=pane 内行号 $2=client
  local kind sid wid
  read -r kind sid wid < <(row_at "$1")
  case "$kind" in
    session) tmux switch-client -c "$2" -t "$sid" ;;
    window) tmux switch-client -c "$2" -t "$sid:$wid" ;;
    new) new_item_popup "$2" window ;;
  esac
}

# 右键菜单：点 session 行 / 窗口行 / 其他位置，各给一套菜单。
# 由脚本弹出（只有脚本知道每行对应什么），拿不到鼠标事件，所以用绝对坐标定位。
# 菜单的 -t 必须是客户端当前看得见的 pane（否则 tmux 不显示菜单），
# 所以每条命令都显式带 -t，名字也由脚本填入而不是用 #{} 格式
# 名字放进菜单的格式字符串和双引号参数里：# 要写成 ##，" 和 \ 要转义
fmt_text() {
  local t="${1//\\/\\\\}"
  t="${t//\"/\\\"}"
  printf '%s' "${t//#/##}"
}

menu() { # $1,$2=鼠标在 pane 内的坐标 $3,$4=pane 左上角 $5=client
  local kind sid wid target name idx items=() prev="" next="" w found=0 width refresh="run-shell -b '$SELF refresh'"
  read -r kind sid wid < <(row_at "$2")
  target="$(tmux display -p -c "$5" '#{pane_id}')"
  case "$kind" in
    session)
      name="$(fmt_text "$(tmux display -p -t "$sid" '#{session_name}')")"
      items=(-T "#[bold] $name "
        "  ⇄ Switch to"   s "switch-client -c '$5' -t '$sid'"
        "  + New Window"  n "new-window -t '$sid:' -c \"#{pane_current_path}\""
        "  ✎ Rename"      r "command-prompt -I \"$name\" -p 'Rename session:' \"rename-session -t '$sid' -- '%%'\""
        ""
        "  ✕ Kill Session" x "confirm-before -p \"Kill session $name? (y/n)\" \"kill-session -t '$sid'\"")
      ;;
    window)
      name="$(fmt_text "$(tmux display -p -t "$sid:$wid" '#{window_name}')")"
      idx="$(tmux display -p -t "$sid:$wid" '#{window_index}')"
      # 同一 session 里的前后窗口，用于上移/下移
      for w in $(tmux list-windows -t "$sid" -F '#{window_id}'); do
        if [ "$w" = "$wid" ]; then found=1
        elif [ "$found" = 0 ]; then prev="$w"
        elif [ -z "$next" ]; then next="$w"
        fi
      done
      # 名字以 - 开头的条目显示为禁用：第一个窗口不能上移，最后一个不能下移
      local up="  ↑ Move Up" down="  ↓ Move Down"
      [ -z "$prev" ] && up="-$up"
      [ -z "$next" ] && down="-$down"
      items=(-T "#[bold] $idx: $name "
        "  ⇄ Switch to" s "switch-client -c '$5' -t '$sid:$wid'"
        "  ✎ Rename"    r "command-prompt -I \"$name\" -p 'Rename window:' \"rename-window -t '$wid' -- '%%'\""
        "$up"   u "swap-window -d -s '$wid' -t '$prev' ; $refresh"
        "$down" d "swap-window -d -s '$wid' -t '$next' ; $refresh"
        ""
        "  ✕ Kill Window" x "confirm-before -p \"Kill window $name? (y/n)\" \"kill-window -t '$wid'\"")
      ;;
    *)
      width="$(tmux show -gqv @sidebar_width)"
      items=(-T "#[bold] Sidebar "
        "  + New Window"  n "run-shell -b '$SELF new-item $5 window'"
        "  + New Session" N "run-shell -b '$SELF new-item $5 session'"
        ""
        "  ⇔ Width…"      w "command-prompt -I '${width:-$DEFAULT_WIDTH}' -p 'Sidebar width:' \"run-shell '$SELF %%'\""
        "  ○ Hide Sidebar" b "run-shell '$SELF'")
      ;;
  esac

  # numeric -y 是菜单底边（不含），高度 = 条目数（含分隔线）+ 上下边框；超出底部时贴底
  local n seps=0 i h y ch
  for ((i = 2; i < ${#items[@]}; i++)); do [ -z "${items[$i]}" ] && seps=$((seps + 1)); done
  n=$(((${#items[@]} - 2 - seps) / 3 + seps))
  h=$((n + 2))
  ch="$(tmux display -p -c "$5" '#{client_height}')"
  y=$(($4 + $2 + h))
  [ "$y" -gt "$ch" ] && y="$ch"
  # 不是从鼠标绑定直接弹出的菜单，tmux 默认当作纯键盘菜单，任何非左键事件（包括松开右键）都会关闭它：
  # -M 打开菜单的鼠标处理；-O 让松开右键时菜单留在原地，再点选项
  tmux display-menu -O -M -c "$5" -t "$target" -x $(($3 + $1)) -y "$y" "${items[@]}"
}

# --- 绘制 ---

fg() { [ -n "$1" ] && printf '\e[38;2;%d;%d;%dm' "0x${1:1:2}" "0x${1:3:2}" "0x${1:5:2}"; }
bg() { [ -n "$1" ] && printf '\e[48;2;%d;%d;%dm' "0x${1:1:2}" "0x${1:3:2}" "0x${1:5:2}"; }

# 补齐/截断到 n 列。名字可能含中文等双宽字符，交给 tmux 按显示宽度处理
fitw() { printf '#{p%d:#{=%d:%s}}' "$1" "$1" "$2"; }

# 仅用于纯 ASCII 文本
fit() { printf '%-*.*s' "$W" "$W" "$1"; }

# 输出整帧。返回 0 = 画好且看得见，3 = 画好但在后台（render 据此决定刷新间隔），
# 1 = 该退出，2 = 这次没拿到信息、跳过。
# 必须在子 shell 里调用：bash 3.2 的 <(...) 每次泄漏一个文件描述符，常驻进程里会耗尽（Too many open files）
# 后台窗口的侧边栏也保持最新：切过去时画面已经是对的，不会先露出旧画面/空白再跳变。
# 高亮按“侧边栏所在的窗口”而不是“session 当前窗口”，这样后台侧边栏的画面在切过去之前就已正确
draw() {
  local mypane mysid mywid npanes visible sync zoom host accent second dim tbg alert
  IFS="$SEP" read -r mypane mysid mywid npanes visible W H < <(
    tmux display -p -t "$TMUX_PANE" "#{pane_id}$SEP#{session_id}$SEP#{window_id}$SEP#{window_panes}$SEP#{&&:#{window_active},#{session_attached}}$SEP#{pane_width}$SEP#{pane_height}"
  )
  # 只有确实读到“窗口里只剩侧边栏”（主 pane 已退出）才退出，让窗口一起关闭；
  # 读失败（tmux 命令出错等）只跳过这一帧，不能把侧边栏自己关掉
  # display-message 找不到目标 pane 时会退回别的 pane 或返回空，不是自己就跳过
  [ "$mypane" = "$TMUX_PANE" ] || return 2
  case "$npanes" in
    '' | *[!0-9]*) return 2 ;;
    1) return 1 ;;
  esac
  [ -n "$H" ] && [ -n "$W" ] || return 2
  IFS="$SEP" read -r sync zoom host accent second dim tbg alert < <(
    tmux display -p -t "$TMUX_PANE" "#{synchronize-panes}$SEP#{window_zoomed_flag}$SEP$(fitw $((W - 1)) host_short)$SEP#{@popup_accent}$SEP#{@popup_dim}$SEP#{@theme_secondary}$SEP#{@theme_bg}$SEP#{@popup_alert}"
  )

  # session 行：" " + 图标 + " " + 名字（当前 session 用打开的文件夹，其他用合上的，Nerd Font 图标）；窗口行："   " + 序号（2 列）+ " " + 名字，共 W 列；有提醒时名字让出 2 列给标记
  local R=$'\e[0m' B=$'\e[1m' lines=() sid wid sname idx name name2 active bell prev=""
  local dimc="${dim:-#7982a9}"
  while IFS="$SEP" read -r sid wid sname idx name name2 active bell; do
    if [ "$sid" != "$prev" ]; then
      if [ "$sid" = "$mysid" ]; then
        lines+=("$B$(fg "$accent")  $sname$R")
      else
        lines+=("$B$(fg "$dimc")  $sname$R")
      fi
      prev="$sid"
    fi
    if [ "$wid" = "$mywid" ] && [ "$sid" = "$mysid" ]; then
      lines+=("$B$(bg "$second")$(fg "$tbg")   $idx $name$R")
    elif [ "$bell" = 1 ]; then
      lines+=("$B$(fg "$alert")   $idx $name2 !$R")
    elif [ "$active" = 1 ] && [ "$sid" != "$mysid" ]; then
      lines+=("$(fg "$second")   $idx $name$R")
    else
      lines+=("$(fg "$dimc")   $idx $name$R")
    fi
  done < <(windows_tree "$(fitw $((W - 3)) session_name)$SEP#{p2:window_index}$SEP$(fitw $((W - 6)) window_name)$SEP$(fitw $((W - 8)) window_name)$SEP#{window_active}$SEP#{window_bell_flag}")
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
  [ "$visible" = 1 ] && return 0
  return 3
}

render() {
  local prev="" frame rc visible=1
  FRAME_FILE="$(mktemp -t tmux-sidebar)"
  trap 'rm -f "$FRAME_FILE"; tmux set -wu -t "$TMUX_PANE" @sidebar_pane 2>/dev/null' EXIT
  trap 'exit' HUP TERM
  # 不 trap WINCH：被捕获的信号会打断 draw 里的 read，读到空值。尺寸变化由 hook 唤醒重绘
  # 备用屏幕：鼠标滚轮不会把侧边栏带进 copy-mode；隐藏光标、关闭回显
  printf '\e[?1049h\e[?25l\e]2;sidebar\e\\'
  stty -echo -icanon 2>/dev/null
  while :; do
    # 在子 shell 里画（见 draw 的注释），出错信息不要画到侧边栏上
    (draw) > "$FRAME_FILE" 2>/dev/null
    rc=$?
    case $rc in
      0 | 3) frame="$(<"$FRAME_FILE")"; visible=$((rc == 0)) ;;
      1) break ;;
      *) # 这次没读到完整信息，保持原样，下一轮再画；但若自己的 pane 已不存在（或 tmux 已退出）就退出
        # display-message 找不到目标时不报错，只返回空，所以要比对返回的 pane id
        [ "$(tmux display -p -t "$TMUX_PANE" '#{pane_id}' 2>/dev/null)" = "$TMUX_PANE" ] || break
        frame="$prev" ;;
    esac
    # 画完把光标停在左上角：pane 变矮时 tmux 先从光标下方删行，树不会被整体上推
    if [ "$frame" != "$prev" ]; then
      printf '\e[H%s\e[J\e[H' "$frame"
      prev="$frame"
    fi
    # 兜底定时重绘（看得见的每秒，后台的每 5 秒）；hook 通过 send-keys 唤醒以立即刷新
    # 注意 bash 3.2 的 read 超时也返回 1，没法据此判断终端是否已关；
    # pane 被关后 read 会立刻返回，靠上面“pane 不存在就退出”的检查结束，不会空转
    read -rsn1 -t "$([ "$visible" = 1 ] && echo 1 || echo 5)" _
  done
}

# 入口放进函数，并与 exit 写在同一行：bash 边执行边读脚本，
# 运行中的进程（侧边栏常驻）在脚本被修改后不会接着读到改过的内容
main() {
  case "$1" in
    render) render ;;
    refresh) refresh ;;
    click) click "$2" "$3"; exit 0 ;;
    menu) menu "$2" "$3" "$4" "$5" "$6"; exit 0 ;;
    new-item) new_item_popup "$2" "$3"; exit 0 ;;
    fix-width)
      # 尺寸/布局变化：修正宽度，并唤醒该窗口的侧边栏按新尺寸重绘（不给窗口则全部）
      enabled || exit 0
      fix_width "$2"
      if [ -n "$2" ]; then
        p="$(sidebar_of "$2")"
        [ -z "$p" ] || tmux send-keys -t "$p" r
      else
        refresh
      fi
      exit 0
      ;;
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
}

main "$@"; exit
