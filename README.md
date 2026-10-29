# Dotfiles

> [!TIP]
> There isn’t a one-size-fits-all configuration file. You should explore dotfiles shared by others and continuously build your own personalized dotfile.

<img width="1005" alt="image" src="https://github.com/user-attachments/assets/36147a29-6c22-40e1-85b4-875cfa773b63">

<img width="1098" alt="image" src="https://github.com/user-attachments/assets/30b98735-9bf8-475d-acff-6f7716a3c16f">

## Features
- **out-of-the-box**
- **minimal side effects**
  - most operations happen in `~/.oh-my-dotfiles`
  - `XDG_CONFIG_HOME` is set to `~/.oh-my-dotfiles/etc/`, so XDG-aware tools (tmux, neovim, etc.) read configs directly from the repo
  - third-party binaries will be downloaded into `~/.oh-my-dotfiles/bin` or `~/.oh-my-dotfiles/deps` (See [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/index.html))
- neovim
  - AI assistant ([copilot.lua](https://github.com/zbirenbaum/copilot.lua))
  - lsp + lsp-manager ([mason](https://github.com/williamboman/mason.nvim))
  - plugin-manager ([lazy-nvim](https://github.com/folke/lazy.nvim))
  - notification-manager ([nvim-notify](https://github.com/rcarriga/nvim-notify))
  - code highlighting ([nvim-treesitter](https://github.com/nvim-treesitter/nvim-treesitter))
  - shortcut key hint ([which-key](https://github.com/folke/which-key.nvim))
- [batcat](https://github.com/sharkdp/bat)
- [tmux](https://github.com/tmux/tmux)
  - True color (256color + RGB)
  - <kbd>CTRL + A</kbd> Prefix
  - Mouse enabled
  - Vim-style pane navigation (<kbd>h</kbd><kbd>j</kbd><kbd>k</kbd><kbd>l</kbd>) and resizing (<kbd>H</kbd><kbd>J</kbd><kbd>K</kbd><kbd>L</kbd>)
  - Intuitive split keys: <kbd>|</kbd> horizontal, <kbd>-</kbd> vertical
  - <kbd>Alt + Arrow</kbd> switch panes, <kbd>Shift + Arrow</kbd> switch windows (no prefix needed)
  - Right-click context menus for pane/window/session management
  - Pane sync toggle (<kbd>S</kbd>), SYNC/ZOOM status indicators
  - Vi copy mode with <kbd>C-v</kbd> rectangle selection
  - Auto dark/light theme following macOS appearance (Tokyo Night)
  - Cheatsheet on new session
- [fzf](https://github.com/junegunn/fzf)
- [zoxide](https://github.com/ajeetdsouza/zoxide)
- [zsh](https://www.zsh.org/) + [starship](https://starship.rs/) + [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) + [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting)

## Get Started

Paste this one line. It clones the repo, wires up the shell and SSH entry points
(keeping existing configuration), opens the `omd` package selector, then starts zsh.
[Go](https://go.dev/dl/) is only needed when no prebuilt `omd` matches your
platform and revision (see below).

```sh
[ -d ~/.oh-my-dotfiles ] || git clone https://github.com/Junyi-99/omd.git ~/.oh-my-dotfiles; for f in .zshenv:zshenv .zshrc:zshrc; do touch ~/${f%:*}; grep -qxF "source ~/.oh-my-dotfiles/etc/zsh/${f#*:}" ~/${f%:*} || echo "source ~/.oh-my-dotfiles/etc/zsh/${f#*:}" >> ~/${f%:*}; done; mkdir -p ~/.ssh; touch ~/.ssh/config; grep -qxF 'Include ~/.oh-my-dotfiles/etc/ssh/*.local' ~/.ssh/config || echo 'Include ~/.oh-my-dotfiles/etc/ssh/*.local' >> ~/.ssh/config; ~/.oh-my-dotfiles/bin/omd && exec zsh
```

Defaults are preselected for your OS on the first run; press Enter in the TUI to install.

## omd

The native TUI uses Go 1.27.1, Bubble Tea v2.0.10, and Lip Gloss v2.0.6.
`cmd/omd/go.mod` and `cmd/omd/go.sum` pin the toolchain requirement and dependency versions.
The launcher `bin/omd` caches the binary in `cache/omd/omd`, keyed on the git
tree of `cmd/omd`. When the cache is stale it downloads the prebuilt binary for
Linux (x86_64, arm64, arm) or macOS (x86_64, arm64) from the rolling `omd`
GitHub Release, which CI publishes for every change to `cmd/omd` together with
a `SHA256SUMS` file; the download is accepted only if that file names the same
`cmd/omd` tree and the checksum matches. Otherwise, and whenever `cmd/omd` has
uncommitted edits, it builds from source with Go (toolchain and modules are
fetched automatically into `cache/`). Shell startup loads `etc/omd/env.sh` and
starts a background update check at the configured interval.

```sh
omd                      # select packages in the TUI
omd install              # install saved selections; upgrade outdated release tools and Node.js
omd outdated             # list tools (GitHub releases, Node.js) and their available updates
omd install --quiet      # log normal output; still report failures
omd env                  # regenerate shell integration without installing
omd update               # update the oh-my-dotfiles repository, then list outdated tools
omd auto-update          # run the scheduled update check
omd version              # show the actual compiled versions
```

| Key | Action |
| --- | --- |
| ↑/↓ or j/k | Move; the list scrolls with the cursor |
| Tab or 1/2 | Switch between Basics and CLI Agents |
| Space / a | Select or deselect one package / the current group |
| m / g / f | Apply minimal / agent / full preset |
| d | Reset to the first-run default selection |
| s | Save selections and regenerate shell integration |
| Enter | Install selected tools and return to the shell with OMD environment loaded |
| r | Refresh installed status |
| u | Cycle auto-update frequency: off / daily / weekly / monthly |
| q / Esc | Quit and discard unsaved selections |

CLI Agents contains Claude Code and Codex CLI. Disabling a package removes its
shell integration, not its installed files.
Each package shows availability (`missing` / `installed`) and its selection
state. `Selected` means Enter will install it if missing; packages with shell
integration also enable that integration. `Unselected` means OMD skips its
installation and omits its integration. Neither state uninstalls existing files.
Select `Kitty config symlink` and press Enter to link `etc/kitty` to
`~/.config/kitty`. An existing config directory is moved to a timestamped
`.omd-backup-*` path first; an existing symlink to another location is preserved
and reported as a conflict.
Selecting Neovim also adds the `vim` → `nvim` shell alias when Neovim is installed,
and pulls in Node.js and the tree-sitter CLI (from npm), which nvim-treesitter
needs to build its parsers.
`tmux` needs 3.8 or newer (`etc/tmux/tmux.conf` uses floating panes); an older
one shows as `missing` and is upgraded. macOS and Linux with Homebrew use
Homebrew. Otherwise Linux builds the latest GitHub release from source into
`~/.local`, installing the compiler and libevent/ncurses headers with apt (or
dnf, yum, apk, pacman, zypper).

On the first run (no saved selection), omd preselects the default
selection (packages marked `Default` in `packages.go`, independent of the m/g/f presets) for the
current OS: on Linux zsh, the zsh plugins, Starship, Node.js, Neovim, the
tree-sitter CLI, tmux, fzf, zoxide, bat, duf, uv, Claude Code and Codex CLI;
on macOS the same plus Homebrew and the Kitty config symlink. Syncthing and
Apptainer are never preselected.

On macOS, CLI tools install with Homebrew. Linux does not use Homebrew:
fzf, zoxide, bat, duf, uv, Starship, and Neovim are unpacked from their GitHub
releases into `dep/<package>` and linked into `bin/`. The archive matches the
CPU omd was built for (x86_64, arm64, or 32-bit arm, e.g. Raspberry Pi), and
the installed tag is recorded in `dep/<package>/.omd-version`. The TUI marks
outdated tools with `update <tag>`; `omd outdated` lists them, and Enter or
`omd install` upgrades the selected ones. Tools from apt or another package
manager are detected as installed and left to that package manager.

Installers keep their interactive terminal prompts and log to
`cache/omd/install.log`. Failures return a nonzero status, including in quiet
mode. Packages install after their dependencies (for example Homebrew before the tools
that need it, Node.js before the npm tools); missing prerequisites are reported
explicitly. Scheduled updates quietly fast-forward the oh-my-dotfiles repository.
They refuse unpushed commits and leave uncommitted changes alone (git aborts only if they collide with incoming files); use `omd install` to install
enabled tools.
The default frequency is weekly; change it with `u` in the TUI and save with `s`.

`omd update` refuses unpushed commits; uncommitted changes are kept unless they collide with incoming files. Use `--force`
only to deliberately discard tracked local changes and local commits.

### Development

The code follows [Effective Go](https://go.dev/doc/effective_go): named fields,
explicit errors, standard formatting, and small functions with concrete types.

```sh
go fmt ./cmd/omd
go vet ./cmd/omd
go build -o cache/omd/omd ./cmd/omd
```
