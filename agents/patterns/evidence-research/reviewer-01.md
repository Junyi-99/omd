# reviewer-01 独立审阅

分类：`research-backed-visualization`
方向：把论文、照片和官方文档转成可追溯的交互式架构学习工具。
领域：硬件架构资料整理、结构化规格、分层可视化。
频率：仅 1 个独立 session。全文件扫描得到 42 条含文本的 user 记录，其中含命令回显、图片路径、3 条 task-notification 和 1 条注入的 XLSX skill；这些都不当成独立需求或复现。日志中以下次数是同一项目内的需求锚点，不代表跨项目习惯。
为什么值得持久化：可复用的是“证据 → 数据模型 → 可视化 → 验证”的方法，尤其是区分硬件层级、原始数字和推断；不是 TPUv2 的具体参数。证据强度足以提出候选，但不足以断言用户反复需要完整 skill。
审阅范围与局限：manifest 含 4 个 memory 和一个约 45 MB JSONL；已流式扫描全文件的全部 user 文本，按主题读取 assistant 结果和过程。没有运行历史代码，assistant 的成功报告仅为历史声称，不视为本次验证。未读取其他 reviewer 结果。

本文证据前缀 L 均指绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-architecture/d88da7fc-99bc-415f-945c-438a4a6ecb51.jsonl` 的 JSONL 行号。

## 1. 多来源规格归一化，保留来源及硬件计量层级（优先候选）

**触发场景**：从论文表格、官方博客或产品文档构建硬件规格、架构图，出现 chip/core/board 等单位混用或不同来源数字冲突。

**可复用步骤**：
1. 先定义 system、board、chip、core、unit 的实体及数量关系，确认图中一个 node 到底指什么。
2. 表格文本抽取错位时渲染原页核对；记录表号、图号、PDF 页码。
3. 将原始值、单位、比较符号、来源分别录入 manifest；推导值放独立字段，不能覆盖原始值。
4. 对跨来源数字换算到相同层级，验算乘积和总量；不一致则保留差异和说明。
5. 让拓扑、详情图、术语及表格消费相同数据源，审阅时检查没有硬编码漂移。

**原始证据**：
- L502 user：“准备一个 manifest……数据就用 paper 的 table3”。同条明确每 host 4 chips、每 chip 2 cores。
- L596 user：“一个 node 应该画成 chip，不能画成 board，也不能画成 TensorCore”。
- L860 user：“数字加进去。并且说清楚单位差别”。
- L506 assistant：“文字提取出来的表格行列有点乱”；转为渲染 Table 3 和 Fig. 3 核对。
- L556 assistant：“torus size and link speed come from one place”。
- L880 assistant 报告将博客原值、derived 和 vsTable3 分开保存，不覆盖论文数字。

**频率/信心**：3 个明确 user 锚点（L502、596、860），均属 1 session 的连续校正；对方法需求信心高，对跨项目复现信心低。

**建议 skill description**：当需要把多个技术来源转成硬件规格或架构可视化时，建立带来源的规格 manifest，明确实体层级和单位，核对表格、换算与冲突，并驱动展示复用同一数据。

**边界及反证**：不能把 TPUv2 数字或云实例命名规则泛化到所有代际；“<611 mm²”不能静默当作精确面积。历史实现使用面积上界绘图，需要明确标示。用户没有在第二个独立项目复现这一需求。

## 2. 带定位引用、稳定标识和多来源释义的术语库（候选）

**触发场景**：阅读技术资料时建立可继续扩充的术语页面，需要从定义回到原始证据，并避免同义或历史名称造成重复。

**可复用步骤**：
1. 全文识别术语，按硬件、互连、数值、软件等组织；通用前置知识单独标记，避免硬件词汇表无限膨胀。
2. 保存短释义及来源定位；分清逐字引用、删节与综合释义。
3. 为术语生成全局唯一、稳定的 slug，并校验重复。
4. 同一术语的新来源追加到现有条目，保留版本及语境，不新增重复词条。
5. 验证引用可打开且页码准确；如需深链，检查路由与锚点是否冲突。

**原始证据**：
- L132 user：“专门写 Terminology”。
- L143 user：“这一个 statement，加一个 reference 到 [本地论文]”。
- L172 user：“提取出所有的 terminology”。
- L220 user：“每一个 terminology 给一个 slug。slug 全局唯一”。
- L217 assistant 报告 57 条、7 类、逐条页码，并承认部分定义合并了跨段内容。
- L230 assistant 报告唯一性检查，同时指出 hash 路由导致无法直接跳转术语。
- L923 assistant 将同一术语多来源释义放入 `more`，并拒绝混入与新代际 MXU 尺寸矛盾的通用句子。

**频率/信心**：4 个明确 user 建设步骤（L132、143、172、220），1 session 内同一功能的渐进完善；L883 是继续补资料，不另算独立复现。候选信心中等。

**建议 skill description**：当用户要从论文或技术文档建设可扩充术语库时，提取分组术语，为定义附定位引用和唯一 slug，保留多来源语境并验证链接及标识。

**边界及反证**：这可能只是该学习网站的一个功能；不应因出现术语便自动建库。不要把一次抽出的 57 个词沉淀为 skill，也不要声称综合改写等于原文。

## 3. 证据约束的分层架构可视化与语义缩放（候选）

**触发场景**：用户想从拓扑进入板卡、芯片和微架构，且物理布局与逻辑结构需要不同展示方式。

**可复用步骤**：
1. 先搭空框架，允许按学习进度增量补资料；先确定层级及视图任务。
2. 物理视图参照照片/版图，记录比例标定及示意尺寸；逻辑视图表达功能和连接关系。
3. 为来源已知和推断的连接分别编码，推断附假设说明；数量关系只用于一致性核对，不当作直接布线证据。
4. 微架构采用可缩放的层次模型，根据屏幕尺寸展开下一层，裁剪视野外结构，避免一次渲染全部重复单元。
5. 提供可继续放大的明确线索、常驻标签与 hover 解释；逐层测试缩放、拖动、切换、窄面板及控制台错误。
6. 来源未提供更深细节时停止，保留空白而非补造 datapath。

**原始证据**：
- L13 user：“先不要往里面放内容……一步步往里面放东西”。
- L303 user 要求选中拓扑节点后右侧显示另一个 3D 详情；L604 要求在一定缩放时显示 label。
- L680 user：“允许切换 physical view / logical view”。
- L1017 user 要求 PCIe 从边缘接口出线，并询问 CDFP 是否对应出口；L1148 assistant 明确未找到公开证实，标为 inferred。
- L1382 user：“我想要的就是这个效果”，所贴内容描述 infinite canvas + semantic zoom。
- L1501 assistant 报告按屏幕宽度展开、仅展开视野内结构，逐层截图，并明确不编造更深 datapath。
- L1564 user：“内部还有东西的……用紫色框表示。这样我知道可以 zoom in”。

**频率/信心**：5 个不同可复用交互需求锚点（L303、604、680、1382、1564），仍为 1 session 的一条产品演进链；不能算 5 次独立工作流。方法候选信心中等，独立 skill 的必要性待其他 session 佐证。

**建议 skill description**：当用户建设可逐层探索的技术架构学习工具时，分离物理和逻辑视图，以有来源的层次数据驱动拓扑、详情和语义缩放，清楚标记推断并验证逐层交互。

**边界及反证**：具体紫色、260px 阈值及 SVG/3D 选型是当前实现选择，不是固定规范。推断出的板内拓扑不等于真实 PCB 走线。该模式可能过大，适合先合并到“技术资料驱动架构可视化”工作流，而非独立创建通用前端 skill。

## 4. 共享 GPU 性能测试的排他性取证（来自 memory 的独立候选）

**触发场景**：共享 GPU 主机上的延迟/性能 benchmark 需要证明结果没有受其他用户或自己的并行任务干扰。

**可复用步骤**：每主机串行执行；启动前确认上个进程已释放显存；持续轮询全部 GPU，空闲也写 heartbeat；每轮开始和结束保存 GPU 进程及 PID/用户/参数快照；日志追加且每行带时间；CSV 行关联 run_id 和起止时间；只有证据覆盖完整窗口才标记 clean。

**原始证据**：绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/memory/shared-gpu-benchmark-hygiene.md` L11–18 记录一次审计无法证明数据干净，尤其“silence proved nothing”；L20–27 给出上述取证流程。

