package znet

import (
	"net"
	"zinx/ziface"
)

/*
*链接模块
 */

type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前链接状态
	isClosed bool

	//当前链接所绑定的业务处理方法API
	handleAPI ziface.HandleFunc

	//告知当前链接已经退出/停止的channel
	ExitChan chan bool
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callBackAPI ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleAPI: callBackAPI,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
