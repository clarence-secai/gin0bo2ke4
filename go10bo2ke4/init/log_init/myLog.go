package log_init

import (
	"fmt"
	"go10bo2ke4/init/runtime_init"
	"go10bo2ke4/util"
	"log"
	"path"
	"runtime"
	"time"
)
var logger *log.Logger
//创建日志文件
func GetLog(){
	//读取path_config.json配置文件到结构体
	runtimePath := runtime_init.MyPath()
	//创建运行时所需的目录和文件
	logFileName := time.Now().Format("20060102")//用当天的日期作为日志的文件名
	file := util.MustOpenFile(runtimePath.MyLogPath,logFileName+runtimePath.ExtensionName)
	logger = log.New(file, "", log.LstdFlags)

}
//提示类报错
func LogInfo(err error){
	logger.SetPrefix(ErrPoint("info"))
	logger.Println(err)
}
//警告类报错
func LogWarn(err error){
	logger.SetPrefix(ErrPoint("warn"))
	logger.Println(err)
}
//严重类报错
func LogFatal(err error){
	logger.SetPrefix(ErrPoint("fatal"))
	logger.Fatalln(err)
}
//报错所在的文件、行等信息作为每一条错误日志的开头，方便定位分析
func ErrPoint(logType string)string{
	_, file, line, ok := runtime.Caller(2)//todo:这一行是日志定位报错地点的核心
	if !ok {
		fmt.Println("没能成功显示代码报错位置")
		return "未查清报错的代码行"
	}
	baseFile := path.Base(file)
	return fmt.Sprintf("[[[%s|在%s文件|第%d行]]]:",logType,baseFile,line)
}
