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

package repository

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
	"github.com/coze-dev/coze-studio/backend/domain/script/internal/dal"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
)

func NewScriptProjectRepo(db *gorm.DB, idGen idgen.IDGenerator) ScriptProjectRepo {
	return dal.NewScriptProjectRepo(db, idGen)
}

func NewWorkflowInstanceRepo(db *gorm.DB, idGen idgen.IDGenerator) WorkflowInstanceRepo {
	return dal.NewWorkflowInstanceRepo(db, idGen)
}

func NewWorkflowRunRepo(db *gorm.DB, idGen idgen.IDGenerator) WorkflowRunRepo {
	return dal.NewWorkflowRunRepo(db, idGen)
}

func NewNodeOutputRepo(db *gorm.DB, idGen idgen.IDGenerator) NodeOutputRepo {
	return dal.NewNodeOutputRepo(db, idGen)
}

func NewAssetRepo(db *gorm.DB, idGen idgen.IDGenerator) AssetRepo {
	return dal.NewAssetRepo(db, idGen)
}

type ScriptProjectRepo interface {
	Create(ctx context.Context, project *entity.ScriptProject) error
	ListByOwner(ctx context.Context, opts *entity.ScriptProjectListOption) ([]*entity.ScriptProject, error)
	GetByID(ctx context.Context, projectID string) (*entity.ScriptProject, error)
}

type WorkflowInstanceRepo interface {
	Create(ctx context.Context, instance *entity.WorkflowInstance) error
	GetByID(ctx context.Context, instanceID string) (*entity.WorkflowInstance, error)
	UpdateLastContext(ctx context.Context, instanceID string, lastContext json.RawMessage) error
}

type WorkflowRunRepo interface {
	Create(ctx context.Context, run *entity.WorkflowRun) error
	ListByInstance(ctx context.Context, instanceID string) ([]*entity.WorkflowRun, error)
}

type NodeOutputRepo interface {
	Create(ctx context.Context, output *entity.NodeOutput) error
	ListByRun(ctx context.Context, runID string) ([]*entity.NodeOutput, error)
}

type AssetRepo interface {
	Create(ctx context.Context, asset *entity.Asset) error
	ListByProject(ctx context.Context, projectID string) ([]*entity.Asset, error)
	ListByInstance(ctx context.Context, instanceID string) ([]*entity.Asset, error)
}
