package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

// The shell points XDG_CONFIG_HOME and XDG_DATA_HOME into the repo (etc/ and
// data/), which hides whatever tools had already stored under the standard
// locations: gh logins, git config, application data. linkLegacyXDG links each
// entry of the standard directories into the repo so those tools keep finding
// it, and shells that do not load omd (bash, cron) keep sharing the same files.
//
// An entry the repo already has (a tracked config such as etc/nvim, or one a
// tool created after the switch) is left alone and reported as a conflict.
func (m *manager) linkLegacyXDG() (linked, conflicts []string, err error) {
	pairs := [][2]string{
		{filepath.Join(m.home, ".config"), filepath.Join(m.root, "etc")},
		{filepath.Join(m.home, ".local", "share"), filepath.Join(m.root, "data")},
	}
	var errs []error
	for _, pair := range pairs {
		legacy, repo := pair[0], pair[1]
		if same, _ := sameDir(legacy, repo); same {
			continue
		}
		entries, readErr := os.ReadDir(legacy)
		if errors.Is(readErr, os.ErrNotExist) {
			continue
		}
		if readErr != nil {
			errs = append(errs, readErr)
			continue
		}
		if err := os.MkdirAll(repo, 0755); err != nil {
			errs = append(errs, err)
			continue
		}
		for _, e := range entries {
			src := filepath.Join(legacy, e.Name())
			dst := filepath.Join(repo, e.Name())
			if target, err := os.Readlink(dst); err == nil && target == src {
				continue
			}
			if _, err := os.Lstat(dst); err == nil {
				conflicts = append(conflicts, fmt.Sprintf("%s (omd uses %s)", src, dst))
				continue
			}
			if err := os.Symlink(src, dst); err != nil {
				errs = append(errs, err)
				continue
			}
			linked = append(linked, dst+" -> "+src)
		}
	}
	slices.Sort(conflicts)
	return linked, conflicts, errors.Join(errs...)
}

// reportLegacyXDG links legacy XDG entries and describes the result.
func (m *manager) reportLegacyXDG(out io.Writer) error {
	linked, conflicts, err := m.linkLegacyXDG()
	for _, l := range linked {
		fmt.Fprintf(out, "Linked %s\n", l)
	}
	if len(conflicts) > 0 {
		fmt.Fprintln(out, "Not linked, the repo already has an entry of the same name; merge them by hand if needed:")
		for _, c := range conflicts {
			fmt.Fprintf(out, "  %s\n", c)
		}
	}
	return err
}

func sameDir(a, b string) (bool, error) {
	ai, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false, err
	}
	return os.SameFile(ai, bi), nil
}
