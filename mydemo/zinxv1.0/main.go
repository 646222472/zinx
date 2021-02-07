package main

import (
	"fmt"
	"re_zinx/ziface"
	"re_zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 在处理conn业务的主方法Hook
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")

	// 先读取客户端的数据，再回写 ping... ping... ping...
	fmt.Printf(
		"Recv from client: MsgID=%d, data=%s\n",
		request.GetMsgID(), request.GetData(),
	)

	err := request.GetConnection().SendMsg(1, []byte("ping... ping... ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloRouter 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// Handle 在处理conn业务的主方法Hook
func (pr *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")

	// 先读取客户端的数据，再回写 hello... hello... hello...
	fmt.Printf(
		"Recv from client: MsgID=%d, data=%s\n",
		request.GetMsgID(), request.GetData(),
	)

	err := request.GetConnection().SendMsg(1, []byte("Hello Welcome to Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionAfter 创建链接之后的执行的钩子函数
func DoConnectionAfter(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionAfter is call")
	if err := conn.SendMsg(202, []byte("DoConnection After")); err != nil {
		fmt.Println(err)
	}

	// 给当前链接设置一些属性
	fmt.Println("Set Conn Name, Hoe ...")
	conn.SetProPerty("Name", "雨非雨")
	conn.SetProPerty("GitHub", "http://github.com/yx_yufeiyu")
}

// DoDestroyBegin 销毁之前执行的钩子函数
func DoDestroyBegin(conn ziface.IConnection) {
	fmt.Println("===> DoDestroyBegin is call")
	fmt.Println("ConnID is ", conn.GetConnID(), " is lost")

	// 链接销毁之前获取链接属性属性
	if name, err := conn.GetProPerty("Name"); err == nil {
		fmt.Println("Name is ", name)
	}
	if github, err := conn.GetProPerty("GitHub"); err == nil {
		fmt.Println("GitHub is ", github)
	}
}

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.1]")

	// 注册链接 Hook 钩子函数
	s.SetOnConnStart(DoConnectionAfter)
	s.SetOnConnStop(DoDestroyBegin)

	// 2.给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 3.启动Server
	s.Serve()
}
