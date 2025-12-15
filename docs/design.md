# AI Comic Novel Studio：在 Coze Studio 节点体系上构建“创作工作流”

本文档的最终目标：使用 **coze-studio 现有的工作流（Workflow）和节点（Node）系统**，搭建面向小说/剧本/视频分镜（含音频）的灵活创作流程；同时让每个工作流与资源（图片/文本/列表/音频/视频）强绑定，并提供可复用的资产管理能力。实现上尽量复用现有代码结构与节点逻辑，避免“另起炉灶”。

关键约束：**不要把任何一个节点当成特殊案例**。世界观、角色、分镜等“创作管理节点”应具备相同的待遇：持久化、可复访、可引用、可产出/关联资产、可被其他节点消费。

---

## 一、你想要的“创作工作流”长什么样

### 1.1 核心体验

- 用户在一个工作流画布里自由添加“创作管理节点”（世界观/角色/分镜/音频方案等），像搭积木一样组合流程。
- 每个管理节点都可在节点侧边栏中编辑结构化数据，并能关联资源（图像、音频、文本段落、列表/表格）。
- 保存工作流后，再次打开：节点数据与资源引用完整恢复；下游节点可直接引用上游节点输出的数据。
- 生成类能力（LLM、图片生成、HTTP/插件等）尽量复用现有节点；“创作管理节点”负责把数据整理成可被生成节点消费的结构。

### 1.2 典型流程（示例）

- 小说/剧本：创作目标 → 世界观 → 角色库 → 大纲/章节 → 场景拆解 → 剧本文本 → 修订/定稿 → 输出/发布
- 分镜视频（含音频）：创作目标 → 世界观 → 角色库 → 分镜/镜头表 → 镜头视觉参考 → 旁白/对白 → 配音方案 → BGM/音效方案 → 合成/输出

---

## 二、先读现有代码：当前节点体系如何工作（必须对齐现状）

> 这一节用于确保“新增节点”走的是 **既有的节点库/注册/保存/执行** 链路，而不是文档里凭空发明新系统。

### 2.1 后端（Go）关键结构与约束

- 节点类型与模板来源：`backend/domain/workflow/entity/node_meta.go`
  - `NodeTypeMeta` 以 **数值 ID**（例如 1/2/3/30…）作为节点类型的对外标识（NodeTemplateList 返回 `node_id`）。
  - `NodeType` 的内部 Key 是字符串（例如 `InputReceiver`、`TextProcessor`），通过 `IDStrToNodeType` 做映射。
- 工作流画布数据结构：`backend/domain/workflow/entity/vo/canvas.go`
  - `Node.Type` 是字符串，但语义上是 NodeTypeMeta.ID 的字符串形式（例如 `"30"`）。
  - `Node.Data.Inputs` 结构是“通用字段 + 各节点专用结构体指针”的组合，保存/校验/执行会依赖该结构。
- 执行与适配：`backend/domain/workflow/internal/canvas/adaptor/to_schema.go`
  - 每种节点类型必须有对应 `NodeAdaptor`，否则会在编译 Schema 时直接报 `unsupported block type`。
  - 适配器注册在 `to_schema.go` 内的 `RegisterNodeAdaptor(...)` 代码段。
- 节点模板列表接口：`backend/application/workflow/workflow.go` 的 `GetNodeTemplateList(...)`
  - 前端节点面板展示的分类与模板来自后端 `ListNodeMeta(...)`。

### 2.2 前端（TS/React）关键结构与约束

- 节点注册位置：`frontend/packages/workflow/playground/src/node-registries/*`
  - 每个节点提供一个 `WorkflowNodeRegistry`，包含 `type`、`meta`、`formMeta` 等。
  - 入口汇总在：`frontend/packages/workflow/playground/src/node-registries/index.ts`。
- 保存/提交时的类型与数据格式化：`frontend/packages/workflow/nodes/src/workflow-json-format.ts`
  - 会在 `formatNodeOnSubmit` 中把变量/输入值转换为 DTO，并提交到后端。
  - 这意味着“新增节点”必须兼容现有的 DTO/Inputs 表达方式（尤其是 `InputParameters`、ValueExpression、文件类型）。
