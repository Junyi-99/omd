#!/bin/bash
# 左侧竖向侧边栏：在每个窗口最左边开一个窄 pane，竖着列出所有 session 及其窗口，替代底部状态栏
#
# 用法：
#   sidebar.sh                  开/关侧边栏（全局生效；开启时隐藏底部状态栏）
#   sidebar.sh <width>          以指定宽度开启；已开启时调整宽度
#
# 内部（由 tmux.conf 的 hook / 鼠标绑定，以及本脚本弹出的菜单调用）：
#   hook <window_id>            若已开启，给该窗口补上侧边栏并刷新
#   refresh                     唤醒所有侧边栏立即重绘
#   fix-width [window_id]       把侧边栏宽度拉回 @sidebar_width（不给窗口则处理全部）
#   adopt-width <window_id>     鼠标拖完边框后，以该窗口侧边栏的新宽度作为全局宽度
#   click <x> <y> <client>      左键：最后一行的“≡ menu”按钮弹出 Sidebar 菜单；session / 主机行点最左 3 列折叠，其余行见 row_at
#   collapse <key>              折叠/展开一个分组（key = session id 或主机名），菜单调用
#   menu <mouse_x> <mouse_y> <pane_left> <pane_top> <client>   右键菜单
#   new-item <client> window|session   新建窗口/会话的命名弹窗（菜单调用）
#   render                      侧边栏 pane 里运行的绘制循环
#   poll-loop / poll-host <host>                远程主机轮询（见“远程主机”）
#   radd/rremove <host>, rshell <host> <client>, rnew <host> <client> <name>,
#   rattach <host> <hex 名字> <client>, rkill <host> <hex 名字>   远程会话操作（菜单调用）

SELF="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/$(basename "${BASH_SOURCE[0]}")"
CONFIG_DIR="$(dirname "$SELF")"
DEFAULT_WIDTH=24
SEP=$'\x1f'    # 非空白分隔符，空字段不会被 read 合并
REMOTES_FILE="$CONFIG_DIR/sidebar-remotes"   # 远程主机列表：每行一个 ssh host 别名，# 开头为注释
CACHE_DIR="/tmp/tmux-sidebar-$UID"          # 远程轮询结果、轮询进程锁、ssh 复用连接
POLL_INTERVAL=15

export LC_CTYPE="${LC_CTYPE:-en_US.UTF-8}"

# --- 开关与宽度 ---

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
  ensure_poller
}

disable() {
  local p
  tmux set -gu @sidebar_enabled \; set -g status on
  for p in $(all_sidebars); do tmux kill-pane -t "$p"; done
}

# --- 远程主机 ---
# 侧边栏底部列出 $REMOTES_FILE 里每台主机上的 tmux session。由唯一一个后台轮询进程
# 每 $POLL_INTERVAL 秒 ssh 过去执行 tmux list-sessions，结果写到 $CACHE_DIR/remote-<host>
# （第一行状态 ok / down（连不上）/ nosrv（没有 tmux server 或没装 tmux），之后每行 名字|窗口数|attached 数），
# 侧边栏只读缓存文件，不自己 ssh。点远程 session 在本地开窗口 ssh -t 过去 tmux new -A（有则 attach），
# 窗口带 @remote=<host>:<session>，再点会切回这个窗口。

remote_hosts() { [ -f "$REMOTES_FILE" ] && sed -e 's/#.*//' -e 's/[[:space:]]//g' "$REMOTES_FILE" | grep -v '^$'; }

# 轮询和点开的交互 ssh 共用一条复用连接：attach 时不用再握手/认证，几乎瞬间进入
SSH_CM=(-o ControlMaster=auto -o ControlPath="$CACHE_DIR/cm-%C" -o ControlPersist=120)

# ssh 远端执行命令用的是非登录 shell，macOS 上 Homebrew 的 PATH（~/.zprofile 里的 brew shellenv）不会加载，
# 直接写 tmux 会 command not found（rc 127，被当成 nosrv）。所以远端命令统一补上常见安装路径再调 tmux。
REMOTE_TMUX='PATH="$PATH:/opt/homebrew/bin:/usr/local/bin:$HOME/.local/bin" tmux'

