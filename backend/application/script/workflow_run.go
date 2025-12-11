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
	"fmt"
	"strconv"

	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	"github.com/coze-dev/coze-studio/backend/domain/script/entity"
	"github.com/coze-dev/coze-studio/backend/domain/script/repository"
	workflowEntity "github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type WorkflowRunService struct {
	runRepo              repository.WorkflowRunRepo
	nodeOutputRepo       repository.NodeOutputRepo
	workflowInstanceRepo repository.WorkflowInstanceRepo
	assetSvc             *AssetService
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

func (s *WorkflowRunService) ExecuteWorkflowInstance(ctx context.Context, instanceID string, input json.RawMessage) (string, error) {
	instance, err := s.workflowInstanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return "", err
	}

	if instance.WorkflowID == "" {
		return "", fmt.Errorf("workflow_id missing for workflow instance %s", instanceID)
	}

	workflowID, err := strconv.ParseInt(instance.WorkflowID, 10, 64)
	if err != nil {
		return "", err
	}

	parameters, err := resolveWorkflowInput(input)
	if err != nil {
		return "", err
	}
	if len(parameters) == 0 {
		parameters, err = resolveWorkflowInput(instance.LastContext)
		if err != nil {
			return "", err
		}
	}

	exeCfg := workflowModel.ExecuteConfig{
		ID:            workflowID,
		From:          workflowModel.FromLatestVersion,
		Operator:      0,
		Mode:          workflowModel.ExecuteModeRelease,
		ConnectorID:   consts.CozeConnectorID,
		ConnectorUID:  instanceID,
		TaskType:      workflowModel.TaskTypeForeground,
		SyncPattern:   workflowModel.SyncPatternSync,
		InputFailFast: true,
		BizType:       workflowModel.BizTypeWorkflow,
	}

	exe, _, err := appworkflow.GetWorkflowDomainSVC().SyncExecute(ctx, exeCfg, parameters)
	if err != nil {
		return "", err
	}

	run := &entity.WorkflowRun{
		InstanceID:    instanceID,
		Status:        mapWorkflowStatus(exe.Status),
		InputContext:  stringPtrToRawMessage(exe.Input),
		OutputContext: stringPtrToRawMessage(exe.Output),
		StartedAt:     exe.CreatedAt,
		FinishedAt:    exe.UpdatedAt,
	}

	nodeOutputs := convertNodeExecutions(exe.NodeExecutions)
	if err := s.persistAssets(ctx, instance, exe.NodeExecutions, nodeOutputs); err != nil {
		return "", err
	}
	var lastContext []byte
	if exe.Output != nil {
		lastContext = []byte(*exe.Output)
	}

	if err := s.RunWorkflowInstance(ctx, run, nodeOutputs, lastContext); err != nil {
		return "", err
	}

	return run.ID, nil
}

func (s *WorkflowRunService) ExecuteNodeInInstance(ctx context.Context, instanceID, nodeID string, input json.RawMessage) (string, error) {
	instance, err := s.workflowInstanceRepo.GetByID(ctx, instanceID)
	if err != nil {
		return "", err
	}

	if instance.WorkflowID == "" {
		return "", fmt.Errorf("workflow_id missing for workflow instance %s", instanceID)
	}

	workflowID, err := strconv.ParseInt(instance.WorkflowID, 10, 64)
	if err != nil {
		return "", err
	}

	parameters, err := resolveWorkflowInput(input)
	if err != nil {
		return "", err
	}
	if len(parameters) == 0 {
		parameters, err = resolveWorkflowInput(instance.LastContext)
		if err != nil {
			return "", err
		}
	}

	exeCfg := workflowModel.ExecuteConfig{
		ID:            workflowID,
		From:          workflowModel.FromLatestVersion,
		Operator:      0,
		Mode:          workflowModel.ExecuteModeNodeDebug,
		ConnectorID:   consts.CozeConnectorID,
		ConnectorUID:  instanceID,
		TaskType:      workflowModel.TaskTypeForeground,
		SyncPattern:   workflowModel.SyncPatternAsync,
		InputFailFast: true,
		BizType:       workflowModel.BizTypeWorkflow,
	}

	executeID, err := appworkflow.GetWorkflowDomainSVC().AsyncExecuteNode(ctx, nodeID, exeCfg, parameters)
	if err != nil {
		return "", err
	}

	nodeExe, innerNodeExe, err := appworkflow.GetWorkflowDomainSVC().GetNodeExecution(ctx, executeID, nodeID)
	if err != nil {
		return "", err
	}

	executions := []*workflowEntity.NodeExecution{nodeExe}
	if innerNodeExe != nil {
		executions = append(executions, innerNodeExe)
	}

	nodeOutputs := convertNodeExecutions(executions)
	if err := s.persistAssets(ctx, instance, executions, nodeOutputs); err != nil {
		return "", err
	}

	run := &entity.WorkflowRun{
		InstanceID:    instanceID,
		Status:        mapNodeStatus(nodeExe.Status),
		InputContext:  stringPtrToRawMessage(nodeExe.Input),
		OutputContext: stringPtrToRawMessage(nodeExe.Output),
		StartedAt:     nodeExe.CreatedAt,
		FinishedAt:    nodeExe.UpdatedAt,
	}

	var lastContext []byte
	if nodeExe.Output != nil {
		lastContext = []byte(*nodeExe.Output)
	} else if nodeExe.Input != nil {
		lastContext = []byte(*nodeExe.Input)
	}

	if err := s.RunWorkflowInstance(ctx, run, nodeOutputs, lastContext); err != nil {
		return "", err
	}

	return run.ID, nil
}

