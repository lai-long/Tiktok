package main

import (
	"Tiktok/biz/middleware"
	ws "Tiktok/internal/ws"
	wsService "Tiktok/internal/ws/service"
	"Tiktok/pkg/config"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("/home/lai-long/Tiktok/.env"); err != nil {
		log.Println("load env error:", err)
	}
	cfgPath := os.Getenv("CONFIG_PATH")
	log.Printf("Loading configuration from %s", cfgPath)
	cfg, err := config.Load([]string{cfgPath})
	if err != nil {
		log.Fatal("加载config.yaml错误", err)
	}
	log.Println(cfg)
	sentinelPath := os.Getenv("SENTINEL_PATH")
	log.Printf("Loading sentinel from %s", sentinelPath)
	err = config.LoadRules([]string{sentinelPath})
	if err != nil {
		log.Fatal("加载sentinel rules错误", err)
	}

	rdb := cache.InitRedis()
	re := cache.NewRedis(rdb)
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Println("main redis close err", err)
		}
	}()

	ddb := dao.InitDb()
	defer func() {
		err := ddb.Close()
		if err != nil {
			log.Println("main database close err", err)
		}
	}()
	mysqlDb := dao.NewMySQLdb(ddb)

	wsSvc := wsService.NewWebsocketService(mysqlDb, re)
	wsHandler := ws.NewWebsocketSever(mysqlDb, re, wsSvc)

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
