package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"marking/response"
	"marking/util"
	"net/http"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	fmt.Println(c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, Authorization")
	c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	c.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func VerifyToken(c *gin.Context) {
	//获取请求path
	obj := c.Request.URL.Path
	//获取请求方法
	act := c.Request.Method
	var tokenString string
	tokenString = c.GetHeader("Authorization")
	t, claims, err := util.ParseToken(tokenString)
	if err != nil || !t.Valid {
		response.Fail(c.Writer, "token lapsed", 403)
		c.Abort()
		return
	}
	sub := claims.UID
	//引入casbin
	e := util.CasbinServiceApp.Casbin()
	success, _ := e.Enforce(fmt.Sprintf("%d", sub), obj, act)
	if sub == 0 || success {
		c.Next()
	} else {
		response.Fail(c.Writer, "权限不足", 403)
		c.Abort()
		return
	}
}
