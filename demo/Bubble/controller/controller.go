package controller

import (
	"bubble/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// controller负责处理 HTTP 请求和响应

func IndexHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func CreatTodo(c *gin.Context) {
	var item models.Todo
	_ = c.BindJSON(&item)
	if err := models.AddOneTodo(&item); err != nil {
		c.JSON(http.StatusOK, gin.H{ // 操作虽然失败，但是响应是成功的，返回200
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, item)
	}
}

func GetTodoList(c *gin.Context) {
	var itemList []models.Todo
	itemList, err := models.GetAllTodo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, itemList)
	}
}

func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	// 先查库，再修改
	if _, err := models.GetOneTodo(id); err != nil { // 可以通过Error查看数据库操作是否成功
		c.JSON(http.StatusOK, gin.H{
			"error": "ID 不存在",
		})
	} else {
		var item models.Todo
		_ = c.BindJSON(&item) // 传来的JSON里只提供了status，那么其他字段会被赋为零值
		// status的值具体是什么、要不要对数据库里的值作反转——这些是前端考虑的，我们只负责拿到数据完成相应的动作
		if err := models.UpdateOneTodo(id, &item); err != nil { // Updates接受结构体对象并且值更新非零值
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": item,
			})
		}
	}
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	// 先查库，再修改
	if _, err := models.GetOneTodo(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "ID 不存在",
		})
	} else {
		if err = models.DelOneTodo(id); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":   err.Error(),
				"message": "删除失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"delete": id,
			})
		}
	}
}