# BatchMode：不弹密码/确认提示，连不上就算 down
ssh_poll() { ssh -o BatchMode=yes -o ConnectTimeout=4 "${SSH_CM[@]}" "$@"; }

# 交互 ssh 的命令行（放进 new-window 里执行），参数已 shell 转义
ssh_cmd() {
  local a q=""
  for a in "${SSH_CM[@]}"; do q+="$(printf %q "$a") "; done
  printf 'ssh %s%s' "$q" "$*"
}

poll_host() {
  local h="$1" f="$CACHE_DIR/remote-$1" rc
  mkdir -p "$CACHE_DIR"
  # 输出写文件而不是 $(...)：ControlPersist 的后台 master 会继承管道，命令替换会一直等到它退出
  ssh_poll "$h" "$REMOTE_TMUX list-sessions -F '#{session_name}|#{session_windows}|#{session_attached}'" > "$f.out" 2>/dev/null < /dev/null
  rc=$?
  {
    case $rc in
      0) echo ok; cat "$f.out" ;;
      255) echo down ;;
      *) echo nosrv ;;
    esac
  } > "$f.tmp" && mv "$f.tmp" "$f"
  rm -f "$f.out"
}

poll_all() {
  local IFS=$' \t\n' h
  for h in $(remote_hosts); do poll_host "$h" & done
  wait
}

