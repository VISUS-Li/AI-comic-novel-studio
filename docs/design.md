下面这个文档会按「专业设计文档」+「清晰开发步骤」来写，你可以**直接复制整篇**当成后续 AI 生成代码的提示词使用。
我会分成几大块：

1. 背景与目标
2. 业务需求梳理
3. 总体架构设计（结合 coze-studio 原有结构）
4. 领域模型 & 数据结构设计
5. 工作流与节点设计
6. 前后端改造点梳理
7. 分阶段开发计划（Step 1 / Step 2 / …，你可以逐步执行并标记完成）

---

# 一、背景与目标

## 1.1 项目背景

在现有开源项目 **coze-studio** 的基础上进行二次开发，新增一个面向 **小说 / 剧本创作** 的能力，包括：

* 世界观、大纲、人物设定的生成与编辑
* 分场景 / 分镜设计
* 分镜图片、视频、音频的生成与管理
* 支持多种创作模式（文字优先、镜头优先、角色优先等）
* 使用 **coze-studio 现有的工作流（Workflow）和节点（Node）系统** 来搭建灵活的创作流程

目标是将 coze-studio 打造成一个 **“有状态的剧本创作工作台”**，每个项目及其工作流可以长期迭代、持续创作，而不是一次性的流水线。

## 1.2 总体目标

1. **新增“剧本项目（Script Project）”概念**：
   每个剧本项目包含世界观、人物、大纲、剧本文本，以及相关的图片 / 视频 / 音频资产。

2. **工作流有状态**：
   工作流不仅是结构定义（Definition），还具备**会话状态**，可以保存/恢复上一次运行的上下文和生成内容。

3. **资产管理体系**：

    * 所有资源（文本片段、图片、视频、音频）与：

        * 剧本项目
        * 工作流实例
        * 甚至具体场景/镜头节点
          强关联。
    * 在打开某个项目和其工作流时，可以看到、管理、复用历史资产。

4. **创作工作流可灵活配置**：
   不同创作模式对应不同的工作流模板，支持自由拖拽节点调整流程。

5. **单节点 + 全流程运行**：
   每个节点支持单独运行以调试和微调；也支持一键运行整个工作流。

---

# 二、业务需求梳理

## 2.1 核心业务对象

* **剧本项目（Script Project）**

    * 一部小说/剧本/剧集工程。
    * 关联：世界观设定、人物卡、章节大纲、剧本文本、资产、工作流实例等。

* **工作流定义（Workflow Definition）**

    * coze-studio 原有概念：节点图。
    * 用于描述某种创作模式下的节点流程（例如“文字优先工作流”、“分镜优先工作流”）。

* **工作流实例（Workflow Instance / Session）**

    * 绑定到某个剧本项目 + 某个 workflow 定义。
    * 有自己的 `last_context`（上次运行的上下文）、运行历史、已产生的资产。
    * 用户下次进入时可以继续创作。

* **资产（Asset）**

    * 文本（可选）、图片、视频、音频等。
    * 绑定：项目、工作流实例、节点、业务标识（场景/镜头 ID）。

* **工作流节点（Node）**

    * 代表单个创作步骤，如：

        * 世界观生成
        * 人物卡生成
        * 大纲生成
        * 场景细化
        * 分镜生成
        * 图片生成
        * 视频生成
        * 音频生成
        * 文本局部编辑（Text Edit）

## 2.2 用户场景

1. 创建剧本项目 → 选择一种创作模式 → 自动生成对应的工作流实例。
2. 在项目中：

    * 编辑世界观、人物、大纲、剧本。
    * 通过工作流节点生成分镜、图片、视频、音频。
3. 再次打开项目 & 工作流：

    * 能看到之前生成的世界观、大纲、分镜图、视频等。
    * 继续在这些基础上微调、增删、再生成。
4. 在文本编辑过程中：

    * 选中任意文本块（例如某段大纲），调用 AI 进行局部修改（润色/改风格/转成剧本格式等）。
    * 可视化 diff 或直接替换。

---

# 三、总体架构设计（结合 coze-studio 结构）

