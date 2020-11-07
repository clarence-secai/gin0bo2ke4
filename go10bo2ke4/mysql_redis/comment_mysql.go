package mysql_redis

import (
	"go10bo2ke4/init/db_init"
	"go10bo2ke4/init/log_init"
)

type Comment struct{
	CommentID int //`gorm:"type:bigint(20);primaryKey;autoIncrement"`//todo 为何autoMigrate或autoIncrement:true标签未起效？
	Sentence string //`gorm:"type:varchar(400)"`
	CreatedOn string //`gorm:"type:varchar(20)"`
	CreatedBy string //`gorm:"type:varchar(20)"`
	UUID string //`gorm:"type:varchar(50)"`  //跟博客表联系起来的字段
}
func CreateCommentTable(){
	db_init.MyDb.Exec("drop table if exists comments")
	//todo 手动建表时表、列命名对应上gorm的蛇形规则，让gorm的crud能用
	db_init.MyDb.Exec(`create table comments (
	comment_id bigint(20) primaryKey auto_Increment, 
	sentence varchar(400) default null,
	created_on varchar(20) default null,
	created_by varchar(20) default null,
	uuid varchar(50) default null)`)
	//err := db_init.MyDb.AutoMigrate(&Comment{})
	//if err != nil {
	//	panic("迁移创建评论表出错")
	//}
}
func AddComment(co Comment)error{
	err := db_init.MyDb.Create(&co).Error
	if err != nil {
		log_init.LogInfo(err)
		return err
	}
	return nil
}
func GetComment(uuid string)(*Comment,error){
	var co Comment
	err := db_init.MyDb.Where("uuid=?", uuid).Find(&co).Error
	if err != nil {
		log_init.LogInfo(err)
		return nil,err
	}
	return &co,nil
}
func DeleteCommentByCommentId(id int)int64{
	res := db_init.MyDb.Where("Comment_iD=?",id).Delete(&Comment{})
	return res.RowsAffected
}
func DeleteCommentByUUID(uuid string)(int64,error){
	res := db_init.MyDb.Where("uuid=?",uuid).Delete(&Comment{})
	return res.RowsAffected,res.Error
}
func GetCommentById(id int)(*Comment,error){
	var co Comment
	res := db_init.MyDb.Where("Comment_iD=?", id).Find(&co)
	if res.Error != nil {
		log_init.LogInfo(res.Error)
		return nil,res.Error
	}
	return &co,nil
}