package ai

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

type Agent struct {
	McpClient []*McpClient
	LLM       *ChatOpenAI
	Model     string
	Ctx       context.Context
	RAGCtx    string
}

func NewAgent(ctx context.Context, model string, mcpClients []*McpClient) *Agent {
	tools := make([]mcp.Tool, 0)
	for _, cli := range mcpClients {
		err := cli.Start()
		if err != nil {
			log.Println(err)
		}
		err = cli.SetTools()
		if err != nil {
			log.Println(err)
		}
		tools = append(tools, cli.GetTools()...)
		log.Println("工具:", cli.GetTools(), "1,cmd", cli.Cmd)
	}
	llm := NewChatOpenAI(ctx, model, WithTools(tools))
	log.Println("init agent success")
	return &Agent{
		McpClient: mcpClients,
		LLM:       llm,
		Model:     model,
		Ctx:       ctx,
	}
}
func (a *Agent) StartAction(prompt string) string {
	if a.LLM == nil {
		log.Println("llm is nil")
		return ""
	}
	responses, callTools := a.LLM.Chat(prompt)
	log.Printf("DEBUG: responses=%s callTools count=%d", responses, len(callTools))
	if len(callTools) > 0 {
		for _, toolCall := range callTools {
			for _, mcpCli := range a.McpClient {
				tools := mcpCli.GetTools()
				for _, tool := range tools {
					if tool.Name == toolCall.Function.Name {
						log.Println("tool call found using tools:", tool.Name)
						toolResult, err := mcpCli.CallTool(tool.Name, toolCall.Function.Arguments)
						if err != nil {
							log.Println("call tools failed", err)
						}
						a.LLM.Message = append(a.LLM.Message, openai.ToolMessage(toolResult, toolCall.ID))
					}
				}
			}
		}
		responses, _ = a.LLM.Chat("")
	}
	a.StopAction()
	return responses
}
func (a *Agent) StopAction() {
	for _, cli := range a.McpClient {
		_ = cli.Close()
	}
}
