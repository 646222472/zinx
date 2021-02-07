package znet

import (
	"fmt"
	"io"
	"net"
	"sync"

	"e.coding.net/646222472/framework/zinc/utils"
	"e.coding.net/646222472/framework/zinc/ziface"
)

// Connection 链接模块
type Connection struct {
	// 当前 Conn 隶属于哪个 Server
	TCPServer ziface.IServer
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接ID
	ConnID uint32
	// 当前链接的状态
	isClosed bool
	// 告知当前链接已经停止/退出 channel （由 Reader 告知 Writer 退出）
	ExitChan chan bool
	// 无缓冲的管道，用于 Goroutine 之间的消息通信
	msgChan chan []byte
	// 消息管理 MsgID 和对应的处理业务 API 关系
	MsgHandler ziface.IMsgHandler
	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的修改锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接模块的方法
func NewConnection(tcpServer ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:    tcpServer,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitChan:     make(chan bool, 1),
		msgChan:      make(chan []byte),
		MsgHandler:   msgHandler,
		property:     make(map[string]interface{}),
		propertyLock: sync.RWMutex{},
	}

	// 将 conn 加入到 ConnManager 中
	c.TCPServer.GetConnMgr().Add(c)

	return c
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running ...]")

	defer fmt.Println("connID=", c.ConnID, " 【Reader is exit, remote addr is】 ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包、解包对象
		dp := NewDataPack()

		// 读取客户端的 Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 拆包，得到 msgID 和 msgDataLen 放在 msg 消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// 根据 datalen 再次读取 Data， 放在 msg.Data 中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给 Worker 工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的 MsgID 找到对应的处理业务的 API 方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写消息的 Goroutine 用户将消息发送客户端，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")

	defer fmt.Println(c.RemoteAddr(), "[conn write exit!]")

	// 不断的阻塞等待 channel 的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出，此时 Writer 也要退出
		}
	}
}

// Start 启动链接  让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID=", c.ConnID)
	// 启动从当前链接的读数据的业务
	go c.StartReader()

	// 启动从当前链接的写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的  创建链接之后需要调用的处理业务，执行对应的 hook 函数
	c.TCPServer.CallOnConnStart(c)
}

// Stop 停止链接  结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop() ... ConnID=", c.ConnID)
	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 按照开发者传递进来的  销毁链接之前需要执行对应的 hook 函数
	c.TCPServer.CallOnConnStop(c)

	// 关闭 Socket 链接
	c.Conn.Close()

	// 告知 Writer 关闭
	c.ExitChan <- true

	// 将当前链接从 ConnMgr 中摘除掉
	c.TCPServer.GetConnMgr().Remove(c)

	// 关闭 退出 的channel，回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前链接所绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 提供一个 SendMsg 方法，将我们给客户端的消息先进行封包，再进行发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return fmt.Errorf("%s", "Connection closed when send msg")
	}

	// 将 Data 进行封包 ｜MsgDataLen ｜ MsgID ｜ Data ｜
	dp := NewDataPack()
	binaryData, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgID)
		return fmt.Errorf("%s", "Pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryData

	return nil
}

// SetProPerty 设置链接属性
func (c *Connection) SetProPerty(key string, value interface{}) {
	// 使用写保护锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProPerty 获取链接属性
func (c *Connection) GetProPerty(key string) (interface{}, error) {
	// 使用读保护锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("%s", "property NOT FOUND")
}

// RemoveProPerty 移除链接属性
func (c *Connection) RemoveProPerty(key string) {
	// 使用写保护锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
