package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"stackgen-todo/core/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func EngineRoutes(router *gin.RouterGroup, todoList *controller.Todo) {

	router.POST("/todo", func(c *gin.Context) {
		CreateToDoItem(c, todoList)
	})

	router.GET("/todo/:id", func(c *gin.Context) {
		GetTodoItem(c, todoList)
	})

	router.PUT("/todo/:id", func(c *gin.Context) {
		UpdateTodoItem(c, todoList)
	})

	router.GET("/todo", func(c *gin.Context) {
		GetAllTodoItems(c, todoList)
	})

	router.PUT("/todoAttachment/:id", func(c *gin.Context) {
		AddAttachment(c, todoList)
	})

	// //  GET metrics route for getting metrics data
	// router.GET("/metrics", func(c *gin.Context) {
	// 	GetMetrics(c, metrics)
	// })
}

func CreateToDoItem(c *gin.Context, todoList *controller.Todo) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil && err != io.EOF {
		fmt.Println("Error while reading request body: ", err)
		c.JSON(400, "Bad Request")
		return
	}
	fmt.Println(string(body))
	jsonBody := &controller.Item{}
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		fmt.Println("Error while decoding request body: ", err)
		c.JSON(400, "Bad Request")
		return
	}

	resp, err := controller.CreateToDoItem(jsonBody, todoList)
	if err != nil {
		fmt.Println("Error ", err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp)
}

func GetTodoItem(c *gin.Context, todoList *controller.Todo) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, "Id not provided")
		return
	}

	resp, err := controller.GetTodoItem(cast.ToInt(id), todoList)
	if err != nil {
		fmt.Println("Error ", err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp)
}

func AddAttachment(c *gin.Context, todoList *controller.Todo) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, "Id not provided")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
		c.JSON(400, "Bad Request")
		return
	}

	files := form.File["attachments"]

	if len(files) == 0 {
		c.JSON(400, "No files uploaded")
		return
	}

	err = controller.AddAttachment(cast.ToInt(id), todoList, files)
	if err != nil {
		fmt.Println("Error ", err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, nil)
}

func GetAllTodoItems(c *gin.Context, todoList *controller.Todo) {
	resp, err := controller.GetAllItems(todoList)
	if err != nil {
		fmt.Println("Error ", err)
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp)
}

func UpdateTodoItem(c *gin.Context, todoList *controller.Todo) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, "Id not provided")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil && err != io.EOF {
		fmt.Println("Error while reading request body: ", err)
		c.JSON(400, "Bad Request")
		return
	}
	fmt.Println(string(body))
	jsonBody := &controller.Item{}
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		fmt.Println("Error while decoding request body: ", err)
		c.JSON(400, "Bad Request")
		return
	}

	controller.UpdateTodoItem(cast.ToInt(id), jsonBody, todoList)
	c.JSON(200, nil)
}
