package ai

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpClient struct {
	Ctx    context.Context
	Client *client.Client
	Tools  []mcp.Tool
	Cmd    string
	Args   []string
	Env    []string
}

func NewMcpClient(Ctx context.Context, cmd string, args []string, env []string) *McpClient {
	stdio := transport.NewStdio(cmd, env, args...)
	clients := client.NewClient(stdio)
	return &McpClient{
		Ctx:    Ctx,
		Client: clients,
		Cmd:    cmd,
		Args:   args,
		Env:    env,
	}
}
func (c *McpClient) Start() error {
	err := c.Client.Start(c.Ctx)
	if err != nil {
		return err
	}
	mcpInitReq := mcp.InitializeRequest{
		Request: mcp.Request{},
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "Tiktok chat",
				Version: "0.0.1",
			},
		},
	}
	_, err = c.Client.Initialize(c.Ctx, mcpInitReq)
	if err != nil {
		return err
	}
	log.Println("successfully initialized")
	return nil
}
func (c *McpClient) SetTools() error {
	toolsReq := mcp.ListToolsRequest{
		PaginatedRequest: mcp.PaginatedRequest{},
		Header:           nil,
	}
	tools, err := c.Client.ListTools(c.Ctx, toolsReq)
	if err != nil {
		return err
	}
	c.Tools = tools.Tools
	log.Println("successfully listed tools")
	return nil
}
func (c *McpClient) Close() error {
	return c.Client.Close()
}

func (c *McpClient) CallTool(toolName string, args interface{}) (string, error) {
	var argument map[string]any
	switch v := args.(type) {
	case map[string]any:
		argument = v
	case string:
		json.Unmarshal([]byte(v), &argument)
	default:
		argument = map[string]any{}
	}
	log.Println("正在call tool:", toolName, argument)
	result, err := c.Client.CallTool(c.Ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: argument,
		},
	})
	log.Println("call tool result:", result)
	if err != nil {
		return "", err
	}
	return mcp.GetTextFromContent(result.Result), nil
}

func (c *McpClient) GetTools() []mcp.Tool {
	return c.Tools
}
