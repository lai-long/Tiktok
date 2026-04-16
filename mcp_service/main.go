package main

import (
	"Tiktok/mcp_service/tools"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"calculator-server",
		"1.0.0",
	)
	registerTools(s)
	if err := server.ServeStdio(s); err != nil {
		log.Printf("mcp Server error: %v\n", err)
	}
}
func registerTools(s *server.MCPServer) {
	addTool := mcp.NewTool("add", mcp.WithDescription("把两个数加起来返回它们的和"), mcp.WithNumber("number1", mcp.Required(), mcp.Description("第一个数字")),
		mcp.WithNumber("number2", mcp.Required(), mcp.Description("第二个数字")))
	s.AddTool(addTool, tools.AddTool)
}
