package tcp

import (
	"context"
	"net"
)

// tcp模块只管tcp连接，业务相关的事情交给handler

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