## 3.1 coze-studio 原有结构（简要）

* 后端：`backend/`

    * `api/`：Hertz HTTP API 层（handler、model、router）
    * `application/`：应用服务（编排 domain + infra）
    * `domain/`：领域模型和业务逻辑
    * `infra/`：基础设施实现（DB、缓存、外部服务）
    * `types/`：常量、错误码、DDL 等

* 前端：`frontend/`

    * `apps/coze-studio/`：主应用
    * `packages/`：多个业务/基础包（workflow、studio 等）

## 3.2 新增/扩展模块概览

从 DDD 角度设计：

### 后端新增领域模块

* `backend/domain/script/`

    * `script_project.go`：剧本项目聚合根
    * `workflow_instance.go`：工作流实例模型
    * `asset.go`：资产模型

* `backend/application/script/`

    * `ScriptProjectService`
    * `WorkflowInstanceService`
    * `AssetService`

同时扩展现有 workflow 执行逻辑，支持 “实例化工作流” 的操作：

* 在 `backend/application/workflow/` 中增加：

    * `RunWorkflowInstance(instanceID, options)`
    * `RunNodeInInstance(instanceID, nodeID, inputOverrides)`

### 前端新增模块

* 新增一个剧本工作台包，例如：

    * `frontend/packages/script-studio/`（命名可调整）
* 在 `frontend/apps/coze-studio/` 中新增页面/路由：

    * `/script-projects`
    * `/script-projects/:projectId/workflows/:instanceId`

该页面集成：

* 文本编辑器（世界观/人物/大纲/剧本）
* 工作流画布（复用现有 workflow 包）
* 资产管理侧栏（展示图片/视频/音频/文本资产）

---

# 四、领域模型 & 数据结构设计

以下设计可以映射为 DB 表或 Thrift/IDL 结构，后续可由 AI 生成具体代码。

## 4.1 剧本项目（ScriptProject）

**领域意图**：代表一个独立的小说/剧本工程。

```text
ScriptProject
- id: string
- name: string
- type: enum("novel", "script", "comic", ...)
- owner_id: string
- description: string
- worldbuilding: Text / JSON  （世界观设定）
- created_at: datetime
- updated_at: datetime
```

可选扩展：将 worldbuilding / characters / outline 拆成子表，这里先简化。

## 4.2 工作流实例（WorkflowInstance）

**领域意图**：将抽象的 workflow 定义实例化为项目中的一个有状态“创作会话”。

```text
WorkflowInstance
- id: string
- project_id: string         # ScriptProject.id
- workflow_id: string        # 原 coze-studio Workflow Definition ID
- name: string               # 如「S01E01 文字优先工作流」
- mode: enum("outline_first", "shot_first", "character_first", ...)
- last_context: JSON         # 最近一次运行后的全局上下文（worldbuilding / outline / scenes）
- status: enum("active", "archived")
- created_at: datetime
- updated_at: datetime
```

## 4.3 工作流运行记录（WorkflowRun）

**领域意图**：记录每次运行实例时的执行过程和结果，便于回溯/对比。

```text
WorkflowRun
- id: string
- instance_id: string            # WorkflowInstance.id
- started_at: datetime
- finished_at: datetime
- status: enum("running", "success", "failed")
- input_context: JSON            # 运行前的上下文快照
- output_context: JSON           # 运行后的上下文快照
- error_message: string (optional)
```

## 4.4 节点输出（NodeOutput）

**领域意图**：记录某次运行中，单个节点的输出及其关联的资产。

```text
NodeOutput
- id: string
- run_id: string                 # WorkflowRun.id
- node_id: string                # Workflow Definition 中的节点 ID
- output_data: JSON              # 节点的业务输出（文本结构、大纲、分镜信息等）
- asset_ids: string[]            # 关联到 Asset.id 的列表
- created_at: datetime
```

## 4.5 资产（Asset）

**领域意图**：统一管理所有生成的资源（文本/图片/视频/音频）。

