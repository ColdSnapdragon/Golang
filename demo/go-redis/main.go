package main

import (
	"fmt"
	"go-redis/config"
	"go-redis/lib/logger"
	"go-redis/tcp"
	"os"
)

const configFile = "redis.conf"

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "go-redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	stat, err := os.Stat(configFile)
	if err == nil && !stat.IsDir() { // 存在配置文件
		config.SetupConfig(configFile)
	}
	// 否则，用config中的默认值

	logger.Info(tcp.ListenAndServeWithSignal(
		&tcp.Config{Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)},
		&tcp.EchoHandler{},
	))
}
