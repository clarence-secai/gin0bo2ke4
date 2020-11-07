package util

import (
	"github.com/gin-gonic/gin"
	"go10bo2ke4/myerr"
)

type Data struct{
	Code int `json:"code"`
	CodeMessage string  `json:"codeMessage"`
	Content interface{} `json:"content"`
}
func Response(c *gin.Context, httpCode int,code int,content interface{}){
	d := Data{
		Code:code,
		CodeMessage: myerr.ErrMap[code],
		Content: content,
	}
	c.JSON(httpCode,d)
}
func Rrespon(c *gin.Context,head string){
	c.Header("Authorization",head)
}
//todo: 切记，不可在一个func(c *gin.Default)类型的handler里使用两个响应前端函数，因为这样
// 就需要分别将c传给这两个响应前端的函数，这就会导致不被认为是同一次响应。
// 像下面的ResponseHeader
// 则是使用的同一个c下两个响应前端的代码，都会被当作是同一次完整的乡音头响应体的c响应
func ResponseHeader(c *gin.Context, httpCode int,code int,head string,content interface{}){
	c.Header("Authorization",head)
	d := Data{
		Code:code,
		CodeMessage: myerr.ErrMap[code],
		Content: content,
	}
	c.JSON(httpCode,d)
}


