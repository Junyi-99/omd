package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func testManager(t *testing.T) *manager {
	t.Helper()
	return &manager{
		root:     filepath.Join(t.TempDir(), "dotfiles with ' quotes"),
		home:     t.TempDir(),
		platform: "linux",
		arch:     "arm64",
		github:   "http://github.invalid",
		packages: packagesFor("linux"),
	}
}

func TestConfigMigration(t *testing.T) {
	m := testManager(t)
	legacy := "# Existing Python configuration\nstarship=enabled\ncopilot-cli=enabled\ngemini-cli=disabled\ncursor-cli=enabled\nkimi-code=disabled\n"
	if err := atomicWrite(m.configPath(), legacy); err != nil {
		t.Fatal(err)
	}
	selected, err := m.loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(selected) != 1 || !selected["starship"] {
		t.Fatalf("migration = %v", selected)
	}
	if err := m.saveConfig(selected); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(m.configPath())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(saved), "copilot") || strings.Contains(string(saved), "gemini") {
		t.Fatal("removed packages survived save")
	}
	if err := atomicWrite(m.configPath(), "starship=enabeld\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.loadConfig(); err == nil {
		t.Fatal("invalid configuration accepted")
	}
}

func TestInstallFailures(t *testing.T) {
	for _, tc := range []struct {
		name      string
		installer Installer
	}{
		{name: "failed command", installer: Installer{Commands: []string{"echo failure detail >&2; exit 7"}}},
		{name: "false success", installer: Installer{Commands: []string{"true"}}},
		{name: "no installer"},
		{name: "missing prerequisite", installer: Installer{Commands: []string{"true"}, Requires: []string{"omd-nonexistent-prerequisite"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := testManager(t)
			m.packages = []Package{{
				ID: "probe", Name: "Probe", Check: "false", Env: "unsafe snippet\n",
				Installers: map[string]Installer{"all": tc.installer},
			}}
			var output bytes.Buffer
			err := m.install(map[string]bool{"probe": true}, true, &output)
			if err == nil || !strings.Contains(err.Error(), "Probe") {
				t.Fatalf("install error = %v", err)
			}
			if output.Len() != 0 {
				t.Fatalf("quiet mode printed success output: %s", &output)
			}
			env, err := os.ReadFile(filepath.Join(m.root, "etc", "omd", "env.sh"))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(env), "unsafe snippet") {
				t.Fatal("failed package leaked into generated environment")
			}
		})
	}
}

