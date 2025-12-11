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

import React from 'react';

import type {
  Asset,
  ScriptProject,
  TextEditResult,
  WorkflowInstanceDetail,
} from './script-workflow';

const cardStyle: React.CSSProperties = {
  padding: '16px',
  borderRadius: '12px',
  background: '#fff',
  border: '1px solid rgba(22, 24, 35, 0.12)',
  boxShadow: '0 1px 2px rgba(22, 24, 35, 0.08)',
};

const gridStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
  gap: '12px',
  marginTop: '12px',
};

const runItemStyle: React.CSSProperties = {
  border: '1px solid rgba(22, 24, 35, 0.08)',
  borderRadius: '8px',
  padding: '8px',
};

const assetCardStyle: React.CSSProperties = {
  border: '1px solid rgba(22, 24, 35, 0.08)',
  borderRadius: '8px',
  padding: '12px',
};

const previewPreStyle: React.CSSProperties = {
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  maxHeight: '240px',
  overflow: 'auto',
  margin: 0,
  fontSize: '13px',
  color: 'rgba(22, 24, 35, 0.8)',
};

const textAreaStyle: React.CSSProperties = {
  width: '100%',
  borderRadius: '8px',
  border: '1px solid rgba(22, 24, 35, 0.18)',
  padding: '10px',
  resize: 'vertical',
};

const textInputStyle: React.CSSProperties = {
  width: '100%',
  marginTop: '8px',
  padding: '8px',
  borderRadius: '8px',
  border: '1px solid rgba(22, 24, 35, 0.18)',
};

const textEditButtonStyle: React.CSSProperties = {
  width: '100%',
  marginTop: '12px',
  padding: '10px',
  borderRadius: '8px',
  border: 'none',
  background: '#111',
  color: '#fff',
  cursor: 'pointer',
};

const textResultPreStyle: React.CSSProperties = {
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-word',
  background: '#f7f7f7',
  padding: '10px',
  borderRadius: '8px',
  fontSize: '13px',
};

const paragraphStyle: React.CSSProperties = {
  margin: '8px 0',
};

const tightParagraphStyle: React.CSSProperties = {
  margin: '4px 0',
};

const errorBannerStyle: React.CSSProperties = {
  borderColor: '#ffb4a9',
  color: '#b71c1c',
};

interface SectionCardProps {
  title: string;
  children: React.ReactNode;
}

export const SectionCard = ({ title, children }: SectionCardProps) => (
  <div style={cardStyle}>
    <h3>{title}</h3>
    {children}
  </div>
);

export const ProjectInfo = ({ project }: { project: ScriptProject | null }) =>
  project ? (
    <>
      <p style={paragraphStyle}>名称：{project.name}</p>
      <p style={paragraphStyle}>类型：{project.type}</p>
      <p style={paragraphStyle}>Owner：{project.owner_id}</p>
      <p style={paragraphStyle}>
        创建时间：{new Date(project.created_at).toLocaleString()}
      </p>
      <p style={paragraphStyle}>{project.description || '暂无描述'}</p>
    </>
  ) : (
    <p>无项目信息</p>
  );

export const InstanceInfo = ({
  detail,
}: {
  detail: WorkflowInstanceDetail | null;
}) =>
  detail?.instance ? (
    <>
      <p style={tightParagraphStyle}>名称：{detail.instance.name}</p>
      <p style={tightParagraphStyle}>模式：{detail.instance.mode}</p>
      <p style={tightParagraphStyle}>Status：{detail.instance.status}</p>
      <p style={tightParagraphStyle}>
        Workflow ID：{detail.instance.workflow_id}
      </p>
      <p style={tightParagraphStyle}>
        更新时间：{new Date(detail.instance.updated_at).toLocaleString()}
      </p>
    </>
  ) : (
    <p>暂无实例数据</p>
  );

export const HistoryRuns = ({
  runs,
}: {
  runs?: WorkflowInstanceDetail['recent_runs'];
}) =>
  runs?.length ? (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
      {runs.map(run => (
        <div key={run.id} style={runItemStyle}>
          <p style={tightParagraphStyle}>Run ID: {run.id}</p>
          <p style={tightParagraphStyle}>Status: {run.status}</p>
          <p style={tightParagraphStyle}>
            Start: {new Date(run.started_at).toLocaleString()}
          </p>
        </div>
      ))}
    </div>
  ) : (
    <p>无历史运行</p>
  );

export const ContextSnapshot = ({
  previewContext,
}: {
  previewContext: string;
}) => <pre style={previewPreStyle}>{previewContext}</pre>;

export const AssetGrid = ({ assets }: { assets: Asset[] }) =>
  assets.length === 0 ? (
    <p>未生成资产</p>
  ) : (
    <div style={gridStyle}>
      {assets.map(asset => (
        <div key={asset.id} style={assetCardStyle}>
          <p style={{ margin: '4px 0', fontWeight: 600 }}>{asset.type}</p>
          <p style={{ margin: '4px 0', color: 'rgba(22, 24, 35, 0.6)' }}>
            Node: {asset.node_id || 'N/A'}
          </p>
          {asset.text_content ? (
            <p style={{ margin: '4px 0' }}>{asset.text_content}</p>
          ) : null}
          <a href={asset.url} target="_blank" rel="noreferrer">
            查看资源
          </a>
        </div>
      ))}
    </div>
  );

interface TextEditSectionProps {
  originalText: string;
  editIntent: string;
  loading: boolean;
  error: string | null;
  result: TextEditResult | null;
  onOriginalChange: (next: string) => void;
  onIntentChange: (next: string) => void;
  onSubmit: () => void;
}

export const TextEditSection = ({
  originalText,
  editIntent,
  loading,
  error,
  result,
  onOriginalChange,
  onIntentChange,
  onSubmit,
}: TextEditSectionProps) => (
  <>
    <p style={paragraphStyle}>
      选中片段 + 编辑意图 → 触发后端 `text-edit` API。
    </p>
    <textarea
      value={originalText}
      onChange={event => onOriginalChange(event.target.value)}
      rows={5}
      style={textAreaStyle}
    />
    <input
      value={editIntent}
      onChange={event => onIntentChange(event.target.value)}
      placeholder="输入编辑意图（例如：润色/转为剧本格式）"
      style={textInputStyle}
    />
    <button
      type="button"
      onClick={onSubmit}
      disabled={loading || !originalText}
      style={textEditButtonStyle}
    >
      {loading ? 'AI 编辑中...' : '调用 AI 编辑'}
    </button>
    {error ? (
      <p style={{ marginTop: '8px', color: '#d32f2f' }}>{String(error)}</p>
    ) : null}
    {result ? (
      <div style={{ marginTop: '12px' }}>
        <p style={{ margin: '4px 0', fontWeight: 600 }}>编辑结果</p>
        <pre style={textResultPreStyle}>{result.edited_text}</pre>
        {result.diff ? (
          <>
            <p style={{ margin: '8px 0 4px', fontWeight: 600 }}>Diff</p>
            <pre style={textResultPreStyle}>{result.diff}</pre>
          </>
        ) : null}
      </div>
    ) : null}
  </>
);

export const ErrorBanner = ({ message }: { message: string }) => (
  <div style={{ ...cardStyle, ...errorBannerStyle }}>
    <strong>错误</strong>
    <p>{message}</p>
  </div>
);