```text
Asset
- id: string
- project_id: string
- workflow_instance_id: string (nullable)
- node_id: string (nullable)
- type: enum("text", "image", "video", "audio")
- url: string                    # image/video/audio 对应资源地址；文本可以为空
- text_content: text (optional)  # 文本类资产可存这里
- meta: JSON                     # 附加信息，例如：
                                 # { "scene_id": "...", "shot_id": "...", "prompt": "...", "model": "...", "duration": 3.5 }
- created_at: datetime
- updated_at: datetime
```

---

# 五、工作流与节点设计

## 5.1 节点家族概览

各节点类型的职责与输入输出，可以通过 IDL + TS 类型来定义，这里只写概念。

### 1. 结构类节点（文本结构）

* `WorldbuildingNode`

    * 输入：题材、风格、用户已有设定片段
    * 输出：世界观设定文本/JSON（可存为 text Asset 或直接写入 last_context）

* `CharacterSheetNode`

    * 输入：世界观设定、人物需求
    * 输出：角色卡数组（人物名、性格、背景、动机、人物弧光等）

* `OutlineNode`

    * 输入：世界观、人物卡、篇幅/集数要求
    * 输出：分幕/分章节大纲（JSON + 文本描述）

* `SceneDetailNode`

    * 输入：某一条大纲（Chapter/Scene）
    * 输出：场景详细描述（地点/时间/冲突/节奏）

* `ScriptFormatNode`

    * 输入：剧情文本
    * 输出：剧本格式文本（场标/角色台词/动作提示等）

### 2. 视觉与分镜节点

* `ShotBreakdownNode`

    * 输入：Scene 描述或剧本片段
    * 输出：镜头列表（镜头号、景别、运动、画面重点、时长等）

* `StoryboardTextNode`

    * 输入：单个镜头信息
    * 输出：面向图像模型的 prompt（构图、光线、风格）

* `ImageGenerationNode`

    * 输入：图像 prompt + 模型配置
    * 输出：图片 URL（Asset type=image）

* `KeyframeImageNode`

    * 输入：首尾帧描述
    * 输出：首/尾帧图片 URL

### 3. 视频节点

* `VideoFromImagesNode`

    * 输入：有序图片列表（Asset IDs）+ 时长信息
    * 输出：视频 URL（Asset type=video）

* `TextToVideoNode`

    * 输入：场景/镜头文字描述
    * 输出：视频 URL

* `VideoAssemblyNode`

    * 输入：多个视频 Asset + 转场参数
    * 输出：合成视频 URL

### 4. 音频节点

* `DialogTTsNode`

    * 输入：角色台词结构 + 声音设定
    * 输出：语音文件 URL（按角色分轨）

* `NarrationTtsNode`

    * 输入：旁白文本
    * 输出：旁白音频 URL

* `BGMMusicNode`

    * 输入：场景氛围标签
    * 输出：BGM 音频 URL

* `AudioMixNode`

    * 输入：对白、旁白、BGM 音轨 + 混音策略
    * 输出：混音音频 URL

### 5. 文本编辑节点

* `TextEditNode`

    * 输入：

        * 目标文本片段
        * 编辑意图（润色/缩写/扩写/改风格/改格式等）
        * 上下文（世界观/人物/当前场景）
    * 输出：

        * 修改后的文本
        * 可选 diff 信息

前端文本编辑器可以调用 TextEditNode 对局部文本进行一次“单节点运行”，不必跑整条 workflow。

---

# 六、前后端改造点梳理

## 6.1 后端改造点

1. **新增 script 领域**

    * 目录：`backend/domain/script`
    * 内容：

        * `ScriptProject` 聚合根
        * `WorkflowInstance` 模型
        * `Asset` 模型
    * 定义相应的 Repository 接口（`ScriptProjectRepo`、`WorkflowInstanceRepo`、`AssetRepo`）。

2. **新增 script 应用服务**

    * 目录：`backend/application/script`
    * 服务：

        * `ScriptProjectService`：CRUD + 列表 + 绑定工作流实例
        * `WorkflowInstanceService`：

            * 创建新实例（项目 + workflow 定义）
            * 获取实例详情（含 last_context、最近运行态）
        * `AssetService`：

            * 创建资产
            * 按项目/工作流实例查询资产