func TestInstallOrderAndIdempotence(t *testing.T) {
	m := testManager(t)
	m.packages = []Package{
		{
			ID: "dependent", Name: "Dependent", Check: `test -f "$OMD_HOME/dependent"`,
			Installers: map[string]Installer{"all": {Commands: []string{`test -f "$OMD_HOME/brew"; touch "$OMD_HOME/dependent"`}}},
		},
		{
			ID: "homebrew", Name: "Homebrew", Check: `test -f "$OMD_HOME/brew"`,
			Installers: map[string]Installer{"all": {Commands: []string{`echo installed >> "$OMD_HOME/brew"`}}},
		},
	}
	selected := map[string]bool{"dependent": true, "homebrew": true}
	for range 2 {
		if err := m.install(selected, true, io.Discard); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(filepath.Join(m.root, "brew"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "installed\n" {
		t.Fatalf("installer ran more than once: %q", data)
	}
}

func TestEnvironmentOrder(t *testing.T) {
	m := testManager(t)
	m.packages = []Package{
		{ID: "zsh-syntax-highlighting", Check: "true", Env: "# highlighting-last\n"},
		{ID: "starship", Check: "true", Env: "# prompt\n"},
		{ID: "homebrew", Check: "true", Env: "# brew-first\n"},
	}
	selected := map[string]bool{"zsh-syntax-highlighting": true, "starship": true, "homebrew": true}
	if err := m.generateEnv(selected); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(m.root, "etc", "omd", "env.sh"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if strings.Index(text, "# brew-first") > strings.Index(text, "# prompt") || !strings.HasSuffix(text, "# highlighting-last\n") {
		t.Fatalf("incorrect environment order:\n%s", text)
	}
}

func TestKittyLinkPreservesDirectory(t *testing.T) {
	m := testManager(t)
	target := filepath.Join(m.root, "etc", "kitty")
	link := filepath.Join(m.home, ".config", "kitty")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(link, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(link, "kitty.conf"), []byte("keep me"), 0600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := m.linkConfigs(&output); err != nil {
		t.Fatal(err)
	}
	// Resolve both sides: macOS temp dirs live under the /var → /private/var symlink.
	wanted, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(link)
	if err != nil || resolved != wanted {
		t.Fatalf("kitty link = %q, %v", resolved, err)
	}
	backups, err := filepath.Glob(link + ".omd-backup-*")
	if err != nil || len(backups) != 1 {
		t.Fatalf("backup paths = %v, %v", backups, err)
	}
	if contents, err := os.ReadFile(filepath.Join(backups[0], "kitty.conf")); err != nil || string(contents) != "keep me" {
		t.Fatalf("existing config was not preserved in backup: %q, %v", contents, err)
	}
	if err := m.linkConfigs(io.Discard); err != nil {
		t.Fatal(err)
	}
}

func TestTUIResizeNavigationAndSave(t *testing.T) {
	m := newModel(testManager(t), make(map[string]bool))
	for _, size := range []tea.WindowSizeMsg{
		{Width: 100, Height: 30}, {Width: 80, Height: 24},
		{Width: 48, Height: 16}, {Width: 20, Height: 4},
	} {
		next, _ := m.Update(size)
		m = next.(model)
		next, _ = m.updateKey("end")
		m = next.(model)
		view := ansi.Strip(m.View().Content)
		if len(strings.Split(view, "\n")) > size.Height {
			t.Fatalf("view exceeds height %d", size.Height)
		}
		for _, line := range strings.Split(view, "\n") {
			if ansi.StringWidth(line) > size.Width {
				t.Fatalf("line exceeds width %d: %q", size.Width, line)
			}
		}
		if size.Width >= 48 && (!strings.Contains(view, "Apptainer") || !strings.Contains(view, "›")) {
			t.Fatal("focused package is outside the viewport")
		}
	}
	next, _ := m.updateKey("space")
	m = next.(model)
	if !m.selected["apptainer"] {
		t.Fatal("space did not toggle the focused package")
	}
	// Save an empty selection without installing; it must still write a config.
	m.selected = make(map[string]bool)
	msg := m.save(false)().(savedMsg)
	if msg.err != nil || msg.install {
		t.Fatalf("save-only result = %+v", msg)
	}
	if _, err := m.manager.loadConfig(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(m.manager.root, "cache", "omd", "install.log")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("save-only unexpectedly invoked the installer")
	}
	var agents []string
	for _, p := range registry() {
		if p.Tab == "CLI Agents" {
			agents = append(agents, p.ID)
		}
	}
	if !slices.Equal(agents, []string{"claude-code", "codex-cli"}) {
		t.Fatalf("agent catalog = %v", agents)
	}
}

func TestUpdateRefusesDirtyTree(t *testing.T) {
	m := testManager(t)
	if err := os.MkdirAll(m.root, 0755); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", m.root}, args...)...)
		if data, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, data)
		}
	}
	git("init", "-b", "main")
	git("-c", "commit.gpgsign=false", "-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "--allow-empty", "-m", "initial")
	git("branch", "upstream")
	git("branch", "--set-upstream-to=upstream")
	path := filepath.Join(m.root, "local-work")
	if err := os.WriteFile(path, []byte("keep me"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := m.update(false); err == nil || !strings.Contains(err.Error(), "uncommitted") {
		t.Fatalf("update error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("local work was removed")
	}
}

func TestReleaseAssetsPerArchitecture(t *testing.T) {
	want := map[string]map[string]string{
		"fzf": {
			"amd64": "fzf-1.2.3-linux_amd64.tar.gz",
			"arm64": "fzf-1.2.3-linux_arm64.tar.gz",
			"arm":   "fzf-1.2.3-linux_armv6.tar.gz",
		},
		"zoxide": {
			"amd64": "zoxide-1.2.3-x86_64-unknown-linux-musl.tar.gz",
			"arm64": "zoxide-1.2.3-aarch64-unknown-linux-musl.tar.gz",
			"arm":   "zoxide-1.2.3-arm-unknown-linux-musleabihf.tar.gz",
		},
		"batcat": {
			"amd64": "bat-v1.2.3-x86_64-unknown-linux-musl.tar.gz",
			"arm64": "bat-v1.2.3-aarch64-unknown-linux-musl.tar.gz",
		},
		"duf": {
			"amd64": "duf_1.2.3_linux_x86_64.tar.gz",
			"arm64": "duf_1.2.3_linux_arm64.tar.gz",
		},
		"starship": {
			"amd64": "starship-x86_64-unknown-linux-musl.tar.gz",
			"arm64": "starship-aarch64-unknown-linux-musl.tar.gz",
		},
		"neovim": {
			"amd64": "nvim-linux-x86_64.tar.gz",
			"arm64": "nvim-linux-arm64.tar.gz",
		},
	}
	for _, p := range packagesFor("linux") {
		rel := p.installerFor("linux").Release
		if rel == nil {
			if strings.Contains(strings.Join(p.installerFor("linux").Commands, " "), "brew") {
				t.Errorf("%s uses Homebrew on Linux", p.ID)
			}
			continue
		}
		for arch, asset := range want[p.ID] {
			got, err := rel.assetName("v1.2.3", arch)
			if err != nil || got != asset {
				t.Errorf("%s/%s asset = %q, %v; want %q", p.ID, arch, got, err, asset)
			}
		}
		if _, err := rel.assetName("v1.2.3", "riscv64"); !errors.Is(err, errUnsupportedArch) {
			t.Errorf("%s accepted an unknown architecture: %v", p.ID, err)
		}
	}
	for _, p := range packagesFor("linux") {
		if p.ID == "homebrew" {
			t.Fatal("Homebrew is offered on Linux")
		}
	}
}

func TestReleaseInstallAndUpgrade(t *testing.T) {
	latest := "v1.0.0"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/owner/tool/releases/latest":
			http.Redirect(w, r, "/owner/tool/releases/tag/"+latest, http.StatusFound)
		case "/owner/tool/releases/download/" + latest + "/tool-" + latest + "-aarch64.tar.gz":
			w.Write(tarball(t, "tool-dir/bin/tool", "#!/bin/sh\necho "+latest+"\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	m := testManager(t)
	m.github = srv.URL
	m.packages = []Package{{
		ID: "tool", Name: "Tool", Check: "command -v tool",
		Installers: map[string]Installer{"linux": {Release: &Release{
			Repo: "owner/tool", Asset: "tool-{tag}-{arch}.tar.gz",
			Arch: map[string]string{"arm64": "aarch64"}, Bins: []string{"bin/tool"},
		}}},
	}}
	selected := map[string]bool{"tool": true}
	run := func() string {
		t.Helper()
		if err := m.install(selected, true, io.Discard); err != nil {
			t.Fatal(err)
		}
		out, err := exec.Command(filepath.Join(m.root, "bin", "tool")).Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	if got := run(); got != "v1.0.0" {
		t.Fatalf("installed tool printed %q", got)
	}
	if installed, latest, err := m.checkRelease(context.Background(), m.packages[0]); err != nil || installed != "v1.0.0" || latest != "v1.0.0" {
		t.Fatalf("check = %q, %q, %v", installed, latest, err)
	}
	latest = "v1.1.0"
	var report bytes.Buffer
	if err := m.outdated(&report); err != nil || !strings.Contains(report.String(), "v1.0.0") || !strings.Contains(report.String(), "→ v1.1.0") {
		t.Fatalf("outdated = %q, %v", &report, err)
	}
	if got := run(); got != "v1.1.0" {
		t.Fatalf("upgraded tool printed %q", got)
	}
	leftovers, _ := filepath.Glob(filepath.Join(m.root, "dep", ".omd-*"))
	if len(leftovers) != 0 {
		t.Fatalf("staging directories left behind: %v", leftovers)
	}
}

func TestExtractRejectsEscapingPaths(t *testing.T) {
	for _, name := range []string{"../evil", "/etc/evil", "a/../../evil"} {
		if err := extractTarGz(bytes.NewReader(tarball(t, name, "x")), t.TempDir()); err == nil {
			t.Errorf("%s was extracted", name)
		}
	}
}

func tarball(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0755, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := errors.Join(tw.Close(), gz.Close()); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
