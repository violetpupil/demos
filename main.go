package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"
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

type Claims struct {
	Aud    string `mapstructure:"aud"`
	Azp    string `mapstructure:"azp"`
	IAt    int    `mapstructure:"iat"`
	Exp    int    `mapstructure:"exp"`
	Iss    string `mapstructure:"iss"`
	Jti    string `mapstructure:"jti"`
	Locale string `mapstructure:"locale"`
	Nbf    int    `mapstructure:"nbf"`

	// 请仅使用 Google ID 令牌 sub 字段作为用户的标识符
	// 一个 Google 帐号在不同时间点可能有多个电子邮件地址
	Email         string `mapstructure:"email"`
	EmailVerified bool   `mapstructure:"email_verified"`
	FamilyName    string `mapstructure:"family_name"`
	GivenName     string `mapstructure:"given_name"`
	Name          string `mapstructure:"name"`
	Picture       string `mapstructure:"picture"`
	Sub           string `mapstructure:"sub"`
}

// signIn 处理token
// https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
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

	var claims Claims
	err = mapstructure.Decode(payload.Claims, &claims)
	if err != nil {
		logger.Errorln("decode claims error", err)
		return err
	}
	logger.Infof("claims %+v", claims)
	return nil
}
