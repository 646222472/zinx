package ziface

import "net"

// IConnection 定义链接模块的抽象层
type IConnection interface {
	// 启动链接  让当前链接准备开始工作
	Start()

	// 停止链接  结束当前链接的工作
	Stop()

	// 获取当前链接所绑定的socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前链接模块的链接ID
	GetConnID() uint32

	// 获取远程客户端的 TCP状态 IP Port
	RemoteAddr() net.Addr

	// 发送数据，将我们给客户端的消息先进行封包，再进行发送
	SendMsg(uint32, []byte) error

	// 设置链接属性
	SetProPerty(string, interface{})

	// 获取链接属性
	GetProPerty(string) (interface{}, error)

	// 移除链接属性
	RemoveProPerty(string)
}

// HandleFunc 定义一个处理链接业务的抽象方法
// *net.TCPConn 客户端链接， []byte发送给客户端的数据，int发送的数据长度
type HandleFunc func(*net.TCPConn, []byte, int) error
