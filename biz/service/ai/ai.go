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
			Value:  *schemas.NewEnvVar(config.Cfg.API.APIKey),
			Models: schemas.WhiteList{"MiniMax-M2.7"},
			Weight: 1.0,
		},
	}, nil
}
func (a *MyAccount) GetConfigForProvider(_ schemas.ModelProvider) (*schemas.ProviderConfig, error) {
	return &schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{
			BaseURL: config.Cfg.API.BaseURL,
		},
		ConcurrencyAndBufferSize: schemas.DefaultConcurrencyAndBufferSize,
	}, nil
}

type ChatClient struct {
	ctx     context.Context
	message []schemas.ChatMessage
	client  *bifrost.Bifrost
}

var ConnectType = map[string]schemas.MCPConnectionType{
	"stdio":     schemas.MCPConnectionTypeSTDIO,
	"http":      schemas.MCPConnectionTypeHTTP,
	"sse":       schemas.MCPConnectionTypeSSE,
	"inprocess": schemas.MCPConnectionTypeInProcess,
}

func NewChatClient(ctx context.Context) (*ChatClient, error) {
	clientCfg := make([]*schemas.MCPClientConfig, len(config.Cfg.Mcp.Clients))
	for i := range config.Cfg.Mcp.Clients {
		clientCfg[i] = &schemas.MCPClientConfig{
			ID:             config.Cfg.Mcp.Clients[i].ID,
			Name:           config.Cfg.Mcp.Clients[i].Name,
			ConnectionType: ConnectType[config.Cfg.Mcp.Clients[i].ConnectionType],
			StdioConfig: &schemas.MCPStdioConfig{
				Command: config.Cfg.Mcp.Clients[i].Command,
				Args:    config.Cfg.Mcp.Clients[i].Args,
			},
			ToolsToExecute:     config.Cfg.Mcp.Clients[i].ToolsToExecute,
			ToolsToAutoExecute: config.Cfg.Mcp.Clients[i].ToolsToAutoExecute,
		}
	}
	toolManagerCfg := &schemas.MCPToolManagerConfig{
		ToolExecutionTimeout: time.Duration(config.Cfg.Mcp.ToolManagerConfig.MaxTime) * time.Second,
		MaxAgentDepth:        int(config.Cfg.Mcp.ToolManagerConfig.MaxDepth),
	}
	client, err := bifrost.Init(ctx, schemas.BifrostConfig{
		Account: &MyAccount{},
		MCPConfig: &schemas.MCPConfig{
			ClientConfigs:     clientCfg,
			ToolManagerConfig: toolManagerCfg,
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
