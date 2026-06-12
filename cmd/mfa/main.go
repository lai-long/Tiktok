package main

import (
	"Tiktok/internal/config"
	handler "Tiktok/internal/mfa"
	"Tiktok/internal/mfa/service"
	"Tiktok/internal/middleware"
	"Tiktok/kitex_gen/mfa/mfaservice"
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
	if err := logger.InitLogger(cfg.Log.Level, "mfa", cfg.Log.Path); err != nil {
		logger.Fatal("初始化日志错误", zap.Error(err))
	}
	logger.Info("logger initialized", logger.WithServiceName("mfa"), zap.String("level", cfg.Log.Level))
	sentinelPath := os.Getenv("SENTINEL_PATH")
	logger.Info("loading sentinel", zap.String("sentinel_path", sentinelPath))
	err = config.LoadRules([]string{sentinelPath})
	if err != nil {
		logger.Fatal("加载sentinel rules错误", zap.Error(err))
	}

	db := dao.InitDb()
	mysqlDb := dao.NewMySQLdb(db)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("db close error", zap.Error(err))
		}
	}()

	mfaRepo := service.NewMfaRepo(mysqlDb)
	mfaService := handler.NewMfaService(mfaRepo)

	r, err := etcd.NewEtcdRegistry([]string{config.GetCfg().EtcdAddr})
	if err != nil {
		logger.Error("registry error", zap.Error(err))
	}
	addr := &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8890,
	}

	svr := mfaservice.NewServer(mfaService,
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithMiddleware(middleware.SentinelMiddleware),
		server.WithMiddleware(middleware.TracingMiddleware("mfa")),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "mfaService",
		}),
	)

	err = svr.Run()

	if err != nil {
		logger.Error("server error", zap.Error(err))
	}
}
