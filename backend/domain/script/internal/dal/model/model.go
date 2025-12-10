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

package model

import (
	"time"

	"gorm.io/datatypes"
)

type ScriptProject struct {
	ID            string         `gorm:"column:id;primaryKey;type:varchar(64)"`
	Name          string         `gorm:"column:name;type:varchar(255);not null"`
	Type          string         `gorm:"column:type;type:varchar(64);not null;default:'novel'"`
	OwnerID       string         `gorm:"column:owner_id;type:varchar(64);not null"`
	Description   string         `gorm:"column:description;type:text"`
	Worldbuilding datatypes.JSON `gorm:"column:worldbuilding;type:json"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
}

func (ScriptProject) TableName() string {
	return "script_projects"
}

type WorkflowInstance struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64)"`
	ProjectID   string         `gorm:"column:project_id;type:varchar(64);not null"`
	WorkflowID  string         `gorm:"column:workflow_id;type:varchar(64);not null"`
	Name        string         `gorm:"column:name;type:varchar(255);not null"`
	Mode        string         `gorm:"column:mode;type:varchar(64);not null"`
	LastContext datatypes.JSON `gorm:"column:last_context;type:json"`
	Status      string         `gorm:"column:status;type:varchar(32);not null;default:'active'"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

func (WorkflowInstance) TableName() string {
	return "workflow_instances"
}

type WorkflowRun struct {
	ID            string         `gorm:"column:id;primaryKey;type:varchar(64)"`
	InstanceID    string         `gorm:"column:instance_id;type:varchar(64);not null"`
	StartedAt     time.Time      `gorm:"column:started_at"`
	FinishedAt    *time.Time     `gorm:"column:finished_at"`
	Status        string         `gorm:"column:status;type:varchar(32);not null"`
	InputContext  datatypes.JSON `gorm:"column:input_context;type:json"`
	OutputContext datatypes.JSON `gorm:"column:output_context;type:json"`
	ErrorMessage  string         `gorm:"column:error_message;type:text"`
}

func (WorkflowRun) TableName() string {
	return "workflow_runs"
}

type NodeOutput struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64)"`
	RunID      string         `gorm:"column:run_id;type:varchar(64);not null"`
	NodeID     string         `gorm:"column:node_id;type:varchar(128);not null"`
	OutputData datatypes.JSON `gorm:"column:output_data;type:json"`
	AssetIDs   datatypes.JSON `gorm:"column:asset_ids;type:json"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
}

func (NodeOutput) TableName() string {
	return "node_outputs"
}

type Asset struct {
	ID                 string         `gorm:"column:id;primaryKey;type:varchar(64)"`
	ProjectID          string         `gorm:"column:project_id;type:varchar(64);not null"`
	WorkflowInstanceID string         `gorm:"column:workflow_instance_id;type:varchar(64)"`
	NodeID             string         `gorm:"column:node_id;type:varchar(128)"`
	Type               string         `gorm:"column:type;type:varchar(32);not null"`
	URL                string         `gorm:"column:url;type:text;not null"`
	TextContent        string         `gorm:"column:text_content;type:text"`
	Meta               datatypes.JSON `gorm:"column:meta;type:json"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
}

func (Asset) TableName() string {
	return "assets"
}
