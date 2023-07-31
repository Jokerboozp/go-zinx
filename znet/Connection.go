package znet

import (
	"fmt"
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

	//告知当前链接已经退出/停止的channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running")
	defer fmt.Println("connID=", c.ConnID, "is stop")
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read err,", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		////调用当前链接所绑定的handleAPI
		//err = c.handleAPI(c.Conn, buf, cnt)
		//if err != nil {
		//	fmt.Println("handleAPI err,", err)
		//	break
		//}
	}
}

func (c *Connection) Start() {
	fmt.Println("connection start,connection id:", c.ConnID)

	//启动从当前链接的读数据业务
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("connection stop,connection id :", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	//关闭
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("close conn err", err)
		return
	}
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
