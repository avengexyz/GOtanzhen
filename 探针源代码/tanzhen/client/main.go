package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/spf13/viper"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default() //带log等中间件
	// r := gin.New() //不带中间件的路由引擎

	//相应路由规则 执行的函数
	r.POST("/data", getdata)
	r.POST("/ajax", getinfo)

	//获取配置项
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 还可以在工作目录中查找配置
	_ = viper.ReadInConfig()      // 查找并读取配置文件
	port := viper.GetString("port")

	// 启动HTTP服务，默认在127.0.0.1:8080启动服务
	r.Run(":" + port) //也可以写全乎了比如"127.0.0.1:80"或者域名
}

func getdata(c *gin.Context) {
	//调用获取系统信息方法
	data := get()
	// c.JSON：返回JSON格式的数据
	c.JSON(200, gin.H{
		"ststus": 200,
		"data":   data,
		"msg":    "获取数据成功",
	})
}
func getinfo(c *gin.Context) {
	//调用获取系统信息方法
	ajax := getajax()
	// c.JSON：返回JSON格式的数据
	c.JSON(200, gin.H{
		"ststus": 200,
		"data":   ajax,
		"msg":    "获取数据成功",
	})
}

//Data 服务器信息结构体
type Data struct {
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

//get 获取系统信息方法
func get() Data {
	//当前时间
	服务器当前时间 := time.Now().Format("2006-01-02 15:04:05")
	//host主机信息
	主机信息, _ := host.Info()
	运行时间 := 主机信息.Uptime                                                      //正常运行时间
	开机时间 := time.Unix(int64(主机信息.BootTime), 0).Format("2006-01-02 15:04:05") //开机时间
	操作系统 := 主机信息.OS                                                          //操作系统
	cpu架构 := 主机信息.KernelArch                                                 //cpu架构
	//cpu信息s
	cpu信息, _ := cpu.Info()
	cpu型号 := cpu信息[0].ModelName                  //cpu型号
	cpu核数 := cpu信息[0].Cores                      //cpu核数
	cpu频率 := cpu信息[0].Mhz                        //cpu频率
	cpu百分比, _ := cpu.Percent(time.Second, false) //cpu百分比
	//内存信息
	内存信息, _ := mem.VirtualMemory()
	总内存 := strconv.FormatFloat(float64(内存信息.Total)/1073741824, 'f', 2, 64)      //总内存
	可用内存 := strconv.FormatFloat(float64(内存信息.Available)/1073741824, 'f', 2, 64) //可用内存
	已用内存 := strconv.FormatFloat(float64(内存信息.Used)/1073741824, 'f', 2, 64)      //已用内存
	内存百分比 := 内存信息.UsedPercent                                                   //内存百分比
	//硬盘信息
	所有硬盘信息, _ := disk.Partitions(true)                                       //获取所有硬盘信息
	硬盘信息, _ := disk.Usage(所有硬盘信息[0].Device)                                  //指定某路径的硬盘使用情况(取第一个,其余挂载暂时不考虑)
	硬盘总大小 := strconv.FormatFloat(float64(硬盘信息.Total)/1073741824, 'f', 2, 64) //硬盘总大小
	可用硬盘大小 := strconv.FormatFloat(float64(硬盘信息.Free)/1073741824, 'f', 2, 64) //可用硬盘大小
	已用硬盘大小 := strconv.FormatFloat(float64(硬盘信息.Used)/1073741824, 'f', 2, 64) //已用硬盘大小
	// strconv.FormatFloat(float64(硬盘信息.Used)/1073741824, 'f', 2, 64)
	硬盘百分比 := 硬盘信息.UsedPercent //硬盘百分比
	//网络信息
	//获取网络读写字节／包的个数
	网络信息, _ := net.IOCounters(false)
	发送的字节数 := 网络信息[0].BytesSent //发送的字节数
	收到的字节数 := 网络信息[0].BytesRecv //收到的字节数

	return Data{
		A服务器当前时间: 服务器当前时间,
		A运行时间:    运行时间,
		A开机时间:    开机时间,
		A操作系统:    操作系统,
		Acpu架构:   cpu架构,
		Acpu型号:   cpu型号,
		Acpu核数:   cpu核数,
		Acpu频率:   cpu频率,
		Acpu百分比:  int64(cpu百分比[0]),
		A总内存:     总内存,
		A可用内存:    可用内存,
		A已用内存:    已用内存,
		A内存百分比:   int64(内存百分比),
		A硬盘总大小:   硬盘总大小,
		A可用硬盘大小:  可用硬盘大小,
		A已用硬盘大小:  已用硬盘大小,
		A硬盘百分比:   int64(硬盘百分比),
		A发送的字节数:  发送的字节数,
		A收到的字节数:  收到的字节数,
	}
}

type ajax struct {
	Acpu百分比 int64 `json:"cpu百分比"`
	A内存百分比  int64 `json:"内存百分比"`
	A硬盘百分比  int64 `json:"硬盘百分比"`
}

func getajax() ajax {
	//cpu信息
	cpu百分比, _ := cpu.Percent(time.Second, false) //cpu百分比
	//内存信息
	内存信息, _ := mem.VirtualMemory()
	内存百分比 := 内存信息.UsedPercent //内存百分比
	//硬盘信息
	所有硬盘信息, _ := disk.Partitions(true)      //获取所有硬盘信息
	硬盘信息, _ := disk.Usage(所有硬盘信息[0].Device) //指定某路径的硬盘使用情况(取第一个,其余挂载暂时不考虑)
	硬盘百分比 := 硬盘信息.UsedPercent               //硬盘百分比

	return ajax{
		Acpu百分比: int64(cpu百分比[0]),
		A内存百分比:  int64(内存百分比),
		A硬盘百分比:  int64(硬盘百分比),
	}
}
