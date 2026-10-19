---
name: GitHub code search does not support pushed: / date qualifiers
description: GitHub's code search (type=code) ignores repo-level qualifiers like `pushed:>YYYY-MM-DD`; don't try to filter code search results by commit/push date
type: reference
originSessionId: 684260bd-b5d2-47ef-a332-8f3b3ef464cd
---
GitHub's code search UI (`https://github.com/search?type=code`) does **not** support the `pushed:` qualifier — that qualifier is only valid in repository search (`type=repositories`). Adding `+pushed:>YYYY-MM-DD` to a code query is silently dropped or breaks the query.

**Verified empirically** in this repo during the SB migration work: adding `+pushed:>2026-01-17` to scan queries did not narrow results to recent repos.

**How to apply:** If the user wants date-filtered results on code search, either:
- Do repo-level filtering first (query `type=repositories` with `pushed:>...`), collect repo list, then scope code search to those repos via `repo:owner/name`.
- Or filter post-scan by inspecting repo metadata (last push date) after finding hits.

Don't try `pushed:` in a code search URL — it won't work.
