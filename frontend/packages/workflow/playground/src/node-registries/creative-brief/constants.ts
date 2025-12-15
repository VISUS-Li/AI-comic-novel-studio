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

import { nanoid } from 'nanoid';
import {
  ValueExpressionType,
  ViewVariableType,
  type InputValueVO,
  type OutputValueVO,
} from '@coze-workflow/base';

export const INPUT_PATH = 'inputs.inputParameters';

export const REFERENCE_FILE_TYPES = [
  ViewVariableType.Image,
  ViewVariableType.Audio,
  ViewVariableType.Video,
  ViewVariableType.Doc,
  ViewVariableType.File,
];

export const BRIEF_FIELD_CONFIG = [
  { name: 'title', defaultContent: '' },
  { name: 'audience', defaultContent: '' },
  { name: 'style', defaultContent: '' },
  { name: 'constraints', defaultContent: [] as unknown[] },
  { name: 'references', defaultContent: [] as unknown[] },
] as const;

export type BriefFieldName = (typeof BRIEF_FIELD_CONFIG)[number]['name'];

const createInput = (
  name: BriefFieldName,
  content: unknown,
  rawMetaType?: ViewVariableType,
): InputValueVO => ({
  key: nanoid(),
  name,
  input: {
    type: ValueExpressionType.LITERAL,
    content: typeof content === 'string' ? content : '',
    rawMeta: rawMetaType ? { type: rawMetaType } : undefined,
  },
});

export const createDefaultInputs = (): InputValueVO[] => [
  createInput('title', ''),
  createInput('audience', ''),
  createInput('style', ''),
  createInput('constraints', [], ViewVariableType.ArrayString),
  createInput('references', [], ViewVariableType.ArrayFile),
];

export const createDefaultOutputs = (): OutputValueVO[] => [
  {
    key: nanoid(),
    name: 'brief',
    type: ViewVariableType.Object,
    children: [
      { key: nanoid(), name: 'title', type: ViewVariableType.String },
      { key: nanoid(), name: 'audience', type: ViewVariableType.String },
      { key: nanoid(), name: 'style', type: ViewVariableType.String },
      {
        key: nanoid(),
        name: 'constraints',
        type: ViewVariableType.ArrayString,
      },
      { key: nanoid(), name: 'references', type: ViewVariableType.ArrayFile },
    ],
  },
];
