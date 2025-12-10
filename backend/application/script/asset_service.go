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

type AssetService struct {
	repo repository.AssetRepo
}

func (s *AssetService) Create(ctx context.Context, asset *entity.Asset) error {
	return s.repo.Create(ctx, asset)
}

func (s *AssetService) ListByProject(ctx context.Context, projectID string) ([]*entity.Asset, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *AssetService) ListByWorkflowInstance(ctx context.Context, instanceID string) ([]*entity.Asset, error) {
	return s.repo.ListByInstance(ctx, instanceID)
}
