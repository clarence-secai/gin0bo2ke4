package mysql_redis

import (
	"fmt"
	"go10bo2ke4/init/db_init"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/util"
	"gorm.io/gorm"
	"time"
)

type Essay struct{
	ID int //gorm默认为主键和自增 id
	Title string  `gorm:"size:64"` //等同于`gorm:"varchar(64)"` 即不多于64个字符
	CreateOn string `gorm:"size:64"`
	Content string `gorm:"type:text"`//注意这里前端在填的文章内容中不能含奇怪的格式
	UUID string `gorm:"type:varchar(50)"`  //和文章评论的表格关联
}
//todo 下面这个CreateTable很重要，自己写完后再看徐新华蓝色书中发现他是手动创建的表
func CreateEssayTable(){
	db_init.MyDb.Exec("drop table if exists essays")
	//todo 迁移创建数据表，如已经有该表则该行代码什么也不会做，不用担心会重复创建表
	err := db_init.MyDb.AutoMigrate(&Essay{})
	if err != nil {
		fmt.Println("迁移创建博客数据表出错=", err)
		panic(err)
	}
}

//新增一篇博客
type NewArticle struct{
	Title string 	`form:"title"`
	Content string	`form:"content"`
}
func AddEssay(newArticle *NewArticle)error{
	newEssay := Essay{
		Title: (*newArticle).Title,
		Content: (*newArticle).Content,
		CreateOn: time.Now().Format("2006/01/02/15:04"),
		UUID:util.CreateUUID(),
	}

	err := db_init.MyDb.Create(&newEssay).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log_init.LogInfo(err)
		return err
	}
	return nil
}

//博客目录
type List struct{
	ID int
	Title string
}
func GetCatalog()(*[]List,int,error){
	var catalog []List
	err := db_init.MyDb.Table("essays").Find(&catalog).Error
	if err != nil {
		log_init.LogInfo(err)//在发生错误的源头将错误写入日志
		return nil,0,err
	}
	listNum := len(catalog)
	return &catalog,listNum,nil
}

//查询谋篇博客，不含评论
func GetEssayById(id int)(*Essay,error){
	var essay Essay
	err := db_init.MyDb.Find(&essay, id).Error
	if err != nil {
		log_init.LogInfo(err)
		return nil,err
	}
	return &essay,nil
}

//更新博客
type ModifyArticle struct{
	ID int	`form:"id"`
	Title string `form:"title"`
	Content string `form:"content"`
}
func UpdateEssay(m *ModifyArticle)error{
	updateEssay := Essay{
		ID:(*m).ID,
		Title: (*m).Title,
		Content:(*m).Content,
	}
	err := db_init.MyDb.Model(&Essay{}).Where("id=?",updateEssay.ID).Updates(map[string]interface{}{
		"title":updateEssay.Title,
		"content":updateEssay.Content,
		}).Error
	if err != nil {
		log_init.LogInfo(err)
		return err
	}
	return nil
}

//删除某篇博客
func DeleteEssay(id int)error{
	err := db_init.MyDb.Delete(Essay{}, id).Error
	if err != nil {
		log_init.LogInfo(err)
		return err
	}
	return nil
}
func GetEssayByUUID(uuid string)*Essay{
	var essay Essay
	db_init.MyDb.Find(&essay,"uuid=?",uuid)
	return &essay
}
//根据评论的id找到所属博客的id
func GetEssayIdByCommentId(commentId int)(int,error){
	co, err := GetCommentById(commentId)
	if err != nil {
		return 0,err
	}
	return GetEssayByUUID(co.UUID).ID,nil

}