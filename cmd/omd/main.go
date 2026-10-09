package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "omd:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	m, err := newManager()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		selected, err := m.loadConfig()
		if errors.Is(err, os.ErrNotExist) {
			selected, err = make(map[string]bool), nil
		}
		if err != nil {
			return err
		}
		frequency, err := m.loadUpdateFrequency()
		if err != nil {
			return err
		}
		_, err = tea.NewProgram(newModel(m, selected, frequency)).Run()
		return err
	}
	switch args[0] {
	case "help", "--help", "-h":
		fmt.Print(`omd — Oh My Dotfiles

Usage:
  omd                     Open the package selection TUI
  omd install [--quiet]    Install saved selections, upgrade outdated release tools,
                          and regenerate shell env
  omd outdated            List tools installed from GitHub releases with newer versions
  omd env                 Regenerate shell env without installing anything
  omd update [--force]     Update the oh-my-dotfiles repository
  omd auto-update          Run the scheduled background update if due
  omd color                Detect terminal color support
  omd version             Show the compiled Go and UI dependency versions

TUI: Space toggle · Tab switch · m/g/f presets · s save · Enter save + install
Disabling a package stops its shell integration; it does not uninstall it.
--force discards tracked local changes and local commits during update.
`)
		return nil
	case "version", "--version":
		fmt.Println("omd ·", runtime.Version())
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, dep := range info.Deps {
				if dep.Path == "charm.land/bubbletea/v2" || dep.Path == "charm.land/lipgloss/v2" {
					fmt.Println(dep.Path, dep.Version)
				}
			}
		}
		return nil
	case "outdated":
		if len(args) != 1 {
			return fmt.Errorf("outdated takes no arguments")
		}
		return m.outdated(os.Stdout)
	case "color":
		if len(args) != 1 {
			return fmt.Errorf("color takes no arguments")
		}
		return showColorSupport(os.Stdout)
	case "install", "env":
		fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
		quiet := fs.Bool("quiet", false, "hide successful install output; retain errors and install.log")
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return err
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("unexpected arguments: %v", fs.Args())
		}
		selected, err := m.loadConfig()
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no packages.conf; run omd to select packages first")
		}
		if err != nil {
			return err
		}
		if args[0] == "env" {
			return m.generateEnv(selected)
		}
		return m.install(selected, *quiet, os.Stdout)
	case "update":
		fs := flag.NewFlagSet("update", flag.ContinueOnError)
		force := fs.Bool("force", false, "discard tracked changes and local commits")
		fs.BoolVar(force, "f", false, "alias for --force")
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return err
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("unexpected arguments: %v", fs.Args())
		}
		return m.update(*force)
	case "auto-update":
		if len(args) != 1 {
			return fmt.Errorf("auto-update takes no arguments")
		}
		return m.autoUpdate()
	default:
		return fmt.Errorf("unknown command %q; run omd --help", args[0])
	}
}
