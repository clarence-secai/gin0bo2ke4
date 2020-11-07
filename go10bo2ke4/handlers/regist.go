package handlers

import (
	"github.com/gin-gonic/gin"
	"go10bo2ke4/myerr"
	"go10bo2ke4/mysql_redis"
	"go10bo2ke4/util"
	"net/http"
)

func Register(c *gin.Context){
	E := c.PostForm("email")
	P := c.PostForm("passWord")
	N := c.PostForm("userName")
	err := mysql_redis.AddUser(E,P,N)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"注册出错")
	}
	util.Response(c,http.StatusOK,myerr.OK,"注册成功")
}
