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

package dal

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
	"github.com/coze-dev/coze-studio/backend/domain/script/internal/dal/model"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
)

type scriptProjectRepo struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewScriptProjectRepo(db *gorm.DB, idGen idgen.IDGenerator) *scriptProjectRepo {
	return &scriptProjectRepo{
		db:    db,
		idGen: idGen,
	}
}

func (r *scriptProjectRepo) Create(ctx context.Context, project *entity.ScriptProject) error {
	if project.ID == "" {
		id, err := r.idGen.GenID(ctx)
		if err != nil {
			return err
		}
		project.ID = strconv.FormatInt(id, 10)
	}
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now
	pm := scriptProjectModelFromEntity(project)
	return r.db.WithContext(ctx).Create(pm).Error
}

func (r *scriptProjectRepo) ListByOwner(ctx context.Context, opts *entity.ScriptProjectListOption) ([]*entity.ScriptProject, error) {
	q := r.db.WithContext(ctx).Model(&model.ScriptProject{})
	if opts != nil && opts.OwnerID != "" {
		q = q.Where("owner_id = ?", opts.OwnerID)
	}
	if opts != nil && opts.Offset > 0 {
		q = q.Offset(opts.Offset)
	}
	if opts != nil && opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	var rows []*model.ScriptProject
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.ScriptProject, 0, len(rows))
	for _, row := range rows {
		result = append(result, scriptProjectEntityFromModel(row))
	}
	return result, nil
}

func (r *scriptProjectRepo) GetByID(ctx context.Context, projectID string) (*entity.ScriptProject, error) {
	var row model.ScriptProject
	if err := r.db.WithContext(ctx).
		Where("id = ?", projectID).
		First(&row).Error; err != nil {
		return nil, err
	}
	return scriptProjectEntityFromModel(&row), nil
}

type workflowInstanceRepo struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewWorkflowInstanceRepo(db *gorm.DB, idGen idgen.IDGenerator) *workflowInstanceRepo {
	return &workflowInstanceRepo{
		db:    db,
		idGen: idGen,
	}
}

func (r *workflowInstanceRepo) Create(ctx context.Context, instance *entity.WorkflowInstance) error {
	if instance.ID == "" {
		id, err := r.idGen.GenID(ctx)
		if err != nil {
			return err
		}
		instance.ID = strconv.FormatInt(id, 10)
	}
	now := time.Now()
	instance.CreatedAt = now
	instance.UpdatedAt = now
	return r.db.WithContext(ctx).Create(workflowInstanceModelFromEntity(instance)).Error
}

func (r *workflowInstanceRepo) GetByID(ctx context.Context, instanceID string) (*entity.WorkflowInstance, error) {
	var row model.WorkflowInstance
	if err := r.db.WithContext(ctx).Where("id = ?", instanceID).First(&row).Error; err != nil {
		return nil, err
	}
	return workflowInstanceEntityFromModel(&row), nil
}

func (r *workflowInstanceRepo) UpdateLastContext(ctx context.Context, instanceID string, lastContext json.RawMessage) error {
	return r.db.WithContext(ctx).Model(&model.WorkflowInstance{}).
		Where("id = ?", instanceID).
		Update("last_context", datatypes.JSON(lastContext)).Error
}

type workflowRunRepo struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewWorkflowRunRepo(db *gorm.DB, idGen idgen.IDGenerator) *workflowRunRepo {
	return &workflowRunRepo{db: db, idGen: idGen}
}