3. **扩展 workflow 执行逻辑**

    * 在 `backend/application/workflow` 中增加：

        * `RunWorkflowInstance(instanceID, options)`：

            * 创建 WorkflowRun
            * 从 WorkflowInstance.last_context 初始化上下文
            * 调用现有 workflow runtime 串行/并行执行节点
            * 每个节点完成后：

                * 通过 `AssetService` 创建资产（如果有图片/视频/音频/大文本资源）
                * 记录 NodeOutput
            * 运行结束后更新 WorkflowInstance.last_context
        * `RunNodeInInstance(instanceID, nodeID, inputOverrides)`：

            * 支持单节点运行，用于调试和局部生成。

4. **文本编辑 API**

    * 在 `backend/api/` 新增一个 controller，例如 `script_text_edit_handler.go`：

        * 接收：project_id、instance_id、原始文本、编辑意图、上下文信息
        * 内部调用 `TextEditNode` 对应的 workflow 节点执行逻辑。

5. **资产存储**

    * 使用现有的 MinIO / 对象存储方案：

        * `Asset.url` 存储对象地址
    * `meta` 字段用于存储场景/镜头等业务信息。

## 6.2 前端改造点

1. **新增 Script Studio 包**

    * 目录建议：`frontend/packages/script-studio/`
    * 功能：

        * 管理 ScriptProject 列表 & 详情
        * 管理 WorkflowInstance 的打开/关闭
        * 提供资产管理组件（AssetPanel）
        * 提供文本编辑器组件（带 AI 编辑入口）

2. **整合到主应用**

    * 在 `frontend/apps/coze-studio/src/routes` 中新增路由：

        * `/script-projects`
        * `/script-projects/:projectId/workflows/:instanceId`
    * 在对应页面中：

        * 左：文本编辑器（多 Tab：世界观/人物/大纲/剧本）
        * 中：workflow 画布（复用 `packages/workflow`）
        * 右：资产面板（按类型分类：文本/图片/视频/音频）

3. **状态管理**

    * 使用 Zustand（与项目现有状态管理一致）：

        * 全局 store 包含：

            * `currentProject`
            * `currentWorkflowInstance`
            * `lastContext`
            * `nodeOutputs`
            * `assets`
    * 打开工作流时：

        * 通过 API 拉取 WorkflowInstance + last_context + 资产列表
        * 填充到 store，并渲染 UI。

4. **文本编辑器中的 AI 集成**

    * 在文本选区右键菜单或工具栏中增加“AI 编辑”按钮：

        * 选择编辑意图（比如“润色”、“改剧本格式”）
        * 调用后端 TextEdit API
        * 返回结果后：

            * 可选择直接替换
            * 或显示 diff 供用户确认

5. **单节点运行的前端交互**

    * 在 workflow 节点 UI 上：

        * 增加右键菜单项「仅运行此节点」
    * 前端调用 `RunNodeInInstance` API：

        * 将结果写入 store，更新对应 NodeOutput 和相关资产。
    * 在节点详情面板中展示本节点最近输出（文本/图片/视频/音频）。

---

# 七、分阶段开发计划（可逐步执行的 Step 计划）

下面是为“AI 辅助编码”设计的分步骤计划。
你可以按顺序执行：**完成 Step N → 标记完成 → 用这份文档 + 当前进度让 AI 继续生成 Step N+1 的代码。**

---

## Step 1：基础领域 & 数据建模（后端）

**目标**：让后端具备 ScriptProject / WorkflowInstance / Asset 的基本数据结构与持久化能力。

**主要任务：**

1. 在 `backend/types/ddl/` 中：

    * 增加剧本相关表结构的 DDL 定义：

        * `script_projects`
        * `workflow_instances`
        * `workflow_runs`
        * `node_outputs`
        * `assets`

2. 在 `backend/domain/script/` 中：

    * 定义：

        * `ScriptProject` 领域模型
        * `WorkflowInstance`
        * `WorkflowRun`
        * `NodeOutput`
        * `Asset`
    * 定义对应的 repository 接口。

