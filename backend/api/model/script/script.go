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
	"encoding/json"
	"time"
)

type ScriptProjectVO struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	OwnerID       string          `json:"owner_id"`
	Description   string          `json:"description"`
	Worldbuilding json.RawMessage `json:"worldbuilding"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type CreateScriptProjectRequest struct {
	Name          string          `json:"name" binding:"required"`
	Type          string          `json:"type"`
	OwnerID       string          `json:"owner_id" binding:"required"`
	Description   string          `json:"description"`
	Worldbuilding json.RawMessage `json:"worldbuilding"`
}

type CreateScriptProjectResponse struct {
	Data *ScriptProjectVO `json:"data"`
}

type ListScriptProjectsRequest struct {
	OwnerID string `form:"owner_id"`
	Limit   int    `form:"limit"`
	Offset  int    `form:"offset"`
}

type ListScriptProjectsResponse struct {
	Data []*ScriptProjectVO `json:"data"`
}

type GetScriptProjectResponse struct {
	Data *ScriptProjectVO `json:"data"`
}

type WorkflowInstanceVO struct {
	ID          string          `json:"id"`
	ProjectID   string          `json:"project_id"`
	WorkflowID  string          `json:"workflow_id"`
	Name        string          `json:"name"`
	Mode        string          `json:"mode"`
	LastContext json.RawMessage `json:"last_context"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateWorkflowInstanceRequest struct {
	WorkflowID string `json:"workflow_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Mode       string `json:"mode"`
}

type CreateWorkflowInstanceResponse struct {
	Data *WorkflowInstanceVO `json:"data"`
}

type WorkflowRunVO struct {
	ID         string     `json:"id"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Status     string     `json:"status"`
}

type WorkflowInstanceDetail struct {
	Instance   *WorkflowInstanceVO `json:"instance"`
	RecentRuns []*WorkflowRunVO    `json:"recent_runs"`
}

type GetWorkflowInstanceDetailResponse struct {
	Data *WorkflowInstanceDetail `json:"data"`
}

type ListAssetsRequest struct {
	ProjectID  string `form:"project_id"`
	InstanceID string `form:"instance_id"`
}

type AssetVO struct {
	ID                 string          `json:"id"`
	ProjectID          string          `json:"project_id"`
	WorkflowInstanceID string          `json:"workflow_instance_id"`
	NodeID             string          `json:"node_id"`
	Type               string          `json:"type"`
	URL                string          `json:"url"`
	TextContent        string          `json:"text_content"`
	Meta               json.RawMessage `json:"meta"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type ListAssetsResponse struct {
	Data []*AssetVO `json:"data"`
}

type RunWorkflowInstanceRequest struct {
	Status        string            `json:"status"`
	InputContext  json.RawMessage   `json:"input_context"`
	OutputContext json.RawMessage   `json:"output_context"`
	LastContext   json.RawMessage   `json:"last_context"`
	NodeOutputs   []*NodeOutputInfo `json:"node_outputs"`
}

type RunWorkflowInstanceResponse struct {
	RunID string `json:"run_id"`
}

type NodeOutputInfo struct {
	NodeID     string          `json:"node_id" binding:"required"`
	OutputData json.RawMessage `json:"output_data"`
	AssetIDs   []string        `json:"asset_ids"`
}

type RunNodeRequest struct {
	Status      string          `json:"status"`
	OutputData  json.RawMessage `json:"output_data"`
	LastContext json.RawMessage `json:"last_context"`
	AssetIDs    []string        `json:"asset_ids"`
}

type RunNodeResponse struct {
	RunID string `json:"run_id"`
}