- 节点模板信息获取：`frontend/packages/workflow/playground/src/workflow-playground-context.ts`
  - `getNodeTemplateInfoByType(type)` 会通过后端模板列表映射节点标题/图标/描述。
  - 所以想让新增节点在节点面板里“像原生节点一样出现”，需要后端提供对应 NodeTypeMeta（而不仅是前端自定义渲染）。
- 文件/资源输入能力：`frontend/packages/workflow/playground/src/node-registries/common/fields/value-expression-input.tsx`
  - `ValueExpressionInput` 支持 `availableFileTypes`，可用于选择/引用 image/audio/video/doc 等资源。
  - 注意在实现前端代码时，要严格按照eslint的约束，避免在我git commit提交时提交失败。

---

## 三、要新增哪些“创作管理节点”（全部同等能力）

> 这里列出的节点不是把世界观当特殊；相反，世界观只是其中一个。所有节点遵守同一套通用能力，并各有自己的特有字段/资源展示。

### 3.1 通用能力（所有创作管理节点必须具备）

1. **结构化数据编辑**：在节点侧边栏提供“表单 + 富文本 + 列表编辑”能力；数据以工作流画布 JSON 持久化（即节点 inputs/outputs 的一部分）。
2. **资源绑定（图片/文本/列表/音频/视频）**：节点可维护一个“资源引用列表”，用于展示、选择、替换、复用。
   - 资源引用在节点数据里以“文件类型 ValueExpression”或“资源 ID/URL”方式持久化。
3. **可被下游节点消费**：节点输出提供稳定的结构（JSON object / array），下游 LLM/图片生成/插件/HTTP 节点可直接引用。
4. **版本与修订（轻量）**：至少提供：
   - 节点级“最近修改时间/人”元信息（可用工作流保存历史间接满足）
   - 数据 diff/回滚（可选，后续增强）
5. **一致的权限/只读策略**：可在“发布/只读/模板”状态下禁用编辑（复用工作流已有 readonly 机制）。

### 3.2 节点清单（建议首批实现）

| 节点 | 主要输出（供下游消费） | 特有字段/行为 | 典型资源类型 |
| --- | --- | --- | --- |
| 创作目标（Project Brief） | `brief`（object） | 题材、受众、风格、约束、参考作品 | 文本、参考图/链接 |
| 世界观/设定（Setting） | `setting`（object） | 世界规则、势力/地理/时间线；条目化管理 | 文本、参考图、文档 |
| 角色库（Characters） | `characters`（array<object>） | 多角色列表、标签、关系图；角色头像/立绘 | 图片、文本、列表 |
| 剧情大纲（Outline） | `outline`（object/array） | 章节/幕结构、关键事件、节奏点 | 文本、列表 |
| 场景/章节拆解（Scenes） | `scenes`（array<object>） | 场景卡片：时间地点人物冲突目标 | 文本、列表、参考图 |
| 分镜/镜头表（Storyboard） | `shots`（array<object>） | 每镜头：景别、运动、构图、情绪、时长、台词关联 | 图片（参考帧）、列表、文本 |
| 台词/旁白（Dialogue） | `dialogue`（array/object） | 角色台词、旁白、情绪标注；与镜头/场景关联 | 文本、列表 |
| 音频方案（Audio Plan） | `audio_plan`（object） | 配音角色、音色、BGM/SFX 清单、时轴标注 | 音频、列表、文本 |
| 资产面板（Assets Panel） | `assets`（array<object>） | 按类型/标签聚合展示整个工作流相关资源；支持回链到来源节点 | 图片/音频/视频/文本 |

说明：上述每个节点都具备“结构化数据 + 资源引用 + 输出给下游”的共性；差异主要在字段结构、资源类型、以及与其他节点的关联粒度（角色与分镜会更强调图片列表与引用）。

---

## 四、实现方案（在现有节点体系上“增量”实现）

### 4.1 总体策略：管理节点负责“数据与资产”，生成节点复用现有能力

- 生成类节点：优先复用现有 `LLM`、`ImageGenerate`、`Http`、`Plugin`、`TextProcessor` 等节点。
- 创作管理节点：新增一组“只做数据整理/输出”的节点类型，使它们在运行时可以把节点内保存的数据输出给下游节点。

### 4.2 后端需要做什么（必须做到，否则节点无法被执行/保存）

1. **新增 NodeTypeMeta（让节点出现在面板中）**
   - 修改 `backend/domain/workflow/entity/node_meta.go`：
     - 增加一个新分类（例如 `creative`）。
     - 为每个创作管理节点分配稳定的数值 ID、名称、描述、颜色、图标 URI。
