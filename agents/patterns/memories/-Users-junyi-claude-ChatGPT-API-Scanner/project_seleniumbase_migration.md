---
name: Pending migration from selenium to SeleniumBase UC mode
description: Next planned direction for the Chrome driver — replace raw selenium with SeleniumBase UC mode to handle Google's anti-automation login block
type: project
originSessionId: 684260bd-b5d2-47ef-a332-8f3b3ef464cd
---
Next step after uv migration: switch the Chrome driver in `src/main.py` from raw `selenium.webdriver.Chrome` to **SeleniumBase UC mode** (https://github.com/seleniumbase/SeleniumBase).

**Why:**
- Google blocks plain selenium on GitHub login with "This browser or app may not be secure"
- `undetected-chromedriver` was the historical fix but has been unmaintained since 2024-02-17 (last release 3.5.5) and breaks on Chrome 147+ with `NoSuchWindowException: target window already closed`
- We rejected PR #26 which tried the uc route for exactly this reason
- SeleniumBase is actively maintained and its UC mode tracks new Chrome versions

**How to apply:**
- Replace `selenium` dep with `seleniumbase` in `pyproject.toml` `[project].dependencies`
- Rewrite `APIKeyLeakageScanner.login_to_github` in `src/main.py` to use SeleniumBase's UC-mode driver construction (see https://github.com/seleniumbase/SeleniumBase/blob/master/help_docs/uc_mode.md)
- Keep the existing cookies.pkl flow (first-run manual login, subsequent runs reuse cookies)
- Verify before committing: real interactive login on Chrome 147+
