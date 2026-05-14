package main

import (
	handler "Tiktok/internal/social"
	"Tiktok/internal/social/service"
	socialservice "Tiktok/kitex_gen/social/socialservice"
	"Tiktok/pkg/config"
	"Tiktok/pkg/dal/dao"
	"log"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	etcd "github.com/kitex-contrib/registry-etcd"
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
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Println(err)
	}
	db := dao.InitDb()
	mysqlDb := dao.NewMySQLdb(db)

	socialRepo := service.NewSocialRepo(mysqlDb)
	socialService := handler.NewSocialServiceImpl(socialRepo)

	addr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8894}
	svr := socialservice.NewServer(socialService,
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "socialService",
		}),
	)
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db close err", err)
		}
	}()
	log.Printf("Social Kitex server started at %s", addr.String())
	if err := svr.Run(); err != nil {
		log.Println("Kitex server error:", err)
	}
}
