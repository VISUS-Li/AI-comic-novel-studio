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

package scriptapi

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	scriptmodel "github.com/coze-dev/coze-studio/backend/api/model/script"
	scriptapp "github.com/coze-dev/coze-studio/backend/application/script"
	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
)

func respondWithError(c *app.RequestContext, status int, err error) {
	c.String(status, err.Error())
}

func CreateScriptProject(ctx context.Context, c *app.RequestContext) {
	var req scriptmodel.CreateScriptProjectRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	project := &entity.ScriptProject{
		Name:          req.Name,
		Type:          req.Type,
		OwnerID:       req.OwnerID,
		Description:   req.Description,
		Worldbuilding: req.Worldbuilding,
	}
	if project.Type == "" {
		project.Type = "novel"
	}

	if err := scriptapp.ScriptSVC.ScriptProject.Create(ctx, project); err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.CreateScriptProjectResponse{
		Data: toScriptProjectVO(project),
	})
}

func ListScriptProjects(ctx context.Context, c *app.RequestContext) {
	var req scriptmodel.ListScriptProjectsRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	opts := &entity.ScriptProjectListOption{
		OwnerID: req.OwnerID,
		Limit:   req.Limit,
		Offset:  req.Offset,
	}

	projects, err := scriptapp.ScriptSVC.ScriptProject.ListByOwner(ctx, opts)
	if err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	details := make([]*scriptmodel.ScriptProjectVO, 0, len(projects))
	for _, proj := range projects {
		details = append(details, toScriptProjectVO(proj))
	}

	c.JSON(consts.StatusOK, &scriptmodel.ListScriptProjectsResponse{Data: details})
}

func GetScriptProject(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("projectId")
	if projectID == "" {
		respondWithError(c, consts.StatusBadRequest, errProjectIDRequired())
		return
	}

	project, err := scriptapp.ScriptSVC.ScriptProject.Get(ctx, projectID)
	if err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.GetScriptProjectResponse{
		Data: toScriptProjectVO(project),
	})
}

func CreateWorkflowInstance(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("projectId")
	if projectID == "" {
		respondWithError(c, consts.StatusBadRequest, errProjectIDRequired())
		return
	}

	var req scriptmodel.CreateWorkflowInstanceRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	instance := &entity.WorkflowInstance{
		ProjectID:  projectID,
		WorkflowID: req.WorkflowID,
		Name:       req.Name,
		Mode:       req.Mode,
		Status:     "active",
	}
	if instance.Mode == "" {
		instance.Mode = "outline_first"
	}

	if err := scriptapp.ScriptSVC.WorkflowInstance.Create(ctx, instance); err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.CreateWorkflowInstanceResponse{
		Data: toWorkflowInstanceVO(instance),
	})
}

func GetWorkflowInstanceDetail(ctx context.Context, c *app.RequestContext) {
	instanceID := c.Param("instanceId")
	if instanceID == "" {
		respondWithError(c, consts.StatusBadRequest, errInstanceIDRequired())
		return
	}

	instance, err := scriptapp.ScriptSVC.WorkflowInstance.Get(ctx, instanceID)
	if err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.GetWorkflowInstanceDetailResponse{
		Data: &scriptmodel.WorkflowInstanceDetail{
			Instance:   toWorkflowInstanceVO(instance),
			RecentRuns: nil,
		},
	})
}

func ListAssets(ctx context.Context, c *app.RequestContext) {
	var req scriptmodel.ListAssetsRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	if req.ProjectID == "" && req.InstanceID == "" {
		respondWithError(c, consts.StatusBadRequest, errProjectOrInstanceRequired())
		return
	}

	var (
		assets []*entity.Asset
		err    error
	)
	if req.InstanceID != "" {
		assets, err = scriptapp.ScriptSVC.Asset.ListByWorkflowInstance(ctx, req.InstanceID)
	} else {
		assets, err = scriptapp.ScriptSVC.Asset.ListByProject(ctx, req.ProjectID)
	}
	if err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	result := make([]*scriptmodel.AssetVO, 0, len(assets))
	for _, asset := range assets {
		result = append(result, toAssetVO(asset))
	}

	c.JSON(consts.StatusOK, &scriptmodel.ListAssetsResponse{Data: result})
}

