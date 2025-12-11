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

package script

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
	"github.com/coze-dev/coze-studio/backend/domain/script/repository"
)

type WorkflowRunService struct {
	runRepo              repository.WorkflowRunRepo
	nodeOutputRepo       repository.NodeOutputRepo
	workflowInstanceRepo repository.WorkflowInstanceRepo
}

func (s *WorkflowRunService) RunWorkflowInstance(ctx context.Context, run *entity.WorkflowRun, nodeOutputs []*entity.NodeOutput, lastContext []byte) error {
	if err := s.runRepo.Create(ctx, run); err != nil {
		return err
	}

	for _, output := range nodeOutputs {
		output.RunID = run.ID
		if err := s.nodeOutputRepo.Create(ctx, output); err != nil {
			return err
		}
	}

	if len(lastContext) > 0 {
		if err := s.workflowInstanceRepo.UpdateLastContext(ctx, run.InstanceID, lastContext); err != nil {
			return err
		}
	}

	return nil
}

func (s *WorkflowRunService) RunNodeInInstance(ctx context.Context, run *entity.WorkflowRun, output *entity.NodeOutput, lastContext []byte) error {
	if err := s.runRepo.Create(ctx, run); err != nil {
		return err
	}

	output.RunID = run.ID
	if err := s.nodeOutputRepo.Create(ctx, output); err != nil {
		return err
	}

	if len(lastContext) > 0 {
		return s.workflowInstanceRepo.UpdateLastContext(ctx, run.InstanceID, lastContext)
	}

	return nil
}
