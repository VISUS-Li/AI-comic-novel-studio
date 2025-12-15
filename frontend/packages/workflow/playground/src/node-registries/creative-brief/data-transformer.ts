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

import { type InputValueVO, type NodeDataDTO } from '@coze-workflow/base';

import { type FormData } from './types';
import { createDefaultInputs, createDefaultOutputs } from './constants';

const normalizeInputs = (params?: InputValueVO[] | null): InputValueVO[] => {
  const defaultInputs = createDefaultInputs();
  if (!params?.length) {
    return defaultInputs;
  }

  const paramMap = new Map(params.map(item => [item?.name, item] as const));

  return defaultInputs.map(item => {
    const match = paramMap.get(item.name);
    if (!match) {
      return item;
    }

    return {
      ...item,
      ...match,
      name: item.name,
      key: match.key ?? item.key,
      input: {
        ...item.input,
        ...match.input,
        rawMeta: match.input?.rawMeta ?? item.input?.rawMeta,
      },
    };
  });
};

export const transformOnInit = (value: NodeDataDTO) => ({
  ...(value ?? {}),
  inputs: {
    ...(value?.inputs ?? {}),
    inputParameters: normalizeInputs(
      (value?.inputs as { inputParameters?: InputValueVO[] })?.inputParameters,
    ),
  },
  outputs:
    Array.isArray(value?.outputs) && value.outputs.length
      ? (value.outputs as NodeDataDTO['outputs'])
      : createDefaultOutputs(),
});

export const transformOnSubmit = (value: FormData): NodeDataDTO =>
  value as unknown as NodeDataDTO;
