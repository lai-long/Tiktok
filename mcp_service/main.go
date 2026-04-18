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
	addTool := mcp.NewTool("add", mcp.WithDescription("把两个数加起来返回和"),
		mcp.WithNumber("number1", mcp.Required(), mcp.Description("第一个数字")),
		mcp.WithNumber("number2", mcp.Required(), mcp.Description("第二个数字")))
	s.AddTool(addTool, tools.AddTool)

	geocodeTool := mcp.NewTool("geocode", mcp.WithDescription("将结构化地址转换为高德经纬度坐标"),
		mcp.WithString("address", mcp.Required(), mcp.Description("结构化地址，如：北京市朝阳区阜通东大街6号")))
	s.AddTool(geocodeTool, tools.GeocodeTool)

	regeoTool := mcp.NewTool("regeo", mcp.WithDescription("将经纬度坐标转换为详细结构化地址"),
		mcp.WithString("location", mcp.Required(), mcp.Description("经纬度坐标，格式：经度,纬度，如：116.480881,39.989410")),
		mcp.WithString("extensions", mcp.Description("返回结果控制，可选 base/all，默认 base。设置为 all 时返回 POI 等更多信息")))
	s.AddTool(regeoTool, tools.RegeoTool)

	weatherTool := mcp.NewTool("weather", mcp.WithDescription("查询天气信息，支持城市名或 adcode，自动识别"),
		mcp.WithString("city", mcp.Required(), mcp.Description("城市名（如北京）或 adcode（如110000）")),
		mcp.WithString("extensions", mcp.Description("气象类型，可选 base/all，默认 base。base 返回实况天气，all 返回预报天气")))
	s.AddTool(weatherTool, tools.WeatherTool)
}