func RunWorkflowInstance(ctx context.Context, c *app.RequestContext) {
	instanceID := c.Param("instanceId")
	if instanceID == "" {
		respondWithError(c, consts.StatusBadRequest, errInstanceIDRequired())
		return
	}

	var req scriptmodel.RunWorkflowInstanceRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	run := &entity.WorkflowRun{
		InstanceID:    instanceID,
		Status:        ternaryNonEmpty(req.Status, "running"),
		InputContext:  req.InputContext,
		OutputContext: req.OutputContext,
		StartedAt:     time.Now(),
	}
	if run.Status != "running" {
		now := time.Now()
		run.FinishedAt = &now
	}

	nodeOutputs := make([]*entity.NodeOutput, 0, len(req.NodeOutputs))
	for _, output := range req.NodeOutputs {
		node := &entity.NodeOutput{
			RunID:      run.ID,
			NodeID:     output.NodeID,
			OutputData: output.OutputData,
		}
		if len(output.AssetIDs) > 0 {
			encoded, err := json.Marshal(output.AssetIDs)
			if err != nil {
				respondWithError(c, consts.StatusBadRequest, err)
				return
			}
			node.AssetIDs = encoded
		}
		nodeOutputs = append(nodeOutputs, node)
	}

	lastContext := req.LastContext
	if len(lastContext) == 0 {
		lastContext = req.OutputContext
	}

	if err := scriptapp.ScriptSVC.WorkflowRun.RunWorkflowInstance(ctx, run, nodeOutputs, lastContext); err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.RunWorkflowInstanceResponse{
		RunID: run.ID,
	})
}

func RunNodeInInstance(ctx context.Context, c *app.RequestContext) {
	instanceID := c.Param("instanceId")
	if instanceID == "" {
		respondWithError(c, consts.StatusBadRequest, errInstanceIDRequired())
		return
	}

	nodeID := c.Param("nodeId")
	if nodeID == "" {
		respondWithError(c, consts.StatusBadRequest, &APIError{"node_id is required"})
		return
	}

	var req scriptmodel.RunNodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		respondWithError(c, consts.StatusBadRequest, err)
		return
	}

	run := &entity.WorkflowRun{
		InstanceID:    instanceID,
		Status:        ternaryNonEmpty(req.Status, "running"),
		InputContext:  nil,
		OutputContext: req.OutputData,
		StartedAt:     time.Now(),
	}
	if run.Status != "running" {
		now := time.Now()
		run.FinishedAt = &now
	}

	output := &entity.NodeOutput{
		NodeID:     nodeID,
		OutputData: req.OutputData,
	}
	if len(req.AssetIDs) > 0 {
		encoded, err := json.Marshal(req.AssetIDs)
		if err != nil {
			respondWithError(c, consts.StatusBadRequest, err)
			return
		}
		output.AssetIDs = encoded
	}

	lastContext := req.LastContext
	if len(lastContext) == 0 {
		lastContext = req.OutputData
	}

	if err := scriptapp.ScriptSVC.WorkflowRun.RunNodeInInstance(ctx, run, output, lastContext); err != nil {
		respondWithError(c, consts.StatusInternalServerError, err)
		return
	}

	c.JSON(consts.StatusOK, &scriptmodel.RunNodeResponse{RunID: run.ID})
}

func toScriptProjectVO(p *entity.ScriptProject) *scriptmodel.ScriptProjectVO {
	if p == nil {
		return nil
	}
	return &scriptmodel.ScriptProjectVO{
		ID:            p.ID,
		Name:          p.Name,
		Type:          p.Type,
		OwnerID:       p.OwnerID,
		Description:   p.Description,
		Worldbuilding: p.Worldbuilding,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toWorkflowInstanceVO(i *entity.WorkflowInstance) *scriptmodel.WorkflowInstanceVO {
	if i == nil {
		return nil
	}
	return &scriptmodel.WorkflowInstanceVO{
		ID:          i.ID,
		ProjectID:   i.ProjectID,
		WorkflowID:  i.WorkflowID,
		Name:        i.Name,
		Mode:        i.Mode,
		LastContext: i.LastContext,
		Status:      i.Status,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func toAssetVO(a *entity.Asset) *scriptmodel.AssetVO {
	if a == nil {
		return nil
	}
	return &scriptmodel.AssetVO{
		ID:                 a.ID,
		ProjectID:          a.ProjectID,
		WorkflowInstanceID: a.WorkflowInstanceID,
		NodeID:             a.NodeID,
		Type:               a.Type,
		URL:                a.URL,
		TextContent:        a.TextContent,
		Meta:               a.Meta,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func errProjectIDRequired() error {
	return &APIError{"project_id is required"}
}

func errInstanceIDRequired() error {
	return &APIError{"instance_id is required"}
}

func errProjectOrInstanceRequired() error {
	return &APIError{"project_id or instance_id is required"}
}

type APIError struct {
	msg string
}

func (e *APIError) Error() string {
	return e.msg
}

func ternaryNonEmpty(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
