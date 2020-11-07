//本文件内函数用于文件创建等
package util

import (
	"fmt"
	"os"
)
//从相对路径获取全路径
func WdirAndPath(path string)string{
	dir, err := os.Getwd()   //todo:获取工作目录是回从相应盘符根目录一直到当前项目的根目录(本例中即go10bo2ke4)
	fmt.Println("dir--util=",dir)
	if err != nil {

	}
	return dir + "/" + path
}
//打开相应全路径的文件，没有则先创建再打开
func MustOpenFile(path string,name string)*os.File{ //path是相对路径，name是文件名
	fullPath := WdirAndPath(path)
	_, err := os.Stat(fullPath)
	exist := os.IsNotExist(err)
	//如果该路径文件不存在则创建（exist为true表示路径文件不存在），这个是创建目录路径而非文件
	if exist {
		os.MkdirAll(fullPath,0666)
	}
	pathANDname := fullPath +"/"+ name
	//以创建的方式打开该日志txt文件，这个是创建文件，而不能创建目录路径
	file, err := os.OpenFile(pathANDname,os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
	if err !=nil {
		fmt.Print("打开路径失败")
		return nil
	}
	return file

}
