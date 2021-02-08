package znet

import (
	"fmt"
	"sync"

	"github.com/646222472/zinx/ziface"
)

// ConnManager 实现链接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的链接集合
	connLock    sync.RWMutex                  // 保护链接集合的的读写锁
}

// NewConnManager 初始化当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 加入 ConnManager 中
	cm.connections[conn.GetConnID()] = conn
	fmt.Printf("ConnID=%d add to ConnManager successfully:conn num := %d\n", conn.GetConnID(), cm.Len())
}

// Remove 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除链接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Printf("ConnID=%d remove to ConnManager successfully:conn num := %d\n", conn.GetConnID(), cm.Len())
}

// Get 根据链接ID查找链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map， 加写锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}

	return nil, fmt.Errorf("%s", "connection NOT FOUND")
}

// Len 获取链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除所有的链接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源 map， 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()

		// 删除
		delete(cm.connections, connID)
	}

	fmt.Printf("Clear All connections succ! conn num=%d", cm.Len())
}
