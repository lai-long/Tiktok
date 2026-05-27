package main

import (
	"Tiktok/biz/middleware"
	"Tiktok/internal/config"
	ws "Tiktok/internal/ws"
	wsService "Tiktok/internal/ws/service"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"Tiktok/pkg/logger"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/zap"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		logger.Fatal("加载.env错误", zap.Error(err))
	}
	cfgPath := os.Getenv("CONFIG_PATH")
	logger.Info("loading configuration", zap.String("config_path", cfgPath))
	cfg, err := config.Load([]string{cfgPath})
	if err != nil {
		logger.Fatal("加载config.yaml错误", zap.Error(err))
	}
	logger.Info("config loaded", logger.WithServiceName(logger.ServiceName))
	if err := logger.InitLogger(cfg.Log.Level, "ws", cfg.Log.Path); err != nil {
		logger.Fatal("初始化日志错误", zap.Error(err))
	}
	logger.Info("logger initialized", logger.WithServiceName("ws"), zap.String("level", cfg.Log.Level))
	sentinelPath := os.Getenv("SENTINEL_PATH")
	logger.Info("loading sentinel", zap.String("sentinel_path", sentinelPath))
	err = config.LoadRules([]string{sentinelPath})
	if err != nil {
		logger.Fatal("加载sentinel rules错误", zap.Error(err))
	}

	rdb := cache.InitRedis()
	re := cache.NewRedis(rdb)
	defer func() {
		err := rdb.Close()
		if err != nil {
			logger.Error("main redis close error", zap.Error(err))
		}
	}()

	ddb := dao.InitDb()
	defer func() {
		err := ddb.Close()
		if err != nil {
			logger.Error("main database close error", zap.Error(err))
		}
	}()
	mysqlDb := dao.NewMySQLdb(ddb)

	wsSvc := wsService.NewWebsocketService(mysqlDb, re)
	wsHandler := ws.NewWebsocketServer(mysqlDb, re, wsSvc)

	h := server.Default(
		server.WithHostPorts(":8881"),
		server.WithMaxRequestBodySize(10*1024*1024),
	)
	h.Use(middleware.AuthMiddleware)

	r := h.Group("/")
	r.GET("/ws", wsHandler.WebSocketHandler)

	go wsSvc.Start()

	h.Spin()
}