poller_running() {
  local pid
  pid="$(cat "$CACHE_DIR/poller.pid" 2>/dev/null)"
  [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null
}

ensure_poller() {
  [ -n "$(remote_hosts)" ] || return 0
  poller_running && return 0
  tmux run-shell -b "'$SELF' poll-loop > /dev/null 2>&1"
}

# 侧边栏关闭、tmux 退出、或主机列表清空后自动结束；mkdir 锁保证只有一个
poll_loop() {
  mkdir -p "$CACHE_DIR"
  if ! mkdir "$CACHE_DIR/poller.lock" 2>/dev/null; then
    poller_running && return 0
    rm -rf "$CACHE_DIR/poller.lock"
    mkdir "$CACHE_DIR/poller.lock" 2>/dev/null || return 0
  fi
  echo $$ > "$CACHE_DIR/poller.pid"
  trap 'rm -rf "$CACHE_DIR/poller.lock" "$CACHE_DIR/poller.pid"' EXIT
  while enabled && [ -n "$(remote_hosts)" ]; do
    poll_all
    refresh
    sleep "$POLL_INTERVAL"
  done
}

# 菜单命令里传 session 名字用十六进制，避免引号/空格问题
hexenc() { printf '%s' "$1" | xxd -p | tr -d '\n'; }
hexdec() { printf '%s' "$1" | xxd -r -p; }

# 已经开着的、对应 <host>:<session> 的本地窗口
remote_window() {
  tmux list-windows -a -F "#{window_id}$SEP#{@remote}" | awk -F "$SEP" -v k="$1" '$2 == k { print $1; exit }'
}

# 在 client 当前 session 开一个跑 $3 的窗口，名字/标记为 $2；已有同标记窗口就切过去
open_remote_window() { # $1=client $2=标记 $3=shell 命令
  local w sid
  w="$(remote_window "$2")"
  if [ -n "$w" ]; then
    tmux switch-client -c "$1" -t "$w"
    return
  fi
  sid="$(tmux display -p -c "$1" '#{session_id}')"
  w="$(tmux new-window -t "$sid:" -n "$2" -P -F '#{window_id}' "$3")" || return
  tmux set -w -t "$w" @remote "$2"
}

# 点侧边栏 / 菜单 Attach：只 attach 已有 session。侧边栏列表是轮询缓存，session 刚退出时还会列着，
# 这里不能用 new-session -A，否则点到过期条目会在远端悄悄新建一个同名 session。
# attach 失败（session 没了 / 连不上）就提示一下，并立刻重新轮询这台主机把过期条目清掉。
remote_attach() { # $1=host $2=session $3=client
  open_remote_window "$3" "$1:$2" \
    "$(ssh_cmd "-t $(printf %q "$1") $(printf %q "$REMOTE_TMUX attach-session -t $(printf %q "$2")")") || \
{ $(printf %q "$SELF") poll-host $(printf %q "$1") >/dev/null 2>&1; echo; echo '  session 不存在或连接已断开，2 秒后关闭'; sleep 2; }"
}

# 菜单 New Session：没有则新建，已有则 attach
remote_new() { # $1=host $2=session $3=client
  open_remote_window "$3" "$1:$2" \
    "$(ssh_cmd "-t $(printf %q "$1") $(printf %q "$REMOTE_TMUX new-session -A -s $(printf %q "$2")")")"
}

remote_shell() { open_remote_window "$2" "$1:" "$(ssh_cmd "$(printf %q "$1")")"; }

remote_kill() {
  ssh_poll "$1" "tmux kill-session -t $(printf %q "=$2")" > /dev/null 2>&1 < /dev/null
  poll_host "$1"
  refresh
}

remote_add() {
  [ -n "$1" ] || return
  remote_hosts | grep -qx "$1" || echo "$1" >> "$REMOTES_FILE"
  ensure_poller
  poll_host "$1"
  refresh
}

remote_remove() {
  local tmp
  [ -f "$REMOTES_FILE" ] || return
  tmp="$(mktemp)"
  awk -v h="$1" '{ l = $0; sub(/#.*/, "", l); gsub(/[[:space:]]/, "", l); if (l != h) print }' "$REMOTES_FILE" > "$tmp" &&
    cat "$tmp" > "$REMOTES_FILE"
  rm -f "$tmp" "$CACHE_DIR/remote-$1"
  refresh
}

# 远程区的行：sep / host SEP <host> SEP <状态> [SEP <折叠时隐藏的 session 数>] / rsession SEP <host> SEP <名字> SEP <窗口数> SEP <attached 数> SEP <本地已开窗口 1/0>
remote_rows() {
  # 调用方常写 IFS="$SEP" read … < <(row_at …)，bash 3.2 下这个临时 IFS 会带进 <(…)，这里显式还原
  local IFS=$' \t\n' hosts h f status s n a seen local_sessions rows cnt
  hosts="$(remote_hosts)"
  [ -n "$hosts" ] || return 0
  echo sep
  for h in $hosts; do
    f="$CACHE_DIR/remote-$h"
    seen=""
    # 本地已开着窗口（@remote=<host>:<session>）的远程 session，SEP 分隔
    local_sessions="$(tmux list-windows -a -F '#{@remote}' |
      awk -v h="$h:" 'index($0, h) == 1 && length($0) > length(h) { print substr($0, length(h) + 1) }' |
      sort -u | tr '\n' "$SEP")"
    rows="" cnt=0
    if [ ! -f "$f" ]; then
      status=pending
    else
      {
        read -r status
        if [ "$status" = ok ]; then
          while IFS='|' read -r s n a; do
            [ -n "$s" ] || continue
            if [[ "$SEP$local_sessions" == *"$SEP$s$SEP"* ]]; then
              rows+="rsession$SEP$h$SEP$s$SEP$n$SEP$a${SEP}1"$'\n'
            else
              rows+="rsession$SEP$h$SEP$s$SEP$n$SEP$a${SEP}0"$'\n'
            fi
            seen+="$s$SEP"
            cnt=$((cnt + 1))
          done
        fi
      } < "$f"
    fi
    # 本地已开窗口、但轮询缓存里还没有的远程 session（刚新建 / 主机离线）：也列出来，窗口数未知
    while IFS="$SEP" read -r -d "$SEP" s; do
      [ -n "$s" ] && [[ "$SEP$seen" != *"$SEP$s$SEP"* ]] && { rows+="rsession$SEP$h$SEP$s$SEP?${SEP}0${SEP}1"$'\n'; cnt=$((cnt + 1)); }
    done <<< "$local_sessions"
    if is_collapsed "$h"; then
      echo "host$SEP$h$SEP$status$SEP$cnt"
    else
      echo "host$SEP$h$SEP$status"
      printf '%s' "$rows"
    fi
  done
}

