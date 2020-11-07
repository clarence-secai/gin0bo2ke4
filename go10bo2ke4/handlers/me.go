package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go10bo2ke4/myerr"
	"go10bo2ke4/util"
	"net/http"
)

func Me(c *gin.Context){
	u,ok := c.Get("tokenContent")
	if !ok {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"token向后传递出错")
		return
	}
	//断言，获取用户的信息，判断是否是可以写博客的用户
	user,ok := u.(User)
	if !ok {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"token断言失败")
		return
	}
	fmt.Println("user=",user)
	if user.PassWord != "413188ok" || user.Email != "1787664981@qq.com" {
		util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"你没有权限发博客")
		c.Abort()
		return
	}
	c.Next()
}
