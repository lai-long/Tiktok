package service

import (
	"Tiktok/internal/config"
	"Tiktok/pkg/logger"
	"context"
	"time"

	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"
	"go.uber.org/zap"
)

type MyAccount struct{}

func (a *MyAccount) GetConfiguredProviders() ([]schemas.ModelProvider, error) {
	return []schemas.ModelProvider{schemas.OpenAI}, nil
}
func (a *MyAccount) GetKeysForProvider(ctx context.Context, provider schemas.ModelProvider) ([]schemas.Key, error) {
	return []schemas.Key{
		{
			Value:  *schemas.NewEnvVar(config.GetCfg().API.APIKey),
			Models: schemas.WhiteList{config.GetCfg().API.Model},
			Weight: 1.0,
		},
	}, nil
}
func (a *MyAccount) GetConfigForProvider(_ schemas.ModelProvider) (*schemas.ProviderConfig, error) {
	return &schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{
			BaseURL: config.GetCfg().API.BaseURL,
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
	clientCfg := make([]*schemas.MCPClientConfig, len(config.GetCfg().Mcp.Clients))
	for i := range config.GetCfg().Mcp.Clients {
		connType := ConnectType[config.GetCfg().Mcp.Clients[i].ConnectionType]
		cfg := &schemas.MCPClientConfig{
			ID:                 config.GetCfg().Mcp.Clients[i].ID,
			Name:               config.GetCfg().Mcp.Clients[i].Name,
			ConnectionType:     connType,
			ToolsToExecute:     config.GetCfg().Mcp.Clients[i].ToolsToExecute,
			ToolsToAutoExecute: config.GetCfg().Mcp.Clients[i].ToolsToAutoExecute,
		}
		switch connType {
		case schemas.MCPConnectionTypeSTDIO:
			cfg.StdioConfig = &schemas.MCPStdioConfig{
				Command: config.GetCfg().Mcp.Clients[i].Command,
				Args:    config.GetCfg().Mcp.Clients[i].Args,
			}
		case schemas.MCPConnectionTypeHTTP, schemas.MCPConnectionTypeSSE:
			if config.GetCfg().Mcp.Clients[i].URL != "" {
				cfg.ConnectionString = schemas.NewEnvVar(config.GetCfg().Mcp.Clients[i].URL)
			}
		case schemas.MCPConnectionTypeInProcess:
		}
		clientCfg[i] = cfg
	}
	toolManagerCfg := &schemas.MCPToolManagerConfig{
		ToolExecutionTimeout: time.Duration(config.GetCfg().Mcp.ToolManagerConfig.MaxTime) * time.Second,
		MaxAgentDepth:        int(config.GetCfg().Mcp.ToolManagerConfig.MaxDepth),
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
		Model:    config.GetCfg().API.Model,
		Input:    c.message,
	}
	resp, bifrostErr := c.client.ChatCompletionRequest(bifrostContext, req)
	if bifrostErr != nil {
		return "", bifrostErr.Error.Error
	}
	if len(resp.Choices) == 0 {
		logger.Warn("ai response choices is nil or empty")
		return "", nil
	}
	if resp.Choices[0].Message.Content != nil && resp.Choices[0].Message.Content.ContentStr != nil {
		content = *resp.Choices[0].Message.Content.ContentStr
		logger.Debug("ai response content", zap.String("content", content))
	} else {
		logger.Warn("ai response content is nil")
	}
	return content, nil
}

func (c *ChatClient) Close() {
	c.client.Shutdown()
}

type Agent struct {
	chatCli *ChatClient
	ctx     context.Context
}

func NewAgent(ctx context.Context) *Agent {
	llm, err := NewChatClient(ctx)
	if err != nil {
		logger.Error("mcp NewChatClient failed", zap.Error(err))
		return nil
	}
	if llm == nil {
		logger.Error("mcp NewChatClient returned nil")
		return nil
	}
	return &Agent{
		chatCli: llm,
		ctx:     ctx,
	}
}

func (a *Agent) StartAction(prompt string) string {
	if a.chatCli == nil {
		logger.Error("startAction failed: chat client is nil")
		return ""
	}
	resp, err := a.chatCli.Chat(prompt)
	if err != nil {
		logger.Error("startAction failed", zap.Error(err))
		return ""
	}
	return resp
}

func (a *Agent) StopAction() {
	if a.chatCli == nil {
		logger.Error("stopAction failed: chat client is nil")
		return
	}
	a.chatCli.Close()
}
