#!/bin/bash

function log_message() {
    local message="$1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $message"
}

# 解压文件的函数
extract() {
    local file="$1"
    if [[ -f $file ]]; then
        case $file in
        *.tar.bz2) tar xjvf "$file" ;;
        *.tar.gz) tar xzvf "$file" ;;
        *.bz2) bunzip2 "$file" ;;
        *.rar) rar x "$file" ;;
        *.gz) gunzip "$file" ;;
        *.tar) tar xvf "$file" ;;
        *.tbz2) tar xjvf "$file" ;;
        *.tgz) tar xzvf "$file" ;;
        *.zip) unzip "$file" ;;
        *.Z) uncompress "$file" ;;
        *) log_message "'$file' cannot be extracted via extract()" ;;
        esac
    else
        log_message "'$file' is not a valid file"
    fi
}

# Show a random tip from etc/fortunes once a day. The file uses fortune's
# %-separated format, but is read directly so fortune/strfile are not needed.
function print_greetings() {
    local tips="${XDG_CONFIG_HOME}/fortunes/ubuntu-server-tips"
    local stamp
    stamp="${XDG_CACHE_HOME}/.greeting-$(date +%Y%m%d)"
    if [ ! -f "$tips" ] || ! command -v cowsay &>/dev/null || [ -f "$stamp" ]; then
        return 0
    fi
    find "${XDG_CACHE_HOME}" -maxdepth 1 -name '.greeting-*' -delete 2>/dev/null
    touch "$stamp"

    local tip
    tip=$(awk -v RS='%' -v seed="$RANDOM" '
        { gsub(/^\n+|\n+$/, ""); if ($0 != "") tips[++n] = $0 }
        END { srand(seed); if (n) print tips[int(rand() * n) + 1] }' "$tips")
    if command -v lolcat &>/dev/null; then
        cowsay -t "$tip" | lolcat
    else
        cowsay -t "$tip"
    fi
}
