package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func testManager(t *testing.T, platform string, pkgs []Package) *manager {
	t.Helper()
	dir := t.TempDir()
	return &manager{root: filepath.Join(dir, "root"), home: filepath.Join(dir, "home"), platform: platform, arch: "amd64", packages: pkgs}
}

func ids(pkgs []Package) []string {
	var out []string
	for _, p := range pkgs {
		out = append(out, p.ID)
	}
	return out
}

func TestInstallOrderTopological(t *testing.T) {
	brew := Installer{Requires: []string{"brew"}, Commands: []string{"true"}}
	m := testManager(t, "darwin", []Package{
		{ID: "agent", Needs: []string{"node"}},
		{ID: "tool", Installers: map[string]Installer{"darwin": brew}},
		{ID: "editor", Needs: []string{"node", "ts"}},
		{ID: "ts", Needs: []string{"node"}},
		{ID: "node"},
		{ID: "homebrew", Provides: []string{"brew"}},
		{ID: "unselected"},
	})
	sel := m.withNeeds(map[string]bool{"agent": true, "tool": true, "editor": true, "homebrew": true})
	got := ids(m.installOrder(sel))
	want := []string{"node", "agent", "homebrew", "tool", "ts", "editor"}
	if !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	if sel["unselected"] {
		t.Fatal("unselected package pulled in")
	}
}

func TestInstallOrderRealRegistry(t *testing.T) {
	for _, platform := range []string{"darwin", "linux"} {
		m := testManager(t, platform, packagesFor(platform))
		order := ids(m.installOrder(m.withNeeds(m.defaultSelection())))
		pos := func(id string) int { return slices.Index(order, id) }
		for _, p := range m.packages {
			for _, dep := range m.dependencies(p) {
				if pos(dep) >= 0 && pos(p.ID) >= 0 && pos(dep) > pos(p.ID) {
					t.Errorf("%s: %s installed before its dependency %s", platform, p.ID, dep)
				}
			}
		}
		if platform == "darwin" && pos("homebrew") > pos("starship") {
			t.Errorf("homebrew after a brew-installed tool on darwin: %v", order)
		}
	}
}

func TestConfigRoundTrip(t *testing.T) {
	m := testManager(t, "linux", []Package{{ID: "a"}, {ID: "b"}, {ID: "c"}})
	sel := map[string]bool{"a": true, "b": false, "c": true}
	if err := m.saveConfig(sel); err != nil {
		t.Fatal(err)
	}
	got, err := m.loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	for id, want := range sel {
		if got[id] != want {
			t.Errorf("%s = %v, want %v", id, got[id], want)
		}
	}
	if err := os.WriteFile(m.configPath(), []byte("a=enabled\nmystery=enabled\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := m.loadConfig(); err == nil || !strings.Contains(err.Error(), "mystery") {
		t.Fatalf("unknown ID error = %v", err)
	}
}

func TestGenerateEnvOrder(t *testing.T) {
	ok := "true"
	m := testManager(t, "linux", []Package{
		{ID: "late", Check: ok, Env: "LATE\n", EnvOrder: 1},
		{ID: "mid", Check: ok, Env: "MID\n"},
		{ID: "first", Check: ok, Env: "FIRST\n", EnvOrder: -1},
		{ID: "missing", Check: "false", Env: "MISSING\n"},
		{ID: "off", Check: ok, Env: "OFF\n"},
	})
	sel := map[string]bool{"late": true, "mid": true, "first": true, "missing": true}
	if err := m.generateEnv(sel); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(m.root, "etc", "omd", "env.sh"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if strings.Contains(s, "MISSING") || strings.Contains(s, "OFF") {
		t.Fatalf("unexpected snippets:\n%s", s)
	}
	if !(strings.Index(s, "FIRST") < strings.Index(s, "MID") && strings.Index(s, "MID") < strings.Index(s, "LATE")) {
		t.Fatalf("wrong order:\n%s", s)
	}
}

func TestProbeSemantics(t *testing.T) {
	m := testManager(t, "linux", nil)
	if ok, err := m.probe("true"); !ok || err != nil {
		t.Fatalf("true: %v %v", ok, err)
	}
	if ok, err := m.probe("false"); ok || err != nil {
		t.Fatalf("false must be missing without error: %v %v", ok, err)
	}
	if ok, err := m.probe(""); ok || err != nil {
		t.Fatalf("empty: %v %v", ok, err)
	}
	if testing.Short() {
		return
	}
	ok, err := m.probe("sleep 30")
	if ok || err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("timeout must be an error, not missing: %v %v", ok, err)
	}
	pkgs := []Package{{ID: "x", Name: "x", Check: "true"}, {ID: "y", Name: "y", Check: "false"}}
	got, err := m.probeAll(pkgs)
	if err != nil || !got["x"] || got["y"] {
		t.Fatalf("probeAll = %v, %v", got, err)
	}
}

func TestLinkDir(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer

	if err := linkDir(filepath.Join(dir, "nope"), filepath.Join(dir, "l"), &out); err == nil {
		t.Fatal("missing source should fail")
	}

	// Fresh link creates parents.
	link := filepath.Join(dir, "home", ".config", "x")
	if err := linkDir(src, link, &out); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(link); got != src {
		t.Fatalf("link -> %q", got)
	}
	// Idempotent.
	out.Reset()
	if err := linkDir(src, link, &out); err != nil || out.Len() != 0 {
		t.Fatalf("second run: %v %q", err, out.String())
	}

	// Existing directory is backed up.
	existing := filepath.Join(dir, "existing")
	if err := os.MkdirAll(existing, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(existing, "keep"), []byte("1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := linkDir(src, existing, &out); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(existing); got != src {
		t.Fatalf("existing not linked: %q", got)
	}
	backups, _ := filepath.Glob(existing + ".omd-backup-*")
	if len(backups) != 1 {
		t.Fatalf("backups = %v", backups)
	}
	if _, err := os.Stat(filepath.Join(backups[0], "keep")); err != nil {
		t.Fatal("backup lost contents")
	}

	// Symlink elsewhere is a conflict and preserved.
	other := filepath.Join(dir, "other")
	if err := os.MkdirAll(other, 0755); err != nil {
		t.Fatal(err)
	}
	conflict := filepath.Join(dir, "conflict")
	if err := os.Symlink(other, conflict); err != nil {
		t.Fatal(err)
	}
	if err := linkDir(src, conflict, &out); err == nil || !strings.Contains(err.Error(), "another location") {
		t.Fatalf("conflict error = %v", err)
	}
	if got, _ := os.Readlink(conflict); got != other {
		t.Fatal("conflicting symlink was modified")
	}
}
