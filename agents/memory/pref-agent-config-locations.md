# Agent 工具配置加载位置说明

## 问题概述

Cloud Code、Cursor 等 AI agent 工具会从哪些文件夹加载自定义配置？如果使用 XDG 重定向，它们会使用重定向后的文件吗？

---

## 1. Cursor 配置加载机制

### CLI 配置 (`cli-config.json`)

**标准位置：**
- Linux: `$XDG_CONFIG_HOME/cursor/cli-config.json` (默认: `~/.config/cursor/cli-config.json`)
- macOS: `~/Library/Application Support/Cursor/cli-config.json` 或 `$XDG_CONFIG_HOME/cursor/cli-config.json`
- Windows: `%APPDATA%\Cursor\cli-config.json`

**XDG 支持：**
- ✅ **Cursor CLI 支持 XDG 标准**
- 如果设置了 `XDG_CONFIG_HOME`，CLI 会从该路径读取配置
- 你的配置位置：`$XDG_CONFIG_HOME/cursor/cli-config.json` → `/Users/junyi/.oh-my-dotfiles/etc/cursor/cli-config.json`

### Agent Rules 和 Skills

**Cursor Agent 的自定义规则和技能：**

1. **Rules (`.cursor/rules/` 或 `~/.cursor/rules/`):**
   - 项目级规则：`.cursor/rules/` (在项目根目录)
   - 用户级规则：`$XDG_CONFIG_HOME/cursor/rules/` 或 `~/.cursor/rules/`
   - 如果设置了 `XDG_CONFIG_HOME`，会优先使用该路径

2. **Skills (`~/.cursor/skills/` 或 `$CURSOR_HOME/skills/`):**
   - 默认位置：`~/.cursor/skills/`
   - 如果设置了 `CURSOR_HOME` 环境变量，会使用 `$CURSOR_HOME/skills/`
   - **注意：** Skills 可能不完全遵循 XDG，需要检查 Cursor 的具体实现

3. **Agent Memory (`agent/memory/`):**
   - 这是项目级的，通常在项目根目录的 `agent/` 文件夹
   - 你的项目中有：`agent/memory/pref-dev-tools.md`, `agent/memory/pref-file-organization.md`

### 验证方法

```bash
# 检查 Cursor CLI 配置位置
cursor --config-path  # 如果支持
# 或
echo $XDG_CONFIG_HOME/cursor/cli-config.json

# 检查实际使用的配置
ls -la ~/.config/cursor/  # 默认位置
ls -la $XDG_CONFIG_HOME/cursor/  # 你的 XDG 重定向位置
```

---

## 2. Cloud Code 配置加载机制

### 配置位置

**标准位置：**
- Linux: `$XDG_CONFIG_HOME/cloud-code/` (默认: `~/.config/cloud-code/`)
- macOS: `~/Library/Application Support/cloud-code/` 或 `$XDG_CONFIG_HOME/cloud-code/`
- Windows: `%APPDATA%\cloud-code\`

**XDG 支持：**
- ⚠️ **部分支持** - Cloud Code 作为 IDE 插件，主要遵循 IDE 的配置机制
- 如果作为独立工具，可能支持 XDG
- 需要查看具体版本和实现

### Agent 配置

- Agent 规则和配置通常在 IDE 的设置中
- 项目级配置可能在 `.cloudcode/` 目录
- 用户级配置可能在 `$XDG_CONFIG_HOME/cloud-code/` (如果支持)

---

## 3. XDG 重定向的影响

### 你的当前设置

从 `.xdg.sh` 可以看到：
```bash
export XDG_CONFIG_HOME="${USER_CONFIG_HOME}/etc"
# 即: /Users/junyi/.oh-my-dotfiles/etc
```

### 哪些工具会使用 XDG 重定向？

**✅ 会使用 XDG_CONFIG_HOME 的工具：**
- Cursor CLI (`cli-config.json`)
- 大多数遵循 XDG 标准的 CLI 工具
- Git (如果配置了 `$XDG_CONFIG_HOME/git/config`)
- Neovim (如果配置了 `$XDG_CONFIG_HOME/nvim/`)
- 其他遵循 XDG Base Directory Specification 的工具

**⚠️ 可能不使用 XDG 的工具：**
- macOS GUI 应用（通常使用 `~/Library/Application Support/`）
- 某些 IDE 插件（可能使用 IDE 自己的配置系统）
- 不遵循 XDG 标准的旧工具

### 验证 XDG 是否生效

```bash
# 1. 检查环境变量
echo $XDG_CONFIG_HOME
# 应该输出: /Users/junyi/.oh-my-dotfiles/etc

