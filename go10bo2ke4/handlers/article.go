package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/myerr"
	"go10bo2ke4/mysql_redis"
	"go10bo2ke4/servece"
	"go10bo2ke4/util"
	"net/http"
	"strconv"
)

//获取博客目录
func GetLists(c *gin.Context){
	lists, listNum, err := servece.GetLists()
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,nil)
	}
	//将结果发送给客户端
	util.Response(c,http.StatusOK,myerr.OK,gin.H{
		"listNum":listNum,
		"lists":*lists,
	})
}

//添加一篇博客文章
func AddArticle(c *gin.Context){
	var newArticle mysql_redis.NewArticle
	err := c.ShouldBind(&newArticle)
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusNotAcceptable,myerr.Err_PrefaceData_Get_Fail,nil)
		return
	}
	//新增一篇博客
	err = mysql_redis.AddEssay(&newArticle)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Add_Fail,nil)
	}

	util.Response(c,http.StatusOK,myerr.OK,"新增博客文章成功")
}

//获取一篇博客和相应的评论
func GetArticle(c *gin.Context){
	//获取博客id
	idStr := c.Param("id")
	if idStr==""{
		util.Response(c,http.StatusBadRequest,myerr.Err_PrefaceData_null,"请传送非空id")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"数据类型转换出错")
	}
	//查询博客
	essay, co, err := servece.GetArticle(id)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Select_Fail,nil)
	}

	util.Response(c,http.StatusOK,myerr.OK,gin.H{
		"id":essay.ID,
		"title":essay.Title,
		"createOn":essay.CreateOn,
		"content":essay.Content,

		"commentId":co.CommentID,
		"sentence":co.Sentence,
		"createdOn":co.CreatedOn,
		"createdBy":co.CreatedBy,
	})
}

//修改一篇博客
func UpdateArticle(c *gin.Context){
	var updateArticle mysql_redis.ModifyArticle
	err := c.ShouldBind(&updateArticle)
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusNotAcceptable,myerr.Err_PrefaceData_Get_Fail,nil)
		return
	}
	//修改数据库表中的博客，删除缓存中的该博客
	err = servece.UpdateEssay(&updateArticle)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Update_Fail,nil)
	}
	util.Response(c,http.StatusOK,myerr.OK,nil)
}

//删除某篇博客及相应评论
func DeleteArticle(c *gin.Context){
	//获取即将被删除博客的id
	idStr := c.PostForm("id")
	if idStr==""{
		util.Response(c,http.StatusBadRequest,myerr.Err_PrefaceData_null,"请传送非空id")
		return
	}
	id, err2 := strconv.Atoi(idStr)
	if err2 != nil {
		log_init.LogInfo(err2)
		util.Response(c,http.StatusNotAcceptable,myerr.Err_PrefaceData_Get_Fail,nil)
		return
	}
	//删除博客和评论及缓存中的博客和评论
	num, err2 := servece.DeleteArticle(id)
	if err2 != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Delete_Fail,"删除数据库表数据出错")
		return
	}
	util.Response(c,http.StatusOK,myerr.OK,fmt.Sprintf("成功删除1篇博客和%d条评论",num))
}