func (r *workflowRunRepo) Create(ctx context.Context, run *entity.WorkflowRun) error {
	if run.ID == "" {
		id, err := r.idGen.GenID(ctx)
		if err != nil {
			return err
		}
		run.ID = strconv.FormatInt(id, 10)
	}
	if run.StartedAt.IsZero() {
		run.StartedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(workflowRunModelFromEntity(run)).Error
}

func (r *workflowRunRepo) ListByInstance(ctx context.Context, instanceID string) ([]*entity.WorkflowRun, error) {
	var rows []*model.WorkflowRun
	if err := r.db.WithContext(ctx).
		Where("instance_id = ?", instanceID).
		Order("started_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]*entity.WorkflowRun, 0, len(rows))
	for _, row := range rows {
		list = append(list, workflowRunEntityFromModel(row))
	}
	return list, nil
}

type nodeOutputRepo struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewNodeOutputRepo(db *gorm.DB, idGen idgen.IDGenerator) *nodeOutputRepo {
	return &nodeOutputRepo{db: db, idGen: idGen}
}

func (r *nodeOutputRepo) Create(ctx context.Context, output *entity.NodeOutput) error {
	if output.ID == "" {
		id, err := r.idGen.GenID(ctx)
		if err != nil {
			return err
		}
		output.ID = strconv.FormatInt(id, 10)
	}
	output.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(nodeOutputModelFromEntity(output)).Error
}

func (r *nodeOutputRepo) ListByRun(ctx context.Context, runID string) ([]*entity.NodeOutput, error) {
	var rows []*model.NodeOutput
	if err := r.db.WithContext(ctx).Where("run_id = ?", runID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.NodeOutput, 0, len(rows))
	for _, row := range rows {
		result = append(result, nodeOutputEntityFromModel(row))
	}
	return result, nil
}

type assetRepo struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewAssetRepo(db *gorm.DB, idGen idgen.IDGenerator) *assetRepo {
	return &assetRepo{db: db, idGen: idGen}
}

func (r *assetRepo) Create(ctx context.Context, asset *entity.Asset) error {
	if asset.ID == "" {
		id, err := r.idGen.GenID(ctx)
		if err != nil {
			return err
		}
		asset.ID = strconv.FormatInt(id, 10)
	}
	now := time.Now()
	if asset.CreatedAt.IsZero() {
		asset.CreatedAt = now
	}
	asset.UpdatedAt = now
	return r.db.WithContext(ctx).Create(assetModelFromEntity(asset)).Error
}

func (r *assetRepo) ListByProject(ctx context.Context, projectID string) ([]*entity.Asset, error) {
	var rows []*model.Asset
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Asset, 0, len(rows))
	for _, row := range rows {
		result = append(result, assetEntityFromModel(row))
	}
	return result, nil
}

func (r *assetRepo) ListByInstance(ctx context.Context, instanceID string) ([]*entity.Asset, error) {
	var rows []*model.Asset
	if err := r.db.WithContext(ctx).
		Where("workflow_instance_id = ?", instanceID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*entity.Asset, 0, len(rows))
	for _, row := range rows {
		result = append(result, assetEntityFromModel(row))
	}
	return result, nil
}

func scriptProjectModelFromEntity(e *entity.ScriptProject) *model.ScriptProject {
	return &model.ScriptProject{
		ID:            e.ID,
		Name:          e.Name,
		Type:          e.Type,
		OwnerID:       e.OwnerID,
		Description:   e.Description,
		Worldbuilding: datatypes.JSON(e.Worldbuilding),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func scriptProjectEntityFromModel(m *model.ScriptProject) *entity.ScriptProject {
	return &entity.ScriptProject{
		ID:            m.ID,
		Name:          m.Name,
		Type:          m.Type,
		OwnerID:       m.OwnerID,
		Description:   m.Description,
		Worldbuilding: json.RawMessage(m.Worldbuilding),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func workflowInstanceModelFromEntity(e *entity.WorkflowInstance) *model.WorkflowInstance {
	return &model.WorkflowInstance{
		ID:          e.ID,
		ProjectID:   e.ProjectID,
		WorkflowID:  e.WorkflowID,
		Name:        e.Name,
		Mode:        e.Mode,
		LastContext: datatypes.JSON(e.LastContext),
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func workflowInstanceEntityFromModel(m *model.WorkflowInstance) *entity.WorkflowInstance {
	return &entity.WorkflowInstance{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		WorkflowID:  m.WorkflowID,
		Name:        m.Name,
		Mode:        m.Mode,
		LastContext: json.RawMessage(m.LastContext),
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func workflowRunModelFromEntity(e *entity.WorkflowRun) *model.WorkflowRun {
	return &model.WorkflowRun{
		ID:            e.ID,
		InstanceID:    e.InstanceID,
		StartedAt:     e.StartedAt,
		FinishedAt:    e.FinishedAt,
		Status:        e.Status,
		InputContext:  datatypes.JSON(e.InputContext),
		OutputContext: datatypes.JSON(e.OutputContext),
		ErrorMessage:  e.ErrorMessage,
	}
}

func workflowRunEntityFromModel(m *model.WorkflowRun) *entity.WorkflowRun {
	return &entity.WorkflowRun{
		ID:            m.ID,
		InstanceID:    m.InstanceID,
		StartedAt:     m.StartedAt,
		FinishedAt:    m.FinishedAt,
		Status:        m.Status,
		InputContext:  json.RawMessage(m.InputContext),
		OutputContext: json.RawMessage(m.OutputContext),
		ErrorMessage:  m.ErrorMessage,
	}
}

func nodeOutputModelFromEntity(e *entity.NodeOutput) *model.NodeOutput {
	return &model.NodeOutput{
		ID:         e.ID,
		RunID:      e.RunID,
		NodeID:     e.NodeID,
		OutputData: datatypes.JSON(e.OutputData),
		AssetIDs:   datatypes.JSON(e.AssetIDs),
		CreatedAt:  e.CreatedAt,
	}
}

func nodeOutputEntityFromModel(m *model.NodeOutput) *entity.NodeOutput {
	return &entity.NodeOutput{
		ID:         m.ID,
		RunID:      m.RunID,
		NodeID:     m.NodeID,
		OutputData: json.RawMessage(m.OutputData),
		AssetIDs:   json.RawMessage(m.AssetIDs),
		CreatedAt:  m.CreatedAt,
	}
}

func assetModelFromEntity(e *entity.Asset) *model.Asset {
	return &model.Asset{
		ID:                 e.ID,
		ProjectID:          e.ProjectID,
		WorkflowInstanceID: e.WorkflowInstanceID,
		NodeID:             e.NodeID,
		Type:               e.Type,
		URL:                e.URL,
		TextContent:        e.TextContent,
		Meta:               datatypes.JSON(e.Meta),
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func assetEntityFromModel(m *model.Asset) *entity.Asset {
	return &entity.Asset{
		ID:                 m.ID,
		ProjectID:          m.ProjectID,
		WorkflowInstanceID: m.WorkflowInstanceID,
		NodeID:             m.NodeID,
		Type:               m.Type,
		URL:                m.URL,
		TextContent:        m.TextContent,
		Meta:               json.RawMessage(m.Meta),
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
