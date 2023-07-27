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

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running")
	defer fmt.Println("connID=", c.ConnID, "is stop")
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read err,", err)
			continue
		}
		//调用当前链接所绑定的handleAPI
		err = c.handleAPI(c.Conn, buf, cnt)
		if err != nil {
			fmt.Println("handleAPI err,", err)
			break
		}
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
