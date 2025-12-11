/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useParams } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';

import {
  AssetGrid,
  ContextSnapshot,
  ErrorBanner,
  HistoryRuns,
  InstanceInfo,
  ProjectInfo,
  SectionCard,
  TextEditSection,
} from './script-workflow-sections';

export interface ScriptProject {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  type: string;
  created_at: string;
}

export interface WorkflowInstanceDetail {
  instance: {
    id: string;
    project_id: string;
    workflow_id: string;
    name: string;
    mode: string;
    last_context: string;
    status: string;
    created_at: string;
    updated_at: string;
  };
  recent_runs: {
    id: string;
    started_at: string;
    finished_at?: string;
    status: string;
  }[];
}

export interface Asset {
  id: string;
  type: string;
  url: string;
  node_id: string;
  text_content?: string;
  meta: string;
  created_at: string;
}

export interface TextEditResult {
  edited_text: string;
  diff: string;
}

const layoutStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
  gap: '16px',
  padding: '24px',
};

const PREVIEW_INDENT = 2;

const usePreviewContext = (rawContext?: string) =>
  useMemo(() => {
    if (!rawContext) {
      return '未生成上下文';
    }
    try {
      const parsedContext = JSON.parse(rawContext);
      return JSON.stringify(parsedContext, null, PREVIEW_INDENT);
    } catch (caughtError) {
      const safeError =
        caughtError instanceof Error
          ? caughtError
          : new Error(String(caughtError ?? '未知错误'));
      console.error('解析上下文失败', safeError);
      return rawContext;
    }
  }, [rawContext]);

const safeFetchData = async <T,>(
  url: string,
  failureMessage: string,
): Promise<T | null> => {
  const response = await fetch(url, { credentials: 'include' });
  if (!response.ok) {
    throw new Error(`${failureMessage}（${response.status}）`);
  }
  const payload = (await response.json()) as { data?: T };
  return (payload?.data ?? null) as T | null;
};

const useScriptWorkflowData = (projectId?: string, instanceId?: string) => {
  const [project, setProject] = useState<ScriptProject | null>(null);
  const [detail, setDetail] = useState<WorkflowInstanceDetail | null>(null);
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    if (!projectId || !instanceId) {
      setLoading(false);
      setError('缺失 projectId 或 workflow 实例 id');
      return () => {
        cancelled = true;
      };
    }

    const fetchAll = async () => {
      setLoading(true);
      setError(null);
      try {
        const [projectData, instanceData, assetList] = await Promise.all([
          safeFetchData<ScriptProject>(
            `/api/script/projects/${projectId}`,
            '获取项目失败',
          ),
          safeFetchData<WorkflowInstanceDetail>(
            `/api/script/workflow-instances/${instanceId}`,
            '获取实例信息失败',
          ),
          safeFetchData<Asset[]>(
            `/api/script/assets?instance_id=${instanceId}`,
            '获取资产失败',
          ),
        ]);
        if (!cancelled) {
          setProject(projectData);
          setDetail(instanceData);
          setAssets(assetList ?? []);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : '网络错误');
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    fetchAll();

    return () => {
      cancelled = true;
    };
  }, [projectId, instanceId]);

  return { project, detail, assets, loading, error };
};

const useTextEditState = (
  projectId?: string,
  instanceId?: string,
  lastContext?: string,
) => {
  const [originalText, setOriginalText] = useState('');
  const [editIntent, setEditIntent] = useState('');
  const [textEditResult, setTextEditResult] = useState<TextEditResult | null>(
    null,
  );
  const [textEditLoading, setTextEditLoading] = useState(false);
  const [textEditError, setTextEditError] = useState<string | null>(null);

  useEffect(() => {
    if (lastContext) {
      setOriginalText(lastContext);
    }
  }, [lastContext]);

  const handleTextEdit = async () => {
    if (!projectId || !instanceId) {
      setTextEditError('请选择项目与实例');
      return;
    }
    setTextEditLoading(true);
    setTextEditError(null);
    setTextEditResult(null);
    try {
      const response = await fetch('/api/script/text-edit', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          project_id: projectId,
          instance_id: instanceId,
          original_text: originalText,
          edit_intent: editIntent,
          context: lastContext ?? '',
        }),
      });
      if (!response.ok) {
        throw new Error(`AI 编辑失败（${response.status}）`);
      }
      const payload = (await response.json()) as TextEditResult;
      setTextEditResult(payload);
    } catch (err) {
      setTextEditError(err instanceof Error ? err.message : 'AI 编辑失败');
    } finally {
      setTextEditLoading(false);
    }
  };

  return {
    originalText,
    setOriginalText,
    editIntent,
    setEditIntent,
    handleTextEdit,
    textEditLoading,
    textEditError,
    textEditResult,
  };
};

type TextEditState = ReturnType<typeof useTextEditState>;

interface WorkflowSectionsProps {
  project: ScriptProject | null;
  detail: WorkflowInstanceDetail | null;
  assets: Asset[];
  previewContext: string;
  textEditState: TextEditState;
}

const WorkflowSections = ({
  project,
  detail,
  assets,
  previewContext,
  textEditState,
}: WorkflowSectionsProps) => {
  const {
    originalText,
    setOriginalText,
    editIntent,
    setEditIntent,
    handleTextEdit,
    textEditLoading,
    textEditError,
    textEditResult,
  } = textEditState;

  return (
    <div style={layoutStyle}>
      <SectionCard title="项目信息">
        <ProjectInfo project={project} />
      </SectionCard>
      <SectionCard title="工作流实例">
        <InstanceInfo detail={detail} />
      </SectionCard>
      <SectionCard title="历史运行">
        <HistoryRuns runs={detail?.recent_runs} />
      </SectionCard>
      <SectionCard title="上下文快照">
        <ContextSnapshot previewContext={previewContext} />
      </SectionCard>
      <SectionCard title="资产面板">
        <AssetGrid assets={assets} />
      </SectionCard>
      <SectionCard title="AI 文本编辑">
        <TextEditSection
          originalText={originalText}
          editIntent={editIntent}
          loading={textEditLoading}
          error={textEditError}
          result={textEditResult}
          onOriginalChange={setOriginalText}
          onIntentChange={setEditIntent}
          onSubmit={handleTextEdit}
        />
      </SectionCard>
    </div>
  );
};

export const ScriptWorkflowPage = () => {
  const { projectId, instanceId } = useParams();
  const { project, detail, assets, loading, error } = useScriptWorkflowData(
    projectId,
    instanceId,
  );
  const previewContext = usePreviewContext(detail?.instance?.last_context);
  const textEditState = useTextEditState(
    projectId,
    instanceId,
    detail?.instance?.last_context,
  );

  if (loading) {
    return (
      <div style={{ padding: '24px', color: 'rgba(22, 24, 35, 0.6)' }}>
        正在加载 Script Studio 内容...
      </div>
    );
  }

  return (
    <>
      <WorkflowSections
        project={project}
        detail={detail}
        assets={assets}
        previewContext={previewContext}
        textEditState={textEditState}
      />
      {error ? <ErrorBanner message={error} /> : null}
    </>
  );
};

export default ScriptWorkflowPage;
