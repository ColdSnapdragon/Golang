package routers

import (
	"bubble/controller"
	"github.com/gin-gonic/gin"
)

// 包含路由

// 生成我们注册好的路由

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 本demo就不做前后端分离了，直接返回页面
	r.Static("/static", "static") // 告诉gin去哪找模板文件引用的静态文件(相对当前位置的路径)
	r.LoadHTMLGlob("templates/*") // 告诉gin去哪找模板文件(相对go.mod的路径)
	r.GET("/", controller.IndexHandle)

	v1Group := r.Group("/v1")
	{
		// 新增条目
		v1Group.POST("/todo", controller.CreatTodo)

		// 查看所有条目
		v1Group.GET("/todo", controller.GetTodoList)

		// 查看某一条目

		// 修改条目
		v1Group.PUT("/todo/:id", controller.UpdateTodo)

		// 删除条目
		v1Group.DELETE("/todo/:id", controller.DeleteTodo)
	}
	return r
}
