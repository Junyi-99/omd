package main

import (
	"slices"
	"strings"
	"testing"
)

func TestRegistryInvariants(t *testing.T) {
	platforms := []string{"linux", "darwin"}
	reg := registry()
	ids := map[string]bool{}
	for _, p := range reg {
		if ids[p.ID] {
			t.Errorf("duplicate package ID %q", p.ID)
		}
		ids[p.ID] = true
	}
	for _, p := range reg {
		for _, n := range p.Needs {
			if !ids[n] {
				t.Errorf("%s needs unknown package %q", p.ID, n)
			}
		}
		for key := range p.Installers {
			if key != "all" && !slices.Contains(platforms, key) {
				t.Errorf("%s has installer for unknown platform %q", p.ID, key)
			}
		}
		for _, plat := range p.Platforms {
			if !slices.Contains(platforms, plat) {
				t.Errorf("%s lists unknown platform %q", p.ID, plat)
			}
		}
		for _, plat := range platforms {
			if !p.supports(plat) {
				continue
			}
			inst := p.installerFor(plat)
			if inst.Release == nil && len(inst.Commands) == 0 && p.ID != "kitty-config-link" {
				t.Errorf("%s supports %s but has no installer", p.ID, plat)
			}
		}
	}
}

func TestDefaultSelection(t *testing.T) {
	common := []string{"zsh-autosuggestions", "zsh-syntax-highlighting", "starship", "nodejs", "neovim",
		"tree-sitter-cli", "tmux", "fzf", "zoxide", "batcat", "duf", "uv", "claude-code", "codex-cli"}
	want := map[string][]string{
		"linux":  append([]string{"zsh"}, common...),
		"darwin": append([]string{"homebrew", "kitty-config-link"}, common...),
	}
	for plat, ids := range want {
		m := &manager{packages: packagesFor(plat)}
		var got []string
		for id, on := range m.defaultSelection() {
			if on {
				got = append(got, id)
			}
		}
		slices.Sort(got)
		w := slices.Clone(ids)
		slices.Sort(w)
		if !slices.Equal(got, w) {
			t.Errorf("%s defaults = %v, want %v", plat, got, w)
		}
	}
}

func TestEmbeddedScripts(t *testing.T) {
	for name, s := range map[string]string{
		"nvm_load": nvmScript, "nodejs_install": nodejsInstallScript, "npm_install": npmInstallScript,
		"tmux_check": tmuxCheckScript, "tmux_brew": tmuxBrew, "tmux_linux": tmuxLinux,
		"distro_install": distroInstallScript, "nvmLoad": nvmLoad,
	} {
		if strings.TrimSpace(s) == "" {
			t.Errorf("embedded script %s is empty", name)
		}
	}
}
