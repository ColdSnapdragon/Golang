package tcp

// 测试tcp服务是否可用

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

// Echo服务的处理器
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean // 是否停止处理
}

// 表示一个客户端实体。我们把Conn和其他信息一起保存，放到activeConn，以便管理
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait // 自己封装的可以超时等待的WaitGroup
}

func (e *EchoClient) Close() error {
	_ = e.Conn.Close() // 所有Read和Write操作将立刻中止并返回error
	e.Waiting.WaitWithTimeout(10 * time.Second)
	return nil
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	cli := &EchoClient{
		Conn: conn,
	}

	// 不用依赖额外的协程来中止此函数
	if handler.closing.Get() {
		_ = cli.Close()
		return
	}

	handler.activeConn.Store(cli, struct{}{}) // 因为是地址，cli可作key (可做相等比较，其值可哈希)
	for {
		reader := bufio.NewReader(conn) // conn可以直接读，但是我们转换为reader，以便使用ReadString方法
		msg, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {
			cli.Waiting.Add(1)
			_, _ = conn.Write([]byte(msg))
			cli.Waiting.Done()
		}
		if err != nil {
			if err != io.EOF {
				logger.Warn(err)
			}
			break
		}
	}
}

func (handler *EchoHandler) Close() {
	logger.Info("echo server closing ...")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value any) bool {
		cli := key.(*EchoClient)
		_ = cli.Close()
		return true // 返回false时，停止当前对Map的遍历
	})
}
