package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

type manager struct {
	root, home, platform string
	// arch is runtime.GOARCH; the launcher builds omd natively, so it also
	// matches the userland of a 64-bit kernel running 32-bit binaries.
	arch     string
	github   string
	packages []Package
}

func newManager() (*manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	root := os.Getenv("USER_CONFIG_HOME")
	if root == "" {
		root = filepath.Join(home, ".oh-my-dotfiles")
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &manager{
		root:     root,
		home:     home,
		platform: runtime.GOOS,
		arch:     runtime.GOARCH,
		github:   "https://github.com",
		packages: packagesFor(runtime.GOOS),
	}, nil
}

// Discover new installs for every command, including Homebrew installed earlier
// in this run. No shell startup files or generated env.sh are executed here.
func (m *manager) environment() []string {
	paths := []string{
		filepath.Join(m.root, "bin"),
		filepath.Join(m.home, ".local", "bin"),
	}
	if m.platform == "darwin" {
		paths = append(paths, "/opt/homebrew/bin", "/opt/homebrew/sbin", "/usr/local/bin")
	}
	paths = append(paths, os.Getenv("PATH"))
	env := os.Environ()
	return append(env,
		"OMD_HOME="+m.root,
		"USER_CONFIG_HOME="+m.root,
		"PATH="+strings.Join(paths, string(os.PathListSeparator)),
	)
}

func (m *manager) command(ctx context.Context, script string) *exec.Cmd {
	// Alpine ships no bash; its busybox sh runs these POSIX scripts and supports pipefail.
	args := []string{"--noprofile", "--norc", "-eo", "pipefail", "-c", script}
	shell, err := exec.LookPath("bash")
	if err != nil {
		shell, args = "sh", args[2:]
	}
	cmd := exec.CommandContext(ctx, shell, args...)
	cmd.Env = m.environment()
	return cmd
}

const checkTimeout = 5 * time.Second

// probe runs a Check script. A nonzero exit means "not installed" (false, nil);
// a timeout or a shell that cannot start is an error, so callers never mistake
// a slow check for a missing tool.
func (m *manager) probe(script string) (bool, error) {
	if script == "" {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()
	err := m.command(ctx, script).Run()
	switch {
	case err == nil:
		return true, nil
	case ctx.Err() != nil:
		return false, fmt.Errorf("check timed out after %s", checkTimeout)
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return false, nil
	}
	return false, err
}

// check reports whether the script succeeds; a timeout or error counts as
// false. Use it only for display; installation decisions use probe.
func (m *manager) check(script string) bool {
	ok, _ := m.probe(script)
	return ok
}

// probeAll runs the Check of every package concurrently. The map holds the
// installed state of each package whose check completed; the error joins the
// checks that timed out or could not run.
func (m *manager) probeAll(pkgs []Package) (map[string]bool, error) {
	oks := make([]bool, len(pkgs))
	errs := make([]error, len(pkgs))
	var wg sync.WaitGroup
	for i, p := range pkgs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			if oks[i], err = m.probe(p.Check); err != nil {
				errs[i] = fmt.Errorf("%s: %w", p.Name, err)
			}
		}()
	}
	wg.Wait()
	installed := make(map[string]bool, len(pkgs))
	for i, p := range pkgs {
		if errs[i] == nil {
			installed[p.ID] = oks[i]
		}
	}
	return installed, errors.Join(errs...)
}

// linkDir symlinks link to the directory target. An existing real file or
// directory at link is moved to a timestamped backup first; an existing
// symlink to another location is preserved and reported as a conflict.
func linkDir(target, link string, out io.Writer) error {
	if info, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("source does not exist: %s", target)
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", target)
	}
	info, err := os.Lstat(link)
	if err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			backup, err := configBackupPath(link)
			if err != nil {
				return err
			}
			if err := os.Rename(link, backup); err != nil {
				return fmt.Errorf("could not preserve existing config at %s: %w", backup, err)
			}
			if err := os.Symlink(target, link); err != nil {
				if restoreErr := os.Rename(backup, link); restoreErr != nil {
					return errors.Join(err, fmt.Errorf("could not restore the existing config from %s: %w", backup, restoreErr))
				}
				return err
			}
			fmt.Fprintf(out, "Moved existing config to %s and linked %s → %s\n", backup, link, target)
			return nil
		}
		wanted, err := filepath.EvalSymlinks(target)
		if err != nil {
			return err
		}
		resolved, err := filepath.EvalSymlinks(link)
		if err != nil {
			return fmt.Errorf("could not resolve existing symlink %s: %w", link, err)
		}
		if resolved == wanted {
			return nil
		}
		return fmt.Errorf("destination is already a symlink to another location and was preserved: %s", link)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		return err
	}
	if err := os.Symlink(target, link); err != nil {
		return err
	}
	fmt.Fprintf(out, "Linked %s → %s\n", link, target)
	return nil
}

