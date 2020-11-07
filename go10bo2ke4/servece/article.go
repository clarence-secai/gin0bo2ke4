package servece

import (
	"encoding/json"
	"errors"
	"fmt"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/mysql_redis"
)

func GetLists()(*[]mysql_redis.List,int,error){
	//先到缓存中查找
	exist := mysql_redis.Exist("getLists")
	if exist {
		res := mysql_redis.Get("getLists")
		resStr,_ := res.(string)
		des := make(map[string]interface{})
		json.Unmarshal([]byte(resStr),&des)
		lists := des["lists"].(*[]mysql_redis.List)
		listNum := des["listNum"].(int)
		return lists,listNum,nil
	}
	//缓存中不存在，则到数据库中查找
	lists,listNum,err := mysql_redis.GetCatalog()
	if err != nil {
		return nil,-1,err
	}
	if lists == nil || listNum == 0 {
		return nil,0,nil
	}

	//存入缓存中
	redisCont := map[string]interface{}{
		"listNum":listNum,
		"lists":*lists,
	}
	resByte,err := json.Marshal(redisCont)
	if err != nil {
		log_init.LogInfo(errors.New("序列化存入缓存出错"))
	}else{
		mysql_redis.Set("getLists",string(resByte))
	}

	return lists,listNum,nil
}

func GetArticle(id int)(*mysql_redis.Essay,*mysql_redis.Comment,error){
	//先从缓存中查询获取
	exist := mysql_redis.Exist(fmt.Sprintf("blog_%d",id))
	if exist {
		res := mysql_redis.Get(fmt.Sprintf("blog_%d",id))
		resStr,_ := res.(string)
		des := make(map[string]interface{})
		json.Unmarshal([]byte(resStr),&des)
		essay,ok := (des["essay"]).(*mysql_redis.Essay)
		co,ok := (des["co"]).(*mysql_redis.Comment)
		if !ok {
			log_init.LogInfo(errors.New("断言失败"))
		}
		return essay,co,nil
	}

	//缓存没有就从数据库获取博客文章信息
	essay, err := mysql_redis.GetEssayById(id)
	if err != nil {
		return nil,nil,err
	}
	if essay == nil {
		return nil,nil,errors.New("不存在该博客")
	}
	//获取匹配的评论信息
	co,err := mysql_redis.GetComment(essay.UUID)
	if err != nil {
		return nil,nil,err
	}
	//存入缓存
	redisCon := map[string]interface{}{
		"essay":*essay,
		"co":*co,
	}
	resByte,err := json.Marshal(redisCon)
	if err != nil {
		log_init.LogInfo(errors.New("序列化存入缓存出错"))
	}else{
		mysql_redis.Set(fmt.Sprintf("blog_%d",essay.ID),string(resByte))
	}

	return essay,co,nil
}

//func AddArticle()error{
//}

func UpdateEssay(upCle *mysql_redis.ModifyArticle)error{
	//更新数据库表中的博客
	err := mysql_redis.UpdateEssay(upCle)
	if err != nil {
		return err
	}
	//删除缓存中的该博客
	exist := mysql_redis.Exist(fmt.Sprintf("blog_%d", upCle.ID))
	if exist {
		mysql_redis.Del(fmt.Sprintf("blog_%d", upCle.ID))
	}
	return nil
}

func DeleteArticle(id int)(int64,error){
	//删除博客
	err := mysql_redis.DeleteEssay(id)
	if err != nil {
		log_init.LogInfo(fmt.Errorf("%d",-1))
		return -1,err
	}
	//删除博客下的评论
	essay, err2 := mysql_redis.GetEssayById(id)
	if err2 != nil {
		log_init.LogInfo(fmt.Errorf("%d",-2))
		return -2,err
	}
	num,err := mysql_redis.DeleteCommentByUUID(essay.UUID)
	if err != nil {
		log_init.LogInfo(fmt.Errorf("%d",-3))
		return -3,err
	}
	//删除缓存中的该博客和评论
	exist := mysql_redis.Exist(fmt.Sprintf("blog_%d", id))
	if exist {
		mysql_redis.Del(fmt.Sprintf("blog_%d",id))
		return num,nil
	}
	log_init.LogInfo(fmt.Errorf("%d",-4))
	return -4,err
}
