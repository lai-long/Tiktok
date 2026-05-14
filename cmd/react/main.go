package main

import (
	handler "Tiktok/internal/react"
	"Tiktok/internal/react/service"
	"Tiktok/kitex_gen/react/commentservice"
	"Tiktok/kitex_gen/react/likeservice"
	"Tiktok/pkg/config"
	"Tiktok/pkg/dal/cache"
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
		log.Fatal("etcd registry init error:", err)
	}
	db := dao.InitDb()
	mysqlDb := dao.NewMySQLdb(db)

	rdb := cache.InitRedis()
	redis := cache.NewRedis(rdb)

	commentRepo := service.NewCommentService(mysqlDb)
	commentServiceImpl := handler.NewCommentService(commentRepo)

	likeRepo := service.NewLikeRepo(mysqlDb, mysqlDb, mysqlDb, redis)
	likeServiceImpl := handler.NewLikeService(likeRepo)

	commentAddr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8892}
	likeAddr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8893}

	commentSvc := commentservice.NewServer(commentServiceImpl,
		server.WithServiceAddr(commentAddr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "commentService",
		}),
	)
	likeSvc := likeservice.NewServer(likeServiceImpl,
		server.WithServiceAddr(likeAddr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "likeService",
		}),
	)
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db close err", err)
		}
		if err := rdb.Close(); err != nil {
			log.Println("redis close err", err)
		}
	}()
	go func() {
		log.Printf("Comment Kitex server started at %s", commentAddr.String())
		if err := commentSvc.Run(); err != nil {
			log.Println("Comment Kitex server error:", err)
		}
	}()
	log.Printf("Like Kitex server started at %s", likeAddr.String())
	if err := likeSvc.Run(); err != nil {
		log.Println("Like Kitex server error:", err)
	}
}
