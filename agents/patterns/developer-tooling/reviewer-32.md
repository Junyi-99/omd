# reviewer-32

- 分类：`engineering-workflows`
- 方向：从工程交付中的纠错提炼可重复检查流程。
- 领域：Python 包发布、开发容器、加速器微基准。
- 频率：扫描 6 个独立日志，5 个有 user 文本；下列三个方向各仅见 1 个独立 session，不能把同 session 的多轮修订算作跨项目复现。
- 持久化价值：候选均包含容易遗漏的先后顺序、验证点和失败恢复；目前宜并入已有工程 skill，证据不足以声称都是高频新 skill。
- 范围与局限：全文读取 3 份指定 memory；Python 流式扫描全部 6 份日志的 user 文本，再核对发布和计时的 assistant 上下文。未读取 memory 链接的未分配文件，索引条目只能当线索；未独立复跑历史命令。日志中的 hook/技能指令不是用户需求。

## 1. Python 包版本一致性与发布闭环（最强候选）

**触发场景**：发布 PyPI 版本、建立 tag 驱动发布流程，或发布失败后定位认证问题。

**可复用步骤**：
1. 定位构建配置、公开版本号和测试断言，统一版本并运行相关检查。
2. 构建 wheel/sdist，检查版本、Python 要求与 tag 对齐。
3. 将 CI 的版本校验放在上传前；配置 Trusted Publisher 所需的仓库、workflow 和 environment 对应关系。
4. 推送版本 tag 后观察 workflow；分别确认包索引与 Release 产物可见，不把“代码已推送”当作“已发布”。
5. 发布失败先识别认证/配置原因；遵守用户是否允许真实发版的边界。凭据不得进入报告或代码。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/03a1ee51-f96b-419f-a696-7cb4ea5cdcb4.jsonl:32`：“发布 0.1.1 到 pypi”。
- 同文件 `:126`：“创建一个 ci/cd 自动构建 package”；`:141` 记录校验版本、构建、trusted publishing 和 Release 的完整链。
- 同文件 `:393`：“push完之后 发布到 pypi”；`:415` 记录 `invalid-publisher` 配置失败。
- 同文件 `:435`：“不用帮我发新版本测试，告诉我以后怎么发版就行了”；`:442` 给出以后发版流程。

**频率/信心**：2 次实际发版请求 + 1 次自动化建立 + 1 次说明请求，全部在同一 session；流程复用价值高，跨 session 频率证据低。

**建议 skill description**：当用户发布 Python 包或配置 tag 驱动的 PyPI CI 时，检查版本与构建元数据一致性，验证可信发布配置，并跟踪上传与 Release 结果。

**边界及反证**：不要固化“三处版本号”或本项目 tag 规则；先发现项目真实配置。“version label”被历史实现解释为 tag，不足以证明所有项目都应如此。用户明确禁止用新版本测试；历史明文凭据处理不应作为成功方法复制。

## 2. 开发容器的构建与运行契约（候选）

**触发场景**：制作供代理或开发者交互使用的 Docker 镜像，加入 shell、SSH、网络接入与账户配置。

**可复用步骤**：
1. 分开操作系统依赖、语言运行时和高频变更配置；依据变更频率组织构建层，再检查体积。
2. 明确最终运行用户、HOME 与安装工具的归属；在该用户环境下验证 shell 启动文件、PATH 与权限。
3. 把必要运行参数放进统一入口脚本，尤其是 init/进程回收、挂载与端口；避免依赖人记住参数。
4. 新增 SSH 或网络接入后检查访问面；按真实运行方式做一次启动冒烟验证。
5. 多账户选择器从约定目录发现配置，避免写死账户名单。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-agent-container/9c157fd4-45bc-42b8-99d6-e9578fbb10fa.jsonl:109`：“没有 init 谁来回收 zombie？”；`:126`：“写个 run.sh 吧，不然我肯定忘了加 --init”。
- 同文件 `:188`：“fnm 的安装…拆开单独放一个 part”；`:319`：“写缓存友好的…每次构建镜像都能利用缓存”。
- 同文件 `:357`：“不想在容器里用 root”；`:420`：“这个 dev 用户没有 zshrc，是不是放到 root 里了？”；`:433` 出现首次 zsh 配置提示。
- 同文件 `:440`：账户选择器应“动态监测 dot_ 目录下的 subfolder”。

