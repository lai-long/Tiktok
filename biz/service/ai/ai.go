package ai

import (
	"Tiktok/pkg/config"
	"context"
	"log"
	"time"

	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"
)

type MyAccount struct{}

func (a *MyAccount) GetConfiguredProviders() ([]schemas.ModelProvider, error) {
	return []schemas.ModelProvider{schemas.OpenAI}, nil
}
func (a *MyAccount) GetKeysForProvider(ctx context.Context, provider schemas.ModelProvider) ([]schemas.Key, error) {
	return []schemas.Key{
		{
			Value:  *schemas.NewEnvVar(config.Cfg.Api.ApiKey),
			Models: schemas.WhiteList{"MiniMax-M2.7"},
			Weight: 1.0,
		},
	}, nil
}
func (a *MyAccount) GetConfigForProvider(_ schemas.ModelProvider) (*schemas.ProviderConfig, error) {
	return &schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{
			BaseURL: config.Cfg.Api.BaseUrl,
		},
		ConcurrencyAndBufferSize: schemas.DefaultConcurrencyAndBufferSize,
	}, nil
}

type ChatClient struct {
	ctx     context.Context
	message []schemas.ChatMessage
	client  *bifrost.Bifrost
}

func NewChatClient(ctx context.Context) (*ChatClient, error) {
	client, err := bifrost.Init(ctx, schemas.BifrostConfig{
		Account: &MyAccount{},
		MCPConfig: &schemas.MCPConfig{
			ClientConfigs: []*schemas.MCPClientConfig{
				{
					ID:             "mcp_client",
					Name:           "mcp_client",
					ConnectionType: schemas.MCPConnectionTypeSTDIO,
					StdioConfig: &schemas.MCPStdioConfig{
						Command: "/home/lai-long/Tiktok/mcp_service/mcp_service",
						Args:    []string{},
						Envs:    []string{},
					},
					ToolsToExecute:     schemas.WhiteList{"*"},
					ToolsToAutoExecute: schemas.WhiteList{"*"},
				},
			},
			ToolManagerConfig: &schemas.MCPToolManagerConfig{
				ToolExecutionTimeout: 200 * time.Second,
				MaxAgentDepth:        10,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	message := []schemas.ChatMessage{
		{
			Role: schemas.ChatMessageRoleSystem,
			Content: &schemas.ChatMessageContent{
				ContentStr: schemas.Ptr("如果你成功通过工具调用获取到了结果，那么只需要返回结果，不要再返回toolCall了"),
			},
		},
	}
	return &ChatClient{ctx: ctx, message: message, client: client}, nil
}

func (c *ChatClient) Chat(prompt string) (content string, err error) {
	c.message = append(c.message, schemas.ChatMessage{
		Role: schemas.ChatMessageRoleUser,
		Content: &schemas.ChatMessageContent{
			ContentStr: schemas.Ptr(prompt),
		},
	})
	bifrostContext := schemas.NewBifrostContext(c.ctx, schemas.NoDeadline)
	req := &schemas.BifrostChatRequest{
		Provider: schemas.OpenAI,
		Model:    "MiniMax-M2.7",
		Input:    c.message,
	}
	resp, bifrostErr := c.client.ChatCompletionRequest(bifrostContext, req)
	if bifrostErr != nil {
		log.Printf("ChatCompletionRequest err: %+v, Error: %+v\n", bifrostErr, bifrostErr.Error)
		return "", bifrostErr.Error.Error
	}
	if len(resp.Choices) == 0 {
		log.Println("choices is nil or empty")
		return "", nil
	}
	if resp.Choices[0].Message.Content != nil && resp.Choices[0].Message.Content.ContentStr != nil {
		content = *resp.Choices[0].Message.Content.ContentStr
		log.Println("ai resp content:", content)
	} else {
		log.Println("short of ai resp content:")
	}
	return content, nil
}

func (c *ChatClient) Close() {
	c.client.Shutdown()
}