func configBackupPath(path string) (string, error) {
	base := path + ".omd-backup-" + time.Now().Format("20060102-150405")
	candidate := base
	for suffix := 1; ; suffix++ {
		if _, err := os.Lstat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate, nil
		} else if err != nil {
			return "", err
		}
		candidate = fmt.Sprintf("%s-%d", base, suffix)
	}
}

// withNeeds returns selected plus the prerequisites of every selected package,
// transitively.
func (m *manager) withNeeds(selected map[string]bool) map[string]bool {
	all := maps.Clone(selected)
	for changed := true; changed; {
		changed = false
		for _, p := range m.packages {
			if !all[p.ID] {
				continue
			}
			for _, id := range p.Needs {
				if !all[id] {
					all[id], changed = true, true
				}
			}
		}
	}
	return all
}

// dependencies lists the IDs that must be installed before p: its Needs, plus
// any package that Provides a command its installer Requires (such as brew).
func (m *manager) dependencies(p Package) []string {
	deps := slices.Clone(p.Needs)
	for _, req := range p.installerFor(m.platform).Requires {
		for _, q := range m.packages {
			if q.ID != p.ID && slices.Contains(q.Provides, req) {
				deps = append(deps, q.ID)
			}
		}
	}
	return deps
}

// installOrder returns the selected packages with dependencies first. It is a
// stable depth-first topological sort: independent packages keep registry order.
func (m *manager) installOrder(selected map[string]bool) []Package {
	byID := make(map[string]Package, len(m.packages))
	for _, p := range m.packages {
		byID[p.ID] = p
	}
	var ordered []Package
	state := make(map[string]int) // 1 visiting, 2 done
	var visit func(p Package)
	visit = func(p Package) {
		if state[p.ID] != 0 {
			return
		}
		state[p.ID] = 1
		for _, id := range m.dependencies(p) {
			if dep, ok := byID[id]; ok && selected[id] {
				visit(dep)
			}
		}
		state[p.ID] = 2
		ordered = append(ordered, p)
	}
	for _, p := range m.packages {
		if selected[p.ID] {
			visit(p)
		}
	}
	return ordered
}

func (m *manager) install(selected map[string]bool, quiet bool, out io.Writer) error {
	if m.platform != "darwin" && m.platform != "linux" {
		return fmt.Errorf("unsupported platform: %s", m.platform)
	}
	for _, dir := range []string{"bin", "dep", "cache/omd"} {
		if err := os.MkdirAll(filepath.Join(m.root, dir), 0755); err != nil {
			return err
		}
	}
	logPath := filepath.Join(m.root, "cache", "omd", "install.log")
	log, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer log.Close()
	var output io.Writer = log
	if !quiet {
		output = io.MultiWriter(out, log)
	}
	selected = m.withNeeds(selected)
	ordered := m.installOrder(selected)
	// One concurrent pass decides install vs upgrade; installPackage verifies
	// afterwards. A check that timed out is reported, never treated as missing.
	present, checkErr := m.probeAll(ordered)
	var failures []error
	if checkErr != nil {
		failures = append(failures, checkErr)
		fmt.Fprintf(output, "✗ %v\n", checkErr)
	}
	installed, upgraded, existing := 0, 0, 0
	for _, p := range ordered {
		have, checked := present[p.ID]
		if !checked {
			continue
		}
		if have {
			ok, err := m.upgradeTool(p, output)
			switch {
			case err != nil:
				failures = append(failures, fmt.Errorf("%s: %w", p.Name, err))
				fmt.Fprintf(output, "✗ %s: %v\n", p.Name, err)
			case ok:
				upgraded++
			default:
				existing++
			}
			continue
		}
		fmt.Fprintf(output, "\n→ Installing %s\n", p.Name)
		err := m.installPackage(p, output)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", p.Name, err))
			fmt.Fprintf(output, "✗ %s: %v\n", p.Name, err)
			continue
		}
		installed++
		fmt.Fprintf(output, "✓ %s installed\n", p.Name)
	}
	if err := m.generateEnv(selected); err != nil {
		failures = append(failures, err)
	}
	if len(failures) > 0 {
		return fmt.Errorf("installation incomplete (%d failures); log: %s\n%w", len(failures), logPath, errors.Join(failures...))
	}
	if !quiet {
		fmt.Fprintf(out, "\n✓ %d installed, %d upgraded, %d already up to date. Shell integration generated.\n", installed, upgraded, existing)
	}
	return nil
}