# 2. 检查 Cursor 配置是否在重定向位置
ls -la $XDG_CONFIG_HOME/cursor/
# 应该看到: cli-config.json

# 3. 检查是否有默认位置的配置文件（可能被忽略）
ls -la ~/.config/cursor/ 2>/dev/null || echo "默认位置不存在，使用 XDG 重定向"
```

---

## 4. 最佳实践建议

### 确保 Agent 工具使用 XDG 重定向

1. **在 shell 初始化时设置 XDG 变量**
   - ✅ 你已经在 `etc/zsh/zshenv` 中设置了
   - 确保所有 shell 会话都加载了这个文件

2. **创建符号链接（如果需要兼容性）**
   ```bash
   # 如果某些工具仍使用默认位置，创建符号链接
   mkdir -p ~/.config
   ln -sfn $XDG_CONFIG_HOME/cursor ~/.config/cursor
   ```

3. **检查工具文档**
   - 查看 Cursor/Cloud Code 官方文档确认配置路径
   - 某些工具可能有特定的环境变量覆盖 XDG

### 项目级 vs 用户级配置

- **项目级配置** (`.cursor/rules/`, `agent/memory/`): 不受 XDG 影响，在项目目录中
- **用户级配置** (`cli-config.json`, 全局 rules): 受 XDG 影响

---

## 5. 实际验证步骤

### 测试 Cursor CLI 配置加载

```bash
# 1. 确认 XDG 环境变量
source ~/.zshenv  # 或重新打开终端
echo $XDG_CONFIG_HOME

# 2. 检查配置文件位置
cat $XDG_CONFIG_HOME/cursor/cli-config.json

# 3. 测试 Cursor CLI 是否读取了配置
cursor --version  # 如果支持，查看配置路径
```

### 测试 Agent Rules 加载

```bash
# 1. 检查用户级 rules
ls -la $XDG_CONFIG_HOME/cursor/rules/ 2>/dev/null || echo "用户级 rules 目录不存在"

# 2. 检查项目级 rules
ls -la .cursor/rules/ 2>/dev/null || echo "项目级 rules 目录不存在"

# 3. 在 Cursor 中测试 agent 是否加载了规则
```

---

## 6. 总结

### Cursor
- ✅ **CLI 配置** (`cli-config.json`): 支持 XDG，会使用 `$XDG_CONFIG_HOME/cursor/cli-config.json`
- ✅ **Agent Rules**: 支持 XDG，用户级规则在 `$XDG_CONFIG_HOME/cursor/rules/`
- ⚠️ **Skills**: 可能不完全支持 XDG，默认在 `~/.cursor/skills/`

### Cloud Code
- ⚠️ **部分支持 XDG**，主要取决于具体实现和版本
- 作为 IDE 插件时，可能使用 IDE 的配置系统

### 你的设置
- ✅ XDG 重定向已正确配置：`XDG_CONFIG_HOME=/Users/junyi/.oh-my-dotfiles/etc`
- ✅ Cursor CLI 配置已在重定向位置：`etc/cursor/cli-config.json`
- ✅ 应该会被 Cursor CLI 正确读取

### 建议
1. 在 Cursor 中测试 agent 功能，确认配置被正确加载
2. 如果发现某些配置未生效，检查是否有工具使用了硬编码路径
3. 考虑为不兼容的工具创建符号链接作为后备方案
