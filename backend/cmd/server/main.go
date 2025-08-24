package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"erp-system/internal/conf"
	"erp-system/internal/server"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "erp-system"
	// Version is the version of the compiled software.
	Version string = "v1.0.0"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Errorf("Failed to load config: %v", err)
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Errorf("Failed to scan config: %v", err)
		panic(err)
	}

	// 调试配置内容
	logger.Log(log.LevelInfo, "msg", "Config path: "+flagconf)
	logger.Log(log.LevelInfo, "msg", fmt.Sprintf("Loaded config - Server: %+v", bc.Server))
	logger.Log(log.LevelInfo, "msg", fmt.Sprintf("Loaded config - Data: %+v", bc.Data))

	// 添加配置验证
	if bc.Server == nil {
		logger.Log(log.LevelError, "msg", "Server configuration is missing")
		panic("Server configuration is missing")
	}
	if bc.Server.Http == nil {
		panic("HTTP server configuration is missing")
	}
	if bc.Server.Grpc == nil {
		panic("gRPC server configuration is missing")
	}

	// 初始化应用
	app, cleanup, err := server.InitializeApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