func (m *manager) installPackage(p Package, out io.Writer) error {
	if installer := p.installerFor(m.platform); installer.Link != nil {
		link := installer.Link
		if err := linkDir(filepath.Join(m.root, link.Source), filepath.Join(m.home, link.Dest), out); err != nil {
			return err
		}
	} else if installer.Release != nil {
		if err := m.installRelease(p, installer.Release, "", out); err != nil {
			return err
		}
	} else {
		commands := installer.Commands
		if len(commands) == 0 {
			return errors.New("no automatic installer; install this tool manually")
		}
		for _, dep := range installer.Requires {
			if !m.check("command -v " + shellQuote(dep)) {
				return fmt.Errorf("missing prerequisite %s; install/enable it first", dep)
			}
		}
		for _, script := range commands {
			cmd := m.command(context.Background(), script)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, out, out
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	ok, err := m.probe(p.Check)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("installer exited successfully, but the tool is still unavailable")
	}
	return nil
}

// updater knows how to find and apply a newer version of one kind of
// omd-managed tool: a GitHub release or a package with Version scripts.
type updater interface {
	// versions returns the installed and latest versions; installed is empty
	// when omd has nothing to update.
	versions(ctx context.Context, m *manager, p Package) (installed, latest string, err error)
	upgrade(m *manager, p Package, latest string, out io.Writer) error
}

// updaterFor returns how p is updated, or nil for tools left to their own
// package manager.
func (m *manager) updaterFor(p Package) updater {
	if p.Version != nil {
		return p.Version
	}
	if rel := m.releaseFor(p); rel != nil {
		return releaseUpdater{rel}
	}
	return nil
}

func (v *Version) versions(ctx context.Context, m *manager, p Package) (installed, latest string, err error) {
	if ok, err := m.probe(p.Check); !ok {
		return "", "", err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	run := func(script string) (string, error) {
		b, err := m.command(ctx, script).Output()
		return strings.TrimSpace(string(b)), err
	}
	if installed, err = run(v.Installed); err != nil {
		return "", "", err
	}
	if installed == "" {
		// Present but not from this package manager, such as a native install.
		return "", "", nil
	}
	if latest, err = run(v.Latest); err != nil {
		return "", "", err
	}
	return installed, latest, nil
}

func (v *Version) upgrade(m *manager, _ Package, _ string, out io.Writer) error {
	cmd := m.command(context.Background(), v.Upgrade)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, out, out
	return cmd.Run()
}

// releaseUpdater updates a tool installed from a GitHub release.
type releaseUpdater struct{ rel *Release }

func (u releaseUpdater) versions(ctx context.Context, m *manager, p Package) (string, string, error) {
	return m.checkRelease(ctx, p)
}

func (u releaseUpdater) upgrade(m *manager, p Package, latest string, out io.Writer) error {
	return m.installRelease(p, u.rel, latest, out)
}

// checkTool returns the installed and latest versions of an omd-managed tool;
// latest is empty for anything omd does not update.
func (m *manager) checkTool(ctx context.Context, p Package) (installed, latest string, err error) {
	u := m.updaterFor(p)
	if u == nil {
		return "", "", nil
	}
	return u.versions(ctx, m, p)
}

// upgradeTool updates an omd-managed tool when a newer version exists.
// Tools from other sources, such as apt, are left to their package manager.
func (m *manager) upgradeTool(p Package, out io.Writer) (bool, error) {
	installed, latest, err := m.checkTool(context.Background(), p)
	if err != nil {
		// Being offline must not turn an installed tool into a failure.
		fmt.Fprintf(out, "! %s: could not check for updates: %v\n", p.Name, err)
		return false, nil
	}
	if latest == "" || latest == installed {
		return false, nil
	}
	fmt.Fprintf(out, "\n→ Upgrading %s %s → %s\n", p.Name, installed, latest)
	if err := m.updaterFor(p).upgrade(m, p, latest, out); err != nil {
		return false, err
	}
	fmt.Fprintf(out, "✓ %s upgraded\n", p.Name)
	return true, nil
}

// outdated lists the version of every enabled omd-managed tool, or of all of
// them when selected is nil.
func (m *manager) outdated(out io.Writer, selected map[string]bool) error {
	var failures []error
	found := false
	type result struct {
		installed, latest string
		err               error
	}
	results := make([]result, len(m.packages))
	var wg sync.WaitGroup
	for i, p := range m.packages {
		if selected != nil && !selected[p.ID] {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := &results[i]
			r.installed, r.latest, r.err = m.checkTool(context.Background(), p)
		}()
	}
	wg.Wait()
	for i, p := range m.packages {
		if selected != nil && !selected[p.ID] {
			continue
		}
		installed, latest, err := results[i].installed, results[i].latest, results[i].err
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", p.Name, err))
			continue
		}
		if latest == "" {
			continue
		}
		found = true
		state := "up to date"
		if latest != installed {
			state = "→ " + latest + " (omd install upgrades selected tools)"
		}
		fmt.Fprintf(out, "%-16s %-10s %s\n", p.Name, installed, state)
	}
	if !found && len(failures) == 0 {
		fmt.Fprintln(out, "No omd-managed tools to check.")
	}
	return errors.Join(failures...)
}
