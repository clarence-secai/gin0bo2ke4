package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/myerr"
	"go10bo2ke4/mysql_redis"
	"go10bo2ke4/util"
	"net/http"
	"strconv"
	"time"
)

//在相应博客下新增一条评论
func AddComment(c *gin.Context){
	//获取被评论的博客的ID
	articleId, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusNotAcceptable,myerr.Err_PrefaceData_Get_Fail,nil)
		return
	}
	//获取评论者评论的内容
	sentence := c.PostForm("sentence")
	if sentence == ""{
		util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"评论不能为空")
		return
	}
	//获取评论者的用户信息，从verifyJwt中间件向后传递的用户信息获取
	user,ok := c.Get("tokenContent")
	if !ok {
		log_init.LogInfo(errors.New("获取中间件向后传递的值出错"))
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,nil)
	}
	u,ok := user.(User)//断言，获取评论者的用户信息
	if !ok {
		log_init.LogInfo(errors.New("断言失败"))
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,nil)
	}
	//查询数据库获取评论者的完整用户信息
	dbU := mysql_redis.User{
		Email:    u.Email,
		PassWord: u.PassWord,
	}
	getUser,err := mysql_redis.GetUser(&dbU)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Select_Fail,nil)
		return
	}

	//根据前端传来的博客文章的id查询到相应的博客文章，便于下面使用其UUID
	essay, err := mysql_redis.GetEssayById(articleId)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Select_Fail,nil)
		return
	}

	co := mysql_redis.Comment{
		Sentence: sentence,
		CreatedOn: time.Now().Format("2006-01-02/15:04"),
		CreatedBy:getUser.UserName,
		UUID:essay.UUID,
	}

	//向数据库插入一条与相应博客文章UUID对应的评论
	err = mysql_redis.AddComment(co)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Add_Fail,"添加评论进数据库出错")
		return
	}

	//删除缓存中该篇博客及其评论
	exist := mysql_redis.Exist(fmt.Sprintf("blog_%d", articleId))
	if exist{
		mysql_redis.Del(fmt.Sprintf("blog_%d", articleId))
	}

	util.Response(c,http.StatusOK,myerr.OK,"成功添加评论")
}

func DeleteComment(c *gin.Context){
	idCommentStr := c.PostForm("commentId")
	commentId,err := strconv.Atoi(idCommentStr)
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusNotAcceptable,myerr.Err_PrefaceData_Get_Fail,nil)
		return
	}

	//删除缓存中该篇博客及其评论
	articleId, err := mysql_redis.GetEssayIdByCommentId(commentId)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"删除缓存出错")
		return
	}
	exist := mysql_redis.Exist(fmt.Sprintf("blog_%d", articleId))
	if exist{
		mysql_redis.Del(fmt.Sprintf("blog_%d", articleId))
	}

	//删除数据库表中的某博客下的该条评论
	res := mysql_redis.DeleteCommentByCommentId(commentId)
	if res == 1 {
		util.Response(c,http.StatusOK,myerr.OK,nil)
	}else{
		log_init.LogInfo(errors.New("删除数据库评论出错"))
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,nil)
	}
}
