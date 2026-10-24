package main

import (
	"slices"
	"strings"
)

// Package describes selection, detection, installation, and shell integration.
// Commands use OMD_HOME; generated shell snippets use USER_CONFIG_HOME.
type Package struct {
	ID          string
	Name        string
	Description string
	Tab         string
	Category    string
	Presets     []string
	// Platforms limits the package to these GOOS values; empty means all.
	Platforms []string
	Check     string
	Env       string
	// Needs lists package IDs that install with this one.
	Needs      []string
	Installers map[string]Installer
	// Version, when set, lets omd compare an installed tool that is not an
	// omd-managed release against its latest version and upgrade it.
	Version *Version
}

// Version holds shell scripts: Installed and Latest print one version each,
// Upgrade moves the tool to Latest.
type Version struct {
	Installed, Latest, Upgrade string
}

// Installer pairs shell commands with the prerequisites they need, or names
// a GitHub release to unpack instead.
type Installer struct {
	Commands []string
	Requires []string
	Release  *Release
}

// nvmLoad makes nvm's node and npm visible; omd's own PATH knows nothing of nvm.
// NVM_DIR is exported so install.sh uses it instead of its XDG_CONFIG_HOME default.
const nvmLoad = `export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"; [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"; `

// npmInstaller installs a global npm package, failing clearly without Node.js.
func npmInstaller(pkg string, flags ...string) Installer {
	return Installer{Commands: []string{nvmLoad +
		`command -v npm >/dev/null || { echo "npm not found: install Node.js (nvm) first" >&2; exit 1; }; npm install -g ` + strings.Join(append(flags, pkg), " ")}}
}

// claudeScriptsFlag passes --allow-scripts only to an npm that has the option;
// older versions warn about the unknown config.
const claudeScriptsFlag = `$(npm config ls -l 2>/dev/null | grep -q '^allow-scripts ' && echo --allow-scripts=@anthropic-ai/claude-code)`

// tmuxVersion is the oldest tmux etc/tmux/tmux.conf fully works with (floating
// panes). Homebrew ships it; a distribution package may be older.
const tmuxVersion = "3.8"

// tmuxCheck accepts tmux at tmuxVersion or newer. Without Homebrew (Linux hosts
// using apt) any installed tmux is accepted, as nothing newer is offered.
// `tmux -V` prints "tmux 3.5a" or "tmux next-3.9"; only major.minor are compared.
var tmuxCheck = `command -v tmux >/dev/null && { ! command -v brew >/dev/null || tmux -V | awk -v min=` + tmuxVersion +
	` '{ sub(/^[^0-9]*/, "", $2); split($2, v, "."); split(min, m, ".");` +
	` exit !(v[1]+0 > m[1]+0 || (v[1]+0 == m[1]+0 && v[2]+0 >= m[2]+0)) }'; }`

// tmuxBrew upgrades an older Homebrew tmux instead of reporting it as already installed.
const tmuxBrew = `brew list --versions tmux >/dev/null 2>&1 && brew upgrade tmux || brew install tmux`

// tmuxLinux prefers Homebrew when present, otherwise the distribution package,
// whatever its version; tmux is never built from source.
var tmuxLinux = `if command -v brew >/dev/null; then
    ` + tmuxBrew + `
    exit 0
fi
[ "$(id -u)" = 0 ] && sudo= || sudo=sudo
if command -v apt-get >/dev/null; then
    $sudo apt-get update && $sudo apt-get install -y tmux
elif command -v dnf >/dev/null; then
    $sudo dnf install -y tmux
elif command -v yum >/dev/null; then
    $sudo yum install -y tmux
elif command -v apk >/dev/null; then
    $sudo apk add tmux
elif command -v pacman >/dev/null; then
    $sudo pacman -S --noconfirm tmux
elif command -v zypper >/dev/null; then
    $sudo zypper install -y tmux
else
    echo "no Homebrew or known package manager found; install tmux manually" >&2
    exit 1
fi
`

func brewInstaller(formula string) Installer {
	return Installer{
		Commands: []string{"brew install " + formula},
		Requires: []string{"brew"},
	}
}

