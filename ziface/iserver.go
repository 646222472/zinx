package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 运行服务器
	Serve()

	// 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(uint32, IRouter)

	// 获取当前 Server 的链接管理器
	GetConnMgr() IConnManager

	// 注册 OnConnStart 钩子函数的方法
	SetOnConnStart(func(connection IConnection))

	// 调用 OnConnStart 钩子函数的方法
	CallOnConnStart(connection IConnection)

	// 注册 OnConnStop 钩子函数的方法
	SetOnConnStop(func(connection IConnection))

	// 调用 OnConnStop 钩子函数的方法
	CallOnConnStop(connection IConnection)
}
