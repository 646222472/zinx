package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 模拟的服务器
	// 1、创建 SocketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server listen error:", err)
		return
	}
	// 创建一个go承载，负责从客户端处理业务
	go func() {
		// 2、从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accpet error", err)
			}

			go func(conn net.Conn) {
				// 处理客户端有请求
				// ---->拆包的过程<----
				// 定义一个拆包的对象 dp
				dp := NewDataPack()
				for {
					// 1.读取 head
					headData := make([]byte, int(dp.GetHeadLen()))
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error", err)
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						// Msg 中有数据，需要第二次读取
						// 2.根据 head 中的 datalen 再读取 data 内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Printf("--->Recv MsgID:%d, datalen:%d, data:%s\n", msg.ID, msg.DataLen, msg.Data)
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err ", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()
	// 模拟粘包过程，封装两个 msg 一同发送
	// 封装第一个 msg1 包
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error ", err)
	}

	// 封装第一个 msg2 包
	msg2 := &Message{
		ID:      2,
		DataLen: 7,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error ", err)
	}

	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
