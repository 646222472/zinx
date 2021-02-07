package znet

import (
	"fmt"
	"net"
	"re_zinx/utils"
	"re_zinx/ziface"
)

// Server iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP的版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口Port
	Port int
	// 当前 Server 的消息管理模块，用来绑定 MsgID 和对应的处理业务 API 关系
	MsgHandler ziface.IMsgHandler
	// 该 Server 的链接管理器
	ConnMgr ziface.IConnManager
	// 该 Server 创建链接之后自动调用 Hook 函数 -- OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 该 Server 销毁链接之前自动调用 Hook 函数 -- OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf(
		"[Zinx] Server Name:%s, listenner at IP:%s, Port:%d is starting ...\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TCPPort,
	)
	fmt.Printf(
		"[Zinx] Version %s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)
	fmt.Printf("[Start] Server Listenner at IP :%s, Port%d, is starting\n", s.IP, s.Port)

	go func() {
		// 0 开启消息队列及 Worker 工作池
		s.MsgHandler.StartWorkerPool()

		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
			return
		}

		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Printf("start Zinx Server succ, %s succ, Listenning ...\n", s.Name)
		var cid uint32
		cid = 0

		// 3 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accpet err", err)
				continue
			}

			// 判断当前系统中的链接数量是否大于最大链接数 MaxConn
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大链接的错误包
				fmt.Printf("Too many connections MaxConn=%d", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务模块
			go dealConn.Start()
		}

	}()

}

// Stop 停止服务器
func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或者回收
	fmt.Printf("[STOP] Zinx server name %s", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	// 启动Server的服务功能
	s.Start()

	// TODO 做一些服务启动服务之后的额外业务

	// 阻塞状态
	select {}
}

// AddRouter 添加路由功能
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

// GetConnMgr 获取当前 Server 的链接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化Server的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

// SetOnConnStart 注册 OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// CallOnConnStart 调用 OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--->Call OnConnStart")
		s.OnConnStart(conn)
	}
}

// SetOnConnStop 注册 OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStop 调用 OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--->Call OnConnStop")
		s.OnConnStop(conn)
	}
}
