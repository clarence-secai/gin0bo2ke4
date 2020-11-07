package handlers

import (
	"github.com/gin-gonic/gin"
	"go10bo2ke4/init/runtime_init"
	"go10bo2ke4/myerr"
	"go10bo2ke4/util"
	"net/http"
)

func UploadPicture(c *gin.Context){
	fileHeader, err := c.FormFile("uploadPicture")
	if err != nil {
		util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"上传头像失败")
		return
	}
	//修改头像文件名，跟用户一一对应
	//pictureName := fileHeader.Filename
	//extenName := path.Base(pictureName)
	c.PostForm("tokenContent")

	//读取配置文件中存放头像的路径
	respectivePath := runtime_init.MyPath().HeadImagePath
	des := util.WdirAndPath(respectivePath)//由相对路径变成使用该项目机器的全路径
	//保存头像图片
	err = c.SaveUploadedFile(fileHeader, des)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"上传头像失败")
	}
}
