# reviewer-31

- 分类：`engineering-workflows`
- 方向：提炼能减少 agent 返工与错误验收的工程方法；优先最小设计、按真实状态验收。
- 领域：Python 数据建模、开发容器、TPU 资源运维、Rust 发布。
- 频率：5 个独立日志 session、3 个 memory 文件；其中最小设计在 2 个独立 session 出现。其余多为单 session 方法候选或跨项目 memory 旁证，不能视为高频惯例。
- 持久化价值：把“先证明必要性、再抽象”和“验证最终运行环境/终态”写成决策步骤，可以避免相同类别的过度设计与假成功；项目名、版本、账号及命名偏好不应进入通用 skill。
- 范围与局限：全文读取全部 memory；Python 流式解析全部日志并扫描全部 user 文本，再选择 assistant 上下文核实。部分 user 文本是粘贴旧回答、系统通知、重复提交；不计作独立复现。未验证历史模型清单、软件版本或云端当前状态。

## 1. 先消除结构性约束，再证明抽象的必要性（推荐候选）

- 触发：用户认为 wrapper、校验、辅助模块过度设计，或数据模型需要同步两份键集合。
- 可复用步骤：
  1. 找到需要维护的真实不变量与全部调用点，区分转换语法、共享业务规则、依赖隔离。
  2. 对“两个平行表必须同键”优先改成单一嵌套结构，删除结构已保证的校验。
  3. 对薄 wrapper 尝试内联；若它承载共享文件清单或阻止重依赖泄漏，则保留最小模块。
  4. 保留结构无法保证的跨条目规则，例如全局唯一性；用负向输入确认不会静默覆盖或去重。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/a0b05c86-c409-42bf-a5f9-1762a47b6396.jsonl:564`：“你直接结构化的放不就行了？”；`:580` 将厂商 URL 与 models 合并到 Vendor，删掉对键校验；`:603` 解释 YAML 重复 key 静默覆盖的反例。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/3050adc8-9a3a-4210-b3d7-853e29b482f7.jsonl:9`：“感觉过度设计了”；`:68` 再问能否调用处包 Path；`:83` 更新后的两调用点共享 FILES，直接内联会复制清单或引入 jax 依赖。
- 频率/信心：2 个独立 session；至少 3 次明确设计质疑，但同 session 追问不增计复现。对用户需求高信心，对具体方案中高信心。
- 建议 skill description：当 agent 新增薄封装或维护冗余数据约束时，追踪调用与不变量，通过结构重设计、依赖分析和负向验证选择最小实现。
- 边界及反证：不能机械“删除 wrapper”。deciax 在 git pull 后出现真实共享职责，旧时无调用的判断不再适用；日志恰好证明必须以当前代码为准。Python 版本偏好不是方法本体。

## 2. 云资源重试以身份与终态验收（单 session 候选）

- 触发：跨 zone 等待容量，后台创建持续失败，用户反复查询状态或希望仅变化时通知。
- 可复用步骤：
  1. 明确账号、项目、资源规格、计费模式与候选区域，避免把 spot 和 on-demand 混用。
  2. 用有上限、可观察的重试进程；保留最近错误、成功计数、进程存活状态。
  3. 将容量不足等常态与权限/配额新错误、进程退出区分，只对有意义变化通知。
  4. 通知判据同时匹配目标资源身份与成功终态；核对创建命令结果，不能把列表中的任意资源或 CREATING 视为建成。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-spot/67278d73-9bd7-42ab-ae34-cb79f0ba57b7.jsonl:55`：“我要 on-demand”；`:238` 询问账号；`:301`：“改成有变化通知我”。
  - 同文件 `:343`：误报其他名称的 TPU；`:368`：“节点短暂进入 CREATING”，随后容量不足被拆掉，SUCCESS 仍为 0。
- 频率/信心：1 个独立 session，2 次不同原因误报；多次“现在咋样”只是一条长任务。失败判据证据强，修复落地证据不足。
- 建议 skill description：等待云资源容量时，建立可恢复的重试与变化监控，并用精确资源身份及最终就绪状态确认成功。
- 边界及反证：日志结尾仍未取得资源，不能宣称此流程已成功验证；不得把某区域/配额、账号或项目固化。mamba memory 第 7 行只是 spot 恢复索引，并非第二次完整案例。

## 3. 开发容器按真实入口与用户做运行验收（单 session 候选）

- 触发：构建包含交互 shell、SSH、Node 工具和普通用户的开发镜像，且要求重建快、体积小。
- 可复用步骤：
  1. 构建并运行验证；按工具安装与配置的变化频率安排层，外置包管理器缓存。
  2. 分别测 PID 1、停止信号与子进程回收；把必需运行参数写进启动脚本。
  3. 对本地 shell 与 SSH 两条入口，分别验证实际普通用户、PATH、HOME 和启动文件。
  4. 用全新家目录做首次启动测试，避免已有文件掩盖初始化提示；在改配置后测相关功能。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-agent-container/47bea3a7-cdcc-4c2a-8e67-34319e8b5371.jsonl:137`：“写个 run.sh 吧，不然我肯定忘了加 --init”；`:133` 比较 PID 1/--init 场景。
  - 同文件 `:409`：“写缓存友好的”；`:431` 描述 BuildKit cache mounts 与改 entrypoint 后的缓存观察。
  - 同文件 `:559` 询问 dev 的 zshrc；`:581` 提供真实首次启动提示；`:584` 确认无用户启动文件会触发 Debian 初始化逻辑。
