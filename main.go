package main

import (
	"web-service-gin/logic"

	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()
	router.GET("/books", logic.GetBook)
	router.GET("/books/:id", logic.GetBookByID)
	router.DELETE("/books/:id", logic.DeleteBook)
	router.POST("/books", logic.PostBook)
	router.PUT("/books/:id",logic.UpdateBook)
	router.LoadHTMLGlob("view/books.html")
	router.GET("/view", func(c *gin.Context) {
		c.HTML(200, "books.html", nil)
	})

	router.Run("localhost:8082")
}
