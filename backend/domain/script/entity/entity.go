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

package entity

import (
	"encoding/json"
	"time"
)

type ScriptProject struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	OwnerID       string          `json:"owner_id"`
	Description   string          `json:"description"`
	Worldbuilding json.RawMessage `json:"worldbuilding"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ScriptProjectListOption struct {
	OwnerID string
	Limit   int
	Offset  int
}

type WorkflowInstance struct {
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

type WorkflowRun struct {
	ID            string          `json:"id"`
	InstanceID    string          `json:"instance_id"`
	StartedAt     time.Time       `json:"started_at"`
	FinishedAt    *time.Time      `json:"finished_at,omitempty"`
	Status        string          `json:"status"`
	InputContext  json.RawMessage `json:"input_context"`
	OutputContext json.RawMessage `json:"output_context"`
	ErrorMessage  string          `json:"error_message,omitempty"`
}

type NodeOutput struct {
	ID         string          `json:"id"`
	RunID      string          `json:"run_id"`
	NodeID     string          `json:"node_id"`
	OutputData json.RawMessage `json:"output_data"`
	AssetIDs   json.RawMessage `json:"asset_ids"`
	CreatedAt  time.Time       `json:"created_at"`
}

type Asset struct {
	ID                 string          `json:"id"`
	ProjectID          string          `json:"project_id"`
	WorkflowInstanceID string          `json:"workflow_instance_id,omitempty"`
	NodeID             string          `json:"node_id,omitempty"`
	Type               string          `json:"type"`
	URL                string          `json:"url"`
	TextContent        string          `json:"text_content,omitempty"`
	Meta               json.RawMessage `json:"meta"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
