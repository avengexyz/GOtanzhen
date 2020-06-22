package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default() //带log等中间件
	// r := gin.New()//不带中间件的路由引擎

	//加载静态文件
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico") //单个文件
	//加载模板文件
	r.LoadHTMLGlob("views/*")

	// 允许使用跨域请求  全局中间件
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	// r.Use(Cors())

	//相应路由规则 执行的函数
	r.GET("/", Index)
	r.GET("/admin", session, Admin)
	r.GET("/login", Login)
	r.POST("/adminlogin", Adminloginapi)
	//后台注销
	r.GET("/logout", Logout)

	r.GET("/list", List)
	r.GET("/ajax", Listinfo)
	r.POST("/check", Check)

	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	port := viper.GetString("port")

	// 启动HTTP服务，默认在127.0.0.1:8080启动服务
	r.Run(":" + port) //也可以写全乎了比如"127.0.0.1:80"或者域名

}

//Index 前台主页方法
func Index(c *gin.Context) {
	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	status := viper.GetString("switch")
	if status != "true" {
		//如果有错误 说明没有cookie 跳转到登录页
		c.Redirect(http.StatusFound, "/login")
		return
	}
	//调用获取客户端信息方法
	data := getlist()
	c.HTML(http.StatusOK, "index.html", data)
}

//Admin 后台主页方法
func Admin(c *gin.Context) {
	//调用获取客户端信息方法
	data := getlist()
	c.HTML(http.StatusOK, "adminindex.html", data)
}

//Login 登录主页方法
func Login(c *gin.Context) {
	//获取cookie
	cookie, _ := c.Cookie("session")
	//如果有session就跳转到index
	if len(cookie) > 0 {
		//如果有session 跳转到登录页
		c.Redirect(http.StatusFound, "/admin")
	}
	c.HTML(http.StatusOK, "login.html", "")
}

//List 获取列表接口
func List(c *gin.Context) {
	//调用获取客户端信息方法
	data := getlist()
	//返回前端所需json数据
	c.JSON(200, data)
}

//Listinfo ajax获取列表接口
func Listinfo(c *gin.Context) {
	//调用获取客户端信息方法
	data := getajax()
	//返回前端所需json数据
	c.JSON(200, data)
}

//Adminloginapi 登录api
func Adminloginapi(c *gin.Context) {
	//获取用户输入
	password := c.PostForm("password")
	//如果输入是空
	if len(password) == 0 {
		c.JSON(200, gin.H{
			"status": 202,
			"msg":    "参数错误",
		})
		return
	}
	//根据用户名查询账号信息
	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	session := viper.GetString("key")
	//如果没有用户名或者密码不匹配
	if session != password {
		c.JSON(200, gin.H{
			"status": 203,
			"msg":    "用户名或密码错误",
		})
		return
	}
	//写入session
	c.SetCookie("session", password, 86400, "/", "", false, true)
	//否则就正常
	c.JSON(200, gin.H{
		"status": 200,
		"msg":    "登录成功",
	})
}

//session session验证中间件
func session(c *gin.Context) {
	//获取cookie
	cookie, err := c.Cookie("session")
	if err != nil {
		//如果有错误 说明没有cookie 跳转到登录页
		c.Redirect(http.StatusFound, "/login")
		//不调用该请求的剩余处理程序
		c.Abort()
	}
	//根据cookie从数据库获取session
	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	session := viper.GetString("key")
	if cookie != session {
		//清空错误session
		c.SetCookie("session", cookie, -1, "/", "", false, true)
		//如果本地cookie和session对不上 跳转到登录页
		c.Redirect(http.StatusFound, "/login")
		//不调用该请求的剩余处理程序
		c.Abort()
	}

	//全都正确则继续执行后面的方法
	// 调用该请求的剩余处理程序
	c.Next()
}

//Logout 后台注销方法
func Logout(c *gin.Context) {
	//删除session
	c.SetCookie("session", "", -1, "/", "", false, true)
	//跳转至登录页
	c.Redirect(http.StatusFound, "/login")

}

