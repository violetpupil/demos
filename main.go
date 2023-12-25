package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// Body token表单请求体
type Body struct {
	Credential string `form:"credential"`
	GCsrfToken string `form:"g_csrf_token"`
}

// SignIn 处理token请求
// @router /signIn [post]
func SignIn(c *gin.Context) {
	var body Body
	err := c.ShouldBind(&body)
	if err != nil {
		logrus.Errorln("bind error", err)
		goto FIN
	}
	signIn(body)

FIN:
	if err != nil {
		c.String(http.StatusOK, "something error")
	} else {
		c.String(http.StatusOK, "ok")
	}
}

func signIn(body Body) {
	logrus.Infof("sign in body %+v\n", body)
}
