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
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"os"
	"os/exec"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// GetUnstagedChanges 获取当前Git仓库中所有没有提交进暂存区的改动
// 返回git diff的完整输出，包含具体的改动内容
func GetUnstagedChanges() (string, error) {
	// 执行 git diff 命令获取工作目录与暂存区的差异
	cmd := exec.Command("git", "diff")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行git diff失败: %v", err)
	}

	// 如果没有改动，检查是否有未跟踪的文件
	if len(output) == 0 {
		// 检查未跟踪的文件
		statusCmd := exec.Command("git", "status", "--porcelain")
		statusOutput, err := statusCmd.Output()
		if err != nil {
			return "", fmt.Errorf("执行git status失败: %v", err)
		}
		
		if len(statusOutput) == 0 {
			return "", nil
		}
		
		// 解析未跟踪的文件
		var untrackedFiles []string
		scanner := bufio.NewScanner(strings.NewReader(string(statusOutput)))
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) >= 2 && line[0:2] == "??" {
				filePath := strings.TrimSpace(line[2:])
				untrackedFiles = append(untrackedFiles, filePath)
			}
		}
		
		if len(untrackedFiles) > 0 {
			return fmt.Sprintf("未跟踪的文件:\n%s", strings.Join(untrackedFiles, "\n")), nil
		}
		
		return "", nil
	}

	return string(output), nil
}

func main() {
	changes, err := GetUnstagedChanges()
	if err != nil {
		fmt.Printf("获取未暂存变更失败: %v\n", err)
		return
	}

	if changes == "" {
		fmt.Printf("当前没有未暂存的变更，无需进行代码审查\n")
		return
	}
	
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	temperature := float32(0.7)

	ctx := context.Background()

	// 创建并配置 ChatModel
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     openAIBaseURL,
		Model:       openAIModelName,
		APIKey:      openAIAPIKey,
		Temperature: &temperature,
	})
	if err != nil {
		return
	}

	// 初始化 tools 配置
	toolsConfig := compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
		},
	}

	// 创建 ReAct Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      toolsConfig,
		MaxStep:          20, // 设置最大推理步数，允许10轮对话（10次ChatModel + 10次Tools）
		// 如果开启 toolcall，可能需要重新实现启用这个注释
		// StreamToolCallChecker: func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
		// 	defer sr.Close()
		// 	for {
		// 		msg, err := sr.Recv()
		// 		if err != nil {
		// 			if errors.Is(err, io.EOF) {
		// 				// finish
		// 				break
		// 			}
		// 			return false, err
		// 		}
		// 		if len(msg.ToolCalls) > 0 {
		// 			return true, nil
		// 		}
		// 	}
		// 	return false, nil
		// },
	})

	if err != nil {
		return
	}

	// 不再使用glamour渲染器，直接输出原始内容

	fmt.Println("开始代码审查...\n")
	
	// 使用Stream方法进行流式输出
	msgReader, err := agent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: CodeReviewSystemPrompt,
		},
		{
			Role:    schema.User,
			Content: "代码变更：\n```diff\n" + changes + "\n```",
		},
	})

	if err != nil {
		return
	}
	
	for {
		// msg type is *schema.Message
		msg, err := msgReader.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			// error
			fmt.Printf("failed to recv: %v\n", err)
			return
		}

		fmt.Print(msg.Content)
	}
}
