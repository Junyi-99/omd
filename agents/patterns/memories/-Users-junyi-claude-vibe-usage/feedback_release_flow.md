---
name: Release flow for vibe-usage
description: When user says "发新版本" / "make a release" without specifying version, default to patch bump, tag, wait for GH Actions, then edit release notes
type: feedback
originSessionId: 65363cf2-cb80-4525-92a4-124f82fd5956
---
When the user asks for a new release without specifying a version number, the default flow is:

1. Bump the **patch** version in `Cargo.toml` (e.g. 0.2.9 → 0.2.10)
2. Run `cargo check` to sync `Cargo.lock`
3. Commit (message like `release: bump version to X.Y.Z`)
4. Push to `main`, tag `vX.Y.Z`, push tag
5. Wait for GitHub Actions release workflow to finish (~2 minutes, builds 4-platform binaries + updates Homebrew)
6. Edit the release notes via `gh release edit vX.Y.Z --notes "..."`

**Why:** Stated by user on 2026-04-17 after the v0.2.9 release. They want the release to be mostly autonomous — don't ask which version component to bump, assume patch.

**How to apply:**
- If user says "发新版本" / "make a release" / "bump and release" with no version specified → patch bump, don't ask.
- Still commit any outstanding feature work first (separate from the version bump commit) if the working tree is dirty.
- If the CI fails (like v0.2.9's `npm ci` / `package-lock.json` mismatch), fix it with an extra commit and move the tag; don't create a new version.
- Release notes should cover user-visible changes since the previous tag — read `git log vPREV..HEAD` to compose.
