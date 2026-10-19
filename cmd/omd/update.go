package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/x/term"
)

func (m *manager) update(force bool) error {
	return m.updateWithOptions(force, false)
}

func (m *manager) updateWithOptions(force, quiet bool) error {
	git := func(args ...string) (string, error) {
		ctx := context.Background()
		if quiet {
			// Unattended runs must not hang on the network or a credential prompt.
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, gitTimeout)
			defer cancel()
		}
		cmd := exec.CommandContext(ctx, "git", append([]string{"-C", m.root}, args...)...)
		if quiet {
			cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		}
		b, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, b)
		}
		return strings.TrimSpace(string(b)), nil
	}
	branch, err := git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return err
	}
	upstream, err := git("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		return err
	}
	before, err := git("rev-parse", "HEAD")
	if err != nil {
		return err
	}
	remote, err := git("config", "--get", "branch."+branch+".remote")
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Printf("Fetching %s...\n", upstream)
	}
	if _, err = git("fetch", "--quiet", remote); err != nil {
		return err
	}
	ahead, err := git("rev-list", "--count", upstream+"..HEAD")
	if err != nil {
		return err
	}
	if ahead != "0" && !force {
		refuse := fmt.Errorf("%s local commits are not on %s; push them before updating", ahead, upstream)
		if quiet || !term.IsTerminal(os.Stdin.Fd()) {
			return refuse
		}
		behind, err := git("rev-list", "--count", "HEAD.."+upstream)
		if err != nil {
			return err
		}
		if behind != "0" {
			fmt.Printf("Diverged from %s; merging (no rebase)...\n", upstream)
			if _, err := git("merge", "--no-rebase", "--no-edit", upstream); err != nil {
				return err
			}
		}
		fmt.Printf("Local commits are not on %s. Push now? [y/N] ", upstream)
		var answer string
		fmt.Scanln(&answer)
		if !strings.EqualFold(strings.TrimSpace(answer), "y") {
			return refuse
		}
		if _, err := git("push"); err != nil {
			return err
		}
	}
	if force {
		if !quiet {
			fmt.Println("--force: discarding tracked local changes and commits that are not upstream")
		}
		_, err = git("reset", "--hard", upstream)
	} else {
		_, err = git("merge", "--ff-only", upstream)
	}
	if err != nil {
		return err
	}
	after, err := git("rev-parse", "HEAD")
	if err != nil {
		return err
	}
	if !quiet {
		if before == after {
			fmt.Println("oh-my-dotfiles is already up to date.")
		} else {
			fmt.Printf("Updated oh-my-dotfiles to %s.\n", after[:min(12, len(after))])
		}
	}
	return nil
}