// Release archives name the CPU differently per project. 32-bit arm uses the
// ARMv6 builds, which also run on ARMv7 and later Raspberry Pis.
var (
	goArch     = map[string]string{"amd64": "amd64", "arm64": "arm64", "arm": "armv6"}
	muslTarget = map[string]string{
		"amd64": "x86_64-unknown-linux-musl",
		"arm64": "aarch64-unknown-linux-musl",
		"arm":   "arm-unknown-linux-musleabihf",
	}
)

func (p Package) supports(platform string) bool {
	return len(p.Platforms) == 0 || slices.Contains(p.Platforms, platform)
}

func (p Package) installerFor(platform string) Installer {
	if installer, ok := p.Installers[platform]; ok {
		return installer
	}
	return p.Installers["all"]
}

// packagesFor returns the packages that can be installed on platform.
func packagesFor(platform string) []Package {
	var result []Package
	for _, p := range registry() {
		if p.supports(platform) {
			result = append(result, p)
		}
	}
	return result
}

func registry() []Package {
	allPresets := []string{"minimal", "agent", "full"}
	fullPreset := []string{"full"}
	agentPresets := []string{"agent", "full"}

	return []Package{
		{
			ID:          "zsh-autosuggestions",
			Name:        "zsh-autosuggestions",
			Description: "Suggestions from your shell history",
			Tab:         "Basics",
			Category:    "Shell",
			Presets:     allPresets,
			Check:       `test -f "$OMD_HOME/dep/zsh-autosuggestions/zsh-autosuggestions.zsh"`,
			Installers: map[string]Installer{
				"all": {
					Commands: []string{`git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions "$OMD_HOME/dep/zsh-autosuggestions"`},
					Requires: []string{"git"},
				},
			},
			Env: `source "$USER_CONFIG_HOME/dep/zsh-autosuggestions/zsh-autosuggestions.zsh"
ZSH_AUTOSUGGEST_STRATEGY=(match_prev_cmd completion)
ZSH_AUTOSUGGEST_CLEAR_WIDGETS+=(bracketed-paste)
bindkey '^B' autosuggest-toggle
`,
		},
		{
			ID:          "zsh-syntax-highlighting",
			Name:        "zsh-syntax-highlighting",
			Description: "Highlight commands as you type",
			Tab:         "Basics",
			Category:    "Shell",
			Presets:     allPresets,
			Check:       `test -f "$OMD_HOME/dep/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"`,
			Installers: map[string]Installer{
				"all": {
					Commands: []string{`git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting "$OMD_HOME/dep/zsh-syntax-highlighting"`},
					Requires: []string{"git"},
				},
			},
			Env: `source "$USER_CONFIG_HOME/dep/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
`,
		},
		{
			ID:          "starship",
			Name:        "Starship",
			Description: "A fast, context-aware prompt",
			Tab:         "Basics",
			Category:    "Shell",
			Presets:     fullPreset,
			Check:       "command -v starship",
			Installers: map[string]Installer{
				"darwin": brewInstaller("starship"),
				"linux": {Release: &Release{
					Repo:  "starship/starship",
					Asset: "starship-{arch}.tar.gz",
					Arch:  muslTarget,
					Bins:  []string{"starship"},
				}},
			},
			Env: `command -v starship >/dev/null && eval "$(starship init zsh)"
`,
		},
		{
			ID:          "homebrew",
			Name:        "Homebrew",
			Description: "Package manager for macOS; Linux tools install from GitHub releases",
			Tab:         "Basics",
			Category:    "Package managers",
			Presets:     fullPreset,
			Platforms:   []string{"darwin"},
			Check:       "command -v brew",
			Installers: map[string]Installer{
				"darwin": {
					Commands: []string{`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`},
					Requires: []string{"curl", "git"},
				},
			},
			Env: `for _omd_brew in /opt/homebrew/bin/brew /usr/local/bin/brew; do
    if [[ -x "$_omd_brew" ]]; then
        eval "$("$_omd_brew" shellenv)"
        break
    fi
done
unset _omd_brew
`,
		},
		{
			ID:          "pyenv",
			Name:        "pyenv",
			Description: "Manage Python versions",
			Tab:         "Basics",
			Category:    "Package managers",
			Presets:     fullPreset,
			Check:       "command -v pyenv",
			Installers: map[string]Installer{
				"darwin": brewInstaller("pyenv"),
				"linux": {
					Commands: []string{"curl -fsSL https://pyenv.run | bash"},
					Requires: []string{"curl", "git"},
				},
			},
			Env: `export PYENV_ROOT="$HOME/.pyenv"
[[ -d "$PYENV_ROOT/bin" ]] && path=("$PYENV_ROOT/bin" $path)
command -v pyenv >/dev/null && eval "$(pyenv init - zsh)"
`,
		},
		{
			ID:          "nodejs",
			Name:        "Node.js (nvm)",
			Description: "Node.js 24 LTS managed by nvm; required by the CLI agents",
			Tab:         "Basics",
			Category:    "Package managers",
			Presets:     agentPresets,
			Check:       nvmLoad + "nvm which 24",
			Version: &Version{
				Installed: nvmLoad + "nvm version 24",
				Latest:    nvmLoad + "nvm version-remote 24",
				Upgrade:   nvmLoad + `old=$(nvm version 24); nvm install "$(nvm version-remote 24)" --reinstall-packages-from="$old"`,
			},
			Installers: map[string]Installer{
				"all": {
					// PROFILE=/dev/null: the Env snippet below loads nvm, not ~/.zshrc edits.
					Commands: []string{nvmLoad + `mkdir -p "$NVM_DIR" && curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.8/install.sh | PROFILE=/dev/null bash && . "$NVM_DIR/nvm.sh" && nvm install 24`},
					Requires: []string{"curl", "git"},
				},
			},
			// Sourcing nvm.sh costs ~0.3s even with --no-use, so load it on the
			// first nvm call and put the default version's bin on PATH directly.
			// The glob is zsh-only.
			Env: `export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
if [[ -s "$NVM_DIR/nvm.sh" ]]; then
    nvm() { unset -f nvm; source "$NVM_DIR/nvm.sh" --no-use; nvm "$@"; }
    _nvm_default=$(<"$NVM_DIR/alias/default" 2>/dev/null)
    _nvm_nodes=("$NVM_DIR"/versions/node/v${_nvm_default#v}*(N-/))
    (( $#_nvm_nodes )) && path=("${${(@On)_nvm_nodes}[1]}/bin" $path)
    unset _nvm_default _nvm_nodes
fi
`,
		},
		{
			ID:          "neovim",
			Name:        "Neovim",
			Description: "Terminal editor; selection also aliases vim to nvim",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v nvim",
			// nvim-treesitter (main) builds parsers with tree-sitter-cli.
			Needs:       []string{"nodejs", "tree-sitter-cli"},
			Installers: map[string]Installer{
				"darwin": brewInstaller("neovim"),
				"linux": {Release: &Release{
					Repo:  "neovim/neovim",
					Asset: "nvim-linux-{arch}.tar.gz",
					Arch:  map[string]string{"amd64": "x86_64", "arm64": "arm64"},
					Bins:  []string{"bin/nvim"},
				}},
			},
			Env: `alias vim="nvim"
`,
		},
		{
			ID:          "tree-sitter-cli",
			Name:        "tree-sitter CLI",
			Description: "Builds Neovim's treesitter parsers; installed with npm",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       nvmLoad + "command -v tree-sitter",
			Needs:       []string{"nodejs"},
			Installers: map[string]Installer{
				// The install script downloads the prebuilt binary.
				"all": npmInstaller("tree-sitter-cli", "--allow-scripts=tree-sitter-cli"),
			},
		},
		{
			ID:          "tmux",
			Name:        "tmux",
			Description: "Terminal multiplexer, " + tmuxVersion + "+ via Homebrew; Linux without Homebrew uses apt and accepts an older version",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       tmuxCheck,
			Installers: map[string]Installer{
				"darwin": {Commands: []string{tmuxBrew}, Requires: []string{"brew"}},
				"linux":  {Commands: []string{tmuxLinux}},
			},
		},
		{
			ID:          "fzf",
			Name:        "fzf",
			Description: "Interactive fuzzy finder",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v fzf",
			Installers: map[string]Installer{
				"darwin": brewInstaller("fzf"),
				"linux": {Release: &Release{
					Repo:  "junegunn/fzf",
					Asset: "fzf-{version}-linux_{arch}.tar.gz",
					Arch:  goArch,
					Bins:  []string{"fzf"},
				}},
			},
		},
		{
			ID:          "zoxide",
			Name:        "zoxide",
			Description: "Jump to frequently used directories",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v zoxide",
			Installers: map[string]Installer{
				"darwin": brewInstaller("zoxide"),
				"linux": {Release: &Release{
					Repo:  "ajeetdsouza/zoxide",
					Asset: "zoxide-{version}-{arch}.tar.gz",
					Arch:  muslTarget,
					Bins:  []string{"zoxide"},
				}},
			},
			Env: `command -v zoxide >/dev/null && eval "$(zoxide init zsh)"
`,
		},
		{
			ID:          "batcat",
			Name:        "bat",
			Description: "Read files with syntax highlighting",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v bat || command -v batcat",
			Installers: map[string]Installer{
				"darwin": brewInstaller("bat"),
				"linux": {Release: &Release{
					Repo:  "sharkdp/bat",
					Asset: "bat-{tag}-{arch}.tar.gz",
					Arch:  muslTarget,
					Bins:  []string{"bat"},
				}},
			},
			// Debian and Ubuntu ship the binary as batcat.
			Env: `if command -v bat >/dev/null; then
    alias cat="bat -p"
elif command -v batcat >/dev/null; then
    alias cat="batcat -p"
fi
`,
		},
		{
			ID:          "duf",
			Name:        "duf",
			Description: "Disk usage at a glance",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v duf",
			Installers: map[string]Installer{
				"darwin": brewInstaller("duf"),
				"linux": {Release: &Release{
					Repo:  "muesli/duf",
					Asset: "duf_{version}_linux_{arch}.tar.gz",
					Arch:  map[string]string{"amd64": "x86_64", "arm64": "arm64", "arm": "armv6"},
					Bins:  []string{"duf"},
				}},
			},
		},
		{
			ID:          "syncthing",
			Name:        "Syncthing",
			Description: "Continuous peer-to-peer file synchronization",
			Tab:         "Basics",
			Category:    "Tools",
			Presets:     fullPreset,
			Check:       "command -v syncthing",
			Installers: map[string]Installer{
				"darwin": brewInstaller("syncthing"),
				"linux": {Release: &Release{
					Repo:  "syncthing/syncthing",
					Asset: "syncthing-linux-{arch}-{tag}.tar.gz",
					Arch:  map[string]string{"amd64": "amd64", "arm64": "arm64", "arm": "arm"},
					Bins:  []string{"syncthing"},
				}},
			},
		},
		{
			ID:          "kitty-config-link",
			Name:        "Kitty config symlink",
			Description: "Link etc/kitty to ~/.config/kitty; backs up an existing config first",
			Tab:         "Basics",
			Category:    "Configuration",
			Check:       `test -L "$HOME/.config/kitty" && test "$(readlink "$HOME/.config/kitty")" = "$OMD_HOME/etc/kitty"`,
		},
		{
			ID:          "claude-code",
			Name:        "Claude Code",
			Description: "Anthropic's command-line assistant",
			Tab:         "CLI Agents",
			Category:    "Agents",
			Presets:     agentPresets,
			Check:       nvmLoad + "command -v claude",
			Needs:       []string{"nodejs"},
			Installers: map[string]Installer{
				"all": npmInstaller("@anthropic-ai/claude-code", claudeScriptsFlag),
			},
		},
		{
			ID:          "codex-cli",
			Name:        "Codex CLI",
			Description: "OpenAI's command-line assistant",
			Tab:         "CLI Agents",
			Category:    "Agents",
			Presets:     agentPresets,
			Check:       nvmLoad + "command -v codex",
			Needs:       []string{"nodejs"},
			Installers: map[string]Installer{
				"all": npmInstaller("@openai/codex"),
			},
		},
		{
			ID:          "apptainer",
			Name:        "Apptainer",
			Description: "Linux container runtime; install it with your distribution's package manager",
			Tab:         "Basics",
			Category:    "Tools",
			Platforms:   []string{"linux"},
			Check:       "command -v apptainer",
			Installers: map[string]Installer{
				"linux": {
					Commands: []string{`echo 'Install Apptainer with your package manager: https://apptainer.org/docs/admin/main/installation.html' >&2; exit 1`},
				},
			},
		},
	}
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
