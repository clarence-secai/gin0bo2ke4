package runtime_init

import (
	"encoding/json"
	"fmt"
	"go10bo2ke4/util"
	"os"
	"strings"
)

//结构体字段对应着json配置文件
type RuntimePath struct{
	RunTimePath string
	HeadImagePath string
	MyLogPath string
	ExtensionName string
}

func MyPath() RuntimePath {
	path := util.WdirAndPath("config/path_config.json")
	//把反斜杠均变成顺斜杠,其实反斜杠也打得开
	path= strings.Replace(path, "\\", "/", -1)
	//fmt.Println("path--runtime_init=",path)
	file, err := os.Open(path)
	if err != nil {
		panic("打开运行时路径配置文件失败")
	}
	//读取运行时路径配置文件
	var myPath RuntimePath
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&myPath)
	fmt.Println("mypath=",myPath)
	if err != nil {
		fmt.Println("运行时配置文件解码失败")
	}
	return myPath
}
//读取运行时所需创建存放日志和头像的目录的配置文件到结构体