3. 在 `backend/infra/impl` 中：

    * 实现 `ScriptProjectRepo`、`WorkflowInstanceRepo`、`WorkflowRunRepo`、`NodeOutputRepo`、`AssetRepo` 的 MySQL 持久化。

4. 在 `backend/application/script/` 中：

    * 实现：

        * `ScriptProjectService.Create/List/Get`
        * `WorkflowInstanceService.Create/Get`
---
        * `AssetService.Create/ListByProject/ListByInstance`

**Current progress**

* Added the SQL DDL for `script_projects`, `workflow_instances`, `workflow_runs`, `node_outputs`, and `assets` in `backend/types/ddl/script_tables.sql`.
* Implemented the domain models + gorm mappings under `backend/domain/script/internal/dal/model` plus DAO logic and repository interfaces in `backend/domain/script/repository`.
* Exposed ScriptProject / WorkflowInstance / Asset services in `backend/application/script` to satisfy the Step 1 API contract (`Create`, `List`, `Get` entrypoints) and to act as the basis for future API handlers.

---

## Step 2：Script 相关 API（后端）

**目标**：暴露基础 HTTP API，供前端管理剧本项目和工作流实例。

**主要任务：**

1. 在 `backend/api/model/script` 下定义请求/响应模型：

    * CreateScriptProjectRequest/Response
    * ListScriptProjectsResponse
    * GetScriptProjectResponse
    * CreateWorkflowInstanceRequest/Response
    * GetWorkflowInstanceDetailResponse（包含 last_context 和最近 run 摘要）
    * ListAssetsResponse

2. 在 `backend/api/handler/script` 中实现 handler：

    * `CreateScriptProject`
    * `ListScriptProjects`
    * `GetScriptProject`
    * `CreateWorkflowInstance`
    * `GetWorkflowInstanceDetail`
    * `ListAssetsByProject`
    * `ListAssetsByWorkflowInstance`

3. 在 `backend/api/router` 中注册对应路由：

    * `/api/script/projects`
    * `/api/script/projects/:projectId`
    * `/api/script/projects/:projectId/workflow-instances`
    * `/api/script/workflow-instances/:instanceId`
    * `/api/script/assets`（带 query 参数：project_id / instance_id）

---

## Step 3：实例化工作流执行（RunWorkflowInstance）

**目标**：把 coze-studio 的工作流运行能力包装成「带状态的实例运行」。

**主要任务：**

1. 在 `backend/application/workflow` 中：

    * 新增 `RunWorkflowInstance(instanceID, options)`：

        * 根据 instanceID 读取 WorkflowInstance + 对应 Workflow Definition。
        * 从 `WorkflowInstance.last_context` 初始化上下文。
        * 调用现有 workflow runtime 顺序/依赖执行各节点。
        * 在节点执行完成后:

            * 使用 `AssetService` 创建资产（如有 image/video/audio/text）。
            * 使用 `NodeOutputRepo` 记录节点输出。
        * 执行完成后更新：

            * `WorkflowRun`（状态、output_context）
            * `WorkflowInstance.last_context = 最新上下文`

2. 在 `backend/api/handler/script` 中新增：

    * `RunWorkflowInstanceHandler`：

        * 路由示例：`POST /api/script/workflow-instances/:instanceId/run`
        * 字段：是否全量运行 / 指定起止节点等 options。

---

## Step 4：单节点运行（RunNodeInInstance）

**目标**：支持在某个工作流实例中，仅执行指定节点（用于调试和局部生成）。

**主要任务：**

1. 在 `application/workflow` 中：

    * 新增 `RunNodeInInstance(instanceID, nodeID, inputOverrides)`：

        * 从实例 `last_context` 或最近一次 run 的输出取输入。
        * 仅执行指定 node 的逻辑。
        * 记录 NodeOutput 与 Asset。
        * 视情况更新 `last_context`（可由参数控制）。

2. 在 `api/handler/script` 中新增：

    * `RunNodeInInstanceHandler`：

        * 路由示例：`POST /api/script/workflow-instances/:instanceId/nodes/:nodeId/run`

---

## Step 5：Script Studio 前端基础页面

