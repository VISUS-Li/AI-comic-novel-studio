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
	"encoding/json"

	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
	"github.com/coze-dev/coze-studio/backend/domain/script/repository"
)

type WorkflowInstanceService struct {
	repo repository.WorkflowInstanceRepo
}

func (s *WorkflowInstanceService) Create(ctx context.Context, inst *entity.WorkflowInstance) error {
	return s.repo.Create(ctx, inst)
}

func (s *WorkflowInstanceService) Get(ctx context.Context, instanceID string) (*entity.WorkflowInstance, error) {
	return s.repo.GetByID(ctx, instanceID)
}

func (s *WorkflowInstanceService) UpdateLastContext(ctx context.Context, instanceID string, lastContext json.RawMessage) error {
	return s.repo.UpdateLastContext(ctx, instanceID, lastContext)
}
