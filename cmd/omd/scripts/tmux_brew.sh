# Upgrade an older Homebrew tmux instead of reporting it as already installed.
brew list --versions tmux >/dev/null 2>&1 && brew upgrade tmux || brew install tmux
