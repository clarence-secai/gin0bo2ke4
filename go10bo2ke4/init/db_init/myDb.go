package db_init

import (
	"encoding/json"
	"fmt"
	"go10bo2ke4/util"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

//结构体字段对应着json配置文件
type sql struct{
	User string
	Password string
	Url string
	Port string   //todo:这里一定不能写int类型，配置文件json中3306也最好加上双引号
	DataBaseName string
	ChartSet string
}
//读取db_config.json配置文件并初始化数据库
var MyDb *gorm.DB
func GetDb(){
	//todo:下一行的天的短路径要去掉该项目的第一个目录即go10bo2ke4/ ，从下一个开始写，
	// 因为获取工作目录的函数回获取到被下载到的系统中根目录到/go10bo2ke4为止
	fullPath := util.WdirAndPath("config/db_config.json")
	//读取配置文件
	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("初始化db配置文件失败")
	}
	decoder:= json.NewDecoder(file)
	var meMysql sql
	decoder.Decode(&meMysql)
	fmt.Println("memysql=",meMysql)

	//todo:获取MyDb的过程并不会自动创建/和?之间的数据库，故下面需要自己手动建库
	//手动通过海豚建立数据库

	//打开数据库
	des := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",meMysql.User,
		meMysql.Password,meMysql.Url,meMysql.Port,meMysql.DataBaseName,meMysql.ChartSet)
	MyDb, err = gorm.Open(mysql.Open(des), &gorm.Config{})
	if err != nil {
		fmt.Println("打开数据库失败")
		return
	}




}