//Check 验证方法
func Check(c *gin.Context) {
	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	confkey := viper.GetString("key")
	//接收表单
	key := c.PostForm("key")
	//判断表单数据是否正确
	if key != confkey {
		c.JSON(200, gin.H{
			"status": 201,
			"msg":    "验证码错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"msg":    "验证码正确",
	})

}

//getlist 获取客户端信息
func getlist() []info {
	//获取配置文件中的信息
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	//将读取到的配置信息赋值给变量
	ip := viper.GetStringSlice("ip")
	name := viper.GetStringSlice("name")
	//定义一个数据切片用于存储请求返回的json
	var infodata []info
	//循环请求数据
	for i := 0; i < len(ip); i++ {
		//访问链接获取数据
		resp, err := http.Post("http://"+ip[i]+"/data", "application/x-www-form-urlencoded", strings.NewReader(""))
		if err != nil {
			fmt.Println("http.Post:", err)
		}
		defer resp.Body.Close()
		//读取返回的数据
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ioutil.ReadAll错误:", err)
		}
		//声明一个info结构体用于存取请求的json数据
		var info info
		// err = json.Unmarshal([]byte(string(body)), &info)
		//对请求返回的json数据反序列化
		_ = json.Unmarshal(body, &info)
		info.Data.A服务器名称 = name[i]
		//将反序列化后的数据 写入info切片
		infodata = append(infodata, info)
	}

	return infodata
}

type info struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}

type Data struct {
	A服务器名称   string  `json:"服务器名称"`
	A服务器当前时间 string  `json:"服务器当前时间"`
	A运行时间    uint64  `json:"运行时间"`
	A开机时间    string  `json:"开机时间"`
	A操作系统    string  `json:"操作系统"`
	Acpu架构   string  `json:"cpu架构"`
	Acpu型号   string  `json:"cpu型号"`
	Acpu核数   int32   `json:"cpu核数"`
	Acpu频率   float64 `json:"cpu频率"`
	Acpu百分比  int64   `json:"cpu百分比"`
	A总内存     string  `json:"总内存"`
	A可用内存    string  `json:"可用内存"`
	A已用内存    string  `json:"已用内存"`
	A内存百分比   int64   `json:"内存百分比"`
	A硬盘总大小   string  `json:"硬盘总大小"`
	A可用硬盘大小  string  `json:"可用硬盘大小"`
	A已用硬盘大小  string  `json:"已用硬盘大小"`
	A硬盘百分比   int64   `json:"硬盘百分比"`
	A发送的字节数  uint64  `json:"发送的字节数"`
	A收到的字节数  uint64  `json:"收到的字节数"`
}

//getajax 获取部分客户端信息
func getajax() []info1 {
	//获取配置文件中的信息
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	//将读取到的配置信息赋值给变量
	ip := viper.GetStringSlice("ip")
	name := viper.GetStringSlice("name")
	//定义一个数据切片用于存储请求返回的json
	var infodata []info1
	//循环请求数据
	for i := 0; i < len(ip); i++ {
		//访问链接获取数据
		resp, err := http.Post("http://"+ip[i]+"/data", "application/x-www-form-urlencoded", strings.NewReader(""))
		if err != nil {
			fmt.Println("http.Post:", err)
		}
		defer resp.Body.Close()
		//读取返回的数据
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ioutil.ReadAll错误:", err)
		}
		//声明一个info结构体用于存取请求的json数据
		var info info1
		// err = json.Unmarshal([]byte(string(body)), &info)
		//对请求返回的json数据反序列化
		_ = json.Unmarshal(body, &info)
		info.Data.A服务器名称 = name[i]
		//将反序列化后的数据 写入info切片
		infodata = append(infodata, info)
	}

	return infodata
}

type info1 struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   Data1  `json:"data"`
}

type Data1 struct {
	A服务器名称  string `json:"服务器名称"`
	Acpu百分比 int64  `json:"cpu百分比"`
	A内存百分比  int64  `json:"内存百分比"`
	A硬盘百分比  int64  `json:"硬盘百分比"`
}

//Cors 跨域
// func Cors() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		method := c.Request.Method               //请求方法
// 		origin := c.Request.Header.Get("Origin") //请求头部
// 		var headerKeys []string                  // 声明请求头keys
// 		for k, _ := range c.Request.Header {
// 			headerKeys = append(headerKeys, k)
// 		}
// 		headerStr := strings.Join(headerKeys, ", ")
// 		if headerStr != "" {
// 			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
// 		} else {
// 			headerStr = "access-control-allow-origin, access-control-allow-headers"
// 		}
// 		if origin != "" {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
// 			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
// 			//header的类型
// 			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
// 			//允许跨域设置 可以返回其他子段
// 			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
// 			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
// 			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
// 			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
// 		}

// 		//放行所有OPTIONS方法
// 		if method == "OPTIONS" {
// 			c.JSON(http.StatusOK, "Options Request!")
// 		}
// 		// 处理请求
// 		c.Next() //  处理请求
// 	}
// }
