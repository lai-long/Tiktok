package main

import (
	"Tiktok/internal/config"
	"Tiktok/internal/middleware"
	handler "Tiktok/internal/react"
	"Tiktok/internal/react/service"
	"Tiktok/kitex_gen/react/commentservice"
	"Tiktok/kitex_gen/react/likeservice"
	"Tiktok/pkg/dal/cache"
	"Tiktok/pkg/dal/dao"
	"Tiktok/pkg/logger"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
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
	if err := logger.InitLogger(cfg.Log.Level, "react", cfg.Log.Path); err != nil {
		logger.Fatal("初始化日志错误", zap.Error(err))
	}
	logger.Info("logger initialized", logger.WithServiceName("react"), zap.String("level", cfg.Log.Level))
	sentinelPath := os.Getenv("SENTINEL_PATH")
	logger.Info("loading sentinel", zap.String("sentinel_path", sentinelPath))
	err = config.LoadRules([]string{sentinelPath})
	if err != nil {
		logger.Fatal("加载sentinel rules错误", zap.Error(err))
	}
	r, err := etcd.NewEtcdRegistry([]string{config.Cfg.EtcdAddr})
	if err != nil {
		logger.Fatal("etcd registry init error", zap.Error(err))
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
		server.WithMiddleware(middleware.SentinelMiddleware),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "commentService",
		}),
	)
	likeSvc := likeservice.NewServer(likeServiceImpl,
		server.WithServiceAddr(likeAddr),
		server.WithRegistry(r),
		server.WithMiddleware(middleware.SentinelMiddleware),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "likeService",
		}),
	)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("db close error", zap.Error(err))
		}
		if err := rdb.Close(); err != nil {
			logger.Error("redis close error", zap.Error(err))
		}
	}()
	go func() {
		logger.Info("comment server started", zap.String("addr", commentAddr.String()))
		if err := commentSvc.Run(); err != nil {
			logger.Error("Comment Kitex server error", zap.Error(err))
		}
	}()
	logger.Info("like server started", zap.String("addr", likeAddr.String()))
	if err := likeSvc.Run(); err != nil {
		logger.Error("Like Kitex server error", zap.Error(err))
	}
}
