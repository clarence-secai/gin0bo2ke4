package handlers
import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go10bo2ke4/init/log_init"
	"go10bo2ke4/myerr"
	"go10bo2ke4/mysql_redis"
	"go10bo2ke4/util"
	"net/http"
	"time"
)

//todo:更多jwt-go用法，参见jwt-go的官方文档https://godoc.org/github.com/dgrijalva/jwt-go#example-New--Hmac

//一个jwt token是由两个.分隔成三部分是字母数字串，大约像下面这个样子
//jolei847wksieh.oisne028doduwjrosue8wodj.dietbdiep8465jeuf93
//第一部分称为Header头部，第二部分称为Payload负载，第三部分称为Signature签名
//下面的代码中把第二部分Payload的原始值叫做initPayLoad，第三方包jwt-go中则
//把它叫做claim

var MySecretKey = "111"   //可以任意自定义
type User struct{
	Email string
	PassWord string
}
//前端登录时制作一个token给前端
func LoginJwt(c *gin.Context){
	E := c.PostForm("email")
	P := c.PostForm("passWord")
	if E=="" || P==""{
		util.Response(c,http.StatusBadRequest,myerr.Err_PrefaceData_null,"请输入再提交")
		return
	}
	user := mysql_redis.User{Email: E,PassWord: P}
	//检查数据库来判断该用户是否已注册
	userFull,err := mysql_redis.GetUser(&user)
	if err != nil {
		util.Response(c,http.StatusInternalServerError,myerr.Err_Select_Fail,nil)
		return
	}
	if !(userFull.ID > 0) {
		util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"请先注册")
		return
	}
	//已注册，为此次登录制作jwt token
	loginUser := User{Email:E,PassWord:P}
	jwt_token, err := GenerateJwt(loginUser)
	if err != nil {
		log_init.LogInfo(err)
		util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"制作token出错")
		return
	}
	//依据不同的登录用户，决定前端页面是否展现增查改删博客的功能按钮
	if E == "1787664981@qq.com" && P == "413188ok" {
	//todo:方式一，将token放进响应体【切记，属于同一次响应的响应头内容和响应体内容
	// 不可分开用俩util.Response各用一次c来响应，否则响应头和响应体不会被一起响应给前端
		//util.Response(c,http.StatusOK,myerr.OK,gin.H{
		//	//"jwt_token":jwt_token,
		//	"userType":"author",
		//})
	//todo: 方式二，将token放进Authorization响应头，其他信息放进响应体，通过
	// util.ResponseHeader 内的c.Header() c.Json() 一起组成一次完整响应发给前端
		util.ResponseHeader(c,http.StatusOK,myerr.OK,jwt_token,gin.H{"userType":"author"})
	}else{
		//util.Response(c,http.StatusOK,myerr.OK,gin.H{
		//	//"jwt_token":jwt_token,
		//	"userType":"normalUser",
		//})
		util.ResponseHeader(c,http.StatusOK,myerr.OK,jwt_token,gin.H{"userType":"normalUser"})
	}
}
//校验前端携带来的jwt_token，并获取token中存放的唯一辨识客户的信息
func VerifyJwt(c *gin.Context){
	//从前端发来的请求头中Authorization字段中获取携带来的jwt_token，但这就要求当初在
	//jwt_token := c.Request.Header.Get("Authorization")
	//将jwt token交给前端时就是给到响应头的Authorization字段的
	jwt_token := c.GetHeader("Authorization")
	initPayLoad, err := ParseJwt(jwt_token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors{
		case jwt.ValidationErrorExpired:
			util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"token超时，请重新登录")
		default:
			util.Response(c,http.StatusInternalServerError,myerr.Err_ServerErr,"服务端解析token出错")
		}
		c.Abort()
	}
	if initPayLoad != nil {
		u := User{initPayLoad.Email,initPayLoad.PassWord}
		c.Set("tokenContent",u)//todo 将解析jwt token得到的数据设置给后面的处理器来
		// 获取使用，这就解决了处理器函数因不能有返回值而无法让获取到的前端数据留存着被后续使用的尴尬
		c.Next()
		return
	}
	util.Response(c,http.StatusBadRequest,myerr.Err_Request_Err,"请先登录")
	c.Abort()
	return
}


