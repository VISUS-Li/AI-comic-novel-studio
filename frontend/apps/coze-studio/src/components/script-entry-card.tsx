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

import { Link, useParams } from 'react-router-dom';
import React from 'react';

const cardStyle: React.CSSProperties = {
  display: 'block',
  padding: '16px',
  borderRadius: '12px',
  background: '#ffffff',
  border: '1px solid rgba(22, 24, 35, 0.12)',
  boxShadow: '0 1px 2px rgba(22, 24, 35, 0.08)',
  textDecoration: 'none',
  color: 'inherit',
  maxWidth: '420px',
};

const titleStyle: React.CSSProperties = {
  margin: 0,
  fontSize: '16px',
  fontWeight: 600,
};

const descStyle: React.CSSProperties = {
  margin: '8px 0 0',
  fontSize: '14px',
  color: 'rgba(22, 24, 35, 0.65)',
};

export const ScriptEntryCard = () => {
  const { space_id } = useParams();
  const target = space_id
    ? `/space/${space_id}/script-projects`
    : '/script-projects';

  return (
    <Link to={target} style={cardStyle}>
      <h3 style={titleStyle}>进入 Script Studio</h3>
      <p style={descStyle}>
        统一管理剧本项目与工作流实例，查看运行结果并继续创作。
      </p>
    </Link>
  );
};
