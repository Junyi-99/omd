---
name: Commit-push-review workflow
description: User prefers commit+push after each logical phase, with codex review between phases
type: feedback
---

After completing each logical phase of work, commit and push immediately. Use /codex for code review between phases before proceeding to the next.

**Why:** User wants incremental delivery with review gates. Aligns with CLAUDE.md rule "After every successful build: Commit all changes and push to remote immediately."

**How to apply:** After each phase/milestone: (1) build to verify, (2) commit+push, (3) /codex review, (4) fix issues if any, (5) proceed to next phase.
