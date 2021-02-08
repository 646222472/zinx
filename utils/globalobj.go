package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/646222472/zinx/ziface"
)

// GlobalOjb 存储一切有关Zinx框架的全局参数，供其它模块使用
// 一切参数是可以通过zinx.json由用户进行配置
type GlobalOjb struct {
	// Server
	TCPServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TCPPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	// Zinx
	Version          string //当前Zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作 Worker 池 Goroutine 的数量
	MaxWorkerTaskLen uint32 // 框架允许用户最多开辟多少个 Worker（限定条件）
}

// GlobalObject 定义一个全局的对外GlobalObj
var GlobalObject *GlobalOjb

// Reload 从 zinx.json 中加载用户自定义的参数
func (g *GlobalOjb) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalOjb{
		Name:             "ZinxServerApp",
		Version:          "V0.5",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   // 框架中 WorkerPool 中 Worker 的数量
		MaxWorkerTaskLen: 1024, // 每个 Worker 对应的消息队列中 task 数量的最大值
	}

	// 应该尝试从 conf/zinx.json 中加载一些用户自定义的参数
	GlobalObject.Reload()
}
