#!/usr/bin/env python3
"""Small tmux creation dialog. Uses argv (never shell interpolation) for names."""
import curses
import os
import subprocess
import sys
import unicodedata

KIND = sys.argv[1] if len(sys.argv) > 1 else "window"
START = os.getcwd()


def tmux(*args):
    return subprocess.run(["tmux", *args], text=True, capture_output=True)


def clipped(text, width):
    result, used = "", 0
    for char in text:
        size = 2 if unicodedata.east_asian_width(char) in "WF" else 1
        if used + size > max(0, width):
            break
        result += char
        used += size
    return result, used


def dialog(screen):
    curses.use_default_colors()
    for index, option, fallback in [(1, "accent", 111), (2, "dim", 103),
                                    (3, "fg", 189), (4, "alert", 210)]:
        color = tmux("show-option", "-gqv", "@popup_" + option).stdout.strip()
        value = fallback
        if color.startswith("#") and len(color) == 7:
            rgb = [int(color[i:i + 2], 16) for i in (1, 3, 5)]
            value = 16 + 36 * round(rgb[0] / 51) + 6 * round(rgb[1] / 51) + round(rgb[2] / 51)
        curses.init_pair(index, value if curses.COLORS >= 256 else index, -1)
    curses.curs_set(1)
    name, error = "", ""

    while True:
        screen.erase()
        height, width = screen.getmaxyx()
        left = 3 if width >= 40 else 1
        inner = max(1, width - 2 * left)

        def put(row, text, color=3, bold=False):
            if 0 <= row < height:
                try:
                    screen.addstr(row, left, clipped(text, inner)[0],
                                  curses.color_pair(color) | (curses.A_BOLD if bold else 0))
                except curses.error:
                    pass

        home = os.path.expanduser("~")
        path = "~" + START[len(home):] if START == home or START.startswith(home + "/") else START
        put(1, "Create a new " + KIND, 1, True)
        put(3, "Directory  " + path, 2)
        put(5, KIND.capitalize() + " name", 3, True)
        shown = name
        while clipped(shown, inner - 3)[0] != shown and shown:
            shown = shown[1:]
        put(6, "> " + (shown or "Automatic name"), 3 if name else 2)
        put(7, "─" * inner, 1)
        put(8, error, 4)
        put(height - 2, "Enter  create    Esc  cancel    Ctrl+U  clear", 2)
        screen.move(min(6, height - 1), min(left + 2 + clipped(shown, inner - 3)[1], width - 1))
        screen.refresh()
        try:
            key = screen.get_wch()
        except KeyboardInterrupt:
            return
        if key in ("\x1b", "\x03", "\x04"):
            return
        if key in ("\n", "\r", curses.KEY_ENTER):
            if KIND == "session" and any(c in name for c in ".:"):
                error = "Session names cannot contain . or :"
                continue
            if os.environ.get("DRYRUN"):
                return
            if KIND == "session":
                args = ["new-session", "-d", "-P", "-F", "#{session_id}", "-c", START]
                if name:
                    args += ["-s", name]
            else:
                # popup 里的 tmux 命令默认作用于发起它的 client/session
                args = ["new-window", "-c", START]
                if name:
                    args += ["-n", name]
            result = tmux(*args)
            if result.returncode:
                error = result.stderr.strip() or "Could not create " + KIND
                continue
            if KIND == "session":
                result = tmux("switch-client", "-t", result.stdout.strip())
                if result.returncode:
                    error = "Created; switch failed. Esc to close."
                    # Do not create another session if switching the client failed.
                    put(8, error, 4)
                    screen.refresh()
                    screen.get_wch()
            return
        if key in (curses.KEY_BACKSPACE, "\x7f", "\b"):
            name = name[:-1]
        elif key == "\x15":
            name = ""
        elif isinstance(key, str) and key.isprintable() and len(name) < 128:
            name += key
        error = ""


if __name__ == "__main__":
    if KIND not in ("window", "session"):
        raise SystemExit("Expected window or session")
    curses.wrapper(dialog)
