/*
 * Copyright 2025 Sg
 * Copyright 2024 CloudWeGo Authors
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// 辅助函数，类似于 gptr.Of
func ptrOf[T any](v T) *T {
	return &v
}

func main() {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	temperature := float32(0.7)

	ctx := context.Background()

	// 创建 DuckDuckGo 工具
	searchTool, err := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{})
	if err != nil {
		log.Printf("NewTextSearchTool failed, err=%v", err)
		return
	}

	// 创建并配置 ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     openAIBaseURL,
		Model:       openAIModelName,
		APIKey:      openAIAPIKey,
		Temperature: &temperature,
	})
	if err != nil {
		log.Printf("NewChatModel failed, err=%v", err)
		return
	}

	// 初始化 tools 配置
	toolsConfig := compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			searchTool,
		},
	}

	// 创建 ReAct Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      toolsConfig,
		MaxStep:          20, // 设置最大推理步数，允许10轮对话（10次ChatModel + 10次Tools）
	})

	if err != nil {
		log.Printf("react.NewAgent failed, err=%v", err)
		return
	}

	log.Println("=== 代码审查开始 ===")
	// HACK: 使用Generate方法获取完整响应，因为Stream会因为模型供应商对于ToolCall的支持而提前终止
	resp, err := agent.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "请搜索cloudwego/eino的仓库地址，然后告诉我仓库的地址",
		},
	})
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}
	fmt.Println(resp.Content)
}
