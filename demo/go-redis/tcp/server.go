package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	var closed chan struct{}
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			close(closed)
		}
	}()

	logger.Info("start listen ...")
	ListenAndServe(listener, handler, closed)
	return nil
}

// ListenAndServeWithSignal创建监听套接字后，调用ListenAndServe处理客户端连接

func ListenAndServe(
	listener net.Listener, // 无需用指针
	handler tcp.Handler,
	closed <-chan struct{}, // 表示server端关闭
) {
	go func() {
		<-closed
		_ = handler.Close()
		_ = listener.Close() // 任何accept操作会被中止并返回error
	}()

	defer func() {
		_ = handler.Close()
		_ = listener.Close()
	}()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(err)
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(context.Background(), conn)
		}()
	}
	wg.Wait()
	logger.Info("server close")
}