# --- 行布局、点击、菜单 ---

# 折叠的分组：@sidebar_collapsed 里空格分隔的 key（session id 如 $3，或远程主机名）。
# 折叠的 session 不列窗口、折叠的主机不列远程 session，标题行右侧显示 ▸ 数量
load_collapsed() { COLLAPSED=" $(tmux show -gqv @sidebar_collapsed) "; }
is_collapsed() { [[ "$COLLAPSED" == *" $1 "* ]]; }

toggle_collapse() {
  load_collapsed
  if is_collapsed "$1"; then
    COLLAPSED="${COLLAPSED/ $1 / }"
  else
    COLLAPSED+="$1 "
  fi
  tmux set -g @sidebar_collapsed "$(echo $COLLAPSED)"
  refresh
}

# 侧边栏的行布局：每个 session 一行标题，下面是它的窗口，然后一行“+ new window”，最后是远程区。
# list-windows -a 按 session 再按窗口序号排列；draw 和 click 共用这个顺序
#   $1 = 每个窗口额外输出的字段（format），输出：session_id SEP window_id SEP @remote SEP 额外字段
# 带 @remote 的本地窗口（ssh 到远程 session 的那些）不列在本地区，只在远程区高亮对应的 session / 主机
windows_tree() { tmux list-windows -a -F "#{session_id}$SEP#{window_id}$SEP#{@remote}$SEP$1"; }

# pane 内第 $1 行对应什么（字段用 SEP 分隔）：session <sid> / window <sid> <wid> / new（“+ new window”行）/
# sep / host <host> <状态> / rsession <host> <名字> …（远程区，见 remote_rows）/ 空
row_at() {
  local y=0 sid wid prev="" row
  load_collapsed
  while IFS="$SEP" read -r sid wid rem _; do
    if [ "$sid" != "$prev" ]; then
      [ "$y" -eq "$1" ] && { echo "session$SEP$sid"; return; }
      y=$((y + 1))
      prev="$sid"
    fi
    [ -n "$rem" ] && continue
    is_collapsed "$sid" && continue
    [ "$y" -eq "$1" ] && { echo "window$SEP$sid$SEP$wid"; return; }
    y=$((y + 1))
  done < <(windows_tree '')
  [ "$y" -eq "$1" ] && { echo new; return; }
  y=$((y + 1))
  while IFS= read -r row; do
    [ "$y" -eq "$1" ] && { echo "$row"; return; }
    y=$((y + 1))
  done < <(remote_rows)
}

new_item_popup() { # $1=client $2=window|session
  local dir
  dir="$(tmux display -p -c "$1" '#{pane_current_path}')"
  tmux display-popup -c "$1" -E -w 64 -h 14 -x C -y C -T " New $2 " -d "$dir" "'$CONFIG_DIR/new-window-prompt.sh' $2"
}

click() { # $1,$2=pane 内坐标 $3=client
  local kind sid wid p left top h
  p="$(sidebar_of "$(tmux display -p -c "$3" '#{window_id}')")"
  IFS="$SEP" read -r left top h < <(tmux display -p -t "$p" "#{pane_left}$SEP#{pane_top}$SEP#{pane_height}")
  # 底栏最后一行是“≡ menu”按钮：弹出 Sidebar 菜单（menu 里行号超出列表即走默认分支）
  if [ "$2" = "$((h - 1))" ]; then menu "$1" "$2" "$left" "$top" "$3"; return; fi
  IFS="$SEP" read -r kind sid wid _ < <(row_at "$2")
  case "$kind" in
    session | host) [ "$1" -lt 3 ] && { toggle_collapse "$sid"; return; } ;;   # 点图标：折叠/展开
  esac
  case "$kind" in
    session) tmux switch-client -c "$3" -t "$sid" ;;
    window) tmux switch-client -c "$3" -t "$sid:$wid" ;;
    new) new_item_popup "$3" window ;;
    host) remote_shell "$sid" "$3" ;;          # sid=host
    rsession) remote_attach "$sid" "$wid" "$3" ;; # sid=host wid=session
  esac
}

