# Prefer Homebrew; otherwise build the latest GitHub release into ~/.local.
if command -v brew >/dev/null; then
    brew list --versions tmux >/dev/null 2>&1 && brew upgrade tmux || brew install tmux
    exit 0
fi
[ "$(id -u)" = 0 ] && sudo= || sudo=sudo
if command -v apt-get >/dev/null; then
    $sudo apt-get update && $sudo apt-get install -y build-essential bison pkg-config libevent-dev libncurses-dev curl
elif command -v dnf >/dev/null; then
    $sudo dnf install -y gcc make bison pkgconf-pkg-config libevent-devel ncurses-devel curl
elif command -v yum >/dev/null; then
    $sudo yum install -y gcc make bison pkgconfig libevent-devel ncurses-devel curl
elif command -v apk >/dev/null; then
    $sudo apk add build-base bison pkgconf libevent-dev ncurses-dev curl
elif command -v pacman >/dev/null; then
    $sudo pacman -S --noconfirm --needed base-devel libevent ncurses curl
elif command -v zypper >/dev/null; then
    $sudo zypper install -y gcc make bison pkg-config libevent-devel ncurses-devel curl
else
    echo "no Homebrew or known package manager; install tmux build dependencies manually" >&2
    exit 1
fi || { echo "installing tmux build dependencies failed" >&2; exit 1; }
tag=$(curl -fsSLI -o /dev/null -w '%{url_effective}' https://github.com/tmux/tmux/releases/latest | sed 's|.*/tag/||')
[ -n "$tag" ] || { echo "cannot find the latest tmux release" >&2; exit 1; }
tmp=$(mktemp -d) && trap 'rm -rf "$tmp"' EXIT
curl -fsSL "https://github.com/tmux/tmux/releases/download/$tag/tmux-$tag.tar.gz" | tar -xz -C "$tmp" &&
    cd "$tmp/tmux-$tag" &&
    ./configure --prefix="$HOME/.local" && make -j"$(nproc)" && make install