**目标**：在前端展示剧本项目列表 & 工作流实例，并能打开一个实例页面。

**主要任务：**

1. 新增前端包 `frontend/packages/script-studio/`：

    * 定义 ScriptProject 类型 / WorkflowInstance 类型 / Asset 类型。
    * 封装调用 script 相关 API 的 client。

2. 在 `apps/coze-studio` 中新增路由：

    * `/script-projects`
    * `/script-projects/:projectId/workflows/:instanceId`

3. 页面内容：

    * 列表页：展示 ScriptProject 列表，可新建项目。
    * 详情页：

        * 左侧暂时简单放几个文本框或占位组件（后续换成专业编辑器）。
        * 中间嵌入现有 workflow 画布组件，加载对应 workflow 定义。
        * 右侧用简单列表展示当前实例的资产（从 `ListAssetsByWorkflowInstance` 获取）。

---

## Step 6：资产管理 UI（资产面板）

**目标**：在前端工作流页面中，提供一个可视化资产面板，按类型查看资产。

**主要任务：**

1. 在 `script-studio` 包中：

    * 实现一个 `AssetPanel` 组件：

        * Tabs：文本 / 图片 / 视频 / 音频。
        * 每个 Tab 中展示对应类型的 Asset 卡片（缩略图、标题、时间、场景/镜头标签等）。

2. 在 `/script-projects/:projectId/workflows/:instanceId` 页面中：

    * 引入 AssetPanel，并在初始化时：

        * 调用 `ListAssetsByWorkflowInstance`。
        * 使用 Zustand 或其他 store 存储 assets。
    * 提供刷新机制（当工作流运行完毕或者单节点运行完成后刷新资产列表）。

---

## Step 7：文本编辑器 + TextEditNode 集成

**目标**：在剧本文本编辑页面中，支持选中文本片段 → 使用 AI 节点进行局部编辑。

**主要任务：**

1. 后端：

    * 为 `TextEditNode` 增加应用层执行逻辑（可以是 workflow 中的一个节点，也可以是独立方法）。
    * 提供 API：

        * `POST /api/script/text-edit`
        * 请求参数：`project_id`, `instance_id`, `original_text`, `edit_intent`, `context`（世界观/人物/当前场景标识）。
        * 返回：`edited_text`, `diff`（可选）。

2. 前端：

    * 在 Script Studio 页面中集成一个富文本编辑器或 Markdown 编辑器（可用现有组件或第三方）。
    * 在文本选区右键菜单/工具栏中增加“AI 编辑”入口：

        * 弹出对话框，让用户选择编辑意图。
        * 调用 `text-edit` API。
        * 将返回的文本插入或替换选中的内容。

---

## Step 8：剧本专用节点定义与注册（前后端）

**目标**：正式把 Worldbuilding/Outline/Shot/Image/Video/Audio 等节点类型注册到 workflow 系统中。

**主要任务：**

1. 后端：

    * 在 workflow 域中增加节点类型定义：

        * `WorldbuildingNode`, `CharacterSheetNode`, `OutlineNode`, `SceneDetailNode`, `ShotBreakdownNode`, `StoryboardTextNode`, `ImageGenerationNode`, `VideoFromImagesNode`, `DialogTTsNode`, `BGMMusicNode`, `AudioMixNode` 等（按优先级逐步实现）。
    * 为每个节点提供：

        * 输入/输出 schema（IDL / Go struct）
        * 执行业务逻辑（调用 LLM、图片/视频/音频生成服务）

2. 前端：

    * 在 workflow 节点注册中心中，为这些节点配置：

        * 节点类型标识字符串
        * 配置表单 UI（参数设置）
        * 节点展示样式（图标、颜色等）

3. 预置几个工作流模板：

    * 文本优先工作流
    * 分镜优先工作流
    * 人物优先工作流

---

## Step 9（可选）：进一步拆分文本结构为独立子域

当主线功能跑通后，你可以考虑把世界观/人物/大纲/场景抽成更细的领域与表结构，以获得更高的可视化能力与复用度。不过这属于下一阶段演进。

---

