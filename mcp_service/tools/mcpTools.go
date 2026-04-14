package tools

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func AddTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Println("成功调用add tool")
	args, _ := req.Params.Arguments.(map[string]any)
	num1, _ := args["number1"].(float64)
	log.Println(args["number1"].(float64))
	num2, _ := args["number2"].(float64)
	log.Println(args["number2"].(float64))
	result := num1 + num2
	log.Println("result:", result)
	return mcp.NewToolResultJSON(map[string]any{
		"result":  result,
		"message": "成功调用add工具结果位于result",
	})
}
