package main

import (
	"fmt"
	"go10bo2ke4/init/db_init"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/init/redis_init"
	"go10bo2ke4/mysql_redis"
	"go10bo2ke4/router"
	"net/http"
	"time"
)
//初始化错误日志、MySQL数据库、redis缓存数据库
func init(){
	var a string
	fmt.Printf("你创建了名为 bo2ke4 的数据库吗？y/n")
	fmt.Scanln(&a)
	if a !="y"{
		fmt.Println("请先创建名为 bo2ke4 的数据库，再重新运行本程序")
		return
	}
	log_init.GetLog()
	db_init.GetDb()
	//上一行已经获得db，下面创建数据表，避免手动通过海豚客户端创建表
	mysql_redis.CreateEssayTable()
	mysql_redis.CreateUserTable()
	mysql_redis.CreateCommentTable()
	redis_init.GetRedis()
}
func main() {
	fmt.Println("init初始化完毕，开始main")

	engine:= router.Routers()
	//方式一
	//engine.Run(":8888")
	//方式二，使用官方文档中的，自行配置一些参数。当然也可一并写入配置文件中，提高扩展性
	server := http.Server{
		Addr:":8888",
		Handler:engine,	//gin框架的engine旗下有 ServeHTTP(ResponseWriter, *Request)方法，故可赋值给http原生包下的该字段
		ReadHeaderTimeout:10*time.Second,
		WriteTimeout:time.Second*10,
		IdleTimeout:1*time.Minute,
	}
	server.ListenAndServe()
}


