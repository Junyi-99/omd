分类：`portable-bootstrap-validation`
方向：把安装说明变成跨发行版验收，并把编译依赖转成可校验的发布产物。
领域：开发环境安装、二进制分发；另有研究文档编辑和代码审查修复候选。
频率：主要方向仅 1 个独立 session；其他候选也各 1 个 session。没有足够跨 session 证据证明任何一个已经是高频技能。
持久化价值：安装矩阵、异步更新验收和预编译分发包含多步可复用程序及实际失败案例，值得作为候选沉淀；偏好、硬件事实不应独立做 skill。
范围与局限：完整流式扫描 manifest 的 6 个 JSONL，提取全部 user 文本，排除工具结果和系统任务通知；3 个 memory 全文读取。针对候选核对 assistant 上下文。日志报告的测试结果未在本次重跑；同一会话内多次容器执行不算独立需求。

## 1. README 驱动的跨平台安装与自动更新验收（首选候选）

- 触发：用户要求确认陌生机器能否按照 README 一键安装，或检查 shell 启动时的自动更新。
- 可复用步骤：
  1. 确认新克隆实际取得的分支和 README；用干净容器按文档安装，避免宿主已装依赖掩盖问题。
  2. 建立发行版 × 架构矩阵，记录镜像真实架构；TTY 安装交互与预置配置路径分别说明。
  3. 验证命令可执行、插件已加载、交互 shell 无错误；不能只检查安装命令退出码。
  4. 让副本落后一个 commit，验证启动更新、间隔节流、禁用开关和脏目录拒绝更新。
  5. 后台更新使用状态轮询和截止时间；区分测试竞争、模拟器限制与产品失败，再修复并复测最小相关矩阵。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/db587ed6-88e2-4e98-a0ee-3c3757607140.jsonl:65`：“在不同的 image 下，能否一键把 dotfiles 装好。并且自动更新”。
  - 同绝对路径 `:381`：assistant 定位到两个自动更新进程抢锁，断言未等待导致误报。
  - 同绝对路径 `:533`：汇总 Ubuntu、Debian、Fedora、Arch、Alpine；确认 Alpine 硬编码 bash、最小镜像缺 hostname 两个真实问题。
  - 同绝对路径 `:721`：“Alpine ... exits 0 ... Fedora ... no error”，修复后相关镜像复测。
- 频率/信心：1 次独立需求/1 session；同会话包含初测、修复复测和 Go-less 复测。方法证据强，高频判断弱。
- 建议 skill description：当用户要求验证安装文档、跨 Linux 发行版兼容性或 shell 自动更新时，用隔离环境构建验收矩阵，检查实际工具可用性和更新状态机，区分测试竞争与产品缺陷，并复测修复。
- 边界及反证：不把 Docker/Rosetta 下 pacman sandbox 特例写成生产安装要求；不假装人工写配置等于 TUI 已测。Go 安装前置暴露出“一键安装”最初尚未闭环。

## 2. 按源码版本绑定、校验并缓存预编译工具（可并入候选 1 的配方）

- 触发：安装器自身需要编译器，用户希望普通使用者无需安装构建工具。
- 可复用步骤：CI 构建目标 OS/架构资产及校验清单 → 清单携带源码身份 → 启动器按源码身份缓存 → 下载后核对身份与摘要 → 原子落盘 → 本地源码改动、离线或不支持架构时按已定义策略回退 → 在无编译器环境和篡改摘要场景验证。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/db587ed6-88e2-4e98-a0ee-3c3757607140.jsonl:762`：“不是需要 golang 编译 omd 吗”。
  - 同绝对路径 `:769`：“github action 里面自动构建 ... 脚本自动下载 ... 校验二进制”。
  - 同绝对路径 `:772`：assistant 将缓存键和清单绑定 `cmd/omd` git tree，并选 SHA-256。
  - 同绝对路径 `:838`、`:850`：错误摘要拒绝、无临时文件残留；无 Go Ubuntu/Alpine 安装通过，并指出 release 发布存在延迟窗口。
- 频率/信心：1 次/1 session，与候选 1 为同一任务演进，不额外计算独立复现；实现步骤证据强。
- 建议 skill description：当命令行安装器需要消除终端用户的编译依赖时，设计 CI 多平台二进制分发、源码身份绑定、摘要验证、缓存失效和本地构建回退，并验证无工具链环境。
- 边界及反证：同源下载摘要只能验证资产与清单一致，不能当成独立发布者认证；摘要算法选择不必独立建 skill。rolling release 延迟造成“无 Go 首次运行”仍可能暂时失败，技能应要求定义该窗口行为。

