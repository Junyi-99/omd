# Prefer Homebrew; otherwise install the prebuilt release from tmux/tmux-builds
# into ~/.local/bin (no root, compiler or libevent/ncurses needed).
if command -v brew >/dev/null; then
    brew list --versions tmux >/dev/null 2>&1 && brew upgrade tmux || brew install tmux
    exit 0
fi
case "$(uname -m)" in
    x86_64|amd64) arch=x86_64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) echo "no prebuilt tmux for $(uname -m)" >&2; exit 1 ;;
esac
tag=$(curl -fsSLI -o /dev/null -w '%{url_effective}' https://github.com/tmux/tmux-builds/releases/latest | sed 's|.*/tag/||')
[ -n "$tag" ] || { echo "cannot find the latest tmux release" >&2; exit 1; }
tmp=$(mktemp -d) && trap 'rm -rf "$tmp"' EXIT
curl -fsSL "https://github.com/tmux/tmux-builds/releases/download/$tag/tmux-${tag#v}-linux-$arch.tar.gz" | tar -xz -C "$tmp" &&
    mkdir -p "$HOME/.local/bin" &&
    install -m 755 "$tmp/tmux" "$HOME/.local/bin/tmux"