type MyClaim struct{
	 Email string `json:"email"`//两个字段的json并非必须，为和jwt.StandardClaims一直故最好也写上json
	 PassWord string `json:"passWord"` //这两个字段是自定义的，切记不能是interface{}类型，写一般的字段最保险
	jwt.StandardClaims //todo:这里只能匿名继承该字段，不能起一个字段名来有名继承，否
	 //则jwt.NewWithClaims(jwt.SigningMethodES256,payLoad)会报payLoad未实现接口
}
func GenerateJwt(u User)(string,error){
	//配置好jwt token三部分完整字母数字串的第二部分的原始值
	initPayLoad := MyClaim{
		u.Email,
		u.PassWord,
		jwt.StandardClaims{
			ExpiresAt:time.Now().Add(time.Hour*2).Unix(),//jwt token在客户端的失效时长
			IssuedAt: time.Now().Unix(),
			Issuer: "bo2ke4",
		},
	}
	//搭配好jwt token共三部分完整字母数字串的第一(jwt.SigningMethodHS256参数就是)、二部分
	//的原始值initPayLoad,赋值给rowToken，准备下一步制作jwt token
	rowToken := jwt.NewWithClaims(jwt.SigningMethodHS256, initPayLoad)

	//将rowToken内的第一部分Header和第二部分Payload(这里也叫claim)编码后用.勾连，被通过MysecretKey算出
	//签名作为第三部分Signature，再和第一部分Header、第二部分Payload一起用.勾连，成为最终的jwt token。
	//注意，使用jwt.SigningMethodHS256后SignedString内的参数必须是[]byte。
	jwt_token, err := rowToken.SignedString([]byte(MySecretKey))
	if err != nil {
		fmt.Print("制作token出错：",err)
		return "",err
	}
	return jwt_token,nil
}
//函数必须是下面这个模样才能传进jwt.ParseWithClaims作为第三个参数，返回值必须是
//原先自定义好的密钥MySecretKey，和制作jwt token时使用的同一个密钥。
func f(rowToken *jwt.Token)(interface{},error){
	//parMySecretKey, _ := jwt.ParseRSAPrivateKeyFromPEM(MySecretKey)
	return []byte(MySecretKey),nil
}
func ParseJwt(jwt_token string)(*MyClaim,error){
	//通过f中的密钥MySecretKey解析出前端传来的jwt token字母数字串中的第一、二、三部分的
	//原始值组成的rowToken。解析返回的err是已经设置好的错误内容，通过
	//err.(*jwt.ValidationError).Errors 可以取到错误值，来判断是解析jwt token出错还是jwt token过期了
	//jwt.ValidationErrorExpired:  //todo:表示jwt_token过期了

	rowToken, err := jwt.ParseWithClaims(jwt_token, &MyClaim{}, f)

	//如果为空则是因为传进来的参数jwt_token本身是空串，因此需要这一步判断一下
	if rowToken != nil {
		//取到rowToken下的Claims字段(即jwt token字母数字串的第二部分原始值)后，断言获得自定义
		//的MyClaim类型变量，也就是jwt token字母数字串中的第二部分的原始值
		initPayLoad,ok := rowToken.Claims.(*MyClaim) //todo:这里之所以断言MyClaim指针，可让initPayLoad类型不同
													// 是方便110行在initPayLoad无值时返回nil，断言MyClaim也并非不行
		if ok && rowToken.Valid {
			return initPayLoad,nil
		}
	}
	return nil,err  //err不为空时，返回无值的initPayLoad和err
}