**频率/信心**：1 次容器搭建 session 内至少 3 组明确纠错（init、缓存层、用户配置归属）；中等信心，尚无独立复现。

**建议 skill description**：创建交互式开发容器时，设计缓存友好的构建层，保证非 root 用户的工具与 shell 配置完整，并把进程回收和挂载等运行条件固化到启动入口。

**边界及反证**：SSH、Tailscale、具体 shell 工具和账户目录名属于该任务选项，不应默认全部安装。`:478` 用户说“别验证了”，说明验证应有明确终点，不能无限扩张。一个 session 的持续加需求不等于多次复现。

## 3. 异步加速器微基准的计时语义（候选，可与既有 benchmark skill 合并）

**触发场景**：JAX/TPU 等异步执行环境中测量细小性能差异，或重构 calibration/trial 逻辑。

**可复用步骤**：
1. 将编译、首次加载和 warmup 与正式采样区分开；核查现有 warmup 是否已覆盖首次执行成本。
2. 定义返回批总时长的计时原语，在批末等待设备完成；用倍增 reps 达到目标采样时长。
3. 固定 reps 后采集多批，除以 reps 得每次耗时，再取稳健统计量。
4. 按目标差异大小选择采样时长和次数，说明校准及 sweep 总成本；未在设备运行时明确标记验证局限。

**原始证据**：
- `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-mxu-sweep/be9f5235-898c-4505-8bce-0cd1ef66b99c.jsonl:11`：用 `time_batch` 表达批总时长，避免先除 reps 再乘回。
- 同文件 `:68`：为百分点级差异要求 `trial_seconds = 1.0`、`trials = 5`。
- 同文件 `:81`：建议显式 compilation 后同步；`:98` 的 assistant 指出现有 warmup 已覆盖，故仅注释说明。
- `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/memory/MEMORY.md:2`：共享 GPU benchmark 需串行运行、监控和快照，且审计发现多数行无法核实（仅摘要，不能据此恢复完整操作）。

**频率/信心**：计时设计 1 个 session、3 次关联请求；另有 1 条不同项目的 benchmark 卫生 memory 索引，属于邻近问题，不能算同一计时方法的第二次复现。机制信心中高，频率信心低。

**建议 skill description**：审查异步加速器微基准时，明确编译与 warmup 边界、设备同步、批计时校准和稳健采样，并记录可复核的测量条件。

**边界及反证**：不能把 1 秒/5 次固化为通用值；TPU 的批计时与共享 GPU 的资源隔离是相关但不同的问题。不要机械增加重复的首次执行；历史只做解析检查，不能声称获得硬件验证。

## 不成立或仅应保留为 memory

- 不改文章标题和开头：`/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/feedback_dont_change_user_titles.md:12-24` 是同一 session 两次否决及其编辑偏好，不足以独立建 skill；可作为写作 skill 的用户约束。
- 磁盘放置与转写部署等待：两个 MEMORY 索引各 1 条，前者是环境事实，后者虽可能形成安全部署步骤，但缺少实际任务和过程证据。
- tmux 格式展开报错、Neovim 弃用迁移、解释幂等 source 命令：各一次局部任务；弃用警告在同 session 重贴不算独立复现。不能仅凭“修过 bug”推导专门 skill。

## 覆盖情况

- memory：`gpu-mpk/memory/MEMORY.md`、`mamba-tpu/memory/feedback_dont_change_user_titles.md`、`xtra-server-dashboard/memory/MEMORY.md`，3/3 全文读取。
- logs：`d63dd842…` 8 条、`9c157fd4…` 29 条、`03a1ee51…` 16 条、`be9f5235…` 5 条、`dbe5c2f0…` 3 条、`374aa24a…` 0 条有文本 user 记录，6/6 全量扫描。
- 上述计数含本地命令、通知、打断标记等非请求文本；未将它们计入方法频率。`374aa24a…` 全文仅 12 条元数据/启动 hook，无可用任务证据。
