package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()
	e.Any("/signIn", SignIn)
	fmt.Println(e.Run(":8888"))
}

func SignIn(c *gin.Context) {}
