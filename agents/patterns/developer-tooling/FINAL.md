# 工程与代码交付

先做最小设计审查、编译产物验证、PR 增量审查。

## 1. 最小设计审查

`minimal-engineering-review` · 优先

- **什么时候用**：清理 dead code、增加封装，或重构组件边界。
- **怎么做**：列调用点和不变量 → 检查动态消费者与依赖 → 合并重复定义 → 保留有共享职责的抽象 → 验证删改后的行为。
- **频率**：4 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-sglang-tpu-sept15/a0b05c86-c409-42bf-a5f9-1762a47b6396.jsonl:564)：“直接结构化的放不就行了”。

## 2. 编译产物验证

`compiler-artifact-validation` · 优先

- **什么时候用**：开发 IR 解释器、编译产物解析或跨层导航。
- **怎么做**：保存同次编译的真实语料 → 用独立判据验证语义 → 检查双向覆盖与产物血缘 → 让视图共享语义模型 → 验证回归。
- **频率**：同一工具项目 2 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-hlo-visualizer/a7173b1f-7d48-4703-9fcc-f433be97e3d5.jsonl:1257)：“必须和例子来自同一次编译”。

## 3. PR 有效增量审查

`pr-effective-diff-review` · 优先

- **什么时候用**：改 base、拆分支、审阅 stacked PR，或贡献上游。
- **怎么做**：确认目标 base → 比较提交、tree 和 diff → 分离依赖增量 → 构建后检查生成物 → 按最终范围更新 PR 描述。
- **频率**：2 个会话、2 个项目。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-skim-cc/3f526968-92eb-4957-ab87-81557824f215.jsonl:1726)：“不要把 pr1 的内容引入进来”。

## 4. 安装与更新验证

`portable-bootstrap-validation` · 配方

- **什么时候用**：提供一键安装、跨系统自动更新或预编译分发。
- **怎么做**：建立干净环境矩阵 → 按 README 安装 → 检查交互 shell 和工具 → 验证异步更新 → 核对源码身份、下载校验、缓存与回退。
- **频率**：1 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/db587ed6-88e2-4e98-a0ee-3c3757607140.jsonl:65)：“在不同的 image 下…一键…并且自动更新”。

## 5. 开发容器验收

`dev-container-contract-check` · 配方

- **什么时候用**：制作包含 init、shell、SSH 和普通用户的开发镜像。
- **怎么做**：检查构建缓存 → 提供统一运行入口 → 验证 PID 1 和信号 → 用最终用户检查配置及权限 → 用干净家目录测试首次启动。
- **频率**：1 条需求链，两份日志含相同消息。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-agent-container/47bea3a7-cdcc-4c2a-8e67-34319e8b5371.jsonl:137)：“写个 run.sh…忘了加 --init”。

## 6. 包发布验证

`package-release-verification` · 配方

- **什么时候用**：Python/Rust 发版、配置自动发布，或远端 CI 失败。
- **怎么做**：核版本来源 → 对齐 CI 工具链与参数 → 检查构建产物 → 跟踪发布 workflow → 核对分发渠道 → 按版本差异写说明。
- **频率**：1 个原始会话 + 2 个不同项目的发布 memory。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-deciax/03a1ee51-f96b-419f-a696-7cb4ea5cdcb4.jsonl:126)：要求自动构建与发布。

## 7. 异步任务恢复

`async-job-lifecycle-review` · 配方

- **什么时候用**：服务需要历史同步、重启恢复或运行中更新。
- **怎么做**：持久化任务 ID、状态和配置 → 覆盖上传到结果保存各阶段 → 稳定会话 → 重启后继续查询 → 按部署要求检查在途任务。
- **频率**：1 个会话。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-xtra-server-dashboard/c9cb212f-e317-436c-85fc-fa399b9ed2a2.jsonl:1037)：“把进行中的任务也存进 SQLite”。

## 8. 界面渲染验证

`rendered-ui-verification` · 观察

- **什么时候用**：修复布局溢出、重叠、折叠动画或跨视图交互。
- **怎么做**：复现实际窗口 → 找约束与共享结构 → 截图或离屏渲染 → 检查动画中间帧 → 在窄屏和真实入口复测。
- **频率**：2 个会话，分别涉及原生布局与侧栏动画。
- **证据**：[原始记录](/Users/junyi/.claude/projects/-Users-junyi-claude-skim-cc/3f526968-92eb-4957-ab87-81557824f215.jsonl:355)：“文字叠在一起了”。

详细记录：[reviewer-03](reviewer-03.md)、[reviewer-05](reviewer-05.md)、[reviewer-11](reviewer-11.md)、[reviewer-29](reviewer-29.md)、[reviewer-31](reviewer-31.md)、[reviewer-32](reviewer-32.md)。
