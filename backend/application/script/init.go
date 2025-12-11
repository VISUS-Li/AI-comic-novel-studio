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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/script/repository"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
)

var ScriptSVC *ScriptApplicationService

type ScriptApplicationService struct {
	ScriptProject    *ScriptProjectService
	WorkflowInstance *WorkflowInstanceService
	Asset            *AssetService
	WorkflowRun      *WorkflowRunService
	Text             *TextEditService
}

func InitService(db *gorm.DB, idGen idgen.IDGenerator) *ScriptApplicationService {
	projectRepo := repository.NewScriptProjectRepo(db, idGen)
	workflowRepo := repository.NewWorkflowInstanceRepo(db, idGen)
	assetRepo := repository.NewAssetRepo(db, idGen)
	runRepo := repository.NewWorkflowRunRepo(db, idGen)
	nodeOutputRepo := repository.NewNodeOutputRepo(db, idGen)

	assetSvc := &AssetService{repo: assetRepo}
	svc := &ScriptApplicationService{
		ScriptProject:    &ScriptProjectService{repo: projectRepo},
		WorkflowInstance: &WorkflowInstanceService{repo: workflowRepo},
		Asset:            assetSvc,
		WorkflowRun:      &WorkflowRunService{runRepo: runRepo, nodeOutputRepo: nodeOutputRepo, workflowInstanceRepo: workflowRepo, assetSvc: assetSvc},
		Text:             &TextEditService{},
	}
	ScriptSVC = svc
	return svc
}
