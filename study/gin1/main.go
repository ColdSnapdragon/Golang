package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // 创建一个默认的路由引擎
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{ // ctx.JSON方法将任何一个结构体对象(这里用gin.H)序列化为JSON格式，并将其作为响应体返回给客户端
			"message": "hello",
		})
	})

	// 另一种方式：对结构体做json序列化
	type msg struct {
		Name string
		nick string // 只有导出的结构体成员(大写字母开头)才会被编码
		Age  int    `json:"age"` // 添加结构体的成员Tag
		// json开头键名对应的值用于控制encoding/json包的编码和解码的行为，这里用于修改名字
	}
	r.GET("/message", func(c *gin.Context) {
		c.JSON(http.StatusOK, msg{
			Name: "Blover",
			nick: "zrt", // 不会参与json序列化
			Age:  20,
		})
	})

	// GET请求URL?后面的是querystring参数
	// 介绍三种常见分析querystring参数的函数
	r.GET("/web", func(c *gin.Context) {
		name := c.Query("name")            // 返回字符串，无值时返回""
		age := c.DefaultQuery("age", "20") // 无值时返回第二个参数
		nick, ok := c.GetQuery("nick")     // 第二个返回值表示是否得值
		if !ok {
			nick = "none"
		}
		// 127.0.0.1:9090/web?name=blover&age=20
		c.JSON(http.StatusOK, gin.H{ // map转json是有序的
			"name": name,
			"age":  age,
			"nick": nick,
		})
	})

	// POST将数据放在HTTP请求的消息体中，而不是放在URL，可以传递更多数据
	// *一个请求对应一个相应*。不同请求访问同一个url的结果可以不一样
	r.POST("/web", func(c *gin.Context) {
		name := c.PostForm("name")
		passwd := c.PostForm("password")
		// 剩余两种函数的作用与上文相同，不赘述
		// c.DefaultPostForm("name")
		// c.GetPostForm("name")
		c.JSON(http.StatusOK, gin.H{
			"name":     name,
			"password": passwd,
		})
		// 在postman的POST请求中填写Body->form-data
	})

	// 匹配并获取url路径参数
	r.GET("/get/:age1/:age2", func(c *gin.Context) {
		// 127.0.0.1:9090/get/ice/2023
		a1 := c.Param("age1")
		a2 := c.Param("age2")
		//fmt.Println(a1, a2)
		c.JSON(http.StatusOK, [...]string{a1, a2})
	})

	type One struct {
		Nick  string
		Level int `json:"lv" form:"level"` // 对于json数据，识别标志为"lv"
	}
	// 注意GET请求没有Body，只能绑定string-query的参数值
	r.POST("/bind", func(c *gin.Context) {
		var o One
		// ShouldBind函数根据传入对象的各字段名(的tag)，去请求中查找值(且类型必须能与字段类型匹配)
		// 可以自动选择绑定方式(string-query JSON Form-Data XML)
		if err := c.ShouldBind(&o); err != nil { // 传指针
			fmt.Println(err)
		}
		fmt.Printf("%#v\n", o)
	})

	// 处理任何匹配不到路由的请求
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "This is a blank page (404)",
		})
	})
	// 处理发到所选路径的任何请求(Any函数的源码中，就是把该请求转发给所有请求类型)
	r.Any("/any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"method": c.Request.Method, // 请求类型(GET,POST,...)
		})
	})

	// 请求重定向
	r.GET("/baidu", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})
	// 请求转发(内部)
	r.GET("/forward", func(c *gin.Context) {
		r.HandleContext(c) // 转发给另一个路由(自己)处理
	})

	homeGroup := r.Group("/home")
	{
		// homeGroup基本可当做正常路由使用(甚至继续套娃)
		homeGroup.GET("/mine", func(c *gin.Context) { // 表示/home/id
			c.JSON(http.StatusOK, gin.H{
				"Path": c.Request.URL.Path,
				"Host": c.Request.URL.Host,
				"Id":   123,
			})
		})
	}

	// 如果对同一个路径注册了多个GET方法，那么只处理最后一个

	err := r.Run(":9090") // 无参时默认本地8080
	// 当前goroutine阻塞进入监听
	fmt.Println(err)
}
