package main

import (
	"Tiktok/internal/config"
	"Tiktok/internal/middleware"
	handler "Tiktok/internal/user"
	"Tiktok/internal/user/service"
	userservice "Tiktok/kitex_gen/user/userservice"
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
	if err := logger.InitLogger(cfg.Log.Level, "user", cfg.Log.Path); err != nil {
		logger.Fatal("初始化日志错误", zap.Error(err))
	}
	logger.Info("logger initialized", logger.WithServiceName("user"), zap.String("level", cfg.Log.Level))
	sentinelPath := os.Getenv("SENTINEL_PATH")
	logger.Info("loading sentinel", zap.String("sentinel_path", sentinelPath))
	err = config.LoadRules([]string{sentinelPath})
	if err != nil {
		logger.Fatal("加载sentinel rules错误", zap.Error(err))
	}
	r, err := etcd.NewEtcdRegistry([]string{config.GetCfg().EtcdAddr})
	if err != nil {
		logger.Error("etcd registry error", zap.Error(err))
	}
	db := dao.InitDb()
	mysqlDb := dao.NewMySQLdb(db)

	rdb := cache.InitRedis()
	redis := cache.NewRedis(rdb)

	userRepo := service.NewUserRepo(mysqlDb, redis)
	userService := handler.NewUserService(userRepo)

	addr := &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8889}
	svr := userservice.NewServer(userService,
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithMiddleware(middleware.SentinelMiddleware),
		server.WithMiddleware(middleware.TracingMiddleware("user")),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "userService",
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
	logger.Info("kitex server started", zap.String("addr", addr.String()))
	if err := svr.Run(); err != nil {
		logger.Error("Kitex server error", zap.Error(err))
	}
}
