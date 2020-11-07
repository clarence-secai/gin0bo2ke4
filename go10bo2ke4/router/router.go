package router

import (
	"github.com/gin-gonic/gin"
	"go10bo2ke4/handlers"
)

func Routers()*gin.Engine{
	engine := gin.Default()
	//显示博客目录
	engine.GET("/index",handlers.GetLists)
	//显示某篇博客及相应的评论
	engine.GET("/index/:id",handlers.GetArticle)

	//注册
	engine.POST("/register",handlers.Register)
	//上传头像
	//engine.POST("/uploadPicture",handlers.UploadPicture)
	//登录	并返给前端jwt token，并让前端区分游客和博主自己
	engine.POST("/login",handlers.LoginJwt)

	//全局中间件，校验jwt token并取得jwt token里的唯一辨识用户的信息
	engine.Use(handlers.VerifyJwt)

	//判断必须是可以发博客的才有接下来的增删改
	groupArticle := engine.Group("/article", handlers.Me)
	//添加一篇博客文章
	groupArticle.POST("/addArticle",handlers.AddArticle)
	//修改谋篇博客文章
	groupArticle.PUT("/updateArticle",handlers.UpdateArticle)
	//删除谋篇博客文章
	groupArticle.DELETE("/deleteArticle",handlers.DeleteArticle)

	//在某篇博客下面进行评论
	engine.POST("/comment/addComment",handlers.AddComment)
	//删除谋篇博客下自己发的评论
	engine.DELETE("/comment/deleteComment",handlers.DeleteComment)

	return engine
}
