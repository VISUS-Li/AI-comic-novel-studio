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

import { useNavigate, useParams } from 'react-router-dom';
import React, { useEffect, useMemo, useState } from 'react';

interface ScriptProject {
  id: string;
  name: string;
  type: string;
  owner_id: string;
  description: string;
  created_at: string;
  updated_at: string;
}

const cardStyle: React.CSSProperties = {
  border: '1px solid rgba(22, 24, 35, 0.12)',
  borderRadius: '12px',
  padding: '16px',
  background: '#fff',
  boxShadow: '0 1px 2px rgba(22, 24, 35, 0.08)',
  marginBottom: '16px',
};

const headerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  marginBottom: '16px',
  alignItems: 'center',
};

const listStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
  gap: '16px',
};

const badgeStyle: React.CSSProperties = {
  padding: '2px 8px',
  borderRadius: '999px',
  background: 'rgba(22, 24, 35, 0.08)',
  fontSize: '12px',
};

const projectOwnerStyle: React.CSSProperties = {
  margin: '8px 0',
  color: 'rgba(22, 24, 35, 0.6)',
};

const projectDescStyle: React.CSSProperties = {
  margin: '8px 0',
  color: 'rgba(22, 24, 35, 0.8)',
};

const formInputStyle: React.CSSProperties = {
  width: '100%',
  padding: '8px',
  borderRadius: '8px',
  border: '1px solid rgba(22, 24, 35, 0.18)',
  marginBottom: '8px',
};

const formButtonStyle: React.CSSProperties = {
  width: '100%',
  padding: '10px',
  borderRadius: '8px',
  border: 'none',
  background: '#111',
  color: '#fff',
  cursor: 'pointer',
};

const errorBoxStyle: React.CSSProperties = {
  marginBottom: '16px',
  color: '#d32f2f',
  background: '#fdecea',
  padding: '12px',
  borderRadius: '8px',
};

const refreshButtonStyle: React.CSSProperties = {
  padding: '8px 16px',
  borderRadius: '999px',
  border: '1px solid rgba(22, 24, 35, 0.12)',
  background: '#fff',
  cursor: 'pointer',
};

const fetchProjects = async () => {
  const response = await fetch('/api/script/projects', {
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error(`无法加载剧本项目（${response.status}）`);
  }
  const payload = await response.json();
  return payload?.data ?? [];
};

const useScriptProjects = () => {
  const [projects, setProjects] = useState<ScriptProject[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchProjects();
        if (!cancelled) {
          setProjects(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : '未知错误');
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return { projects, loading, error };
};

interface ProjectCardProps {
  project: ScriptProject;
  instanceId: string;
  onInstanceIdChange: (next: string) => void;
  onOpenInstance: () => void;
}

const ProjectCard = ({
  project,
  instanceId,
  onInstanceIdChange,
  onOpenInstance,
}: ProjectCardProps) => (
  <div style={cardStyle}>
    <div style={{ display: 'flex', justifyContent: 'space-between' }}>
      <h3 style={{ margin: 0 }}>{project.name}</h3>
      <span style={badgeStyle}>{project.type}</span>
    </div>
    <p style={projectOwnerStyle}>Owner: {project.owner_id}</p>
    <p style={projectDescStyle}>
      {project.description || '暂未设置描述'}
      {' · '}
      {new Date(project.created_at).toLocaleString()}
    </p>
    <div style={{ marginTop: '12px' }}>
      <input
        type="text"
        placeholder="输入 Workflow Instance ID"
        value={instanceId}
        onChange={event => onInstanceIdChange(event.target.value.trim())}
        style={formInputStyle}
      />
      <button type="button" onClick={onOpenInstance} style={formButtonStyle}>
        打开实例
      </button>
    </div>
  </div>
);

interface ProjectsContentProps {
  projects: ScriptProject[];
  loading: boolean;
  displayError: string | null;
  instanceTips: Record<string, string>;
  onInstanceIdChange: (projectId: string, next: string) => void;
  onOpenInstance: (projectId: string) => void;
  onRefresh: () => void;
}

const ProjectsContent = ({
  projects,
  loading,
  displayError,
  instanceTips,
  onInstanceIdChange,
  onOpenInstance,
  onRefresh,
}: ProjectsContentProps) => {
  const displayErrorMessage = displayError ? String(displayError) : null;

  return (
    <div style={{ padding: '24px' }}>
      <div style={headerStyle}>
        <div>
          <h2 style={{ margin: 0 }}>Script Projects</h2>
          <p style={{ margin: 0, color: 'rgba(22, 24, 35, 0.6)' }}>
            绠＄悊鍓ф湰椤圭洰骞堕€氳繃瀹炰緥鍖栧伐浣滄祦缁х画鍒涗綔銆?
          </p>
        </div>
        <button type="button" onClick={onRefresh} style={refreshButtonStyle}>
          鍒锋柊
        </button>
      </div>

      {displayErrorMessage ? (
        <div style={errorBoxStyle}>{displayErrorMessage}</div>
      ) : null}

      {loading ? (
        <div
          style={{ padding: '16px', background: '#fff', borderRadius: '12px' }}
        >
          姝ｅ湪鍔犺浇...
        </div>
      ) : projects.length ? (
        <div style={listStyle}>
          {projects.map(project => (
            <ProjectCard
              key={project.id}
              project={project}
              instanceId={instanceTips[project.id] ?? ''}
              onInstanceIdChange={next => onInstanceIdChange(project.id, next)}
              onOpenInstance={() => onOpenInstance(project.id)}
            />
          ))}
        </div>
      ) : (
        <div
          style={{ padding: '16px', background: '#fff', borderRadius: '12px' }}
        >
          暂无项目
        </div>
      )}
    </div>
  );
};

export const ScriptProjectsPage = () => {
  const { projects, loading, error: fetchError } = useScriptProjects();
  const [instanceTips, setInstanceTips] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState<string | null>(null);
  const { space_id } = useParams();
  const navigate = useNavigate();

  const basePath = useMemo(
    () => (space_id ? `/space/${space_id}` : ''),
    [space_id],
  );
  const displayError = formError ?? fetchError;

  const handleJump = (projectId: string) => {
    const targetInstance = instanceTips[projectId];
    if (!targetInstance) {
      setFormError('请先输入需要打开的 workflow 实例 ID');
      return;
    }
    setFormError(null);
    navigate(
      `${basePath}/script-projects/${projectId}/workflows/${targetInstance}`,
    );
  };

  const handleInstanceIdChange = (projectId: string, next: string) => {
    setInstanceTips(prev => ({ ...prev, [projectId]: next }));
  };

  const handleRefresh = () => window.location.reload();

  return (
    <ProjectsContent
      projects={projects}
      loading={loading}
      displayError={displayError}
      instanceTips={instanceTips}
      onInstanceIdChange={handleInstanceIdChange}
      onOpenInstance={handleJump}
      onRefresh={handleRefresh}
    />
  );
};

export default ScriptProjectsPage;
