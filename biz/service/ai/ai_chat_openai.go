package ai

import "C"
import (
	"Tiktok/pkg/config"
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type ChatOpenAI struct {
	Ctx          context.Context
	ModelName    string
	SystemPrompt string
	RagContext   string
	Tools        []mcp.Tool
	LLM          openai.Client
	Message      []openai.ChatCompletionMessageParamUnion
}
type LLMOption func(opt *ChatOpenAI)

func WithSystemPrompt(prompt string) LLMOption {
	return func(c *ChatOpenAI) {
		c.SystemPrompt = prompt
	}
}
func WithRagContext(rag string) LLMOption {
	return func(c *ChatOpenAI) {
		c.RagContext = rag
	}
}
func WithTools(tools []mcp.Tool) LLMOption {
	return func(c *ChatOpenAI) {
		c.Tools = tools
	}
}
func NewChatOpenAI(ctx context.Context, Model string, opts ...LLMOption) *ChatOpenAI {
	if Model == "" {
		log.Println("[ERROR] model name is required")
	}
	apiKep := config.Cfg.Api.ApiKey
	baseUrl := config.Cfg.Api.BaseUrl
	options := []option.RequestOption{
		option.WithAPIKey(apiKep),
		option.WithBaseURL(baseUrl),
	}
	client := openai.NewClient(options...)
	llm := &ChatOpenAI{
		Ctx:       ctx,
		ModelName: Model,
		LLM:       client,
	}
	for _, opt := range opts {
		opt(llm)
	}
	if llm.SystemPrompt != "" {
		llm.Message = append(llm.Message, openai.SystemMessage(llm.SystemPrompt))
	}
	if llm.RagContext != "" {
		llm.Message = append(llm.Message, openai.UserMessage(llm.RagContext))
	}
	log.Println("[INFO] create new chat " + Model)
	return llm
}
func (c *ChatOpenAI) Chat(prompt string) (result string, toolCalls []openai.ToolCallUnion) {
	if prompt != "" {
		c.Message = append(c.Message, openai.UserMessage(prompt))
	}
	stream := c.LLM.Chat.Completions.NewStreaming(c.Ctx, openai.ChatCompletionNewParams{
		Messages: c.Message,
		Model:    c.ModelName,
		Seed:     openai.Int(0),
		Tools:    ConvertMcpToOpenAITool(c.Tools),
	})
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		acc.AddChunk(stream.Current())
		if content, ok := acc.JustFinishedContent(); ok {
			result = content
		}
		log.Println("what tools====", acc.Choices[0].Message.ToolCalls)
		if acc.Choices[0].Message.ToolCalls != nil {
			toolTemp := openai.ToolCallUnion{
				ID: acc.Choices[0].Message.ToolCalls[0].ID,
				Function: openai.FunctionToolCallFunction{
					Arguments: acc.Choices[0].Message.ToolCalls[0].Function.Arguments,
					Name:      acc.Choices[0].Message.ToolCalls[0].Function.Name,
				},
			}
			toolCalls = append(toolCalls, toolTemp)
		}
	}
	if stream.Err() != nil {
		log.Println("error:", stream.Err())
	}
	log.Println("历史消息：", c.Message)
	log.Println(result + "---result")
	return result, toolCalls
}
func ConvertMcpToOpenAITool(mcpTools []mcp.Tool) []openai.ChatCompletionToolParam {
	openAITools := make([]openai.ChatCompletionToolParam, len(mcpTools))
	for i, tool := range mcpTools {
		log.Println("tool.name:", tool.Name, "tool.InputSchema", tool.InputSchema)
		openAITools[i] = openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters: shared.FunctionParameters{
					"type":       tool.InputSchema.Type,
					"properties": tool.InputSchema.Properties,
					"required":   tool.InputSchema.Required,
				},
			},
		}
	}
	return openAITools
}
