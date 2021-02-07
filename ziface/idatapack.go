package ziface

// IDataPack 封包，拆包的模块
// 直接面向 TCP 连接中的数据流，用于处理 TCP 粘包问题
type IDataPack interface {
	// 获取头部长度的方法
	GetHeadLen() uint32

	// 封包方法
	Pack(msg IMessage) ([]byte, error)

	// 拆包方法
	UnPack([]byte) (IMessage, error)
}
