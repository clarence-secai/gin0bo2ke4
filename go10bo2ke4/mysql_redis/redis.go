package mysql_redis

import (
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/init/redis_init"
)

func Set(key string,value interface{}){
	conn := redis_init.MyRedisPool.Get()
	defer conn.Close()  //todo:用完一个conn后还给连接池pool
	conn.Do("set",key,value)
	conn.Do("expire",key,1200)
}
func Get(key string)interface{}{
	conn := redis_init.MyRedisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("get", key)
	if err != nil {
		log_init.LogFatal(err)
	}
	return reply
}
func Exist(key string)bool{
	conn := redis_init.MyRedisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("exists", key)
	if err != nil {
		log_init.LogInfo(err)
	}
	if reply==1{
		return true
	}
	return false
}
func Del(key string)bool{
	conn := redis_init.MyRedisPool.Get()
	reply, err := conn.Do("del", key)
	if err != nil {
		log_init.LogFatal(err)
	}
	if reply== 0{
		return false
	}
	return true
}