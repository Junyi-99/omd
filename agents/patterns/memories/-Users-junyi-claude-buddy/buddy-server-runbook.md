---
name: buddy-server-runbook
description: How the Mac-side receiver ("受控 CLI", server/buddy_server.py) is run and where the Xcode toolchain lives on this machine
metadata:
  type: project
---

"受控 CLI" means `server/buddy_server.py`, the Mac receiver. As of 2026-08-25 it runs detached via nohup (user closed the tmux session and asked for plain background): PID in `~/Library/Logs/buddy_server.pid`, log at `~/Library/Logs/buddy_server.log`, venv at `server/venv` (Python 3.14). Bonjour name `robin`, port 8765; the phone reaches it over Tailscale too (100.x address seen). Don't move it to launchd without warning: the Accessibility (TCC) grant is attributed to the launching terminal, and a launchd-spawned python would need its own grant before keystrokes work. The tmux session named `buddy` is the Claude Code window itself — never reuse it.

Toolchain: `xcode-select` points to CommandLineTools; xcodebuild only works with `DEVELOPER_DIR="/Applications/Xcode-beta 2.app/Contents/Developer"` (Xcode 27 beta). Its Simulator.app GUI is missing — `simctl` boot/install/launch/screenshot works headless, but nothing can tap the UI.

**Why:** the venv is gitignored and vanished once already; re-discovering the DEVELOPER_DIR workaround and the missing Simulator GUI cost several probes.

**How to apply:** check `pgrep -f buddy_server.py` first; to start: `cd server && nohup ./venv/bin/python -u buddy_server.py >> ~/Library/Logs/buddy_server.log 2>&1 & echo $! > ~/Library/Logs/buddy_server.pid` (recreate venv from requirements.txt if missing); to stop: `kill $(cat ~/Library/Logs/buddy_server.pid)`. To build: prefix xcodebuild with that DEVELOPER_DIR and use `-destination 'platform=iOS Simulator,name=iPhone 17 Pro'`. Verify UI via screenshot only; no clicking.
