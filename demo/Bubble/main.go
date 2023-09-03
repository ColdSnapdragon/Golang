package main

import (
	"bubble/dao"
	"bubble/models"
	"bubble/routers"
)

func main() {

	// 先连接数据库，若失败则退出
	if err := dao.ConnectMySQL(); err != nil {
		panic(err)
	}
	if err := dao.DB.AutoMigrate(&models.Todo{}); err != nil {
		panic(err)
	}

	r := routers.SetupRouter()

	_ = r.Run(":9090")
}
