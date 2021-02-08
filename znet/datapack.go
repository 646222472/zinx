package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/646222472/zinx/utils"
	"github.com/646222472/zinx/ziface"
)

// DataPack 封包，拆包的模块
// 直接面向 TCP 连接中的数据流，用于处理 TCP 粘包问题
type DataPack struct {
}

// NewDataPack 封包，拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取头部长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法 |dataLen(4)|MsgId(4)|MsgData|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 bytes 字节流的缓冲
	dataBuffer := bytes.NewBuffer([]byte{})

	// 将 DataLen 写入 dataBuffer 中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将 MsgId 写入 dataBuffer 中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 将 Data 写入 dataBuffer 中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

// UnPack 拆包方法:1、将包中的 Head 信息读出；2、再根据 Head 信息中的 DataLen，再进行一次读取
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个 ioReader (数据从输入的二进制数据流得到)
	dataBuffer := bytes.NewReader(binaryData)

	// 接收拆包的数据
	msg := &Message{}

	// 只解压 Head 信息，得到 DataLen 和 MsgID

	// 读取 DataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取 MsgID
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	// 判断 datalen 是否已经超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return nil, fmt.Errorf("%s", "too large message data recv !!!")
	}

	return msg, nil
}
