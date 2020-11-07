package redis_init

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go10bo2ke4/util"
	"os"
	"time"
)

//结构体字段对应着json配置文件
type Redi struct{
	NetWork string
	Address string
	Port string
	MaxActive int
	MaxIdle int
	IdleTimeout time.Duration
}
var MyRedisPool redis.Pool
func GetRedis(){
	fullPath := util.WdirAndPath("config/redis_config.json")
	//读取redis配置文件
	file, err := os.Open(fullPath)
	if err != nil {
		panic("打开redis配置文件失败")
	}
	decoder := json.NewDecoder(file)
	var myRedis Redi
	err = decoder.Decode(&myRedis)
	if err != nil {
		panic("解码redis配置文件失败")
	}

	//初始化redis连接池
	MyRedisPool = redis.Pool{
		MaxActive:   myRedis.MaxActive,
		MaxIdle:     myRedis.MaxIdle,
		IdleTimeout: myRedis.IdleTimeout*time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(myRedis.NetWork, myRedis.Address+":"+myRedis.Port)
			if err != nil {
				fmt.Print("获取redis链接失败")
				return nil,err
			}
			return conn,nil
		},
	}
}