**频率/信心**：1 份反馈 memory，记载 1 个独立 origin session；未分配原始日志，不推断发生次数。方法可复用性高，历史核实强度中等。

**建议 skill description**：在共享 GPU 主机运行性能基准时，串行调度并持续记录资源占用，以每轮证据关联输出结果，识别受干扰的测量并避免无证据宣称排他运行。

**边界及反证**：仅适用于共享性能测量，不能要求普通开发任务都独占主机；检测不到进程不等于没有其他干扰源。不可把 memory 中“多数行不能证明干净”改写成“多数行被污染”。

## 未提升为 pattern 的事项

- 图片重命名只出现 L927 一次；3 条代理完成通知是该次工作的子结果，不能当作 3 次复现。适合资产整理步骤，暂不足以独立成 skill。
- XLSX 创建 L1505 以及 cooling system 更新 L1600/L1604 是一次表格功能及重复澄清；不据此推断常见 Excel 工作流。L1518 是工具注入 skill 文本，不是用户习惯证据。
- Prefill 术语 memory 是明确的文档范围纠偏，支持 pattern 1 的语义精度原则，但“禁用 serving”的项目规则应留在 memory；外部引文本身可以保留该词（该文件 L26–28）。
- 部署等待 memory（绝对路径 `/Users/junyi/.claude/projects/-Users-junyi-claude-xtra-server-dashboard/memory/audio-transcribe-deploy-wait-for-jobs.md` L11–15）给出一次正在转写时重启的教训。可抽象为发布前检查实际活动任务、准备改动后等待并最终复查，但只有一次且绑定具体服务/DB，不独立升为 skill。不要将本地可恢复测试误作当前任务可中断的授权。
- React/Tailwind/TypeScript 和 ponytail 只在开场指定一次，属于项目约束，不把技术栈偏好变成方法。

## 覆盖情况

- Memory：下列 4 个文件均已全文读取：
  - `/Users/junyi/.claude/projects/-Users-junyi--oh-my-dotfiles/memory/MEMORY.md`（1 行）：索引及一次历史合并状态，不是 skill 证据。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-gpu-mpk/memory/shared-gpu-benchmark-hygiene.md`（29 行）。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-mamba-tpu/memory/feedback_prefill_terminology.md`（28 行）。
  - `/Users/junyi/.claude/projects/-Users-junyi-claude-xtra-server-dashboard/memory/audio-transcribe-deploy-wait-for-jobs.md`（15 行）。
- Log：`/Users/junyi/.claude/projects/-Users-junyi-claude-tpu-architecture/d88da7fc-99bc-415f-945c-438a4a6ecb51.jsonl`：已流式解析全部记录并扫描全部 42 条含文本 user 记录；按主题补读 L140–238、500–603、840–926、1017–1174、1382–1504、1564–1609 的 assistant 文本。
- 仅写本报告；未创建 skill、未移动文件、未读取其他 reviewer 结果。
