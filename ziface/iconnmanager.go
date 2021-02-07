package ziface

// IConnManager 链接管理模块抽象层
type IConnManager interface {
	// 添加链接
	Add(IConnection)

	// 删除链接
	Remove(IConnection)

	// 根据链接ID查找链接
	Get(uint32) (IConnection, error)

	// 获取链接总数
	Len() int

	// 清除所有的链接
	ClearConn()
}
