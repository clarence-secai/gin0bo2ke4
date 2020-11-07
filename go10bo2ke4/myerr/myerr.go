package myerr

const (
	OK                       =2000
	Err_Request_Err = 4000
	Err_PrefaceData_Get_Fail = 4001
	Err_PrefaceData_null = 4002
	Err_ServerErr			=5000
	Err_Select_Fail          = 5001
	Err_Add_Fail      = 5002
	Err_Update_Fail   = 5003
	Err_Delete_Fail 	 = 5004
)
var ErrMap = map[int]string{
	2000:"操作成功",
	4000:"请求本身有误",
	4001:"前端发送数据获取出错",
	4002:"前端发送的数据是空的",
	5000:"服务端处理数据出错",
	5001:"查找数据表出错",
	5002:"添加进数据库出错",
	5003:"更新数据库出错",
	5004:"删除数据库出错",
}