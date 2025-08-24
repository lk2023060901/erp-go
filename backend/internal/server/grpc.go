package server

import (
	"erp-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// GRPCServer GRPC服务器
type GRPCServer struct {
	Server *grpc.Server
}

// NewGRPCServer 创建GRPC服务器
func NewGRPCServer(c *conf.Server, logger log.Logger) *GRPCServer {
	var opts = []grpc.ServerOption{}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != "" {
		timeout := c.Grpc.GetTimeout()
		opts = append(opts, grpc.Timeout(timeout))
	}

	srv := grpc.NewServer(opts...)

	return &GRPCServer{
		Server: srv,
	}
}
