package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
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
	var cookieCsrfToken string
	err := c.ShouldBind(&body)
	if err != nil {
		logrus.Errorln("bind error", err)
		goto FIN
	}
	cookieCsrfToken, err = c.Cookie("g_csrf_token")
	if err != nil {
		logrus.Errorln("get cookie error", err)
		goto FIN
	}
	err = signIn(body, cookieCsrfToken)
	if err != nil {
		logrus.Errorln("process error", err)
		goto FIN
	}

FIN:
	if err != nil {
		c.String(http.StatusOK, "something error")
	} else {
		c.String(http.StatusOK, "ok")
	}
}

var (
	ErrCsrfTokenValidateFail = errors.New("csrf token validate fail")
	ErrIdtokenValidateFail   = errors.New("idtoken validate fail")
)

func signIn(body Body, cookieCsrfToken string) error {
	logger := logrus.WithFields(logrus.Fields{
		"body":            fmt.Sprintf("%+v", body),
		"cookieCsrfToken": cookieCsrfToken,
	})
	logger.Infoln("sign in argument")

	if cookieCsrfToken == "" || body.GCsrfToken == "" ||
		cookieCsrfToken != body.GCsrfToken {
		return ErrCsrfTokenValidateFail
	}

	payload, err := idtoken.Validate(context.Background(), body.Credential, Config.ClientID)
	if err != nil {
		logger.Errorln("idtoken validate error", err)
		return ErrIdtokenValidateFail
	}
	logger.Infof("payload %+v", payload)
	return nil
}
