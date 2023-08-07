package znet

import (
	"fmt"
	"net"
	"zinx/utils"
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

	//告知当前链接已经退出/停止的channel,由reader告知writer退出
	ExitChan chan bool

	//无缓冲的管道，由于读写goroutine之间的消息通信
	MsgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMessageHandler
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, handle ziface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handle,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	return c
}

// StartWriter 写消息。专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit")
	//不断的阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.MsgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出，此时writer也要退出
			return
		}
	}
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running")
	defer fmt.Println("connID=", c.ConnID, "is stop")
	defer c.Stop()

	for {
		//读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read err,", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := c.Conn.Read(headData)
		if err != nil {
			fmt.Println("read msg head err,", err)
			break
		}

		//拆包，得到msgID和msgDatalen放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err,", err)
			break
		}

		//根据datalen再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := c.Conn.Read(data)
			if err != nil {
				fmt.Println("read msg data err,", err)
				break
			}
		}

		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//判断是否已经开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给worker工作池即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}

		//执行注册的路由方法
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)

		//从路由中找到注册绑定的Conn对应的router调用
		//根据绑定好的msgID，找到对应处理API业务 执行
		go c.MsgHandler.DoMsgHandler(&req)

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
	go c.StartWriter()
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
	//告知writer关闭
	c.ExitChan <- true
	close(c.ExitChan)
	close(c.MsgChan)
}

// SendMsg 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return fmt.Errorf("connection closed when send msg")
	}

	//将data进行封包
	dp := NewDataPack()

	//MsgDataLen|MsgID|Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack err msg id:", msgId)
		return fmt.Errorf("pack err msg id:%d", msgId)
	}

	//将数据发送给管道
	c.MsgChan <- binaryMsg

	return nil
}

func NewMsgPackage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
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