# 名字放进菜单的格式字符串和双引号参数里：# 要写成 ##，" 和 \ 要转义
fmt_text() {
  local t="${1//\\/\\\\}"
  t="${t//\"/\\\"}"
  printf '%s' "${t//#/##}"
}

# 右键菜单：点 session 行 / 窗口行 / 其他位置，各给一套菜单。
# 由脚本弹出（只有脚本知道每行对应什么），拿不到鼠标事件，所以用绝对坐标定位。
# 菜单的 -t 必须是客户端当前看得见的 pane（否则 tmux 不显示菜单），
# 所以每条命令都显式带 -t，名字也由脚本填入而不是用 #{} 格式
menu() { # $1,$2=鼠标在 pane 内的坐标 $3,$4=pane 左上角 $5=client
  local kind sid wid target name idx items=() prev="" next="" w found=0 width hex refresh="run-shell -b '$SELF refresh'"
  IFS="$SEP" read -r kind sid wid _ < <(row_at "$2")
  load_collapsed
  target="$(tmux display -p -c "$5" '#{pane_id}')"
  case "$kind" in
    session)
      name="$(fmt_text "$(tmux display -p -t "$sid" '#{session_name}')")"
      items=(-T "#[bold] $name "
        "  ⇄ Switch to"   s "switch-client -c '$5' -t '$sid'"
        "$(is_collapsed "$sid" && echo "  ▾ Expand" || echo "  ▸ Collapse")" c "run-shell -b '$SELF collapse $sid'"
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
    host) # sid=host
      items=(-T "#[bold] $sid "
        "  ⌁ Open Shell"   o "run-shell -b '$SELF rshell $sid $5'"
        "$(is_collapsed "$sid" && echo "  ▾ Expand" || echo "  ▸ Collapse")" c "run-shell -b '$SELF collapse $sid'"
        "  + New Session…" n "command-prompt -p 'New session on $sid:' \"run-shell -b '$SELF rnew $sid $5 %%'\""
        "  ↻ Refresh"      r "run-shell -b '$SELF poll-host $sid'"
        ""
        "  ✕ Remove Host"  x "confirm-before -p 'Remove $sid from sidebar? (y/n)' \"run-shell -b '$SELF rremove $sid'\"")
      ;;
    rsession) # sid=host wid=session
      name="$(fmt_text "$wid")"
      hex="$(hexenc "$wid")"
      items=(-T "#[bold] $sid: $name "
        "  ⇄ Attach"       a "run-shell -b '$SELF rattach $sid $hex $5'"
        ""
        "  ✕ Kill Session" x "confirm-before -p \"Kill $name on $sid? (y/n)\" \"run-shell -b '$SELF rkill $sid $hex'\"")
      ;;
    *)
      width="$(tmux show -gqv @sidebar_width)"
      items=(-T "#[bold] Sidebar "
        "  + New Window"  n "run-shell -b '$SELF new-item $5 window'"
        "  + New Session" N "run-shell -b '$SELF new-item $5 session'"
        "  ⊕ Add Remote Host…" h "command-prompt -p 'SSH host:' \"run-shell -b '$SELF radd %%'\""
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

# 按字符数补齐/截断到 W 列（不处理双宽字符）
fit() {
  local s="$1"
  [ ${#s} -gt "$W" ] && s="${s:0:$W}"
  printf '%s%*s' "$s" $((W - ${#s})) ''
}

# 字面文本放进 tmux 格式：# , } 要转义
fmt_lit() {
  local t="${1//#/##}" c='}' r='#}'
  t="${t//,/#,}"
  printf '%s' "${t//$c/$r}"
}

# 左边 $1 截断/补齐，右边 $2 贴右，共 W 列
fitr() {
  local l="$1" r="$2" n=$((W - ${#2}))
  [ ${#l} -gt "$n" ] && l="${l:0:$((n - 1))}…"
  printf '%s%*s%s' "$l" $((n - ${#l})) '' "$r"
}

# 输出整帧。返回 0 = 画好且看得见，3 = 画好但在后台（render 据此决定刷新间隔），
# 1 = 该退出，2 = 这次没拿到信息、跳过。
# 必须在子 shell 里调用：bash 3.2 的 <(...) 每次泄漏一个文件描述符，常驻进程里会耗尽（Too many open files）
# 后台窗口的侧边栏也保持最新：切过去时画面已经是对的，不会先露出旧画面/空白再跳变。
# 高亮按“侧边栏所在的窗口”而不是“session 当前窗口”，这样后台侧边栏的画面在切过去之前就已正确
draw() {
  local mypane mysid mywid npanes visible sync zoom host accent second dim tbg alert myremote
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
  IFS="$SEP" read -r sync zoom host accent second dim tbg alert myremote < <(
    tmux display -p -t "$TMUX_PANE" "#{synchronize-panes}$SEP#{window_zoomed_flag}$SEP$(fitw $((W - 8)) host_short)$SEP#{@popup_accent}$SEP#{@popup_dim}$SEP#{@theme_secondary}$SEP#{@theme_bg}$SEP#{@popup_alert}$SEP#{@remote}"
  )

  # session 行：" " + 图标 + " " + 名字（当前 session 用打开的文件夹，其他用合上的，Nerd Font 图标）；窗口行："   " + 序号（2 列）+ " " + 名字，共 W 列；有提醒时名字让出 2 列给标记
  local R=$'\e[0m' B=$'\e[1m' lines=() sid wid sname sname2 idx name name2 active bell prev="" hidx cnt=0
  local dimc="${dim:-#7982a9}"
  load_collapsed
  # 折叠的 session：标题行先占位，窗口数在下一个 session 开始 / 列表结束时填上
  fill_collapsed() { [ -n "$hidx" ] && lines[$hidx]="${lines[$hidx]} ▸$(printf %2d "$cnt")$R"; hidx=""; cnt=0; }
  while IFS="$SEP" read -r sid wid rem sname sname2 idx name name2 active bell; do
    if [ "$sid" != "$prev" ]; then
      fill_collapsed
      if is_collapsed "$sid"; then
        hidx=${#lines[@]}
        if [ "$sid" = "$mysid" ]; then
          lines+=("$B$(fg "$accent")  $sname2")
        else
          lines+=("$B$(fg "$dimc")  $sname2")
        fi
      elif [ "$sid" = "$mysid" ]; then
        lines+=("$B$(fg "$accent")  $sname$R")
      else
        lines+=("$B$(fg "$dimc")  $sname$R")
      fi
      prev="$sid"
    fi
    [ -n "$rem" ] && continue
    is_collapsed "$sid" && { cnt=$((cnt + 1)); continue; }
    if [ "$wid" = "$mywid" ] && [ "$sid" = "$mysid" ]; then
      lines+=("$B$(bg "$second")$(fg "$tbg")   $idx $name$R")
    elif [ "$bell" = 1 ]; then
      lines+=("$B$(fg "$alert")   $idx $name2 !$R")
    elif [ "$active" = 1 ] && [ "$sid" != "$mysid" ]; then
      lines+=("$(fg "$second")   $idx $name$R")
    else
      lines+=("$(fg "$dimc")   $idx $name$R")
    fi
  done < <(windows_tree "$(fitw $((W - 3)) session_name)$SEP$(fitw $((W - 7)) session_name)$SEP#{p2:window_index}$SEP$(fitw $((W - 6)) window_name)$SEP$(fitw $((W - 8)) window_name)$SEP#{window_active}$SEP#{window_bell_flag}")
  fill_collapsed
  lines+=("$(fg "$second")$(fit " +  new window")$R")

  # 远程区：分隔线、主机行（ 󰒋 host + 状态）、session 行（名字 + 窗口数）
  # 远程 session 名可能含中文等双宽字符：先占位，最后用一次 tmux 调用按显示宽度截断/补齐
  local rk rh rs rn ra rl rule rfmt="" ridx=() rpre=() rsuf=() rnames j
  while IFS="$SEP" read -r rk rh rs rn ra rl; do
    case "$rk" in
      sep)
        rule=" ─ remote "
        while [ ${#rule} -lt $((W - 1)) ]; do rule+="─"; done
        lines+=("$(fg "$dimc")$rule $R") ;;
      host)
        if [ -n "$rn" ]; then   # 折叠：rn = 隐藏的 session 数
          lines+=("$B$([ "$myremote" = "$rh:" ] && bg "$second" && fg "$tbg" || fg "$accent")$(fitr " 󰒋 $rh" "▸$(printf %2d "$rn") ")$R")
          continue
        fi
        case "$rs" in
          ok) lines+=("$B$([ "$myremote" = "$rh:" ] && bg "$second" && fg "$tbg" || fg "$accent")$(fit " 󰒋 $rh")$R") ;;
          down) lines+=("$B$(fg "$dimc")$(fitr " 󰒋 $rh" " offline")$R") ;;
          nosrv) lines+=("$B$(fg "$dimc")$(fitr " 󰒋 $rh" " no tmux")$R") ;;
          *) lines+=("$B$(fg "$dimc")$(fitr " 󰒋 $rh" " …")$R") ;;
        esac ;;
      rsession)
        if [ "$myremote" = "$rh:$rs" ]; then
          rpre+=("$B$(bg "$second")$(fg "$tbg") ● ")
        elif [ "$rl" = 1 ]; then
          rpre+=("$(fg "$second") ● ")
        else
          rpre+=("$(fg "$dimc")   ")
        fi
        rsuf+=(" ${rn}w$R")
        ridx+=(${#lines[@]})
        lines+=("")
        j=$((W - 3 - ${#rn} - 2))
        [ "$j" -lt 1 ] && j=1
        rfmt+="#{p$j:#{=/$j/…:#{l:$(fmt_lit "$rs")}}}$SEP" ;;
    esac
  done < <(remote_rows)
  if [ ${#ridx[@]} -gt 0 ]; then
    IFS="$SEP" read -r -a rnames < <(tmux display -p -t "$TMUX_PANE" "$rfmt")
    for j in "${!ridx[@]}"; do
      lines[${ridx[$j]}]="${rpre[$j]}${rnames[$j]}${rsuf[$j]}"
    done
  fi

  # 底部：SYNC / ZOOM 标记 + 主机名；列表过长时截断列表，保留底部
  local footer=() i
  [ "$sync" = 1 ] && footer+=("$B$(bg "$alert")$(fg '#ffffff') SYNC $R")
  [ "$zoom" = 1 ] && footer+=("$B$(bg '#2ecc71')$(fg '#ffffff') ZOOM $R")
  footer+=("$B$(fg "$accent") $host$R$(fg "$second")≡ menu $R")
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

# --- 入口 ---

# 入口放进函数，并与 exit 写在同一行：bash 边执行边读脚本，
# 运行中的进程（侧边栏常驻）在脚本被修改后不会接着读到改过的内容
main() {
  case "$1" in
    render) render ;;
    refresh) refresh ;;
    click) click "$2" "$3" "$4"; exit 0 ;;
    collapse) toggle_collapse "$2"; exit 0 ;;
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
      ensure_poller
      ;;
    poll-loop) poll_loop ;;
    poll-host) poll_host "$2"; refresh; exit 0 ;;
    radd) remote_add "$2"; exit 0 ;;
    rremove) remote_remove "$2"; exit 0 ;;
    rshell) remote_shell "$2" "$3"; exit 0 ;;
    rnew) [ -n "$4" ] && remote_new "$2" "$4" "$3"; exit 0 ;;
    rattach) remote_attach "$2" "$(hexdec "$3")" "$4"; exit 0 ;;
    rkill) remote_kill "$2" "$(hexdec "$3")"; exit 0 ;;
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
