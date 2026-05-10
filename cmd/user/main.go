package main

import (
	handler "Tiktok/internal/user"
	"Tiktok/internal/user/service"
	userservice "Tiktok/kitex_gen/user/userservice"
	"Tiktok/pkg/config"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"log"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/joho/godotenv"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/cloudwego/kitex/server"
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

	rdb := cache.InitRedis()
	redis := cache.NewRedis(rdb)

	userRepo := service.NewUserRepo(mysqlDb, mysqlDb, redis)
	userService := handler.NewUserService(userRepo)

	addr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8889}
	svr := userservice.NewServer(userService,
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "userService",
		}),
	)

	log.Printf("User Kitex server started at %s", addr.String())
	if err := db.Close(); err != nil {
		log.Println("db close err", err)
	}
	if err := rdb.Close(); err != nil {
		log.Println("redis close err", err)
	}
	if err := svr.Run(); err != nil {
		log.Fatal("Kitex server error:", err)
	}
}
