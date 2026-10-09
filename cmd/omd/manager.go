package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
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
		filepath.Join(m.home, ".pyenv", "bin"),
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

func (m *manager) check(script string) bool {
	if script == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.command(ctx, script).Run() == nil
}

func (m *manager) linkConfigs(out io.Writer) error {
	target := filepath.Join(m.root, "etc", "kitty")
	if info, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Kitty source config does not exist: %s", target)
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("Kitty source config is not a directory: %s", target)
	}
	link := filepath.Join(m.home, ".config", "kitty")
	info, err := os.Lstat(link)
	if err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			backup, err := configBackupPath(link)
			if err != nil {
				return err
			}
			if err := os.Rename(link, backup); err != nil {
				return fmt.Errorf("could not preserve existing Kitty config at %s: %w", backup, err)
			}
			if err := os.Symlink(target, link); err != nil {
				if restoreErr := os.Rename(backup, link); restoreErr != nil {
					return errors.Join(err, fmt.Errorf("could not restore the existing Kitty config from %s: %w", backup, restoreErr))
				}
				return err
			}
			fmt.Fprintf(out, "Moved existing Kitty config to %s and linked %s → %s\n", backup, link, target)
			return nil
		}
		wanted, err := filepath.EvalSymlinks(target)
		if err != nil {
			return err
		}
		resolved, err := filepath.EvalSymlinks(link)
		if err != nil {
			return fmt.Errorf("could not resolve existing Kitty symlink %s: %w", link, err)
		}
		if resolved == wanted {
			return nil
		}
		return fmt.Errorf("Kitty destination is already a symlink to another location and was preserved: %s", link)
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
	ordered := slices.Clone(m.packages)
	slices.SortStableFunc(ordered, func(a, b Package) int {
		if a.ID == b.ID {
			return 0
		}
		if a.ID == "homebrew" {
			return -1
		}
		if b.ID == "homebrew" {
			return 1
		}
		return 0
	})
	var failures []error
	installed, upgraded, existing := 0, 0, 0
	for _, p := range ordered {
		if !selected[p.ID] {
			continue
		}
		if m.check(p.Check) {
			ok, err := m.upgradeRelease(p, output)
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
	if p.ID == "kitty-config-link" {
		if err := m.linkConfigs(out); err != nil {
			return err
		}
	} else if installer := p.installerFor(m.platform); installer.Release != nil {
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
	if !m.check(p.Check) {
		return errors.New("installer exited successfully, but the tool is still unavailable")
	}
	return nil
}

// upgradeRelease reinstalls an omd-managed release when a newer one exists.
// Tools from other sources, such as apt, are left to their package manager.
func (m *manager) upgradeRelease(p Package, out io.Writer) (bool, error) {
	installed, latest, err := m.checkRelease(context.Background(), p)
	if err != nil {
		// Being offline must not turn an installed tool into a failure.
		fmt.Fprintf(out, "! %s: could not check for updates: %v\n", p.Name, err)
		return false, nil
	}
	if latest == "" || latest == installed {
		return false, nil
	}
	fmt.Fprintf(out, "\n→ Upgrading %s %s → %s\n", p.Name, installed, latest)
	if err := m.installRelease(p, m.releaseFor(p), latest, out); err != nil {
		return false, err
	}
	fmt.Fprintf(out, "✓ %s upgraded\n", p.Name)
	return true, nil
}

// outdated lists every omd-managed release with its installed and latest version.
func (m *manager) outdated(out io.Writer) error {
	var failures []error
	found := false
	for _, p := range m.packages {
		installed, latest, err := m.checkRelease(context.Background(), p)
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
		fmt.Fprintf(out, "%-12s %-10s %s\n", p.Name, installed, state)
	}
	if !found && len(failures) == 0 {
		fmt.Fprintln(out, "No tools installed from GitHub releases by omd.")
	}
	return errors.Join(failures...)
}
