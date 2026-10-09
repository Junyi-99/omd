package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

func showColorSupport(out io.Writer) error {
	terminal := os.Getenv("TERM")
	colorTerminal := os.Getenv("COLORTERM")
	colors := terminfoColors(terminal)
	if colors <= 0 {
		colors = inferColors(terminal)
	}

	trueColor := supportsTrueColor(terminal, colorTerminal)
	if trueColor {
		colors = 1 << 24
	}
	if os.Getenv("NO_COLOR") != "" {
		colors = 0
	}

	if _, err := fmt.Fprintf(out, "TTY: %s\n", yesNo(term.IsTerminal(os.Stdout.Fd()))); err != nil {
		return err
	}
	if terminal == "" {
		terminal = "(unset)"
	}
	if colorTerminal == "" {
		colorTerminal = "(unset)"
	}
	if _, err := fmt.Fprintf(out, "TERM: %s\nCOLORTERM: %s\n", terminal, colorTerminal); err != nil {
		return err
	}
	if colors == 0 {
		if os.Getenv("NO_COLOR") != "" {
			_, err := fmt.Fprintln(out, "Color support: disabled by NO_COLOR")
			return err
		}
		_, err := fmt.Fprintln(out, "Color support: none detected")
		return err
	}
	if trueColor && os.Getenv("NO_COLOR") == "" {
		_, err := fmt.Fprintln(out, "Color support: truecolor (24-bit, 16.7 million colors)")
		return err
	}
	_, err := fmt.Fprintf(out, "Color support: %s (%d colors)\n", colorDepth(colors), colors)
	return err
}

func terminfoColors(terminal string) int {
	if terminal == "" || terminal == "dumb" {
		return 0
	}
	output, err := exec.Command("tput", "colors").Output()
	if err != nil {
		return 0
	}
	colors, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil || colors < 0 {
		return 0
	}
	return colors
}

func inferColors(terminal string) int {
	terminal = strings.ToLower(terminal)
	switch {
	case terminal == "" || terminal == "dumb":
		return 0
	case strings.Contains(terminal, "256color"):
		return 256
	case strings.Contains(terminal, "16color"):
		return 16
	default:
		return 8
	}
}

func supportsTrueColor(terminal, colorTerminal string) bool {
	terminal = strings.ToLower(terminal)
	colorTerminal = strings.ToLower(colorTerminal)
	return colorTerminal == "truecolor" || colorTerminal == "24bit" ||
		strings.Contains(terminal, "truecolor") || strings.Contains(terminal, "-direct") ||
		strings.Contains(terminal, "kitty") || strings.Contains(terminal, "wezterm") ||
		strings.Contains(terminal, "alacritty") || strings.Contains(terminal, "foot")
}

func colorDepth(colors int) string {
	switch {
	case colors >= 1<<24:
		return "truecolor"
	case colors >= 256:
		return "256-color"
	case colors >= 16:
		return "16-color"
	case colors > 0:
		return "basic"
	default:
		return "none"
	}
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
