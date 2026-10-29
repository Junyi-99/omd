# Install a package with the distribution's package manager.
# Inputs: OMD_APT_PKGS (apt name), OMD_PKGS (name elsewhere).
[ "$(id -u)" = 0 ] && sudo= || sudo=sudo
if command -v apt-get >/dev/null; then $sudo apt-get update && $sudo apt-get install -y $OMD_APT_PKGS
elif command -v dnf >/dev/null; then $sudo dnf install -y $OMD_PKGS
elif command -v yum >/dev/null; then $sudo yum install -y $OMD_PKGS
elif command -v apk >/dev/null; then $sudo apk add $OMD_PKGS
elif command -v pacman >/dev/null; then $sudo pacman -S --noconfirm --needed $OMD_PKGS
elif command -v zypper >/dev/null; then $sudo zypper install -y $OMD_PKGS
else echo "no known package manager; install $OMD_PKGS manually" >&2; exit 1
fi
