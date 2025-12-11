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
	"fmt"
	"strings"

	scriptmodel "github.com/coze-dev/coze-studio/backend/api/model/script"
)

type TextEditService struct{}

func (s *TextEditService) Edit(ctx context.Context, req *scriptmodel.TextEditRequest) (*scriptmodel.TextEditResponse, error) {
	trimmed := strings.TrimSpace(req.OriginalText)
	if trimmed == "" {
		return nil, fmt.Errorf("original_text is required")
	}

	edited := trimmed
	if req.EditIntent != "" {
		edited = fmt.Sprintf("%s\n\n[AI 编辑意图：%s]\n%s", trimmed, req.EditIntent, trimmed)
	}

	diff := buildDiff(trimmed, edited)

	return &scriptmodel.TextEditResponse{
		EditedText: edited,
		Diff:       diff,
	}, nil
}

func buildDiff(original, edited string) string {
	if original == edited {
		return ""
	}
	return fmt.Sprintf("--- 原始\n+++ 编辑后\n@@\n- %s\n+ %s", original, edited)
}
