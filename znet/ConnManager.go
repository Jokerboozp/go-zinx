package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/**
链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (c *ConnManager) Add(connection ziface.IConnection) {
	//保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn加入到ConnManager
	c.connections[connection.GetConnID()] = connection
	fmt.Println("connID = ", connection.GetConnID(), " connection add to connManager success: conn num = ", c.Len())
}

// Remove 删除链接
func (c *ConnManager) Remove(connection ziface.IConnection) {
	//保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除链接信息
	delete(c.connections, connection.GetConnID())
	fmt.Println("connID = ", connection.GetConnID(), " connection delete success: conn num = ", c.Len())
}

// Get 根据connID获取链接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// Len 得到当前链接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// ClearConn 清除并终止所有的链接
func (c *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}
	fmt.Println("clear all connections.conn num = ", c.Len())
}