## 3. 技术综述去防御性赘语（编辑操作候选，非单纯文风记忆）

- 触发：用户要求批量删掉研究文档中的过度防御、范围声明和“不是 X 而是 Y”赘语。
- 可复用步骤：明确文件范围 → 搜索疑似措辞 → 逐句区分作者实测限制与编辑者自我免责 → 保留正向事实并压缩赘语 → 检查表格、链接、引用及数值 → 第二遍审阅遗漏与语义损失。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-papers/5d61343c-9b64-4605-861c-9754d9c97605.jsonl:9`：“清理一下 ... ‘并非’、‘不代表’、‘未确立’、‘不要求’ ... 过渡防御”。
  - 同绝对路径 `:114`：“不用管 subfolder 里面的 md”。
  - 同绝对路径 `:237`：“再看看还有这种过渡防御的句子吗”。
  - 同绝对路径 `:337`：任务结果报告复核表格列数、引用链接，同时保留“不支持 CUDA Graphs”等系统事实。
- 频率/信心：1 session，2 次审阅请求属于同一轮清理与追查遗漏，不是 2 个独立复现；候选置信度中。
- 建议 skill description：当用户明确要求精简技术综述的防御性表述时，按指定范围逐句压缩，保留影响结论的事实限制、证据等级及引用结构，并检查编辑后语义。
- 边界及反证：本记录后续清理也删除过 reviewer “未知”说明和证据可访问性信息，存在损失证据透明度的风险，不能把“所有未知都删”固化。仅“用户喜欢直接表达”应存风格偏好，不是新技能。

## 未升级的线索

- PR 回复驱动修复：`/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/9b992e11-8208-4303-b833-87ad152a2e12.jsonl:9,195,200,303,329` 有“读 review → 最小改动 → 与改前失败集比较 → allocator probe → 例子解释”的良好方法，但只有 1 session；可补入现有代码审查/调试技能，暂不新建。不在 commit 中写 session 属偏好。
- PDF 评论定位源文件：`/Users/junyi/.claude/projects/-Users-junyi-claude-skim-cc-examples-demo-paper/91912def-80cd-4a56-8bd0-365f5d42de27.jsonl:5,41,51` 包含选中文本优先、SyncTeX 辅助、只改源文件的自动化协议；后三次携带同一评论历史，仅 1 次真实编辑，不能算 3 次复现。适合产品既有 worker 协议，不据此另做个人 skill。
- Profile 展示：`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep-results/4be343d5-a57b-4ba4-8a9d-089090298527.jsonl:9,73,88` 从单 run 改为结果根目录以显示全部 experiments，只有 1 session；是简短操作配方，尚不足以独立建 skill。
- deck memory 的直接推 main 是项目授权偏好；TPU v7 memory 是 API surface 与当时配额事实，应保留带日期参考记忆，不能泛化成跨项目权限或长期可用性结论；wiki memory 只是索引。

## 覆盖情况

全部文件均已扫描；以下路径均为 manifest 原文件：

- `/Users/junyi/.claude/projects/-Users-junyi-claude-deck/memory/deck-push-direct-to-main.md`：全文 14 行。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-llm-wiki/memory/MEMORY.md`：全文 3 行；没有沿索引读取未分配文件。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/memory/tpu-v7-is-compute-machine-type.md`：全文 23 行；未复制项目标识及配额详情。
- `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/db587ed6-88e2-4e98-a0ee-3c3757607140.jsonl`：全部 user turns；核对安装、更新与发布总结。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-papers/5d61343c-9b64-4605-861c-9754d9c97605.jsonl`：全部 user turns，含第二次清理；任务结果只作为执行上下文。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-kimi-k3-kda/9b992e11-8208-4303-b833-87ad152a2e12.jsonl`：全部 user turns；核对测试结果与参数解释。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep-results/4be343d5-a57b-4ba4-8a9d-089090298527.jsonl`：全部 user turns 与展示结果。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-skim-cc-examples-demo-paper/91912def-80cd-4a56-8bd0-365f5d42de27.jsonl`：全部三次评论载荷，区分历史重放与新请求。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-skim-cc/8c02c41a-5dec-476f-a509-e9608a614f49.jsonl`：全部扫描，仅本地 `/clear`，没有实质任务。