2. **新增节点执行适配器（让节点类型可被编译/执行）**
   - 在 `backend/domain/workflow/internal/canvas/adaptor/to_schema.go` 注册新节点的 `NodeAdaptor`。
   - 新 adaptor 的核心目标：
     - 从画布 node 的 `Data.Inputs.InputParameters` 读取字段值（常量/引用/文件）。
     - 通过 `convert.SetInputsForNodeSchema(...)`、`convert.SetOutputTypesForNodeSchema(...)` 填充 schema。
3. **新增一个通用执行实现（Identity/Pass-through）**
   - 新增一个通用 node（例如 `backend/domain/workflow/internal/nodes/creativedoc`）：
     - `Invoke(ctx, input)` 直接把输入字段透传为输出字段（或组装成一个 object 输出）。
     - 这样管理节点就能在执行时把编辑好的结构化数据“作为变量”提供给下游节点。
4. **校验策略**
   - 复用现有 CanvasValidator：对管理节点只校验“必须字段存在/类型正确”，避免过度约束。

### 4.3 前端需要做什么（节点库内新增，兼容现有 workflow 逻辑）

1. **新增节点 registry**
   - 在 `frontend/packages/workflow/playground/src/node-registries/` 为每个新节点创建目录与 `node-registry.ts`。
   - 在 `frontend/packages/workflow/playground/src/node-registries/index.ts` 导出这些 registries。
2. **复用现有表单与 ValueExpression 能力**
   - 文本/对象/数组字段：使用现有表单系统（`formMeta` + setters）。
   - 文件字段：使用 `ValueExpressionInput` 并传入 `availableFileTypes`，支持 image/audio/video/doc。
   - 列表字段：提供“列表编辑器”（数组增删改）+ 每项支持文本与文件引用。
3. **资产展示（节点内 + 全局资产面板）**
   - 节点内：展示“该节点引用的资源列表”（缩略图/播放控件/文本预览）。
   - 全局：提供一个“资产面板节点/侧边栏”，聚合整个工作流画布中出现的资源引用，并支持按来源节点回链。

### 4.4 数据保存与再次打开（绑定工作流即可）

- 节点的结构化数据与资源引用统一保存在工作流画布 JSON（node `data.inputs` / `data.outputs`）中，由现有 SaveWorkflow 流程持久化。
- 再次打开工作流：前端从画布 JSON 恢复表单值与资源引用，节点即可展示上次编辑内容与相关资产。

---

## 五、落地顺序（推荐分三步，避免一次性大改）

1. **Step 1：先打通一个最小闭环**
   - 后端新增 1 个“创作管理节点”类型（例如 `Project Brief`），实现 adaptor + pass-through 执行。
   - 前端新增对应 registry，支持文本 + 列表 + 资源引用（至少图片）。
2. **Step 2：批量补齐节点清单**
   - 逐步加入 Setting / Characters / Storyboard / Audio Plan。
   - 强化列表编辑、资源引用、节点间引用体验。
3. **Step 3：资产面板与复用**
   - 资产面板节点：聚合全工作流资源，支持搜索/标签/回链；支持把资产引用插入任意节点字段。

---

## 六、验收标准（以“可复访 + 可执行 + 可复用”为主）

1. 任意新增的创作管理节点（不区分世界观/角色/分镜）都能：编辑 → 保存 → 再次打开恢复。
2. 角色/分镜/音频等节点可以展示并持久化图片/音频/列表/文本资源引用。
3. 管理节点执行后能把结构化输出提供给下游节点（LLM/图片生成/插件等）作为输入变量引用。
4. 新节点出现在节点面板中，符合现有分类/模板机制（NodeTemplateList + node-registries）。

---

## 近期更新（Step1 落地记录）

- 新增“创作管理”节点：实现 `CreativeBrief` / Project Brief 节点（NodeType ID 1001，新建 `creative` 分类），通过 adaptor 注册和 pass-through 执行，输出 `brief` 对象（title/audience/style/constraints/references）。
- 前端注册与面板：`frontend/packages/workflow/playground/src/node-registries/creative-brief/*` 实现节点 registry/表单/内容，支持文本、数组与资源引用；已加入 V2 节点常量、可用节点列表和 node-registries 集成，保证节点可见、可编辑。
