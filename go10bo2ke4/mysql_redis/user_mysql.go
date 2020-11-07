package mysql_redis

import (
	"fmt"
	"go10bo2ke4/init/db_init"
	"go10bo2ke4/init/log_init"
)

type User struct{
	ID int
	Email string `gorm:"type:varchar(20)"`
	PassWord string `gorm:"type:varchar(20)"`
	UserName string `gorm:"type:varchar(20)"`
}
func CreateUserTable(){
	db_init.MyDb.Exec("drop table if exists users")
	err := db_init.MyDb.AutoMigrate(&User{})
	if err != nil {
		fmt.Print("迁移创建用户数据表出错")
		panic(err)
	}
}
func AddUser(e string,p string,N string)error{
	newUser := User{Email:e,PassWord: p,UserName: N}
	err := db_init.MyDb.Create(&newUser).Error
	if err != nil {
		log_init.LogInfo(err)
		return err
	}
	return nil
}
func GetUser(u *User)(*User,error){
	var getUser User
	res := db_init.MyDb.Where("email=? and pass_word=?",u.Email,u.PassWord).Find(&getUser)
	if res.Error != nil {
		log_init.LogInfo(res.Error)
		return nil,res.Error
	}
	return &getUser,res.Error
}