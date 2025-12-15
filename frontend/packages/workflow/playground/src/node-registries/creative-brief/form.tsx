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

import { ViewVariableType } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import { ValueExpressionInputField } from '@/node-registries/common/fields';
import { NodeConfigForm } from '@/node-registries/common/components';
import { Label, Section } from '@/form';

import { OutputsField } from '../common/fields';
import {
  BRIEF_FIELD_CONFIG,
  INPUT_PATH,
  REFERENCE_FILE_TYPES,
  type BriefFieldName,
} from './constants';

const getInputPath = (fieldName: BriefFieldName) => {
  const index = BRIEF_FIELD_CONFIG.findIndex(item => item.name === fieldName);
  return `${INPUT_PATH}.${index}.input`;
};

const stringOnly = ViewVariableType.getComplement([ViewVariableType.String]);
const constraintTypes = ViewVariableType.getComplement([
  ViewVariableType.ArrayString,
]);
const referenceTypes = ViewVariableType.getComplement([
  ViewVariableType.ArrayImage,
  ViewVariableType.ArrayAudio,
  ViewVariableType.ArrayVideo,
  ViewVariableType.ArrayDoc,
  ViewVariableType.ArrayFile,
]);

export const FormRender = () => (
  <NodeConfigForm>
    <Section
      title={I18n.t('workflow_project_brief_basic', {}, '基础信息')}
      tooltip={I18n.t(
        'workflow_project_brief_basic_tip',
        {},
        '明确作品的基础方向，便于后续节点复用',
      )}
    >
      <div className="flex flex-col gap-[12px]">
        <div className="flex flex-col gap-[4px]">
          <Label required>
            {I18n.t('workflow_project_brief_title', {}, '标题/项目名')}
          </Label>
          <ValueExpressionInputField
            name={getInputPath('title')}
            disabledTypes={stringOnly}
            inputPlaceholder={I18n.t(
              'workflow_project_brief_title_placeholder',
              {},
              '例如：太空科幻悬疑短剧',
            )}
          />
        </div>
        <div className="flex flex-col gap-[4px]">
          <Label>
            {I18n.t('workflow_project_brief_audience', {}, '目标受众')}
          </Label>
          <ValueExpressionInputField
            name={getInputPath('audience')}
            disabledTypes={stringOnly}
            inputPlaceholder={I18n.t(
              'workflow_project_brief_audience_placeholder',
              {},
              '例如：18-30 岁科幻迷、偏爱悬疑解谜',
            )}
          />
        </div>
        <div className="flex flex-col gap-[4px]">
          <Label>
            {I18n.t('workflow_project_brief_style', {}, '风格/基调')}
          </Label>
          <ValueExpressionInputField
            name={getInputPath('style')}
            disabledTypes={stringOnly}
            inputPlaceholder={I18n.t(
              'workflow_project_brief_style_placeholder',
              {},
              '例如：冷峻、写实、充满压迫感',
            )}
          />
        </div>
      </div>
    </Section>

    <Section
      title={I18n.t('workflow_project_brief_resources', {}, '约束与参考')}
      tooltip={I18n.t(
        'workflow_project_brief_resources_tip',
        {},
        '沉淀要遵守的限制、灵感清单与素材引用',
      )}
    >
      <div className="flex flex-col gap-[12px]">
        <div className="flex flex-col gap-[4px]">
          <Label>
            {I18n.t('workflow_project_brief_constraints', {}, '约束/禁用项')}
          </Label>
          <ValueExpressionInputField
            name={getInputPath('constraints')}
            disabledTypes={constraintTypes}
            inputPlaceholder={I18n.t(
              'workflow_project_brief_constraints_placeholder',
              {},
              '例如：避免出现明显暴力镜头；每段不超过 300 字',
            )}
          />
        </div>
        <div className="flex flex-col gap-[4px]">
          <Label
            tooltip={I18n.t(
              'workflow_project_brief_references_tip',
              {},
              '可添加灵感图、分镜样例、音频氛围等资源',
            )}
          >
            {I18n.t('workflow_project_brief_references', {}, '参考素材')}
          </Label>
          <ValueExpressionInputField
            name={getInputPath('references')}
            availableFileTypes={REFERENCE_FILE_TYPES}
            disabledTypes={referenceTypes}
            inputPlaceholder={I18n.t(
              'workflow_project_brief_references_placeholder',
              {},
              '支持上传或引用图片/音频/视频/文档资源列表',
            )}
          />
        </div>
      </div>
    </Section>

    <OutputsField
      title={I18n.t('workflow_detail_node_output')}
      tooltip={I18n.t(
        'workflow_project_brief_outputs_tip',
        {},
        '固定输出 brief 对象，供下游节点直接引用字段',
      )}
      name="outputs"
      topLevelReadonly
      customReadonly
    />
  </NodeConfigForm>
);
