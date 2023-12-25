package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.LoadHTMLGlob("*.html")
	e.GET("/", Index)
	e.POST("/signIn", SignIn)
	fmt.Println(e.Run(":8888"))
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"clientId": Config.ClientID})
}

func SignIn(c *gin.Context) {}