- 频率/信心：1 个独立 session；多个相关迭代，不能算多次独立复现。运行验收方法中高信心。
- 建议 skill description：构建开发容器时，兼顾可复用构建缓存，并用真实登录入口、普通用户及首次启动场景检查运行行为。
- 边界及反证：先前 assistant 曾以函数可用断言无需用户 zshrc，随后真实启动暴露问题，因此不能只检查配置文件位置。固定工具清单和选账户 UI 属单项目需求。

## 4. 发布前镜像 CI 条件，发布后核对工作流（memory 支持候选）

- 触发：Rust 项目发布或本地检查通过而远端 CI 失败。
- 可复用步骤：读取实际 workflow 与工具链版本；运行同参数检查及必要格式检查；分离功能提交与版本提交；按项目规则发 tag；追踪对应运行结果；按上一 tag 之后的用户可见变化编写说明。
- 原始证据：
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-deck/memory/deck-verify-with-ci-clippy-flags.md:15`–`:23`：旧本地工具链且缺少 `-D warnings` 导致“三次 pushes”失败，要求推后检查 workflow。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-vibe-usage/memory/feedback_release_flow.md:9`–`:22`：版本/lock 同步、tag、等待 release workflow、说明依据 tag diff。
- 频率/信心：2 个不同项目 memory，各标识 1 个来源 session；本批未分配来源原日志，不能将“三次 pushes”算成 3 个 session。流程中信心，用户具体发布偏好仅单来源。
- 建议 skill description：发布 Rust 项目时，复现 CI 的工具链及检查参数，跟踪 tag 对应工作流，并基于版本差异整理发布说明。
- 边界及反证：默认 patch、直推 main、移动 tag 是项目政策，不能泛化授权；deck 明确 release 与普通 CI 独立，不能把 release 成功等同主分支 CI 通过。

## 未提升为 skill

- 模型目录维护：1 个 session，有从厂商 HF /models 页面更新、排除微调、保留完整模型路径等具体要求；缺乏第二次独立执行证据，且日期/清单是时效数据。
- genesis 自指分层：同一需求在一个 session 重复两次（日志 :11/:32），没有完成或重建验证，不能据此认定幂等方法有效。
- mamba MEMORY 的写作保护、prefill 术语、数值溢出和 spot 区域信息：仅索引，未分配其链接正文；属于偏好/事实/项目状态，不足单独建立通用方法。

## 覆盖

- 全文读取 3/3 memory：deck CI、mamba MEMORY、vibe-usage release flow。
- 流式扫描 5/5 JSONL 全部 user turns：a0b05c86（模型清单）、47bea3a7（容器）、67278d73（TPU）、d5c59458（genesis）、3050adc8（deciax）；对上述证据附近 assistant 文本回读。
- 未读取其他 reviewer 结果；未创建 skill，未变更其他文件。