func resolveWorkflowInput(raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var decoded any
	if err := sonic.Unmarshal(raw, &decoded); err != nil {
		return map[string]any{
			"payload": string(raw),
		}, nil
	}
	if m, ok := decoded.(map[string]any); ok {
		return m, nil
	}
	return map[string]any{
		"value": decoded,
	}, nil
}

func convertNodeExecutions(nodes []*workflowEntity.NodeExecution) []*entity.NodeOutput {
	result := make([]*entity.NodeOutput, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		data := normalizeNodeOutput(node.Output)
		if len(data) == 0 {
			continue
		}
		result = append(result, &entity.NodeOutput{
			NodeID:     node.NodeID,
			OutputData: data,
		})
	}
	return result
}

func (s *WorkflowRunService) persistAssets(ctx context.Context, instance *entity.WorkflowInstance, nodes []*workflowEntity.NodeExecution, outputs []*entity.NodeOutput) error {
	if s.assetSvc == nil || instance == nil {
		return nil
	}

	outputMap := make(map[string]*entity.NodeOutput, len(outputs))
	for _, output := range outputs {
		outputMap[output.NodeID] = output
	}

	for _, node := range nodes {
		if node == nil {
			continue
		}
		output, ok := outputMap[node.NodeID]
		if !ok {
			continue
		}
		text := string(output.OutputData)
		meta := map[string]any{
			"node_type": node.NodeType,
			"node_name": node.NodeName,
		}
		metaBytes, err := json.Marshal(meta)
		if err != nil {
			metaBytes = []byte("{}")
		}
		asset := &entity.Asset{
			ProjectID:          instance.ProjectID,
			WorkflowInstanceID: instance.ID,
			NodeID:             node.NodeID,
			Type:               "text",
			URL:                fmt.Sprintf("text://%s/%s", instance.ID, node.NodeID),
			TextContent:        text,
			Meta:               json.RawMessage(metaBytes),
		}
		if err := s.assetSvc.Create(ctx, asset); err != nil {
			return err
		}
	}

	return nil
}

func normalizeNodeOutput(value *string) json.RawMessage {
	if value == nil || len(*value) == 0 {
		return nil
	}
	raw := []byte(*value)
	if json.Valid(raw) {
		return json.RawMessage(raw)
	}
	encoded, err := json.Marshal(*value)
	if err != nil {
		return json.RawMessage([]byte(`""`))
	}
	return json.RawMessage(encoded)
}

func stringPtrToRawMessage(v *string) json.RawMessage {
	if v == nil {
		return nil
	}
	return json.RawMessage([]byte(*v))
}

func mapWorkflowStatus(status workflowEntity.WorkflowExecuteStatus) string {
	switch status {
	case workflowEntity.WorkflowRunning:
		return "running"
	case workflowEntity.WorkflowSuccess:
		return "success"
	case workflowEntity.WorkflowFailed:
		return "failed"
	case workflowEntity.WorkflowCancel:
		return "cancelled"
	case workflowEntity.WorkflowInterrupted:
		return "interrupted"
	default:
		return "unknown"
	}
}

func mapNodeStatus(status workflowEntity.NodeExecuteStatus) string {
	switch status {
	case workflowEntity.NodeRunning:
		return "running"
	case workflowEntity.NodeSuccess:
		return "success"
	case workflowEntity.NodeFailed:
		return "failed"
	default:
		return "unknown"
	}
}
